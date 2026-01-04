package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"testing"
)

func TestCreateSkillWithCategories(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test dependencies
	source := getOrCreateSource("TSTCRT", "TestCreate")
	category := getOrCreateCategory("Alltag", source.ID)
	difficulty := getOrCreateDifficulty("normal")

	// Prepare create request
	req := SkillUpdateRequest{
		Skill: models.Skill{
			Name:             "Neue Fertigkeit",
			GameSystem:       "midgard",
			Beschreibung:     "Test Fertigkeit",
			Initialwert:      5,
			BasisWert:        0,
			Bonuseigenschaft: "In",
			Improvable:       true,
			InnateSkill:      false,
			SourceID:         source.ID,
			PageNumber:       42,
		},
		CategoryDifficulties: []CategoryDifficultyPair{
			{
				CategoryID:   category.ID,
				DifficultyID: difficulty.ID,
			},
		},
	}

	// Test creating new skill
	skillID, err := CreateSkillWithCategories(req)
	if err != nil {
		t.Fatalf("CreateSkillWithCategories failed: %v", err)
	}

	if skillID == 0 {
		t.Fatalf("Expected non-zero skill ID, got 0")
	}

	// Verify skill was created
	var skill models.Skill
	if err := database.DB.First(&skill, skillID).Error; err != nil {
		t.Fatalf("Failed to retrieve created skill: %v", err)
	}

	if skill.Name != "Neue Fertigkeit" {
		t.Errorf("Expected name 'Neue Fertigkeit', got '%s'", skill.Name)
	}

	if skill.Initialwert != 5 {
		t.Errorf("Expected initialwert 5, got %d", skill.Initialwert)
	}

	if skill.BasisWert != 0 {
		t.Errorf("Expected basiswert 0, got %d", skill.BasisWert)
	}

	// Verify category-difficulty relationship
	var scd models.SkillCategoryDifficulty
	if err := database.DB.Where("skill_id = ?", skillID).First(&scd).Error; err != nil {
		t.Fatalf("Failed to retrieve skill category difficulty: %v", err)
	}

	if scd.SkillCategoryID != category.ID {
		t.Errorf("Expected category ID %d, got %d", category.ID, scd.SkillCategoryID)
	}

	if scd.SkillDifficultyID != difficulty.ID {
		t.Errorf("Expected difficulty ID %d, got %d", difficulty.ID, scd.SkillDifficultyID)
	}
}

func TestCreateSkillWithMultipleCategories(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test dependencies
	source := getOrCreateSource("TSTMLT", "TestMultiple")
	category1 := getOrCreateCategory("Körper", source.ID)
	category2 := getOrCreateCategory("Geist", source.ID)
	difficulty1 := getOrCreateDifficulty("leicht")
	difficulty2 := getOrCreateDifficulty("schwer")

	// Prepare create request with multiple categories
	req := SkillUpdateRequest{
		Skill: models.Skill{
			Name:        "Multi-Kategorie Fertigkeit",
			GameSystem:  "midgard",
			Initialwert: 10,
			Improvable:  true,
			SourceID:    source.ID,
		},
		CategoryDifficulties: []CategoryDifficultyPair{
			{
				CategoryID:   category1.ID,
				DifficultyID: difficulty1.ID,
			},
			{
				CategoryID:   category2.ID,
				DifficultyID: difficulty2.ID,
			},
		},
	}

	// Test creating skill with multiple categories
	skillID, err := CreateSkillWithCategories(req)
	if err != nil {
		t.Fatalf("CreateSkillWithCategories failed: %v", err)
	}

	// Verify both category-difficulty relationships exist
	var scds []models.SkillCategoryDifficulty
	if err := database.DB.Where("skill_id = ?", skillID).Find(&scds).Error; err != nil {
		t.Fatalf("Failed to retrieve skill category difficulties: %v", err)
	}

	if len(scds) != 2 {
		t.Fatalf("Expected 2 category-difficulty relationships, got %d", len(scds))
	}
}

func TestCreateSkillValidation(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Test creating skill without name
	req := SkillUpdateRequest{
		Skill: models.Skill{
			GameSystem:  "midgard",
			Initialwert: 5,
		},
		CategoryDifficulties: []CategoryDifficultyPair{},
	}

	_, err := CreateSkillWithCategories(req)
	if err == nil {
		t.Error("Expected error when creating skill without name, got nil")
	}
}
