package models

import (
	"bamort/database"

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
		&Skill{},
		&WeaponSkill{},
		&Spell{},
		&Equipment{},
		&Weapon{},
		&Container{},
		&Transportation{},
		&Believe{},
	)
	if err != nil {
		return err
	}
	return nil
}
