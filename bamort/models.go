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

/*
Define models for each table
Add other models for Ausruestung, Fertigkeiten, etc., following the same pattern.
*/
