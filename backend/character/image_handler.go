package character

import (
	"bamort/database"
	"bamort/logger"
	"bamort/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ImageUpdateRequest struct {
	Image string `json:"image" binding:"required"`
}

func UpdateCharacterImage(c *gin.Context) {
	id := c.Param("id")
	logger.Debug("UpdateCharacterImage called for character ID: %s", id)

	var character models.Char
	err := character.FirstID(id)
	if err != nil {
		logger.Error("Character not found: %s", err.Error())
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Check ownership
	if !checkCharacterOwnership(c, &character) {
		return
	}

	var request ImageUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("Invalid request data: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Invalid request data")
		return
	}

	character.Image = request.Image

	if err := database.DB.Save(&character).Error; err != nil {
		logger.Error("Failed to update character image: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to update character image")
		return
	}

	logger.Info("Character image updated successfully for ID: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Image updated successfully", "image": character.Image})
}
