package maintenance

import (
	"bamort/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	charGrp := r.Group("/maintenance")
	charGrp.Use(user.RequireMaintainer())
	{
		charGrp.GET("/gsm-believes", GetBelieves)
		charGrp.PUT("/gsm-believes/:id", UpdateBelieve)
		charGrp.GET("/game-systems", GetGameSystems)
		charGrp.PUT("/game-systems/:id", UpdateGameSystem)
		charGrp.GET("/gsm-lit-sources", GetLitSources)
		charGrp.PUT("/gsm-lit-sources/:id", UpdateLitSource)
		charGrp.GET("/gsm-misc", GetMisc)
		charGrp.PUT("/gsm-misc/:id", UpdateMisc)
		charGrp.GET("/skill-improvement-cost2", GetSkillImprovementCost2)
		charGrp.PUT("/skill-improvement-cost2/:id", UpdateSkillImprovementCost2)
		charGrp.GET("/setupcheck", SetupCheck)
		charGrp.GET("/setupcheck-dev", SetupCheckDev)
		charGrp.GET("/mktestdata", MakeTestdataFromLive)
		charGrp.GET("/reconndb", ReconnectDataBase) // Datenbank neu verbinden
		charGrp.GET("/reloadenv", ReloadENV)
		charGrp.POST("/transfer-sqlite-to-mariadb", TransferSQLiteToMariaDB) // Transfer data from SQLite to MariaDB
		//charGrp.POST("/populate-class-learning-points", PopulateClassLearningPoints) // Populate class learning points from hardcoded data
		/*
			//nur zur einmaligen Ausführung, um das Lernkosten-System zu initialisieren
			charGrp.POST("/initialize-learning-costs", InitializeLearningCosts)
			// Zur Überprüfung der Lernkosten-Daten
			charGrp.GET("/learning-costs-summary", gsmaster.GetLearningCostsSummaryHandler)
		*/
	}
}
