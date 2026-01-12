package models

import (
	"bamort/database"
	"strings"
)

type SkFertigkeit struct {
	BamortCharTrait
	Beschreibung    string    `json:"beschreibung"`
	Fertigkeitswert int       `json:"fertigkeitswert"`
	BasisWert       int       `json:"basiswert"`
	Bonus           int       `json:"bonus,omitempty"`
	Pp              int       `json:"pp,omitempty"` //Praxispunkte
	Bemerkung       string    `json:"bemerkung"`
	Improvable      bool      `json:"improvable"`
	Category        string    `json:"category"`
	LearningCost    LearnCost `gorm:"-" json:"learncost"`
}

type SkWaffenfertigkeit struct {
	SkFertigkeit
}

type SkAngeboreneFertigkeit struct {
	SkFertigkeit
}

type SkZauber struct {
	BamortCharTrait
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

func (object *SkFertigkeit) TableName() string {
	dbPrefix := "skill"
	return dbPrefix + "_" + "skills"
}
func (object *SkWaffenfertigkeit) TableName() string {
	dbPrefix := "skill"
	return dbPrefix + "_" + "weaponskills"
}
func (object *SkZauber) TableName() string {
	dbPrefix := "skill"
	return dbPrefix + "_" + "spells"
}

func (object *SkFertigkeit) GetSkillByName() *Skill {
	var gsmsk Skill
	gsmsk.First(object.Name)
	if gsmsk.ID == 0 {
		return nil
	}
	return &gsmsk
}

func (object *SkWaffenfertigkeit) GetSkillByName() *Skill {
	// For weapon skills, we need to look in the WeaponSkill table
	var weaponSkill WeaponSkill
	err := weaponSkill.First(object.Name)
	if err != nil || weaponSkill.ID == 0 {
		return nil
	}
	// Return the embedded Skill from WeaponSkill
	return &weaponSkill.Skill
}

func (object *SkFertigkeit) GetCategory() string {
	// Always fetch category from gsmaster, ignoring the category field in skill_skills
	var gsmsk Skill
	gsmsk.First(object.Name)
	if gsmsk.ID == 0 {
		return "Unkategorisiert"
	}

	// Fetch category from learning_skill_category_difficulties table
	// Order by ID to get the lowest ID when multiple categories exist
	var scd SkillCategoryDifficulty
	err := database.DB.Where("skill_id = ?", gsmsk.ID).
		Order("id ASC").
		Preload("SkillCategory").
		First(&scd).Error

	if err != nil {
		// If not found in learning table, fall back to Unkategorisiert
		return "Unkategorisiert"
	}

	// Use the SCategory field which contains the category name
	category := strings.TrimSpace(scd.SCategory)
	if category == "" {
		return "Unkategorisiert"
	}
	return category
}

func (object *SkWaffenfertigkeit) GetCategory() string {
	// Weapon skills don't use the learning_skill_category_difficulties table
	// They have their category directly in gsm_weaponskills
	var weaponSkill WeaponSkill
	err := weaponSkill.First(object.Name)
	if err != nil || weaponSkill.ID == 0 {
		return "Waffenfertigkeiten"
	}
	// Trim whitespace from category to handle inconsistent data
	category := strings.TrimSpace(weaponSkill.Category)
	if category == "" {
		return "Waffenfertigkeiten"
	}
	return category
}
