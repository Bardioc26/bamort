package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"testing"
)

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}

// TestGetSkillCategory tests the GetSkillCategory function
func TestGetSkillCategory(t *testing.T) {
	tests := []struct {
		name      string
		skillName string
		expected  []string // Allow multiple valid categories for skills that appear in multiple places
	}{
		{
			name:      "Skill in multiple categories",
			skillName: "Klettern", // appears in Alltag, Halbwelt, and Körper
			expected:  []string{"Alltag", "Halbwelt", "Körper"},
		},
		{
			name:      "Skill in Freiland category",
			skillName: "Überleben",
			expected:  []string{"Freiland"},
		},
		{
			name:      "Skill in Waffen category",
			skillName: "Stichwaffen",
			expected:  []string{"Waffen"},
		},
		{
			name:      "Skill in Wissen category",
			skillName: "Alchimie",
			expected:  []string{"Wissen"},
		},
		{
			name:      "Skill unique to one category",
			skillName: "Gerätekunde", // only in Alltag sehr schwer
			expected:  []string{"Alltag"},
		},
		{
			name:      "Non-existent skill",
			skillName: "NichtExistierendeFertigkeit",
			expected:  []string{"Unbekannt"},
		},
		{
			name:      "Empty skill name",
			skillName: "",
			expected:  []string{"Unbekannt"},
		},
		{
			name:      "Case sensitive test",
			skillName: "klettern",            // lowercase
			expected:  []string{"Unbekannt"}, // should not match "Klettern"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSkillCategory(tt.skillName)

			// Check if result is in the list of expected values
			found := false
			for _, expected := range tt.expected {
				if result == expected {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("GetSkillCategory(%q) = %q, want one of %v", tt.skillName, result, tt.expected)
			}
		})
	}
}

// TestGetSkillDifficulty tests the GetSkillDifficulty function
func TestGetSkillDifficulty(t *testing.T) {
	tests := []struct {
		name      string
		category  string
		skillName string
		expected  string
	}{
		{
			name:      "Skill with specific category",
			category:  "Alltag",
			skillName: "Klettern",
			expected:  "leicht",
		},
		{
			name:      "Skill with specific category - normal difficulty",
			category:  "Alltag",
			skillName: "Schreiben",
			expected:  "normal",
		},
		{
			name:      "Skill with specific category - schwer difficulty",
			category:  "Alltag",
			skillName: "Erste Hilfe",
			expected:  "schwer",
		},
		{
			name:      "Skill with specific category - sehr schwer difficulty",
			category:  "Alltag",
			skillName: "Gerätekunde",
			expected:  "sehr schwer",
		},
		{
			name:      "Skill without category - should return first occurrence",
			category:  "",
			skillName: "Klettern",
			expected:  "leicht", // appears as leicht in all categories where it exists
		},
		{
			name:      "Skill without category - another skill",
			category:  "",
			skillName: "Überleben",
			expected:  "leicht", // in Freiland
		},
		{
			name:      "Skill in wrong category",
			category:  "Waffen",
			skillName: "Klettern", // Klettern is not in Waffen category
			expected:  "Unbekannt",
		},
		{
			name:      "Non-existent category",
			category:  "NichtExistierendeKategorie",
			skillName: "Klettern",
			expected:  "Unbekannt",
		},
		{
			name:      "Non-existent skill with valid category",
			category:  "Alltag",
			skillName: "NichtExistierendeFertigkeit",
			expected:  "Unbekannt",
		},
		{
			name:      "Non-existent skill without category",
			category:  "",
			skillName: "NichtExistierendeFertigkeit",
			expected:  "Unbekannt",
		},
		{
			name:      "Empty skill name with category",
			category:  "Alltag",
			skillName: "",
			expected:  "Unbekannt",
		},
		{
			name:      "Empty skill name without category",
			category:  "",
			skillName: "",
			expected:  "Unbekannt",
		},
		{
			name:      "Skill that appears in multiple categories - specific category",
			category:  "Halbwelt",
			skillName: "Klettern", // also exists in Alltag and Körper
			expected:  "leicht",
		},
		{
			name:      "Skill in Freiland schwer",
			category:  "Freiland",
			skillName: "Schleichen",
			expected:  "schwer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSkillDifficulty(tt.category, tt.skillName)
			if result != tt.expected {
				t.Errorf("GetSkillDifficulty(%q, %q) = %q, want %q", tt.category, tt.skillName, result, tt.expected)
			}
		})
	}
}

// TestContains tests the contains helper function
func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "Item exists in slice",
			slice:    []string{"Klettern", "Reiten", "Seilkunst"},
			item:     "Klettern",
			expected: true,
		},
		{
			name:     "Item does not exist in slice",
			slice:    []string{"Klettern", "Reiten", "Seilkunst"},
			item:     "Schwimmen",
			expected: false,
		},
		{
			name:     "Empty slice",
			slice:    []string{},
			item:     "Klettern",
			expected: false,
		},
		{
			name:     "Empty item in slice",
			slice:    []string{"", "Klettern", "Reiten"},
			item:     "",
			expected: true,
		},
		{
			name:     "Empty item not in slice",
			slice:    []string{"Klettern", "Reiten"},
			item:     "",
			expected: false,
		},
		{
			name:     "Case sensitive check",
			slice:    []string{"Klettern", "Reiten"},
			item:     "klettern",
			expected: false,
		},
		{
			name:     "Single item slice - match",
			slice:    []string{"Klettern"},
			item:     "Klettern",
			expected: true,
		},
		{
			name:     "Single item slice - no match",
			slice:    []string{"Klettern"},
			item:     "Reiten",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("contains(%v, %q) = %v, want %v", tt.slice, tt.item, result, tt.expected)
			}
		})
	}
}

// TestCalcSkillLernCost tests the CalcSkillLernCost function
func TestCalcSkillLernCost(t *testing.T) {
	tests := []struct {
		name         string
		costResult   *SkillCostResultNew
		expectError  bool
		expectedLE   int
		expectedEP   int
		expectedGold int
	}{
		{
			name: "Valid calculation for Assassine Alltag leicht",
			costResult: &SkillCostResultNew{
				CharacterClass: "As",
				Category:       "Alltag",
				Difficulty:     "leicht",
			},
			expectError:  false,
			expectedLE:   1,   // LearnCost for leicht in Alltag
			expectedEP:   60,  // 20 (EP per TE for As/Alltag) * 1 (LE) * 3
			expectedGold: 200, // 1 (LE) * 200
		},
		{
			name: "Valid calculation for Krieger Waffen schwer",
			costResult: &SkillCostResultNew{
				CharacterClass: "Kr",
				Category:       "Waffen",
				Difficulty:     "schwer",
			},
			expectError:  false,
			expectedLE:   6,    // LearnCost for schwer in Waffen
			expectedEP:   180,  // 10 (EP per TE for Kr/Waffen) * 6 (LE) * 3
			expectedGold: 1200, // 6 (LE) * 200
		},
		{
			name: "Valid calculation for Magier Wissen normal",
			costResult: &SkillCostResultNew{
				CharacterClass: "Ma",
				Category:       "Wissen",
				Difficulty:     "normal",
			},
			expectError:  false,
			expectedLE:   2,   // LearnCost for normal in Wissen
			expectedEP:   60,  // 10 (EP per TE for Ma/Wissen) * 2 (LE) * 3
			expectedGold: 400, // 2 (LE) * 200
		},
		{
			name: "Invalid character class",
			costResult: &SkillCostResultNew{
				CharacterClass: "InvalidClass",
				Category:       "Alltag",
				Difficulty:     "leicht",
			},
			expectError: true,
		},
		{
			name: "Invalid category",
			costResult: &SkillCostResultNew{
				CharacterClass: "As",
				Category:       "InvalidCategory",
				Difficulty:     "leicht",
			},
			expectError: true,
		},
		{
			name: "Invalid difficulty",
			costResult: &SkillCostResultNew{
				CharacterClass: "As",
				Category:       "Alltag",
				Difficulty:     "InvalidDifficulty",
			},
			expectError: true,
		},
		{
			name: "Valid but category not in character class",
			costResult: &SkillCostResultNew{
				CharacterClass: "As",
				Category:       "Schilde und Parierwaﬀen", // This category might not have EP costs for As
				Difficulty:     "normal",
			},
			expectError: true, // Should fail because EP costs not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CalcSkillLernCost(tt.costResult, nil) // nil reward for original tests

			if tt.expectError {
				if err == nil {
					t.Errorf("CalcSkillLernCost() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("CalcSkillLernCost() unexpected error: %v", err)
				return
			}

			if tt.costResult.LE != tt.expectedLE {
				t.Errorf("CalcSkillLernCost() LE = %d, want %d", tt.costResult.LE, tt.expectedLE)
			}

			if tt.costResult.EP != tt.expectedEP {
				t.Errorf("CalcSkillLernCost() EP = %d, want %d", tt.costResult.EP, tt.expectedEP)
			}

			if tt.costResult.GoldCost != tt.expectedGold {
				t.Errorf("CalcSkillLernCost() GoldCost = %d, want %d", tt.costResult.GoldCost, tt.expectedGold)
			}
		})
	}
}

// TestLearningCostsDataIntegrity tests the integrity of the learning costs data structure
func TestLearningCostsDataIntegrity(t *testing.T) {
	// Test that learningCostsData is not nil
	if learningCostsData == nil {
		t.Fatal("learningCostsData is nil")
	}

	// Test that main maps are not nil
	if learningCostsData.EPPerTE == nil {
		t.Error("EPPerTE map is nil")
	}

	if learningCostsData.SpellEPPerLE == nil {
		t.Error("SpellEPPerLE map is nil")
	}

	if learningCostsData.ImprovementCost == nil {
		t.Error("ImprovementCost map is nil")
	}

	// Test that we have data for expected character classes
	expectedClasses := []string{"As", "Bb", "Gl", "Hä", "Kr", "Sp", "Wa", "Ba", "Or", "Dr", "Hx", "Ma", "PB", "PS", "Sc"}
	for _, class := range expectedClasses {
		if _, exists := learningCostsData.EPPerTE[class]; !exists {
			t.Errorf("Missing EPPerTE data for character class: %s", class)
		}
	}

	// Test that we have data for expected categories
	expectedCategories := []string{"Alltag", "Freiland", "Halbwelt", "Kampf", "Körper", "Sozial", "Unterwelt", "Waffen", "Wissen"}
	for _, category := range expectedCategories {
		if _, exists := learningCostsData.ImprovementCost[category]; !exists {
			t.Errorf("Missing ImprovementCost data for category: %s", category)
		}
	}

	// Test that difficulty levels exist where expected
	for category, difficulties := range learningCostsData.ImprovementCost {
		if len(difficulties) == 0 {
			t.Errorf("Category %s has no difficulty levels", category)
		}

		// Check that each difficulty has proper structure
		for difficulty, data := range difficulties {
			if data.Skills == nil {
				t.Errorf("Category %s, difficulty %s has nil Skills slice", category, difficulty)
			}
			if data.TrainCosts == nil {
				t.Errorf("Category %s, difficulty %s has nil TrainCosts map", category, difficulty)
			}
		}
	}
}

// TestSkillCoverage tests that all skills in the data structure can be found by the functions
func TestSkillCoverage(t *testing.T) {
	skillsFound := make(map[string]bool)

	// Collect all skills from the data structure
	for category, difficulties := range learningCostsData.ImprovementCost {
		for difficulty, data := range difficulties {
			for _, skill := range data.Skills {
				if skill != "" { // Skip empty skill names
					skillsFound[skill] = false

					// Test that GetSkillCategory can find this skill
					foundCategory := GetSkillCategory(skill)
					if foundCategory == "Unbekannt" {
						t.Errorf("GetSkillCategory could not find skill: %s (should be in %s)", skill, category)
					}

					// Test that GetSkillDifficulty can find this skill without category
					foundDifficulty := GetSkillDifficulty("", skill)
					if foundDifficulty == "Unbekannt" {
						t.Errorf("GetSkillDifficulty could not find skill: %s (should have difficulty %s)", skill, difficulty)
					}

					// Test that GetSkillDifficulty can find this skill with category
					foundDifficultyWithCategory := GetSkillDifficulty(category, skill)
					if foundDifficultyWithCategory == "Unbekannt" {
						t.Errorf("GetSkillDifficulty could not find skill: %s in category %s (should have difficulty %s)", skill, category, difficulty)
					}

					skillsFound[skill] = true
				}
			}
		}
	}

	t.Logf("Tested coverage for %d unique skills", len(skillsFound))
}

// TestFindBestCategoryForSkill tests the findBestCategoryForSkill function
func TestFindBestCategoryForSkill(t *testing.T) {
	tests := []struct {
		name             string
		skillName        string
		characterClass   string
		currentLevel     int
		expectedCategory string
		expectError      bool
	}{
		{
			name:           "Klettern - should choose cheapest category",
			skillName:      "Klettern",
			characterClass: "Kr", // Krieger
			currentLevel:   13,   // Level 13->14
			// Klettern ist in: Alltag (leicht), Halbwelt (leicht), Körper (leicht)
			// Für Kr: Alltag=20 EP/TE, Halbwelt=30 EP/TE, Körper=20 EP/TE
			// Level 13->14 kostet in allen 1 TE, also 20*1=20 EP für Alltag und Körper, 30*1=30 EP für Halbwelt
			// Sollte Alltag oder Körper wählen (beide gleich günstig)
			expectedCategory: "Alltag", // oder "Körper" - beide sind gleich günstig
			expectError:      false,
		},
		{
			name:           "Non-existent skill",
			skillName:      "NichtExistierendeFertigkeit",
			characterClass: "Kr",
			currentLevel:   10,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, difficulty, err := findBestCategoryForSkillImprovement(tt.skillName, tt.characterClass, tt.currentLevel)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Für Klettern sind mehrere Kategorien gleich günstig, also akzeptieren wir alle
			if tt.skillName == "Klettern" {
				validCategories := []string{"Alltag", "Körper"} // Beide haben 20 EP/TE für Kr
				found := false
				for _, validCat := range validCategories {
					if category == validCat {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected category to be one of %v, got %s", validCategories, category)
				}
			} else {
				if category != tt.expectedCategory {
					t.Errorf("Expected category %s, got %s", tt.expectedCategory, category)
				}
			}

			t.Logf("Skill %s for class %s at level %d: category=%s, difficulty=%s",
				tt.skillName, tt.characterClass, tt.currentLevel, category, difficulty)
		})
	}
}

// TestCalcSkillLernCostWithRewards tests the reward logic in CalcSkillLernCost
func TestCalcSkillLernCostWithRewards(t *testing.T) {
	tests := []struct {
		name           string
		skillName      string
		characterClass string
		reward         *string
		expectedGold   int
		expectedEPMult float64 // multiplier for EP (1.0 = normal, 0.5 = half)
	}{
		{
			name:           "Default reward - normal costs",
			skillName:      "Klettern",
			characterClass: "Kr", // Use abbreviation
			reward:         stringPtr("default"),
			expectedGold:   200, // 1 LE * 200 Gold per LE
			expectedEPMult: 1.0,
		},
		{
			name:           "NoGold reward - no gold cost",
			skillName:      "Klettern",
			characterClass: "Kr", // Use abbreviation
			reward:         stringPtr("noGold"),
			expectedGold:   0, // Should be 0 with noGold reward
			expectedEPMult: 1.0,
		},
		{
			name:           "No reward - normal costs",
			skillName:      "Klettern",
			characterClass: "Kr", // Use abbreviation
			reward:         nil,
			expectedGold:   200, // 1 LE * 200 Gold per LE
			expectedEPMult: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create cost result
			costResult := &SkillCostResultNew{
				CharacterClass: tt.characterClass,
				SkillName:      tt.skillName,
				Category:       GetSkillCategory(tt.skillName),
				Difficulty:     GetSkillDifficulty(GetSkillCategory(tt.skillName), tt.skillName),
			}

			// Calculate normal costs first to get baseline EP
			baselineResult := &SkillCostResultNew{
				CharacterClass: tt.characterClass,
				SkillName:      tt.skillName,
				Category:       costResult.Category,
				Difficulty:     costResult.Difficulty,
			}
			err := CalcSkillLernCost(baselineResult, stringPtr("default"))
			if err != nil {
				t.Fatalf("Failed to calculate baseline costs: %v", err)
			}

			// Calculate costs with reward
			err = CalcSkillLernCost(costResult, tt.reward)
			if err != nil {
				t.Fatalf("Failed to calculate costs: %v", err)
			}

			// Check gold cost
			if costResult.GoldCost != tt.expectedGold {
				t.Errorf("Expected gold cost %d, got %d", tt.expectedGold, costResult.GoldCost)
			}

			// Check EP cost
			expectedEP := int(float64(baselineResult.EP) * tt.expectedEPMult)
			if costResult.EP != expectedEP {
				t.Errorf("Expected EP %d (baseline %d * %.1f), got %d", expectedEP, baselineResult.EP, tt.expectedEPMult, costResult.EP)
			}

			// LE should always be the same regardless of reward
			if costResult.LE != baselineResult.LE {
				t.Errorf("LE should be unchanged by rewards. Expected %d, got %d", baselineResult.LE, costResult.LE)
			}
		})
	}
}

// TestCalcSpellLernCostWithRewards tests the reward logic in CalcSpellLernCost
/*
func TestCalcSpellLernCostWithRewards(t *testing.T) {
	costResult := &SkillCostResultNew{
		CharacterClass: "Ma", // Use abbreviation
		SkillName:      "TestSpell",
		Category:       "Hellsicht", // Use existing category
		Difficulty:     "Schwer",
	}

	// Test with noGold reward
	err := CalcSpellLernCost(costResult, stringPtr("noGold"))
	if err != nil {
		t.Fatalf("Failed to calculate spell costs: %v", err)
	}

	if costResult.GoldCost != 0 {
		t.Errorf("Expected gold cost 0 with noGold reward, got %d", costResult.GoldCost)
	}
}
*/

// TestCalcSkillImproveCostWithRewards tests the reward logic in CalcSkillImproveCost
func TestCalcSkillImproveCostWithRewards(t *testing.T) {
	tests := []struct {
		name           string
		skillName      string
		characterClass string
		currentLevel   int // represents the level the character currently has must be incremented by 1 when calculating the costs
		ppUsed         int
		reward         *string
		expectedEP     int
		expectedGold   int
	}{
		{
			name:           "Normal improvement to 13 without reward",
			skillName:      "Klettern",
			characterClass: "Kr",
			currentLevel:   12,
			ppUsed:         0,
			reward:         nil,
			expectedEP:     20, // Kr has 20 EP/TE for Alltag, level 12->13 costs 0 TE, so 20*0=0
			expectedGold:   20, // 0 TE * 20 Gold per TE
		},
		{
			name:           "Normal improvement to 14 without reward",
			skillName:      "Klettern",
			characterClass: "Kr",
			currentLevel:   13,
			ppUsed:         0,
			reward:         nil,
			expectedEP:     40, // Kr has 20 EP/TE for Alltag, level 13->14 costs 1 TE, so 20*1=20
			expectedGold:   40, // 1 TE * 20 Gold per TE
		},
		{
			name:           "Improvement with halveep reward",
			skillName:      "Klettern",
			characterClass: "Kr",
			currentLevel:   13,
			ppUsed:         0,
			reward:         stringPtr("halveep"),
			expectedEP:     20, // Kr has 20 EP/TE for Alltag, level 13->14 costs 1 TE, so 20*1=20, halved = 10
			expectedGold:   40, // Gold cost not affected by halveep
		},

		{
			name:           "Improvement to 15 without reward",
			skillName:      "Klettern",
			characterClass: "Kr",
			currentLevel:   14,
			ppUsed:         0,
			reward:         nil,
			expectedEP:     100, // Kr has 20 EP/TE for Alltag, level 14->15 costs 2 TE, minus 1 PP = 1 TE, so 20*1=20
			expectedGold:   100, // 1 TE * 20 Gold per TE
		},
		{
			name:           "Improvement to 15 with PP used",
			skillName:      "Klettern",
			characterClass: "Kr",
			currentLevel:   14,
			ppUsed:         1,
			reward:         nil,
			expectedEP:     80, // Kr has 20 EP/TE for Alltag, level 14->15 costs 2 TE, minus 1 PP = 1 TE, so 20*1=20
			expectedGold:   80, // 1 TE * 20 Gold per TE
		},
		{
			name:           "Improvement with halveepnoGold reward",
			skillName:      "Klettern",
			characterClass: "Kr",
			currentLevel:   15,
			ppUsed:         0,
			reward:         stringPtr("halveepnoGold"),
			expectedEP:     100, // Kr has 20 EP/TE for Alltag, level 15->16 costs 5 TE, so 20*5=100, halved = 50
			expectedGold:   0,   // Should be 0 with halveepnoGold reward
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			costResult := &SkillCostResultNew{
				CharacterClass: tt.characterClass,
				SkillName:      tt.skillName,
				PPUsed:         tt.ppUsed,
				// Lassen Sie Kategorie und Schwierigkeit leer, damit die Funktion die beste auswählt
			}

			err := CalcSkillImproveCost(costResult, tt.currentLevel, tt.reward)
			if err != nil {
				t.Fatalf("Failed to calculate improvement costs: %v", err)
			}

			// Log the chosen category for debugging
			t.Logf("Skill: %s, Class %s, Chosen category: %s, difficulty: %s", costResult.SkillName, costResult.CharacterClass, costResult.Category, costResult.Difficulty)

			if costResult.EP != tt.expectedEP {
				t.Errorf("Expected EP %d, got %d", tt.expectedEP, costResult.EP)
			}

			if costResult.GoldCost != tt.expectedGold {
				t.Errorf("Expected gold cost %d, got %d", tt.expectedGold, costResult.GoldCost)
			}
		})
	}
}

// TestGetSpellInfo tests the GetSpellInfo function
func TestGetSpellInfo(t *testing.T) {

	// Initialize test database with migration (but no test data since we don't have the preparedTestDB file)
	database.SetupTestDB(true, false) // Use in-memory SQLite, no test data loading
	defer database.ResetTestDB()
	MigrateStructure()

	// Create minimal test spell data for our test
	testSpells := []models.Spell{
		{
			LookupList: models.LookupList{
				GameSystem:   "midgard",
				Name:         "Schlummer",
				Beschreibung: "Test spell for GetSpellInfo",
				Quelle:       "Test",
			},
			Stufe:    1,
			Category: "Beherrschen",
		},
		{
			LookupList: models.LookupList{
				GameSystem:   "midgard",
				Name:         "Erkennen von Krankheit",
				Beschreibung: "Test spell for GetSpellInfo",
				Quelle:       "Test",
			},
			Stufe:    2,
			Category: "Dweomerzauber",
		},
		{
			LookupList: models.LookupList{
				GameSystem:   "midgard",
				Name:         "Das Loblied",
				Beschreibung: "Test spell for GetSpellInfo",
				Quelle:       "Test",
			},
			Stufe:    3,
			Category: "Zauberlied",
		},
	}

	// Insert test data directly
	for _, spell := range testSpells {
		if err := database.DB.Create(&spell).Error; err != nil {
			t.Fatalf("Failed to create test spell: %v", err)
		}
	}

	tests := []struct {
		spellName      string
		expectedSchool string
		expectedLevel  int
		expectError    bool
	}{
		{
			spellName:      "Schlummer",
			expectedSchool: "Beherrschen",
			expectedLevel:  1,
			expectError:    false,
		},
		{
			spellName:      "Erkennen von Krankheit",
			expectedSchool: "Dweomerzauber",
			expectedLevel:  2,
			expectError:    false,
		},
		{
			spellName:      "Das Loblied",
			expectedSchool: "Zauberlied",
			expectedLevel:  3,
			expectError:    false,
		},
		{
			spellName:      "Unknown Spell",
			expectedSchool: "", // Should error for unknown spell
			expectedLevel:  0,  // Should error for unknown spell
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.spellName, func(t *testing.T) {
			school, level, err := GetSpellInfo(tt.spellName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for unknown spell, but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Failed to get spell info: %v", err)
			}

			if school != tt.expectedSchool {
				t.Errorf("Expected school %s, got %s", tt.expectedSchool, school)
			}

			if level != tt.expectedLevel {
				t.Errorf("Expected level %d, got %d", tt.expectedLevel, level)
			}
		})
	}
}

// TestCalcSpellLernCostWithRewards tests the reward logic in CalcSpellLernCost
func TestCalcSpellLernCostWithRewards(t *testing.T) {
	// Initialize test database with migration (but no test data since we don't have the preparedTestDB file)
	database.SetupTestDB(true, false) // Use in-memory SQLite, no test data loading
	defer database.ResetTestDB()
	MigrateStructure()

	// Create minimal test spell data for our test
	testSpells := []models.Spell{
		{
			LookupList: models.LookupList{
				GameSystem:   "midgard",
				Name:         "Schlummer",
				Beschreibung: "Test spell for GetSpellInfo",
				Quelle:       "Test",
			},
			Stufe:    1,
			Category: "Beherrschen",
		},
		{
			LookupList: models.LookupList{
				GameSystem:   "midgard",
				Name:         "Erkennen von Krankheit",
				Beschreibung: "Test spell for GetSpellInfo",
				Quelle:       "Test",
			},
			Stufe:    2,
			Category: "Dweomer",
		},
		{
			LookupList: models.LookupList{
				GameSystem:   "midgard",
				Name:         "Das Loblied",
				Beschreibung: "Test spell for GetSpellInfo",
				Quelle:       "Test",
			},
			Stufe:    3,
			Category: "Zauberlied",
		},
	}
	// Insert test data directly
	for _, spell := range testSpells {
		if err := database.DB.Create(&spell).Error; err != nil {
			t.Fatalf("Failed to create test spell: %v", err)
		}
	}

	tests := []struct {
		name           string
		spellName      string
		characterClass string
		reward         *string
		expectedEP     int
		expectedGold   int
	}{
		{
			name:           "Simple spell for Magier without but specialized",
			spellName:      "Schlummer",
			characterClass: "Ma",
			reward:         nil,
			expectedEP:     30,  // Ma has 60 EP/LE for Beherrschen, Furcht is level 1 = 1 LE, so 1*60=60
			expectedGold:   100, // 1 LE * 100 Gold per LE
		},
		{
			name:           "Spell with spruchrolle no reward",
			spellName:      "Erkennen von Krankheit",
			characterClass: "Ma",
			reward:         nil,
			expectedEP:     120, // 60/3 for spruchrolle
			expectedGold:   100, // Fixed 20 Gold for spruchrolle
		},
		{
			name:           "Spell with spruchrolle reward",
			spellName:      "Erkennen von Krankheit",
			characterClass: "Ma",
			reward:         stringPtr("spruchrolle"),
			expectedEP:     40, // 60/3 for spruchrolle
			expectedGold:   20, // Fixed 20 Gold for spruchrolle
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			costResult := &SkillCostResultNew{
				CharacterClass: tt.characterClass,
				SkillName:      tt.spellName,
				CharacterID:    "test-character",
			}

			err := CalcSpellLernCost(costResult, tt.reward)
			if err != nil {
				t.Fatalf("Failed to calculate spell costs: %v", err)
			}

			if costResult.EP != tt.expectedEP {
				t.Errorf("Expected EP %d, got %d", tt.expectedEP, costResult.EP)
			}

			if costResult.GoldCost != tt.expectedGold {
				t.Errorf("Expected gold cost %d, got %d", tt.expectedGold, costResult.GoldCost)
			}
		})
	}
}

// TestGetSpecialization tests the GetSpecialization function
func TestGetSpecialization(t *testing.T) {
	tests := []struct {
		name         string
		characterID  string
		expectedSpec string
	}{
		{
			name:         "Default specialization",
			characterID:  "123",
			expectedSpec: "Beherrschen",
		},
		{
			name:         "Another character",
			characterID:  "456",
			expectedSpec: "Beherrschen", // Currently returns default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSpecialization(tt.characterID)
			if result != tt.expectedSpec {
				t.Errorf("Expected specialization %s, got %s", tt.expectedSpec, result)
			}
		})
	}
}

// TestFindBestCategoryForSkillLearning tests the findBestCategoryForSkillLearning function
func TestFindBestCategoryForSkillLearning(t *testing.T) {
	tests := []struct {
		name           string
		skillName      string
		characterClass string
		expectedCat    string
		expectedDiff   string
		expectError    bool
	}{
		{
			name:           "Klettern for Assassine - should find best category",
			skillName:      "Klettern",
			characterClass: "As",
			expectedCat:    "Körper", // Should prefer Körper (10 EP/TE * 1 LE * 3 = 30 EP) over Alltag (20 EP/TE * 1 LE * 3 = 60 EP)
			expectedDiff:   "leicht",
			expectError:    false,
		},
		{
			name:           "Schleichen for Spitzbube - should find Unterwelt",
			skillName:      "Schleichen",
			characterClass: "Sp",
			expectedCat:    "Unterwelt", // Sp has 10 EP/TE for Unterwelt vs 30 EP/TE for Freiland
			expectedDiff:   "normal",
			expectError:    false,
		},
		{
			name:           "Invalid skill",
			skillName:      "NonExistentSkill",
			characterClass: "As",
			expectError:    true,
		},
		{
			name:           "Invalid character class",
			skillName:      "Klettern",
			characterClass: "InvalidClass",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category, difficulty, err := findBestCategoryForSkillLearning(tt.skillName, tt.characterClass)

			if tt.expectError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if category != tt.expectedCat {
				t.Errorf("Expected category %s, got %s", tt.expectedCat, category)
			}

			if difficulty != tt.expectedDiff {
				t.Errorf("Expected difficulty %s, got %s", tt.expectedDiff, difficulty)
			}
		})
	}
}

// TestGetLernCostNextLevel tests the GetLernCostNextLevel function
func TestGetLernCostNextLevel(t *testing.T) {
	tests := []struct {
		name          string
		request       *LernCostRequest
		costResult    *SkillCostResultNew
		reward        *string
		level         int
		characterTyp  string
		expectError   bool
		expectedEP    int
		expectedGold  int
		expectedElfEP int // Expected EP bonus for Elves
	}{
		{
			name: "Learn skill as Human",
			request: &LernCostRequest{
				Action: "learn",
				Type:   "skill",
				Reward: stringPtr("default"),
			},
			costResult: &SkillCostResultNew{
				CharacterClass: "As",
				SkillName:      "Klettern",
				Category:       "Körper",
				Difficulty:     "leicht",
			},
			level:        1,
			characterTyp: "Mensch",
			expectError:  false,
			expectedEP:   30,  // 10 * 1 * 3
			expectedGold: 200, // 1 * 200
		},
		{
			name: "Learn skill as Human Kr",
			request: &LernCostRequest{
				Action: "learn",
				Type:   "skill",
				Reward: stringPtr("default"),
			},
			costResult: &SkillCostResultNew{
				CharacterClass: "Kr",
				SkillName:      "Abrichten",
				Category:       "Körper",
				Difficulty:     "leicht",
			},
			level:        1,
			characterTyp: "Mensch",
			expectError:  true, // TODO Abrichten kommt im Mysterium mit dem Tiermeister
			expectedEP:   30,   // 10 * 1 * 3
			expectedGold: 200,  // 1 * 200
		},
		{
			name: "Learn skill as Elf - should have EP bonus",
			request: &LernCostRequest{
				Action: "learn",
				Type:   "skill",
				Reward: stringPtr("default"),
			},
			costResult: &SkillCostResultNew{
				CharacterClass: "As",
				SkillName:      "Klettern",
				Category:       "Körper",
				Difficulty:     "leicht",
			},
			level:         1,
			characterTyp:  "Elf",
			expectError:   false,
			expectedEP:    30,
			expectedElfEP: 6, // Additional 6 EP for Elves
			expectedGold:  200,
		},
		{
			name: "Improve skill as human",
			request: &LernCostRequest{
				Action:       "improve",
				Type:         "skill",
				CurrentLevel: 12,
				Reward:       stringPtr("default"),
			},
			costResult: &SkillCostResultNew{
				CharacterClass: "As",
				SkillName:      "Klettern",
				Category:       "Körper",
				Difficulty:     "leicht",
			},
			level:        13,
			characterTyp: "Mensch",
			expectError:  false,
			expectedEP:   10, // 10 * 1 (TE cost for level 13)
			expectedGold: 20, // 1 * 20
		},

		{
			name: "Improve skill as Elf",
			request: &LernCostRequest{
				Action:       "improve",
				Type:         "skill",
				CurrentLevel: 12,
				Reward:       stringPtr("default"),
			},
			costResult: &SkillCostResultNew{
				CharacterClass: "As",
				SkillName:      "Klettern",
				Category:       "Körper",
				Difficulty:     "leicht",
			},
			level:        13,
			characterTyp: "Elf",
			expectError:  false,
			expectedEP:   10, // 10 * 1 (TE cost for level 13)
			expectedGold: 20, // 1 * 20
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetLernCostNextLevel(tt.request, tt.costResult, tt.reward, tt.level, tt.characterTyp)

			if tt.expectError {
				if err == nil {
					t.Error("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			expectedTotalEP := tt.expectedEP
			if tt.characterTyp == "Elf" && (tt.request.Action == "learn") {
				expectedTotalEP += tt.expectedElfEP
			}

			if tt.costResult.EP != expectedTotalEP {
				t.Errorf("Expected EP %d, got %d", expectedTotalEP, tt.costResult.EP)
			}

			if tt.costResult.GoldCost != tt.expectedGold {
				t.Errorf("Expected gold cost %d, got %d", tt.expectedGold, tt.costResult.GoldCost)
			}
		})
	}
}
