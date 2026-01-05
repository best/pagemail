package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         string         `gorm:"not null;default:user" json:"role"`
	IsActive     bool           `gorm:"not null;default:true" json:"is_active"`
	AvatarKey    *string        `gorm:"type:text" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	LastLoginAt  *time.Time     `json:"last_login_at,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type RolePermission struct {
	Role         string     `gorm:"primary_key" json:"role"`
	PermissionID uuid.UUID  `gorm:"type:uuid;primary_key" json:"permission_id"`
	Permission   Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

type SystemSetting struct {
	Key       string    `gorm:"primaryKey" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SMTPProfile struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID      *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	User        *User      `gorm:"foreignKey:UserID" json:"-"`
	Name        string     `gorm:"not null" json:"name"`
	Host        string     `gorm:"not null" json:"host"`
	Port        int        `gorm:"not null" json:"port"`
	Username    string     `json:"username"`
	PasswordEnc []byte     `json:"-"`
	FromName    string     `json:"from_name"`
	FromEmail   string     `gorm:"not null" json:"from_email"`
	UseTLS      bool       `gorm:"not null;default:true" json:"use_tls"`
	IsDefault   bool       `gorm:"not null;default:false" json:"is_default"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (s *SMTPProfile) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type WebhookEndpoint struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Name      string    `gorm:"not null" json:"name"`
	URL       string    `gorm:"not null" json:"url"`
	Secret    string    `json:"-"`
	Headers   string    `gorm:"type:text" json:"headers"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (w *WebhookEndpoint) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

const (
	FormatPDF  = 1
	FormatHTML = 2
	FormatPNG  = 4
)

const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
)

type CaptureTask struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key" json:"id"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	User           User            `gorm:"foreignKey:UserID" json:"-"`
	URL            string          `gorm:"not null" json:"url"`
	Status         string          `gorm:"not null;default:pending;index" json:"status"`
	Formats        int             `gorm:"not null;default:1" json:"formats"`
	CookiesEnc     []byte          `json:"-"`
	UserAgent      string          `json:"user_agent,omitempty"`
	ViewportWidth  int             `gorm:"default:1920" json:"viewport_width"`
	ViewportHeight int             `gorm:"default:1080" json:"viewport_height"`
	WaitTimeoutMs  int             `gorm:"default:30000" json:"wait_timeout_ms"`
	Attempts       int             `gorm:"not null;default:0" json:"attempts"`
	MaxAttempts    int             `gorm:"not null;default:3" json:"max_attempts"`
	ErrorMessage   string          `json:"error_message,omitempty"`
	CreatedAt      time.Time       `gorm:"index" json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	CompletedAt    *time.Time      `json:"completed_at,omitempty"`
	Outputs        []CaptureOutput `gorm:"foreignKey:TaskID" json:"outputs,omitempty"`
	Deliveries     []Delivery      `gorm:"foreignKey:TaskID" json:"deliveries,omitempty"`
}

func (c *CaptureTask) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type CaptureOutput struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	TaskID         uuid.UUID `gorm:"type:uuid;not null;index" json:"task_id"`
	Format         string    `gorm:"not null" json:"format"`
	StorageBackend string    `gorm:"not null" json:"storage_backend"`
	ObjectKey      string    `gorm:"not null" json:"object_key"`
	ContentType    string    `gorm:"not null" json:"content_type"`
	SizeBytes      int64     `gorm:"not null" json:"size_bytes"`
	SHA256         string    `json:"sha256,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

func (c *CaptureOutput) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

const (
	ChannelEmail   = "email"
	ChannelWebhook = "webhook"
)

const (
	DeliveryStatusPending = "pending"
	DeliveryStatusSent    = "sent"
	DeliveryStatusFailed  = "failed"
)

type Delivery struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	TaskID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"task_id"`
	Channel      string     `gorm:"not null" json:"channel"`
	TargetConfig string     `gorm:"type:text;not null" json:"target_config"`
	Status       string     `gorm:"not null;default:pending" json:"status"`
	Attempts     int        `gorm:"not null;default:0" json:"attempts"`
	MaxAttempts  int        `gorm:"not null;default:3" json:"max_attempts"`
	LastError    string     `json:"last_error,omitempty"`
	NextRetryAt  *time.Time `json:"next_retry_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

func (d *Delivery) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

const (
	JobTypeCapture = "capture"
	JobTypeDeliver = "deliver"
)

const (
	JobStatusPending = "pending"
	JobStatusRunning = "running"
	JobStatusSuccess = "succeeded"
	JobStatusFailed  = "failed"
)

type Job struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Type        string     `gorm:"not null;index" json:"type"`
	Payload     string     `gorm:"type:text;not null" json:"payload"`
	Status      string     `gorm:"not null;default:pending;index" json:"status"`
	Priority    int        `gorm:"not null;default:0" json:"priority"`
	RunAt       time.Time  `gorm:"not null;index" json:"run_at"`
	LockedBy    string     `json:"locked_by,omitempty"`
	LockedAt    *time.Time `json:"locked_at,omitempty"`
	LeaseUntil  *time.Time `json:"lease_until,omitempty"`
	Attempts    int        `gorm:"not null;default:0" json:"attempts"`
	MaxAttempts int        `gorm:"not null;default:3" json:"max_attempts"`
	LastError   string     `json:"last_error,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}

type AuditLog struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	ActorID      *uuid.UUID `gorm:"type:uuid;index" json:"actor_id,omitempty"`
	ActorEmail   string     `json:"actor_email,omitempty"`
	Action       string     `gorm:"not null;index" json:"action"`
	ResourceType string     `gorm:"index" json:"resource_type,omitempty"`
	ResourceID   *uuid.UUID `gorm:"type:uuid" json:"resource_id,omitempty"`
	Details      string     `gorm:"type:text" json:"details,omitempty"`
	IPAddress    string     `json:"ip_address,omitempty"`
	UserAgent    string     `json:"user_agent,omitempty"`
	TraceID      string     `json:"trace_id,omitempty"`
	CreatedAt    time.Time  `gorm:"index" json:"created_at"`
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
