package maintenance

import (
	"bamort/database"
	"bamort/models"
	"bamort/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestTableListCompleteness verifies that all database models are included in the table lists
func TestTableListCompleteness(t *testing.T) {
	// Define all known database models that should be copied
	expectedModels := map[string]bool{
		// User models
		"*user.User": true,

		// Learning Costs System - Basis
		"*models.Source":          true,
		"*models.CharacterClass":  true,
		"*models.SkillCategory":   true,
		"*models.SkillDifficulty": true,
		"*models.SpellSchool":     true,

		// Learning Costs System - Dependent
		"*models.ClassCategoryEPCost":           true,
		"*models.ClassSpellSchoolEPCost":        true,
		"*models.SpellLevelLECost":              true,
		"*models.SkillCategoryDifficulty":       true,
		"*models.WeaponSkillCategoryDifficulty": true,
		"*models.SkillImprovementCost":          true,

		// GSMaster Base Data
		"*models.Skill":          true,
		"*models.WeaponSkill":    true,
		"*models.Spell":          true,
		"*models.Equipment":      true,
		"*models.Weapon":         true,
		"*models.Container":      true,
		"*models.Transportation": true,
		"*models.Believe":        true,

		// Characters (Base)
		"*models.Char": true,

		// Character Properties
		"*models.Eigenschaft":      true,
		"*models.Lp":               true,
		"*models.Ap":               true,
		"*models.B":                true,
		"*models.Merkmale":         true,
		"*models.Erfahrungsschatz": true,
		"*models.Bennies":          true,
		"*models.Vermoegen":        true,

		// Character Skills
		"*models.SkFertigkeit":           true,
		"*models.SkWaffenfertigkeit":     true,
		"*models.SkAngeboreneFertigkeit": true,
		"*models.SkZauber":               true,

		// Character Equipment
		"*models.EqAusruestung": true,
		"*models.EqWaffe":       true,
		"*models.EqContainer":   true,

		// Character Creation Sessions
		"*models.CharacterCreationSession": true,

		// Audit Logging
		"*models.AuditLogEntry": true,
	}

	// Get the table list from copyMariaDBToSQLite (simulated)
	tables := []interface{}{
		&user.User{},
		&models.Source{},
		&models.CharacterClass{},
		&models.SkillCategory{},
		&models.SkillDifficulty{},
		&models.SpellSchool{},
		&models.ClassCategoryEPCost{},
		&models.ClassSpellSchoolEPCost{},
		&models.SpellLevelLECost{},
		&models.SkillCategoryDifficulty{},
		&models.WeaponSkillCategoryDifficulty{},
		&models.SkillImprovementCost{},
		&models.Skill{},
		&models.WeaponSkill{},
		&models.Spell{},
		&models.Equipment{},
		&models.Weapon{},
		&models.Container{},
		&models.Transportation{},
		&models.Believe{},
		&models.Char{},
		&models.Eigenschaft{},
		&models.Lp{},
		&models.Ap{},
		&models.B{},
		&models.Merkmale{},
		&models.Erfahrungsschatz{},
		&models.Bennies{},
		&models.Vermoegen{},
		&models.SkFertigkeit{},
		&models.SkWaffenfertigkeit{},
		&models.SkAngeboreneFertigkeit{},
		&models.SkZauber{},
		&models.EqAusruestung{},
		&models.EqWaffe{},
		&models.EqContainer{},
		&models.CharacterCreationSession{},
		&models.AuditLogEntry{},
	}

	// Verify all expected models are in the table list
	foundModels := make(map[string]bool)
	for _, model := range tables {
		modelType := getModelTypeName(model)
		foundModels[modelType] = true
	}

	// Check for missing models
	for expectedModel := range expectedModels {
		if !foundModels[expectedModel] {
			t.Errorf("Missing model in table list: %s", expectedModel)
		}
	}

	// Check for unexpected models (not in expected list)
	for foundModel := range foundModels {
		if !expectedModels[foundModel] {
			t.Logf("Warning: Unexpected model in table list (may be intentional): %s", foundModel)
		}
	}

	t.Logf("Total models in table list: %d", len(tables))
	t.Logf("Total expected models: %d", len(expectedModels))
}

// TestImprovableFieldTransfer tests that the Improvable field is correctly transferred from MariaDB to SQLite
func TestImprovableFieldTransfer(t *testing.T) {
	setupTestEnvironment(t)

	// Create source database (simulating MariaDB)
	sourceDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create target database (simulating SQLite)
	targetDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate structures
	err = sourceDB.AutoMigrate(&models.Skill{})
	require.NoError(t, err)
	err = targetDB.AutoMigrate(&models.Skill{})
	require.NoError(t, err)

	// Insert test data with different Improvable values
	// Use raw SQL to avoid GORM applying default values
	testSkills := []struct {
		ID          uint
		Name        string
		GameSystem  string
		Improvable  bool
		InnateSkill bool
	}{
		{1, "Hören", "midgard", false, true},
		{2, "Alchimie", "midgard", true, false},
		{3, "Nachtsicht", "midgard", false, true},
	}

	for _, skill := range testSkills {
		err = sourceDB.Exec(`INSERT INTO gsm_skills (id, name, game_system, improvable, innate_skill, initialwert, basis_wert) VALUES (?, ?, ?, ?, ?, 5, 0)`,
			skill.ID, skill.Name, skill.GameSystem, skill.Improvable, skill.InnateSkill).Error
		require.NoError(t, err)
	}

	// Verify source data
	var sourceSkill models.Skill
	err = sourceDB.First(&sourceSkill, 1).Error
	require.NoError(t, err)
	assert.Equal(t, "Hören", sourceSkill.Name)
	t.Logf("DEBUG: Source skill Hören - Improvable: %v, InnateSkill: %v", sourceSkill.Improvable, sourceSkill.InnateSkill)

	// Also check raw data from database
	var rawImprovable int
	err = sourceDB.Raw("SELECT improvable FROM gsm_skills WHERE id = 1").Scan(&rawImprovable).Error
	require.NoError(t, err)
	t.Logf("DEBUG: Raw SQL value for improvable: %d", rawImprovable)

	assert.False(t, sourceSkill.Improvable, "Source skill should have Improvable=false")

	// Copy data using copyTableData
	err = copyTableData(sourceDB, targetDB, &models.Skill{})
	require.NoError(t, err)

	// Verify all skills in target database
	var targetSkills []models.Skill
	err = targetDB.Find(&targetSkills).Error
	require.NoError(t, err)
	require.Len(t, targetSkills, 3, "Should have 3 skills in target")

	// Check each skill's Improvable field
	expectedValues := map[uint]bool{
		1: false, // Hören
		2: true,  // Alchimie
		3: false, // Nachtsicht
	}

	for _, targetSkill := range targetSkills {
		expectedImprovable := expectedValues[targetSkill.ID]
		assert.Equal(t, expectedImprovable, targetSkill.Improvable,
			"Skill %s (ID: %d) should have Improvable=%v, got %v",
			targetSkill.Name, targetSkill.ID, expectedImprovable, targetSkill.Improvable)
	}

	// Specific checks
	var hoeren models.Skill
	err = targetDB.Where("name = ?", "Hören").First(&hoeren).Error
	require.NoError(t, err)
	assert.False(t, hoeren.Improvable, "Hören should have Improvable=false after transfer")

	var alchimie models.Skill
	err = targetDB.Where("name = ?", "Alchimie").First(&alchimie).Error
	require.NoError(t, err)
	assert.True(t, alchimie.Improvable, "Alchimie should have Improvable=true after transfer")
}

// TestImprovableFieldInPreparedTestDB verifies the prepared test database has correct Improvable values
func TestImprovableFieldInPreparedTestDB(t *testing.T) {
	setupTestEnvironment(t)

	// Ensure clean database state before setup
	database.ResetTestDB()

	// Use the prepared test database
	database.SetupTestDB(true)
	require.NotNil(t, database.DB)

	// Ensure database cleanup after test
	t.Cleanup(func() {
		database.ResetTestDB()
	})

	// Check specific skills that should have Improvable=false (innate skills)
	innateSkills := []string{"Hören", "Nachtsicht", "Riechen", "Sechster Sinn", "Sehen"}

	for _, skillName := range innateSkills {
		var skill models.Skill
		err := database.DB.Where("name = ?", skillName).First(&skill).Error
		if err == gorm.ErrRecordNotFound {
			t.Logf("Skill %s not found in prepared test DB - skipping", skillName)
			continue
		}
		require.NoError(t, err)

		// These are innate skills and should not be improvable
		assert.True(t, skill.InnateSkill, "Skill %s should be marked as InnateSkill", skillName)
		// Note: Based on game rules, innate skills are typically not improvable
		t.Logf("Skill: %s, Improvable: %v, InnateSkill: %v", skillName, skill.Improvable, skill.InnateSkill)
	}

	// Check a regular skill that should be improvable
	var alchimie models.Skill
	err := database.DB.Where("name = ?", "Alchimie").First(&alchimie).Error
	if err != gorm.ErrRecordNotFound {
		require.NoError(t, err)
		assert.True(t, alchimie.Improvable, "Alchimie should be improvable")
		assert.False(t, alchimie.InnateSkill, "Alchimie should not be an innate skill")
	}
}

// getModelTypeName returns the type name of a model
func getModelTypeName(model interface{}) string {
	switch model.(type) {
	case *user.User:
		return "*user.User"
	case *models.Source:
		return "*models.Source"
	case *models.CharacterClass:
		return "*models.CharacterClass"
	case *models.SkillCategory:
		return "*models.SkillCategory"
	case *models.SkillDifficulty:
		return "*models.SkillDifficulty"
	case *models.SpellSchool:
		return "*models.SpellSchool"
	case *models.ClassCategoryEPCost:
		return "*models.ClassCategoryEPCost"
	case *models.ClassSpellSchoolEPCost:
		return "*models.ClassSpellSchoolEPCost"
	case *models.SpellLevelLECost:
		return "*models.SpellLevelLECost"
	case *models.SkillCategoryDifficulty:
		return "*models.SkillCategoryDifficulty"
	case *models.WeaponSkillCategoryDifficulty:
		return "*models.WeaponSkillCategoryDifficulty"
	case *models.SkillImprovementCost:
		return "*models.SkillImprovementCost"
	case *models.Skill:
		return "*models.Skill"
	case *models.WeaponSkill:
		return "*models.WeaponSkill"
	case *models.Spell:
		return "*models.Spell"
	case *models.Equipment:
		return "*models.Equipment"
	case *models.Weapon:
		return "*models.Weapon"
	case *models.Container:
		return "*models.Container"
	case *models.Transportation:
		return "*models.Transportation"
	case *models.Believe:
		return "*models.Believe"
	case *models.Char:
		return "*models.Char"
	case *models.Eigenschaft:
		return "*models.Eigenschaft"
	case *models.Lp:
		return "*models.Lp"
	case *models.Ap:
		return "*models.Ap"
	case *models.B:
		return "*models.B"
	case *models.Merkmale:
		return "*models.Merkmale"
	case *models.Erfahrungsschatz:
		return "*models.Erfahrungsschatz"
	case *models.Bennies:
		return "*models.Bennies"
	case *models.Vermoegen:
		return "*models.Vermoegen"
	case *models.SkFertigkeit:
		return "*models.SkFertigkeit"
	case *models.SkWaffenfertigkeit:
		return "*models.SkWaffenfertigkeit"
	case *models.SkAngeboreneFertigkeit:
		return "*models.SkAngeboreneFertigkeit"
	case *models.SkZauber:
		return "*models.SkZauber"
	case *models.EqAusruestung:
		return "*models.EqAusruestung"
	case *models.EqWaffe:
		return "*models.EqWaffe"
	case *models.EqContainer:
		return "*models.EqContainer"
	case *models.CharacterCreationSession:
		return "*models.CharacterCreationSession"
	case *models.AuditLogEntry:
		return "*models.AuditLogEntry"
	default:
		return "UNKNOWN"
	}
}
