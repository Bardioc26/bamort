package models

import (
	"bamort/database"
	"fmt"

	"gorm.io/gorm"
)

type LookupList struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	System       string `gorm:"index" json:"system"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Quelle       string `json:"quelle"`
}

type LookupSkill struct {
	LookupList
	Initialkeitswert int    `json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
}

type LookupWaeponSkill struct {
	LookupSkill
}

type LookupSpell struct {
	LookupList
	Bonus        int `json:"bonus"`
	Stufe        int
	AP           int
	Reichweite   int
	Wirkungsziel string
}

type LookupEquipment struct {
	LookupList
	Gewicht float64 `json:"gewicht"`
	Wert    float64 `json:"wert"`
}

type LookupContainer struct {
	LookupEquipment
	Tragkraft float64 `json:"tragkraft"`
	Volumen   float64 `json:"volumen"`
}

type LookupTransportation struct {
	LookupContainer
}

func (stamm *LookupList) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *LookupList) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save Lookup: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *LookupSkill) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (stamm *LookupSkill) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupSkill: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *LookupWaeponSkill) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (stamm *LookupWaeponSkill) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupWaeponSkill: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *LookupSpell) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *LookupSpell) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupSpell: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *LookupEquipment) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *LookupEquipment) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupEquipment: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *LookupContainer) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *LookupContainer) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupContainer: %w", err)
		}
		return nil
	})

	return err
}
func (stamm *LookupTransportation) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *LookupTransportation) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save Lookup: %w", err)
		}
		return nil
	})

	return err
}
