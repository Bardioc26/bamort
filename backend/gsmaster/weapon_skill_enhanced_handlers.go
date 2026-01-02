package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetEnhancedMDWeaponSkills returns weapon skills with their source information
func GetEnhancedMDWeaponSkills(c *gin.Context) {
	type EnhancedWeaponSkill struct {
		models.WeaponSkill
		SourceCode string `json:"source_code,omitempty"`
	}

	type Response struct {
		WeaponSkills []EnhancedWeaponSkill `json:"weaponskills"`
		Sources      []models.Source       `json:"sources"`
	}

	var weaponSkills []models.WeaponSkill
	if err := database.DB.Find(&weaponSkills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve weapon skills"})
		return
	}

	var sources []models.Source
	if err := database.DB.Find(&sources).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sources"})
		return
	}

	// Build source code map
	sourceMap := make(map[uint]string)
	for _, source := range sources {
		sourceMap[source.ID] = source.Code
	}

	// Enhance weapon skills with source code
	enhancedWeaponSkills := make([]EnhancedWeaponSkill, len(weaponSkills))
	for i, ws := range weaponSkills {
		enhancedWeaponSkills[i] = EnhancedWeaponSkill{
			WeaponSkill: ws,
			SourceCode:  sourceMap[ws.SourceID],
		}
	}

	c.JSON(http.StatusOK, Response{
		WeaponSkills: enhancedWeaponSkills,
		Sources:      sources,
	})
}

// GetEnhancedMDWeaponSkill returns a single weapon skill with source information
func GetEnhancedMDWeaponSkill(c *gin.Context) {
	type EnhancedWeaponSkill struct {
		models.WeaponSkill
		SourceCode string `json:"source_code,omitempty"`
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var weaponSkill models.WeaponSkill
	if err := database.DB.First(&weaponSkill, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Weapon skill not found"})
		return
	}

	var source models.Source
	var sourceCode string
	if weaponSkill.SourceID > 0 {
		if err := database.DB.First(&source, weaponSkill.SourceID).Error; err == nil {
			sourceCode = source.Code
		}
	}

	enhanced := EnhancedWeaponSkill{
		WeaponSkill: weaponSkill,
		SourceCode:  sourceCode,
	}

	c.JSON(http.StatusOK, enhanced)
}

// UpdateEnhancedMDWeaponSkill updates a weapon skill
func UpdateEnhancedMDWeaponSkill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var weaponSkill models.WeaponSkill
	if err := database.DB.First(&weaponSkill, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Weapon skill not found"})
		return
	}

	// Bind the request to weapon skill
	if err := c.ShouldBindJSON(&weaponSkill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID stays the same
	weaponSkill.ID = uint(id)

	// Update the weapon skill
	if err := database.DB.Save(&weaponSkill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update weapon skill"})
		return
	}

	c.JSON(http.StatusOK, weaponSkill)
}
