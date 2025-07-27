package models

import (
	"bamort/database"
	"fmt"

	"gorm.io/gorm"
)

type BamortBase struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type BamortCharTrait struct {
	BamortBase
	CharacterID uint `gorm:"index" json:"character_id"`
}

type Magisch struct {
	IstMagisch  bool `json:"ist_magisch"`
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
}

type LookupList struct {
	ID           uint   `gorm:"primaryKey" json:"id"` //`gorm:"default:uuid_generate_v3()"` // db func
	GameSystem   string `gorm:"column:game_system;index;default:midgard" json:"game_system"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Quelle       string `json:"quelle"`
}

func (object *LookupList) Create() error {
	gameSystem := "midgard"
	object.GameSystem = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&object).Error; err != nil {
			return fmt.Errorf("failed to save Lookup: %w", err)
		}
		return nil
	})

	return err
}

func (object *LookupList) First(value string) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND name!='Placeholder' AND name = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *LookupList) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND name!='Placeholder' AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *LookupList) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}
