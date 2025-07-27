package gsmaster

import (
	"time"

	"gorm.io/gorm"
)

// CharacterClass represents a character class with its abbreviation and full name
type CharacterClass struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Code        string `gorm:"index:idx_character_class_code,unique,length:2;size:2;not null" json:"code"`
	Name        string `gorm:"index:idx_character_class_name,unique,length:50;not null" json:"name"`
	Description string `json:"description"`
	GameSystem  string `gorm:"column:game_system;index;size:50;default:midgard" json:"game_system"`
	Quelle      string `json:"quelle"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SkillCategory represents a skill category like "Alltag", "Freiland", etc.
type SkillCategory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"index:idx_skill_category_name,unique,length:50;size:50;not null" json:"name"`
	GameSystem  string `gorm:"column:game_system;index;size:50;default:midgard" json:"game_system"`
	Quelle      string `json:"quelle"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SkillDifficulty represents difficulty levels within a category
type SkillDifficulty struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	Name       string        `gorm:"index:idx_skill_difficulty_name;not null" json:"name"`
	CategoryID uint          `gorm:"not null" json:"category_id"`
	Category   SkillCategory `gorm:"foreignKey:CategoryID" json:"category"`
	LearnCost  int           `gorm:"not null" json:"learn_cost"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ClassCategoryEPCost represents EP costs per training unit for class/category combinations
type ClassCategoryEPCost struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ClassID    uint           `gorm:"not null" json:"class_id"`
	Class      CharacterClass `gorm:"foreignKey:ClassID" json:"class"`
	CategoryID uint           `gorm:"not null" json:"category_id"`
	Category   SkillCategory  `gorm:"foreignKey:CategoryID" json:"category"`
	EPPerTE    int            `gorm:"not null" json:"ep_per_te"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// SkillDefinitionNew represents individual skills and their category/difficulty
type SkillDefinitionNew struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	Name             string          `gorm:"index:idx_skill_definition_name;not null" json:"name"`
	GameSystem       string          `gorm:"column:game_system;index;size:50;default:midgard" json:"game_system"`
	Quelle           string          `json:"quelle"`
	CategoryID       uint            `gorm:"not null" json:"category_id"`
	Category         SkillCategory   `gorm:"foreignKey:CategoryID" json:"category"`
	DifficultyID     uint            `gorm:"not null" json:"difficulty_id"`
	Difficulty       SkillDifficulty `gorm:"foreignKey:DifficultyID" json:"difficulty"`
	Initialwert      int             `gorm:"default:5" json:"initialwert"`
	Bonuseigenschaft string          `json:"bonuseigenschaft,omitempty"`
	Improvable       bool            `gorm:"default:true" json:"improvable"`
	InnateSkill      bool            `gorm:"default:false" json:"innateskill"`
	Spezialized      bool            `gorm:"default:false" json:"specialized"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// SkillImprovementCost represents training costs for skill improvements by level
type SkillImprovementCost struct {
	ID           uint            `gorm:"primaryKey" json:"id"`
	DifficultyID uint            `gorm:"not null" json:"difficulty_id"`
	Difficulty   SkillDifficulty `gorm:"foreignKey:DifficultyID" json:"difficulty"`
	CurrentLevel int             `gorm:"not null" json:"current_level"`
	TECost       int             `gorm:"not null" json:"te_cost"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SpellSchool represents spell schools like "Beherrschen", "Bewegen", etc.
type SpellSchool struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"index:idx_spell_school_name,unique,length:50;size:50;not null" json:"name"`
	GameSystem  string `gorm:"column:game_system;index;size:50;default:midgard" json:"game_system"`
	Quelle      string `json:"quelle"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ClassSpellSchoolCost represents EP costs per learning unit for class/spell school combinations
type ClassSpellSchoolCost struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ClassID          uint           `gorm:"not null" json:"class_id"`
	Class            CharacterClass `gorm:"foreignKey:ClassID" json:"class"`
	SchoolID         uint           `gorm:"not null" json:"school_id"`
	School           SpellSchool    `gorm:"foreignKey:SchoolID" json:"school"`
	EPPerLE          int            `gorm:"not null" json:"ep_per_le"`
	IsSpecialization bool           `gorm:"default:false" json:"is_specialization"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// SpellLevelCost represents learning costs for spells by level
type SpellLevelCost struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	SpellLevel int  `gorm:"index:idx_spell_level,unique;not null" json:"spell_level"`
	LERequired int  `gorm:"not null" json:"le_required"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// TableName methods to ensure correct table names
func (CharacterClass) TableName() string {
	return "character_classes"
}

func (SkillCategory) TableName() string {
	return "skill_categories"
}

func (SkillDifficulty) TableName() string {
	return "skill_difficulties"
}

func (ClassCategoryEPCost) TableName() string {
	return "class_category_ep_costs"
}

func (SkillDefinitionNew) TableName() string {
	return "skill_definitions"
}

func (SkillImprovementCost) TableName() string {
	return "skill_improvement_costs"
}

func (SpellSchool) TableName() string {
	return "spell_schools"
}

func (ClassSpellSchoolCost) TableName() string {
	return "class_spell_school_costs"
}

func (SpellLevelCost) TableName() string {
	return "spell_level_costs"
}

// MigrateLearningCostsTables creates all learning costs tables
func MigrateLearningCostsTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&CharacterClass{},
		&SkillCategory{},
		&SkillDifficulty{},
		&ClassCategoryEPCost{},
		&SkillDefinitionNew{},
		&SkillImprovementCost{},
		&SpellSchool{},
		&ClassSpellSchoolCost{},
		&SpellLevelCost{},
	)
}
