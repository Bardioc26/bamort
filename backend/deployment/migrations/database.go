package migrations

import (
	"gorm.io/gorm"
)

// MigrateStructure migrates all deployment-related structures to the database
func MigrateStructure(db *gorm.DB) error {
	// Migrate deployment package structures (schema_version and migration_history tables)
	return db.AutoMigrate(
		&SchemaVersion{},
		&MigrationHistory{},
	)
}
