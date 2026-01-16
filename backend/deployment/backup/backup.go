package backup

import (
	"bamort/config"
	"bamort/logger"
	"bamort/transfer"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// BackupService handles database backups
type BackupService struct {
	BackupDir string
}

// BackupMetadata contains metadata about a backup
type BackupMetadata struct {
	Timestamp       time.Time
	Version         string
	MigrationNumber int
	Method          string // "json" or "sqldump"
	FilePath        string
	SizeBytes       int64
}

// NewBackupService creates a new backup service
func NewBackupService() *BackupService {
	backupDir := filepath.Join(".", "backups")
	return &BackupService{
		BackupDir: backupDir,
	}
}

// EnsureBackupDir ensures the backup directory exists
func (s *BackupService) EnsureBackupDir() error {
	if err := os.MkdirAll(s.BackupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	return nil
}

// CreateJSONBackup creates a JSON backup using the existing transfer package
func (s *BackupService) CreateJSONBackup(version string, migrationNumber int) (*BackupMetadata, error) {
	if err := s.EnsureBackupDir(); err != nil {
		return nil, err
	}

	timestamp := time.Now()
	filename := fmt.Sprintf("backup_%s_v%s_m%d.json",
		timestamp.Format("20060102_150405"),
		version,
		migrationNumber,
	)
	filepath := filepath.Join(s.BackupDir, filename)

	logger.Info("Creating JSON backup: %s", filename)

	// Use the existing export functionality
	result, err := transfer.ExportDatabase(s.BackupDir)
	if err != nil {
		return nil, fmt.Errorf("database export failed: %w", err)
	}

	// Rename the export file to our backup filename
	if err := os.Rename(result.FilePath, filepath); err != nil {
		return nil, fmt.Errorf("failed to rename export file: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat backup file: %w", err)
	}

	metadata := &BackupMetadata{
		Timestamp:       timestamp,
		Version:         version,
		MigrationNumber: migrationNumber,
		Method:          "json",
		FilePath:        filepath,
		SizeBytes:       fileInfo.Size(),
	}

	logger.Info("JSON backup created: %s (%d bytes)", filename, metadata.SizeBytes)
	return metadata, nil
}

// CreateMariaDBDump creates a MariaDB dump backup (only works in production with MySQL)
func (s *BackupService) CreateMariaDBDump(version string, migrationNumber int) (*BackupMetadata, error) {
	if config.Cfg.DatabaseType != "mysql" {
		return nil, fmt.Errorf("MariaDB dump only available for MySQL databases")
	}

	if err := s.EnsureBackupDir(); err != nil {
		return nil, err
	}

	timestamp := time.Now()
	filename := fmt.Sprintf("backup_%s_v%s_m%d.sql",
		timestamp.Format("20060102_150405"),
		version,
		migrationNumber,
	)
	filepath := filepath.Join(s.BackupDir, filename)

	logger.Info("Creating MariaDB dump: %s", filename)

	// Use docker exec to create mysqldump
	// This assumes we're running in docker-compose environment
	cmd := exec.Command("docker", "exec", "bamort-mariadb-dev",
		"mysqldump",
		"-u", "bamort",
		"-pbG4)efozrc",
		"--single-transaction",
		"--routines",
		"--triggers",
		"bamort",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("mysqldump failed: %w - Output: %s", err, string(output))
	}

	// Write dump to file
	if err := os.WriteFile(filepath, output, 0644); err != nil {
		return nil, fmt.Errorf("failed to write dump file: %w", err)
	}

	// Get file size
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat backup file: %w", err)
	}

	metadata := &BackupMetadata{
		Timestamp:       timestamp,
		Version:         version,
		MigrationNumber: migrationNumber,
		Method:          "sqldump",
		FilePath:        filepath,
		SizeBytes:       fileInfo.Size(),
	}

	logger.Info("MariaDB dump created: %s (%d bytes)", filename, metadata.SizeBytes)
	return metadata, nil
}

// CleanupOldBackups removes backups older than the retention period
func (s *BackupService) CleanupOldBackups(retentionDays int) error {
	logger.Info("Cleaning up backups older than %d days", retentionDays)

	entries, err := os.ReadDir(s.BackupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No backup directory yet
		}
		return fmt.Errorf("failed to read backup directory: %w", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	deletedCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(s.BackupDir, entry.Name())
		fileInfo, err := entry.Info()
		if err != nil {
			logger.Warn("Failed to get info for %s: %v", entry.Name(), err)
			continue
		}

		if fileInfo.ModTime().Before(cutoffTime) {
			logger.Info("Deleting old backup: %s (age: %v)", entry.Name(), time.Since(fileInfo.ModTime()))
			if err := os.Remove(filePath); err != nil {
				logger.Warn("Failed to delete %s: %v", entry.Name(), err)
			} else {
				deletedCount++
			}
		}
	}

	logger.Info("Cleanup complete: deleted %d old backup(s)", deletedCount)
	return nil
}

// ListBackups returns a list of all backups
func (s *BackupService) ListBackups() ([]string, error) {
	entries, err := os.ReadDir(s.BackupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []string
	for _, entry := range entries {
		if !entry.IsDir() {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}
