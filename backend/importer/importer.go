package importer

import (
	"bamort/character"
	"bamort/gsmaster"
	"bamort/models"
	"fmt"
)

func ImportChar(char CharacterImport) (*character.Char, error) {
	return nil, fmt.Errorf("char could not be imported %s", "Weil Wegen Kommt noch")
}

func CheckSkill(fertigkeit *Fertigkeit, autocreate bool) (*models.Skill, error) {
	stammF := models.Skill{}
	//err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, fertigkeit.Name).Error
	err := stammF.First(fertigkeit.Name)
	if err == nil {
		// Fertigkeit found
		return &stammF, nil
	}
	if !autocreate {
		return nil, fmt.Errorf("does not exist in Fertigkeit importer")
	}
	stammF.GameSystem = "midgard"
	stammF.Name = fertigkeit.Name
	if stammF.Name != "Sprache" {
		stammF.Beschreibung = fertigkeit.Beschreibung
	}
	if fertigkeit.Fertigkeitswert < 12 {
		stammF.Initialwert = 5
	} else {
		stammF.Initialwert = 12
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

func CheckSpell(zauber *Zauber, autocreate bool) (*gsmaster.Spell, error) {
	stammF := gsmaster.Spell{}

	//err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, zauber.Name).Error
	err := stammF.First(zauber.Name)
	if err == nil {
		// zauber found
		return &stammF, nil
	}
	if !autocreate {
		return nil, fmt.Errorf("does not exist in zauber importer")
	}
	stammF.GameSystem = "midgard"
	stammF.Name = zauber.Name
	stammF.Beschreibung = zauber.Beschreibung
	stammF.AP = "1"
	stammF.Stufe = 1
	stammF.Wirkungsziel = "Zauberer"
	stammF.Reichweite = "15 m"

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
