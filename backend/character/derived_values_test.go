package character

import (
	"bamort/database"
	"bamort/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDerivedValuesStorage(t *testing.T) {
	// Setup test environment
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original != "" {
			os.Setenv("ENVIRONMENT", original)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})

	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	err := models.MigrateStructure()
	assert.NoError(t, err)

	t.Run("Derived values are stored in database", func(t *testing.T) {
		// Create a character with derived values
		char := models.Char{
			BamortBase: models.BamortBase{
				Name: "Test Character",
			},
			UserID:           1,
			Rasse:            "Mensch",
			Typ:              "Krieger",
			Grad:             1,
			ResistenzKoerper: 12,
			ResistenzGeist:   11,
			Abwehr:           13,
			Zaubern:          10,
			Raufen:           14,
		}

		// Save character
		err := char.Create()
		assert.NoError(t, err, "Should create character successfully")

		// Reload character from database
		var loadedChar models.Char
		err = database.DB.First(&loadedChar, char.ID).Error
		assert.NoError(t, err, "Should load character from database")

		// Verify derived values are persisted
		assert.Equal(t, 12, loadedChar.ResistenzKoerper, "ResistenzKoerper should be persisted")
		assert.Equal(t, 11, loadedChar.ResistenzGeist, "ResistenzGeist should be persisted")
		assert.Equal(t, 13, loadedChar.Abwehr, "Abwehr should be persisted")
		assert.Equal(t, 10, loadedChar.Zaubern, "Zaubern should be persisted")
		assert.Equal(t, 14, loadedChar.Raufen, "Raufen should be persisted")
	})
}

func TestCalculateBonuses(t *testing.T) {
	// Setup test environment
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original != "" {
			os.Setenv("ENVIRONMENT", original)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})

	t.Run("Calculate bonuses from character attributes", func(t *testing.T) {
		char := models.Char{
			BamortBase: models.BamortBase{ID: 1},
			Rasse:      "Mensch",
			Typ:        "Krieger",
			Eigenschaften: []models.Eigenschaft{
				{CharacterID: 1, Name: "St", Value: 85},
				{CharacterID: 1, Name: "Gs", Value: 70},
				{CharacterID: 1, Name: "Gw", Value: 65},
				{CharacterID: 1, Name: "Ko", Value: 75},
				{CharacterID: 1, Name: "In", Value: 50},
				{CharacterID: 1, Name: "Zt", Value: 30},
			},
		}

		bonuses := char.CalculateBonuses()

		// Ausdauer Bonus: Ko/10 + St/20 = 75/10 + 85/20 = 7 + 4 = 11
		assert.Equal(t, 11, bonuses.AusdauerBonus, "AusdauerBonus should be calculated correctly")

		// Schadens Bonus: St/20 + Gs/30 - 3 = 85/20 + 70/30 - 3 = 4 + 2 - 3 = 3
		assert.Equal(t, 3, bonuses.SchadensBonus, "SchadensBonus should be calculated correctly")

		// Angriffs Bonus: Gs 70 -> bonus 2 (61-80 range)
		assert.Equal(t, 2, bonuses.AngriffsBonus, "AngriffsBonus should be calculated correctly")

		// Abwehr Bonus: Gw 65 -> bonus 2 (61-80 range)
		assert.Equal(t, 2, bonuses.AbwehrBonus, "AbwehrBonus should be calculated correctly")

		// Zauber Bonus: Zt 30 -> bonus 0 (21-40 range)
		assert.Equal(t, 0, bonuses.ZauberBonus, "ZauberBonus should be calculated correctly")

		// Resistenz Bonus Körper: Ko 75 -> bonus 2 (61-80), Mensch Krieger +1 = 3
		assert.Equal(t, 3, bonuses.ResistenzBonusKoerper, "ResistenzBonusKoerper should be calculated correctly")

		// Resistenz Bonus Geist: In 50 -> bonus 1 (41-60), Mensch = 1
		assert.Equal(t, 1, bonuses.ResistenzBonusGeist, "ResistenzBonusGeist should be calculated correctly")
	})

	t.Run("Calculate bonuses for Zwerg Kämpfer", func(t *testing.T) {
		char := models.Char{
			Rasse: "Zwerge",
			Typ:   "Krieger",
			Eigenschaften: []models.Eigenschaft{
				{Name: "Ko", Value: 85},
				{Name: "In", Value: 50},
			},
		}

		bonuses := char.CalculateBonuses()

		// Zwerge get +3 base, Kämpfer +1 = 4
		assert.Equal(t, 4, bonuses.ResistenzBonusKoerper, "Zwerge Kämpfer ResistenzBonusKoerper should be 4")

		// Zwerge get +3 base
		assert.Equal(t, 3, bonuses.ResistenzBonusGeist, "Zwerge ResistenzBonusGeist should be 3")
	})

	t.Run("Calculate bonuses for Elf Magier", func(t *testing.T) {
		char := models.Char{
			Rasse: "Elfen",
			Typ:   "Magier",
			Eigenschaften: []models.Eigenschaft{
				{Name: "Ko", Value: 50},
				{Name: "In", Value: 85},
			},
		}

		bonuses := char.CalculateBonuses()

		// Elfen get +2 base, Magier (Zauberer) +2 = 4
		assert.Equal(t, 4, bonuses.ResistenzBonusKoerper, "Elfen Magier ResistenzBonusKoerper should be 4")

		// Elfen get +2 base, Magier (Zauberer) +2 = 4
		assert.Equal(t, 4, bonuses.ResistenzBonusGeist, "Elfen Magier ResistenzBonusGeist should be 4")
	})
}
