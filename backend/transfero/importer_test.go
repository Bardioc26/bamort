package transfero

import (
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"os"
	"testing"
)

func setupImportTestEnvironment(t *testing.T) {
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

func TestImportCharacter(t *testing.T) {
	setupImportTestEnvironment(t)

	// First export character 18
	exportData, err := ExportCharacter(uint(18))
	if err != nil {
		t.Fatalf("Failed to export character: %v", err)
	}

	// Change IDs to simulate import into new system
	originalID := exportData.Character.ID
	exportData.Character.ID = 0 // Reset ID for new character
	exportData.Character.Name = "Imported " + exportData.Character.Name

	// Import the character
	importedCharID, err := ImportCharacter(exportData, 1) // UserID 1
	if err != nil {
		t.Fatalf("ImportCharacter failed: %v", err)
	}

	if importedCharID == 0 {
		t.Error("Expected non-zero character ID after import")
	}

	if importedCharID == originalID {
		t.Error("Imported character should have new ID, not original ID")
	}

	// Verify imported character exists
	var importedChar models.Char
	err = database.DB.Preload("Fertigkeiten").Preload("Zauber").Preload("Waffen").First(&importedChar, importedCharID).Error
	if err != nil {
		t.Fatalf("Failed to load imported character: %v", err)
	}

	// Verify skills were imported
	if len(importedChar.Fertigkeiten) != len(exportData.Character.Fertigkeiten) {
		t.Errorf("Expected %d skills, got %d", len(exportData.Character.Fertigkeiten), len(importedChar.Fertigkeiten))
	}
}

func TestImportCharacterWithExistingGSMData(t *testing.T) {
	setupImportTestEnvironment(t)

	// Export character 18
	exportData, err := ExportCharacter(uint(18))
	if err != nil {
		t.Fatalf("Failed to export character: %v", err)
	}

	// Count existing skills in GSM before import
	var skillCountBefore int64
	database.DB.Model(&models.Skill{}).Count(&skillCountBefore)

	exportData.Character.ID = 0
	exportData.Character.Name = "Test Import"

	// Import character
	_, err = ImportCharacter(exportData, 1)
	if err != nil {
		t.Fatalf("ImportCharacter failed: %v", err)
	}

	// Count skills after - should be same (no duplicates)
	var skillCountAfter int64
	database.DB.Model(&models.Skill{}).Count(&skillCountAfter)

	if skillCountAfter < skillCountBefore {
		t.Error("Skills should not be deleted during import")
	}
}

func TestImportCharacterUpdatesIncompleteGSMData(t *testing.T) {
	setupImportTestEnvironment(t)

	// Create a skill with incomplete data (no description)
	incompleteSkill := models.Skill{
		ID:           1000,
		Name:         "TestSkillIncomplete",
		GameSystemId: 1,
		Beschreibung: "",
		Category:     "",
	}
	database.DB.Create(&incompleteSkill)

	// Create export data with complete version of same skill
	exportData := &CharacterExport{
		Character: models.Char{
			BamortBase: models.BamortBase{
				Name: "Test Char",
			},
			Typ:    "Krieger",
			Rasse:  "Mensch",
			Grad:   1,
			UserID: 1,
		},
		GSMSkills: []models.Skill{
			{
				Name:         "TestSkillIncomplete",
				GameSystemId: 1,
				Beschreibung: "Complete description",
				Category:     "Alltag",
				Difficulty:   "normal",
			},
		},
	}

	_, err := ImportCharacter(exportData, 1)
	if err != nil {
		t.Fatalf("ImportCharacter failed: %v", err)
	}

	// Verify skill was updated
	var updatedSkill models.Skill
	database.DB.Where("name = ?", "TestSkillIncomplete").First(&updatedSkill)

	if updatedSkill.Beschreibung != "Complete description" {
		t.Error("Expected skill description to be updated")
	}
	if updatedSkill.Category != "Alltag" {
		t.Error("Expected skill category to be updated")
	}
}

func TestImportCharacterSetsSourceIDDefault(t *testing.T) {
	setupImportTestEnvironment(t)

	exportData := &CharacterExport{
		Character: models.Char{
			BamortBase: models.BamortBase{
				Name: "Test Char",
			},
			Typ:    "Krieger",
			Rasse:  "Mensch",
			Grad:   1,
			UserID: 1,
		},
		GSMSkills: []models.Skill{
			{
				Name:         "TestSkill",
				GameSystemId: 1,
				SourceID:     0, // Should be set to 1
			},
		},
		GSMSpells: []models.Spell{
			{
				Name:         "TestSpell",
				GameSystemId: 1,
				SourceID:     0, // Should be set to 2
			},
		},
	}

	_, err := ImportCharacter(exportData, 1)
	if err != nil {
		t.Fatalf("ImportCharacter failed: %v", err)
	}

	// Verify skill source_id was set to 1
	var skill models.Skill
	database.DB.Where("name = ?", "TestSkill").First(&skill)
	if skill.SourceID != 1 {
		t.Errorf("Expected skill source_id to be 1, got %d", skill.SourceID)
	}

	// Verify spell source_id was set to 2
	var spell models.Spell
	database.DB.Where("name = ?", "TestSpell").First(&spell)
	if spell.SourceID != 2 {
		t.Errorf("Expected spell source_id to be 2, got %d", spell.SourceID)
	}
}

func TestImportCharacterIncludesAuditLog(t *testing.T) {
	setupImportTestEnvironment(t)

	// Create export data with audit log entries
	exportData := &CharacterExport{
		Character: models.Char{
			BamortBase: models.BamortBase{
				Name: "Test Char",
			},
			Typ:    "Krieger",
			Rasse:  "Mensch",
			Grad:   1,
			UserID: 1,
		},
		AuditLogEntries: []models.AuditLogEntry{
			{
				FieldName:  "experience_points",
				OldValue:   100,
				NewValue:   150,
				Difference: 50,
				Reason:     "skill_learning",
			},
		},
	}

	charID, err := ImportCharacter(exportData, 1)
	if err != nil {
		t.Fatalf("ImportCharacter failed: %v", err)
	}

	// Verify audit log entry was imported
	var auditEntries []models.AuditLogEntry
	database.DB.Where("character_id = ?", charID).Find(&auditEntries)

	if len(auditEntries) == 0 {
		t.Error("Expected audit log entries to be imported")
	}
}

func TestImportCharacterAsJSON(t *testing.T) {
	setupImportTestEnvironment(t)

	// Export character 18 to JSON
	exportData, err := ExportCharacter(uint(18))
	if err != nil {
		t.Fatalf("Failed to export character: %v", err)
	}

	jsonData, err := json.Marshal(exportData)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	// Unmarshal and import
	var importData CharacterExport
	err = json.Unmarshal(jsonData, &importData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	importData.Character.ID = 0
	importData.Character.Name = "JSON Imported"

	charID, err := ImportCharacter(&importData, 1)
	if err != nil {
		t.Fatalf("ImportCharacter from JSON failed: %v", err)
	}

	if charID == 0 {
		t.Error("Expected non-zero character ID")
	}
}
