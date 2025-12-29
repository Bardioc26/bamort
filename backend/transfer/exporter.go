package transfer

import (
	"bamort/database"
	"bamort/models"
	"fmt"
)

// CharacterExport contains all data needed to export and import a character
type CharacterExport struct {
	Character       models.Char            `json:"character"`
	GSMSkills       []models.Skill         `json:"gsm_skills"`
	GSMWeaponSkills []models.WeaponSkill   `json:"gsm_weapon_skills"`
	GSMSpells       []models.Spell         `json:"gsm_spells"`
	GSMWeapons      []models.Weapon        `json:"gsm_weapons"`
	GSMEquipment    []models.Equipment     `json:"gsm_equipment"`
	GSMContainers   []models.Container     `json:"gsm_containers"`
	LearningData    LearningDataExport     `json:"learning_data"`
	AuditLogEntries []models.AuditLogEntry `json:"audit_log_entries"`
}

// LearningDataExport contains all learning-related master data
type LearningDataExport struct {
	Sources                   []models.Source                  `json:"sources"`
	CharacterClasses          []models.CharacterClass          `json:"character_classes"`
	SkillCategories           []models.SkillCategory           `json:"skill_categories"`
	SkillDifficulties         []models.SkillDifficulty         `json:"skill_difficulties"`
	SpellSchools              []models.SpellSchool             `json:"spell_schools"`
	ClassCategoryEPCosts      []models.ClassCategoryEPCost     `json:"class_category_ep_costs"`
	ClassSpellSchoolEPCosts   []models.ClassSpellSchoolEPCost  `json:"class_spell_school_ep_costs"`
	SpellLevelLECosts         []models.SpellLevelLECost        `json:"spell_level_le_costs"`
	SkillCategoryDifficulties []models.SkillCategoryDifficulty `json:"skill_category_difficulties"`
	SkillImprovementCosts     []models.SkillImprovementCost    `json:"skill_improvement_costs"`
}

// ExportCharacter exports a complete character with all related data
func ExportCharacter(characterID uint) (*CharacterExport, error) {
	var char models.Char

	// Load character with all relations
	err := database.DB.
		Preload("User").
		Preload("Lp").
		Preload("Ap").
		Preload("B").
		Preload("Merkmale").
		Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Zauber").
		Preload("Bennies").
		Preload("Vermoegen").
		Preload("Erfahrungsschatz").
		Preload("Waffen").
		Preload("Behaeltnisse").
		Preload("Transportmittel").
		Preload("Ausruestung").
		First(&char, characterID).Error

	if err != nil {
		return nil, fmt.Errorf("failed to load character: %w", err)
	}

	export := &CharacterExport{
		Character: char,
	}

	// Collect GSM skill data
	export.GSMSkills = make([]models.Skill, 0)
	export.GSMWeaponSkills = make([]models.WeaponSkill, 0)

	skillNames := make(map[string]bool)
	for _, skill := range char.Fertigkeiten {
		if !skillNames[skill.Name] {
			var gsmSkill models.Skill
			err := gsmSkill.First(skill.Name)
			if err == nil && gsmSkill.ID != 0 {
				export.GSMSkills = append(export.GSMSkills, gsmSkill)
				skillNames[skill.Name] = true
			}
		}
	}

	weaponSkillNames := make(map[string]bool)
	for _, skill := range char.Waffenfertigkeiten {
		if !weaponSkillNames[skill.Name] {
			var weaponSkill models.WeaponSkill
			err := weaponSkill.First(skill.Name)
			if err == nil && weaponSkill.ID != 0 {
				export.GSMWeaponSkills = append(export.GSMWeaponSkills, weaponSkill)
				weaponSkillNames[skill.Name] = true
			}
		}
	}

	// Collect GSM spell data
	export.GSMSpells = make([]models.Spell, 0)
	spellNames := make(map[string]bool)
	for _, spell := range char.Zauber {
		if !spellNames[spell.Name] {
			var gsmSpell models.Spell
			err := gsmSpell.First(spell.Name)
			if err == nil && gsmSpell.ID != 0 {
				export.GSMSpells = append(export.GSMSpells, gsmSpell)
				spellNames[spell.Name] = true
			}
		}
	}

	// Collect GSM weapon data
	export.GSMWeapons = make([]models.Weapon, 0)
	weaponNames := make(map[string]bool)
	for _, weapon := range char.Waffen {
		if !weaponNames[weapon.Name] {
			var gsmWeapon models.Weapon
			err := gsmWeapon.First(weapon.Name)
			if err == nil && gsmWeapon.ID != 0 {
				export.GSMWeapons = append(export.GSMWeapons, gsmWeapon)
				weaponNames[weapon.Name] = true
			}
		}
	}

	// Collect GSM equipment data
	export.GSMEquipment = make([]models.Equipment, 0)
	equipmentNames := make(map[string]bool)
	for _, equip := range char.Ausruestung {
		if !equipmentNames[equip.Name] {
			var gsmEquip models.Equipment
			err := gsmEquip.First(equip.Name)
			if err == nil && gsmEquip.ID != 0 {
				export.GSMEquipment = append(export.GSMEquipment, gsmEquip)
				equipmentNames[equip.Name] = true
			}
		}
	}

	// Collect GSM container data
	export.GSMContainers = make([]models.Container, 0)
	containerNames := make(map[string]bool)
	for _, container := range char.Behaeltnisse {
		if !containerNames[container.Name] {
			var gsmContainer models.Container
			err := gsmContainer.First(container.Name)
			if err == nil && gsmContainer.ID != 0 {
				export.GSMContainers = append(export.GSMContainers, gsmContainer)
				containerNames[container.Name] = true
			}
		}
	}
	for _, container := range char.Transportmittel {
		if !containerNames[container.Name] {
			var gsmContainer models.Container
			err := gsmContainer.First(container.Name)
			if err == nil && gsmContainer.ID != 0 {
				export.GSMContainers = append(export.GSMContainers, gsmContainer)
				containerNames[container.Name] = true
			}
		}
	}

	// Load learning data
	export.LearningData = LearningDataExport{
		Sources:                   make([]models.Source, 0),
		CharacterClasses:          make([]models.CharacterClass, 0),
		SkillCategories:           make([]models.SkillCategory, 0),
		SkillDifficulties:         make([]models.SkillDifficulty, 0),
		SpellSchools:              make([]models.SpellSchool, 0),
		ClassCategoryEPCosts:      make([]models.ClassCategoryEPCost, 0),
		ClassSpellSchoolEPCosts:   make([]models.ClassSpellSchoolEPCost, 0),
		SpellLevelLECosts:         make([]models.SpellLevelLECost, 0),
		SkillCategoryDifficulties: make([]models.SkillCategoryDifficulty, 0),
		SkillImprovementCosts:     make([]models.SkillImprovementCost, 0),
	}

	database.DB.Find(&export.LearningData.Sources)
	database.DB.Preload("Source").Find(&export.LearningData.CharacterClasses)
	database.DB.Preload("Source").Find(&export.LearningData.SkillCategories)
	database.DB.Find(&export.LearningData.SkillDifficulties)
	database.DB.Preload("Source").Find(&export.LearningData.SpellSchools)
	database.DB.Preload("CharacterClass").Preload("SkillCategory").Find(&export.LearningData.ClassCategoryEPCosts)
	database.DB.Preload("CharacterClass").Preload("SpellSchool").Find(&export.LearningData.ClassSpellSchoolEPCosts)
	database.DB.Find(&export.LearningData.SpellLevelLECosts)
	database.DB.Preload("Skill").Preload("SkillCategory").Preload("SkillDifficulty").Find(&export.LearningData.SkillCategoryDifficulties)
	database.DB.Preload("SkillCategoryDifficulty").Find(&export.LearningData.SkillImprovementCosts)

	// Load audit log entries
	export.AuditLogEntries = make([]models.AuditLogEntry, 0)
	database.DB.Where("character_id = ?", characterID).Order("timestamp ASC").Find(&export.AuditLogEntries)

	return export, nil
}
