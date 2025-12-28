package audit

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"pagemail/internal/models"
)

type Logger struct {
	db *gorm.DB
}

func NewLogger(db *gorm.DB) *Logger {
	return &Logger{db: db}
}

type LogEntry struct {
	ActorID      *uuid.UUID
	ActorEmail   string
	Action       string
	ResourceType string
	ResourceID   *uuid.UUID
	Details      interface{}
	IPAddress    string
	UserAgent    string
	TraceID      string
}

func (l *Logger) Log(entry *LogEntry) error {
	var detailsJSON string
	if entry.Details != nil {
		data, err := json.Marshal(entry.Details)
		if err == nil {
			detailsJSON = string(data)
		}
	}

	log := models.AuditLog{
		ActorID:      entry.ActorID,
		ActorEmail:   entry.ActorEmail,
		Action:       entry.Action,
		ResourceType: entry.ResourceType,
		ResourceID:   entry.ResourceID,
		Details:      detailsJSON,
		IPAddress:    entry.IPAddress,
		UserAgent:    entry.UserAgent,
		TraceID:      entry.TraceID,
	}

	return l.db.Create(&log).Error
}

func (l *Logger) LogFromContext(c *gin.Context, action string, resourceType string, resourceID *uuid.UUID, details interface{}) error {
	var actorID *uuid.UUID
	if userIDStr, exists := c.Get("user_id"); exists {
		if id, err := uuid.Parse(userIDStr.(string)); err == nil {
			actorID = &id
		}
	}

	actorEmail, _ := c.Get("user_email")

	entry := &LogEntry{
		ActorID:      actorID,
		ActorEmail:   actorEmail.(string),
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
		IPAddress:    c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		TraceID:      c.GetString("trace_id"),
	}

	return l.Log(entry)
}

const (
	ActionUserCreate     = "user.create"
	ActionUserUpdate     = "user.update"
	ActionUserDelete     = "user.delete"
	ActionUserLogin      = "user.login"
	ActionSettingsUpdate = "settings.update"
	ActionSMTPCreate     = "smtp.create"
	ActionSMTPUpdate     = "smtp.update"
	ActionSMTPDelete     = "smtp.delete"
	ActionWebhookCreate  = "webhook.create"
	ActionWebhookUpdate  = "webhook.update"
	ActionWebhookDelete  = "webhook.delete"
	ActionCaptureCreate  = "capture.create"
	ActionCaptureDelete  = "capture.delete"
	ActionDeliveryCreate = "delivery.create"
)
