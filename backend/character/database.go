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
func MigrateStructure() error {
	err := database.DB.AutoMigrate(
		&Char{},
		&Eigenschaft{},
		&Lp{},
		&Ap{},
		&B{},
		&Merkmale{},
		&Erfahrungsschatz{},
		&Bennies{},
	)
	if err != nil {
		return err
	}
	return nil
}
