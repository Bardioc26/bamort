package gsmaster

import (
	"bamort/importer"
	"fmt"
)

func CheckSkill(fertigkeit *importer.Fertigkeit, autocreate bool) (*LookupSkill, error) {
	stammF := LookupSkill{}
	//err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, fertigkeit.Name).Error
	err := stammF.First(fertigkeit.Name)
	if err == nil {
		// Fertigkeit found
		return &stammF, nil
	}
	if !autocreate {
		return nil, fmt.Errorf("does not exist in Fertigkeit importer")
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

func CheckSpell(zauber *importer.Zauber, autocreate bool) (*LookupSpell, error) {
	stammF := LookupSpell{}

	//err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, zauber.Name).Error
	err := stammF.First(zauber.Name)
	if err == nil {
		// zauber found
		return &stammF, nil
	}
	if !autocreate {
		return nil, fmt.Errorf("does not exist in zauber importer")
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
