package main

type User struct {
	UserID       uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string `gorm:"unique"`
}

type Character struct {
	CharacterID uint `gorm:"primaryKey"`
	UserID      uint `gorm:"index"`
	Name        string
	Rasse       string
	Typ         string
	Alter       int
	Anrede      string
	Grad        int
	Groesse     int
	Gewicht     int
	Glaube      string
	Hand        string
	Image       string
}

type Eigenschaften struct {
	EigenschaftenID uint `gorm:"primaryKey"`
	CharacterID     uint `gorm:"index"`
	Au              int
	Gs              int
	Gw              int
	In              int
	Ko              int
	Pa              int
	St              int
	Wk              int
	Zt              int
}

type Ausruestung struct {
	AusruestungID uint    `gorm:"primaryKey"`
	CharacterID   uint    `gorm:"index"`
	Name          string  `json:"name"`
	Anzahl        int     `json:"anzahl"`
	Gewicht       float64 `json:"gewicht"`
	Wert          float64 `json:"wert"`
	BeinhaltetIn  *string `json:"beinhaltet_in"`
	Beschreibung  *string `json:"beschreibung"`
	Bonus         *int    `json:"bonus"`
	IstMagisch    bool    `json:"ist_magisch"`
	Abw           *int    `json:"abw"`
	Ausgebrannt   bool    `json:"ausgebrannt"`
}

/*
Define models for each table
Add other models for Ausruestung, Fertigkeiten, etc., following the same pattern.
*/
