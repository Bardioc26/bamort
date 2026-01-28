package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// EquipmentWithCategories represents equipment with enhanced information
type EquipmentWithCategories struct {
	models.Equipment
}

// GetEquipmentWithCategories retrieves equipment with all its information
func GetEquipmentWithCategories(equipmentID uint, gs *models.GameSystem) (*EquipmentWithCategories, error) {
	var equipment models.Equipment
	query := database.DB
	if gs != nil {
		query = query.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID)
	}

	if err := query.First(&equipment, equipmentID).Error; err != nil {
		return nil, err
	}

	result := &EquipmentWithCategories{
		Equipment: equipment,
	}

	return result, nil
}

// GetAllEquipmentWithCategories retrieves all equipment
func GetAllEquipmentWithCategories(gs *models.GameSystem) ([]EquipmentWithCategories, error) {
	var equipments []models.Equipment
	query := database.DB
	if gs != nil {
		query = query.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID)
	}

	if err := query.Find(&equipments).Error; err != nil {
		return nil, err
	}

	result := make([]EquipmentWithCategories, len(equipments))
	for i, equipment := range equipments {
		equipmentWithCats, err := GetEquipmentWithCategories(equipment.ID, gs)
		if err != nil {
			return nil, err
		}
		result[i] = *equipmentWithCats
	}

	return result, nil
}

// EquipmentUpdateRequest represents the request to update equipment
type EquipmentUpdateRequest struct {
	models.Equipment
}

// UpdateEquipmentWithCategories updates equipment
func UpdateEquipmentWithCategories(equipmentID uint, req EquipmentUpdateRequest, gs *models.GameSystem) error {
	// Start transaction
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Update equipment info
		query := tx.Model(&models.Equipment{}).Where("id = ?", equipmentID)
		if gs != nil {
			req.Equipment.GameSystem = gs.Name
			req.Equipment.GameSystemId = gs.ID
			query = query.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID)
		}

		if err := query.Updates(req.Equipment).Error; err != nil {
			return err
		}

		return nil
	})
}

// ===== Handler Functions =====

// GetEnhancedMDEquipment returns equipment with enhanced information
func GetEnhancedMDEquipment(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	equipments, err := GetAllEquipmentWithCategories(gs)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve equipment: "+err.Error())
		return
	}

	// Also get learning sources for the dropdowns
	var sources []models.Source
	database.DB.Where("is_active = ?", true).Find(&sources)

	c.JSON(http.StatusOK, gin.H{
		"equipment": equipments,
		"sources":   sources,
	})
}

// GetEnhancedMDEquipmentItem returns a single equipment with enhanced information
func GetEnhancedMDEquipmentItem(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	equipment, err := GetEquipmentWithCategories(uint(id), gs)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Equipment not found")
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// UpdateEnhancedMDEquipmentItem updates equipment
func UpdateEnhancedMDEquipmentItem(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req EquipmentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	req.Equipment.GameSystem = gs.Name
	req.Equipment.GameSystemId = gs.ID

	// Ensure the ID matches
	req.Equipment.ID = uint(id)

	if err := UpdateEquipmentWithCategories(uint(id), req, gs); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update equipment: "+err.Error())
		return
	}

	// Return updated equipment
	equipment, err := GetEquipmentWithCategories(uint(id), gs)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve updated equipment")
		return
	}

	c.JSON(http.StatusOK, equipment)
}
