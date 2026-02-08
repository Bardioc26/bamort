package transfero

import (
	"bamort/config"
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	database.SetupTestDB(true, true)
	models.MigrateStructure()
	/*
		t.Cleanup(func() {
			database.ResetTestDB()
		})
	*/
	return database.DB
}

func TestExportDatabase_Success(t *testing.T) {
	//setupTestEnvironment(t)
	db := setupTestDB(t)
	if db.Error != nil {
		t.Fatalf("Failed to setup test DB: %v", db.Error)
	}

	// Create temporary export directory
	exportDir := t.TempDir()
	config.LoadConfig()
	exportDir = config.Cfg.ExportTempDir

	// Export database
	result, err := ExportDatabase(exportDir)

	// Assertions
	require.NoError(t, err, "ExportDatabase should succeed")
	assert.NotNil(t, result, "Result should not be nil")
	assert.NotEmpty(t, result.Filename, "Filename should be set")
	assert.NotEmpty(t, result.FilePath, "FilePath should be set")
	assert.True(t, result.RecordCount > 0, "Should export some records")

	// Verify file exists
	assert.FileExists(t, result.FilePath, "Export file should exist")

	// Verify file is valid JSON
	data, err := os.ReadFile(result.FilePath)
	require.NoError(t, err, "Should be able to read export file")

	var export DatabaseExport
	err = json.Unmarshal(data, &export)
	require.NoError(t, err, "Export file should be valid JSON")
}

func TestExportDatabase_InvalidDirectory(t *testing.T) {
	//setupTestEnvironment(t)
	db := setupTestDB(t)
	if db.Error != nil {
		t.Fatalf("Failed to setup test DB: %v", db.Error)
	}

	// Try to export to non-existent directory
	result, err := ExportDatabase("/invalid/path/that/does/not/exist")

	assert.Error(t, err, "Should fail with invalid directory")
	assert.Nil(t, result, "Result should be nil on error")
}

func TestImportDatabase_Success(t *testing.T) {
	//setupTestEnvironment(t)
	db := setupTestDB(t)
	if db.Error != nil {
		t.Fatalf("Failed to setup test DB: %v", db.Error)
	}

	// First export database
	exportDir := t.TempDir()
	exportResult, err := ExportDatabase(exportDir)
	require.NoError(t, err, "Export should succeed")

	// Clear database
	db.Exec("DELETE FROM characters")
	db.Exec("DELETE FROM users")

	// Import database
	result, err := ImportDatabase(exportResult.FilePath)

	// Assertions
	require.NoError(t, err, "ImportDatabase should succeed")
	assert.NotNil(t, result, "Result should not be nil")
	assert.True(t, result.RecordCount > 0, "Should import some records")
}

func TestImportDatabase_FileNotFound(t *testing.T) {
	//setupTestEnvironment(t)
	db := setupTestDB(t)
	if db.Error != nil {
		t.Fatalf("Failed to setup test DB: %v", db.Error)
	}

	result, err := ImportDatabase("/nonexistent/file.json")

	assert.Error(t, err, "Should fail with non-existent file")
	assert.Nil(t, result, "Result should be nil on error")
}

func TestImportDatabase_InvalidJSON(t *testing.T) {
	//setupTestEnvironment(t)
	db := setupTestDB(t)
	if db.Error != nil {
		t.Fatalf("Failed to setup test DB: %v", db.Error)
	}

	// Create invalid JSON file
	tmpFile := filepath.Join(t.TempDir(), "invalid.json")
	err := os.WriteFile(tmpFile, []byte("invalid json content"), 0644)
	require.NoError(t, err)

	result, err := ImportDatabase(tmpFile)

	assert.Error(t, err, "Should fail with invalid JSON")
	assert.Nil(t, result, "Result should be nil on error")
}

func TestExportImportRoundtrip(t *testing.T) {
	//setupTestEnvironment(t)
	db := setupTestDB(t)
	if db.Error != nil {
		t.Fatalf("Failed to setup test DB: %v", db.Error)
	}

	exportDir := t.TempDir()

	// Export
	exportResult, err := ExportDatabase(exportDir)
	require.NoError(t, err, "Export should succeed")

	originalCount := exportResult.RecordCount

	// Import
	importResult, err := ImportDatabase(exportResult.FilePath)
	require.NoError(t, err, "Import should succeed")

	// Record counts should match
	assert.Equal(t, originalCount, importResult.RecordCount,
		"Import should restore same number of records")
}
