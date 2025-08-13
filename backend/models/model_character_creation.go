package models

import (
	"bamort/user"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// CharacterCreationSession speichert den Fortschritt der Charakter-Erstellung
type CharacterCreationSession struct {
	ID            string                  `json:"id" gorm:"primaryKey"`
	UserID        uint                    `json:"user_id" gorm:"index;not null"`
	User          user.User               `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
