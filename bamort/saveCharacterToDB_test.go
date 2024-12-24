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
	db.AutoMigrate(&Character{},
		&Fertigkeit{}, &Zauber{}, &Lp{},
		&Eigenschaft{}, &Merkmale{},
		&Bennies{},
		&Gestalt{},
		&Ap{}, &B{},
		&Erfahrungsschatz{},
		&MagischTransport{},
		&Transportation{},
		&MagischAusruestung{},
		&Ausruestung{},
		&MagischBehaelter{},
		&Behaeltniss{},
		&MagischWaffe{},
		&Waffe{},
		&Waffenfertigkeit{},
	)
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
		Bennies: Bennies{
			Gg: 1,
			Gp: 0,
			Sg: 2,
		},
		Gestalt: Gestalt{
			Breite:  "schmal",
			Groesse: "klein",
		},
		Ap: Ap{
			Max:   50,
			Value: 40,
		},
		B: B{

			Max:   25,
			Value: 20,
		},
		Erfahrungsschatz: Erfahrungsschatz{
			Value: 2768,
		},
		Transportmittel: []Transportation{
			{Name: "Karren",
				Beschreibung: "ein Karren",
				Gewicht:      100, Tragkraft: 300, Wert: 55,
				Magisch: MagischTransport{IstMagisch: true, Abw: 30, Ausgebrannt: false},
			},
		},
		Ausruestung: []Ausruestung{
			{Name: "Staff", Beschreibung: "Magic Staff", Anzahl: 1, Gewicht: 2.5, Wert: 500,
				Magisch: MagischAusruestung{IstMagisch: true, Abw: 10, Ausgebrannt: false},
			},
		},
		Behaeltnisse: []Behaeltniss{
			{Name: "Backpack", Beschreibung: "Leather backpack",
				Gewicht: 1.5, Tragkraft: 10, Volumen: 20, Wert: 50,
				//Magisch: MagischBehaelter{IstMagisch: false},
			},
		},
		Waffen: []Waffe{
			{Name: "Schwert", Beschreibung: "Ein schwert", Abwb: 0, Anb: 0, Gewicht: 1.5, NameFuerSpezialisierung: "Schwert", Schb: 0, Wert: 3,
				Magisch: MagischWaffe{IstMagisch: false}},
		},
		Waffenfertigkeiten: []Waffenfertigkeit{
			{Name: "Einhandschlagwaffe", Beschreibung: "z.B. für Kurzschwerter", Bonus: 0,
				Fertigkeitswert: 12, Pp: 1, Quelle: "Kod 256"},
		},
		Spezialisierung: []string{
			"Bogen", "Streitaxt",
		},
	}

	fmt.Println(character)

	// Call the function being tested
	err := saveCharacterToDB(character)
	assert.NoError(t, err, "Expected no error when saving character to DB")
	//fmt.Println(character)

	// Verify that the character was saved
	var savedCharacter Character
	//err = DB.Preload("Eigenschaften").Preload("Ausruestung").Preload("Behaeltnisse").
	//	Preload("Fertigkeiten").Preload("Merkmale").Preload("Lp").Preload("Ap").
	err = DB.
		Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Zauber").
		Preload("Lp").
		Preload("Merkmale").
		Preload("Bennies").
		Preload("Gestalt").
		Preload("Ap").
		Preload("B").
		Preload("Erfahrungsschatz").
		//Preload("Magisch").
		Preload("Transportmittel").
		Preload("Ausruestung").
		Preload("Behaeltnisse").
		Preload("Waffen").
		Preload("Waffenfertigkeiten").
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
	assert.Equal(t, 1, len(savedCharacter.Ausruestung))
	assert.Equal(t, "Staff", savedCharacter.Ausruestung[0].Name)
	assert.Equal(t, "Blau", savedCharacter.Merkmale.Augenfarbe)
	assert.Equal(t, "Backpack", savedCharacter.Behaeltnisse[0].Name)
	assert.Equal(t, "Schwert", savedCharacter.Waffen[0].Name)
	assert.Equal(t, 40, savedCharacter.Ap.Value)
	assert.Equal(t, "Einhandschlagwaffe", savedCharacter.Waffenfertigkeiten[0].Name)
	assert.Equal(t, 2, len(savedCharacter.Spezialisierung))
	assert.Equal(t, "Bogen", savedCharacter.Spezialisierung[0])
	assert.Equal(t, "Streitaxt", savedCharacter.Spezialisierung[1])

	/*
	 */
}
