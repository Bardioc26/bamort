package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"log"
)

// EnhanceLearningDataWithExistingTables ergänzt die neuen Lernkosten-Tabellen mit Daten aus den bereits vorhandenen Tabellen
func EnhanceLearningDataWithExistingTables() error {
	log.Println("Starting enhancement of learning data with existing tables...")

	// 1. Fertigkeiten aus der bestehenden Skill-Tabelle in SkillCategoryDifficulty einbinden
	if err := linkExistingSkillsToCategories(); err != nil {
		return fmt.Errorf("failed to link existing skills to categories: %w", err)
	}

	// 2. Zauber aus der bestehenden Spell-Tabelle mit Zauberschulen verknüpfen
	if err := linkExistingSpellsToSchools(); err != nil {
		return fmt.Errorf("failed to link existing spells to schools: %w", err)
	}

	// 3. Charaktere aus der bestehenden Character-Tabelle mit Charakterklassen verknüpfen
	if err := linkExistingCharactersToClasses(); err != nil {
		return fmt.Errorf("failed to link existing characters to classes: %w", err)
	}

	// 4. Fehlende Fertigkeiten aus der bestehenden Tabelle zu den Kategorien hinzufügen
	if err := addMissingSkillsToCategories(); err != nil {
		return fmt.Errorf("failed to add missing skills to categories: %w", err)
	}

	// 5. Verknüpfe bestehende Fertigkeiten und Zauber mit Standardquellen
	if err := linkExistingItemsToDefaultSources(); err != nil {
		return fmt.Errorf("failed to link existing items to default sources: %w", err)
	}

	log.Println("Enhancement completed successfully!")
	return nil
}

// linkExistingSkillsToCategories verknüpft vorhandene Fertigkeiten mit Kategorien und Schwierigkeiten
func linkExistingSkillsToCategories() error {
	var skills []models.Skill
	if err := database.DB.Find(&skills).Error; err != nil {
		return fmt.Errorf("failed to fetch existing skills: %w", err)
	}

	for _, skill := range skills {
		// Prüfe ob bereits Verknüpfungen für diese Fertigkeit existieren
		var existingLinks []models.SkillCategoryDifficulty
		if err := database.DB.Where("skill_id = ?", skill.ID).Find(&existingLinks).Error; err != nil {
			return fmt.Errorf("failed to check existing links for skill %s: %w", skill.Name, err)
		}

		// Wenn bereits Verknüpfungen existieren, überspringe
		if len(existingLinks) > 0 {
			log.Printf("Skill %s already has category links, skipping", skill.Name)
			continue
		}

		// Suche nach der Fertigkeit in den lerningCostsData
		category, difficulty := findSkillInLearningData(skill.Name)
		if category == "" {
			log.Printf("Warning: Skill %s not found in learning data, using default category", skill.Name)
			// Verwende eine Standard-Kategorie basierend auf der bestehenden Category
			if skill.Category != "" {
				category = skill.Category
				difficulty = skill.Difficulty
				if difficulty == "" {
					difficulty = "normal"
				}
			} else {
				category = "Unbekannt"
				difficulty = "normal"
			}
		}

		// Erstelle Verknüpfung
		if err := createSkillCategoryLink(skill.ID, skill.Name, category, difficulty); err != nil {
			log.Printf("Warning: Failed to create category link for skill %s: %v", skill.Name, err)
		}
	}

	return nil
}

// linkExistingSpellsToSchools verknüpft vorhandene Zauber mit Zauberschulen
func linkExistingSpellsToSchools() error {
	var spells []models.Spell
	if err := database.DB.Find(&spells).Error; err != nil {
		return fmt.Errorf("failed to fetch existing spells: %w", err)
	}

	for _, spell := range spells {
		// Prüfe ob die Zauberschule existiert
		if spell.Category == "" {
			log.Printf("Warning: Spell %s has no category, skipping", spell.Name)
			continue
		}

		var spellSchool models.SpellSchool
		if err := spellSchool.FirstByName(spell.Category); err != nil {
			// Zauberschule existiert nicht, erstelle sie
			newSchool := models.SpellSchool{
				Name:         spell.Category,
				GameSystemId: 1,
			}
			if err := newSchool.Create(); err != nil {
				log.Printf("Warning: Failed to create spell school %s for spell %s: %v", spell.Category, spell.Name, err)
				continue
			}
			log.Printf("Created new spell school: %s", spell.Category)
		}

		log.Printf("Spell %s linked to school %s", spell.Name, spell.Category)
	}

	return nil
}

// linkExistingCharactersToClasses verknüpft vorhandene Charaktere mit Charakterklassen
func linkExistingCharactersToClasses() error {
	var characters []models.Char
	if err := database.DB.Find(&characters).Error; err != nil {
		return fmt.Errorf("failed to fetch existing characters: %w", err)
	}

	// Lade alle verfügbaren Charakterklassen aus der Datenbank
	var characterClasses []models.CharacterClass
	if err := database.DB.Find(&characterClasses).Error; err != nil {
		return fmt.Errorf("failed to fetch character classes from database: %w", err)
	}

	// Erstelle Mapping-Maps für beide Richtungen (Code->Klasse und Name->Klasse)
	codeToClass := make(map[string]models.CharacterClass)
	nameToClass := make(map[string]models.CharacterClass)

	for _, class := range characterClasses {
		codeToClass[class.Code] = class
		nameToClass[class.Name] = class
	}

	for _, character := range characters {
		var characterClass models.CharacterClass
		var found bool

		// Versuche zuerst über Code zu finden (falls character.Typ bereits ein Code ist)
		if class, exists := codeToClass[character.Typ]; exists {
			characterClass = class
			found = true
		} else if class, exists := nameToClass[character.Typ]; exists {
			// Wenn nicht über Code gefunden, versuche über vollständigen Namen
			characterClass = class
			found = true
		}

		if !found {
			log.Printf("Warning: Character class '%s' not found in database for character %s", character.Typ, character.Name)
			continue
		}

		log.Printf("Character %s linked to class %s (%s)", character.Name, characterClass.Code, characterClass.Name)
	}

	return nil
}

// addMissingSkillsToCategories fügt Fertigkeiten hinzu, die in der Datenbank existieren aber nicht in den Lernkosten-Daten
func addMissingSkillsToCategories() error {
	var skills []models.Skill
	if err := database.DB.Find(&skills).Error; err != nil {
		return fmt.Errorf("failed to fetch skills: %w", err)
	}

	for _, skill := range skills {
		// Prüfe ob bereits Kategorien-Verknüpfungen existieren
		var count int64
		if err := database.DB.Model(&models.SkillCategoryDifficulty{}).Where("skill_id = ?", skill.ID).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count category links for skill %s: %w", skill.Name, err)
		}

		if count == 0 {
			// Keine Kategorien-Verknüpfung gefunden, erstelle eine Standard-Verknüpfung
			category := "Unbekannt"
			difficulty := "normal"

			// Verwende bestehende Kategorie falls vorhanden
			if skill.Category != "" {
				category = skill.Category
			}
			if skill.Difficulty != "" {
				difficulty = skill.Difficulty
			}

			if err := createSkillCategoryLink(skill.ID, skill.Name, category, difficulty); err != nil {
				log.Printf("Warning: Failed to create default category link for skill %s: %v", skill.Name, err)
			} else {
				log.Printf("Created default category link for skill %s: %s - %s", skill.Name, category, difficulty)
			}
		}
	}

	return nil
}

// Hilfsfunktionen

// findSkillInLearningData sucht eine Fertigkeit in den learningCostsData und gibt Kategorie und Schwierigkeit zurück
func findSkillInLearningData(skillName string) (string, string) {
	for categoryName, difficulties := range learningCostsData.ImprovementCost {
		for difficultyName, data := range difficulties {
			for _, skill := range data.Skills {
				if skill == skillName {
					return categoryName, difficultyName
				}
			}
		}
	}
	return "", ""
}

// createSkillCategoryLink erstellt eine Verknüpfung zwischen Fertigkeit, Kategorie und Schwierigkeit
func createSkillCategoryLink(skillID uint, skillName, categoryName, difficultyName string) error {
	// Hole oder erstelle die Kategorie
	var skillCategory models.SkillCategory
	gs := GetGameSystem(0, "midgard")
	if err := skillCategory.FirstByName(categoryName); err != nil {
		// Kategorie existiert nicht, erstelle sie
		skillCategory = models.SkillCategory{
			Name:         categoryName,
			GameSystemId: gs.ID,
		}
		if err := skillCategory.Create(); err != nil {
			return fmt.Errorf("failed to create skill category %s: %w", categoryName, err)
		}
		log.Printf("Created new skill category: %s", categoryName)
	}

	// Hole oder erstelle die Schwierigkeit
	var skillDifficulty models.SkillDifficulty
	if err := skillDifficulty.FirstByName(difficultyName); err != nil {
		// Schwierigkeit existiert nicht, erstelle sie
		skillDifficulty = models.SkillDifficulty{
			Name:         difficultyName,
			GameSystemId: gs.ID,
		}
		if err := skillDifficulty.Create(); err != nil {
			return fmt.Errorf("failed to create skill difficulty %s: %w", difficultyName, err)
		}
		log.Printf("Created new skill difficulty: %s", difficultyName)
	}

	// Bestimme LearnCost basierend auf den learningCostsData
	learnCost := 2 // Standard-Wert
	if categoryData, exists := learningCostsData.ImprovementCost[categoryName]; exists {
		if difficultyData, exists := categoryData[difficultyName]; exists {
			learnCost = difficultyData.LearnCost
		}
	}

	// Erstelle die Verknüpfung
	categoryDifficulty := models.SkillCategoryDifficulty{
		SkillID:           skillID,
		SkillCategoryID:   skillCategory.ID,
		SkillDifficultyID: skillDifficulty.ID,
		LearnCost:         learnCost,
	}

	if err := categoryDifficulty.Create(); err != nil {
		return fmt.Errorf("failed to create skill category difficulty link: %w", err)
	}

	log.Printf("Created skill category link: %s - %s - %s (Learn: %d LE)", skillName, categoryName, difficultyName, learnCost)
	return nil
}

// CreateLearningCostsTables erstellt alle neuen Tabellen für das Lernkosten-System
func CreateLearningCostsTables() error {
	log.Println("Creating learning costs tables...")

	// Liste aller neuen Modelle
	models := []interface{}{
		&models.Source{},
		&models.CharacterClass{},
		&models.SkillCategory{},
		&models.SkillDifficulty{},
		&models.SpellSchool{},
		&models.ClassCategoryEPCost{},
		&models.ClassSpellSchoolEPCost{},
		&models.SpellLevelLECost{},
		&models.SkillCategoryDifficulty{},
		&models.SkillImprovementCost{},
	}

	// Erstelle oder migriere alle Tabellen
	for _, model := range models {
		if err := database.DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate table for model %T: %w", model, err)
		}
	}

	log.Println("Learning costs tables created successfully!")
	return nil
}

// linkExistingItemsToDefaultSources verknüpft bestehende Fertigkeiten und Zauber mit Standardquellen
func linkExistingItemsToDefaultSources() error {
	// Hole die Standardquellen
	var kodSource models.Source
	if err := kodSource.FirstByCode("KOD"); err != nil {
		return fmt.Errorf("KOD source not found: %w", err)
	}

	var arkSource models.Source
	if err := arkSource.FirstByCode("ARK"); err != nil {
		return fmt.Errorf("ARK source not found: %w", err)
	}

	// Verknüpfe Fertigkeiten mit KOD (falls noch keine Quelle zugeordnet)
	var skills []models.Skill
	if err := database.DB.Where("source_id IS NULL OR source_id = 0").Find(&skills).Error; err != nil {
		return fmt.Errorf("failed to fetch skills without source: %w", err)
	}

	for _, skill := range skills {
		if err := database.DB.Model(&skill).Update("source_id", kodSource.ID).Error; err != nil {
			log.Printf("Warning: Failed to update source for skill %s: %v", skill.Name, err)
		} else {
			log.Printf("Linked skill %s to source KOD", skill.Name)
		}
	}

	// Verknüpfe Zauber mit ARK (falls noch keine Quelle zugeordnet)
	var spells []models.Spell
	if err := database.DB.Where("source_id IS NULL OR source_id = 0").Find(&spells).Error; err != nil {
		return fmt.Errorf("failed to fetch spells without source: %w", err)
	}

	for _, spell := range spells {
		if err := database.DB.Model(&spell).Update("source_id", arkSource.ID).Error; err != nil {
			log.Printf("Warning: Failed to update source for spell %s: %v", spell.Name, err)
		} else {
			log.Printf("Linked spell %s to source ARK", spell.Name)
		}
	}

	return nil
}
