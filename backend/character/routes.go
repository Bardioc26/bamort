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
	charGrp.GET("/:id/experience-wealth", GetCharacterExperienceAndWealth) // NewSystem
	charGrp.PUT("/:id/experience", UpdateCharacterExperience)              // NewSystem
	charGrp.PUT("/:id/wealth", UpdateCharacterWealth)                      // NewSystem

	// Audit-Log für Änderungen
	charGrp.GET("/:id/audit-log", GetCharacterAuditLog)   // Alle Änderungen oder gefiltert nach Feld (?field=experience_points)
	charGrp.GET("/:id/audit-log/stats", GetAuditLogStats) // Statistiken über Änderungen

	// im Frontend wir nur noch der neue Endpunkt benutzt
	//charGrp.POST("/lerncost", GetLernCost)              // alter Hauptendpunkt für alle Kostenberechnungen (verwendet lerningCostsData)
	charGrp.POST("/lerncost-new", GetLernCostNewSystem) // neuer Hauptendpunkt für alle Kostenberechnungen (verwendet neue Datenbank)
	charGrp.POST("/improve-skill-new", ImproveSkill)    // Fertigkeit verbessern

	// Lernen und Verbessern (mit automatischem Audit-Log)
	charGrp.POST("/:id/learn-skill-new", LearnSkill) // Fertigkeit lernen (neues System)
	//charGrp.POST("/:id/learn-skill", LearnSkillOld)  // Fertigkeit lernen (altes System)
	charGrp.POST("/:id/learn-spell-new", LearnSpell) // Zauber lernen (neues System)
	//charGrp.POST("/:id/learn-spell", LearnSpellOld)  // Zauber lernen (altes System)

	// Fertigkeiten-Information
	//charGrp.GET("/:id/available-skills", GetAvailableSkillsOld)               // Verfügbare Fertigkeiten mit Kosten (bereits gelernte ausgeschlossen)
	charGrp.POST("/available-skills-new", GetAvailableSkillsNewSystem)        // Verfügbare Fertigkeiten mit Kosten (bereits gelernte ausgeschlossen)
	charGrp.POST("/available-skills-creation", GetAvailableSkillsForCreation) // Verfügbare Fertigkeiten mit Lernkosten für Charaktererstellung
	charGrp.POST("/available-spells-creation", GetAvailableSpellsForCreation) // Verfügbare Zauber mit Lernkosten für Charaktererstellung
	charGrp.POST("/available-spells-new", GetAvailableSpellsNewSystem)        // Verfügbare Zauber mit Kosten (bereits gelernte ausgeschlossen)
	charGrp.GET("/spell-details", GetSpellDetails)                            // Detaillierte Informationen zu einem bestimmten Zauber

	// Belohnungsarten für verschiedene Lernszenarien
	charGrp.GET("/:id/reward-types", GetRewardTypesStatic) // Verfügbare Belohnungsarten je nach Kontext

	// Praxispunkte-Verwaltung
	charGrp.GET("/:id/practice-points", GetPracticePoints)     // NewSystem
	charGrp.PUT("/:id/practice-points", UpdatePracticePoints)  // NewSystem
	charGrp.POST("/:id/practice-points/add", AddPracticePoint) // NewSystem
	charGrp.POST("/:id/practice-points/use", UsePracticePoint) // NewSystem

	// System-Information
	//charGrp.GET("/character-classes", GetCharacterClassesHandlerOld)
	charGrp.GET("/skill-categories", GetSkillCategoriesHandlerStatic)

	// Character Creation
	charGrp.GET("/create-sessions", ListCharacterSessions)                          // Aktive Sessions für Benutzer auflisten
	charGrp.POST("/create-session", CreateCharacterSession)                         // Neue Charakter-Erstellungssession
	charGrp.GET("/create-session/:sessionId", GetCharacterSession)                  // Session-Daten abrufen
	charGrp.PUT("/create-session/:sessionId/basic", UpdateCharacterBasicInfo)       // Grundinformationen speichern
	charGrp.PUT("/create-session/:sessionId/attributes", UpdateCharacterAttributes) // Grundwerte speichern
	charGrp.PUT("/create-session/:sessionId/derived", UpdateCharacterDerivedValues) // Abgeleitete Werte speichern
	charGrp.PUT("/create-session/:sessionId/skills", UpdateCharacterSkills)         // Fertigkeiten speichern
	charGrp.POST("/create-session/:sessionId/finalize", FinalizeCharacterCreation)  // Charakter-Erstellung abschließen
	charGrp.DELETE("/create-session/:sessionId", DeleteCharacterSession)            // Session löschen

	// Reference Data für Character Creation
	charGrp.GET("/races", GetRaces)                                          // Verfügbare Rassen
	charGrp.GET("/classes", GetCharacterClasses)                             // Verfügbare Klassen
	charGrp.GET("/classes/learning-points", GetCharacterClassLearningPoints) // Lernpunkte für Charakterklasse
	charGrp.GET("/origins", GetOrigins)                                      // Verfügbare Herkünfte
	charGrp.GET("/beliefs", SearchBeliefs)                                   // Glaube-Suche

	// Derived Values Calculation
	charGrp.POST("/calculate-static-fields", CalculateStaticFields) // Berechnung ohne Würfelwürfe
	charGrp.POST("/calculate-rolled-field", CalculateRolledField)   // Berechnung mit Würfelwürfen
}
