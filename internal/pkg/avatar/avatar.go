package avatar

import (
	"bytes"
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
		return nil, fmt.Errorf("file too large, max %d bytes", maxBytes)
	}

	contentType := http.DetectContentType(data)
	ext, ok := AllowedContentTypes[contentType]
	if !ok {
		return nil, fmt.Errorf("unsupported image type: %s", contentType)
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid image: %w", err)
	}
	if cfg.Width*cfg.Height > MaxPixels {
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
