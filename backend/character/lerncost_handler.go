package character

import (
	"bamort/gsmaster"
	"bamort/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SkillCostRequest struct {
	Name         string `json:"name" binding:"required"`                          // Name der Fertigkeit
	CurrentLevel int    `json:"current_level,omitempty"`                          // Aktueller Wert (nur für Verbesserung)
	Type         string `json:"type" binding:"required,oneof=skill spell weapon"` // 'skill', 'spell' oder 'weapon'
	Action       string `json:"action" binding:"required,oneof=learn improve"`    // 'learn' oder 'improve'
	TargetLevel  int    `json:"target_level,omitempty"`                           // Zielwert (optional, für Kostenberechnung bis zu einem bestimmten Level)
	UsePP        int    `json:"use_pp,omitempty"`                                 // Anzahl der zu verwendenden Praxispunkte

	// Belohnungsoptionen
	Reward *RewardOptions `json:"reward,omitempty"` // Belohnungsoptionen
}

// RewardOptions definiert die verschiedenen Belohnungsmöglichkeiten
type RewardOptions struct {
	Type         string `json:"type,omitempty" binding:"omitempty,oneof=free_learning free_spell_learning half_ep_improvement gold_for_ep"` // Art der Belohnung
	UseGoldForEP bool   `json:"use_gold_for_ep,omitempty"`                                                                                  // 10 GS statt 1 EP verwenden
	MaxGoldEP    int    `json:"max_gold_ep,omitempty"`                                                                                      // Maximale EP die durch Gold ersetzt werden (automatisch: Hälfte der Kosten)
}

type SkillCostResponse struct {
	*models.LearnCost
	SkillName    string `json:"skill_name"`
	SkillType    string `json:"skill_type"`
	Action       string `json:"action"`
	CharacterID  uint   `json:"character_id"`
	CurrentLevel int    `json:"current_level,omitempty"`
	TargetLevel  int    `json:"target_level,omitempty"`
	Category     string `json:"category,omitempty"`
	Difficulty   string `json:"difficulty,omitempty"`
	CanAfford    bool   `json:"can_afford"`
	Notes        string `json:"notes,omitempty"`
	PPUsed       int    `json:"pp_used,omitempty"`      // Anzahl der verwendeten Praxispunkte
	PPAvailable  int    `json:"pp_available,omitempty"` // Verfügbare Praxispunkte für diese Kategorie

	// Belohnungsdetails
	RewardApplied      string            `json:"reward_applied,omitempty"`       // Art der angewendeten Belohnung
	OriginalCostStruct *models.LearnCost `json:"original_cost_struct,omitempty"` // Ursprüngliche Kosten ohne Belohnung
	Savings            *models.LearnCost `json:"savings,omitempty"`              // Ersparnisse durch Belohnung
	GoldUsedForEP      int               `json:"gold_used_for_ep,omitempty"`     // Gold das für EP verwendet wurde
	PPReduction        int               `json:"pp_reduction,omitempty"`         // Reduktion der Kosten durch PP
	OriginalCost       int               `json:"original_cost,omitempty"`        // Ursprüngliche Kosten (vor PP-Reduktion)
	FinalCost          int               `json:"final_cost,omitempty"`           // Endgültige Kosten (nach PP-Reduktion)
}

type MultiLevelCostResponse struct {
	SkillName      string              `json:"skill_name"`
	SkillType      string              `json:"skill_type"`
	CharacterID    uint                `json:"character_id"`
	CurrentLevel   int                 `json:"current_level"`
	TargetLevel    int                 `json:"target_level"`
	LevelCosts     []SkillCostResponse `json:"level_costs"`
	TotalCost      *models.LearnCost   `json:"total_cost"`
	CanAffordTotal bool                `json:"can_afford_total"`
}

// GetLernCostNewSystem verwendet das neue Datenbank-Lernkosten-System
// und produziert die gleichen Ergebnisse wie GetLernCost.
//
// Wie es funktionert:
// - Für "learn" Aktion: Nur eine Berechnung, da Lernkosten einmalig sind
// - Für "improve" Aktion: Berechne für jedes Level von current+1 bis 18
// - Berücksichtigt Praxispunkte (PP) und Gold-für-EP Konvertierung
// - Wendet Belohnungen an (kostenloses Lernen, halbe EP, etc.)
// - Gibt eine Liste von Kosten pro Level zurück
// Schritt für Schritt:
//  1. Hole Charakter und Klassenabkürzung
//  2. Normalisiere Fertigkeits-/Zaubername
//  3. Initialisiere einzusetzende/verbleibende PP und Gold
//  4. Je nach Aktion:
//     4.1 "learn": Hole Lerninformationen und berechne Kosten
//     4.1.1 "spell": Hole Zauber-Lerninformationen und berechne Kosten
//     4.1.2 "skill": Hole Fertigkeits-Lerninformationen und berechne Kosten
//     4.2 "improve": Für jedes Level, hole Verbesserungsinformationen und berechne Kosten
//     (nur Fertigkeiten, keine Zauber)
//
// 5. Wende Belohnungen an
// 6. Wende PP und Gold-für-EP an
// 7. Sammle Ergebnisse und sende als JSON-Antwort
// GetLernCostNewSystem godoc
// @Summary Calculate learning costs
// @Description Calculates the experience point cost for learning or improving a skill/spell
// @Tags Characters
// @Accept json
// @Produce json
// @Param learn_request body object{character_id=int,skill_id=int,spell_id=int,current_value=int,reward_type=string} true "Learning cost request"
// @Success 200 {object} object "Learning cost calculation result"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/characters/lerncost-new [post]
func GetLernCostNewSystem(c *gin.Context) {
	// Request-Parameter abrufen
	var request gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// 1. Hole Charakter
	charID := fmt.Sprintf("%d", request.CharId)
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(character.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviationNewSystem(character.Typ)
	} else {
		characterClass = character.Typ
	}

	//2. Normalize skill/spell name (trim whitespace, proper case)
	skillName := strings.TrimSpace(request.Name)

	var response []gsmaster.SkillCostResultNew
	remainingPP := request.UsePP
	remainingGold := request.UseGold

	// Für "learn" Aktion: nur eine Berechnung, da Lernkosten einmalig sind
	if request.Action == "learn" {
		// 4.1 "learn": Hole Lerninformationen und berechne Kosten
		if request.Type == "spell" {
			// 4.1.1 "spell": Hole Zauber-Lerninformationen und berechne Kosten
			// Spell learning logic
			spellInfo, err := models.GetSpellLearningInfoNewSystem(skillName, characterClass)
			if err != nil {
				respondWithError(c, http.StatusBadRequest, fmt.Sprintf("Zauber '%s' nicht gefunden oder nicht für Klasse '%s' verfügbar: %v", skillName, characterClass, err))
				return
			}

			levelResult := gsmaster.SkillCostResultNew{
				CharacterID:    charID,
				CharacterClass: characterClass,
				SkillName:      skillName,
				Category:       spellInfo.SchoolName,
				Difficulty:     fmt.Sprintf("Stufe %d", spellInfo.SpellLevel),
				TargetLevel:    1, // Lernkosten sind für das Erlernen des Zaubers (Level 1)
			}

			err = calculateSpellLearnCostNewSystem(&request, &levelResult, &remainingPP, &remainingGold, spellInfo)
			if err != nil {
				respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
				return
			}

			response = append(response, levelResult)
		} else {
			// 4.1.2 "skill": Hole Fertigkeits-Lerninformationen und berechne Kosten
			skillInfo, err := models.GetSkillCategoryAndDifficultyNewSystem(skillName, characterClass)
			if err != nil {
				respondWithError(c, http.StatusBadRequest, fmt.Sprintf("Fertigkeit '%s' nicht gefunden oder nicht für Klasse '%s' verfügbar: %v", skillName, characterClass, err))
				return
			}

			levelResult := gsmaster.SkillCostResultNew{
				CharacterID:    charID,
				CharacterClass: characterClass,
				SkillName:      skillName,
				Category:       skillInfo.CategoryName,
				Difficulty:     skillInfo.DifficultyName,
				TargetLevel:    1, // Lernkosten sind für das Erlernen der Fertigkeit (Level 1)
			}

			err = calculateSkillLearnCostNewSystem(&request, &levelResult, &remainingPP, &remainingGold, skillInfo)
			if err != nil {
				respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
				return
			}

			response = append(response, levelResult)
		}
	} else {
		// Für "improve" Aktion: berechne für jedes Level von current+1 bis 18
		// Improvement only works on skills, not spells
		skillInfo, err := models.GetSkillCategoryAndDifficultyNewSystem(skillName, characterClass)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, fmt.Sprintf("Fertigkeit '%s' nicht gefunden oder nicht für Klasse '%s' verfügbar: %v", skillName, characterClass, err))
			return
		}

		for i := request.CurrentLevel + 1; i <= 18; i++ {
			levelResult := gsmaster.SkillCostResultNew{
				CharacterID:    charID,
				CharacterClass: characterClass,
				SkillName:      skillName,
				Category:       skillInfo.CategoryName,
				Difficulty:     skillInfo.DifficultyName,
				TargetLevel:    i,
			}

			err := CalculateSkillImproveCostNewSystem(&request, &levelResult, i, &remainingPP, &remainingGold, skillInfo)
			if err != nil {
				respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
				return
			}
			// für die nächste Runde die PP und Gold reduzieren die zum Lernen genutzt werden sollen
			if levelResult.PPUsed > 0 {
				request.UsePP -= levelResult.PPUsed
				// Sicherstellen, dass PP nicht unter 0 fallen
				if request.UsePP < 0 {
					request.UsePP = 0
				}
			}
			if levelResult.GoldUsed > 0 {
				request.UseGold -= levelResult.GoldUsed
				// Sicherstellen, dass Gold nicht unter 0 fällt
				if request.UseGold < 0 {
					request.UseGold = 0
				}
			}
			response = append(response, levelResult)
		}
	}

	c.JSON(http.StatusOK, response)
}

// calculateCostNewSystem berechnet die Kosten für ein Level mit dem neuen Datenbank-System
func CalculateSkillImproveCostNewSystem(request *gsmaster.LernCostRequest, result *gsmaster.SkillCostResultNew, targetLevel int, remainingPP *int, remainingGold *int, skillInfo *models.SkillLearningInfo) error {
	// 1. Hole die TE-Kosten für die Verbesserung vom aktuellen Level
	teRequired, err := models.GetImprovementCost(skillInfo.SkillName, skillInfo.CategoryName, skillInfo.DifficultyName, targetLevel)
	if err != nil {
		return fmt.Errorf("verbesserungskosten nicht gefunden für %s (Level %d): %v", skillInfo.SkillName, targetLevel, err)
	}

	// 2. Hole die EP-Kosten pro TE für diese Klasse und Kategorie
	if skillInfo.EPPerTE == 0 {
		epPerTE, err := models.GetEPPerTEForClassAndCategory(result.CharacterClass, skillInfo.CategoryName)
		if err != nil {
			return fmt.Errorf("EP-Kosten pro TE nicht gefunden für Klasse %s, Kategorie %s: %v", result.CharacterClass, skillInfo.CategoryName, err)
		}
		skillInfo.EPPerTE = epPerTE
	}

	// 3. Setze die ursprünglichen TE-Kosten
	trainCost := teRequired

	// 4. Anwenden von Praxispunkten (PP) - Exakt wie im alten System
	ppUsed := 0
	if *remainingPP > 0 {
		if trainCost < *remainingPP {
			ppUsed = trainCost // Maximal so viele PP verwenden wie TE benötigt werden
			trainCost = 0      // Wenn PP alle TE abdecken, setze trainCost auf 0
		} else if *remainingPP > 0 {
			ppUsed = *remainingPP // Verwende alle verfügbaren PP
			trainCost -= ppUsed   // Reduziere TE um verwendete PP
		}

		result.PPUsed = ppUsed
		*remainingPP -= ppUsed

		if *remainingPP < 0 {
			*remainingPP = 0
		}
	}

	// 5. Berechne Kosten nach PP-Anwendung (wie im alten System)
	result.LE = trainCost
	result.EP = skillInfo.EPPerTE * trainCost
	result.GoldCost = trainCost * 20 // Wie im alten System: 20 Gold pro TE

	// 6. Anwenden von Belohnungen
	if request.Reward != nil {
		applyRewardNewSystem(result, request.Reward, result.EP)
	}

	// 7. Anwenden von Gold für EP (falls verfügbar) - Beschränkt auf EP/2
	goldUsed := 0
	if *remainingGold > 0 {
		// 10 Gold = 1 EP, aber maximal EP/2 kann durch Gold ersetzt werden
		maxEPFromGold := result.EP / 2
		epFromGold := *remainingGold / 10

		if epFromGold > maxEPFromGold {
			// Beschränke auf maximal EP/2
			epFromGold = maxEPFromGold
			goldUsed = epFromGold * 10
		} else {
			// Verwende das verfügbare Gold
			goldUsed = *remainingGold
		}

		// Reduziere EP um die durch Gold ersetzte Menge
		result.EP -= epFromGold
		result.GoldUsed = goldUsed
		*remainingGold -= goldUsed

		if *remainingGold < 0 {
			*remainingGold = 0
		}
	}

	return nil
}

// calculateSkillLearnCostNewSystem berechnet die Kosten für das Erlernen einer Fertigkeit (Action: "learn", Type: "skill")
func calculateSkillLearnCostNewSystem(request *gsmaster.LernCostRequest, result *gsmaster.SkillCostResultNew, remainingPP *int, remainingGold *int, skillInfo *models.SkillLearningInfo) error {
	// 1. Hole die EP-Kosten pro TE für diese Klasse und Kategorie
	epPerTE, err := models.GetEPPerTEForClassAndCategory(result.CharacterClass, skillInfo.CategoryName)
	if err != nil {
		return fmt.Errorf("EP-Kosten pro TE nicht gefunden für Klasse %s, Kategorie %s: %v", result.CharacterClass, skillInfo.CategoryName, err)
	}

	// 2. Verwende die Lernkosten (LE) direkt aus der skillInfo - diese enthält bereits alle benötigten Informationen
	learnCost := skillInfo.LearnCost

	// 3. Berechne Kosten nach Lernregeln (wie im alten System)
	result.LE = learnCost
	result.EP = epPerTE * result.LE * 3 // Faktor 3 beim Lernen!
	result.GoldCost = result.LE * 200   // 200 Gold pro LE (nicht 20 Gold pro TE)

	// 4. Anwenden von Belohnungen
	if request.Reward != nil {
		applyRewardNewSystem(result, request.Reward, result.EP)
	}

	// 5. Für Skill-Lernen: Keine PP oder Gold-für-EP Anwendung erlaubt (im Gegensatz zu Spell-Lernen)
	// PP und Gold bleiben unverändert, da sie bei Skill-Lernen nicht verwendet werden

	return nil
}

// applyRewardNewSystem wendet Belohnungen auf die Kosten an (neues System)
func applyRewardNewSystem(result *gsmaster.SkillCostResultNew, reward *string, originalEP int) {
	if reward == nil || *reward == "" {
		return
	}

	switch *reward {
	case "noGold":
		// Kostenlose Fertigkeiten: Nur Geld ist 0, EP bleiben
		result.GoldCost = 0

	case "halveep":
		// Halbe EP für Verbesserungen
		result.EP = result.EP / 2

	case "halveepnoGold":
		// Halbe EP und kein Gold
		result.EP = result.EP / 2
		result.GoldCost = 0

	case "default":
		// Keine Änderungen
		break

	default:
		// Unbekannte Belohnung - ignorieren
		break
	}
}

// calculateSpellLearnCostNewSystem berechnet die Kosten für das Erlernen eines Zaubers (Action: "learn", Type: "spell")
func calculateSpellLearnCostNewSystem(request *gsmaster.LernCostRequest, result *gsmaster.SkillCostResultNew, remainingPP *int, remainingGold *int, spellInfo *models.SpellLearningInfo) error {
	// 1. Setze die grundlegenden Zauber-Informationen
	result.Category = spellInfo.SchoolName
	result.Difficulty = fmt.Sprintf("Stufe %d", spellInfo.SpellLevel)

	// 2. Berechne die LE-Kosten basierend auf der Zaubergrad
	leRequired := spellInfo.LERequired

	// 3. Anwenden von PP (Practice Points): 1 PP = 1 LE Reduktion (bei Zauber-Lernen erlaubt)
	ppUsed := 0
	if *remainingPP > 0 {
		if leRequired <= *remainingPP {
			ppUsed = leRequired // Maximal so viele PP verwenden wie LE benötigt werden
			leRequired = 0      // Wenn PP alle LE abdecken
		} else {
			ppUsed = *remainingPP // Verwende alle verfügbaren PP
			leRequired -= ppUsed  // Reduziere LE um verwendete PP
		}

		result.PPUsed = ppUsed
		*remainingPP -= ppUsed

		if *remainingPP < 0 {
			*remainingPP = 0
		}
	}

	// 4. Setze die finalen LE-Kosten
	result.LE = leRequired

	// 5. Berechne EP-Kosten basierend auf LE und EP-pro-LE für diese Klasse/Schule
	result.EP = result.LE * spellInfo.EPPerLE

	// 6. Berechne Gold-Kosten (Beispiel: 100 Gold pro LE wie im alten System)
	result.GoldCost = result.LE * 100

	// 7. Anwenden von Belohnungen (spruchrolle, halveep, etc.)
	if request.Reward != nil {
		applySpellRewardNewSystem(result, request.Reward)
	}

	// 8. Gold-für-EP Konvertierung für Zauber-Lernen (erlaubt)
	goldUsed := 0
	if *remainingGold > 0 {
		// 10 Gold = 1 EP, aber maximal EP/2 kann durch Gold ersetzt werden
		maxEPFromGold := result.EP / 2
		epFromGold := *remainingGold / 10

		if epFromGold > maxEPFromGold {
			// Beschränke auf maximal EP/2
			epFromGold = maxEPFromGold
			goldUsed = epFromGold * 10
		} else {
			// Verwende das verfügbare Gold
			goldUsed = *remainingGold
		}

		// Reduziere EP um die durch Gold ersetzte Menge
		result.EP -= epFromGold
		result.GoldUsed = goldUsed
		*remainingGold -= goldUsed

		if *remainingGold < 0 {
			*remainingGold = 0
		}
	}

	return nil
}

// applySpellRewardNewSystem wendet zauber-spezifische Belohnungen an
func applySpellRewardNewSystem(result *gsmaster.SkillCostResultNew, reward *string) {
	if reward == nil || *reward == "" {
		return
	}

	switch *reward {
	case "spruchrolle":
		// Spruchrolle: 20 Gold für jeden Versuch und 1/3 EP-Kosten bei Erfolg
		result.GoldCost = 20
		result.EP = result.EP / 3

	case "halveep":
		// Halbe EP für Zauber-Lernen
		result.EP = result.EP / 2

	case "halveepnoGold":
		// Halbe EP und kein Gold
		result.EP = result.EP / 2
		result.GoldCost = 0

	case "noGold":
		// Nur Geld ist 0, EP bleiben
		result.GoldCost = 0

	case "default":
		// Keine Änderungen
		break

	default:
		// Unbekannte Belohnung - ignorieren
		break
	}
}

// Helper function to get current skill level from character
func getCurrentSkillLevel(character *models.Char, skillName, skillType string) int {
	switch skillType {
	case "skill":
		for _, skill := range character.Fertigkeiten {
			if skill.Name == skillName {
				return skill.Fertigkeitswert
			}
		}
	case "weapon":
		for _, skill := range character.Waffenfertigkeiten {
			if skill.Name == skillName {
				return skill.Fertigkeitswert
			}
		}
	case "spell":
		// Spells don't have levels in the same way
		return 0
	}
	return -1
}

// applyReward wendet Belohnungen auf die Kosten an
func applyReward(cost *models.LearnCost, request *SkillCostRequest) *models.LearnCost {
	if request.Reward == nil || request.Reward.Type == "" {
		return cost
	}

	newCost := *cost // Kopie der ursprünglichen Kosten

	switch request.Reward.Type {
	case "free_learning":
		// Kostenlose Fertigkeiten: Nur Geld ist 0, EP/LE bleiben
		if request.Type == "skill" && request.Action == "learn" {
			newCost.Money = 0
		}

	case "free_spell_learning":
		// Kostenlose Zauber: Nur LE ist 0, EP/Geld bleiben
		if request.Type == "spell" && request.Action == "learn" {
			newCost.LE = 0
		}

	case "half_ep_improvement":
		// Halbe EP für Verbesserungen
		if request.Action == "improve" {
			newCost.Ep = newCost.Ep / 2
		}

	case "gold_for_ep":
		// Gold statt EP verwenden (10 GS = 1 EP)
		if request.Reward.UseGoldForEP && newCost.Ep > 0 {
			maxGoldEP := request.Reward.MaxGoldEP
			if maxGoldEP == 0 {
				// Standard: Maximal die Hälfte der EP durch Gold ersetzen
				maxGoldEP = newCost.Ep / 2
			}

			// Beschränke auf verfügbare EP
			if maxGoldEP > newCost.Ep {
				maxGoldEP = newCost.Ep
			}

			// Ersetze EP durch Gold (10 GS pro EP)
			newCost.Ep -= maxGoldEP
			newCost.Money += maxGoldEP * 10
		}
	}

	return &newCost
}

// Helper structures and functions
type skillInfo struct {
	Category   string
	Difficulty string
}

func getSpellInfo(spellName string) skillInfo {
	var spell models.Spell
	if err := spell.First(spellName); err != nil {
		return skillInfo{Category: "unknown", Difficulty: "unknown"}
	}
	return skillInfo{Category: spell.Category, Difficulty: strconv.Itoa(spell.Stufe)}
}

func canCharacterAfford(character *models.Char, cost *models.LearnCost) bool {
	// Check if character has enough EP
	if character.Erfahrungsschatz.EP < cost.Ep {
		return false
	}

	// Check if character has enough money
	// Assuming money is stored in Bennies (Gold pieces)
	totalMoney := character.Bennies.Gg + character.Bennies.Gp + character.Bennies.Sg
	return totalMoney >= cost.Money
}

func generateNotes(character *models.Char, request *SkillCostRequest, cost *models.LearnCost) string {
	var notes []string

	if request.Action == "learn" {
		notes = append(notes, "Neue Fertigkeit erlernen")
	} else {
		notes = append(notes, fmt.Sprintf("Verbesserung von %d auf %d", request.CurrentLevel, request.CurrentLevel+1))
	}

	// Add character class specific notes
	if character.Typ != "" {
		notes = append(notes, fmt.Sprintf("Kosten für %s", character.Typ))
	}

	// Add PP usage notes
	if request.UsePP > 0 {
		notes = append(notes, fmt.Sprintf("Verwendung von %d Praxispunkten", request.UsePP))
	}

	// Add affordability note
	if !canCharacterAfford(character, cost) {
		notes = append(notes, "Nicht genügend EP oder Gold vorhanden")
	}

	return strings.Join(notes, ". ")
}

// getPPForSkill ermittelt die verfügbaren Praxispunkte für eine spezifische Fertigkeit
func getPPForSkill(character *models.Char, skillName string) int {
	// Ermittle die tatsächliche Fertigkeit (bei Zaubern die Zaubergruppe)
	targetSkillName := getSpellCategoryNewSystem(skillName)

	for _, fertigkeit := range character.Fertigkeiten {
		if fertigkeit.Name == targetSkillName {
			return fertigkeit.Pp
		}
	}
	return 0
}

// applyPPReduction reduziert die Kosten entsprechend der verwendeten Praxispunkte
func applyPPReduction(request *SkillCostRequest, cost *models.LearnCost, availablePP int) (int, int, int) {
	if request.UsePP <= 0 {
		return cost.Ep, cost.LE, 0
	}

	// Maximal so viele PP verwenden, wie verfügbar sind
	ppToUse := request.UsePP
	if ppToUse > availablePP {
		ppToUse = availablePP
	}

	originalEP := cost.Ep
	originalLE := cost.LE

	var finalEP, finalLE int
	var reduction int

	if request.Action == "improve" {
		// Für Verbesserungen: 1 TE für 1 PP
		// Jeder PP ersetzt 1 TE, daher wird die entsprechende EP-Menge reduziert
		reduction = ppToUse // PP-Punkte direkt als Reduktion verwenden
		finalEP = originalEP - reduction
		finalLE = originalLE
		if finalEP < 0 {
			finalEP = 0
		}
	} else if request.Action == "learn" && request.Type == "spell" {
		// Für Zauber lernen: 1 LE für 1 PP
		reduction = ppToUse // PP-Punkte direkt als Reduktion verwenden
		finalLE = originalLE - reduction
		finalEP = originalEP
		if finalLE < 0 {
			finalLE = 0
		}
	} else {
		// Für andere Lernfälle: keine PP-Reduktion
		finalEP = originalEP
		finalLE = originalLE
		reduction = 0
	}

	return finalEP, finalLE, reduction
}
