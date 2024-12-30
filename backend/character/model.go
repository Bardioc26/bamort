package character

import (
	"bamort/database"
	"bamort/models"
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
	ID uint `gorm:"primaryKey" json:"id"`

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
type Char struct {
	models.BamortBase
	Rasse              string                    `json:"rasse"`
	Typ                string                    `json:"typ"`
	Alter              int                       `json:"alter"`
	Anrede             string                    `json:"anrede"`
	Grad               int                       `json:"grad"`
	Groesse            int                       `json:"groesse"`
	Gewicht            int                       `json:"gewicht"`
	Glaube             string                    `json:"glaube"`
	Hand               string                    `json:"hand"`
	Lp                 Lp                        `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"lp"`
	Ap                 Ap                        `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"ap"`
	B                  B                         `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"b"`
	Merkmale           models.Merkmale           `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"merkmale"`
	Eigenschaften      []Eigenschaft             `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"eigenschaften"`
	Fertigkeiten       []models.Fertigkeit       `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"fertigkeiten"`
	Waffenfertigkeiten []models.Waffenfertigkeit `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"waffenfertigkeiten"`
	Zauber             []models.Zauber           `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"zauber"`
	Spezialisierung    database.StringArray      `gorm:"type:TEXT"  json:"spezialisierung"`
	Bennies            models.Bennies            `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"bennies"`
	Erfahrungsschatz   models.Erfahrungsschatz   `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"erfahrungsschatz"`
	Waffen             []models.Waffe            `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"waffen"`
	Behaeltnisse       []models.Behaeltniss      `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"behaeltnisse"`
	Transportmittel    []models.Transportation   `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"transportmittel"`
	Ausruestung        []models.Ausruestung      `gorm:"foreignKey:CharacterID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"ausruestung"`
	Image              string                    `json:"image,omitempty"`
}

func (object *Char) First(name string) error {
	err := database.DB.
		Preload("Lp").
		Preload("Ap").
		Preload("B").
		Preload("Merkmale").
		Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Zauber").
		Preload("Bennies").
		Preload("Erfahrungsschatz").
		Preload("Waffen").
		Preload("Behaeltnisse").
		Preload("Transportmittel").
		Preload("Ausruestung").
		First(&object, " name = ?", name).Error
	if err != nil {
		// zauber found
		return err
	}
	return nil
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
