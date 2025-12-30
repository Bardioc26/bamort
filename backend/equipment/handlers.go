package equipment

import (
	"bamort/database"
	"bamort/models"
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
	var ausruestung models.EqAusruestung
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

	var ausruestung []models.EqAusruestung
	if err := database.DB.Where("character_id = ?", characterID).Find(&ausruestung).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve Ausruestung")
		return
	}

	c.JSON(http.StatusOK, ausruestung)
}

func UpdateAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")
	var ausruestung models.EqAusruestung

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
	if err := database.DB.Delete(&models.EqAusruestung{}, ausruestungID).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete Ausruestung")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ausruestung deleted successfully"})
}

/*
Endpoints for Managing Weapons (Waffen)
*/

func CreateWaffe(c *gin.Context) {
	var waffe models.EqWaffe
	if err := c.ShouldBindJSON(&waffe); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Create(&waffe).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create Waffe")
		return
	}

	c.JSON(http.StatusCreated, waffe)
}

func ListWaffen(c *gin.Context) {
	characterID := c.Param("character_id")

	var waffen []models.EqWaffe
	if err := database.DB.Where("character_id = ?", characterID).Find(&waffen).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve Waffen")
		return
	}

	c.JSON(http.StatusOK, waffen)
}

func UpdateWaffe(c *gin.Context) {
	waffeID := c.Param("waffe_id")
	var waffe models.EqWaffe

	if err := database.DB.First(&waffe, waffeID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Waffe not found")
		return
	}

	if err := c.ShouldBindJSON(&waffe); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := database.DB.Save(&waffe).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update Waffe")
		return
	}

	c.JSON(http.StatusOK, waffe)
}

func DeleteWaffe(c *gin.Context) {
	waffeID := c.Param("waffe_id")
	if err := database.DB.Delete(&models.EqWaffe{}, waffeID).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete Waffe")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Waffe deleted successfully"})
}
