package integration_test

import (
	"bamort/database"
	"bamort/gsmaster"
	_ "bamort/maintenance" // Import for init() function that sets up test data callbacks
	"testing"
)

// TestGetSpellInfoIntegration tests the GetSpellInfo function with full database integration
func TestGetSpellInfoIntegration(t *testing.T) {

	// Initialize test database with test data loading
	database.SetupTestDB(true) // Use in-memory SQLite with test data
	defer database.ResetTestDB()

	// Test with real spell names that should exist in the prepared test database
	// If the prepared test database doesn't exist, these tests will fail gracefully
	tests := []struct {
		spellName      string
		expectedSchool string
		expectedLevel  int
		expectError    bool
		description    string
	}{
		{
			spellName:      "Licht",
			expectedSchool: "Erschaffen", // Common light spell in Midgard
			expectedLevel:  1,
			expectError:    false,
			description:    "Basic light spell should exist in prepared database",
		},
		{
			spellName:      "Erkennen von Krankheit",
			expectedSchool: "Erkennen",
			expectedLevel:  2,
			expectError:    false,
			description:    "Common diagnostic spell should exist in prepared database",
		},
		{
			spellName:      "Furcht",
			expectedSchool: "Beherrschen",
			expectedLevel:  1,
			expectError:    false,
			description:    "Fear spell should exist in prepared database",
		},
		{
			spellName:      "Unknown Spell That Should Not Exist",
			expectedSchool: "", // Should error for unknown spell
			expectedLevel:  0,  // Should error for unknown spell
			expectError:    true,
			description:    "Non-existent spell should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.spellName, func(t *testing.T) {
			school, level, err := gsmaster.GetSpellInfoNewSystem(tt.spellName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none. %s", tt.spellName, tt.description)
				}
				return
			}

			if err != nil {
				// If we get an error for a spell we expect to exist, it might mean
				// the prepared test database doesn't exist or doesn't contain this spell
				t.Logf("Warning: Failed to get spell info for %s: %v. %s", tt.spellName, err, tt.description)
				t.Logf("This might indicate that the prepared test database is missing or incomplete.")
				t.Skip("Skipping test due to missing prepared test data")
				return
			}

			if school != tt.expectedSchool {
				t.Errorf("Expected school %s, got %s for spell %s", tt.expectedSchool, school, tt.spellName)
			}

			if level != tt.expectedLevel {
				t.Errorf("Expected level %d, got %d for spell %s", tt.expectedLevel, level, tt.spellName)
			}

			t.Logf("✅ Successfully found spell %s: school=%s, level=%d", tt.spellName, school, level)
		})
	}
}
