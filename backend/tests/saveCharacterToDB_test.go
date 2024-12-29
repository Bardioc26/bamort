package tests

import (
	"bamort/models"

	"gorm.io/gorm"
)

func initTestDB4Char() *gorm.DB {
	db := SetupTestDB()
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
		&models.LookupSkill{},
	)
	return db
}

/*
func TestSaveCharacterToDB(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Char()
	database.DB = testDB // Assign test DB to global DB

	// Define a sample character for testing
	char := &models.Char{}
	char.Name=    "Test Character"
	char.Rasse =   "Elf"
	char.Typ =     "Mage"
	char.Alter =   100
	char.Anrede =  "Lord"
	char.Grad =    5
	char.Groesse = 180
	char.Gewicht = 70
	char.Glaube =  "None"
	char.Hand =    "Right"
	char.Lp.Max=100
	char.Lp.Value=80
	char.Eigenschaften = []models.Eigenschaft{
		{Name = "Au", Value = 50}
			{Name = "St", Value = 80}
			{Name = "Zt", Value = 100}
		}
	char.Fertigkeiten = []models.Fertigkeit{
		{Name = "Stehlen", Beschreibung = "jemandem etwas wegnehmen ohne das der es merkt", Fertigkeitswert = 6}
			{Name = "Geländelauf", Beschreibung = "Lauf um Hindernisse herum", Fertigkeitswert = 12}
		}
	char.Zauber = []models.Zauber{
		{Name = "Fireball", Beschreibung = "Cast a fireball", Bonus = 0, Quelle = "Ark 20"}
		}
	char.Lp = models.Lp{
		Max =   100
			Value = 80
		}
	char.Merkmale = models.Merkmale{
		Augenfarbe = "Blau"
			Haarfarbe =  "Blonde"
			Sonstige =   "Scar on the left cheek"
		}
	char.Bennies = models.Bennies{
		Gg = 1
			Gp = 0
			Sg = 2
		}
	char.Gestalt = models.Gestalt{
		Breite =  "schmal"
			Groesse = "klein"
		}
	char.Ap = models.Ap{
		Max =   50
			Value = 40
		}
	char.B = models.B{
		Max =   25
			Value = 20
		}
	char.Erfahrungsschatz = models.Erfahrungsschatz{
		Value = 2768
		}
	char.Transportmittel = []models.Transportation{
		{Name = "Karren"
				Beschreibung = "ein Karren"
				Gewicht =      100, Tragkraft = 300, Wert = 55
				Magisch = models.MagischTransport{IstMagisch = true, Abw = 30, Ausgebrannt = false}
			}
		}
	char.Ausruestung = []models.Ausruestung{
		{Name = "Staff", Beschreibung = "Magic Staff", Anzahl = 1, Gewicht = 2.5, Wert = 500
				Magisch = models.MagischAusruestung{IstMagisch = true, Abw = 10, Ausgebrannt = false}
			}
		}
	char.Behaeltnisse = []models.Behaeltniss{
		{Name = "Backpack", Beschreibung = "Leather backpack"
				Gewicht = 1.5, Tragkraft = 10, Volumen = 20, Wert = 50
				//Magisch = MagischBehaelter{IstMagisch = false}
			}
		}
	char.Waffen = []models.Waffe{
		{Name = "Schwert", Beschreibung = "Ein schwert", Abwb = 0, Anb = 0, Gewicht = 1.5, NameFuerSpezialisierung = "Schwert", Schb = 0, Wert = 3
				Magisch = models.MagischWaffe{IstMagisch = false}}
		}
	char.Waffenfertigkeiten = []models.Waffenfertigkeit{
		{Name = "Einhandschlagwaffe", Beschreibung = "z.B. für Kurzschwerter", Bonus = 0
				Fertigkeitswert = 12, Pp = 1, Quelle = "Kod 256"}
		}
	char.Spezialisierung = []string{
		"Bogen", "Streitaxt"
		}
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
*/
