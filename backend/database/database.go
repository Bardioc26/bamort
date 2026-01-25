package database

import "gorm.io/gorm"

func MigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = DB
	}

	err := targetDB.AutoMigrate(
		&SchemaVersion{},
		&MigrationHistory{},
	)
	if err != nil {
		return err
	}
	return nil
}

func MigrateDataIfNeeded(db *gorm.DB) error {
	// Implement data migration logic here if needed
	schemaVersion := SchemaVersion{}
	err := db.First(&schemaVersion, "version = ?", "0.1.37").Error
	if err != nil {
		// No initial version found, assume no migration needed
		schemaVersion.Version = "0.1.37"
		schemaVersion.MigrationNumber = 1
		schemaVersion.BackendVersion = "0.1.37"
		schemaVersion.Description = "Initial schema version"
		err = db.Create(&schemaVersion).Error
		if err != nil {
			return err
		}
	}
	return nil
}
