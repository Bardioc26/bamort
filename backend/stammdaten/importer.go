package stammdaten

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func CheckFertigkeit(fertigkeit *models.ImFertigkeit, autocreate bool) (*models.ImStammFertigkeit, error) {
	stammF := models.ImStammFertigkeit{}
	//err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, fertigkeit.Name).Error
	err := stammF.First(fertigkeit.Name)
	if err == nil {
		// Fertigkeit found
		return &stammF, nil
	}
	if !autocreate {
		return nil, fmt.Errorf("does not exist in Fertigkeit Stammdaten")
	}
	stammF.System = "midgard"
	stammF.Name = fertigkeit.Name
	if stammF.Name != "Sprache" {
		stammF.Beschreibung = fertigkeit.Beschreibung
	}
	if fertigkeit.Fertigkeitswert < 12 {
		stammF.Initialkeitswert = 5
	} else {
		stammF.Initialkeitswert = 12
	}
	stammF.Bonuseigenschaft = "keine"
	stammF.Quelle = fertigkeit.Quelle
	//fmt.Println(stammF)

	err = stammF.Create()
	if err != nil {
		// Fertigkeit found
		return nil, err
	}

	//err = database.DB.First(&stammF, "system=? AND name = ?", gameSystem, fertigkeit.Name).Error
	err = stammF.First(fertigkeit.Name)
	if err != nil {
		// Fertigkeit found
		return nil, err
	}
	return &stammF, nil
}

func CheckZauber(zauber *models.ImZauber, autocreate bool) (*models.ImStammZauber, error) {
	stammF := models.ImStammZauber{}
	gameSystem := "none"
	if strings.HasPrefix(zauber.ID, "moam") {
		gameSystem = "midgard"
	}
	err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, zauber.Name).Error
	if err == nil {
		// zauber found
		return &stammF, nil
	}
	if !autocreate {
		return nil, fmt.Errorf("does not exist in zauber Stammdaten")
	}
	stammF.System = "midgard"
	stammF.Name = zauber.Name
	stammF.Beschreibung = zauber.Beschreibung
	stammF.AP = 1
	stammF.Stufe = 1
	stammF.Wirkungsziel = "Zauberer"
	stammF.Reichweite = 15

	stammF.Quelle = zauber.Quelle
	//fmt.Println(stammF)
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stammF).Error; err != nil {
			return fmt.Errorf("failed to save zauber Stammdaten: %w", err)
		}
		return nil
	})
	if err != nil {
		// zauber found
		return nil, err
	}

	err = database.DB.First(&stammF, "system=? AND name = ?", gameSystem, zauber.Name).Error
	if err != nil {
		// zauber found
		return nil, err
	}
	return &stammF, nil
}
