package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportJSONToDB(t *testing.T) {
	// Setup test database
	//testDB := SetupTestDB()
	//DB = testDB // Assign test DB to global DB

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

	/*
	 */
}
