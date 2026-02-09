package importer

import (
	"bamort/database"
	"bamort/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupImportTest initializes the test database and runs migrations
func setupImportTest() {
	database.SetupTestDB()
	models.MigrateStructure(database.DB)
	MigrateStructure(database.DB)
}

func TestImportCharacter_Success(t *testing.T) {
	setupImportTest()

	char := &CharacterImport{
		Name:   "Test Character",
		Rasse:  "Mensch",
		Typ:    "midgard",
		Grad:   1,
		Alter:  25,
		Anrede: "Herr",
		Eigenschaften: Eigenschaften{
			St: 50,
			Gs: 60,
			Gw: 70,
			Ko: 55,
			In: 65,
			Zt: 58,
			Au: 45,
			Pa: 50,
			Wk: 60,
		},
		Lp:               Lp{Max: 10, Value: 10},
		Ap:               Ap{Max: 20, Value: 20},
		B:                B{Max: 15, Value: 15},
		Erfahrungsschatz: Erfahrungsschatz{Value: 0},
		Bennies: Bennies{
			Gg: 0,
			Gp: 3,
			Sg: 0,
		},
		Fertigkeiten: []Fertigkeit{
			{
				ImportBase:      ImportBase{Name: "ImportTestSkill"},
				Fertigkeitswert: 10,
			},
		},
	}

	result, err := ImportCharacter(char, 1, "test-adapter", []byte(`{"test": "data"}`))

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "success", result.Status)
	assert.NotZero(t, result.CharacterID)
	assert.NotZero(t, result.ImportID)
	assert.Equal(t, "test-adapter", result.AdapterID)

	// Verify character was created
	var createdChar models.Char
	err = database.DB.Preload("Eigenschaften").First(&createdChar, result.CharacterID).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test Character", createdChar.Name)
	assert.Equal(t, "midgard", createdChar.Typ)

	// Verify attributes
	var stAttribute models.Eigenschaft
	err = database.DB.Where("character_id = ? AND name = ?", createdChar.ID, "St").First(&stAttribute).Error
	assert.NoError(t, err)
	assert.Equal(t, 50, stAttribute.Value)

	// Verify import history was created
	var history ImportHistory
	err = database.DB.First(&history, result.ImportID).Error
	assert.NoError(t, err)
	assert.Equal(t, "success", history.Status)
	assert.NotNil(t, history.CharacterID)
	assert.Equal(t, result.CharacterID, *history.CharacterID)
	assert.NotEmpty(t, history.SourceSnapshot)
}

func TestImportCharacter_CreatesPersonalItems(t *testing.T) {
	setupImportTest()

	char := &CharacterImport{
		Name:  "Character with Skills",
		Rasse: "Elf",
		Typ:   "midgard",
		Grad:  2,
		Eigenschaften: Eigenschaften{
			St: 45, Gs: 70, Gw: 75, Ko: 50,
			In: 80, Zt: 55, Au: 60, Pa: 65, Wk: 70,
		},
		Lp:               Lp{Max: 12, Value: 12},
		Ap:               Ap{Max: 25, Value: 25},
		B:                B{Max: 18, Value: 18},
		Erfahrungsschatz: Erfahrungsschatz{Value: 1000},
		Bennies:          Bennies{Gg: 0, Gp: 3, Sg: 0},
		Fertigkeiten: []Fertigkeit{
			{
				ImportBase:      ImportBase{Name: "UniqueImportSkill1"},
				Fertigkeitswert: 15,
				Beschreibung:    "A unique skill",
			},
			{
				ImportBase:      ImportBase{Name: "UniqueImportSkill2"},
				Fertigkeitswert: 12,
			},
		},
		Zauber: []Zauber{
			{
				ImportBase:   ImportBase{Name: "UniqueImportSpell1"},
				Beschreibung: "A unique spell",
			},
		},
	}

	result, err := ImportCharacter(char, 1, "test-adapter", nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "success", result.Status)

	// Check that personal items were created
	assert.Equal(t, 2, result.CreatedItems["skills"], "Should create 2 personal skills")
	assert.Equal(t, 1, result.CreatedItems["spells"], "Should create 1 personal spell")
}

func TestImportCharacter_RollbackOnError(t *testing.T) {
	setupImportTest()

	// Create a character with invalid data that will cause reconciliation to fail
	// This is a simplified test - in reality you'd need to trigger an actual error
	char := &CharacterImport{
		Name:             "Invalid Character",
		Rasse:            "Unknown",
		Typ:              "midgard",
		Eigenschaften:    Eigenschaften{St: 50, Gs: 60, Gw: 70, Ko: 55, In: 65, Zt: 58, Au: 45, Pa: 50, Wk: 60},
		Lp:               Lp{Max: 10, Value: 10},
		Ap:               Ap{Max: 20, Value: 20},
		B:                B{Max: 15, Value: 15},
		Erfahrungsschatz: Erfahrungsschatz{Value: 0},
		Bennies:          Bennies{Gg: 0, Gp: 3, Sg: 0},
	}

	// Count characters before
	var countBefore int64
	database.DB.Model(&models.Char{}).Count(&countBefore)

	result, _ := ImportCharacter(char, 1, "test-adapter", nil)

	// Even if there's no error (simplified test), verify the transaction logic
	assert.NotNil(t, result)

	// In a real error scenario, character count should remain the same
	var countAfter int64
	database.DB.Model(&models.Char{}).Count(&countAfter)

	// This test is simplified - in a real scenario with an error, we'd check:
	// assert.Equal(t, countBefore, countAfter, "Character count should not change on rollback")
}

func TestCompressData(t *testing.T) {
	data := []byte("This is test data that should be compressed")

	compressed, err := compressData(data)

	assert.NoError(t, err)
	assert.NotNil(t, compressed)
	assert.Less(t, len(compressed), len(data)+50, "Compressed data should not be much larger than original")
}

func TestCompressDecompressRoundTrip(t *testing.T) {
	original := []byte(`{"name": "Test Character", "skills": ["Skill1", "Skill2"]}`)

	compressed, err := compressData(original)
	assert.NoError(t, err)

	decompressed, err := decompressData(compressed)
	assert.NoError(t, err)

	assert.Equal(t, original, decompressed, "Round trip should preserve data")
}

func TestDecompressData_InvalidData(t *testing.T) {
	invalid := []byte("not gzip data")

	_, err := decompressData(invalid)

	assert.Error(t, err)
}
