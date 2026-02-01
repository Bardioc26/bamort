package transfer

import (
	"bamort/config"
	"bamort/database"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportDatabaseHandler_Success(t *testing.T) {
	setupTestEnvironment(t)
	db := setupTestDB(t)
	database.DB = db

	// Setup temporary export directory
	exportDir := t.TempDir()
	config.Cfg.ExportTempDir = exportDir

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/transfer/database/export", ExportDatabaseHandler)

	// Create request
	req, _ := http.NewRequest("POST", "/api/transfer/database/export", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Response should be valid JSON")

	assert.Equal(t, "Database exported successfully", response["message"])
	assert.NotEmpty(t, response["filename"])
	assert.NotEmpty(t, response["filepath"])
	assert.Greater(t, response["record_count"], float64(0))
}

func TestImportDatabaseHandler_Success(t *testing.T) {
	setupTestEnvironment(t)
	db := setupTestDB(t)
	database.DB = db

	// First, export the database
	exportDir := t.TempDir()
	exportResult, err := ExportDatabase(exportDir)
	require.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/transfer/database/import", ImportDatabaseHandler)

	// Create request
	requestBody := map[string]string{
		"filepath": exportResult.FilePath,
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/transfer/database/import", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "Response should be valid JSON")

	assert.Equal(t, "Database imported successfully", response["message"])
	assert.Greater(t, response["record_count"], float64(0))
	assert.Equal(t, float64(exportResult.RecordCount), response["record_count"])
}

func TestImportDatabaseHandler_MissingFilepath(t *testing.T) {
	setupTestEnvironment(t)
	db := setupTestDB(t)
	database.DB = db

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/transfer/database/import", ImportDatabaseHandler)

	// Create request without filepath
	requestBody := map[string]string{}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/transfer/database/import", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 Bad Request")
}

func TestImportDatabaseHandler_InvalidFilepath(t *testing.T) {
	setupTestEnvironment(t)
	db := setupTestDB(t)
	database.DB = db

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/transfer/database/import", ImportDatabaseHandler)

	// Create request with invalid filepath
	requestBody := map[string]string{
		"filepath": "/nonexistent/path/to/file.json",
	}
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/api/transfer/database/import", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code, "Should return 500 Internal Server Error")
}
