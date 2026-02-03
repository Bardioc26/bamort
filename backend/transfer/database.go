package transfer

import (
	"bamort/appsystem"
	"bamort/database"
	"bamort/models"
	"bamort/user"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

// DatabaseExport contains all database tables for export/import
type DatabaseExport struct {
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`

	// User data
	Users []user.User `json:"users"`

	// Character data
	Characters                []models.Char                     `json:"characters"`
	Eigenschaften             []models.Eigenschaft              `json:"eigenschaften"`
	Lps                       []models.Lp                       `json:"lps"`
	Aps                       []models.Ap                       `json:"aps"`
	Bs                        []models.B                        `json:"bs"`
	Merkmale                  []models.Merkmale                 `json:"merkmale"`
	Erfahrungsschatze         []models.Erfahrungsschatz         `json:"erfahrungsschatze"`
	Bennies                   []models.Bennies                  `json:"bennies"`
	Vermoegen                 []models.Vermoegen                `json:"vermoegen"`
	CharacterCreationSessions []models.CharacterCreationSession `json:"character_creation_sessions"`

	// Skills
	SkFertigkeiten       []models.SkFertigkeit       `json:"sk_fertigkeiten"`
	SkWaffenfertigkeiten []models.SkWaffenfertigkeit `json:"sk_waffenfertigkeiten"`
	SkZauber             []models.SkZauber           `json:"sk_zauber"`

	// Equipment
	EqAusruestungen []models.EqAusruestung `json:"eq_ausruestungen"`
	EqWaffen        []models.EqWaffe       `json:"eq_waffen"`
	EqContainers    []models.EqContainer   `json:"eq_containers"`

	// GSMaster data
	GsmSkills          []models.Skill          `json:"gsm_skills"`
	GsmWeaponSkills    []models.WeaponSkill    `json:"gsm_weapon_skills"`
	GsmSpells          []models.Spell          `json:"gsm_spells"`
	GsmEquipment       []models.Equipment      `json:"gsm_equipment"`
	GsmWeapons         []models.Weapon         `json:"gsm_weapons"`
	GsmContainers      []models.Container      `json:"gsm_containers"`
	GsmTransportations []models.Transportation `json:"gsm_transportations"`
	GsmBelieves        []models.Believe        `json:"gsm_believes"`
	Sources            []models.Source         `json:"gsm_lit_sources"`
	CharacterClasses   []models.CharacterClass `json:"gsm_character_classes"`

	// Learning data
	SkillCategories           []models.SkillCategory           `json:"learning_skill_categories"`
	SkillDifficulties         []models.SkillDifficulty         `json:"learning_skill_difficulties"`
	SpellSchools              []models.SpellSchool             `json:"learning_spell_schools"`
	ClassCategoryEPCosts      []models.ClassCategoryEPCost     `json:"learning_class_category_ep_costs"`
	ClassSpellSchoolEPCosts   []models.ClassSpellSchoolEPCost  `json:"learning_class_spell_school_ep_costs"`
	SpellLevelLECosts         []models.SpellLevelLECost        `json:"learning_spell_level_le_costs"`
	SkillCategoryDifficulties []models.SkillCategoryDifficulty `json:"learning_skill_category_difficulties"`
	SkillImprovementCosts     []models.SkillImprovementCost    `json:"learning_skill_improvement_costs"`
	AuditLogEntries           []models.AuditLogEntry           `json:"audit_log_entries"`
}

// ExportResult contains information about the export operation
type ExportResult struct {
	Filename    string `json:"filename"`
	FilePath    string `json:"filepath"`
	RecordCount int    `json:"record_count"`
	Timestamp   string `json:"timestamp"`
}

// ImportResult contains information about the import operation
type ImportResult struct {
	RecordCount int    `json:"record_count"`
	Timestamp   string `json:"timestamp"`
}

// ExportDatabase exports all database content to a JSON file
func ExportDatabase(exportDir string) (*ExportResult, error) {
	// Create export directory if it doesn't exist
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create export directory: %w", err)
	}

	export := DatabaseExport{
		Version:   appsystem.GetVersion(),
		Timestamp: time.Now(),
	}

	// Export all tables
	if err := database.DB.Find(&export.Users).Error; err != nil {
		return nil, fmt.Errorf("failed to export users: %w", err)
	}

	if err := database.DB.Find(&export.Characters).Error; err != nil {
		return nil, fmt.Errorf("failed to export characters: %w", err)
	}

	database.DB.Find(&export.Eigenschaften)
	database.DB.Find(&export.Lps)
	database.DB.Find(&export.Aps)
	database.DB.Find(&export.Bs)
	database.DB.Find(&export.Merkmale)
	database.DB.Find(&export.Erfahrungsschatze)
	database.DB.Find(&export.Bennies)
	database.DB.Find(&export.Vermoegen)
	database.DB.Find(&export.CharacterCreationSessions)

	database.DB.Find(&export.SkFertigkeiten)
	database.DB.Find(&export.SkWaffenfertigkeiten)
	database.DB.Find(&export.SkZauber)

	database.DB.Find(&export.EqAusruestungen)
	database.DB.Find(&export.EqWaffen)
	database.DB.Find(&export.EqContainers)

	database.DB.Find(&export.GsmSkills)
	database.DB.Find(&export.GsmWeaponSkills)
	database.DB.Find(&export.GsmSpells)
	database.DB.Find(&export.GsmEquipment)
	database.DB.Find(&export.GsmWeapons)
	database.DB.Find(&export.GsmContainers)
	database.DB.Find(&export.GsmTransportations)
	database.DB.Find(&export.GsmBelieves)

	database.DB.Find(&export.Sources)
	database.DB.Find(&export.CharacterClasses)
	database.DB.Find(&export.SkillCategories)
	database.DB.Find(&export.SkillDifficulties)
	database.DB.Find(&export.SpellSchools)
	database.DB.Find(&export.ClassCategoryEPCosts)
	database.DB.Find(&export.ClassSpellSchoolEPCosts)
	database.DB.Find(&export.SpellLevelLECosts)
	database.DB.Find(&export.SkillCategoryDifficulties)
	database.DB.Find(&export.SkillImprovementCosts)
	database.DB.Find(&export.AuditLogEntries)

	// Count total records
	recordCount := len(export.Users) + len(export.Characters) +
		len(export.Eigenschaften) + len(export.Lps) + len(export.Aps) +
		len(export.Bs) + len(export.Merkmale) + len(export.Erfahrungsschatze) +
		len(export.Bennies) + len(export.Vermoegen) +
		len(export.SkFertigkeiten) + len(export.SkWaffenfertigkeiten) + len(export.SkZauber) +
		len(export.EqAusruestungen) + len(export.EqWaffen) + len(export.EqContainers) +
		len(export.GsmSkills) + len(export.GsmWeaponSkills) + len(export.GsmSpells) +
		len(export.GsmEquipment) + len(export.GsmWeapons) + len(export.GsmContainers) +
		len(export.GsmTransportations) + len(export.GsmBelieves) +
		len(export.Sources) + len(export.CharacterClasses) + len(export.SkillCategories) +
		len(export.SkillDifficulties) + len(export.SpellSchools) +
		len(export.ClassCategoryEPCosts) + len(export.ClassSpellSchoolEPCosts) +
		len(export.SpellLevelLECosts) + len(export.SkillCategoryDifficulties) +
		len(export.SkillImprovementCosts) + len(export.AuditLogEntries) +
		len(export.CharacterCreationSessions)

	// Generate filename with timestamp
	filename := fmt.Sprintf("database_export_%s.json", time.Now().Format("20060102_150405"))
	filepath := filepath.Join(exportDir, filename)

	// Marshal to JSON
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal export data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write export file: %w", err)
	}

	return &ExportResult{
		Filename:    filename,
		FilePath:    filepath,
		RecordCount: recordCount,
		Timestamp:   time.Now().Format(time.RFC3339),
	}, nil
}

// ImportDatabase imports all database content from a JSON file
func ImportDatabase(filePath string) (*ImportResult, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	// Unmarshal JSON
	var export DatabaseExport
	if err := json.Unmarshal(data, &export); err != nil {
		return nil, fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	// Import all tables in transaction
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Import users (upsert to handle existing IDs)
		for _, item := range export.Users {
			if err := tx.Save(&item).Error; err != nil {
				return fmt.Errorf("failed to import user: %w", err)
			}
		}

		// Import characters (upsert)
		for _, item := range export.Characters {
			if err := tx.Save(&item).Error; err != nil {
				return fmt.Errorf("failed to import character: %w", err)
			}
		}

		// Import character-related data (upsert)
		for _, item := range export.Eigenschaften {
			tx.Save(&item)
		}
		for _, item := range export.Lps {
			tx.Save(&item)
		}
		for _, item := range export.Aps {
			tx.Save(&item)
		}
		for _, item := range export.Bs {
			tx.Save(&item)
		}
		for _, item := range export.Merkmale {
			tx.Save(&item)
		}
		for _, item := range export.Erfahrungsschatze {
			tx.Save(&item)
		}
		for _, item := range export.Bennies {
			tx.Save(&item)
		}
		for _, item := range export.Vermoegen {
			tx.Save(&item)
		}
		for _, item := range export.CharacterCreationSessions {
			tx.Save(&item)
		}

		// Import skills
		for _, item := range export.SkFertigkeiten {
			tx.Save(&item)
		}
		for _, item := range export.SkWaffenfertigkeiten {
			tx.Save(&item)
		}
		for _, item := range export.SkZauber {
			tx.Save(&item)
		}

		// Import equipment
		for _, item := range export.EqAusruestungen {
			tx.Save(&item)
		}
		for _, item := range export.EqWaffen {
			tx.Save(&item)
		}
		for _, item := range export.EqContainers {
			tx.Save(&item)
		}

		// Import GSMaster data
		for _, item := range export.GsmSkills {
			tx.Save(&item)
		}
		for _, item := range export.GsmWeaponSkills {
			tx.Save(&item)
		}
		for _, item := range export.GsmSpells {
			tx.Save(&item)
		}
		for _, item := range export.GsmEquipment {
			tx.Save(&item)
		}
		for _, item := range export.GsmWeapons {
			tx.Save(&item)
		}
		for _, item := range export.GsmContainers {
			tx.Save(&item)
		}
		for _, item := range export.GsmTransportations {
			tx.Save(&item)
		}
		for _, item := range export.GsmBelieves {
			tx.Save(&item)
		}

		// Import learning data
		for _, item := range export.Sources {
			tx.Save(&item)
		}
		for _, item := range export.CharacterClasses {
			tx.Save(&item)
		}
		for _, item := range export.SkillCategories {
			tx.Save(&item)
		}
		for _, item := range export.SkillDifficulties {
			tx.Save(&item)
		}
		for _, item := range export.SpellSchools {
			tx.Save(&item)
		}
		for _, item := range export.ClassCategoryEPCosts {
			tx.Save(&item)
		}
		for _, item := range export.ClassSpellSchoolEPCosts {
			tx.Save(&item)
		}
		for _, item := range export.SpellLevelLECosts {
			tx.Save(&item)
		}
		for _, item := range export.SkillCategoryDifficulties {
			tx.Save(&item)
		}
		for _, item := range export.SkillImprovementCosts {
			tx.Save(&item)
		}
		for _, item := range export.AuditLogEntries {
			tx.Save(&item)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	recordCount := len(export.Users) + len(export.Characters) +
		len(export.Eigenschaften) + len(export.Lps) + len(export.Aps) +
		len(export.Bs) + len(export.Merkmale) + len(export.Erfahrungsschatze) +
		len(export.Bennies) + len(export.Vermoegen) +
		len(export.SkFertigkeiten) + len(export.SkWaffenfertigkeiten) + len(export.SkZauber) +
		len(export.EqAusruestungen) + len(export.EqWaffen) + len(export.EqContainers) +
		len(export.GsmSkills) + len(export.GsmWeaponSkills) + len(export.GsmSpells) +
		len(export.GsmEquipment) + len(export.GsmWeapons) + len(export.GsmContainers) +
		len(export.GsmTransportations) + len(export.GsmBelieves) +
		len(export.Sources) + len(export.CharacterClasses) + len(export.SkillCategories) +
		len(export.SkillDifficulties) + len(export.SpellSchools) +
		len(export.ClassCategoryEPCosts) + len(export.ClassSpellSchoolEPCosts) +
		len(export.SpellLevelLECosts) + len(export.SkillCategoryDifficulties) +
		len(export.SkillImprovementCosts) + len(export.AuditLogEntries) +
		len(export.CharacterCreationSessions)

	return &ImportResult{
		RecordCount: recordCount,
		Timestamp:   time.Now().Format(time.RFC3339),
	}, nil
}
