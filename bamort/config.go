package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "bamort:bG4)efozrc@tcp(192.168.0.5:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = database
}

/*
Replace user, password, and dbname with your MySQL credentials and database name.
*/

func saveCharacterToDB(character *Character) error {
	// Use GORM to save the character and its relationships
	err := DB.Transaction(func(tx *gorm.DB) error {
		// Save the main character
		if err := tx.Create(character).Error; err != nil {
			return err
		}

		// Save Eigenschaften (Attributes)
		for i := range character.Eigenschaften {
			character.Eigenschaften[i].CharacterID = character.ID
		}
		if err := tx.Create(&character.Eigenschaften).Error; err != nil {
			return err
		}

		// Save Ausruestung (Equipment)
		for i := range character.Ausruestung {
			character.Ausruestung[i].CharacterID = character.ID
		}
		if err := tx.Create(&character.Ausruestung).Error; err != nil {
			return err
		}

		// Save Fertigkeiten (Skills)
		for i := range character.Fertigkeiten {
			character.Fertigkeiten[i].CharacterID = character.ID
		}
		if err := tx.Create(&character.Fertigkeiten).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
