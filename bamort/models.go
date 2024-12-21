package main

type User struct {
	UserID       uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique"`
	PasswordHash string
	Email        string `gorm:"unique"`
}

type Character struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	Name          string        `json:"name"`
	Rasse         string        `json:"rasse"`
	Typ           string        `json:"typ"`
	Alter         int           `json:"alter"`
	Anrede        string        `json:"anrede"`
	Grad          int           `json:"grad"`
	Groesse       int           `json:"groesse"`
	Gewicht       int           `json:"gewicht"`
	Glaube        string        `json:"glaube"`
	Hand          string        `json:"hand"`
	Eigenschaften []Eigenschaft `gorm:"foreignKey:CharacterID" json:"eigenschaften"`
	Ausruestung   []Ausruestung `gorm:"foreignKey:CharacterID" json:"ausruestung"`
	Fertigkeiten  []Fertigkeit  `gorm:"foreignKey:CharacterID" json:"fertigkeiten"`
}

type Eigenschaft struct {
	ID          uint   `gorm:"primaryKey"`
	CharacterID uint   `gorm:"index" json:"character_id"`
	Name        string `json:"name"`
	Value       int    `json:"value"`
}

type Ausruestung struct {
	ID           uint    `gorm:"primaryKey"`
	CharacterID  uint    `gorm:"index" json:"character_id"`
	Name         string  `json:"name"`
	Anzahl       int     `json:"anzahl"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
	Beschreibung string  `json:"beschreibung"`
}

type Fertigkeit struct {
	ID              uint   `gorm:"primaryKey"`
	CharacterID     uint   `gorm:"index" json:"character_id"`
	Name            string `json:"name"`
	Fertigkeitswert int    `json:"fertigkeitswert"`
	Beschreibung    string `json:"beschreibung"`
}

/*
Define models for each table
Add other models for Ausruestung, Fertigkeiten, etc., following the same pattern.
*/
