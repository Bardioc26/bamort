package character

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/logger"
	"bamort/models"
	"sort"
	"strconv"
	"strings"
	"time"

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
	logger.Warn("HTTP Fehler %d: %s", status, message)
	c.JSON(status, gin.H{"error": message})
}

func ListCharacters(c *gin.Context) {
	logger.Debug("ListCharacters aufgerufen")

	type AllCharacters struct {
		SelfOwned []models.CharList `json:"self_owned"`
		Others    []models.CharList `json:"others"`
	}
	allCharacters := AllCharacters{}

	logger.Debug("Lade Charaktere aus der Datenbank...")
	//if err := database.DB.Find(&characters).Error; err != nil {
	listOfChars, err := models.FindCharListByUserID(c.GetUint("userID"))
	if err != nil {

		logger.Error("Fehler beim Laden der Charaktere: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve characters")
		return
	}

	logger.Debug("Gefundene Charaktere: %d", len(listOfChars))
	allCharacters.SelfOwned = listOfChars

	listPublic, err := models.FindPublicCharList()

	if err != nil {
		logger.Error("Fehler beim Laden der öffentlichen Charaktere: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve public characters")
		return
	}
	allCharacters.Others = listPublic

	logger.Info("Charakterliste erfolgreich geladen: %d Charaktere", len(listOfChars))
	c.JSON(http.StatusOK, allCharacters)
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
		userID := c.GetUint("userID")
		character.Vermoegen = models.Vermoegen{
			BamortCharTrait: models.BamortCharTrait{
				CharacterID: character.ID,
				UserID:      userID,
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
	// Parse LernCostRequest aus POST body
	var baseRequest gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&baseRequest); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// For character creation (char_id = 0), we don't need to load an existing character
	var character models.Char
	learnedSkills := make(map[string]bool)

	if baseRequest.CharId != 0 {
		// Load existing character and their learned skills
		if err := database.DB.Preload("Fertigkeiten").Preload("Erfahrungsschatz").Preload("Vermoegen").First(&character, baseRequest.CharId).Error; err != nil {
			respondWithError(c, http.StatusNotFound, "Character not found")
			return
		}

		// Create map of learned skills for existing character
		for _, skill := range character.Fertigkeiten {
			learnedSkills[skill.Name] = true
		}
	}
	// For character creation (char_id = 0), learnedSkills remains empty

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
		request.CharId = baseRequest.CharId // Use the char_id from the request (0 for character creation)
		request.Name = skill.Name
		request.CurrentLevel = 0 // Nicht gelernt
		request.TargetLevel = 1  // Auf Level 1 lernen
		request.Type = "skill"
		request.Action = "learn"

		// Erstelle SkillCostResultNew
		characterClass := ""
		characterID := "0"

		if baseRequest.CharId != 0 {
			// Use existing character data
			characterID = fmt.Sprintf("%d", character.ID)
			characterClass = getCharacterClassOld(&character)
		}
		// For character creation, we don't have a character class yet, use empty string

		levelResult := gsmaster.SkillCostResultNew{
			CharacterID:    characterID,
			CharacterClass: characterClass,
			SkillName:      skill.Name,
			TargetLevel:    1,
		}

		remainingPP := request.UsePP
		remainingGold := request.UseGold

		// Hole die vollständigen Skill-Informationen für die Kostenberechnung
		skillLearningInfo, err := models.GetSkillCategoryAndDifficultyNewSystem(skill.Name, characterClass)
		if err != nil {
			// Fallback für unbekannte Skills
			skillLearningInfo = &models.SkillLearningInfo{
				SkillName:    skill.Name,
				CategoryName: skill.Category,
				LearnCost:    50, // Standard-Lernkosten
			}
		}

		// For character creation (CharId = 0), use learning costs instead of improvement costs
		var epCost, goldCost int
		if baseRequest.CharId == 0 {
			// Character creation: use basic learning costs from skillLearningInfo
			learnCost := skillLearningInfo.LearnCost
			if learnCost == 0 {
				learnCost = 50 // Default learning cost
			}

			// For character creation, costs are much lower - just the basic learning cost
			epCost = learnCost * 2   // Simple formula: learning cost * 2 for EP
			goldCost = learnCost * 5 // Simple formula: learning cost * 5 for gold
		} else {
			// Existing character improvement: use the full system
			err = calculateSkillLearnCostNewSystem(&request, &levelResult, &remainingPP, &remainingGold, skillLearningInfo)
			epCost = 10000   // Fallback-Wert for improvements
			goldCost = 50000 // Fallback-Wert for improvements
			if err == nil {
				epCost = levelResult.EP
				goldCost = levelResult.GoldCost
			}
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

// getCharacterClassCode converts a character class name to its code using the database
func getCharacterClassCode(className string) (string, error) {
	var characterClass models.CharacterClass
	err := characterClass.FirstByName(className)
	if err != nil {
		err := characterClass.FirstByCode(className)
		if err != nil {
			return "", fmt.Errorf("character class '%s' not found: %w", className, err)
		}
	}
	return characterClass.Code, nil
}

func GetAvailableSpellsForCreation(c *gin.Context) {
	var request struct {
		CharacterClass string `json:"characterClass" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warn("HTTP Fehler 400: Ungültige Anfrageparameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ungültige Anfrageparameter",
			"details": err.Error(),
		})
		return
	}
	logger.Info("GetAvailableSpellsForCreation - CharacterClass: %s", request.CharacterClass)

	// Convert character class name to code
	characterClassCode, err := getCharacterClassCode(request.CharacterClass)
	if err != nil {
		logger.Error("Fehler beim Konvertieren der Charakterklasse: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Unbekannte Charakterklasse: %s", request.CharacterClass),
		})
		return
	}

	// Get all available spells with their learning costs
	spellsByCategory, err := GetAllSpellsWithLE(characterClassCode, 2)
	if err != nil {
		logger.Error("Fehler beim Abrufen der Zauber: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fehler beim Abrufen der Zauber",
		})
		return
	}

	logger.Info("GetAvailableSpellsForCreation - Gefundene Kategorien: %d", len(spellsByCategory))

	if len(spellsByCategory) == 0 {
		logger.Warn("GetAvailableSpellsForCreation - Keine Zauber für Klasse %s gefunden", request.CharacterClass)
		c.JSON(http.StatusNotFound, gin.H{
			"spells_by_category": map[string][]gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"spells_by_category": spellsByCategory,
	})
}

// GetAvailableSkillsForCreation returns skills with learning costs for character creation
func GetAvailableSkillsForCreation(c *gin.Context) {
	var request struct {
		CharacterClass string `json:"characterClass" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Warn("HTTP Fehler 400: Ungültige Anfrageparameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Ungültige Anfrageparameter",
			"details": err.Error(),
		})
		return
	}

	logger.Info("GetAvailableSkillsForCreation - CharacterClass: %s", request.CharacterClass)

	// Get all available skills with their learning costs
	skillsByCategory, err := GetAllSkillsWithLE()
	if err != nil {
		logger.Error("Fehler beim Abrufen der Fertigkeiten: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Fehler beim Abrufen der Fertigkeiten",
		})
		return
	}

	logger.Info("GetAvailableSkillsForCreation - Gefundene Kategorien: %d", len(skillsByCategory))

	c.JSON(http.StatusOK, gin.H{
		"skills_by_category": skillsByCategory,
	})
}

func GetAllSpellsWithLE(characterClass string, maxLevel int) (map[string][]gin.H, error) {
	// Create mapping of character classes to allowed learning categories
	allowedCategories := getCharacterClassSpellSchoolMapping()
	allowedLearningCategories := getCharacterClassSpellLearningCategoriesMapping()

	// Check if character class has allowed spell schools

	allowedSchools, exists := allowedCategories[characterClass]
	if !exists {
		return map[string][]gin.H{}, nil // Return empty map if class can't learn spells
	}
	allowedSpellType, exists := allowedLearningCategories[characterClass]
	if !exists {
		return map[string][]gin.H{}, nil // Return empty map if class can't learn spells
	}
	// Extract allowed school names and spell types from maps
	var allowedSchoolNames []string
	for school, allowed := range allowedSchools {
		if allowed {
			allowedSchoolNames = append(allowedSchoolNames, school)
		}
	}

	var allowedSpellTypeNames []string
	for spellType, allowed := range allowedSpellType {
		if allowed {
			allowedSpellTypeNames = append(allowedSpellTypeNames, spellType)
		}
	}

	// Get all spells from database with level filter
	var spells []models.Spell
	err := database.DB.Where("stufe <= ? AND category in (?) and learning_category in (?)  AND category IS NOT NULL AND learning_category IS NOT NULL", maxLevel, allowedSchoolNames, allowedSpellTypeNames).Find(&spells).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch spells: %w", err)
	}

	// Group spells by category (using LearningCategory field)
	spellsByCategory := make(map[string][]gin.H)

	for _, spell := range spells {
		// Check if this character class can learn this spell school
		//if !allowedSchools[spell.Category] || !allowedSpellType[spell.LearningCategory] {
		//	continue // Skip spells from schools this class can't learn
		//}

		// Calculate learning cost for this spell
		leCost := getSpellLECost(spell.Stufe)

		spellData := gin.H{
			"id":                spell.ID,
			"name":              spell.Name,
			"level":             spell.Stufe,
			"school":            spell.Category,         // Display category
			"learning_category": spell.LearningCategory, // Internal category for learning
			"description":       spell.Beschreibung,
			"le_cost":           leCost,
			"ap":                spell.AP,
			"art":               spell.Art,
			"zauberdauer":       spell.Zauberdauer,
			"reichweite":        spell.Reichweite,
			"wirkungsziel":      spell.Wirkungsziel,
			"wirkungsbereich":   spell.Wirkungsbereich,
			"wirkungsdauer":     spell.Wirkungsdauer,
		}

		// Use LearningCategory for grouping
		category := spell.LearningCategory
		spellsByCategory[category] = append(spellsByCategory[category], spellData)
	}

	// Sort spells within each category by level, then by name
	for category := range spellsByCategory {
		spells := spellsByCategory[category]
		sort.Slice(spells, func(i, j int) bool {
			levelI := spells[i]["level"].(int)
			levelJ := spells[j]["level"].(int)
			if levelI != levelJ {
				return levelI < levelJ
			}
			nameI := spells[i]["name"].(string)
			nameJ := spells[j]["name"].(string)
			return nameI < nameJ
		})
		spellsByCategory[category] = spells
	}

	return spellsByCategory, nil
}
func getCharacterClassSpellLearningCategoriesMapping() map[string]map[string]bool {
	return map[string]map[string]bool{
		"Ma": { // Magier
			"Spruch": true,
			//"Salz":      true,
			//"Runenstab": true,
		},
		"Hx": { // Hexer
			"Spruch": true,
			"Salz":   true,
			//"Runenstab": true,
		},
		"Dr": { // Druide
			"Spruch": true,
			//"Salz":      true,
			//"Runenstab": true,
			"Dweomer": true,
		},
		"Sc": { // Schamane
			"Spruch": true,
			//"Salz":      true,
			//"Runenstab": true,
			"Dweomer": true,
		},
		"PB": { // Priester Beschützer
			"Salz":      true,
			"Wundertat": true,
		},
		"PS": { // Priester Streiter
			"Salz":      true,
			"Wundertat": true,
		},
		"Ba": { // Barde
			"Lied": true,
		},
		"Or": { // Ordenskrieger
			"Wundertat": true,
		},
	}
}

// getCharacterClassSpellSchoolMapping returns the mapping of character classes to allowed spell schools
func getCharacterClassSpellSchoolMapping() map[string]map[string]bool {
	return map[string]map[string]bool{
		"Ma": { // Magier
			"Beherrschen": true,
			"Bewegen":     true,
			"Dweomer":     true,
			"Erkennen":    true,
			"Erschaffen":  true,
			"Verändern":   true,
			"Zerstören":   true,
		},
		"Hx": { // Hexer
			"Beherrschen": true,
			"Zerstören":   true,
			"Erkennen":    true,
			"Verändern":   true,
			"Erschaffen":  true,
			"Bewegen":     true,
			"Formen":      true,
		},
		"Dr": { // Druide
			"Bewegen":    true,
			"Erkennen":   true,
			"Erschaffen": true,
			"Verändern":  true,
		},
		"Sc": { // Schamane
			"Beherrschen": true,
			"Erkennen":    true,
			"Verändern":   true,
		},
		"PB": { // Priester Beschützer
			"Dweomer":   true,
			"Erkennen":  true,
			"Verändern": true,
		},
		"PS": { // Priester Streiter
			"Dweomer":   true,
			"Erkennen":  true,
			"Verändern": true,
		},
		"Ba": { // Barde
			"Beherrschen": true,
			"Dweomer":     true,
			"Erkennen":    true,
		},
		"Or": { // Ordenskrieger
			"Dweomer":  true,
			"Erkennen": true,
		},
	}
}

// getSpellLECost returns the learning cost in LE for a given spell level from the database
func getSpellLECost(level int) int {
	var spellLECost models.SpellLevelLECost

	// Query the database for the LE cost for this level
	err := database.DB.Where("level = ? AND game_system = ?", level, "midgard").First(&spellLECost).Error
	if err != nil {
		// If not found in database, fall back to standard Midgard costs
		spellLECosts := map[int]int{
			1:  1,
			2:  2,
			3:  3,
			4:  4,
			5:  5,
			6:  6,
			7:  8,
			8:  10,
			9:  12,
			10: 15,
			11: 18,
			12: 21,
		}

		if cost, exists := spellLECosts[level]; exists {
			return cost
		}

		// Final fallback for unknown levels
		return level
	}

	return spellLECost.LERequired
}

func GetAllSkillsWithLE() (map[string][]gin.H, error) {
	// Get all skill categories from database
	var skillCategories []models.SkillCategory
	if err := database.DB.Find(&skillCategories).Error; err != nil {
		return nil, err
	}

	skillsByCategory := make(map[string][]gin.H)

	// For each category, find all skills that can be learned in that category
	for _, category := range skillCategories {
		skillsByCategory[category.Name] = []gin.H{}

		// Query all skill-category-difficulty combinations for this category
		var skillCategoryDifficulties []models.SkillCategoryDifficulty
		err := database.DB.Preload("Skill").Preload("SkillDifficulty").
			Where("skill_category_id = ?", category.ID).
			Find(&skillCategoryDifficulties).Error

		if err != nil {
			continue // Skip this category if there's an error
		}

		// For each skill in this category, add it with its LE cost and difficulty
		for _, scd := range skillCategoryDifficulties {
			// Skip Placeholder skills
			if category.Name == "Unbekannt" || scd.Skill.Name == "Placeholder" || scd.Skill.InnateSkill {
				continue
			}

			skillInfo := gin.H{
				"name":       scd.Skill.Name,
				"leCost":     scd.LearnCost,
				"difficulty": scd.SkillDifficulty.Name,
			}

			skillsByCategory[category.Name] = append(skillsByCategory[category.Name], skillInfo)
		}
	}

	// Add weapon skills to "Kampf" category
	weaponSkills, err := GetWeaponSkillsWithLE()
	if err == nil {
		if _, exists := skillsByCategory["Waffen"]; !exists {
			skillsByCategory["Waffen"] = []gin.H{}
		}
		skillsByCategory["Waffen"] = append(skillsByCategory["Waffen"], weaponSkills...)
	}

	return skillsByCategory, nil
}

// GetWeaponSkillsWithLE returns all weapon skills with their learning costs
func GetWeaponSkillsWithLE() ([]gin.H, error) {
	// Query weapon skills from gsm_weaponskills table
	var weaponSkills []struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}

	err := database.DB.Table("gsm_weaponskills").
		Select("id, name").
		Where("name IS NOT NULL AND name != ''").
		Find(&weaponSkills).Error

	if err != nil {
		return nil, err
	}

	var result []gin.H
	for _, weapon := range weaponSkills {
		// Get learning cost for this weapon skill
		// For now, use default costs, but this could be enhanced to query from a learning costs table
		leCost := getDefaultWeaponSkillCost(weapon.Name)

		skillInfo := gin.H{
			"name":       weapon.Name,
			"leCost":     leCost,
			"difficulty": "Normal", // Default difficulty for weapon skills
			"type":       "weapon", // Mark as weapon skill
		}

		result = append(result, skillInfo)
	}

	return result, nil
}

// getDefaultWeaponSkillCost returns default learning costs for weapon skills
func getDefaultWeaponSkillCost(weaponName string) int {
	// Basic weapon skill costs - these could be made configurable or stored in database
	switch weaponName {
	case "Schilde":
		return 1 // Shields are easier to learn
	case "Einhandschwerter", "Zweihandschwerter", "Fechtwaffen":
		return 4 // Sword skills are more expensive
	case "Bögen", "Armbrüste":
		return 3 // Ranged weapons
	case "Spießwaffen", "Stielwurfwaffen":
		return 2 // Polearms and throwing weapons
	default:
		return 2 // Default cost for other weapon skills
	}
}

// GetAllSkillsWithLearningCosts returns all skills with their basic learning costs for all possible categories
func GetAllSkillsWithLearningCosts(characterClass string) (map[string][]gin.H, error) {
	skills, err := models.SelectSkills("", "")
	if err != nil {
		return nil, err
	}

	skillsByCategory := make(map[string][]gin.H)

	// Define all possible categories for skills
	allCategories := []string{"Alltag", "Kampf", "Körper", "Sozial", "Wissen", "Halbwelt", "Unterwelt", "Freiland", "Sonstige"}

	for _, skill := range skills {
		// Skip Placeholder skills
		if skill.Name == "Placeholder" {
			continue
		}

		// First, always add to the skill's original category
		originalCategory := skill.Category
		if originalCategory == "" {
			originalCategory = "Sonstige"
		}

		// Try to get the best category and learning cost for this skill and character class
		bestCategory, difficulty, err := gsmaster.FindBestCategoryForSkillLearningOld(skill.Name, characterClass)

		var learnCost int
		if err == nil && bestCategory != "" {
			// Use the difficulty as a basis for learning cost
			switch difficulty {
			case "Leicht":
				learnCost = 1
			case "Normal":
				learnCost = 2
			case "Schwer":
				learnCost = 4
			case "Sehr Schwer":
				learnCost = 10
			default:
				learnCost = 50 // Default fallback
			}

			// Add to the best category
			skillInfo := gin.H{
				"name":      skill.Name,
				"learnCost": learnCost,
			}
			skillsByCategory[bestCategory] = append(skillsByCategory[bestCategory], skillInfo)

			// If the best category is different from original, also add to original with higher cost
			if bestCategory != originalCategory {
				skillInfoOriginal := gin.H{
					"name":      skill.Name,
					"learnCost": learnCost * 2, // Higher cost for non-optimal category
				}
				skillsByCategory[originalCategory] = append(skillsByCategory[originalCategory], skillInfoOriginal)
			}
		} else {
			// Fallback: add to original category only
			skillInfo := gin.H{
				"name":      skill.Name,
				"learnCost": 50, // Default learning cost
			}
			skillsByCategory[originalCategory] = append(skillsByCategory[originalCategory], skillInfo)
		}

		// Try to add skill to other logical categories with higher costs
		// This allows more flexibility in character creation
		for _, category := range allCategories {
			if category == bestCategory || category == originalCategory {
				continue // Already added
			}

			// Only add to certain categories if it makes sense
			if shouldSkillBeInCategory(skill.Name, category) {
				higherCost := learnCost
				if higherCost == 0 {
					higherCost = 50
				}
				higherCost = higherCost * 3 // Much higher cost for cross-category learning

				skillInfo := gin.H{
					"name":      skill.Name,
					"learnCost": higherCost,
				}
				skillsByCategory[category] = append(skillsByCategory[category], skillInfo)
			}
		}
	}

	return skillsByCategory, nil
}

// shouldSkillBeInCategory determines if a skill should be available in a given category
func shouldSkillBeInCategory(skillName, category string) bool {
	// Define which skills can appear in which categories
	skillCategoryMap := map[string][]string{
		// Physical skills can appear in multiple categories
		"Athletik":  {"Körper", "Kampf", "Freiland"},
		"Klettern":  {"Körper", "Freiland", "Alltag"},
		"Schwimmen": {"Körper", "Freiland", "Alltag"},
		"Laufen":    {"Körper", "Kampf", "Freiland"},
		"Akrobatik": {"Körper", "Kampf"},

		// Combat skills
		"Dolch":   {"Kampf", "Halbwelt"},
		"Schwert": {"Kampf"},
		"Bogen":   {"Kampf", "Freiland"},

		// Social skills
		"Menschenkenntnis": {"Sozial", "Halbwelt"},
		"Verführen":        {"Sozial", "Halbwelt"},
		"Anführen":         {"Sozial", "Kampf"},

		// Knowledge skills
		"Schreiben":  {"Wissen", "Alltag"},
		"Sprache":    {"Wissen", "Sozial"},
		"Naturkunde": {"Wissen", "Freiland"},

		// Stealth and underworld
		"Schleichen": {"Halbwelt", "Freiland", "Kampf"},
		"Tarnen":     {"Halbwelt", "Freiland", "Kampf"},
		"Stehlen":    {"Halbwelt"},

		// Survival and wilderness
		"Überleben":    {"Freiland", "Alltag"},
		"Spurensuche":  {"Freiland", "Halbwelt"},
		"Orientierung": {"Freiland", "Alltag"},
	}

	categories, exists := skillCategoryMap[skillName]
	if !exists {
		return false // Only add skills we explicitly define
	}

	for _, cat := range categories {
		if cat == category {
			return true
		}
	}
	return false
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

// Character Creation Session Management

// CreateCharacterSession erstellt eine neue Charakter-Erstellungssession
func CreateCharacterSession(c *gin.Context) {
	logger.Debug("CreateCharacterSession aufgerufen")

	// Debug: Alle Kontext-Keys anzeigen
	keys := make([]string, 0)
	for key := range c.Keys {
		keys = append(keys, fmt.Sprintf("%s=%v", key, c.Keys[key]))
	}
	logger.Debug("CreateCharacterSession: Verfügbare Kontext-Keys: [%s]", strings.Join(keys, ", "))

	userID := c.GetUint("userID")
	logger.Debug("CreateCharacterSession: UserID = %d", userID)

	if userID == 0 {
		logger.Warn("CreateCharacterSession: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	sessionID := fmt.Sprintf("char_create_%d_%d", userID, time.Now().Unix())
	logger.Debug("CreateCharacterSession: Generierte SessionID = %s", sessionID)

	session := models.CharacterCreationSession{
		ID:            sessionID,
		UserID:        userID,
		Name:          "",
		Geschlecht:    "",
		Rasse:         "",
		Typ:           "",
		Herkunft:      "",
		Stand:         "",
		Glaube:        "",
		Attributes:    models.AttributesData{},
		DerivedValues: models.DerivedValuesData{},
		Skills:        models.CharacterCreationSkills{},
		Spells:        models.CharacterCreationSpells{},
		SkillPoints:   models.SkillPointsData{},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ExpiresAt:     time.Now().AddDate(0, 0, 14), // 14 Tage
		CurrentStep:   1,
	}
	logger.Debug("CreateCharacterSession: Session-Struktur erstellt, ExpiresAt = %s", session.ExpiresAt.Format(time.RFC3339))

	// Session in Datenbank speichern
	logger.Debug("CreateCharacterSession: Speichere Session in Datenbank...")
	err := database.DB.Create(&session).Error
	if err != nil {
		logger.Error("CreateCharacterSession: Fehler beim Erstellen der Session: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	logger.Info("CreateCharacterSession: Session erfolgreich erstellt - SessionID: %s, UserID: %d", sessionID, userID)
	c.JSON(http.StatusCreated, gin.H{
		"session_id": sessionID,
		"expires_at": session.ExpiresAt,
	})
}

// ListCharacterSessions gibt alle aktiven Sessions für einen Benutzer zurück
func ListCharacterSessions(c *gin.Context) {
	logger.Debug("ListCharacterSessions aufgerufen")

	userID := c.GetUint("userID")
	logger.Debug("ListCharacterSessions: UserID = %d", userID)

	if userID == 0 {
		logger.Warn("ListCharacterSessions: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Sessions aus Datenbank laden
	logger.Debug("ListCharacterSessions: Lade Sessions für UserID %d aus Datenbank...", userID)
	sessions, err := models.GetUserSessions(database.DB, userID)
	if err != nil {
		logger.Error("ListCharacterSessions: Fehler beim Laden der Sessions: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load sessions"})
		return
	}

	logger.Debug("ListCharacterSessions: Gefundene Sessions: %d", len(sessions))

	// Sessions für Frontend formatieren
	var formattedSessions []gin.H
	for i, session := range sessions {
		// Schritt-Text bestimmen
		progressText := getProgressText(session.CurrentStep)
		logger.Debug("ListCharacterSessions: Formatiere Session %d - ID: %s, Step: %d, Name: %s",
			i+1, session.ID, session.CurrentStep, session.Name)

		formattedSessions = append(formattedSessions, gin.H{
			"session_id":    session.ID,
			"name":          session.Name,
			"rasse":         session.Rasse,
			"typ":           session.Typ,
			"current_step":  session.CurrentStep,
			"total_steps":   5,
			"created_at":    session.CreatedAt,
			"updated_at":    session.UpdatedAt,
			"expires_at":    session.ExpiresAt,
			"progress_text": progressText,
		})
	}

	logger.Info("ListCharacterSessions: Sessions erfolgreich geladen für UserID %d - Anzahl: %d", userID, len(formattedSessions))
	c.JSON(http.StatusOK, gin.H{
		"sessions": formattedSessions,
		"count":    len(formattedSessions),
	})
}

// getProgressText gibt den Schritt-Text für die Frontend-Anzeige zurück
func getProgressText(step int) string {
	logger.Debug("getProgressText: Ermittle Text für Schritt %d", step)

	var text string
	switch step {
	case 1:
		text = "Grundinformationen"
	case 2:
		text = "Attribute"
	case 3:
		text = "Abgeleitete Werte"
	case 4:
		text = "Fertigkeiten"
	case 5:
		text = "Zauber"
	default:
		text = "Unbekannt"
		logger.Warn("getProgressText: Unbekannter Schritt %d", step)
	}

	logger.Debug("getProgressText: Schritt %d = '%s'", step, text)
	return text
}

// GetCharacterSession gibt Session-Daten zurück
func GetCharacterSession(c *gin.Context) {
	logger.Debug("GetCharacterSession aufgerufen")

	sessionID := c.Param("sessionId")
	userID := c.GetUint("userID")
	logger.Debug("GetCharacterSession: SessionID = %s, UserID = %d", sessionID, userID)

	if userID == 0 {
		logger.Warn("GetCharacterSession: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Session aus Datenbank laden
	logger.Debug("GetCharacterSession: Lade Session aus Datenbank...")
	var session models.CharacterCreationSession
	err := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		logger.Error("GetCharacterSession: Session nicht gefunden - SessionID: %s, UserID: %d, Error: %s",
			sessionID, userID, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	logger.Debug("GetCharacterSession: Session gefunden - Name: %s, Step: %d, ExpiresAt: %s",
		session.Name, session.CurrentStep, session.ExpiresAt.Format(time.RFC3339))

	// Prüfen ob Session noch gültig ist
	if session.ExpiresAt.Before(time.Now()) {
		logger.Warn("GetCharacterSession: Session abgelaufen - SessionID: %s, ExpiresAt: %s",
			sessionID, session.ExpiresAt.Format(time.RFC3339))

		// Abgelaufene Session löschen
		logger.Debug("GetCharacterSession: Lösche abgelaufene Session...")
		database.DB.Delete(&session)

		c.JSON(http.StatusGone, gin.H{"error": "Session expired"})
		return
	}

	logger.Info("GetCharacterSession: Session erfolgreich geladen - SessionID: %s, UserID: %d, Step: %d",
		sessionID, userID, session.CurrentStep)
	c.JSON(http.StatusOK, session)
}

// UpdateCharacterBasicInfo Request
type UpdateBasicInfoRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=50"`
	Geschlecht string `json:"geschlecht" binding:"required"`
	Rasse      string `json:"rasse" binding:"required"`
	Typ        string `json:"typ" binding:"required"`
	Herkunft   string `json:"herkunft" binding:"required"`
	Stand      string `json:"stand" binding:"required"`
	Glaube     string `json:"glaube"`
}

// UpdateCharacterBasicInfo speichert Grundinformationen
func UpdateCharacterBasicInfo(c *gin.Context) {
	logger.Debug("UpdateCharacterBasicInfo aufgerufen")

	sessionID := c.Param("sessionId")
	userID := c.GetUint("userID")
	logger.Debug("UpdateCharacterBasicInfo: SessionID = %s, UserID = %d", sessionID, userID)

	if userID == 0 {
		logger.Warn("UpdateCharacterBasicInfo: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request UpdateBasicInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("UpdateCharacterBasicInfo: Ungültige Eingabedaten - %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Ungültige Eingabedaten: "+err.Error())
		return
	}

	logger.Debug("UpdateCharacterBasicInfo: Request-Daten - Name: %s, Geschlecht: %s, Rasse: %s, Typ: %s, Herkunft: %s, Stand: %s, Glaube: %s",
		request.Name, request.Geschlecht, request.Rasse, request.Typ, request.Herkunft, request.Stand, request.Glaube)

	// Session aus Datenbank laden
	logger.Debug("UpdateCharacterBasicInfo: Lade Session aus Datenbank...")
	var session models.CharacterCreationSession
	err := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		logger.Error("UpdateCharacterBasicInfo: Session nicht gefunden - SessionID: %s, UserID: %d, Error: %s",
			sessionID, userID, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	logger.Debug("UpdateCharacterBasicInfo: Aktueller Session-Status - Step: %d, Name: %s",
		session.CurrentStep, session.Name)

	// Grundinformationen aktualisieren
	session.Name = request.Name
	session.Geschlecht = request.Geschlecht
	session.Rasse = request.Rasse
	session.Typ = request.Typ
	session.Herkunft = request.Herkunft
	session.Stand = request.Stand
	session.Glaube = request.Glaube
	session.CurrentStep = 2
	session.UpdatedAt = time.Now()

	logger.Debug("UpdateCharacterBasicInfo: Session aktualisiert, setze CurrentStep auf 2")

	// Session in Datenbank aktualisieren
	logger.Debug("UpdateCharacterBasicInfo: Speichere Session in Datenbank...")
	err = database.DB.Save(&session).Error
	if err != nil {
		logger.Error("UpdateCharacterBasicInfo: Fehler beim Speichern der Session: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	logger.Info("UpdateCharacterBasicInfo: Grundinformationen erfolgreich gespeichert - SessionID: %s, Name: %s",
		sessionID, request.Name)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Grundinformationen gespeichert",
		"session_id":   sessionID,
		"current_step": 2,
	})
}

// UpdateAttributesRequest
type UpdateAttributesRequest struct {
	ST int `json:"st" binding:"required,min=1,max=100"` // Stärke
	GS int `json:"gs" binding:"required,min=1,max=100"` // Geschicklichkeit
	GW int `json:"gw" binding:"required,min=1,max=100"` // Gewandtheit
	KO int `json:"ko" binding:"required,min=1,max=100"` // Konstitution
	IN int `json:"in" binding:"required,min=1,max=100"` // Intelligenz
	ZT int `json:"zt" binding:"required,min=1,max=100"` // Zaubertalent
	AU int `json:"au" binding:"required,min=1,max=100"` // Ausstrahlung
}

// UpdateCharacterAttributes speichert Grundwerte
func UpdateCharacterAttributes(c *gin.Context) {
	logger.Debug("UpdateCharacterAttributes aufgerufen")

	sessionID := c.Param("sessionId")
	userID := c.GetUint("userID")
	logger.Debug("UpdateCharacterAttributes: SessionID = %s, UserID = %d", sessionID, userID)

	if userID == 0 {
		logger.Warn("UpdateCharacterAttributes: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request UpdateAttributesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("UpdateCharacterAttributes: Ungültige Attributswerte - %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Ungültige Attributswerte: "+err.Error())
		return
	}

	logger.Debug("UpdateCharacterAttributes: Attribute - ST:%d GS:%d GW:%d KO:%d IN:%d ZT:%d AU:%d",
		request.ST, request.GS, request.GW, request.KO, request.IN, request.ZT, request.AU)

	// Session aus Datenbank laden
	logger.Debug("UpdateCharacterAttributes: Lade Session aus Datenbank...")
	var session models.CharacterCreationSession
	err := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		logger.Error("UpdateCharacterAttributes: Session nicht gefunden - SessionID: %s, UserID: %d, Error: %s",
			sessionID, userID, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	logger.Debug("UpdateCharacterAttributes: Session geladen - CurrentStep: %d", session.CurrentStep)

	// Attribute aktualisieren
	session.Attributes = models.AttributesData{
		ST: request.ST,
		GS: request.GS,
		GW: request.GW,
		KO: request.KO,
		IN: request.IN,
		ZT: request.ZT,
		AU: request.AU,
	}
	session.CurrentStep = 3
	session.UpdatedAt = time.Now()

	logger.Debug("UpdateCharacterAttributes: Attribute gesetzt, CurrentStep auf 3 aktualisiert")

	// Session in Datenbank aktualisieren
	logger.Debug("UpdateCharacterAttributes: Speichere Session in Datenbank...")
	err = database.DB.Save(&session).Error
	if err != nil {
		logger.Error("UpdateCharacterAttributes: Fehler beim Speichern der Session: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	logger.Info("UpdateCharacterAttributes: Grundwerte erfolgreich gespeichert - SessionID: %s", sessionID)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Grundwerte gespeichert",
		"session_id":   sessionID,
		"current_step": 3,
	})
}

// UpdateDerivedValuesRequest
type UpdateDerivedValuesRequest struct {
	PA                    int `json:"pa" binding:"required,min=1,max=100"`               // Persönliche Ausstrahlung
	WK                    int `json:"wk" binding:"required,min=1,max=100"`               // Willenskraft
	LP_Max                int `json:"lp_max" binding:"required,min=1,max=50"`            // Lebenspunkte Maximum
	AP_Max                int `json:"ap_max" binding:"required,min=1,max=200"`           // Abenteuerpunkte Maximum
	B_Max                 int `json:"b_max" binding:"required,min=1,max=50"`             // Belastung Maximum
	ResistenzKoerper      int `json:"resistenz_koerper" binding:"required,min=1,max=20"` // Resistenz Körper
	ResistenzGeist        int `json:"resistenz_geist" binding:"required,min=1,max=20"`   // Resistenz Geist
	ResistenzBonusKoerper int `json:"resistenz_bonus_koerper" binding:"min=-5,max=5"`    // Resistenz Bonus Körper
	ResistenzBonusGeist   int `json:"resistenz_bonus_geist" binding:"min=-5,max=5"`      // Resistenz Bonus Geist
	Abwehr                int `json:"abwehr" binding:"required,min=1,max=20"`            // Abwehr
	AbwehrBonus           int `json:"abwehr_bonus" binding:"min=-5,max=5"`               // Abwehr Bonus
	AusdauerBonus         int `json:"ausdauer_bonus" binding:"min=-50,max=50"`           // Ausdauer Bonus
	AngriffsBonus         int `json:"angriffs_bonus" binding:"min=-5,max=5"`             // Angriffs Bonus
	Zaubern               int `json:"zaubern" binding:"required,min=1,max=20"`           // Zaubern
	ZauberBonus           int `json:"zauber_bonus" binding:"min=-5,max=5"`               // Zauber Bonus
	Raufen                int `json:"raufen" binding:"required,min=1,max=20"`            // Raufen
	SchadensBonus         int `json:"schadens_bonus" binding:"min=-10,max=10"`           // Schadens Bonus
	SG                    int `json:"sg" binding:"min=0,max=50"`                         // Schicksalsgunst
	GG                    int `json:"gg" binding:"min=0,max=50"`                         // Göttliche Gnade
	GP                    int `json:"gp" binding:"min=0,max=50"`                         // Glückspunkte
}

// UpdateCharacterDerivedValues speichert abgeleitete Werte
func UpdateCharacterDerivedValues(c *gin.Context) {
	logger.Debug("UpdateCharacterDerivedValues aufgerufen")

	sessionID := c.Param("sessionId")
	userID := c.GetUint("userID")
	logger.Debug("UpdateCharacterDerivedValues: SessionID = %s, UserID = %d", sessionID, userID)

	if userID == 0 {
		logger.Warn("UpdateCharacterDerivedValues: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request UpdateDerivedValuesRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("UpdateCharacterDerivedValues: Ungültige abgeleitete Werte - %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Ungültige abgeleitete Werte: "+err.Error())
		return
	}

	logger.Debug("UpdateCharacterDerivedValues: Werte - LP_Max:%d AP_Max:%d B_Max:%d PA:%d WK:%d SG:%d GG:%d GP:%d",
		request.LP_Max, request.AP_Max, request.B_Max, request.PA, request.WK, request.SG, request.GG, request.GP)

	// Session aus Datenbank laden
	logger.Debug("UpdateCharacterDerivedValues: Lade Session aus Datenbank...")
	var session models.CharacterCreationSession
	err := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		logger.Error("UpdateCharacterDerivedValues: Session nicht gefunden - SessionID: %s, UserID: %d, Error: %s",
			sessionID, userID, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	logger.Debug("UpdateCharacterDerivedValues: Session geladen - CurrentStep: %d", session.CurrentStep)

	// Abgeleitete Werte aktualisieren
	session.DerivedValues = models.DerivedValuesData{
		PA:                    request.PA,
		WK:                    request.WK,
		LPMax:                 request.LP_Max,
		APMax:                 request.AP_Max,
		BMax:                  request.B_Max,
		ResistenzKoerper:      request.ResistenzKoerper,
		ResistenzGeist:        request.ResistenzGeist,
		ResistenzBonusKoerper: request.ResistenzBonusKoerper,
		ResistenzBonusGeist:   request.ResistenzBonusGeist,
		Abwehr:                request.Abwehr,
		AbwehrBonus:           request.AbwehrBonus,
		AusdauerBonus:         request.AusdauerBonus,
		AngriffsBonus:         request.AngriffsBonus,
		Zaubern:               request.Zaubern,
		ZauberBonus:           request.ZauberBonus,
		Raufen:                request.Raufen,
		SchadensBonus:         request.SchadensBonus,
		SG:                    request.SG,
		GG:                    request.GG,
		GP:                    request.GP,
	}
	session.CurrentStep = 4
	session.UpdatedAt = time.Now()

	logger.Debug("UpdateCharacterDerivedValues: Abgeleitete Werte gesetzt, CurrentStep auf 4 aktualisiert")

	// Session in Datenbank aktualisieren
	logger.Debug("UpdateCharacterDerivedValues: Speichere Session in Datenbank...")
	err = database.DB.Save(&session).Error
	if err != nil {
		logger.Error("UpdateCharacterDerivedValues: Fehler beim Speichern der Session: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	logger.Info("UpdateCharacterDerivedValues: Abgeleitete Werte erfolgreich gespeichert - SessionID: %s", sessionID)
	c.JSON(http.StatusOK, gin.H{
		"message":      "Abgeleitete Werte gespeichert",
		"session_id":   sessionID,
		"current_step": 4,
	})
}

// UpdateSkillsRequest
type UpdateSkillsRequest struct {
	Skills      models.CharacterCreationSkills `json:"skills"`
	Spells      models.CharacterCreationSpells `json:"spells"`
	SkillPoints models.SkillPointsData         `json:"skill_points"` // Verbleibende Punkte pro Kategorie
}

// UpdateCharacterSkills speichert Fertigkeiten und Zauber
func UpdateCharacterSkills(c *gin.Context) {
	logger.Debug("UpdateCharacterSkills aufgerufen")

	sessionID := c.Param("sessionId")
	userID := c.GetUint("userID")
	logger.Debug("UpdateCharacterSkills: SessionID = %s, UserID = %d", sessionID, userID)

	if userID == 0 {
		logger.Warn("UpdateCharacterSkills: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request UpdateSkillsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("UpdateCharacterSkills: Ungültige Fertigkeitsdaten - %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Ungültige Fertigkeitsdaten: "+err.Error())
		return
	}

	logger.Debug("UpdateCharacterSkills: Skills-Anzahl: %d, Spells-Anzahl: %d",
		len(request.Skills), len(request.Spells))

	// Session aus Datenbank laden
	logger.Debug("UpdateCharacterSkills: Lade Session aus Datenbank...")
	var session models.CharacterCreationSession
	err := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		logger.Error("UpdateCharacterSkills: Session nicht gefunden - SessionID: %s, UserID: %d, Error: %s",
			sessionID, userID, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	logger.Debug("UpdateCharacterSkills: Session geladen - CurrentStep: %d", session.CurrentStep)

	// Fertigkeiten und Zauber aktualisieren
	session.Skills = request.Skills
	session.Spells = request.Spells
	session.SkillPoints = request.SkillPoints
	session.CurrentStep = 5
	session.UpdatedAt = time.Now()

	logger.Debug("UpdateCharacterSkills: Skills/Spells gesetzt, CurrentStep auf 5 aktualisiert")

	// Session in Datenbank aktualisieren
	logger.Debug("UpdateCharacterSkills: Speichere Session in Datenbank...")
	err = database.DB.Save(&session).Error
	if err != nil {
		logger.Error("UpdateCharacterSkills: Fehler beim Speichern der Session: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
		return
	}

	logger.Info("UpdateCharacterSkills: Fertigkeiten erfolgreich gespeichert - SessionID: %s, Skills: %d, Spells: %d",
		sessionID, len(request.Skills), len(request.Spells))
	c.JSON(http.StatusOK, gin.H{
		"message":      "Fertigkeiten gespeichert",
		"session_id":   sessionID,
		"current_step": 5,
	})
}

// FinalizeCharacterCreation schließt die Charakter-Erstellung ab
func FinalizeCharacterCreation(c *gin.Context) {
	logger.Debug("FinalizeCharacterCreation aufgerufen")

	sessionID := c.Param("sessionId")
	userID := c.GetUint("userID")
	logger.Debug("FinalizeCharacterCreation: SessionID = %s, UserID = %d", sessionID, userID)

	if userID == 0 {
		logger.Warn("FinalizeCharacterCreation: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Session laden
	logger.Debug("FinalizeCharacterCreation: Lade Session aus Datenbank...")
	var session models.CharacterCreationSession
	err := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	if err != nil {
		logger.Error("FinalizeCharacterCreation: Session nicht gefunden - SessionID: %s, UserID: %d, Error: %s",
			sessionID, userID, err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	logger.Debug("FinalizeCharacterCreation: Session geladen - Name: %s, CurrentStep: %d",
		session.Name, session.CurrentStep)

	// Session validieren
	if session.CurrentStep < 5 {
		logger.Warn("FinalizeCharacterCreation: Charakter-Erstellung unvollständig - CurrentStep: %d (erwartet: 5)",
			session.CurrentStep)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Character creation not complete"})
		return
	}

	logger.Debug("FinalizeCharacterCreation: Erstelle Charakter-Struktur...")
	// Character erstellen
	char := models.Char{
		BamortBase: models.BamortBase{
			Name: session.Name,
		},
		UserID: userID,
		Rasse:  session.Rasse,
		Typ:    session.Typ,
		Glaube: session.Glaube,
		Public: false, // Default to private
		Grad:   1,     // Default starting grade

		// Lebenspunkte
		Lp: models.Lp{
			Max:   session.DerivedValues.LPMax,
			Value: session.DerivedValues.LPMax,
		},

		// Ausdauerpunkte
		Ap: models.Ap{
			Max:   session.DerivedValues.APMax,
			Value: session.DerivedValues.APMax,
		},

		// Bewegung
		B: models.B{
			Max:   session.DerivedValues.BMax,
			Value: session.DerivedValues.BMax,
		},
		Vermoegen: models.Vermoegen{
			BamortCharTrait: models.BamortCharTrait{
				UserID: userID,
			},
			Goldstücke: 80,
		},

		// Bennies (Glückspunkte, etc.)
		Bennies: models.Bennies{
			BamortCharTrait: models.BamortCharTrait{
				UserID: userID,
			},
			Gg: session.DerivedValues.GG,
			Gp: session.DerivedValues.GP,
			Sg: session.DerivedValues.SG,
		},
	}

	// Eigenschaften (Attribute) hinzufügen
	char.Eigenschaften = []models.Eigenschaft{
		{UserID: userID, Name: "St", Value: session.Attributes.ST},
		{UserID: userID, Name: "Gs", Value: session.Attributes.GS},
		{UserID: userID, Name: "Gw", Value: session.Attributes.GW},
		{UserID: userID, Name: "Ko", Value: session.Attributes.KO},
		{UserID: userID, Name: "In", Value: session.Attributes.IN},
		{UserID: userID, Name: "Zt", Value: session.Attributes.ZT},
		{UserID: userID, Name: "Au", Value: session.Attributes.AU},
		{UserID: userID, Name: "pA", Value: session.DerivedValues.PA}, // PA kommt aus derived values
		{UserID: userID, Name: "Wk", Value: session.DerivedValues.WK}, // WK kommt aus derived values
	}

	logger.Debug("FinalizeCharacterCreation: Charakter-Struktur erstellt mit %d Eigenschaften",
		len(char.Eigenschaften))

	// Fertigkeiten aus der Session übertragen
	logger.Debug("FinalizeCharacterCreation: Übertrage %d Fertigkeiten", len(session.Skills))
	for _, skill := range session.Skills {
		// Suche den Initialwert der Fertigkeit aus der Datenbank
		dbSkill := models.Skill{}
		err := dbSkill.First(skill.Name)
		if err != nil {
			logger.Warn("FinalizeCharacterCreation: Konnte Fertigkeit '%s' nicht in der Datenbank finden, verwende Level %d", skill.Name, skill.Level)
		}

		// Verwende den Initialwert aus der Datenbank wenn verfügbar, sonst fallback auf Session-Level
		initialValue := skill.Level // Fallback
		if err == nil {
			initialValue = dbSkill.Initialwert
			logger.Debug("FinalizeCharacterCreation: Verwende Initialwert %d für Fertigkeit '%s'", initialValue, skill.Name)
		}

		// Unterscheide zwischen normalen Fertigkeiten und Waffenfertigkeiten
		if skill.Category == "Waffen" || skill.Category == "waffen" {
			dbWPSkill := models.WeaponSkill{}
			err := dbWPSkill.First(skill.Name)
			if err != nil {
				logger.Warn("FinalizeCharacterCreation: Konnte WaffenFertigkeit '%s' nicht in der Datenbank finden, verwende Level %d", skill.Name, skill.Level)
			}
			// Verwende den Initialwert aus der Datenbank wenn verfügbar, sonst fallback auf Session-Level
			initialValue := skill.Level // Fallback
			if err == nil {
				initialValue = dbWPSkill.Initialwert
				logger.Debug("FinalizeCharacterCreation: Verwende Initialwert %d für Fertigkeit '%s'", initialValue, skill.Name)
			}
			// Waffenfertigkeit
			weaponSkill := models.SkWaffenfertigkeit{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{
							Name: skill.Name,
						},
						UserID: userID,
					},
					Fertigkeitswert: initialValue,
					Improvable:      true,
					Category:        skill.Category,
				},
			}
			char.Waffenfertigkeiten = append(char.Waffenfertigkeiten, weaponSkill)
		} else {
			// Normale Fertigkeit
			normalSkill := models.SkFertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: skill.Name,
					},
					CharacterID: char.ID,
					UserID:      userID,
				},
				Fertigkeitswert: initialValue,
				Improvable:      true,
				Category:        skill.Category,
			}
			char.Fertigkeiten = append(char.Fertigkeiten, normalSkill)
		}
	}

	// Zauber aus der Session übertragen
	logger.Debug("FinalizeCharacterCreation: Übertrage %d Zauber", len(session.Spells))
	for _, spell := range session.Spells {
		zauber := models.SkZauber{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: spell.Name,
				},
				CharacterID: char.ID,
				UserID:      userID,
			},
		}
		char.Zauber = append(char.Zauber, zauber)
	}

	logger.Debug("FinalizeCharacterCreation: Charakter vollständig erstellt - %d Fertigkeiten, %d Waffenfertigkeiten, %d Zauber",
		len(char.Fertigkeiten), len(char.Waffenfertigkeiten), len(char.Zauber))

	// Character in Datenbank speichern
	logger.Debug("FinalizeCharacterCreation: Speichere Charakter in Datenbank...")
	err = char.Create()
	if err != nil {
		logger.Error("FinalizeCharacterCreation: Fehler beim Erstellen des Charakters: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character: " + err.Error()})
		return
	}

	logger.Debug("FinalizeCharacterCreation: Charakter erfolgreich erstellt mit ID: %d", char.ID)

	// Session löschen
	logger.Debug("FinalizeCharacterCreation: Lösche Session aus Datenbank...")
	database.DB.Delete(&session)

	logger.Info("FinalizeCharacterCreation: Charakter-Erstellung abgeschlossen - CharacterID: %d, SessionID: %s, Name: %s",
		char.ID, sessionID, session.Name)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Charakter erfolgreich erstellt",
		"character_id": char.ID,
		"session_id":   sessionID,
	})
}

// DeleteCharacterSession löscht eine Session
func DeleteCharacterSession(c *gin.Context) {
	logger.Debug("DeleteCharacterSession aufgerufen")

	sessionID := c.Param("sessionId")
	userID := c.GetUint("userID")
	logger.Debug("DeleteCharacterSession: SessionID = %s, UserID = %d", sessionID, userID)

	if userID == 0 {
		logger.Warn("DeleteCharacterSession: Unauthorized - UserID ist 0")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Session aus Datenbank löschen (nur eigene Sessions)
	logger.Debug("DeleteCharacterSession: Lösche Session aus Datenbank...")
	result := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).Delete(&models.CharacterCreationSession{})
	if result.Error != nil {
		logger.Error("DeleteCharacterSession: Fehler beim Löschen der Session: %s", result.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	if result.RowsAffected == 0 {
		logger.Warn("DeleteCharacterSession: Session nicht gefunden oder bereits gelöscht - SessionID: %s, UserID: %d",
			sessionID, userID)
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	logger.Info("DeleteCharacterSession: Session erfolgreich gelöscht - SessionID: %s, UserID: %d, RowsAffected: %d",
		sessionID, userID, result.RowsAffected)

	c.JSON(http.StatusOK, gin.H{
		"message":    "Session gelöscht",
		"session_id": sessionID,
	})
}

// Reference Data Handlers

// GetRaces gibt verfügbare Rassen zurück
func GetRaces(c *gin.Context) {
	// TODO: Aus Datenbank laden
	races := []string{
		"Mensch", "Elf", "Halbling", "Zwerg", "Gnom",
	}

	c.JSON(http.StatusOK, gin.H{"races": races})
}

// GetCharacterClasses gibt verfügbare Klassen zurück
func GetCharacterClasses(c *gin.Context) {
	// Get game system from query parameter, default to "midgard"
	gameSystem := c.DefaultQuery("game_system", "midgard")

	// Load character classes from database
	classes, err := models.GetCharacterClassesByActiveSources(gameSystem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load character classes"})
		return
	}

	// Extract class names for the frontend
	var classNames []string
	for _, class := range classes {
		classNames = append(classNames, class.Name)
	}

	// If no classes found in database, fall back to hardcoded list
	if len(classNames) == 0 {
		classNames = []string{
			"Assassine", "Barbar", "Barde",
			"Heiler", "Händler", "Kämpfer", "Krieger",
			"Magier", "Ordenskrieger", "Priester Beschützer", "Schamane", "Skalde",
			"Magister", "Thaumaturg", "Waldläufer", "Zauberer",
		}
	}

	c.JSON(http.StatusOK, gin.H{"classes": classNames})
}

// GetOrigins gibt verfügbare Herkünfte zurück
func GetOrigins(c *gin.Context) {
	// TODO: Aus Datenbank laden
	origins := []string{
		"Alba", "Aran", "Buluga", "Chryseia",
		"Eschar", "Fuardain", "Ikenga", "KanThaiPan", "Küstenstaaten",
		"Medjis", "Moravod", "Nahuatlan", "Rawindra", "Scharidis",
		"Tegarisch Steppe", "Valian", "Waeland", "Ywerddon",
	}

	c.JSON(http.StatusOK, gin.H{"origins": origins})
}

// SearchBeliefs sucht Glaubensrichtungen
func SearchBeliefs(c *gin.Context) {
	query := c.Query("q")

	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mindestens 2 Zeichen erforderlich"})
		return
	}

	// Get game system from query parameter, default to "midgard"
	gameSystem := c.DefaultQuery("game_system", "midgard")

	// Load beliefs from database
	believes, err := models.GetBelievesByActiveSources(gameSystem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load beliefs from database"})
		return
	}

	// Extract belief names and filter by query
	var allBeliefs []string
	for _, belief := range believes {
		allBeliefs = append(allBeliefs, belief.Name)
	}

	// If no beliefs found in database, fall back to hardcoded list
	if len(allBeliefs) == 0 {
		allBeliefs = []string{
			"Apshai", "Arthusos", "Beschützer", "Dwyllas", "Elfen",
			"Fruchtbarkeitsgöttin", "Gaia", "Grafschafter", "Heiler",
			"Jäger", "Kämpfer", "Lichbringer", "Meeresherr", "Natur",
			"Ostküste", "Priester", "Rechtschaffener", "Schutzpatron",
			"Stammesgeist", "Totengott", "Unterwelt", "Vater", "Weisheit",
			"Xan", "Ylhoon", "Zauberer",
		}
	}

	var results []string
	queryLower := strings.ToLower(query)
	for _, belief := range allBeliefs {
		if strings.Contains(strings.ToLower(belief), queryLower) {
			results = append(results, belief)
		}
	}

	c.JSON(http.StatusOK, gin.H{"beliefs": results})
}

// SkillCategoryWithPoints repräsentiert eine Kategorie mit verfügbaren Lernpunkten
type SkillCategoryWithPoints struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Points      int    `json:"points"`
	MaxPoints   int    `json:"max_points"`
}

// LearningPointsData repräsentiert die Lernpunkte und typischen Fertigkeiten einer Charakterklasse
type LearningPointsData struct {
	ClassName      string         `json:"class_name"`
	ClassCode      string         `json:"class_code"`
	LearningPoints map[string]int `json:"learning_points"` // Kategorie -> Lernpunkte
	WeaponPoints   int            `json:"weapon_points"`   // Waffenlernpunkte
	SpellPoints    int            `json:"spell_points"`    // Zauberlerneinheiten (falls vorhanden)
	TypicalSkills  []TypicalSkill `json:"typical_skills"`  // Typische Fertigkeiten
	TypicalSpells  []string       `json:"typical_spells"`  // Typische Zauber (falls vorhanden)
}

// TypicalSkill repräsentiert eine typische Fertigkeit mit Bonus
type TypicalSkill struct {
	Name      string `json:"name"`
	Bonus     int    `json:"bonus"`
	Attribute string `json:"attribute"` // Zugehöriges Attribut (z.B. "Gs", "In")
	Notes     string `json:"notes"`     // Zusätzliche Notizen
}

// GetCharacterClassLearningPoints gibt die Lernpunkte und typischen Fertigkeiten für eine Charakterklasse zurück
func GetCharacterClassLearningPoints(c *gin.Context) {
	className := c.Query("class")
	if className == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Charakterklassen-Name ist erforderlich (Parameter 'class')"})
		return
	}

	stand := c.Query("stand") // Optional: Unfreie, Volk, Mittelschicht, Adel

	// Hole die Lernpunkte-Daten für die Klasse
	learningData, err := getLearningPointsForClass(className, stand)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Charakterklasse nicht gefunden oder nicht unterstützt: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, learningData)
}

// getLearningPointsForClass gibt die Lernpunkte-Daten für eine bestimmte Charakterklasse zurück
func getLearningPointsForClass(className string, stand string) (*LearningPointsData, error) {
	// Mapping der Klassennamen zu Codes (falls notwendig)
	classMapping := map[string]string{
		"Assassine":           "As",
		"Barbar":              "Bb",
		"Glücksritter":        "Gl",
		"Händler":             "Hä",
		"Krieger":             "Kr",
		"Spitzbube":           "Sp",
		"Waldläufer":          "Wa",
		"Barde":               "Ba",
		"Ordenskrieger":       "Or",
		"Druide":              "Dr",
		"Hexer":               "Hx",
		"Magier":              "Ma",
		"Priester Beschützer": "PB",
		"Priester Streiter":   "PS",
		"Schamane":            "Sc",
	}

	classCode := classMapping[className]
	if classCode == "" {
		classCode = className // Falls der Name bereits ein Code ist
	}

	// Definiere die Lernpunkte-Daten basierend auf Lerntabelle_Erstellung.md
	var data *LearningPointsData

	switch classCode {
	case "As", "Assassine":
		data = &LearningPointsData{
			ClassName: "Assassine",
			ClassCode: "As",
			LearningPoints: map[string]int{
				"Alltag":    1,
				"Halbwelt":  2,
				"Sozial":    4,
				"Unterwelt": 8,
				"Waffen":    24,
			},
			//WeaponPoints: 24,
			TypicalSkills: []TypicalSkill{
				{Name: "Meucheln", Bonus: 8, Attribute: "Gs", Notes: ""},
			},
		}
	case "Bb", "Barbar":
		data = &LearningPointsData{
			ClassName: "Barbar",
			ClassCode: "Bb",
			LearningPoints: map[string]int{
				"Alltag":   2,
				"Freiland": 4,
				"Kampf":    1,
				"Körper":   2,
				"Waffen":   24,
			},
			//WeaponPoints: 24,
			TypicalSkills: []TypicalSkill{
				{Name: "Spurensuche", Bonus: 8, Attribute: "In", Notes: "in Heimatlandschaft"},
				{Name: "Überleben", Bonus: 8, Attribute: "In", Notes: "in Heimatlandschaft"},
			},
		}
	case "Gl", "Glücksritter":
		data = &LearningPointsData{
			ClassName: "Glücksritter",
			ClassCode: "Gl",
			LearningPoints: map[string]int{
				"Alltag":   2,
				"Halbwelt": 3,
				"Sozial":   8,
				"Waffen":   24,
			},
			//WeaponPoints: 24,
			TypicalSkills: []TypicalSkill{
				{Name: "Fechten", Bonus: 5, Attribute: "Gs", Notes: "oder beidhändiger Kampf+5 (Gs)"},
			},
		}
	case "Hä", "Händler":
		data = &LearningPointsData{
			ClassName: "Händler",
			ClassCode: "Hä",
			LearningPoints: map[string]int{
				"Alltag": 4,
				"Sozial": 8,
				"Wissen": 4,
				"Waffen": 20,
			},
			//WeaponPoints: 20,
			TypicalSkills: []TypicalSkill{
				{Name: "Geschäftssinn", Bonus: 8, Attribute: "In", Notes: ""},
			},
		}
	case "Kr", "Krieger":
		data = &LearningPointsData{
			ClassName: "Krieger",
			ClassCode: "Kr",
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Kampf":  3,
				"Körper": 1,
				"Waffen": 36,
			},
			//WeaponPoints: 36,
			TypicalSkills: []TypicalSkill{
				{Name: "Kampf in Vollrüstung", Bonus: 5, Attribute: "St", Notes: ""},
			},
		}
	case "Sp", "Spitzbube":
		data = &LearningPointsData{
			ClassName: "Spitzbube",
			ClassCode: "Sp",
			LearningPoints: map[string]int{
				"Alltag":    2,
				"Halbwelt":  6,
				"Unterwelt": 12,
				"Waffen":    20,
			},
			//WeaponPoints: 20,
			TypicalSkills: []TypicalSkill{
				{Name: "Fallenmechanik", Bonus: 8, Attribute: "Gs", Notes: "oder Geschäftssinn+8 (In)"},
			},
		}
	case "Wa", "Waldläufer":
		data = &LearningPointsData{
			ClassName: "Waldläufer",
			ClassCode: "Wa",
			LearningPoints: map[string]int{
				"Alltag":   1,
				"Freiland": 11,
				"Körper":   4,
				"Waffen":   20,
			},
			//WeaponPoints: 20,
			TypicalSkills: []TypicalSkill{
				{Name: "Scharfschießen", Bonus: 5, Attribute: "Gs", Notes: ""},
			},
		}
	case "Ba", "Barde":
		data = &LearningPointsData{
			ClassName: "Barde",
			ClassCode: "Ba",
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Sozial": 4,
				"Wissen": 4,
				"Waffen": 16,
			},
			//WeaponPoints: 16,
			SpellPoints: 3,
			TypicalSkills: []TypicalSkill{
				{Name: "Musizieren", Bonus: 12, Attribute: "Gs", Notes: ""},
				{Name: "Landeskunde", Bonus: 8, Attribute: "In", Notes: "für Heimat"},
			},
			TypicalSpells: []string{"Zauberlieder"},
		}
	case "Or", "Ordenskrieger":
		data = &LearningPointsData{
			ClassName: "Ordenskrieger",
			ClassCode: "Or",
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Kampf":  3,
				"Wissen": 2,
				"Waffen": 18,
			},
			//WeaponPoints: 18,
			SpellPoints: 3,
			TypicalSkills: []TypicalSkill{
				{Name: "Athletik", Bonus: 8, Attribute: "St", Notes: "oder Meditieren+8 (Wk)"},
			},
			TypicalSpells: []string{"Wundertaten"},
		}
	case "Dr", "Druide":
		data = &LearningPointsData{
			ClassName: "Druide",
			ClassCode: "Dr",
			LearningPoints: map[string]int{
				"Alltag":   2,
				"Freiland": 4,
				"Wissen":   2,
				"Waffen":   6,
			},
			//WeaponPoints: 6,
			SpellPoints: 5,
			TypicalSkills: []TypicalSkill{
				{Name: "Pflanzenkunde", Bonus: 8, Attribute: "In", Notes: ""},
				{Name: "Schreiben", Bonus: 12, Attribute: "In", Notes: "für Ogam-Zeichen"},
			},
			TypicalSpells: []string{"Dweomer", "Tiere rufen"},
		}
	case "Hx", "Hexer":
		data = &LearningPointsData{
			ClassName: "Hexer",
			ClassCode: "Hx",
			LearningPoints: map[string]int{
				"Alltag": 3,
				"Sozial": 2,
				"Wissen": 2,
				"Waffen": 2,
			},
			//WeaponPoints: 2,
			SpellPoints: 6,
			TypicalSkills: []TypicalSkill{
				{Name: "Gassenwissen", Bonus: 8, Attribute: "In", Notes: "oder Verführen+8 (pA)"},
			},
			TypicalSpells: []string{"Beherrschen", "Verändern", "Verwünschen", "Binden des Vertrauten"},
		}
	case "Ma", "Magier":
		data = &LearningPointsData{
			ClassName: "Magier",
			ClassCode: "Ma",
			LearningPoints: map[string]int{
				"Alltag": 1,
				"Wissen": 5,
				"Waffen": 2,
			},
			//WeaponPoints: 2,
			SpellPoints: 7,
			TypicalSkills: []TypicalSkill{
				{Name: "Zauberkunde", Bonus: 8, Attribute: "In", Notes: ""},
				{Name: "Schreiben", Bonus: 12, Attribute: "In", Notes: "für Muttersprache"},
			},
			TypicalSpells: []string{"beliebig außer Dweomer, Wundertaten, Zauberlieder", "Erkennen von Zauberei"},
		}
	case "PB", "Priester Beschützer":
		data = &LearningPointsData{
			ClassName: "Priester Beschützer",
			ClassCode: "PB",
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Sozial": 2,
				"Wissen": 3,
				"Waffen": 6,
			},
			//WeaponPoints: 6,
			SpellPoints: 5,
			TypicalSkills: []TypicalSkill{
				{Name: "Menschenkenntnis", Bonus: 8, Attribute: "In", Notes: ""},
				{Name: "Schreiben", Bonus: 12, Attribute: "In", Notes: "für Muttersprache"},
			},
			TypicalSpells: []string{"Wundertaten", "Heilen von Wunden"},
		}
	case "PS", "Priester Streiter":
		data = &LearningPointsData{
			ClassName: "Priester Streiter",
			ClassCode: "PS",
			LearningPoints: map[string]int{
				"Alltag": 3,
				"Kampf":  2,
				"Wissen": 2,
				"Waffen": 8,
			},
			//WeaponPoints: 8,
			SpellPoints: 5,
			TypicalSkills: []TypicalSkill{
				{Name: "Erste Hilfe", Bonus: 8, Attribute: "Gs", Notes: ""},
				{Name: "Schreiben", Bonus: 12, Attribute: "In", Notes: "für Muttersprache"},
			},
			TypicalSpells: []string{"Wundertaten", "Bannen von Finsterwerk", "Strahlender Panzer"},
		}
	case "Sc", "Schamane":
		data = &LearningPointsData{
			ClassName: "Schamane",
			ClassCode: "Sc",
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Körper": 4,
				"Wissen": 2,
				"Waffen": 6,
			},
			//WeaponPoints: 6,
			SpellPoints: 5,
			TypicalSkills: []TypicalSkill{
				{Name: "Tierkunde", Bonus: 8, Attribute: "In", Notes: ""},
				{Name: "Überleben", Bonus: 8, Attribute: "In", Notes: "in Heimatlandschaft"},
			},
			TypicalSpells: []string{"Dweomer", "Wundertaten", "Austreibung des Bösen", "Bannen von Gift"},
		}
	default:
		return nil, fmt.Errorf("unbekannte Charakterklasse: %s", className)
	}

	// Bonus-Lernpunkte basierend auf Stand hinzufügen
	if stand != "" && data != nil {
		standBonus := getStandBonusPoints(stand)
		// Füge die Stand-Bonuspunkte zu den normalen Lernpunkten hinzu
		for category, bonus := range standBonus {
			if currentPoints, exists := data.LearningPoints[category]; exists {
				data.LearningPoints[category] = currentPoints + bonus
			} else {
				// Falls die Kategorie noch nicht existiert, füge sie hinzu
				data.LearningPoints[category] = bonus
			}
		}
		// Speichere die Stand-Bonuspunkte auch separat für Referenz
		//data.StandPoints = standBonus
	}

	return data, nil
}

// getStandBonusPoints gibt die Bonus-Lernpunkte basierend auf dem Stand zurück
func getStandBonusPoints(stand string) map[string]int {
	switch stand {
	case "Unfreie":
		return map[string]int{"Halbwelt": 2}
	case "Volk":
		return map[string]int{"Alltag": 2}
	case "Mittelschicht":
		return map[string]int{"Wissen": 2}
	case "Adel":
		return map[string]int{"Sozial": 2}
	default:
		return make(map[string]int)
	}
}
