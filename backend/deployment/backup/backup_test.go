package backup

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBackupService(t *testing.T) {
	service := NewBackupService()

	assert.NotNil(t, service)
	assert.NotEmpty(t, service.BackupDir)
	assert.Contains(t, service.BackupDir, "backups")
}

func TestEnsureBackupDir(t *testing.T) {
	tempDir := t.TempDir()
	service := &BackupService{
		BackupDir: filepath.Join(tempDir, "test-backups"),
	}

	err := service.EnsureBackupDir()
	assert.NoError(t, err)

	info, err := os.Stat(service.BackupDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestListBackups_NoDirectory(t *testing.T) {
	service := &BackupService{
		BackupDir: filepath.Join(t.TempDir(), "nonexistent"),
	}

	backups, err := service.ListBackups()
	assert.NoError(t, err)
	assert.Empty(t, backups)
}

func TestListBackups_WithFiles(t *testing.T) {
	tempDir := t.TempDir()
	service := &BackupService{
		BackupDir: tempDir,
	}

	testFiles := []string{
		"backup1.json",
		"backup2.sql",
	}

	for _, filename := range testFiles {
		fp := filepath.Join(tempDir, filename)
		err := os.WriteFile(fp, []byte("test"), 0644)
		require.NoError(t, err)
	}

	backups, err := service.ListBackups()
	assert.NoError(t, err)
	assert.Len(t, backups, 2)
}

func TestCleanupOldBackups(t *testing.T) {
	tempDir := t.TempDir()
	service := &BackupService{
		BackupDir: tempDir,
	}

	oldFile := filepath.Join(tempDir, "old_backup.json")
	newFile := filepath.Join(tempDir, "new_backup.json")

	require.NoError(t, os.WriteFile(oldFile, []byte("old"), 0644))
	require.NoError(t, os.WriteFile(newFile, []byte("new"), 0644))

	oldTime := time.Now().AddDate(0, 0, -31)
	require.NoError(t, os.Chtimes(oldFile, oldTime, oldTime))

	err := service.CleanupOldBackups(30)
	assert.NoError(t, err)

	_, err = os.Stat(oldFile)
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(newFile)
	assert.NoError(t, err)
}

func TestBackupMetadata(t *testing.T) {
	metadata := &BackupMetadata{
		Timestamp:       time.Now(),
		Version:         "0.4.0",
		MigrationNumber: 1,
		Method:          "json",
		FilePath:        "/path/to/backup.json",
		SizeBytes:       1024,
	}

	assert.Equal(t, "0.4.0", metadata.Version)
	assert.Equal(t, 1, metadata.MigrationNumber)
	assert.Equal(t, "json", metadata.Method)
	assert.Equal(t, int64(1024), metadata.SizeBytes)
}
