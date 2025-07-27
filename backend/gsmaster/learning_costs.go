package gsmaster

import (
	"bamort/models"
	"errors"
	"fmt"
)

// LearningCostsTable strukturiert die Daten aus Lerntabellen.md
type LearningCostsTable struct {
	// EP-Kosten für 1 Trainingseinheit (TE) pro Charakterklasse und Fertigkeitskategorie
	EPPerTE map[string]map[string]int

	// EP-Kosten für 1 Lerneinheit (LE) für Zauber pro Charakterklasse und Zauberschule
	SpellEPPerLE map[string]map[string]int

	// LE-Kosten für Fertigkeiten basierend auf Schwierigkeit
	BaseLearnCost map[string]map[string]int

	// TE-Kosten für Verbesserungen basierend auf Kategorie, Schwierigkeit und aktuellem Wert
	ImprovementCost map[string]map[string]map[int]int
}

// Initialisierung der Lerntabellen mit allen 15 Charakterklassen aus Lerntabellen.md
var learningCosts = &LearningCostsTable{
	// EP-Kosten für 1 Trainingseinheit (TE) nach Klasse und Kategorie
	EPPerTE: map[string]map[string]int{
		"As": { // Assassine
			"Alltag":    20,
			"Freiland":  20,
			"Halbwelt":  20,
			"Kampf":     30,
			"Körper":    10,
			"Sozial":    20,
			"Unterwelt": 10,
			"Waffen":    20,
			"Wissen":    20,
		},
		"Bb": { // Barbar
			"Alltag":    20,
			"Freiland":  10,
			"Halbwelt":  30,
			"Kampf":     10,
			"Körper":    10,
			"Sozial":    30,
			"Unterwelt": 30,
			"Waffen":    20,
			"Wissen":    40,
		},
		"Gl": { // Glücksritter
			"Alltag":    20,
			"Freiland":  30,
			"Halbwelt":  10,
			"Kampf":     20,
			"Körper":    30,
			"Sozial":    10,
			"Unterwelt": 30,
			"Waffen":    20,
			"Wissen":    20,
		},
		"Hä": { // Händler
			"Alltag":    10,
			"Freiland":  20,
			"Halbwelt":  20,
			"Kampf":     20,
			"Körper":    20,
			"Sozial":    10,
			"Unterwelt": 40,
			"Waffen":    20,
			"Wissen":    20,
		},
		"Kr": { // Krieger
			"Alltag":    20,
			"Freiland":  30,
			"Halbwelt":  30,
			"Kampf":     10,
			"Körper":    20,
			"Sozial":    20,
			"Unterwelt": 30,
			"Waffen":    10,
			"Wissen":    40,
		},
		"Sp": { // Spitzbube
			"Alltag":    20,
			"Freiland":  30,
			"Halbwelt":  10,
			"Kampf":     40,
			"Körper":    10,
			"Sozial":    10,
			"Unterwelt": 10,
			"Waffen":    20,
			"Wissen":    30,
		},
		"Wa": { // Waldläufer
			"Alltag":    20,
			"Freiland":  10,
			"Halbwelt":  20,
			"Kampf":     20,
			"Körper":    10,
			"Sozial":    30,
			"Unterwelt": 30,
			"Waffen":    20,
			"Wissen":    30,
		},
		"Ba": { // Barde
			"Alltag":    10,
			"Freiland":  20,
			"Halbwelt":  20,
			"Kampf":     40,
			"Körper":    20,
			"Sozial":    30,
			"Unterwelt": 40,
			"Waffen":    40,
			"Wissen":    10,
		},
		"Or": { // Ordenskrieger
			"Alltag":    20,
			"Freiland":  30,
			"Halbwelt":  40,
			"Kampf":     10,
			"Körper":    20,
			"Sozial":    20,
			"Unterwelt": 40,
			"Waffen":    10,
			"Wissen":    20,
		},
		"Dr": { // Druide
			"Alltag":    20,
			"Freiland":  10,
			"Halbwelt":  30,
			"Kampf":     40,
			"Körper":    20,
			"Sozial":    30,
			"Unterwelt": 40,
			"Waffen":    40,
			"Wissen":    10,
		},
		"Hx": { // Hexer
			"Alltag":    20,
			"Freiland":  20,
			"Halbwelt":  30,
			"Kampf":     40,
			"Körper":    30,
			"Sozial":    20,
			"Unterwelt": 30,
			"Waffen":    40,
			"Wissen":    20,
		},
		"Ma": { // Magier
			"Alltag":    20,
			"Freiland":  30,
			"Halbwelt":  40,
			"Kampf":     40,
			"Körper":    30,
			"Sozial":    30,
			"Unterwelt": 40,
			"Waffen":    40,
			"Wissen":    10,
		},
		"PB": { // Priester Beschützer
			"Alltag":    20,
			"Freiland":  30,
			"Halbwelt":  30,
			"Kampf":     40,
			"Körper":    30,
			"Sozial":    10,
			"Unterwelt": 40,
			"Waffen":    40,
			"Wissen":    20,
		},
		"PS": { // Priester Streiter
			"Alltag":    20,
			"Freiland":  30,
			"Halbwelt":  40,
			"Kampf":     30,
			"Körper":    30,
			"Sozial":    30,
			"Unterwelt": 40,
			"Waffen":    30,
			"Wissen":    20,
		},
		"Sc": { // Schamane
			"Alltag":    20,
			"Freiland":  10,
			"Halbwelt":  40,
			"Kampf":     40,
			"Körper":    20,
			"Sozial":    20,
			"Unterwelt": 40,
			"Waffen":    40,
			"Wissen":    20,
		},
	},

	// EP-Kosten für 1 Lerneinheit (LE) für Zauber pro Charakterklasse und Zauberschule
	SpellEPPerLE: map[string]map[string]int{
		"Dr": { // Druide
			"Beherrschen": 90,
			"Bewegen":     60,
			"Erkennen":    120,
			"Erschaffen":  90,
			"Formen":      60,
			"Verändern":   90,
			"Zerstören":   120,
			"Wunder":      0, // Nicht verfügbar
			"Dweomer":     30,
			"Lied":        0, // Nicht verfügbar
		},
		"Hx": { // Hexer
			"Beherrschen": 30,
			"Bewegen":     90,
			"Erkennen":    90,
			"Erschaffen":  90,
			"Formen":      60,
			"Verändern":   30,
			"Zerstören":   60,
			"Wunder":      0, // Nicht verfügbar
			"Dweomer":     90,
			"Lied":        0, // Nicht verfügbar
		},
		"Ma": { // Magier (* = Spezialgebiet für 30 EP)
			"Beherrschen": 60, // *30
			"Bewegen":     60, // *30
			"Erkennen":    60, // *30
			"Erschaffen":  60, // *30
			"Formen":      60, // *30
			"Verändern":   60, // *30
			"Zerstören":   60, // *30
			"Wunder":      0,  // Nicht verfügbar
			"Dweomer":     120,
			"Lied":        0, // Nicht verfügbar
		},
		"PB": { // Priester Beschützer
			"Beherrschen": 90,
			"Bewegen":     90,
			"Erkennen":    60,
			"Erschaffen":  90,
			"Formen":      90,
			"Verändern":   90,
			"Zerstören":   90,
			"Wunder":      30,
			"Dweomer":     120,
			"Lied":        0, // Nicht verfügbar
		},
		"PS": { // Priester Streiter
			"Beherrschen": 90,
			"Bewegen":     90,
			"Erkennen":    90,
			"Erschaffen":  90,
			"Formen":      90,
			"Verändern":   90,
			"Zerstören":   60,
			"Wunder":      30,
			"Dweomer":     120,
			"Lied":        0, // Nicht verfügbar
		},
		"Sc": { // Schamane
			"Beherrschen": 90,
			"Bewegen":     90,
			"Erkennen":    60,
			"Erschaffen":  60,
			"Formen":      90,
			"Verändern":   90,
			"Zerstören":   90,
			"Wunder":      30,
			"Dweomer":     120,
			"Lied":        0, // Nicht verfügbar
		},
		"Ba": { // Barde
			"Beherrschen": 0, // Nicht verfügbar
			"Bewegen":     0, // Nicht verfügbar
			"Erkennen":    0, // Nicht verfügbar
			"Erschaffen":  0, // Nicht verfügbar
			"Formen":      0, // Nicht verfügbar
			"Verändern":   0, // Nicht verfügbar
			"Zerstören":   0, // Nicht verfügbar
			"Wunder":      0, // Nicht verfügbar
			"Dweomer":     0, // Nicht verfügbar
			"Lied":        30,
		},
	},

	// LE-Kosten für Fertigkeiten basierend auf Schwierigkeit (aus Lerntabellen.md)
	BaseLearnCost: map[string]map[string]int{
		"leicht": {
			"LE": 1, // Alltag: leicht (1 LE), Körper: leicht (1 LE), Kampf: leicht (1 LE)
		},
		"normal": {
			"LE": 2, // Alltag: normal (1 LE), Freiland: normal (2 LE), Halbwelt: normal (2 LE), Kampf: normal (2 LE), Sozial: normal (2 LE), Körper: normal (1 LE), Waffen: normal (4 LE), Wissen: normal (2 LE) - Durchschnitt
		},
		"schwer": {
			"LE": 4, // Alltag: schwer (2 LE), Freiland: schwer (4 LE), Halbwelt: schwer (2 LE), Sozial: schwer (4 LE), Körper: schwer (2 LE), Waffen: schwer (6 LE), Wissen: schwer (2 LE) - Durchschnitt
		},
		"sehr schwer": {
			"LE": 10, // Alltag: sehr schwer (10 LE), Halbwelt: sehr schwer (10 LE), Kampf: sehr schwer (10 LE), Waffen: sehr schwer (8 LE) - Durchschnitt
		},
	},

	// TE-Kosten für Verbesserungen basierend auf Kategorie, Schwierigkeit und aktuellem Wert
	ImprovementCost: map[string]map[string]map[int]int{
		"Alltag": {
			"leicht": {
				0: 1, 1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1, 9: 1, 10: 1, 11: 1, 12: 1, 13: 2, 14: 2, 15: 2, 16: 2, 17: 2, 18: 2, 19: 2, 20: 2,
			},
			"normal": {
				0: 1, 1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1, 9: 1, 10: 1, 11: 2, 12: 2, 13: 2, 14: 2, 15: 2, 16: 2, 17: 2, 18: 2, 19: 2, 20: 2,
			},
			"schwer": {
				0: 1, 1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 1, 8: 1, 9: 2, 10: 2, 11: 2, 12: 2, 13: 2, 14: 2, 15: 2, 16: 2, 17: 2, 18: 2, 19: 2, 20: 2,
			},
			"sehr schwer": {
				0: 1, 1: 1, 2: 1, 3: 1, 4: 1, 5: 1, 6: 1, 7: 2, 8: 2, 9: 2, 10: 2, 11: 2, 12: 2, 13: 2, 14: 2, 15: 2, 16: 2, 17: 2, 18: 2, 19: 2, 20: 2,
			},
		},
		// Weitere Kategorien können hier hinzugefügt werden
	},
}

// GetLearningCosts gibt die konfigurierten Lernkosten zurück
func GetLearningCosts() *LearningCostsTable {
	return learningCosts
}

// CalculateSkillLearningCosts berechnet die Kosten für das Lernen einer Fertigkeit
func CalculateSkillLearningCosts(characterClass, category, difficulty string) (*SkillCostResult, error) {
	//Überprüfe ob die tabelle vorhanden ist in der die EP-Kosten pro LE for die einzelnen Kategorien für jede Charakterklasse definiert sind
	if learningCosts.EPPerTE == nil {
		return nil, errors.New("keine EP-per-TE-Definition gefunden")
	}

	// Konvertiere Vollnamen der Charakterklasse zu Abkürzungen falls nötig
	classKey := GetClassAbbreviation(characterClass)

	// Hole die EP-Kosten pro TE für die angegebene Charakterklasse
	classData, exists := learningCosts.EPPerTE[classKey]
	if !exists {
		return nil, fmt.Errorf("unbekannte Charakterklasse: %s (gesucht als %s)", characterClass, classKey)
	}

	// Ermittle die EP pro TE für die angegebene Kategorie
	epPerTE, exists := classData[category]
	if !exists {
		return nil, fmt.Errorf("unbekannte Kategorie '%s' für Klasse %s", category, characterClass)
	}

	// 1 Lerneinheit(LE) kostet 3 mal so viel wie eine Trainingseinheit (TE) +6 EP when der Charakter ein Elf ist

	if learningCosts.BaseLearnCost == nil {
		return nil, errors.New("keine LE-Definition gefunden")
	}

	difficultyData, exists := learningCosts.BaseLearnCost[difficulty]
	if !exists {
		return nil, fmt.Errorf("unbekannte Schwierigkeit: %s", difficulty)
	}

	le := difficultyData["LE"]
	totalEP := epPerTE * le

	return &SkillCostResult{
		CharacterClass: characterClass,
		SkillName:      "",
		Category:       category,
		Difficulty:     difficulty,
		EP:             totalEP,
		LE:             le,
		GoldCost:       totalEP, // 1 EP = 1 GS
		Details: map[string]interface{}{
			"ep_per_te": epPerTE,
			"le_needed": le,
		},
	}, nil
}

// CalculateSpellLearningCosts berechnet die Kosten für das Lernen eines Zaubers
func CalculateSpellLearningCosts(characterClass, spellSchool string, leNeeded int) (*SkillCostResult, error) {
	if learningCosts.SpellEPPerLE == nil {
		return nil, errors.New("keine Zauber-EP-Definition gefunden")
	}

	// Konvertiere Vollnamen zu Abkürzungen falls nötig
	classKey := GetClassAbbreviation(characterClass)

	classData, exists := learningCosts.SpellEPPerLE[classKey]
	if !exists {
		return nil, fmt.Errorf("Charakterklasse %s kann keine Zauber lernen", characterClass)
	}

	epPerLE, exists := classData[spellSchool]
	if !exists || epPerLE == 0 {
		return nil, fmt.Errorf("Charakterklasse %s kann keine Zauber der Schule %s lernen", characterClass, spellSchool)
	}

	totalEP := epPerLE * leNeeded

	return &SkillCostResult{
		CharacterClass: characterClass,
		SkillName:      "",
		Category:       spellSchool,
		Difficulty:     "",
		EP:             totalEP,
		LE:             leNeeded,
		GoldCost:       totalEP, // 1 EP = 1 GS
		Details: map[string]interface{}{
			"ep_per_le":    epPerLE,
			"le_needed":    leNeeded,
			"spell_school": spellSchool,
		},
	}, nil
}

// SkillCostResult definiert das Ergebnis einer Kostenberechnung
type SkillCostResult struct {
	CharacterClass string                 `json:"character_class"`
	SkillName      string                 `json:"skill_name"`
	Category       string                 `json:"category"`
	Difficulty     string                 `json:"difficulty"`
	EP             int                    `json:"ep"`
	LE             int                    `json:"le"`
	GoldCost       int                    `json:"gold_cost"`
	Details        map[string]interface{} `json:"details"`
}

// SkillCostResultNew definiert das Ergebnis einer Kostenberechnung
type SkillCostResultNew struct {
	CharacterID    string                 `json:"character_id"`
	CharacterClass string                 `json:"character_class"`
	SkillName      string                 `json:"skill_name"`
	Category       string                 `json:"category"`
	Difficulty     string                 `json:"difficulty"`
	EP             int                    `json:"ep"`
	LE             int                    `json:"le"`
	GoldCost       int                    `json:"gold_cost"`
	PPUsed         int                    `json:"pp_used"`
	GoldUsed       int                    `json:"gold_used"`
	TargetLevel    int                    `json:"target_level"`
	Details        map[string]interface{} `json:"details"`
}

// CalculateDetailedSkillLearningCost berechnet die Kosten für das Lernen einer Fertigkeit mit Details
func CalculateDetailedSkillLearningCost(skillName, characterClass string) (*models.LearnCost, error) {
	// Fallback-Werte für Skills ohne definierte Kategorie/Schwierigkeit
	category := GetDefaultCategory(skillName)
	difficulty := GetDefaultDifficulty(skillName)

	result, err := CalculateSkillLearningCosts(characterClass, category, difficulty)
	if err != nil {
		return nil, err
	}

	// Konvertiere SkillCostResult zu LearnCost
	return &models.LearnCost{
		Stufe: 0, // Lernen startet bei Stufe 0
		LE:    result.LE,
		Ep:    result.EP,
		Money: result.GoldCost,
	}, nil
}

// CalculateDetailedSkillImprovementCost berechnet die Kosten für die Verbesserung einer Fertigkeit
func CalculateDetailedSkillImprovementCost(skillName, characterClass string, currentLevel int) (*models.LearnCost, error) {
	// Fallback-Werte für Skills ohne definierte Kategorie/Schwierigkeit
	category := GetDefaultCategory(skillName)
	difficulty := GetDefaultDifficulty(skillName)

	// Verwende die Lernkosten als Basis für Verbesserungen
	// In einer vollständigen Implementierung würden hier die ImprovementCost-Tabellen verwendet
	baseCost, err := CalculateSkillLearningCosts(characterClass, category, difficulty)
	if err != nil {
		return nil, err
	}

	// Vereinfachte Verbesserungslogik: höhere Level = höhere Kosten
	improvementFactor := 1.0
	if currentLevel > 10 {
		improvementFactor = 1.5
	} else if currentLevel > 15 {
		improvementFactor = 2.0
	}

	improvedEP := int(float64(baseCost.EP) * improvementFactor)

	// Konvertiere zu LearnCost
	return &models.LearnCost{
		Stufe: currentLevel + 1, // Ziel-Stufe
		LE:    1,                // TE für Verbesserung (meist 1)
		Ep:    improvedEP,
		Money: improvedEP,
	}, nil
}

// CalculateDetailedSpellLearningCost berechnet die Kosten für das Lernen eines Zaubers
func CalculateDetailedSpellLearningCost(spellName, characterClass string) (*models.LearnCost, error) {
	// Standard-Zauberschule bestimmen
	spellSchool := getDefaultSpellSchool(spellName)

	// Standard LE für Zauber
	standardLE := 4

	result, err := CalculateSpellLearningCosts(characterClass, spellSchool, standardLE)
	if err != nil {
		return nil, err
	}

	// Konvertiere SkillCostResult zu LearnCost
	return &models.LearnCost{
		Stufe: 0, // Lernen startet bei Stufe 0
		LE:    result.LE,
		Ep:    result.EP,
		Money: result.GoldCost,
	}, nil
}

// SkillCategoryOption definiert eine Kategorie-Schwierigkeit-Kombination für eine Fertigkeit
type SkillCategoryOption struct {
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
	LE         int    `json:"le"`
}

// GetAvailableSkillCategories gibt alle verfügbaren Kategorie-Kombinationen für eine Fertigkeit zurück
func GetAvailableSkillCategories(skillName string) []SkillCategoryOption {
	// Basierend auf den offiziellen Lerntabellen.md - alle möglichen Kombinationen
	skillOptions := map[string][]SkillCategoryOption{
		"Klettern": {
			{Category: "Alltag", Difficulty: "leicht", LE: 1},
			{Category: "Halbwelt", Difficulty: "leicht", LE: 1},
			{Category: "Körper", Difficulty: "leicht", LE: 1},
		},
		"Glücksspiel": {
			{Category: "Alltag", Difficulty: "leicht", LE: 1},
			{Category: "Halbwelt", Difficulty: "leicht", LE: 1},
		},
		"Reiten": {
			{Category: "Alltag", Difficulty: "leicht", LE: 1},
			{Category: "Kampf", Difficulty: "leicht", LE: 1},
		},
		"Anführen": {
			{Category: "Kampf", Difficulty: "normal", LE: 2},
			{Category: "Sozial", Difficulty: "leicht", LE: 2},
		},
		"Etikette": {
			{Category: "Alltag", Difficulty: "schwer", LE: 2},
			{Category: "Sozial", Difficulty: "leicht", LE: 2},
		},
		"Gassenwissen": {
			{Category: "Halbwelt", Difficulty: "schwer", LE: 2},
			{Category: "Sozial", Difficulty: "normal", LE: 2},
			{Category: "Unterwelt", Difficulty: "leicht", LE: 2},
		},
		"Betäuben": {
			{Category: "Halbwelt", Difficulty: "sehr schwer", LE: 10},
			{Category: "Kampf", Difficulty: "schwer", LE: 10},
		},
		"Athletik": {
			{Category: "Kampf", Difficulty: "normal", LE: 2},
			{Category: "Körper", Difficulty: "schwer", LE: 2},
		},
		"Balancieren": {
			{Category: "Halbwelt", Difficulty: "leicht", LE: 1},
			{Category: "Körper", Difficulty: "leicht", LE: 1},
		},
		"Akrobatik": {
			{Category: "Halbwelt", Difficulty: "normal", LE: 2},
			{Category: "Körper", Difficulty: "schwer", LE: 2},
		},
		"Schleichen": {
			{Category: "Freiland", Difficulty: "schwer", LE: 4},
			{Category: "Unterwelt", Difficulty: "normal", LE: 4},
		},
		"Spurensuche": {
			{Category: "Freiland", Difficulty: "schwer", LE: 4},
			{Category: "Unterwelt", Difficulty: "normal", LE: 4},
		},
		"Tarnen": {
			{Category: "Freiland", Difficulty: "schwer", LE: 4},
			{Category: "Unterwelt", Difficulty: "normal", LE: 4},
		},
		"Stehlen": {
			{Category: "Halbwelt", Difficulty: "schwer", LE: 2},
			{Category: "Unterwelt", Difficulty: "leicht", LE: 2},
		},
		"Verhören": {
			{Category: "Sozial", Difficulty: "normal", LE: 2},
			{Category: "Unterwelt", Difficulty: "normal", LE: 4},
		},
		"Menschenkenntnis": {
			{Category: "Sozial", Difficulty: "schwer", LE: 4},
			{Category: "Unterwelt", Difficulty: "schwer", LE: 4},
		},
		"Schreiben": {
			{Category: "Alltag", Difficulty: "normal", LE: 1},
			{Category: "Wissen", Difficulty: "leicht", LE: 1},
		},
		"Sprache": {
			{Category: "Alltag", Difficulty: "normal", LE: 1},
			{Category: "Wissen", Difficulty: "leicht", LE: 1},
		},
		"Erste Hilfe": {
			{Category: "Alltag", Difficulty: "schwer", LE: 2},
			{Category: "Wissen", Difficulty: "normal", LE: 2},
		},
		"Meditieren": {
			{Category: "Körper", Difficulty: "schwer", LE: 2},
			{Category: "Wissen", Difficulty: "normal", LE: 2},
		},
		"Naturkunde": {
			{Category: "Freiland", Difficulty: "normal", LE: 2},
			{Category: "Wissen", Difficulty: "schwer", LE: 2},
		},
		"Pflanzenkunde": {
			{Category: "Freiland", Difficulty: "normal", LE: 2},
			{Category: "Wissen", Difficulty: "schwer", LE: 2},
		},
		"Tierkunde": {
			{Category: "Freiland", Difficulty: "normal", LE: 2},
			{Category: "Wissen", Difficulty: "schwer", LE: 2},
		},
	}

	if options, exists := skillOptions[skillName]; exists {
		return options
	}

	// Fallback: verwende Standard-Mapping (erste gefundene Kategorie)
	category := GetDefaultCategory(skillName)
	difficulty := GetDefaultDifficulty(skillName)

	// Bestimme LE basierend auf Schwierigkeit
	le := 2 // Standard
	switch difficulty {
	case "leicht":
		le = 1
	case "normal":
		le = 2
	case "schwer":
		le = 4
	case "sehr schwer":
		le = 10
	}

	return []SkillCategoryOption{
		{Category: category, Difficulty: difficulty, LE: le},
	}
}

// GetDefaultCategory gibt die erste (bevorzugte) Kategorie für eine Fertigkeit zurück
func GetDefaultCategory(skillName string) string {
	// WICHTIG: Wir verwenden bewusst die erste gefundene Kategorie als Standard.
	// Für das Lernen ist es unerheblich, aber später wird es für andere Dinge wichtig werden.
	// Die Reihenfolge der Kategorien ist nach Wichtigkeit/Häufigkeit sortiert.
	categoryMap := map[string]string{
		"Stichwaffen":             "Waffen",
		"Einhandschlagwaffen":     "Waffen",
		"Zweihandschlagwaffen":    "Waffen",
		"Geländelauf":             "Körper",
		"Klettern":                "Körper",
		"Schwimmen":               "Körper",
		"Reiten":                  "Körper",
		"Balancieren":             "Körper",
		"Tauchen":                 "Körper",
		"Akrobatik":               "Körper",
		"Athletik":                "Körper",
		"Laufen":                  "Körper",
		"Meditieren":              "Körper",
		"Seilkunst":               "Alltag",
		"Bootfahren":              "Alltag",
		"Glücksspiel":             "Alltag",
		"Wagenlenken":             "Alltag",
		"Schreiben":               "Alltag",
		"Sprache":                 "Alltag",
		"Erste Hilfe":             "Alltag",
		"Etikette":                "Alltag",
		"Gerätekunde":             "Alltag",
		"Geschäftssinn":           "Alltag",
		"Musizieren":              "Alltag",
		"Schleichen":              "Freiland",
		"Spurensuche":             "Freiland",
		"Tarnen":                  "Freiland",
		"Überleben":               "Freiland",
		"Naturkunde":              "Freiland",
		"Pflanzenkunde":           "Freiland",
		"Tierkunde":               "Freiland",
		"Gassenwissen":            "Unterwelt",
		"Stehlen":                 "Unterwelt",
		"Fallen entdecken":        "Unterwelt",
		"Schlösser öffnen":        "Unterwelt",
		"Fallenmechanik":          "Unterwelt",
		"Meucheln":                "Unterwelt",
		"Anführen":                "Sozial",
		"Verführen":               "Sozial",
		"Verstellen":              "Sozial",
		"Beredsamkeit":            "Sozial",
		"Verhören":                "Sozial",
		"Menschenkenntnis":        "Sozial",
		"Lesen von Zauberschrift": "Wissen",
		"Alchimie":                "Wissen",
		"Heilkunde":               "Wissen",
		"Landeskunde":             "Wissen",
		"Zauberkunde":             "Wissen",
	}

	if category, exists := categoryMap[skillName]; exists {
		return category
	}

	// Standard-Fallback
	return "Alltag"
}

func GetDifficulty(skillName string, category string) string {
	// aktuell nur ein Wrapper das die Stantruktur noch keine Fettigkeiten in mehreren Kategorien enthält
	return GetDefaultDifficulty(skillName)
}

// GetDefaultDifficulty gibt die erste (bevorzugte) Schwierigkeit für eine Fertigkeit zurück
func GetDefaultDifficulty(skillName string) string {
	// WICHTIG: Korrespondiert mit getDefaultCategory() - verwendet die Schwierigkeit
	// der ersten (bevorzugten) Kategorie für konsistente Ergebnisse.
	// Schwierigkeitszuordnung basierend auf dem Fertigkeitsnamen und Lerntabellen.md
	difficultyMap := map[string]string{
		// Waffen (aus Waffen-Kategorie)
		"Stichwaffen":          "leicht", // 2 LE
		"Einhandschlagwaffen":  "normal", // 4 LE
		"Zweihandschlagwaffen": "schwer", // 6 LE

		// Körper-Fertigkeiten
		"Geländelauf": "leicht", // Körper: leicht (1 LE)
		"Klettern":    "leicht", // Körper: leicht (1 LE)
		"Schwimmen":   "leicht", // Körper: leicht (1 LE)
		"Balancieren": "leicht", // Körper: leicht (1 LE)
		"Reiten":      "leicht", // Alltag: leicht (1 LE)
		"Tauchen":     "normal", // Körper: normal (1 LE)
		"Akrobatik":   "schwer", // Körper: schwer (2 LE)
		"Athletik":    "schwer", // Körper: schwer (2 LE)
		"Laufen":      "schwer", // Körper: schwer (2 LE)
		"Meditieren":  "schwer", // Körper: schwer (2 LE)

		// Alltag-Fertigkeiten
		"Seilkunst":     "leicht",      // Alltag: leicht (1 LE)
		"Bootfahren":    "leicht",      // Alltag: leicht (1 LE)
		"Glücksspiel":   "leicht",      // Alltag: leicht (1 LE)
		"Wagenlenken":   "leicht",      // Alltag: leicht (1 LE)
		"Musizieren":    "leicht",      // Alltag: leicht (1 LE)
		"Schreiben":     "normal",      // Alltag: normal (1 LE)
		"Sprache":       "normal",      // Alltag: normal (1 LE)
		"Erste Hilfe":   "schwer",      // Alltag: schwer (2 LE)
		"Etikette":      "schwer",      // Alltag: schwer (2 LE)
		"Gerätekunde":   "sehr schwer", // Alltag: sehr schwer (10 LE)
		"Geschäftssinn": "sehr schwer", // Alltag: sehr schwer (10 LE)

		// Freiland-Fertigkeiten
		"Überleben":     "leicht", // Freiland: leicht (1 LE)
		"Naturkunde":    "normal", // Freiland: normal (2 LE)
		"Pflanzenkunde": "normal", // Freiland: normal (2 LE)
		"Tierkunde":     "normal", // Freiland: normal (2 LE)
		"Schleichen":    "schwer", // Freiland: schwer (4 LE)
		"Spurensuche":   "schwer", // Freiland: schwer (4 LE)
		"Tarnen":        "schwer", // Freiland: schwer (4 LE)

		// Unterwelt-Fertigkeiten
		"Gassenwissen":     "leicht", // Unterwelt: leicht (2 LE)
		"Stehlen":          "leicht", // Unterwelt: leicht (2 LE)
		"Fallen entdecken": "normal", // Unterwelt: normal (4 LE)
		"Schlösser öffnen": "normal", // Unterwelt: normal (4 LE)
		"Fallenmechanik":   "schwer", // Unterwelt: schwer (10 LE)
		"Meucheln":         "schwer", // Unterwelt: schwer (10 LE)

		// Sozial-Fertigkeiten
		"Anführen":         "leicht", // Sozial: leicht (2 LE)
		"Verführen":        "leicht", // Sozial: leicht (2 LE)
		"Verstellen":       "leicht", // Sozial: leicht (2 LE)
		"Beredsamkeit":     "normal", // Sozial: normal (2 LE)
		"Verhören":         "normal", // Sozial: normal (2 LE)
		"Menschenkenntnis": "schwer", // Sozial: schwer (4 LE)

		// Wissen-Fertigkeiten
		"Lesen von Zauberschrift": "leicht", // Wissen: leicht (1 LE)
		"Alchimie":                "schwer", // Wissen: schwer (2 LE)
		"Heilkunde":               "schwer", // Wissen: schwer (2 LE)
		"Landeskunde":             "schwer", // Wissen: schwer (2 LE)
		"Zauberkunde":             "schwer", // Wissen: schwer (2 LE)
	}

	if difficulty, exists := difficultyMap[skillName]; exists {
		return difficulty
	}

	// Standard-Fallback
	return "normal"
}

// getDefaultSpellSchool gibt eine Standard-Zauberschule für einen Zauber zurück
func getDefaultSpellSchool(spellName string) string {
	// Vereinfachte Zuordnung von Zauber zu Schulen
	spellSchoolMap := map[string]string{
		"Licht":          "Erschaffen",
		"Sehen":          "Erkennen",
		"Heilen":         "Verändern",
		"Schutz":         "Beherrschen",
		"Unsichtbarkeit": "Verändern",
		"Feuerlanze":     "Zerstören",
		"Eislanze":       "Zerstören",
		"Blitze":         "Zerstören",
		"Schweben":       "Bewegen",
		"Teleportation":  "Bewegen",
	}

	if school, exists := spellSchoolMap[spellName]; exists {
		return school
	}

	// Standard-Fallback
	return "Verändern"
}

// GetClassAbbreviation konvertiert Charakterklassen-Vollnamen zu Abkürzungen
func GetClassAbbreviation(characterClass string) string {
	// Mapping von Vollnamen zu Abkürzungen
	classMap := map[string]string{
		// Abenteurer-Klassen
		"Assassine":    "As",
		"Barbar":       "Bb",
		"Glücksritter": "Gl",
		"Händler":      "Hä",
		"Krieger":      "Kr",
		"Spitzbube":    "Sp",
		"Waldläufer":   "Wa",
		// Zauberer-Klassen
		"Barde":               "Ba",
		"Ordenskrieger":       "Or",
		"Druide":              "Dr",
		"Hexer":               "Hx",
		"Magier":              "Ma",
		"Priester":            "PB", // Standard Priester = Beschützer
		"Priester Beschützer": "PB",
		"Priester Streiter":   "PS",
		"Schamane":            "Sc",
	}

	// Prüfe ob es bereits eine Abkürzung ist
	if len(characterClass) <= 2 {
		return characterClass
	}

	// Suche nach Vollname
	if abbrev, exists := classMap[characterClass]; exists {
		return abbrev
	}

	// Fallback: originale Eingabe zurückgeben
	return characterClass
}
