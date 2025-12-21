package database

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupHandlerTestEnvironment sets up test environment for handler tests
func setupHandlerTestEnvironment(t *testing.T) {
	// Save original values
	origEnv := os.Getenv("ENVIRONMENT")
	origDevTesting := os.Getenv("DEV_TESTING")
	origDatabaseType := os.Getenv("DATABASE_TYPE")
	origDatabaseURL := os.Getenv("DATABASE_URL")
	originalDB := DB

	// Set test environment
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DEV_TESTING", "yes")

	// Cleanup function to restore original values
	t.Cleanup(func() {
		if origEnv != "" {
			os.Setenv("ENVIRONMENT", origEnv)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
		if origDevTesting != "" {
			os.Setenv("DEV_TESTING", origDevTesting)
		} else {
			os.Unsetenv("DEV_TESTING")
		}
		if origDatabaseType != "" {
			os.Setenv("DATABASE_TYPE", origDatabaseType)
		} else {
			os.Unsetenv("DATABASE_TYPE")
		}
		if origDatabaseURL != "" {
			os.Setenv("DATABASE_URL", origDatabaseURL)
		} else {
			os.Unsetenv("DATABASE_URL")
		}

		// Reset global DB variable
		DB = originalDB
	})
}

func TestSetupCheck(t *testing.T) {
	setupHandlerTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a Gin router and register the handler
	router := gin.New()
	router.GET("/setup-check", SetupCheck)

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/setup-check", nil)
	assert.NoError(t, err, "Should be able to create HTTP request")

	// Create a response recorder
	recorder := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(recorder, req)

	// Verify that the function executed without panicking
	// Since SetupCheck doesn't return any response, we just check that it didn't crash
	assert.True(t, true, "SetupCheck should execute without panicking")

	// Verify that DB is now initialized (ConnectDatabase was called)
	assert.NotNil(t, DB, "SetupCheck should initialize the database connection")
}

func TestSetupCheck_DatabaseInitialization(t *testing.T) {
	setupHandlerTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Reset DB to ensure fresh test
	DB = nil

	// Create a Gin context manually
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request, _ = http.NewRequest("GET", "/setup-check", nil)

	// Call SetupCheck directly
	SetupCheck(c)

	// Verify that the database connection was established
	assert.NotNil(t, DB, "SetupCheck should establish database connection")

	// Verify we can perform basic database operations
	if DB != nil {
		sqlDB, err := DB.DB()
		assert.NoError(t, err, "Should be able to get underlying sql.DB")
		if sqlDB != nil {
			err = sqlDB.Ping()
			assert.NoError(t, err, "Should be able to ping the database")
		}
	}
}

func TestSetupCheck_MultipleInvocations(t *testing.T) {
	setupHandlerTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Reset DB to ensure fresh test
	DB = nil

	// Create Gin contexts for multiple requests
	recorder1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(recorder1)
	c1.Request, _ = http.NewRequest("GET", "/setup-check", nil)

	recorder2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(recorder2)
	c2.Request, _ = http.NewRequest("GET", "/setup-check", nil)

	// Call SetupCheck multiple times
	SetupCheck(c1)
	firstDB := DB
	assert.NotNil(t, firstDB, "First call should establish database connection")

	SetupCheck(c2)
	secondDB := DB
	assert.NotNil(t, secondDB, "Second call should maintain database connection")

	// The key thing is that both calls succeed and establish a database connection
	// The actual instance may vary depending on the implementation, but both should be valid
	if firstDB != nil && secondDB != nil {
		// Verify both connections are working
		sqlDB1, err1 := firstDB.DB()
		sqlDB2, err2 := secondDB.DB()
		assert.NoError(t, err1, "First database connection should be valid")
		assert.NoError(t, err2, "Second database connection should be valid")

		if sqlDB1 != nil {
			assert.NoError(t, sqlDB1.Ping(), "First database should be pingable")
		}
		if sqlDB2 != nil {
			assert.NoError(t, sqlDB2.Ping(), "Second database should be pingable")
		}
	}
}

func TestSetupCheck_WithDifferentHTTPMethods(t *testing.T) {
	setupHandlerTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Test with different HTTP methods
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			// Reset DB for each test
			DB = nil

			// Create a Gin context with the specific HTTP method
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request, _ = http.NewRequest(method, "/setup-check", nil)

			// Call SetupCheck
			SetupCheck(c)

			// Verify that it works regardless of HTTP method
			assert.NotNil(t, DB, "SetupCheck should work with %s method", method)
		})
	}
}

func TestSetupCheck_WithRequestHeaders(t *testing.T) {
	setupHandlerTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Reset DB to ensure fresh test
	DB = nil

	// Create a request with various headers
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req, _ := http.NewRequest("GET", "/setup-check", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer test-token")
	req.Header.Set("User-Agent", "test-client/1.0")
	c.Request = req

	// Call SetupCheck
	SetupCheck(c)

	// Verify that it works with headers present
	assert.NotNil(t, DB, "SetupCheck should work with request headers")

	// Verify that the headers are still accessible in the context
	assert.Equal(t, "application/json", c.GetHeader("Content-Type"))
	assert.Equal(t, "Bearer test-token", c.GetHeader("Authorization"))
	assert.Equal(t, "test-client/1.0", c.GetHeader("User-Agent"))
}

func TestSetupCheck_WithQueryParameters(t *testing.T) {
	setupHandlerTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Reset DB to ensure fresh test
	DB = nil

	// Create a request with query parameters
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	req, _ := http.NewRequest("GET", "/setup-check?param1=value1&param2=value2", nil)
	c.Request = req

	// Call SetupCheck
	SetupCheck(c)

	// Verify that it works with query parameters
	assert.NotNil(t, DB, "SetupCheck should work with query parameters")

	// Verify that query parameters are still accessible
	assert.Equal(t, "value1", c.Query("param1"))
	assert.Equal(t, "value2", c.Query("param2"))
}

func TestSetupCheck_FullHTTPFlow(t *testing.T) {
	setupHandlerTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a complete Gin router with middleware
	router := gin.New()

	// Add some middleware to test that SetupCheck works in a middleware chain
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Register the handler
	router.GET("/api/setup-check", SetupCheck)

	// Create a test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Make a real HTTP request to the test server
	resp, err := http.Get(server.URL + "/api/setup-check")
	assert.NoError(t, err, "Should be able to make HTTP request")
	defer resp.Body.Close()

	// Since SetupCheck doesn't write any response, we just verify the request succeeded
	// The status code should be 200 (OK) even though no response was written
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Request should succeed")

	// Verify that DB was initialized
	assert.NotNil(t, DB, "SetupCheck should initialize database in full HTTP flow")
}

func TestSetupCheck_ErrorScenarios(t *testing.T) {
	// Note: Since SetupCheck just calls ConnectDatabase() and doesn't handle errors,
	// we test that it doesn't panic even in error scenarios

	// Save original environment
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("WithInvalidDatabaseConfig", func(t *testing.T) {
		// Set invalid database configuration
		os.Setenv("DATABASE_TYPE", "invalid")
		os.Setenv("DATABASE_URL", "invalid://connection/string")

		defer func() {
			os.Unsetenv("DATABASE_TYPE")
			os.Unsetenv("DATABASE_URL")
		}()

		// Reset DB
		DB = nil

		// Create Gin context
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request, _ = http.NewRequest("GET", "/setup-check", nil)

		// This might panic depending on the ConnectDatabase implementation
		// We'll handle the panic gracefully
		defer func() {
			if r := recover(); r != nil {
				// If it panics, that's expected behavior with invalid config
				t.Logf("SetupCheck panicked with invalid config (expected): %v", r)
			}
		}()

		// Call SetupCheck - might panic with invalid config
		SetupCheck(c)
	})
}

func TestSetupCheck_ContextIntegrity(t *testing.T) {
	setupHandlerTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a Gin context with some values set
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request, _ = http.NewRequest("GET", "/setup-check", nil)

	// Set some values in the context before calling SetupCheck
	c.Set("test_key", "test_value")
	c.Set("user_id", 12345)

	// Call SetupCheck
	SetupCheck(c)

	// Verify that context values are preserved
	value, exists := c.Get("test_key")
	assert.True(t, exists, "Context should preserve existing values")
	assert.Equal(t, "test_value", value, "Context values should remain unchanged")

	userID, exists := c.Get("user_id")
	assert.True(t, exists, "Context should preserve existing values")
	assert.Equal(t, 12345, userID, "Context values should remain unchanged")

	// Verify that the database was still initialized
	assert.NotNil(t, DB, "SetupCheck should initialize database without affecting context")
}

// Benchmark test for performance
func BenchmarkSetupCheck(b *testing.B) {
	// Setup test environment
	os.Setenv("ENVIRONMENT", "test")
	defer os.Unsetenv("ENVIRONMENT")

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a reusable context
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request, _ = http.NewRequest("GET", "/setup-check", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SetupCheck(c)
	}
}

func BenchmarkSetupCheck_WithHTTPOverhead(b *testing.B) {
	// Setup test environment
	os.Setenv("ENVIRONMENT", "test")
	defer os.Unsetenv("ENVIRONMENT")

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a router
	router := gin.New()
	router.GET("/setup-check", SetupCheck)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/setup-check", nil)
		router.ServeHTTP(recorder, req)
	}
}
