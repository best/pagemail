package avatar

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	_ "golang.org/x/image/webp" // register WebP decoder
)

var (
	ErrFileTooLarge    = errors.New("file too large")
	ErrUnsupportedType = errors.New("unsupported image type")
)

const (
	MaxUploadBytes = 5 << 20 // 5MB
	MaxDimension   = 256
	MaxPixels      = 20 << 20 // 20 million pixels (decompression bomb protection)
	JPEGQuality    = 82
)

var AllowedContentTypes = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

type ProcessResult struct {
	Data        []byte
	ContentType string
	Extension   string
}

func ValidateAndProcess(reader io.Reader, maxBytes int64) (*ProcessResult, error) {
	data, err := io.ReadAll(io.LimitReader(reader, maxBytes+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	if int64(len(data)) > maxBytes {
		return nil, fmt.Errorf("%w: max %d bytes", ErrFileTooLarge, maxBytes)
	}

	contentType := http.DetectContentType(data)
	ext, ok := AllowedContentTypes[contentType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedType, contentType)
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid image: %w", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return nil, fmt.Errorf("invalid image dimensions")
	}
	w, h := int64(cfg.Width), int64(cfg.Height)
	if w > MaxPixels/h || w*h > MaxPixels {
		return nil, fmt.Errorf("image dimensions too large")
	}

	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	resized := imaging.Fit(img, MaxDimension, MaxDimension, imaging.Lanczos)

	var buf bytes.Buffer
	var outputContentType string

	switch ext {
	case "png":
		err = png.Encode(&buf, resized)
		outputContentType = "image/png"
	default:
		err = jpeg.Encode(&buf, resized, &jpeg.Options{Quality: JPEGQuality})
		outputContentType = "image/jpeg"
		ext = "jpg"
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return &ProcessResult{
		Data:        buf.Bytes(),
		ContentType: outputContentType,
		Extension:   ext,
	}, nil
}

func BuildStorageKey(userID uuid.UUID, ext string) string {
	now := time.Now().UTC()
	timestamp := now.Format("20060102150405") + fmt.Sprintf("%06d", now.Nanosecond()/1000)
	return fmt.Sprintf("avatars/%s/%s_%s.%s",
		now.Format("2006/01/02"),
		timestamp,
		userID.String(),
		ext)
}
