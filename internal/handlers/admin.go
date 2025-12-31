package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"pagemail/internal/models"
	"pagemail/internal/pkg/errors"
)

func (h *Handler) AdminListUsers(c *gin.Context) {
	page, limit := parsePagination(c)

	var users []models.User
	var total int64

	h.db.Model(&models.User{}).Count(&total)
	h.db.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&users)

	result := make([]gin.H, len(users))
	for i := range users {
		result[i] = gin.H{
			"id":            users[i].ID,
			"email":         users[i].Email,
			"role":          users[i].Role,
			"is_active":     users[i].IsActive,
			"last_login_at": users[i].LastLoginAt,
			"created_at":    users[i].CreatedAt,
			"updated_at":    users[i].UpdatedAt,
		}
	}

	paginatedResponse(c, result, total, page, limit)
}

func (h *Handler) AdminGetUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		errors.NotFound("User not found").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":            user.ID,
		"email":         user.Email,
		"role":          user.Role,
		"is_active":     user.IsActive,
		"last_login_at": user.LastLoginAt,
		"created_at":    user.CreatedAt,
		"updated_at":    user.UpdatedAt,
	})
}

type AdminUpdateUserRequest struct {
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

func (h *Handler) AdminUpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		errors.NotFound("User not found").Respond(c)
		return
	}

	var req AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}

	if req.Role != nil {
		if *req.Role != "admin" && *req.Role != "user" {
			errors.BadRequest("Invalid role").Respond(c)
			return
		}
		user.Role = *req.Role
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := h.db.Save(&user).Error; err != nil {
		errors.InternalError("Failed to update user").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"role":      user.Role,
		"is_active": user.IsActive,
	})
}

func (h *Handler) AdminDeleteUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("user_id")

	if userID == currentUserID {
		errors.BadRequest("Cannot delete yourself").Respond(c)
		return
	}

	result := h.db.Delete(&models.User{}, "id = ?", userID)
	if result.RowsAffected == 0 {
		errors.NotFound("User not found").Respond(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func (h *Handler) ListAuditLogs(c *gin.Context) {
	page, limit := parsePagination(c)

	var logs []models.AuditLog
	var total int64

	h.db.Model(&models.AuditLog{}).Count(&total)
	h.db.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&logs)

	result := make([]gin.H, len(logs))
	for i := range logs {
		result[i] = gin.H{
			"id":            logs[i].ID,
			"actor_id":      logs[i].ActorID,
			"actor_email":   logs[i].ActorEmail,
			"action":        logs[i].Action,
			"resource_type": logs[i].ResourceType,
			"resource_id":   logs[i].ResourceID,
			"details":       logs[i].Details,
			"ip_address":    logs[i].IPAddress,
			"created_at":    logs[i].CreatedAt,
		}
	}

	paginatedResponse(c, result, total, page, limit)
}

func (h *Handler) GetSystemStats(c *gin.Context) {
	var userCount, taskCount, completedCount, failedCount int64

	h.db.Model(&models.User{}).Count(&userCount)
	h.db.Model(&models.CaptureTask{}).Count(&taskCount)
	h.db.Model(&models.CaptureTask{}).Where("status = ?", "completed").Count(&completedCount)
	h.db.Model(&models.CaptureTask{}).Where("status = ?", "failed").Count(&failedCount)

	c.JSON(http.StatusOK, gin.H{
		"users": gin.H{
			"total": userCount,
		},
		"tasks": gin.H{
			"total":     taskCount,
			"completed": completedCount,
			"failed":    failedCount,
		},
	})
}

func (h *Handler) GetGlobalSMTP(c *gin.Context) {
	var admin models.User
	if err := h.db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		errors.NotFound("Admin user not found").Respond(c)
		return
	}

	var profiles []models.SMTPProfile
	h.db.Where("user_id = ? AND is_default = ?", admin.ID, true).Find(&profiles)

	if len(profiles) == 0 {
		c.JSON(http.StatusOK, nil)
		return
	}

	p := profiles[0]
	c.JSON(http.StatusOK, gin.H{
		"id":         p.ID,
		"name":       p.Name,
		"host":       p.Host,
		"port":       p.Port,
		"username":   p.Username,
		"from_email": p.FromEmail,
		"from_name":  p.FromName,
		"use_tls":    p.UseTLS,
	})
}

func (h *Handler) UpdateGlobalSMTP(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Use SMTP profile management instead"})
}

func (h *Handler) GetStorageConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"storage_type":            h.cfg.Storage.Backend,
		"storage_path":            h.cfg.Storage.LocalPath,
		"s3_bucket":               h.cfg.Storage.S3Bucket,
		"s3_region":               h.cfg.Storage.S3Region,
		"default_formats":         []string{"pdf"},
		"max_concurrent_captures": 5,
	})
}

func (h *Handler) UpdateStorageConfig(c *gin.Context) {
	// Note: This would require config file modification
	// For now, return the current config
	c.JSON(http.StatusOK, gin.H{
		"message": "Storage configuration should be updated via environment variables or config file",
	})
}

func (h *Handler) GetSiteConfig(c *gin.Context) {
	config, err := h.loadSiteConfig()
	if err != nil {
		errors.InternalError("Failed to load site config").Respond(c)
		return
	}
	c.JSON(http.StatusOK, config)
}

func (h *Handler) UpdateSiteConfig(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.BadRequest(err.Error()).Respond(c)
		return
	}
	if len(req) == 0 {
		errors.BadRequest("No configuration values provided").Respond(c)
		return
	}

	for key := range req {
		if !isAllowedSiteConfigKey(key) {
			errors.BadRequest("Unsupported site config key: " + key).Respond(c)
			return
		}
	}

	now := time.Now()
	err := h.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range req {
			setting := models.SystemSetting{
				Key:       key,
				Value:     value,
				UpdatedAt: now,
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "key"}},
				DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
			}).Create(&setting).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		errors.InternalError("Failed to update site config").Respond(c)
		return
	}

	config, err := h.loadSiteConfig()
	if err != nil {
		errors.InternalError("Failed to load site config").Respond(c)
		return
	}
	c.JSON(http.StatusOK, config)
}
