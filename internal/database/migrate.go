package database

import (
	"embed"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// MigrateDatabase runs database migrations
func MigrateDatabase() error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Get underlying sql.DB from GORM
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Create postgres driver instance with NoLock config to avoid closing connection
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{
		MigrationsTable: "schema_migrations",
	})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create source driver from embedded filesystem
	sourceDriver, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	// Note: Don't call m.Close() to avoid closing the shared database connection

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Database migrations: No changes to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back the last migration
func RollbackMigrations(steps int) error {
	if DB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	// Note: Don't call m.Close() to avoid closing the shared database connection

	// Get current version
	version, dirty, err := m.Version()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if dirty {
		return fmt.Errorf("database is in dirty state, please fix manually")
	}

	// Calculate target version
	targetVersion := int(version) - steps
	if targetVersion < 0 {
		targetVersion = 0
	}

	// Migrate to target version
	if err := m.Migrate(uint(targetVersion)); err != nil {
		if err == migrate.ErrNoChange {
			log.Printf("Database rollback: No changes to apply (already at version %d)", targetVersion)
			return nil
		}
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	log.Printf("Database rolled back to version %d", targetVersion)
	return nil
}

// GetMigrationStatus returns current migration status
func GetMigrationStatus() (version uint, dirty bool, err error) {
	if DB == nil {
		return 0, false, fmt.Errorf("database connection not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return 0, false, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return 0, false, fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}
	// Note: Don't call m.Close() to avoid closing the shared database connection

	return m.Version()
}
