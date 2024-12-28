package main

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func CheckFertigkeit(fertigkeit *Fertigkeit, autocreate bool) (*StammFertigkeit, error) {
	stammF := StammFertigkeit{}

	if strings.HasPrefix(fertigkeit.ImportID, "moam") {
		err := DB.First(&stammF, "system=? AND name = ?", "midgard", fertigkeit.Name).Error
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
		err = DB.Transaction(func(tx *gorm.DB) error {
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
	err := DB.First(&stammF, "system=? AND name = ?", "midgard", fertigkeit.Name).Error
	if err != nil {
		// Fertigkeit found
		return nil, err
	}
	return &stammF, nil
}
