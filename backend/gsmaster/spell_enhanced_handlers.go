package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SpellWithCategories represents a spell with enhanced information
type SpellWithCategories struct {
	models.Spell
}

// GetSpellWithCategories retrieves a spell with all its information
func GetSpellWithCategories(spellID uint) (*SpellWithCategories, error) {
	var spell models.Spell
	if err := database.DB.First(&spell, spellID).Error; err != nil {
		return nil, err
	}

	result := &SpellWithCategories{
		Spell: spell,
	}

	return result, nil
}

// GetAllSpellsWithCategories retrieves all spells
func GetAllSpellsWithCategories() ([]SpellWithCategories, error) {
	var spells []models.Spell
	if err := database.DB.Find(&spells).Error; err != nil {
		return nil, err
	}

	result := make([]SpellWithCategories, len(spells))
	for i, spell := range spells {
		spellWithCats, err := GetSpellWithCategories(spell.ID)
		if err != nil {
			return nil, err
		}
		result[i] = *spellWithCats
	}

	return result, nil
}

// SpellUpdateRequest represents the request to update a spell
type SpellUpdateRequest struct {
	models.Spell
}

// UpdateSpellWithCategories updates a spell
func UpdateSpellWithCategories(spellID uint, req SpellUpdateRequest) error {
	// Start transaction
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Update spell info
		if err := tx.Model(&models.Spell{}).Where("id = ?", spellID).Updates(req.Spell).Error; err != nil {
			return err
		}

		return nil
	})
}

// ===== Handler Functions =====

// GetEnhancedMDSpells returns spells with enhanced information
func GetEnhancedMDSpells(c *gin.Context) {
	spells, err := GetAllSpellsWithCategories()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve spells: "+err.Error())
		return
	}

	// Also get learning sources and categories for the dropdowns
	var sources []models.Source
	database.DB.Where("is_active = ?", true).Find(&sources)

	// Get spell categories
	var spell models.Spell
	categories, err := spell.GetSpellCategories()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve spell categories: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"spells":     spells,
		"sources":    sources,
		"categories": categories,
	})
}

// GetEnhancedMDSpell returns a single spell with enhanced information
func GetEnhancedMDSpell(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	spell, err := GetSpellWithCategories(uint(id))
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Spell not found")
		return
	}

	c.JSON(http.StatusOK, spell)
}

// UpdateEnhancedMDSpell updates a spell
func UpdateEnhancedMDSpell(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req SpellUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Ensure the ID matches
	req.Spell.ID = uint(id)

	if err := UpdateSpellWithCategories(uint(id), req); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update spell: "+err.Error())
		return
	}

	// Return updated spell
	spell, err := GetSpellWithCategories(uint(id))
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve updated spell")
		return
	}

	c.JSON(http.StatusOK, spell)
}
