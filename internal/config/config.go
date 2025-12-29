package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	DB         DBConfig
	JWT        JWTConfig
	Encryption EncryptionConfig
	Storage    StorageConfig
	SMTP       SMTPConfig
	Capture    CaptureConfig
	Queue      QueueConfig
	RateLimit  RateLimitConfig
	Log        LogConfig
}

type ServerConfig struct {
	Addr string `mapstructure:"SERVER_ADDR" validate:"required"`
	Env  string `mapstructure:"SERVER_ENV" validate:"required,oneof=development staging production"`
}

type DBConfig struct {
	Driver     string `mapstructure:"DB_DRIVER" validate:"required,oneof=sqlite postgres"`
	URL        string `mapstructure:"DB_URL"`
	SQLitePath string `mapstructure:"DB_SQLITE_PATH"`
}

type JWTConfig struct {
	Secret        string        `mapstructure:"JWT_SECRET" validate:"required,min=16"`
	AccessExpiry  time.Duration `mapstructure:"JWT_ACCESS_EXPIRY"`
	RefreshExpiry time.Duration `mapstructure:"JWT_REFRESH_EXPIRY"`
}

type EncryptionConfig struct {
	Key string `mapstructure:"ENCRYPTION_KEY" validate:"required,min=32"`
}

type StorageConfig struct {
	Backend        string `mapstructure:"STORAGE_BACKEND" validate:"required,oneof=local s3"`
	LocalPath      string `mapstructure:"STORAGE_LOCAL_PATH"`
	S3Endpoint     string `mapstructure:"STORAGE_S3_ENDPOINT"`
	S3Region       string `mapstructure:"STORAGE_S3_REGION"`
	S3Bucket       string `mapstructure:"STORAGE_S3_BUCKET"`
	S3AccessKey    string `mapstructure:"STORAGE_S3_ACCESS_KEY"`
	S3SecretKey    string `mapstructure:"STORAGE_S3_SECRET_KEY"`
	S3UsePathStyle bool   `mapstructure:"STORAGE_S3_USE_PATH_STYLE"`
}

type SMTPConfig struct {
	Host      string `mapstructure:"SMTP_HOST"`
	Port      int    `mapstructure:"SMTP_PORT"`
	Username  string `mapstructure:"SMTP_USERNAME"`
	Password  string `mapstructure:"SMTP_PASSWORD"`
	FromName  string `mapstructure:"SMTP_FROM_NAME"`
	FromEmail string `mapstructure:"SMTP_FROM_EMAIL"`
	UseTLS    bool   `mapstructure:"SMTP_USE_TLS"`
}

type CaptureConfig struct {
	ViewportWidth  int `mapstructure:"CAPTURE_VIEWPORT_WIDTH"`
	ViewportHeight int `mapstructure:"CAPTURE_VIEWPORT_HEIGHT"`
	WaitTimeout    int `mapstructure:"CAPTURE_WAIT_TIMEOUT"`
	Workers        int `mapstructure:"CAPTURE_WORKERS" validate:"min=1,max=10"`
}

type QueueConfig struct {
	PollInterval  int `mapstructure:"QUEUE_POLL_INTERVAL" validate:"min=1"`
	MaxRetries    int `mapstructure:"QUEUE_MAX_RETRIES" validate:"min=1"`
	LeaseDuration int `mapstructure:"QUEUE_LEASE_DURATION" validate:"min=60"`
}

type RateLimitConfig struct {
	RPS   float64 `mapstructure:"RATE_LIMIT_RPS" validate:"min=1"`
	Burst int     `mapstructure:"RATE_LIMIT_BURST" validate:"min=1"`
}

type LogConfig struct {
	Level  string `mapstructure:"LOG_LEVEL" validate:"oneof=debug info warn error"`
	Format string `mapstructure:"LOG_FORMAT" validate:"oneof=json console"`
}

func Load() (*Config, error) {
	viper.AutomaticEnv()
	setDefaults()

	// Try to read .env file if it exists
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig() // Ignore error - env vars will be used

	var cfg Config
	cfg.Server.Addr = viper.GetString("SERVER_ADDR")
	cfg.Server.Env = viper.GetString("SERVER_ENV")
	cfg.DB.Driver = viper.GetString("DB_DRIVER")
	cfg.DB.URL = viper.GetString("DB_URL")
	cfg.DB.SQLitePath = viper.GetString("DB_SQLITE_PATH")
	cfg.JWT.Secret = viper.GetString("JWT_SECRET")
	cfg.JWT.AccessExpiry = viper.GetDuration("JWT_ACCESS_EXPIRY")
	cfg.JWT.RefreshExpiry = viper.GetDuration("JWT_REFRESH_EXPIRY")
	cfg.Encryption.Key = viper.GetString("ENCRYPTION_KEY")
	cfg.Storage.Backend = viper.GetString("STORAGE_BACKEND")
	cfg.Storage.LocalPath = viper.GetString("STORAGE_LOCAL_PATH")
	cfg.Storage.S3Endpoint = viper.GetString("STORAGE_S3_ENDPOINT")
	cfg.Storage.S3Region = viper.GetString("STORAGE_S3_REGION")
	cfg.Storage.S3Bucket = viper.GetString("STORAGE_S3_BUCKET")
	cfg.Storage.S3AccessKey = viper.GetString("STORAGE_S3_ACCESS_KEY")
	cfg.Storage.S3SecretKey = viper.GetString("STORAGE_S3_SECRET_KEY")
	cfg.Storage.S3UsePathStyle = viper.GetBool("STORAGE_S3_USE_PATH_STYLE")
	cfg.SMTP.Host = viper.GetString("SMTP_HOST")
	cfg.SMTP.Port = viper.GetInt("SMTP_PORT")
	cfg.SMTP.Username = viper.GetString("SMTP_USERNAME")
	cfg.SMTP.Password = viper.GetString("SMTP_PASSWORD")
	cfg.SMTP.FromName = viper.GetString("SMTP_FROM_NAME")
	cfg.SMTP.FromEmail = viper.GetString("SMTP_FROM_EMAIL")
	cfg.SMTP.UseTLS = viper.GetBool("SMTP_USE_TLS")
	cfg.Capture.ViewportWidth = viper.GetInt("CAPTURE_VIEWPORT_WIDTH")
	cfg.Capture.ViewportHeight = viper.GetInt("CAPTURE_VIEWPORT_HEIGHT")
	cfg.Capture.WaitTimeout = viper.GetInt("CAPTURE_WAIT_TIMEOUT")
	cfg.Capture.Workers = viper.GetInt("CAPTURE_WORKERS")
	cfg.Queue.PollInterval = viper.GetInt("QUEUE_POLL_INTERVAL")
	cfg.Queue.MaxRetries = viper.GetInt("QUEUE_MAX_RETRIES")
	cfg.Queue.LeaseDuration = viper.GetInt("QUEUE_LEASE_DURATION")
	cfg.RateLimit.RPS = viper.GetFloat64("RATE_LIMIT_RPS")
	cfg.RateLimit.Burst = viper.GetInt("RATE_LIMIT_BURST")
	cfg.Log.Level = viper.GetString("LOG_LEVEL")
	cfg.Log.Format = viper.GetString("LOG_FORMAT")

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("SERVER_ADDR", ":8080")
	viper.SetDefault("SERVER_ENV", "development")
	viper.SetDefault("DB_DRIVER", "sqlite")
	viper.SetDefault("DB_SQLITE_PATH", "./pagemail.db")
	viper.SetDefault("JWT_ACCESS_EXPIRY", "1h")
	viper.SetDefault("JWT_REFRESH_EXPIRY", "168h")
	viper.SetDefault("STORAGE_BACKEND", "local")
	viper.SetDefault("STORAGE_LOCAL_PATH", "./storage")
	viper.SetDefault("STORAGE_S3_REGION", "us-east-1")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USE_TLS", true)
	viper.SetDefault("CAPTURE_VIEWPORT_WIDTH", 1920)
	viper.SetDefault("CAPTURE_VIEWPORT_HEIGHT", 1080)
	viper.SetDefault("CAPTURE_WAIT_TIMEOUT", 30000)
	viper.SetDefault("CAPTURE_WORKERS", 3)
	viper.SetDefault("QUEUE_POLL_INTERVAL", 5)
	viper.SetDefault("QUEUE_MAX_RETRIES", 3)
	viper.SetDefault("QUEUE_LEASE_DURATION", 300)
	viper.SetDefault("RATE_LIMIT_RPS", 10)
	viper.SetDefault("RATE_LIMIT_BURST", 20)
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "console")
}

func validateConfig(cfg *Config) error {
	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		var errMsgs []string
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation errors: %s", strings.Join(errMsgs, ", "))
	}

	if cfg.DB.Driver == "sqlite" && cfg.DB.SQLitePath == "" {
		return fmt.Errorf("DB_SQLITE_PATH is required when DB_DRIVER is sqlite")
	}

	if cfg.DB.Driver == "postgres" && cfg.DB.URL == "" {
		return fmt.Errorf("DB_URL is required when DB_DRIVER is postgres")
	}

	if cfg.Storage.Backend == "s3" {
		if cfg.Storage.S3Bucket == "" {
			return fmt.Errorf("STORAGE_S3_BUCKET is required when STORAGE_BACKEND is s3")
		}
	}

	return nil
}
