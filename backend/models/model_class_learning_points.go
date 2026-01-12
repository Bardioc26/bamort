package models

import (
	"bamort/database"

	"gorm.io/gorm"
)

// ClassLearningPoints stores the learning points distribution for a character class
type ClassLearningPoints struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CharacterClassID uint           `gorm:"uniqueIndex:idx_class_category;not null" json:"character_class_id"`
	SkillCategoryID  uint           `gorm:"uniqueIndex:idx_class_category;not null" json:"skill_category_id"`
	Points           int            `gorm:"not null;default:0" json:"points"`
	CharacterClass   CharacterClass `gorm:"foreignKey:CharacterClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"character_class"`
	SkillCategory    SkillCategory  `gorm:"foreignKey:SkillCategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill_category"`
}

// ClassSpellPoints stores the spell learning points for a character class
type ClassSpellPoints struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CharacterClassID uint           `gorm:"uniqueIndex;not null" json:"character_class_id"`
	SpellPoints      int            `gorm:"not null;default:0" json:"spell_points"` // Zauberlerneinheiten
	CharacterClass   CharacterClass `gorm:"foreignKey:CharacterClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"character_class"`
}

// ClassTypicalSkill stores typical skills for a character class with bonuses
type ClassTypicalSkill struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CharacterClassID uint           `gorm:"not null;index" json:"character_class_id"`
	SkillID          uint           `gorm:"not null;index" json:"skill_id"`
	Bonus            int            `gorm:"not null;default:0" json:"bonus"`
	Attribute        string         `gorm:"size:10" json:"attribute,omitempty"` // z.B. "Gs", "In", "St"
	Notes            string         `json:"notes,omitempty"`                    // Zusätzliche Notizen
	CharacterClass   CharacterClass `gorm:"foreignKey:CharacterClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"character_class"`
	Skill            Skill          `gorm:"foreignKey:SkillID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"skill"`
}

// ClassTypicalSpell stores typical spells for a character class
type ClassTypicalSpell struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	CharacterClassID uint           `gorm:"not null;index" json:"character_class_id"`
	SpellID          uint           `gorm:"not null;index" json:"spell_id"`
	Notes            string         `json:"notes,omitempty"` // z.B. "beliebig außer Dweomer"
	CharacterClass   CharacterClass `gorm:"foreignKey:CharacterClassID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"character_class"`
	Spell            Spell          `gorm:"foreignKey:SpellID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"spell"`
}

// TableName specifies the table name for ClassLearningPoints
func (ClassLearningPoints) TableName() string {
	return "class_learning_points"
}

// TableName specifies the table name for ClassSpellPoints
func (ClassSpellPoints) TableName() string {
	return "class_spell_points"
}

// TableName specifies the table name for ClassTypicalSkill
func (ClassTypicalSkill) TableName() string {
	return "class_typical_skills"
}

// TableName specifies the table name for ClassTypicalSpell
func (ClassTypicalSpell) TableName() string {
	return "class_typical_spells"
}

// GetLearningPointsForClass retrieves all learning points for a character class
func GetLearningPointsForClass(classID uint) ([]ClassLearningPoints, error) {
	var points []ClassLearningPoints
	err := database.DB.Preload("SkillCategory").Where("character_class_id = ?", classID).Find(&points).Error
	return points, err
}

// GetSpellPointsForClass retrieves spell points for a character class
func GetSpellPointsForClass(classID uint) (*ClassSpellPoints, error) {
	var points ClassSpellPoints
	err := database.DB.Where("character_class_id = ?", classID).First(&points).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &ClassSpellPoints{SpellPoints: 0}, nil
		}
		return nil, err
	}
	return &points, nil
}

// GetTypicalSkillsForClass retrieves typical skills for a character class
func GetTypicalSkillsForClass(classID uint) ([]ClassTypicalSkill, error) {
	var skills []ClassTypicalSkill
	err := database.DB.Preload("Skill").Where("character_class_id = ?", classID).Find(&skills).Error
	return skills, err
}

// GetTypicalSpellsForClass retrieves typical spells for a character class
func GetTypicalSpellsForClass(classID uint) ([]ClassTypicalSpell, error) {
	var spells []ClassTypicalSpell
	err := database.DB.Preload("Spell").Where("character_class_id = ?", classID).Find(&spells).Error
	return spells, err
}
