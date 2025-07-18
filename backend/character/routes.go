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
	charGrp.GET("/:id/improve", GetSkillNextLevelCosts)
	charGrp.GET("/:id/improve/skill", GetSkillAllLevelCosts)
	charGrp.GET("/:id/learn", GetLearnCost)
	charGrp.POST("/:id/skill-cost", GetSkillCost) // Neuer Endpunkt für detaillierte Kostenberechnung

	// Praxispunkte-Verwaltung
	charGrp.GET("/:id/practice-points", GetPracticePoints)
	charGrp.PUT("/:id/practice-points", UpdatePracticePoints)
	charGrp.POST("/:id/practice-points/add", AddPracticePoint)
	charGrp.POST("/:id/practice-points/use", UsePracticePoint)
}
