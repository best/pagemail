package api

import (
	"fmt"
	"net/http"
	"os"
	"pagemail/internal/converter"
	"pagemail/internal/database"
	"pagemail/internal/mailer"
	"pagemail/internal/models"
	"pagemail/internal/scraper"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type ScrapeRequest struct {
	URL    string `json:"url" binding:"required,url"`
	Email  string `json:"email" binding:"required,email"`
	Format string `json:"format" binding:"required,oneof=html pdf screenshot"`
}

type ScrapeResponse struct {
	RequestID uint   `json:"request_id"`
	Message   string `json:"message"`
	Status    string `json:"status"`
}

func handleScrapeRequest(c *gin.Context) {
	var req ScrapeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithValidationError(c, err.Error())
		return
	}

	// Get user ID if authenticated
	userID, _ := c.Get("user_id")
	var userIDPtr *uint
	if uid, ok := userID.(uint); ok {
		userIDPtr = &uid
	}

	// Create request record
	request := models.Request{
		UserID: userIDPtr,
		URL:    req.URL,
		Email:  req.Email,
		Format: req.Format,
		Status: "pending",
	}

	if err := database.DB.Create(&request).Error; err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeDatabaseError, "Failed to create scrape request")
		return
	}

	// Process request asynchronously
	go processRequest(&request)

	RespondWithAccepted(c, ScrapeResponse{
		RequestID: request.ID,
		Message:   "Scrape request accepted and is being processed",
		Status:    "pending",
	})
}

func processRequest(request *models.Request) {
	// Update status to processing
	request.Status = "processing"
	database.DB.Save(request)

	// Initialize scraper and converter
	scraperManager := scraper.NewManager()
	defer scraperManager.Close()

	converterManager := converter.NewManager()

	var content []byte
	var err error

	// Determine scraping strategy
	useChrome, _ := scraperManager.GetRecommendedStrategy(request.URL)

	// Get content based on format
	switch request.Format {
	case "screenshot":
		content, err = scraperManager.Screenshot(request.URL)
	default:
		// For HTML and PDF, we need HTML content first
		content, err = scraperManager.ScrapeHTML(request.URL, useChrome)
	}

	if err != nil {
		// Update request with error
		request.Status = "failed"
		request.ErrorMsg = err.Error()
		database.DB.Save(request)
		return
	}

	// Generate filename with proper extension
	ext := converterManager.GetOutputExtension(request.Format)
	filename := fmt.Sprintf("%s_%d_%d%s", request.Format, request.ID, time.Now().Unix(), ext)

	// Get files directory from environment variable or use default
	filesDir := os.Getenv("FILES_DIR")
	if filesDir == "" {
		filesDir = "files"
	}
	
	// Ensure files directory exists
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		os.MkdirAll(filesDir, 0755)
	}

	filePath := filepath.Join(filesDir, filename)

	// Convert and save content
	if err := converterManager.ConvertContent(content, request.Format, request.URL, filePath); err != nil {
		request.Status = "failed"
		request.ErrorMsg = fmt.Sprintf("Failed to convert content: %v", err)
		database.DB.Save(request)
		return
	}

	// Update request as completed
	now := time.Now()
	request.Status = "completed"
	request.FilePath = filePath
	request.CompletedAt = &now
	database.DB.Save(request)

	// Send email with attachment
	mailerService := mailer.NewMailer()
	if err := mailerService.SendPageMail(request.Email, request.URL, request.Format, filePath); err != nil {
		// Log email error but don't fail the request
		fmt.Printf("Warning: Failed to send email for request %d: %v\n", request.ID, err)
		// Could optionally update request status to "completed_no_email"
	}

	fmt.Printf("Request %d completed: %s -> %s (%s)\n", request.ID, request.URL, filePath, request.Format)
}

func handleRequestHistory(c *gin.Context) {
	// Get validated user ID from middleware
	userID, exists := c.Get("validated_user_id")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized)
		return
	}

	var requests []models.Request
	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(50).
		Find(&requests).Error; err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeDatabaseError, "Failed to fetch scrape history")
		return
	}

	// Remove file paths and error messages from response for security
	for i := range requests {
		requests[i].FilePath = ""
		if requests[i].Status != "failed" {
			requests[i].ErrorMsg = ""
		}
	}

	RespondWithSuccess(c, gin.H{
		"scrapes": requests,
		"total":   len(requests),
	})
}

func handleScrapeDetail(c *gin.Context) {
	scrapeID := c.Param("scrape_id")
	if scrapeID == "" {
		RespondWithError(c, http.StatusBadRequest, ErrorCodeMissingParameter, "scrape_id parameter is required")
		return
	}

	var request models.Request
	if err := database.DB.Where("id = ?", scrapeID).First(&request).Error; err != nil {
		if err.Error() == "record not found" {
			RespondWithError(c, http.StatusNotFound, ErrorCodeUserNotFound, "Scrape not found")
		} else {
			RespondWithError(c, http.StatusInternalServerError, ErrorCodeDatabaseError)
		}
		return
	}

	// Check if user has access to this scrape
	userID, exists := c.Get("user_id")
	if exists {
		// Authenticated user - check if they own this scrape
		if request.UserID != nil && *request.UserID == userID.(uint) {
			// User owns this scrape, return full details
		} else {
			// User doesn't own this scrape, deny access
			RespondWithError(c, http.StatusForbidden, ErrorCodeUnauthorized, "Access denied")
			return
		}
	} else {
		// Guest user - can only see their own scrapes, but we have no way to verify
		// For security, don't allow guest access to specific scrape details
		RespondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, "Authentication required")
		return
	}

	// Remove sensitive information
	request.FilePath = ""
	if request.Status != "failed" {
		request.ErrorMsg = ""
	}

	RespondWithSuccess(c, request)
}
