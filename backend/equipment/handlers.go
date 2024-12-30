package equipment

import (
	"bamort/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Endpoints for Managing Ausruestung
1. Create Ausruestung

Allows users to add new equipment items for a specific character.
*/

func CreateAusruestung(c *gin.Context) {
	var ausruestung Ausruestung
	if err := c.ShouldBindJSON(&ausruestung); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&ausruestung).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Ausruestung"})
		return
	}

	c.JSON(http.StatusCreated, ausruestung)
}

func GetAusruestung(c *gin.Context) {
	characterID := c.Param("character_id")

	var ausruestung []Ausruestung
	if err := database.DB.Where("character_id = ?", characterID).Find(&ausruestung).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Ausruestung"})
		return
	}

	c.JSON(http.StatusOK, ausruestung)
}

func UpdateAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")
	var ausruestung Ausruestung

	if err := database.DB.First(&ausruestung, ausruestungID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ausruestung not found"})
		return
	}

	if err := c.ShouldBindJSON(&ausruestung); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&ausruestung).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Ausruestung"})
		return
	}

	c.JSON(http.StatusOK, ausruestung)
}

func DeleteAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")
	if err := database.DB.Delete(&Ausruestung{}, ausruestungID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Ausruestung"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ausruestung deleted successfully"})
}
