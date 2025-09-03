package api

import (
	"net/http"
	"pagemail/internal/auth"

	"github.com/gin-gonic/gin"
)

func handleUsageInfo(c *gin.Context) {
	// Get validated user ID from middleware
	userID, exists := c.Get("validated_user_id")
	if !exists {
		RespondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized)
		return
	}

	// Convert to proper type for usage function
	uid := userID.(uint)
	userIDPtr := &uid

	// Get usage information
	usageInfo := auth.GetUsageInfo(userIDPtr)
	if usageInfo == nil {
		RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Failed to retrieve usage information")
		return
	}

	RespondWithSuccess(c, gin.H{
		"usage": usageInfo,
	})
}
