package auth

import (
	"fmt"
	"pagemail/internal/database"
	"pagemail/internal/models"
	"time"

	"github.com/google/uuid"
)

// EmailVerificationService handles email verification logic
type EmailVerificationService struct{}

// NewEmailVerificationService creates a new email verification service
func NewEmailVerificationService() *EmailVerificationService {
	return &EmailVerificationService{}
}

// GenerateVerificationToken generates a secure verification token
func (s *EmailVerificationService) GenerateVerificationToken() (string, error) {
	// Use UUID v4 for security
	token := uuid.New().String()
	return token, nil
}

// CanSendVerificationEmail checks if verification email can be sent based on rate limits
func (s *EmailVerificationService) CanSendVerificationEmail(email, ipAddress string) (bool, string, error) {
	now := time.Now()

	// Check email-based rate limit (5 minutes)
	emailCooldown := now.Add(-5 * time.Minute)
	var emailCount int64
	if err := database.DB.Model(&models.EmailVerification{}).
		Where("email = ? AND sent_at > ?", email, emailCooldown).
		Count(&emailCount).Error; err != nil {
		return false, "", fmt.Errorf("failed to check email rate limit: %w", err)
	}

	if emailCount > 0 {
		return false, "同一邮箱5分钟内只能发送1次验证邮件", nil
	}

	// Check IP-based rate limit (1 hour, max 5 emails)
	ipCooldown := now.Add(-1 * time.Hour)
	var ipCount int64
	if err := database.DB.Model(&models.EmailVerification{}).
		Where("ip_address = ? AND sent_at > ?", ipAddress, ipCooldown).
		Count(&ipCount).Error; err != nil {
		return false, "", fmt.Errorf("failed to check IP rate limit: %w", err)
	}

	if ipCount >= 5 {
		return false, "同一IP地址1小时内最多发送5次验证邮件", nil
	}

	return true, "", nil
}

// RecordVerificationEmailSent records that a verification email was sent
func (s *EmailVerificationService) RecordVerificationEmailSent(email, ipAddress string) error {
	verification := models.EmailVerification{
		Email:     email,
		IPAddress: ipAddress,
		SentAt:    time.Now(),
	}

	if err := database.DB.Create(&verification).Error; err != nil {
		return fmt.Errorf("failed to record verification email: %w", err)
	}

	return nil
}

// SetVerificationToken sets verification token for a user
func (s *EmailVerificationService) SetVerificationToken(userID uint, token string) error {
	// Set token expiry to 24 hours from now
	expires := time.Now().Add(24 * time.Hour)

	if err := database.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"email_verify_token":   token,
			"email_verify_expires": expires,
		}).Error; err != nil {
		return fmt.Errorf("failed to set verification token: %w", err)
	}

	return nil
}

// VerifyEmail verifies a user's email using the verification token
func (s *EmailVerificationService) VerifyEmail(token string) (*models.User, error) {
	var user models.User

	// Find user by token and check if token is not expired
	if err := database.DB.Where("email_verify_token = ? AND email_verify_expires > ?",
		token, time.Now()).First(&user).Error; err != nil {
		return nil, fmt.Errorf("invalid or expired verification token")
	}

	// Update user as verified and active
	updates := map[string]interface{}{
		"email_verified":       true,
		"is_active":            true,
		"email_verify_token":   nil,
		"email_verify_expires": nil,
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to verify email: %w", err)
	}

	// Refresh user data
	if err := database.DB.First(&user, user.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh user data: %w", err)
	}

	return &user, nil
}

// CleanupExpiredTokens removes expired verification tokens (should be run periodically)
func (s *EmailVerificationService) CleanupExpiredTokens() error {
	now := time.Now()

	// Clear expired tokens
	if err := database.DB.Model(&models.User{}).
		Where("email_verify_expires IS NOT NULL AND email_verify_expires < ?", now).
		Updates(map[string]interface{}{
			"email_verify_token":   nil,
			"email_verify_expires": nil,
		}).Error; err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	// Remove old verification records (older than 7 days)
	weekAgo := now.Add(-7 * 24 * time.Hour)
	if err := database.DB.Where("sent_at < ?", weekAgo).
		Delete(&models.EmailVerification{}).Error; err != nil {
		return fmt.Errorf("failed to cleanup old verification records: %w", err)
	}

	return nil
}
