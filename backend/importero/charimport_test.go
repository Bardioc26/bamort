package importero

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportVTT2Char(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	char, err := ImportVTTJSON(fileName, 1)
	assert.NoError(t, err, "expected no error when saving imported Char")
	var chr2 models.Char
	chr2.First(char.Name)
	assert.GreaterOrEqual(t, char.ID, chr2.ID)
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
		assert.Equal(t, 17, models.Lp.Value)
		assert.Equal(t, 96, models.Eigenschaften.Gs)
		assert.Equal(t, 74, models.Eigenschaften.Au)
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
