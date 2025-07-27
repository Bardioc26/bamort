package gsmaster

import (
	"fmt"

	"gorm.io/gorm"
)

// MigrateLearningCostsData migrates the static learningCostsData into the database
func MigrateLearningCostsData(db *gorm.DB) error {
	// First, ensure all tables exist
	if err := MigrateLearningCostsTables(db); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}

	// Clear existing data to avoid duplicates
	if err := clearLearningCostsData(db); err != nil {
		return fmt.Errorf("failed to clear existing data: %w", err)
	}

	// Migrate character classes
	if err := migrateCharacterClasses(db); err != nil {
		return fmt.Errorf("failed to migrate character classes: %w", err)
	}

	// Migrate skill categories
	if err := migrateSkillCategories(db); err != nil {
		return fmt.Errorf("failed to migrate skill categories: %w", err)
	}

	// Migrate spell schools
	if err := migrateSpellSchools(db); err != nil {
		return fmt.Errorf("failed to migrate spell schools: %w", err)
	}

	// Migrate skill difficulties
	if err := migrateSkillDifficulties(db); err != nil {
		return fmt.Errorf("failed to migrate skill difficulties: %w", err)
	}

	// Migrate EP costs for class/category combinations
	if err := migrateClassCategoryEPCosts(db); err != nil {
		return fmt.Errorf("failed to migrate class category EP costs: %w", err)
	}

	// Migrate spell EP costs for class/school combinations
	if err := migrateClassSpellSchoolCosts(db); err != nil {
		return fmt.Errorf("failed to migrate class spell school costs: %w", err)
	}

	// Migrate skill definitions
	if err := migrateSkillDefinitions(db); err != nil {
		return fmt.Errorf("failed to migrate skill definitions: %w", err)
	}

	// Migrate skill improvement costs
	if err := migrateSkillImprovementCosts(db); err != nil {
		return fmt.Errorf("failed to migrate skill improvement costs: %w", err)
	}

	// Migrate spell level costs
	if err := migrateSpellLevelCosts(db); err != nil {
		return fmt.Errorf("failed to migrate spell level costs: %w", err)
	}

	return nil
}

// clearLearningCostsData clears all existing learning costs data
func clearLearningCostsData(db *gorm.DB) error {
	tables := []string{
		"skill_improvement_costs",
		"skill_definitions",
		"class_spell_school_costs",
		"class_category_ep_costs",
		"spell_level_costs",
		"skill_difficulties",
		"spell_schools",
		"skill_categories",
		"character_classes",
	}

	for _, table := range tables {
		if err := db.Exec("DELETE FROM " + table).Error; err != nil {
			return fmt.Errorf("failed to clear table %s: %w", table, err)
		}
	}
	return nil
}

// migrateCharacterClasses migrates character class data
func migrateCharacterClasses(db *gorm.DB) error {
	classes := []CharacterClass{
		{GameSystem: "midgard", Name: "Assassine", Description: "Assassine", Quelle: "Midgard5", Code: "As"},
		{GameSystem: "midgard", Name: "Barbar", Description: "Barbar", Quelle: "Midgard5", Code: "Ba"},
		{GameSystem: "midgard", Name: "Druide", Description: "Druide", Quelle: "Midgard5", Code: "Dr"},
		{GameSystem: "midgard", Name: "Händler", Description: "Händler", Quelle: "Midgard5", Code: "Hä"},
		{GameSystem: "midgard", Name: "Hexer", Description: "Hexer", Quelle: "Midgard5", Code: "Hx"},
		{GameSystem: "midgard", Name: "Krieger", Description: "Krieger", Quelle: "Midgard5", Code: "Kr"},
		{GameSystem: "midgard", Name: "Magier", Description: "Magier", Quelle: "Midgard5", Code: "Ma"},
		{GameSystem: "midgard", Name: "Ordenskrieger", Description: "Ordenskrieger", Quelle: "Midgard5", Code: "Or"},
		{GameSystem: "midgard", Name: "Priester", Description: "Priester", Quelle: "Midgard5", Code: "Pr"},
		{GameSystem: "midgard", Name: "Schamane", Description: "Schamane", Quelle: "Midgard5", Code: "Sc"},
		{GameSystem: "midgard", Name: "Seefahrer", Description: "Seefahrer", Quelle: "Midgard5", Code: "Sf"},
		{GameSystem: "midgard", Name: "Spitzbube", Description: "Spitzbube", Quelle: "Midgard5", Code: "Sp"},
		{GameSystem: "midgard", Name: "Thaumaturg", Description: "Thaumaturg", Quelle: "Midgard5", Code: "Th"},
		{GameSystem: "midgard", Name: "Waldläufer", Description: "Waldläufer", Quelle: "Midgard5", Code: "Wl"},
		{GameSystem: "midgard", Name: "Zauberer", Description: "Zauberer", Quelle: "Midgard5", Code: "Zb"},
	}

	for _, class := range classes {
		if err := db.Create(&class).Error; err != nil {
			return fmt.Errorf("failed to create character class %s: %w", class.Name, err)
		}
	}
	return nil
}

// migrateSkillCategories migrates skill category data
func migrateSkillCategories(db *gorm.DB) error {
	categories := []SkillCategory{
		{GameSystem: "midgard", Name: "Alltag", Description: "Alltägliche Fertigkeiten", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Freiland", Description: "Freilandfertigkeiten", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Halbwelt", Description: "Halbweltfertigkeiten", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Kampf", Description: "Kampffertigkeiten", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Körper", Description: "Körperfertigkeiten", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Sozial", Description: "Sozialfertigkeiten", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Unterwelt", Description: "Unterweltfertigkeiten", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Waffen", Description: "Waffenfertigkeiten", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Wissen", Description: "Wissensfertigkeiten", Quelle: "Midgard5"},
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			return fmt.Errorf("failed to create skill category %s: %w", category.Name, err)
		}
	}
	return nil
}

// migrateSpellSchools migrates spell school data
func migrateSpellSchools(db *gorm.DB) error {
	schools := []SpellSchool{
		{GameSystem: "midgard", Name: "Beherrschen", Description: "Beherrschungsmagie", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Bewegen", Description: "Bewegungsmagie", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Dweomer", Description: "Dweomermagie", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Erkennen", Description: "Erkenntnismagie", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Erschaffen", Description: "Erschaffungsmagie", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Formen", Description: "Formungsmagie", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Verändern", Description: "Veränderungsmagie", Quelle: "Midgard5"},
		{GameSystem: "midgard", Name: "Zerstören", Description: "Zerstörungsmagie", Quelle: "Midgard5"},
	}

	for _, school := range schools {
		if err := db.Create(&school).Error; err != nil {
			return fmt.Errorf("failed to create spell school %s: %w", school.Name, err)
		}
	}
	return nil
}

// migrateSkillDifficulties migrates skill difficulty data
func migrateSkillDifficulties(db *gorm.DB) error {
	// First, we need to get the category IDs
	var categories []SkillCategory
	if err := db.Find(&categories).Error; err != nil {
		return fmt.Errorf("failed to retrieve skill categories: %w", err)
	}

	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	difficulties := []SkillDifficulty{
		{Name: "leicht", CategoryID: categoryMap["Alltag"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Alltag"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Alltag"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Alltag"], LearnCost: 4},
		{Name: "leicht", CategoryID: categoryMap["Freiland"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Freiland"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Freiland"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Freiland"], LearnCost: 4},
		{Name: "leicht", CategoryID: categoryMap["Halbwelt"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Halbwelt"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Halbwelt"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Halbwelt"], LearnCost: 4},
		{Name: "leicht", CategoryID: categoryMap["Kampf"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Kampf"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Kampf"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Kampf"], LearnCost: 4},
		{Name: "leicht", CategoryID: categoryMap["Körper"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Körper"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Körper"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Körper"], LearnCost: 4},
		{Name: "leicht", CategoryID: categoryMap["Sozial"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Sozial"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Sozial"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Sozial"], LearnCost: 4},
		{Name: "leicht", CategoryID: categoryMap["Unterwelt"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Unterwelt"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Unterwelt"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Unterwelt"], LearnCost: 4},
		{Name: "leicht", CategoryID: categoryMap["Waffen"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Waffen"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Waffen"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Waffen"], LearnCost: 4},
		{Name: "leicht", CategoryID: categoryMap["Wissen"], LearnCost: 1},
		{Name: "standard", CategoryID: categoryMap["Wissen"], LearnCost: 2},
		{Name: "schwer", CategoryID: categoryMap["Wissen"], LearnCost: 3},
		{Name: "sehr schwer", CategoryID: categoryMap["Wissen"], LearnCost: 4},
	}

	for _, difficulty := range difficulties {
		if err := db.Create(&difficulty).Error; err != nil {
			return fmt.Errorf("failed to create skill difficulty %s for category %d: %w", difficulty.Name, difficulty.CategoryID, err)
		}
	}
	return nil
}

// migrateClassCategoryEPCosts migrates EP cost data for class/category combinations
func migrateClassCategoryEPCosts(db *gorm.DB) error {
	// Get class and category IDs
	var classes []CharacterClass
	var categories []SkillCategory
	
	if err := db.Find(&classes).Error; err != nil {
		return fmt.Errorf("failed to retrieve character classes: %w", err)
	}
	if err := db.Find(&categories).Error; err != nil {
		return fmt.Errorf("failed to retrieve skill categories: %w", err)
	}

	classMap := make(map[string]uint)
	for _, class := range classes {
		classMap[class.Name] = class.ID
	}
	
	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	// Define EP costs based on learningCostsData
	epCosts := []ClassCategoryEPCost{
		// Krieger
		{ClassID: classMap["Krieger"], CategoryID: categoryMap["Alltag"], EPPerTE: 20},
		{ClassID: classMap["Krieger"], CategoryID: categoryMap["Freiland"], EPPerTE: 30},
		{ClassID: classMap["Krieger"], CategoryID: categoryMap["Halbwelt"], EPPerTE: 30},
		{ClassID: classMap["Krieger"], CategoryID: categoryMap["Kampf"], EPPerTE: 10},
		{ClassID: classMap["Krieger"], CategoryID: categoryMap["Körper"], EPPerTE: 20},
		{ClassID: classMap["Krieger"], CategoryID: categoryMap["Sozial"], EPPerTE: 20},
		{ClassID: classMap["Krieger"], CategoryID: categoryMap["Waffen"], EPPerTE: 5},
		{ClassID: classMap["Krieger"], CategoryID: categoryMap["Wissen"], EPPerTE: 30},
		
		// Assassine
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Alltag"], EPPerTE: 20},
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Freiland"], EPPerTE: 20},
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Halbwelt"], EPPerTE: 20},
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Kampf"], EPPerTE: 30},
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Körper"], EPPerTE: 10},
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Sozial"], EPPerTE: 20},
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Unterwelt"], EPPerTE: 10},
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Waffen"], EPPerTE: 20},
		{ClassID: classMap["Assassine"], CategoryID: categoryMap["Wissen"], EPPerTE: 20},
		
		// Magier
		{ClassID: classMap["Magier"], CategoryID: categoryMap["Alltag"], EPPerTE: 20},
		{ClassID: classMap["Magier"], CategoryID: categoryMap["Freiland"], EPPerTE: 30},
		{ClassID: classMap["Magier"], CategoryID: categoryMap["Kampf"], EPPerTE: 40},
		{ClassID: classMap["Magier"], CategoryID: categoryMap["Körper"], EPPerTE: 30},
		{ClassID: classMap["Magier"], CategoryID: categoryMap["Sozial"], EPPerTE: 20},
		{ClassID: classMap["Magier"], CategoryID: categoryMap["Waffen"], EPPerTE: 40},
		{ClassID: classMap["Magier"], CategoryID: categoryMap["Wissen"], EPPerTE: 10},
	}

	for _, cost := range epCosts {
		if err := db.Create(&cost).Error; err != nil {
			return fmt.Errorf("failed to create class category EP cost: %w", err)
		}
	}
	return nil
}

// migrateClassSpellSchoolCosts migrates spell EP cost data for class/school combinations
func migrateClassSpellSchoolCosts(db *gorm.DB) error {
	// Get class and school IDs
	var classes []CharacterClass
	var schools []SpellSchool
	
	if err := db.Find(&classes).Error; err != nil {
		return fmt.Errorf("failed to retrieve character classes: %w", err)
	}
	if err := db.Find(&schools).Error; err != nil {
		return fmt.Errorf("failed to retrieve spell schools: %w", err)
	}

	classMap := make(map[string]uint)
	for _, class := range classes {
		classMap[class.Name] = class.ID
	}
	
	schoolMap := make(map[string]uint)
	for _, school := range schools {
		schoolMap[school.Name] = school.ID
	}

	// Define spell EP costs based on learningCostsData
	spellCosts := []ClassSpellSchoolCost{
		// Magier - basic costs
		{ClassID: classMap["Magier"], SchoolID: schoolMap["Beherrschen"], EPPerLE: 10, IsSpecialization: false},
		{ClassID: classMap["Magier"], SchoolID: schoolMap["Bewegen"], EPPerLE: 10, IsSpecialization: false},
		{ClassID: classMap["Magier"], SchoolID: schoolMap["Dweomer"], EPPerLE: 15, IsSpecialization: false},
		{ClassID: classMap["Magier"], SchoolID: schoolMap["Erkennen"], EPPerLE: 10, IsSpecialization: false},
		{ClassID: classMap["Magier"], SchoolID: schoolMap["Erschaffen"], EPPerLE: 10, IsSpecialization: false},
		{ClassID: classMap["Magier"], SchoolID: schoolMap["Formen"], EPPerLE: 10, IsSpecialization: false},
		{ClassID: classMap["Magier"], SchoolID: schoolMap["Verändern"], EPPerLE: 10, IsSpecialization: false},
		{ClassID: classMap["Magier"], SchoolID: schoolMap["Zerstören"], EPPerLE: 10, IsSpecialization: false},
	}

	for _, cost := range spellCosts {
		if err := db.Create(&cost).Error; err != nil {
			return fmt.Errorf("failed to create class spell school cost: %w", err)
		}
	}
	return nil
}

// migrateSkillDefinitions migrates individual skill definitions from learningCostsData
func migrateSkillDefinitions(db *gorm.DB) error {
	// Get category and difficulty IDs
	var categories []SkillCategory
	var difficulties []SkillDifficulty
	
	if err := db.Find(&categories).Error; err != nil {
		return fmt.Errorf("failed to retrieve skill categories: %w", err)
	}
	if err := db.Find(&difficulties).Error; err != nil {
		return fmt.Errorf("failed to retrieve skill difficulties: %w", err)
	}

	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}
	
	difficultyMap := make(map[string]map[string]uint)
	for _, diff := range difficulties {
		catName := ""
		for _, cat := range categories {
			if cat.ID == diff.CategoryID {
				catName = cat.Name
				break
			}
		}
		if difficultyMap[catName] == nil {
			difficultyMap[catName] = make(map[string]uint)
		}
		difficultyMap[catName][diff.Name] = diff.ID
	}

	// Define skills based on learningCostsData
	skills := []SkillDefinitionNew{
		// Alltag
		{GameSystem: "midgard", Name: "Klettern", Quelle: "Midgard5", CategoryID: categoryMap["Alltag"], DifficultyID: difficultyMap["Alltag"]["leicht"]},
		{GameSystem: "midgard", Name: "Laufen", Quelle: "Midgard5", CategoryID: categoryMap["Alltag"], DifficultyID: difficultyMap["Alltag"]["leicht"]},
		{GameSystem: "midgard", Name: "Schwimmen", Quelle: "Midgard5", CategoryID: categoryMap["Alltag"], DifficultyID: difficultyMap["Alltag"]["standard"]},
		{GameSystem: "midgard", Name: "Springen", Quelle: "Midgard5", CategoryID: categoryMap["Alltag"], DifficultyID: difficultyMap["Alltag"]["standard"]},
		
		// Körper
		{GameSystem: "midgard", Name: "Athletik", Quelle: "Midgard5", CategoryID: categoryMap["Körper"], DifficultyID: difficultyMap["Körper"]["standard"]},
		{GameSystem: "midgard", Name: "Balancieren", Quelle: "Midgard5", CategoryID: categoryMap["Körper"], DifficultyID: difficultyMap["Körper"]["standard"]},
		{GameSystem: "midgard", Name: "Tauchen", Quelle: "Midgard5", CategoryID: categoryMap["Körper"], DifficultyID: difficultyMap["Körper"]["schwer"]},
		
		// Freiland
		{GameSystem: "midgard", Name: "Bogenschießen", Quelle: "Midgard5", CategoryID: categoryMap["Freiland"], DifficultyID: difficultyMap["Freiland"]["standard"]},
		{GameSystem: "midgard", Name: "Erste Hilfe", Quelle: "Midgard5", CategoryID: categoryMap["Freiland"], DifficultyID: difficultyMap["Freiland"]["standard"]},
		{GameSystem: "midgard", Name: "Reiten", Quelle: "Midgard5", CategoryID: categoryMap["Freiland"], DifficultyID: difficultyMap["Freiland"]["standard"]},
		{GameSystem: "midgard", Name: "Spurenlesen", Quelle: "Midgard5", CategoryID: categoryMap["Freiland"], DifficultyID: difficultyMap["Freiland"]["schwer"]},
		{GameSystem: "midgard", Name: "Überleben", Quelle: "Midgard5", CategoryID: categoryMap["Freiland"], DifficultyID: difficultyMap["Freiland"]["schwer"]},
		
		// Wissen
		{GameSystem: "midgard", Name: "Alchemie", Quelle: "Midgard5", CategoryID: categoryMap["Wissen"], DifficultyID: difficultyMap["Wissen"]["sehr schwer"]},
		{GameSystem: "midgard", Name: "Heilkunde", Quelle: "Midgard5", CategoryID: categoryMap["Wissen"], DifficultyID: difficultyMap["Wissen"]["sehr schwer"]},
		{GameSystem: "midgard", Name: "Naturkunde", Quelle: "Midgard5", CategoryID: categoryMap["Wissen"], DifficultyID: difficultyMap["Wissen"]["standard"]},
		{GameSystem: "midgard", Name: "Zauberkunde", Quelle: "Midgard5", CategoryID: categoryMap["Wissen"], DifficultyID: difficultyMap["Wissen"]["sehr schwer"]},
		
		// Waffen
		{GameSystem: "midgard", Name: "Dolch", Quelle: "Midgard5", CategoryID: categoryMap["Waffen"], DifficultyID: difficultyMap["Waffen"]["leicht"]},
		{GameSystem: "midgard", Name: "Einhandschwert", Quelle: "Midgard5", CategoryID: categoryMap["Waffen"], DifficultyID: difficultyMap["Waffen"]["standard"]},
		{GameSystem: "midgard", Name: "Zweihandschwert", Quelle: "Midgard5", CategoryID: categoryMap["Waffen"], DifficultyID: difficultyMap["Waffen"]["schwer"]},
		{GameSystem: "midgard", Name: "Wurfwaffen", Quelle: "Midgard5", CategoryID: categoryMap["Waffen"], DifficultyID: difficultyMap["Waffen"]["standard"]},
	}

	for _, skill := range skills {
		if err := db.Create(&skill).Error; err != nil {
			return fmt.Errorf("failed to create skill definition %s: %w", skill.Name, err)
		}
	}
	return nil
}

// migrateSkillImprovementCosts migrates skill improvement cost data
func migrateSkillImprovementCosts(db *gorm.DB) error {
	// Get difficulty IDs
	var difficulties []SkillDifficulty
	if err := db.Find(&difficulties).Error; err != nil {
		return fmt.Errorf("failed to retrieve skill difficulties: %w", err)
	}

	// Create improvement costs for each difficulty level and skill level
	var costs []SkillImprovementCost
	for _, diff := range difficulties {
		for level := 1; level <= 20; level++ {
			teCost := level * diff.LearnCost // Simple formula: level * base difficulty cost
			costs = append(costs, SkillImprovementCost{
				DifficultyID: diff.ID,
				CurrentLevel: level,
				TECost:       teCost,
			})
		}
	}

	for _, cost := range costs {
		if err := db.Create(&cost).Error; err != nil {
			return fmt.Errorf("failed to create skill improvement cost: %w", err)
		}
	}
	return nil
}

// migrateSpellLevelCosts migrates spell level cost data
func migrateSpellLevelCosts(db *gorm.DB) error {
	costs := []SpellLevelCost{
		{SpellLevel: 1, LERequired: 1},
		{SpellLevel: 2, LERequired: 2},
		{SpellLevel: 3, LERequired: 3},
		{SpellLevel: 4, LERequired: 4},
		{SpellLevel: 5, LERequired: 5},
		{SpellLevel: 6, LERequired: 6},
		{SpellLevel: 7, LERequired: 7},
		{SpellLevel: 8, LERequired: 8},
		{SpellLevel: 9, LERequired: 9},
		{SpellLevel: 10, LERequired: 10},
	}

	for _, cost := range costs {
		if err := db.Create(&cost).Error; err != nil {
			return fmt.Errorf("failed to create spell level cost for level %d: %w", cost.SpellLevel, err)
		}
	}
	return nil
}
