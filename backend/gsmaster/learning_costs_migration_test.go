package gsmaster

import (
	"bamort/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateLearningCostsData(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, false) // Use in-memory SQLite, no test data loading
	defer database.ResetTestDB()

	// Run migration
	err := MigrateLearningCostsData(database.DB)
	assert.NoError(t, err, "Migration should succeed")

	// Verify character classes were created
	var classCount int64
	err = database.DB.Model(&CharacterClass{}).Count(&classCount).Error
	assert.NoError(t, err)
	assert.Greater(t, classCount, int64(10), "Should have multiple character classes")

	// Verify skill categories were created
	var categoryCount int64
	err = database.DB.Model(&SkillCategory{}).Count(&categoryCount).Error
	assert.NoError(t, err)
	assert.Greater(t, categoryCount, int64(5), "Should have multiple skill categories")

	// Verify spell schools were created
	var schoolCount int64
	err = database.DB.Model(&SpellSchool{}).Count(&schoolCount).Error
	assert.NoError(t, err)
	assert.Greater(t, schoolCount, int64(5), "Should have multiple spell schools")

	// Verify skill definitions were created
	var skillDefCount int64
	err = database.DB.Model(&SkillDefinitionNew{}).Count(&skillDefCount).Error
	assert.NoError(t, err)
	assert.Greater(t, skillDefCount, int64(10), "Should have multiple skill definitions")

	// Test specific character class lookup
	var assassinClass CharacterClass
	err = database.DB.Where("code = ?", "As").First(&assassinClass).Error
	assert.NoError(t, err)
	assert.Equal(t, "Assassine", assassinClass.Name)

	// Test specific skill category lookup
	var alltagCategory SkillCategory
	err = database.DB.Where("name = ?", "Alltag").First(&alltagCategory).Error
	assert.NoError(t, err)
	assert.Equal(t, "Alltägliche Fertigkeiten", alltagCategory.Description)

	// Test skill definition with relationships
	var skillDef SkillDefinitionNew
	err = database.DB.Preload("Category").Preload("Difficulty").
		Joins("JOIN skill_categories ON skill_definitions.category_id = skill_categories.id").
		Where("skill_categories.name = ? AND skill_definitions.name = ?", "Alltag", "Klettern").
		First(&skillDef).Error
	assert.NoError(t, err)
	assert.Equal(t, "Klettern", skillDef.Name)
	assert.Equal(t, "Alltag", skillDef.Category.Name)
}

func TestMigrationIdempotency(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, false)
	defer database.ResetTestDB()

	// Run migration twice
	err := MigrateLearningCostsData(database.DB)
	assert.NoError(t, err, "First migration should succeed")

	err = MigrateLearningCostsData(database.DB)
	assert.NoError(t, err, "Second migration should succeed (idempotent)")

	// Verify we don't have duplicates
	var classCount int64
	err = database.DB.Model(&CharacterClass{}).Count(&classCount).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(15), classCount, "Should have exactly 15 character classes after double migration")
}

func TestLearningCostsServiceIntegration(t *testing.T) {
	// Setup test database with migration
	database.SetupTestDB(true, false)
	defer database.ResetTestDB()

	err := MigrateLearningCostsData(database.DB)
	assert.NoError(t, err)

	// Initialize service
	service := NewLearningCostsService(database.DB)

	// Test skill category lookup
	category, err := service.GetSkillCategoryDB("Klettern")
	assert.NoError(t, err)
	assert.Equal(t, "Alltag", category)

	// Test skill difficulty lookup
	difficulty, err := service.GetSkillDifficultyDB("Alltag", "Klettern")
	assert.NoError(t, err)
	assert.Equal(t, "leicht", difficulty)

	// Test EP cost lookup
	epCost, err := service.GetEPCostForClassCategory("As", "Alltag")
	assert.NoError(t, err)
	assert.Equal(t, 20, epCost) // Should match the static data
}
