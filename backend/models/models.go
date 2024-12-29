package models

import (
	"bamort/database"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID       uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"name"`
	PasswordHash string    `json:"password"`
	Email        string    `gorm:"unique" json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type BamortBase struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	System string `json:"gamingsystem"`
}

type BamortCharTrait struct {
	BamortBase
	CharacterID uint `gorm:"index" json:"character_id"`
}

type BamortCharTraitMaxVal struct {
	BamortCharTrait
	Max   int `json:"max"`
	Value int `json:"value"`
}

type Char struct {
	BamortBase

	Rasse              string               `json:"rasse"`
	Typ                string               `json:"typ"`
	Alter              int                  `json:"alter"`
	Anrede             string               `json:"anrede"`
	Grad               int                  `json:"grad"`
	Groesse            int                  `json:"groesse"`
	Gewicht            int                  `json:"gewicht"`
	Glaube             string               `json:"glaube"`
	Hand               string               `json:"hand"`
	Lp                 Lp                   `gorm:"foreignKey:CharacterID" json:"lp"`
	Ap                 Ap                   `gorm:"foreignKey:CharacterID" json:"ap"`
	B                  B                    `gorm:"foreignKey:CharacterID" json:"b"`
	Merkmale           Merkmale             `gorm:"foreignKey:CharacterID" json:"merkmale"`
	Eigenschaften      []Eigenschaft        `gorm:"foreignKey:CharacterID" json:"eigenschaften"`
	Fertigkeiten       []Fertigkeit         `gorm:"foreignKey:CharacterID" json:"fertigkeiten"`
	Waffenfertigkeiten []Waffenfertigkeit   `gorm:"foreignKey:CharacterID" json:"waffenfertigkeiten"`
	Zauber             []Zauber             `gorm:"foreignKey:CharacterID" json:"zauber"`
	Spezialisierung    database.StringArray `gorm:"type:TEXT"  json:"spezialisierung"`
	Bennies            Bennies              `gorm:"foreignKey:CharacterID" json:"bennies"`
	Erfahrungsschatz   Erfahrungsschatz     `gorm:"foreignKey:CharacterID" json:"erfahrungsschatz"`
	Waffen             []Waffe              `gorm:"foreignKey:CharacterID" json:"waffen"`
	Behaeltnisse       []Behaeltniss        `gorm:"foreignKey:CharacterID" json:"behaeltnisse"`
	Transportmittel    []Transportation     `gorm:"foreignKey:CharacterID" json:"transportmittel"`
	Ausruestung        []Ausruestung        `gorm:"foreignKey:CharacterID" json:"ausruestung"`
	Image              string               `json:"image,omitempty"`
}

// Au, Gs, Gw ,In, Ko, Pa, St, Wk, Zt
type Eigenschaft struct {
	BamortCharTrait
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Ausruestung struct {
	BamortCharTrait
	Magisch
	Beschreibung string  `json:"beschreibung"`
	Anzahl       int     `json:"anzahl"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Bonus        int     `json:"bonus,omitempty"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
}

type Fertigkeit struct {
	BamortCharTrait
	Beschreibung    string `json:"beschreibung"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Bonus           int    `json:"bonus,omitempty"`
	Pp              int    `json:"pp,omitempty"`
	Bemerkung       string `json:"bemerkung"`
}

type Waffenfertigkeit struct {
	Fertigkeit
}

type Zauber struct {
	BamortCharTrait
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

type Waffe struct {
	BamortCharTrait
	Beschreibung string  `json:"beschreibung"`
	Abwb         int     `json:"abwb"`
	Anb          int     `json:"anb"`
	Anzahl       int     `json:"anzahl"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Gewicht      float64 `json:"gewicht"`
	Magisch
	NameFuerSpezialisierung string  `json:"nameFuerSpezialisierung"`
	Schb                    int     `json:"schb"`
	Wert                    float64 `json:"wert"`
}

type Merkmale struct {
	BamortCharTrait
	Augenfarbe string `json:"augenfarbe"`
	Haarfarbe  string `json:"haarfarbe"`
	Sonstige   string `json:"sonstige"`
	Breite     string `json:"breite"`
	Groesse    string `json:"groesse"`
}

type Gestalt struct {
	BamortCharTrait
}

type Erfahrungsschatz struct {
	BamortCharTrait
	Value int `json:"value"`
}

type Bennies struct {
	BamortCharTrait
	Gg int `json:"gg"`
	Gp int `json:"gp"`
	Sg int `json:"sg"`
}

type Lp struct {
	BamortCharTraitMaxVal
}

type Ap struct {
	BamortCharTraitMaxVal
}

type B struct {
	BamortCharTraitMaxVal
}

type Behaeltniss struct {
	BamortCharTrait
	Beschreibung string  `json:"beschreibung"`
	BeinhaltetIn string  `json:"beinhaltet_in"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
	Tragkraft    float64 `json:"tragkraft"`
	Volumen      float64 `json:"volumen"`
	Magisch
}

type Transportation struct {
	Behaeltniss
	//Magisch   Magisch `gorm:"polymorphic:Item;polymorphicValue:Transportmittel" json:"magisch"`
}

type Magisch struct {
	IstMagisch  bool `json:"ist_magisch"`
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
}

func (object *Char) First(name string) error {
	gameSystem := "midgard"
	err := database.DB.First(&object, "system=? AND name = ?", gameSystem, name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
}

func (object *Char) Create() error {
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
