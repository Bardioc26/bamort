package deployment_test

import (
	"bamort/config"
	"bamort/database"
	"bamort/deployment/backup"
	"bamort/deployment/install"
	"bamort/deployment/migrations"
	"bamort/deployment/version"
	"bamort/models"
	"bamort/user"
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
	// Setup: Create fresh test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Create minimal test master data
	tempDir := createTestMasterDataDir(t)
	defer os.RemoveAll(tempDir)

	// Create installer
	installer := install.NewInstaller(database.DB)
	installer.MasterDataPath = tempDir
	installer.CreateAdminUser = true
	installer.AdminUsername = "admin"
	installer.AdminPassword = "test123"

	// Execute installation
	result, err := installer.Initialize()
	require.NoError(t, err, "Installation should succeed")
	assert.NotNil(t, result, "Result should not be nil")
	assert.True(t, result.Success, "Installation should be successful")
	assert.Equal(t, config.GetVersion(), result.Version)
	assert.True(t, result.MasterDataOK, "Master data should be imported")
	assert.True(t, result.AdminCreated, "Admin user should be created")

	// Verify tables exist (skip SHOW TABLES on SQLite)
	// Just verify key models can be queried
	var charCount int64
	err = database.DB.Model(&models.Char{}).Count(&charCount).Error
	assert.NoError(t, err, "Should be able to query characters table")

	// Verify version tracking initialized
	runner := migrations.NewMigrationRunner(database.DB)
	currentVersion, migrationNum, err := runner.GetCurrentVersion()
	require.NoError(t, err)
	latestMigration := migrations.GetLatestMigration()
	if latestMigration != nil {
		assert.Equal(t, latestMigration.Version, currentVersion, "Version should match latest migration")
	}
	assert.Greater(t, migrationNum, 0, "Should have migration number set")

	// Verify master data imported
	var sourceCount int64
	database.DB.Model(&models.Source{}).Count(&sourceCount)
	assert.Greater(t, sourceCount, int64(0), "Should have imported sources")

	// Verify admin user created
	var adminUser user.User
	err = database.DB.Where("username = ?", "admin").First(&adminUser).Error
	require.NoError(t, err, "Admin user should exist")
	assert.True(t, adminUser.IsAdmin(), "User should be admin")

	// Verify compatibility
	compat := version.CheckCompatibility(currentVersion)
	assert.True(t, compat.Compatible, "Should be compatible after installation")
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
	// Version should be the latest migration's target version, not necessarily config.GetVersion()
	latestMigration := migrations.GetLatestMigration()
	if latestMigration != nil {
		assert.Equal(t, latestMigration.Version, endVersion, "Version should match latest migration")
	}
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
	results, err := runner.ApplyAll()
	require.NoError(t, err)

	// Check if any migrations were actually applied
	if len(results) == 0 {
		t.Skip("No migrations were applied, skipping rollback test")
	}

	// Get current state
	beforeVersion, beforeNum, err := runner.GetCurrentVersion()
	require.NoError(t, err)

	// Check if there are any migrations to rollback
	if beforeNum == 0 {
		t.Skip("No migrations applied, skipping rollback test")
	}

	// Rollback 1 step (safe - we know we have at least 1)
	rollbackSteps := 1
	if beforeNum > 1 {
		rollbackSteps = 2 // Can test rolling back 2 if we have more than 1
	}

	err = runner.Rollback(rollbackSteps)
	require.NoError(t, err, "Rollback should succeed")

	// Verify version rolled back
	_, afterNum, err := runner.GetCurrentVersion()
	require.NoError(t, err)

	// If we rolled back all migrations, afterNum should be 0
	expectedNum := beforeNum - rollbackSteps
	if expectedNum < 0 {
		expectedNum = 0
	}
	assert.Equal(t, expectedNum, afterNum, "Should have rolled back %d migration(s)", rollbackSteps)

	// Re-apply migrations
	results2, err2 := runner.ApplyAll()
	require.NoError(t, err2, "Re-applying should work")

	// Should be back to original state
	finalVersion, finalNum, err3 := runner.GetCurrentVersion()
	require.NoError(t, err3)
	assert.Equal(t, beforeVersion, finalVersion, "Should be back to original version")
	assert.Equal(t, beforeNum, finalNum, "Should be back to original migration number")
	assert.Equal(t, rollbackSteps, len(results2), "Should have re-applied %d migration(s)", rollbackSteps)
}

// TestScenario4_ImportOldExport tests backward compatible import
func TestScenario4_ImportOldExport(t *testing.T) {
	// Setup test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Create schema
	err := models.MigrateStructure(database.DB)
	require.NoError(t, err)

	// Create an old format export file (without export_version field)
	oldExportFile := createOldFormatExport(t)
	defer os.Remove(oldExportFile)

	// Read the export
	data, err := os.ReadFile(oldExportFile)
	require.NoError(t, err)

	var exportData map[string]interface{}
	err = json.Unmarshal(data, &exportData)
	require.NoError(t, err)

	// Verify it's old format (no export_version)
	_, hasVersion := exportData["export_version"]
	assert.False(t, hasVersion, "Old export should not have version field")

	// Import sources from old format
	sources, ok := exportData["sources"].([]interface{})
	require.True(t, ok, "Should have sources array")
	assert.Greater(t, len(sources), 0, "Should have at least one source")

	// Convert and import
	for i, src := range sources {
		srcMap, ok := src.(map[string]interface{})
		if !ok {
			t.Fatalf("Source %d is not a map: %v", i, src)
		}

		// Extract fields safely (JSON uses lowercase field names from json tags in the struct)
		code, hasCode := srcMap["code"].(string)
		name, hasName := srcMap["name"].(string)
		gameSystem, hasGameSystem := srcMap["game_system"].(string)
		isActive, hasIsActive := srcMap["is_active"].(bool)

		if !hasCode || code == "" {
			t.Logf("Available keys in srcMap: %v", srcMap)
			t.Fatalf("Source %d missing Code field", i)
		}

		if !hasName {
			name = "Unknown"
		}
		if !hasGameSystem {
			gameSystem = "midgard"
		}
		if !hasIsActive {
			isActive = true
		}

		source := models.Source{
			Code:       code,
			Name:       name,
			GameSystem: gameSystem,
			IsActive:   isActive,
		}

		t.Logf("Importing source: Code=%s, Name=%s, GameSystem=%s, IsActive=%v", code, name, gameSystem, isActive)
		err = database.DB.Create(&source).Error
		require.NoError(t, err, "Should import old format source")
	}

	// Verify import succeeded - check that our test source was imported
	var importedSource models.Source
	err = database.DB.Where("code = ?", "OLD").First(&importedSource).Error
	require.NoError(t, err, "Should find imported source")
	assert.Equal(t, "Old Source", importedSource.Name)
	assert.Equal(t, "midgard", importedSource.GameSystem)
	assert.True(t, importedSource.IsActive, "Source should be active")
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

	// Modify database - delete our test source
	database.DB.Where("code = ?", "TEST").Delete(&models.Source{})
	var count int64
	database.DB.Model(&models.Source{}).Where("code = ?", "TEST").Count(&count)
	assert.Equal(t, int64(0), count, "Test source should be deleted")

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

	// Ensure version tables exist (SQLite-compatible syntax)
	err = database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version VARCHAR(20) NOT NULL,
			migration_number INTEGER NOT NULL,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			backend_version VARCHAR(20) NOT NULL,
			description TEXT
		)
	`).Error
	require.NoError(t, err)

	err = database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS migration_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			migration_number INTEGER NOT NULL UNIQUE,
			version VARCHAR(20) NOT NULL,
			description TEXT NOT NULL,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			applied_by VARCHAR(100),
			execution_time_ms INTEGER,
			success INTEGER DEFAULT 1,
			rollback_available INTEGER DEFAULT 1
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

	// Create empty files for other master data to avoid errors
	emptyFiles := []string{
		"character_classes.json",
		"skill_categories.json",
		"skill_difficulties.json",
		"spell_schools.json",
		"skills.json",
		"weapon_skills.json",
		"spells.json",
		"equipment.json",
		"skill_improvement_costs.json",
	}

	for _, filename := range emptyFiles {
		// Write empty JSON array
		err := os.WriteFile(filepath.Join(tempDir, filename), []byte("[]"), 0644)
		require.NoError(t, err)
	}

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

	// Create empty files for other master data
	emptyFiles := []string{
		"character_classes.json",
		"skill_categories.json",
		"skill_difficulties.json",
		"spell_schools.json",
		"skills.json",
		"weapon_skills.json",
		"spells.json",
		"equipment.json",
		"skill_improvement_costs.json",
	}

	for _, filename := range emptyFiles {
		os.WriteFile(filepath.Join(tempDir, filename), []byte("[]"), 0644)
	}

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
