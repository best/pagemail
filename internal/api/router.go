package api

import (
	"net/http"
	"os"
	"pagemail/internal/auth"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS middleware - supports both development and production
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	// Configure CORS origins based on environment
	if gin.Mode() == gin.DebugMode {
		// Development mode - allow frontend dev server
		corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	} else {
		// Production mode - allow same origin and common deployment URLs
		corsConfig.AllowAllOrigins = true
	}

	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", handleHealthCheck)

	// Static file serving for production
	setupStaticFileServing(router)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/register", handleRegister)
			authRoutes.POST("/login", handleLogin)
			authRoutes.POST("/verification", handleVerifyEmail)
			authRoutes.POST("/verification/resend", handleResendVerification)
		}

		// Protected user routes (with user ID validation)
		userRoutes := v1.Group("/users")
		userRoutes.Use(auth.AuthMiddleware())
		userRoutes.Use(ValidateUserIDMiddleware())
		{
			userRoutes.GET("/:user_id", handleProfile)
			userRoutes.GET("/:user_id/scrapes", handleRequestHistory)
			userRoutes.GET("/:user_id/usage", handleUsageInfo)
		}

		// Public scrape routes (with optional auth + rate limiting)
		scrapeRoutes := v1.Group("/scrapes")
		scrapeRoutes.Use(auth.OptionalAuthMiddleware())
		scrapeRoutes.Use(auth.RateLimitMiddleware())
		{
			scrapeRoutes.POST("", handleScrapeRequest)
			scrapeRoutes.GET("/:scrape_id", handleScrapeDetail)
		}
	}

	return router
}

// setupStaticFileServing configures static file serving for the frontend
func setupStaticFileServing(router *gin.Engine) {
	// Check if we're in production mode (static files should exist)
	webDir := "./frontend"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		// No static files available, skip static serving
		return
	}

	// Serve static assets (CSS, JS, images, etc.)
	router.Static("/_next", filepath.Join(webDir, "_next"))
	router.StaticFile("/favicon.ico", filepath.Join(webDir, "favicon.ico"))

	// Serve other static files from public directory
	if publicDir := filepath.Join(webDir, "public"); dirExists(publicDir) {
		router.Static("/public", publicDir)
	}

	// Handle SPA routing - serve index.html for all non-API routes
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Don't handle API routes
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/health") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// Serve static files directly if they exist
		if strings.Contains(path, ".") {
			filePath := filepath.Join(webDir, path)
			if fileExists(filePath) {
				c.File(filePath)
				return
			}
		}

		// For all other routes, serve the SPA index.html
		indexPath := filepath.Join(webDir, "index.html")
		if fileExists(indexPath) {
			c.File(indexPath)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Frontend not available"})
		}
	})
}

// Helper functions
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
