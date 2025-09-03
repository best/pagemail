package api

import (
	"fmt"
	"net/http"
	"os"
	"pagemail/internal/auth"
	"pagemail/internal/database"
	"pagemail/internal/mailer"
	"pagemail/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithValidationError(c, err.Error())
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		RespondWithError(c, http.StatusConflict, ErrorCodeUserAlreadyExists)
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to hash password")
		return
	}

	// Create user (inactive until email verified)
	user := models.User{
		Email:         req.Email,
		Password:      hashedPassword,
		IsActive:      false,
		EmailVerified: false,
		DailyLimit:    10,
		MonthlyLimit:  300,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeUserCreationFailed)
		return
	}

	// Generate verification token
	verificationService := auth.NewEmailVerificationService()
	token, err := verificationService.GenerateVerificationToken()
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to generate verification token")
		return
	}

	// Set verification token
	if err := verificationService.SetVerificationToken(user.ID, token); err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to set verification token")
		return
	}

	// Send verification email
	clientIP := c.ClientIP()
	canSend, message, err := verificationService.CanSendVerificationEmail(req.Email, clientIP)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to check send limits")
		return
	}

	if !canSend {
		RespondWithError(c, http.StatusTooManyRequests, ErrorCodeTooManyRequests, message)
		return
	}

	// Build verification URL
	baseURL := getBaseURL()
	verificationURL := fmt.Sprintf("%s/api/auth/verify/%s", baseURL, token)
	
	// Send email
	mailService := mailer.NewMailer()
	if err := mailService.SendVerificationEmail(req.Email, verificationURL); err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeEmailServiceError, "Failed to send verification email")
		return
	}

	// Record email sent
	if err := verificationService.RecordVerificationEmailSent(req.Email, clientIP); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to record verification email: %v\n", err)
	}

	user.Password = "" // Don't return password
	RespondWithCreated(c, gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
		"user":    user,
	})
}

func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithValidationError(c, err.Error())
		return
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			RespondWithError(c, http.StatusUnauthorized, ErrorCodeInvalidCredentials)
			return
		}
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeDatabaseError)
		return
	}

	// Check if user is active
	if !user.IsActive {
		RespondWithError(c, http.StatusUnauthorized, ErrorCodeAccountDeactivated)
		return
	}

	// Check if email is verified
	if !user.EmailVerified {
		RespondWithError(c, http.StatusForbidden, ErrorCodeEmailNotVerified)
		return
	}

	// Verify password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		RespondWithError(c, http.StatusUnauthorized, ErrorCodeInvalidCredentials)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to generate token")
		return
	}

	user.Password = "" // Don't return password
	RespondWithSuccess(c, AuthResponse{
		Token: token,
		User:  user,
	})
}

func handleProfile(c *gin.Context) {
	// Get validated user ID from middleware
	userID, exists := c.Get("validated_user_id")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized)
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		RespondWithError(c, http.StatusNotFound, ErrorCodeUserNotFound)
		return
	}

	user.Password = "" // Don't return password
	RespondWithSuccess(c, user)
}

// Helper function to get base URL
func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default for development
	}
	return baseURL
}

// Email verification request struct
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// Email verification handlers
func handleVerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithValidationError(c, err.Error())
		return
	}

	verificationService := auth.NewEmailVerificationService()
	user, err := verificationService.VerifyEmail(req.Token)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, ErrorCodeInvalidVerificationToken, err.Error())
		return
	}

	// Generate JWT token for immediate login
	jwtToken, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to generate token")
		return
	}

	user.Password = "" // Don't return password
	RespondWithSuccess(c, gin.H{
		"message": "Email verified successfully. You can now login.",
		"token":   jwtToken,
		"user":    user,
	})
}

type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func handleResendVerification(c *gin.Context) {
	var req ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithValidationError(c, err.Error())
		return
	}

	// Find unverified user
	var user models.User
	if err := database.DB.Where("email = ? AND email_verified = ?", req.Email, false).
		First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			RespondWithError(c, http.StatusNotFound, ErrorCodeUserNotFound, "User not found or already verified")
			return
		}
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeDatabaseError)
		return
	}

	// Check rate limits
	verificationService := auth.NewEmailVerificationService()
	clientIP := c.ClientIP()
	canSend, message, err := verificationService.CanSendVerificationEmail(req.Email, clientIP)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to check send limits")
		return
	}

	if !canSend {
		RespondWithError(c, http.StatusTooManyRequests, ErrorCodeTooManyRequests, message)
		return
	}

	// Generate new verification token
	token, err := verificationService.GenerateVerificationToken()
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to generate verification token")
		return
	}

	// Set verification token
	if err := verificationService.SetVerificationToken(user.ID, token); err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to set verification token")
		return
	}

	// Build verification URL
	baseURL := getBaseURL()
	verificationURL := fmt.Sprintf("%s/api/auth/verify/%s", baseURL, token)
	
	// Send email
	mailService := mailer.NewMailer()
	if err := mailService.SendVerificationEmail(req.Email, verificationURL); err != nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeEmailServiceError, "Failed to send verification email")
		return
	}

	// Record email sent
	if err := verificationService.RecordVerificationEmailSent(req.Email, clientIP); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to record verification email: %v\n", err)
	}

	RespondWithSuccess(c, gin.H{
		"message": "Verification email sent. Please check your inbox.",
	})
}