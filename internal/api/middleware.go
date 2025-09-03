package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ValidateUserIDMiddleware ensures that the user can only access their own resources
func ValidateUserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the user_id from URL parameters
		userIDParam := c.Param("user_id")
		if userIDParam == "" {
			RespondWithError(c, http.StatusBadRequest, ErrorCodeMissingParameter, "user_id parameter is required")
			c.Abort()
			return
		}

		// Parse user_id from URL
		urlUserID, err := strconv.ParseUint(userIDParam, 10, 32)
		if err != nil {
			RespondWithError(c, http.StatusBadRequest, ErrorCodeValidationFailed, "Invalid user_id format")
			c.Abort()
			return
		}

		// Get authenticated user ID from context (set by auth middleware)
		authUserID, exists := c.Get("user_id")
		if !exists {
			RespondWithError(c, http.StatusUnauthorized, ErrorCodeUnauthorized, "User not authenticated")
			c.Abort()
			return
		}

		// Convert auth user ID to uint for comparison
		authUID, ok := authUserID.(uint)
		if !ok {
			RespondWithError(c, http.StatusInternalServerError, ErrorCodeInternalError, "Invalid user authentication data")
			c.Abort()
			return
		}

		// Verify that the URL user_id matches the authenticated user's ID
		if uint(urlUserID) != authUID {
			RespondWithError(c, http.StatusForbidden, ErrorCodeUnauthorized, "Access denied: cannot access other user's resources")
			c.Abort()
			return
		}

		// Store parsed user_id in context for handlers to use
		c.Set("validated_user_id", uint(urlUserID))
		c.Next()
	}
}
