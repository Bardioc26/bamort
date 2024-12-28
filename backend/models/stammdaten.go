package models

import (
	"bamort/database"
	"fmt"

	"gorm.io/gorm"
)

/*type ImStammFertigkeit struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	System           string `gorm:"index" json:"system"`
	Name             string `json:"name"`
	Beschreibung     string `json:"beschreibung"`
	Initialkeitswert int    `json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
	Quelle           string `json:"quelle"`
}
*/

type ImStamm struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	System       string `gorm:"index" json:"system"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Quelle       string `json:"quelle"`
}

type ImStammFertigkeit struct {
	ImStamm
	Initialkeitswert int    `json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
}
type ImStammWaffenFertigkeit struct {
	ImStammFertigkeit
}

type ImStammZauber struct {
	ImStamm
	Bonus        int `json:"bonus"`
	Stufe        int
	AP           int
	Reichweite   int
	Wirkungsziel string
}

func (stamm *ImStamm) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *ImStamm) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save Stammdaten: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *ImStammFertigkeit) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (stamm *ImStammFertigkeit) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save Fertigkeit Stammdaten: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *ImStammWaffenFertigkeit) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (stamm *ImStammWaffenFertigkeit) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save Fertigkeit Stammdaten: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *ImStammZauber) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *ImStammZauber) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save Zauber Stammdaten: %w", err)
		}
		return nil
	})

	return err
}
