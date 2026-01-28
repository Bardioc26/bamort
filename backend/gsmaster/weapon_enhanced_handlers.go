package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetEnhancedMDWeapons returns weapons with their source information
func GetEnhancedMDWeapons(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	type EnhancedWeapon struct {
		models.Weapon
		SourceCode string `json:"source_code,omitempty"`
	}

	type Response struct {
		Weapons []EnhancedWeapon `json:"weapons"`
		Sources []models.Source  `json:"sources"`
	}

	var weapons []models.Weapon
	if err := database.DB.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID).Find(&weapons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve weapons"})
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

	// Enhance weapons with source code
	enhancedWeapons := make([]EnhancedWeapon, len(weapons))
	for i, w := range weapons {
		enhancedWeapons[i] = EnhancedWeapon{
			Weapon:     w,
			SourceCode: sourceMap[w.SourceID],
		}
	}

	c.JSON(http.StatusOK, Response{
		Weapons: enhancedWeapons,
		Sources: sources,
	})
}

// GetEnhancedMDWeapon returns a single weapon with source information
func GetEnhancedMDWeapon(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	type EnhancedWeapon struct {
		models.Weapon
		SourceCode string `json:"source_code,omitempty"`
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var weapon models.Weapon
	if err := database.DB.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID).First(&weapon, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Weapon not found"})
		return
	}

	var source models.Source
	var sourceCode string
	if weapon.SourceID > 0 {
		if err := database.DB.First(&source, weapon.SourceID).Error; err == nil {
			sourceCode = source.Code
		}
	}

	enhanced := EnhancedWeapon{
		Weapon:     weapon,
		SourceCode: sourceCode,
	}

	c.JSON(http.StatusOK, enhanced)
}

// UpdateEnhancedMDWeapon updates a weapon
func UpdateEnhancedMDWeapon(c *gin.Context) {
	gs, ok := resolveGameSystem(c)
	if !ok {
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var weapon models.Weapon
	if err := database.DB.Where("game_system=? OR game_system_id=?", gs.Name, gs.ID).First(&weapon, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Weapon not found"})
		return
	}

	// Bind the request to weapon
	if err := c.ShouldBindJSON(&weapon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID stays the same
	weapon.ID = uint(id)
	weapon.GameSystem = gs.Name
	weapon.GameSystemId = gs.ID

	// Update the weapon
	if err := weapon.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update weapon"})
		return
	}

	c.JSON(http.StatusOK, weapon)
}
