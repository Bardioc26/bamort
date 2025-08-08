package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"fmt"
	"os"

	"gorm.io/gorm/clause"
)

// ExportData represents the combined data to be exported
type exportData struct {
	Skills          []models.Skill       `json:"skills"`
	WeaponSkills    []models.WeaponSkill `json:"weapon_skills"`
	Spells          []models.Spell       `json:"spells"`
	Equipments      []models.Equipment   `json:"equipments"`
	Weapons         []models.Weapon      `json:"weapons"`
	Containers      []models.Container   `json:"containers"`
	Transportations []models.Container   `json:"transportations"`
	Believes        []models.Believe     `json:"believes"`
}

func Export(filePath string) error {
	var skills []models.Skill
	var weaponSkills []models.WeaponSkill
	var spells []models.Spell
	var equipments []models.Equipment
	var weapons []models.Weapon
	var containers []models.Container
	var transportations []models.Container
	var believes []models.Believe

	// Fetch all data from the respective tables
	if err := database.DB.Find(&skills).Error; err != nil {
		return fmt.Errorf("failed to retrieve equipment: %w", err)
	}
	if err := database.DB.Find(&weaponSkills).Error; err != nil {
		return fmt.Errorf("failed to retrieve equipment: %w", err)
	}
	if err := database.DB.Find(&spells).Error; err != nil {
		return fmt.Errorf("failed to retrieve equipment: %w", err)
	}
	if err := database.DB.Find(&equipments).Error; err != nil {
		return fmt.Errorf("failed to retrieve equipment: %w", err)
	}
	if err := database.DB.Find(&weapons).Error; err != nil {
		return fmt.Errorf("failed to retrieve equipment: %w", err)
	}
	if err := database.DB.Find(&containers).Error; err != nil {
		return fmt.Errorf("failed to retrieve equipment: %w", err)
	}
	if err := database.DB.Find(&transportations).Error; err != nil {
		return fmt.Errorf("failed to retrieve equipment: %w", err)
	}
	if err := database.DB.Find(&believes).Error; err != nil {
		return fmt.Errorf("failed to retrieve equipment: %w", err)
	}

	// Combine the data into a single structure
	exportData := exportData{
		Skills:          skills,
		WeaponSkills:    weaponSkills,
		Spells:          spells,
		Equipments:      equipments,
		Weapons:         weapons,
		Containers:      containers,
		Transportations: transportations,
		Believes:        believes,
	}

	// Create the JSON file for writing
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	// Write the combined data as JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print the JSON
	if err := encoder.Encode(exportData); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	fmt.Printf("Data exported to %s successfully\n", filePath)
	return nil

}

// ReimportFromJSON reads a JSON file and reimports the data into the database
func Import(filePath string) error {

	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %w", err)
	}
	defer file.Close()

	// Decode the JSON file into the ExportData structure
	var data exportData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON file: %w", err)
	}

	// Save the data back to the database
	if len(data.Skills) > 0 {
		if err := database.DB.Clauses(clause.OnConflict{
			UpdateAll: true, // Update all fields if there's a conflict
		}).Create(&data.Skills).Error; err != nil {
			return fmt.Errorf("failed to reimport skills: %w", err)
		}
	}
	if len(data.WeaponSkills) > 0 {
		if err := database.DB.Clauses(clause.OnConflict{
			UpdateAll: true, // Update all fields if there's a conflict
		}).Create(&data.WeaponSkills).Error; err != nil {
			return fmt.Errorf("failed to reimport WeaponSkills: %w", err)
		}
	}
	if len(data.Spells) > 0 {
		if err := database.DB.Clauses(clause.OnConflict{
			UpdateAll: true, // Update all fields if there's a conflict
		}).Create(&data.Spells).Error; err != nil {
			return fmt.Errorf("failed to reimport Spells: %w", err)
		}
	}
	if len(data.Equipments) > 0 {
		if err := database.DB.Clauses(clause.OnConflict{
			UpdateAll: true, // Update all fields if there's a conflict
		}).Create(&data.Equipments).Error; err != nil {
			return fmt.Errorf("failed to reimport equipment: %w", err)
		}
	}
	if len(data.Weapons) > 0 {
		if err := database.DB.Clauses(clause.OnConflict{
			UpdateAll: true, // Update all fields if there's a conflict
		}).Create(&data.Weapons).Error; err != nil {
			return fmt.Errorf("failed to reimport Weapons: %w", err)
		}
	}
	if len(data.Containers) > 0 {
		if err := database.DB.Clauses(clause.OnConflict{
			UpdateAll: true, // Update all fields if there's a conflict
		}).Create(&data.Containers).Error; err != nil {
			return fmt.Errorf("failed to reimport Containers: %w", err)
		}
	}
	if len(data.Transportations) > 0 {
		if err := database.DB.Clauses(clause.OnConflict{
			UpdateAll: true, // Update all fields if there's a conflict
		}).Create(&data.Transportations).Error; err != nil {
			return fmt.Errorf("failed to reimport Transportations: %w", err)
		}
	}
	if len(data.Believes) > 0 {
		if err := database.DB.Clauses(clause.OnConflict{
			UpdateAll: true, // Update all fields if there's a conflict
		}).Create(&data.Believes).Error; err != nil {
			return fmt.Errorf("failed to reimport Believes: %w", err)
		}
	}

	fmt.Printf("Data imported from %s successfully\n", filePath)
	return nil

}
