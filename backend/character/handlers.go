package character

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"bamort/skills"

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
	var characters []Char
	var listOfChars []CharList
	if err := database.DB.Find(&characters).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve characters")
		return
	}
	for i := range characters {
		listOfChars = append(listOfChars, CharList{
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
	var character Char
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
	var character Char
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
	var character Char

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
	var character Char
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
func AddFertigkeit(charID uint, fertigkeit *skills.Fertigkeit) error {
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

func ToFeChar(object *Char) *FeChar {
	feC := &FeChar{
		Char: *object,
	}
	skills, innateSkills, categories := splitSkills(object.Fertigkeiten)
	feC.Fertigkeiten = skills
	feC.InnateSkills = innateSkills
	feC.CategorizedSkills = categories
	return feC
}

func splitSkills(object []skills.Fertigkeit) ([]skills.Fertigkeit, []skills.Fertigkeit, map[string][]skills.Fertigkeit) {
	var normSkills []skills.Fertigkeit
	var innateSkills []skills.Fertigkeit
	//var categories map[string][]skills.Fertigkeit
	categories := make(map[string][]skills.Fertigkeit)
	for _, skill := range object {
		gsmsk := skill.GetGsm()
		if gsmsk.Improvable {
			category := "Unkategorisiert"
			if gsmsk.ID != 0 && gsmsk.Category != "" {
				category = gsmsk.Category
			}
			normSkills = append(normSkills, skill)
			if _, exists := categories[category]; !exists {
				categories[category] = make([]skills.Fertigkeit, 0)
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
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	// Load the skill from the request
	var s skills.Fertigkeit
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var skill gsmaster.Skill
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
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve character")
		return
	}

	// Load the spell from the request
	var s skills.Zauber
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	var spell gsmaster.Spell
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
	var s skills.Fertigkeit
	if err := c.ShouldBindJSON(&s); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Return the updated character
	c.JSON(http.StatusOK, character.Fertigkeiten)
}

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
	var character Char

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
		ExperiencePoints: character.Erfahrungsschatz.Value,
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
	var character Char

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
		oldValue = character.Erfahrungsschatz.Value
	}

	// Aktualisiere oder erstelle Erfahrungsschatz
	if character.Erfahrungsschatz.ID == 0 {
		// Erstelle neuen Erfahrungsschatz
		character.Erfahrungsschatz = Erfahrungsschatz{
			BamortCharTrait: models.BamortCharTrait{
				CharacterID: character.ID,
			},
			Value: req.ExperiencePoints,
		}
		if err := database.DB.Create(&character.Erfahrungsschatz).Error; err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to create experience record")
			return
		}
	} else {
		// Aktualisiere existierenden Erfahrungsschatz
		character.Erfahrungsschatz.Value = req.ExperiencePoints
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
	var character Char

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
		character.Vermoegen = Vermoegen{
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

// Learn and Improve handlers with automatic audit logging

// LearnSkillRequest definiert die Struktur für das Lernen einer Fertigkeit
type LearnSkillRequest struct {
	Name  string `json:"name" binding:"required"`
	Notes string `json:"notes,omitempty"`
	UsePP int    `json:"use_pp,omitempty"`
}

// ImproveSkillRequest definiert die Struktur für das Verbessern einer Fertigkeit
type ImproveSkillRequest struct {
	Name         string `json:"name" binding:"required"`
	CurrentLevel int    `json:"current_level,omitempty"`
	Notes        string `json:"notes,omitempty"`
	UsePP        int    `json:"use_pp,omitempty"`
}

// LearnSpellRequest definiert die Struktur für das Lernen eines Zaubers
type LearnSpellRequest struct {
	Name  string `json:"name" binding:"required"`
	Notes string `json:"notes,omitempty"`
}

// ImproveSpellRequest definiert die Struktur für das Verbessern eines Zaubers
type ImproveSpellRequest struct {
	Name         string `json:"name" binding:"required"`
	CurrentLevel int    `json:"current_level,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

// LearnSkill lernt eine neue Fertigkeit und erstellt Audit-Log-Einträge
func LearnSkill(c *gin.Context) {
	charID := c.Param("id")
	var character Char

	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	var request LearnSkillRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Berechne Kosten mit GetSkillCost
	costRequest := SkillCostRequest{
		Name:   request.Name,
		Type:   "skill",
		Action: "learn",
		UsePP:  request.UsePP,
	}

	cost, _, _, err := calculateSingleCost(&character, &costRequest)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
		return
	}

	// Prüfe, ob genügend EP vorhanden sind
	currentEP := character.Erfahrungsschatz.Value
	if currentEP < cost.Ep {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Erfahrungspunkte vorhanden")
		return
	}

	// Prüfe, ob genügend Gold vorhanden ist
	currentGold := character.Vermoegen.Goldstücke
	if currentGold < cost.Money {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Gold vorhanden")
		return
	}

	// EP abziehen und Audit-Log erstellen
	newEP := currentEP - cost.Ep
	if cost.Ep > 0 {
		notes := fmt.Sprintf("Fertigkeit '%s' gelernt", request.Name)
		if request.Notes != "" {
			notes += " - " + request.Notes
		}

		err = CreateAuditLogEntry(character.ID, "experience_points", currentEP, newEP, ReasonSkillLearning, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		character.Erfahrungsschatz.Value = newEP
	}

	// Gold abziehen und Audit-Log erstellen
	newGold := currentGold - cost.Money
	if cost.Money > 0 {
		notes := fmt.Sprintf("Gold für Fertigkeit '%s' ausgegeben", request.Name)
		if request.Notes != "" {
			notes += " - " + request.Notes
		}

		err = CreateAuditLogEntry(character.ID, "gold", currentGold, newGold, ReasonSkillLearning, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		character.Vermoegen.Goldstücke = newGold
	}

	// TODO: Hier sollte die Fertigkeit dem Charakter hinzugefügt werden
	// Das hängt davon ab, wie Fertigkeiten in der Datenbank gespeichert werden

	// Charakter speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Fertigkeit erfolgreich gelernt",
		"skill_name":     request.Name,
		"ep_cost":        cost.Ep,
		"gold_cost":      cost.Money,
		"remaining_ep":   newEP,
		"remaining_gold": newGold,
	})
}

// ImproveSkill verbessert eine bestehende Fertigkeit und erstellt Audit-Log-Einträge
func ImproveSkill(c *gin.Context) {
	charID := c.Param("id")
	var character Char

	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	var request ImproveSkillRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Aktuellen Level ermitteln, falls nicht angegeben
	currentLevel := request.CurrentLevel
	if currentLevel <= 0 {
		currentLevel = getCurrentSkillLevel(&character, request.Name, "skill")
		if currentLevel == -1 {
			respondWithError(c, http.StatusBadRequest, "Fertigkeit nicht bei diesem Charakter vorhanden")
			return
		}
	}

	// Berechne Kosten mit GetSkillCost
	costRequest := SkillCostRequest{
		Name:         request.Name,
		Type:         "skill",
		Action:       "improve",
		CurrentLevel: currentLevel,
		UsePP:        request.UsePP,
	}

	cost, _, _, err := calculateSingleCost(&character, &costRequest)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
		return
	}

	// Prüfe, ob genügend EP vorhanden sind
	currentEP := character.Erfahrungsschatz.Value
	if currentEP < cost.Ep {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Erfahrungspunkte vorhanden")
		return
	}

	// Prüfe, ob genügend Gold vorhanden ist
	currentGold := character.Vermoegen.Goldstücke
	if currentGold < cost.Money {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Gold vorhanden")
		return
	}

	// EP abziehen und Audit-Log erstellen
	newEP := currentEP - cost.Ep
	if cost.Ep > 0 {
		notes := fmt.Sprintf("Fertigkeit '%s' von %d auf %d verbessert", request.Name, currentLevel, currentLevel+1)
		if request.Notes != "" {
			notes += " - " + request.Notes
		}

		err = CreateAuditLogEntry(character.ID, "experience_points", currentEP, newEP, ReasonSkillImprovement, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		character.Erfahrungsschatz.Value = newEP
	}

	// Gold abziehen und Audit-Log erstellen
	newGold := currentGold - cost.Money
	if cost.Money > 0 {
		notes := fmt.Sprintf("Gold für Verbesserung von '%s' ausgegeben", request.Name)
		if request.Notes != "" {
			notes += " - " + request.Notes
		}

		err = CreateAuditLogEntry(character.ID, "gold", currentGold, newGold, ReasonSkillImprovement, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		character.Vermoegen.Goldstücke = newGold
	}

	// TODO: Hier sollte die Fertigkeit des Charakters verbessert werden
	// Das hängt davon ab, wie Fertigkeiten in der Datenbank gespeichert werden

	// Charakter speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Fertigkeit erfolgreich verbessert",
		"skill_name":     request.Name,
		"from_level":     currentLevel,
		"to_level":       currentLevel + 1,
		"ep_cost":        cost.Ep,
		"gold_cost":      cost.Money,
		"remaining_ep":   newEP,
		"remaining_gold": newGold,
	})
}

// LearnSpell lernt einen neuen Zauber und erstellt Audit-Log-Einträge
func LearnSpell(c *gin.Context) {
	charID := c.Param("id")
	var character Char

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
	currentEP := character.Erfahrungsschatz.Value
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
		character.Erfahrungsschatz.Value = newEP
	}

	// TODO: Hier sollte der Zauber dem Charakter hinzugefügt werden

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

// ImproveSpell verbessert einen bestehenden Zauber und erstellt Audit-Log-Einträge
func ImproveSpell(c *gin.Context) {
	charID := c.Param("id")
	var character Char

	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	var request ImproveSpellRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	// Aktuellen Level ermitteln, falls nicht angegeben
	currentLevel := request.CurrentLevel
	if currentLevel <= 0 {
		currentLevel = getCurrentSkillLevel(&character, request.Name, "spell")
		if currentLevel == -1 {
			respondWithError(c, http.StatusBadRequest, "Zauber nicht bei diesem Charakter vorhanden")
			return
		}
	}

	// Berechne Kosten mit GetSkillCost
	costRequest := SkillCostRequest{
		Name:         request.Name,
		Type:         "spell",
		Action:       "improve",
		CurrentLevel: currentLevel,
	}

	cost, _, _, err := calculateSingleCost(&character, &costRequest)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
		return
	}

	// Prüfe, ob genügend EP vorhanden sind
	currentEP := character.Erfahrungsschatz.Value
	if currentEP < cost.Ep {
		respondWithError(c, http.StatusBadRequest, "Nicht genügend Erfahrungspunkte vorhanden")
		return
	}

	// EP abziehen und Audit-Log erstellen
	newEP := currentEP - cost.Ep
	if cost.Ep > 0 {
		notes := fmt.Sprintf("Zauber '%s' von %d auf %d verbessert", request.Name, currentLevel, currentLevel+1)
		if request.Notes != "" {
			notes += " - " + request.Notes
		}

		err = CreateAuditLogEntry(character.ID, "experience_points", currentEP, newEP, ReasonSpellImprovement, 0, notes)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Fehler beim Erstellen des Audit-Log-Eintrags")
			return
		}
		character.Erfahrungsschatz.Value = newEP
	}

	// TODO: Hier sollte der Zauber des Charakters verbessert werden

	// Charakter speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern des Charakters")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Zauber erfolgreich verbessert",
		"spell_name":   request.Name,
		"from_level":   currentLevel,
		"to_level":     currentLevel + 1,
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
		// Neue Fertigkeit lernen - meist nur EP oder Gold
		rewardTypes = append(rewardTypes,
			gin.H{"value": "ep", "label": "Erfahrungspunkte verwenden", "description": "Verwende EP zum Lernen"},
			gin.H{"value": "gold", "label": "Gold verwenden", "description": "Bezahle einen Lehrer mit Gold"},
		)

	case "spell":
		// Zauber - mehr Optionen including Ritual
		rewardTypes = append(rewardTypes,
			gin.H{"value": "ep", "label": "Erfahrungspunkte verwenden", "description": "Verwende EP zum Verbessern"},
			gin.H{"value": "gold", "label": "Gold verwenden", "description": "Bezahle einen Zauberlehrer"},
			gin.H{"value": "pp", "label": "Praxispunkte verwenden", "description": "Nutze gesammelte Praxis"},
			gin.H{"value": "mixed", "label": "Gemischt (EP + PP)", "description": "Kombiniere EP und PP für reduzierten Aufwand"},
		)

	case "improve":
		// Fertigkeit verbessern - Standard-Optionen
		rewardTypes = append(rewardTypes,
			gin.H{"value": "ep", "label": "Erfahrungspunkte verwenden", "description": "Verwende EP zum Verbessern"},
			gin.H{"value": "gold", "label": "Gold verwenden", "description": "Bezahle einen Lehrer"},
			gin.H{"value": "pp", "label": "Praxispunkte verwenden", "description": "Nutze gesammelte Praxis"},
			gin.H{"value": "mixed", "label": "Gemischt (EP + PP)", "description": "Kombiniere EP und PP für reduzierten Aufwand"},
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
