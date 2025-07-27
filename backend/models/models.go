package models

import (
	"bamort/database"
	"fmt"

	"gorm.io/gorm"
)

type BamortBase struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type BamortCharTrait struct {
	BamortBase
	CharacterID uint `gorm:"index" json:"character_id"`
}

type Magisch struct {
	IstMagisch  bool `json:"ist_magisch"`
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
}

type LookupList struct {
	ID           uint   `gorm:"primaryKey" json:"id"` //`gorm:"default:uuid_generate_v3()"` // db func
	GameSystem   string `gorm:"column:game_system;index;default:midgard" json:"game_system"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Quelle       string `json:"quelle"`
}

type LearnCost struct {
	Stufe int `json:"stufe"`
	LE    int `json:"le"`
	TE    int `json:"te"`
	Ep    int `json:"ep"`
	Money int `json:"money"`
	PP    int `json:"pp"`
}

type Skill struct {
	LookupList
	Initialwert      int    `gorm:"default:5" json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
	Improvable       bool   `gorm:"default:true" json:"improvable"`
	InnateSkill      bool   `gorm:"default:false" json:"innateskill"`
	Category         string `json:"category"`
	Difficulty       string `json:"difficulty"`
}

type WeaponSkill struct {
	Skill
}

type Spell struct {
	LookupList
	Bonus           int    `json:"bonus"`
	Stufe           int    `json:"level"`
	AP              string `gorm:"default:1"  json:"ap"`
	Art             string `gorm:"default:Gestenzauber" json:"art"`
	Zauberdauer     string `gorm:"default:10 sec" json:"zauberdauer"`
	Reichweite      string `json:"reichweite"` // in m
	Wirkungsziel    string `json:"wirkungsziel"`
	Wirkungsbereich string `json:"wirkungsbereich"`
	Wirkungsdauer   string `json:"wirkungsdauer"`
	Ursprung        string `json:"ursprung"`
	Category        string `gorm:"default:normal" json:"category"`
}

func (object *LookupList) Create() error {
	gameSystem := "midgard"
	object.GameSystem = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&object).Error; err != nil {
			return fmt.Errorf("failed to save Lookup: %w", err)
		}
		return nil
	})

	return err
}

func (object *LookupList) First(value string) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND name!='Placeholder' AND name = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *LookupList) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND name!='Placeholder' AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *LookupList) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Skill) TableName() string {
	var dbPrefix = "gsm"
	return dbPrefix + "_" + "skills"
}

func (stamm *Skill) Create() error {
	gameSystem := "midgard"
	stamm.GameSystem = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupSkill: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *Skill) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "game_system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (object *Skill) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Skill) Select(fieldName string, value string) ([]Skill, error) {
	gameSystem := "midgard"
	var skills []Skill
	err := database.DB.Find(&skills, "game_system=? AND name != 'Placeholder' AND "+fieldName+" = ?", gameSystem, value).Error
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func SelectSkills(fieldName string, value string) ([]Skill, error) {
	gameSystem := "midgard"
	var skills []Skill
	if fieldName == "" {
		err := database.DB.Find(&skills, "game_system=? AND name != 'Placeholder'", gameSystem).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := database.DB.Find(&skills, "game_system=? AND name != 'Placeholder' AND "+fieldName+" = ?", gameSystem, value).Error
		if err != nil {
			return nil, err
		}
	}
	return skills, nil
}

func (object *Skill) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Skill) Delete() error {
	result := database.DB.Delete(&object, object.ID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no record found with ID %v", object.ID)
	}
	return nil
}

func (object *Skill) GetSkillCategories() ([]string, error) {
	var categories []string
	gameSystem := "midgard"

	result := database.DB.Model(&Skill{}).
		Where("game_system = ? and category is not null", gameSystem).
		Distinct().
		Pluck("category", &categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (object *WeaponSkill) TableName() string {
	var dbPrefix = "gsm"
	return dbPrefix + "_" + "weaponskills"
}

func (stamm *WeaponSkill) Create() error {
	gameSystem := "midgard"
	stamm.GameSystem = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupWeaponSkill: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *WeaponSkill) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "game_system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (object *WeaponSkill) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *WeaponSkill) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Spell) TableName() string {
	var dbPrefix = "gsm"
	return dbPrefix + "_" + "spells"
}

func (stamm *Spell) Create() error {
	gameSystem := "midgard"
	stamm.GameSystem = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupSpell: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *Spell) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "game_system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Spell) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Spell) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Spell) GetSpellCategories() ([]string, error) {
	var categories []string
	gameSystem := "midgard"

	result := database.DB.Model(&Spell{}).
		Where("game_system=? and category is not null", gameSystem).
		Distinct().
		Pluck("category", &categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}
