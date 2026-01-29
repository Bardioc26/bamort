package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"log"
)

// MigrateLearningCostsToDatabase überträgt die Daten aus learningCostsData in die Datenbank
func MigrateLearningCostsToDatabase() error {
	log.Println("Starting migration of learning costs data to database...")

	// 1. Quellen migrieren
	if err := migrateSources(); err != nil {
		return fmt.Errorf("failed to migrate sources: %w", err)
	}

	// 2. Charakterklassen migrieren
	if err := migrateCharacterClasses(); err != nil {
		return fmt.Errorf("failed to migrate character classes: %w", err)
	}

	// 2. Fertigkeitskategorien migrieren
	if err := migrateSkillCategories(); err != nil {
		return fmt.Errorf("failed to migrate skill categories: %w", err)
	}

	// 3. Schwierigkeitsgrade migrieren
	if err := migrateSkillDifficulties(); err != nil {
		return fmt.Errorf("failed to migrate skill difficulties: %w", err)
	}

	// 4. Zauberschulen migrieren
	if err := migrateSpellSchools(); err != nil {
		return fmt.Errorf("failed to migrate spell schools: %w", err)
	}

	// 5. EP-Kosten für Kategorien migrieren
	if err := migrateClassCategoryEPCosts(); err != nil {
		return fmt.Errorf("failed to migrate class category EP costs: %w", err)
	}

	// 5a. EP-Kosten für "Unbekannt"-Kategorie hinzufügen
	if err := migrateUnknownCategoryEPCosts(); err != nil {
		return fmt.Errorf("failed to migrate unknown category EP costs: %w", err)
	}

	// 6. EP-Kosten für Zauberschulen migrieren
	if err := migrateClassSpellSchoolEPCosts(); err != nil {
		return fmt.Errorf("failed to migrate class spell school EP costs: %w", err)
	}

	// 7. Zauber-Level LE-Kosten migrieren
	if err := migrateSpellLevelLECosts(); err != nil {
		return fmt.Errorf("failed to migrate spell level LE costs: %w", err)
	}

	// 8. Fertigkeits-Kategorien-Schwierigkeiten migrieren
	if err := migrateSkillCategoryDifficulties(); err != nil {
		return fmt.Errorf("failed to migrate skill category difficulties: %w", err)
	}

	// 9. Verbesserungskosten migrieren
	if err := migrateSkillImprovementCosts(); err != nil {
		return fmt.Errorf("failed to migrate skill improvement costs: %w", err)
	}

	log.Println("Migration completed successfully!")
	return nil
}

// migrateSources erstellt die Standardquellen
func migrateSources() error {
	sources := []models.Source{
		{
			Code:        "KOD",
			Name:        "Kodex",
			FullName:    "Midgard Regelwerk - Kodex",
			Edition:     "5. Edition",
			Publisher:   "Pegasus Spiele",
			IsCore:      true,
			IsActive:    true,
			GameSystem:  "midgard",
			Description: "Grundregelwerk für Midgard",
		},
		{
			Code:        "ARK",
			Name:        "Arkanum",
			FullName:    "Midgard Arkanum",
			Edition:     "5. Edition",
			Publisher:   "Pegasus Spiele",
			IsCore:      false,
			IsActive:    true,
			GameSystem:  "midgard",
			Description: "Erweiterungsregelwerk für Zauber und Magie",
		},
		{
			Code:        "MYS",
			Name:        "Mysterium",
			FullName:    "Midgard Mysterium",
			Edition:     "5. Edition",
			Publisher:   "Pegasus Spiele",
			IsCore:      false,
			IsActive:    true,
			GameSystem:  "midgard",
			Description: "Erweiterungsregelwerk für Geheimnisse und Mysterien",
		},
		{
			Code:        "UNB",
			Name:        "Unbekannt",
			FullName:    "Unbekannte Quelle",
			IsCore:      false,
			IsActive:    true,
			GameSystem:  "midgard",
			Description: "Für Inhalte ohne bekannte Quelle",
		},
	}

	for _, source := range sources {
		// Prüfe ob die Quelle bereits existiert
		var existing models.Source
		if err := existing.FirstByCode(source.Code); err != nil {
			// Quelle existiert nicht, erstelle sie
			if err := source.Create(); err != nil {
				return fmt.Errorf("failed to create source %s: %w", source.Code, err)
			}
			log.Printf("Created source: %s - %s", source.Code, source.Name)
		}
	}

	return nil
}

// migrateCharacterClasses erstellt Charakterklassen-Einträge
func migrateCharacterClasses() error {
	// Hole die KOD-Quelle für Charakterklassen (alle Grundklassen sind im Kodex)
	var kodSource models.Source
	if err := kodSource.FirstByCode("KOD"); err != nil {
		return fmt.Errorf("KOD source not found: %w", err)
	}

	characterClasses := map[string]string{
		"As": "Assassine",
		"Bb": "Barbar",
		"Gl": "Glücksritter",
		"Hä": "Händler",
		"Kr": "Krieger",
		"Sp": "Spitzbube",
		"Wa": "Waldläufer",
		"Ba": "Barde",
		"Or": "Ordenskrieger",
		"Dr": "Druide",
		"Hx": "Hexer",
		"Ma": "Magier",
		"PB": "Priester Beschützer",
		"PS": "Priester Streiter",
		"Sc": "Schamane",
	}
	gs := GetGameSystem(0, "")

	for code, name := range characterClasses {
		class := models.CharacterClass{
			Code:         code,
			Name:         name,
			SourceID:     kodSource.ID,
			GameSystemId: gs.ID,
		}

		// Prüfe ob die Klasse bereits existiert
		var existing models.CharacterClass
		if err := existing.FirstByCode(code); err != nil {
			// Klasse existiert nicht, erstelle sie
			if err := class.Create(); err != nil {
				return fmt.Errorf("failed to create character class %s: %w", code, err)
			}
			log.Printf("Created character class: %s - %s (Source: %s)", code, name, kodSource.Code)
		}
	}

	return nil
}

// migrateSkillCategories erstellt Fertigkeitskategorien
func migrateSkillCategories() error {
	// Hole die KOD-Quelle für Kategorien
	var kodSource models.Source
	if err := kodSource.FirstByCode("KOD"); err != nil {
		return fmt.Errorf("KOD source not found: %w", err)
	}

	categories := []string{
		"Unbekannt", // Für Fertigkeiten ohne bekannte Kategorie
		"Alltag", "Freiland", "Halbwelt", "Kampf", "Körper",
		"Sozial", "Unterwelt", "Waffen", "Wissen", "Schilde und Parierwaﬀen",
	}
	gs := GetGameSystem(0, "")

	for _, categoryName := range categories {
		sourceID := kodSource.ID
		// "Unbekannt" bekommt die UNB-Quelle
		if categoryName == "Unbekannt" {
			var unbSource models.Source
			if err := unbSource.FirstByCode("UNB"); err != nil {
				return fmt.Errorf("UNB source not found: %w", err)
			}
			sourceID = unbSource.ID
		}

		category := models.SkillCategory{
			Name:         categoryName,
			SourceID:     sourceID,
			GameSystemId: gs.ID,
		}

		// Prüfe ob die Kategorie bereits existiert
		var existing models.SkillCategory
		if err := existing.FirstByName(categoryName); err != nil {
			// Kategorie existiert nicht, erstelle sie
			if err := category.Create(); err != nil {
				return fmt.Errorf("failed to create skill category %s: %w", categoryName, err)
			}
			log.Printf("Created skill category: %s", categoryName)
		}
	}

	return nil
}

// migrateSkillDifficulties erstellt Schwierigkeitsgrade
func migrateSkillDifficulties() error {
	difficulties := []string{"leicht", "normal", "schwer", "sehr schwer"}
	gs := GetGameSystem(0, "")

	for _, difficultyName := range difficulties {
		difficulty := models.SkillDifficulty{
			Name:         difficultyName,
			GameSystemId: gs.ID,
		}

		// Prüfe ob die Schwierigkeit bereits existiert
		var existing models.SkillDifficulty
		if err := existing.FirstByName(difficultyName); err != nil {
			// Schwierigkeit existiert nicht, erstelle sie
			if err := difficulty.Create(); err != nil {
				return fmt.Errorf("failed to create skill difficulty %s: %w", difficultyName, err)
			}
			log.Printf("Created skill difficulty: %s", difficultyName)
		}
	}

	return nil
}

// migrateSpellSchools erstellt Zauberschulen
func migrateSpellSchools() error {
	// Hole die KOD-Quelle für Basis-Zauberschulen, ARK für erweiterte
	var kodSource models.Source
	if err := kodSource.FirstByCode("KOD"); err != nil {
		return fmt.Errorf("KOD source not found: %w", err)
	}

	var arkSource models.Source
	if err := arkSource.FirstByCode("ARK"); err != nil {
		return fmt.Errorf("ARK source not found: %w", err)
	}
	gs := GetGameSystem(0, "")

	schools := map[string]uint{
		// Basis-Zauberschulen (Kodex)
		"Beherrschen": kodSource.ID,
		"Bewegen":     kodSource.ID,
		"Erkennen":    kodSource.ID,
		"Erschaffen":  kodSource.ID,
		"Formen":      kodSource.ID,
		"Verändern":   kodSource.ID,
		"Zerstören":   kodSource.ID,
		"Wunder":      kodSource.ID,
		// Erweiterte Zauberschulen (Arkanum)
		"Dweomer": arkSource.ID,
		"Lied":    arkSource.ID,
	}

	for schoolName, sourceID := range schools {
		school := models.SpellSchool{
			Name:         schoolName,
			SourceID:     sourceID,
			GameSystemId: gs.ID,
		}

		// Prüfe ob die Schule bereits existiert
		var existing models.SpellSchool
		if err := existing.FirstByName(schoolName); err != nil {
			// Schule existiert nicht, erstelle sie
			if err := school.Create(); err != nil {
				return fmt.Errorf("failed to create spell school %s: %w", schoolName, err)
			}
			log.Printf("Created spell school: %s", schoolName)
		}
	}

	return nil
}

// migrateClassCategoryEPCosts migriert EP-Kosten für Kategorien
func migrateClassCategoryEPCosts() error {
	for classCode, categories := range learningCostsData.EPPerTE {
		// Hole die Charakterklasse
		var characterClass models.CharacterClass
		if err := characterClass.FirstByCode(classCode); err != nil {
			return fmt.Errorf("character class %s not found: %w", classCode, err)
		}

		for categoryName, epCost := range categories {
			// Hole die Kategorie
			var skillCategory models.SkillCategory
			if err := skillCategory.FirstByName(categoryName); err != nil {
				return fmt.Errorf("skill category %s not found: %w", categoryName, err)
			}

			// Erstelle EP-Kosten-Eintrag
			cost := models.ClassCategoryEPCost{
				CharacterClassID: characterClass.ID,
				SkillCategoryID:  skillCategory.ID,
				EPPerTE:          epCost,
			}

			if err := cost.Create(); err != nil {
				return fmt.Errorf("failed to create EP cost for class %s category %s: %w", classCode, categoryName, err)
			}
			log.Printf("Created EP cost: %s - %s = %d EP/TE", classCode, categoryName, epCost)
		}
	}

	return nil
}

// migrateClassSpellSchoolEPCosts migriert EP-Kosten für Zauberschulen
func migrateClassSpellSchoolEPCosts() error {
	for classCode, schools := range learningCostsData.SpellEPPerLE {
		// Hole die Charakterklasse
		var characterClass models.CharacterClass
		if err := characterClass.FirstByCode(classCode); err != nil {
			return fmt.Errorf("character class %s not found: %w", classCode, err)
		}

		for schoolName, epCost := range schools {
			// Überspringe Schulen mit 0 EP-Kosten (nicht verfügbar)
			if epCost == 0 {
				continue
			}

			// Hole die Zauberschule
			var spellSchool models.SpellSchool
			if err := spellSchool.FirstByName(schoolName); err != nil {
				return fmt.Errorf("spell school %s not found: %w", schoolName, err)
			}

			// Erstelle EP-Kosten-Eintrag
			cost := models.ClassSpellSchoolEPCost{
				CharacterClassID: characterClass.ID,
				SpellSchoolID:    spellSchool.ID,
				EPPerLE:          epCost,
			}

			if err := cost.Create(); err != nil {
				return fmt.Errorf("failed to create spell EP cost for class %s school %s: %w", classCode, schoolName, err)
			}
			log.Printf("Created spell EP cost: %s - %s = %d EP/LE", classCode, schoolName, epCost)
		}
	}

	return nil
}

// migrateSpellLevelLECosts migriert LE-Kosten pro Zauber-Level
func migrateSpellLevelLECosts() error {
	gs := GetGameSystem(0, "")
	for level, leCost := range learningCostsData.SpellLEPerLevel {
		cost := models.SpellLevelLECost{
			Level:        level,
			LERequired:   leCost,
			GameSystemId: gs.ID,
		}

		if err := cost.Create(); err != nil {
			return fmt.Errorf("failed to create spell level LE cost for level %d: %w", level, err)
		}
		log.Printf("Created spell level LE cost: Level %d = %d LE", level, leCost)
	}

	return nil
}

// migrateSkillCategoryDifficulties migriert Fertigkeits-Kategorien-Schwierigkeiten
func migrateSkillCategoryDifficulties() error {
	for categoryName, difficulties := range learningCostsData.ImprovementCost {
		// Hole die Kategorie
		var skillCategory models.SkillCategory
		if err := skillCategory.FirstByName(categoryName); err != nil {
			return fmt.Errorf("skill category %s not found: %w", categoryName, err)
		}

		for difficultyName, data := range difficulties {
			// Hole die Schwierigkeit
			var skillDifficulty models.SkillDifficulty
			if err := skillDifficulty.FirstByName(difficultyName); err != nil {
				return fmt.Errorf("skill difficulty %s not found: %w", difficultyName, err)
			}

			// Für jede Fertigkeit in dieser Kategorie/Schwierigkeit
			for _, skillName := range data.Skills {
				// Hole die Fertigkeit aus der bestehenden Skill-Tabelle
				var skill models.Skill
				if err := skill.First(skillName); err != nil {
					// Wenn die Fertigkeit nicht existiert, überspringe sie oder logge eine Warnung
					log.Printf("Warning: Skill %s not found in database, skipping", skillName)
					continue
				}

				// Erstelle Kategorie-Schwierigkeits-Eintrag
				categoryDifficulty := models.SkillCategoryDifficulty{
					SkillID:           skill.ID,
					SkillCategoryID:   skillCategory.ID,
					SkillDifficultyID: skillDifficulty.ID,
					LearnCost:         data.LearnCost,
				}

				if err := categoryDifficulty.Create(); err != nil {
					return fmt.Errorf("failed to create skill category difficulty for %s: %w", skillName, err)
				}
				log.Printf("Created skill category difficulty: %s - %s - %s (Learn: %d LE)",
					skillName, categoryName, difficultyName, data.LearnCost)
			}
		}
	}

	return nil
}

// migrateSkillImprovementCosts migriert Verbesserungskosten
func migrateSkillImprovementCosts() error {
	for categoryName, difficulties := range learningCostsData.ImprovementCost {
		// Hole die Kategorie
		var skillCategory models.SkillCategory
		if err := skillCategory.FirstByName(categoryName); err != nil {
			return fmt.Errorf("skill category %s not found: %w", categoryName, err)
		}

		for difficultyName, data := range difficulties {
			// Hole die Schwierigkeit
			var skillDifficulty models.SkillDifficulty
			if err := skillDifficulty.FirstByName(difficultyName); err != nil {
				return fmt.Errorf("skill difficulty %s not found: %w", difficultyName, err)
			}

			// Für jede Fertigkeit in dieser Kategorie/Schwierigkeit
			for _, skillName := range data.Skills {
				// Hole die Fertigkeit
				var skill models.Skill
				if err := skill.First(skillName); err != nil {
					log.Printf("Warning: Skill %s not found in database, skipping improvement costs", skillName)
					continue
				}

				// Finde den entsprechenden SkillCategoryDifficulty-Eintrag
				var categoryDifficulty models.SkillCategoryDifficulty
				if err := database.DB.Where("skill_id = ? AND skill_category_id = ? AND skill_difficulty_id = ?",
					skill.ID, skillCategory.ID, skillDifficulty.ID).First(&categoryDifficulty).Error; err != nil {
					log.Printf("Warning: SkillCategoryDifficulty not found for %s, skipping improvement costs", skillName)
					continue
				}

				// Für jeden Level in TrainCosts
				for level, teCost := range data.TrainCosts {
					improvementCost := models.SkillImprovementCost{
						SkillCategoryDifficultyID: categoryDifficulty.ID,
						CurrentLevel:              level,
						TERequired:                teCost,
					}

					if err := improvementCost.Create(); err != nil {
						return fmt.Errorf("failed to create improvement cost for %s level %d: %w", skillName, level, err)
					}
					log.Printf("Created improvement cost: %s - %s - %s Level %d = %d TE",
						skillName, categoryName, difficultyName, level, teCost)
				}
			}
		}
	}

	return nil
}

// migrateUnknownCategoryEPCosts fügt EP-Kosten für die "Unbekannt"-Kategorie hinzu
func migrateUnknownCategoryEPCosts() error {
	// Hole die "Unbekannt"-Kategorie
	var unknownCategory models.SkillCategory
	if err := unknownCategory.FirstByName("Unbekannt"); err != nil {
		return fmt.Errorf("unknown category not found: %w", err)
	}

	// Hole alle Charakterklassen
	var characterClasses []models.CharacterClass
	if err := database.DB.Find(&characterClasses).Error; err != nil {
		return fmt.Errorf("failed to fetch character classes: %w", err)
	}

	// Standard EP-Kosten für "Unbekannt": 50 EP/TE
	standardEPCost := 50

	for _, characterClass := range characterClasses {
		// Prüfe ob bereits EP-Kosten für diese Klasse und "Unbekannt" existieren
		var existingCost models.ClassCategoryEPCost
		err := database.DB.Where("character_class_id = ? AND skill_category_id = ?",
			characterClass.ID, unknownCategory.ID).First(&existingCost).Error

		if err == nil {
			// EP-Kosten existieren bereits, überspringe
			log.Printf("EP cost for class %s and category 'Unbekannt' already exists, skipping", characterClass.Code)
			continue
		}

		// Erstelle EP-Kosten-Eintrag für "Unbekannt"
		cost := models.ClassCategoryEPCost{
			CharacterClassID: characterClass.ID,
			SkillCategoryID:  unknownCategory.ID,
			EPPerTE:          standardEPCost,
		}

		if err := cost.Create(); err != nil {
			return fmt.Errorf("failed to create EP cost for class %s category 'Unbekannt': %w", characterClass.Code, err)
		}
		log.Printf("Created EP cost: %s - Unbekannt = %d EP/TE", characterClass.Code, standardEPCost)
	}

	return nil
}
