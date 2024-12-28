package models

/*type ImStammFertigkeit struct {
	ID               uint   `gorm:"primaryKey" json:"id"`
	System           string `gorm:"index" json:"system"`
	Name             string `json:"name"`
	Beschreibung     string `json:"beschreibung"`
	Initialkeitswert int    `json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
	Quelle           string `json:"quelle"`
}
*/

type ImStamm struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	System       string `gorm:"index" json:"system"`
	Name         string `json:"name"`
	Beschreibung string `json:"beschreibung"`
	Quelle       string `json:"quelle"`
}

type ImStammFertigkeit struct {
	ImStamm
	Initialkeitswert int    `json:"initialwert"`
	Bonuseigenschaft string `json:"bonuseigenschaft,omitempty"`
}

type ImStammZauber struct {
	ImStamm
	Bonus        int `json:"bonus"`
	Stufe        int
	AP           int
	Reichweite   int
	Wirkungsziel string
}
