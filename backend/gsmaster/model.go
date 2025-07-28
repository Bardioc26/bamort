package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"fmt"

	"gorm.io/gorm"
)

var dbPrefix = "gsm"

type Believe struct {
	models.LookupList
}

func (object *Believe) TableName() string {
	return dbPrefix + "_" + "believes"
}

func (stamm *Believe) Create() error {
	gameSystem := "midgard"
	stamm.GameSystem = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save Lookup: %w", err)
		}
		return nil
	})

	return err
}

func (object *Believe) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *Believe) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "game_system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Believe) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}
