package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"pagemail/internal/audit"
	"pagemail/internal/models"
	"pagemail/internal/pkg/errors"
	"pagemail/internal/queue"
)

const (
	formatPDF        = "pdf"
	formatHTML       = "html"
	formatScreenshot = "screenshot"
)

type CreateCaptureRequest struct {
	URL            string          `json:"url" binding:"required,url"`
	Formats        []string        `json:"formats" binding:"required,min=1"`
	Cookies        string          `json:"cookies"`
	DeliveryConfig *DeliveryConfig `json:"delivery_config"`
}

type DeliveryConfig struct {
	Type string `json:"type" binding:"required,oneof=email webhook"`
	ID   string `json:"id" binding:"required"`
}

func formatsToInt(formats []string) int {
	result := 0
	for _, f := range formats {
		switch f {
		case formatPDF:
			result |= models.FormatPDF
		case formatHTML:
			result |= models.FormatHTML
		case formatScreenshot:
			result |= models.FormatPNG
		}
	}
	return result
}

func intToFormats(flags int) []string {
	var formats []string
	if flags&models.FormatPDF != 0 {
		formats = append(formats, formatPDF)
	}
	if flags&models.FormatHTML != 0 {
		formats = append(formats, formatHTML)
	}
	if flags&models.FormatPNG != 0 {
		formats = append(formats, formatScreenshot)
	}
	return formats
}

func (h *Handler) CreateCapture(c *gin.Context) {
	var req CreateCaptureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}

	for _, f := range req.Formats {
		if f != "pdf" && f != "html" && f != "screenshot" {
			errors.BadRequest("Invalid format: " + f).Respond(c)
			return
		}
	}

	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	task := models.CaptureTask{
		UserID:  uid,
		URL:     req.URL,
		Formats: formatsToInt(req.Formats),
		Status:  models.TaskStatusPending,
	}

	if req.Cookies != "" {
		task.CookiesEnc = []byte(req.Cookies) // TODO: encrypt
	}

	if err := h.db.Create(&task).Error; err != nil {
		errors.InternalError("Failed to create capture task").Respond(c)
		return
	}

	payload := map[string]interface{}{
		"task_id": task.ID.String(),
		"url":     req.URL,
		"cookies": req.Cookies,
		"formats": req.Formats,
	}

	if err := queue.EnqueueJob(h.db, models.JobTypeCapture, payload); err != nil {
		errors.InternalError("Failed to enqueue job").Respond(c)
		return
	}

	h.logAudit(c, audit.ActionCaptureCreate, "capture", &task.ID, audit.ResourceDetails{
		URL: task.URL, Formats: req.Formats,
	})

	c.JSON(http.StatusCreated, gin.H{
		"id":         task.ID,
		"url":        task.URL,
		"formats":    req.Formats,
		"status":     task.Status,
		"created_at": task.CreatedAt,
	})
}

func (h *Handler) ListCaptures(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)
	page, limit := parsePagination(c)
	status := c.Query("status")

	var tasks []models.CaptureTask
	var total int64

	query := h.db.Model(&models.CaptureTask{}).Where("user_id = ?", uid)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&tasks)

	result := make([]gin.H, len(tasks))
	for i := range tasks {
		result[i] = gin.H{
			"id":         tasks[i].ID,
			"url":        tasks[i].URL,
			"formats":    intToFormats(tasks[i].Formats),
			"status":     tasks[i].Status,
			"error":      tasks[i].ErrorMessage,
			"created_at": tasks[i].CreatedAt,
			"updated_at": tasks[i].UpdatedAt,
		}
	}

	paginatedResponse(c, result, total, page, limit)
}

func (h *Handler) GetCapture(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var task models.CaptureTask
	if err := h.db.Where("id = ? AND user_id = ?", taskID, uid).First(&task).Error; err != nil {
		errors.NotFound("Capture task not found").Respond(c)
		return
	}

	var outputs []models.CaptureOutput
	h.db.Where("task_id = ?", task.ID).Find(&outputs)

	outputResult := make([]gin.H, len(outputs))
	for i := range outputs {
		outputResult[i] = gin.H{
			"id":     outputs[i].ID,
			"format": outputs[i].Format,
			"size":   outputs[i].SizeBytes,
			"path":   outputs[i].ObjectKey,
		}
	}

	var deliveries []models.Delivery
	h.db.Where("task_id = ?", task.ID).Order("created_at DESC").Find(&deliveries)

	deliveryHistory := make([]gin.H, len(deliveries))
	for i := range deliveries {
		deliveryHistory[i] = gin.H{
			"channel":      deliveries[i].Channel,
			"status":       deliveries[i].Status,
			"attempt_time": deliveries[i].CreatedAt,
			"error":        deliveries[i].LastError,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":               task.ID,
		"url":              task.URL,
		"formats":          intToFormats(task.Formats),
		"status":           task.Status,
		"error_message":    task.ErrorMessage,
		"outputs":          outputResult,
		"delivery_history": deliveryHistory,
		"created_at":       task.CreatedAt,
		"updated_at":       task.UpdatedAt,
	})
}

func (h *Handler) RetryCapture(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var task models.CaptureTask
	if err := h.db.Where("id = ? AND user_id = ?", taskID, uid).First(&task).Error; err != nil {
		errors.NotFound("Capture task not found").Respond(c)
		return
	}

	if task.Status != models.TaskStatusFailed {
		errors.BadRequest("Only failed tasks can be retried").Respond(c)
		return
	}

	task.Status = models.TaskStatusPending
	task.ErrorMessage = ""
	h.db.Save(&task)

	formats := intToFormats(task.Formats)
	cookies := ""
	if task.CookiesEnc != nil {
		cookies = string(task.CookiesEnc) // TODO: decrypt
	}

	payload := map[string]interface{}{
		"task_id": task.ID.String(),
		"url":     task.URL,
		"cookies": cookies,
		"formats": formats,
	}

	if err := queue.EnqueueJob(h.db, models.JobTypeCapture, payload); err != nil {
		log.Error().Err(err).Str("task_id", task.ID.String()).Msg("Failed to enqueue capture job")
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     task.ID,
		"status": task.Status,
	})
}

func (h *Handler) DeleteCapture(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var task models.CaptureTask
	if err := h.db.Where("id = ? AND user_id = ?", taskID, uid).First(&task).Error; err != nil {
		errors.NotFound("Capture task not found").Respond(c)
		return
	}

	// Delete related records first (foreign key constraints)
	h.db.Where("task_id = ?", task.ID).Delete(&models.CaptureOutput{})
	h.db.Where("task_id = ?", task.ID).Delete(&models.Delivery{})

	if err := h.db.Delete(&task).Error; err != nil {
		errors.InternalError("Failed to delete task").Respond(c)
		return
	}

	h.logAudit(c, audit.ActionCaptureDelete, "capture", &task.ID, audit.ResourceDetails{URL: task.URL})

	c.JSON(http.StatusOK, gin.H{"message": "Capture task deleted"})
}

func (h *Handler) ListCaptureOutputs(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var task models.CaptureTask
	if err := h.db.Where("id = ? AND user_id = ?", taskID, uid).First(&task).Error; err != nil {
		errors.NotFound("Capture task not found").Respond(c)
		return
	}

	var outputs []models.CaptureOutput
	h.db.Where("task_id = ?", task.ID).Find(&outputs)

	result := make([]gin.H, len(outputs))
	for i := range outputs {
		result[i] = gin.H{
			"id":         outputs[i].ID,
			"format":     outputs[i].Format,
			"size":       outputs[i].SizeBytes,
			"created_at": outputs[i].CreatedAt,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) DownloadOutput(c *gin.Context) {
	taskID := c.Param("id")
	outputID := c.Param("oid")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var task models.CaptureTask
	if err := h.db.Where("id = ? AND user_id = ?", taskID, uid).First(&task).Error; err != nil {
		errors.NotFound("Capture task not found").Respond(c)
		return
	}

	var output models.CaptureOutput
	if err := h.db.Where("id = ? AND task_id = ?", outputID, task.ID).First(&output).Error; err != nil {
		errors.NotFound("Output not found").Respond(c)
		return
	}

	reader, info, err := h.storage.Download(c, output.ObjectKey)
	if err != nil {
		errors.InternalError("Failed to download file").Respond(c)
		return
	}
	defer reader.Close()

	c.Header("Content-Type", output.ContentType)
	c.Header("Content-Disposition", "attachment; filename="+output.Format+getExtension(output.Format))
	c.Header("Content-Length", fmt.Sprintf("%d", info.Size))
	c.DataFromReader(http.StatusOK, info.Size, output.ContentType, reader, nil)
}

func (h *Handler) PreviewOutput(c *gin.Context) {
	taskID := c.Param("id")
	outputID := c.Param("oid")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var task models.CaptureTask
	if err := h.db.Where("id = ? AND user_id = ?", taskID, uid).First(&task).Error; err != nil {
		errors.NotFound("Capture task not found").Respond(c)
		return
	}

	var output models.CaptureOutput
	if err := h.db.Where("id = ? AND task_id = ?", outputID, task.ID).First(&output).Error; err != nil {
		errors.NotFound("Output not found").Respond(c)
		return
	}

	if output.Format != formatPDF && output.Format != formatScreenshot {
		errors.NewProblemDetail(http.StatusUnsupportedMediaType, "Unsupported Media Type", "Preview is only supported for PDF and screenshot formats").Respond(c)
		return
	}

	reader, info, err := h.storage.Download(c, output.ObjectKey)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			errors.NotFound("File not found in storage").Respond(c)
		} else {
			errors.InternalError("Failed to download file").Respond(c)
		}
		return
	}
	defer reader.Close()

	contentType := getContentType(output.Format)
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "inline; filename="+output.Format+getExtension(output.Format))
	c.Header("Content-Length", fmt.Sprintf("%d", info.Size))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "private, no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.DataFromReader(http.StatusOK, info.Size, contentType, reader, nil)
}

func getContentType(format string) string {
	switch format {
	case formatPDF:
		return "application/pdf"
	case formatScreenshot:
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

func getExtension(format string) string {
	switch format {
	case "pdf":
		return ".pdf"
	case "html":
		return ".html"
	case "screenshot":
		return ".png"
	default:
		return ""
	}
}

func (h *Handler) DeliverCapture(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"detail": "Use delivery config in task creation"})
}

func (h *Handler) ListDeliveries(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var task models.CaptureTask
	if err := h.db.Where("id = ? AND user_id = ?", taskID, uid).First(&task).Error; err != nil {
		errors.NotFound("Capture task not found").Respond(c)
		return
	}

	var deliveries []models.Delivery
	h.db.Where("task_id = ?", task.ID).Order("created_at DESC").Find(&deliveries)

	result := make([]gin.H, len(deliveries))
	for i := range deliveries {
		result[i] = gin.H{
			"id":         deliveries[i].ID,
			"channel":    deliveries[i].Channel,
			"status":     deliveries[i].Status,
			"error":      deliveries[i].LastError,
			"created_at": deliveries[i].CreatedAt,
		}
	}

	c.JSON(http.StatusOK, result)
}
