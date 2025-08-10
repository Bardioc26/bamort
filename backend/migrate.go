package main

import (
	"bamort/database"
	"bamort/models"
	"fmt"
)

func main() {
	fmt.Println("Starte Migration...")

	database.ConnectDatabase()

	err := database.DB.AutoMigrate(&models.CharacterCreationSession{})
	if err != nil {
		fmt.Printf("Migration Fehler: %v\n", err)
	} else {
		fmt.Println("Migration erfolgreich!")
	}
}
