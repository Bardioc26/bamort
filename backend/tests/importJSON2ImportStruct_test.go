package tests

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/importer"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initTestDB4Import() *gorm.DB {
	db := SetupTestDB()
	db.AutoMigrate(
		&gsmaster.Skill{},
	)
	return db
}

func readImportChar(fileName string) (*importer.CharacterImport, error) {
	// loading file to Modell
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	character := importer.CharacterImport{}
	err = json.Unmarshal(fileContent, &character)
	return &character, err
}

func testChar(t *testing.T, object *importer.CharacterImport) {
	assert.Equal(t, "moam-character-41421", object.ID)
	assert.Equal(t, "Harsk Hammerhuter, Zen", object.Name)
	assert.Equal(t, "Zwerg", object.Rasse)
	assert.Equal(t, "Krieger", object.Typ)
	assert.Equal(t, 39, object.Alter)
	assert.Equal(t, "er", object.Anrede)
	assert.Equal(t, 3, object.Grad)
	assert.Equal(t, 140, object.Groesse)
	assert.Equal(t, 82, object.Gewicht)
	assert.Equal(t, "Torkin", object.Glaube)
	assert.Equal(t, "rechts", object.Hand)
	assert.Equal(t, 17, object.Lp.Max)
	assert.Equal(t, 17, object.Lp.Value)
	assert.Equal(t, 31, object.Ap.Max)
	assert.Equal(t, 31, object.Ap.Value)
	assert.Equal(t, 18, object.B.Max)
	assert.Equal(t, 0, object.B.Value)
	assert.Equal(t, 74, object.Eigenschaften.Au)
	assert.Equal(t, 96, object.Eigenschaften.Gs)
	assert.Equal(t, 70, object.Eigenschaften.Gw)
	assert.Equal(t, 65, object.Eigenschaften.In)
	assert.Equal(t, 85, object.Eigenschaften.Ko)
	assert.Equal(t, 75, object.Eigenschaften.Pa)
	assert.Equal(t, 95, object.Eigenschaften.St)
	assert.Equal(t, 71, object.Eigenschaften.Wk)
	assert.Equal(t, 35, object.Eigenschaften.Zt)
	assert.Equal(t, "blau", object.Merkmale.Augenfarbe)
	assert.Equal(t, "sandfarben", object.Merkmale.Haarfarbe)
	assert.Equal(t, "", object.Merkmale.Sonstige)
	assert.Equal(t, 0, object.Bennies.Gg)
	assert.Equal(t, 1, object.Bennies.Sg)
	assert.Equal(t, 0, object.Bennies.Gp)
	assert.Equal(t, "breit", object.Gestalt.Breite)
	assert.Equal(t, "klein", object.Gestalt.Groesse)
	assert.Equal(t, 325, object.Erfahrungsschatz.Value)
	assert.Equal(t, 3, len(object.Spezialisierung))
	assert.Equal(t, "Kriegshammer", object.Spezialisierung[0])
	assert.Equal(t, "Armbrust:schwer", object.Spezialisierung[1])
	assert.Equal(t, "Stielhammer", object.Spezialisierung[2])
	assert.Contains(t, object.Image, "data:image;base64,")

}

func testSkill(t *testing.T, objects []importer.Fertigkeit) {
	assert.Equal(t, 19, len(objects))
	i := 0
	assert.Equal(t, "moam-ability-horen", objects[i].ID)
	assert.Equal(t, "Hören", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 6, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 99", objects[i].Quelle)
	i = 6
	assert.Equal(t, "moam-ability-759918", objects[i].ID)
	assert.Equal(t, "Athletik", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 9, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 104", objects[i].Quelle)
	i = 16
	assert.Equal(t, "moam-ability-759920", objects[i].ID)
	assert.Equal(t, "Sprache", objects[i].Name)
	assert.Equal(t, "Albisch", objects[i].Beschreibung)
	assert.Equal(t, 8, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 127", objects[i].Quelle)

}

func testWaeponSkill(t *testing.T, objects []importer.Waffenfertigkeit) {
	assert.Equal(t, 8, len(objects))
	i := 0
	assert.Equal(t, "moam-ability-759916", objects[i].ID)
	assert.Equal(t, "Armbrüste", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 8, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 144", objects[i].Quelle)
	i = 2
	assert.Equal(t, "moam-ability-759912", objects[i].ID)
	assert.Equal(t, "Schilde", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 3, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 145", objects[i].Quelle)

}

func testSpell(t *testing.T, objects []importer.Zauber) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-spell-134630", objects[i].ID)
	assert.Equal(t, "Angst", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, "ARK5 63", objects[i].Quelle)
}

func testWaepon(t *testing.T, objects []importer.Waffe) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-weapon-126819", objects[i].ID)
	assert.Equal(t, "Armbrust:schwer", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 0, objects[i].Abwb)
	assert.Equal(t, 0, objects[i].Anb)
	assert.Equal(t, 1, objects[i].Anzahl)
	assert.Equal(t, "moam-container-47363", objects[i].BeinhaltetIn)
	assert.Equal(t, 5.0, objects[i].Gewicht)
	assert.Equal(t, false, objects[i].Magisch.IstMagisch)
	assert.Equal(t, 0, objects[i].Magisch.Abw)
	assert.Equal(t, false, objects[i].Magisch.Ausgebrannt)
	assert.Equal(t, "Armbrust:schwer", objects[i].NameFuerSpezialisierung)
	assert.Equal(t, 0, objects[i].Schb)
	assert.Equal(t, 40.0, objects[i].Wert)

}

func testEquipment(t *testing.T, objects []importer.Ausruestung) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-armor-48616", objects[i].ID)
	assert.Equal(t, "Lederrüstung", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 1, objects[i].Anzahl)
	assert.Equal(t, "", objects[i].BeinhaltetIn)
	assert.Equal(t, 13.0, objects[i].Gewicht)
	assert.Equal(t, "", objects[i].BeinhaltetIn)
	assert.Equal(t, false, objects[i].Magisch.IstMagisch)
	assert.Equal(t, 0, objects[i].Magisch.Abw)
	assert.Equal(t, false, objects[i].Magisch.Ausgebrannt)
	assert.Equal(t, 30.0, objects[i].Wert)
	assert.Equal(t, 0, objects[i].Bonus)
}

func testContainer(t *testing.T, objects []importer.Behaeltniss) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-container-47363", objects[i].ID)
	assert.Equal(t, "Lederrucksack", objects[i].Name)
	assert.Equal(t, "für 25 kg", objects[i].Beschreibung)
	assert.Equal(t, 4.0, objects[i].Wert)
	assert.Equal(t, 0.50, objects[i].Gewicht)
	assert.Equal(t, 25.0, objects[i].Volumen)
	assert.Equal(t, 25.0, objects[i].Tragkraft)
	assert.Empty(t, "", objects[i].BeinhaltetIn) //Value in json is null
	assert.Equal(t, false, objects[i].Magisch.IstMagisch)
	assert.Equal(t, 0, objects[i].Magisch.Abw)
	assert.Equal(t, false, objects[i].Magisch.Ausgebrannt)

}

func testTransportation(t *testing.T, objects []importer.Transportation) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-container-47000", objects[i].ID)
	assert.Equal(t, "Karren", objects[i].Name)
	assert.Equal(t, "für 250 kg", objects[i].Beschreibung)
	assert.Equal(t, 14.0, objects[i].Wert)
	assert.Equal(t, 40, objects[i].Gewicht)
	assert.Equal(t, 250.0, objects[i].Tragkraft)
	assert.Empty(t, "", objects[i].BeinhaltetIn) //Value in json is null
	assert.Equal(t, false, objects[i].Magisch.IstMagisch)
	assert.Equal(t, 0, objects[i].Magisch.Abw)
	assert.Equal(t, false, objects[i].Magisch.Ausgebrannt)

}

func TestImportVTTStructure(t *testing.T) {
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	testChar(t, character)
	testSkill(t, character.Fertigkeiten)
	testWaeponSkill(t, character.Waffenfertigkeiten)
	testSpell(t, character.Zauber)
	testEquipment(t, character.Ausruestung)
	testWaepon(t, character.Waffen)
	testContainer(t, character.Behaeltnisse)
	testTransportation(t, character.Transportmittel)
	//fmt.Println(character)
}

func TestImportSkill2GSMaster(t *testing.T) {
	testDB := initTestDB4Import()
	database.DB = testDB // Assign test DB to global DB
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	skill, erro := importer.TransformImportFertigkeit2GSDMaster(&character.Fertigkeiten[0])

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(skill.ID), 1)
	assert.Equal(t, "Hören", skill.Name)
	assert.Equal(t, "", skill.Beschreibung)
	assert.Equal(t, 6, skill.Initialkeitswert)
	assert.Equal(t, "check", skill.Bonuseigenschaft)
	assert.Equal(t, "KOD5 99", skill.Quelle)
	assert.Equal(t, false, skill.Improvable)
	assert.Equal(t, "midgard", skill.System)
	//}
	skill2 := gsmaster.Skill{}
	erro = skill2.First("Hören")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(skill.ID))

	skill3 := gsmaster.Skill{}
	erro = skill3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Hören", skill3.Name)

	assert.Equal(t, skill2.ID, skill3.ID)
	assert.Equal(t, skill2.Name, skill3.Name)
	assert.Equal(t, skill2.Beschreibung, skill3.Beschreibung)
	assert.Equal(t, skill2.Initialkeitswert, skill3.Initialkeitswert)
	assert.Equal(t, skill2.Bonuseigenschaft, skill3.Bonuseigenschaft)
	assert.Equal(t, skill2.Quelle, skill3.Quelle)
	assert.Equal(t, skill2.Improvable, skill3.Improvable)
	assert.Equal(t, skill2.System, skill3.System)
}

/*
func TestImportFertigkeitenStammdatenSingle(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Import()
	database.DB = testDB // Assign test DB to global DB

	// loading file to Modell
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	assert.Equal(t, "../testdata/VTT_Import1.json", fileName)
	fileContent, err := os.ReadFile(fileName)
	assert.NoError(t, err, "Expected no error when reading file "+fileName)
	character := importer.CharacterImport{}
	err = json.Unmarshal(fileContent, &character)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")

	//checke Fertigkeit auf vorhandensein in den Stammdaten
	fertigkeit := character.Fertigkeiten[1]
	stammF, err := importer.CheckSkill(&fertigkeit, false)
	assert.Error(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	if stammF == nil && err != nil {
		stammF, err = importer.CheckSkill(&fertigkeit, true)
	}
	assert.NoError(t, err, "Expected to finds the Fertigkeit Stammdaten in the database")
	assert.Equal(t, fertigkeit.Name, stammF.Name)
	assert.Equal(t, fertigkeit.Beschreibung, stammF.Beschreibung)
	assert.Equal(t, fertigkeit.Quelle, stammF.Quelle)
	assert.Equal(t, 5, stammF.Initialkeitswert)
	assert.Equal(t, "keine", stammF.Bonuseigenschaft)
	assert.Equal(t, "midgard", stammF.System)
	assert.Equal(t, 1, int(stammF.ID))

	// und noch mal
	//checke Fertigkeit auf vorhandensein in den Stammdaten
	//fertigkeit := character.Fertigkeiten[1]

	stammF, err = importer.CheckSkill(&fertigkeit, false)
	assert.NoError(t, err, "expexted no Error exist in Fertigkeit Stammdaten")
	if stammF == nil && err != nil {
		stammF, err = importer.CheckSkill(&fertigkeit, true)
	}
	assert.NoError(t, err, "Expected to finds the Fertigkeit Stammdaten in the database")
	assert.Equal(t, fertigkeit.Name, stammF.Name)
	assert.Equal(t, fertigkeit.Beschreibung, stammF.Beschreibung)
	assert.Equal(t, fertigkeit.Quelle, stammF.Quelle)
	assert.Equal(t, 5, stammF.Initialkeitswert)
	assert.Equal(t, "keine", stammF.Bonuseigenschaft)
	assert.Equal(t, "midgard", stammF.System)
	assert.Equal(t, 1, int(stammF.ID))
}

func TestImportFertigkeitenStammdatenMulti(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Import()
	database.DB = testDB // Assign test DB to global DB

	// loading file to Modell
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	assert.Equal(t, "../testdata/VTT_Import1.json", fileName)
	fileContent, err := os.ReadFile(fileName)
	assert.NoError(t, err, "Expected no error when reading file "+fileName)
	character := importer.CharacterImport{}
	err = json.Unmarshal(fileContent, &character)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")

	//for index, fertigkeit := range character.Fertigkeiten {
	for _, fertigkeit := range character.Fertigkeiten {
		fmt.Println(fertigkeit.Name)
		stammF, err := importer.CheckSkill(&fertigkeit, true)
		assert.NoError(t, err, "Expected to finds the Fertigkeit Stammdaten in the database")
		assert.Equal(t, fertigkeit.Name, stammF.Name, "Name should be equal")
		if fertigkeit.Name != "Sprache" {
			assert.Equal(t, fertigkeit.Beschreibung, stammF.Beschreibung, "Beschreibung should be equal")
		} else {
			assert.Equal(t, "", stammF.Beschreibung, "Beschreibung should be equal")
		}
		assert.Equal(t, fertigkeit.Quelle, stammF.Quelle, "Quelle should be equal")
		//assert.Equal(t, 5, stammF.Initialkeitswert, "Initialkeitswert should be equal")
		assert.Equal(t, "keine", stammF.Bonuseigenschaft, "Bonuseigenschaft should be equal")
		assert.Equal(t, "midgard", stammF.System, "System should be equal")
		//assert.NotEmpty(t, index+1, int(stammF.ID), "ID should be equal")
	}
}

func TestImportZauberStammdatenSingle(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Import()
	database.DB = testDB // Assign test DB to global DB

	// loading file to Modell
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	assert.Equal(t, "../testdata/VTT_Import1.json", fileName)
	fileContent, err := os.ReadFile(fileName)
	assert.NoError(t, err, "Expected no error when reading file "+fileName)
	character := importer.CharacterImport{}
	err = json.Unmarshal(fileContent, &character)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")

	//checke zauber auf vorhandensein in den Stammdaten
	zauber := character.Zauber[0]
	stammF, err := importer.CheckSpell(&zauber, false)
	assert.Error(t, err, "expexted Error does not exist in zauber Stammdaten")
	if stammF == nil && err != nil {
		stammF, err = importer.CheckSpell(&zauber, true)
	}
	assert.NoError(t, err, "Expected to finds the zauber Stammdaten in the database")
	assert.Equal(t, zauber.Name, stammF.Name)
	assert.Equal(t, zauber.Beschreibung, stammF.Beschreibung)
	assert.Equal(t, zauber.Quelle, stammF.Quelle)
	assert.Equal(t, 1, stammF.AP)
	assert.Equal(t, 1, stammF.Stufe)
	assert.Equal(t, "Zauberer", stammF.Wirkungsziel)
	assert.Equal(t, 15, stammF.Reichweite)
	assert.Equal(t, "midgard", stammF.System)
	assert.Equal(t, 1, int(stammF.ID))

	// und noch mal
	//checke zauber auf vorhandensein in den Stammdaten
	//zauber := character.zauberen[1]
	stammF, err = importer.CheckSpell(&zauber, false)
	assert.NoError(t, err, "expexted no Error exist in zauber Stammdaten")
	if stammF == nil && err != nil {
		stammF, err = importer.CheckSpell(&zauber, true)
	}
	assert.NoError(t, err, "Expected to finds the zauber Stammdaten in the database")
	assert.Equal(t, zauber.Name, stammF.Name)
	assert.Equal(t, zauber.Beschreibung, stammF.Beschreibung)
	assert.Equal(t, zauber.Quelle, stammF.Quelle)
	assert.Equal(t, 1, stammF.AP)
	assert.Equal(t, 1, stammF.Stufe)
	assert.Equal(t, "Zauberer", stammF.Wirkungsziel)
	assert.Equal(t, 15, stammF.Reichweite)
	assert.Equal(t, "midgard", stammF.System)
	assert.Equal(t, 1, int(stammF.ID))
}
*/
