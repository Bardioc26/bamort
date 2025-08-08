package character

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"strconv"
	"strings"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Character Handlers

type LearnRequestStruct struct {
	SkillType string `json:"skillType"`
	Name      string `json:"name"`
	Stufe     int    `json:"stufe"`
}

func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func ListCharacters(c *gin.Context) {
	var characters []models.Char
	var listOfChars []models.CharList
	if err := database.DB.Find(&characters).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve characters")
		return
	}
	for i := range characters {
		listOfChars = append(listOfChars, models.CharList{
			BamortBase: models.BamortBase{
				ID:   characters[i].ID,
				Name: characters[i].Name,
			},
			Rasse: characters[i].Rasse,
			Typ:   characters[i].Typ,
			Grad:  characters[i].Grad,
			Owner: "test",
		})
	}
	c.JSON(http.StatusOK, listOfChars)
}

func CreateCharacter(c *gin.Context) {
	var character models.Char
	if err := c.ShouldBindJSON(&character); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Create(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create character")
		return
	}

	c.JSON(http.StatusCreated, character)
}
func GetCharacter(c *gin.Context) {
	id := c.Param("id")
	var character models.Char
	err := character.FirstID(id)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}
	feChar := ToFeChar(&character)
	c.JSON(http.StatusOK, feChar)
}
func UpdateCharacter(c *gin.Context) {
	id := c.Param("id")
	var character models.Char

	// First, find the existing character
	err := character.FirstID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Bind the updated data
	if err := c.ShouldBindJSON(&character); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Save the updated character
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update character")
		return
	}

	c.JSON(http.StatusOK, character)
}
func DeleteCharacter(c *gin.Context) {
	id := c.Param("id")
	var character models.Char
	err := character.FirstID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}
	err = character.Delete()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete character")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Character deleted successfully"})
}

// Add Fertigkeit by putting it directly to the DB
func AddFertigkeit(charID uint, fertigkeit *models.SkFertigkeit) error {
	// Set the foreign key for the new Eigenschaft
	fertigkeit.CharacterID = charID

	// Save the new Eigenschaft to the database
	if err := database.DB.Create(&fertigkeit).Error; err != nil {
		return fmt.Errorf("failed to add Eigenschaft: %w", err)
	}
	return nil
}

// Append the new Fertigkeit to the slice of the characters property
//character.Fertigkeiten = append(character.Fertigkeiten, fertigkeit)

func ToFeChar(object *models.Char) *models.FeChar {
	feC := &models.FeChar{
		Char: *object,
	}
	skills, innateSkills, categories := splitSkills(object.Fertigkeiten)
	feC.Fertigkeiten = skills
	feC.InnateSkills = innateSkills
	feC.CategorizedSkills = categories
	return feC
}

func splitSkills(object []models.SkFertigkeit) ([]models.SkFertigkeit, []models.SkFertigkeit, map[string][]models.SkFertigkeit) {
	var normSkills []models.SkFertigkeit
	var innateSkills []models.SkFertigkeit
	//var categories map[string][]models.Fertigkeit
	categories := make(map[string][]models.SkFertigkeit)
	for _, skill := range object {
		gsmsk := skill.GetSkillByName()
		if gsmsk.Improvable {
			category := "Unkategorisiert"
			if gsmsk.ID != 0 && gsmsk.Category != "" {
				category = gsmsk.Category
			}
			normSkills = append(normSkills, skill)
			if _, exists := categories[category]; !exists {
				categories[category] = make([]models.SkFertigkeit, 0)
			}
			categories[category] = append(categories[category], skill)
		} else {
			innateSkills = append(innateSkills, skill)
		}
	}

	return normSkills, innateSkills, categories
}

// GetLearnSkillCostOld is deprecated. Use GetLearnSkillCost instead.
// This function uses the old hardcoded learning cost system.
func GetLearnSkillCostOld(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	// Load the skill from the request
	var s models.SkFertigkeit
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var skill models.Skill
	if err := skill.First(s.Name); err != nil {
		respondWithError(c, http.StatusBadRequest, "can not find speel in gsmaster: "+err.Error())
		return
	}

	cost, err := gsmaster.CalculateSkillLearnCostOld(skill.Name, character.Typ)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
		return
	}

	// Return the updated character
	c.JSON(http.StatusOK, cost)
}

// GetLearnSpellCostOld is deprecated. Use GetLearnSpellCost instead.
// This function uses the old hardcoded learning cost system.
func GetLearnSpellCostOld(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	// Load the spell from the request
	var s models.SkZauber
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	var spell models.Spell
	if err := spell.First(s.Name); err != nil {
		respondWithError(c, http.StatusBadRequest, "can not find speel in gsmaster: "+err.Error())
		return
	}
	sd := gsmaster.SpellDefinition{
		Name:   spell.Name,
		Stufe:  spell.Stufe,
		School: spell.Category,
	}

	cost, err := gsmaster.CalculateSpellLearnCostOld(spell.Name, character.Typ)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "error getting costs to learn spell: "+err.Error())
		return
	}

	sd.CostEP = cost
	// Return the updated character
	c.JSON(http.StatusOK, sd)
}

// ExperienceAndWealthResponse repräsentiert die Antwort für EP und Vermögen
type ExperienceAndWealthResponse struct {
	ExperiencePoints int `json:"experience_points"`
	Wealth           struct {
		Goldstücke   int `json:"gold_coins"`   // GS
		Silberstücke int `json:"silver_coins"` // SS
		Kupferstücke int `json:"copper_coins"` // KS
		TotalInGS    int `json:"total_in_ss"`  // Gesamt in Silberstücken
	} `json:"wealth"`
}

// GetCharacterExperienceAndWealth gibt nur die EP und Vermögensdaten eines Charakters zurück
func GetCharacterExperienceAndWealth(c *gin.Context) {
	id := c.Param("id")
	var character models.Char

	// Lade nur die benötigten Felder
	err := database.DB.
		Preload("Erfahrungsschatz").
		Preload("Vermoegen").
		First(&character, id).Error
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Berechne Gesamtvermögen in Silbergroschen
	// Annahme: 1 GS = 10 SS, 1 SS = 10 KS (typische Midgard Währung)
	gs := character.Vermoegen.Goldstücke
	ss := character.Vermoegen.Silberstücke
	ks := character.Vermoegen.Kupferstücke
	totalInSS := (gs * 10) + ss + (ks / 10)

	response := ExperienceAndWealthResponse{
		ExperiencePoints: character.Erfahrungsschatz.EP,
	}
	response.Wealth.Goldstücke = gs
	response.Wealth.Silberstücke = ss
	response.Wealth.Kupferstücke = ks
	response.Wealth.TotalInGS = totalInSS

	c.JSON(http.StatusOK, response)
}

// UpdateExperienceRequest repräsentiert die Anfrage für EP-Update
type UpdateExperienceRequest struct {
	ExperiencePoints int    `json:"experience_points" binding:"required,min=0"`
	Reason           string `json:"reason,omitempty"` // Grund der Änderung
	Notes            string `json:"notes,omitempty"`  // Zusätzliche Notizen
}

// UpdateCharacterExperience aktualisiert die Erfahrungspunkte eines Charakters
// TODO Wenn EP verändert werden ändert sich auch ES
func UpdateCharacterExperience(c *gin.Context) {
	id := c.Param("id")
	var character models.Char

	// Lade den Charakter
	err := database.DB.
		Preload("Erfahrungsschatz").
		First(&character, id).Error
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Parse Request
	var req UpdateExperienceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Standard-Grund setzen, falls nicht angegeben
	if req.Reason == "" {
		req.Reason = string(ReasonManual)
	}

	// Alten Wert für Audit-Log speichern
	oldValue := 0
	if character.Erfahrungsschatz.ID != 0 {
		oldValue = character.Erfahrungsschatz.EP
	}

	// Aktualisiere oder erstelle Erfahrungsschatz
	if character.Erfahrungsschatz.ID == 0 {
		// Erstelle neuen Erfahrungsschatz
		character.Erfahrungsschatz = models.Erfahrungsschatz{
			BamortCharTrait: models.BamortCharTrait{
				CharacterID: character.ID,
			},
			EP: req.ExperiencePoints,
		}
		if err := database.DB.Create(&character.Erfahrungsschatz).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to create experience record")
			return
		}
	} else {
		// Aktualisiere existierenden Erfahrungsschatz
		character.Erfahrungsschatz.EP = req.ExperiencePoints
		if err := database.DB.Save(&character.Erfahrungsschatz).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to update experience")
			return
		}
	}

	// Audit-Log-Eintrag erstellen (nur wenn sich der Wert geändert hat)
	if oldValue != req.ExperiencePoints {
		// TODO: User-ID aus dem Authentifizierungs-Context holen
		userID := uint(0) // Placeholder

		err = CreateAuditLogEntry(
			character.ID,
			"experience_points",
			oldValue,
			req.ExperiencePoints,
			AuditLogReason(req.Reason),
			userID,
			req.Notes,
		)
		if err != nil {
			// Log-Fehler sollten die Hauptoperation nicht blockieren
			// TODO: Proper logging implementieren
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":           "Experience updated successfully",
		"experience_points": req.ExperiencePoints,
		"audit_logged":      oldValue != req.ExperiencePoints,
	})
}

// UpdateWealthRequest repräsentiert die Anfrage für Vermögens-Update
type UpdateWealthRequest struct {
	Goldstücke   *int   `json:"goldstücke,omitempty"`
	Silberstücke *int   `json:"silberstücke,omitempty"`
	Kupferstücke *int   `json:"kupferstücke,omitempty"`
	Reason       string `json:"reason,omitempty"` // Grund der Änderung
	Notes        string `json:"notes,omitempty"`  // Zusätzliche Notizen
}

// UpdateCharacterWealth aktualisiert das Vermögen eines Charakters
func UpdateCharacterWealth(c *gin.Context) {
	id := c.Param("id")
	var character models.Char

	// Lade den Charakter
	err := database.DB.
		Preload("Vermoegen").
		First(&character, id).Error
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Parse Request
	var req UpdateWealthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Standard-Grund setzen, falls nicht angegeben
	if req.Reason == "" {
		req.Reason = string(ReasonManual)
	}

	// Alte Werte für Audit-Log speichern
	oldGold := 0
	oldSilver := 0
	oldCopper := 0
	if character.Vermoegen.ID != 0 {
		oldGold = character.Vermoegen.Goldstücke
		oldSilver = character.Vermoegen.Silberstücke
		oldCopper = character.Vermoegen.Kupferstücke
	}

	// Aktualisiere oder erstelle Vermögen
	if character.Vermoegen.ID == 0 {
		// Erstelle neues Vermögen
		character.Vermoegen = models.Vermoegen{
			BamortCharTrait: models.BamortCharTrait{
				CharacterID: character.ID,
			},
			Goldstücke:   getValueOrDefault(req.Goldstücke, 0),
			Silberstücke: getValueOrDefault(req.Silberstücke, 0),
			Kupferstücke: getValueOrDefault(req.Kupferstücke, 0),
		}
		if err := database.DB.Create(&character.Vermoegen).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to create wealth record")
			return
		}
	} else {
		// Aktualisiere existierendes Vermögen
		if req.Goldstücke != nil {
			character.Vermoegen.Goldstücke = *req.Goldstücke
		}
		if req.Silberstücke != nil {
			character.Vermoegen.Silberstücke = *req.Silberstücke
		}
		if req.Kupferstücke != nil {
			character.Vermoegen.Kupferstücke = *req.Kupferstücke
		}
		if err := database.DB.Save(&character.Vermoegen).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to update wealth")
			return
		}
	}

	// Audit-Log-Einträge erstellen (nur für geänderte Werte)
	// TODO: User-ID aus dem Authentifizierungs-Context holen
	userID := uint(0) // Placeholder

	if req.Goldstücke != nil && oldGold != character.Vermoegen.Goldstücke {
		CreateAuditLogEntry(
			character.ID,
			"gold",
			oldGold,
			character.Vermoegen.Goldstücke,
			AuditLogReason(req.Reason),
			userID,
			req.Notes,
		)
	}

	if req.Silberstücke != nil && oldSilver != character.Vermoegen.Silberstücke {
		CreateAuditLogEntry(
			character.ID,
			"silver",
			oldSilver,
			character.Vermoegen.Silberstücke,
			AuditLogReason(req.Reason),
			userID,
			req.Notes,
		)
	}

	if req.Kupferstücke != nil && oldCopper != character.Vermoegen.Kupferstücke {
		CreateAuditLogEntry(
			character.ID,
			"copper",
			oldCopper,
			character.Vermoegen.Kupferstücke,
			AuditLogReason(req.Reason),
			userID,
			req.Notes,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Wealth updated successfully",
		"wealth": gin.H{
			"goldstücke":   character.Vermoegen.Goldstücke,
			"silberstücke": character.Vermoegen.Silberstücke,
			"kupferstücke": character.Vermoegen.Kupferstücke,
		},
	})
}

// getValueOrDefault gibt den Wert zurück oder einen Default-Wert falls nil
func getValueOrDefault(value *int, defaultValue int) int {
	if value != nil {
		return *value
	}
	return defaultValue
}

// updateOrCreateSkill aktualisiert eine vorhandene Fertigkeit oder erstellt eine neue
func updateOrCreateSkill(character *models.Char, skillName string, newLevel int) error {
	// Suche erst in normalen Fertigkeiten
	for i := range character.Fertigkeiten {
		if character.Fertigkeiten[i].Name == skillName {
			character.Fertigkeiten[i].Fertigkeitswert = newLevel
			return database.DB.Save(&character.Fertigkeiten[i]).Error
		}
	}

	// Suche in Waffenfertigkeiten
	for i := range character.Waffenfertigkeiten {
		if character.Waffenfertigkeiten[i].Name == skillName {
			character.Waffenfertigkeiten[i].Fertigkeitswert = newLevel
			return database.DB.Save(&character.Waffenfertigkeiten[i]).Error
		}
	}

	// Fertigkeit nicht gefunden - erstelle neue normale Fertigkeit
	newSkill := models.SkFertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: skillName,
			},
			CharacterID: character.ID,
		},
		Fertigkeitswert: newLevel,
		Improvable:      true,
	}

	if err := database.DB.Create(&newSkill).Error; err != nil {
		return err
	}

	// Füge zur Charakter-Liste hinzu
	character.Fertigkeiten = append(character.Fertigkeiten, newSkill)
	return nil
}

// addSpellToCharacter fügt einen neuen Zauber zum Charakter hinzu
func addSpellToCharacter(character *models.Char, spellName string) error {
	// Prüfe, ob Zauber bereits existiert
	for _, spell := range character.Zauber {
		if spell.Name == spellName {
			// Zauber bereits vorhanden, nichts zu tun
			return nil
		}
	}

	// Erstelle neuen Zauber
	newSpell := models.SkZauber{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: spellName,
			},
			CharacterID: character.ID,
		},
	}

	if err := database.DB.Create(&newSpell).Error; err != nil {
		return err
	}

	// Füge zur Charakter-Liste hinzu
	character.Zauber = append(character.Zauber, newSpell)
	return nil
}

// Learn and Improve handlers with automatic audit logging

// LearnSpellRequest definiert die Struktur für das Lernen eines Zaubers
type LearnSpellRequest struct {
	Name  string `json:"name" binding:"required"`
	Notes string `json:"notes,omitempty"`
}

// calculateMultiLevelCostsOld is deprecated. Use the new database-based learning cost system instead.
// This function uses the old hardcoded learning cost system.
// calculateMultiLevelCostsOld berechnet die Kosten für mehrere Level-Verbesserungen mit gsmaster.GetLernCostNextLevel
func calculateMultiLevelCostsOld(character *models.Char, skillName string, currentLevel int, levelsToLearn []int, rewardType string, usePP, useGold int) (*models.LearnCost, error) {
	if len(levelsToLearn) == 0 {
		return nil, fmt.Errorf("keine Level zum Lernen angegeben")
	}

	// Sortiere die Level aufsteigend
	sortedLevels := make([]int, len(levelsToLearn))
	copy(sortedLevels, levelsToLearn)
	for i := 0; i < len(sortedLevels)-1; i++ {
		for j := i + 1; j < len(sortedLevels); j++ {
			if sortedLevels[i] > sortedLevels[j] {
				sortedLevels[i], sortedLevels[j] = sortedLevels[j], sortedLevels[i]
			}
		}
	}

	// Erstelle LernCostRequest
	var rewardTypePtr *string
	if rewardType != "" {
		rewardTypePtr = &rewardType
	}

	request := gsmaster.LernCostRequest{
		CharId:       uint(character.ID),
		Name:         skillName,
		CurrentLevel: currentLevel,
		Type:         "skill",
		Action:       "improve",
		TargetLevel:  sortedLevels[len(sortedLevels)-1], // Höchstes Level als Ziel
		UsePP:        usePP,
		UseGold:      useGold,
		Reward:       rewardTypePtr,
	}

	totalCost := &models.LearnCost{
		Stufe: sortedLevels[len(sortedLevels)-1],
		LE:    0,
		Ep:    0,
		Money: 0,
	}

	remainingPP := usePP
	remainingGold := useGold

	// Berechne Kosten für jedes Level
	for _, targetLevel := range sortedLevels {
		classAbr := getCharacterClassOld(character)
		cat, difficulty, _ := gsmaster.FindBestCategoryForSkillLearningOld(skillName, classAbr)
		levelResult := gsmaster.SkillCostResultNew{
			CharacterID:    fmt.Sprintf("%d", character.ID),
			CharacterClass: classAbr,
			SkillName:      skillName,
			Category:       cat,
			Difficulty:     gsmaster.GetSkillDifficultyOld(difficulty, skillName),
			TargetLevel:    targetLevel,
		}

		// Temporäre Request für dieses Level
		tempRequest := request
		tempRequest.CurrentLevel = targetLevel - 1
		tempRequest.UsePP = remainingPP
		tempRequest.UseGold = remainingGold

		err := gsmaster.GetLernCostNextLevelOld(&tempRequest, &levelResult, rewardTypePtr, targetLevel, character.Typ)
		if err != nil {
			return nil, fmt.Errorf("fehler bei Level %d: %v", targetLevel, err)
		}

		// Aktualisiere verbleibende Ressourcen
		if levelResult.PPUsed > 0 {
			remainingPP -= levelResult.PPUsed
			if remainingPP < 0 {
				remainingPP = 0
			}
		}
		if levelResult.GoldUsed > 0 {
			remainingGold -= levelResult.GoldUsed
			if remainingGold < 0 {
				remainingGold = 0
			}
		}

		totalCost.Ep += levelResult.EP
		totalCost.Money += levelResult.GoldCost
		totalCost.LE += levelResult.LE
	}

	return totalCost, nil
}

// getCharacterClassOld is deprecated. Use character.Klasse directly or appropriate database lookups.
// This function provides backwards compatibility for character class access.
// getCharacterClassOld gibt die Charakterklassen-Abkürzung zurück
func getCharacterClassOld(character *models.Char) string {
	if len(character.Typ) > 3 {
		return gsmaster.GetClassAbbreviationOld(character.Typ)
	}
	return character.Typ
}

// LearnSkillOld is deprecated. Use LearnSkill instead.
// This function uses the old hardcoded learning cost system.
// LearnSkillOld lernt eine neue Fertigkeit und erstellt Audit-Log-Einträge
func LearnSkillOld(c *gin.Context) {
	charID := c.Param("id")
	var character models.Char

	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Verwende gsmaster.LernCostRequest direkt
	var request gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Setze Charakter-ID und Action für learning
	request.CharId = character.ID
	request.Action = "learn"
	if request.Type == "" {
		request.Type = "skill" // Default zu skill für Learning
	}
	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(character.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviationOld(character.Typ)
	} else {
		characterClass = character.Typ
	}

	// Bestimme das finale Level
	finalLevel := request.TargetLevel
	if finalLevel <= 0 {
		finalLevel = 1 // Standard für neue Fertigkeit
	}

	// Für Learning müssen wir von Level 0 (nicht gelernt) auf finalLevel lernen
	var totalEP, totalGold, totalPP int
	var err error

	// Loop für jeden Level von 0 bis finalLevel (für neue Fertigkeiten)
	for tempLevel := 0; tempLevel < finalLevel; tempLevel++ {
		nextLevel := tempLevel + 1

		// Erstelle temporären Request für diesen Level
		tempRequest := request
		tempRequest.CurrentLevel = tempLevel
		tempRequest.TargetLevel = nextLevel

		// Für das erste Level (0->1) ist es ein "learn", für weitere Level "improve"
		if tempLevel == 0 {
			tempRequest.Action = "learn"
		} else {
			tempRequest.Action = "improve"
		}

		// Berechne Kosten für diesen einen Level
		var costResult gsmaster.SkillCostResultNew
		costResult.CharacterID = fmt.Sprintf("%d", character.ID)
		costResult.CharacterClass = characterClass
		costResult.SkillName = request.Name

		err = gsmaster.GetLernCostNextLevelOld(&tempRequest, &costResult, request.Reward, nextLevel, character.Rasse)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, fmt.Sprintf("Fehler bei Level %d: %v", nextLevel, err))
			return
		}

		// Addiere die Kosten
		totalEP += costResult.EP
		totalGold += costResult.GoldCost
		totalPP += costResult.PPUsed
	}

	// Prüfe, ob genügend EP vorhanden sind
	currentEP := character.Erfahrungsschatz.EP
	if currentEP < totalEP {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Erfahrungspunkte vorhanden")
		return
	}

	// Prüfe, ob genügend Gold vorhanden ist
	currentGold := character.Vermoegen.Goldstücke
	if currentGold < totalGold {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Gold vorhanden")
		return
	}

	// Prüfe, ob genügend PP vorhanden sind (PP der jeweiligen Fertigkeit) - für neue Fertigkeiten normalerweise 0
	currentPP := 0
	for _, skill := range character.Fertigkeiten {
		if skill.Name == request.Name {
			currentPP = skill.Pp
			break
		}
	}
	// Falls nicht in normalen Fertigkeiten gefunden, prüfe Waffenfertigkeiten
	if currentPP == 0 {
		for _, skill := range character.Waffenfertigkeiten {
			if skill.Name == request.Name {
				currentPP = skill.Pp
				break
			}
		}
	}
	if totalPP > 0 && currentPP < totalPP {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Praxispunkte vorhanden")
		return
	}

	// EP abziehen und Audit-Log erstellen
	newEP := currentEP - totalEP
	if totalEP > 0 {
		var notes string
		if finalLevel > 1 {
			notes = fmt.Sprintf("Fertigkeit '%s' bis Level %d gelernt", request.Name, finalLevel)
		} else {
			notes = fmt.Sprintf("Fertigkeit '%s' gelernt", request.Name)
		}

		err = CreateAuditLogEntry(character.ID, "experience_points", currentEP, newEP, ReasonSkillLearning, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		character.Erfahrungsschatz.EP = newEP
		if err := database.DB.Save(&character.Erfahrungsschatz).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern der Erfahrungspunkte")
			return
		}
	}

	// Gold abziehen und Audit-Log erstellen
	newGold := currentGold - totalGold
	if totalGold > 0 {
		notes := fmt.Sprintf("Gold für Fertigkeit '%s' ausgegeben", request.Name)

		err = CreateAuditLogEntry(character.ID, "gold", currentGold, newGold, ReasonSkillLearning, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		character.Vermoegen.Goldstücke = newGold
		if err := database.DB.Save(&character.Vermoegen).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Vermögens")
			return
		}
	}

	// PP abziehen (falls vorhanden und erforderlich)
	if totalPP > 0 {
		// Suche die richtige Fertigkeit und ziehe PP ab
		for i, skill := range character.Fertigkeiten {
			if skill.Name == request.Name {
				character.Fertigkeiten[i].Pp -= totalPP
				break
			}
		}
		// Falls nicht in normalen Fertigkeiten gefunden, prüfe Waffenfertigkeiten
		for i, skill := range character.Waffenfertigkeiten {
			if skill.Name == request.Name {
				character.Waffenfertigkeiten[i].Pp -= totalPP
				break
			}
		}
	}

	// Erstelle die neue Fertigkeit mit dem finalen Level
	if err := updateOrCreateSkill(&character, request.Name, finalLevel); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Hinzufügen der Fertigkeit: "+err.Error())
		return
	}

	// Charakter speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	// Response für Multi-Level oder Single-Level
	response := gin.H{
		"message":        "Fertigkeit erfolgreich gelernt",
		"skill_name":     request.Name,
		"final_level":    finalLevel,
		"ep_cost":        totalEP,
		"gold_cost":      totalGold,
		"remaining_ep":   newEP,
		"remaining_gold": newGold,
	}

	// Füge Multi-Level-spezifische Informationen hinzu
	if finalLevel > 1 {
		// Erstelle Array der gelernten Level für Kompatibilität
		var levelsLearned []int
		for i := 1; i <= finalLevel; i++ {
			levelsLearned = append(levelsLearned, i)
		}
		response["levels_learned"] = levelsLearned
		response["level_count"] = finalLevel
		response["multi_level"] = true
	}

	c.JSON(http.StatusOK, response)
}

// LearnSkill lernt eine neue Fertigkeit und erstellt Audit-Log-Einträge
func LearnSkill(c *gin.Context) {
	charID := c.Param("id")
	var character models.Char

	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Verwende gsmaster.LernCostRequest direkt
	var request gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Setze Charakter-ID und Action für learning
	request.CharId = character.ID
	request.Action = "learn"
	if request.Type == "" {
		request.Type = "skill" // Default zu skill für Learning
	}

	// 1. Charakter laden
	char, err := loadCharacterForImprovement(request.CharId)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// 2. Skill validieren (für Learning beginnen wir bei Level 0)
	characterClass, skillInfo, currentLevel, err := validateSkillForLearning(char, &request)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Bestimme das finale Level
	finalLevel := request.TargetLevel
	if finalLevel <= 0 {
		finalLevel = 1 // Standard für neue Fertigkeit
	}

	// 3. Kosten berechnen (von Level 0 bis finalLevel)
	response, totalEP, totalGold, totalPP, err := calculateLearningCosts(char, &request, characterClass, skillInfo, currentLevel, finalLevel)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Ressourcen validieren
	err = validateResources(char, request.Name, totalEP, totalGold, totalPP)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 5. Ressourcen abziehen
	newEP, newGold, err := deductResourcesForLearning(char, request.Name, finalLevel, totalEP, totalGold, totalPP)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 6. Skill hinzufügen/erstellen
	if err := updateOrCreateSkill(char, request.Name, finalLevel); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Hinzufügen der Fertigkeit: "+err.Error())
		return
	}

	// 7. Charakter speichern
	if err := database.DB.Save(char).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	// 8. Response erstellen
	responseData := gin.H{
		"message":        "Fertigkeit erfolgreich gelernt",
		"skill_name":     request.Name,
		"final_level":    finalLevel,
		"ep_cost":        totalEP,
		"gold_cost":      totalGold,
		"remaining_ep":   newEP,
		"remaining_gold": newGold,
		"cost_details":   response,
	}

	// Füge Multi-Level-spezifische Informationen hinzu
	if finalLevel > 1 {
		// Erstelle Array der gelernten Level für Kompatibilität
		var levelsLearned []int
		for i := 1; i <= finalLevel; i++ {
			levelsLearned = append(levelsLearned, i)
		}
		responseData["levels_learned"] = levelsLearned
		responseData["level_count"] = finalLevel
		responseData["multi_level"] = true
	}

	c.JSON(http.StatusOK, responseData)
}

// ImproveSkill verbessert eine bestehende Fertigkeit und erstellt Audit-Log-Einträge
// validateSkillForLearning validiert Skill-Namen für neue Fertigkeiten (Learning)
func validateSkillForLearning(char *models.Char, request *gsmaster.LernCostRequest) (string, *models.SkillLearningInfo, int, error) {
	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(char.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviationNewSystem(char.Typ)
	} else {
		characterClass = char.Typ
	}

	// Normalize skill/spell name (trim whitespace, proper case)
	skillName := strings.TrimSpace(request.Name)

	skillInfo, err := models.GetSkillCategoryAndDifficultyNewSystem(skillName, characterClass)
	if err != nil {
		return "", nil, 0, fmt.Errorf("fertigkeit '%s' nicht gefunden oder nicht für Klasse '%s' verfügbar: %v", skillName, characterClass, err)
	}

	// Für Learning starten wir bei Level 0
	currentLevel := 0

	// Prüfe, ob die Fertigkeit bereits existiert
	existingLevel := getCurrentSkillLevel(char, request.Name, "skill")
	if existingLevel > 0 {
		return "", nil, 0, fmt.Errorf("fertigkeit '%s' ist bereits auf Level %d - verwende ImproveSkill stattdessen", request.Name, existingLevel)
	}

	return characterClass, skillInfo, currentLevel, nil
}

// calculateLearningCosts berechnet die Kosten für das Erlernen einer neuen Fertigkeit
func calculateLearningCosts(char *models.Char, request *gsmaster.LernCostRequest, characterClass string, skillInfo *models.SkillLearningInfo, currentLevel, finalLevel int) ([]gsmaster.SkillCostResultNew, int, int, int, error) {
	var response []gsmaster.SkillCostResultNew
	var totalEP, totalGold, totalPP int

	// Loop für jeden Level von 0 bis finalLevel (für neue Fertigkeiten)
	for tempLevel := currentLevel; tempLevel < finalLevel; tempLevel++ {
		nextLevel := tempLevel + 1

		// Erstelle temporären Request für diesen Level
		tempRequest := *request
		tempRequest.CurrentLevel = tempLevel
		tempRequest.TargetLevel = nextLevel

		// Für das erste Level (0->1) ist es ein "learn", für weitere Level "improve"
		if tempLevel == 0 {
			tempRequest.Action = "learn"
		} else {
			tempRequest.Action = "improve"
		}

		// Erstelle cost result structure
		costResult := gsmaster.SkillCostResultNew{
			CharacterID:    fmt.Sprintf("%d", char.ID),
			CharacterClass: characterClass,
			SkillName:      request.Name,
			TargetLevel:    nextLevel,
		}

		// Verwende die gleiche Kostenfunktion wie für Improvements
		err := CalculateSkillImproveCostNewSystem(&tempRequest, &costResult, nextLevel, &tempRequest.UsePP, &tempRequest.UseGold, skillInfo)
		if err != nil {
			return nil, 0, 0, 0, fmt.Errorf("fehler bei der Kostenberechnung: %v", err)
		}

		// für die nächste Runde die PP und Gold reduzieren die zum Lernen genutzt werden sollen
		if costResult.PPUsed > 0 {
			request.UsePP -= costResult.PPUsed
			if request.UsePP < 0 {
				request.UsePP = 0
			}
		}

		if costResult.GoldUsed > 0 {
			request.UseGold -= costResult.GoldUsed
			if request.UseGold < 0 {
				request.UseGold = 0
			}
		}

		response = append(response, costResult)

		// Addiere die Kosten
		totalEP += costResult.EP
		totalGold += costResult.GoldCost
		totalPP += costResult.PPUsed
	}

	return response, totalEP, totalGold, totalPP, nil
}

// deductResourcesForLearning zieht die Ressourcen für das Lernen ab und erstellt Audit-Log-Einträge
func deductResourcesForLearning(char *models.Char, skillName string, finalLevel, totalEP, totalGold, totalPP int) (int, int, error) {
	return deductResourcesWithAuditReason(char, skillName, finalLevel, totalEP, totalGold, totalPP, ReasonSkillLearning)
}

// deductResourcesWithAuditReason zieht EP, Gold und PP ab und erstellt entsprechende Audit-Log-Einträge
func deductResourcesWithAuditReason(char *models.Char, itemName string, finalLevel, totalEP, totalGold, totalPP int, auditReason AuditLogReason) (int, int, error) {
	currentEP := char.Erfahrungsschatz.EP
	currentGold := char.Vermoegen.Goldstücke

	// EP abziehen und Audit-Log erstellen
	newEP := currentEP - totalEP
	if totalEP > 0 {
		var notes string
		if finalLevel > 1 {
			notes = fmt.Sprintf("Fertigkeit '%s' bis Level %d gelernt", itemName, finalLevel)
		} else if auditReason == ReasonSpellLearning {
			notes = fmt.Sprintf("Zauber '%s' gelernt", itemName)
		} else {
			notes = fmt.Sprintf("Fertigkeit '%s' gelernt", itemName)
		}

		err := CreateAuditLogEntry(char.ID, "experience_points", currentEP, newEP, auditReason, 0, notes)
		if err != nil {
			return 0, 0, fmt.Errorf("fehler beim Erstellen des Audit-Log-Eintrags: %v", err)
		}
		char.Erfahrungsschatz.EP = newEP
		if err := database.DB.Save(&char.Erfahrungsschatz).Error; err != nil {
			return 0, 0, fmt.Errorf("fehler beim Speichern der Erfahrungspunkte: %v", err)
		}
	}

	// Gold abziehen und Audit-Log erstellen
	newGold := currentGold - totalGold
	if totalGold > 0 {
		var notes string
		if auditReason == ReasonSpellLearning {
			notes = fmt.Sprintf("Gold für Zauber '%s' ausgegeben", itemName)
		} else {
			notes = fmt.Sprintf("Gold für Fertigkeit '%s' ausgegeben", itemName)
		}

		err := CreateAuditLogEntry(char.ID, "gold", currentGold, newGold, auditReason, 0, notes)
		if err != nil {
			return 0, 0, fmt.Errorf("fehler beim Erstellen des Audit-Log-Eintrags: %v", err)
		}
		char.Vermoegen.Goldstücke = newGold
		if err := database.DB.Save(&char.Vermoegen).Error; err != nil {
			return 0, 0, fmt.Errorf("fehler beim Speichern des Vermögens: %v", err)
		}
	}

	// PP abziehen (falls vorhanden und erforderlich)
	if totalPP > 0 {
		// Suche die richtige Fertigkeit und ziehe PP ab
		for i := range char.Fertigkeiten {
			if char.Fertigkeiten[i].Name == itemName {
				char.Fertigkeiten[i].Pp -= totalPP
				if err := database.DB.Save(&char.Fertigkeiten[i]).Error; err != nil {
					return 0, 0, fmt.Errorf("fehler beim Aktualisieren der Praxispunkte: %v", err)
				}
				break
			}
		}
		// Falls nicht in normalen Fertigkeiten gefunden, prüfe Waffenfertigkeiten
		for i := range char.Waffenfertigkeiten {
			if char.Waffenfertigkeiten[i].Name == itemName {
				char.Waffenfertigkeiten[i].Pp -= totalPP
				if err := database.DB.Save(&char.Waffenfertigkeiten[i]).Error; err != nil {
					return 0, 0, fmt.Errorf("fehler beim Aktualisieren der Praxispunkte: %v", err)
				}
				break
			}
		}
	}

	return newEP, newGold, nil
}

// ImproveSkill verbessert eine bestehende Fertigkeit und erstellt Audit-Log-Einträge
// loadCharacterForImprovement lädt einen Charakter mit allen benötigten Beziehungen
func loadCharacterForImprovement(characterID uint) (*models.Char, error) {
	var char models.Char
	err := database.DB.
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Erfahrungsschatz").
		Preload("Vermoegen").
		Preload("Zauber").
		First(&char, characterID).Error
	return &char, err
}

// validateSkillForImprovement validiert Skill-Namen und ermittelt aktuelle Level
func validateSkillForImprovement(char *models.Char, request *gsmaster.LernCostRequest) (string, *models.SkillLearningInfo, int, error) {
	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(char.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviationNewSystem(char.Typ)
	} else {
		characterClass = char.Typ
	}

	// Normalize skill/spell name (trim whitespace, proper case)
	skillName := strings.TrimSpace(request.Name)

	skillInfo, err := models.GetSkillCategoryAndDifficultyNewSystem(skillName, characterClass)
	if err != nil {
		return "", nil, 0, fmt.Errorf("Fertigkeit '%s' nicht gefunden oder nicht für Klasse '%s' verfügbar: %v", skillName, characterClass, err)
	}

	// Aktuellen Level ermitteln, falls nicht angegeben oder zu klein
	currentLevel := request.CurrentLevel
	if currentLevel <= 3 {
		currentLevel = getCurrentSkillLevel(char, request.Name, "skill")
		if currentLevel == -1 {
			return "", nil, 0, fmt.Errorf("Fertigkeit nicht bei diesem Charakter vorhanden")
		}
		request.CurrentLevel = currentLevel
	}

	return characterClass, skillInfo, currentLevel, nil
}

// calculateImprovementCosts berechnet die Gesamtkosten für Multi-Level-Verbesserungen
func calculateImprovementCosts(char *models.Char, request *gsmaster.LernCostRequest, characterClass string, skillInfo *models.SkillLearningInfo, currentLevel, finalLevel int) ([]gsmaster.SkillCostResultNew, int, int, int, error) {
	var response []gsmaster.SkillCostResultNew
	var totalEP, totalGold, totalPP int

	// Loop für jeden Level von currentLevel bis finalLevel
	tempLevel := currentLevel
	for tempLevel < finalLevel {
		nextLevel := tempLevel + 1
		// Erstelle temporären Request für diesen Level
		tempRequest := *request
		tempRequest.CurrentLevel = tempLevel
		tempRequest.TargetLevel = nextLevel

		// Berechne Kosten für diesen einen Level
		var costResult gsmaster.SkillCostResultNew
		costResult.CharacterID = fmt.Sprintf("%d", char.ID)
		costResult.CharacterClass = characterClass
		costResult.SkillName = request.Name

		err := CalculateSkillImproveCostNewSystem(&tempRequest, &costResult, nextLevel, &tempRequest.UsePP, &tempRequest.UseGold, skillInfo)
		if err != nil {
			return nil, 0, 0, 0, fmt.Errorf("Fehler bei der Kostenberechnung: %v", err)
		}

		// für die nächste Runde die PP und Gold reduzieren die zum Lernen genutzt werden sollen
		if costResult.PPUsed > 0 {
			request.UsePP -= costResult.PPUsed
			if request.UsePP < 0 {
				request.UsePP = 0
			}
		}

		if costResult.GoldUsed > 0 {
			request.UseGold -= costResult.GoldUsed
			if request.UseGold < 0 {
				request.UseGold = 0
			}
		}

		response = append(response, costResult)

		// Addiere die Kosten
		totalEP += costResult.EP
		totalGold += costResult.GoldCost
		totalPP += costResult.PPUsed

		tempLevel++
	}

	return response, totalEP, totalGold, totalPP, nil
}

// validateResources prüft, ob genügend Ressourcen vorhanden sind
func validateResources(char *models.Char, skillName string, totalEP, totalGold, totalPP int) error {
	// Prüfe, ob genügend EP vorhanden sind
	currentEP := char.Erfahrungsschatz.EP
	if currentEP < totalEP {
		return fmt.Errorf("Nicht genügend Erfahrungspunkte vorhanden")
	}

	// Prüfe, ob genügend Gold vorhanden ist
	currentGold := char.Vermoegen.Goldstücke
	if currentGold < totalGold {
		return fmt.Errorf("Nicht genügend Gold vorhanden")
	}

	// Prüfe, ob genügend PP vorhanden sind (PP der jeweiligen Fertigkeit)
	currentPP := 0
	for _, skill := range char.Fertigkeiten {
		if skill.Name == skillName {
			currentPP = skill.Pp
			break
		}
	}
	// Falls nicht in normalen Fertigkeiten gefunden, prüfe Waffenfertigkeiten
	if currentPP == 0 {
		for _, skill := range char.Waffenfertigkeiten {
			if skill.Name == skillName {
				currentPP = skill.Pp
				break
			}
		}
	}
	if totalPP > 0 && currentPP < totalPP {
		return fmt.Errorf("Nicht genügend Praxispunkte vorhanden")
	}

	return nil
}

// deductResources zieht die Kosten von den Charakterressourcen ab
// TODO Fehlerbehandlung (Falls Tabelle nicht vorhanden ist)
func deductResources(char *models.Char, skillName string, currentLevel, finalLevel, totalEP, totalGold, totalPP int) (int, int, error) {
	currentEP := char.Erfahrungsschatz.EP
	currentGold := char.Vermoegen.Goldstücke

	// EP abziehen und Audit-Log erstellen
	newEP := currentEP - totalEP
	if totalEP > 0 {
		// Erstelle Notiz für Multi-Level Improvement
		levelCount := finalLevel - currentLevel
		var notes string
		if levelCount > 1 {
			notes = fmt.Sprintf("Fertigkeit '%s' von %d auf %d verbessert (%d Level)", skillName, currentLevel, finalLevel, levelCount)
		} else {
			notes = fmt.Sprintf("Fertigkeit '%s' von %d auf %d verbessert", skillName, currentLevel, finalLevel)
		}

		err := CreateAuditLogEntry(char.ID, "experience_points", currentEP, newEP, ReasonSkillImprovement, 0, notes)
		if err != nil {
			return newEP, 0, fmt.Errorf("Fehler beim Erstellen des Audit-Log-Eintrags: %v", err)
		}
		char.Erfahrungsschatz.EP = newEP
		if err := database.DB.Save(&char.Erfahrungsschatz).Error; err != nil {
			return newEP, 0, fmt.Errorf("Fehler beim Speichern der Erfahrungspunkte: %v", err)
		}
	}

	// Gold abziehen und Audit-Log erstellen
	newGold := currentGold - totalGold
	if totalGold > 0 {
		notes := fmt.Sprintf("Gold für Verbesserung von '%s' ausgegeben", skillName)

		err := CreateAuditLogEntry(char.ID, "gold", currentGold, newGold, ReasonSkillImprovement, 0, notes)
		if err != nil {
			return newEP, newGold, fmt.Errorf("Fehler beim Erstellen des Audit-Log-Eintrags: %v", err)
		}
		char.Vermoegen.Goldstücke = newGold
		if err := database.DB.Save(&char.Vermoegen).Error; err != nil {
			return newEP, newGold, fmt.Errorf("Fehler beim Speichern des Vermögens: %v", err)
		}
	}

	// PP abziehen wenn verwendet (PP der jeweiligen Fertigkeit)
	if totalPP > 0 {
		// Finde die richtige Fertigkeit und ziehe PP ab
		for i := range char.Fertigkeiten {
			if char.Fertigkeiten[i].Name == skillName {
				char.Fertigkeiten[i].Pp -= totalPP
				if err := database.DB.Save(&char.Fertigkeiten[i]).Error; err != nil {
					return newEP, newGold, fmt.Errorf("Fehler beim Aktualisieren der Praxispunkte: %v", err)
				}
				break
			}
		}
		// Falls nicht in normalen Fertigkeiten gefunden, prüfe Waffenfertigkeiten
		for i := range char.Waffenfertigkeiten {
			if char.Waffenfertigkeiten[i].Name == skillName {
				char.Waffenfertigkeiten[i].Pp -= totalPP
				if err := database.DB.Save(&char.Waffenfertigkeiten[i]).Error; err != nil {
					return newEP, newGold, fmt.Errorf("Fehler beim Aktualisieren der Praxispunkte: %v", err)
				}
				break
			}
		}
	}

	return newEP, newGold, nil
}

func ImproveSkill(c *gin.Context) {
	var request gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// 1. Charakter laden
	char, err := loadCharacterForImprovement(request.CharId)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// 2. Skill validieren und Level ermitteln
	characterClass, skillInfo, currentLevel, err := validateSkillForImprovement(char, &request)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Bestimme das finale Level
	finalLevel := request.TargetLevel
	if finalLevel <= 0 {
		finalLevel = currentLevel + 1
	}

	// 3. Kosten berechnen
	response, totalEP, totalGold, totalPP, err := calculateImprovementCosts(char, &request, characterClass, skillInfo, currentLevel, finalLevel)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Ressourcen validieren
	err = validateResources(char, request.Name, totalEP, totalGold, totalPP)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 5. Ressourcen abziehen
	newEP, newGold, err := deductResources(char, request.Name, currentLevel, finalLevel, totalEP, totalGold, totalPP)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 6. Skill-Level aktualisieren
	if err := updateOrCreateSkill(char, request.Name, finalLevel); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren der Fertigkeit: "+err.Error())
		return
	}

	// 7. Charakter speichern
	if err := database.DB.Save(char).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	// 8. Response erstellen
	responseData := gin.H{
		"message":        "Fertigkeit erfolgreich verbessert",
		"skill_name":     request.Name,
		"from_level":     currentLevel,
		"to_level":       finalLevel,
		"ep_cost":        totalEP,
		"gold_cost":      totalGold,
		"remaining_ep":   newEP,
		"remaining_gold": newGold,
		"cost_details":   response,
	}

	// Füge Multi-Level-spezifische Informationen hinzu
	levelCount := finalLevel - currentLevel
	if levelCount > 1 {
		var levelsLearned []int
		for i := currentLevel + 1; i <= finalLevel; i++ {
			levelsLearned = append(levelsLearned, i)
		}
		responseData["levels_learned"] = levelsLearned
		responseData["level_count"] = levelCount
		responseData["multi_level"] = true
	}

	c.JSON(http.StatusOK, responseData)
}

// ImproveSkillOld is deprecated. Use ImproveSkill instead.
// This function uses the old hardcoded learning cost system.
// ImproveSkillOld verbessert eine bestehende Fertigkeit und erstellt Audit-Log-Einträge
func ImproveSkillOld(c *gin.Context) {
	// Verwende gsmaster.LernCostRequest direkt
	var request gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Hole Charakter über die ID aus dem Request
	var char models.Char
	err := database.DB.
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Erfahrungsschatz").
		Preload("Vermoegen").
		First(&char, request.CharId).Error
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(char.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviationOld(char.Typ)
	} else {
		characterClass = char.Typ
	}

	// Aktuellen Level ermitteln, falls nicht angegeben
	currentLevel := request.CurrentLevel
	if currentLevel <= 0 {
		currentLevel = getCurrentSkillLevel(&char, request.Name, "skill")
		if currentLevel == -1 {
			respondWithError(c, http.StatusBadRequest, "Fertigkeit nicht bei diesem Charakter vorhanden")
			return
		}
		request.CurrentLevel = currentLevel
	}

	// Bestimme das finale Level
	finalLevel := request.TargetLevel
	if finalLevel <= 0 {
		finalLevel = currentLevel + 1
	}

	// Initialisiere Gesamtkosten
	var totalEP, totalGold, totalPP int

	// Loop für jeden Level von currentLevel bis finalLevel
	tempLevel := currentLevel
	for tempLevel < finalLevel {
		nextLevel := tempLevel + 1

		// Erstelle temporären Request für diesen Level
		tempRequest := request
		tempRequest.CurrentLevel = tempLevel
		tempRequest.TargetLevel = nextLevel

		// Berechne Kosten für diesen einen Level
		var costResult gsmaster.SkillCostResultNew
		costResult.CharacterID = fmt.Sprintf("%d", char.ID)
		costResult.CharacterClass = characterClass
		costResult.SkillName = request.Name

		err = gsmaster.GetLernCostNextLevelOld(&tempRequest, &costResult, request.Reward, nextLevel, char.Rasse)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, fmt.Sprintf("Fehler bei Level %d: %v", nextLevel, err))
			return
		}

		// Addiere die Kosten
		totalEP += costResult.EP
		totalGold += costResult.GoldCost
		totalPP += costResult.PPUsed

		tempLevel++
	}

	// Prüfe, ob genügend EP vorhanden sind
	currentEP := char.Erfahrungsschatz.EP
	if currentEP < totalEP {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Erfahrungspunkte vorhanden")
		return
	}

	// Prüfe, ob genügend Gold vorhanden ist
	currentGold := char.Vermoegen.Goldstücke
	if currentGold < totalGold {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Gold vorhanden")
		return
	}

	// Prüfe, ob genügend PP vorhanden sind (PP der jeweiligen Fertigkeit)
	currentPP := 0
	for _, skill := range char.Fertigkeiten {
		if skill.Name == request.Name {
			currentPP = skill.Pp
			break
		}
	}
	// Falls nicht in normalen Fertigkeiten gefunden, prüfe Waffenfertigkeiten
	if currentPP == 0 {
		for _, skill := range char.Waffenfertigkeiten {
			if skill.Name == request.Name {
				currentPP = skill.Pp
				break
			}
		}
	}
	if totalPP > 0 && currentPP < totalPP {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Praxispunkte vorhanden")
		return
	}

	// EP abziehen und Audit-Log erstellen
	newEP := currentEP - totalEP
	if totalEP > 0 {
		// Erstelle Notiz für Multi-Level Improvement
		levelCount := finalLevel - currentLevel
		var notes string
		if levelCount > 1 {
			notes = fmt.Sprintf("Fertigkeit '%s' von %d auf %d verbessert (%d Level)", request.Name, currentLevel, finalLevel, levelCount)
		} else {
			notes = fmt.Sprintf("Fertigkeit '%s' von %d auf %d verbessert", request.Name, currentLevel, finalLevel)
		}

		err = CreateAuditLogEntry(char.ID, "experience_points", currentEP, newEP, ReasonSkillImprovement, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		char.Erfahrungsschatz.EP = newEP
		if err := database.DB.Save(&char.Erfahrungsschatz).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern der Erfahrungspunkte")
			return
		}
	}

	// Gold abziehen und Audit-Log erstellen
	newGold := currentGold - totalGold
	if totalGold > 0 {
		notes := fmt.Sprintf("Gold für Verbesserung von '%s' ausgegeben", request.Name)

		err = CreateAuditLogEntry(char.ID, "gold", currentGold, newGold, ReasonSkillImprovement, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		char.Vermoegen.Goldstücke = newGold
		if err := database.DB.Save(&char.Vermoegen).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Vermögens")
			return
		}
	}

	// PP abziehen wenn verwendet (PP der jeweiligen Fertigkeit)
	if totalPP > 0 {
		// Finde die richtige Fertigkeit und ziehe PP ab
		for i := range char.Fertigkeiten {
			if char.Fertigkeiten[i].Name == request.Name {
				char.Fertigkeiten[i].Pp -= totalPP
				if err := database.DB.Save(&char.Fertigkeiten[i]).Error; err != nil {
					respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren der Praxispunkte")
					return
				}
				break
			}
		}
		// Falls nicht in normalen Fertigkeiten gefunden, prüfe Waffenfertigkeiten
		for i := range char.Waffenfertigkeiten {
			if char.Waffenfertigkeiten[i].Name == request.Name {
				char.Waffenfertigkeiten[i].Pp -= totalPP
				if err := database.DB.Save(&char.Waffenfertigkeiten[i]).Error; err != nil {
					respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren der Praxispunkte")
					return
				}
				break
			}
		}
	}

	// Aktualisiere die Fertigkeit mit dem neuen Level
	if err := updateOrCreateSkill(&char, request.Name, finalLevel); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren der Fertigkeit: "+err.Error())
		return
	}

	// Charakter speichern
	if err := database.DB.Save(&char).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	// Response für Multi-Level oder Single-Level
	response := gin.H{
		"message":        "Fertigkeit erfolgreich verbessert",
		"skill_name":     request.Name,
		"from_level":     currentLevel,
		"to_level":       finalLevel,
		"ep_cost":        totalEP,
		"gold_cost":      totalGold,
		"remaining_ep":   newEP,
		"remaining_gold": newGold,
	}

	// Füge Multi-Level-spezifische Informationen hinzu
	levelCount := finalLevel - currentLevel
	if levelCount > 1 {
		// Erstelle Array der gelernten Level für Kompatibilität
		var levelsLearned []int
		for i := currentLevel + 1; i <= finalLevel; i++ {
			levelsLearned = append(levelsLearned, i)
		}
		response["levels_learned"] = levelsLearned
		response["level_count"] = levelCount
		response["multi_level"] = true
	}

	c.JSON(http.StatusOK, response)
}

// LearnSpellOld is deprecated. Use LearnSpell instead.
// This function uses the old hardcoded learning cost system.
// LearnSpellOld lernt einen neuen Zauber und erstellt Audit-Log-Einträge
func LearnSpellOld(c *gin.Context) {
	charID := c.Param("id")
	var character models.Char

	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	var request LearnSpellRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Berechne Kosten mit GetSkillCost
	costRequest := SkillCostRequest{
		Name:   request.Name,
		Type:   "spell",
		Action: "learn",
	}

	cost, _, _, err := calculateSingleCostOld(&character, &costRequest)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
		return
	}

	// Prüfe, ob genügend EP vorhanden sind
	currentEP := character.Erfahrungsschatz.EP
	if currentEP < cost.Ep {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Erfahrungspunkte vorhanden")
		return
	}

	// EP abziehen und Audit-Log erstellen
	newEP := currentEP - cost.Ep
	if cost.Ep > 0 {
		notes := fmt.Sprintf("Zauber '%s' gelernt", request.Name)
		if request.Notes != "" {
			notes += " - " + request.Notes
		}

		err = CreateAuditLogEntry(character.ID, "experience_points", currentEP, newEP, ReasonSpellLearning, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		character.Erfahrungsschatz.EP = newEP
	}

	// Füge den Zauber zum Charakter hinzu
	if err := addSpellToCharacter(&character, request.Name); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Hinzufügen des Zaubers: "+err.Error())
		return
	}

	// Charakter speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Zauber erfolgreich gelernt",
		"spell_name":   request.Name,
		"ep_cost":      cost.Ep,
		"remaining_ep": newEP,
	})
}

// validateSpellForLearning validiert Zauber-Namen für neue Zauber (Learning)
func validateSpellForLearning(char *models.Char, request *gsmaster.LernCostRequest) (string, *models.SpellLearningInfo, int, error) {
	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(char.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviationNewSystem(char.Typ)
	} else {
		characterClass = char.Typ
	}

	// Normalize spell name (trim whitespace, proper case)
	spellName := strings.TrimSpace(request.Name)

	spellInfo, err := models.GetSpellLearningInfoNewSystem(spellName, characterClass)
	if err != nil {
		return "", nil, 0, fmt.Errorf("zauber '%s' nicht gefunden oder nicht für Klasse '%s' verfügbar: %v", spellName, characterClass, err)
	}

	// Für Learning starten wir bei Level 0
	currentLevel := 0

	// Prüfe, ob der Zauber bereits existiert
	for _, spell := range char.Zauber {
		if spell.Name == request.Name {
			return "", nil, 0, fmt.Errorf("zauber '%s' ist bereits gelernt - Zauber können nicht verbessert werden", request.Name)
		}
	}

	return characterClass, spellInfo, currentLevel, nil
}

// calculateSpellLearningCosts berechnet die Kosten für das Erlernen eines neuen Zaubers
func calculateSpellLearningCosts(char *models.Char, request *gsmaster.LernCostRequest, characterClass string, spellInfo *models.SpellLearningInfo, currentLevel, finalLevel int) ([]gsmaster.SkillCostResultNew, int, error) {
	var response []gsmaster.SkillCostResultNew
	var totalEP int

	// Erstelle cost result structure für Zauber
	costResult := gsmaster.SkillCostResultNew{
		CharacterID:    fmt.Sprintf("%d", char.ID),
		CharacterClass: characterClass,
		SkillName:      request.Name,
		TargetLevel:    finalLevel,
	}

	remainingPP := 0
	remainingGold := 0

	// Verwende die Spell-spezifische Kostenfunktion
	err := calculateSpellLearnCostNewSystem(request, &costResult, &remainingPP, &remainingGold, spellInfo)
	if err != nil {
		return nil, 0, fmt.Errorf("fehler bei der Kostenberechnung: %v", err)
	}

	response = append(response, costResult)
	totalEP = costResult.EP

	// Zauber haben normalerweise keine Gold- oder PP-Kosten
	return response, totalEP, nil
}

// LearnSpell lernt einen neuen Zauber und erstellt Audit-Log-Einträge
func LearnSpell(c *gin.Context) {
	char_ID := c.Param("id")
	/*
		var character models.Char

		if err := character.FirstID(charID); err != nil {
			respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
			return
		}
	*/
	charIDInt, err := strconv.Atoi(char_ID)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Charakter-ID")
		return
	}
	charID := uint(charIDInt)

	var lernRequest gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&lernRequest); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Setze die CharId aus der URL, falls sie nicht im Request enthalten ist
	if lernRequest.CharId == 0 {
		lernRequest.CharId = charID
	}

	// Setze Standard-Werte für Spell Learning falls nicht gesetzt
	if lernRequest.Type == "" {
		lernRequest.Type = "spell"
	}
	if lernRequest.Action == "" {
		lernRequest.Action = "learn"
	}
	if lernRequest.CurrentLevel == 0 && lernRequest.TargetLevel == 0 {
		lernRequest.CurrentLevel = 0 // Zauber sind nicht gelernt
		lernRequest.TargetLevel = 1  // Zauber werden auf Level 1 gelernt
	}

	// 1. Charakter laden
	char, err := loadCharacterForImprovement(lernRequest.CharId)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// 2. Zauber validieren (für Learning beginnen wir bei Level 0)
	characterClass, spellInfo, currentLevel, err := validateSpellForLearning(char, &lernRequest)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	finalLevel := 1 // Zauber werden immer auf Level 1 gelernt

	// 3. Kosten berechnen (von Level 0 bis 1)
	response, totalEP, err := calculateSpellLearningCosts(char, &lernRequest, characterClass, spellInfo, currentLevel, finalLevel)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Ressourcen validieren (nur EP für Zauber)
	err = validateResources(char, lernRequest.Name, totalEP, 0, 0) // Gold=0, PP=0 für Zauber
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// 5. Ressourcen abziehen
	newEP, _, err := deductResourcesWithAuditReason(char, lernRequest.Name, 1, totalEP, 0, 0, ReasonSpellLearning)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 6. Zauber hinzufügen
	if err := addSpellToCharacter(char, lernRequest.Name); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Hinzufügen des Zaubers: "+err.Error())
		return
	}

	// 7. Charakter speichern
	if err := database.DB.Save(char).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	// 8. Response erstellen (kompatibel mit alter Version)
	responseData := gin.H{
		"message":      "Zauber erfolgreich gelernt",
		"spell_name":   lernRequest.Name,
		"ep_cost":      totalEP,
		"remaining_ep": newEP,
		"cost_details": response,
	}

	c.JSON(http.StatusOK, responseData)
}

// GetRewardTypesOld is deprecated. Use GetRewardTypes instead.
// This function provides hardcoded reward type mappings.
// GetRewardTypesOld liefert verfügbare Belohnungsarten für ein bestimmtes Lernszenario
func GetRewardTypesOld(c *gin.Context) {
	characterID := c.Param("id")
	learningType := c.Query("learning_type") // 'improve', 'learn', 'spell'
	skillName := c.Query("skill_name")
	skillType := c.Query("skill_type") // 'skill', 'weapon', 'spell'

	// Basis-Belohnungsarten
	rewardTypes := []gin.H{}

	// Je nach Lerntyp verschiedene Belohnungsarten anbieten
	switch learningType {
	case "learn":
		// Neue Fertigkeit lernen - noGold Belohnung verfügbar
		rewardTypes = append(rewardTypes,
			gin.H{"value": "default", "label": "Standard (EP + Gold)", "description": "Normale EP- und Goldkosten"},
			gin.H{"value": "noGold", "label": "Ohne Gold (nur EP)", "description": "Keine Goldkosten, nur EP als Belohnung"},
		)

	case "spell":
		// Zauber lernen - halveepnoGold verfügbar
		rewardTypes = append(rewardTypes,
			gin.H{"value": "default", "label": "Standard (EP)", "description": "Normale EP-Kosten"},
			gin.H{"value": "halveepnoGold", "label": "Halbe EP ohne Gold", "description": "Halbe EP-Kosten, kein Gold als Belohnung"},
		)

	case "improve":
		// Fertigkeit verbessern - halveepnoGold verfügbar
		rewardTypes = append(rewardTypes,
			gin.H{"value": "default", "label": "Standard (EP + Gold)", "description": "Normale EP- und Goldkosten"},
			gin.H{"value": "halveepnoGold", "label": "Halbe EP ohne Gold", "description": "Halbe EP-Kosten, kein Gold als Belohnung"},
		)

		// Spezielle Optionen für bestimmte Fertigkeiten
		if skillType == "weapon" {
			// Waffenfertigkeiten könnten spezielle Trainingsmethoden haben
			rewardTypes = append(rewardTypes,
				gin.H{"value": "training", "label": "Training mit Meister", "description": "Intensives Training mit einem Waffenmeister"},
			)
		}
	default:
	}

	c.JSON(http.StatusOK, gin.H{
		"reward_types":  rewardTypes,
		"learning_type": learningType,
		"skill_name":    skillName,
		"skill_type":    skillType,
		"character_id":  characterID,
	})
}

// GetAvailableSkillsNewSystem gibt alle verfügbaren Fertigkeiten mit Lernkosten zurück (POST mit LernCostRequest)
func GetAvailableSkillsNewSystem(c *gin.Context) {
	characterID := c.Param("id")

	// Parse LernCostRequest aus POST body
	var baseRequest gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&baseRequest); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	var character models.Char
	if err := database.DB.Preload("Fertigkeiten").Preload("Erfahrungsschatz").Preload("Vermoegen").First(&character, characterID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Hole alle verfügbaren Fertigkeiten aus der gsmaster Datenbank, aber filtere Placeholder aus
	var allSkills []models.Skill

	allSkills, err := models.SelectSkills("", "")
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve skills from gsmaster")
		return
	}
	/*if err := database.DB.Where("name != ?", "Placeholder").Find(&allSkills).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve skills")
		return
	}
	*/

	// Erstelle eine Map der bereits gelernten Fertigkeiten
	learnedSkills := make(map[string]bool)
	for _, skill := range character.Fertigkeiten {
		learnedSkills[skill.Name] = true
	}

	// Organisiere Fertigkeiten nach Kategorien
	skillsByCategory := make(map[string][]gin.H)

	for _, skill := range allSkills {
		// Überspringe bereits gelernte Fertigkeiten
		if learnedSkills[skill.Name] {
			continue
		}
		// Überspringe Placeholder-Fertigkeiten (zusätzliche Sicherheit)
		if skill.Name == "Placeholder" {
			continue
		}

		// Erstelle LernCostRequest für diese Fertigkeit basierend auf der Basis-Anfrage
		request := baseRequest
		request.CharId = character.ID
		request.Name = skill.Name
		request.CurrentLevel = 0 // Nicht gelernt
		request.TargetLevel = 1  // Auf Level 1 lernen
		request.Type = "skill"
		request.Action = "learn"

		// Erstelle SkillCostResultNew
		levelResult := gsmaster.SkillCostResultNew{
			CharacterID:    fmt.Sprintf("%d", character.ID),
			CharacterClass: getCharacterClassOld(&character),
			SkillName:      skill.Name,
			TargetLevel:    1,
		}

		remainingPP := request.UsePP
		remainingGold := request.UseGold

		// Hole die vollständigen Skill-Informationen für die Kostenberechnung
		skillLearningInfo, err := models.GetSkillCategoryAndDifficultyNewSystem(skill.Name, getCharacterClassOld(&character))
		if err != nil {
			// Fallback für unbekannte Skills
			skillLearningInfo = &models.SkillLearningInfo{
				SkillName:    skill.Name,
				CategoryName: skill.Category,
				LearnCost:    50, // Standard-Lernkosten
			}
		}

		// Berechne Lernkosten mit calculateSkillLearnCostNewSystem
		err = calculateSkillLearnCostNewSystem(&request, &levelResult, &remainingPP, &remainingGold, skillLearningInfo)
		epCost := 10000   // Fallback-Wert
		goldCost := 50000 // Fallback-Wert
		if err == nil {
			epCost = levelResult.EP
			goldCost = levelResult.GoldCost
		}

		skillInfo := gin.H{
			"name":     skill.Name,
			"epCost":   epCost,
			"goldCost": goldCost,
		}

		category := skill.Category
		if category == "" {
			category = "Sonstige"
		}

		skillsByCategory[category] = append(skillsByCategory[category], skillInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"skills_by_category": skillsByCategory,
	})
}

// GetAvailableSpellsNewSystem gibt alle verfügbaren Zauber mit Lernkosten zurück (POST mit LernCostRequest)
func GetAvailableSpellsNewSystem(c *gin.Context) {
	//characterID := c.Param("id")

	// Parse LernCostRequest aus POST body
	var baseRequest gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&baseRequest); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	var character models.Char
	if err := database.DB.Preload("Zauber").Preload("Erfahrungsschatz").Preload("Vermoegen").First(&character, baseRequest.CharId).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	charakteClass := getCharacterClassOld(&character)
	// Hole alle verfügbaren Zauber aus der gsmaster Datenbank, aber filtere Placeholder aus
	var allSpells []models.Spell

	allSpells, err := models.SelectSpells("", "")
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve spells from gsmaster")
		return
	}

	// Erstelle eine Map der bereits gelernten Zauber
	learnedSpells := make(map[string]bool)
	for _, spell := range character.Zauber {
		learnedSpells[spell.Name] = true
	}

	// Organisiere Zauber nach Schulen (analog zu Kategorien bei Fertigkeiten)
	spellsBySchool := make(map[string][]gin.H)

	for _, spell := range allSpells {
		// Überspringe bereits gelernte Zauber
		if learnedSpells[spell.Name] {
			continue
		}
		// Überspringe Placeholder-Zauber (zusätzliche Sicherheit)
		if spell.Name == "Placeholder" {
			continue
		}

		// Erstelle LernCostRequest für diesen Zauber basierend auf der Basis-Anfrage
		request := baseRequest
		request.CharId = character.ID
		request.Name = spell.Name
		request.CurrentLevel = 0 // Nicht gelernt
		request.TargetLevel = 1  // Auf Level 1 lernen
		request.Type = "spell"
		request.Action = "learn"

		// Erstelle SkillCostResultNew
		levelResult := gsmaster.SkillCostResultNew{
			CharacterID:    fmt.Sprintf("%d", character.ID),
			CharacterClass: charakteClass,
			SkillName:      spell.Name,
			TargetLevel:    1,
		}

		remainingPP := request.UsePP
		remainingGold := request.UseGold

		// Hole die vollständigen Spell-Informationen für die Kostenberechnung
		spellLearningInfo, err := models.GetSpellLearningInfoNewSystem(spell.Name, charakteClass)
		if err != nil {
			// Fallback für unbekannte Zauber
			spellLearningInfo = &models.SpellLearningInfo{
				SpellName:  spell.Name,
				SpellLevel: spell.Stufe,
				SchoolName: spell.Category,
				LERequired: 20, // Standard-Lernkosten für Zauber
			}
		}

		// Berechne Lernkosten mit calculateSpellLearnCostNewSystem
		err = calculateSpellLearnCostNewSystem(&request, &levelResult, &remainingPP, &remainingGold, spellLearningInfo)
		epCost := 10000   // Fallback-Wert
		goldCost := 50000 // Fallback-Wert
		if err == nil {
			epCost = levelResult.EP
			goldCost = levelResult.GoldCost
		}

		spellInfo := gin.H{
			"name":     spell.Name,
			"level":    spell.Stufe,
			"epCost":   epCost,
			"goldCost": goldCost,
		}

		school := spell.Category
		if school == "" {
			school = "Sonstige"
		}

		spellsBySchool[school] = append(spellsBySchool[school], spellInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"spells_by_school": spellsBySchool,
	})
}

// GetSpellDetails gibt detaillierte Informationen zu einem bestimmten Zauber zurück
func GetSpellDetails(c *gin.Context) {
	spellName := c.Query("name")
	if spellName == "" {
		respondWithError(c, http.StatusBadRequest, "Zaubername ist erforderlich")
		return
	}

	// Lade den Zauber aus der Datenbank
	var spell models.Spell
	if err := database.DB.Where("name = ?", spellName).First(&spell).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Zauber nicht gefunden")
		return
	}

	// Erstelle Response mit allen verfügbaren Details
	spellDetails := gin.H{
		"id":                spell.ID,
		"name":              spell.Name,
		"beschreibung":      spell.Beschreibung,
		"level":             spell.Stufe,
		"bonus":             spell.Bonus,
		"ap":                spell.AP,
		"art":               spell.Art,
		"zauberdauer":       spell.Zauberdauer,
		"reichweite":        spell.Reichweite,
		"wirkungsziel":      spell.Wirkungsziel,
		"wirkungsbereich":   spell.Wirkungsbereich,
		"wirkungsdauer":     spell.Wirkungsdauer,
		"ursprung":          spell.Ursprung,
		"category":          spell.Category,
		"learning_category": spell.LearningCategory,
		"quelle":            spell.Quelle,
		"page_number":       spell.PageNumber,
		"game_system":       spell.GameSystem,
	}

	c.JSON(http.StatusOK, gin.H{
		"spell": spellDetails,
	})
}

// GetAvailableSkillsOld is deprecated. Use GetAvailableSkillsNewSystem instead.
// This function uses the old hardcoded learning cost system.
// GetAvailableSkillsOld gibt alle verfügbaren Fertigkeiten mit Lernkosten zurück
func GetAvailableSkillsOld(c *gin.Context) {
	characterID := c.Param("id")
	rewardType := c.Query("reward_type")

	var character models.Char
	if err := database.DB.Preload("Fertigkeiten").Preload("Erfahrungsschatz").Preload("Vermoegen").First(&character, characterID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Hole alle verfügbaren Fertigkeiten aus der gsmaster Datenbank, aber filtere Placeholder aus
	var allSkills []models.Skill

	allSkills, err := models.SelectSkills("", "")
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve skills from gsmaster")
		return
	}
	/*if err := database.DB.Where("name != ?", "Placeholder").Find(&allSkills).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve skills")
		return
	}
	*/

	// Erstelle eine Map der bereits gelernten Fertigkeiten
	learnedSkills := make(map[string]bool)
	for _, skill := range character.Fertigkeiten {
		learnedSkills[skill.Name] = true
	}

	// Organisiere Fertigkeiten nach Kategorien
	skillsByCategory := make(map[string][]gin.H)

	for _, skill := range allSkills {
		// Überspringe bereits gelernte Fertigkeiten
		if learnedSkills[skill.Name] {
			continue
		}

		// Überspringe Placeholder-Fertigkeiten (zusätzliche Sicherheit)
		if skill.Name == "Placeholder" {
			continue
		}

		// Berechne Lernkosten mit GetLernCostNextLevel
		epCost, goldCost := calculateSkillLearningCostsOld(skill, character, rewardType)

		skillInfo := gin.H{
			"name":     skill.Name,
			"epCost":   epCost,
			"goldCost": goldCost,
		}

		category := skill.Category
		if category == "" {
			category = "Sonstige"
		}

		skillsByCategory[category] = append(skillsByCategory[category], skillInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"skills_by_category": skillsByCategory,
	})
}

// calculateSkillLearningCostsOld is deprecated. Use calculateSkillLearnCostNewSystem instead.
// This function uses the old hardcoded learning cost system.
// calculateSkillLearningCostsOld berechnet die EP- und Goldkosten für das Lernen einer Fertigkeit mit GetLernCostNextLevel
func calculateSkillLearningCostsOld(skill models.Skill, character models.Char, rewardType string) (int, int) {
	// Erstelle LernCostRequest für das Lernen (Level 0 -> 1)
	var rewardTypePtr *string
	if rewardType != "" && rewardType != "default" {
		rewardTypePtr = &rewardType
	}

	request := gsmaster.LernCostRequest{
		CharId:       character.ID,
		Name:         skill.Name,
		CurrentLevel: 0, // Nicht gelernt
		TargetLevel:  1, // Auf Level 1 lernen
		Type:         "skill",
		Action:       "learn",
		UsePP:        0,
		UseGold:      0,
		Reward:       rewardTypePtr,
	}

	// Erstelle SkillCostResultNew
	costResult := gsmaster.SkillCostResultNew{
		CharacterID:    fmt.Sprintf("%d", character.ID),
		CharacterClass: getCharacterClassOld(&character),
		SkillName:      skill.Name,
		Category:       skill.Category,
		Difficulty:     skill.Difficulty,
		TargetLevel:    1,
	}

	// Berechne Kosten mit GetLernCostNextLevel
	err := gsmaster.GetLernCostNextLevelOld(&request, &costResult, rewardTypePtr, 1, character.Typ)
	if err != nil {
		// Fallback zu Standard-Kosten bei Fehler
		epCost := 100
		goldCost := 50

		return epCost, goldCost
	}

	return costResult.EP, costResult.GoldCost
}
