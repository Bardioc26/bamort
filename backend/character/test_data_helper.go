package character

import (
	"bamort/database"
	"bamort/models"
)

// createTestSkillData erstellt Test-Daten für Skills und Spells
func createTestSkillData() error {
	// Test-Fertigkeit erstellen
	testSkill := models.SkFertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Menschenkenntnis",
			},
			CharacterID: 0, // Global skill
		},
		Beschreibung:    "Test-Fertigkeit für PP Tests",
		Fertigkeitswert: 0,
		Improvable:      true,
		Category:        "Sozial",
	}
	if err := database.DB.Create(&testSkill).Error; err != nil {
		return err
	}

	// Test-Zauber erstellen
	testSpell := models.SkZauber{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Macht über das Selbst",
			},
			CharacterID: 0, // Global spell
		},
		Beschreibung: "Test-Zauber für PP Tests",
		Quelle:       "Beherrschen",
	}
	if err := database.DB.Create(&testSpell).Error; err != nil {
		return err
	}

	// GSM Test-Skill erstellen
	gsmSkill := models.Skill{
		GameSystemId: 1,
		Name:         "Menschenkenntnis",
		Beschreibung: "Test Skill",
		Quelle:       "Test",
		Category:     "Sozial",
		Initialwert:  5,
		Improvable:   true,
		Difficulty:   "normal",
	}
	if err := database.DB.Create(&gsmSkill).Error; err != nil {
		return err
	}

	// GSM Test-Spell erstellen
	gsmSpell := models.Spell{
		GameSystemId: 1,
		Name:         "Macht über das Selbst",
		Beschreibung: "Test Spell",
		Quelle:       "Test",
		Stufe:        1,
		AP:           "1",
		Category:     "Beherrschen",
	}
	if err := database.DB.Create(&gsmSpell).Error; err != nil {
		return err
	}

	return nil
}

// cleanupTestSkillData entfernt Test-Daten
func cleanupTestSkillData() {
	database.DB.Where("name = ?", "Menschenkenntnis").Delete(&models.SkFertigkeit{})
	database.DB.Where("name = ?", "Macht über das Selbst").Delete(&models.SkZauber{})
	database.DB.Where("name = ?", "Menschenkenntnis").Delete(&models.Skill{})
	database.DB.Where("name = ?", "Macht über das Selbst").Delete(&models.Spell{})
}
