package models

type BamortBase struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type BamortCharTrait struct {
	BamortBase
	CharacterID uint `gorm:"index" json:"character_id"`
}

type Magisch struct {
	IstMagisch  bool `json:"ist_magisch"`
	Abw         int  `json:"abw"`
	Ausgebrannt bool `json:"ausgebrannt"`
}
