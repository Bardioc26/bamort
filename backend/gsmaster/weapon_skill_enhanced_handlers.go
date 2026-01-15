package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetEnhancedMDWeaponSkills returns weapon skills with their source information and difficulty
func GetEnhancedMDWeaponSkills(c *gin.Context) {
	type EnhancedWeaponSkill struct {
		models.WeaponSkill
		SourceCode      string   `json:"source_code,omitempty"`
		Difficulty      string   `json:"difficulty,omitempty"`
		CategoryName    string   `json:"category_name,omitempty"`
		AllDifficulties []string `json:"all_difficulties,omitempty"`
	}

	type Response struct {
		WeaponSkills []EnhancedWeaponSkill    `json:"weaponskills"`
		Sources      []models.Source          `json:"sources"`
		Difficulties []models.SkillDifficulty `json:"difficulties"`
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

	var difficulties []models.SkillDifficulty
	if err := database.DB.Find(&difficulties).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve difficulties"})
		return
	}

	// Build source code map
	sourceMap := make(map[uint]string)
	for _, source := range sources {
		sourceMap[source.ID] = source.Code
	}

	// Query all weapon skill difficulties from learning_weaponskill_category_difficulties
	var categoryDifficulties []models.WeaponSkillCategoryDifficulty
	database.DB.Find(&categoryDifficulties)

	// Build difficulty maps
	difficultyMap := make(map[uint]string)
	categoryMap := make(map[uint]string)
	allDifficultiesMap := make(map[uint][]string)

	for _, cd := range categoryDifficulties {
		// Get the first (or preferred) difficulty for the weapon skill
		if _, exists := difficultyMap[cd.WeaponSkillID]; !exists {
			difficultyMap[cd.WeaponSkillID] = cd.SDifficulty
			categoryMap[cd.WeaponSkillID] = cd.SCategory
		}
		// Collect all difficulties for this weapon skill
		allDifficultiesMap[cd.WeaponSkillID] = append(allDifficultiesMap[cd.WeaponSkillID], cd.SDifficulty)
	}

	// Enhance weapon skills with source code and difficulty
	enhancedWeaponSkills := make([]EnhancedWeaponSkill, len(weaponSkills))
	for i, ws := range weaponSkills {
		enhancedWeaponSkills[i] = EnhancedWeaponSkill{
			WeaponSkill:     ws,
			SourceCode:      sourceMap[ws.SourceID],
			Difficulty:      difficultyMap[ws.ID],
			CategoryName:    categoryMap[ws.ID],
			AllDifficulties: allDifficultiesMap[ws.ID],
		}
	}

	c.JSON(http.StatusOK, Response{
		WeaponSkills: enhancedWeaponSkills,
		Sources:      sources,
		Difficulties: difficulties,
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

// UpdateEnhancedMDWeaponSkill updates a weapon skill and its difficulty relationship
func UpdateEnhancedMDWeaponSkill(c *gin.Context) {
	type UpdateRequest struct {
		models.WeaponSkill
		Difficulty string `json:"difficulty"`
		Category   string `json:"category"`
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

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the weapon skill basic fields
	weaponSkill.Name = req.Name
	weaponSkill.Beschreibung = req.Beschreibung
	weaponSkill.Initialwert = req.Initialwert
	weaponSkill.BasisWert = req.BasisWert
	weaponSkill.Bonuseigenschaft = req.Bonuseigenschaft
	weaponSkill.Improvable = req.Improvable
	weaponSkill.InnateSkill = req.InnateSkill
	weaponSkill.SourceID = req.SourceID
	weaponSkill.PageNumber = req.PageNumber
	weaponSkill.GameSystem = req.GameSystem
	weaponSkill.ID = uint(id)

	if err := database.DB.Save(&weaponSkill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update weapon skill"})
		return
	}

	// Update difficulty relationship if provided
	if req.Difficulty != "" {
		category := req.Category
		if category == "" {
			category = "Waffen" // Default category for weapon skills
		}

		// Look up SkillCategory ID
		var skillCategory models.SkillCategory
		if err := database.DB.Where("name = ?", category).First(&skillCategory).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found: " + category})
			return
		}

		// Look up SkillDifficulty ID
		var skillDifficulty models.SkillDifficulty
		if err := database.DB.Where("name = ?", req.Difficulty).First(&skillDifficulty).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Difficulty not found: " + req.Difficulty})
			return
		}

		// Find or create the WeaponSkillCategoryDifficulty entry
		var wscd models.WeaponSkillCategoryDifficulty
		err := database.DB.
			Where("weapon_skill_id = ? AND skill_category_id = ?", weaponSkill.ID, skillCategory.ID).
			First(&wscd).Error

		if err != nil {
			// Create new entry
			wscd = models.WeaponSkillCategoryDifficulty{
				WeaponSkillID:     weaponSkill.ID,
				SkillCategoryID:   skillCategory.ID,
				SkillDifficultyID: skillDifficulty.ID,
				SCategory:         category,
				SDifficulty:       req.Difficulty,
				LearnCost:         10, // Default learn cost for weapon skills
			}
			if err := database.DB.Create(&wscd).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create difficulty relationship: " + err.Error()})
				return
			}
		} else {
			// Update existing entry
			wscd.SkillDifficultyID = skillDifficulty.ID
			wscd.SDifficulty = req.Difficulty
			if err := database.DB.Save(&wscd).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update difficulty relationship: " + err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, weaponSkill)
}
