package gsmaster

import (
	"fmt"

	"gorm.io/gorm"
)

// LearningCostsService provides database-backed learning costs functionality
type LearningCostsService struct {
	db *gorm.DB
}

// NewLearningCostsService creates a new learning costs service
func NewLearningCostsService(db *gorm.DB) *LearningCostsService {
	return &LearningCostsService{db: db}
}

// GetSkillCategoryDB returns the category of a skill from the database
func (s *LearningCostsService) GetSkillCategoryDB(skillName string) (string, error) {
	var skillDef SkillDefinitionNew

	err := s.db.Preload("Category").Where("name = ?", skillName).First(&skillDef).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "Unbekannt", nil
		}
		return "", fmt.Errorf("error finding skill category: %w", err)
	}

	return skillDef.Category.Name, nil
}

// GetSkillDifficultyDB returns the difficulty of a skill from the database
func (s *LearningCostsService) GetSkillDifficultyDB(category, skillName string) (string, error) {
	query := s.db.Preload("Category").Preload("Difficulty")

	if category != "" {
		// Search within specific category
		query = query.Joins("JOIN skill_categories ON skill_definitions.category_id = skill_categories.id").
			Where("skill_definitions.name = ? AND skill_categories.name = ?", skillName, category)
	} else {
		// Search all categories, return first match
		query = query.Where("name = ?", skillName)
	}

	var skillDef SkillDefinitionNew
	err := query.First(&skillDef).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "Unbekannt", nil
		}
		return "", fmt.Errorf("error finding skill difficulty: %w", err)
	}

	return skillDef.Difficulty.Name, nil
}

// GetEPCostForClassCategory returns EP cost per TE for a class/category combination
func (s *LearningCostsService) GetEPCostForClassCategory(classCode, categoryName string) (int, error) {
	var epCost ClassCategoryEPCost

	err := s.db.Joins("JOIN character_classes ON class_category_ep_costs.class_id = character_classes.id").
		Joins("JOIN skill_categories ON class_category_ep_costs.category_id = skill_categories.id").
		Where("character_classes.code = ? AND skill_categories.name = ?", classCode, categoryName).
		First(&epCost).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no EP cost found for class %s and category %s", classCode, categoryName)
		}
		return 0, fmt.Errorf("error finding EP cost: %w", err)
	}

	return epCost.EPPerTE, nil
}

// GetLearnCostForSkill returns the learn cost (LE) for a skill
func (s *LearningCostsService) GetLearnCostForSkill(skillName, categoryName string) (int, error) {
	var skillDef SkillDefinitionNew

	query := s.db.Preload("Difficulty")
	if categoryName != "" {
		query = query.Joins("JOIN skill_categories ON skill_definitions.category_id = skill_categories.id").
			Where("skill_definitions.name = ? AND skill_categories.name = ?", skillName, categoryName)
	} else {
		query = query.Where("name = ?", skillName)
	}

	err := query.First(&skillDef).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("skill %s not found", skillName)
		}
		return 0, fmt.Errorf("error finding skill: %w", err)
	}

	return skillDef.Difficulty.LearnCost, nil
}

// GetImprovementCost returns the TE cost for improving a skill to the next level
func (s *LearningCostsService) GetImprovementCost(skillName, categoryName string, currentLevel int) (int, error) {
	// First get the skill's difficulty
	var skillDef SkillDefinitionNew

	query := s.db.Preload("Difficulty")
	if categoryName != "" {
		query = query.Joins("JOIN skill_categories ON skill_definitions.category_id = skill_categories.id").
			Where("skill_definitions.name = ? AND skill_categories.name = ?", skillName, categoryName)
	} else {
		query = query.Where("name = ?", skillName)
	}

	err := query.First(&skillDef).Error
	if err != nil {
		return 0, fmt.Errorf("skill %s not found: %w", skillName, err)
	}

	// Get the improvement cost for this difficulty and level
	var improvementCost SkillImprovementCost
	err = s.db.Where("difficulty_id = ? AND current_level = ?", skillDef.DifficultyID, currentLevel+1).
		First(&improvementCost).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no improvement cost found for level %d", currentLevel+1)
		}
		return 0, fmt.Errorf("error finding improvement cost: %w", err)
	}

	return improvementCost.TECost, nil
}

// GetSpellEPCost returns EP cost per LE for a class/spell school combination
func (s *LearningCostsService) GetSpellEPCost(classCode, schoolName string) (int, error) {
	var spellCost ClassSpellSchoolCost

	err := s.db.Joins("JOIN character_classes ON class_spell_school_costs.class_id = character_classes.id").
		Joins("JOIN spell_schools ON class_spell_school_costs.school_id = spell_schools.id").
		Where("character_classes.code = ? AND spell_schools.name = ?", classCode, schoolName).
		First(&spellCost).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no spell EP cost found for class %s and school %s", classCode, schoolName)
		}
		return 0, fmt.Errorf("error finding spell EP cost: %w", err)
	}

	return spellCost.EPPerLE, nil
}

// GetSpellLevelCost returns the LE required for a spell level
func (s *LearningCostsService) GetSpellLevelCost(spellLevel int) (int, error) {
	var levelCost SpellLevelCost

	err := s.db.Where("spell_level = ?", spellLevel).First(&levelCost).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("no LE cost found for spell level %d", spellLevel)
		}
		return 0, fmt.Errorf("error finding spell level cost: %w", err)
	}

	return levelCost.LERequired, nil
}

// FindBestCategoryForSkillLearningDB finds the cheapest category for learning a skill
func (s *LearningCostsService) FindBestCategoryForSkillLearningDB(skillName, characterClass string) (string, string, error) {
	type CategoryOption struct {
		Category   string
		Difficulty string
		EPCost     int
	}

	var options []CategoryOption

	// Get all categories where this skill exists
	var skillDefs []SkillDefinitionNew
	err := s.db.Preload("Category").Preload("Difficulty").
		Where("name = ?", skillName).Find(&skillDefs).Error
	if err != nil {
		return "", "", fmt.Errorf("error finding skill definitions: %w", err)
	}

	if len(skillDefs) == 0 {
		return "", "", fmt.Errorf("skill %s not found", skillName)
	}

	// Calculate costs for each category
	for _, skillDef := range skillDefs {
		epPerTE, err := s.GetEPCostForClassCategory(characterClass, skillDef.Category.Name)
		if err != nil {
			continue // Skip if no EP cost found for this class/category
		}

		totalEP := epPerTE * skillDef.Difficulty.LearnCost * 3
		options = append(options, CategoryOption{
			Category:   skillDef.Category.Name,
			Difficulty: skillDef.Difficulty.Name,
			EPCost:     totalEP,
		})
	}

	if len(options) == 0 {
		return "", "", fmt.Errorf("no valid category found for skill %s and class %s", skillName, characterClass)
	}

	// Find the cheapest option
	bestOption := options[0]
	for _, option := range options[1:] {
		if option.EPCost < bestOption.EPCost {
			bestOption = option
		}
	}

	return bestOption.Category, bestOption.Difficulty, nil
}

// FindBestCategoryForSkillImprovementDB finds the cheapest category for improving a skill
func (s *LearningCostsService) FindBestCategoryForSkillImprovementDB(skillName, characterClass string, level int) (string, string, error) {
	type CategoryOption struct {
		Category   string
		Difficulty string
		EPCost     int
	}

	var options []CategoryOption

	// Get all categories where this skill exists
	var skillDefs []SkillDefinitionNew
	err := s.db.Preload("Category").Preload("Difficulty").
		Where("name = ?", skillName).Find(&skillDefs).Error
	if err != nil {
		return "", "", fmt.Errorf("error finding skill definitions: %w", err)
	}

	if len(skillDefs) == 0 {
		return "", "", fmt.Errorf("skill %s not found", skillName)
	}

	// Calculate costs for each category
	for _, skillDef := range skillDefs {
		epPerTE, err := s.GetEPCostForClassCategory(characterClass, skillDef.Category.Name)
		if err != nil {
			continue // Skip if no EP cost found for this class/category
		}

		teCost, err := s.GetImprovementCost(skillName, skillDef.Category.Name, level-1)
		if err != nil {
			continue // Skip if no improvement cost found
		}

		totalEP := epPerTE * teCost
		options = append(options, CategoryOption{
			Category:   skillDef.Category.Name,
			Difficulty: skillDef.Difficulty.Name,
			EPCost:     totalEP,
		})
	}

	if len(options) == 0 {
		return "", "", fmt.Errorf("no valid category found for skill %s and class %s at level %d", skillName, characterClass, level)
	}

	// Find the cheapest option
	bestOption := options[0]
	for _, option := range options[1:] {
		if option.EPCost < bestOption.EPCost {
			bestOption = option
		}
	}

	return bestOption.Category, bestOption.Difficulty, nil
}

// GetAllSkillsInCategory returns all skills in a specific category
func (s *LearningCostsService) GetAllSkillsInCategory(categoryName string) ([]string, error) {
	var skillDefs []SkillDefinitionNew

	err := s.db.Joins("JOIN skill_categories ON skill_definitions.category_id = skill_categories.id").
		Where("skill_categories.name = ?", categoryName).
		Find(&skillDefs).Error
	if err != nil {
		return nil, fmt.Errorf("error finding skills in category %s: %w", categoryName, err)
	}

	var skills []string
	for _, skillDef := range skillDefs {
		skills = append(skills, skillDef.Name)
	}

	return skills, nil
}

// GetAllCategories returns all skill categories
func (s *LearningCostsService) GetAllCategories() ([]SkillCategory, error) {
	var categories []SkillCategory
	err := s.db.Find(&categories).Error
	return categories, err
}
