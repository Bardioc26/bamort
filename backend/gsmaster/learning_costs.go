package gsmaster

import (
	"bamort/database"
	"bamort/models"
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

// SkillCategoryOption definiert eine Kategorie-Schwierigkeit-Kombination für eine Fertigkeit
type SkillCategoryOption struct {
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
	LE         int    `json:"le"`
}

func GetClassAbbreviationNewSystem(characterClass string) string {
	// Try to find by code first (e.g., "Kr" -> "Kr")
	var charClass models.CharacterClass
	if err := charClass.FirstByName(characterClass); err == nil {
		return charClass.Code
	}

	// Try to find by name (e.g., "Krieger" -> "Kr")
	if err := database.DB.Where("name = ?", characterClass).First(&charClass).Error; err == nil {
		return charClass.Code
	}
	return ""
}
