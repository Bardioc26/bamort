package models

import (
	"bamort/database"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func CheckFertigkeit(fertigkeit *ImFertigkeit, autocreate bool) (*ImStammFertigkeit, error) {
	stammF := ImStammFertigkeit{}

	if strings.HasPrefix(fertigkeit.ID, "moam") {
		err := database.DB.First(&stammF, "system=? AND name = ?", "midgard", fertigkeit.Name).Error
		if err == nil {
			// Fertigkeit found
			return &stammF, nil
		}
		if !autocreate {
			return nil, fmt.Errorf("does not exist in Fertigkeit Stammdaten")
		}
		stammF.System = "midgard"
		stammF.Name = fertigkeit.Name
		stammF.Beschreibung = fertigkeit.Beschreibung
		if fertigkeit.Fertigkeitswert < 12 {
			stammF.Initialkeitswert = 5
		} else {
			stammF.Initialkeitswert = 12
		}
		stammF.Bonuseigenschaft = "keine"
		stammF.Quelle = fertigkeit.Quelle
		//fmt.Println(stammF)
		err = database.DB.Transaction(func(tx *gorm.DB) error {
			// Save the main character record
			if err := tx.Create(&stammF).Error; err != nil {
				return fmt.Errorf("failed to save Fertigkeit Stammdaten: %w", err)
			}
			return nil
		})
		if err != nil {
			// Fertigkeit found
			return nil, err
		}

	}
	err := database.DB.First(&stammF, "system=? AND name = ?", "midgard", fertigkeit.Name).Error
	if err != nil {
		// Fertigkeit found
		return nil, err
	}
	return &stammF, nil
}

type ImStammFertigkeit struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	System           string `gorm:"index" json:"system"`
	Name             string `json:"name"`
	Beschreibung     string `json:"beschreibung"`
	Initialkeitswert int    `json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
	Quelle           string `json:"quelle"`
}
