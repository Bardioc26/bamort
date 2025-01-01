package gsmaster

import (
	"bamort/database"
	"fmt"

	"gorm.io/gorm"
)

type LookupList struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	System       string `gorm:"index" json:"system"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Quelle       string `json:"quelle"`
}

type Skill struct {
	LookupList
	Initialkeitswert int    `json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
	Improvable       bool   `json:"improvable"`
}

type WaeponSkill struct {
	Skill
}

type Spell struct {
	LookupList
	Bonus        int `json:"bonus"`
	Stufe        int
	AP           int
	Reichweite   int
	Wirkungsziel string
}

type Equipment struct {
	LookupList
	Gewicht float64 `json:"gewicht"`
	Wert    float64 `json:"wert"`
}

type Waepon struct {
	Equipment
	SkillRequired string `json:"skill_required"`
	Damage        string `json:"damage"`
}

type Container struct {
	Equipment
	Tragkraft float64 `json:"tragkraft"`
	Volumen   float64 `json:"volumen"`
}

type Transportation struct {
	Container
}

type Believe struct {
	LookupList
}

func (object *LookupList) Create() error {
	gameSystem := "midgard"
	object.System = gameSystem
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
	err := database.DB.First(&object, "system=? AND name = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *LookupList) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND id = ?", gameSystem, value).Error
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

func (stamm *Skill) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
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
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (object *Skill) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Skill) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *WaeponSkill) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupWaeponSkill: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *WaeponSkill) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (object *WaeponSkill) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *WaeponSkill) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *Spell) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
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
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Spell) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND id = ?", gameSystem, value).Error
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

func (stamm *Equipment) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupEquipment: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *Equipment) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Equipment) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Equipment) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *Waepon) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupEquipment: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *Waepon) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Waepon) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Waepon) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *Container) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupContainer: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *Container) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Container) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Container) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *Transportation) Create() error {
	gameSystem := "midgard"
	stamm.System = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save Lookup: %w", err)
		}
		return nil
	})

	return err
}

func (object *Transportation) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *Transportation) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Transportation) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}
