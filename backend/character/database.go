package character

import (
	"bamort/database"

	"fmt"

	"gorm.io/gorm"
)

func SaveCharacterToDB(character *Char) error {
	// Use GORM transaction to ensure atomicity
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(character).Error; err != nil {
			return fmt.Errorf("failed to save character: %w", err)
		}
		return nil
	})
}

func MigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := targetDB.AutoMigrate(
		&Char{},
		&Eigenschaft{},
		&Lp{},
		&Ap{},
		&B{},
		&Merkmale{},
		&Erfahrungsschatz{},
		&Bennies{},
		&Vermoegen{},
	)
	if err != nil {
		return err
	}
	return nil
}
