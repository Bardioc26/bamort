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

// checkEquipmentOwnership verifies that the logged-in user owns the equipment's character
func checkEquipmentOwnership(c *gin.Context, characterID uint) bool {
	userID := c.GetUint("userID")
	var character models.Char
	if err := database.DB.Select("id", "user_id").First(&character, characterID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return false
	}
	if character.UserID != userID {
		respondWithError(c, http.StatusForbidden, "You are not authorized to modify this character's equipment")
		return false
	}
	return true
}

// CreateAusruestung godoc
// @Summary Create equipment
// @Description Creates a new equipment item for a character (owner only)
// @Tags Equipment
// @Accept json
// @Produce json
// @Param equipment body models.EqAusruestung true "Equipment data"
// @Success 201 {object} models.EqAusruestung "Created equipment"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Security BearerAuth
// @Router /api/equipment [post]
func CreateAusruestung(c *gin.Context) {
	var ausruestung models.EqAusruestung
	if err := c.ShouldBindJSON(&ausruestung); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check ownership
	if !checkEquipmentOwnership(c, ausruestung.CharacterID) {
		return
	}

	if err := database.DB.Create(&ausruestung).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create Ausruestung")
		return
	}

	c.JSON(http.StatusCreated, ausruestung)
}

// ListAusruestung godoc
// @Summary List character equipment
// @Description Returns all equipment items for a specific character
// @Tags Equipment
// @Produce json
// @Param character_id path int true "Character ID"
// @Success 200 {array} models.EqAusruestung "List of equipment"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Security BearerAuth
// @Router /api/equipment/character/{character_id} [get]
func ListAusruestung(c *gin.Context) {
	characterID := c.Param("character_id")

	var ausruestung []models.EqAusruestung
	if err := database.DB.Where("character_id = ?", characterID).Find(&ausruestung).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve Ausruestung")
		return
	}

	c.JSON(http.StatusOK, ausruestung)
}

// UpdateAusruestung godoc
// @Summary Update equipment
// @Description Updates an existing equipment item (owner only)
// @Tags Equipment
// @Accept json
// @Produce json
// @Param ausruestung_id path int true "Equipment ID"
// @Param equipment body models.EqAusruestung true "Updated equipment data"
// @Success 200 {object} models.EqAusruestung "Updated equipment"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Security BearerAuth
// @Router /api/equipment/{ausruestung_id} [put]
func UpdateAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")
	var ausruestung models.EqAusruestung

	if err := database.DB.First(&ausruestung, ausruestungID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Ausruestung not found")
		return
	}

	// Check ownership
	if !checkEquipmentOwnership(c, ausruestung.CharacterID) {
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

// DeleteAusruestung godoc
// @Summary Delete equipment
// @Description Deletes an equipment item (owner only)
// @Tags Equipment
// @Produce json
// @Param ausruestung_id path int true "Equipment ID"
// @Success 200 {object} map[string]string "Equipment deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Failure 404 {object} map[string]string "Equipment not found"
// @Security BearerAuth
// @Router /api/equipment/{ausruestung_id} [delete]
func DeleteAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")

	var ausruestung models.EqAusruestung
	if err := database.DB.First(&ausruestung, ausruestungID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Ausruestung not found")
		return
	}

	// Check ownership
	if !checkEquipmentOwnership(c, ausruestung.CharacterID) {
		return
	}

	if err := database.DB.Delete(&ausruestung).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete Ausruestung")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ausruestung deleted successfully"})
}

/*
Endpoints for Managing Weapons (Waffen)
*/

// CreateWaffe godoc
// @Summary Create weapon
// @Description Creates a new weapon for a character (owner only)
// @Tags Equipment
// @Accept json
// @Produce json
// @Param weapon body models.EqWaffe true "Weapon data"
// @Success 201 {object} models.EqWaffe "Created weapon"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Security BearerAuth
// @Router /api/weapons [post]
func CreateWaffe(c *gin.Context) {
	var waffe models.EqWaffe
	if err := c.ShouldBindJSON(&waffe); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check ownership
	if !checkEquipmentOwnership(c, waffe.CharacterID) {
		return
	}

	if err := database.DB.Create(&waffe).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create Waffe")
		return
	}

	c.JSON(http.StatusCreated, waffe)
}

// ListWaffen godoc
// @Summary List character weapons
// @Description Returns all weapons for a specific character
// @Tags Equipment
// @Produce json
// @Param character_id path int true "Character ID"
// @Success 200 {array} models.EqWaffe "List of weapons"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Security BearerAuth
// @Router /api/weapons/character/{character_id} [get]
func ListWaffen(c *gin.Context) {
	characterID := c.Param("character_id")

	var waffen []models.EqWaffe
	if err := database.DB.Where("character_id = ?", characterID).Find(&waffen).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve Waffen")
		return
	}

	c.JSON(http.StatusOK, waffen)
}

// UpdateWaffe godoc
// @Summary Update weapon
// @Description Updates an existing weapon (owner only)
// @Tags Equipment
// @Accept json
// @Produce json
// @Param waffe_id path int true "Weapon ID"
// @Param weapon body models.EqWaffe true "Updated weapon data"
// @Success 200 {object} models.EqWaffe "Updated weapon"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Security BearerAuth
// @Router /api/weapons/{waffe_id} [put]
func UpdateWaffe(c *gin.Context) {
	waffeID := c.Param("waffe_id")
	var waffe models.EqWaffe

	if err := database.DB.First(&waffe, waffeID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Waffe not found")
		return
	}

	// Check ownership
	if !checkEquipmentOwnership(c, waffe.CharacterID) {
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

// DeleteWaffe godoc
// @Summary Delete weapon
// @Description Deletes a weapon (owner only)
// @Tags Equipment
// @Produce json
// @Param waffe_id path int true "Weapon ID"
// @Success 200 {object} map[string]string "Weapon deleted successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied - owner only"
// @Failure 404 {object} map[string]string "Weapon not found"
// @Security BearerAuth
// @Router /api/weapons/{waffe_id} [delete]
func DeleteWaffe(c *gin.Context) {
	waffeID := c.Param("waffe_id")

	var waffe models.EqWaffe
	if err := database.DB.First(&waffe, waffeID).Error; err != nil {
		respondWithError(c, http.StatusNotFound, "Waffe not found")
		return
	}

	// Check ownership
	if !checkEquipmentOwnership(c, waffe.CharacterID) {
		return
	}

	if err := database.DB.Delete(&waffe).Error; err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to delete Waffe")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Waffe deleted successfully"})
}
