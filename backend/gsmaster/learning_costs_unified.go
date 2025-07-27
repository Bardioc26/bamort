package gsmaster

import (
	"bamort/database"
	"fmt"
)

// LearningCostsConfig controls whether to use database or static data
type LearningCostsConfig struct {
	UseDatabase bool
	Service     *LearningCostsService
}

var learningCostsConfig = &LearningCostsConfig{
	UseDatabase: false, // Start with static data, can be switched later
}

// InitializeLearningCosts initializes the learning costs system
func InitializeLearningCosts(useDatabase bool) error {
	learningCostsConfig.UseDatabase = useDatabase

	if useDatabase {
		// Initialize database service
		if database.DB == nil {
			return fmt.Errorf("database not initialized")
		}

		// Migrate tables and data
		if err := MigrateLearningCostsData(database.DB); err != nil {
			return fmt.Errorf("failed to migrate learning costs data: %w", err)
		}

		learningCostsConfig.Service = NewLearningCostsService(database.DB)
		fmt.Println("Learning costs system initialized with database backend")
	} else {
		fmt.Println("Learning costs system initialized with static data backend")
	}

	return nil
}

// SwitchToDatabase switches from static data to database
func SwitchToDatabase() error {
	if database.DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// Migrate data to database
	if err := MigrateLearningCostsData(database.DB); err != nil {
		return fmt.Errorf("failed to migrate learning costs data: %w", err)
	}

	learningCostsConfig.UseDatabase = true
	learningCostsConfig.Service = NewLearningCostsService(database.DB)
	fmt.Println("Switched to database backend for learning costs")

	return nil
}

// SwitchToStatic switches from database to static data
func SwitchToStatic() {
	learningCostsConfig.UseDatabase = false
	learningCostsConfig.Service = nil
	fmt.Println("Switched to static data backend for learning costs")
}

// GetSkillCategoryNew returns the category of a skill (unified interface)
func GetSkillCategoryNew(skillName string) (string, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.GetSkillCategoryDB(skillName)
	}

	// Fallback to static data
	category := GetSkillCategory(skillName)
	if category == "Unbekannt" {
		return category, fmt.Errorf("skill %s not found", skillName)
	}
	return category, nil
}

// GetSkillDifficultyNew returns the difficulty of a skill (unified interface)
func GetSkillDifficultyNew(category, skillName string) (string, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.GetSkillDifficultyDB(category, skillName)
	}

	// Fallback to static data
	difficulty := GetSkillDifficulty(category, skillName)
	if difficulty == "Unbekannt" {
		return difficulty, fmt.Errorf("skill %s difficulty not found", skillName)
	}
	return difficulty, nil
}

// FindBestCategoryForSkillLearningNew finds the cheapest category for learning a skill (unified interface)
func FindBestCategoryForSkillLearningNew(skillName, characterClass string) (string, string, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.FindBestCategoryForSkillLearningDB(skillName, characterClass)
	}

	// Fallback to static data
	return findBestCategoryForSkillLearning(skillName, characterClass)
}

// FindBestCategoryForSkillImprovementNew finds the cheapest category for improving a skill (unified interface)
func FindBestCategoryForSkillImprovementNew(skillName, characterClass string, level int) (string, string, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.FindBestCategoryForSkillImprovementDB(skillName, characterClass, level)
	}

	// Fallback to static data
	return findBestCategoryForSkillImprovement(skillName, characterClass, level)
}

// GetEPCostForClassCategoryNew returns EP cost per TE for a class/category combination (unified interface)
func GetEPCostForClassCategoryNew(classCode, categoryName string) (int, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.GetEPCostForClassCategory(classCode, categoryName)
	}

	// Fallback to static data
	epPerTE, exists := learningCostsData.EPPerTE[classCode][categoryName]
	if !exists {
		return 0, fmt.Errorf("EP cost not found for class %s and category %s", classCode, categoryName)
	}
	return epPerTE, nil
}

// GetLearnCostForSkillNew returns the learn cost (LE) for a skill (unified interface)
func GetLearnCostForSkillNew(skillName, categoryName string) (int, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.GetLearnCostForSkill(skillName, categoryName)
	}

	// Fallback to static data
	if categoryName == "" {
		// Find the category first
		category := GetSkillCategory(skillName)
		if category == "Unbekannt" {
			return 0, fmt.Errorf("skill %s not found", skillName)
		}
		categoryName = category
	}

	difficulty := GetSkillDifficulty(categoryName, skillName)
	if difficulty == "Unbekannt" {
		return 0, fmt.Errorf("skill %s difficulty not found in category %s", skillName, categoryName)
	}

	diffData, exists := learningCostsData.ImprovementCost[categoryName][difficulty]
	if !exists {
		return 0, fmt.Errorf("difficulty data not found for category %s and difficulty %s", categoryName, difficulty)
	}

	return diffData.LearnCost, nil
}

// GetImprovementCostNew returns the TE cost for improving a skill to the next level (unified interface)
func GetImprovementCostNew(skillName, categoryName string, currentLevel int) (int, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.GetImprovementCost(skillName, categoryName, currentLevel)
	}

	// Fallback to static data
	if categoryName == "" {
		// Find the category first
		category := GetSkillCategory(skillName)
		if category == "Unbekannt" {
			return 0, fmt.Errorf("skill %s not found", skillName)
		}
		categoryName = category
	}

	difficulty := GetSkillDifficulty(categoryName, skillName)
	if difficulty == "Unbekannt" {
		return 0, fmt.Errorf("skill %s difficulty not found in category %s", skillName, categoryName)
	}

	diffData, exists := learningCostsData.ImprovementCost[categoryName][difficulty]
	if !exists {
		return 0, fmt.Errorf("difficulty data not found for category %s and difficulty %s", categoryName, difficulty)
	}

	teCost, exists := diffData.TrainCosts[currentLevel+1]
	if !exists {
		return 0, fmt.Errorf("no improvement cost found for level %d", currentLevel+1)
	}

	return teCost, nil
}

// GetSpellEPCostNew returns EP cost per LE for a class/spell school combination (unified interface)
func GetSpellEPCostNew(classCode, schoolName string) (int, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.GetSpellEPCost(classCode, schoolName)
	}

	// Fallback to static data
	spellEP, exists := learningCostsData.SpellEPPerLE[classCode][schoolName]
	if !exists {
		return 0, fmt.Errorf("spell EP cost not found for class %s and school %s", classCode, schoolName)
	}

	if spellEP == 0 {
		return 0, fmt.Errorf("spell school %s not available for class %s", schoolName, classCode)
	}

	return spellEP, nil
}

// GetSpellLevelCostNew returns the LE required for a spell level (unified interface)
func GetSpellLevelCostNew(spellLevel int) (int, error) {
	if learningCostsConfig.UseDatabase && learningCostsConfig.Service != nil {
		return learningCostsConfig.Service.GetSpellLevelCost(spellLevel)
	}

	// Fallback to static data
	leCost, exists := learningCostsData.SpellLEPerLevel[spellLevel]
	if !exists {
		return 0, fmt.Errorf("no LE cost found for spell level %d", spellLevel)
	}

	return leCost, nil
}
