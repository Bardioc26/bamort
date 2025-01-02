package skills

import "bamort/models"

var dbPrefix = "skill"

type Fertigkeit struct {
	models.BamortCharTrait
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
