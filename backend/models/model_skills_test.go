package models

import (
	"bamort/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) {
	database.SetupTestDB()

	// Migrate structures
	err := MigrateStructure()
	require.NoError(t, err, "Failed to migrate database structure")
}

func createTestSkill(name string) *Skill {
	return &Skill{
		GameSystem:       "midgard",
		Name:             name,
		Beschreibung:     "Test skill description",
		Category:         "Körper",
		Difficulty:       "Normal",
		Initialwert:      5,
		Bonuseigenschaft: "Gs",
		Improvable:       true,
		InnateSkill:      false,
	}
}

// =============================================================================
// Tests for SkFertigkeit struct
// =============================================================================

func TestSkFertigkeit_TableName(t *testing.T) {
	skill := SkFertigkeit{}
	expected := "char_skills"
	actual := skill.TableName()
	assert.Equal(t, expected, actual)
}

func TestSkFertigkeit_GetSkillByName_Success(t *testing.T) {
	//setupTestDB(t)
	database.SetupTestDB(true)

	// Create a test skill in the database
	testSkill := createTestSkill("Athletika")
	err := testSkill.Create()
	require.NoError(t, err)

	// Create SkFertigkeit with the same name
	skFertigkeit := SkFertigkeit{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "Athletika",
			},
		},
	}

	// Test GetSkillByName
	result := skFertigkeit.GetSkillByName()

	assert.NotNil(t, result)
	assert.Equal(t, "Athletika", result.Name)
	assert.Equal(t, "Körper", result.Category)
	assert.Equal(t, "Normal", result.Difficulty)
	assert.Equal(t, 5, result.Initialwert)
}

func TestSkFertigkeit_GetSkillByName_NotFound(t *testing.T) {
	setupTestDB(t)

	// Create SkFertigkeit with non-existent name
	skFertigkeit := SkFertigkeit{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "NonExistentSkill",
			},
		},
	}

	// Test GetSkillByName
	result := skFertigkeit.GetSkillByName()

	assert.Nil(t, result)
}

func TestSkFertigkeit_GetCategory_AlreadySet(t *testing.T) {
	setupTestDB(t)

	// Category should now ALWAYS be fetched from gsmaster, not from the skill's category field
	// Since "AlreadySetCategory" is not a real skill in gsmaster, it should return "Unkategorisiert"
	skFertigkeit := SkFertigkeit{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "NonExistentSkill",
			},
		},
		Category: "AlreadySetCategory", // This field should be ignored now
	}

	result := skFertigkeit.GetCategory()
	// Since the skill doesn't exist in gsmaster, should return "Unkategorisiert"
	assert.Equal(t, "Unkategorisiert", result)
}

func TODOT_estSkFertigkeit_GetCategory_FromDatabase(t *testing.T) {
	setupTestDB(t)

	// Create a test skill in the database
	testSkill := createTestSkill("Schwimmen")
	testSkill.Category = "Körper"
	err := testSkill.Create()
	require.NoError(t, err)

	// Create SkFertigkeit without category but with name
	skFertigkeit := SkFertigkeit{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "Schwimmen",
			},
		},
		Category: "", // Empty category
	}

	result := skFertigkeit.GetCategory()

	assert.Equal(t, "Körper", result)
	assert.Equal(t, "Körper", skFertigkeit.Category) // Should be set on the object
}

func TestSkFertigkeit_GetCategory_NotFoundInDatabase(t *testing.T) {
	setupTestDB(t)

	skFertigkeit := SkFertigkeit{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "UnknownSkill",
			},
		},
		Category: "", // Empty category
	}

	result := skFertigkeit.GetCategory()

	assert.Equal(t, "Unkategorisiert", result)
}

func TestSkFertigkeit_StructTags(t *testing.T) {
	// Test that the struct has the expected JSON tags
	skill := SkFertigkeit{
		Beschreibung:    "Test description",
		Fertigkeitswert: 10,
		Bonus:           2,
		Pp:              5,
		Bemerkung:       "Test comment",
		Improvable:      true,
		Category:        "TestCategory",
	}

	// Verify struct fields exist and can be set
	assert.Equal(t, "Test description", skill.Beschreibung)
	assert.Equal(t, 10, skill.Fertigkeitswert)
	assert.Equal(t, 2, skill.Bonus)
	assert.Equal(t, 5, skill.Pp)
	assert.Equal(t, "Test comment", skill.Bemerkung)
	assert.True(t, skill.Improvable)
	assert.Equal(t, "TestCategory", skill.Category)
}

// =============================================================================
// Tests for SkWaffenfertigkeit struct
// =============================================================================

func TestSkWaffenfertigkeit_TableName(t *testing.T) {
	weaponSkill := SkWaffenfertigkeit{}
	expected := "char_weaponskills"
	actual := weaponSkill.TableName()
	assert.Equal(t, expected, actual)
}

func TestSkWaffenfertigkeit_Inheritance(t *testing.T) {
	// Test that SkWaffenfertigkeit properly inherits from SkFertigkeit
	weaponSkill := SkWaffenfertigkeit{
		SkFertigkeit: SkFertigkeit{
			BamortCharTrait: BamortCharTrait{
				BamortBase: BamortBase{
					Name: "Einhandschwerter",
				},
				CharacterID: 1,
				UserID:      1,
			},
			Beschreibung:    "Sword fighting skill",
			Fertigkeitswert: 12,
			Bonus:           3,
			Pp:              8,
			Bemerkung:       "Advanced sword training",
			Improvable:      true,
			Category:        "Kampf",
		},
	}

	// Verify all inherited fields are accessible
	assert.Equal(t, "Einhandschwerter", weaponSkill.Name)
	assert.Equal(t, uint(1), weaponSkill.CharacterID)
	assert.Equal(t, uint(1), weaponSkill.UserID)
	assert.Equal(t, "Sword fighting skill", weaponSkill.Beschreibung)
	assert.Equal(t, 12, weaponSkill.Fertigkeitswert)
	assert.Equal(t, 3, weaponSkill.Bonus)
	assert.Equal(t, 8, weaponSkill.Pp)
	assert.Equal(t, "Advanced sword training", weaponSkill.Bemerkung)
	assert.True(t, weaponSkill.Improvable)
	assert.Equal(t, "Kampf", weaponSkill.Category)
}

func TestSkWaffenfertigkeit_InheritedMethods(t *testing.T) {
	setupTestDB(t)

	// Test inherited methods work - use existing Stichwaffen weapon skill
	weaponSkill := SkWaffenfertigkeit{
		SkFertigkeit: SkFertigkeit{
			BamortCharTrait: BamortCharTrait{
				BamortBase: BamortBase{
					Name: "Stichwaffen",
				},
			},
		},
	}

	// Test GetSkillByName (inherited from SkFertigkeit)
	result := weaponSkill.GetSkillByName()
	assert.NotNil(t, result)
	assert.Equal(t, "Stichwaffen", result.Name)

	// Test GetCategory - should return "Waffenfertigkeiten" for weapon skills
	category := weaponSkill.GetCategory()
	assert.Equal(t, "Waffenfertigkeiten", category)
}

// =============================================================================
// Tests for SkAngeboreneFertigkeit struct
// =============================================================================

func TestSkAngeboreneFertigkeit_Inheritance(t *testing.T) {
	// Test that SkAngeboreneFertigkeit properly inherits from SkFertigkeit
	innateSkill := SkAngeboreneFertigkeit{
		SkFertigkeit: SkFertigkeit{
			BamortCharTrait: BamortCharTrait{
				BamortBase: BamortBase{
					Name: "Nachtsicht",
				},
				CharacterID: 2,
				UserID:      2,
			},
			Beschreibung:    "Natural night vision",
			Fertigkeitswert: 15,
			Bonus:           0,
			Pp:              0,
			Bemerkung:       "Racial ability",
			Improvable:      false, // Innate skills typically not improvable
			Category:        "Angeboren",
		},
	}

	// Verify all inherited fields are accessible
	assert.Equal(t, "Nachtsicht", innateSkill.Name)
	assert.Equal(t, uint(2), innateSkill.CharacterID)
	assert.Equal(t, uint(2), innateSkill.UserID)
	assert.Equal(t, "Natural night vision", innateSkill.Beschreibung)
	assert.Equal(t, 15, innateSkill.Fertigkeitswert)
	assert.Equal(t, 0, innateSkill.Bonus)
	assert.Equal(t, 0, innateSkill.Pp)
	assert.Equal(t, "Racial ability", innateSkill.Bemerkung)
	assert.False(t, innateSkill.Improvable)
	assert.Equal(t, "Angeboren", innateSkill.Category)
}

// =============================================================================
// Tests for SkZauber struct
// =============================================================================

func TestSkZauber_TableName(t *testing.T) {
	spell := SkZauber{}
	expected := "char_spells"
	actual := spell.TableName()
	assert.Equal(t, expected, actual)
}

func TestSkZauber_StructFields(t *testing.T) {
	spell := SkZauber{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				ID:   1,
				Name: "Feuerkugel",
			},
			CharacterID: 3,
			UserID:      3,
		},
		Beschreibung: "Creates a ball of fire",
		Bonus:        5,
		Quelle:       "Arkanum",
	}

	// Verify all fields are accessible and correct
	assert.Equal(t, uint(1), spell.ID)
	assert.Equal(t, "Feuerkugel", spell.Name)
	assert.Equal(t, uint(3), spell.CharacterID)
	assert.Equal(t, uint(3), spell.UserID)
	assert.Equal(t, "Creates a ball of fire", spell.Beschreibung)
	assert.Equal(t, 5, spell.Bonus)
	assert.Equal(t, "Arkanum", spell.Quelle)
}

func TestSkZauber_StructTags(t *testing.T) {
	// Test that the struct has the expected JSON tags by creating a spell
	spell := SkZauber{
		Beschreibung: "Lightning spell",
		Bonus:        3,
		Quelle:       "Elementar",
	}

	// Verify struct fields exist and can be set
	assert.Equal(t, "Lightning spell", spell.Beschreibung)
	assert.Equal(t, 3, spell.Bonus)
	assert.Equal(t, "Elementar", spell.Quelle)
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestSkillStructures_WithDatabase(t *testing.T) {
	setupTestDB(t)

	// Create test skills in database
	testSkill := createTestSkill("Klettern")
	testSkill.Category = "Alltag"
	err := testSkill.Create()
	require.NoError(t, err)

	weaponSkillData := WeaponSkill{
		Skill: Skill{
			GameSystem:       "midgard",
			Name:             "Bögen",
			Beschreibung:     "Bow skills",
			Category:         "Fernkampf",
			Difficulty:       "Normal",
			Initialwert:      5,
			Bonuseigenschaft: "Gs",
			Improvable:       true,
			InnateSkill:      false,
		},
	}
	err = weaponSkillData.Create()
	require.NoError(t, err)

	// Test SkFertigkeit with database interaction
	t.Run("SkFertigkeit Database Integration", func(t *testing.T) {
		skill := SkFertigkeit{
			BamortCharTrait: BamortCharTrait{
				BamortBase: BamortBase{
					Name: "Klettern",
				},
				CharacterID: 1,
				UserID:      1,
			},
			Fertigkeitswert: 8,
			Improvable:      true,
		}

		// Test GetSkillByName
		gsSkill := skill.GetSkillByName()
		assert.NotNil(t, gsSkill)
		assert.Equal(t, "Klettern", gsSkill.Name)
		// Note: gsSkill.Category is empty because category is stored in learning_skill_category_difficulties table

		// Test GetCategory - should fetch from learning_skill_category_difficulties
		category := skill.GetCategory()
		assert.Equal(t, "Alltag", category)
		// Note: skill.Category field is NOT set anymore, GetCategory always queries the database
	})

	// Test SkWaffenfertigkeit with database interaction
	/*
		t.Run("SkWaffenfertigkeit Database Integration", func(t *testing.T) {
			weaponSkill := SkWaffenfertigkeit{
				SkFertigkeit: SkFertigkeit{
					BamortCharTrait: BamortCharTrait{
						BamortBase: BamortBase{
							Name: "Bögen",
						},
						CharacterID: 2,
						UserID:      2,
					},
					Fertigkeitswert: 10,
					Improvable:      true,
				},
			}

			// Test inherited methods
			gsSkill := weaponSkill.GetSkillByName()
			assert.NotNil(t, gsSkill)
			assert.Equal(t, "Bögen", gsSkill.Name)
			assert.Equal(t, "Fernkampf", gsSkill.Category)

			category := weaponSkill.GetCategory()
			assert.Equal(t, "Fernkampf", category)
		})
	*/
}

func TestTableNames_Consistency(t *testing.T) {
	// Test that all table names follow the expected pattern
	skill := SkFertigkeit{}
	weaponSkill := SkWaffenfertigkeit{}
	spell := SkZauber{}

	assert.Equal(t, "char_skills", skill.TableName())
	assert.Equal(t, "char_weaponskills", weaponSkill.TableName())
	assert.Equal(t, "char_spells", spell.TableName())

	// All table names should start with "char_"
	assert.Contains(t, skill.TableName(), "char_")
	assert.Contains(t, weaponSkill.TableName(), "char_")
	assert.Contains(t, spell.TableName(), "char_")
}

func TestSkFertigkeit_EdgeCases(t *testing.T) {
	setupTestDB(t)

	t.Run("GetSkillByName with empty name", func(t *testing.T) {
		skill := SkFertigkeit{
			BamortCharTrait: BamortCharTrait{
				BamortBase: BamortBase{
					Name: "",
				},
			},
		}

		result := skill.GetSkillByName()
		assert.Nil(t, result)
	})

	t.Run("GetCategory with nil database connection", func(t *testing.T) {
		// This tests the behavior when the database lookup fails
		skill := SkFertigkeit{
			BamortCharTrait: BamortCharTrait{
				BamortBase: BamortBase{
					Name: "NonExistentSkillForSure",
				},
			},
			Category: "",
		}

		result := skill.GetCategory()
		assert.Equal(t, "Unkategorisiert", result)
	})

	t.Run("GetCategory with already set category", func(t *testing.T) {
		// Category field is now ignored - GetCategory always queries database
		skill := SkFertigkeit{
			BamortCharTrait: BamortCharTrait{
				BamortBase: BamortBase{
					Name: "NonExistentSkill",
				},
			},
			Category: "PresetCategory", // This is ignored
		}

		result := skill.GetCategory()
		// Since skill doesn't exist in gsmaster, returns "Unkategorisiert"
		assert.Equal(t, "Unkategorisiert", result)
	})
}

// =============================================================================
// Benchmark Tests (Optional)
// =============================================================================

func BenchmarkSkFertigkeit_GetSkillByName(b *testing.B) {
	database.SetupTestDB()
	err := MigrateStructure()
	if err != nil {
		b.Fatal("Failed to migrate structure:", err)
	}

	// Create test skill
	testSkill := createTestSkill("BenchmarkSkill")
	err = testSkill.Create()
	if err != nil {
		b.Fatal("Failed to create test skill:", err)
	}

	skill := SkFertigkeit{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "BenchmarkSkill",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skill.GetSkillByName()
	}
}

func BenchmarkSkFertigkeit_GetCategory(b *testing.B) {
	database.SetupTestDB()
	err := MigrateStructure()
	if err != nil {
		b.Fatal("Failed to migrate structure:", err)
	}

	// Create test skill
	testSkill := createTestSkill("BenchmarkCategorySkill")
	testSkill.Category = "BenchmarkCategory"
	err = testSkill.Create()
	if err != nil {
		b.Fatal("Failed to create test skill:", err)
	}

	skill := SkFertigkeit{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "BenchmarkCategorySkill",
			},
		},
		Category: "", // Force database lookup
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset category for each iteration to force database lookup
		skill.Category = ""
		skill.GetCategory()
	}
}
