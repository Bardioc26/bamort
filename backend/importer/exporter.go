package importer

import (
	"bamort/models"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ExportCharToVTT converts a BaMoRT character to VTT JSON format
func ExportCharToVTT(char *models.Char) (*CharacterImport, error) {
	vtt := &CharacterImport{}

	// Basic info
	vtt.ID = fmt.Sprintf("bamort-character-%d", char.ID)
	vtt.Name = char.Name
	vtt.Rasse = char.Rasse
	vtt.Typ = char.Typ
	vtt.Alter = char.Alter
	vtt.Anrede = char.Anrede
	vtt.Grad = char.Grad
	vtt.Groesse = char.Groesse
	vtt.Gewicht = char.Gewicht
	vtt.Glaube = char.Glaube
	vtt.Hand = char.Hand
	vtt.Image = char.Image

	// LP
	vtt.Lp.Max = char.Lp.Max
	vtt.Lp.Value = char.Lp.Value

	// AP
	vtt.Ap.Max = char.Ap.Max
	vtt.Ap.Value = char.Ap.Value

	// B
	vtt.B.Max = char.B.Max
	vtt.B.Value = char.B.Value

	// Eigenschaften - convert from array to struct
	for _, e := range char.Eigenschaften {
		switch e.Name {
		case "Au":
			vtt.Eigenschaften.Au = e.Value
		case "Gs":
			vtt.Eigenschaften.Gs = e.Value
		case "Gw":
			vtt.Eigenschaften.Gw = e.Value
		case "In":
			vtt.Eigenschaften.In = e.Value
		case "Ko":
			vtt.Eigenschaften.Ko = e.Value
		case "PA":
			vtt.Eigenschaften.Pa = e.Value
		case "St":
			vtt.Eigenschaften.St = e.Value
		case "Wk":
			vtt.Eigenschaften.Wk = e.Value
		case "Zt":
			vtt.Eigenschaften.Zt = e.Value
		}
	}

	// Merkmale
	vtt.Merkmale.Augenfarbe = char.Merkmale.Augenfarbe
	vtt.Merkmale.Haarfarbe = char.Merkmale.Haarfarbe
	vtt.Merkmale.Sonstige = char.Merkmale.Sonstige

	// Gestalt
	vtt.Gestalt.Breite = char.Merkmale.Breite
	vtt.Gestalt.Groesse = char.Merkmale.Groesse

	// Bennies
	vtt.Bennies.Gg = char.Bennies.Gg
	vtt.Bennies.Sg = char.Bennies.Sg
	vtt.Bennies.Gp = char.Bennies.Gp

	// Erfahrungsschatz
	vtt.Erfahrungsschatz.Value = char.Erfahrungsschatz.ES

	// Spezialisierung
	vtt.Spezialisierung = char.Spezialisierung

	// Fertigkeiten
	for i, f := range char.Fertigkeiten {
		vtt.Fertigkeiten = append(vtt.Fertigkeiten, Fertigkeit{
			ImportBase: ImportBase{
				ID:   fmt.Sprintf("bamort-skill-%d-%d", char.ID, i),
				Name: f.Name,
			},
			Beschreibung:    f.Beschreibung,
			Fertigkeitswert: f.Fertigkeitswert,
			Bonus:           f.Bonus,
			Pp:              f.Pp,
			Quelle:          "", // Not stored in character skill
		})
	}

	// Waffenfertigkeiten
	for i, w := range char.Waffenfertigkeiten {
		vtt.Waffenfertigkeiten = append(vtt.Waffenfertigkeiten, Waffenfertigkeit{
			ImportBase: ImportBase{
				ID:   fmt.Sprintf("bamort-weaponskill-%d-%d", char.ID, i),
				Name: w.Name,
			},
			Beschreibung:    w.Beschreibung,
			Fertigkeitswert: w.Fertigkeitswert,
			Bonus:           w.Bonus,
			Pp:              w.Pp,
			Quelle:          "", // Not stored in character skill
		})
	}

	// Zauber
	for i, z := range char.Zauber {
		vtt.Zauber = append(vtt.Zauber, Zauber{
			ImportBase: ImportBase{
				ID:   fmt.Sprintf("bamort-spell-%d-%d", char.ID, i),
				Name: z.Name,
			},
			Beschreibung: z.Beschreibung,
			Bonus:        z.Bonus,
			Quelle:       z.Quelle,
		})
	}

	// Waffen
	for i, w := range char.Waffen {
		vtt.Waffen = append(vtt.Waffen, Waffe{
			ImportBase: ImportBase{
				ID:   fmt.Sprintf("bamort-weapon-%d-%d", char.ID, i),
				Name: w.Name,
			},
			Beschreibung:            w.Beschreibung,
			Gewicht:                 w.Gewicht,
			Wert:                    w.Wert,
			Anzahl:                  w.Anzahl,
			Anb:                     w.Anb,
			Schb:                    w.Schb,
			Abwb:                    w.Abwb,
			NameFuerSpezialisierung: w.NameFuerSpezialisierung,
			BeinhaltetIn:            fmt.Sprintf("bamort-container-%d", w.ContainedIn),
			ContainedIn:             w.ContainedIn,
			Magisch: Magisch{
				IstMagisch:  w.IstMagisch,
				Abw:         w.Abw,
				Ausgebrannt: w.Ausgebrannt,
			},
		})
	}

	// Ausrüstung
	for i, a := range char.Ausruestung {
		vtt.Ausruestung = append(vtt.Ausruestung, Ausruestung{
			ImportBase: ImportBase{
				ID:   fmt.Sprintf("bamort-equipment-%d-%d", char.ID, i),
				Name: a.Name,
			},
			Beschreibung: a.Beschreibung,
			Gewicht:      a.Gewicht,
			Wert:         a.Wert,
			Anzahl:       a.Anzahl,
			Bonus:        a.Bonus,
			BeinhaltetIn: fmt.Sprintf("bamort-container-%d", a.ContainedIn),
			ContainedIn:  a.ContainedIn,
			Magisch: Magisch{
				IstMagisch:  a.IstMagisch,
				Abw:         a.Abw,
				Ausgebrannt: a.Ausgebrannt,
			},
		})
	}

	// Behältnisse
	for i, b := range char.Behaeltnisse {
		vtt.Behaeltnisse = append(vtt.Behaeltnisse, Behaeltniss{
			ImportBase: ImportBase{
				ID:   fmt.Sprintf("bamort-container-%d-%d", char.ID, i),
				Name: b.Name,
			},
			Beschreibung: b.Beschreibung,
			Gewicht:      b.Gewicht,
			Wert:         b.Wert,
			Tragkraft:    b.Tragkraft,
			Volumen:      b.Volumen,
			BeinhaltetIn: fmt.Sprintf("bamort-container-%d", b.ContainedIn),
			ContainedIn:  b.ContainedIn,
			Magisch: Magisch{
				IstMagisch:  b.IstMagisch,
				Abw:         b.Abw,
				Ausgebrannt: b.Ausgebrannt,
			},
		})
	}

	// Transportmittel
	for i, tm := range char.Transportmittel {
		vtt.Transportmittel = append(vtt.Transportmittel, Transportation{
			ImportBase: ImportBase{
				ID:   fmt.Sprintf("bamort-transport-%d-%d", char.ID, i),
				Name: tm.Name,
			},
			Beschreibung: tm.Beschreibung,
			Gewicht:      int(tm.Gewicht),
			Wert:         tm.Wert,
			Tragkraft:    tm.Tragkraft,
			BeinhaltetIn: fmt.Sprintf("bamort-container-%d", tm.ContainedIn),
			ContainedIn:  tm.ContainedIn,
			Magisch: Magisch{
				IstMagisch:  tm.IstMagisch,
				Abw:         tm.Abw,
				Ausgebrannt: tm.Ausgebrannt,
			},
		})
	}

	return vtt, nil
}

// ExportCharToVTTFile exports a character to VTT JSON file
func ExportCharToVTTFile(char *models.Char, filename string) error {
	vtt, err := ExportCharToVTT(char)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(vtt)
}

// ExportSpellsToCSV exports spell master data to CSV format
func ExportSpellsToCSV(spells []models.Spell, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Header
	header := []string{
		"Nr", "game_system", "name", "Beschreibung", "Quelle", "source_id",
		"page_number", "bonus", "stufe", "ap", "Art", "Zauberdauer",
		"Reichweite", "Wirkungsziel", "Wirkungsbereich", "Wirkungsdauer",
		"Ursprung", "Category", "learning_category", "Agens", "Reagens",
		"Material (Kosten)",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Data rows
	for i, spell := range spells {
		row := []string{
			strconv.Itoa(i + 1),
			spell.GameSystem,
			spell.Name,
			spell.Beschreibung,
			spell.Quelle,
			strconv.Itoa(int(spell.SourceID)),
			strconv.Itoa(spell.PageNumber),
			strconv.Itoa(spell.Bonus),
			strconv.Itoa(spell.Stufe),
			spell.AP,
			spell.Art,
			spell.Zauberdauer,
			spell.Reichweite,
			spell.Wirkungsziel,
			spell.Wirkungsbereich,
			spell.Wirkungsdauer,
			spell.Ursprung,
			spell.Category,
			spell.LearningCategory,
			"", // Agens - not in current model
			"", // Reagens - not in current model
			"", // Material - not in current model
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ExportCharToCSV exports a character to CSV format (MOAM-compatible)
func ExportCharToCSV(char *models.Char, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Helper to write a row
	writeRow := func(fields ...string) error {
		return writer.Write(fields)
	}

	// Header: Name, Typ, Grad, Stand, Glaube, Herkunft
	if err := writeRow("Name", "Typ", "Grad", "Stand", "Glaube", "Herkunft"); err != nil {
		return err
	}
	stand := ""
	herkunft := ""
	if err := writeRow(char.Name, char.Typ, strconv.Itoa(char.Grad), stand, char.Glaube, herkunft); err != nil {
		return err
	}

	// Spezialisierung
	if err := writeRow("Spezialisierung", strings.Join(char.Spezialisierung, ";")); err != nil {
		return err
	}
	if err := writeRow(""); err != nil { // Empty line
		return err
	}

	// Basiseigenschaften
	if err := writeRow("Basiseigenschaften"); err != nil {
		return err
	}
	eigenschaftenMap := make(map[string]int)
	for _, e := range char.Eigenschaften {
		eigenschaftenMap[e.Name] = e.Value
	}
	if err := writeRow("St", strconv.Itoa(eigenschaftenMap["St"])); err != nil {
		return err
	}
	if err := writeRow("Gs", strconv.Itoa(eigenschaftenMap["Gs"])); err != nil {
		return err
	}
	if err := writeRow("Gw", strconv.Itoa(eigenschaftenMap["Gw"])); err != nil {
		return err
	}
	if err := writeRow("Ko", strconv.Itoa(eigenschaftenMap["Ko"])); err != nil {
		return err
	}
	if err := writeRow("In", strconv.Itoa(eigenschaftenMap["In"])); err != nil {
		return err
	}
	if err := writeRow("Zt", strconv.Itoa(eigenschaftenMap["Zt"])); err != nil {
		return err
	}
	if err := writeRow("Au", strconv.Itoa(eigenschaftenMap["Au"])); err != nil {
		return err
	}
	if err := writeRow("pA", strconv.Itoa(eigenschaftenMap["PA"])); err != nil {
		return err
	}
	if err := writeRow("Wk", strconv.Itoa(eigenschaftenMap["Wk"])); err != nil {
		return err
	}

	// LP, AP, B, SchB, AbB, AnB
	if err := writeRow("LP", "AP", "B", "SchB", "AbB", "AnB"); err != nil {
		return err
	}
	if err := writeRow(
		strconv.Itoa(char.Lp.Max),
		strconv.Itoa(char.Ap.Max),
		strconv.Itoa(char.B.Max),
		"0", // SchB - not in model
		"0", // AbB - not in model
		"0", // AnB - not in model
	); err != nil {
		return err
	}
	if err := writeRow(""); err != nil {
		return err
	}

	// Raufen, Abwehr (not in current model, use defaults)
	if err := writeRow("Raufen", "Abwehr"); err != nil {
		return err
	}
	if err := writeRow("7", "14"); err != nil {
		return err
	}

	// Resistenz
	if err := writeRow("Resistenz Geist", "Resistenz Körper"); err != nil {
		return err
	}
	if err := writeRow("0", "0"); err != nil {
		return err
	}
	if err := writeRow("Bonus Resistenz Geist", "Bonus Resistenz Körper"); err != nil {
		return err
	}
	if err := writeRow("0", "0"); err != nil {
		return err
	}
	if err := writeRow(""); err != nil {
		return err
	}

	// Waffen
	if err := writeRow("Waffe", "Erfolgswert", "Angriffsbonus", "Schadensbonus", "Abwehrbonus", "Praxispunkte"); err != nil {
		return err
	}
	for _, w := range char.Waffen {
		// Find corresponding weapon skill
		erfolgswert := 0
		for _, ws := range char.Waffenfertigkeiten {
			if strings.Contains(w.Name, ws.Name) || strings.Contains(ws.Name, w.Name) {
				erfolgswert = ws.Fertigkeitswert
				break
			}
		}
		if err := writeRow(
			w.Name,
			strconv.Itoa(erfolgswert),
			strconv.Itoa(w.Anb),
			strconv.Itoa(w.Schb),
			strconv.Itoa(w.Abwb),
			"0",
		); err != nil {
			return err
		}
	}
	if err := writeRow(""); err != nil {
		return err
	}

	// Rüstung (armor from equipment)
	if err := writeRow("Rüstung", "RK", "Rüstungsbonus"); err != nil {
		return err
	}
	for _, a := range char.Ausruestung {
		if strings.Contains(strings.ToLower(a.Name), "rüstung") || strings.Contains(strings.ToLower(a.Name), "armor") {
			if err := writeRow(a.Name, strconv.Itoa(a.Bonus), "0"); err != nil {
				return err
			}
		}
	}
	if err := writeRow(""); err != nil {
		return err
	}

	// Fertigkeiten
	if err := writeRow("Fertigkeit", "Erfolgswert", "Bonus", "Praxispunkte"); err != nil {
		return err
	}
	for _, f := range char.Fertigkeiten {
		desc := f.Beschreibung
		if desc != "" {
			desc = " (" + desc + ")"
		}
		if err := writeRow(
			f.Name+desc,
			strconv.Itoa(f.Fertigkeitswert),
			strconv.Itoa(f.Bonus),
			strconv.Itoa(f.Pp),
		); err != nil {
			return err
		}
	}
	if err := writeRow(""); err != nil {
		return err
	}

	// Zauber section (if character has spells)
	if len(char.Zauber) > 0 {
		if err := writeRow("Zaubern", "ZauB"); err != nil {
			return err
		}
		if err := writeRow("0", "0"); err != nil { // Placeholder values
			return err
		}
		if err := writeRow("Zauber", "Bonus", "Praxispunkte"); err != nil {
			return err
		}
		for _, z := range char.Zauber {
			if err := writeRow(z.Name, strconv.Itoa(z.Bonus), "0"); err != nil {
				return err
			}
		}
		if err := writeRow(""); err != nil {
			return err
		}
	}

	// Ausrüstung (equipment list)
	if err := writeRow("Ausrüstung"); err != nil {
		return err
	}
	var equipmentNames []string
	for _, a := range char.Ausruestung {
		equipmentNames = append(equipmentNames, a.Name)
	}
	for _, b := range char.Behaeltnisse {
		equipmentNames = append(equipmentNames, b.Name)
	}
	for _, t := range char.Transportmittel {
		equipmentNames = append(equipmentNames, t.Name)
	}
	if len(equipmentNames) > 0 {
		if err := writeRow(strings.Join(equipmentNames, ";")); err != nil {
			return err
		}
	}
	if err := writeRow(""); err != nil {
		return err
	}

	// Erfahrung
	if err := writeRow("Erfahrung"); err != nil {
		return err
	}
	if err := writeRow("Erfahrungsschatz", "EP", "Gold"); err != nil {
		return err
	}
	if err := writeRow(
		strconv.Itoa(char.Erfahrungsschatz.ES),
		"0", // EP not in model
		"0", // Gold not in model
	); err != nil {
		return err
	}

	return nil
}
