package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"pagemail/internal/config"
	"pagemail/internal/models"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DB.Driver {
	case "sqlite":
		dialector = sqlite.Open(cfg.DB.SQLitePath)
	case "postgres":
		dialector = postgres.Open(cfg.DB.URL)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.DB.Driver)
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	if cfg.Server.Env == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info().Str("driver", cfg.DB.Driver).Msg("Database connected")
	return db, nil
}

func Migrate(db *gorm.DB) error {
	log.Info().Msg("Running database migrations")

	err := db.AutoMigrate(
		&models.User{},
		&models.Permission{},
		&models.RolePermission{},
		&models.SMTPProfile{},
		&models.WebhookEndpoint{},
		&models.CaptureTask{},
		&models.CaptureOutput{},
		&models.Delivery{},
		&models.Job{},
		&models.AuditLog{},
	)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err := seedPermissions(db); err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	log.Info().Msg("Database migrations completed")
	return nil
}

func seedPermissions(db *gorm.DB) error {
	permissions := []models.Permission{
		{Name: "admin.users.read", Description: "Read user list"},
		{Name: "admin.users.write", Description: "Create/update users"},
		{Name: "admin.users.delete", Description: "Delete users"},
		{Name: "admin.smtp.global.read", Description: "Read global SMTP settings"},
		{Name: "admin.smtp.global.write", Description: "Update global SMTP settings"},
		{Name: "admin.storage.read", Description: "Read storage settings"},
		{Name: "admin.storage.write", Description: "Update storage settings"},
		{Name: "admin.audit.read", Description: "Read audit logs"},
		{Name: "capture.create", Description: "Create capture tasks"},
		{Name: "capture.read.own", Description: "Read own capture tasks"},
		{Name: "capture.delete.own", Description: "Delete own capture tasks"},
		{Name: "smtp.manage.own", Description: "Manage own SMTP profiles"},
		{Name: "webhook.manage.own", Description: "Manage own webhooks"},
	}

	for _, p := range permissions {
		var existing models.Permission
		if err := db.Where("name = ?", p.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&p).Error; err != nil {
					return err
				}
			}
		}
	}

	adminPerms := []string{
		"admin.users.read", "admin.users.write", "admin.users.delete",
		"admin.smtp.global.read", "admin.smtp.global.write",
		"admin.storage.read", "admin.storage.write",
		"admin.audit.read",
		"capture.create", "capture.read.own", "capture.delete.own",
		"smtp.manage.own", "webhook.manage.own",
	}

	userPerms := []string{
		"capture.create", "capture.read.own", "capture.delete.own",
		"smtp.manage.own", "webhook.manage.own",
	}

	if err := assignRolePermissions(db, "admin", adminPerms); err != nil {
		return err
	}

	if err := assignRolePermissions(db, "user", userPerms); err != nil {
		return err
	}

	return nil
}

func assignRolePermissions(db *gorm.DB, role string, permNames []string) error {
	for _, name := range permNames {
		var perm models.Permission
		if err := db.Where("name = ?", name).First(&perm).Error; err != nil {
			continue
		}

		var existing models.RolePermission
		if err := db.Where("role = ? AND permission_id = ?", role, perm.ID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				rp := models.RolePermission{
					Role:         role,
					PermissionID: perm.ID,
				}
				if err := db.Create(&rp).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}
