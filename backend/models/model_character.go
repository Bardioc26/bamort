package models

import (
	"bamort/database"
	"bamort/user"
	"fmt"

	"gorm.io/gorm"
)

// Au, Gs, Gw ,In, Ko, Pa, St, Wk, Zt
type Eigenschaft struct {
	ID          uint   `gorm:"index" json:"id"`
	CharacterID uint   `gorm:"primaryKey" json:"character_id"`
	UserID      uint   `gorm:"index" json:"user_id"`
	Name        string `gorm:"primaryKey" json:"name"`
	Value       int    `json:"value"`
}
type Lp struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	CharacterID uint `gorm:"index" json:"character_id"`
	Max         int  `json:"max"`
	Value       int  `json:"value"`
}

type Ap struct {
	ID uint `gorm:"primaryKey" json:"id"`

	CharacterID uint `gorm:"index" json:"character_id"`
	Max         int  `json:"max"`
	Value       int  `json:"value"`
}

type B struct {
	ID uint `gorm:"primaryKey" json:"id"`

	CharacterID uint `gorm:"index" json:"character_id"`
	Max         int  `json:"max"`
	Value       int  `json:"value"`
}

/*
	type Gestalt struct {
		models.BamortCharTrait
	}
*/

type Merkmale struct {
	BamortCharTrait
	Augenfarbe string `json:"augenfarbe"`
	Haarfarbe  string `json:"haarfarbe"`
	Sonstige   string `json:"sonstige"`
	Breite     string `json:"breite"`
	Groesse    string `json:"groesse"`
}

type Erfahrungsschatz struct {
	BamortCharTrait
	ES int `json:"es"` // Erfahrungsschatz
	EP int `json:"ep"` // Erfahrungspunkte
}

type Bennies struct {
	BamortCharTrait
	Gg int `json:"gg"` // Göttliche Gnade
	Gp int `json:"gp"` // Glückspunkte
	Sg int `json:"sg"` // Schicksalsgunst
}

type Vermoegen struct {
	BamortCharTrait
	Goldstuecke   int `json:"goldstücke"`   // GS
	Silberstuecke int `json:"silberstücke"` // SS
	Kupferstuecke int `json:"kupferstücke"` // KS
}

type Char struct {
	BamortBase
	UserID      uint      `gorm:"index;not null;default:1" json:"user_id"`
	User        user.User `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Rasse       string    `json:"rasse"`
	Typ         string    `json:"typ"`
	Alter       int       `json:"alter"`
	Anrede      string    `json:"anrede"`
	Grad        int       `json:"grad"`
	Gender      string    `json:"gender"`
	SocialClass string    `json:"social_class"`
	Groesse     int       `json:"groesse"`
	Gewicht     int       `json:"gewicht"`
	Herkunft    string    `json:"origin"`
	Glaube      string    `json:"glaube"`
	Hand        string    `json:"hand"`
	Public      bool      `json:"public"`
	// Static derived values (can increase with grade)
	ResistenzKoerper   int                  `json:"resistenz_koerper"`
	ResistenzGeist     int                  `json:"resistenz_geist"`
	Abwehr             int                  `json:"abwehr"`
	Zaubern            int                  `json:"zaubern"`
	Raufen             int                  `json:"raufen"`
	Lp                 Lp                   `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"lp"`
	Ap                 Ap                   `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"ap"`
	B                  B                    `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"b"`
	Merkmale           Merkmale             `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"merkmale"`
	Eigenschaften      []Eigenschaft        `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"eigenschaften"`
	Fertigkeiten       []SkFertigkeit       `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"fertigkeiten"`
	Waffenfertigkeiten []SkWaffenfertigkeit `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"waffenfertigkeiten"`
	Zauber             []SkZauber           `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"zauber"`
	Spezialisierung    database.StringArray `gorm:"type:TEXT"  json:"spezialisierung"`
	Bennies            Bennies              `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bennies"`
	Vermoegen          Vermoegen            `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"vermoegen"`
	Erfahrungsschatz   Erfahrungsschatz     `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"erfahrungsschatz"`
	Waffen             []EqWaffe            `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"waffen"`
	Behaeltnisse       []EqContainer        `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"behaeltnisse"`
	Transportmittel    []EqContainer        `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"transportmittel"`
	Ausruestung        []EqAusruestung      `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"ausruestung"`
	Image              string               `json:"image,omitempty"`
}
type CharList struct {
	BamortBase
	UserID uint   `json:"user_id"`
	Rasse  string `json:"rasse"`
	Typ    string `json:"typ"`
	Grad   int    `json:"grad"`
	Owner  string `json:"owner"`
	Public bool   `json:"public"`
}

type FeChar struct {
	Char
	CategorizedSkills map[string][]SkFertigkeit `json:"categorizedskills"`
	InnateSkills      []SkFertigkeit            `json:"innateskills"`
}

func (object *Char) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "chars"
}

func (object *Char) First(charName string) error {
	err := database.DB.
		Preload("User").
		Preload("Lp").
		Preload("Ap").
		Preload("B").
		Preload("Merkmale").
		Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Zauber").
		Preload("Bennies").
		Preload("Vermoegen").
		Preload("Erfahrungsschatz").
		Preload("Waffen").
		Preload("Behaeltnisse").
		Preload("Transportmittel").
		Preload("Ausruestung").
		First(&object, " name = ?", charName).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Char) FirstID(charID string) error {
	err := database.DB.
		Preload("User").
		Preload("Lp").
		Preload("Ap").
		Preload("B").
		Preload("Merkmale").
		Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Zauber").
		Preload("Bennies").
		Preload("Vermoegen").
		Preload("Erfahrungsschatz").
		Preload("Waffen").
		Preload("Behaeltnisse").
		Preload("Transportmittel").
		Preload("Ausruestung").
		First(&object, charID).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

// FindByUserID finds all characters belonging to a specific user
func (object *Char) FindByUserID(userID uint) ([]Char, error) {
	var chars []Char
	err := database.DB.
		Preload("User").
		Preload("Lp").
		Preload("Ap").
		Preload("B").
		Preload("Merkmale").
		Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Zauber").
		Preload("Bennies").
		Preload("Vermoegen").
		Preload("Erfahrungsschatz").
		Preload("Waffen").
		Preload("Behaeltnisse").
		Preload("Transportmittel").
		Preload("Ausruestung").
		Where("user_id = ?", userID).
		Find(&chars).Error
	if err != nil {
		return nil, err
	}
	return chars, nil
}

func FindPublicCharList() ([]CharList, error) {
	var chars []CharList
	err := database.DB.Table("char_chars").
		Select("char_chars.id, char_chars.name, char_chars.user_id, char_chars.rasse, char_chars.typ, char_chars.grad, char_chars.public, users.username as owner").
		Joins("LEFT JOIN users ON char_chars.user_id = users.user_id").
		Where("char_chars.public = ?", true).
		Find(&chars).Error
	if err != nil {
		return nil, err
	}
	return chars, nil
}

// FindCharListByUserID finds all characters belonging to a specific user for listing (minimal data)
func FindCharListByUserID(userID uint) ([]CharList, error) {
	var chars []CharList
	err := database.DB.Table("char_chars").
		Select("char_chars.id, char_chars.name, char_chars.user_id, char_chars.rasse, char_chars.typ, char_chars.grad, char_chars.public, users.username as owner").
		Joins("LEFT JOIN users ON char_chars.user_id = users.user_id").
		Where("char_chars.user_id = ?", userID).
		Find(&chars).Error
	if err != nil {
		return nil, err
	}
	return chars, nil
}

func (object *Char) Create() error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character record
		if err := tx.Create(&object).Error; err != nil {
			return fmt.Errorf("failed to save Lookup: %w", err)
		}
		return nil
	})

	return err
}

func (object *Char) Delete() error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// delete the main character record
		//should cascade for all elements
		if err := tx.Delete(&object).Error; err != nil {
			return fmt.Errorf("failed to delete char: %w", err)
		}
		return nil
	})

	return err
}
func (object *Eigenschaft) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "eigenschaften"
}
func (object *Lp) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "health"
}
func (object *Ap) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "endurances"
}
func (object *B) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "motionranges"
}
func (object *Merkmale) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "characteristics"
}
func (object *Erfahrungsschatz) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "experiances"
}
func (object *Bennies) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "bennies"
}
func (object *Vermoegen) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "wealth"
}

// DerivedBonuses contains all calculated bonus values
type DerivedBonuses struct {
	AusdauerBonus         int
	SchadensBonus         int
	AngriffsBonus         int
	AbwehrBonus           int
	ZauberBonus           int
	ResistenzBonusKoerper int
	ResistenzBonusGeist   int
}

// GetAttributeValue returns the value of an attribute by name
func (char *Char) GetAttributeValue(name string) int {
	for _, attr := range char.Eigenschaften {
		if attr.Name == name {
			return attr.Value
		}
	}
	return 0
}

// CalculateBonuses calculates all derived bonuses from attributes
func (char *Char) CalculateBonuses() DerivedBonuses {
	st := char.GetAttributeValue("St")
	gs := char.GetAttributeValue("Gs")
	gw := char.GetAttributeValue("Gw")
	ko := char.GetAttributeValue("Ko")
	in := char.GetAttributeValue("In")
	zt := char.GetAttributeValue("Zt")

	bonuses := DerivedBonuses{
		// Ausdauer Bonus: Ko/10 + St/20
		AusdauerBonus: (ko / 10) + (st / 20),

		// Schadens Bonus: St/20 + Gs/30 - 3
		SchadensBonus: (st / 20) + (gs / 30) - 3,

		// Angriffs Bonus basierend auf GS
		AngriffsBonus: calculateAttributeBonus(gs),

		// Abwehr Bonus basierend auf GW
		AbwehrBonus: calculateAttributeBonus(gw),

		// Zauber Bonus basierend auf Zt
		ZauberBonus: calculateAttributeBonus(zt),
	}

	// Resistenz Bonus Körper
	bonuses.ResistenzBonusKoerper = calculateResistenzBonusKoerper(ko, char.Rasse, char.Typ)

	// Resistenz Bonus Geist
	bonuses.ResistenzBonusGeist = calculateResistenzBonusGeist(in, char.Rasse, char.Typ)

	return bonuses
}

// Helper functions for bonus calculation

func calculateAttributeBonus(value int) int {
	if value <= 5 {
		return -5
	} else if value <= 20 {
		return -1
	} else if value <= 40 {
		return 0
	} else if value <= 60 {
		return 1
	} else if value <= 80 {
		return 2
	} else if value <= 95 {
		return 3
	} else {
		return 4
	}
}

func calculateResistenzBonusKoerper(ko int, rasse string, typ string) int {
	bonus := 0

	if rasse == "Mensch" || rasse == "Menschen" {
		bonus = calculateAttributeBonus(ko)
	} else {
		switch rasse {
		case "Elfen":
			bonus = 2
		case "Gnome", "Halblinge":
			bonus = 4
		case "Zwerge":
			bonus = 3
		}
	}

	// Klassenmodifikator
	if isKaempfer(typ) {
		bonus += 1
	} else if isZauberer(typ) {
		bonus += 2
	}

	return bonus
}

func calculateResistenzBonusGeist(in int, rasse string, typ string) int {
	bonus := 0

	if rasse == "Mensch" || rasse == "Menschen" {
		bonus = calculateAttributeBonus(in)
	} else {
		switch rasse {
		case "Elfen":
			bonus = 2
		case "Gnome", "Halblinge":
			bonus = 4
		case "Zwerge":
			bonus = 3
		}
	}

	// Klassenmodifikator (nur Zauberer bekommen Geist-Bonus)
	if isZauberer(typ) {
		bonus += 2
	}

	return bonus
}

func isKaempfer(typ string) bool {
	kaempferClasses := []string{
		"Krieger", "Barbar", "Spitzbube", "Assassine",
		"Streiter", "Waldläufer", "Krieger Magier",
	}
	for _, class := range kaempferClasses {
		if typ == class {
			return true
		}
	}
	return false
}

func isZauberer(typ string) bool {
	zaubererClasses := []string{
		"Magier", "Hexer", "Thaumaturg", "Krieger Magier",
		"Priester", "Druide", "Schamane",
	}
	for _, class := range zaubererClasses {
		if typ == class {
			return true
		}
	}
	return false
}
