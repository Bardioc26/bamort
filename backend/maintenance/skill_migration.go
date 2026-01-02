package maintenance

import (
	"bamort/logger"
	"bamort/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// MigrateSkillCategoriesToRelations migrates existing Category/Difficulty fields to the relational model
func MigrateSkillCategoriesToRelations(db *gorm.DB) error {
	logger.Info("Starting migration of skill categories to relational model...")

	// Get all skills with existing category/difficulty data
	var skills []models.Skill
	if err := db.Where("category IS NOT NULL AND category != ''").Find(&skills).Error; err != nil {
		return fmt.Errorf("failed to fetch skills: %w", err)
	}

	logger.Info("Found %d skills to migrate", len(skills))

	migrated := 0
	skipped := 0
	errors := 0

	for _, skill := range skills {
		if err := migrateSkillCategoryDifficulty(db, &skill); err != nil {
			logger.Error("Failed to migrate skill %d (%s): %s", skill.ID, skill.Name, err.Error())
			errors++
			continue
		}
		migrated++
	}

	logger.Info("Migration completed: %d migrated, %d skipped, %d errors", migrated, skipped, errors)
	return nil
}

func migrateSkillCategoryDifficulty(db *gorm.DB, skill *models.Skill) error {
	// Check if already migrated
	var count int64
	if err := db.Model(&models.SkillCategoryDifficulty{}).Where("skill_id = ?", skill.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing migration: %w", err)
	}

	if count > 0 {
		logger.Debug("Skill %d (%s) already has category difficulties, skipping", skill.ID, skill.Name)
		return nil
	}

	// Find or create skill category
	var skillCategory models.SkillCategory
	categoryName := strings.TrimSpace(skill.Category)
	if categoryName == "" {
		categoryName = "Alltag" // Default category
	}

	err := db.Where("name = ? AND game_system = ?", categoryName, skill.GameSystem).First(&skillCategory).Error
	if err == gorm.ErrRecordNotFound {
		// Create new category if it doesn't exist
		// First, find a valid source (we'll use the skill's source or a default)
		var source models.Source
		if skill.SourceID > 0 {
			db.First(&source, skill.SourceID)
		} else {
			// Use default source (KOD)
			db.Where("code = ?", "KOD").First(&source)
		}

		if source.ID == 0 {
			return fmt.Errorf("no valid source found for category creation")
		}

		skillCategory = models.SkillCategory{
			Name:       categoryName,
			GameSystem: skill.GameSystem,
			SourceID:   source.ID,
			Quelle:     source.Code,
		}
		if err := db.Create(&skillCategory).Error; err != nil {
			return fmt.Errorf("failed to create skill category: %w", err)
		}
		logger.Debug("Created new skill category: %s", categoryName)
	} else if err != nil {
		return fmt.Errorf("failed to find skill category: %w", err)
	}

	// Find or create skill difficulty
	var skillDifficulty models.SkillDifficulty
	difficultyName := strings.TrimSpace(skill.Difficulty)
	if difficultyName == "" {
		difficultyName = "normal" // Default difficulty
	}

	err = db.Where("name = ? AND game_system = ?", difficultyName, skill.GameSystem).First(&skillDifficulty).Error
	if err == gorm.ErrRecordNotFound {
		skillDifficulty = models.SkillDifficulty{
			Name:       difficultyName,
			GameSystem: skill.GameSystem,
		}
		if err := db.Create(&skillDifficulty).Error; err != nil {
			return fmt.Errorf("failed to create skill difficulty: %w", err)
		}
		logger.Debug("Created new skill difficulty: %s", difficultyName)
	} else if err != nil {
		return fmt.Errorf("failed to find skill difficulty: %w", err)
	}

	// Create SkillCategoryDifficulty relationship
	scd := models.SkillCategoryDifficulty{
		SkillID:           skill.ID,
		SkillCategoryID:   skillCategory.ID,
		SkillDifficultyID: skillDifficulty.ID,
		LearnCost:         getDefaultLearnCost(difficultyName),
		SDifficulty:       difficultyName,
		SCategory:         categoryName,
	}

	if err := db.Create(&scd).Error; err != nil {
		return fmt.Errorf("failed to create skill category difficulty: %w", err)
	}

	logger.Debug("Migrated skill %d (%s): category=%s, difficulty=%s", skill.ID, skill.Name, categoryName, difficultyName)
	return nil
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
