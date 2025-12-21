package models

import "strings"

type SkFertigkeit struct {
	BamortCharTrait
	Beschreibung    string    `json:"beschreibung"`
	Fertigkeitswert int       `json:"fertigkeitswert"`
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
	if object.Category != "" {
		return object.Category
	}
	var gsmsk Skill
	gsmsk.First(object.Name)
	if gsmsk.ID == 0 {
		return "Unkategorisiert"
	}
	// Trim whitespace from category to handle inconsistent data
	category := strings.TrimSpace(gsmsk.Category)
	object.Category = category
	return object.Category
}

func (object *SkWaffenfertigkeit) GetCategory() string {
	if object.Category != "" {
		return object.Category
	}
	// For weapon skills, we need to look in the WeaponSkill table
	var weaponSkill WeaponSkill
	err := weaponSkill.First(object.Name)
	if err != nil || weaponSkill.ID == 0 {
		return "Unkategorisiert"
	}
	// Trim whitespace from category to handle inconsistent data
	category := strings.TrimSpace(weaponSkill.Category)
	object.Category = category
	return object.Category
}
