package models

import (
	"gorm.io/gorm"
)

// importerMigrateStructure is a placeholder - actual importer migrations
// are now handled by the importer package itself via importer.MigrateStructure()
// which is called from cmd/main.go after models.MigrateStructure()
func importerMigrateStructure(db ...*gorm.DB) error {
	// No longer needed - migrations handled by importer package
	return nil
}
