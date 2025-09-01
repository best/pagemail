package scraper

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type ChromeScraper struct {
	timeout time.Duration
	options []chromedp.ExecAllocatorOption
}

func NewChromeScraper() *ChromeScraper {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-background-timer-throttling", false),
		chromedp.Flag("disable-backgrounding-occluded-windows", false),
		chromedp.Flag("disable-renderer-backgrounding", false),
		chromedp.WindowSize(1920, 1080),
	)

	return &ChromeScraper{
		timeout: 60 * time.Second,
		options: opts,
	}
}

func (c *ChromeScraper) ScrapeHTML(targetURL string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Create allocator context
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, c.options...)
	defer cancelAlloc()

	// Create chrome context
	chromeCtx, cancelChrome := chromedp.NewContext(allocCtx)
	defer cancelChrome()

	var htmlContent string

	// Run chrome actions
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second), // Wait for dynamic content
		chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
	)

	if err != nil {
		return nil, fmt.Errorf("chrome scraping failed: %w", err)
	}

	return []byte(htmlContent), nil
}

func (c *ChromeScraper) Screenshot(targetURL string) ([]byte, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Create allocator context
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, c.options...)
	defer cancelAlloc()

	// Create chrome context
	chromeCtx, cancelChrome := chromedp.NewContext(allocCtx)
	defer cancelChrome()

	var screenshotData []byte

	// Run chrome actions for screenshot
	err := chromedp.Run(chromeCtx,
		chromedp.Navigate(targetURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(3*time.Second), // Wait for dynamic content and images
		chromedp.FullScreenshot(&screenshotData, 90), // Quality 90
	)

	if err != nil {
		return nil, fmt.Errorf("chrome screenshot failed: %w", err)
	}

	return screenshotData, nil
}