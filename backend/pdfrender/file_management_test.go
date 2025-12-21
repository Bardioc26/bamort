package pdfrender

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestEnvironment(t *testing.T) {
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original != "" {
			os.Setenv("ENVIRONMENT", original)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})
}

func TestSanitizeFilename(t *testing.T) {
	setupTestEnvironment(t)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple name",
			input:    "Character",
			expected: "Character",
		},
		{
			name:     "Name with spaces",
			input:    "Fanjo Vetrani",
			expected: "Fanjo_Vetrani",
		},
		{
			name:     "Name with umlauts",
			input:    "Müller Ökonom",
			expected: "Mueller_Oekonom",
		},
		{
			name:     "Name with special chars",
			input:    "Test/Character\\Name:With*Special?Chars",
			expected: "Test_Character_Name_With_Special_Chars",
		},
		{
			name:     "Multiple consecutive spaces",
			input:    "Test   Name",
			expected: "Test_Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateExportFilename(t *testing.T) {
	setupTestEnvironment(t)

	// Test with a fixed timestamp
	timestamp := time.Date(2023, 12, 25, 14, 30, 45, 0, time.UTC)
	filename := GenerateExportFilename("Fanjo Vetrani", timestamp)

	expectedPrefix := "Fanjo_Vetrani_20231225_143045"
	if filename != expectedPrefix+".pdf" {
		t.Errorf("GenerateExportFilename() = %q, want prefix %q", filename, expectedPrefix)
	}

	// Test with special characters
	filename2 := GenerateExportFilename("Test/Char*Name", timestamp)
	expectedPrefix2 := "Test_Char_Name_20231225_143045"
	if filename2 != expectedPrefix2+".pdf" {
		t.Errorf("GenerateExportFilename() = %q, want prefix %q", filename2, expectedPrefix2)
	}
}

func TestEnsureExportTempDir(t *testing.T) {
	setupTestEnvironment(t)

	// Use a test-specific temp directory
	testDir := filepath.Join(os.TempDir(), "bamort_test_xporttemp")
	defer os.RemoveAll(testDir)

	// First call should create the directory
	err := EnsureExportTempDir(testDir)
	if err != nil {
		t.Fatalf("EnsureExportTempDir() error = %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(testDir)
	if err != nil {
		t.Fatalf("Directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Path exists but is not a directory")
	}

	// Second call should succeed (directory already exists)
	err = EnsureExportTempDir(testDir)
	if err != nil {
		t.Errorf("EnsureExportTempDir() on existing directory error = %v", err)
	}
}

func TestCleanupOldFiles(t *testing.T) {
	setupTestEnvironment(t)

	// Create test directory
	testDir := filepath.Join(os.TempDir(), "bamort_test_cleanup")
	os.RemoveAll(testDir)
	defer os.RemoveAll(testDir)

	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create test files with different ages
	now := time.Now()

	// Old file (8 days old)
	oldFile := filepath.Join(testDir, "old_file.pdf")
	if err := os.WriteFile(oldFile, []byte("old"), 0644); err != nil {
		t.Fatalf("Failed to create old file: %v", err)
	}
	oldTime := now.Add(-8 * 24 * time.Hour)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to set old file time: %v", err)
	}

	// Recent file (3 days old)
	recentFile := filepath.Join(testDir, "recent_file.pdf")
	if err := os.WriteFile(recentFile, []byte("recent"), 0644); err != nil {
		t.Fatalf("Failed to create recent file: %v", err)
	}
	recentTime := now.Add(-3 * 24 * time.Hour)
	if err := os.Chtimes(recentFile, recentTime, recentTime); err != nil {
		t.Fatalf("Failed to set recent file time: %v", err)
	}

	// New file (1 hour old)
	newFile := filepath.Join(testDir, "new_file.pdf")
	if err := os.WriteFile(newFile, []byte("new"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// Cleanup files older than 7 days
	count, err := CleanupOldFiles(testDir, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("CleanupOldFiles() error = %v", err)
	}

	if count != 1 {
		t.Errorf("CleanupOldFiles() deleted %d files, want 1", count)
	}

	// Verify old file is deleted
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("Old file should be deleted")
	}

	// Verify recent file still exists
	if _, err := os.Stat(recentFile); err != nil {
		t.Error("Recent file should still exist")
	}

	// Verify new file still exists
	if _, err := os.Stat(newFile); err != nil {
		t.Error("New file should still exist")
	}
}
