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
	UserID      uint `gorm:"index" json:"user_id"`
}

type Magisch struct {
	IstMagisch  bool `json:"ist_magisch"`
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
}

type LookupList struct {
	ID           uint   `gorm:"primaryKey" json:"id"` //`gorm:"default:uuid_generate_v3()"` // db func
	GameSystem   string `gorm:"column:game_system;index;default:midgard" json:"game_system"`
	Name         string `gorm:"type:varchar(255);index" json:"name"`
	Beschreibung string `json:"beschreibung"`
	Quelle       string `json:"quelle"`                           // Deprecated: Für Rückwärtskompatibilität
	SourceID     uint   `gorm:"index" json:"source_id,omitempty"` // Verweis auf strukturierte Quelle
	PageNumber   int    `json:"page_number,omitempty"`            // Seitenzahl im Quellenbuch
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
	ID               uint   `gorm:"primaryKey" json:"id"`
	GameSystem       string `gorm:"column:game_system;index;default:midgard" json:"game_system"`
	GameSystemId     uint   `json:"game_system_id,omitempty"`
	Name             string `gorm:"type:varchar(255);index" json:"name"`
	Beschreibung     string `json:"beschreibung"`
	Quelle           string `json:"quelle"`                           // Deprecated: Für Rückwärtskompatibilität
	SourceID         uint   `gorm:"index" json:"source_id,omitempty"` // Verweis auf strukturierte Quelle
	PageNumber       int    `json:"page_number,omitempty"`            // Seitenzahl im Quellenbuch
	Initialwert      int    `gorm:"default:5" json:"initialwert"`
	BasisWert        int    `gorm:"default:0" json:"basiswert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
	Improvable       bool   `json:"improvable"`
	InnateSkill      bool   `json:"innateskill"`
	Category         string `json:"category"`
	Difficulty       string `json:"difficulty"`
}

type WeaponSkill struct {
	Skill
}

type Spell struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	GameSystem       string `gorm:"column:game_system;index;default:midgard" json:"game_system"`
	GameSystemId     uint   `json:"game_system_id,omitempty"`
	Name             string `gorm:"type:varchar(255);index" json:"name"`
	Beschreibung     string `json:"beschreibung"`
	Quelle           string `json:"quelle"`                           // Deprecated: Für Rückwärtskompatibilität
	SourceID         uint   `gorm:"index" json:"source_id,omitempty"` // Verweis auf strukturierte Quelle
	PageNumber       int    `json:"page_number,omitempty"`            // Seitenzahl im Quellenbuch
	Bonus            int    `json:"bonus"`
	Stufe            int    `json:"level"`
	AP               string `gorm:"default:1"  json:"ap"`
	Art              string `gorm:"default:Gestenzauber" json:"art"`
	Zauberdauer      string `gorm:"default:10 sec" json:"zauberdauer"`
	Reichweite       string `json:"reichweite"` // in m
	Wirkungsziel     string `json:"wirkungsziel"`
	Wirkungsbereich  string `json:"wirkungsbereich"`
	Wirkungsdauer    string `json:"wirkungsdauer"`
	Ursprung         string `json:"ursprung"`
	Category         string `gorm:"default:normal" json:"category"` // spell_school
	LearningCategory string `gorm:"type:varchar(25);index" json:"learning_category"`
}

type Equipment struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	GameSystem   string  `gorm:"column:game_system;index;default:midgard" json:"game_system"`
	Name         string  `gorm:"type:varchar(255);index" json:"name"`
	Beschreibung string  `json:"beschreibung"`
	Quelle       string  `json:"quelle"`                           // Deprecated: Für Rückwärtskompatibilität
	SourceID     uint    `gorm:"index" json:"source_id,omitempty"` // Verweis auf strukturierte Quelle
	PageNumber   int     `json:"page_number,omitempty"`            // Seitenzahl im Quellenbuch
	Gewicht      float64 `json:"gewicht"`                          // in kg
	Wert         float64 `json:"wert"`                             // in Gold
	PersonalItem bool    `gorm:"default:false" json:"personal_item"`
}

type Weapon struct {
	Equipment
	SkillRequired string `json:"skill_required"`
	Damage        string `json:"damage"`
	RangeNear     int    `json:"range_near"`   // Range in meters for "Nah" (near)
	RangeMiddle   int    `json:"range_middle"` // Range in meters for "Mittel" (medium)
	RangeFar      int    `json:"range_far"`    // Range in meters for "Fern" (far)
}

type Container struct {
	Equipment
	Tragkraft float64 `json:"tragkraft"` // in kg
	Volumen   float64 `json:"volumen"`   // in Liter
}

type Transportation struct {
	Container
}

type Believe struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	GameSystem   string `gorm:"column:game_system;index;default:midgard" json:"game_system"`
	Name         string `gorm:"type:varchar(255);index" json:"name"`
	Beschreibung string `json:"beschreibung"`
	Quelle       string `json:"quelle"`                           // Deprecated: Für Rückwärtskompatibilität
	SourceID     uint   `gorm:"index" json:"source_id,omitempty"` // Verweis auf strukturierte Quelle
	PageNumber   int    `json:"page_number,omitempty"`            // Seitenzahl im Quellenbuch
	GameSystemId uint   `json:"game_system_id,omitempty"`
}

// MiscLookup represents miscellaneous lookup values like gender, race, origin, etc.
type MiscLookup struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	GameSystemId uint   `json:"game_system_id,omitempty"`
	Key          string `gorm:"column:key;type:varchar(50);index;not null" json:"key"`
	Value        string `gorm:"type:varchar(255);not null" json:"value"`
	SourceID     uint   `json:"source_id,omitempty"`
	PageNumber   int    `json:"page_number,omitempty"`
}

func (object *Skill) TableName() string {
	dbPrefix := "gsm"
	return dbPrefix + "_" + "skills"
}

func (stamm *Skill) Create() error {
	if err := database.DB.AutoMigrate(&Skill{}); err != nil {
		return fmt.Errorf("failed to migrate skills: %w", err)
	}
	stamm.ensureGameSystem()

	// Set default values for boolean fields if not explicitly set
	// Note: We cannot rely on GORM defaults anymore since they interfere with copying operations
	if !stamm.Improvable && !stamm.InnateSkill {
		// If both are false, set Improvable to true as default (most skills are improvable)
		stamm.Improvable = true
	}

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
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gs := GetGameSystem(stamm.GameSystemId, stamm.GameSystem)
	// Order by ID to get the record with the lowest ID when multiple categories exist
	err := database.DB.Where("(game_system=? OR game_system_id=?) AND name = ?", gs.Name, gs.ID, name).Order("id ASC").First(&stamm).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (object *Skill) FirstId(value uint) error {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)
	err := database.DB.First(&object, "(game_system=? OR game_system_id=?) AND id = ?", gs.Name, gs.ID, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Skill) Select(fieldName string, value string) ([]Skill, error) {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)
	var skills []Skill
	err := database.DB.Find(&skills, "(game_system=? OR game_system_id=?) AND "+fieldName+" = ?", gs.Name, gs.ID, value).Error
	if err != nil {
		return nil, err
	}
	return skills, nil
}

func SelectSkills(opts ...string) ([]Skill, error) {

	fieldName := ""
	value := ""
	gs := GetGameSystem(0, "midgard")

	if len(opts) > 1 {
		fieldName = opts[0]
		value = opts[1]
	}

	var skills []Skill
	if fieldName == "" {
		err := database.DB.Find(&skills, "(game_system=? OR game_system_id=?)", gs.Name, gs.ID).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := database.DB.Find(&skills, "(game_system=? OR game_system_id=?) AND "+fieldName+" = ?", gs.Name, gs.ID, value).Error
		if err != nil {
			return nil, err
		}
	}
	return skills, nil
}

func (object *Skill) Save() error {
	object.ensureGameSystem()
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
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)

	result := database.DB.Model(&SkillCategory{}).
		Where("game_system = ?", gs.Name).
		Pluck("name", &categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (object *WeaponSkill) TableName() string {
	dbPrefix := "gsm"
	return dbPrefix + "_" + "weaponskills"
}

func (stamm *WeaponSkill) Create() error {
	if err := database.DB.AutoMigrate(&Skill{}, &WeaponSkill{}); err != nil {
		return fmt.Errorf("failed to migrate weapon skills: %w", err)
	}
	stamm.ensureGameSystem()

	// Set default values for boolean fields if not explicitly set
	if !stamm.Improvable && !stamm.InnateSkill {
		stamm.Improvable = true
	}

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
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gs := GetGameSystem(stamm.GameSystemId, stamm.GameSystem)
	// Order by ID to get the record with the lowest ID when multiple categories exist
	err := database.DB.Where("(game_system=? OR game_system_id=?) AND name = ?", gs.Name, gs.ID, name).Order("id ASC").First(&stamm).Error
	if err != nil {
		// Fertigkeit found
		return err
	}
	return nil
}

func (object *WeaponSkill) FirstId(value uint) error {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)
	err := database.DB.First(&object, "(game_system=? OR game_system_id=?) AND id = ?", gs.Name, gs.ID, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *WeaponSkill) Save() error {
	object.ensureGameSystem()
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Spell) TableName() string {
	dbPrefix := "gsm"
	return dbPrefix + "_" + "spells"
}

func (stamm *Spell) Create() error {
	stamm.ensureGameSystem()
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
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gs := GetGameSystem(stamm.GameSystemId, stamm.GameSystem)
	err := database.DB.First(&stamm, "(game_system=? OR game_system_id=?) AND name = ?", gs.Name, gs.ID, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Spell) FirstId(value uint) error {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)
	err := database.DB.First(&object, "(game_system=? OR game_system_id=?) AND id = ?", gs.Name, gs.ID, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Spell) Save() error {
	object.ensureGameSystem()
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

// SelectSpells gibt alle Zauber zurück, optional gefiltert nach einem Feld
func SelectSpells(opts ...string) ([]Spell, error) {
	fieldName := ""
	value := ""
	gameSystem := "midgard"

	if len(opts) > 1 {
		fieldName = opts[0]
		value = opts[1]
	}

	var spells []Spell
	if fieldName == "" {
		err := database.DB.Find(&spells, "game_system=?", gameSystem).Error
		if err != nil {
			return nil, err
		}
	} else {
		err := database.DB.Find(&spells, "game_system=? AND "+fieldName+" = ?", gameSystem, value).Error
		if err != nil {
			return nil, err
		}
	}
	return spells, nil
}

func (object *Spell) GetSpellCategories() ([]string, error) {
	var categories []string
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)

	result := database.DB.Model(&SpellSchool{}).
		Where("game_system = ? OR game_system_id = ?", gs.Name, gs.ID).
		Pluck("name", &categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (object *Spell) ensureGameSystem() {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)

	if object.GameSystemId == 0 {
		object.GameSystemId = gs.ID
	}

	if object.GameSystem == "" {
		object.GameSystem = gs.Name
	}
}

func (object *Skill) ensureGameSystem() {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)

	if object.GameSystemId == 0 {
		object.GameSystemId = gs.ID
	}

	if object.GameSystem == "" {
		object.GameSystem = gs.Name
	}
}

func (object *Skill) BeforeCreate(tx *gorm.DB) error {
	object.ensureGameSystem()
	return nil
}

func (object *Skill) BeforeSave(tx *gorm.DB) error {
	object.ensureGameSystem()
	return nil
}

func (object *Spell) BeforeCreate(tx *gorm.DB) error {
	object.ensureGameSystem()
	return nil
}

func (object *Spell) BeforeSave(tx *gorm.DB) error {
	object.ensureGameSystem()
	return nil
}

func (object *Equipment) TableName() string {
	dbPrefix := "gsm"
	return dbPrefix + "_" + "equipments"
}

func (stamm *Equipment) Create() error {
	gameSystem := "midgard"
	stamm.GameSystem = gameSystem
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
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "game_system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Equipment) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND id = ?", gameSystem, value).Error
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

func (object *Weapon) TableName() string {
	dbPrefix := "gsm"
	return dbPrefix + "_" + "weapons"
}

// IsRanged returns true if the weapon has at least one range value > 0
func (w *Weapon) IsRanged() bool {
	return w.RangeNear > 0 || w.RangeMiddle > 0 || w.RangeFar > 0
}

func (stamm *Weapon) Create() error {
	gameSystem := "midgard"
	stamm.GameSystem = gameSystem
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&stamm).Error; err != nil {
			return fmt.Errorf("failed to save LookupEquipment: %w", err)
		}
		return nil
	})

	return err
}

func (stamm *Weapon) First(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "game_system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Weapon) FirstId(id uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND id = ?", gameSystem, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (object *Weapon) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Container) TableName() string {
	dbPrefix := "gsm"
	return dbPrefix + "_" + "containers"
}

func (stamm *Container) Create() error {
	gameSystem := "midgard"
	stamm.GameSystem = gameSystem
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
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "game_system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Container) FirstId(value uint) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "game_system=? AND id = ?", gameSystem, value).Error
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

func (object *Transportation) TableName() string {
	dbPrefix := "gsm"
	return dbPrefix + "_" + "transportations"
}

func (stamm *Transportation) Create() error {
	gameSystem := "midgard"
	stamm.GameSystem = gameSystem
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
	err := database.DB.First(&object, "game_system=? AND id = ?", gameSystem, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (stamm *Transportation) First(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gameSystem := "midgard"
	err := database.DB.First(&stamm, "game_system=? AND name = ?", gameSystem, name).Error
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

func (object *Believe) TableName() string {
	dbPrefix := "gsm"
	return dbPrefix + "_" + "believes"
}

func (object *Believe) ensureGameSystem() {
	gs := GetGameSystem(object.GameSystemId, "")
	object.GameSystemId = gs.ID
	object.GameSystem = gs.Name
}

func (object *Believe) BeforeCreate(tx *gorm.DB) error {
	object.ensureGameSystem()
	return nil
}

func (object *Believe) BeforeSave(tx *gorm.DB) error {
	object.ensureGameSystem()
	return nil
}

func (object *Believe) Create() error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&object).Error; err != nil {
			return fmt.Errorf("failed to save Lookup: %w", err)
		}
		return nil
	})

	return err
}

func (object *Believe) Save() error {

	err := database.DB.Save(&object).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Believe) FirstId(value uint) error {
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)
	object.GameSystemId = gs.ID
	object.GameSystem = gs.Name

	err := database.DB.First(&object, "(game_system=? OR game_system_id=?) AND id = ?", object.GameSystem, object.GameSystemId, value).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Believe) First(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	gs := GetGameSystem(object.GameSystemId, object.GameSystem)
	object.GameSystemId = gs.ID
	object.GameSystem = gs.Name

	err := database.DB.First(&object, "(game_system=? OR game_system_id=?) AND name = ?", object.GameSystem, object.GameSystemId, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

// GetBelievesByActiveSources gibt Glaubensrichtungen nach aktiven Quellen zurück
func GetBelievesByActiveSources(gameSystem string) ([]Believe, error) {
	var believes []Believe
	gs := GetGameSystem(0, gameSystem)
	if gs.ID == 0 {
		// return empty slice if no valid game system found
		return believes, fmt.Errorf("No GameSystem ID or Name found for %s", gameSystem)
	}
	err := database.DB.
		Joins("LEFT JOIN gsm_lit_sources ON gsm_believes.source_id = gsm_lit_sources.id").
		Where("(gsm_believes.game_system = ? or gsm_believes.game_system_id=?) AND (gsm_lit_sources.is_active = ? OR gsm_believes.source_id IS NULL)", gs.Name, gs.ID, true).
		Order("gsm_believes.name ASC").
		Find(&believes).Error
	return believes, err
}

func (object *MiscLookup) ensureGameSystem() {
	gs := GetGameSystem(object.GameSystemId, "")
	object.GameSystemId = gs.ID
}

func (object *MiscLookup) BeforeCreate(tx *gorm.DB) error {
	object.ensureGameSystem()
	return nil
}

func (object *MiscLookup) BeforeSave(tx *gorm.DB) error {
	object.ensureGameSystem()
	return nil
}

// TableName specifies the table name for MiscLookup
func (MiscLookup) TableName() string {
	return "gsm_misc"
}

func GetGameSystem(id uint, name string) *GameSystem {
	gs := &GameSystem{}
	if id == 0 && name == "" {
		gs.GetDefault()
		return gs
	}
	if id == 0 && name != "" {
		gs.FirstByCode(name)
		if gs.ID == 0 {
			gs.FirstByName(name)
		}
		if gs.ID == 0 {
			gs.GetDefault()
			gs.Name = name
		}
		return gs
	}
	gs.FirstByID(uint(id))
	if gs.ID == 0 && name != "" {
		gs.Name = name
	}
	return gs
}
