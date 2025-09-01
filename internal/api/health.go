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

	// Database check
	if database.DB != nil {
		sqlDB, err := database.DB.DB()
		if err != nil {
			checks["database"] = "error: " + err.Error()
			overallStatus = "unhealthy"
		} else if err := sqlDB.Ping(); err != nil {
			checks["database"] = "connection_failed: " + err.Error()
			overallStatus = "unhealthy"
		} else {
			checks["database"] = "connected"
		}
	} else {
		checks["database"] = "not_initialized"
		overallStatus = "unhealthy"
	}

	// SMTP check (optional - only if credentials are configured)
	mailerService := mailer.NewMailer()
	if err := mailerService.TestConnection(); err != nil {
		checks["smtp"] = "unavailable: " + err.Error()
		// Don't mark as unhealthy since SMTP might not be configured in dev
	} else {
		checks["smtp"] = "configured"
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