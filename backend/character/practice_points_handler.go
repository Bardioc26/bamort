package character

import (
	"bamort/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPracticePoints gibt die verfügbaren Praxispunkte eines Charakters zurück
func GetPracticePoints(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	c.JSON(http.StatusOK, character.Praxispunkte)
}

// UpdatePracticePoints aktualisiert die Praxispunkte eines Charakters
func UpdatePracticePoints(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Request-Parameter abrufen
	var praxispunkte []Praxispunkt
	if err := c.ShouldBindJSON(&praxispunkte); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Praxispunkt-Daten: "+err.Error())
		return
	}

	// Praxispunkte aktualisieren
	character.Praxispunkte = praxispunkte

	// Charakter in der Datenbank speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern der Praxispunkte: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, character.Praxispunkte)
}

// AddPracticePoint fügt einen Praxispunkt zu einer Kategorie hinzu
func AddPracticePoint(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Request-Parameter abrufen
	type AddPPRequest struct {
		Kategorie string `json:"kategorie" binding:"required"`
		Anzahl    int    `json:"anzahl"`
	}

	var request AddPPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	if request.Anzahl <= 0 {
		request.Anzahl = 1
	}

	// Praxispunkt zur entsprechenden Kategorie hinzufügen
	found := false
	for i := range character.Praxispunkte {
		if character.Praxispunkte[i].Kategorie == request.Kategorie {
			character.Praxispunkte[i].Anzahl += request.Anzahl
			found = true
			break
		}
	}

	if !found {
		// Neue Kategorie hinzufügen
		characterIDUint, _ := strconv.ParseUint(charID, 10, 32)
		newPP := Praxispunkt{
			Kategorie: request.Kategorie,
			Anzahl:    request.Anzahl,
		}
		newPP.CharacterID = uint(characterIDUint)
		character.Praxispunkte = append(character.Praxispunkte, newPP)
	}

	// Charakter in der Datenbank speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern der Praxispunkte: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, character.Praxispunkte)
}

// UsePracticePoint verbraucht Praxispunkte aus einer Kategorie
func UsePracticePoint(c *gin.Context) {
	// Charakter-ID aus der URL abrufen
	charID := c.Param("id")

	// Charakter aus der Datenbank laden
	var character Char
	if err := character.FirstID(charID); err != nil {
		respondWithError(c, http.StatusNotFound, "Charakter nicht gefunden")
		return
	}

	// Request-Parameter abrufen
	type UsePPRequest struct {
		Kategorie string `json:"kategorie" binding:"required"`
		Anzahl    int    `json:"anzahl"`
	}

	var request UsePPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter: "+err.Error())
		return
	}

	if request.Anzahl <= 0 {
		request.Anzahl = 1
	}

	// Praxispunkt von der entsprechenden Kategorie abziehen
	found := false
	for i := range character.Praxispunkte {
		if character.Praxispunkte[i].Kategorie == request.Kategorie {
			if character.Praxispunkte[i].Anzahl >= request.Anzahl {
				character.Praxispunkte[i].Anzahl -= request.Anzahl
				found = true
			} else {
				respondWithError(c, http.StatusBadRequest, "Nicht genügend Praxispunkte verfügbar")
				return
			}
			break
		}
	}

	if !found {
		respondWithError(c, http.StatusBadRequest, "Keine Praxispunkte in dieser Kategorie vorhanden")
		return
	}

	// Charakter in der Datenbank speichern
	if err := database.DB.Save(&character).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Fehler beim Speichern der Praxispunkte: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, character.Praxispunkte)
}
