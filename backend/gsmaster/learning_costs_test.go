package gsmaster

import (
	"testing"

	"bamort/database"
	"bamort/models"
)

func TestGetClassAbbreviationNewSystem(t *testing.T) {
	setupTestEnvironment(t)

	database.SetupTestDB(true)
	defer database.ResetTestDB()

	defaultGS := GetGameSystem(1, "")

	source := models.Source{Code: "TST", Name: "Test Source", FullName: "Test Source", GameSystem: defaultGS.Name, GameSystemId: defaultGS.ID}
	if err := database.DB.Create(&source).Error; err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	defaultClass := models.CharacterClass{Code: "TC", Name: "Test Class", SourceID: source.ID, GameSystem: defaultGS.Name, GameSystemId: defaultGS.ID}
	if err := database.DB.Create(&defaultClass).Error; err != nil {
		t.Fatalf("failed to create default game system class: %v", err)
	}

	altGS := models.GameSystem{Code: "ALT", Name: "Alternate"}
	if err := database.DB.Create(&altGS).Error; err != nil {
		t.Fatalf("failed to create alternate game system: %v", err)
	}

	altClass := models.CharacterClass{Code: "AC", Name: "Alternate Class", SourceID: source.ID, GameSystem: altGS.Name, GameSystemId: altGS.ID}
	if err := database.DB.Create(&altClass).Error; err != nil {
		t.Fatalf("failed to create alternate game system class: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "returns code for existing class", input: defaultClass.Name, expected: defaultClass.Code},
		{name: "returns empty for unknown class", input: "Unknown Class", expected: ""},
		{name: "returns empty for class in other game system", input: altClass.Name, expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := GetClassAbbreviationNewSystem(tt.input); result != tt.expected {
				t.Fatalf("GetClassAbbreviationNewSystem(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
