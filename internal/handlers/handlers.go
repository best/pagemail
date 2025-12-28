package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"pagemail/internal/config"
	"pagemail/internal/middleware"
	"pagemail/internal/models"
	"pagemail/internal/storage"
)

type Handler struct {
	cfg     *config.Config
	db      *gorm.DB
	storage storage.Storage
}

func New(cfg *config.Config, db *gorm.DB, store storage.Storage) *Handler {
	return &Handler{cfg: cfg, db: db, storage: store}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

func (h *Handler) Ready(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": "database error"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": "database unreachable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://pagemail.app/errors/validation",
			"title":  "Validation Error",
			"status": 400,
			"detail": err.Error(),
		})
		return
	}

	var existingUser models.User
	if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"type":   "https://pagemail.app/errors/conflict",
			"title":  "Conflict",
			"status": 409,
			"detail": "Email already registered",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://pagemail.app/errors/internal",
			"title":  "Internal Error",
			"status": 500,
			"detail": "Failed to process password",
		})
		return
	}

	var userCount int64
	h.db.Model(&models.User{}).Count(&userCount)

	role := "user"
	if userCount == 0 {
		role = "admin"
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsActive:     true,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://pagemail.app/errors/internal",
			"title":  "Internal Error",
			"status": 500,
			"detail": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"type":   "https://pagemail.app/errors/validation",
			"title":  "Validation Error",
			"status": 400,
			"detail": err.Error(),
		})
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://pagemail.app/errors/unauthorized",
			"title":  "Unauthorized",
			"status": 401,
			"detail": "Invalid credentials",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"type":   "https://pagemail.app/errors/forbidden",
			"title":  "Forbidden",
			"status": 403,
			"detail": "Account is deactivated",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"type":   "https://pagemail.app/errors/unauthorized",
			"title":  "Unauthorized",
			"status": 401,
			"detail": "Invalid credentials",
		})
		return
	}

	now := time.Now()
	user.LastLoginAt = &now
	h.db.Save(&user)

	accessToken, err := h.generateToken(user, h.cfg.JWT.AccessExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://pagemail.app/errors/internal",
			"title":  "Internal Error",
			"status": 500,
			"detail": "Failed to generate token",
		})
		return
	}

	refreshToken, err := h.generateToken(user, h.cfg.JWT.RefreshExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"type":   "https://pagemail.app/errors/internal",
			"title":  "Internal Error",
			"status": 500,
			"detail": "Failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    int(h.cfg.JWT.AccessExpiry.Seconds()),
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *Handler) generateToken(user models.User, expiry time.Duration) (string, error) {
	claims := middleware.Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWT.Secret))
}

func (h *Handler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"detail": "Not implemented"})
}

func (h *Handler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")
	var user models.User
	if err := h.db.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"type":   "https://pagemail.app/errors/not-found",
			"title":  "Not Found",
			"status": 404,
			"detail": "User not found",
		})
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
