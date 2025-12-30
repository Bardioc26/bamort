package maintenance

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	charGrp := r.Group("/maintenance")
	charGrp.GET("/setupcheck", SetupCheck)
	charGrp.GET("/setupcheck-dev", SetupCheckDev)
	charGrp.GET("/mktestdata", MakeTestdataFromLive)
	charGrp.GET("/reconndb", ReconnectDataBase) // Datenbank neu verbinden
	charGrp.GET("/reloadenv", ReloadENV)
	charGrp.POST("/transfer-sqlite-to-mariadb", TransferSQLiteToMariaDB) // Transfer data from SQLite to MariaDB
	/*
		//nur zur einmaligen Ausführung, um das Lernkosten-System zu initialisieren
		charGrp.POST("/initialize-learning-costs", InitializeLearningCosts)
		// Zur Überprüfung der Lernkosten-Daten
		charGrp.GET("/learning-costs-summary", gsmaster.GetLearningCostsSummaryHandler)
	*/

}
