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
type Difficulty string
type CharClass string

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
	//AllowedGroups           map[CharClass]map[SkillGroup]bool            `json:"allowedGroups"`
	SpellLearnCost          map[int]int                   `json:"spellLearnCost"`
	SpellEPPerSchoolByClass map[CharClass]map[string]int  `json:"spellEPPerSchoolByClass"`
	AllowedSchools          map[CharClass]map[string]bool `json:"allowedSchools"`
}

// SkillDefinition beschreibt eine Fertigkeit
type SkillDefinition struct {
	Name       string
	Group      SkillGroup
	Difficulty Difficulty
}

// SpellDefinition differs from SkillDefinition,
// here we have a “Stufe” (1..12) plus a “School” (e.g. "Beherr", "Beweg", etc.)
type SpellDefinition struct {
	Name   string `json:"name"`
	Stufe  int    `json:"level"`
	School string `json:"school"` // e.g. "Beherr", "Beweg", "Erken", etc.
	CostEP int    `json:"cost_ep"`
	CostLE int    `json:"cost_le"`
}

// SpellLearnCost: Stufe->Lerneinheiten (1..12)
var SpellLearnCost = map[int]int{}

// Erstmals Lernen: (SkillGroup->Difficulty->LE)
var BaseLearnCost = map[SkillGroup]map[Difficulty]int{}

// Verbesserungs-Kosten in einem Feld differenziert nach Gruppe & Schwierigkeit
// (SkillGroup->Difficulty->(aktuellerLevel+1)->LE)
var ImprovementCost = map[SkillGroup]map[Difficulty]map[int]int{}

// SpellEPPerSchoolByClass: EP-Kosten pro TE, depends on both
// character class and spell school
var SpellEPPerSchoolByClass = map[CharClass]map[string]int{}

var EPPerTE = map[CharClass]map[SkillGroup]int{}

/*
// Eventuell Erlaubnis pro Klasse/Gruppe
var AllowedGroups = map[CharClass]map[SkillGroup]bool}
*/

var Config LevelConfig // holds all loaded data

func loadLevelingConfig(opts ...string) {
	// Adjust path as needed
	filePath := "../testdata/leveldata.json"
	if len(opts) > 0 {
		filePath = opts[0]
	}
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

// CalculateSpellLearnCost combines SpellLearnCost with SpellEPPerSchoolByClass
func CalculateSpellLearnCost(spell string, class string) (int, error) {
	if Config.AllowedSchools == nil {
		loadLevelingConfig()
	}

	var spl Spell
	if err := spl.First(spell); err != nil {
		return 0, errors.New("unbekannter Zauberspruch")
	}

	if !Config.AllowedSchools[CharClass(class)][spl.Category] {
		return 0, fmt.Errorf("die Klasse %s darf die Schule %s nicht lernen", class, spl.Category)
	}
	neededLE, ok := Config.SpellLearnCost[spl.Stufe]
	if !ok {
		return 0, fmt.Errorf("ungültige Zauberstufe: %d", spl.Stufe)
	}

	classMap, ok := Config.SpellEPPerSchoolByClass[CharClass(class)]
	if !ok {
		return 0, fmt.Errorf("keine EP-Tabelle für Klasse: %s", class)
	}

	epPerTE, found := classMap[spl.Category]
	if !found {
		return 0, fmt.Errorf("unbekannte Schule '%s' bei Klasse '%s'", spl.Category, class)
	}

	// Gesamt-EP = benötigte LE * EP pro LE.
	totalEP := neededLE * (epPerTE * 3)
	// +6 EP for elves
	if class == "Elf" {
		totalEP += 6
	}

	return totalEP, nil
}

// CalculateSkillLearnCost: erstmalige Kosten in EP
// Then refer to Config in your calculations:
func CalculateSkillLearnCost(skill string, class string) (int, error) {
	/*
		if !Config.AllowedGroups[class][skill.Group] {
			return 0, fmt.Errorf("die Klasse %s darf %s nicht lernen", class, skill.Group)
		}
	*/
	var skl Skill
	if err := skl.First(skill); err != nil {
		return 0, errors.New("unbekannte Fertigkeit")
	}

	groupMap, ok := Config.BaseLearnCost[SkillGroup(skl.Category)]
	if !ok {
		return 0, errors.New("unbekannte Gruppe")
	}

	baseLE, ok := groupMap[Difficulty(skl.Difficulty)]
	if !ok {
		return 0, errors.New("keine LE-Definition für diese Schwierigkeit")
	}
	epPerTE, ok := Config.EPPerTE[CharClass(class)][SkillGroup(skl.Category)]
	if !ok {
		return 0, fmt.Errorf("keine EP-Kosten für %s bei %s", class, skl.Category)
	}
	totalEP := baseLE * (epPerTE * 3)
	// +6 EP for elves
	if class == "Elf" {
		totalEP += 6
	}

	return totalEP, nil
}

// CalculateImprovementCost: Kosten zum Steigern von +X auf +X+1
func CalculateSkillImprovementCost(skill string, class string, currentSkillLevel int) (int, error) {
	return CalculateImprovementCost(skill, class, currentSkillLevel)
}

func CalculateImprovementCost(skill string, class string, currentSkillLevel int) (int, error) {
	/*
		if !Config.AllowedGroups[class][skill.Group] {
			return 0, fmt.Errorf("die Klasse %s darf %s nicht lernen", class, skill.Group)
		}
	*/
	var skl Skill
	if err := skl.First(skill); err != nil {
		return 0, errors.New("unbekannte Fertigkeit")
	}

	grpMap, ok := Config.ImprovementCost[SkillGroup(skl.Category)]
	if !ok {
		return 0, errors.New("keine Improvement-Daten für diese Gruppe")
	}
	diffMap, ok := grpMap[Difficulty(skl.Difficulty)]
	if !ok {
		return 0, errors.New("keine Improvement-Daten für diese Schwierigkeit")
	}

	neededLE, found := diffMap[fmt.Sprintf("%d", currentSkillLevel+1)]
	if !found {
		return 0, fmt.Errorf("kein Eintrag für Bonus %d→%d", currentSkillLevel, currentSkillLevel+1)
	}
	epPerTE, ok := Config.EPPerTE[CharClass(class)][SkillGroup(skl.Category)]
	if !ok {
		return 0, fmt.Errorf("keine EP-Kosten für %s bei %s", class, skl.Category)
	}
	return neededLE * epPerTE, nil
}

func CalculateLearnCost(skillType string, name string, class string) (int, error) {
	if skillType == "skill" {
		return CalculateSkillLearnCost(name, class)
	} else if skillType == "spell" {
		return CalculateSpellLearnCost(name, class)
	}
	return 0, errors.New("unknown skill type")
}
