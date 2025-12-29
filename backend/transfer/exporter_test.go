package transfer

import (
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"os"
	"testing"
)

func setupTestEnvironment(t *testing.T) {
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	database.SetupTestDB(true, true)
	models.MigrateStructure()
	t.Cleanup(func() {
		database.ResetTestDB()
		if original == "" {
			os.Unsetenv("ENVIRONMENT")
		} else {
			os.Setenv("ENVIRONMENT", original)
		}
	})
}

func TestExportCharacter(t *testing.T) {
	setupTestEnvironment(t)

	// Test with character ID 18 (Fanjo Vetrani - exists in test DB)
	characterID := uint(18)

	exportData, err := ExportCharacter(characterID)
	if err != nil {
		t.Fatalf("ExportCharacter failed: %v", err)
	}

	// Verify basic character data
	if exportData.Character.ID != characterID {
		t.Errorf("Expected character ID %d, got %d", characterID, exportData.Character.ID)
	}

	if exportData.Character.Name == "" {
		t.Error("Character name should not be empty")
	}

	// Verify JSON serialization
	jsonData, err := json.Marshal(exportData)
	if err != nil {
		t.Fatalf("Failed to marshal export data to JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON export should not be empty")
	}

	// Verify we can unmarshal back
	var reimported CharacterExport
	err = json.Unmarshal(jsonData, &reimported)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if reimported.Character.ID != characterID {
		t.Errorf("After JSON round-trip, expected character ID %d, got %d", characterID, reimported.Character.ID)
	}

	// Verify sensitive user data is removed
	if exportData.Character.User.PasswordHash != "" {
		t.Error("Password hash should be empty in export")
	}

	if !exportData.Character.User.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be zero time in export")
	}
}

func TestExportCharacterIncludesSkills(t *testing.T) {
	setupTestEnvironment(t)

	characterID := uint(18)
	exportData, err := ExportCharacter(characterID)
	if err != nil {
		t.Fatalf("ExportCharacter failed: %v", err)
	}

	// Character should have some skills
	if len(exportData.Character.Fertigkeiten) == 0 {
		t.Error("Character should have skills")
	}

	// Verify GSM data for skills is included
	if len(exportData.GSMSkills) == 0 {
		t.Error("GSM skills should be included in export")
	}
}

func TestExportCharacterIncludesSpells(t *testing.T) {
	setupTestEnvironment(t)

	characterID := uint(18)
	exportData, err := ExportCharacter(characterID)
	if err != nil {
		t.Fatalf("ExportCharacter failed: %v", err)
	}

	// Verify GSM data for spells is included if character has spells
	if len(exportData.Character.Zauber) > 0 && len(exportData.GSMSpells) == 0 {
		t.Error("If character has spells, GSM spells should be included")
	}
}

func TestExportCharacterIncludesEquipment(t *testing.T) {
	setupTestEnvironment(t)

	characterID := uint(18)
	exportData, err := ExportCharacter(characterID)
	if err != nil {
		t.Fatalf("ExportCharacter failed: %v", err)
	}

	// Verify GSM data for equipment is included if character has weapons
	if len(exportData.Character.Waffen) > 0 && len(exportData.GSMWeapons) == 0 {
		t.Error("If character has weapons, GSM weapons should be included")
	}
}

func TestExportCharacterIncludesLearningData(t *testing.T) {
	setupTestEnvironment(t)

	characterID := uint(18)
	exportData, err := ExportCharacter(characterID)
	if err != nil {
		t.Fatalf("ExportCharacter failed: %v", err)
	}

	// Verify learning data structures exist (they might be empty)
	if exportData.LearningData.Sources == nil {
		t.Error("Learning sources should be initialized")
	}
}

func TestExportCharacterIncludesAuditLog(t *testing.T) {
	setupTestEnvironment(t)

	characterID := uint(18)
	exportData, err := ExportCharacter(characterID)
	if err != nil {
		t.Fatalf("ExportCharacter failed: %v", err)
	}

	// Audit log entries should be included (even if empty)
	if exportData.AuditLogEntries == nil {
		t.Error("Audit log entries should be initialized")
	}
}

func TestExportNonExistentCharacter(t *testing.T) {
	setupTestEnvironment(t)

	// Try to export non-existent character
	_, err := ExportCharacter(uint(999999))
	if err == nil {
		t.Error("Expected error when exporting non-existent character")
	}
}
