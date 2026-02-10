package character

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
// GetCharacterClassesHandlerOld gibt alle verfügbaren Charakterklassen zurück
func GetCharacterClassesHandlerOld(c *gin.Context) {
	// Vereinfachte Antwort mit den drei Hauptklassen
	classes := map[string]interface{}{
		"Sp": map[string]interface{}{
			"code":        "Sp",
			"name":        "Spitzbube",
			"description": "Experte in Halbwelt-Fertigkeiten",
		},
		"Hx": map[string]interface{}{
			"code":        "Hx",
			"name":        "Hexer",
			"description": "Zauberer mit Beherrschungs- und Veränderungs-Zaubern",
		},
		"PS": map[string]interface{}{
			"code":        "PS",
			"name":        "Priester Streiter",
			"description": "Kämpfender Priester mit Wunder-Zaubern",
		},
	}

	c.JSON(http.StatusOK, gin.H{"character_classes": classes})
}
*/
// GetSkillCategoriesHandlerStatic gibt alle verfügbaren Fertigkeitskategorien zurück
// GetSkillCategoriesHandlerStatic godoc
// @Summary Get skill categories
// @Description Returns list of all skill categories (static data)
// @Tags Reference Data
// @Produce json
// @Success 200 {array} object "List of skill categories"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/characters/skill-categories [get]
func GetSkillCategoriesHandlerStatic(c *gin.Context) {
	categories := map[string]interface{}{
		"Alltag": map[string]interface{}{
			"name":        "Alltag",
			"description": "Tägliche Fertigkeiten wie Reiten, Schwimmen",
		},
		"Freiland": map[string]interface{}{
			"name":        "Freiland",
			"description": "Natur-Fertigkeiten wie Spurensuche, Tarnen",
		},
		"Halbwelt": map[string]interface{}{
			"name":        "Halbwelt",
			"description": "Diebstahlsfertigkeiten wie Stehlen, Gaukeln",
		},
		"Kampf": map[string]interface{}{
			"name":        "Kampf",
			"description": "Kampffertigkeiten wie Reiterkampf, Athletik",
		},
		"Körper": map[string]interface{}{
			"name":        "Körper",
			"description": "Körperliche Fertigkeiten wie Klettern, Akrobatik",
		},
		"Sozial": map[string]interface{}{
			"name":        "Sozial",
			"description": "Soziale Fertigkeiten wie Menschenkenntnis, Etikette",
		},
		"Unterwelt": map[string]interface{}{
			"name":        "Unterwelt",
			"description": "Unterwelt-Fertigkeiten wie Gassenwissen",
		},
		"Waffen": map[string]interface{}{
			"name":        "Waffen",
			"description": "Waffenfertigkeiten wie Dolch, Bogen, Fechten",
		},
		"Wissen": map[string]interface{}{
			"name":        "Wissen",
			"description": "Wissensfertigkeiten wie Geschichte, Rechtskunde",
		},
	}

	c.JSON(http.StatusOK, gin.H{"skill_categories": categories})
}
