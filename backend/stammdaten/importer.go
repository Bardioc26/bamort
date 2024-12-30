package stammdaten

import (
	"bamort/models"
	"fmt"
)

func CheckSkill(fertigkeit *models.ImFertigkeit, autocreate bool) (*models.LookupSkill, error) {
	stammF := models.LookupSkill{}
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

func CheckSpell(zauber *models.ImZauber, autocreate bool) (*models.LookupSpell, error) {
	stammF := models.LookupSpell{}

	//err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, zauber.Name).Error
	err := stammF.First(zauber.Name)
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
	err = stammF.Create()
	if err != nil {
		// spell found
		return nil, err
	}

	//err = database.DB.First(&stammF, "system=? AND name = ?", gameSystem, zauber.Name).Error
	err = stammF.First(zauber.Name)
	if err != nil {
		// spell found
		return nil, err
	}
	return &stammF, nil
}
