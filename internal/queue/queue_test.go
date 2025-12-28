package queue

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"pagemail/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = db.AutoMigrate(&models.Job{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestEnqueueJob(t *testing.T) {
	db := setupTestDB(t)

	tests := []struct {
		name    string
		jobType string
		payload interface{}
		wantErr bool
	}{
		{
			name:    "capture job",
			jobType: models.JobTypeCapture,
			payload: map[string]interface{}{
				"task_id": "123",
				"url":     "https://example.com",
			},
			wantErr: false,
		},
		{
			name:    "deliver job",
			jobType: models.JobTypeDeliver,
			payload: map[string]interface{}{
				"delivery_id": "456",
			},
			wantErr: false,
		},
		{
			name:    "string payload",
			jobType: models.JobTypeCapture,
			payload: "simple string",
			wantErr: false,
		},
		{
			name:    "complex nested payload",
			jobType: models.JobTypeCapture,
			payload: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "value",
				},
				"array": []string{"a", "b", "c"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnqueueJob(db, tt.jobType, tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnqueueJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Verify jobs were created
	var count int64
	db.Model(&models.Job{}).Count(&count)
	if count != 4 {
		t.Errorf("Expected 4 jobs, got %d", count)
	}
}

func TestEnqueueJobDefaults(t *testing.T) {
	db := setupTestDB(t)

	err := EnqueueJob(db, models.JobTypeCapture, map[string]string{"test": "data"})
	if err != nil {
		t.Fatalf("EnqueueJob() error = %v", err)
	}

	var job models.Job
	if err := db.First(&job).Error; err != nil {
		t.Fatalf("Failed to retrieve job: %v", err)
	}

	if job.Status != models.JobStatusPending {
		t.Errorf("Job status = %q, want %q", job.Status, models.JobStatusPending)
	}

	if job.MaxAttempts != 3 {
		t.Errorf("Job MaxAttempts = %d, want 3", job.MaxAttempts)
	}

	if job.Attempts != 0 {
		t.Errorf("Job Attempts = %d, want 0", job.Attempts)
	}

	if job.Type != models.JobTypeCapture {
		t.Errorf("Job Type = %q, want %q", job.Type, models.JobTypeCapture)
	}
}
