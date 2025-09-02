package api

import (
	"pagemail/internal/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	
	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	// Health check endpoint
	router.GET("/health", handleHealthCheck)
	
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", handleRegister)
			authRoutes.POST("/login", handleLogin)
			authRoutes.GET("/verify/:token", handleVerifyEmail)
			authRoutes.POST("/resend-verification", handleResendVerification)
		}
		
		// Protected user routes
		userRoutes := v1.Group("/user")
		userRoutes.Use(auth.AuthMiddleware())
		{
			userRoutes.GET("/profile", handleProfile)
		}
		
		// Page scraping routes (with optional auth + rate limiting)
		pages := v1.Group("/pages")
		pages.Use(auth.OptionalAuthMiddleware())
		pages.Use(auth.RateLimitMiddleware())
		{
			pages.POST("/scrape", handleScrapeRequest)
		}
		
		// Usage info routes (with optional auth)
		usage := v1.Group("/usage")
		usage.Use(auth.OptionalAuthMiddleware())
		{
			usage.GET("/", handleUsageInfo)
		}
		
		// Protected page history routes
		protectedPages := v1.Group("/pages")
		protectedPages.Use(auth.AuthMiddleware())
		{
			protectedPages.GET("/history", handleRequestHistory)
		}
	}
	
	return router
}

