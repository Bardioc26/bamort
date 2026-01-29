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

// TestGetSpellInfo tests the GetSpellInfo function
func TestGetSpellInfo(t *testing.T) {

	// Initialize test database with migration (but no test data since we don't have the preparedTestDB file)
	database.SetupTestDB(true, false) // Use in-memory SQLite, no test data loading
	defer database.ResetTestDB()
	models.MigrateStructure()

	// Create minimal test spell data for our test
	testSpells := []models.Spell{
		{
			GameSystemId: 1,
			Name:         "Schlummer",
			Beschreibung: "Test spell for GetSpellInfo",
			Quelle:       "Test",
			Stufe:        1,
			Category:     "Beherrschen",
		},
		{
			GameSystemId: 1,
			Name:         "Erkennen von Krankheit",
			Beschreibung: "Test spell for GetSpellInfo",
			Quelle:       "Test",
			Stufe:        2,
			Category:     "Dweomerzauber",
		},
		{
			GameSystemId: 1,
			Name:         "Das Loblied",
			Beschreibung: "Test spell for GetSpellInfo",
			Quelle:       "Test",
			Stufe:        3,
			Category:     "Zauberlied",
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
			expectedSchool: "Erkennen",
			expectedLevel:  2,
			expectError:    false,
		},
		{
			spellName:      "Das Loblied",
			expectedSchool: "Verändern",
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
			school, level, err := GetSpellInfoNewSystem(tt.spellName)

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
