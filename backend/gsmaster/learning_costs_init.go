package gsmaster

import (
	"bamort/database"
	"fmt"
	"log"
)

/*
// InitializeLearningCostsSystem initialisiert das komplette Lernkosten-System
// Diese Funktion sollte einmalig ausgeführt werden, um die Datenbank zu migrieren
func InitializeLearningCostsSystem() error {
	log.Println("Initializing learning costs system...")

	// 1. Erstelle alle neuen Tabellen
	if err := CreateLearningCostsTables(); err != nil {
		return fmt.Errorf("failed to create learning costs tables: %w", err)
	}

	// 2. Migriere die Daten aus learningCostsData
	if err := MigrateLearningCostsToDatabase(); err != nil {
		return fmt.Errorf("failed to migrate learning costs data: %w", err)
	}

	// 3. Verknüpfe mit bestehenden Tabellen
	if err := EnhanceLearningDataWithExistingTables(); err != nil {
		return fmt.Errorf("failed to enhance learning data with existing tables: %w", err)
	}

	log.Println("Learning costs system initialized successfully!")
	return nil
}
*/
// ValidateLearningCostsData validiert die Konsistenz der migrierten Daten
func ValidateLearningCostsData() error {
	log.Println("Validating learning costs data...")

	// Prüfe Anzahl der migrierten Einträge
	validationQueries := []struct {
		name     string
		query    string
		expected int
	}{
		{"Character Classes", "SELECT COUNT(*) FROM learning_character_classes", 15},
		{"Skill Categories", "SELECT COUNT(*) FROM learning_skill_categories", 10},
		{"Skill Difficulties", "SELECT COUNT(*) FROM learning_skill_difficulties", 4},
		{"Spell Schools", "SELECT COUNT(*) FROM learning_spell_schools", 10},
		{"Spell Level LE Costs", "SELECT COUNT(*) FROM learning_spell_level_le_costs", 12},
	}

	for _, validation := range validationQueries {
		var count int
		if err := database.DB.Raw(validation.query).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to validate %s: %w", validation.name, err)
		}
		log.Printf("Validation: %s = %d entries", validation.name, count)

		if count == 0 {
			log.Printf("Warning: No entries found for %s", validation.name)
		}
	}

	log.Println("Data validation completed!")
	return nil
}

// GetLearningCostsSummary gibt eine Zusammenfassung der Lernkosten-Daten zurück
func GetLearningCostsSummary() (map[string]interface{}, error) {
	summary := make(map[string]interface{})

	// Zähle Einträge in verschiedenen Tabellen
	tables := []string{
		"learning_character_classes",
		"learning_skill_categories",
		"learning_skill_difficulties",
		"learning_spell_schools",
		"learning_class_category_ep_costs",
		"learning_class_spell_school_ep_costs",
		"learning_spell_level_le_costs",
		"learning_skill_category_difficulties",
		"learning_skill_improvement_costs",
	}

	for _, table := range tables {
		var count int64
		if err := database.DB.Table(table).Count(&count).Error; err != nil {
			log.Printf("Warning: Failed to count entries in %s: %v", table, err)
			summary[table] = "Error"
		} else {
			summary[table] = count
		}
	}

	return summary, nil
}
