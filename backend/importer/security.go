package importer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements per-user rate limiting with sliding time window
type RateLimiter struct {
	requests map[uint][]time.Time // userID -> request timestamps
	mu       sync.RWMutex
	limit    int           // requests per window
	window   time.Duration // time window
}

// NewRateLimiter creates a new rate limiter with specified limit and window
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[uint][]time.Time),
		limit:    limit,
		window:   window,
	}
}

// Middleware returns a Gin middleware function for rate limiting
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c) // Extract from JWT token

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		cutoff := now.Add(-rl.window)

		// Remove expired timestamps
		timestamps := rl.requests[userID]
		valid := make([]time.Time, 0)
		for _, t := range timestamps {
			if t.After(cutoff) {
				valid = append(valid, t)
			}
		}

		// Check limit
		if len(valid) >= rl.limit {
			c.JSON(429, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": rl.window.Seconds(),
			})
			c.Abort()
			return
		}

		// Add current request
		valid = append(valid, now)
		rl.requests[userID] = valid

		c.Next()
	}
}

// getUserID extracts the user ID from the JWT token in the context
// This is a placeholder - actual implementation depends on your auth system
func getUserID(c *gin.Context) uint {
	// TODO: Extract from JWT token or session
	// This is a placeholder implementation
	userIDInterface, exists := c.Get("userID")
	if !exists {
		return 0
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		return 0
	}

	return userID
}

// ValidateFileSizeMiddleware limits upload file size
func ValidateFileSizeMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// ValidateJSONDepth prevents deeply nested JSON attacks
func ValidateJSONDepth(data []byte, maxDepth int) error {
	var depth int
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber() // Prevent float precision issues

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch token {
		case json.Delim('{'), json.Delim('['):
			depth++
			if depth > maxDepth {
				return fmt.Errorf("JSON depth exceeds maximum of %d levels", maxDepth)
			}
		case json.Delim('}'), json.Delim(']'):
			depth--
		}
	}

	return nil
}

// SSRFProtection provides SSRF attack prevention via URL whitelisting
type SSRFProtection struct {
	allowedHosts []string // Whitelist of adapter hosts
}

// NewSSRFProtection creates a new SSRF protection instance
func NewSSRFProtection(allowedHosts []string) *SSRFProtection {
	return &SSRFProtection{allowedHosts: allowedHosts}
}

// ValidateURL checks if a URL is in the whitelist and not an internal IP
func (s *SSRFProtection) ValidateURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Block redirects to internal networks
	if isInternalIP(parsed.Host) {
		return fmt.Errorf("internal network access forbidden")
	}

	// Check whitelist
	allowed := false
	for _, host := range s.allowedHosts {
		if strings.HasPrefix(parsed.Host, host) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("host %s not in whitelist", parsed.Host)
	}

	return nil
}

// isInternalIP checks if a host is an internal/private IP address
func isInternalIP(host string) bool {
	// Remove port if present
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}

	internal := []string{
		"localhost",
		"127.",
		"10.",
		"172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.",
		"172.24.", "172.25.", "172.26.", "172.27.",
		"172.28.", "172.29.", "172.30.", "172.31.",
		"192.168.",
		"169.254.", // Link-local
	}

	for _, prefix := range internal {
		if strings.HasPrefix(host, prefix) {
			return true
		}
	}

	return false
}

// NewSecureHTTPClient creates an HTTP client with security settings
func NewSecureHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Disable redirects
		},
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 2,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
			DisableCompression:  false,
		},
	}
}
