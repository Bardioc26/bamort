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
func GetEquipmentWithCategories(equipmentID uint) (*EquipmentWithCategories, error) {
	var equipment models.Equipment
	if err := database.DB.First(&equipment, equipmentID).Error; err != nil {
		return nil, err
	}

	result := &EquipmentWithCategories{
		Equipment: equipment,
	}

	return result, nil
}

// GetAllEquipmentWithCategories retrieves all equipment
func GetAllEquipmentWithCategories() ([]EquipmentWithCategories, error) {
	var equipments []models.Equipment
	if err := database.DB.Find(&equipments).Error; err != nil {
		return nil, err
	}

	result := make([]EquipmentWithCategories, len(equipments))
	for i, equipment := range equipments {
		equipmentWithCats, err := GetEquipmentWithCategories(equipment.ID)
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
func UpdateEquipmentWithCategories(equipmentID uint, req EquipmentUpdateRequest) error {
	// Start transaction
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Update equipment info
		if err := tx.Model(&models.Equipment{}).Where("id = ?", equipmentID).Updates(req.Equipment).Error; err != nil {
			return err
		}

		return nil
	})
}

// ===== Handler Functions =====

// GetEnhancedMDEquipment returns equipment with enhanced information
func GetEnhancedMDEquipment(c *gin.Context) {
	equipments, err := GetAllEquipmentWithCategories()
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
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	equipment, err := GetEquipmentWithCategories(uint(id))
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Equipment not found")
		return
	}

	c.JSON(http.StatusOK, equipment)
}

// UpdateEnhancedMDEquipmentItem updates equipment
func UpdateEnhancedMDEquipmentItem(c *gin.Context) {
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

	// Ensure the ID matches
	req.Equipment.ID = uint(id)

	if err := UpdateEquipmentWithCategories(uint(id), req); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update equipment: "+err.Error())
		return
	}

	// Return updated equipment
	equipment, err := GetEquipmentWithCategories(uint(id))
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve updated equipment")
		return
	}

	c.JSON(http.StatusOK, equipment)
}
