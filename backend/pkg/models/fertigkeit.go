package models

type Fertigkeit struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	CharacterID     uint   `gorm:"index" json:"character_id"`
	Name            string `json:"name"`
	Beschreibung    string `json:"beschreibung"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Bonus           int    `json:"bonus,omitempty"`
	// ...
}
