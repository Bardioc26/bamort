package character

import (
	"bamort/gsmaster"
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
}

type SkillCostResponse struct {
	*gsmaster.LearnCost
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
	PPUsed       int    `json:"pp_used,omitempty"`       // Anzahl der verwendeten Praxispunkte
	PPAvailable  int    `json:"pp_available,omitempty"`  // Verfügbare Praxispunkte für diese Kategorie
	PPReduction  int    `json:"pp_reduction,omitempty"`  // Reduktion der Kosten durch PP
	OriginalCost int    `json:"original_cost,omitempty"` // Ursprüngliche Kosten (vor PP-Reduktion)
	FinalCost    int    `json:"final_cost,omitempty"`    // Endgültige Kosten (nach PP-Reduktion)
}

type MultiLevelCostResponse struct {
	SkillName      string              `json:"skill_name"`
	SkillType      string              `json:"skill_type"`
	CharacterID    uint                `json:"character_id"`
	CurrentLevel   int                 `json:"current_level"`
	TargetLevel    int                 `json:"target_level"`
	LevelCosts     []SkillCostResponse `json:"level_costs"`
	TotalCost      *gsmaster.LearnCost `json:"total_cost"`
	CanAffordTotal bool                `json:"can_afford_total"`
}

// GetSkillCost berechnet die Kosten zum Erlernen oder Verbessern einer Fertigkeit
func GetSkillCost(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character Char
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
	cost, skillInfo, err := calculateSingleCost(&character, &request)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
		return
	}

	// Originalkosten berechnen (ohne PP-Reduktion)
	originalRequest := request
	originalRequest.UsePP = 0
	originalCost, _, err := calculateSingleCost(&character, &originalRequest)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der ursprünglichen Kostenberechnung: "+err.Error())
		return
	}

	// PP-Informationen sammeln (fertigkeitsspezifisch)
	skillType := getSkillType(request.Type)
	availablePP := getPPForSkill(&character, request.Name, skillType)
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

	// Create response
	response := &SkillCostResponse{
		LearnCost:    cost,
		SkillName:    request.Name,
		SkillType:    request.Type,
		Action:       request.Action,
		CharacterID:  character.ID,
		CurrentLevel: request.CurrentLevel,
		TargetLevel:  request.TargetLevel,
		Category:     skillInfo.Category,
		Difficulty:   skillInfo.Difficulty,
		CanAfford:    canAfford,
		Notes:        generateNotes(&character, &request, cost),
		PPUsed:       ppUsed,
		PPAvailable:  availablePP,
		PPReduction:  ppReduction,
		OriginalCost: originalCost.Ep,
		FinalCost:    cost.Ep,
	}

	c.JSON(http.StatusOK, response)
}

// Helper function to get current skill level from character
func getCurrentSkillLevel(character *Char, skillName, skillType string) int {
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
func calculateSingleCost(character *Char, request *SkillCostRequest) (*gsmaster.LearnCost, *skillInfo, error) {
	var cost *gsmaster.LearnCost
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
		return nil, nil, fmt.Errorf("ungültige Kombination aus Aktion und Typ")
	}

	if err != nil {
		return nil, nil, err
	}

	// Praxispunkte anwenden, falls angefordert (fertigkeitsspezifisch)
	if request.UsePP > 0 {
		skillType := getSkillType(request.Type)
		availablePP := getPPForSkill(character, request.Name, skillType)
		finalEP, finalLE, _ := applyPPReduction(request, cost, availablePP)

		// Erstelle eine neue LearnCost mit den reduzierten Werten
		cost = &gsmaster.LearnCost{
			Stufe: cost.Stufe,
			LE:    finalLE,
			Ep:    finalEP,
			Money: cost.Money, // Geldkosten bleiben unverändert
		}
	}

	return cost, &info, err
}

// Helper function to calculate multi-level costs
func calculateMultiLevelCost(character *Char, request *SkillCostRequest) *MultiLevelCostResponse {
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

		cost, skillInfo, err := calculateSingleCost(character, &tempRequest)
		if err != nil {
			continue
		}

		// Originalkosten berechnen (ohne PP)
		originalRequest := tempRequest
		originalRequest.UsePP = 0
		originalCost, _, _ := calculateSingleCost(character, &originalRequest)

		// PP-Informationen sammeln (fertigkeitsspezifisch)
		skillType := getSkillType(request.Type)
		availablePP := getPPForSkill(character, request.Name, skillType)
		ppUsed := tempRequest.UsePP
		if ppUsed > availablePP {
			ppUsed = availablePP
		}

		ppReduction := 0
		if tempRequest.UsePP > 0 {
			ppReduction = ppUsed
		}

		levelCost := SkillCostResponse{
			LearnCost:    cost,
			SkillName:    request.Name,
			SkillType:    request.Type,
			Action:       "improve",
			CharacterID:  character.ID,
			CurrentLevel: level,
			TargetLevel:  level + 1,
			Category:     skillInfo.Category,
			Difficulty:   skillInfo.Difficulty,
			CanAfford:    canCharacterAfford(character, cost),
			PPUsed:       ppUsed,
			PPAvailable:  availablePP,
			PPReduction:  ppReduction,
			OriginalCost: originalCost.Ep,
			FinalCost:    cost.Ep,
		}

		levelCosts = append(levelCosts, levelCost)
		totalEP += cost.Ep
		totalMoney += cost.Money
	}

	totalCost := &gsmaster.LearnCost{
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
	var skill gsmaster.Skill
	if err := skill.First(skillName); err != nil {
		return skillInfo{Category: "unknown", Difficulty: "unknown"}
	}
	return skillInfo{Category: skill.Category, Difficulty: skill.Difficulty}
}

func getSpellInfo(spellName string) skillInfo {
	var spell gsmaster.Spell
	if err := spell.First(spellName); err != nil {
		return skillInfo{Category: "unknown", Difficulty: "unknown"}
	}
	return skillInfo{Category: spell.Category, Difficulty: strconv.Itoa(spell.Stufe)}
}

func canCharacterAfford(character *Char, cost *gsmaster.LearnCost) bool {
	// Check if character has enough EP
	if character.Erfahrungsschatz.Value < cost.Ep {
		return false
	}

	// Check if character has enough money
	// Assuming money is stored in Bennies (Gold pieces)
	totalMoney := character.Bennies.Gg + character.Bennies.Gp + character.Bennies.Sg
	if totalMoney < cost.Money {
		return false
	}

	return true
}

func generateNotes(character *Char, request *SkillCostRequest, cost *gsmaster.LearnCost) string {
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
func getPPForSkill(character *Char, skillName string, skillType string) int {
	for _, pp := range character.Praxispunkte {
		if pp.SkillName == skillName && pp.SkillType == skillType {
			return pp.Anzahl
		}
	}
	return 0
}

// getSkillType konvertiert den Request-Type in den internen SkillType
func getSkillType(requestType string) string {
	switch requestType {
	case "skill":
		return "fertigkeit"
	case "weapon":
		return "waffenfertigkeit"
	case "spell":
		return "zauber"
	default:
		return ""
	}
}

// applyPPReduction reduziert die Kosten entsprechend der verwendeten Praxispunkte
func applyPPReduction(request *SkillCostRequest, cost *gsmaster.LearnCost, availablePP int) (int, int, int) {
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
