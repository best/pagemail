package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"pagemail/internal/config"
	"pagemail/internal/models"
	"pagemail/internal/storage"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.SMTPProfile{},
		&models.WebhookEndpoint{},
		&models.CaptureTask{},
		&models.CaptureOutput{},
		&models.Delivery{},
		&models.Job{},
		&models.AuditLog{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func setupTestHandler(t *testing.T) (*Handler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret-key",
		},
		Encryption: config.EncryptionConfig{
			Key: "test-encryption-key-32-bytes!!!!",
		},
		Storage: config.StorageConfig{
			Backend:   "local",
			LocalPath: t.TempDir(),
		},
	}

	store, err := storage.NewLocalStorage(cfg.Storage.LocalPath)
	if err != nil {
		t.Fatalf("Failed to create test storage: %v", err)
	}

	h := New(cfg, db, store)
	r := gin.New()

	return h, r
}

func TestHealth(t *testing.T) {
	h, r := setupTestHandler(t)
	r.GET("/health", h.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Health() status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Errorf("Health() status = %q, want %q", resp["status"], "healthy")
	}
}

func TestReady(t *testing.T) {
	h, r := setupTestHandler(t)
	r.GET("/ready", h.Ready)

	req := httptest.NewRequest(http.MethodGet, "/ready", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ready() status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRegister(t *testing.T) {
	h, r := setupTestHandler(t)
	r.POST("/register", h.Register)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid registration (first user becomes admin)",
			body: map[string]string{
				"email":    "admin@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "valid registration (second user)",
			body: map[string]string{
				"email":    "user@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "duplicate email",
			body: map[string]string{
				"email":    "admin@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusConflict,
		},
		{
			name: "invalid email",
			body: map[string]string{
				"email":    "not-an-email",
				"password": "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "short password",
			body: map[string]string{
				"email":    "new@example.com",
				"password": "short",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing fields",
			body:       map[string]string{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Register() status = %d, want %d, body = %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestLogin(t *testing.T) {
	h, r := setupTestHandler(t)
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// Register a user first
	regBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	tests := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name: "valid login",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "wrong password",
			body: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "non-existent user",
			body: map[string]string{
				"email":    "nobody@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing password",
			body: map[string]string{
				"email": "test@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Login() status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				if resp["access_token"] == nil {
					t.Error("Login() response missing access_token")
				}
			}
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	h, r := setupTestHandler(t)
	r.GET("/me", func(c *gin.Context) {
		c.Set("user_id", "invalid-uuid")
		h.GetCurrentUser(c)
	})

	req := httptest.NewRequest(http.MethodGet, "/me", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GetCurrentUser() with invalid ID status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestParsePagination(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		query     string
		wantPage  int
		wantLimit int
	}{
		{"default values", "", 1, 10},
		{"custom page", "page=5", 5, 10},
		{"custom limit", "limit=20", 1, 20},
		{"both custom", "page=3&limit=50", 3, 50},
		{"negative page", "page=-1", 1, 10},
		{"negative limit", "limit=-5", 1, 10},
		{"limit too large", "limit=200", 1, 10},
		{"invalid values", "page=abc&limit=xyz", 1, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			var gotPage, gotLimit int
			r.GET("/test", func(c *gin.Context) {
				gotPage, gotLimit = parsePagination(c)
			})

			req := httptest.NewRequest(http.MethodGet, "/test?"+tt.query, http.NoBody)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if gotPage != tt.wantPage {
				t.Errorf("parsePagination() page = %d, want %d", gotPage, tt.wantPage)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("parsePagination() limit = %d, want %d", gotLimit, tt.wantLimit)
			}
		})
	}
}
