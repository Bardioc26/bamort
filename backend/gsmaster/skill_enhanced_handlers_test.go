package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"os"
	"testing"
)

func setupTestEnvironment(t *testing.T) {
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original != "" {
			os.Setenv("ENVIRONMENT", original)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})
}

// Helper function to get or create a source
func getOrCreateSource(code, name string) models.Source {
	var source models.Source
	if err := database.DB.Where("code = ?", code).First(&source).Error; err != nil {
		source = models.Source{
			Code:       code,
			Name:       name,
			GameSystem: "midgard",
			IsActive:   true,
		}
		database.DB.Create(&source)
	}
	return source
}

// Helper function to get or create a category
func getOrCreateCategory(name string, sourceID uint) models.SkillCategory {
	var category models.SkillCategory
	if err := database.DB.Where("name = ? AND game_system = ?", name, "midgard").First(&category).Error; err != nil {
		category = models.SkillCategory{
			Name:       name,
			GameSystem: "midgard",
			SourceID:   sourceID,
		}
		database.DB.Create(&category)
	}
	return category
}

// Helper function to get or create a difficulty
func getOrCreateDifficulty(name string) models.SkillDifficulty {
	var difficulty models.SkillDifficulty
	if err := database.DB.Where("name = ? AND game_system = ?", name, "midgard").First(&difficulty).Error; err != nil {
		difficulty = models.SkillDifficulty{
			Name:       name,
			GameSystem: "midgard",
		}
		database.DB.Create(&difficulty)
	}
	return difficulty
}

func TestGetSkillWithCategories(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data
	source := getOrCreateSource("TSTSKL", "TestSkill")

	skill := models.Skill{
		Name:             "TestSchwimmen",
		GameSystem:       "midgard",
		Initialwert:      12,
		Improvable:       true,
		Bonuseigenschaft: "Gw",
		SourceID:         source.ID,
	}
	database.DB.Create(&skill)

	category := getOrCreateCategory("Körper", source.ID)
	difficulty := getOrCreateDifficulty("leicht")

	scd := models.SkillCategoryDifficulty{
		SkillID:           skill.ID,
		SkillCategoryID:   category.ID,
		SkillDifficultyID: difficulty.ID,
		LearnCost:         5,
		SCategory:         category.Name,
		SDifficulty:       difficulty.Name,
	}
	database.DB.Create(&scd)

	// Test GetSkillWithCategories
	result, err := GetSkillWithCategories(skill.ID)
	if err != nil {
		t.Fatalf("GetSkillWithCategories failed: %v", err)
	}

	if result.Name != "TestSchwimmen" {
		t.Errorf("Expected skill name 'TestSchwimmen', got '%s'", result.Name)
	}

	if len(result.Categories) != 1 {
		t.Fatalf("Expected 1 category, got %d", len(result.Categories))
	}

	if result.Categories[0].CategoryName != "Körper" {
		t.Errorf("Expected category 'Körper', got '%s'", result.Categories[0].CategoryName)
	}

	if result.Categories[0].DifficultyName != "leicht" {
		t.Errorf("Expected difficulty 'leicht', got '%s'", result.Categories[0].DifficultyName)
	}

	if len(result.Difficulties) != 1 || result.Difficulties[0] != "leicht" {
		t.Errorf("Expected difficulties ['leicht'], got %v", result.Difficulties)
	}
}

func TestGetSkillWithCategories_MultipleCategories(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data
	source := getOrCreateSource("TSTMC", "TestMultiCat")

	skill := models.Skill{
		Name:             "TestReiten",
		GameSystem:       "midgard",
		Initialwert:      5,
		Improvable:       true,
		Bonuseigenschaft: "Gw",
		SourceID:         source.ID,
	}
	database.DB.Create(&skill)

	// Create multiple categories
	category1 := getOrCreateCategory("Bewegung", source.ID)
	category2 := getOrCreateCategory("Reiten", source.ID)
	difficultyNormal := getOrCreateDifficulty("normal")
	difficultySchwer := getOrCreateDifficulty("schwer")

	// Create relationships
	scd1 := models.SkillCategoryDifficulty{
		SkillID:           skill.ID,
		SkillCategoryID:   category1.ID,
		SkillDifficultyID: difficultyNormal.ID,
		LearnCost:         10,
		SCategory:         category1.Name,
		SDifficulty:       difficultyNormal.Name,
	}
	database.DB.Create(&scd1)

	scd2 := models.SkillCategoryDifficulty{
		SkillID:           skill.ID,
		SkillCategoryID:   category2.ID,
		SkillDifficultyID: difficultySchwer.ID,
		LearnCost:         20,
		SCategory:         category2.Name,
		SDifficulty:       difficultySchwer.Name,
	}
	database.DB.Create(&scd2)

	// Test
	result, err := GetSkillWithCategories(skill.ID)
	if err != nil {
		t.Fatalf("GetSkillWithCategories failed: %v", err)
	}

	if len(result.Categories) != 2 {
		t.Fatalf("Expected 2 categories, got %d", len(result.Categories))
	}

	// Check that both categories exist (order may vary)
	foundMovement := false
	foundRiding := false
	for _, cat := range result.Categories {
		if cat.CategoryName == "Bewegung" && cat.DifficultyName == "normal" {
			foundMovement = true
		}
		if cat.CategoryName == "Reiten" && cat.DifficultyName == "schwer" {
			foundRiding = true
		}
	}

	if !foundMovement {
		t.Error("Expected to find 'Bewegung/normal' category")
	}
	if !foundRiding {
		t.Error("Expected to find 'Reiten/schwer' category")
	}
}

func TestUpdateSkillWithCategories(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data
	source := getOrCreateSource("TSTUPD", "TestUpdate")

	skill := models.Skill{
		Name:             "TestKlettern",
		GameSystem:       "midgard",
		Initialwert:      10,
		Improvable:       true,
		Bonuseigenschaft: "Gw",
		SourceID:         source.ID,
	}
	database.DB.Create(&skill)

	category1 := getOrCreateCategory("Körper", source.ID)
	category2 := getOrCreateCategory("Alltag", source.ID)
	difficultyNormal := getOrCreateDifficulty("normal")
	difficultyLeicht := getOrCreateDifficulty("leicht")

	// Create initial relationship
	scd := models.SkillCategoryDifficulty{
		SkillID:           skill.ID,
		SkillCategoryID:   category1.ID,
		SkillDifficultyID: difficultyNormal.ID,
		LearnCost:         10,
		SCategory:         category1.Name,
		SDifficulty:       difficultyNormal.Name,
	}
	database.DB.Create(&scd)

	// Update with new categories
	updateReq := SkillUpdateRequest{
		Skill: models.Skill{
			ID:               skill.ID,
			Name:             "TestKlettern",
			GameSystem:       "midgard",
			Initialwert:      12, // Changed
			Improvable:       true,
			Bonuseigenschaft: "St", // Changed
			SourceID:         source.ID,
		},
		CategoryDifficulties: []CategoryDifficultyPair{
			{
				CategoryID:   category1.ID,
				DifficultyID: difficultyLeicht.ID, // Changed difficulty
				LearnCost:    5,
			},
			{
				CategoryID:   category2.ID, // Added category
				DifficultyID: difficultyNormal.ID,
				LearnCost:    10,
			},
		},
	}

	err := UpdateSkillWithCategories(skill.ID, updateReq)
	if err != nil {
		t.Fatalf("UpdateSkillWithCategories failed: %v", err)
	}

	// Verify update
	result, err := GetSkillWithCategories(skill.ID)
	if err != nil {
		t.Fatalf("GetSkillWithCategories failed: %v", err)
	}

	if result.Initialwert != 12 {
		t.Errorf("Expected initialwert 12, got %d", result.Initialwert)
	}

	if result.Bonuseigenschaft != "St" {
		t.Errorf("Expected bonuseigenschaft 'St', got '%s'", result.Bonuseigenschaft)
	}

	if len(result.Categories) != 2 {
		t.Fatalf("Expected 2 categories after update, got %d", len(result.Categories))
	}

	// Verify old category has new difficulty and new category exists
	foundKoerperLeicht := false
	foundAlltagNormal := false
	for _, cat := range result.Categories {
		if cat.CategoryName == "Körper" && cat.DifficultyName == "leicht" {
			foundKoerperLeicht = true
		}
		if cat.CategoryName == "Alltag" && cat.DifficultyName == "normal" {
			foundAlltagNormal = true
		}
	}

	if !foundKoerperLeicht {
		t.Error("Expected to find 'Körper/leicht' category after update")
	}
	if !foundAlltagNormal {
		t.Error("Expected to find 'Alltag/normal' category after update")
	}
}

func TestGetDefaultLearnCost(t *testing.T) {
	tests := []struct {
		difficulty string
		expected   int
	}{
		{"leicht", 5},
		{"easy", 5},
		{"normal", 10},
		{"standard", 10},
		{"schwer", 20},
		{"hard", 20},
		{"sehr schwer", 30},
		{"very hard", 30},
		{"unknown", 10},
	}

	for _, tt := range tests {
		t.Run(tt.difficulty, func(t *testing.T) {
			result := getDefaultLearnCost(tt.difficulty)
			if result != tt.expected {
				t.Errorf("getDefaultLearnCost(%s) = %d, expected %d",
					tt.difficulty, result, tt.expected)
			}
		})
	}
}
