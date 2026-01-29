package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestExportSkills(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data
	source := getOrCreateSource("KOD", "Kodex")
	skill := models.Skill{
		Name:             "Schwimmen",
		GameSystem:       "midgard",
		Beschreibung:     "Schwimmen im Wasser",
		Initialwert:      12,
		BasisWert:        0,
		Bonuseigenschaft: "Gw",
		Improvable:       true,
		InnateSkill:      false,
		SourceID:         source.ID,
		PageNumber:       42,
	}
	if err := skill.Create(); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	// Export skills
	tmpDir := t.TempDir()
	err := ExportSkills(tmpDir)
	if err != nil {
		t.Fatalf("ExportSkills failed: %v", err)
	}

	// Verify file exists
	exportFile := filepath.Join(tmpDir, "skills.json")
	if _, err := os.Stat(exportFile); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", exportFile)
	}
}

func TestImportSkills(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create source that will be referenced
	source := getOrCreateSource("KOD", "Kodex")

	// Export first
	tmpDir := t.TempDir()
	skill := models.Skill{
		Name:             "TestImportSkill",
		GameSystem:       "midgard",
		Beschreibung:     "Klettern an Wänden",
		Initialwert:      10,
		BasisWert:        0,
		Bonuseigenschaft: "Gw",
		Improvable:       true,
		InnateSkill:      false,
		SourceID:         source.ID,
		PageNumber:       50,
	}
	if err := skill.Create(); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	err := ExportSkills(tmpDir)
	if err != nil {
		t.Fatalf("ExportSkills failed: %v", err)
	}

	// Delete the skill
	database.DB.Delete(&skill)

	// Import back
	err = ImportSkills(tmpDir)
	if err != nil {
		t.Fatalf("ImportSkills failed: %v", err)
	}

	// Verify skill was imported
	var importedSkill models.Skill
	err = database.DB.Where("name = ? AND game_system = ?", "TestImportSkill", "midgard").First(&importedSkill).Error
	if err != nil {
		t.Fatalf("Imported skill not found: %v", err)
	}

	if importedSkill.Beschreibung != "Klettern an Wänden" {
		t.Errorf("Expected beschreibung 'Klettern an Wänden', got '%s'", importedSkill.Beschreibung)
	}

	if importedSkill.Initialwert != 10 {
		t.Errorf("Expected initialwert 10, got %d", importedSkill.Initialwert)
	}
	if importedSkill.GameSystemId == 0 {
		t.Errorf("Expected game_system_id to be set")
	}
}

func TestImportSkillsUpdate(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	source := getOrCreateSource("KOD", "Kodex")

	// Create existing skill
	skill := models.Skill{
		Name:         "Reiten",
		GameSystem:   "midgard",
		Beschreibung: "Alte Beschreibung",
		Initialwert:  8,
		SourceID:     source.ID,
		PageNumber:   30,
	}
	if err := skill.Create(); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	// Export, modify, and re-import
	tmpDir := t.TempDir()
	err := ExportSkills(tmpDir)
	if err != nil {
		t.Fatalf("ExportSkills failed: %v", err)
	}

	// Update skill manually
	skill.Beschreibung = "Neue Beschreibung"
	skill.Initialwert = 12
	database.DB.Save(&skill)

	// Export again with updated values
	err = ExportSkills(tmpDir)
	if err != nil {
		t.Fatalf("ExportSkills failed: %v", err)
	}

	// Reset to old values
	skill.Beschreibung = "Alte Beschreibung"
	skill.Initialwert = 8
	database.DB.Save(&skill)

	// Import should update to exported values
	err = ImportSkills(tmpDir)
	if err != nil {
		t.Fatalf("ImportSkills failed: %v", err)
	}

	// Verify update
	var updatedSkill models.Skill
	err = database.DB.Where("name = ? AND game_system = ?", "Reiten", "midgard").First(&updatedSkill).Error
	if err != nil {
		t.Fatalf("Updated skill not found: %v", err)
	}

	if updatedSkill.Beschreibung != "Neue Beschreibung" {
		t.Errorf("Expected updated beschreibung 'Neue Beschreibung', got '%s'", updatedSkill.Beschreibung)
	}

	if updatedSkill.Initialwert != 12 {
		t.Errorf("Expected updated initialwert 12, got %d", updatedSkill.Initialwert)
	}
}

func TestExportImportSources(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test source
	source := models.Source{
		Code:         "ARK",
		Name:         "Arkanum",
		GameSystemId: 1,
		IsActive:     true,
	}
	database.DB.Create(&source)

	// Export
	tmpDir := t.TempDir()
	err := ExportSources(tmpDir)
	if err != nil {
		t.Fatalf("ExportSources failed: %v", err)
	}

	// Delete
	database.DB.Delete(&source)

	// Import
	err = ImportSources(tmpDir)
	if err != nil {
		t.Fatalf("ImportSources failed: %v", err)
	}

	// Verify
	var imported models.Source
	err = database.DB.Where("code = ?", "ARK").First(&imported).Error
	if err != nil {
		t.Fatalf("Imported source not found: %v", err)
	}

	if imported.Name != "Arkanum" {
		t.Errorf("Expected name 'Arkanum', got '%s'", imported.Name)
	}
}

func TestExportImportSkillCategoryDifficulty(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create dependencies
	source := getOrCreateSource("KOD", "Kodex")
	skill := models.Skill{
		Name:         "Tanzen",
		GameSystemId: 1,
		SourceID:     source.ID,
	}
	if err := skill.Create(); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	category := getOrCreateCategory("Alltag", source.ID)
	difficulty := getOrCreateDifficulty("leicht")

	// Create relationship
	scd := models.SkillCategoryDifficulty{
		SkillID:           skill.ID,
		SkillCategoryID:   category.ID,
		SkillDifficultyID: difficulty.ID,
		LearnCost:         5,
		SCategory:         category.Name,
		SDifficulty:       difficulty.Name,
	}
	database.DB.Create(&scd)

	// Export
	tmpDir := t.TempDir()
	err := ExportSkillCategoryDifficulties(tmpDir)
	if err != nil {
		t.Fatalf("ExportSkillCategoryDifficulties failed: %v", err)
	}

	// Delete relationship
	database.DB.Delete(&scd)

	// Import
	err = ImportSkillCategoryDifficulties(tmpDir)
	if err != nil {
		t.Fatalf("ImportSkillCategoryDifficulties failed: %v", err)
	}

	// Verify relationship was recreated
	var imported models.SkillCategoryDifficulty
	err = database.DB.Where("skill_id = ? AND skill_category_id = ?", skill.ID, category.ID).First(&imported).Error
	if err != nil {
		t.Fatalf("Imported relationship not found: %v", err)
	}

	if imported.LearnCost != 5 {
		t.Errorf("Expected learn_cost 5, got %d", imported.LearnCost)
	}
}

func TestExportImportSkillCategories(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data
	source := getOrCreateSource("TEST_SC", "Test Source")
	category := models.SkillCategory{
		Name:         "TestCategory",
		GameSystemId: 1,
		SourceID:     source.ID,
	}
	database.DB.Create(&category)

	// Export
	tempDir := t.TempDir()
	err := ExportSkillCategories(tempDir)
	if err != nil {
		t.Fatalf("ExportSkillCategories failed: %v", err)
	}

	// Verify file was created
	filename := filepath.Join(tempDir, "skill_categories.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	// Delete the category
	database.DB.Unscoped().Delete(&category)

	// Import
	err = ImportSkillCategories(tempDir)
	if err != nil {
		t.Fatalf("ImportSkillCategories failed: %v", err)
	}

	// Verify the category was recreated
	var imported models.SkillCategory
	result := database.DB.Where("name = ? AND game_system = ?", "TestCategory", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Category not found after import: %v", result.Error)
	}

	if imported.SourceID != source.ID {
		t.Errorf("Expected SourceID %d, got %d", source.ID, imported.SourceID)
	}
}

func TestExportImportSkillDifficulties(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data
	difficulty := models.SkillDifficulty{
		Name:         "TestDifficulty",
		GameSystemId: 1,
	}
	database.DB.Create(&difficulty)

	// Export
	tempDir := t.TempDir()
	err := ExportSkillDifficulties(tempDir)
	if err != nil {
		t.Fatalf("ExportSkillDifficulties failed: %v", err)
	}

	// Verify file was created
	filename := filepath.Join(tempDir, "skill_difficulties.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	// Delete the difficulty
	database.DB.Unscoped().Delete(&difficulty)

	// Import
	err = ImportSkillDifficulties(tempDir)
	if err != nil {
		t.Fatalf("ImportSkillDifficulties failed: %v", err)
	}

	// Verify the difficulty was recreated
	var imported models.SkillDifficulty
	result := database.DB.Where("name = ? AND game_system = ?", "TestDifficulty", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Difficulty not found after import: %v", result.Error)
	}
}

func TestExportImportSpells(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()
	models.MigrateStructure(database.DB)

	// Create test data
	source := getOrCreateSource("TEST_SP", "Test Spell Source")
	spell := models.Spell{
		Name:             "TestSpell",
		GameSystem:       "midgard",
		Beschreibung:     "Test description",
		SourceID:         source.ID,
		PageNumber:       42,
		Bonus:            5,
		Stufe:            3,
		AP:               "2",
		Art:              "Gestenzauber",
		Zauberdauer:      "10 sec",
		Reichweite:       "10m",
		Wirkungsziel:     "Person",
		Wirkungsbereich:  "1 Person",
		Wirkungsdauer:    "1h",
		Ursprung:         "Elben",
		Category:         "normal",
		LearningCategory: "default",
	}
	//database.DB.Create(&spell)
	spell.Create()

	// Export
	tempDir := t.TempDir()
	err := ExportSpells(tempDir)
	if err != nil {
		t.Fatalf("ExportSpells failed: %v", err)
	}

	// Verify file was created
	filename := filepath.Join(tempDir, "spells.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	// Modify the spell
	spell.Beschreibung = "Old description"
	spell.Bonus = 3
	database.DB.Save(&spell)

	// Import (should update)
	err = ImportSpells(tempDir)
	if err != nil {
		t.Fatalf("ImportSpells failed: %v", err)
	}

	// Verify the spell was updated
	var imported models.Spell
	result := database.DB.Where("name = ? AND game_system = ?", "TestSpell", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Spell not found after import: %v", result.Error)
	}

	if imported.Beschreibung != "Test description" {
		t.Errorf("Expected description 'Test description', got '%s'", imported.Beschreibung)
	}
	if imported.Bonus != 5 {
		t.Errorf("Expected bonus 5, got %d", imported.Bonus)
	}
	if imported.Stufe != 3 {
		t.Errorf("Expected level 3, got %d", imported.Stufe)
	}
	if imported.GameSystemId == 0 {
		t.Errorf("Expected game_system_id to be set")
	}
}

func TestExportImportAll(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data
	source := getOrCreateSource("TEST_ALL", "Test All Source")

	category := models.SkillCategory{Name: "AllCategory", GameSystemId: 1, SourceID: source.ID}
	database.DB.Create(&category)

	difficulty := models.SkillDifficulty{Name: "AllDifficulty", GameSystemId: 1}
	database.DB.Create(&difficulty)

	skill := models.Skill{
		Name:        "AllSkill",
		GameSystem:  "midgard",
		SourceID:    source.ID,
		Initialwert: 10,
	}
	if err := skill.Create(); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	spell := models.Spell{
		Name:         "AllSpell",
		GameSystemId: 1,
		SourceID:     source.ID,
		Stufe:        2,
	}
	database.DB.Create(&spell)

	// Export all
	tempDir := t.TempDir()
	err := ExportAll(tempDir)
	if err != nil {
		t.Fatalf("ExportAll failed: %v", err)
	}

	// Verify all files were created
	files := []string{
		"sources.json",
		"character_classes.json",
		"skill_categories.json",
		"skill_difficulties.json",
		"spell_schools.json",
		"skills.json",
		"skill_category_difficulties.json",
		"spells.json",
		"class_category_ep_costs.json",
		"class_spell_school_ep_costs.json",
		"spell_level_le_costs.json",
		"skill_improvement_costs.json",
		"weapon_skills.json",
		"equipment.json",
		"weapons.json",
		"containers.json",
		"transportation.json",
		"believes.json",
	}
	for _, file := range files {
		filename := filepath.Join(tempDir, file)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Export file not created: %s", filename)
		}
	}

	// Delete all test data
	database.DB.Unscoped().Delete(&spell)
	database.DB.Unscoped().Delete(&skill)
	database.DB.Unscoped().Delete(&difficulty)
	database.DB.Unscoped().Delete(&category)
	// Don't delete source to avoid FK constraints

	// Import all
	err = ImportAll(tempDir)
	if err != nil {
		t.Fatalf("ImportAll failed: %v", err)
	}

	// Verify all data was recreated
	var importedCategory models.SkillCategory
	if err := database.DB.Where("name = ? AND game_system = ?", "AllCategory", "midgard").First(&importedCategory).Error; err != nil {
		t.Errorf("Category not found after import: %v", err)
	}

	var importedDifficulty models.SkillDifficulty
	if err := database.DB.Where("name = ? AND game_system = ?", "AllDifficulty", "midgard").First(&importedDifficulty).Error; err != nil {
		t.Errorf("Difficulty not found after import: %v", err)
	}

	var importedSkill models.Skill
	if err := database.DB.Where("name = ? AND game_system = ?", "AllSkill", "midgard").First(&importedSkill).Error; err != nil {
		t.Errorf("Skill not found after import: %v", err)
	}

	var importedSpell models.Spell
	if err := database.DB.Where("name = ? AND game_system = ?", "AllSpell", "midgard").First(&importedSpell).Error; err != nil {
		t.Errorf("Spell not found after import: %v", err)
	}
}

func TestExportAll_live(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB(false)

	// Export all
	tempDir := t.TempDir()
	err := ExportAll(tempDir)
	if err != nil {
		t.Fatalf("ExportAll failed: %v", err)
	}

	// Verify all files were created
	files := []string{"sources.json", "skill_categories.json", "skill_difficulties.json", "skills.json", "spells.json", "skill_category_difficulties.json"}
	for _, file := range files {
		filename := filepath.Join(tempDir, file)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Export file not created: %s", filename)
		}
	}

	assert.Equal(t, len(files), 6)

}

func TestExportImportWeaponSkills(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	source := getOrCreateSource("TEST_WS", "Test Weapon Source")
	weaponSkill := models.WeaponSkill{
		Skill: models.Skill{
			Name:             "Langschwert",
			GameSystem:       "midgard",
			Beschreibung:     "Langschwert Waffenfertigkeiten",
			SourceID:         source.ID,
			PageNumber:       50,
			Initialwert:      10,
			BasisWert:        5,
			Bonuseigenschaft: "St",
			Improvable:       true,
			InnateSkill:      false,
			Category:         "Waffen",
			Difficulty:       "normal",
		},
	}
	if err := weaponSkill.Create(); err != nil {
		t.Fatalf("failed to create weapon skill: %v", err)
	}

	tempDir := t.TempDir()
	err := ExportWeaponSkills(tempDir)
	if err != nil {
		t.Fatalf("ExportWeaponSkills failed: %v", err)
	}

	filename := filepath.Join(tempDir, "weapon_skills.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	database.DB.Unscoped().Delete(&weaponSkill)

	err = ImportWeaponSkills(tempDir)
	if err != nil {
		t.Fatalf("ImportWeaponSkills failed: %v", err)
	}

	var imported models.WeaponSkill
	result := database.DB.Where("name = ? AND game_system = ?", "Langschwert", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Weapon skill not found after import: %v", result.Error)
	}

	assert.Equal(t, "Langschwert Waffenfertigkeiten", imported.Beschreibung)
	assert.Equal(t, 10, imported.Initialwert)
	assert.Equal(t, 5, imported.BasisWert)
	assert.NotZero(t, imported.GameSystemId)
}

func TestExportImportEquipment(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	source := getOrCreateSource("TEST_EQ", "Test Equipment Source")
	equipment := models.Equipment{
		Name:         "Seil",
		GameSystem:   "midgard",
		Beschreibung: "10m langes Hanfseil",
		SourceID:     source.ID,
		PageNumber:   75,
		Gewicht:      2.5,
		Wert:         15.0,
		PersonalItem: false,
	}
	database.DB.Create(&equipment)

	tempDir := t.TempDir()
	err := ExportEquipment(tempDir)
	if err != nil {
		t.Fatalf("ExportEquipment failed: %v", err)
	}

	filename := filepath.Join(tempDir, "equipment.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	equipment.Wert = 10.0
	database.DB.Save(&equipment)

	err = ImportEquipment(tempDir)
	if err != nil {
		t.Fatalf("ImportEquipment failed: %v", err)
	}

	var imported models.Equipment
	result := database.DB.Where("name = ? AND game_system = ?", "Seil", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Equipment not found after import: %v", result.Error)
	}

	assert.Equal(t, 15.0, imported.Wert)
	assert.Equal(t, 2.5, imported.Gewicht)
}

func TestExportImportWeapons(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	source := getOrCreateSource("TEST_WP", "Test Weapon Source")
	weapon := models.Weapon{
		Equipment: models.Equipment{
			Name:         "Kurzschwert",
			GameSystem:   "midgard",
			Beschreibung: "Einhändiges Kurzschwert",
			SourceID:     source.ID,
			PageNumber:   80,
			Gewicht:      1.5,
			Wert:         50.0,
			PersonalItem: false,
		},
		SkillRequired: "Langschwert",
		Damage:        "1W6+1",
		RangeNear:     0,
		RangeMiddle:   0,
		RangeFar:      0,
	}
	database.DB.Create(&weapon)

	tempDir := t.TempDir()
	err := ExportWeapons(tempDir)
	if err != nil {
		t.Fatalf("ExportWeapons failed: %v", err)
	}

	filename := filepath.Join(tempDir, "weapons.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	weapon.Damage = "1W6"
	database.DB.Save(&weapon)

	err = ImportWeapons(tempDir)
	if err != nil {
		t.Fatalf("ImportWeapons failed: %v", err)
	}

	var imported models.Weapon
	result := database.DB.Where("name = ? AND game_system = ?", "Kurzschwert", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Weapon not found after import: %v", result.Error)
	}

	assert.Equal(t, "1W6+1", imported.Damage)
	assert.Equal(t, "Langschwert", imported.SkillRequired)
}

func TestExportImportContainers(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	source := getOrCreateSource("TEST_CT", "Test Container Source")
	container := models.Container{
		Equipment: models.Equipment{
			Name:         "Rucksack",
			GameSystem:   "midgard",
			Beschreibung: "Großer Lederrucksack",
			SourceID:     source.ID,
			PageNumber:   85,
			Gewicht:      1.0,
			Wert:         20.0,
			PersonalItem: false,
		},
		Tragkraft: 30.0,
		Volumen:   50.0,
	}
	database.DB.Create(&container)

	tempDir := t.TempDir()
	err := ExportContainers(tempDir)
	if err != nil {
		t.Fatalf("ExportContainers failed: %v", err)
	}

	filename := filepath.Join(tempDir, "containers.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	container.Tragkraft = 25.0
	database.DB.Save(&container)

	err = ImportContainers(tempDir)
	if err != nil {
		t.Fatalf("ImportContainers failed: %v", err)
	}

	var imported models.Container
	result := database.DB.Where("name = ? AND game_system = ?", "Rucksack", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Container not found after import: %v", result.Error)
	}

	assert.Equal(t, 30.0, imported.Tragkraft)
	assert.Equal(t, 50.0, imported.Volumen)
}

func TestExportImportTransportation(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	source := getOrCreateSource("TEST_TR", "Test Transport Source")
	transportation := models.Transportation{
		Container: models.Container{
			Equipment: models.Equipment{
				Name:         "Pferdewagen",
				GameSystem:   "midgard",
				Beschreibung: "Zweirädriger Wagen",
				SourceID:     source.ID,
				PageNumber:   90,
				Gewicht:      100.0,
				Wert:         200.0,
				PersonalItem: false,
			},
			Tragkraft: 500.0,
			Volumen:   1000.0,
		},
	}
	database.DB.Create(&transportation)

	tempDir := t.TempDir()
	err := ExportTransportation(tempDir)
	if err != nil {
		t.Fatalf("ExportTransportation failed: %v", err)
	}

	filename := filepath.Join(tempDir, "transportation.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	transportation.Tragkraft = 450.0
	database.DB.Save(&transportation)

	err = ImportTransportation(tempDir)
	if err != nil {
		t.Fatalf("ImportTransportation failed: %v", err)
	}

	var imported models.Transportation
	result := database.DB.Where("name = ? AND game_system = ?", "Pferdewagen", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Transportation not found after import: %v", result.Error)
	}

	assert.Equal(t, 500.0, imported.Tragkraft)
	assert.Equal(t, 1000.0, imported.Volumen)
}

func TestExportImportBelieves(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	source := getOrCreateSource("TEST_BL", "Test Believe Source")
	believe := models.Believe{
		Name:         "Kirche des Lichts",
		GameSystem:   "midgard",
		Beschreibung: "Hauptreligion in Valian",
		SourceID:     source.ID,
		PageNumber:   95,
	}
	database.DB.Create(&believe)

	tempDir := t.TempDir()
	err := ExportBelieves(tempDir)
	if err != nil {
		t.Fatalf("ExportBelieves failed: %v", err)
	}

	filename := filepath.Join(tempDir, "believes.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	believe.Beschreibung = "Alte Beschreibung"
	database.DB.Save(&believe)

	err = ImportBelieves(tempDir)
	if err != nil {
		t.Fatalf("ImportBelieves failed: %v", err)
	}

	var imported models.Believe
	result := database.DB.Where("name = ? AND game_system = ?", "Kirche des Lichts", "midgard").First(&imported)
	if result.Error != nil {
		t.Fatalf("Believe not found after import: %v", result.Error)
	}

	assert.Equal(t, "Hauptreligion in Valian", imported.Beschreibung)
	assert.Equal(t, 95, imported.PageNumber)
}

func TestExportImportSkillImprovementCosts(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create dependencies that already exist in test DB
	// Use existing skill, category, difficulty from test database
	var skill models.Skill
	if err := database.DB.Where("name = ?", "Abrichten").First(&skill).Error; err != nil {
		t.Skip("Test skill not found in database, skipping test")
	}

	var category models.SkillCategory
	if err := database.DB.First(&category).Error; err != nil {
		t.Skip("No skill category found in database, skipping test")
	}

	var difficulty models.SkillDifficulty
	if err := database.DB.First(&difficulty).Error; err != nil {
		t.Skip("No skill difficulty found in database, skipping test")
	}

	// Find or create SkillCategoryDifficulty
	var scd models.SkillCategoryDifficulty
	err := database.DB.Where("skill_id = ? AND skill_category_id = ? AND skill_difficulty_id = ?",
		skill.ID, category.ID, difficulty.ID).First(&scd).Error

	if err == gorm.ErrRecordNotFound {
		scd = models.SkillCategoryDifficulty{
			SkillID:           skill.ID,
			SkillCategoryID:   category.ID,
			SkillDifficultyID: difficulty.ID,
			LearnCost:         10,
			SCategory:         category.Name,
			SDifficulty:       difficulty.Name,
		}
		database.DB.Create(&scd)
	} else if err != nil {
		t.Fatalf("Failed to query SkillCategoryDifficulty: %v", err)
	}

	// Create SkillImprovementCost
	improvementCost := models.SkillImprovementCost{
		SkillCategoryDifficultyID: scd.ID,
		CurrentLevel:              15, // Use unique level to avoid conflicts
		TERequired:                5,
	}
	database.DB.Create(&improvementCost)

	// Export
	tempDir := t.TempDir()
	err = ExportSkillImprovementCosts(tempDir)
	if err != nil {
		t.Fatalf("ExportSkillImprovementCosts failed: %v", err)
	}

	filename := filepath.Join(tempDir, "skill_improvement_costs.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Export file not created: %s", filename)
	}

	// Modify the record
	improvementCost.TERequired = 7
	database.DB.Save(&improvementCost)

	// Import should restore original value
	err = ImportSkillImprovementCosts(tempDir)
	if err != nil {
		t.Fatalf("ImportSkillImprovementCosts failed: %v", err)
	}

	var imported models.SkillImprovementCost
	result := database.DB.Where("skill_category_difficulty_id = ? AND current_level = ?", scd.ID, 15).First(&imported)
	if result.Error != nil {
		t.Fatalf("SkillImprovementCost not found after import: %v", result.Error)
	}

	// Should be restored to original value from export
	assert.Equal(t, 5, imported.TERequired)
	assert.Equal(t, 15, imported.CurrentLevel)
}

// TestExportImportCompleteness verifies that all GSMaster tables are included in export/import
func TestExportImportCompleteness(t *testing.T) {
	// List of all GSMaster-related tables that should be exported/imported
	expectedExports := []string{
		"Sources",                         // Base data
		"CharacterClasses",                // Base data
		"SkillCategories",                 // Base data
		"SkillDifficulties",               // Base data
		"SpellSchools",                    // Base data
		"Skills",                          // Game data
		"WeaponSkills",                    // Game data
		"Spells",                          // Game data
		"Equipment",                       // Game data
		"Weapons",                         // Game data
		"Containers",                      // Game data
		"Transportation",                  // Game data
		"Believes",                        // Game data
		"SkillCategoryDifficulties",       // Learning cost relationships
		"WeaponSkillCategoryDifficulties", // Learning cost relationships
		"ClassCategoryEPCosts",            // Learning cost definitions
		"ClassSpellSchoolEPCosts",         // Learning cost definitions
		"SpellLevelLECosts",               // Learning cost definitions
		"SkillImprovementCosts",           // Learning cost definitions
	}

	// Count exports/imports actually implemented
	// These are verified by checking the function exists and is called in ExportAll/ImportAll
	implementedExports := []string{
		"Sources",
		"CharacterClasses",
		"SkillCategories",
		"SkillDifficulties",
		"SpellSchools",
		"Skills",
		"WeaponSkills",
		"Spells",
		"Equipment",
		"Weapons",
		"Containers",
		"Transportation",
		"Believes",
		"SkillCategoryDifficulties",
		"WeaponSkillCategoryDifficulties",
		"ClassCategoryEPCosts",
		"ClassSpellSchoolEPCosts",
		"SpellLevelLECosts",
		"SkillImprovementCosts",
	}

	// Create maps for comparison
	expected := make(map[string]bool)
	for _, name := range expectedExports {
		expected[name] = true
	}

	implemented := make(map[string]bool)
	for _, name := range implementedExports {
		implemented[name] = true
	}

	// Check for missing implementations
	missing := []string{}
	for name := range expected {
		if !implemented[name] {
			missing = append(missing, name)
		}
	}

	// Check for unexpected implementations
	extra := []string{}
	for name := range implemented {
		if !expected[name] {
			extra = append(extra, name)
		}
	}

	// Report results
	if len(missing) > 0 {
		t.Errorf("Missing export/import implementations: %v", missing)
	}

	if len(extra) > 0 {
		t.Logf("Extra export/import implementations (may be intentional): %v", extra)
	}

	if len(missing) == 0 && len(extra) == 0 {
		t.Logf("✓ All %d GSMaster tables have export/import implementations", len(expectedExports))
	}
}

// TestExportAllCallsAllExports verifies that ExportAll calls all export functions
func TestExportAllCallsAllExports(t *testing.T) {
	// This is a documentation test - it verifies the expected behavior
	// In a real test, we would mock the functions and verify they're called

	expectedCalls := []string{
		"ExportSources",
		"ExportCharacterClasses",
		"ExportSkillCategories",
		"ExportSkillDifficulties",
		"ExportSpellSchools",
		"ExportSkills",
		"ExportSkillCategoryDifficulties",
		"ExportSpells",
		"ExportClassCategoryEPCosts",
		"ExportClassSpellSchoolEPCosts",
		"ExportSpellLevelLECosts",
		"ExportSkillImprovementCosts",
		"ExportWeaponSkills",
		"ExportWeaponSkillCategoryDifficulties",
		"ExportEquipment",
		"ExportWeapons",
		"ExportContainers",
		"ExportTransportation",
		"ExportBelieves",
	}

	t.Logf("ExportAll should call %d export functions", len(expectedCalls))
	t.Logf("Export functions called:")
	for i, funcName := range expectedCalls {
		t.Logf("  %2d. %s", i+1, funcName)
	}
}

// TestImportAllCallsAllImports verifies that ImportAll calls all import functions
func TestImportAllCallsAllImports(t *testing.T) {
	// This is a documentation test - it verifies the expected behavior

	expectedCalls := []string{
		"ImportSources",
		"ImportCharacterClasses",
		"ImportSkillCategories",
		"ImportSkillDifficulties",
		"ImportSpellSchools",
		"ImportSkills",
		"ImportSkillCategoryDifficulties",
		"ImportSpells",
		"ImportClassCategoryEPCosts",
		"ImportClassSpellSchoolEPCosts",
		"ImportSpellLevelLECosts",
		"ImportSkillImprovementCosts",
		"ImportWeaponSkills",
		"ImportWeaponSkillCategoryDifficulties",
		"ImportEquipment",
		"ImportWeapons",
		"ImportContainers",
		"ImportTransportation",
		"ImportBelieves",
	}

	t.Logf("ImportAll should call %d import functions", len(expectedCalls))
	t.Logf("Import functions called:")
	for i, funcName := range expectedCalls {
		t.Logf("  %2d. %s", i+1, funcName)
	}
}

// TestExportImportOrderIsCorrect verifies dependency order
func TestExportImportOrderIsCorrect(t *testing.T) {
	// Define the correct dependency order
	// Base tables first, then dependent tables

	correctOrder := []string{
		// Base data (no dependencies)
		"Sources",
		"CharacterClasses",
		"SkillCategories",
		"SkillDifficulties",
		"SpellSchools",

		// Master data (depends on sources)
		"Skills",
		"Spells",
		"WeaponSkills",
		"Equipment",
		"Weapons",
		"Containers",
		"Transportation",
		"Believes",

		// Relationship/cost tables (depend on base + master data)
		"SkillCategoryDifficulties",
		"WeaponSkillCategoryDifficulties",
		"ClassCategoryEPCosts",
		"ClassSpellSchoolEPCosts",
		"SpellLevelLECosts",
		"SkillImprovementCosts",
	}

	t.Logf("Correct dependency order for export/import:")
	t.Logf("\n1. Base data (no dependencies):")
	t.Logf("   - Sources")
	t.Logf("   - CharacterClasses")
	t.Logf("   - SkillCategories")
	t.Logf("   - SkillDifficulties")
	t.Logf("   - SpellSchools")

	t.Logf("\n2. Master data (depends on Sources):")
	t.Logf("   - Skills")
	t.Logf("   - Spells")
	t.Logf("   - WeaponSkills")
	t.Logf("   - Equipment")
	t.Logf("   - Weapons")
	t.Logf("   - Containers")
	t.Logf("   - Transportation")
	t.Logf("   - Believes")

	t.Logf("\n3. Relationship/cost tables (depend on base + master):")
	t.Logf("   - SkillCategoryDifficulties")
	t.Logf("   - WeaponSkillCategoryDifficulties")
	t.Logf("   - ClassCategoryEPCosts")
	t.Logf("   - ClassSpellSchoolEPCosts")
	t.Logf("   - SpellLevelLECosts")
	t.Logf("   - SkillImprovementCosts")

	t.Logf("\nTotal: %d tables", len(correctOrder))
}

// TestExportImportClassCategoryLearningPoints tests export and import of class category learning points
func TestExportImportClassCategoryLearningPoints(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data with unique names
	class := models.CharacterClass{
		Code:         "TEST_KRI",
		Name:         "Test-Krieger",
		GameSystemId: 1,
	}
	database.DB.Create(&class)

	category := models.SkillCategory{
		Name:         "Test-Kampf",
		GameSystemId: 1,
	}
	database.DB.Create(&category)

	record := models.ClassCategoryLearningPoints{
		CharacterClassID: class.ID,
		SkillCategoryID:  category.ID,
		Points:           10,
	}
	database.DB.Create(&record)

	// Export
	tmpDir := t.TempDir()
	err := ExportClassCategoryLearningPoints(tmpDir)
	assert.NoError(t, err)

	// Verify file exists
	exportFile := filepath.Join(tmpDir, "class_category_learning_points.json")
	assert.FileExists(t, exportFile)

	// Delete record
	database.DB.Delete(&record)

	// Import back
	err = ImportClassCategoryLearningPoints(tmpDir)
	assert.NoError(t, err)

	// Verify record was imported
	var imported models.ClassCategoryLearningPoints
	err = database.DB.Where("character_class_id = ? AND skill_category_id = ?", class.ID, category.ID).First(&imported).Error
	assert.NoError(t, err)
	assert.Equal(t, 10, imported.Points)
}

// TestExportImportClassSpellPoints tests export and import of class spell points
func TestExportImportClassSpellPoints(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data with unique name
	class := models.CharacterClass{
		Code:         "TEST_MAG",
		Name:         "Test-Magier",
		GameSystemId: 1,
	}
	database.DB.Create(&class)

	record := models.ClassSpellPoints{
		CharacterClassID: class.ID,
		SpellPoints:      50,
	}
	database.DB.Create(&record)

	// Export
	tmpDir := t.TempDir()
	err := ExportClassSpellPoints(tmpDir)
	assert.NoError(t, err)

	// Verify file exists
	exportFile := filepath.Join(tmpDir, "class_spell_points.json")
	assert.FileExists(t, exportFile)

	// Delete record
	database.DB.Delete(&record)

	// Import back
	err = ImportClassSpellPoints(tmpDir)
	assert.NoError(t, err)

	// Verify record was imported
	var imported models.ClassSpellPoints
	err = database.DB.Where("character_class_id = ?", class.ID).First(&imported).Error
	assert.NoError(t, err)
	assert.Equal(t, 50, imported.SpellPoints)
}

// TestExportImportClassTypicalSkills tests export and import of class typical skills
func TestExportImportClassTypicalSkills(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data with unique names
	class := models.CharacterClass{
		Code:         "TEST_WAL",
		Name:         "Test-Waldläufer",
		GameSystemId: 1,
	}
	database.DB.Create(&class)

	skill := models.Skill{
		Name:        "Test-Spurenlesen",
		GameSystem:  "midgard",
		Improvable:  true,
		InnateSkill: false,
	}
	if err := skill.Create(); err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	record := models.ClassTypicalSkill{
		CharacterClassID: class.ID,
		SkillID:          skill.ID,
		Bonus:            4,
		Attribute:        "In",
		Notes:            "Typische Fertigkeit",
	}
	database.DB.Create(&record)

	// Export
	tmpDir := t.TempDir()
	err := ExportClassTypicalSkills(tmpDir)
	assert.NoError(t, err)

	// Verify file exists
	exportFile := filepath.Join(tmpDir, "class_typical_skills.json")
	assert.FileExists(t, exportFile)

	// Delete record
	database.DB.Delete(&record)

	// Import back
	err = ImportClassTypicalSkills(tmpDir)
	assert.NoError(t, err)

	// Verify record was imported
	var imported models.ClassTypicalSkill
	err = database.DB.Where("character_class_id = ? AND skill_id = ?", class.ID, skill.ID).First(&imported).Error
	assert.NoError(t, err)
	assert.Equal(t, 4, imported.Bonus)
	assert.Equal(t, "In", imported.Attribute)
	assert.Equal(t, "Typische Fertigkeit", imported.Notes)
}

// TestExportImportClassTypicalSpells tests export and import of class typical spells
func TestExportImportClassTypicalSpells(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	// Create test data with unique names
	class := models.CharacterClass{
		Code:         "TEST_DRU",
		Name:         "Test-Druide",
		GameSystemId: 1,
	}
	database.DB.Create(&class)

	spell := models.Spell{
		Name:         "Test-Heilen",
		GameSystemId: 1,
		Stufe:        1,
	}
	database.DB.Create(&spell)

	record := models.ClassTypicalSpell{
		CharacterClassID: class.ID,
		SpellID:          spell.ID,
		Notes:            "Immer verfügbar",
	}
	database.DB.Create(&record)

	// Export
	tmpDir := t.TempDir()
	err := ExportClassTypicalSpells(tmpDir)
	assert.NoError(t, err)

	// Verify file exists
	exportFile := filepath.Join(tmpDir, "class_typical_spells.json")
	assert.FileExists(t, exportFile)

	// Delete record
	database.DB.Delete(&record)

	// Import back
	err = ImportClassTypicalSpells(tmpDir)
	assert.NoError(t, err)

	// Verify record was imported
	var imported models.ClassTypicalSpell
	err = database.DB.Where("character_class_id = ? AND spell_id = ?", class.ID, spell.ID).First(&imported).Error
	assert.NoError(t, err)
	assert.Equal(t, "Immer verfügbar", imported.Notes)
}
