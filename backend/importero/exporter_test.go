package importero

import (
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportChar2VTT(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Import a test character first
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	char, err := ImportVTTJSON(fileName, 1)
	assert.NoError(t, err, "Expected no error when importing char")

	// Export the character back to VTT format
	exportedChar, err := ExportCharToVTT(char)
	assert.NoError(t, err, "Expected no error when exporting char")

	// Basic validations
	assert.Equal(t, char.UserID, uint(1))
	assert.Equal(t, char.Name, exportedChar.Name)
	assert.Equal(t, char.Rasse, exportedChar.Rasse)
	assert.Equal(t, char.Typ, exportedChar.Typ)
	assert.Equal(t, char.Alter, exportedChar.Alter)
	assert.Equal(t, char.Grad, exportedChar.Grad)
	assert.Equal(t, char.Groesse, exportedChar.Groesse)
	assert.Equal(t, char.Gewicht, exportedChar.Gewicht)
	assert.Equal(t, char.Glaube, exportedChar.Glaube)
	assert.Equal(t, char.Hand, exportedChar.Hand)

	// Check LP
	assert.Equal(t, char.Lp.Max, exportedChar.Lp.Max)
	assert.Equal(t, char.Lp.Value, exportedChar.Lp.Value)

	// Check AP
	assert.Equal(t, char.Ap.Max, exportedChar.Ap.Max)
	assert.Equal(t, char.Ap.Value, exportedChar.Ap.Value)

	// Check Eigenschaften
	eigenschaftenMap := getEigenschaftenMap(char)
	assert.Equal(t, eigenschaftenMap["Au"], exportedChar.Eigenschaften.Au)
	assert.Equal(t, eigenschaftenMap["Gs"], exportedChar.Eigenschaften.Gs)
	assert.Equal(t, eigenschaftenMap["Gw"], exportedChar.Eigenschaften.Gw)
	assert.Equal(t, eigenschaftenMap["In"], exportedChar.Eigenschaften.In)
	assert.Equal(t, eigenschaftenMap["Ko"], exportedChar.Eigenschaften.Ko)
	assert.Equal(t, eigenschaftenMap["PA"], exportedChar.Eigenschaften.Pa)
	assert.Equal(t, eigenschaftenMap["St"], exportedChar.Eigenschaften.St)
	assert.Equal(t, eigenschaftenMap["Wk"], exportedChar.Eigenschaften.Wk)
	assert.Equal(t, eigenschaftenMap["Zt"], exportedChar.Eigenschaften.Zt)

	// Check Fertigkeiten exist
	assert.Greater(t, len(exportedChar.Fertigkeiten), 0, "Should have fertigkeiten")

	// Check Waffenfertigkeiten exist
	assert.Greater(t, len(exportedChar.Waffenfertigkeiten), 0, "Should have waffenfertigkeiten")
}

func TestExportChar2VTTRoundTrip(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Import original
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	char1, err := ImportVTTJSON(fileName, 1)
	assert.NoError(t, err, "Expected no error when importing char")

	// Export to VTT
	exportedChar, err := ExportCharToVTT(char1)
	assert.NoError(t, err, "Expected no error when exporting char")

	// Write to temp file
	tempFile, err := os.CreateTemp("", "vtt_export_*.json")
	assert.NoError(t, err, "Expected no error creating temp file")
	defer os.Remove(tempFile.Name())

	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(exportedChar)
	assert.NoError(t, err, "Expected no error encoding JSON")
	tempFile.Close()

	// Re-import the exported file
	char2, err := ImportVTTJSON(tempFile.Name(), 6)
	assert.NoError(t, err, "Expected no error when re-importing char")

	// Compare key fields
	assert.Equal(t, char1.UserID, uint(1), "UserID should match 1 as set in first import")
	assert.Equal(t, char2.UserID, uint(6), "UserID should match 6 as set in re-import")
	assert.Equal(t, char1.Name, char2.Name)
	assert.Equal(t, char1.Rasse, char2.Rasse)
	assert.Equal(t, char1.Typ, char2.Typ)
	assert.Equal(t, char1.Alter, char2.Alter)
	assert.Equal(t, char1.Grad, char2.Grad)
	assert.Equal(t, char1.Lp.Max, char2.Lp.Max)
	assert.Equal(t, char1.Ap.Max, char2.Ap.Max)
}

// Helper function to convert char eigenschaften array to map
func getEigenschaftenMap(char *models.Char) map[string]int {
	m := make(map[string]int)
	for _, e := range char.Eigenschaften {
		m[e.Name] = e.Value
	}
	return m
}

func TestExportSpellsToCSV(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Get some spells from master data
	var spells []models.Spell
	database.DB.Limit(10).Find(&spells)

	if len(spells) == 0 {
		t.Skip("No spells in test database")
	}

	// Export to CSV
	tempFile, err := os.CreateTemp("", "spell_export_*.csv")
	assert.NoError(t, err, "Expected no error creating temp file")
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	err = ExportSpellsToCSV(spells, tempFile.Name())
	assert.NoError(t, err, "Expected no error exporting spells to CSV")

	// Verify file exists and has content
	data, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err, "Expected no error reading CSV file")
	assert.Greater(t, len(data), 0, "CSV file should have content")

	// Verify CSV has header
	content := string(data)
	assert.Contains(t, content, "game_system")
	assert.Contains(t, content, "name")
	assert.Contains(t, content, "Beschreibung")
}

func TestExportCharToCSV(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Import a test character first
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	char, err := ImportVTTJSON(fileName, 1)
	assert.NoError(t, err, "Expected no error when importing char")

	// Export to CSV
	tempFile, err := os.CreateTemp("", "char_export_*.csv")
	assert.NoError(t, err, "Expected no error creating temp file")
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	err = ExportCharToCSV(char, tempFile.Name())
	assert.NoError(t, err, "Expected no error exporting character to CSV")

	// Verify file exists and has content
	data, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err, "Expected no error reading CSV file")
	assert.Greater(t, len(data), 0, "CSV file should have content")

	// Verify CSV has expected sections
	content := string(data)
	assert.Contains(t, content, char.Name, "Should contain character name")
	assert.Contains(t, content, "Basiseigenschaften", "Should contain base attributes section")
	assert.Contains(t, content, "Fertigkeit", "Should contain skills section")
	assert.Contains(t, content, "Waffe", "Should contain weapons section")
	assert.Contains(t, content, "Erfahrung", "Should contain experience section")
}

func TestExportImportWithoutMasterData(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Import a test character first
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	char1, err := ImportVTTJSON(fileName, 1)
	assert.NoError(t, err, "Expected no error when importing char")

	// Export to VTT
	vttChar, err := ExportCharToVTT(char1)
	assert.NoError(t, err, "Expected no error when exporting char")

	// Write to temp file
	tempFile, err := os.CreateTemp("", "vtt_export_*.json")
	assert.NoError(t, err, "Expected no error creating temp file")
	defer os.Remove(tempFile.Name())

	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(vttChar)
	assert.NoError(t, err, "Expected no error encoding JSON")
	tempFile.Close()

	// Clear all master data tables
	database.DB.Exec("DELETE FROM gsm_skills")
	database.DB.Exec("DELETE FROM gsm_weaponskills")
	database.DB.Exec("DELETE FROM gsm_spells")
	database.DB.Exec("DELETE FROM gsm_weapons")
	database.DB.Exec("DELETE FROM gsm_equipments")
	database.DB.Exec("DELETE FROM gsm_containers")
	database.DB.Exec("DELETE FROM gsm_transportations")
	database.DB.Exec("DELETE FROM gsm_believes")
	database.DB.Exec("DELETE FROM sqlite_sequence WHERE name LIKE 'gsm_%'")

	// Re-import without master data
	char2, err := ImportVTTJSON(tempFile.Name(), 1)
	assert.NoError(t, err, "Expected no error when re-importing without master data")
	assert.NotNil(t, char2, "Character should be imported")

	// Verify critical data was preserved
	assert.Equal(t, char1.Name, char2.Name, "Name should match")
	assert.Equal(t, char1.Rasse, char2.Rasse, "Race should match")
	assert.Equal(t, char1.Typ, char2.Typ, "Type should match")
	assert.Equal(t, char1.Grad, char2.Grad, "Grade should match")

	// Verify LP/AP
	assert.Equal(t, char1.Lp.Max, char2.Lp.Max, "LP Max should match")
	assert.Equal(t, char1.Ap.Max, char2.Ap.Max, "AP Max should match")

	// Verify skills were imported
	assert.Greater(t, len(char2.Fertigkeiten), 0, "Should have skills after reimport")
	assert.Greater(t, len(char2.Waffenfertigkeiten), 0, "Should have weapon skills after reimport")

	// Verify weapons were imported
	assert.Greater(t, len(char2.Waffen), 0, "Should have weapons after reimport")

	// Verify master data was created
	var skillCount, weaponSkillCount, weaponCount int64
	database.DB.Model(&models.Skill{}).Count(&skillCount)
	database.DB.Model(&models.WeaponSkill{}).Count(&weaponSkillCount)
	database.DB.Model(&models.Weapon{}).Count(&weaponCount)

	assert.Greater(t, skillCount, int64(0), "Master data should be created for skills")
	assert.Greater(t, weaponSkillCount, int64(0), "Master data should be created for weapon skills")
	assert.Greater(t, weaponCount, int64(0), "Master data should be created for weapons")
}

func TestExportImportPreservesCharacterData(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Import a test character
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	char1, err := ImportVTTJSON(fileName, 1)
	assert.NoError(t, err, "Expected no error when importing char")

	// Store original counts and values
	originalSkillCount := len(char1.Fertigkeiten)
	originalWeaponSkillCount := len(char1.Waffenfertigkeiten)
	originalSpellCount := len(char1.Zauber)
	originalWeaponCount := len(char1.Waffen)
	originalEquipmentCount := len(char1.Ausruestung)
	originalContainerCount := len(char1.Behaeltnisse)
	originalTransportCount := len(char1.Transportmittel)

	// Export to VTT
	vttChar, err := ExportCharToVTT(char1)
	assert.NoError(t, err, "Expected no error when exporting char")

	// Verify export has all data
	assert.Equal(t, originalSkillCount, len(vttChar.Fertigkeiten), "All skills should be exported")
	assert.Equal(t, originalWeaponSkillCount, len(vttChar.Waffenfertigkeiten), "All weapon skills should be exported")
	assert.Equal(t, originalSpellCount, len(vttChar.Zauber), "All spells should be exported")
	assert.Equal(t, originalWeaponCount, len(vttChar.Waffen), "All weapons should be exported")
	assert.Equal(t, originalEquipmentCount, len(vttChar.Ausruestung), "All equipment should be exported")
	assert.Equal(t, originalContainerCount, len(vttChar.Behaeltnisse), "All containers should be exported")
	assert.Equal(t, originalTransportCount, len(vttChar.Transportmittel), "All transportation should be exported")

	// Verify specific skill data is preserved
	if originalSkillCount > 0 {
		assert.NotEmpty(t, vttChar.Fertigkeiten[0].Name, "Skill name should be exported")
		assert.GreaterOrEqual(t, vttChar.Fertigkeiten[0].Fertigkeitswert, 0, "Skill value should be exported")
	}

	// Verify weapon data is preserved
	if originalWeaponCount > 0 {
		assert.NotEmpty(t, vttChar.Waffen[0].Name, "Weapon name should be exported")
		assert.GreaterOrEqual(t, vttChar.Waffen[0].Gewicht, float64(0), "Weapon weight should be exported")
	}

	// Write to temp file
	tempFile, err := os.CreateTemp("", "vtt_export_*.json")
	assert.NoError(t, err, "Expected no error creating temp file")
	defer os.Remove(tempFile.Name())

	encoder := json.NewEncoder(tempFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(vttChar)
	assert.NoError(t, err, "Expected no error encoding JSON")
	tempFile.Close()

	// Clear master data
	database.DB.Exec("DELETE FROM gsm_skills")
	database.DB.Exec("DELETE FROM gsm_weaponskills")
	database.DB.Exec("DELETE FROM gsm_spells")
	database.DB.Exec("DELETE FROM gsm_weapons")
	database.DB.Exec("DELETE FROM gsm_equipments")
	database.DB.Exec("DELETE FROM gsm_containers")
	database.DB.Exec("DELETE FROM gsm_transportations")

	// Re-import
	char2, err := ImportVTTJSON(tempFile.Name(), 1)
	assert.NoError(t, err, "Expected no error when re-importing")

	// Verify all data was preserved
	assert.Equal(t, originalSkillCount, len(char2.Fertigkeiten), "All skills should be reimported")
	assert.Equal(t, originalWeaponSkillCount, len(char2.Waffenfertigkeiten), "All weapon skills should be reimported")
	assert.Equal(t, originalSpellCount, len(char2.Zauber), "All spells should be reimported")
	assert.Equal(t, originalWeaponCount, len(char2.Waffen), "All weapons should be reimported")
	assert.Equal(t, originalEquipmentCount, len(char2.Ausruestung), "All equipment should be reimported")
	assert.Equal(t, originalContainerCount, len(char2.Behaeltnisse), "All containers should be reimported")
	assert.Equal(t, originalTransportCount, len(char2.Transportmittel), "All transportation should be reimported")

	// Verify specific values match
	if originalSkillCount > 0 {
		skill1 := findSkillByName(char1.Fertigkeiten, char1.Fertigkeiten[0].Name)
		skill2 := findSkillByName(char2.Fertigkeiten, char1.Fertigkeiten[0].Name)
		assert.NotNil(t, skill1, "Original skill should exist")
		assert.NotNil(t, skill2, "Reimported skill should exist")
		assert.Equal(t, skill1.Fertigkeitswert, skill2.Fertigkeitswert, "Skill values should match")
	}
}

// Helper function to find skill by name
func findSkillByName(skills []models.SkFertigkeit, name string) *models.SkFertigkeit {
	for i := range skills {
		if skills[i].Name == name {
			return &skills[i]
		}
	}
	return nil
}
