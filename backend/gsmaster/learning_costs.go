package gsmaster

import (
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

// Initialisierung der Lerntabellen
var learningCosts = &LearningCostsTable{
	// EP-Kosten für 1 Trainingseinheit (TE) nach Klasse und Kategorie
	EPPerTE: map[string]map[string]int{
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
	},

	// EP-Kosten für 1 Lerneinheit (LE) für Zauber
	SpellEPPerLE: map[string]map[string]int{
		"Hx": { // Hexer
			"Beherr":    30,
			"Beweg":     90,
			"Erkenn":    90,
			"Erschaff":  90,
			"Formen":    60,
			"Veränd":    30,
			"Zerstören": 60,
			"Wunder":    0,
			"Dweom":     90,
			"Lied":      0,
		},
		"PS": { // Priester Streiter
			"Beherr":    90,
			"Beweg":     90,
			"Erkenn":    90,
			"Erschaff":  60,
			"Formen":    120,
			"Veränd":    120,
			"Zerstören": 60,
			"Wunder":    30,
			"Dweom":     0,
			"Lied":      0,
		},
	},

	// Lernkosten (LE) basierend auf Schwierigkeit
	BaseLearnCost: map[string]map[string]int{
		"Alltag": {
			"leicht":      1,
			"normal":      1,
			"schwer":      2,
			"sehr_schwer": 10,
		},
		"Freiland": {
			"leicht": 1,
			"normal": 2,
			"schwer": 4,
		},
		"Halbwelt": {
			"leicht":      1,
			"normal":      2,
			"schwer":      2,
			"sehr_schwer": 10,
		},
		"Kampf": {
			"leicht":      1,
			"normal":      2,
			"schwer":      10,
			"sehr_schwer": 10,
		},
		"Körper": {
			"leicht": 1,
			"normal": 1,
			"schwer": 2,
		},
		"Sozial": {
			"leicht": 2,
			"normal": 2,
			"schwer": 4,
		},
		"Unterwelt": {
			"leicht": 2,
			"normal": 2,
			"schwer": 4,
		},
		"Waffen": {
			"leicht": 1,
			"normal": 2,
			"schwer": 10,
		},
		"Wissen": {
			"normal": 2,
			"schwer": 4,
		},
	},

	// Verbesserungskosten (TE) basierend auf Kategorie, Schwierigkeit und aktuellem Wert
	ImprovementCost: map[string]map[string]map[int]int{
		"Alltag": {
			"leicht": {
				13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20,
			},
			"normal": {
				9: 1, 10: 1, 11: 1, 12: 1, 13: 2, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20,
			},
			"schwer": {
				9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50,
			},
			"sehr_schwer": {
				9: 5, 10: 5, 11: 10, 12: 10, 13: 20, 14: 20, 15: 50, 16: 50, 17: 100, 18: 100,
			},
		},
		"Freiland": {
			"leicht": {
				9: 1, 10: 1, 11: 1, 12: 2, 13: 2, 14: 2, 15: 5, 16: 5, 17: 10, 18: 10,
			},
			"normal": {
				9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 20, 17: 50, 18: 50,
			},
			"schwer": {
				9: 5, 10: 5, 11: 10, 12: 10, 13: 20, 14: 20, 15: 50, 16: 50, 17: 100, 18: 100,
			},
		},
		"Halbwelt": {
			"leicht": {
				13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20,
			},
			"normal": {
				9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50,
			},
			"schwer": {
				9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 50, 17: 50, 18: 50,
			},
			"sehr_schwer": {
				9: 5, 10: 10, 11: 20, 12: 20, 13: 30, 14: 50, 15: 80, 16: 80, 17: 100, 18: 100,
			},
		},
		"Kampf": {
			"leicht": {
				13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20,
			},
			"normal": {
				9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50,
			},
			"schwer": {
				9: 5, 10: 10, 11: 20, 12: 20, 13: 30, 14: 50, 15: 80, 16: 80, 17: 100, 18: 100,
			},
			"sehr_schwer": {
				6: 2, 7: 5, 8: 10, 9: 10, 10: 20, 11: 20, 12: 30, 13: 50, 14: 50, 15: 100, 16: 100, 17: 150, 18: 200,
			},
		},
		"Körper": {
			"leicht": {
				13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20,
			},
			"normal": {
				9: 1, 10: 1, 11: 2, 12: 2, 13: 5, 14: 10, 15: 10, 16: 20, 17: 20, 18: 50,
			},
			"schwer": {
				9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50,
			},
		},
		"Sozial": {
			"leicht": {
				9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50,
			},
			"normal": {
				9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 50, 17: 50, 18: 50,
			},
			"schwer": {
				9: 5, 10: 5, 11: 10, 12: 10, 13: 20, 14: 20, 15: 50, 16: 50, 17: 100, 18: 100,
			},
		},
		"Unterwelt": {
			"leicht": {
				9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50,
			},
			"normal": {
				9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 50, 17: 50, 18: 50,
			},
			"schwer": {
				9: 5, 10: 5, 11: 10, 12: 10, 13: 20, 14: 20, 15: 50, 16: 50, 17: 100, 18: 100,
			},
		},
		"Waffen": {
			"leicht": {
				13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20,
			},
			"normal": {
				9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50,
			},
			"schwer": {
				6: 2, 7: 5, 8: 10, 9: 10, 10: 20, 11: 20, 12: 30, 13: 50, 14: 50, 15: 100, 16: 100, 17: 150, 18: 200,
			},
		},
		"Wissen": {
			"normal": {
				9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50,
			},
			"schwer": {
				9: 5, 10: 10, 11: 10, 12: 20, 13: 20, 14: 50, 15: 50, 16: 100, 17: 100, 18: 100,
			},
		},
	},
}

// CalculateDetailedSkillLearningCost berechnet die detaillierten Kosten zum Erlernen einer neuen Fertigkeit
func CalculateDetailedSkillLearningCost(skillName, characterClass string) (*LearnCost, error) {
	// Fertigkeit aus der Datenbank holen
	var skill Skill
	if err := skill.First(skillName); err != nil {
		return nil, errors.New("unbekannte Fertigkeit")
	}

	// Prüfen, ob die Fertigkeit für die Charakterklasse erlaubt ist
	// (Hier könnte man optional eine Prüfung implementieren)

	// Basiskosten für die Fertigkeit basierend auf Kategorie und Schwierigkeit
	categoryMap, ok := learningCosts.BaseLearnCost[skill.Category]
	if !ok {
		return nil, fmt.Errorf("keine Lernkosten für Kategorie: %s", skill.Category)
	}

	// LE-Kosten basierend auf Schwierigkeit
	baseLE, ok := categoryMap[skill.Difficulty]
	if !ok {
		return nil, fmt.Errorf("keine LE-Definition für Schwierigkeit: %s", skill.Difficulty)
	}

	// EP pro TE für die Charakterklasse und Kategorie
	classMap, ok := learningCosts.EPPerTE[characterClass]
	if !ok {
		return nil, fmt.Errorf("keine EP-Kosten für Charakterklasse: %s", characterClass)
	}

	epPerTE, ok := classMap[skill.Category]
	if !ok {
		return nil, fmt.Errorf("keine EP-Kosten für Kategorie %s bei Klasse %s", skill.Category, characterClass)
	}

	// Gesamtkosten berechnen (LE * EP pro TE * 3)
	totalEP := baseLE * (epPerTE * 3)

	// Geldkosten berechnen (20 GS je TE, 200 GS je LE)
	moneyGS := baseLE * 200

	// Ergebnis zurückgeben
	return &LearnCost{
		Stufe: 0, // Neue Fertigkeit hat keinen aktuellen Wert
		LE:    baseLE,
		Ep:    totalEP,
		Money: moneyGS,
	}, nil
}

// CalculateDetailedSkillImprovementCost berechnet die detaillierten Kosten zum Verbessern einer Fertigkeit
func CalculateDetailedSkillImprovementCost(skillName, characterClass string, currentLevel int) (*LearnCost, error) {
	// Fertigkeit aus der Datenbank holen
	var skill Skill
	if err := skill.First(skillName); err != nil {
		return nil, errors.New("unbekannte Fertigkeit")
	}

	// Nächster Level, auf den verbessert werden soll
	nextLevel := currentLevel + 1

	// Kosten für die Verbesserung basierend auf Kategorie, Schwierigkeit und aktuellem Wert
	categoryMap, ok := learningCosts.ImprovementCost[skill.Category]
	if !ok {
		return nil, fmt.Errorf("keine Verbesserungskosten für Kategorie: %s", skill.Category)
	}

	difficultyMap, ok := categoryMap[skill.Difficulty]
	if !ok {
		return nil, fmt.Errorf("keine Verbesserungskosten für Schwierigkeit: %s", skill.Difficulty)
	}

	neededTE, ok := difficultyMap[nextLevel]
	if !ok {
		return nil, fmt.Errorf("kein Eintrag für Verbesserung von %d auf %d", currentLevel, nextLevel)
	}

	// EP pro TE für die Charakterklasse und Kategorie
	classMap, ok := learningCosts.EPPerTE[characterClass]
	if !ok {
		return nil, fmt.Errorf("keine EP-Kosten für Charakterklasse: %s", characterClass)
	}

	epPerTE, ok := classMap[skill.Category]
	if !ok {
		return nil, fmt.Errorf("keine EP-Kosten für Kategorie %s bei Klasse %s", skill.Category, characterClass)
	}

	// Gesamtkosten berechnen (TE * EP pro TE)
	totalEP := neededTE * epPerTE

	// Geldkosten berechnen (20 GS je TE)
	moneyGS := neededTE * 20

	// Ergebnis zurückgeben
	return &LearnCost{
		Stufe: nextLevel,
		LE:    0, // Bei Verbesserung ist LE nicht relevant
		Ep:    totalEP,
		Money: moneyGS,
	}, nil
}

// CalculateDetailedSpellLearningCost berechnet die detaillierten Kosten zum Erlernen eines neuen Zaubers
func CalculateDetailedSpellLearningCost(spellName, characterClass string) (*LearnCost, error) {
	// Zauber aus der Datenbank holen
	var spell Spell
	if err := spell.First(spellName); err != nil {
		return nil, errors.New("unbekannter Zauberspruch")
	}

	// Prüfen, ob der Zauber für die Charakterklasse erlaubt ist
	// (Hier könnte man optional eine Prüfung implementieren)

	// EP pro LE für die Charakterklasse und Zauberschule
	classMap, ok := learningCosts.SpellEPPerLE[characterClass]
	if !ok {
		return nil, fmt.Errorf("keine EP-Kosten für Charakterklasse: %s", characterClass)
	}

	epPerLE, ok := classMap[spell.Category]
	if !ok {
		return nil, fmt.Errorf("keine EP-Kosten für Zauberschule %s bei Klasse %s", spell.Category, characterClass)
	}

	// LE basierend auf Zauberstufe
	// Für dieses Beispiel nehmen wir an, dass die Stufe direkt mit den LE korreliert
	// Eine genauere Implementierung würde eine Mapping-Tabelle verwenden
	neededLE := spell.Stufe

	// Gesamtkosten berechnen (LE * EP pro LE)
	totalEP := neededLE * epPerLE

	// Geldkosten berechnen (100 GS je LE)
	moneyGS := neededLE * 100

	// Ergebnis zurückgeben
	return &LearnCost{
		Stufe: spell.Stufe,
		LE:    neededLE,
		Ep:    totalEP,
		Money: moneyGS,
	}, nil
}
