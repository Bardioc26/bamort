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
	Goldstücke   int `json:"goldstücke"`   // GS
	Silberstücke int `json:"silberstücke"` // SS
	Kupferstücke int `json:"kupferstücke"` // KS
}

type Char struct {
	BamortBase
	UserID             uint                 `gorm:"index;not null" json:"user_id"`
	User               user.User            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Rasse              string               `json:"rasse"`
	Typ                string               `json:"typ"`
	Alter              int                  `json:"alter"`
	Anrede             string               `json:"anrede"`
	Grad               int                  `json:"grad"`
	Groesse            int                  `json:"groesse"`
	Gewicht            int                  `json:"gewicht"`
	Glaube             string               `json:"glaube"`
	Hand               string               `json:"hand"`
	Public             bool                 `json:"public"`
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

// FindCharListByUserID finds all characters belonging to a specific user for listing (minimal data)
func FindCharListByUserID(userID uint) ([]CharList, error) {
	var chars []CharList
	err := database.DB.Table("char_chars").
		Select("id, name, user_id, rasse, typ, grad, owner, public").
		Where("user_id = ?", userID).
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
