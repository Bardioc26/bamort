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

	for level := request.CurrentLevel; level < request.TargetLevel; level++ {
		tempRequest := *request
		tempRequest.CurrentLevel = level
		tempRequest.Action = "improve"

		cost, skillInfo, err := calculateSingleCost(character, &tempRequest)
		if err != nil {
			continue
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

	// Add affordability note
	if !canCharacterAfford(character, cost) {
		notes = append(notes, "Nicht genügend EP oder Gold vorhanden")
	}

	return strings.Join(notes, ". ")
}
