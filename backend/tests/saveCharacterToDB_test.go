package main

import (
	"bamort/character"
	"bamort/database"
	"bamort/models"
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
	db.AutoMigrate(&models.Char{},
		&models.Fertigkeit{}, &models.Zauber{}, &models.Lp{},
		&models.Eigenschaft{}, &models.Merkmale{},
		&models.Bennies{},
		&models.Gestalt{},
		&models.Ap{}, &models.B{},
		&models.Erfahrungsschatz{},
		&models.MagischTransport{},
		&models.Transportation{},
		&models.MagischAusruestung{},
		&models.Ausruestung{},
		&models.MagischBehaelter{},
		&models.Behaeltniss{},
		&models.MagischWaffe{},
		&models.Waffe{},
		&models.Waffenfertigkeit{},
		&models.ImStammFertigkeit{},
	)
	return db
}

func TestSaveCharacterToDB(t *testing.T) {
	// Setup test database
	testDB := SetupTestDB()
	database.DB = testDB // Assign test DB to global DB

	// Define a sample character for testing
	char := &models.Char{
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
		Eigenschaften: []models.Eigenschaft{
			{Name: "Au", Value: 50},
			{Name: "St", Value: 80},
			{Name: "Zt", Value: 100},
		},
		Fertigkeiten: []models.Fertigkeit{
			{Name: "Stehlen", Beschreibung: "jemandem etwas wegnehmen ohne das der es merkt", Fertigkeitswert: 6},
			{Name: "Geländelauf", Beschreibung: "Lauf um Hindernisse herum", Fertigkeitswert: 12},
		},
		Zauber: []models.Zauber{
			{Name: "Fireball", Beschreibung: "Cast a fireball", Bonus: 0, Quelle: "Ark 20"},
		},
		Lp: models.Lp{
			Max:   100,
			Value: 80,
		},
		Merkmale: models.Merkmale{
			Augenfarbe: "Blau",
			Haarfarbe:  "Blonde",
			Sonstige:   "Scar on the left cheek",
		},
		Bennies: models.Bennies{
			Gg: 1,
			Gp: 0,
			Sg: 2,
		},
		Gestalt: models.Gestalt{
			Breite:  "schmal",
			Groesse: "klein",
		},
		Ap: models.Ap{
			Max:   50,
			Value: 40,
		},
		B: models.B{

			Max:   25,
			Value: 20,
		},
		Erfahrungsschatz: models.Erfahrungsschatz{
			Value: 2768,
		},
		Transportmittel: []models.Transportation{
			{Name: "Karren",
				Beschreibung: "ein Karren",
				Gewicht:      100, Tragkraft: 300, Wert: 55,
				Magisch: models.MagischTransport{IstMagisch: true, Abw: 30, Ausgebrannt: false},
			},
		},
		Ausruestung: []models.Ausruestung{
			{Name: "Staff", Beschreibung: "Magic Staff", Anzahl: 1, Gewicht: 2.5, Wert: 500,
				Magisch: models.MagischAusruestung{IstMagisch: true, Abw: 10, Ausgebrannt: false},
			},
		},
		Behaeltnisse: []models.Behaeltniss{
			{Name: "Backpack", Beschreibung: "Leather backpack",
				Gewicht: 1.5, Tragkraft: 10, Volumen: 20, Wert: 50,
				//Magisch: MagischBehaelter{IstMagisch: false},
			},
		},
		Waffen: []models.Waffe{
			{Name: "Schwert", Beschreibung: "Ein schwert", Abwb: 0, Anb: 0, Gewicht: 1.5, NameFuerSpezialisierung: "Schwert", Schb: 0, Wert: 3,
				Magisch: models.MagischWaffe{IstMagisch: false}},
		},
		Waffenfertigkeiten: []models.Waffenfertigkeit{
			{Name: "Einhandschlagwaffe", Beschreibung: "z.B. für Kurzschwerter", Bonus: 0,
				Fertigkeitswert: 12, Pp: 1, Quelle: "Kod 256"},
		},
		Spezialisierung: []string{
			"Bogen", "Streitaxt",
		},
	}

	//fmt.Println(char)

	// Call the function being tested
	err := character.SaveCharacterToDB(char)
	assert.NoError(t, err, "Expected no error when saving character to DB")
	//fmt.Println(char)

	// Verify that the character was saved
	var savedChar models.Char
	//err = DB.Preload("Eigenschaften").Preload("Ausruestung").Preload("Behaeltnisse").
	//	Preload("Fertigkeiten").Preload("Merkmale").Preload("Lp").Preload("Ap").
	err = database.DB.
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
		First(&savedChar, "name = ?", "Test Character").Error
	assert.NoError(t, err, "Expected to find the character in the database")
	assert.Equal(t, "Test Character", savedChar.Name)
	assert.Equal(t, "Elf", savedChar.Rasse)
	assert.Equal(t, "Stehlen", savedChar.Fertigkeiten[0].Name)
	assert.Equal(t, "Fireball", savedChar.Zauber[0].Name)
	assert.Equal(t, 80, savedChar.Lp.Value)
	assert.Equal(t, 3, len(savedChar.Eigenschaften))
	assert.Equal(t, "Au", savedChar.Eigenschaften[0].Name)
	assert.Equal(t, 50, savedChar.Eigenschaften[0].Value)
	assert.Equal(t, "Blau", savedChar.Merkmale.Augenfarbe)
	assert.Equal(t, 1, len(savedChar.Ausruestung))
	assert.Equal(t, "Staff", savedChar.Ausruestung[0].Name)
	assert.Equal(t, "Blau", savedChar.Merkmale.Augenfarbe)
	assert.Equal(t, "Backpack", savedChar.Behaeltnisse[0].Name)
	assert.Equal(t, "Schwert", savedChar.Waffen[0].Name)
	assert.Equal(t, 40, savedChar.Ap.Value)
	assert.Equal(t, "Einhandschlagwaffe", savedChar.Waffenfertigkeiten[0].Name)
	assert.Equal(t, 2, len(savedChar.Spezialisierung))
	assert.Equal(t, "Bogen", savedChar.Spezialisierung[0])
	assert.Equal(t, "Streitaxt", savedChar.Spezialisierung[1])
}
