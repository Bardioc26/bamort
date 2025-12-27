package pdfrender

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeTemplates(t *testing.T) {
	// Setup test directories
	tmpDir := t.TempDir()
	defaultDir := filepath.Join(tmpDir, "default_templates")
	targetDir := filepath.Join(tmpDir, "templates")

	// Create default templates directory with test files
	if err := os.MkdirAll(filepath.Join(defaultDir, "TestTemplate"), 0755); err != nil {
		t.Fatalf("Failed to create default template dir: %v", err)
	}
	testFile := filepath.Join(defaultDir, "TestTemplate", "page1.html")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test 1: Copy when target directory doesn't exist
	if err := InitializeTemplates(defaultDir, targetDir); err != nil {
		t.Errorf("InitializeTemplates failed: %v", err)
	}

	// Verify file was copied
	copiedFile := filepath.Join(targetDir, "TestTemplate", "page1.html")
	if _, err := os.Stat(copiedFile); os.IsNotExist(err) {
		t.Error("Expected file was not copied")
	}

	content, err := os.ReadFile(copiedFile)
	if err != nil || string(content) != "test content" {
		t.Error("Copied file content doesn't match")
	}

	// Test 2: Don't overwrite when content is identical
	if err := InitializeTemplates(defaultDir, targetDir); err != nil {
		t.Errorf("InitializeTemplates failed on second run: %v", err)
	}

	// Verify file still has same content
	content, err = os.ReadFile(copiedFile)
	if err != nil || string(content) != "test content" {
		t.Error("File content should remain unchanged")
	}

	// Test 3: Update when default template changed
	updatedContent := []byte("updated test content")
	if err := os.WriteFile(testFile, updatedContent, 0644); err != nil {
		t.Fatalf("Failed to update source file: %v", err)
	}

	if err := InitializeTemplates(defaultDir, targetDir); err != nil {
		t.Errorf("InitializeTemplates failed after source update: %v", err)
	}

	// Verify file was updated
	content, err = os.ReadFile(copiedFile)
	if err != nil || string(content) != "updated test content" {
		t.Error("File should have been updated with new content")
	}

	// Test 4: Handle missing default directory gracefully
	if err := InitializeTemplates("/nonexistent", targetDir); err == nil {
		t.Error("Expected error for nonexistent default directory")
	}
}

func TestInitializeTemplatesWithMultipleFiles(t *testing.T) {
	tmpDir := t.TempDir()
	defaultDir := filepath.Join(tmpDir, "default_templates")
	targetDir := filepath.Join(tmpDir, "templates")

	// Create multiple templates and files
	templates := []string{"Template1", "Template2"}
	for _, tmpl := range templates {
		tmplDir := filepath.Join(defaultDir, tmpl)
		if err := os.MkdirAll(tmplDir, 0755); err != nil {
			t.Fatalf("Failed to create template dir: %v", err)
		}
		for i := 1; i <= 3; i++ {
			file := filepath.Join(tmplDir, filepath.Base(tmpl)+".html")
			if err := os.WriteFile(file, []byte("content "+tmpl), 0644); err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
		}
	}

	// Initialize templates
	if err := InitializeTemplates(defaultDir, targetDir); err != nil {
		t.Fatalf("InitializeTemplates failed: %v", err)
	}

	// Verify all templates were copied
	for _, tmpl := range templates {
		tmplDir := filepath.Join(targetDir, tmpl)
		if _, err := os.Stat(tmplDir); os.IsNotExist(err) {
			t.Errorf("Template directory %s was not copied", tmpl)
		}
	}
}
