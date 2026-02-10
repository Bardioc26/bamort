package character

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PracticePointResponse repräsentiert die Antwort für Praxispunkte einer Fertigkeit
type PracticePointResponse struct {
	SkillName string `json:"skill_name"`
	Amount    int    `json:"amount"`
}

// PracticePointActionResponse repräsentiert die erweiterte Antwort für PP-Aktionen
type PracticePointActionResponse struct {
	Success        bool                    `json:"success"`
	Message        string                  `json:"message"`
	RequestedSkill string                  `json:"requested_skill"` // Ursprünglich angefragter Name
	TargetSkill    string                  `json:"target_skill"`    // Tatsächlich betroffene Fertigkeit
	IsSpell        bool                    `json:"is_spell"`        // Ob es sich um einen Zauber handelt
	PracticePoints []PracticePointResponse `json:"practice_points"` // Aktuelle PP-Liste
}

// GetPracticePoints gibt die verfügbaren Praxispunkte eines Charakters zurück
// GetPracticePoints godoc
// @Summary Get practice points
// @Description Returns all practice points for a character
// @Tags Characters
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {array} PracticePointResponse "Practice points list"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Security BearerAuth
// @Router /api/characters/{id}/practice-points [get]
func GetPracticePoints(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Praxispunkte aus den Fertigkeiten extrahieren
	var practicePoints []PracticePointResponse
	for _, skill := range character.Fertigkeiten {
		if skill.Pp > 0 {
			practicePoints = append(practicePoints, PracticePointResponse{
				SkillName: skill.Name,
				Amount:    skill.Pp,
			})
		}
	}

	c.JSON(http.StatusOK, practicePoints)
}

// UpdatePracticePoints aktualisiert die Praxispunkte eines Charakters
// UpdatePracticePoints godoc
// @Summary Update practice points
// @Description Updates the practice points for a character (owner only)
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param practice_points body []PracticePointResponse true "Practice points data"
// @Success 200 {object} map[string]string "Practice points updated successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Security BearerAuth
// @Router /api/characters/{id}/practice-points [put]
func UpdatePracticePoints(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Check ownership
	if !checkCharacterOwnership(c, &character) {
		return
	}

	// Request-Parameter abrufen
	var practicePoints []PracticePointResponse
	if err := c.ShouldBindJSON(&practicePoints); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Praxispunkt-Daten: "+err.Error())
		return
	}

	// Alle Fertigkeiten durchgehen und Praxispunkte zurücksetzen
	/*
		for i := range character.Fertigkeiten {
			character.Fertigkeiten[i].Pp = 0
		}
	*/

	// Neue Praxispunkte setzen
	for _, pp := range practicePoints {
		for i := range character.Fertigkeiten {
			if character.Fertigkeiten[i].Name == pp.SkillName {
				character.Fertigkeiten[i].Pp = pp.Amount
				break
			}
		}
	}

	// Charakter in der Datenbank speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern der Praxispunkte: "+err.Error())
		return
	}

	// Aktualisierte Praxispunkte zurückgeben
	var updatedPracticePoints []PracticePointResponse
	for _, skill := range character.Fertigkeiten {
		if skill.Pp > 0 {
			updatedPracticePoints = append(updatedPracticePoints, PracticePointResponse{
				SkillName: skill.Name,
				Amount:    skill.Pp,
			})
		}
	}

	c.JSON(http.StatusOK, updatedPracticePoints)
}

// AddPracticePoint fügt einen Praxispunkt zu einer Fertigkeit hinzu
// TODO prüfe speichern der PP für Spells
// AddPracticePoint godoc
// @Summary Add practice point
// @Description Adds a new practice point to a character (owner only)
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param practice_point body PracticePointResponse true "Practice point data"
// @Success 200 {object} PracticePointResponse "Added practice point"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Security BearerAuth
// @Router /api/characters/{id}/practice-points/add [post]
func AddPracticePoint(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Check ownership
	if !checkCharacterOwnership(c, &character) {
		return
	}

	// Request-Parameter abrufen
	type AddPPRequest struct {
		SkillName string `json:"skill_name" binding:"required"`
		Amount    int    `json:"amount"`
	}

	var request AddPPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	if request.Amount <= 0 {
		request.Amount = 1
	}

	// Prüfen, ob es sich um einen Zauber handelt
	var targetSkillName string
	var isSpellFlag bool
	if isSpellNewSystem(request.SkillName) {
		// Bei Zaubern: PP werden der entsprechenden Zaubergruppe zugeordnet
		targetSkillName = getSpellCategoryNewSystem(request.SkillName)
		isSpellFlag = true
	} else {
		// Bei normalen Fertigkeiten: PP werden direkt der Fertigkeit zugeordnet
		targetSkillName = request.SkillName
		isSpellFlag = false
	}

	// Praxispunkt zur entsprechenden Fertigkeit hinzufügen
	found := false
	for i := range character.Fertigkeiten {
		if character.Fertigkeiten[i].Name == targetSkillName {
			character.Fertigkeiten[i].Pp += request.Amount
			found = true
			break
		}
	}

	if !found {
		respondWithError(c, http.StatusBadRequest, "Fertigkeit nicht gefunden: "+targetSkillName)
		return
	}

	// Fertigkeiten explizit speichern
	for i := range character.Fertigkeiten {
		if character.Fertigkeiten[i].Name == targetSkillName {
			if err := database.DB.Save(&character.Fertigkeiten[i]).Error; err != nil {
				respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern der Fertigkeit: "+err.Error())
				return
			}
			break
		}
	}

	// Aktualisierte Praxispunkte sammeln
	var practicePoints []PracticePointResponse
	for _, skill := range character.Fertigkeiten {
		if skill.Pp > 0 {
			practicePoints = append(practicePoints, PracticePointResponse{
				SkillName: skill.Name,
				Amount:    skill.Pp,
			})
		}
	}

	// Erfolgreiche Response mit Details zurückgeben
	var message string
	if isSpellFlag {
		message = "Praxispunkt für Zauber '" + request.SkillName + "' wurde der Zaubergruppe '" + targetSkillName + "' hinzugefügt"
	} else {
		message = "Praxispunkt für Fertigkeit '" + targetSkillName + "' hinzugefügt"
	}

	response := PracticePointActionResponse{
		Success:        true,
		Message:        message,
		RequestedSkill: request.SkillName,
		TargetSkill:    targetSkillName,
		IsSpell:        isSpellFlag,
		PracticePoints: practicePoints,
	}

	c.JSON(http.StatusOK, response)
}

// UsePracticePoint verbraucht Praxispunkte für eine spezifische Fertigkeit
// UsePracticePoint godoc
// @Summary Use practice point
// @Description Marks a practice point as used and applies it to skill improvement (owner only)
// @Tags Characters
// @Accept json
// @Produce json
// @Param id path int true "Character ID"
// @Param use_data body object{practice_point_id=int,skill_id=int} true "Practice point usage data"
// @Success 200 {object} map[string]string "Practice point used successfully"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Security BearerAuth
// @Router /api/characters/{id}/practice-points/use [post]
func UsePracticePoint(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character models.Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Check ownership
	if !checkCharacterOwnership(c, &character) {
		return
	}

	// Request-Parameter abrufen
	type UsePPRequest struct {
		SkillName string `json:"skill_name" binding:"required"`
		Amount    int    `json:"amount"`
	}

	var request UsePPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	if request.Amount <= 0 {
		request.Amount = 1
	}

	// Prüfen, ob es sich um einen Zauber handelt
	var targetSkillName string
	var isSpellFlag bool
	if isSpellNewSystem(request.SkillName) {
		// Bei Zaubern: PP werden von der entsprechenden Zaubergruppe abgezogen
		targetSkillName = getSpellCategoryNewSystem(request.SkillName)
		isSpellFlag = true
	} else {
		// Bei normalen Fertigkeiten: PP werden direkt von der Fertigkeit abgezogen
		targetSkillName = request.SkillName
		isSpellFlag = false
	}

	// Praxispunkt von der entsprechenden Fertigkeit abziehen
	found := false
	for i := range character.Fertigkeiten {
		if character.Fertigkeiten[i].Name == targetSkillName {
			if character.Fertigkeiten[i].Pp >= request.Amount {
				character.Fertigkeiten[i].Pp -= request.Amount
				found = true
			} else {
				respondWithError(c, http.StatusBadRequest, "Nicht genügend Praxispunkte verfügbar")
				return
			}
			break
		}
	}

	if !found {
		respondWithError(c, http.StatusBadRequest, "Fertigkeit nicht gefunden: "+targetSkillName)
		return
	}

	// Fertigkeiten explizit speichern
	for i := range character.Fertigkeiten {
		if character.Fertigkeiten[i].Name == targetSkillName {
			if err := database.DB.Save(&character.Fertigkeiten[i]).Error; err != nil {
				respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern der Fertigkeit: "+err.Error())
				return
			}
			break
		}
	}

	// Erfolgreiche Antwort mit detaillierten Informationen und aktueller PP-Liste
	var practicePoints []PracticePointResponse
	for _, skill := range character.Fertigkeiten {
		if skill.Pp > 0 {
			practicePoints = append(practicePoints, PracticePointResponse{
				SkillName: skill.Name,
				Amount:    skill.Pp,
			})
		}
	}

	response := PracticePointActionResponse{
		Success:        true,
		Message:        fmt.Sprintf("%d Übungspunkte erfolgreich von %s verwendet", request.Amount, targetSkillName),
		RequestedSkill: request.SkillName,
		TargetSkill:    targetSkillName,
		IsSpell:        isSpellFlag,
		PracticePoints: practicePoints,
	}

	c.JSON(http.StatusOK, response)
}
