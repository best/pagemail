package api

import (
	"net/http"
	"pagemail/internal/auth"

	"github.com/gin-gonic/gin"
)

func handleUsageInfo(c *gin.Context) {
	// Get user ID if authenticated
	userID, _ := c.Get("user_id")
	var userIDPtr *uint
	if uid, ok := userID.(uint); ok {
		userIDPtr = &uid
	}

	// Get usage information
	usageInfo := auth.GetUsageInfo(userIDPtr)
	if usageInfo == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve usage information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usage": usageInfo,
	})
}