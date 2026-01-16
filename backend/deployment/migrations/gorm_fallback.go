package migrations

import (
	"bamort/logger"
	"bamort/models"
	"fmt"
)

// RunGORMAutoMigrate runs GORM's AutoMigrate as a safety net after SQL migrations
// This catches any columns or tables that might have been missed in SQL migrations
func (r *MigrationRunner) RunGORMAutoMigrate() error {
	logger.Info("Running GORM AutoMigrate as safety net...")

	// Run models.MigrateStructure() after SQL migrations
	// This catches any columns we might have missed
	if err := models.MigrateStructure(r.DB); err != nil {
		return fmt.Errorf("GORM AutoMigrate failed: %w", err)
	}

	logger.Info("GORM AutoMigrate completed successfully")
	return nil
}
