package models

import (
	"bamort/user"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// CharacterCreationSession speichert den Fortschritt der Charakter-Erstellung
type CharacterCreationSession struct {
	ID            string                   `json:"id" gorm:"primaryKey"`
	UserID        uint                     `json:"user_id" gorm:"index;not null"`
	User          user.User                `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name          string                   `json:"name"`
	Rasse         string                   `json:"rasse"`
	Typ           string                   `json:"typ"`
	Herkunft      string                   `json:"herkunft"`
	Glaube        string                   `json:"glaube"`
	Attributes    AttributesData           `json:"attributes" gorm:"type:json"`
	DerivedValues DerivedValuesData        `json:"derived_values" gorm:"type:json"`
	Skills        []CharacterCreationSkill `json:"skills" gorm:"type:json"`
	Spells        []CharacterCreationSpell `json:"spells" gorm:"type:json"`
	SkillPoints   SkillPointsData          `json:"skill_points" gorm:"type:json"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
	ExpiresAt     time.Time                `json:"expires_at"`
	CurrentStep   int                      `json:"current_step"` // 1=Basic, 2=Attributes, 3=Derived, 4=Skills
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
	PA int `json:"pa"` // Psi-Kraft
	WK int `json:"wk"` // Willenskraft
}

// DerivedValuesData speichert die abgeleiteten Werte
type DerivedValuesData struct {
	LPMax int `json:"lp_max"` // Lebenspunkte Maximum
	APMax int `json:"ap_max"` // Abenteuerpunkte Maximum
	BMax  int `json:"b_max"`  // Belastung Maximum
	SG    int `json:"sg"`     // Schicksalsgunst
	GG    int `json:"gg"`     // Göttliche Gnade
	GP    int `json:"gp"`     // Glückspunkte
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
