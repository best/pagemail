package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"pagemail/internal/database"

	"github.com/joho/godotenv"
)

func main() {
	var (
		action = flag.String("action", "up", "Migration action: up, down, status")
		steps  = flag.Int("steps", 1, "Number of steps to rollback (for down action)")
	)
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	switch *action {
	case "up":
		if err := database.MigrateDatabase(); err != nil {
			log.Fatal("Migration failed:", err)
		}
		fmt.Println("✅ Migrations completed successfully")

	case "down":
		if err := database.RollbackMigrations(*steps); err != nil {
			log.Fatal("Rollback failed:", err)
		}
		fmt.Printf("✅ Rolled back %d migration(s)\n", *steps)

	case "status":
		version, dirty, err := database.GetMigrationStatus()
		if err != nil {
			log.Fatal("Failed to get migration status:", err)
		}
		
		status := "clean"
		if dirty {
			status = "dirty"
		}
		
		fmt.Printf("Current migration version: %d (%s)\n", version, status)

	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: up, down, status")
		os.Exit(1)
	}
}