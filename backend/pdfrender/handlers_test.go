package pdfrender

import (
	"bamort/config"
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// Set templates directory for tests (tests run from pdfrender/ directory)
	config.Cfg.TemplatesDir = "../templates"
	config.Cfg.ExportTempDir = "../xporttemp"
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

	// Use test-specific temp directory
	testDir := "../xporttemp_test_basic"
	config.Cfg.ExportTempDir = testDir
	defer func() {
		os.RemoveAll(testDir)
		config.Cfg.ExportTempDir = "../xporttemp"
	}()

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
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Check content type - now returns JSON
	contentType := w.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Parse JSON response
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Verify filename is returned
	filename, ok := response["filename"]
	if !ok {
		t.Fatal("Response should contain 'filename' field")
	}

	if filename == "" {
		t.Error("Filename should not be empty")
	}

	// Verify PDF file exists
	filePath := filepath.Join(testDir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read PDF file: %v", err)
	}

	// Verify PDF marker
	if string(data[0:4]) != "%PDF" {
		t.Error("File does not start with PDF marker")
	}
}

func TestExportCharacterToPDF_WithTemplate(t *testing.T) {
	// Arrange
	database.SetupTestDB()

	// Use test-specific temp directory
	testDir := "../xporttemp_test_template"
	config.Cfg.ExportTempDir = testDir
	defer func() {
		os.RemoveAll(testDir)
		config.Cfg.ExportTempDir = "../xporttemp"
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/export/:id", ExportCharacterToPDF)

	// Act - Export with specific template
	req, _ := http.NewRequest("GET", "/export/18?template=Default_A4_Quer", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Parse JSON response
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	filename := response["filename"]
	filePath := filepath.Join(testDir, filename)

	// Verify it's a PDF
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read PDF: %v", err)
	}

	if string(data[0:4]) != "%PDF" {
		t.Error("File does not start with PDF marker")
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

func TestExportCharacterToPDF_SavesFileAndReturnsFilename(t *testing.T) {
	// Arrange
	database.SetupTestDB()

	// Use test-specific temp directory
	testDir := "../xporttemp_test"
	config.Cfg.ExportTempDir = testDir
	defer func() {
		// Cleanup
		os.RemoveAll(testDir)
		config.Cfg.ExportTempDir = "../xporttemp"
	}()

	// Load test character
	char := &models.Char{}
	err := char.FirstID("18")
	if err != nil {
		t.Fatalf("Failed to load test character: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/export/:id", ExportCharacterToPDF)

	// Act
	req, _ := http.NewRequest("GET", "/export/18", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Parse response
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	// Check that filename is returned
	filename, ok := response["filename"]
	if !ok {
		t.Fatal("Response should contain 'filename' field")
	}

	if filename == "" {
		t.Error("Filename should not be empty")
	}

	// Verify filename format (should contain character name and timestamp)
	if !strings.Contains(filename, "Fanjo_Vetrani") {
		t.Errorf("Filename should contain sanitized character name, got: %s", filename)
	}

	if !strings.HasSuffix(filename, ".pdf") {
		t.Errorf("Filename should end with .pdf, got: %s", filename)
	}

	// Verify file exists in xporttemp directory
	filePath := filepath.Join(testDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("PDF file should exist at %s", filePath)
	}

	// Verify it's a valid PDF
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read PDF file: %v", err)
	}

	if len(data) == 0 {
		t.Error("PDF file should not be empty")
	}

	if string(data[0:4]) != "%PDF" {
		t.Error("File should be a valid PDF")
	}
}

func TestGetPDFFile(t *testing.T) {
	setupTestEnvironment(t)

	// Create test directory and file
	testDir := "../xporttemp_test_get"
	config.Cfg.ExportTempDir = testDir
	defer func() {
		os.RemoveAll(testDir)
		config.Cfg.ExportTempDir = "../xporttemp"
	}()

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a test PDF file
	testFilename := "Test_Character_20231225_120000.pdf"
	testPDFContent := []byte("%PDF-1.4\nTest PDF Content")
	filePath := filepath.Join(testDir, testFilename)
	if err := os.WriteFile(filePath, testPDFContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/file/:filename", GetPDFFile)

	// Act - Get the PDF file
	req, _ := http.NewRequest("GET", "/file/"+testFilename, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/pdf" {
		t.Errorf("Expected Content-Type 'application/pdf', got '%s'", contentType)
	}

	// Check content
	body := w.Body.Bytes()
	if string(body) != string(testPDFContent) {
		t.Errorf("Expected PDF content, got different content")
	}
}

func TestGetPDFFile_NotFound(t *testing.T) {
	setupTestEnvironment(t)

	testDir := "../xporttemp_test_notfound"
	config.Cfg.ExportTempDir = testDir
	defer func() {
		os.RemoveAll(testDir)
		config.Cfg.ExportTempDir = "../xporttemp"
	}()

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/file/:filename", GetPDFFile)

	// Act - Try to get non-existent file
	req, _ := http.NewRequest("GET", "/file/nonexistent.pdf", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestGetPDFFile_PathTraversal(t *testing.T) {
	setupTestEnvironment(t)

	testDir := "../xporttemp_test_security"
	config.Cfg.ExportTempDir = testDir
	defer func() {
		os.RemoveAll(testDir)
		config.Cfg.ExportTempDir = "../xporttemp"
	}()

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/file/:filename", GetPDFFile)

	// Act - Try path traversal attack
	req, _ := http.NewRequest("GET", "/file/../../../etc/passwd", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert - Should return error (either 404 or 400)
	if w.Code == http.StatusOK {
		t.Error("Should not allow path traversal attacks")
	}
}

func TestCleanupEndpoint(t *testing.T) {
	setupTestEnvironment(t)

	// Create test directory with old and new files
	testDir := "../xporttemp_test_cleanup_endpoint"
	config.Cfg.ExportTempDir = testDir
	defer func() {
		os.RemoveAll(testDir)
		config.Cfg.ExportTempDir = "../xporttemp"
	}()

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	now := time.Now()

	// Create old file (8 days old)
	oldFile := filepath.Join(testDir, "old_file.pdf")
	if err := os.WriteFile(oldFile, []byte("%PDF-old"), 0644); err != nil {
		t.Fatalf("Failed to create old file: %v", err)
	}
	oldTime := now.Add(-8 * 24 * time.Hour)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to set old file time: %v", err)
	}

	// Create recent file (3 days old)
	recentFile := filepath.Join(testDir, "recent_file.pdf")
	if err := os.WriteFile(recentFile, []byte("%PDF-recent"), 0644); err != nil {
		t.Fatalf("Failed to create recent file: %v", err)
	}
	recentTime := now.Add(-3 * 24 * time.Hour)
	if err := os.Chtimes(recentFile, recentTime, recentTime); err != nil {
		t.Fatalf("Failed to set recent file time: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/cleanup", CleanupExportTemp)

	// Act
	req, _ := http.NewRequest("POST", "/cleanup", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check deleted count
	deletedCount, ok := response["deleted"]
	if !ok {
		t.Fatal("Response should contain 'deleted' field")
	}

	// Should have deleted 1 file (the old one)
	if int(deletedCount.(float64)) != 1 {
		t.Errorf("Expected 1 file deleted, got %v", deletedCount)
	}

	// Verify old file is deleted
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("Old file should be deleted")
	}

	// Verify recent file still exists
	if _, err := os.Stat(recentFile); err != nil {
		t.Error("Recent file should still exist")
	}
}
