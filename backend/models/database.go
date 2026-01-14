package models

import (
	"bamort/database"

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

	err := gsMasterMigrateStructure(targetDB)
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
		&SkillImprovementCost{},
		&ClassCategoryLearningPoints{},
		&ClassSpellPoints{},
		&ClassTypicalSkill{},
		&ClassTypicalSpell{},
		&AuditLogEntry{},
	)
	if err != nil {
		return err
	}
	return nil
}
