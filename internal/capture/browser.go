package capture

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/rs/zerolog/log"
)

type BrowserConfig struct {
	Headless       bool
	ViewportWidth  int
	ViewportHeight int
	UserAgent      string
	Timeout        time.Duration
}

type Browser struct {
	browser *rod.Browser
	config  *BrowserConfig
}

func findChromium() string {
	paths := []string{
		"/usr/bin/chromium-browser",
		"/usr/bin/chromium",
		"/usr/bin/google-chrome",
		"/usr/bin/google-chrome-stable",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func NewBrowser(cfg *BrowserConfig) (*Browser, error) {
	chromePath := findChromium()
	if chromePath == "" {
		return nil, fmt.Errorf("chromium not found")
	}

	l := launcher.New().
		Bin(chromePath).
		Headless(cfg.Headless).
		Set("no-sandbox").
		Set("disable-gpu").
		Set("disable-dev-shm-usage").
		Set("disable-setuid-sandbox").
		Set("disable-extensions").
		Set("disable-background-networking").
		Set("disable-sync").
		Set("disable-translate").
		Set("disable-default-apps").
		Set("mute-audio").
		Set("hide-scrollbars").
		Set("disable-features", "VizDisplayCompositor")

	browserURL, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(browserURL)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

	return &Browser{
		browser: browser,
		config:  cfg,
	}, nil
}

func (b *Browser) Close() error {
	return b.browser.Close()
}

type CaptureOptions struct {
	URL            string
	Cookies        []*proto.NetworkCookieParam
	ViewportWidth  int
	ViewportHeight int
	UserAgent      string
	Timeout        time.Duration
	WaitUntil      string
}

type CaptureResult struct {
	HTML       []byte
	PDF        []byte
	Screenshot []byte
	Title      string
	FinalURL   string
}

func (b *Browser) Capture(ctx context.Context, opts *CaptureOptions) (*CaptureResult, error) {
	if err := validateURL(opts.URL); err != nil {
		return nil, err
	}

	page, err := b.browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}
	defer page.Close()

	width := opts.ViewportWidth
	if width == 0 {
		width = b.config.ViewportWidth
	}
	height := opts.ViewportHeight
	if height == 0 {
		height = b.config.ViewportHeight
	}

	if err := page.SetViewport(&proto.EmulationSetDeviceMetricsOverride{
		Width:  width,
		Height: height,
	}); err != nil {
		return nil, fmt.Errorf("failed to set viewport: %w", err)
	}

	if opts.UserAgent != "" {
		if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{
			UserAgent: opts.UserAgent,
		}); err != nil {
			return nil, fmt.Errorf("failed to set user agent: %w", err)
		}
	}

	if len(opts.Cookies) > 0 {
		if err := page.SetCookies(opts.Cookies); err != nil {
			return nil, fmt.Errorf("failed to set cookies: %w", err)
		}
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = b.config.Timeout
	}
	page = page.Timeout(timeout)

	if err := page.Navigate(opts.URL); err != nil {
		return nil, fmt.Errorf("failed to navigate: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("failed to wait for page load: %w", err)
	}

	time.Sleep(2 * time.Second)

	result := &CaptureResult{}

	info, err := page.Info()
	if err == nil {
		result.Title = info.Title
		result.FinalURL = info.URL
	}

	html, err := page.HTML()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get HTML")
	} else {
		result.HTML = []byte(html)
	}

	paperWidth := 8.5
	paperHeight := 11.0
	pdf, err := page.PDF(&proto.PagePrintToPDF{
		PrintBackground: true,
		PaperWidth:      &paperWidth,
		PaperHeight:     &paperHeight,
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to generate PDF")
	} else {
		result.PDF, _ = io.ReadAll(pdf)
	}

	quality := 90
	screenshot, err := page.Screenshot(true, &proto.PageCaptureScreenshot{
		Format:  proto.PageCaptureScreenshotFormatPng,
		Quality: &quality,
	})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to take screenshot")
	} else {
		result.Screenshot = screenshot
	}

	return result, nil
}

func validateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("only http and https URLs are allowed")
	}

	host := parsed.Hostname()

	if isPrivateIP(host) {
		return fmt.Errorf("private IP addresses are not allowed")
	}

	blockedHosts := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"::1",
		"metadata.google.internal",
		"169.254.169.254",
	}

	for _, blocked := range blockedHosts {
		if strings.EqualFold(host, blocked) {
			return fmt.Errorf("blocked host: %s", host)
		}
	}

	return nil
}

func isPrivateIP(host string) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			return false
		}
		ip = ips[0]
	}

	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"169.254.0.0/16",
		"fc00::/7",
		"fe80::/10",
		"::1/128",
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

func ParseCookies(cookieStr, domain string) []*proto.NetworkCookieParam {
	var cookies []*proto.NetworkCookieParam

	pairs := strings.Split(cookieStr, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		cookies = append(cookies, &proto.NetworkCookieParam{
			Name:   name,
			Value:  value,
			Domain: domain,
		})
	}

	return cookies
}
