package models

import (
	"bamort/database"
	"bamort/user"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// CharacterCreationSession speichert den Fortschritt der Charakter-Erstellung
type CharacterCreationSession struct {
	ID            string                  `json:"id" gorm:"primaryKey"`
	UserID        uint                    `json:"user_id" gorm:"index;not null"`
	User          user.User               `json:"user" gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name          string                  `json:"name"`
	Geschlecht    string                  `json:"geschlecht"`
	Rasse         string                  `json:"rasse"`
	Typ           string                  `json:"typ"`
	Herkunft      string                  `json:"herkunft"`
	Stand         string                  `json:"stand"`
	Glaube        string                  `json:"glaube"`
	Attributes    AttributesData          `json:"attributes" gorm:"type:text;serializer:json"`
	DerivedValues DerivedValuesData       `json:"derived_values" gorm:"type:text;serializer:json"`
	Skills        CharacterCreationSkills `json:"skills" gorm:"type:text;serializer:json"`
	Spells        CharacterCreationSpells `json:"spells" gorm:"type:text;serializer:json"`
	SkillPoints   SkillPointsData         `json:"skill_points" gorm:"type:text;serializer:json"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	ExpiresAt     time.Time               `json:"expires_at"`
	CurrentStep   int                     `json:"current_step"` // 1=Basic, 2=Attributes, 3=Derived, 4=Skills
}

// AttributesData speichert die Grundwerte
type AttributesData struct {
	ST int `json:"st"` // Stärke
	GS int `json:"gs"` // Geschicklichkeit
	GW int `json:"gw"` // Gewandtheit
	KO int `json:"ko"` // Konstitution
	IN int `json:"in"` // Intelligenz
	ZT int `json:"zt"` // Zaubertalent
	AU int `json:"au"` // Ausstrahlung
}

// DerivedValuesData speichert die abgeleiteten Werte
type DerivedValuesData struct {
	PA                    int `json:"pa"`                      // Persönliche Ausstrahlung
	WK                    int `json:"wk"`                      // Willenskraft
	LPMax                 int `json:"lp_max"`                  // Lebenspunkte Maximum
	APMax                 int `json:"ap_max"`                  // Abenteuerpunkte Maximum
	BMax                  int `json:"b_max"`                   // Belastung Maximum
	ResistenzKoerper      int `json:"resistenz_koerper"`       // Resistenz Körper
	ResistenzGeist        int `json:"resistenz_geist"`         // Resistenz Geist
	ResistenzBonusKoerper int `json:"resistenz_bonus_koerper"` // Resistenz Bonus Körper
	ResistenzBonusGeist   int `json:"resistenz_bonus_geist"`   // Resistenz Bonus Geist
	Abwehr                int `json:"abwehr"`                  // Abwehr
	AbwehrBonus           int `json:"abwehr_bonus"`            // Abwehr Bonus
	AusdauerBonus         int `json:"ausdauer_bonus"`          // Ausdauer Bonus
	AngriffsBonus         int `json:"angriffs_bonus"`          // Angriffs Bonus
	Zaubern               int `json:"zaubern"`                 // Zaubern
	ZauberBonus           int `json:"zauber_bonus"`            // Zauber Bonus
	Raufen                int `json:"raufen"`                  // Raufen
	SchadensBonus         int `json:"schadens_bonus"`          // Schadens Bonus
	SG                    int `json:"sg"`                      // Schicksalsgunst
	GG                    int `json:"gg"`                      // Göttliche Gnade
	GP                    int `json:"gp"`                      // Glückspunkte
}

// SkillPointsData speichert die verbleibenden Lernpunkte pro Kategorie
type SkillPointsData map[string]int

// CharacterCreationSkill repräsentiert eine ausgewählte Fertigkeit
type CharacterCreationSkill struct {
	Name     string `json:"name"`
	Level    int    `json:"level"`
	Category string `json:"category"`
	Cost     int    `json:"cost"`
}

// CharacterCreationSpell repräsentiert einen ausgewählten Zauber
type CharacterCreationSpell struct {
	Name string `json:"name"`
	Cost int    `json:"cost"`
}

// Slice types for GORM JSON handling
type CharacterCreationSkills []CharacterCreationSkill
type CharacterCreationSpells []CharacterCreationSpell

// ClassCategoryLearningPoints stores the learning points distribution for a character class
type ClassCategoryLearningPoints struct {
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
func (ClassCategoryLearningPoints) TableName() string {
	//class_learning_points
	return "gsm_cc_class_category_points"
}

// TableName specifies the table name for ClassSpellPoints
func (ClassSpellPoints) TableName() string {
	//class_spell_points
	return "gsm_cc_class_spell_points"
}

// TableName specifies the table name for ClassTypicalSkill
func (ClassTypicalSkill) TableName() string {
	//class_typical_skills
	return "gsm_cc_class_typical_skills"
}

// TableName specifies the table name for ClassTypicalSpell
func (ClassTypicalSpell) TableName() string {
	//class_typical_spells
	return "gsm_cc_class_typical_spells"
}

// JSON scanning methods for Skills slice
func (s *CharacterCreationSkills) Scan(value interface{}) error {
	if value == nil {
		*s = make([]CharacterCreationSkill, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		*s = make([]CharacterCreationSkill, 0)
		return nil
	}

	return json.Unmarshal(bytes, s)
}

func (s CharacterCreationSkills) Value() (interface{}, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// JSON scanning methods for Spells slice
func (sp *CharacterCreationSpells) Scan(value interface{}) error {
	if value == nil {
		*sp = make([]CharacterCreationSpell, 0)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		*sp = make([]CharacterCreationSpell, 0)
		return nil
	}

	return json.Unmarshal(bytes, sp)
}

func (sp CharacterCreationSpells) Value() (interface{}, error) {
	if len(sp) == 0 {
		return "[]", nil
	}
	return json.Marshal(sp)
}

func (object *CharacterCreationSession) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "char_creation_session"
}

// Scanning methods for JSON fields in GORM
func (a *AttributesData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, a)
}

func (a AttributesData) Value() (interface{}, error) {
	return json.Marshal(a)
}

func (d *DerivedValuesData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, d)
}

func (d DerivedValuesData) Value() (interface{}, error) {
	return json.Marshal(d)
}

func (s *SkillPointsData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, s)
}

func (s SkillPointsData) Value() (interface{}, error) {
	return json.Marshal(s)
}

// Cleanup expired sessions
func CleanupExpiredSessions(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).Delete(&CharacterCreationSession{}).Error
}

// Find sessions by user ID
func GetUserSessions(db *gorm.DB, userID uint) ([]CharacterCreationSession, error) {
	var sessions []CharacterCreationSession
	err := db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error
	return sessions, err
}

// GetLearningPointsForClass retrieves all learning points for a character class
func GetLearningPointsForClass(classID uint) ([]ClassCategoryLearningPoints, error) {
	var points []ClassCategoryLearningPoints
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
