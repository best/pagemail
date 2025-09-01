package scraper

import (
	"strings"
)

type Scraper interface {
	ScrapeHTML(url string) ([]byte, error)
}

type ScreenshotScraper interface {
	Screenshot(url string) ([]byte, error)
}

type Manager struct {
	httpScraper   *HTTPScraper
	chromeScraper *ChromeScraper
}

func NewManager() *Manager {
	return &Manager{
		httpScraper:   NewHTTPScraper(),
		chromeScraper: NewChromeScraper(),
	}
}

func (m *Manager) ScrapeHTML(url string, useChrome bool) ([]byte, error) {
	if useChrome {
		return m.chromeScraper.ScrapeHTML(url)
	}

	// Try HTTP first, fallback to Chrome if it fails
	content, err := m.httpScraper.ScrapeHTML(url)
	if err != nil {
		// If HTTP fails, try Chrome as fallback
		return m.chromeScraper.ScrapeHTML(url)
	}

	// Check if the content looks like it needs JavaScript rendering
	contentStr := strings.ToLower(string(content))
	if m.needsJavaScript(contentStr) {
		// Use Chrome for JavaScript-heavy sites
		return m.chromeScraper.ScrapeHTML(url)
	}

	return content, nil
}

func (m *Manager) Screenshot(url string) ([]byte, error) {
	return m.chromeScraper.Screenshot(url)
}

func (m *Manager) needsJavaScript(content string) bool {
	// Heuristics to detect if a page needs JavaScript rendering
	jsIndicators := []string{
		"<noscript>",
		"react",
		"angular",
		"vue.js",
		"document.getelementbyid",
		"window.onload",
		"$(document).ready",
		"application/json",
	}

	for _, indicator := range jsIndicators {
		if strings.Contains(content, indicator) {
			return true
		}
	}

	// Check for minimal content (possible SPA)
	if len(strings.TrimSpace(content)) < 1000 {
		return true
	}

	return false
}

func (m *Manager) Close() error {
	if m.httpScraper != nil {
		m.httpScraper.Close()
	}
	return nil
}

// Helper function to determine scraping strategy based on URL
func (m *Manager) GetRecommendedStrategy(url string) (bool, error) {
	// Domain-specific rules for when to use Chrome
	useChrome := false
	
	lowerURL := strings.ToLower(url)
	
	// Sites that definitely need Chrome
	chromeRequired := []string{
		"twitter.com", "x.com",
		"facebook.com",
		"instagram.com",
		"linkedin.com",
		"youtube.com",
		"gmail.com",
		"web.whatsapp.com",
	}
	
	for _, site := range chromeRequired {
		if strings.Contains(lowerURL, site) {
			useChrome = true
			break
		}
	}
	
	return useChrome, nil
}