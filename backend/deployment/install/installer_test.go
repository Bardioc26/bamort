package install

import (
	"bamort/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) {
	database.SetupTestDB()
	t.Cleanup(func() {
		database.ResetTestDB()
	})
}

func TestNewInstaller(t *testing.T) {
	setupTestDB(t)

	installer := NewInstaller(database.DB)

	assert.NotNil(t, installer)
	assert.NotNil(t, installer.DB)
	assert.Equal(t, "./masterdata", installer.MasterDataPath)
	assert.False(t, installer.CreateAdminUser)
	assert.Equal(t, "midgard", installer.GameSystem)
}

func TestInitialize_MinimalSetup(t *testing.T) {
	setupTestDB(t)

	installer := NewInstaller(database.DB)
	installer.MasterDataPath = "./testdata" // Use non-existent path for test

	// Should fail because master data path doesn't exist
	result, err := installer.Initialize()

	// Check that we got to the master data import step before failing
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	assert.Contains(t, err.Error(), "master data")
}

func TestInitializeVersionTracking(t *testing.T) {
	setupTestDB(t)

	installer := NewInstaller(database.DB)

	err := installer.initializeVersionTracking()
	assert.NoError(t, err)

	// Verify version table was created and populated
	var version struct {
		Version         string
		MigrationNumber int
		Description     string
	}

	err = installer.DB.Table("schema_version").
		Order("id DESC").
		Limit(1).
		Scan(&version).Error

	assert.NoError(t, err)
	assert.NotEmpty(t, version.Version)
	assert.Greater(t, version.MigrationNumber, 0)
	assert.Equal(t, "Initial installation", version.Description)
}

func TestCreateAdmin(t *testing.T) {
	setupTestDB(t)

	installer := NewInstaller(database.DB)
	installer.CreateAdminUser = true
	installer.AdminUsername = "testadmin"
	installer.AdminPassword = "testpassword123"

	err := installer.createAdmin()
	assert.NoError(t, err)

	// Verify admin user was created
	var count int64
	installer.DB.Table("users").Where("username = ?", "testadmin").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestCreateAdmin_AlreadyExists(t *testing.T) {
	setupTestDB(t)

	installer := NewInstaller(database.DB)
	installer.CreateAdminUser = true
	installer.AdminUsername = "testadmin"
	installer.AdminPassword = "testpassword123"

	// Create once
	err := installer.createAdmin()
	assert.NoError(t, err)

	// Try to create again - should not error, just skip
	err = installer.createAdmin()
	assert.NoError(t, err)

	// Verify only one user exists
	var count int64
	installer.DB.Table("users").Where("username = ?", "testadmin").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestCreateAdmin_NoPassword(t *testing.T) {
	setupTestDB(t)

	installer := NewInstaller(database.DB)
	installer.AdminUsername = "testadmin"
	installer.AdminPassword = "" // Empty password

	err := installer.createAdmin()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password")
}

func TestCreateDatabaseSchema(t *testing.T) {
	setupTestDB(t)

	installer := NewInstaller(database.DB)

	err := installer.createDatabaseSchema()
	assert.NoError(t, err)

	// Verify some key tables exist
	tables := []string{"users", "chars", "gsm_skills", "gsm_spells"}

	for _, table := range tables {
		var exists bool
		err := installer.DB.Raw("SELECT 1 FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&exists).Error
		assert.NoError(t, err, "Failed to check table %s", table)
	}
}
