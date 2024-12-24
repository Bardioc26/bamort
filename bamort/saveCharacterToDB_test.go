package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the test database")
	}

	// Auto-migrate the schemas for all related models
	db.AutoMigrate(&Character{}, &Fertigkeit{}, &Zauber{}, &Lp{}, &Eigenschaft{}, &Merkmale{}) //, &Ausruestung{}, &Behaeltniss{}, &Waffenfertigkeit{},
	// &Waffe{}, &Ap{}, &B{}, &Transportmittel{}, &Erfahrungsschatz{}, &Magisch{})
	return db
}

func TestSaveCharacterToDB(t *testing.T) {
	// Setup test database
	testDB := SetupTestDB()
	DB = testDB // Assign test DB to global DB

	// Define a sample character for testing
	character := &Character{
		Name:    "Test Character",
		Rasse:   "Elf",
		Typ:     "Mage",
		Alter:   100,
		Anrede:  "Lord",
		Grad:    5,
		Groesse: 180,
		Gewicht: 70,
		Glaube:  "None",
		Hand:    "Right",
		Eigenschaften: []Eigenschaft{
			{Name: "Au", Value: 50},
			{Name: "St", Value: 80},
			{Name: "Zt", Value: 100},
		},
		/*
			Ausruestung: []Ausruestung{
					{Name: "Staff", Beschreibung: "Magic Staff", Anzahl: 1, Gewicht: 2.5, Wert: 500},
				},
				Behaeltnisse: []Behaeltniss{
					{Name: "Backpack", Beschreibung: "Leather backpack", Gewicht: 1.5, Tragkraft: 10, Volumen: 20, Wert: 50},
				},
		*/
		Fertigkeiten: []Fertigkeit{
			{Name: "Stehlen", Beschreibung: "jemandem etwas wegnehmen ohne das der es merkt", Fertigkeitswert: 6},
			{Name: "Geländelauf", Beschreibung: "Lauf um Hindernisse herum", Fertigkeitswert: 12},
		},
		Zauber: []Zauber{
			{Name: "Fireball", Beschreibung: "Cast a fireball", Bonus: 0, Quelle: "Ark 20"},
		},
		Lp: Lp{
			Max:   100,
			Value: 80,
		},
		Merkmale: Merkmale{
			Augenfarbe: "Blau",
			Haarfarbe:  "Blonde",
			Sonstige:   "Scar on the left cheek",
		},
		/*
			Ap: Ap{
				Max:   50,
				Value: 40,
			},
		*/
	}

	fmt.Println(character)

	// Call the function being tested
	err := saveCharacterToDB(character)
	assert.NoError(t, err, "Expected no error when saving character to DB")
	fmt.Println(character)

	// Verify that the character was saved
	var savedCharacter Character
	//err = DB.Preload("Eigenschaften").Preload("Ausruestung").Preload("Behaeltnisse").
	//	Preload("Fertigkeiten").Preload("Merkmale").Preload("Lp").Preload("Ap").
	err = DB.Preload("Fertigkeiten").Preload("Zauber").Preload("Lp").
		Preload("Eigenschaften").Preload("Merkmale").
		First(&savedCharacter, "name = ?", "Test Character").Error
	assert.NoError(t, err, "Expected to find the character in the database")
	assert.Equal(t, "Test Character", savedCharacter.Name)
	assert.Equal(t, "Elf", savedCharacter.Rasse)
	assert.Equal(t, "Stehlen", savedCharacter.Fertigkeiten[0].Name)
	assert.Equal(t, "Fireball", savedCharacter.Zauber[0].Name)
	assert.Equal(t, 80, savedCharacter.Lp.Value)
	assert.Equal(t, 3, len(savedCharacter.Eigenschaften))
	assert.Equal(t, "Au", savedCharacter.Eigenschaften[0].Name)
	assert.Equal(t, 50, savedCharacter.Eigenschaften[0].Value)
	assert.Equal(t, "Blau", savedCharacter.Merkmale.Augenfarbe)
	/*
		assert.Equal(t, 1, len(savedCharacter.Ausruestung))
		assert.Equal(t, "Staff", savedCharacter.Ausruestung[0].Name)
		assert.Equal(t, "Blue", savedCharacter.Merkmale.Augenfarbe)
		assert.Equal(t, 40, savedCharacter.Ap.Value)
	*/
}
