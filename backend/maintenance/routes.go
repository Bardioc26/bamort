package maintenance

import (
	"bamort/gsmaster"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	charGrp := r.Group("/maintenance")
	charGrp.GET("/setupcheck", SetupCheck)
	charGrp.GET("/mktestdata", MakeTestdataFromLive)
	//nur zur einmaligen Ausführung, um das Lernkosten-System zu initialisieren
	charGrp.POST("/initialize-learning-costs", InitializeLearningCosts)
	// Zur Überprüfung der Lernkosten-Daten
	charGrp.GET("/learning-costs-summary", gsmaster.GetLearningCostsSummaryHandler)

}
