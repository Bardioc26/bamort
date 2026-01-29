package models

/*
 * Database Migration Tests
 *
 * This file contains comprehensive tests for database.go migration functions.
 * Tests cover:
 * - MigrateStructure and all sub-migration functions
 * - Database parameter handling (default DB, custom DB, nil handling)
 * - Table structure verification for all model categories
 * - Migration idempotency (can be run multiple times safely)
 * - Integration workflow testing
 * - Error resilience and edge cases
 * - Database consistency and relationship verification
 */

import (
	"bamort/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupDatabaseTestDB(t *testing.T) {
	database.SetupTestDB()
}

// =============================================================================
// Tests for MigrateStructure function
// =============================================================================

func TestMigrateStructure_DefaultDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := MigrateStructure()

	assert.NoError(t, err, "MigrateStructure should succeed with default database")
}

func TestMigrateStructure_CustomDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	// Use the same database instance for consistency
	err := MigrateStructure(database.DB)

	assert.NoError(t, err, "MigrateStructure should succeed with custom database")
}

func TestMigrateStructure_NilDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	// Pass nil DB - should fall back to default
	err := MigrateStructure(nil)

	assert.NoError(t, err, "MigrateStructure should succeed with nil database (fallback to default)")
}

func TestMigrateStructure_VerifyTablesCreated(t *testing.T) {
	setupDatabaseTestDB(t)

	err := MigrateStructure()
	require.NoError(t, err, "MigrateStructure should succeed")

	// Verify that key tables exist by checking if we can perform basic operations
	// Test each migration category

	// Test gsmaster structures
	var skill Skill
	err = database.DB.First(&skill).Error
	// Error is expected if no records exist, but table should exist
	assert.True(t, err == nil || err == gorm.ErrRecordNotFound, "Skills table should exist")

	var spell Spell
	err = database.DB.First(&spell).Error
	assert.True(t, err == nil || err == gorm.ErrRecordNotFound, "Spells table should exist")

	// Test character structures
	var char Char
	err = database.DB.First(&char).Error
	assert.True(t, err == nil || err == gorm.ErrRecordNotFound, "Characters table should exist")

	// Test equipment structures
	var equipment EqAusruestung
	err = database.DB.First(&equipment).Error
	assert.True(t, err == nil || err == gorm.ErrRecordNotFound, "Equipment table should exist")

	// Test skills structures
	var skFertigkeit SkFertigkeit
	err = database.DB.First(&skFertigkeit).Error
	assert.True(t, err == nil || err == gorm.ErrRecordNotFound, "Skill skills table should exist")

	// Test learning structures
	var source Source
	err = database.DB.First(&source).Error
	assert.True(t, err == nil || err == gorm.ErrRecordNotFound, "Learning sources table should exist")
}

// =============================================================================
// Tests for gsMasterMigrateStructure function
// =============================================================================

func TestGsMasterMigrateStructure_DefaultDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := gsMasterMigrateStructure()

	assert.NoError(t, err, "gsMasterMigrateStructure should succeed with default database")
}

func TestGsMasterMigrateStructure_CustomDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := gsMasterMigrateStructure(database.DB)

	assert.NoError(t, err, "gsMasterMigrateStructure should succeed with custom database")
}

func TestGsMasterMigrateStructure_VerifyStructures(t *testing.T) {
	setupDatabaseTestDB(t)

	err := gsMasterMigrateStructure()
	require.NoError(t, err, "gsMasterMigrateStructure should succeed")

	// Verify all gsmaster structures can be accessed
	structures := []interface{}{
		&Skill{},
		&WeaponSkill{},
		&Spell{},
		&Equipment{},
		&Weapon{},
		&Container{},
		&Transportation{},
		&Believe{},
	}

	for _, structure := range structures {
		err = database.DB.First(structure).Error
		assert.True(t, err == nil || err == gorm.ErrRecordNotFound,
			"Structure %T table should exist and be accessible", structure)
	}
}

// =============================================================================
// Tests for characterMigrateStructure function
// =============================================================================

func TestCharacterMigrateStructure_DefaultDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := characterMigrateStructure()

	assert.NoError(t, err, "characterMigrateStructure should succeed with default database")
}

func TestCharacterMigrateStructure_CustomDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := characterMigrateStructure(database.DB)

	assert.NoError(t, err, "characterMigrateStructure should succeed with custom database")
}

func TestCharacterMigrateStructure_VerifyStructures(t *testing.T) {
	setupDatabaseTestDB(t)

	err := characterMigrateStructure()
	require.NoError(t, err, "characterMigrateStructure should succeed")

	// Verify all character structures can be accessed
	structures := []interface{}{
		&Char{},
		&Eigenschaft{},
		&Lp{},
		&Ap{},
		&B{},
		&Merkmale{},
		&Erfahrungsschatz{},
		&Bennies{},
		&Vermoegen{},
		&CharacterCreationSession{},
	}

	for _, structure := range structures {
		err = database.DB.First(structure).Error
		assert.True(t, err == nil || err == gorm.ErrRecordNotFound,
			"Structure %T table should exist and be accessible", structure)
	}
}

// =============================================================================
// Tests for equipmentMigrateStructure function
// =============================================================================

func TestEquipmentMigrateStructure_DefaultDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := equipmentMigrateStructure()

	assert.NoError(t, err, "equipmentMigrateStructure should succeed with default database")
}

func TestEquipmentMigrateStructure_CustomDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := equipmentMigrateStructure(database.DB)

	assert.NoError(t, err, "equipmentMigrateStructure should succeed with custom database")
}

func TestEquipmentMigrateStructure_VerifyStructures(t *testing.T) {
	setupDatabaseTestDB(t)

	err := equipmentMigrateStructure()
	require.NoError(t, err, "equipmentMigrateStructure should succeed")

	// Verify all equipment structures can be accessed
	structures := []interface{}{
		&EqAusruestung{},
		&EqWaffe{},
		&EqContainer{},
	}

	for _, structure := range structures {
		err = database.DB.First(structure).Error
		assert.True(t, err == nil || err == gorm.ErrRecordNotFound,
			"Structure %T table should exist and be accessible", structure)
	}
}

// =============================================================================
// Tests for skillsMigrateStructure function
// =============================================================================

func TestSkillsMigrateStructure_DefaultDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := skillsMigrateStructure()

	assert.NoError(t, err, "skillsMigrateStructure should succeed with default database")
}

func TestSkillsMigrateStructure_CustomDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := skillsMigrateStructure(database.DB)

	assert.NoError(t, err, "skillsMigrateStructure should succeed with custom database")
}

func TestSkillsMigrateStructure_VerifyStructures(t *testing.T) {
	setupDatabaseTestDB(t)

	err := skillsMigrateStructure()
	require.NoError(t, err, "skillsMigrateStructure should succeed")

	// Verify all skills structures can be accessed
	structures := []interface{}{
		&SkFertigkeit{},
		&SkWaffenfertigkeit{},
		&SkZauber{},
	}

	for _, structure := range structures {
		err = database.DB.First(structure).Error
		assert.True(t, err == nil || err == gorm.ErrRecordNotFound,
			"Structure %T table should exist and be accessible", structure)
	}
}

// =============================================================================
// Tests for learningMigrateStructure function
// =============================================================================

func TestLearningMigrateStructure_DefaultDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := learningMigrateStructure()

	assert.NoError(t, err, "learningMigrateStructure should succeed with default database")
}

func TestLearningMigrateStructure_CustomDB_Success(t *testing.T) {
	setupDatabaseTestDB(t)

	err := learningMigrateStructure(database.DB)

	assert.NoError(t, err, "learningMigrateStructure should succeed with custom database")
}

func TestLearningMigrateStructure_VerifyStructures(t *testing.T) {
	setupDatabaseTestDB(t)

	err := learningMigrateStructure()
	require.NoError(t, err, "learningMigrateStructure should succeed")

	// Verify all learning structures can be accessed
	structures := []interface{}{
		&Source{},
		&CharacterClass{},
		&SkillCategory{},
		&SkillDifficulty{},
		&SpellSchool{},
		&ClassCategoryEPCost{},
		&ClassSpellSchoolEPCost{},
		&SpellLevelLECost{},
		&SkillCategoryDifficulty{},
		&SkillImprovementCost{},
		&AuditLogEntry{},
	}

	for _, structure := range structures {
		err = database.DB.First(structure).Error
		assert.True(t, err == nil || err == gorm.ErrRecordNotFound,
			"Structure %T table should exist and be accessible", structure)
	}
}

// =============================================================================
// Integration tests for complete migration workflow
// =============================================================================

func TestMigrationWorkflow_CompleteFlow(t *testing.T) {
	setupDatabaseTestDB(t)

	// Test that all migration functions can be called in sequence without errors
	err := gsMasterMigrateStructure()
	assert.NoError(t, err, "gsMasterMigrateStructure should succeed")

	err = characterMigrateStructure()
	assert.NoError(t, err, "characterMigrateStructure should succeed")

	err = equipmentMigrateStructure()
	assert.NoError(t, err, "equipmentMigrateStructure should succeed")

	err = skillsMigrateStructure()
	assert.NoError(t, err, "skillsMigrateStructure should succeed")

	err = learningMigrateStructure()
	assert.NoError(t, err, "learningMigrateStructure should succeed")

	// Finally run the complete migration
	err = MigrateStructure()
	assert.NoError(t, err, "MigrateStructure should succeed after individual migrations")
}

func TestMigrationWorkflow_Idempotency(t *testing.T) {
	setupDatabaseTestDB(t)

	// Test that running migrations multiple times doesn't cause errors
	err := MigrateStructure()
	assert.NoError(t, err, "First MigrateStructure should succeed")

	err = MigrateStructure()
	assert.NoError(t, err, "Second MigrateStructure should succeed (idempotent)")

	err = MigrateStructure()
	assert.NoError(t, err, "Third MigrateStructure should succeed (idempotent)")
}

func TestMigrationWorkflow_StructureIntegrity(t *testing.T) {
	setupDatabaseTestDB(t)

	err := MigrateStructure()
	require.NoError(t, err, "MigrateStructure should succeed")

	// Test that we can create and retrieve basic records for each major category

	// Test gsmaster category with a skill
	testSkill := &Skill{
		GameSystem:  "midgard",
		Name:        "Test Skill",
		Category:    "Test Category",
		Difficulty:  "normal",
		Initialwert: 5,
		Improvable:  true,
	}
	err = database.DB.Create(testSkill).Error
	assert.NoError(t, err, "Should be able to create a skill")

	var retrievedSkill Skill
	err = database.DB.Where("name = ?", "Test Skill").First(&retrievedSkill).Error
	assert.NoError(t, err, "Should be able to retrieve created skill")
	assert.Equal(t, "Test Skill", retrievedSkill.Name)

	// Test learning category with a source
	testSource := &Source{
		Code:         "TEST",
		Name:         "Test Source",
		GameSystemId: 1,
		IsActive:     true,
	}
	err = database.DB.Create(testSource).Error
	assert.NoError(t, err, "Should be able to create a source")

	var retrievedSource Source
	err = database.DB.Where("code = ?", "TEST").First(&retrievedSource).Error
	assert.NoError(t, err, "Should be able to retrieve created source")
	assert.Equal(t, "TEST", retrievedSource.Code)
}

// =============================================================================
// Tests for database parameter handling
// =============================================================================

func TestMigrationFunctions_ParameterHandling(t *testing.T) {
	setupDatabaseTestDB(t)

	// Test that all functions handle both default and custom DB parameters correctly
	functions := []struct {
		name     string
		function func(...*gorm.DB) error
	}{
		{"MigrateStructure", MigrateStructure},
		{"gsMasterMigrateStructure", gsMasterMigrateStructure},
		{"characterMigrateStructure", characterMigrateStructure},
		{"equipmentMigrateStructure", equipmentMigrateStructure},
		{"skillsMigrateStructure", skillsMigrateStructure},
		{"learningMigrateStructure", learningMigrateStructure},
	}

	for _, fn := range functions {
		t.Run(fn.name, func(t *testing.T) {
			// Test with no parameters (default DB)
			err := fn.function()
			assert.NoError(t, err, "%s should work with default DB", fn.name)

			// Test with explicit DB parameter
			err = fn.function(database.DB)
			assert.NoError(t, err, "%s should work with explicit DB", fn.name)

			// Test with nil parameter (should fall back to default)
			err = fn.function(nil)
			assert.NoError(t, err, "%s should work with nil DB (fallback)", fn.name)
		})
	}
}

// =============================================================================
// Error handling and edge case tests
// =============================================================================

func TestMigration_ErrorResilience(t *testing.T) {
	setupDatabaseTestDB(t)

	// Test that migration can handle being called multiple times
	for i := 0; i < 3; i++ {
		err := MigrateStructure()
		assert.NoError(t, err, "Migration attempt %d should succeed", i+1)
	}

	// Test individual migration functions multiple times
	functions := []func(...*gorm.DB) error{
		gsMasterMigrateStructure,
		characterMigrateStructure,
		equipmentMigrateStructure,
		skillsMigrateStructure,
		learningMigrateStructure,
	}

	for _, fn := range functions {
		for i := 0; i < 2; i++ {
			err := fn()
			assert.NoError(t, err, "Function should be callable multiple times")
		}
	}
}

func TestMigration_DatabaseConsistency(t *testing.T) {
	setupDatabaseTestDB(t)

	// Run complete migration
	err := MigrateStructure()
	require.NoError(t, err, "MigrateStructure should succeed")

	// Verify that database state is consistent by checking if we can perform
	// cross-references between different structure categories

	// Create a character and related data
	testChar := &Char{
		BamortBase: BamortBase{
			Name: "Test Character",
		},
		Typ:    "Kr",
		Grad:   1,
		Public: false,
	}
	err = database.DB.Create(testChar).Error
	assert.NoError(t, err, "Should be able to create character")

	// Create related skill
	testSkill := &SkFertigkeit{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "Test Character Skill",
			},
			CharacterID: testChar.ID,
		},
		Fertigkeitswert: 10,
	}
	err = database.DB.Create(testSkill).Error
	assert.NoError(t, err, "Should be able to create character skill")

	// Verify relationship by simply checking that the character exists
	var retrievedChar Char
	err = database.DB.Where("id = ?", testChar.ID).First(&retrievedChar).Error
	assert.NoError(t, err, "Character retrieval should work")
	assert.Equal(t, "Test Character", retrievedChar.Name, "Character name should match")
}
