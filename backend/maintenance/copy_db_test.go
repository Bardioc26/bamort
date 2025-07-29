package maintenance

import (
	"bamort/models"
	"bamort/user"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestCopyLiveDatabaseToFile_Success tests the main functionality of copyLiveDatabaseToFile
func TestCopyLiveDatabaseToFile_Success(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "test_backup.db")

	// Create test live database with data using GORM AutoMigrate for simplicity
	liveDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test live database")
	defer func() {
		if sqlDB, err := liveDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Migrate only basic structures we can test with
	err = liveDB.AutoMigrate(&user.User{}, &models.Char{})
	require.NoError(t, err, "Failed to migrate test structures")

	// Create simple test data
	testUser := &user.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	err = liveDB.Create(testUser).Error
	require.NoError(t, err, "Failed to create test user")

	testChar := &models.Char{
		BamortBase: models.BamortBase{Name: "Test Hero"},
		Rasse:      "Human",
		Typ:        "Warrior",
		Alter:      25,
		Grad:       1,
	}
	err = liveDB.Create(testChar).Error
	require.NoError(t, err, "Failed to create test character")

	// Execute - test the exported function
	err = CopyLiveDatabaseToFile(liveDB, targetFile)

	// Verify
	require.NoError(t, err, "CopyLiveDatabaseToFile should succeed")

	// Check that target file exists
	assert.FileExists(t, targetFile, "Target SQLite file should be created")

	// Verify target database contains expected data
	targetDB, err := gorm.Open(sqlite.Open(targetFile), &gorm.Config{})
	require.NoError(t, err, "Should be able to open target database")
	defer func() {
		if sqlDB, err := targetDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Verify data was copied correctly
	var sourceUserCount, targetUserCount int64
	err = liveDB.Model(&user.User{}).Count(&sourceUserCount).Error
	require.NoError(t, err, "Failed to count users in source database")

	err = targetDB.Model(&user.User{}).Count(&targetUserCount).Error
	require.NoError(t, err, "Failed to count users in target database")

	assert.Equal(t, sourceUserCount, targetUserCount, "User count should match")
	assert.Greater(t, targetUserCount, int64(0), "Should have copied user data")

	// Verify specific user data
	var sourceUser, targetUser user.User
	err = liveDB.Where("username = ?", "testuser").First(&sourceUser).Error
	require.NoError(t, err, "Failed to find test user in source")

	err = targetDB.Where("username = ?", "testuser").First(&targetUser).Error
	require.NoError(t, err, "Failed to find test user in target")

	assert.Equal(t, sourceUser.Username, targetUser.Username, "Username should match")
	assert.Equal(t, sourceUser.Email, targetUser.Email, "Email should match")
}

// TestCopyLiveDatabaseToFile_BackupExisting tests file backup functionality
func TestCopyLiveDatabaseToFile_BackupExisting(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "test_backup.db")
	backupFile := targetFile + ".backup"

	// Create existing file
	existingContent := "existing database content"
	err := os.WriteFile(targetFile, []byte(existingContent), 0644)
	require.NoError(t, err, "Failed to create existing file")

	// Create minimal test live database
	liveDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test live database")
	defer func() {
		if sqlDB, err := liveDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Just migrate user table for simple test
	err = liveDB.AutoMigrate(&user.User{})
	require.NoError(t, err, "Failed to migrate user structure")

	// Execute
	err = CopyLiveDatabaseToFile(liveDB, targetFile)

	// Verify
	require.NoError(t, err, "CopyLiveDatabaseToFile should succeed")

	// Check that backup file was created
	assert.FileExists(t, backupFile, "Backup file should be created")

	// Verify backup contains original content
	backupContent, err := os.ReadFile(backupFile)
	require.NoError(t, err, "Should be able to read backup file")
	assert.Equal(t, existingContent, string(backupContent), "Backup should contain original content")

	// Verify new target file is a valid SQLite database
	targetDB, err := gorm.Open(sqlite.Open(targetFile), &gorm.Config{})
	require.NoError(t, err, "New target file should be valid SQLite database")
	defer func() {
		if sqlDB, err := targetDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()
}

// TestCopyLiveDatabaseToFile_EmptyDatabase tests with empty database
func TestCopyLiveDatabaseToFile_EmptyDatabase(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "empty_backup.db")

	// Create empty live database (only migrate structures, no data)
	liveDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create empty live database")
	defer func() {
		if sqlDB, err := liveDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Migrate only user table for simple test
	err = liveDB.AutoMigrate(&user.User{})
	require.NoError(t, err, "Failed to migrate user structure")

	// Execute
	err = CopyLiveDatabaseToFile(liveDB, targetFile)

	// Verify
	require.NoError(t, err, "CopyLiveDatabaseToFile should succeed with empty database")
	assert.FileExists(t, targetFile, "Target file should be created")

	// Verify target database has structures but no data
	targetDB, err := gorm.Open(sqlite.Open(targetFile), &gorm.Config{})
	require.NoError(t, err, "Should be able to open target database")
	defer func() {
		if sqlDB, err := targetDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Check that tables exist but are empty
	var userCount int64
	err = targetDB.Model(&user.User{}).Count(&userCount).Error
	require.NoError(t, err, "Should be able to count users")
	assert.Equal(t, int64(0), userCount, "User table should be empty")
}

// BenchmarkCopyLiveDatabaseToFile benchmarks the copy function performance
func BenchmarkCopyLiveDatabaseToFile(b *testing.B) {
	// Setup test database once
	liveDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}
	defer func() {
		if sqlDB, err := liveDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Migrate minimal structures
	liveDB.AutoMigrate(&user.User{}, &models.Char{})

	// Add some test data
	testUser := &user.User{Username: "benchuser", Email: "bench@test.com", PasswordHash: "hash"}
	liveDB.Create(testUser)

	tempDir := b.TempDir()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		targetFile := filepath.Join(tempDir, fmt.Sprintf("benchmark_%d.db", i))

		err := CopyLiveDatabaseToFile(liveDB, targetFile)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}

		// Clean up for next iteration
		os.Remove(targetFile)
	}
}
