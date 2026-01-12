package models

import (
	"bamort/database"
	"fmt"
)

// PopulateClassLearningPointsData inserts all hardcoded learning points data into the database
func PopulateClassLearningPointsData() error {
	// Define all character class learning data
	classData := []struct {
		ClassCode      string
		ClassName      string
		LearningPoints map[string]int
		SpellPoints    int
		TypicalSkills  []struct {
			Name      string
			Bonus     int
			Attribute string
			Notes     string
		}
		TypicalSpells []string
	}{
		{
			ClassCode: "As",
			ClassName: "Assassine",
			LearningPoints: map[string]int{
				"Alltag":    1,
				"Halbwelt":  2,
				"Sozial":    4,
				"Unterwelt": 8,
				"Waffen":    24,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Meucheln", Bonus: 8, Attribute: "Gs", Notes: ""},
			},
		},
		{
			ClassCode: "Bb",
			ClassName: "Barbar",
			LearningPoints: map[string]int{
				"Alltag":   2,
				"Freiland": 4,
				"Kampf":    1,
				"Körper":   2,
				"Waffen":   24,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Spurensuche", Bonus: 8, Attribute: "In", Notes: "in Heimatlandschaft"},
				{Name: "Überleben", Bonus: 8, Attribute: "In", Notes: "in Heimatlandschaft"},
			},
		},
		{
			ClassCode: "Gl",
			ClassName: "Glücksritter",
			LearningPoints: map[string]int{
				"Alltag":   2,
				"Halbwelt": 3,
				"Sozial":   8,
				"Waffen":   24,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Fechten", Bonus: 5, Attribute: "Gs", Notes: "oder beidhändiger Kampf+5 (Gs)"},
			},
		},
		{
			ClassCode: "Hä",
			ClassName: "Händler",
			LearningPoints: map[string]int{
				"Alltag": 4,
				"Sozial": 8,
				"Wissen": 4,
				"Waffen": 20,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Geschäftssinn", Bonus: 8, Attribute: "In", Notes: ""},
			},
		},
		{
			ClassCode: "Kr",
			ClassName: "Krieger",
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Kampf":  3,
				"Körper": 1,
				"Waffen": 36,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Kampf in Vollrüstung", Bonus: 5, Attribute: "St", Notes: ""},
			},
		},
		{
			ClassCode: "Sp",
			ClassName: "Spitzbube",
			LearningPoints: map[string]int{
				"Alltag":    2,
				"Halbwelt":  6,
				"Unterwelt": 12,
				"Waffen":    20,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Fallenmechanik", Bonus: 8, Attribute: "Gs", Notes: "oder Geschäftssinn+8 (In)"},
			},
		},
		{
			ClassCode: "Wa",
			ClassName: "Waldläufer",
			LearningPoints: map[string]int{
				"Alltag":   1,
				"Freiland": 11,
				"Körper":   4,
				"Waffen":   20,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Scharfschießen", Bonus: 5, Attribute: "Gs", Notes: ""},
			},
		},
		{
			ClassCode:   "Ba",
			ClassName:   "Barde",
			SpellPoints: 3,
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Sozial": 4,
				"Wissen": 4,
				"Waffen": 16,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Musizieren", Bonus: 12, Attribute: "Gs", Notes: ""},
				{Name: "Landeskunde", Bonus: 8, Attribute: "In", Notes: "für Heimat"},
			},
			TypicalSpells: []string{"Zauberlieder"},
		},
		{
			ClassCode:   "Or",
			ClassName:   "Ordenskrieger",
			SpellPoints: 3,
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Kampf":  3,
				"Wissen": 2,
				"Waffen": 18,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Athletik", Bonus: 8, Attribute: "St", Notes: "oder Meditieren+8 (Wk)"},
			},
			TypicalSpells: []string{"Wundertaten"},
		},
		{
			ClassCode:   "Dr",
			ClassName:   "Druide",
			SpellPoints: 5,
			LearningPoints: map[string]int{
				"Alltag":   2,
				"Freiland": 4,
				"Wissen":   2,
				"Waffen":   6,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Pflanzenkunde", Bonus: 8, Attribute: "In", Notes: ""},
				{Name: "Schreiben", Bonus: 12, Attribute: "In", Notes: "für Ogam-Zeichen"},
			},
			TypicalSpells: []string{"Dweomer", "Tiere rufen"},
		},
		{
			ClassCode:   "Hx",
			ClassName:   "Hexer",
			SpellPoints: 6,
			LearningPoints: map[string]int{
				"Alltag": 3,
				"Sozial": 2,
				"Wissen": 2,
				"Waffen": 2,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Gassenwissen", Bonus: 8, Attribute: "In", Notes: "oder Verführen+8 (pA)"},
			},
			TypicalSpells: []string{"Beherrschen", "Verändern", "Verwünschen", "Binden des Vertrauten"},
		},
		{
			ClassCode:   "Ma",
			ClassName:   "Magier",
			SpellPoints: 7,
			LearningPoints: map[string]int{
				"Alltag": 1,
				"Wissen": 5,
				"Waffen": 2,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Zauberkunde", Bonus: 8, Attribute: "In", Notes: ""},
				{Name: "Schreiben", Bonus: 12, Attribute: "In", Notes: "für Muttersprache"},
			},
			TypicalSpells: []string{"beliebig außer Dweomer, Wundertaten, Zauberlieder", "Erkennen von Zauberei"},
		},
		{
			ClassCode:   "PB",
			ClassName:   "Priester Beschützer",
			SpellPoints: 5,
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Sozial": 2,
				"Wissen": 3,
				"Waffen": 6,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Menschenkenntnis", Bonus: 8, Attribute: "In", Notes: ""},
				{Name: "Schreiben", Bonus: 12, Attribute: "In", Notes: "für Muttersprache"},
			},
			TypicalSpells: []string{"Wundertaten", "Heilen von Wunden"},
		},
		{
			ClassCode:   "PS",
			ClassName:   "Priester Streiter",
			SpellPoints: 5,
			LearningPoints: map[string]int{
				"Alltag": 3,
				"Kampf":  2,
				"Wissen": 2,
				"Waffen": 8,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Erste Hilfe", Bonus: 8, Attribute: "Gs", Notes: ""},
				{Name: "Schreiben", Bonus: 12, Attribute: "In", Notes: "für Muttersprache"},
			},
			TypicalSpells: []string{"Wundertaten", "Bannen von Finsterwerk", "Strahlender Panzer"},
		},
		{
			ClassCode:   "Sc",
			ClassName:   "Schamane",
			SpellPoints: 5,
			LearningPoints: map[string]int{
				"Alltag": 2,
				"Körper": 4,
				"Wissen": 2,
				"Waffen": 6,
			},
			TypicalSkills: []struct {
				Name      string
				Bonus     int
				Attribute string
				Notes     string
			}{
				{Name: "Tierkunde", Bonus: 8, Attribute: "In", Notes: ""},
				{Name: "Überleben", Bonus: 8, Attribute: "In", Notes: "in Heimatlandschaft"},
			},
			TypicalSpells: []string{"Dweomer", "Wundertaten", "Austreibung des Bösen", "Bannen von Gift"},
		},
	}

	// Process each character class
	for _, cd := range classData {
		// Find or create character class
		var charClass CharacterClass
		result := database.DB.Where("code = ?", cd.ClassCode).First(&charClass)
		if result.Error != nil {
			fmt.Printf("Warning: Character class %s (%s) not found in database, skipping\n", cd.ClassName, cd.ClassCode)
			continue
		}

		// Insert learning points for each category
		for categoryName, points := range cd.LearningPoints {
			var category SkillCategory
			if err := database.DB.Where("name = ?", categoryName).First(&category).Error; err != nil {
				fmt.Printf("Warning: Skill category %s not found, skipping\n", categoryName)
				continue
			}

			learningPoints := ClassLearningPoints{
				CharacterClassID: charClass.ID,
				SkillCategoryID:  category.ID,
				Points:           points,
			}

			if err := database.DB.Where("character_class_id = ? AND skill_category_id = ?",
				charClass.ID, category.ID).FirstOrCreate(&learningPoints).Error; err != nil {
				return fmt.Errorf("failed to insert learning points for %s/%s: %w", cd.ClassName, categoryName, err)
			}
		}

		// Insert spell points if applicable
		if cd.SpellPoints > 0 {
			spellPoints := ClassSpellPoints{
				CharacterClassID: charClass.ID,
				SpellPoints:      cd.SpellPoints,
			}

			if err := database.DB.Where("character_class_id = ?", charClass.ID).
				FirstOrCreate(&spellPoints).Error; err != nil {
				return fmt.Errorf("failed to insert spell points for %s: %w", cd.ClassName, err)
			}
		}

		// Insert typical skills
		for _, ts := range cd.TypicalSkills {
			var skill Skill
			if err := database.DB.Where("name = ?", ts.Name).First(&skill).Error; err != nil {
				fmt.Printf("Warning: Skill %s not found for class %s, skipping\n", ts.Name, cd.ClassName)
				continue
			}

			typicalSkill := ClassTypicalSkill{
				CharacterClassID: charClass.ID,
				SkillID:          skill.ID,
				Bonus:            ts.Bonus,
				Attribute:        ts.Attribute,
				Notes:            ts.Notes,
			}

			if err := database.DB.Where("character_class_id = ? AND skill_id = ?",
				charClass.ID, skill.ID).FirstOrCreate(&typicalSkill).Error; err != nil {
				return fmt.Errorf("failed to insert typical skill %s for %s: %w", ts.Name, cd.ClassName, err)
			}
		}

		// Insert typical spells
		for _, spellName := range cd.TypicalSpells {
			var spell Spell
			if err := database.DB.Where("name = ?", spellName).First(&spell).Error; err != nil {
				fmt.Printf("Warning: Spell %s not found for class %s, creating reference with notes\n", spellName, cd.ClassName)
				// For special cases like "beliebig außer...", we can skip or create a placeholder
				continue
			}

			typicalSpell := ClassTypicalSpell{
				CharacterClassID: charClass.ID,
				SpellID:          spell.ID,
				Notes:            "",
			}

			if err := database.DB.Where("character_class_id = ? AND spell_id = ?",
				charClass.ID, spell.ID).FirstOrCreate(&typicalSpell).Error; err != nil {
				return fmt.Errorf("failed to insert typical spell %s for %s: %w", spellName, cd.ClassName, err)
			}
		}

		fmt.Printf("Populated learning points data for %s (%s)\n", cd.ClassName, cd.ClassCode)
	}

	return nil
}
