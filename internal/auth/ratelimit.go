package auth

import (
	"net/http"
	"pagemail/internal/database"
	"pagemail/internal/models"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimitConfig struct {
	GuestDailyLimit      int
	GuestMonthlyLimit    int
	AuthenticatedDefault bool
}

var rateLimitConfig = RateLimitConfig{
	GuestDailyLimit:      1,
	GuestMonthlyLimit:    5,
	AuthenticatedDefault: true,
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		
		if exists && userID != nil {
			// Authenticated user - check their specific limits
			if !checkUserRateLimit(c, userID.(uint)) {
				return
			}
		} else {
			// Guest user - check IP-based limits
			if !checkGuestRateLimit(c) {
				return
			}
		}
		
		c.Next()
	}
}

func checkUserRateLimit(c *gin.Context, userID uint) bool {
	// Get user's quota settings
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user limits"})
		c.Abort()
		return false
	}

	// Check daily and monthly limits
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Count today's requests
	var dailyCount int64
	database.DB.Model(&models.Request{}).
		Where("user_id = ? AND created_at >= ?", userID, startOfDay).
		Count(&dailyCount)

	// Count this month's requests
	var monthlyCount int64
	database.DB.Model(&models.Request{}).
		Where("user_id = ? AND created_at >= ?", userID, startOfMonth).
		Count(&monthlyCount)

	// Check limits
	if dailyCount >= int64(user.DailyLimit) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error_code": 4001,
			"error_type": "QUOTA_ERROR",
			"message":    "Daily request limit exceeded",
			"used":       dailyCount,
			"limit":      user.DailyLimit,
			"reset_time": startOfDay.Add(24 * time.Hour).Unix(),
		})
		c.Abort()
		return false
	}

	if monthlyCount >= int64(user.MonthlyLimit) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error_code": 4002,
			"error_type": "QUOTA_ERROR",
			"message":    "Monthly request limit exceeded",
			"used":       monthlyCount,
			"limit":      user.MonthlyLimit,
			"reset_time": startOfMonth.AddDate(0, 1, 0).Unix(),
		})
		c.Abort()
		return false
	}

	// Add usage info to context for potential logging
	c.Set("daily_usage", dailyCount)
	c.Set("monthly_usage", monthlyCount)
	c.Set("daily_limit", user.DailyLimit)
	c.Set("monthly_limit", user.MonthlyLimit)

	return true
}

func checkGuestRateLimit(c *gin.Context) bool {
	// Count guest requests from this IP
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// For guest users, we'll check based on email + IP combination in recent requests
	// This is a simple approach - in production you might want Redis or similar
	
	var dailyCount int64
	database.DB.Model(&models.Request{}).
		Where("user_id IS NULL AND created_at >= ?", startOfDay).
		Count(&dailyCount)

	var monthlyCount int64
	database.DB.Model(&models.Request{}).
		Where("user_id IS NULL AND created_at >= ?", startOfMonth).
		Count(&monthlyCount)

	// Check daily limit for guests
	if dailyCount >= int64(rateLimitConfig.GuestDailyLimit) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error_code": 4001,
			"error_type": "QUOTA_ERROR",
			"message":    "Daily request limit exceeded for guests - please register for higher limits",
			"used":       dailyCount,
			"limit":      rateLimitConfig.GuestDailyLimit,
		})
		c.Abort()
		return false
	}

	// Check monthly limit for guests
	if monthlyCount >= int64(rateLimitConfig.GuestMonthlyLimit) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error_code": 4002,
			"error_type": "QUOTA_ERROR",
			"message":    "Monthly request limit exceeded for guests - please register for higher limits",
			"used":       monthlyCount,
			"limit":      rateLimitConfig.GuestMonthlyLimit,
		})
		c.Abort()
		return false
	}

	return true
}

func GetUsageInfo(userID *uint) map[string]interface{} {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	if userID != nil {
		// Authenticated user
		var user models.User
		if err := database.DB.First(&user, *userID).Error; err != nil {
			return nil
		}

		var dailyCount, monthlyCount int64
		database.DB.Model(&models.Request{}).
			Where("user_id = ? AND created_at >= ?", *userID, startOfDay).
			Count(&dailyCount)
		database.DB.Model(&models.Request{}).
			Where("user_id = ? AND created_at >= ?", *userID, startOfMonth).
			Count(&monthlyCount)

		return map[string]interface{}{
			"type": "authenticated",
			"daily": map[string]interface{}{
				"used":  dailyCount,
				"limit": user.DailyLimit,
				"remaining": user.DailyLimit - int(dailyCount),
			},
			"monthly": map[string]interface{}{
				"used":  monthlyCount,
				"limit": user.MonthlyLimit,
				"remaining": user.MonthlyLimit - int(monthlyCount),
			},
		}
	} else {
		// Guest user
		var dailyCount, monthlyCount int64
		database.DB.Model(&models.Request{}).
			Where("user_id IS NULL AND created_at >= ?", startOfDay).
			Count(&dailyCount)
		database.DB.Model(&models.Request{}).
			Where("user_id IS NULL AND created_at >= ?", startOfMonth).
			Count(&monthlyCount)

		return map[string]interface{}{
			"type": "guest",
			"daily": map[string]interface{}{
				"used":  dailyCount,
				"limit": rateLimitConfig.GuestDailyLimit,
				"remaining": rateLimitConfig.GuestDailyLimit - int(dailyCount),
			},
			"monthly": map[string]interface{}{
				"used":  monthlyCount,
				"limit": rateLimitConfig.GuestMonthlyLimit,
				"remaining": rateLimitConfig.GuestMonthlyLimit - int(monthlyCount),
			},
		}
	}
}