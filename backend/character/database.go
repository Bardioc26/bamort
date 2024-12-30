package character

import (
	"bamort/database"
	"bamort/models"
	"database/sql/driver"
	"encoding/json"
	"errors"

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

		/*
			// Save Eigenschaften (Attributes)
			for i := range character.Eigenschaften {
				character.Eigenschaften[i].CharacterID = character.ID
			}
			if len(character.Eigenschaften) > 0 {
				if err := tx.Create(&character.Eigenschaften).Error; err != nil {
					return fmt.Errorf("failed to save eigenschaften: %w", err)
				}
			}

			// Save Ausruestung (Equipment)
			for i := range character.Ausruestung {
				character.Ausruestung[i].CharacterID = character.ID
			}
			if len(character.Ausruestung) > 0 {
				if err := tx.Create(&character.Ausruestung).Error; err != nil {
					return fmt.Errorf("failed to save ausruestung: %w", err)
				}
			}

			// Save Behaeltnisse (Containers)
			for i := range character.Behaeltnisse {
				character.Behaeltnisse[i].CharacterID = character.ID
			}
			if len(character.Behaeltnisse) > 0 {
				if err := tx.Create(&character.Behaeltnisse).Error; err != nil {
					return fmt.Errorf("failed to save behaeltnisse: %w", err)
				}
			}

			// Save Fertigkeiten (Skills)
			for i := range character.Fertigkeiten {
				character.Fertigkeiten[i].CharacterID = character.ID
			}
			if len(character.Fertigkeiten) > 0 {
				if err := tx.Create(&character.Fertigkeiten).Error; err != nil {
					return fmt.Errorf("failed to save fertigkeiten: %w", err)
				}
			}

			// Save Waffenfertigkeiten (Weapon Skills)
			for i := range character.Waffenfertigkeiten {
				character.Waffenfertigkeiten[i].CharacterID = character.ID
			}
			if len(character.Waffenfertigkeiten) > 0 {
				if err := tx.Create(&character.Waffenfertigkeiten).Error; err != nil {
					return fmt.Errorf("failed to save waffenfertigkeiten: %w", err)
				}
			}

			// Save Zauber (Spells)
			for i := range character.Zauber {
				character.Zauber[i].CharacterID = character.ID
			}
			if len(character.Zauber) > 0 {
				if err := tx.Create(&character.Zauber).Error; err != nil {
					return fmt.Errorf("failed to save zauber: %w", err)
				}
			}

			// Save Waffen (Weapons)
			for i := range character.Waffen {
				character.Waffen[i].CharacterID = character.ID
			}
			if len(character.Waffen) > 0 {
				if err := tx.Create(&character.Waffen).Error; err != nil {
					return fmt.Errorf("failed to save waffen: %w", err)
				}
			}

			// Save Merkmale (Characteristics)
			character.Merkmale.CharacterID = character.ID
			if err := tx.Create(&character.Merkmale).Error; err != nil {
				return fmt.Errorf("failed to save merkmale: %w", err)
			}

			// Save Bennies
			character.Bennies.CharacterID = character.ID
			if err := tx.Create(&character.Bennies).Error; err != nil {
				return fmt.Errorf("failed to save bennies: %w", err)
			}

			// Save Gestalt (Appearance)
			character.Gestalt.CharacterID = character.ID
			if err := tx.Create(&character.Gestalt).Error; err != nil {
				return fmt.Errorf("failed to save gestalt: %w", err)
			}

			// Save Lp (Life Points)
			character.Lp.CharacterID = character.ID
			if err := tx.Create(&character.Lp).Error; err != nil {
				return fmt.Errorf("failed to save lp: %w", err)
			}

			// Save Ap (Action Points)
			character.Ap.CharacterID = character.ID
			if err := tx.Create(&character.Ap).Error; err != nil {
				return fmt.Errorf("failed to save ap: %w", err)
			}

			// Save B (Other Points)
			character.B.CharacterID = character.ID
			if err := tx.Create(&character.B).Error; err != nil {
				return fmt.Errorf("failed to save b: %w", err)
			}

			// Save Transportmittel (Transportation)
			for i := range character.Transportmittel {
				character.Transportmittel[i].CharacterID = character.ID
			}
			if len(character.Transportmittel) > 0 {
				if err := tx.Create(&character.Transportmittel).Error; err != nil {
					return fmt.Errorf("failed to save transportmittel: %w", err)
				}
			}

			// Save Erfahrungsschatz (Experience)
			character.Erfahrungsschatz.CharacterID = character.ID
			if err := tx.Create(&character.Erfahrungsschatz).Error; err != nil {
				return fmt.Errorf("failed to save erfahrungsschatz: %w", err)
			}
		*/

		return nil
	})
}

// StringArray is a custom type for []string
type StringArray []string

// Value implements the driver.Valuer interface for database storage
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s) // Serialize []string to JSON
}

// Scan implements the sql.Scanner interface for database retrieval
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to convert database value to []byte")
	}

	return json.Unmarshal(bytes, s) // Deserialize JSON to []string
}
