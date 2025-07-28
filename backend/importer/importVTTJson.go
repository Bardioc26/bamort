package importer

import (
	"bamort/models"
	"encoding/json"
	"os"
)

func readImportChar(fileName string) (*CharacterImport, error) {
	// loading file to Modell
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	character := CharacterImport{}
	err = json.Unmarshal(fileContent, &character)
	return &character, err
}

func ImportVTTJSON(fileName string) (*models.Char, error) {
	//fileName = fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	imp, err := readImportChar(fileName)
	if err != nil {
		return nil, err
	}
	err = CheckFertigkeiten2GSMaster(imp.Fertigkeiten)
	if err != nil {
		return nil, err
	}
	err = CheckWaffenFertigkeiten2GSMaster(imp.Waffenfertigkeiten)
	if err != nil {
		return nil, err
	}
	err = CheckSpells2GSMaster(imp.Zauber)
	if err != nil {
		return nil, err
	}
	err = CheckWeapons2GSMaster(imp.Waffen)
	if err != nil {
		return nil, err
	}
	err = CheckContainers2GSMaster(imp.Behaeltnisse)
	if err != nil {
		return nil, err
	}
	err = CheckTransportations2GSMaster(imp.Transportmittel)
	if err != nil {
		return nil, err
	}
	err = CheckEquipments2GSMaster(imp.Ausruestung)
	if err != nil {
		return nil, err
	}
	err = CheckBelieve2GSMaster(imp)
	if err != nil {
		return nil, err
	}

	char := models.Char{}
	char.Name = imp.Name
	char.Rasse = imp.Rasse
	char.Typ = imp.Typ
	char.Alter = imp.Alter
	char.Anrede = imp.Anrede
	char.Grad = imp.Grad
	char.Groesse = imp.Groesse
	char.Gewicht = imp.Gewicht
	char.Glaube = imp.Glaube
	char.Hand = imp.Hand
	char.Image = imp.Image
	for i := range imp.Fertigkeiten {
		char.Fertigkeiten = append(char.Fertigkeiten, models.SkFertigkeit{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: imp.Fertigkeiten[i].Name},
			},
			Beschreibung:    imp.Fertigkeiten[i].Beschreibung,
			Fertigkeitswert: imp.Fertigkeiten[i].Fertigkeitswert,
			Bonus:           imp.Fertigkeiten[i].Bonus,
			Bemerkung:       imp.Fertigkeiten[i].Beschreibung,
		})
	}
	for i := range imp.Zauber {
		char.Zauber = append(char.Zauber, models.SkZauber{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: imp.Zauber[i].Name},
			},
			Beschreibung: imp.Zauber[i].Beschreibung,
			Bonus:        imp.Zauber[i].Bonus,
			Quelle:       imp.Zauber[i].Quelle,
		})
	}
	char.Lp.Max = imp.Lp.Max
	char.Lp.Value = imp.Lp.Value
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "Au",
		Value: imp.Eigenschaften.Au,
	})
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "Gs",
		Value: imp.Eigenschaften.Gs,
	})
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "Gw",
		Value: imp.Eigenschaften.Gw,
	})
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "In",
		Value: imp.Eigenschaften.In,
	})
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "Ko",
		Value: imp.Eigenschaften.Ko,
	})
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "PA",
		Value: imp.Eigenschaften.Pa,
	})
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "St",
		Value: imp.Eigenschaften.St,
	})
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "Wk",
		Value: imp.Eigenschaften.Wk,
	})
	char.Eigenschaften = append(char.Eigenschaften, models.Eigenschaft{
		Name:  "Zt",
		Value: imp.Eigenschaften.Zt,
	})
	char.Merkmale.Augenfarbe = imp.Merkmale.Augenfarbe
	char.Merkmale.Haarfarbe = imp.Merkmale.Haarfarbe
	char.Merkmale.Sonstige = imp.Merkmale.Sonstige
	char.Bennies = models.Bennies{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{Name: "bennies"},
		},
		Gg: imp.Bennies.Gg,
		Sg: imp.Bennies.Sg,
		Gp: imp.Bennies.Gp,
	}
	char.Merkmale.Breite = imp.Gestalt.Breite
	char.Merkmale.Groesse = imp.Gestalt.Groesse
	char.Ap.Max = imp.Ap.Max
	char.Ap.Value = imp.Ap.Value
	char.B.Max = imp.B.Max
	char.B.Value = imp.B.Value
	char.Erfahrungsschatz.ES = imp.Erfahrungsschatz.Value
	for i := range imp.Transportmittel {
		char.Transportmittel = append(char.Transportmittel, models.EqContainer{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: imp.Transportmittel[i].Name},
			},
			IsTransportation: true,
			Beschreibung:     imp.Transportmittel[i].Beschreibung,
			BeinhaltetIn:     imp.Transportmittel[i].BeinhaltetIn,
			Gewicht:          float64(imp.Transportmittel[i].Gewicht),
			Tragkraft:        imp.Transportmittel[i].Tragkraft,
			Wert:             imp.Transportmittel[i].Wert,
			Magisch: models.Magisch{
				IstMagisch:  imp.Transportmittel[i].Magisch.IstMagisch,
				Abw:         imp.Transportmittel[i].Magisch.Abw,
				Ausgebrannt: imp.Transportmittel[i].Magisch.Ausgebrannt,
			},
			ExtID: imp.Transportmittel[i].ID,
		})
	}
	for i := range imp.Ausruestung {
		char.Ausruestung = append(char.Ausruestung, models.EqAusruestung{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: imp.Ausruestung[i].Name},
			},
			Beschreibung: imp.Ausruestung[i].Beschreibung,
			Anzahl:       imp.Ausruestung[i].Anzahl,
			BeinhaltetIn: imp.Ausruestung[i].BeinhaltetIn,
			Bonus:        imp.Ausruestung[i].Bonus,
			Gewicht:      float64(imp.Ausruestung[i].Gewicht),
			Wert:         imp.Ausruestung[i].Wert,
			Magisch: models.Magisch{
				IstMagisch:  imp.Ausruestung[i].Magisch.IstMagisch,
				Abw:         imp.Ausruestung[i].Magisch.Abw,
				Ausgebrannt: imp.Ausruestung[i].Magisch.Ausgebrannt,
			},
		})
	}
	for i := range imp.Behaeltnisse {
		char.Behaeltnisse = append(char.Behaeltnisse, models.EqContainer{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: imp.Behaeltnisse[i].Name},
			},
			IsTransportation: false,
			Beschreibung:     imp.Behaeltnisse[i].Beschreibung,
			BeinhaltetIn:     imp.Behaeltnisse[i].BeinhaltetIn,
			Tragkraft:        imp.Behaeltnisse[i].Tragkraft,
			Volumen:          imp.Behaeltnisse[i].Volumen,
			Gewicht:          float64(imp.Behaeltnisse[i].Gewicht),
			Wert:             imp.Behaeltnisse[i].Wert,
			Magisch: models.Magisch{
				IstMagisch:  imp.Behaeltnisse[i].Magisch.IstMagisch,
				Abw:         imp.Behaeltnisse[i].Magisch.Abw,
				Ausgebrannt: imp.Behaeltnisse[i].Magisch.Ausgebrannt,
			},
			ExtID: imp.Behaeltnisse[i].ID,
		})
	}
	for i := range imp.Waffen {
		char.Waffen = append(char.Waffen, models.EqWaffe{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: imp.Waffen[i].Name},
			},
			Beschreibung:            imp.Waffen[i].Beschreibung,
			Abwb:                    imp.Waffen[i].Abwb,
			Anb:                     imp.Waffen[i].Anb,
			Anzahl:                  imp.Waffen[i].Anzahl,
			Schb:                    imp.Waffen[i].Schb,
			BeinhaltetIn:            imp.Waffen[i].BeinhaltetIn,
			Gewicht:                 float64(imp.Waffen[i].Gewicht),
			Wert:                    imp.Waffen[i].Wert,
			NameFuerSpezialisierung: imp.Waffen[i].NameFuerSpezialisierung,
			Magisch: models.Magisch{
				IstMagisch:  imp.Waffen[i].Magisch.IstMagisch,
				Abw:         imp.Waffen[i].Magisch.Abw,
				Ausgebrannt: imp.Waffen[i].Magisch.Ausgebrannt,
			},
		})

	}
	for i := range imp.Waffenfertigkeiten {
		char.Waffenfertigkeiten = append(char.Waffenfertigkeiten, models.SkWaffenfertigkeit{
			SkFertigkeit: models.SkFertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{Name: imp.Waffenfertigkeiten[i].Name},
				},
				Beschreibung:    imp.Waffenfertigkeiten[i].Beschreibung,
				Fertigkeitswert: imp.Waffenfertigkeiten[i].Fertigkeitswert,
				Bonus:           imp.Waffenfertigkeiten[i].Bonus,
				//Bemerkung: imp.Waffenfertigkeiten[i],
			},
		})
	}
	for i := range imp.Spezialisierung {
		char.Spezialisierung = append(char.Spezialisierung, imp.Spezialisierung[i])
	}
	err = char.Create()
	if err != nil {
		return nil, err
	}
	// Fix contained in links
	for i := range char.Ausruestung {
		err := char.Ausruestung[i].LinkContainer()
		if err != nil {
			return &char, err
		}
	}
	for i := range char.Waffen {
		err := char.Waffen[i].LinkContainer()
		if err != nil {
			return &char, err
		}
	}
	for i := range char.Behaeltnisse {
		err := char.Behaeltnisse[i].LinkContainer()
		if err != nil {
			return &char, err
		}
	}
	for i := range char.Transportmittel {
		err := char.Transportmittel[i].LinkContainer()
		if err != nil {
			return &char, err
		}
	}
	return &char, nil
}
