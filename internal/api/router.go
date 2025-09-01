package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"service": "pagemail",
		})
	})
	
	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handleRegister)
			auth.POST("/login", handleLogin)
		}
		
		// Page scraping routes
		pages := v1.Group("/pages")
		{
			pages.POST("/scrape", handleScrape)
			pages.GET("/history", handleHistory)
		}
	}
	
	return router
}

func handleRegister(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Register endpoint - to be implemented"})
}

func handleLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Login endpoint - to be implemented"})
}

func handleScrape(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Scrape endpoint - to be implemented"})
}

func handleHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "History endpoint - to be implemented"})
}