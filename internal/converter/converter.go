package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Manager struct {
	pdfConverter  *PDFConverter
	htmlProcessor *HTMLProcessor
}

func NewManager() *Manager {
	return &Manager{
		pdfConverter: NewPDFConverter(),
	}
}

func (m *Manager) ConvertContent(content []byte, format string, originalURL string, outputPath string) error {
	// Ensure output directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	switch strings.ToLower(format) {
	case "html":
		return m.convertToHTML(content, originalURL, outputPath)
	case "pdf":
		return m.convertToPDF(content, originalURL, outputPath)
	case "screenshot":
		// Screenshot content is already in the correct format (PNG)
		return m.saveRawContent(content, outputPath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func (m *Manager) convertToHTML(content []byte, originalURL string, outputPath string) error {
	processor := NewHTMLProcessor(originalURL)
	return processor.ProcessHTML(content, outputPath)
}

func (m *Manager) convertToPDF(content []byte, originalURL string, outputPath string) error {
	if m.pdfConverter == nil {
		return fmt.Errorf("PDF converter not initialized")
	}

	// For PDF conversion, we can either:
	// 1. Convert HTML content to PDF
	// 2. Convert URL directly to PDF (usually better results)
	
	// Try URL-based conversion first if we have a valid URL
	if originalURL != "" && (strings.HasPrefix(originalURL, "http://") || strings.HasPrefix(originalURL, "https://")) {
		err := m.pdfConverter.ConvertURLToPDF(originalURL, outputPath)
		if err == nil {
			return nil
		}
		// If URL conversion fails, fall back to HTML content conversion
	}

	// Fall back to HTML content conversion
	return m.pdfConverter.ConvertHTMLToPDF(content, outputPath)
}

func (m *Manager) saveRawContent(content []byte, outputPath string) error {
	return os.WriteFile(outputPath, content, 0644)
}

func (m *Manager) GetOutputExtension(format string) string {
	switch strings.ToLower(format) {
	case "html":
		return ".html"
	case "pdf":
		return ".pdf"
	case "screenshot":
		return ".png"
	default:
		return ".txt"
	}
}

func (m *Manager) GetMimeType(format string) string {
	switch strings.ToLower(format) {
	case "html":
		return "text/html"
	case "pdf":
		return "application/pdf"
	case "screenshot":
		return "image/png"
	default:
		return "text/plain"
	}
}