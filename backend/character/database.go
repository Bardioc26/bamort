package character

import (
	"bamort/database"
	"bamort/models"

	"fmt"

	"gorm.io/gorm"
)

func SaveCharacterToDB(character *models.Char) error {
	// Use GORM transaction to ensure atomicity
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(character).Error; err != nil {
			return fmt.Errorf("failed to save character: %w", err)
		}
		return nil
	})
}
