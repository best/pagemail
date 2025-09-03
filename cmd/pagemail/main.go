package main

import (
	"log"
	"os"
	"pagemail/internal/api"
	"pagemail/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting PageMail server...")

	// Load environment variables from .env file (if exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	router := api.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	log.Printf("Server listening on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
