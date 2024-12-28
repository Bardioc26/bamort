package models

type ImStammFertigkeit struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	System           string `gorm:"index" json:"system"`
	Name             string `json:"name"`
	Beschreibung     string `json:"beschreibung"`
	Initialkeitswert int    `json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
	Quelle           string `json:"quelle"`
}

type ImStammZauber struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	System       string `gorm:"index" json:"system"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Bonus        int    `json:"bonus"`
	Quelle       string `json:"quelle"`
	Stufe        int
	AP           int
	Reichweite   int
	Wirkungsziel string
}
