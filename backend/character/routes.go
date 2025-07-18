package character

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	charGrp := r.Group("/characters")
	charGrp.GET("", ListCharacters)
	charGrp.POST("", CreateCharacter)
	charGrp.GET("/:id", GetCharacter)
	charGrp.PUT("/:id", UpdateCharacter)
	charGrp.DELETE("/:id", DeleteCharacter)

	// Kostenberechnung (konsolidiert)
	charGrp.POST("/:id/skill-cost", GetSkillCost)            // Hauptendpunkt für alle Kostenberechnungen
	charGrp.GET("/:id/improve", GetSkillNextLevelCosts)      // Legacy - für nächste Stufe
	charGrp.GET("/:id/improve/skill", GetSkillAllLevelCosts) // Legacy - für alle Stufen
	charGrp.GET("/:id/learn", GetLearnCost)                  // Legacy - einfache Lernkosten

	// Praxispunkte-Verwaltung
	charGrp.GET("/:id/practice-points", GetPracticePoints)
	charGrp.PUT("/:id/practice-points", UpdatePracticePoints)
	charGrp.POST("/:id/practice-points/add", AddPracticePoint)
	charGrp.POST("/:id/practice-points/use", UsePracticePoint)

	// System-Information
	charGrp.GET("/character-classes", GetCharacterClassesHandler)
	charGrp.GET("/skill-categories", GetSkillCategoriesHandler)
}
