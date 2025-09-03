package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorType string `json:"error_type"`
	Message   string `json:"message"`
	Details   string `json:"details,omitempty"`
}

// Error codes by category
const (
	// Authentication & Authorization Errors (1000-1999)
	ErrorCodeUnauthorized           = 1001
	ErrorCodeInvalidCredentials     = 1002
	ErrorCodeTokenExpired           = 1003
	ErrorCodeEmailNotVerified       = 1004
	ErrorCodeAccountDeactivated     = 1005
	ErrorCodeInvalidVerificationToken = 1006

	// User Related Errors (2000-2999)
	ErrorCodeUserNotFound       = 2001
	ErrorCodeUserAlreadyExists  = 2002
	ErrorCodeUserCreationFailed = 2003

	// Request & Validation Errors (3000-3999)
	ErrorCodeInvalidRequest     = 3001
	ErrorCodeValidationFailed   = 3002
	ErrorCodeInvalidURL         = 3003
	ErrorCodeInvalidEmail       = 3004
	ErrorCodeInvalidFormat      = 3005
	ErrorCodeMissingParameter   = 3006

	// Quota & Rate Limiting Errors (4000-4999)
	ErrorCodeDailyLimitExceeded   = 4001
	ErrorCodeMonthlyLimitExceeded = 4002
	ErrorCodeRateLimitExceeded    = 4003
	ErrorCodeTooManyRequests      = 4004

	// System & Service Errors (5000-5999)
	ErrorCodeInternalError        = 5001
	ErrorCodeDatabaseError        = 5002
	ErrorCodeEmailServiceError    = 5003
	ErrorCodeScrapingFailed       = 5004
	ErrorCodeFileProcessingFailed = 5005
	ErrorCodeServiceUnavailable   = 5006
)

// Error type strings
const (
	ErrorTypeAuthentication = "AUTHENTICATION_ERROR"
	ErrorTypeUser           = "USER_ERROR"
	ErrorTypeValidation     = "VALIDATION_ERROR"
	ErrorTypeQuota          = "QUOTA_ERROR"
	ErrorTypeSystem         = "SYSTEM_ERROR"
)

// Standard error responses
var errorMessages = map[int]ErrorResponse{
	ErrorCodeUnauthorized: {
		ErrorCode: ErrorCodeUnauthorized,
		ErrorType: ErrorTypeAuthentication,
		Message:   "Authentication required",
	},
	ErrorCodeInvalidCredentials: {
		ErrorCode: ErrorCodeInvalidCredentials,
		ErrorType: ErrorTypeAuthentication,
		Message:   "Invalid email or password",
	},
	ErrorCodeTokenExpired: {
		ErrorCode: ErrorCodeTokenExpired,
		ErrorType: ErrorTypeAuthentication,
		Message:   "Authentication token has expired",
	},
	ErrorCodeEmailNotVerified: {
		ErrorCode: ErrorCodeEmailNotVerified,
		ErrorType: ErrorTypeAuthentication,
		Message:   "Email address not verified",
	},
	ErrorCodeAccountDeactivated: {
		ErrorCode: ErrorCodeAccountDeactivated,
		ErrorType: ErrorTypeAuthentication,
		Message:   "Account has been deactivated",
	},
	ErrorCodeInvalidVerificationToken: {
		ErrorCode: ErrorCodeInvalidVerificationToken,
		ErrorType: ErrorTypeAuthentication,
		Message:   "Invalid or expired verification token",
	},
	ErrorCodeUserNotFound: {
		ErrorCode: ErrorCodeUserNotFound,
		ErrorType: ErrorTypeUser,
		Message:   "User not found",
	},
	ErrorCodeUserAlreadyExists: {
		ErrorCode: ErrorCodeUserAlreadyExists,
		ErrorType: ErrorTypeUser,
		Message:   "User already exists with this email",
	},
	ErrorCodeUserCreationFailed: {
		ErrorCode: ErrorCodeUserCreationFailed,
		ErrorType: ErrorTypeUser,
		Message:   "Failed to create user account",
	},
	ErrorCodeInvalidRequest: {
		ErrorCode: ErrorCodeInvalidRequest,
		ErrorType: ErrorTypeValidation,
		Message:   "Invalid request format",
	},
	ErrorCodeValidationFailed: {
		ErrorCode: ErrorCodeValidationFailed,
		ErrorType: ErrorTypeValidation,
		Message:   "Request validation failed",
	},
	ErrorCodeInvalidURL: {
		ErrorCode: ErrorCodeInvalidURL,
		ErrorType: ErrorTypeValidation,
		Message:   "Invalid URL format",
	},
	ErrorCodeInvalidEmail: {
		ErrorCode: ErrorCodeInvalidEmail,
		ErrorType: ErrorTypeValidation,
		Message:   "Invalid email format",
	},
	ErrorCodeInvalidFormat: {
		ErrorCode: ErrorCodeInvalidFormat,
		ErrorType: ErrorTypeValidation,
		Message:   "Invalid format specified",
	},
	ErrorCodeMissingParameter: {
		ErrorCode: ErrorCodeMissingParameter,
		ErrorType: ErrorTypeValidation,
		Message:   "Required parameter is missing",
	},
	ErrorCodeDailyLimitExceeded: {
		ErrorCode: ErrorCodeDailyLimitExceeded,
		ErrorType: ErrorTypeQuota,
		Message:   "Daily request limit exceeded",
	},
	ErrorCodeMonthlyLimitExceeded: {
		ErrorCode: ErrorCodeMonthlyLimitExceeded,
		ErrorType: ErrorTypeQuota,
		Message:   "Monthly request limit exceeded",
	},
	ErrorCodeRateLimitExceeded: {
		ErrorCode: ErrorCodeRateLimitExceeded,
		ErrorType: ErrorTypeQuota,
		Message:   "Rate limit exceeded",
	},
	ErrorCodeTooManyRequests: {
		ErrorCode: ErrorCodeTooManyRequests,
		ErrorType: ErrorTypeQuota,
		Message:   "Too many requests",
	},
	ErrorCodeInternalError: {
		ErrorCode: ErrorCodeInternalError,
		ErrorType: ErrorTypeSystem,
		Message:   "Internal server error",
	},
	ErrorCodeDatabaseError: {
		ErrorCode: ErrorCodeDatabaseError,
		ErrorType: ErrorTypeSystem,
		Message:   "Database operation failed",
	},
	ErrorCodeEmailServiceError: {
		ErrorCode: ErrorCodeEmailServiceError,
		ErrorType: ErrorTypeSystem,
		Message:   "Email service unavailable",
	},
	ErrorCodeScrapingFailed: {
		ErrorCode: ErrorCodeScrapingFailed,
		ErrorType: ErrorTypeSystem,
		Message:   "Failed to scrape web page",
	},
	ErrorCodeFileProcessingFailed: {
		ErrorCode: ErrorCodeFileProcessingFailed,
		ErrorType: ErrorTypeSystem,
		Message:   "File processing failed",
	},
	ErrorCodeServiceUnavailable: {
		ErrorCode: ErrorCodeServiceUnavailable,
		ErrorType: ErrorTypeSystem,
		Message:   "Service temporarily unavailable",
	},
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, httpStatus int, errorCode int, details ...string) {
	errorResp, exists := errorMessages[errorCode]
	if !exists {
		errorResp = ErrorResponse{
			ErrorCode: ErrorCodeInternalError,
			ErrorType: ErrorTypeSystem,
			Message:   "Unknown error",
		}
	}

	if len(details) > 0 {
		errorResp.Details = details[0]
	}

	c.JSON(httpStatus, errorResp)
}

// RespondWithValidationError sends a validation error with field details
func RespondWithValidationError(c *gin.Context, fieldError string) {
	RespondWithError(c, http.StatusBadRequest, ErrorCodeValidationFailed, fieldError)
}

// RespondWithQuotaError sends a quota exceeded error with quota information
func RespondWithQuotaError(c *gin.Context, errorCode int, used, limit int, resetTime ...int64) {
	response := map[string]interface{}{
		"error_code": errorCode,
		"error_type": ErrorTypeQuota,
		"used":       used,
		"limit":      limit,
	}

	if errorCode == ErrorCodeDailyLimitExceeded {
		response["message"] = "Daily request limit exceeded"
	} else if errorCode == ErrorCodeMonthlyLimitExceeded {
		response["message"] = "Monthly request limit exceeded"
	} else {
		response["message"] = "Request limit exceeded"
	}

	if len(resetTime) > 0 {
		response["reset_time"] = resetTime[0]
	}

	c.JSON(http.StatusTooManyRequests, response)
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// RespondWithCreated sends a standardized created response
func RespondWithCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

// RespondWithAccepted sends a standardized accepted response
func RespondWithAccepted(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAccepted, data)
}