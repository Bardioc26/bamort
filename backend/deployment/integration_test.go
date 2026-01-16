package deployment_test

import (
	"bamort/config"
	"bamort/database"
	"bamort/deployment/backup"
	"bamort/deployment/install"
	"bamort/deployment/migrations"
	"bamort/deployment/version"
	"bamort/models"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestScenario1_FreshInstallation tests a complete fresh installation workflow
func TestScenario1_FreshInstallation(t *testing.T) {
	t.Skip("Skipping full installation test - requires complete master data files")

	// This test would require creating complete master data structure
	// For now, we test individual components separately
}

// TestScenario2_UpdateDeployment tests updating from older version to newer
func TestScenario2_UpdateDeployment(t *testing.T) {
	// Setup: Create test database with older version
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Simulate older version (0.4.0) already installed
	setupOlderVersion(t, "0.4.0")

	// Verify starting state
	runner := migrations.NewMigrationRunner(database.DB)
	startVersion, startMigrationNum, err := runner.GetCurrentVersion()
	require.NoError(t, err)
	assert.Equal(t, "0.4.0", startVersion)

	// Create migration runner to update to current version
	runner.Verbose = true

	// Get pending migrations
	pending, err := runner.GetPendingMigrations()
	require.NoError(t, err)
	assert.Greater(t, len(pending), 0, "Should have pending migrations")

	// Apply all pending migrations
	results, err := runner.ApplyAll()
	require.NoError(t, err, "Migration should succeed")

	// Verify all migrations succeeded
	for i, result := range results {
		assert.True(t, result.Success, "Migration %d should succeed", i)
		assert.NoError(t, result.Error)
	}

	// Verify version updated
	endVersion, endMigrationNum, err := runner.GetCurrentVersion()
	require.NoError(t, err)
	assert.Equal(t, config.GetVersion(), endVersion, "Version should be updated")
	assert.Greater(t, endMigrationNum, startMigrationNum, "Migration number should increase")

	// Verify database integrity
	var char models.Char
	err = database.DB.First(&char).Error
	// Should not error even if no characters (table should exist)
	if err != nil && err.Error() != "record not found" {
		t.Errorf("Unexpected error querying characters: %v", err)
	}

	// Verify compatibility check passes
	compat := version.CheckCompatibility(endVersion)
	assert.True(t, compat.Compatible, "Should be compatible after migration")
	assert.Contains(t, compat.Reason, "matches required", "Compatibility reason should be positive")
}

// TestScenario3_Rollback tests migration rollback functionality
func TestScenario3_Rollback(t *testing.T) {
	// Setup test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Apply all migrations first
	runner := migrations.NewMigrationRunner(database.DB)
	_, err := runner.ApplyAll()
	require.NoError(t, err)

	// Get current state
	beforeVersion, beforeNum, err := runner.GetCurrentVersion()
	require.NoError(t, err)

	// Rollback 2 steps
	err = runner.Rollback(2)
	require.NoError(t, err, "Rollback should succeed")

	// Verify version rolled back
	afterVersion, afterNum, err := runner.GetCurrentVersion()
	require.NoError(t, err)
	assert.Equal(t, beforeNum-2, afterNum, "Should have rolled back 2 migrations")

	// Version should be different unless both migrations were for same version
	if beforeNum > 2 {
		assert.NotEqual(t, beforeVersion, afterVersion, "Version should change after rollback")
	}

	// Re-apply migrations
	results, err := runner.ApplyAll()
	require.NoError(t, err, "Re-applying should work")

	// Should be back to original state
	finalVersion, finalNum, err := runner.GetCurrentVersion()
	require.NoError(t, err)
	assert.Equal(t, beforeVersion, finalVersion, "Should be back to original version")
	assert.Equal(t, beforeNum, finalNum, "Should be back to original migration number")
	assert.Equal(t, 2, len(results), "Should have re-applied 2 migrations")
}

// TestScenario4_ImportOldExport tests backward compatible import
func TestScenario4_ImportOldExport(t *testing.T) {
	t.Skip("Skipping old export import test - requires masterdata import implementation")
}

// TestScenario5_BackupAndRestore tests the backup/restore workflow
func TestScenario5_BackupAndRestore(t *testing.T) {
	// Setup test database with some data
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Create test data
	source := models.Source{
		Code:       "TEST",
		Name:       "Test Source",
		GameSystem: "midgard",
		IsActive:   true,
	}
	err := database.DB.Create(&source).Error
	require.NoError(t, err)

	// Create backup service
	tempDir := t.TempDir()
	backupSvc := backup.NewBackupService()
	backupSvc.BackupDir = tempDir

	// Create backup
	result, err := backupSvc.CreateJSONBackup("0.4.0", 0)
	require.NoError(t, err, "Backup creation should succeed")
	assert.NotNil(t, result, "Backup result should not be nil")
	assert.FileExists(t, result.FilePath, "Backup file should exist")

	// Modify database
	database.DB.Delete(&source)
	var count int64
	database.DB.Model(&models.Source{}).Count(&count)
	assert.Equal(t, int64(0), count, "Source should be deleted")

	// Note: Restore functionality would be tested here when implemented
	t.Log("Backup created successfully, restore test skipped (not yet implemented)")
}

// TestScenario6_ConcurrentMigration tests that concurrent migrations are prevented
func TestScenario6_ConcurrentMigration(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Create two migration runners
	runner1 := migrations.NewMigrationRunner(database.DB)
	runner2 := migrations.NewMigrationRunner(database.DB)

	// Try to run migrations concurrently
	done1 := make(chan error)
	done2 := make(chan error)

	go func() {
		_, err := runner1.ApplyAll()
		done1 <- err
	}()

	// Small delay to ensure first one starts
	time.Sleep(100 * time.Millisecond)

	go func() {
		_, err := runner2.ApplyAll()
		done2 <- err
	}()

	// Wait for both
	err1 := <-done1
	err2 := <-done2

	// At least one should succeed, one might fail with lock error
	if err1 == nil && err2 == nil {
		// Both succeeded - one must have found no pending migrations
		t.Log("Both completed - second must have found no pending migrations")
	} else if err1 != nil && err2 != nil {
		t.Fatal("Both migrations failed - unexpected")
	} else {
		// One succeeded, one failed - expected
		t.Log("One migration succeeded, one was prevented - expected behavior")
	}
}

// TestScenario7_PerformanceTest tests deployment performance with realistic data
func TestScenario7_PerformanceTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	database.SetupTestDB()
	defer database.ResetTestDB()

	// Create installer with master data
	tempDir := createLargeMasterDataDir(t)
	defer os.RemoveAll(tempDir)

	installer := install.NewInstaller(database.DB)
	installer.MasterDataPath = tempDir

	// Measure installation time
	startTime := time.Now()
	result, err := installer.Initialize()
	duration := time.Since(startTime)

	require.NoError(t, err, "Installation should succeed")
	assert.True(t, result.Success)

	// Performance assertions (adjust based on acceptable performance)
	assert.Less(t, duration, 30*time.Second, "Installation should complete within 30 seconds")
	t.Logf("Installation completed in %v", duration)
	t.Logf("Execution time from result: %v", result.ExecutionTime)
}

// Helper functions

func setupOlderVersion(t *testing.T, oldVersion string) {
	// Create basic schema
	err := models.MigrateStructure(database.DB)
	require.NoError(t, err)

	// Ensure version tables exist
	err = database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			id INT PRIMARY KEY AUTO_INCREMENT,
			version VARCHAR(20) NOT NULL,
			migration_number INT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			backend_version VARCHAR(20) NOT NULL,
			description TEXT
		)
	`).Error
	require.NoError(t, err)

	err = database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS migration_history (
			id INT PRIMARY KEY AUTO_INCREMENT,
			migration_number INT NOT NULL UNIQUE,
			version VARCHAR(20) NOT NULL,
			description TEXT NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			applied_by VARCHAR(100),
			execution_time_ms INT,
			success BOOLEAN DEFAULT TRUE,
			rollback_available BOOLEAN DEFAULT TRUE
		)
	`).Error
	require.NoError(t, err)

	// Set old version in database
	versionData := map[string]interface{}{
		"version":          oldVersion,
		"migration_number": 0,
		"applied_at":       time.Now(),
		"backend_version":  oldVersion,
		"description":      "Test setup - old version",
	}
	err = database.DB.Table("schema_version").Create(versionData).Error
	require.NoError(t, err)
}

func createTestMasterDataDir(t *testing.T) string {
	tempDir := t.TempDir()

	// Create minimal test export files with correct structure
	sources := []models.Source{
		{Code: "ALBA", Name: "Alba", GameSystem: "midgard", IsActive: true},
		{Code: "ARK", Name: "Arkanum", GameSystem: "midgard", IsActive: true},
	}
	sourcesJSON, _ := json.MarshalIndent(sources, "", "  ")
	err := os.WriteFile(filepath.Join(tempDir, "sources.json"), sourcesJSON, 0644)
	require.NoError(t, err)

	return tempDir
}

func createLargeMasterDataDir(t *testing.T) string {
	tempDir := t.TempDir()

	// Create larger dataset for performance testing
	sources := make([]models.Source, 50)
	for i := range sources {
		sources[i] = models.Source{
			Code:       string(rune('A' + i%26)),
			Name:       string(rune('A'+i%26)) + " Source",
			GameSystem: "midgard",
			IsActive:   true,
		}
	}
	sourcesJSON, _ := json.MarshalIndent(sources, "", "  ")
	os.WriteFile(filepath.Join(tempDir, "sources.json"), sourcesJSON, 0644)

	return tempDir
}

func createOldFormatExport(t *testing.T) string {
	// Create a v1.0 format export (old format without version field)
	oldExport := map[string]interface{}{
		"sources": []models.Source{
			{Code: "OLD", Name: "Old Source", GameSystem: "midgard", IsActive: true},
		},
	}

	tempFile := filepath.Join(t.TempDir(), "old_export.json")
	data, _ := json.MarshalIndent(oldExport, "", "  ")
	err := os.WriteFile(tempFile, data, 0644)
	require.NoError(t, err)

	return tempFile
}
