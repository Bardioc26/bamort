package gsmaster

import (
	"bamort/models"
	"errors"
	"fmt"
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
// here we have a “Stufe” (1..12) plus a “School” (e.g. "Beherrschen", "Bewegen", etc.)
type SpellDefinition struct {
	Name   string `json:"name"`
	Stufe  int    `json:"level"`
	School string `json:"school"` // e.g. "Beherrschen", "Bewegen", "Erkennen", etc.
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

// CalculateImprovementCost: Kosten zum Steigern von +X auf +X+1
func CalculateSkillImprovementCost(skill string, class string, currentSkillLevel int) (*models.LearnCost, error) {
	return CalculateImprovementCost(skill, class, currentSkillLevel)
}

// Deprecated old static function, now using DB data
func CalculateImprovementCost(skill string, class string, currentSkillLevel int) (*models.LearnCost, error) {
	/*
		if !Config.AllowedGroups[class][skill.Group] {
			return 0, fmt.Errorf("die Klasse %s darf %s nicht lernen", class, skill.Group)
		}
	*/
	lCost := models.LearnCost{}
	var skl models.Skill
	if err := skl.First(skill); err != nil {
		return nil, errors.New("unbekannte Fertigkeit")
	}

	grpMap, ok := Config.ImprovementCost[SkillGroup(skl.Category)]
	if !ok {
		return nil, errors.New("keine Improvement-Daten für diese Gruppe")
	}
	diffMap, ok := grpMap[Difficulty(skl.Difficulty)]
	if !ok {
		return nil, errors.New("keine Improvement-Daten für diese Schwierigkeit")
	}

	neededLE, found := diffMap[fmt.Sprintf("%d", currentSkillLevel+1)]
	if !found {
		return nil, fmt.Errorf("kein Eintrag für Bonus %d→%d", currentSkillLevel, currentSkillLevel+1)
	}
	epPerTE, ok := Config.EPPerTE[CharClass(class)][SkillGroup(skl.Category)]
	if !ok {
		return nil, fmt.Errorf("keine EP-Kosten für %s bei %s", class, skl.Category)
	}
	lCost.LE = neededLE
	lCost.Stufe = currentSkillLevel + 1
	lCost.Ep = neededLE * epPerTE
	lCost.Money = lCost.Ep
	return &lCost, nil
}
