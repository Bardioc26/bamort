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

// GetLernCost
func GetLernCost(c *gin.Context) {
	// Request-Parameter abrufen
	var request gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}
	charID := fmt.Sprintf("%d", request.CharId)
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}
	var costResult gsmaster.SkillCostResultNew
	costResult.CharacterID = charID

	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	if len(character.Typ) > 3 {
		costResult.CharacterClass = gsmaster.GetClassAbbreviation(character.Typ)
	} else {
		costResult.CharacterClass = character.Typ
	}

	// Normalize skill name (trim whitespace, proper case)
	costResult.SkillName = strings.TrimSpace(request.Name)

	// Lasse Kategorie und Schwierigkeit leer, damit CalcSkillLernCost die beste Option wählt
	// costResult.Category = gsmaster.GetSkillCategory(request.Name)
	// costResult.Difficulty = gsmaster.GetSkillDifficulty(costResult.Category, costResult.SkillName)
	var response []gsmaster.SkillCostResultNew

	// Für "learn" Aktion: nur eine Berechnung, da Lernkosten einmalig sind
	if request.Action == "learn" {
		levelResult := gsmaster.SkillCostResultNew{
			CharacterID:    costResult.CharacterID,
			CharacterClass: costResult.CharacterClass,
			SkillName:      costResult.SkillName,
			Category:       costResult.Category,
			Difficulty:     costResult.Difficulty,
			TargetLevel:    1, // Lernkosten sind für das Erlernen der Fertigkeit (Level 1)
		}
		err := gsmaster.GetLernCostNextLevel(&request, &levelResult, request.Reward, 1, character.Typ)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
			return
		}
		response = append(response, levelResult)
	} else {
		// Für "improve" Aktion: berechne für jedes Level von current+1 bis 18
		for i := request.CurrentLevel + 1; i <= 18; i++ {
			levelResult := gsmaster.SkillCostResultNew{
				CharacterID:    costResult.CharacterID,
				CharacterClass: costResult.CharacterClass,
				SkillName:      costResult.SkillName,
				Category:       costResult.Category,
				Difficulty:     costResult.Difficulty,
				TargetLevel:    i,
			}
			err := gsmaster.GetLernCostNextLevel(&request, &levelResult, request.Reward, i, character.Typ)
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

// GetLernCostNewSystem verwendet das neue Datenbank-Lernkosten-System
// und produziert die gleichen Ergebnisse wie GetLernCost.
//
// Unterschiede zum alten System:
// - Verwendet Models aus models/model_learning_costs.go statt der hardkodierten learningCostsData
// - Daten werden aus der Datenbank gelesen (learning_* Tabellen)
// - Unterstützt die gleichen Belohnungen und Parameter wie das alte System
// - API ist vollständig kompatibel mit GetLernCost
//
// Das neue System muss zuerst mit gsmaster.InitializeLearningCostsSystem() initialisiert werden.
func GetLernCostNewSystem(c *gin.Context) {
	// Request-Parameter abrufen
	var request gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	charID := fmt.Sprintf("%d", request.CharId)
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(character.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviationNew(character.Typ)
	} else {
		characterClass = character.Typ
	}

	// Normalize skill/spell name (trim whitespace, proper case)
	skillName := strings.TrimSpace(request.Name)

	var response []gsmaster.SkillCostResultNew
	remainingPP := request.UsePP
	remainingGold := request.UseGold

	// Für "learn" Aktion: nur eine Berechnung, da Lernkosten einmalig sind
	if request.Action == "learn" {
		if request.Type == "spell" {
			// Spell learning logic
			spellInfo, err := models.GetSpellLearningInfo(skillName, characterClass)
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
			// Skill learning logic
			skillInfo, err := models.GetSkillCategoryAndDifficulty(skillName, characterClass)
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
		skillInfo, err := models.GetSkillCategoryAndDifficulty(skillName, characterClass)
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

			err := calculateSkillImproveCostNewSystem(&request, &levelResult, i, &remainingPP, &remainingGold, skillInfo)
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
func calculateSkillImproveCostNewSystem(request *gsmaster.LernCostRequest, result *gsmaster.SkillCostResultNew, targetLevel int, remainingPP *int, remainingGold *int, skillInfo *models.SkillLearningInfo) error {
	// 1. Hole die TE-Kosten für die Verbesserung vom aktuellen Level
	teRequired, err := models.GetImprovementCost(skillInfo.SkillName, skillInfo.CategoryName, skillInfo.DifficultyName, targetLevel)
	if err != nil {
		return fmt.Errorf("Verbesserungskosten nicht gefunden für %s (Level %d): %v", skillInfo.SkillName, targetLevel, err)
	}

	// 2. Hole die EP-Kosten pro TE für diese Klasse und Kategorie
	epPerTE, err := models.GetEPPerTEForClassAndCategory(result.CharacterClass, skillInfo.CategoryName)
	if err != nil {
		return fmt.Errorf("EP-Kosten pro TE nicht gefunden für Klasse %s, Kategorie %s: %v", result.CharacterClass, skillInfo.CategoryName, err)
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
	result.EP = epPerTE * trainCost
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

// GetSkillCost berechnet die Kosten zum Erlernen oder Verbessern einer Fertigkeit
func GetSkillCost(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Request-Parameter abrufen
	var request SkillCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Normalize skill name (trim whitespace, proper case)
	request.Name = strings.TrimSpace(request.Name)

	// Validate current level for improvement
	if request.Action == "improve" && request.CurrentLevel <= 0 {
		// Try to get current level from character's skills
		currentLevel := getCurrentSkillLevel(&character, request.Name, request.Type)
		if currentLevel == -1 {
			respondWithError(c, http.StatusBadRequest, "Fertigkeit nicht bei diesem Charakter vorhanden oder current_level erforderlich")
			return
		}
		request.CurrentLevel = currentLevel
	}

	// Handle multi-level cost calculation
	if request.TargetLevel > 0 && request.Action == "improve" {
		response := calculateMultiLevelCost(&character, &request)
		if response == nil {
			respondWithError(c, http.StatusBadRequest, "Fehler bei der Multi-Level-Kostenberechnung")
			return
		}
		c.JSON(http.StatusOK, response)
		return
	}

	// Single cost calculation
	cost, originalCost, skillInfo, err := calculateSingleCost(&character, &request)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
		return
	}

	// Originalkosten berechnen (ohne PP-Reduktion)
	originalRequest := request
	originalRequest.UsePP = 0
	_, _, _, err = calculateSingleCost(&character, &originalRequest)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der ursprünglichen Kostenberechnung: "+err.Error())
		return
	}

	// PP-Informationen sammeln (fertigkeitsspezifisch)
	availablePP := getPPForSkill(&character, request.Name)
	ppUsed := request.UsePP
	if ppUsed > availablePP {
		ppUsed = availablePP
	}

	// PP-Reduktion berechnen
	var ppReduction int
	if request.UsePP > 0 {
		if request.Action == "improve" {
			ppReduction = ppUsed // PP entsprechen direkt der TE-Reduktion
		} else if request.Action == "learn" && request.Type == "spell" {
			ppReduction = ppUsed // PP entsprechen direkt der LE-Reduktion
		}
	}

	// Check if character can afford it
	canAfford := canCharacterAfford(&character, cost)

	// Belohnungsinformationen berechnen
	var rewardApplied string
	var savings *models.LearnCost
	var goldUsedForEP int

	if request.Reward != nil && request.Reward.Type != "" {
		rewardApplied = request.Reward.Type

		// Ersparnisse berechnen
		savings = &models.LearnCost{
			Ep:    originalCost.Ep - cost.Ep,
			LE:    originalCost.LE - cost.LE,
			Money: originalCost.Money - cost.Money,
		}

		// Gold für EP berechnen
		if request.Reward.Type == "gold_for_ep" && request.Reward.UseGoldForEP {
			goldUsedForEP = (cost.Money - originalCost.Money) / 10
		}
	}

	// Create response
	response := &SkillCostResponse{
		LearnCost:          cost,
		SkillName:          request.Name,
		SkillType:          request.Type,
		Action:             request.Action,
		CharacterID:        character.ID,
		CurrentLevel:       request.CurrentLevel,
		TargetLevel:        request.TargetLevel,
		Category:           skillInfo.Category,
		Difficulty:         skillInfo.Difficulty,
		CanAfford:          canAfford,
		Notes:              generateNotes(&character, &request, cost),
		PPUsed:             ppUsed,
		PPAvailable:        availablePP,
		PPReduction:        ppReduction,
		OriginalCost:       originalCost.Ep,
		FinalCost:          cost.Ep,
		RewardApplied:      rewardApplied,
		OriginalCostStruct: originalCost,
		Savings:            savings,
		GoldUsedForEP:      goldUsedForEP,
	}

	c.JSON(http.StatusOK, response)
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

// Helper function to calculate single cost
func calculateSingleCost(character *models.Char, request *SkillCostRequest) (*models.LearnCost, *models.LearnCost, *skillInfo, error) {
	var cost *models.LearnCost
	var err error
	var info skillInfo

	switch {
	case request.Action == "learn" && request.Type == "skill":
		cost, err = gsmaster.CalculateDetailedSkillLearningCost(request.Name, character.Typ)
		if err == nil {
			info = getSkillInfo(request.Name, request.Type)
		}

	case request.Action == "improve" && request.Type == "skill":
		cost, err = gsmaster.CalculateDetailedSkillImprovementCost(request.Name, character.Typ, request.CurrentLevel)
		if err == nil {
			info = getSkillInfo(request.Name, request.Type)
		}

	case request.Action == "improve" && request.Type == "weapon":
		cost, err = gsmaster.CalculateDetailedSkillImprovementCost(request.Name, character.Typ, request.CurrentLevel)
		if err == nil {
			info = getSkillInfo(request.Name, request.Type)
		}

	case request.Action == "learn" && request.Type == "spell":
		cost, err = gsmaster.CalculateDetailedSpellLearningCost(request.Name, character.Typ)
		if err == nil {
			info = getSpellInfo(request.Name)
		}

	default:
		return nil, nil, nil, fmt.Errorf("ungültige Kombination aus Aktion und Typ")
	}

	if err != nil {
		return nil, nil, nil, err
	}

	// Belohnungen anwenden, falls spezifiziert
	originalCost := *cost // Kopie der ursprünglichen Kosten
	if request.Reward != nil && request.Reward.Type != "" {
		cost = applyReward(cost, request)
	}

	// Praxispunkte anwenden, falls angefordert (fertigkeitsspezifisch)
	if request.UsePP > 0 {
		availablePP := getPPForSkill(character, request.Name)
		finalEP, finalLE, _ := applyPPReduction(request, cost, availablePP)

		// Erstelle eine neue LearnCost mit den reduzierten Werten
		cost = &models.LearnCost{
			Stufe: cost.Stufe,
			LE:    finalLE,
			Ep:    finalEP,
			Money: cost.Money, // Geldkosten bleiben unverändert
		}
	}

	return cost, &originalCost, &info, err
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

// Helper function to calculate multi-level costs
func calculateMultiLevelCost(character *models.Char, request *SkillCostRequest) *MultiLevelCostResponse {
	if request.TargetLevel <= request.CurrentLevel {
		return nil
	}

	var levelCosts []SkillCostResponse
	totalEP := 0
	totalMoney := 0
	remainingPP := request.UsePP

	for level := request.CurrentLevel; level < request.TargetLevel; level++ {
		tempRequest := *request
		tempRequest.CurrentLevel = level
		tempRequest.Action = "improve"

		// Verteile die PP auf die verschiedenen Level
		if remainingPP > 0 {
			tempRequest.UsePP = 1 // Maximal 1 PP pro Level
			remainingPP--
		} else {
			tempRequest.UsePP = 0
		}

		cost, originalCost, skillInfo, err := calculateSingleCost(character, &tempRequest)
		if err != nil {
			continue
		}

		// Originalkosten berechnen (ohne PP)
		originalRequest := tempRequest
		originalRequest.UsePP = 0
		_, _, _, _ = calculateSingleCost(character, &originalRequest)

		// PP-Informationen sammeln (fertigkeitsspezifisch)
		availablePP := getPPForSkill(character, request.Name)
		ppUsed := tempRequest.UsePP
		if ppUsed > availablePP {
			ppUsed = availablePP
		}

		ppReduction := 0
		if tempRequest.UsePP > 0 {
			ppReduction = ppUsed
		}

		// Belohnungsinformationen für Level berechnen
		var rewardApplied string
		var savings *models.LearnCost
		var goldUsedForEP int

		if tempRequest.Reward != nil && tempRequest.Reward.Type != "" {
			rewardApplied = tempRequest.Reward.Type

			// Ersparnisse berechnen
			savings = &models.LearnCost{
				Ep:    originalCost.Ep - cost.Ep,
				LE:    originalCost.LE - cost.LE,
				Money: originalCost.Money - cost.Money,
			}

			// Gold für EP berechnen
			if tempRequest.Reward.Type == "gold_for_ep" && tempRequest.Reward.UseGoldForEP {
				goldUsedForEP = (cost.Money - originalCost.Money) / 10
			}
		}

		levelCost := SkillCostResponse{
			LearnCost:          cost,
			SkillName:          request.Name,
			SkillType:          request.Type,
			Action:             "improve",
			CharacterID:        character.ID,
			CurrentLevel:       level,
			TargetLevel:        level + 1,
			Category:           skillInfo.Category,
			Difficulty:         skillInfo.Difficulty,
			CanAfford:          canCharacterAfford(character, cost),
			PPUsed:             ppUsed,
			PPAvailable:        availablePP,
			PPReduction:        ppReduction,
			OriginalCost:       originalCost.Ep,
			FinalCost:          cost.Ep,
			RewardApplied:      rewardApplied,
			OriginalCostStruct: originalCost,
			Savings:            savings,
			GoldUsedForEP:      goldUsedForEP,
		}

		levelCosts = append(levelCosts, levelCost)
		totalEP += cost.Ep
		totalMoney += cost.Money
	}

	totalCost := &models.LearnCost{
		Stufe: request.TargetLevel,
		LE:    0,
		Ep:    totalEP,
		Money: totalMoney,
	}

	return &MultiLevelCostResponse{
		SkillName:      request.Name,
		SkillType:      request.Type,
		CharacterID:    character.ID,
		CurrentLevel:   request.CurrentLevel,
		TargetLevel:    request.TargetLevel,
		LevelCosts:     levelCosts,
		TotalCost:      totalCost,
		CanAffordTotal: canCharacterAfford(character, totalCost),
	}
}

// Helper structures and functions
type skillInfo struct {
	Category   string
	Difficulty string
}

func getSkillInfo(skillName, skillType string) skillInfo {
	var skill models.Skill
	if err := skill.First(skillName); err != nil {
		return skillInfo{Category: "unknown", Difficulty: "unknown"}
	}

	// Fallback für fehlende Category und Difficulty Werte
	category := skill.Category
	difficulty := skill.Difficulty

	if category == "" {
		// Standard-Kategorien basierend auf Skill-Namen
		category = gsmaster.GetDefaultCategory(skillName)
	}

	if difficulty == "" {
		// Standard-Schwierigkeit für verschiedene Skills
		difficulty = gsmaster.GetDefaultDifficulty(skillName)
	}

	return skillInfo{Category: category, Difficulty: difficulty}
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
	targetSkillName := getSpellCategory(skillName)

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

func CalcSkillLearnCost(req *gsmaster.LernCostRequest, skillCostInfo *gsmaster.SkillCostResultNew) error {
	// Fallback-Werte für Skills ohne definierte Kategorie/Schwierigkeit

	result, err := gsmaster.CalculateSkillLearningCosts(skillCostInfo.CharacterClass, skillCostInfo.Category, skillCostInfo.Difficulty)
	if err != nil {
		return err
	}

	//Stufe: 0, // Lernen startet bei Stufe 0
	skillCostInfo.LE = result.LE
	skillCostInfo.EP = result.EP
	skillCostInfo.GoldCost = result.GoldCost
	return nil
}
