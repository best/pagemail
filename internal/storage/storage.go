package storage

import (
	"context"
	"io"
)

type FileInfo struct {
	Key         string
	Size        int64
	ContentType string
}

type Storage interface {
	Upload(ctx context.Context, key string, reader io.Reader, contentType string) (*FileInfo, error)
	Download(ctx context.Context, key string) (io.ReadCloser, *FileInfo, error)
	Delete(ctx context.Context, key string) error
	GetPresignedURL(ctx context.Context, key string, expirySeconds int) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type Config struct {
	Backend        string
	LocalPath      string
	S3Endpoint     string
	S3Region       string
	S3Bucket       string
	S3AccessKey    string
	S3SecretKey    string
	S3UsePathStyle bool
}

func New(cfg *Config) (Storage, error) {
	switch cfg.Backend {
	case "s3":
		return NewS3Storage(cfg)
	default:
		return NewLocalStorage(cfg.LocalPath)
	}
}
