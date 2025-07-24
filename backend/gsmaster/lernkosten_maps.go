package gsmaster

import "fmt"

// DifficultyData enthält Skills und Trainingskosten für eine Schwierigkeitsstufe
type DifficultyData struct {
	LearnCost  int         `json:"learncosts"`
	Skills     []string    `json:"skills"`
	TrainCosts map[int]int `json:"traincosts"`
}

// LearningCostsTable strukturiert die Daten aus Lerntabellen.md
type LearningCostsTable2 struct {
	// EP-Kosten für 1 Trainingseinheit (TE) pro Charakterklasse und Fertigkeitskategorie
	EPPerTE map[string]map[string]int

	// EP-Kosten für 1 Lerneinheit (LE) für Zauber pro Charakterklasse und Zauberschule
	SpellEPPerLE map[string]map[string]int

	// LE-Kosten für Fertigkeiten basierend auf Schwierigkeit
	BaseLearnCost map[string]map[string]int

	// TE-Kosten für Verbesserungen basierend auf Kategorie, Schwierigkeit und aktuellem Wert
	ImprovementCost map[string]map[string]DifficultyData
}

// learningCostsData enthält alle statischen Lerntabellendaten
var learningCostsData = &LearningCostsTable2{
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

	// TE-Kosten für Verbesserungen basierend auf Kategorie, Schwierigkeit und aktuellem Wert
	ImprovementCost: map[string]map[string]DifficultyData{
		"Alltag": {
			"leicht": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Klettern", "Reiten", "Seilkunst", "Bootfahren", "Glücksspiel", "Wagenlenken", "Musizieren"},
				TrainCosts: map[int]int{9: 0, 10: 0, 11: 0, 12: 0, 13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20},
			},
			"normal": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Schreiben", "Sprache"},
				TrainCosts: map[int]int{9: 1, 10: 1, 11: 1, 12: 1, 13: 2, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20},
			},
			"schwer": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Erste Hilfe", "Etikette"},
				TrainCosts: map[int]int{9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50},
			},
			"sehr schwer": DifficultyData{
				LearnCost:  10,
				Skills:     []string{"Gerätekunde", "Geschäftssinn"},
				TrainCosts: map[int]int{9: 5, 10: 5, 11: 10, 12: 10, 13: 20, 14: 20, 15: 50, 16: 50, 17: 100, 18: 100},
			},
		},
		"Freiland": {
			"leicht": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Überleben"},
				TrainCosts: map[int]int{9: 1, 10: 1, 11: 1, 12: 2, 13: 2, 14: 2, 15: 5, 16: 5, 17: 10, 18: 10},
			},
			"normal": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Naturkunde", "Pflanzenkunde", "Tierkunde"},
				TrainCosts: map[int]int{9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 20, 17: 50, 18: 50},
			},
			"schwer": DifficultyData{
				LearnCost:  4,
				Skills:     []string{"Schleichen", "Spurensuche", "Tarnen"},
				TrainCosts: map[int]int{9: 5, 10: 5, 11: 10, 12: 10, 13: 20, 14: 20, 15: 50, 16: 50, 17: 100, 18: 100},
			},
		},
		"Halbwelt": {
			"leicht": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Klettern", "Glücksspiel", "Balancieren"},
				TrainCosts: map[int]int{9: 0, 10: 0, 11: 0, 12: 0, 13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20},
			},
			"normal": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Akrobatik"},
				TrainCosts: map[int]int{9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50},
			},
			"schwer": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Gassenwissen", "Stehlen"},
				TrainCosts: map[int]int{9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 50, 17: 50, 18: 50},
			},
			"sehr schwer": DifficultyData{
				LearnCost:  10,
				Skills:     []string{"Betäuben"},
				TrainCosts: map[int]int{9: 5, 10: 10, 11: 20, 12: 20, 13: 30, 14: 50, 15: 80, 16: 80, 17: 100, 18: 100},
			},
		},
		"Kampf": {
			"leicht": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Reiten"},
				TrainCosts: map[int]int{6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0, 12: 0, 13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20},
			},
			"normal": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Anführen", "Athletik"},
				TrainCosts: map[int]int{6: 0, 7: 0, 8: 0, 9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50},
			},
			"schwer": DifficultyData{
				LearnCost:  10,
				Skills:     []string{"Betäuben"},
				TrainCosts: map[int]int{6: 0, 7: 0, 8: 0, 9: 5, 10: 10, 11: 20, 12: 20, 13: 30, 14: 50, 15: 80, 16: 80, 17: 100, 18: 100},
			},
			"sehr schwer": DifficultyData{
				LearnCost:  10,
				Skills:     []string{},
				TrainCosts: map[int]int{6: 2, 7: 5, 8: 10, 9: 10, 10: 20, 11: 20, 12: 30, 13: 50, 14: 50, 15: 100, 16: 100, 17: 150, 18: 200},
			},
		},
		"Körper": {
			"leicht": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Geländelauf", "Klettern", "Schwimmen", "Balancieren"},
				TrainCosts: map[int]int{9: 0, 10: 0, 11: 0, 12: 0, 13: 1, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20},
			},
			"normal": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Tauchen"},
				TrainCosts: map[int]int{9: 1, 10: 1, 11: 2, 12: 2, 13: 5, 14: 10, 15: 10, 16: 20, 17: 20, 18: 50},
			},
			"schwer": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Akrobatik", "Athletik", "Laufen", "Meditieren"},
				TrainCosts: map[int]int{9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50},
			},
		},
		"Sozial": {
			"leicht": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Anführen", "Verführen", "Verstellen", "Etikette"},
				TrainCosts: map[int]int{9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50},
			},
			"normal": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Gassenwissen", "Beredsamkeit", "Verhören"},
				TrainCosts: map[int]int{9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 50, 17: 50, 18: 50},
			},
			"schwer": DifficultyData{
				LearnCost:  4,
				Skills:     []string{"Menschenkenntnis"},
				TrainCosts: map[int]int{9: 5, 10: 5, 11: 10, 12: 10, 13: 20, 14: 20, 15: 50, 16: 50, 17: 100, 18: 100},
			},
		},
		"Unterwelt": {
			"leicht": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Gassenwissen", "Stehlen"},
				TrainCosts: map[int]int{9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 50, 17: 50, 18: 50},
			},
			"normal": DifficultyData{
				LearnCost:  4,
				Skills:     []string{"Schleichen", "Spurensuche", "Tarnen", "Fallen entdecken", "Schlösser öffnen"},
				TrainCosts: map[int]int{9: 5, 10: 5, 11: 10, 12: 10, 13: 20, 14: 20, 15: 50, 16: 50, 17: 100, 18: 100},
			},
			"schwer": DifficultyData{
				LearnCost:  10,
				Skills:     []string{"Fallenmechanik", "Meucheln", "Menschenkenntnis"},
				TrainCosts: map[int]int{9: 5, 10: 10, 11: 20, 12: 20, 13: 30, 14: 50, 15: 80, 16: 80, 17: 100, 18: 100},
			},
		},
		"Waffen": {
			"leicht": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Stichwaffen"},
				TrainCosts: map[int]int{6: 1, 7: 1, 8: 1, 9: 2, 10: 2, 11: 5, 12: 10, 13: 20, 14: 50, 15: 100, 16: 100, 17: 150, 18: 150},
			},
			"normal": DifficultyData{
				LearnCost:  4,
				Skills:     []string{"Einhandschlagwaffen"},
				TrainCosts: map[int]int{6: 1, 7: 1, 8: 2, 9: 2, 10: 5, 11: 10, 12: 20, 13: 50, 14: 50, 15: 100, 16: 150, 17: 150, 18: 200},
			},
			"schwer": DifficultyData{
				LearnCost:  6,
				Skills:     []string{"Zweihandschlagwaffen"},
				TrainCosts: map[int]int{6: 1, 7: 2, 8: 2, 9: 5, 10: 5, 11: 10, 12: 20, 13: 50, 14: 100, 15: 150, 16: 200, 17: 300, 18: 300},
			},
			"sehr schwer": DifficultyData{
				LearnCost:  8,
				Skills:     []string{},
				TrainCosts: map[int]int{6: 1, 7: 2, 8: 2, 9: 5, 10: 10, 11: 20, 12: 50, 13: 100, 14: 150, 15: 200, 16: 300, 17: 300, 18: 400},
			},
		},
		"Wissen": {
			"leicht": DifficultyData{
				LearnCost:  1,
				Skills:     []string{"Lesen von Zauberschrift", "Schreiben", "Sprache"},
				TrainCosts: map[int]int{9: 1, 10: 1, 11: 1, 12: 1, 13: 2, 14: 2, 15: 5, 16: 10, 17: 10, 18: 20},
			},
			"normal": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Erste Hilfe", "Meditieren"},
				TrainCosts: map[int]int{9: 2, 10: 2, 11: 5, 12: 5, 13: 10, 14: 10, 15: 20, 16: 20, 17: 50, 18: 50},
			},
			"schwer": DifficultyData{
				LearnCost:  2,
				Skills:     []string{"Alchimie", "Heilkunde", "Landeskunde", "Zauberkunde", "Naturkunde", "Pflanzenkunde", "Tierkunde"},
				TrainCosts: map[int]int{9: 2, 10: 5, 11: 5, 12: 10, 13: 10, 14: 20, 15: 20, 16: 20, 17: 50, 18: 50},
			},
		},
		"Schilde und Parierwaﬀen": {
			"normal": DifficultyData{
				LearnCost:  0, // Not defined in BaseLearnCost, using 0 as default
				Skills:     []string{},
				TrainCosts: map[int]int{2: 1, 3: 2, 4: 10, 5: 30, 6: 50, 7: 100, 8: 150},
			},
		},
	},
}

func GetSkillCategory(skillName string) string {

	for category, difficulties := range learningCostsData.ImprovementCost {
		for _, data := range difficulties {
			if contains(data.Skills, skillName) {
				return category
			}
		}
	}
	return "Unbekannt"
}

func GetSkillDifficulty(category string, skillName string) string {
	// Wenn eine Kategorie angegeben ist, suche nur in dieser Kategorie
	if category != "" {
		difficulties, ok := learningCostsData.ImprovementCost[category]
		if !ok {
			return "Unbekannt" // Kategorie nicht gefunden
		}
		for difficulty, data := range difficulties {
			if contains(data.Skills, skillName) {
				return difficulty
			}
		}
		return "Unbekannt" // Skill in der angegebenen Kategorie nicht gefunden
	}

	// Wenn keine Kategorie angegeben ist, durchsuche alle Kategorien und gib das erste Vorkommen zurück
	for _, difficulties := range learningCostsData.ImprovementCost {
		for difficulty, data := range difficulties {
			if contains(data.Skills, skillName) {
				return difficulty
			}
		}
	}
	return "Unbekannt"
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

//### End of Helper functions ###

func CalcSkillLernCost(costResult *SkillCostResultNew, reward *string) error {
	// Berechne die Lernkosten basierend auf den aktuellen Werten im costResult
	// Hier sollte die Logik zur Berechnung der Lernkosten implementiert werden
	//Finde EP kosten für die Kategorie für die Charakterklasse aus learningCostsData.EPPerTE
	// Konvertiere Vollnamen der Charakterklasse zu Abkürzungen falls nötig
	//classKey := getClassAbbreviation(costResult.CharacterClass)
	classKey := costResult.CharacterClass

	epPerTE, ok := learningCostsData.EPPerTE[classKey][costResult.Category]
	if !ok {
		return fmt.Errorf("EP-Kosten für Kategorie '%s' und Klasse '%s' nicht gefunden", costResult.Category, costResult.CharacterClass)
	}
	// finde LE für den Skill aufgrund der Kategorie und schwierigkeit aus DifficultyData
	learnCost, ok := learningCostsData.ImprovementCost[costResult.Category][costResult.Difficulty]
	if !ok {
		return fmt.Errorf("Lernkosten für Kategorie '%s' und Schwierigkeit '%s' nicht gefunden", costResult.Category, costResult.Difficulty)
	}
	costResult.LE = learnCost.LearnCost
	costResult.EP = epPerTE * costResult.LE * 3
	costResult.GoldCost = costResult.LE * 200 // Beispiel: 200 Gold pro LE

	// Apply reward logic
	if reward != nil && *reward == "noGold" {
		costResult.GoldCost = 0 // Keine Goldkosten für diese Belohnung
	}

	return nil
}

// CalcSkillImproveCost berechnet die Kosten für die Verbesserung einer Fertigkeit
func CalcSkillImproveCost(costResult *SkillCostResultNew, currentLevel int, reward *string) error {
	// Für Skill-Verbesserung könnten die Kosten vom aktuellen Level abhängen

	//Finde EP kosten für die Kategorie für die Charakterklasse aus learningCostsData.EPPerTE
	//classKey := getClassAbbreviation(costResult.CharacterClass)
	classKey := costResult.CharacterClass

	epPerTE, ok := learningCostsData.EPPerTE[classKey][costResult.Category]
	if !ok {
		return fmt.Errorf("EP-Kosten für Kategorie '%s' und Klasse '%s' nicht gefunden", costResult.Category, costResult.CharacterClass)
	}

	diffData := learningCostsData.ImprovementCost[costResult.Category][costResult.Difficulty]
	trainCost := diffData.TrainCosts[currentLevel+1]
	if costResult.PPUsed > 0 {
		trainCost -= costResult.PPUsed // Wenn PP verwendet werden, setze die Kosten auf die PP
	}
	// Apply reward logic
	costResult.LE = trainCost
	costResult.EP = epPerTE * trainCost
	costResult.GoldCost = trainCost * 20 // Beispiel: 20 Gold pro TE

	if reward != nil && *reward == "halveep" {
		costResult.EP = costResult.EP / 2 // Halbiere die EP-Kosten für diese Belohnung
	}
	if reward != nil && *reward == "halveepnoGold" {
		costResult.GoldCost = 0           // Keine Goldkosten für diese Belohnung
		costResult.EP = costResult.EP / 2 // Halbiere die EP-Kosten für diese Belohnung
	}

	return nil
}

// CalcSpellLernCost berechnet die Kosten für das Erlernen eines Zaubers
func CalcSpellLernCost(costResult *SkillCostResultNew, reward *string) error {
	// Für Zauber verwenden wir eine ähnliche Logik wie für Skills
	// TODO: Implementiere spezifische Zauber-Kostenlogik wenn verfügbar
	// Für jetzt verwenden wir die gleiche Logik wie für Skills
	return CalcSkillLernCost(costResult, reward)
}
