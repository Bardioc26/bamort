package masterdata

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

func TestNewMasterDataSync(t *testing.T) {
	setupTestDB(t)

	sync := NewMasterDataSync(database.DB, "./testdata")

	assert.NotNil(t, sync)
	assert.NotNil(t, sync.DB)
	assert.Equal(t, "./testdata", sync.ImportDir)
	assert.False(t, sync.DryRun)
	assert.False(t, sync.Verbose)
}

func TestSyncAll_DryRun(t *testing.T) {
	setupTestDB(t)

	sync := NewMasterDataSync(database.DB, "./testdata")
	sync.DryRun = true
	sync.Verbose = true

	// In dry-run mode, should not error even if directory doesn't exist
	err := sync.SyncAll()
	assert.NoError(t, err)
}

func TestSyncAll_InvalidDirectory(t *testing.T) {
	setupTestDB(t)

	sync := NewMasterDataSync(database.DB, "/nonexistent/path")
	sync.Verbose = true

	// Should error when trying to import from non-existent directory
	err := sync.SyncAll()
	assert.Error(t, err)
}
