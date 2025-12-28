package models

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = db.AutoMigrate(
		&User{},
		&SMTPProfile{},
		&WebhookEndpoint{},
		&CaptureTask{},
		&CaptureOutput{},
		&Delivery{},
		&Job{},
		&AuditLog{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestUserBeforeCreate(t *testing.T) {
	db := setupTestDB(t)

	user := User{
		Email:        "test@example.com",
		PasswordHash: "hash",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID == uuid.Nil {
		t.Error("User ID should be auto-generated")
	}
}

func TestUserBeforeCreateWithID(t *testing.T) {
	db := setupTestDB(t)

	existingID := uuid.New()
	user := User{
		ID:           existingID,
		Email:        "test2@example.com",
		PasswordHash: "hash",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID != existingID {
		t.Error("User ID should not be changed if already set")
	}
}

func TestCaptureTaskFormats(t *testing.T) {
	tests := []struct {
		name    string
		formats int
		hasPDF  bool
		hasHTML bool
		hasPNG  bool
	}{
		{"none", 0, false, false, false},
		{"pdf only", FormatPDF, true, false, false},
		{"html only", FormatHTML, false, true, false},
		{"png only", FormatPNG, false, false, true},
		{"all", FormatPDF | FormatHTML | FormatPNG, true, true, true},
		{"pdf and html", FormatPDF | FormatHTML, true, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.formats&FormatPDF != 0) != tt.hasPDF {
				t.Errorf("PDF check failed for %d", tt.formats)
			}
			if (tt.formats&FormatHTML != 0) != tt.hasHTML {
				t.Errorf("HTML check failed for %d", tt.formats)
			}
			if (tt.formats&FormatPNG != 0) != tt.hasPNG {
				t.Errorf("PNG check failed for %d", tt.formats)
			}
		})
	}
}

func TestTaskStatusConstants(t *testing.T) {
	statuses := []string{
		TaskStatusPending,
		TaskStatusRunning,
		TaskStatusCompleted,
		TaskStatusFailed,
	}

	seen := make(map[string]bool)
	for _, s := range statuses {
		if seen[s] {
			t.Errorf("Duplicate status constant: %s", s)
		}
		seen[s] = true
	}
}

func TestJobStatusConstants(t *testing.T) {
	statuses := []string{
		JobStatusPending,
		JobStatusRunning,
		JobStatusSuccess,
		JobStatusFailed,
	}

	seen := make(map[string]bool)
	for _, s := range statuses {
		if seen[s] {
			t.Errorf("Duplicate status constant: %s", s)
		}
		seen[s] = true
	}
}

func TestDeliveryChannelConstants(t *testing.T) {
	if ChannelEmail != "email" {
		t.Errorf("ChannelEmail = %q, want %q", ChannelEmail, "email")
	}
	if ChannelWebhook != "webhook" {
		t.Errorf("ChannelWebhook = %q, want %q", ChannelWebhook, "webhook")
	}
}

func TestCaptureTaskCreate(t *testing.T) {
	db := setupTestDB(t)

	// Create user first
	user := User{Email: "test@example.com", PasswordHash: "hash"}
	db.Create(&user)

	task := CaptureTask{
		UserID:  user.ID,
		URL:     "https://example.com",
		Formats: FormatPDF | FormatHTML,
	}

	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	if task.ID == uuid.Nil {
		t.Error("Task ID should be auto-generated")
	}

	if task.Status != TaskStatusPending {
		t.Errorf("Task status = %q, want %q", task.Status, TaskStatusPending)
	}
}

func TestSMTPProfileNullableUserID(t *testing.T) {
	db := setupTestDB(t)

	// Global SMTP profile (no user)
	globalProfile := SMTPProfile{
		Name:      "Global SMTP",
		Host:      "smtp.example.com",
		Port:      587,
		FromEmail: "noreply@example.com",
	}

	if err := db.Create(&globalProfile).Error; err != nil {
		t.Fatalf("Failed to create global SMTP profile: %v", err)
	}

	if globalProfile.UserID != nil {
		t.Error("Global SMTP profile should have nil UserID")
	}

	// User-specific SMTP profile
	user := User{Email: "test@example.com", PasswordHash: "hash"}
	db.Create(&user)

	userProfile := SMTPProfile{
		UserID:    &user.ID,
		Name:      "User SMTP",
		Host:      "smtp.user.com",
		Port:      587,
		FromEmail: "user@example.com",
	}

	if err := db.Create(&userProfile).Error; err != nil {
		t.Fatalf("Failed to create user SMTP profile: %v", err)
	}

	if userProfile.UserID == nil || *userProfile.UserID != user.ID {
		t.Error("User SMTP profile should have correct UserID")
	}
}
