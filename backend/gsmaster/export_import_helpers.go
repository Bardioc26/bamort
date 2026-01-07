package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"fmt"

	"gorm.io/gorm"
)

// LookupMap builders - reusable functions to build ID<->Code/Name maps

// buildSourceMap creates a map from source ID to source code
func buildSourceMap() map[uint]string {
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[uint]string)
	for _, s := range sources {
		sourceMap[s.ID] = s.Code
	}
	return sourceMap
}

// buildSourceMapReverse creates a map from source code to source ID
func buildSourceMapReverse() map[string]uint {
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[string]uint)
	for _, s := range sources {
		sourceMap[s.Code] = s.ID
	}
	return sourceMap
}

// buildCategoryMap creates a nested map: game_system -> name -> id
func buildCategoryMap() map[string]map[string]uint {
	var categories []models.SkillCategory
	database.DB.Find(&categories)
	categoryMap := make(map[string]map[string]uint)
	for _, c := range categories {
		if categoryMap[c.GameSystem] == nil {
			categoryMap[c.GameSystem] = make(map[string]uint)
		}
		categoryMap[c.GameSystem][c.Name] = c.ID
	}
	return categoryMap
}

// buildDifficultyMap creates a nested map: game_system -> name -> id
func buildDifficultyMap() map[string]map[string]uint {
	var difficulties []models.SkillDifficulty
	database.DB.Find(&difficulties)
	difficultyMap := make(map[string]map[string]uint)
	for _, d := range difficulties {
		if difficultyMap[d.GameSystem] == nil {
			difficultyMap[d.GameSystem] = make(map[string]uint)
		}
		difficultyMap[d.GameSystem][d.Name] = d.ID
	}
	return difficultyMap
}

// buildCharacterClassMap creates a map from character class code to ID
func buildCharacterClassMap() map[string]uint {
	var classes []models.CharacterClass
	database.DB.Find(&classes)
	classMap := make(map[string]uint)
	for _, c := range classes {
		classMap[c.Code] = c.ID
	}
	return classMap
}

// buildSpellSchoolMap creates a nested map: game_system -> name -> id
func buildSpellSchoolMap() map[string]map[string]uint {
	var schools []models.SpellSchool
	database.DB.Find(&schools)
	schoolMap := make(map[string]map[string]uint)
	for _, s := range schools {
		if schoolMap[s.GameSystem] == nil {
			schoolMap[s.GameSystem] = make(map[string]uint)
		}
		schoolMap[s.GameSystem][s.Name] = s.ID
	}
	return schoolMap
}

// buildSkillMap creates a nested map: game_system -> name -> id
func buildSkillMap() map[string]map[string]uint {
	var skills []models.Skill
	database.DB.Find(&skills)
	skillMap := make(map[string]map[string]uint)
	for _, s := range skills {
		if skillMap[s.GameSystem] == nil {
			skillMap[s.GameSystem] = make(map[string]uint)
		}
		skillMap[s.GameSystem][s.Name] = s.ID
	}
	return skillMap
}

// buildWeaponSkillMap creates a nested map: game_system -> name -> id
func buildWeaponSkillMap() map[string]map[string]uint {
	var weaponSkills []models.WeaponSkill
	database.DB.Find(&weaponSkills)
	weaponSkillMap := make(map[string]map[string]uint)
	for _, ws := range weaponSkills {
		if weaponSkillMap[ws.GameSystem] == nil {
			weaponSkillMap[ws.GameSystem] = make(map[string]uint)
		}
		weaponSkillMap[ws.GameSystem][ws.Name] = ws.ID
	}
	return weaponSkillMap
}

// Generic import helper for entities with name + game_system natural key
type ImportConfig struct {
	EntityName string // For error messages, e.g., "skill category"
}

// findOrCreateByNameAndSystem is a helper for import operations
// It looks up an entity by name and game_system, creates if not found
func findOrCreateByNameAndSystem(
	name string,
	gameSystem string,
	model interface{},
	entityName string,
) error {
	result := database.DB.Where("name = ? AND game_system = ?", name, gameSystem).First(model)

	if result.Error == gorm.ErrRecordNotFound {
		if err := database.DB.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create %s %s: %w", entityName, name, err)
		}
	} else if result.Error != nil {
		return fmt.Errorf("failed to query %s %s: %w", entityName, name, result.Error)
	} else {
		if err := database.DB.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update %s %s: %w", entityName, name, err)
		}
	}

	return nil
}

// findOrCreateByCode is a helper for import operations with code as natural key
func findOrCreateByCode(
	code string,
	model interface{},
	entityName string,
) error {
	result := database.DB.Where("code = ?", code).First(model)

	if result.Error == gorm.ErrRecordNotFound {
		if err := database.DB.Create(model).Error; err != nil {
			return fmt.Errorf("failed to create %s %s: %w", entityName, code, err)
		}
	} else if result.Error != nil {
		return fmt.Errorf("failed to query %s %s: %w", entityName, code, result.Error)
	} else {
		if err := database.DB.Save(model).Error; err != nil {
			return fmt.Errorf("failed to update %s %s: %w", entityName, code, err)
		}
	}

	return nil
}
