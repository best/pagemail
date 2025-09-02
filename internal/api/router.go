package api

import (
	"net/http"
	"os"
	"path/filepath"
	"pagemail/internal/auth"
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

// setupStaticFileServing configures static file serving for the frontend
func setupStaticFileServing(router *gin.Engine) {
	// Check if we're in production mode (static files should exist)
	webDir := "./web"
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

