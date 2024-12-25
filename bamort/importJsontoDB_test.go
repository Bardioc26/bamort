package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestImportVTT(t *testing.T) {
	// Setup test database
	//testDB := SetupTestDB()
	//DB = testDB // Assign test DB to global DB

	// loading file to Modell
	fileName := fmt.Sprintf("./uploads/%s", "test.json")
	fileContent, err := os.ReadFile(fileName)
	assert.NoError(t, err, "Expected no error when reading file "+fileName)
	character := CharacterImport{}
	err = json.Unmarshal(fileContent, &character)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")

	assert.Equal(t, "Harsk Hammerhuter, Zen", character.Name)
	assert.Equal(t, "Zwerg", character.Rasse)
	assert.Equal(t, "Hören", character.Fertigkeiten[0].Name)
	assert.Equal(t, 0, len(character.Zauber))
	assert.Equal(t, 17, character.Lp.Value)
	assert.Equal(t, 96, character.Eigenschaften.Gs)
	assert.Equal(t, 74, character.Eigenschaften.Au)
	assert.Equal(t, 21, len(character.Ausruestung))
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

func TestImportFertigkeitenStammdaten(t *testing.T) {
	// Setup test database
	testDB := SetupTestDB()
	DB = testDB // Assign test DB to global DB
	// loading file to Modell
	fileName := fmt.Sprintf("./uploads/%s", "test.json")
	fileContent, err := os.ReadFile(fileName)
	assert.NoError(t, err, "Expected no error when reading file "+fileName)
	character := CharacterImport{}
	err = json.Unmarshal(fileContent, &character)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")

	//fertigkeit := character.Fertigkeiten[1]
	for _, fertigkeit := range character.Fertigkeiten {
		//fmt.Printf("Name: %s, Beschreibung: %s, Wert: %d\n",
		//	fertigkeit.Name, fertigkeit.Beschreibung, fertigkeit.Fertigkeitswert)
		//fmt.Println(fertigkeit)
		stammF := StammFertigkeit{}
		if strings.HasPrefix(fertigkeit.ImportID, "moam") {
			err = DB.First(&stammF, "system=? AND name = ?", "midgard", fertigkeit.Name).Error
			assert.Error(t, err, "Expected not to find the Fertigkeit Stammdaten in the database")
			stammF.System = "midgard"
			stammF.Name = fertigkeit.Name
			stammF.Beschreibung = fertigkeit.Beschreibung
			if fertigkeit.Fertigkeitswert < 12 {
				stammF.Initialkeitswert = 5
			} else {
				stammF.Initialkeitswert = 12
			}
			stammF.Bonuseigenschaft = "keine"
			stammF.Quelle = fertigkeit.Quelle
			//fmt.Println(stammF)
			err = DB.Transaction(func(tx *gorm.DB) error {
				// Save the main character record
				if err := tx.Create(&stammF).Error; err != nil {
					return fmt.Errorf("failed to save Fertigkeit Stammdaten: %w", err)
				}
				return nil
			})
			assert.NoError(t, err, "Expected no error saving Fertigkeit Stammdaten in the database")
		}
		err = DB.First(&stammF, "system=? AND name = ?", "midgard", fertigkeit.Name).Error
		fmt.Println(stammF)
		assert.NoError(t, err, "Expected to find the Fertigkeit Stammdaten in the database")

	}

}
