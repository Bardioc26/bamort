package gamesystem

import (
	"bamort/database"
	"bamort/models"

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

	err := targetDB.AutoMigrate(
		&models.GameSystem{},
	)
	if err != nil {
		return err
	}
	return nil
}

func MigrateDataIfNeeded(db *gorm.DB) error {
	// Implement data migration logic here if needed
	gameSystem := models.GameSystem{}
	err := db.First(&gameSystem, "ID = ?", 1).Error
	if err != nil {
		// No initial version found, assume no migration needed
		gameSystem.Code = "M5"
		gameSystem.Name = "M-System"
		gameSystem.Description = "Version 5 des Rollenspiels"
		gameSystem.IsActive = true
		err = db.Create(&gameSystem).Error
		if err != nil {
			return err
		}
	}
	return nil
}
