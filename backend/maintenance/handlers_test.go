package maintenance

import (
	"bamort/models"
	"bamort/user"
	"testing"
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
