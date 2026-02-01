package maintenance

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

func TestMigrateSkillCategoriesToRelations(t *testing.T) {
	setupTestEnvironment(t)

	// Reset and setup fresh test database to avoid interference from other tests
	database.ResetTestDB()
	database.SetupTestDB()
	testDB := database.DB

	// Ensure database is valid
	if testDB == nil {
		t.Fatal("Database connection is nil")
	}

	// Create a test source - use unique code to avoid conflicts
	source := models.Source{
		Code:         "TSTMIG1",
		Name:         "Test Migration Source",
		GameSystemId: 1,
		IsActive:     true,
	}
	if err := testDB.Create(&source).Error; err != nil {
		t.Fatalf("Failed to create test source: %v", err)
	}

	// Create test skills with old-style category/difficulty - use unique names to avoid conflicts
	testSkills := []models.Skill{
		{
			Name:             "TestMigSkill_Schwimmen",
			Category:         "Körper",
			Difficulty:       "leicht",
			GameSystemId: 1,
			Initialwert:      12,
			Improvable:       true,
			Bonuseigenschaft: "Gw",
			SourceID:         source.ID,
		},
		{
			Name:             "TestMigSkill_Klettern",
			Category:         "Körper",
			Difficulty:       "normal",
			GameSystemId: 1,
			Initialwert:      10,
			Improvable:       true,
			Bonuseigenschaft: "Gw",
			SourceID:         source.ID,
		},
		{
			Name:             "TestMigSkill_LesenSchreiben",
			Category:         "Wissen",
			Difficulty:       "schwer",
			GameSystemId: 1,
			Initialwert:      0,
			Improvable:       true,
			Bonuseigenschaft: "In",
			SourceID:         source.ID,
		},
	}

	for i := range testSkills {
		if err := testDB.Create(&testSkills[i]).Error; err != nil {
			t.Fatalf("Failed to create test skill: %v", err)
		}
	}

	// Run migration
	if err := MigrateSkillCategoriesToRelations(testDB); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Verify migration results
	for _, skill := range testSkills {
		var scds []models.SkillCategoryDifficulty
		err := testDB.Preload("SkillCategory").Preload("SkillDifficulty").
			Where("skill_id = ?", skill.ID).Find(&scds).Error
		if err != nil {
			t.Errorf("Failed to find migrated data for skill %s: %v", skill.Name, err)
			continue
		}

		if len(scds) == 0 {
			t.Errorf("No SkillCategoryDifficulty created for skill %s", skill.Name)
			continue
		}

		scd := scds[0]
		if scd.SkillCategory.Name != skill.Category {
			t.Errorf("Category mismatch for %s: expected %s, got %s",
				skill.Name, skill.Category, scd.SkillCategory.Name)
		}

		if scd.SkillDifficulty.Name != skill.Difficulty {
			t.Errorf("Difficulty mismatch for %s: expected %s, got %s",
				skill.Name, skill.Difficulty, scd.SkillDifficulty.Name)
		}
	}

	// Test idempotency - running migration again should not create duplicates
	if err := MigrateSkillCategoriesToRelations(testDB); err != nil {
		t.Fatalf("Second migration failed: %v", err)
	}

	for _, skill := range testSkills {
		var count int64
		testDB.Model(&models.SkillCategoryDifficulty{}).Where("skill_id = ?", skill.ID).Count(&count)
		if count != 1 {
			t.Errorf("Idempotency check failed for %s: expected 1 record, got %d", skill.Name, count)
		}
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
		{"unknown", 10}, // default
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

func TestMigrateSkillCategoryDifficulty_NoCategory(t *testing.T) {
	setupTestEnvironment(t)

	// Reset and setup fresh test database to avoid interference from other tests
	database.ResetTestDB()
	database.SetupTestDB()
	testDB := database.DB

	// Ensure database is valid
	if testDB == nil {
		t.Fatal("Database connection is nil")
	}

	// Use existing source or create one with a unique code
	var source models.Source
	err := testDB.Where("code = ?", "KOD").First(&source).Error
	if err != nil {
		// Create a test source if KOD doesn't exist
		source = models.Source{
			Code:       "TSTMIG2",
			Name:       "Test Migration Source 2",
			GameSystemId: 1,
			IsActive:   true,
		}
		if err := testDB.Create(&source).Error; err != nil {
			t.Fatalf("Failed to create test source: %v", err)
		}
	}

	// Create skill without category (should use default) - use unique name
	skill := models.Skill{
		Name:        "TestMigSkill_NoCategory",
		Category:    "", // Empty category
		Difficulty:  "", // Empty difficulty
		GameSystemId: 1,
		Initialwert: 10,
		SourceID:    source.ID,
	}
	if err := testDB.Create(&skill).Error; err != nil {
		t.Fatalf("Failed to create test skill: %v", err)
	}

	// Migrate
	if err := migrateSkillCategoryDifficulty(testDB, &skill); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	// Verify defaults were used
	var scd models.SkillCategoryDifficulty
	err = testDB.Preload("SkillCategory").Preload("SkillDifficulty").
		Where("skill_id = ?", skill.ID).First(&scd).Error
	if err != nil {
		t.Fatalf("Failed to find migrated data: %v", err)
	}

	if scd.SkillCategory.Name != "Alltag" {
		t.Errorf("Expected default category 'Alltag', got %s", scd.SkillCategory.Name)
	}

	if scd.SkillDifficulty.Name != "normal" {
		t.Errorf("Expected default difficulty 'normal', got %s", scd.SkillDifficulty.Name)
	}
}
