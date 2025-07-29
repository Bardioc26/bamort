package gsmaster

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitializeLearningCostsHandler HTTP-Handler zur Initialisierung des Lernkosten-Systems
func InitializeLearningCostsHandler(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "OK",
		"details": "skipped initialization, already done",
	})
	return

	err := InitializeLearningCostsSystem()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initialize learning costs system",
			"details": err.Error(),
		})
		return
	}

	// Validierung ausführen
	if err := ValidateLearningCostsData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Learning costs system initialized but validation failed",
			"details": err.Error(),
		})
		return
	}

	// Zusammenfassung erstellen
	summary, err := GetLearningCostsSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get summary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Learning costs system initialized successfully",
		"summary": summary,
	})
}

// GetLearningCostsSummaryHandler HTTP-Handler für die Zusammenfassung
func GetLearningCostsSummaryHandler(c *gin.Context) {
	summary, err := GetLearningCostsSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get learning costs summary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
	})
}
