package gsmaster

/*
Fertigkeiten werden in Fertigkeitsgruppen eingeteilt zum Beispiel Alltag, Freiland und Kampf
innerhalb der Gruppe saind die Fertigkeiten eingeordnet nach Schwierigkeit. Also leicht, normal, schwer und sehr schwer.
Je schwerer eine Fertigkeit zu lernen ist um so mehr Lernpunkte müssen führ ihre Verbesserung investiert werden.
Auch pro Verbesserungsstufe sind je nach Schwierigkeit unterschiedliche viele Lernpunkte zu investieren.

Beispiel
Tabelle Lern und Verbesserungskosten:
Alltag
Lernen:
leicht: 1LE, {Bootfahren, Kochen, Reiten, Schwimmen, Tanzen, ...}
normal: 1LE, {schreiben, Sprache, Lesen, ...}
schwer: 2LE, {Erste Hilfe, Etikette, Fälschen, ...}
sehr schwer: 10LE, {Gerätekunde, Geschäftssinn, ...}
Verbessern:

	+9;		+10;	+11;	+12;	+13;...

leicht: 	-;		-;		-;		-;		1;...
normal: 	1;		1;		1;		1;		2;...
schwer: 	2;		2;		5;		5;		10;...
sehr schwer: 5;		5;		10;		10;		20;...

Freiland
Lernen:
leicht: 1LE, {Überleben, ...}
normal: 1LE, {Naturkunde, ...}
schwer: 2LE, {Tarnen, ...}
Verbessern:

	+9;		+10;	+11;	+12;	+13;...

leicht: 	1;		1;		1;		2;		2;...
normal: 	2;		5;		5;		10;		15;...
schwer: 	5;		5;		10;		10;		20;...

Tabelle EP Kosten pro Lerneeinheit

	Alltag;	Freiland;	Kampf;	...

Krieger		20;		20;			10;	...
Magier		30;		30;			-;	...
Schurke		10;		30;			30;	...
...

Die pro Lerneinheit aufzuwendenden Erfahrungspunkte sind für die einzelnen Fertigkeitsgruppen pro charakterklasse unterschiedlich. Für einige Charakterklassen sind einzelne Fertigkeitsgruppen nicht erlernbar.

Erstelle aus diesen Informationen eine Datenstruktur, die es ermöglicht die Lernpunkte für eine Fertigkeit zu berechnen.
Ziehe dazu die Dateien model.go in den Verzeichnissen backend/gsmaster, backend/skills, und backend/character zu Rate.
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type SkillGroup string

const (
	GroupAlltag   SkillGroup = "Alltag"
	GroupFreiland SkillGroup = "Freiland"
	GroupKampf    SkillGroup = "Kampf"
)

type Difficulty string

const (
	DiffLeicht     Difficulty = "leicht"
	DiffNormal     Difficulty = "normal"
	DiffSchwer     Difficulty = "schwer"
	DiffSehrSchwer Difficulty = "sehr_schwer"
)

type LevelConfig struct {
	BaseLearnCost   map[SkillGroup]map[Difficulty]int            `json:"baseLearnCost"`
	ImprovementCost map[SkillGroup]map[Difficulty]map[string]int `json:"improvementCost"`
	EPPerTE         map[CharClass]map[SkillGroup]int             `json:"epPerTE"`
	AllowedGroups   map[CharClass]map[SkillGroup]bool            `json:"allowedGroups"`
}

// SkillDefinition beschreibt eine Fertigkeit
type SkillDefinition struct {
	Name       string
	Group      SkillGroup
	Difficulty Difficulty
}

// Erstmals Lernen: (SkillGroup->Difficulty->LE)
var BaseLearnCost = map[SkillGroup]map[Difficulty]int{
	GroupAlltag: {
		DiffLeicht:     1,
		DiffNormal:     1,
		DiffSchwer:     2,
		DiffSehrSchwer: 10,
	},
	GroupFreiland: {
		DiffLeicht:     1,
		DiffNormal:     1,
		DiffSchwer:     2,
		DiffSehrSchwer: 10,
	},
	GroupKampf: {
		DiffLeicht:     2,
		DiffNormal:     2,
		DiffSchwer:     4,
		DiffSehrSchwer: 12,
	},
}

// Verbesserungs-Kosten in einem Feld differenziert nach Gruppe & Schwierigkeit
// (SkillGroup->Difficulty->(aktuellerLevel+1)->LE)
var ImprovementCost = map[SkillGroup]map[Difficulty]map[int]int{
	GroupAlltag: {
		DiffLeicht: {
			9: 0, 10: 0, 11: 0, 12: 0, 13: 1, 14: 1,
		},
		DiffNormal: {
			9: 1, 10: 1, 11: 1, 12: 1, 13: 2,
		},
		DiffSchwer: {
			9: 2, 10: 2, 11: 5, 12: 5, 13: 10,
		},
		DiffSehrSchwer: {
			9: 5, 10: 5, 11: 10, 12: 10, 13: 20,
		},
	},
	GroupFreiland: {
		DiffLeicht: {
			9: 1, 10: 1, 11: 1, 12: 2, 13: 2,
		},
		DiffNormal: {
			9: 2, 10: 5, 11: 5, 12: 10, 13: 15,
		},
		DiffSchwer: {
			9: 5, 10: 5, 11: 10, 12: 10, 13: 20,
		},
		// z.B. keine Daten für sehr_schwer => 10er-Blöcke oder leer
	},
	// Gruppe Kampf nur Beispiele
	GroupKampf: {
		DiffLeicht: {
			9: 1, 10: 1, 11: 2, 12: 2, 13: 4,
		},
		DiffNormal: {
			9: 2, 10: 2, 11: 3,
		},
	},
}

// Beispiel: EP-Kosten pro LE je Klasse & Gruppe
type CharClass string

const (
	ClassKrieger CharClass = "Krieger"
	ClassMagier  CharClass = "Magier"
	ClassSchurke CharClass = "Schurke"
)

var EPPerLE = map[CharClass]map[SkillGroup]int{
	ClassKrieger: {
		GroupAlltag:   20,
		GroupFreiland: 20,
		GroupKampf:    10,
	},
	ClassMagier: {
		GroupAlltag:   30,
		GroupFreiland: 30,
		// kein Kampf => 0 oder kein Eintrag
	},
	ClassSchurke: {
		GroupAlltag:   10,
		GroupFreiland: 30,
		GroupKampf:    30,
	},
}

// Eventuell Erlaubnis pro Klasse/Gruppe
var AllowedGroups = map[CharClass]map[SkillGroup]bool{
	ClassKrieger: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    true,
	},
	ClassMagier: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    false,
	},
	ClassSchurke: {
		GroupAlltag:   true,
		GroupFreiland: true,
		GroupKampf:    true,
	},
}

var Config LevelConfig // holds all loaded data

func init() {
	// Adjust path as needed
	filePath := "/data/dev/bamort/config/leveldata.json"
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to open JSON file: %w", err))
	}
	defer file.Close()

	// Decode the JSON file into the ExportData structure
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
		panic(fmt.Errorf("failed to decode JSON file: %w", err))
	}
}

// CalculateLearnCost: erstmalige Kosten in EP
// Then refer to Config in your calculations:
func CalculateLearnCost(skill SkillDefinition, class CharClass) (int, error) {
	if !Config.AllowedGroups[class][skill.Group] {
		return 0, fmt.Errorf("die Klasse %s darf %s nicht lernen", class, skill.Group)
	}

	groupMap, ok := Config.BaseLearnCost[skill.Group]
	if !ok {
		return 0, errors.New("unbekannte Gruppe")
	}
	baseLE, ok := groupMap[skill.Difficulty]
	if !ok {
		return 0, errors.New("keine LE-Definition für diese Schwierigkeit")
	}
	epPerTE, ok := Config.EPPerTE[class][skill.Group]
	if !ok {
		return 0, fmt.Errorf("keine EP-Kosten für %s bei %s", class, skill.Group)
	}
	return baseLE * (epPerTE * 3), nil
}

// CalculateImprovementCost: Kosten zum Steigern von +X auf +X+1
func CalculateImprovementCost(skill SkillDefinition, class CharClass, currentSkillLevel int) (int, error) {
	if !Config.AllowedGroups[class][skill.Group] {
		return 0, fmt.Errorf("die Klasse %s darf %s nicht lernen", class, skill.Group)
	}
	grpMap, ok := Config.ImprovementCost[skill.Group]
	if !ok {
		return 0, errors.New("keine Improvement-Daten für diese Gruppe")
	}
	diffMap, ok := grpMap[skill.Difficulty]
	if !ok {
		return 0, errors.New("keine Improvement-Daten für diese Schwierigkeit")
	}

	neededLE, found := diffMap[fmt.Sprintf("%d", currentSkillLevel+1)]
	if !found {
		return 0, fmt.Errorf("kein Eintrag für Bonus %d→%d", currentSkillLevel, currentSkillLevel+1)
	}
	epPerTE, ok := Config.EPPerTE[class][skill.Group]
	if !ok {
		return 0, fmt.Errorf("keine EP-Kosten für %s bei %s", class, skill.Group)
	}
	return neededLE * epPerTE, nil
}
