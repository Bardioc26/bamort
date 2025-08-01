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

	// Erfahrung und Vermögen
	charGrp.GET("/:id/experience-wealth", GetCharacterExperienceAndWealth)
	charGrp.PUT("/:id/experience", UpdateCharacterExperience)
	charGrp.PUT("/:id/wealth", UpdateCharacterWealth)

	// Audit-Log für Änderungen
	charGrp.GET("/:id/audit-log", GetCharacterAuditLog)   // Alle Änderungen oder gefiltert nach Feld (?field=experience_points)
	charGrp.GET("/:id/audit-log/stats", GetAuditLogStats) // Statistiken über Änderungen

	// im Frontend wir nur noch der neue Endpunkt benutzt
	//charGrp.POST("/lerncost", GetLernCost)              // alter Hauptendpunkt für alle Kostenberechnungen (verwendet lerningCostsData)
	charGrp.POST("/lerncost-new", GetLernCostNewSystem) // neuer Hauptendpunkt für alle Kostenberechnungen (verwendet neue Datenbank)
	charGrp.POST("/improve-skill", ImproveSkillOld)     // Fertigkeit verbessern

	// Lernen und Verbessern (mit automatischem Audit-Log)
	charGrp.POST("/:id/learn-skill", LearnSkillOld) // Fertigkeit lernen
	charGrp.POST("/:id/learn-spell", LearnSpellOld) // Zauber lernen

	// Fertigkeiten-Information
	charGrp.GET("/:id/available-skills", GetAvailableSkillsOld)        // Verfügbare Fertigkeiten mit Kosten (bereits gelernte ausgeschlossen)
	charGrp.POST("/available-skills-new", GetAvailableSkillsNewSystem) // Verfügbare Fertigkeiten mit Kosten (bereits gelernte ausgeschlossen)
	charGrp.POST("/available-spells-new", GetAvailableSpellsNewSystem) // Verfügbare Zauber mit Kosten (bereits gelernte ausgeschlossen)

	// Belohnungsarten für verschiedene Lernszenarien
	charGrp.GET("/:id/reward-types", GetRewardTypesOld) // Verfügbare Belohnungsarten je nach Kontext

	// Praxispunkte-Verwaltung
	charGrp.GET("/:id/practice-points", GetPracticePoints)
	charGrp.PUT("/:id/practice-points", UpdatePracticePoints)
	charGrp.POST("/:id/practice-points/add", AddPracticePoint)
	charGrp.POST("/:id/practice-points/use", UsePracticePoint)

	// System-Information
	charGrp.GET("/character-classes", GetCharacterClassesHandler)
	charGrp.GET("/skill-categories", GetSkillCategoriesHandler)
}
