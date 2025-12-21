package user

import (
	"bamort/database"
	"fmt"

	"gorm.io/gorm"
)

func MigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	// Check if we have a valid database connection
	if targetDB == nil {
		return fmt.Errorf("no database connection available for migration")
	}

	err := targetDB.AutoMigrate(
		&User{},
	)
	if err != nil {
		return err
	}
	return nil
}
