package pdfrender

import (
	"bamort/config"
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set templates directory for tests (tests run from pdfrender/ directory)
	config.Cfg.TemplatesDir = "../templates"
}

func TestListTemplates(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/templates", ListTemplates)

	// Act
	req, _ := http.NewRequest("GET", "/templates", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []TemplateInfo
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Should return at least one template
	if len(response) == 0 {
		t.Error("Expected at least one template in response")
	}

	// Verify structure of first template
	if len(response) > 0 {
		tmpl := response[0]
		if tmpl.ID == "" {
			t.Error("Template ID should not be empty")
		}
		if tmpl.Name == "" {
			t.Error("Template Name should not be empty")
		}
	}
}

func TestExportCharacterToPDF(t *testing.T) {
	// Arrange
	database.SetupTestDB()

	// Load test character
	char := &models.Char{}
	err := char.FirstID("18")
	if err != nil {
		t.Fatalf("Failed to load test character: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/export/:id", ExportCharacterToPDF)

	// Act - Export with default template
	req, _ := http.NewRequest("GET", "/export/18", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/pdf" {
		t.Errorf("Expected Content-Type 'application/pdf', got '%s'", contentType)
	}

	// Check PDF content
	body := w.Body.Bytes()
	if len(body) == 0 {
		t.Error("Expected non-empty PDF content")
	}

	// Verify PDF marker
	if string(body[0:4]) != "%PDF" {
		t.Error("Response does not start with PDF marker")
	}
}

func TestExportCharacterToPDF_WithTemplate(t *testing.T) {
	// Arrange
	database.SetupTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/export/:id", ExportCharacterToPDF)

	// Act - Export with specific template
	req, _ := http.NewRequest("GET", "/export/18?template=Default_A4_Quer", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verify it's a PDF
	body := w.Body.Bytes()
	if string(body[0:4]) != "%PDF" {
		t.Error("Response does not start with PDF marker")
	}
}

func TestExportCharacterToPDF_CharacterNotFound(t *testing.T) {
	// Arrange
	database.SetupTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/export/:id", ExportCharacterToPDF)

	// Act - Try to export non-existent character
	req, _ := http.NewRequest("GET", "/export/99999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}
