package importer

import (
	"bamort/database"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTestEnvironment sets up the test environment
func setupTestEnvironment(t *testing.T) {
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original == "" {
			os.Unsetenv("ENVIRONMENT")
		} else {
			os.Setenv("ENVIRONMENT", original)
		}
	})
}

// setupTestHandlerDB initializes the test database for handler tests
func setupTestHandlerDB(t *testing.T) {
	setupTestEnvironment(t)
	// For now, just ensure migrations run - full DB setup will be added later
	if database.DB != nil {
		err := MigrateStructure(database.DB)
		assert.NoError(t, err)
	}

	t.Cleanup(func() {
		// Clean up test data
		if database.DB != nil {
			database.DB.Exec("DELETE FROM import_histories")
			database.DB.Exec("DELETE FROM master_data_imports")
		}
	})
}

// createTestFile creates a multipart file for testing
func createTestFile(t *testing.T, fieldName, filename string, content []byte) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	assert.NoError(t, err)

	_, err = part.Write(content)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	return body, writer.FormDataContentType()
}

// mockAdapter creates a mock adapter for testing
func mockAdapter() *AdapterMetadata {
	return &AdapterMetadata{
		ID:                  "test-adapter-v1",
		Name:                "Test Adapter",
		Version:             "1.0",
		BmrtVersions:        []string{"1.0"},
		SupportedExtensions: []string{".json"},
		BaseURL:             "http://localhost:8999",
		Capabilities:        []string{"import", "export", "detect"},
		Healthy:             true,
	}
}

// TestDetectHandler tests the format detection endpoint
func TestDetectHandler(t *testing.T) {
	setupTestHandlerDB(t)
	gin.SetMode(gin.TestMode)

	// Setup mock registry
	registry := NewAdapterRegistry()
	registry.Register(*mockAdapter())
	InitializeRegistry(registry)

	// Create test file
	testData := []byte(`{"name": "Test Character", "type": "test"}`)
	body, contentType := createTestFile(t, "file", "test.json", testData)

	// Create request
	req := httptest.NewRequest("POST", "/api/import/detect", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	// Setup context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler (with mock detect - would need adapter server for real test)
	// For now, we verify the handler validates input correctly
	DetectHandler(c)

	// Without a real adapter server, we expect an error
	// This test mainly validates the request parsing
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestImportHandler_NoFile tests import with missing file
func TestImportHandler_NoFile(t *testing.T) {
	setupTestHandlerDB(t)
	gin.SetMode(gin.TestMode)

	registry := NewAdapterRegistry()
	InitializeRegistry(registry)

	req := httptest.NewRequest("POST", "/api/import/import", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", uint(1))

	ImportHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "No file uploaded")
}

// TestImportHandler_NoAuth tests import without authentication
func TestImportHandler_NoAuth(t *testing.T) {
	setupTestHandlerDB(t)
	gin.SetMode(gin.TestMode)

	registry := NewAdapterRegistry()
	InitializeRegistry(registry)

	testData := []byte(`{"name": "Test"}`)
	body, contentType := createTestFile(t, "file", "test.json", testData)

	req := httptest.NewRequest("POST", "/api/import/import", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	// No userID set - simulates unauthenticated request

	ImportHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestListAdaptersHandler tests listing registered adapters
func TestListAdaptersHandler(t *testing.T) {
	setupTestHandlerDB(t)
	gin.SetMode(gin.TestMode)

	registry := NewAdapterRegistry()
	registry.Register(*mockAdapter())
	InitializeRegistry(registry)

	req := httptest.NewRequest("GET", "/api/import/adapters", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	ListAdaptersHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "adapters")
	assert.Contains(t, response, "count")
}

// TestImportHistoryHandler tests import history retrieval
func TestImportHistoryHandler(t *testing.T) {
	setupTestHandlerDB(t)
	gin.SetMode(gin.TestMode)

	// Skip if database not available
	if database.DB == nil {
		t.Skip("Database not available")
	}

	registry := NewAdapterRegistry()
	InitializeRegistry(registry)

	req := httptest.NewRequest("GET", "/api/import/history", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", uint(1))

	ImportHistoryHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "histories")
	assert.Contains(t, response, "total")
	assert.Contains(t, response, "page")
}

// TestImportHistoryHandler_Pagination tests pagination
func TestImportHistoryHandler_Pagination(t *testing.T) {
	setupTestHandlerDB(t)
	gin.SetMode(gin.TestMode)

	// Skip if database not available
	if database.DB == nil {
		t.Skip("Database not available")
	}

	registry := NewAdapterRegistry()
	InitializeRegistry(registry)

	req := httptest.NewRequest("GET", "/api/import/history?page=2&per_page=10", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", uint(1))

	ImportHistoryHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(2), response["page"])
	assert.Equal(t, float64(10), response["per_page"])
}

// TestExportHandler_NotFound tests export with non-existent character
func TestExportHandler_NotFound(t *testing.T) {
	setupTestHandlerDB(t)
	gin.SetMode(gin.TestMode)

	// Skip if database not available
	if database.DB == nil {
		t.Skip("Database not available")
	}

	registry := NewAdapterRegistry()
	InitializeRegistry(registry)

	req := httptest.NewRequest("POST", "/api/import/export/999999", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", uint(1))
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999999"}}

	ExportHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestValidateFileSizeMiddleware tests file size validation
func TestValidateFileSizeMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create middleware with 100 byte limit
	middleware := ValidateFileSizeMiddleware(100)

	// Create large content
	largeContent := bytes.Repeat([]byte("a"), 200)

	req := httptest.NewRequest("POST", "/test", bytes.NewReader(largeContent))
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	middleware(c)

	// Try to read more than limit - should error
	_, err := io.ReadAll(c.Request.Body)
	assert.Error(t, err)
}
