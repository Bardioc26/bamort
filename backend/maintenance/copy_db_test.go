package maintenance

import (
	"bamort/config"
	"bamort/database"
	"bamort/models"
	"bamort/user"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCopyLiveDatabaseToFile(t *testing.T) {
	// Setup
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "empty_backup.db")
	envpath, _ := filepath.Abs("../.env.test")
	os.Setenv("CONFIG_FILE", envpath)

	config.Cfg = config.LoadConfig()

	// Reset database connection to ensure we use environment config
	database.DB = nil
	database.ConnectDatabase()
	liveDB := database.DB
	require.NotNil(t, liveDB, "Live database should be connected")
	defer func() {
		if sqlDB, err := liveDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Execute
	err := CopyLiveDatabaseToFile(liveDB, targetFile)

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
	assert.GreaterOrEqual(t, userCount, int64(2), "User table should have more that 2 users")

	// Copy target file to database.PreparedTestDB for permanent storage
	// Close the database connection before copying
	if sqlDB, err := targetDB.DB(); err == nil {
		sqlDB.Close()
	}

	// Ensure the directory for PreparedTestDB exists
	preparedDir := filepath.Dir(database.PreparedTestDB)
	err = os.MkdirAll(preparedDir, 0755)
	require.NoError(t, err, "Failed to create directory for PreparedTestDB")

	// Copy the target file to PreparedTestDB using direct file copy
	sourceFile, err := os.Open(targetFile)
	require.NoError(t, err, "Failed to open source file")
	defer sourceFile.Close()

	destFile, err := os.Create(database.PreparedTestDB)
	require.NoError(t, err, "Failed to create destination file")
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	require.NoError(t, err, "Failed to copy file")

	t.Logf("Test database successfully copied to: %s", database.PreparedTestDB)

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
