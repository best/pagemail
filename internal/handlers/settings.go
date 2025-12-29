package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pagemail/internal/models"
	"pagemail/internal/pkg/crypto"
	"pagemail/internal/pkg/errors"
)

type CreateSMTPRequest struct {
	Name      string `json:"name" binding:"required"`
	Host      string `json:"host" binding:"required"`
	Port      int    `json:"port" binding:"required"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FromEmail string `json:"from_email" binding:"required,email"`
	FromName  string `json:"from_name"`
	UseTLS    bool   `json:"use_tls"`
	IsDefault bool   `json:"is_default"`
}

func (h *Handler) ListSMTPProfiles(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var profiles []models.SMTPProfile
	if err := h.db.Where("user_id = ?", &uid).Find(&profiles).Error; err != nil {
		errors.InternalError("Failed to fetch SMTP profiles").Respond(c)
		return
	}

	result := make([]gin.H, len(profiles))
	for i := range profiles {
		result[i] = gin.H{
			"id":         profiles[i].ID,
			"name":       profiles[i].Name,
			"host":       profiles[i].Host,
			"port":       profiles[i].Port,
			"username":   profiles[i].Username,
			"from_email": profiles[i].FromEmail,
			"from_name":  profiles[i].FromName,
			"use_tls":    profiles[i].UseTLS,
			"is_default": profiles[i].IsDefault,
			"created_at": profiles[i].CreatedAt,
			"updated_at": profiles[i].UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateSMTPProfile(c *gin.Context) {
	var req CreateSMTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}

	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	encryptor, err := crypto.NewEncryptor(h.cfg.Encryption.Key)
	if err != nil {
		errors.InternalError("Encryption error").Respond(c)
		return
	}

	var encryptedPassword string
	if req.Password != "" {
		encrypted, err := encryptor.Encrypt([]byte(req.Password))
		if err != nil {
			errors.InternalError("Failed to encrypt password").Respond(c)
			return
		}
		encryptedPassword = string(encrypted)
	}

	if req.IsDefault {
		h.db.Model(&models.SMTPProfile{}).Where("user_id = ?", uid).Update("is_default", false)
	}

	profile := models.SMTPProfile{
		UserID:      &uid,
		Name:        req.Name,
		Host:        req.Host,
		Port:        req.Port,
		Username:    req.Username,
		PasswordEnc: []byte(encryptedPassword),
		FromEmail:   req.FromEmail,
		FromName:    req.FromName,
		UseTLS:      req.UseTLS,
		IsDefault:   req.IsDefault,
	}

	if err := h.db.Create(&profile).Error; err != nil {
		errors.InternalError("Failed to create SMTP profile").Respond(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         profile.ID,
		"name":       profile.Name,
		"host":       profile.Host,
		"port":       profile.Port,
		"username":   profile.Username,
		"from_email": profile.FromEmail,
		"from_name":  profile.FromName,
		"use_tls":    profile.UseTLS,
		"is_default": profile.IsDefault,
		"created_at": profile.CreatedAt,
	})
}

func (h *Handler) UpdateSMTPProfile(c *gin.Context) {
	profileID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var profile models.SMTPProfile
	if err := h.db.Where("id = ? AND user_id = ?", profileID, &uid).First(&profile).Error; err != nil {
		errors.NotFound("SMTP profile not found").Respond(c)
		return
	}

	var req CreateSMTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}

	profile.Name = req.Name
	profile.Host = req.Host
	profile.Port = req.Port
	profile.Username = req.Username
	profile.FromEmail = req.FromEmail
	profile.FromName = req.FromName
	profile.UseTLS = req.UseTLS

	if req.Password != "" {
		encryptor, _ := crypto.NewEncryptor(h.cfg.Encryption.Key)
		encrypted, _ := encryptor.Encrypt([]byte(req.Password))
		profile.PasswordEnc = encrypted
	}

	if req.IsDefault && !profile.IsDefault {
		h.db.Model(&models.SMTPProfile{}).Where("user_id = ? AND id != ?", uid, profile.ID).Update("is_default", false)
	}
	profile.IsDefault = req.IsDefault

	if err := h.db.Save(&profile).Error; err != nil {
		errors.InternalError("Failed to update SMTP profile").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         profile.ID,
		"name":       profile.Name,
		"host":       profile.Host,
		"port":       profile.Port,
		"username":   profile.Username,
		"from_email": profile.FromEmail,
		"from_name":  profile.FromName,
		"use_tls":    profile.UseTLS,
		"is_default": profile.IsDefault,
		"updated_at": profile.UpdatedAt,
	})
}

func (h *Handler) DeleteSMTPProfile(c *gin.Context) {
	profileID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	result := h.db.Where("id = ? AND user_id = ?", profileID, &uid).Delete(&models.SMTPProfile{})
	if result.RowsAffected == 0 {
		errors.NotFound("SMTP profile not found").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMTP profile deleted"})
}

func (h *Handler) TestSMTPProfile(c *gin.Context) {
	profileID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var profile models.SMTPProfile
	if err := h.db.Where("id = ? AND user_id = ?", profileID, &uid).First(&profile).Error; err != nil {
		errors.NotFound("SMTP profile not found").Respond(c)
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}

	// TODO: implement actual SMTP test using notify.SMTPSender
	c.JSON(http.StatusOK, gin.H{"message": "Test email sent to " + req.Email})
}

type CreateWebhookRequest struct {
	Name     string `json:"name" binding:"required"`
	URL      string `json:"url" binding:"required,url"`
	Secret   string `json:"secret"`
	IsActive bool   `json:"is_active"`
}

func (h *Handler) ListWebhooks(c *gin.Context) {
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var webhooks []models.WebhookEndpoint
	if err := h.db.Where("user_id = ?", uid).Find(&webhooks).Error; err != nil {
		errors.InternalError("Failed to fetch webhooks").Respond(c)
		return
	}

	result := make([]gin.H, len(webhooks))
	for i := range webhooks {
		result[i] = gin.H{
			"id":         webhooks[i].ID,
			"name":       webhooks[i].Name,
			"url":        webhooks[i].URL,
			"is_active":  webhooks[i].IsActive,
			"created_at": webhooks[i].CreatedAt,
			"updated_at": webhooks[i].UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) CreateWebhook(c *gin.Context) {
	var req CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}

	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var encryptedSecret string
	if req.Secret != "" {
		encryptor, _ := crypto.NewEncryptor(h.cfg.Encryption.Key)
		encrypted, _ := encryptor.Encrypt([]byte(req.Secret))
		encryptedSecret = string(encrypted)
	}

	webhook := models.WebhookEndpoint{
		UserID:   uid,
		Name:     req.Name,
		URL:      req.URL,
		Secret:   encryptedSecret,
		IsActive: req.IsActive,
	}

	if err := h.db.Create(&webhook).Error; err != nil {
		errors.InternalError("Failed to create webhook").Respond(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         webhook.ID,
		"name":       webhook.Name,
		"url":        webhook.URL,
		"is_active":  webhook.IsActive,
		"created_at": webhook.CreatedAt,
	})
}

func (h *Handler) UpdateWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var webhook models.WebhookEndpoint
	if err := h.db.Where("id = ? AND user_id = ?", webhookID, uid).First(&webhook).Error; err != nil {
		errors.NotFound("Webhook not found").Respond(c)
		return
	}

	var req CreateWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}

	webhook.Name = req.Name
	webhook.URL = req.URL
	webhook.IsActive = req.IsActive

	if req.Secret != "" {
		encryptor, _ := crypto.NewEncryptor(h.cfg.Encryption.Key)
		encrypted, _ := encryptor.Encrypt([]byte(req.Secret))
		webhook.Secret = string(encrypted)
	}

	if err := h.db.Save(&webhook).Error; err != nil {
		errors.InternalError("Failed to update webhook").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         webhook.ID,
		"name":       webhook.Name,
		"url":        webhook.URL,
		"is_active":  webhook.IsActive,
		"updated_at": webhook.UpdatedAt,
	})
}

func (h *Handler) DeleteWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	result := h.db.Where("id = ? AND user_id = ?", webhookID, uid).Delete(&models.WebhookEndpoint{})
	if result.RowsAffected == 0 {
		errors.NotFound("Webhook not found").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook deleted"})
}

func (h *Handler) TestWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	userID := c.GetString("user_id")
	uid, _ := uuid.Parse(userID)

	var webhook models.WebhookEndpoint
	if err := h.db.Where("id = ? AND user_id = ?", webhookID, uid).First(&webhook).Error; err != nil {
		errors.NotFound("Webhook not found").Respond(c)
		return
	}

	// TODO: implement actual webhook test using notify.WebhookSender
	c.JSON(http.StatusOK, gin.H{"message": "Test webhook sent"})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		errors.NotFound("User not found").Respond(c)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		errors.Unauthorized("Current password is incorrect").Respond(c)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		errors.InternalError("Failed to process password").Respond(c)
		return
	}

	user.PasswordHash = string(hashedPassword)
	if err := h.db.Save(&user).Error; err != nil {
		errors.InternalError("Failed to update password").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func (h *Handler) UpdateCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		errors.NotFound("User not found").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func parsePagination(c *gin.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return
}

func paginatedResponse(c *gin.Context, data interface{}, total int64, page, limit int) {
	totalPages := (int(total) + limit - 1) / limit
	c.JSON(http.StatusOK, gin.H{
		"data": data,
		"meta": gin.H{
			"total":       total,
			"page":        page,
			"per_page":    limit,
			"total_pages": totalPages,
		},
	})
}
