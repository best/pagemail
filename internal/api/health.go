package api

import (
	"net/http"
	"pagemail/internal/database"
	"pagemail/internal/mailer"

	"github.com/gin-gonic/gin"
)

type HealthStatus struct {
	Status   string            `json:"status"`
	Service  string            `json:"service"`
	Checks   map[string]string `json:"checks"`
	Version  string            `json:"version"`
}

func handleHealthCheck(c *gin.Context) {
	checks := make(map[string]string)
	overallStatus := "healthy"

	// Database check - use a simple query instead of Ping
	if database.DB != nil {
		var result int
		err := database.DB.Raw("SELECT 1").Scan(&result).Error
		if err != nil {
			checks["database"] = "query_failed: " + err.Error()
			overallStatus = "unhealthy"
		} else {
			checks["database"] = "connected"
		}
	} else {
		checks["database"] = "not_initialized"
		overallStatus = "unhealthy"
	}

	// SMTP check - only test if credentials are configured
	mailerService := mailer.NewMailer()
	if mailerService.IsConfigured() {
		if err := mailerService.TestConnection(); err != nil {
			checks["smtp"] = "unavailable: " + err.Error()
			overallStatus = "unhealthy"
		} else {
			checks["smtp"] = "connected"
		}
	} else {
		checks["smtp"] = "not_configured"
	}

	response := HealthStatus{
		Status:  overallStatus,
		Service: "pagemail",
		Checks:  checks,
		Version: "1.0.0",
	}

	statusCode := http.StatusOK
	if overallStatus == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}