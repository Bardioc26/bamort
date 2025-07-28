package main

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
	"bamort/models"
	"bamort/skills"
	"bamort/user"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Beispiel 1: Verwende die Standard-Datenbank (wie bisher)
	fmt.Println("=== Standard Datenbank Migration ===")
	err := user.MigrateStructure()
	if err != nil {
		fmt.Printf("Fehler bei Standard-Migration: %v\n", err)
		return
	}
	fmt.Println("Standard-Migration erfolgreich")

	// Beispiel 2: Verwende eine zweite Datenbank
	fmt.Println("\n=== Zweite Datenbank Migration ===")

	// Erstelle eine zweite Datenbankverbindung
	secondDB, err := gorm.Open(sqlite.Open("second_database.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Fehler beim Öffnen der zweiten DB: %v\n", err)
		return
	}

	// Migriere alle Packages zur zweiten Datenbank
	err = database.MigrateStructure(secondDB)
	if err != nil {
		fmt.Printf("Fehler bei database Migration: %v\n", err)
		return
	}

	err = user.MigrateStructure(secondDB)
	if err != nil {
		fmt.Printf("Fehler bei user Migration: %v\n", err)
		return
	}

	err = character.MigrateStructure(secondDB)
	if err != nil {
		fmt.Printf("Fehler bei character Migration: %v\n", err)
		return
	}

	err = models.MigrateStructure(secondDB)
	if err != nil {
		fmt.Printf("Fehler bei models Migration: %v\n", err)
		return
	}

	err = equipment.MigrateStructure(secondDB)
	if err != nil {
		fmt.Printf("Fehler bei equipment Migration: %v\n", err)
		return
	}

	err = skills.MigrateStructure(secondDB)
	if err != nil {
		fmt.Printf("Fehler bei skills Migration: %v\n", err)
		return
	}

	fmt.Println("Migration zur zweiten Datenbank erfolgreich")

	// Beispiel 3: Helper-Funktion für komplette Migration
	fmt.Println("\n=== Helper Funktion Beispiel ===")
	err = migrateAllToSecondDB(secondDB)
	if err != nil {
		fmt.Printf("Fehler bei kompletter Migration: %v\n", err)
		return
	}
	fmt.Println("Komplette Migration erfolgreich")
}

// Helper-Funktion für komplette Migration zu einer spezifischen Datenbank
func migrateAllToSecondDB(db *gorm.DB) error {
	migrators := []func(db ...*gorm.DB) error{
		database.MigrateStructure,
		user.MigrateStructure,
		character.MigrateStructure,
		models.MigrateStructure,
		equipment.MigrateStructure,
		skills.MigrateStructure,
	}

	for _, migrate := range migrators {
		if err := migrate(db); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}
	return nil
}
