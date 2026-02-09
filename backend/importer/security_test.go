package importer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_AllowsWithinLimit(t *testing.T) {
	limiter := NewRateLimiter(5, time.Minute)

	// Simulate 5 requests (should all pass)
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		limiter.Middleware()(c)

		assert.False(t, c.IsAborted(), "Request %d should not be aborted", i+1)
	}
}

func TestRateLimiter_BlocksOverLimit(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute)

	// Simulate 3 allowed requests
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		limiter.Middleware()(c)

		assert.False(t, c.IsAborted(), "Request %d should not be aborted", i+1)
	}

	// 4th request should be blocked
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))

	limiter.Middleware()(c)

	assert.True(t, c.IsAborted(), "4th request should be aborted")
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}

func TestRateLimiter_SeparateUsersIndependent(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	// User 1: 2 requests
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(1))

		limiter.Middleware()(c)

		assert.False(t, c.IsAborted())
	}

	// User 2: should still have 2 requests available
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("userID", uint(2))

		limiter.Middleware()(c)

		assert.False(t, c.IsAborted(), "User 2 request %d should not be aborted", i+1)
	}
}

func TestValidateJSONDepth_Valid(t *testing.T) {
	json := []byte(`{"a": {"b": {"c": "value"}}}`)

	err := ValidateJSONDepth(json, 10)

	assert.NoError(t, err)
}

func TestValidateJSONDepth_ExceedsLimit(t *testing.T) {
	json := []byte(`{"a": {"b": {"c": {"d": "value"}}}}`)

	err := ValidateJSONDepth(json, 2)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "depth exceeds maximum")
}

func TestValidateJSONDepth_Array(t *testing.T) {
	json := []byte(`[[[["value"]]]]`)

	err := ValidateJSONDepth(json, 3)

	assert.Error(t, err)
}

func TestSSRFProtection_ValidURL(t *testing.T) {
	protection := NewSSRFProtection([]string{"adapter-foundry:8181", "adapter-vtt:8182"})

	err := protection.ValidateURL("http://adapter-foundry:8181/metadata")

	assert.NoError(t, err)
}

func TestSSRFProtection_NotInWhitelist(t *testing.T) {
	protection := NewSSRFProtection([]string{"adapter-foundry:8181"})

	err := protection.ValidateURL("http://evil-site.com/metadata")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not in whitelist")
}

func TestSSRFProtection_InternalIP(t *testing.T) {
	protection := NewSSRFProtection([]string{"localhost:8181"})

	err := protection.ValidateURL("http://localhost:8181/metadata")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "internal network access forbidden")
}

func TestIsInternalIP_Localhost(t *testing.T) {
	assert.True(t, isInternalIP("localhost"))
	assert.True(t, isInternalIP("127.0.0.1"))
	assert.True(t, isInternalIP("127.0.0.1:8080"))
}

func TestIsInternalIP_PrivateRanges(t *testing.T) {
	assert.True(t, isInternalIP("10.0.0.1"))
	assert.True(t, isInternalIP("172.16.0.1"))
	assert.True(t, isInternalIP("192.168.1.1"))
	assert.True(t, isInternalIP("169.254.1.1"))
}

func TestIsInternalIP_PublicIP(t *testing.T) {
	assert.False(t, isInternalIP("8.8.8.8"))
	assert.False(t, isInternalIP("google.com"))
	assert.False(t, isInternalIP("example.com:443"))
}

func TestNewSecureHTTPClient_DisablesRedirects(t *testing.T) {
	client := NewSecureHTTPClient(5 * time.Second)

	// Test that redirects are disabled
	assert.NotNil(t, client.CheckRedirect)

	err := client.CheckRedirect(nil, nil)
	assert.Equal(t, http.ErrUseLastResponse, err)
}

func TestNewSecureHTTPClient_HasTimeout(t *testing.T) {
	timeout := 10 * time.Second
	client := NewSecureHTTPClient(timeout)

	assert.Equal(t, timeout, client.Timeout)
}
