package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SkillWithCategories represents a skill with all its categories and difficulties
type SkillWithCategories struct {
	models.Skill
	Categories   []SkillCategoryInfo `json:"categories"`
	Difficulties []string            `json:"difficulties"`
}

// SkillCategoryInfo contains category details for a skill
type SkillCategoryInfo struct {
	CategoryID     uint   `json:"category_id"`
	CategoryName   string `json:"category_name"`
	DifficultyID   uint   `json:"difficulty_id"`
	DifficultyName string `json:"difficulty_name"`
	LearnCost      int    `json:"learn_cost"`
}

// GetSkillWithCategories retrieves a skill with all its category-difficulty relationships
func GetSkillWithCategories(skillID uint) (*SkillWithCategories, error) {
	var skill models.Skill
	if err := database.DB.First(&skill, skillID).Error; err != nil {
		return nil, err
	}

	// Get all category-difficulty relationships
	var scds []models.SkillCategoryDifficulty
	err := database.DB.Preload("SkillCategory").Preload("SkillDifficulty").
		Where("skill_id = ?", skillID).Find(&scds).Error
	if err != nil {
		return nil, err
	}

	result := &SkillWithCategories{
		Skill:        skill,
		Categories:   make([]SkillCategoryInfo, len(scds)),
		Difficulties: make([]string, len(scds)),
	}

	for i, scd := range scds {
		result.Categories[i] = SkillCategoryInfo{
			CategoryID:     scd.SkillCategoryID,
			CategoryName:   scd.SkillCategory.Name,
			DifficultyID:   scd.SkillDifficultyID,
			DifficultyName: scd.SkillDifficulty.Name,
			LearnCost:      scd.LearnCost,
		}
		result.Difficulties[i] = scd.SkillDifficulty.Name
	}

	return result, nil
}

// GetAllSkillsWithCategories retrieves all skills with their categories
func GetAllSkillsWithCategories() ([]SkillWithCategories, error) {
	var skills []models.Skill
	if err := database.DB.Find(&skills).Error; err != nil {
		return nil, err
	}

	result := make([]SkillWithCategories, len(skills))
	for i, skill := range skills {
		skillWithCats, err := GetSkillWithCategories(skill.ID)
		if err != nil {
			return nil, err
		}
		result[i] = *skillWithCats
	}

	return result, nil
}

// SkillUpdateRequest represents the request to update a skill with categories
type SkillUpdateRequest struct {
	models.Skill
	CategoryDifficulties []CategoryDifficultyPair `json:"category_difficulties"`
}

// CategoryDifficultyPair represents a category-difficulty mapping
type CategoryDifficultyPair struct {
	CategoryID   uint `json:"category_id"`
	DifficultyID uint `json:"difficulty_id"`
	LearnCost    int  `json:"learn_cost,omitempty"`
}

// CreateSkillWithCategories creates a new skill with category-difficulty relationships
func CreateSkillWithCategories(req SkillUpdateRequest) (uint, error) {
	// Validate required fields
	if req.Skill.Name == "" {
		return 0, fmt.Errorf("skill name is required")
	}

	var skillID uint

	// Start transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Create skill
		if err := tx.Create(&req.Skill).Error; err != nil {
			return err
		}

		skillID = req.Skill.ID

		// Create category-difficulty relationships
		for _, cd := range req.CategoryDifficulties {
			// Get category and difficulty names for denormalized fields
			var category models.SkillCategory
			if err := tx.First(&category, cd.CategoryID).Error; err != nil {
				return fmt.Errorf("category not found: %w", err)
			}

			var difficulty models.SkillDifficulty
			if err := tx.First(&difficulty, cd.DifficultyID).Error; err != nil {
				return fmt.Errorf("difficulty not found: %w", err)
			}

			learnCost := cd.LearnCost
			if learnCost == 0 {
				// Use default based on difficulty
				learnCost = getDefaultLearnCost(difficulty.Name)
			}

			scd := models.SkillCategoryDifficulty{
				SkillID:           skillID,
				SkillCategoryID:   cd.CategoryID,
				SkillDifficultyID: cd.DifficultyID,
				LearnCost:         learnCost,
				SCategory:         category.Name,
				SDifficulty:       difficulty.Name,
			}

			if err := tx.Create(&scd).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return skillID, nil
}

// UpdateSkillWithCategories updates a skill and its category-difficulty relationships
func UpdateSkillWithCategories(skillID uint, req SkillUpdateRequest) error {
	// Start transaction
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// Update skill basic info - use Select to explicitly include boolean fields
		// This ensures false values are also updated (GORM skips zero values by default in Updates)
		if err := tx.Model(&models.Skill{}).Where("id = ?", skillID).
			Select("name", "beschreibung", "game_system", "initialwert", "basis_wert", 
				"bonuseigenschaft", "improvable", "innate_skill", "source_id", "page_number").
			Updates(req.Skill).Error; err != nil {
			return err
		}

		// Delete existing category-difficulty relationships
		if err := tx.Where("skill_id = ?", skillID).Delete(&models.SkillCategoryDifficulty{}).Error; err != nil {
			return err
		}

		// Create new relationships
		for _, cd := range req.CategoryDifficulties {
			// Get category and difficulty names for denormalized fields
			var category models.SkillCategory
			if err := tx.First(&category, cd.CategoryID).Error; err != nil {
				return fmt.Errorf("category not found: %w", err)
			}

			var difficulty models.SkillDifficulty
			if err := tx.First(&difficulty, cd.DifficultyID).Error; err != nil {
				return fmt.Errorf("difficulty not found: %w", err)
			}

			learnCost := cd.LearnCost
			if learnCost == 0 {
				// Use default based on difficulty
				learnCost = getDefaultLearnCost(difficulty.Name)
			}

			scd := models.SkillCategoryDifficulty{
				SkillID:           skillID,
				SkillCategoryID:   cd.CategoryID,
				SkillDifficultyID: cd.DifficultyID,
				LearnCost:         learnCost,
				SCategory:         category.Name,
				SDifficulty:       difficulty.Name,
			}

			if err := tx.Create(&scd).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// getDefaultLearnCost returns default LE cost based on difficulty
func getDefaultLearnCost(difficulty string) int {
	switch strings.ToLower(difficulty) {
	case "leicht", "easy":
		return 5
	case "normal", "standard":
		return 10
	case "schwer", "hard":
		return 20
	case "sehr schwer", "very hard":
		return 30
	default:
		return 10
	}
}

// ===== Handler Functions =====

// GetEnhancedMDSkills returns skills with their full category/difficulty information
func GetEnhancedMDSkills(c *gin.Context) {
	skills, err := GetAllSkillsWithCategories()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve skills: "+err.Error())
		return
	}

	// Also get learning sources and difficulties for the dropdowns
	var sources []models.Source
	database.DB.Where("is_active = ?", true).Find(&sources)

	var categories []models.SkillCategory
	database.DB.Find(&categories)

	var difficulties []models.SkillDifficulty
	database.DB.Find(&difficulties)

	c.JSON(http.StatusOK, gin.H{
		"skills":       skills,
		"sources":      sources,
		"categories":   categories,
		"difficulties": difficulties,
	})
}

// GetEnhancedMDSkill returns a single skill with category/difficulty information
func GetEnhancedMDSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	skill, err := GetSkillWithCategories(uint(id))
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Skill not found")
		return
	}

	c.JSON(http.StatusOK, skill)
}

// UpdateEnhancedMDSkill updates a skill with its categories
func UpdateEnhancedMDSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req SkillUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Ensure the ID matches
	req.Skill.ID = uint(id)

	if err := UpdateSkillWithCategories(uint(id), req); err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to update skill: "+err.Error())
		return
	}

	// Return updated skill
	skill, err := GetSkillWithCategories(uint(id))
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve updated skill")
		return
	}

	c.JSON(http.StatusOK, skill)
}

// CreateEnhancedMDSkill creates a new skill with categories
func CreateEnhancedMDSkill(c *gin.Context) {
	var req SkillUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	// Create the skill
	skillID, err := CreateSkillWithCategories(req)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create skill: "+err.Error())
		return
	}

	// Return created skill
	skill, err := GetSkillWithCategories(skillID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to retrieve created skill")
		return
	}

	c.JSON(http.StatusCreated, skill)
}
