package engine_test

import (
	"net/url"
	"strings"
	"testing"
)

// isValidHostingURL validates a URL to ensure it does not match disallowed patterns:
// - Contains "http:" protocol (only https is allowed)
// - Contains "localhost" (with or without a port)
// - Internal or non-routable IP addresses (e.g., 192.168.x.x, 10.x.x.x, 172.16.x.x to 172.31.x.x)
// - Non-routable IPs like 127.x.x.x, 0.0.0.0, or IPv6 loopback (::1)
func isValidHostingURL(hostingURL string) bool {
	if hostingURL == "" {
		return false
	}

	parsedURL, err := url.Parse(hostingURL)
	if err != nil {
		return false
	}

	// Disallow http:
	if parsedURL.Scheme == "http" {
		return false
	}

	// Require a valid scheme
	if parsedURL.Scheme == "" {
		return false
	}

	hostname := parsedURL.Hostname()

	// Disallow localhost (case-insensitive)
	if strings.EqualFold(hostname, "localhost") {
		return false
	}

	// Check for non-routable IPv4 addresses
	if isNonRoutableIPv4(hostname) {
		return false
	}

	// Check for IPv6 loopback
	if hostname == "::1" || hostname == "[::1]" {
		return false
	}

	// Also check Host field for IPv6 without brackets
	if parsedURL.Host == "::1" {
		return false
	}

	return true
}

func isNonRoutableIPv4(ip string) bool {
	// Pattern matching for non-routable IPv4 addresses
	patterns := []struct {
		prefix string
		check  func(string) bool
	}{
		{"127.", func(ip string) bool { return len(ip) >= 4 && ip[:4] == "127." }},
		{"10.", func(ip string) bool { return len(ip) >= 3 && ip[:3] == "10." }},
		{"192.168.", func(ip string) bool { return len(ip) >= 8 && ip[:8] == "192.168." }},
		{"0.0.0.0", func(ip string) bool { return ip == "0.0.0.0" }},
	}

	for _, pattern := range patterns {
		if pattern.check(ip) {
			return true
		}
	}

	// Check for 172.16.x.x to 172.31.x.x
	if len(ip) >= 7 && ip[:4] == "172." {
		// Extract second octet
		secondOctet := ""
		for i := 4; i < len(ip) && ip[i] != '.'; i++ {
			secondOctet += string(ip[i])
		}
		if octet := parseOctet(secondOctet); octet >= 16 && octet <= 31 {
			return true
		}
	}

	return false
}

func parseOctet(s string) int {
	if s == "" {
		return -1
	}
	val := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return -1
		}
		val = val*10 + int(ch-'0')
		if val > 255 {
			return -1
		}
	}
	return val
}

func TestIsValidHostingURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		// Valid URLs
		{"valid public HTTPS URL", "https://example.com", true},
		{"valid public HTTPS URL with port", "https://example.com:8080", true},
		{"valid public HTTPS URL with path", "https://example.com/path", true},
		{"valid IP HTTPS URL", "https://1.2.3.4", true},

		// Invalid URLs - HTTP protocol
		{"HTTP protocol", "http://example.com", false},
		{"HTTP with port", "http://example.com:8080", false},

		// Invalid URLs - localhost
		{"localhost", "https://localhost", false},
		{"localhost with port", "https://localhost:8080", false},
		{"localhost uppercase", "https://LOCALHOST", false},
		{"localhost mixed case", "https://LocalHost:3000", false},

		// Invalid URLs - loopback addresses
		{"loopback 127.0.0.1", "https://127.0.0.1", false},
		{"loopback 127.0.0.1 with port", "https://127.0.0.1:8080", false},
		{"loopback 127.1.2.3", "https://127.1.2.3", false},
		{"IPv6 loopback", "https://[::1]", false},
		{"IPv6 loopback without brackets", "https://::1", false},

		// Invalid URLs - private IP ranges
		{"private IP 10.x", "https://10.0.0.1", false},
		{"private IP 10.x with port", "https://10.255.255.255:8080", false},
		{"private IP 192.168.x", "https://192.168.1.1", false},
		{"private IP 192.168.x with port", "https://192.168.0.1:3000", false},
		{"private IP 172.16.x", "https://172.16.0.1", false},
		{"private IP 172.31.x", "https://172.31.255.255", false},
		{"private IP 172.15.x (valid)", "https://172.15.0.1", true},
		{"private IP 172.32.x (valid)", "https://172.32.0.1", true},

		// Invalid URLs - non-routable
		{"non-routable 0.0.0.0", "https://0.0.0.0", false},

		// Invalid URLs - malformed
		{"empty string", "", false},
		{"invalid URL", "not-a-url", false},
		{"missing protocol", "example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidHostingURL(tt.url)
			if got != tt.expected {
				t.Errorf("isValidHostingURL(%q) = %v, want %v", tt.url, got, tt.expected)
			}
		})
	}
}
