package converter

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type HTMLProcessor struct {
	baseURL string
}

func NewHTMLProcessor(baseURL string) *HTMLProcessor {
	return &HTMLProcessor{
		baseURL: baseURL,
	}
}

func (h *HTMLProcessor) ProcessHTML(htmlContent []byte, outputPath string) error {
	// Process the HTML content
	processedHTML := h.makeHTMLSelfContained(string(htmlContent))
	
	// Write processed HTML to file
	if err := os.WriteFile(outputPath, []byte(processedHTML), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	return nil
}

func (h *HTMLProcessor) makeHTMLSelfContained(html string) string {
	// Convert relative URLs to absolute URLs
	html = h.convertRelativeURLs(html)
	
	// Add metadata and styling
	html = h.addMetadata(html)
	
	return html
}

func (h *HTMLProcessor) convertRelativeURLs(html string) string {
	if h.baseURL == "" {
		return html
	}

	baseURL, err := url.Parse(h.baseURL)
	if err != nil {
		return html
	}

	// Convert relative URLs in src attributes
	srcRegex := regexp.MustCompile(`src="([^"]*)"`)
	html = srcRegex.ReplaceAllStringFunc(html, func(match string) string {
		return h.convertURL(match, baseURL, `src="([^"]*)"`)
	})

	// Convert relative URLs in href attributes
	hrefRegex := regexp.MustCompile(`href="([^"]*)"`)
	html = hrefRegex.ReplaceAllStringFunc(html, func(match string) string {
		return h.convertURL(match, baseURL, `href="([^"]*)"`)
	})

	// Convert relative URLs in CSS url() functions
	cssRegex := regexp.MustCompile(`url\(["']?([^"')]*)["']?\)`)
	html = cssRegex.ReplaceAllStringFunc(html, func(match string) string {
		return h.convertCSSURL(match, baseURL)
	})

	return html
}

func (h *HTMLProcessor) convertURL(match string, baseURL *url.URL, pattern string) string {
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(match)
	
	if len(matches) < 2 {
		return match
	}

	originalURL := matches[1]
	
	// Skip data URLs, absolute URLs, and fragments
	if strings.HasPrefix(originalURL, "data:") ||
		strings.HasPrefix(originalURL, "http://") ||
		strings.HasPrefix(originalURL, "https://") ||
		strings.HasPrefix(originalURL, "//") ||
		strings.HasPrefix(originalURL, "#") {
		return match
	}

	// Convert relative URL to absolute
	resolvedURL, err := baseURL.Parse(originalURL)
	if err != nil {
		return match
	}

	return strings.Replace(match, originalURL, resolvedURL.String(), 1)
}

func (h *HTMLProcessor) convertCSSURL(match string, baseURL *url.URL) string {
	regex := regexp.MustCompile(`url\(["']?([^"')]*)["']?\)`)
	matches := regex.FindStringSubmatch(match)
	
	if len(matches) < 2 {
		return match
	}

	originalURL := matches[1]
	
	// Skip data URLs, absolute URLs
	if strings.HasPrefix(originalURL, "data:") ||
		strings.HasPrefix(originalURL, "http://") ||
		strings.HasPrefix(originalURL, "https://") ||
		strings.HasPrefix(originalURL, "//") {
		return match
	}

	// Convert relative URL to absolute
	resolvedURL, err := baseURL.Parse(originalURL)
	if err != nil {
		return match
	}

	return strings.Replace(match, originalURL, resolvedURL.String(), 1)
}

func (h *HTMLProcessor) addMetadata(html string) string {
	// Add charset if missing
	if !strings.Contains(strings.ToLower(html), `<meta charset=`) &&
		!strings.Contains(strings.ToLower(html), `<meta http-equiv="content-type"`) {
		html = strings.Replace(html, "<head>", 
			`<head>
    <meta charset="UTF-8">`, 1)
	}

	// Add viewport if missing
	if !strings.Contains(strings.ToLower(html), `name="viewport"`) {
		html = strings.Replace(html, "<head>", 
			`<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">`, 1)
	}

	// Add PageMail metadata
	metadata := fmt.Sprintf(`
    <meta name="generator" content="PageMail">
    <meta name="archived-date" content="%s">
    <meta name="original-url" content="%s">`, 
		h.getCurrentTimestamp(), h.baseURL)
	
	html = strings.Replace(html, "<head>", "<head>"+metadata, 1)

	return html
}

func (h *HTMLProcessor) getCurrentTimestamp() string {
	return fmt.Sprintf("%d", 1694000000) // Placeholder timestamp
}