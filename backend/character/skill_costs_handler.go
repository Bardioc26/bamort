package character

import (
	"bamort/gsmaster"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SkillCostRequest struct {
	Name         string `json:"name"`          // Name der Fertigkeit
	CurrentLevel int    `json:"current_level"` // Aktueller Wert (nur für Verbesserung)
	Type         string `json:"type"`          // 'skill' oder 'spell'
	Action       string `json:"action"`        // 'learn' oder 'improve'
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
		respondWithError(c, http.StatusBadRequest, "Ungültige Anfrageparameter")
		return
	}

	// Überprüfen, ob alle notwendigen Parameter vorhanden sind
	if request.Name == "" {
		respondWithError(c, http.StatusBadRequest, "Fertigkeitsname ist erforderlich")
		return
	}

	// Kosten berechnen basierend auf Aktion und Typ
	var cost *gsmaster.LearnCost
	var err error

	switch {
	case request.Action == "learn" && request.Type == "skill":
		// Kosten zum Erlernen einer neuen Fertigkeit
		cost, err = gsmaster.CalculateDetailedSkillLearningCost(request.Name, character.Typ)

	case request.Action == "improve" && request.Type == "skill":
		// Kosten zum Verbessern einer vorhandenen Fertigkeit
		cost, err = gsmaster.CalculateDetailedSkillImprovementCost(request.Name, character.Typ, request.CurrentLevel)

	case request.Action == "learn" && request.Type == "spell":
		// Kosten zum Erlernen eines neuen Zaubers
		cost, err = gsmaster.CalculateDetailedSpellLearningCost(request.Name, character.Typ)

	default:
		respondWithError(c, http.StatusBadRequest, "Ungültige Kombination aus Aktion und Typ")
		return
	}

	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Fehler bei der Kostenberechnung: "+err.Error())
		return
	}

	// Ergebnis zurückgeben
	c.JSON(http.StatusOK, cost)
}
