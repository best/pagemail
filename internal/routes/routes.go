package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"pagemail/internal/audit"
	"pagemail/internal/config"
	"pagemail/internal/handlers"
	"pagemail/internal/middleware"
	"pagemail/internal/storage"
)

func Setup(cfg *config.Config, db *gorm.DB, store storage.Storage) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg))
	r.Use(middleware.TraceID())

	auditLogger := audit.NewLogger(db)
	h := handlers.New(cfg, db, store, auditLogger)

	v1 := r.Group("/v1")
	v1.GET("/health", h.Health)
	v1.GET("/ready", h.Ready)
	v1.GET("/config/site", h.GetPublicSiteConfig)

	auth := v1.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/logout", middleware.Auth(cfg), h.Logout)

	users := v1.Group("/users")
	users.Use(middleware.Auth(cfg))
	users.GET("/me", h.GetCurrentUser)
	users.PATCH("/me", h.UpdateCurrentUser)
	users.PATCH("/me/password", h.ChangePassword)
	users.PUT("/me/avatar", h.UploadAvatar)
	users.DELETE("/me/avatar", h.DeleteAvatar)
	users.GET("/:id/avatar", h.GetAvatar)

	captures := v1.Group("/captures")
	captures.Use(middleware.Auth(cfg))
	captures.POST("", h.CreateCapture)
	captures.GET("", h.ListCaptures)
	captures.GET("/:id", h.GetCapture)
	captures.POST("/:id/retry", h.RetryCapture)
	captures.DELETE("/:id", h.DeleteCapture)
	captures.GET("/:id/outputs", h.ListCaptureOutputs)
	captures.GET("/:id/outputs/:oid/download", h.DownloadOutput)
	captures.GET("/:id/outputs/:oid/preview", h.PreviewOutput)
	captures.POST("/:id/deliver", h.DeliverCapture)

	deliveries := v1.Group("/deliveries")
	deliveries.Use(middleware.Auth(cfg))
	deliveries.GET("", h.ListDeliveries)

	smtp := v1.Group("/smtp")
	smtp.Use(middleware.Auth(cfg))
	smtp.GET("/profiles", h.ListSMTPProfiles)
	smtp.POST("/profiles", h.CreateSMTPProfile)
	smtp.PUT("/profiles/:id", h.UpdateSMTPProfile)
	smtp.DELETE("/profiles/:id", h.DeleteSMTPProfile)
	smtp.POST("/profiles/:id/test", h.TestSMTPProfile)

	webhooks := v1.Group("/webhooks")
	webhooks.Use(middleware.Auth(cfg))
	webhooks.GET("", h.ListWebhooks)
	webhooks.POST("", h.CreateWebhook)
	webhooks.PUT("/:id", h.UpdateWebhook)
	webhooks.DELETE("/:id", h.DeleteWebhook)
	webhooks.POST("/:id/test", h.TestWebhook)

	admin := v1.Group("/admin")
	admin.Use(middleware.Auth(cfg), middleware.RequireAdmin())
	admin.GET("/users", h.AdminListUsers)
	admin.GET("/users/:id", h.AdminGetUser)
	admin.PATCH("/users/:id", h.AdminUpdateUser)
	admin.DELETE("/users/:id", h.AdminDeleteUser)
	admin.GET("/config/site", h.GetSiteConfig)
	admin.PUT("/config/site", h.UpdateSiteConfig)
	admin.GET("/smtp/global", h.GetGlobalSMTP)
	admin.PUT("/smtp/global", h.UpdateGlobalSMTP)
	admin.GET("/storage", h.GetStorageConfig)
	admin.PUT("/storage", h.UpdateStorageConfig)
	admin.GET("/audit-logs", h.ListAuditLogs)
	admin.GET("/stats", h.GetSystemStats)

	return r
}
