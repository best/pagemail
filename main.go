package main

import (
	"log"
	"pagemail/internal/api"
	"pagemail/internal/database"
)

func main() {
	log.Println("Starting PageMail server...")
	
	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	router := api.SetupRouter()
	
	port := ":8080"
	log.Printf("Server listening on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}