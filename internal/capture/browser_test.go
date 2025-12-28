package capture

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		// Valid URLs
		{"valid https", "https://example.com", false},
		{"valid http", "http://example.com", false},
		{"valid with path", "https://example.com/path/to/page", false},
		{"valid with query", "https://example.com?q=search", false},
		{"valid with port", "https://example.com:8080/api", false},

		// Invalid schemes
		{"ftp scheme", "ftp://example.com", true},
		{"file scheme", "file:///etc/passwd", true},
		{"javascript scheme", "javascript:alert(1)", true},
		{"data scheme", "data:text/html,<h1>hello</h1>", true},
		{"no scheme", "example.com", true},

		// Blocked hosts - SSRF protection
		{"localhost", "http://localhost/admin", true},
		{"127.0.0.1", "http://127.0.0.1/", true},
		{"0.0.0.0", "http://0.0.0.0/", true},
		{"IPv6 localhost", "http://[::1]/", true},
		{"metadata endpoint AWS", "http://169.254.169.254/", true},
		{"metadata endpoint GCP", "http://metadata.google.internal/", true},

		// Private IPs - SSRF protection
		{"private 10.x.x.x", "http://10.0.0.1/admin", true},
		{"private 172.16.x.x", "http://172.16.0.1/admin", true},
		{"private 192.168.x.x", "http://192.168.1.1/admin", true},

		// Invalid URLs
		{"empty string", "", true},
		{"malformed", "not-a-url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		isPrivate bool
	}{
		// Private IPv4 ranges
		{"10.0.0.1", "10.0.0.1", true},
		{"10.255.255.255", "10.255.255.255", true},
		{"172.16.0.1", "172.16.0.1", true},
		{"172.31.255.255", "172.31.255.255", true},
		{"192.168.0.1", "192.168.0.1", true},
		{"192.168.255.255", "192.168.255.255", true},
		{"127.0.0.1", "127.0.0.1", true},
		{"127.0.0.2", "127.0.0.2", true},

		// Link-local
		{"169.254.169.254", "169.254.169.254", true},
		{"169.254.0.1", "169.254.0.1", true},

		// Public IPs
		{"8.8.8.8", "8.8.8.8", false},
		{"1.1.1.1", "1.1.1.1", false},
		{"142.250.80.46", "142.250.80.46", false},

		// Edge cases for non-private
		{"172.32.0.1 (not private)", "172.32.0.1", false},
		{"11.0.0.1 (not private)", "11.0.0.1", false},
		{"192.169.0.1 (not private)", "192.169.0.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPrivateIP(tt.host)
			if got != tt.isPrivate {
				t.Errorf("isPrivateIP(%q) = %v, want %v", tt.host, got, tt.isPrivate)
			}
		})
	}
}

func TestParseCookies(t *testing.T) {
	tests := []struct {
		name       string
		cookieStr  string
		domain     string
		wantCount  int
		wantFirst  string
	}{
		{
			name:      "single cookie",
			cookieStr: "session=abc123",
			domain:    "example.com",
			wantCount: 1,
			wantFirst: "session",
		},
		{
			name:      "multiple cookies",
			cookieStr: "session=abc123; user_id=42; theme=dark",
			domain:    "example.com",
			wantCount: 3,
			wantFirst: "session",
		},
		{
			name:      "cookies with spaces",
			cookieStr: "  name = value  ;  foo = bar  ",
			domain:    "example.com",
			wantCount: 2,
			wantFirst: "name",
		},
		{
			name:      "empty string",
			cookieStr: "",
			domain:    "example.com",
			wantCount: 0,
		},
		{
			name:      "only semicolons",
			cookieStr: ";;;",
			domain:    "example.com",
			wantCount: 0,
		},
		{
			name:      "invalid format (no equals)",
			cookieStr: "invalidcookie",
			domain:    "example.com",
			wantCount: 0,
		},
		{
			name:      "cookie with equals in value",
			cookieStr: "token=abc=def=ghi",
			domain:    "example.com",
			wantCount: 1,
			wantFirst: "token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookies := ParseCookies(tt.cookieStr, tt.domain)
			if len(cookies) != tt.wantCount {
				t.Errorf("ParseCookies() returned %d cookies, want %d", len(cookies), tt.wantCount)
			}
			if tt.wantCount > 0 && tt.wantFirst != "" {
				if cookies[0].Name != tt.wantFirst {
					t.Errorf("First cookie name = %q, want %q", cookies[0].Name, tt.wantFirst)
				}
				if cookies[0].Domain != tt.domain {
					t.Errorf("Cookie domain = %q, want %q", cookies[0].Domain, tt.domain)
				}
			}
		})
	}
}
