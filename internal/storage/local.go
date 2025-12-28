package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) (*FileInfo, error) {
	fullPath := filepath.Join(s.basePath, key)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	size, err := io.Copy(file, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &FileInfo{
		Key:         key,
		Size:        size,
		ContentType: contentType,
	}, nil
}

func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, *FileInfo, error) {
	fullPath := filepath.Join(s.basePath, key)

	stat, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("file not found: %s", key)
		}
		return nil, nil, fmt.Errorf("failed to stat file: %w", err)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	info := &FileInfo{
		Key:         key,
		Size:        stat.Size(),
		ContentType: "application/octet-stream",
	}

	return file, info, nil
}

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(s.basePath, key)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *LocalStorage) GetPresignedURL(ctx context.Context, key string, expirySeconds int) (string, error) {
	return "", fmt.Errorf("presigned URLs not supported for local storage")
}

func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	fullPath := filepath.Join(s.basePath, key)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
