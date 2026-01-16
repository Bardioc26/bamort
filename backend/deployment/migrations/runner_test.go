package migrations

import (
	"bamort/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) {
	database.SetupTestDB()
	t.Cleanup(func() {
		database.ResetTestDB()
	})
}

func TestNewMigrationRunner(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)

	assert.NotNil(t, runner)
	assert.NotNil(t, runner.DB)
	assert.False(t, runner.DryRun)
	assert.False(t, runner.Verbose)
}

func TestGetCurrentVersion_NoMigrations(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)
	version, number, err := runner.GetCurrentVersion()

	assert.NoError(t, err)
	assert.Equal(t, "", version)
	assert.Equal(t, 0, number)
}

func TestGetPendingMigrations_AllPending(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)
	pending, err := runner.GetPendingMigrations()

	assert.NoError(t, err)
	assert.Len(t, pending, len(AllMigrations))
}

func TestApplyMigration_Success(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)
	runner.Verbose = true

	// Apply first migration
	migration := AllMigrations[0]
	result, err := runner.ApplyMigration(migration)

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, migration.Number, result.Number)
	assert.Greater(t, result.ExecutionTimeMs, int64(0))
	assert.Len(t, result.SQLExecuted, len(migration.UpSQL))

	// Verify version was recorded
	version, number, err := runner.GetCurrentVersion()
	assert.NoError(t, err)
	assert.Equal(t, migration.Version, version)
	assert.Equal(t, migration.Number, number)
}

func TestApplyMigration_DryRun(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)
	runner.DryRun = true
	runner.Verbose = true

	migration := AllMigrations[0]
	result, err := runner.ApplyMigration(migration)

	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Success)

	// Verify nothing was actually applied
	version, number, err := runner.GetCurrentVersion()
	assert.NoError(t, err)
	assert.Equal(t, "", version)
	assert.Equal(t, 0, number)
}

func TestApplyAll_Success(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)
	runner.Verbose = true

	results, err := runner.ApplyAll()

	assert.NoError(t, err)
	require.NotNil(t, results)
	assert.Len(t, results, len(AllMigrations))

	// Verify all migrations succeeded
	for _, result := range results {
		assert.True(t, result.Success)
		assert.NoError(t, result.Error)
	}

	// Verify final version
	version, number, err := runner.GetCurrentVersion()
	assert.NoError(t, err)
	lastMigration := AllMigrations[len(AllMigrations)-1]
	assert.Equal(t, lastMigration.Version, version)
	assert.Equal(t, lastMigration.Number, number)
}

func TestApplyAll_NoPending(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)

	// Apply all first
	_, err := runner.ApplyAll()
	assert.NoError(t, err)

	// Try to apply again - should have no pending
	results, err := runner.ApplyAll()
	assert.NoError(t, err)
	assert.Nil(t, results)
}

func TestGetPendingMigrations_SomeApplied(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)

	// Apply first migration
	migration := AllMigrations[0]
	_, err := runner.ApplyMigration(migration)
	assert.NoError(t, err)

	// Check pending - should be all except first
	pending, err := runner.GetPendingMigrations()
	assert.NoError(t, err)
	assert.Len(t, pending, len(AllMigrations)-1)

	// Verify first pending is second migration
	if len(pending) > 0 {
		assert.Equal(t, AllMigrations[1].Number, pending[0].Number)
	}
}

func TestRollback_Success(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)

	// Apply first migration
	migration := AllMigrations[0]
	_, err := runner.ApplyMigration(migration)
	assert.NoError(t, err)

	// Verify it was applied
	version, number, err := runner.GetCurrentVersion()
	assert.NoError(t, err)
	assert.Equal(t, migration.Version, version)
	assert.Equal(t, migration.Number, number)

	// Rollback
	err = runner.Rollback(1)
	assert.NoError(t, err)

	// Verify rollback
	version, number, err = runner.GetCurrentVersion()
	assert.NoError(t, err)
	assert.Equal(t, "", version)
	assert.Equal(t, 0, number)
}

func TestRollback_NoMigrations(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)

	// Try to rollback when nothing is applied
	err := runner.Rollback(1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no migrations to rollback")
}

func TestRollback_InvalidSteps(t *testing.T) {
	setupTestDB(t)

	runner := NewMigrationRunner(database.DB)

	err := runner.Rollback(0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")

	err = runner.Rollback(-1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be positive")
}

func TestGetMigrationByNumber(t *testing.T) {
	migration := GetMigrationByNumber(1)
	assert.NotNil(t, migration)
	assert.Equal(t, 1, migration.Number)

	migration = GetMigrationByNumber(9999)
	assert.Nil(t, migration)
}

func TestGetLatestMigration(t *testing.T) {
	migration := GetLatestMigration()
	assert.NotNil(t, migration)
	assert.Equal(t, AllMigrations[len(AllMigrations)-1].Number, migration.Number)
}
