package models

import (
	"bamort/database"
	"bamort/logger"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func MigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := gameSystemMigrateStructure(targetDB)
	if err != nil {
		return err
	}
	err = gsMasterMigrateStructure(targetDB)
	if err != nil {
		return err
	}
	err = characterMigrateStructure(targetDB)
	if err != nil {
		return err
	}
	err = equipmentMigrateStructure(targetDB)
	if err != nil {
		return err
	}
	err = skillsMigrateStructure(targetDB)
	if err != nil {
		return err
	}
	err = importerMigrateStructure(targetDB)
	if err != nil {
		return err
	}
	err = learningMigrateStructure(targetDB)
	if err != nil {
		return err
	}

	return nil
}

func gameSystemMigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := targetDB.AutoMigrate(
		&GameSystem{},
	)
	if err != nil {
		return err
	}
	return nil
}

func gsMasterMigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := targetDB.AutoMigrate(
		&Skill{},
		&WeaponSkill{},
		&Spell{},
		&Equipment{},
		&Weapon{},
		&Container{},
		&Transportation{},
		&Believe{},
		&MiscLookup{},
	)
	if err != nil {
		return err
	}
	return nil
}

func characterMigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := targetDB.AutoMigrate(
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
		&CharShare{},
	)
	if err != nil {
		return err
	}
	return nil
}

func equipmentMigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := targetDB.AutoMigrate(
		&EqAusruestung{},
		&EqWaffe{},
		&EqContainer{},
		//Transportation{},
	)
	if err != nil {
		return err
	}
	return nil
}

func skillsMigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := targetDB.AutoMigrate(
		&SkFertigkeit{},
		&SkWaffenfertigkeit{},
		&SkZauber{},
	)
	if err != nil {
		return err
	}
	return nil
}

func learningMigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := targetDB.AutoMigrate(
		&Source{},
		&CharacterClass{},
		&SkillCategory{},
		&SkillDifficulty{},
		&SpellSchool{},
		&ClassCategoryEPCost{},
		&ClassSpellSchoolEPCost{},
		&SpellLevelLECost{},
		&SkillCategoryDifficulty{},
		&WeaponSkillCategoryDifficulty{},
		&ClassCategoryLearningPoints{},
		&ClassSpellPoints{},
		&ClassTypicalSkill{},
		&ClassTypicalSpell{},
		&AuditLogEntry{},
		&SkillImprovementCost{},
	)
	if err != nil {
		return err
	}
	return nil
}

func MigrateDataIfNeeded(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}
	tables := []string{
		"gsm_lit_sources",
		"gsm_believes",
		"gsm_equipments",
		"gsm_weapons",
		"gsm_containers",
		"gsm_transportations",
		"gsm_skills",
		"gsm_weaponskills",
		"gsm_spells",
		"char_chars",
		"gsm_character_classes",
		"gsm_misc",
		"learning_skill_categories",
		"learning_skill_difficulties",
		"learning_spell_level_le_costs",
		"learning_spell_schools",
	}
	sql := `
		UPDATE TABLE 
		SET game_system_id = 1
		WHERE game_system_id IS NULL or game_system_id = 0;	
	`
	logger.Debug("Führe SQL-Update aus: %s", strings.ReplaceAll(sql, "\n", " "))
	errors := []error{}
	for _, currentTable := range tables {
		//currentTable := tables[tables]
		execSQL := strings.ReplaceAll(sql, "TABLE", currentTable)
		result := targetDB.Exec(execSQL)
		if result.Error != nil {
			logger.Error("Fehler beim SQL-Update der Tabelle %s: %s", currentTable, result.Error.Error())
			errors = append(errors, fmt.Errorf("failed to update table %s: %w", currentTable, result.Error))
		}
	}
	if len(errors) == 0 {
		return nil
	}
	return fmt.Errorf("The following Errors occured when migrating data: %v", errors)
}
