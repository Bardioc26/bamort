package tests

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
	"bamort/models"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initTestDB4Character() *gorm.DB {
	db := SetupTestDB()
	db.AutoMigrate(
		&character.Char{},
		&character.Lp{},
		&character.Ap{},
		&character.B{},
		&models.Merkmale{},
		&character.Eigenschaft{},
		&models.Fertigkeit{},
		&models.Waffenfertigkeit{},
		&models.Zauber{},
		&models.Bennies{},
		&models.Erfahrungsschatz{},
		&equipment.Waffe{},
		&equipment.Behaeltniss{},
		&equipment.Transportation{},
		&equipment.Ausruestung{},
	)
	return db
}

// ReadImageAsBase64 reads an image file and returns it as a Base64 string
// with the prefix "data:mimeType;base64,"
func ReadImageAsBase64(filePath string) (string, error) {
	// Read the image file into bytes
	imageBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// Detect the file extension
	ext := filepath.Ext(filePath)
	mimeType := "image" // Default MIME type
	switch ext {
	case ".jpg", ".jpeg":
		mimeType = "image/jpeg"
	case ".png":
		mimeType = "image/png"
	case ".gif":
		mimeType = "image/gif"
	case ".bmp":
		mimeType = "image/bmp"
	case ".svg":
		mimeType = "image/svg+xml"
	default:
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}

	// Encode the bytes to a Base64 string
	base64String := base64.StdEncoding.EncodeToString(imageBytes)

	// Prefix with the MIME type
	fullBase64String := fmt.Sprintf("data:%s;base64,%s", mimeType, base64String)

	return fullBase64String, nil
}

func createChar() *character.Char {

	char := character.Char{}
	char.Name = "Harsk Hammerhuter, Zen"
	char.Rasse = "Zwerg"
	char.Typ = "Krieger"
	char.Alter = 39
	char.Anrede = "er"
	char.Grad = 3
	char.Groesse = 140
	char.Gewicht = 82
	char.Glaube = "Torkin"
	char.Hand = "rechts"
	char.Lp.Max = 17
	char.Lp.Value = 17
	char.Ap.Max = 31
	char.Ap.Value = 31
	char.B.Max = 18
	char.B.Value = 18
	char.Merkmale.Augenfarbe = "blau"
	char.Merkmale.Haarfarbe = "sandfarben"
	char.Merkmale.Sonstige = ""
	char.Eigenschaften = []character.Eigenschaft{
		{Name: "Au", Value: 74},
		{Name: "Gs", Value: 96},
		{Name: "Gw", Value: 70},
		{Name: "In", Value: 65},
		{Name: "Ko", Value: 85},
		{Name: "PA", Value: 75},
		{Name: "St", Value: 95},
		{Name: "Wk", Value: 71},
		{Name: "Zt", Value: 35},
	}
	char.Fertigkeiten = []models.Fertigkeit{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Hören",
				},
			},
			Beschreibung:    "xx",
			Fertigkeitswert: 6,
			Bonus:           0,
			Pp:              0,
			Bemerkung:       "",
		},
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Sprache",
				},
			},
			Beschreibung:    "Albisch",
			Fertigkeitswert: 8,
			Bonus:           0,
			Pp:              0,
			Bemerkung:       "",
		},
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Sprache",
				},
			},
			Beschreibung:    "Comentang",
			Fertigkeitswert: 12,
			Bonus:           0,
			Pp:              0,
			Bemerkung:       "",
		},
	}
	char.Waffenfertigkeiten = []models.Waffenfertigkeit{
		{
			Fertigkeit: models.Fertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Armbr\u00fcste",
					},
				},
				Beschreibung:    "",
				Fertigkeitswert: 8,
				Bonus:           0,
				Pp:              0,
				Bemerkung:       "",
			},
		},
		{
			Fertigkeit: models.Fertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Einhandschlagwaffen",
					},
				},
				Beschreibung:    "",
				Fertigkeitswert: 8,
				Bonus:           0,
				Pp:              0,
				Bemerkung:       "",
			},
		},
		{
			Fertigkeit: models.Fertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Schilde",
					},
				},
				Beschreibung:    "",
				Fertigkeitswert: 3,
				Bonus:           0,
				Pp:              0,
				Bemerkung:       "",
			},
		},
	}
	char.Zauber = []models.Zauber{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Angst",
				},
			},
			Beschreibung: "",
			Bonus:        0,
		},
	}
	char.Spezialisierung = database.StringArray{
		"Kriegshammer", "Armbrust:schwer", "Stielhammer",
	}
	char.Bennies = models.Bennies{
		Sg: 1,
		Gg: 0,
		Gp: 0,
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "",
			},
		},
	}
	char.Erfahrungsschatz = models.Erfahrungsschatz{
		Value: 325,
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "",
			},
		},
	}
	char.Waffen = []equipment.Waffe{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Armbrust:schwer",
				},
			},
			Beschreibung:            "Eine Armbrust schwer zu spannen",
			Abwb:                    0,
			Anb:                     0,
			Schb:                    0,
			Anzahl:                  1,
			Gewicht:                 5,
			NameFuerSpezialisierung: "Armbrust:schwer",
			Wert:                    40,
			Magisch: models.Magisch{
				IstMagisch:  false,
				Abw:         0,
				Ausgebrannt: false,
			},
			BeinhaltetIn: "moam-container-47363",
		},
	}
	char.Behaeltnisse = []equipment.Behaeltniss{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Lederrucksack",
				},
			},
			Beschreibung: "f\u00fcr 25 kg",
			Wert:         4,
			Tragkraft:    25,
			Volumen:      25,
			Gewicht:      0.5,
			Magisch: models.Magisch{
				IstMagisch:  false,
				Abw:         0,
				Ausgebrannt: false,
			},
			BeinhaltetIn: "moam-container-47363",
		},
	}
	char.Transportmittel = []equipment.Transportation{
		{

			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Karren",
				},
			},
			Beschreibung: "für 500 kg",
			Wert:         40,
			Tragkraft:    500,
			Volumen:      250,
			Gewicht:      55.5,
			Magisch: models.Magisch{
				IstMagisch:  true,
				Abw:         30,
				Ausgebrannt: false,
			},
			//BeinhaltetIn: "moam-container-47363",

		},
	}
	char.Ausruestung = []equipment.Ausruestung{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Lederr\u00fcstung",
				},
			},
			Beschreibung: "",
			Wert:         30,
			Anzahl:       1,
			Gewicht:      13.0,
			Bonus:        0,
			Magisch: models.Magisch{
				IstMagisch:  false,
				Abw:         0,
				Ausgebrannt: false,
			},
			//BeinhaltetIn: "moam-container-47363",
		},
	}
	fileName := fmt.Sprintf("../testdata/%s", "Krampus.png")
	char.Image, _ = ReadImageAsBase64(fileName)

	return &char
}
func createFertigkeit(sel int) *models.Fertigkeit {

	liste := []models.Fertigkeit{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Hören",
				},
			},
			Beschreibung:    "xx",
			Fertigkeitswert: 6,
			Bonus:           0,
			Pp:              0,
			Bemerkung:       "",
		},
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Sprache",
				},
			},
			Beschreibung:    "Albisch",
			Fertigkeitswert: 8,
			Bonus:           0,
			Pp:              0,
			Bemerkung:       "",
		},
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Sprache",
				},
			},
			Beschreibung:    "Comentang",
			Fertigkeitswert: 12,
			Bonus:           0,
			Pp:              0,
			Bemerkung:       "",
		},
	}

	if sel > len(liste)-1 {
		sel = len(liste) - 1
	}
	return &liste[sel]
}
func createWaffenfertigkeit(sel int) *models.Waffenfertigkeit {

	liste := []models.Waffenfertigkeit{
		{
			Fertigkeit: models.Fertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Einhandschwerter",
					},
				},
				Beschreibung:    "",
				Fertigkeitswert: 5,
				Bonus:           0,
				Pp:              0,
				Bemerkung:       "",
			},
		},
		{
			Fertigkeit: models.Fertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Zweihandschlagwaffen",
					},
				},
				Beschreibung:    "",
				Fertigkeitswert: 5,
				Bonus:           0,
				Pp:              0,
				Bemerkung:       "",
			},
		},
		{
			Fertigkeit: models.Fertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Stichwaffen",
					},
				},
				Beschreibung:    "",
				Fertigkeitswert: 5,
				Bonus:           0,
				Pp:              0,
				Bemerkung:       "",
			},
		},
	}

	if sel > len(liste)-1 {
		sel = len(liste) - 1
	}
	return &liste[sel]
}
func createZauber(sel int) *models.Zauber {

	liste := []models.Zauber{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Feuerkugel",
				},
			},
			Beschreibung: "",
			Bonus:        0,
		},
	}

	if sel > len(liste)-1 {
		sel = len(liste) - 1
	}
	return &liste[sel]
}

func createWaffen(sel int) *equipment.Waffe {

	liste := []equipment.Waffe{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Armbrust:leicht",
				},
			},
			Beschreibung:            "Eine Armbrust leicht zu spannen",
			Abwb:                    0,
			Anb:                     0,
			Schb:                    0,
			Anzahl:                  1,
			Gewicht:                 5,
			NameFuerSpezialisierung: "Armbrust:leicht",
			Wert:                    40,
			Magisch: models.Magisch{
				IstMagisch:  false,
				Abw:         0,
				Ausgebrannt: false,
			},
			BeinhaltetIn: "",
		},
	}
	if sel > len(liste)-1 {
		sel = len(liste) - 1
	}
	return &liste[sel]
}
func createBehaeltniss(sel int) *equipment.Behaeltniss {
	liste := []equipment.Behaeltniss{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Lederrucksack",
				},
			},
			Beschreibung: "f\u00fcr 25 kg",
			Wert:         4,
			Tragkraft:    25,
			Volumen:      25,
			Gewicht:      0.5,
			Magisch: models.Magisch{
				IstMagisch:  false,
				Abw:         0,
				Ausgebrannt: false,
			},
			BeinhaltetIn: "moam-container-47363",
		},
	}

	if sel > len(liste)-1 {
		sel = len(liste) - 1
	}
	return &liste[sel]
}
func createTransportmittel(sel int) *equipment.Transportation {
	liste := []equipment.Transportation{
		{

			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Karren",
				},
			},
			Beschreibung: "für 500 kg",
			Wert:         40,
			Tragkraft:    500,
			Volumen:      250,
			Gewicht:      55.5,
			Magisch: models.Magisch{
				IstMagisch:  true,
				Abw:         30,
				Ausgebrannt: false,
			},
			//BeinhaltetIn: "moam-container-47363",

		},
	}

	if sel > len(liste)-1 {
		sel = len(liste) - 1
	}
	return &liste[sel]
}
func createAusruestung(sel int) *equipment.Ausruestung {
	liste := []equipment.Ausruestung{
		{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Lederr\u00fcstung",
				},
			},
			Beschreibung: "",
			Wert:         30,
			Anzahl:       1,
			Gewicht:      13.0,
			Bonus:        0,
			Magisch: models.Magisch{
				IstMagisch:  false,
				Abw:         0,
				Ausgebrannt: false,
			},
			//BeinhaltetIn: "moam-container-47363",
		},
	}

	if sel > len(liste)-1 {
		sel = len(liste) - 1
	}
	return &liste[sel]
}

/*
func setSpezialisierung(liste  *database.StringArray, pos int,name string) {
	liste[pos] =name
}
*/

func charTests(t *testing.T, char *character.Char) {
	assert.LessOrEqual(t, 0, int(char.ID))
	assert.Equal(t, "Harsk Hammerhuter, Zen", char.Name)
	assert.Equal(t, "Zwerg", char.Rasse)
	assert.Equal(t, "Krieger", char.Typ)
	assert.Equal(t, 39, char.Alter)
	assert.Equal(t, "er", char.Anrede)
	assert.Equal(t, 3, char.Grad)
	assert.Equal(t, 140, char.Groesse)
	assert.Equal(t, 82, char.Gewicht)
	assert.Equal(t, "Torkin", char.Glaube)
	assert.Equal(t, "rechts", char.Hand)
	assert.Equal(t, 17, char.Lp.Max)
	assert.Equal(t, 17, char.Lp.Value)
	assert.Equal(t, 31, char.Ap.Max)
	assert.Equal(t, 31, char.Ap.Value)
	assert.Equal(t, 18, char.B.Max)
	assert.Equal(t, 18, char.B.Value)
	assert.Equal(t, "blau", char.Merkmale.Augenfarbe)
	assert.Equal(t, "sandfarben", char.Merkmale.Haarfarbe)
	assert.Equal(t, "", char.Merkmale.Sonstige)
	for i := range char.Eigenschaften {
		switch char.Eigenschaften[i].Name {
		case "Au":
			assert.Equal(t, 74, char.Eigenschaften[i].Value)
		case "Gs":
			assert.Equal(t, 96, char.Eigenschaften[i].Value)
		case "Gw":
			assert.Equal(t, 70, char.Eigenschaften[i].Value)
		case "In":
			assert.Equal(t, 65, char.Eigenschaften[i].Value)
		case "Ko":
			assert.Equal(t, 85, char.Eigenschaften[i].Value)
		case "PA":
			assert.Equal(t, 75, char.Eigenschaften[i].Value)
		case "St":
			assert.Equal(t, 95, char.Eigenschaften[i].Value)
		case "Wk":
			assert.Equal(t, 71, char.Eigenschaften[i].Value)
		case "Zt":
			assert.Equal(t, 35, char.Eigenschaften[i].Value)

		}
	}
	assert.LessOrEqual(t, 3, len(char.Fertigkeiten))
	for i := range char.Fertigkeiten {
		assert.LessOrEqual(t, 0, int(char.Fertigkeiten[i].CharacterID))
		assert.Equal(t, char.ID, char.Fertigkeiten[i].CharacterID)
		switch char.Fertigkeiten[i].Name {
		case "Hören":
			assert.Equal(t, "Hören", char.Fertigkeiten[i].Name)
			assert.Equal(t, "xx", char.Fertigkeiten[i].Beschreibung)
			assert.Equal(t, 6, char.Fertigkeiten[i].Fertigkeitswert)
			assert.Equal(t, 0, char.Fertigkeiten[i].Bonus)
			assert.Equal(t, 0, char.Fertigkeiten[i].Pp)
			assert.Equal(t, "", char.Fertigkeiten[i].Bemerkung)
		case "Sprache":
			switch char.Fertigkeiten[i].Beschreibung {
			case "Albisch":
				assert.Equal(t, "Sprache", char.Fertigkeiten[i].Name)
				assert.Equal(t, "Albisch", char.Fertigkeiten[i].Beschreibung)
				assert.Equal(t, 8, char.Fertigkeiten[i].Fertigkeitswert)
				assert.Equal(t, 0, char.Fertigkeiten[i].Bonus)
				assert.Equal(t, 0, char.Fertigkeiten[i].Pp)
				assert.Equal(t, "", char.Fertigkeiten[i].Bemerkung)
			case "Comentang":
				assert.Equal(t, "Sprache", char.Fertigkeiten[i].Name)
				assert.Equal(t, "Comentang", char.Fertigkeiten[i].Beschreibung)
				assert.Equal(t, 12, char.Fertigkeiten[i].Fertigkeitswert)
				assert.Equal(t, 0, char.Fertigkeiten[i].Bonus)
				assert.Equal(t, 0, char.Fertigkeiten[i].Pp)
				assert.Equal(t, "", char.Fertigkeiten[i].Bemerkung)
			}
		}
	}
	assert.LessOrEqual(t, 1, len(char.Waffenfertigkeiten))
	for i := range char.Waffenfertigkeiten {
		assert.LessOrEqual(t, 0, int(char.Waffenfertigkeiten[i].CharacterID))
		assert.Equal(t, char.ID, char.Waffenfertigkeiten[i].CharacterID)
		switch char.Waffenfertigkeiten[i].Name {
		case "Armbrüste":
			assert.Equal(t, "Armbrüste", char.Waffenfertigkeiten[i].Name)
			assert.Equal(t, "", char.Waffenfertigkeiten[i].Beschreibung)
			assert.Equal(t, 8, char.Waffenfertigkeiten[i].Fertigkeitswert)
			assert.Equal(t, 0, char.Waffenfertigkeiten[i].Bonus)
			assert.Equal(t, 0, char.Waffenfertigkeiten[i].Pp)
			assert.Equal(t, "", char.Waffenfertigkeiten[i].Bemerkung)
		case "Einhandschlagwaffen":
			assert.Equal(t, "Einhandschlagwaffen", char.Waffenfertigkeiten[i].Name)
			assert.Equal(t, "", char.Waffenfertigkeiten[i].Beschreibung)
			assert.Equal(t, 8, char.Waffenfertigkeiten[i].Fertigkeitswert)
			assert.Equal(t, 0, char.Waffenfertigkeiten[i].Bonus)
			assert.Equal(t, 0, char.Waffenfertigkeiten[i].Pp)
			assert.Equal(t, "", char.Waffenfertigkeiten[i].Bemerkung)
		case "Schilde":
			assert.Equal(t, "Schilde", char.Waffenfertigkeiten[i].Name)
			assert.Equal(t, "", char.Waffenfertigkeiten[i].Beschreibung)
			assert.Equal(t, 3, char.Waffenfertigkeiten[i].Fertigkeitswert)
			assert.Equal(t, 0, char.Waffenfertigkeiten[i].Bonus)
			assert.Equal(t, 0, char.Waffenfertigkeiten[i].Pp)
			assert.Equal(t, "", char.Waffenfertigkeiten[i].Bemerkung)
		}
	}
	assert.LessOrEqual(t, 1, len(char.Zauber))
	for i := range char.Zauber {
		assert.LessOrEqual(t, 0, int(char.Zauber[i].CharacterID))
		assert.Equal(t, char.ID, char.Zauber[i].CharacterID)
		switch char.Zauber[i].Name {
		case "Angst":
			assert.Equal(t, "Angst", char.Zauber[i].Name)
			assert.Equal(t, "", char.Zauber[i].Beschreibung)
			assert.Equal(t, 0, char.Zauber[i].Bonus)
		}
	}
	assert.LessOrEqual(t, 3, len(char.Spezialisierung))
	assert.Equal(t, "Kriegshammer", char.Spezialisierung[0])
	assert.Equal(t, "Armbrust:schwer", char.Spezialisierung[1])
	assert.Equal(t, "Stielhammer", char.Spezialisierung[2])

	assert.Equal(t, 1, char.Bennies.Sg)
	assert.Equal(t, 0, char.Bennies.Gg)
	assert.Equal(t, 0, char.Bennies.Gp)
	assert.LessOrEqual(t, 0, int(char.Bennies.CharacterID))
	assert.Equal(t, 1, int(char.Bennies.ID))
	assert.Equal(t, 325, char.Erfahrungsschatz.Value)
	assert.LessOrEqual(t, 0, int(char.Erfahrungsschatz.CharacterID))
	assert.Equal(t, 1, int(char.Erfahrungsschatz.ID))

	assert.LessOrEqual(t, 1, len(char.Waffen))
	for i := range char.Waffen {
		assert.LessOrEqual(t, 0, int(char.Waffen[i].CharacterID))
		assert.Equal(t, char.ID, char.Waffen[i].CharacterID)
		switch char.Waffen[i].Name {
		case "Armbrust:schwer":
			assert.Equal(t, "Armbrust:schwer", char.Waffen[i].Name)
			assert.Equal(t, "Eine Armbrust schwer zu spannen", char.Waffen[i].Beschreibung)
			assert.Equal(t, 0, char.Waffen[i].Abwb)
			assert.Equal(t, 0, char.Waffen[i].Anb)
			assert.Equal(t, 0, char.Waffen[i].Schb)
			assert.Equal(t, 1, char.Waffen[i].Anzahl)
			assert.Equal(t, 5.0, char.Waffen[i].Gewicht)
			assert.Equal(t, "Armbrust:schwer", char.Waffen[i].NameFuerSpezialisierung)
			assert.Equal(t, "moam-container-47363", char.Waffen[i].BeinhaltetIn)
			assert.Equal(t, false, char.Waffen[i].IstMagisch)
			assert.Equal(t, 0, char.Waffen[i].Abw)
			assert.Equal(t, false, char.Waffen[i].Ausgebrannt)
		}
	}
	assert.LessOrEqual(t, 1, len(char.Behaeltnisse))
	for i := range char.Behaeltnisse {
		assert.LessOrEqual(t, 0, int(char.Behaeltnisse[i].CharacterID))
		assert.Equal(t, char.ID, char.Behaeltnisse[i].CharacterID)
		switch char.Behaeltnisse[i].Name {
		case "Lederrucksack":
			assert.Equal(t, "Lederrucksack", char.Behaeltnisse[i].Name)
			assert.Equal(t, "für 25 kg", char.Behaeltnisse[i].Beschreibung)
			assert.Equal(t, 4.0, char.Behaeltnisse[i].Wert)
			assert.Equal(t, 25.0, char.Behaeltnisse[i].Tragkraft)
			assert.Equal(t, 25.0, char.Behaeltnisse[i].Volumen)
			assert.Equal(t, 0.5, char.Behaeltnisse[i].Gewicht)
			assert.Equal(t, false, char.Behaeltnisse[i].IstMagisch)
			assert.Equal(t, 0, char.Behaeltnisse[i].Abw)
			assert.Equal(t, false, char.Behaeltnisse[i].Ausgebrannt)
		}
	}
	assert.LessOrEqual(t, 1, len(char.Transportmittel))
	for i := range char.Transportmittel {
		assert.LessOrEqual(t, 0, int(char.Transportmittel[i].CharacterID))
		assert.Equal(t, char.ID, char.Transportmittel[i].CharacterID)
		switch char.Transportmittel[i].Name {
		case "Karren":
			assert.LessOrEqual(t, 0, int(char.Transportmittel[i].CharacterID))
			assert.Equal(t, "Karren", char.Transportmittel[i].Name)
			assert.Equal(t, "für 500 kg", char.Transportmittel[i].Beschreibung)
			assert.Equal(t, 40.0, char.Transportmittel[i].Wert)
			assert.Equal(t, 500.0, char.Transportmittel[i].Tragkraft)
			assert.Equal(t, 250.0, char.Transportmittel[i].Volumen)
			assert.Equal(t, 55.5, char.Transportmittel[i].Gewicht)
			assert.Equal(t, true, char.Transportmittel[i].IstMagisch)
			assert.Equal(t, 30, char.Transportmittel[i].Abw)
			assert.Equal(t, false, char.Transportmittel[i].Ausgebrannt)
		}
	}
	assert.LessOrEqual(t, 1, len(char.Ausruestung))
	for i := range char.Ausruestung {
		assert.LessOrEqual(t, 0, int(char.Ausruestung[i].CharacterID))
		assert.Equal(t, char.ID, char.Ausruestung[i].CharacterID)
		switch char.Ausruestung[i].Name {
		case "Lederrüstung":
			assert.Equal(t, "Lederrüstung", char.Ausruestung[i].Name)
			assert.Equal(t, "", char.Ausruestung[i].Beschreibung)
			assert.Equal(t, 30.0, char.Ausruestung[i].Wert)
			assert.Equal(t, 13.0, char.Ausruestung[i].Gewicht)
			assert.Equal(t, 1, char.Ausruestung[i].Anzahl)
			assert.Equal(t, false, char.Ausruestung[i].IstMagisch)
			assert.Equal(t, 0, char.Ausruestung[i].Abw)
			assert.Equal(t, false, char.Ausruestung[i].Ausgebrannt)
		}
	}

	assert.Contains(t, char.Image, "data:image/png;base64,")

}

func TestCreateChar(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Character()
	database.DB = testDB // Assign test DB to global DB
	char := createChar()
	//char.Name = "Harsk Hammerhuter, Zen2"
	err := char.First(char.Name)
	//assert.Error(t, err, "expected error character not found")
	if err != nil && err.Error() == "record not found" {
		err = char.Create()
		assert.NoError(t, err, "expected no error creating the char")
	}
	char2 := character.Char{}
	char2.Name = "Harsk Hammerhuter, Zen"
	err = char2.First(char2.Name)
	assert.NoError(t, err, "expected no error creating the char")
	charTests(t, &char2)
}

func TestReadChar(t *testing.T) {
	// Setup test database
	//testDB := initTestDB4Character()
	//database.DB = testDB // Assign test DB to global DB

	TestCreateChar(t)
	char := character.Char{}
	char.Name = "Harsk Hammerhuter, Zen"
	err := char.First(char.Name)
	assert.NoError(t, err, "expected NO error character not found")
	if err != nil {
		//try again
		err := char.First(char.Name)
		assert.NoError(t, err, "expected no error creating the char")

	}
	charTests(t, &char)
}

func TestAddAusrüstung(t *testing.T) {
	// Setup test database
	//testDB := initTestDB4Character()
	//database.DB = testDB // Assign test DB to global DB
	/*
		TestCreateChar(t)
		char := character.Char{}
		char.Name = "Harsk Hammerhuter, Zen"
		err := char.First(char.Name)
		assert.NoError(t, err, "expexted no Errorselecting record from DataBase")
		item := createAusruestung(0)
		mcr := character.Char{}
		err = mcr.AddAusruestung(item.Name)
		assert.NoError(t, err, "No error expected when adding Ausrüstung to char")
	*/

}

func TestImportVTT2Char(t *testing.T) {
	// Setup test database
	//testDB := SetupTestDB()
	//DB = testDB // Assign test DB to global DB
	/*
		// loading file to Modell
		fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
		assert.Equal(t, "../testdata/VTT_Import1.json", fileName)
		fileContent, err := os.ReadFile(fileName)
		assert.NoError(t, err, "Expected no error when reading file "+fileName)
		character := models.ImCharacterImport{}
		err = json.Unmarshal(fileContent, &character)
		assert.NoError(t, err, "Expected no error when Unmarshal filecontent")

		assert.Equal(t, "Harsk Hammerhuter, Zen", character.Name)
		assert.Equal(t, "Zwerg", character.Rasse)
		assert.Equal(t, "Hören", character.Fertigkeiten[0].Name)
		assert.Equal(t, 1, len(character.Zauber))
		assert.Equal(t, 17, character.Lp.Value)
		assert.Equal(t, 96, character.Eigenschaften.Gs)
		assert.Equal(t, 74, character.Eigenschaften.Au)
		assert.Equal(t, 1, len(character.Ausruestung))
		assert.Equal(t, "Lederrüstung", character.Ausruestung[0].Name)
		assert.Equal(t, "blau", character.Merkmale.Augenfarbe)
		assert.Equal(t, "Lederrucksack", character.Behaeltnisse[0].Name)
		assert.Equal(t, "Armbrust:schwer", character.Waffen[0].Name)
		assert.Equal(t, 31, character.Ap.Value)
		assert.Equal(t, "Armbrüste", character.Waffenfertigkeiten[0].Name)
		assert.Equal(t, 3, len(character.Spezialisierung))
		assert.Equal(t, "Kriegshammer", character.Spezialisierung[0])
		assert.Equal(t, "Armbrust:schwer", character.Spezialisierung[1])
		//fmt.Println(character)
	*/
}
