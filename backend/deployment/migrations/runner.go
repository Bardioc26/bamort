package migrations

import (
	"bamort/config"
	"bamort/logger"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MigrationRunner handles database migration execution
type MigrationRunner struct {
	DB      *gorm.DB
	DryRun  bool
	Verbose bool
}

// MigrationResult contains the result of a migration execution
type MigrationResult struct {
	Number          int
	Description     string
	Success         bool
	ExecutionTimeMs int64
	Error           error
	SQLExecuted     []string
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *gorm.DB) *MigrationRunner {
	return &MigrationRunner{
		DB:      db,
		DryRun:  false,
		Verbose: false,
	}
}

// GetCurrentVersion returns the current database version and migration number
func (r *MigrationRunner) GetCurrentVersion() (string, int, error) {
	var version struct {
		Version         string
		MigrationNumber int
	}

	err := r.DB.Raw(`
		SELECT version, migration_number 
		FROM schema_version 
		ORDER BY id DESC 
		LIMIT 1
	`).Scan(&version).Error

	if err == gorm.ErrRecordNotFound || err != nil {
		// No migrations applied yet or table doesn't exist
		return "", 0, nil
	}

	return version.Version, version.MigrationNumber, nil
}

// GetPendingMigrations returns all migrations that haven't been applied yet
func (r *MigrationRunner) GetPendingMigrations() ([]Migration, error) {
	_, currentNumber, err := r.GetCurrentVersion()
	if err != nil {
		return nil, err
	}

	var pending []Migration
	for _, m := range AllMigrations {
		if m.Number > currentNumber {
			pending = append(pending, m)
		}
	}

	return pending, nil
}

// ApplyMigration applies a single migration
func (r *MigrationRunner) ApplyMigration(m Migration) (*MigrationResult, error) {
	startTime := time.Now()
	result := &MigrationResult{
		Number:      m.Number,
		Description: m.Description,
	}

	if r.Verbose {
		logger.Info("Applying migration %d: %s", m.Number, m.Description)
	}

	// Transaction for safety
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Execute SQL statements
		for _, sql := range m.UpSQL {
			if r.Verbose {
				logger.Debug("Executing SQL: %s", sql)
			}

			if r.DryRun {
				logger.Info("[DRY RUN] Would execute: %s", sql)
				result.SQLExecuted = append(result.SQLExecuted, sql)
				continue
			}

			if err := tx.Exec(sql).Error; err != nil {
				return fmt.Errorf("SQL failed: %s - Error: %w", sql, err)
			}
			result.SQLExecuted = append(result.SQLExecuted, sql)
		}

		// Execute data migration function if exists
		if m.DataFunc != nil && !r.DryRun {
			if r.Verbose {
				logger.Debug("Executing data migration function")
			}
			if err := m.DataFunc(tx); err != nil {
				return fmt.Errorf("data migration failed: %w", err)
			}
		}

		// Record migration in history
		if !r.DryRun {
			now := time.Now().Unix()
			history := map[string]interface{}{
				"migration_number":   m.Number,
				"version":            m.Version,
				"description":        m.Description,
				"applied_at":         now,
				"applied_by":         "migration-runner",
				"execution_time_ms":  time.Since(startTime).Milliseconds(),
				"success":            true,
				"rollback_available": len(m.DownSQL) > 0,
			}

			if err := tx.Table("migration_history").Create(history).Error; err != nil {
				return fmt.Errorf("failed to record migration: %w", err)
			}

			// Update schema_version
			version := map[string]interface{}{
				"version":          m.Version,
				"migration_number": m.Number,
				"applied_at":       now,
				"backend_version":  config.GetVersion(),
				"description":      m.Description,
			}

			if err := tx.Table("schema_version").Create(version).Error; err != nil {
				return fmt.Errorf("failed to update version: %w", err)
			}
		}

		return nil
	})

	result.ExecutionTimeMs = time.Since(startTime).Milliseconds()

	if err != nil {
		result.Success = false
		result.Error = err
		return result, err
	}

	result.Success = true
	if r.Verbose {
		logger.Info("Migration %d completed in %dms", m.Number, result.ExecutionTimeMs)
	}

	return result, nil
}

// ApplyAll applies all pending migrations
func (r *MigrationRunner) ApplyAll() ([]*MigrationResult, error) {
	pending, err := r.GetPendingMigrations()
	if err != nil {
		return nil, err
	}

	if len(pending) == 0 {
		logger.Info("No pending migrations")
		return nil, nil
	}

	logger.Info("Found %d pending migrations", len(pending))

	var results []*MigrationResult
	for _, migration := range pending {
		logger.Info("Applying migration %d: %s", migration.Number, migration.Description)

		result, err := r.ApplyMigration(migration)
		results = append(results, result)

		if err != nil {
			if migration.Critical {
				logger.Error("Critical migration failed, stopping: %v", err)
				return results, err
			}
			logger.Warn("Non-critical migration failed: %v", err)
		}
	}

	logger.Info("All pending migrations completed")
	return results, nil
}

// Rollback rolls back the last N migrations
func (r *MigrationRunner) Rollback(steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be positive")
	}

	// Get migration history
	var history []struct {
		MigrationNumber int
		Version         string
		Description     string
	}

	err := r.DB.Raw(`
		SELECT migration_number, version, description
		FROM migration_history
		WHERE success = TRUE
		ORDER BY migration_number DESC
		LIMIT ?
	`, steps).Scan(&history).Error

	if err != nil {
		// Check if table doesn't exist - means no migrations applied
		if err == gorm.ErrRecordNotFound || err.Error() == "no such table: migration_history" {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get migration history: %w", err)
	}

	if len(history) == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	logger.Info("Rolling back %d migration(s)", len(history))

	// Rollback in reverse order
	for _, h := range history {
		migration := GetMigrationByNumber(h.MigrationNumber)
		if migration == nil {
			return fmt.Errorf("migration %d not found", h.MigrationNumber)
		}

		if len(migration.DownSQL) == 0 {
			return fmt.Errorf("migration %d has no rollback SQL", h.MigrationNumber)
		}

		logger.Info("Rolling back migration %d: %s", migration.Number, migration.Description)

		err := r.DB.Transaction(func(tx *gorm.DB) error {
			// Remove from migration history FIRST (before dropping tables)
			if err := tx.Exec("DELETE FROM migration_history WHERE migration_number = ?", migration.Number).Error; err != nil {
				return fmt.Errorf("failed to remove from history: %w", err)
			}

			// Update schema_version (remove this version entry)
			if err := tx.Exec(`
				DELETE FROM schema_version 
				WHERE migration_number = ?
			`, migration.Number).Error; err != nil {
				return fmt.Errorf("failed to update version: %w", err)
			}

			// Execute rollback SQL (drop tables)
			for _, sql := range migration.DownSQL {
				if r.Verbose {
					logger.Debug("Executing rollback SQL: %s", sql)
				}

				if err := tx.Exec(sql).Error; err != nil {
					return fmt.Errorf("rollback SQL failed: %s - Error: %w", sql, err)
				}
			}

			return nil
		})

		if err != nil {
			return err
		}

		logger.Info("Migration %d rolled back successfully", migration.Number)
	}

	logger.Info("Rollback completed")
	return nil
}
