package models

type BamortBase struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
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

type Bennies struct {
	BamortCharTrait
	Gg int `json:"gg"`
	Gp int `json:"gp"`
	Sg int `json:"sg"`
}

type Magisch struct {
	IstMagisch  bool `json:"ist_magisch"`
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
}
