package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"pagemail/internal/config"
	"pagemail/internal/db"
	"pagemail/internal/queue"
	"pagemail/internal/routes"
	"pagemail/internal/storage"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	setupLogger(cfg)

	log.Info().
		Str("version", Version).
		Str("build_time", BuildTime).
		Str("env", cfg.Server.Env).
		Msg("Starting Pagemail")

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	if err := db.Migrate(database); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	store, err := storage.New(&storage.Config{
		Backend:        cfg.Storage.Backend,
		LocalPath:      cfg.Storage.LocalPath,
		S3Endpoint:     cfg.Storage.S3Endpoint,
		S3Region:       cfg.Storage.S3Region,
		S3Bucket:       cfg.Storage.S3Bucket,
		S3AccessKey:    cfg.Storage.S3AccessKey,
		S3SecretKey:    cfg.Storage.S3SecretKey,
		S3UsePathStyle: cfg.Storage.S3UsePathStyle,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize storage")
	}

	router := routes.Setup(cfg, database, store)

	dispatcher := queue.NewDispatcher(cfg, database, store)
	go dispatcher.Start()

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("addr", cfg.Server.Addr).Msg("HTTP server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	dispatcher.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}

func setupLogger(cfg *config.Config) {
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.Log.Format == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}

	zerolog.TimeFieldFormat = time.RFC3339
}
