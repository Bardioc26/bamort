package transfer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bamort/database"
	"bamort/models"

	"github.com/gin-gonic/gin"
)

func setupHandlerTestEnvironment(t *testing.T) *gin.Engine {
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	database.SetupTestDB(true, true)
	models.MigrateStructure()

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	t.Cleanup(func() {
		database.ResetTestDB()
		if original == "" {
			os.Unsetenv("ENVIRONMENT")
		} else {
			os.Setenv("ENVIRONMENT", original)
		}
	})

	return r
}

func TestExportCharacterHandlerAPI(t *testing.T) {
	r := setupHandlerTestEnvironment(t)

	// Register routes
	api := r.Group("/api")
	RegisterRoutes(api)

	// Test export endpoint
	req, _ := http.NewRequest("GET", "/api/transfer/export/18", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify response is valid JSON
	var exportData CharacterExport
	err := json.Unmarshal(w.Body.Bytes(), &exportData)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if exportData.Character.ID != 18 {
		t.Errorf("Expected character ID 18, got %d", exportData.Character.ID)
	}
}

func TestDownloadCharacterHandlerAPI(t *testing.T) {
	r := setupHandlerTestEnvironment(t)

	api := r.Group("/api")
	RegisterRoutes(api)

	// Test download endpoint
	req, _ := http.NewRequest("GET", "/api/transfer/download/18", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify content-type is application/json
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Verify Content-Disposition header exists
	disposition := w.Header().Get("Content-Disposition")
	if disposition == "" {
		t.Error("Expected Content-Disposition header to be set")
	}
}

func TestImportCharacterHandlerAPI(t *testing.T) {
	r := setupHandlerTestEnvironment(t)

	api := r.Group("/api")
	RegisterRoutes(api)

	// First export a character
	exportData, err := ExportCharacter(uint(18))
	if err != nil {
		t.Fatalf("Failed to export character: %v", err)
	}

	// Modify for import
	exportData.Character.ID = 0
	exportData.Character.Name = "API Imported Character"

	// Convert to JSON
	jsonData, err := json.Marshal(exportData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Test import endpoint
	req, _ := http.NewRequest("POST", "/api/transfer/import", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Mock user_id in context (normally set by auth middleware)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	ImportCharacterHandler(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify response contains character_id
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, exists := response["character_id"]; !exists {
		t.Error("Expected character_id in response")
	}
}

func TestExportNonExistentCharacterAPI(t *testing.T) {
	r := setupHandlerTestEnvironment(t)

	api := r.Group("/api")
	RegisterRoutes(api)

	// Test with non-existent character ID
	req, _ := http.NewRequest("GET", "/api/transfer/export/999999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestImportInvalidJSONAPI(t *testing.T) {
	r := setupHandlerTestEnvironment(t)

	api := r.Group("/api")
	RegisterRoutes(api)

	// Test with invalid JSON
	req, _ := http.NewRequest("POST", "/api/transfer/import", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("user_id", uint(1))

	ImportCharacterHandler(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
