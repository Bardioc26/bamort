package tests

import (
	"bamort/importer"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
func initTestDB4Import() *gorm.DB {
	db := SetupTestDB()
	db.AutoMigrate(
			&gsmaster.Skill{}, //needed for gsmaster.CheckFertigkeit
			&gsmaster.Spell{}, //needed for gsmaster.CheckZauber

	)
	return db
}
*/

func readImportChat(fileName string) (*importer.CharacterImport, error) {
	// loading file to Modell
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	character := importer.CharacterImport{}
	err = json.Unmarshal(fileContent, &character)
	return &character, err
}

func TestImportVTTStructure(t *testing.T) {
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChat(fileName)
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
