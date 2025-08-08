package models

type SkFertigkeit struct {
	BamortCharTrait
	Beschreibung    string    `json:"beschreibung"`
	Fertigkeitswert int       `json:"fertigkeitswert"`
	Bonus           int       `json:"bonus,omitempty"`
	Pp              int       `json:"pp,omitempty"` //Praxispunkte
	Bemerkung       string    `json:"bemerkung"`
	Improvable      bool      `json:"improvable"`
	Category        string    `json:"category"`
	LearningCost    LearnCost `gorm:"-" json:"learncost"`
}

type SkWaffenfertigkeit struct {
	SkFertigkeit
}

type SkAngeboreneFertigkeit struct {
	SkFertigkeit
}

type SkZauber struct {
	BamortCharTrait
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
}

func (object *SkFertigkeit) TableName() string {
	dbPrefix := "skill"
	return dbPrefix + "_" + "skills"
}
func (object *SkWaffenfertigkeit) TableName() string {
	dbPrefix := "skill"
	return dbPrefix + "_" + "weaponskills"
}
func (object *SkZauber) TableName() string {
	dbPrefix := "skill"
	return dbPrefix + "_" + "spells"
}

func (object *SkFertigkeit) GetSkillByName() *Skill {
	var gsmsk Skill
	gsmsk.First(object.Name)
	if gsmsk.ID == 0 {
		return nil
	}
	return &gsmsk
}

func (object *SkFertigkeit) GetCategory() string {
	if object.Category != "" {
		return object.Category
	}
	var gsmsk Skill
	gsmsk.First(object.Name)
	if gsmsk.ID == 0 {
		return "Unkategorisiert"
	}
	object.Category = gsmsk.Category
	return object.Category
}
