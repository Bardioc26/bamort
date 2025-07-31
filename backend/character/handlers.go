package character

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"

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
		gsmsk := skill.GetGsm()
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

func GetLearnSkillCost(c *gin.Context) {
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

	cost, err := gsmaster.CalculateSkillLearnCost(skill.Name, character.Typ)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
		return
	}

	// Return the updated character
	c.JSON(http.StatusOK, cost)
}

func GetLearnSpellCost(c *gin.Context) {
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

	cost, err := gsmaster.CalculateSpellLearnCost(spell.Name, character.Typ)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "error getting costs to learn spell: "+err.Error())
		return
	}

	sd.CostEP = cost
	// Return the updated character
	c.JSON(http.StatusOK, sd)
}

/*
func GetSkillNextLevelCosts(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	for int, skill := range character.Fertigkeiten {
		lCost, err := gsmaster.CalculateSkillImprovementCost(skill.Name, character.Typ, skill.Fertigkeitswert)
		if err != nil {
			respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
			return
		}
		character.Fertigkeiten[int].LearningCost = *lCost
	}

	// Load the skill from the request
	var s models.Fertigkeit
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Return the updated character
	c.JSON(http.StatusOK, character.Fertigkeiten)
}
*/
/*
func GetSkillAllLevelCosts(c *gin.Context) {
	// Get the character ID from the request
	charID := c.Param("id")

	// Load the character from the database
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}
	var s LearnRequestStruct
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	if s.Name == "" {
		respondWithError(c, http.StatusBadRequest, "no name given")
	}

	costArr := make([]gsmaster.LearnCost, 0)
	notfound := true
	for _, skill := range character.Fertigkeiten {
		if skill.Name == s.Name {
			for i := skill.Fertigkeitswert; i <= 20; i++ {
				lCost, err := gsmaster.CalculateSkillImprovementCost(skill.Name, character.Typ, skill.Fertigkeitswert)
				if err != nil {
					respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
					return
				}
				costArr = append(costArr, *lCost)
			}
			notfound = false
			break
		}
	}
	if notfound {
		for _, skill := range character.Waffenfertigkeiten {
			if skill.Name == s.Name {
				for i := skill.Fertigkeitswert; i <= 20; i++ {
					lCost, err := gsmaster.CalculateSkillImprovementCost(skill.Name, character.Typ, skill.Fertigkeitswert)
					if err != nil {
						respondWithError(c, http.StatusBadRequest, "error getting costs to learn skill: "+err.Error())
						return
					}
					costArr = append(costArr, *lCost)
				}
				break
			}
		}
	}

	// Return the updated character
	c.JSON(http.StatusOK, costArr)
}
*/
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

// calculateMultiLevelCosts berechnet die Kosten für mehrere Level-Verbesserungen mit gsmaster.GetLernCostNextLevel
func calculateMultiLevelCosts(character *models.Char, skillName string, currentLevel int, levelsToLearn []int, rewardType string, usePP, useGold int) (*models.LearnCost, error) {
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
		classAbr := getCharacterClass(character)
		cat, difficulty, _ := gsmaster.FindBestCategoryForSkillLearning(skillName, classAbr)
		levelResult := gsmaster.SkillCostResultNew{
			CharacterID:    fmt.Sprintf("%d", character.ID),
			CharacterClass: classAbr,
			SkillName:      skillName,
			Category:       cat,
			Difficulty:     gsmaster.GetSkillDifficulty(difficulty, skillName),
			TargetLevel:    targetLevel,
		}

		// Temporäre Request für dieses Level
		tempRequest := request
		tempRequest.CurrentLevel = targetLevel - 1
		tempRequest.UsePP = remainingPP
		tempRequest.UseGold = remainingGold

		err := gsmaster.GetLernCostNextLevel(&tempRequest, &levelResult, rewardTypePtr, targetLevel, character.Typ)
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

// getCharacterClass gibt die Charakterklassen-Abkürzung zurück
func getCharacterClass(character *models.Char) string {
	if len(character.Typ) > 3 {
		return gsmaster.GetClassAbbreviation(character.Typ)
	}
	return character.Typ
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
	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(character.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviation(character.Typ)
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

		err = gsmaster.GetLernCostNextLevel(&tempRequest, &costResult, request.Reward, nextLevel, character.Rasse)
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

// ImproveSkill verbessert eine bestehende Fertigkeit und erstellt Audit-Log-Einträge
func ImproveSkill(c *gin.Context) {
	// Verwende gsmaster.LernCostRequest direkt
	var request gsmaster.LernCostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Hole Charakter über die ID aus dem Request
	var character models.Char
	err := database.DB.
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Erfahrungsschatz").
		Preload("Vermoegen").
		First(&character, request.CharId).Error
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Verwende Klassenabkürzung wenn der Typ länger als 3 Zeichen ist
	var characterClass string
	if len(character.Typ) > 3 {
		characterClass = gsmaster.GetClassAbbreviation(character.Typ)
	} else {
		characterClass = character.Typ
	}

	// Aktuellen Level ermitteln, falls nicht angegeben
	currentLevel := request.CurrentLevel
	if currentLevel <= 0 {
		currentLevel = getCurrentSkillLevel(&character, request.Name, "skill")
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
		costResult.CharacterID = fmt.Sprintf("%d", character.ID)
		costResult.CharacterClass = characterClass
		costResult.SkillName = request.Name

		err = gsmaster.GetLernCostNextLevel(&tempRequest, &costResult, request.Reward, nextLevel, character.Rasse)
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

	// Prüfe, ob genügend PP vorhanden sind (PP der jeweiligen Fertigkeit)
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
		// Erstelle Notiz für Multi-Level Improvement
		levelCount := finalLevel - currentLevel
		var notes string
		if levelCount > 1 {
			notes = fmt.Sprintf("Fertigkeit '%s' von %d auf %d verbessert (%d Level)", request.Name, currentLevel, finalLevel, levelCount)
		} else {
			notes = fmt.Sprintf("Fertigkeit '%s' von %d auf %d verbessert", request.Name, currentLevel, finalLevel)
		}

		err = CreateAuditLogEntry(character.ID, "experience_points", currentEP, newEP, ReasonSkillImprovement, 0, notes)
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
		notes := fmt.Sprintf("Gold für Verbesserung von '%s' ausgegeben", request.Name)

		err = CreateAuditLogEntry(character.ID, "gold", currentGold, newGold, ReasonSkillImprovement, 0, notes)
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

	// PP abziehen wenn verwendet (PP der jeweiligen Fertigkeit)
	if totalPP > 0 {
		// Finde die richtige Fertigkeit und ziehe PP ab
		for i := range character.Fertigkeiten {
			if character.Fertigkeiten[i].Name == request.Name {
				character.Fertigkeiten[i].Pp -= totalPP
				if err := database.DB.Save(&character.Fertigkeiten[i]).Error; err != nil {
					respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren der Praxispunkte")
					return
				}
				break
			}
		}
		// Falls nicht in normalen Fertigkeiten gefunden, prüfe Waffenfertigkeiten
		for i := range character.Waffenfertigkeiten {
			if character.Waffenfertigkeiten[i].Name == request.Name {
				character.Waffenfertigkeiten[i].Pp -= totalPP
				if err := database.DB.Save(&character.Waffenfertigkeiten[i]).Error; err != nil {
					respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren der Praxispunkte")
					return
				}
				break
			}
		}
	}

	// Aktualisiere die Fertigkeit mit dem neuen Level
	if err := updateOrCreateSkill(&character, request.Name, finalLevel); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Aktualisieren der Fertigkeit: "+err.Error())
		return
	}

	// Charakter speichern
	if err := database.DB.Save(&character).Error; err != nil {
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

// LearnSpell lernt einen neuen Zauber und erstellt Audit-Log-Einträge
func LearnSpell(c *gin.Context) {
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

	cost, _, _, err := calculateSingleCost(&character, &costRequest)
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

// GetRewardTypes liefert verfügbare Belohnungsarten für ein bestimmtes Lernszenario
func GetRewardTypes(c *gin.Context) {
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

// GetAvailableSkills gibt alle verfügbaren Fertigkeiten mit Lernkosten zurück
func GetAvailableSkills(c *gin.Context) {
	characterID := c.Param("id")
	//rewardType := c.Query("reward_type")

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
		request := gsmaster.LernCostRequest{
			CharId:       character.ID,
			Name:         skill.Name,
			CurrentLevel: 0, // Nicht gelernt
			TargetLevel:  1, // Auf Level 1 lernen
			Type:         "skill",
			Action:       "learn",
			UsePP:        0,                       // Keine PP für neue Fertigkeiten
			UseGold:      0,                       // Keine Gold für neue Fertigkeiten
			Reward:       &[]string{"default"}[0], // Belohnungstyp aus Query-Parameter
		}
		levelResult := &gsmaster.SkillCostResultNew{}
		remainingPP := request.UsePP
		remainingGold := request.UseGold
		spellInfo := &models.SkillLearningInfo{}
		// Berechne Lernkosten mit GetLernCostNextLevel
		//epCost, goldCost := calculateSkillLearningCosts(skill, character, rewardType)
		err := calculateSkillLearnCostNewSystem(&request, levelResult, &remainingPP, &remainingGold, spellInfo)
		epCost := 10000
		goldCost := 50000
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

// GetAvailableSkills gibt alle verfügbaren Fertigkeiten mit Lernkosten zurück
func GetAvailableSkillsNewSystem(c *gin.Context) {
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
		epCost, goldCost := calculateSkillLearningCosts(skill, character, rewardType)

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

// calculateSkillLearningCosts berechnet die EP- und Goldkosten für das Lernen einer Fertigkeit mit GetLernCostNextLevel
func calculateSkillLearningCosts(skill models.Skill, character models.Char, rewardType string) (int, int) {
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
		CharacterClass: getCharacterClass(&character),
		SkillName:      skill.Name,
		Category:       skill.Category,
		Difficulty:     skill.Difficulty,
		TargetLevel:    1,
	}

	// Berechne Kosten mit GetLernCostNextLevel
	err := gsmaster.GetLernCostNextLevel(&request, &costResult, rewardTypePtr, 1, character.Typ)
	if err != nil {
		// Fallback zu Standard-Kosten bei Fehler
		epCost := 100
		goldCost := 50

		return epCost, goldCost
	}

	return costResult.EP, costResult.GoldCost
}
