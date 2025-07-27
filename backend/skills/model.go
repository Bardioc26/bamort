package skills

import (
	"bamort/gsmaster"
	"bamort/models"
)

var dbPrefix = "skill"

type Fertigkeit struct {
	models.BamortCharTrait
	Beschreibung    string             `json:"beschreibung"`
	Fertigkeitswert int                `json:"fertigkeitswert"`
	Bonus           int                `json:"bonus,omitempty"`
	Pp              int                `json:"pp,omitempty"` //Praxispunkte
	Bemerkung       string             `json:"bemerkung"`
	Improvable      bool               `json:"improvable"`
	Category        string             `json:"category"`
	LearningCost    gsmaster.LearnCost `gorm:"-" json:"learncost"`
}

type Waffenfertigkeit struct {
	Fertigkeit
}

type AngeboreneFertigkeit struct {
	Fertigkeit
}

type Zauber struct {
	models.BamortCharTrait
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

func (object *Fertigkeit) TableName() string {
	return dbPrefix + "_" + "skills"
}
func (object *Waffenfertigkeit) TableName() string {
	return dbPrefix + "_" + "weaponskills"
}
func (object *Zauber) TableName() string {
	return dbPrefix + "_" + "spells"
}

func (object *Fertigkeit) GetGsm() *models.Skill {
	var gsmsk models.Skill
	gsmsk.First(object.Name)
	if gsmsk.ID == 0 {
		return nil
	}
	return &gsmsk
}

func (object *Fertigkeit) GetCategory() string {
	if object.Category != "" {
		return object.Category
	}
	var gsmsk models.Skill
	gsmsk.First(object.Name)
	if gsmsk.ID == 0 {
		return "Unkategorisiert"
	}
	object.Category = gsmsk.Category
	return object.Category
}
