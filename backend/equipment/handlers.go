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

func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func CreateAusruestung(c *gin.Context) {
	var ausruestung Ausruestung
	if err := c.ShouldBindJSON(&ausruestung); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Create(&ausruestung).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create Ausruestung")
		return
	}

	c.JSON(http.StatusCreated, ausruestung)
}

func ListAusruestung(c *gin.Context) {
	characterID := c.Param("character_id")

	var ausruestung []Ausruestung
	if err := database.DB.Where("character_id = ?", characterID).Find(&ausruestung).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve Ausruestung")
		return
	}

	c.JSON(http.StatusOK, ausruestung)
}

func UpdateAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")
	var ausruestung Ausruestung

	if err := database.DB.First(&ausruestung, ausruestungID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Ausruestung not found")
		return
	}

	if err := c.ShouldBindJSON(&ausruestung); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Save(&ausruestung).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update Ausruestung")
		return
	}

	c.JSON(http.StatusOK, ausruestung)
}

func DeleteAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")
	if err := database.DB.Delete(&Ausruestung{}, ausruestungID).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete Ausruestung")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ausruestung deleted successfully"})
}
