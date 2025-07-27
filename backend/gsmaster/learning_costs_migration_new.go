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

// migrateCharacterClasses migrates character class data directly from learningCostsData
func migrateCharacterClasses(db *gorm.DB) error {
	// Generate classes directly from learningCostsData.EPPerTE keys
	codeToClassInfo := map[string]struct {
		Name        string
		Description string
	}{
		"As": {"Assassine", "Assassine"},
		"Bb": {"Barbar", "Barbar"},
		"Gl": {"Glücksritter", "Glücksritter"},
		"Hä": {"Händler", "Händler"},
		"Kr": {"Krieger", "Krieger"},
		"Sp": {"Spitzbube", "Spitzbube"},
		"Wa": {"Waldläufer", "Waldläufer"},
		"Ba": {"Barde", "Barde"},
		"Or": {"Ordenskrieger", "Ordenskrieger"},
		"Dr": {"Druide", "Druide"},
		"Hx": {"Hexer", "Hexer"},
		"Ma": {"Magier", "Magier"},
		"PB": {"Priester Beschützer", "Priester Beschützer"},
		"PS": {"Priester Streiter", "Priester Streiter"},
		"Sc": {"Schamane", "Schamane"},
	}

	var classes []CharacterClass
	for classCode := range learningCostsData.EPPerTE {
		if info, exists := codeToClassInfo[classCode]; exists {
			classes = append(classes, CharacterClass{
				GameSystem:  "midgard",
				Name:        info.Name,
				Description: info.Description,
				Quelle:      "Midgard5",
				Code:        classCode,
			})
		}
	}

	for _, class := range classes {
		if err := db.Create(&class).Error; err != nil {
			return fmt.Errorf("failed to create character class %s: %w", class.Name, err)
		}
	}
	return nil
}

// migrateSkillCategories migrates skill category data directly from learningCostsData
func migrateSkillCategories(db *gorm.DB) error {
	// Generate categories from learningCostsData.ImprovementCost keys
	categoryDescriptions := map[string]string{
		"Alltag":    "Alltagsfertigkeiten",
		"Freiland":  "Freilandfertigkeiten",
		"Halbwelt":  "Halbweltfertigkeiten",
		"Kampf":     "Kampffertigkeiten",
		"Körper":    "Körperfertigkeiten",
		"Sozial":    "Sozialfertigkeiten",
		"Unterwelt": "Unterweltfertigkeiten",
		"Waffen":    "Waffenfertigkeiten",
		"Wissen":    "Wissensfertigkeiten",
	}

	var categories []SkillCategory
	for categoryName := range learningCostsData.ImprovementCost {
		if description, exists := categoryDescriptions[categoryName]; exists {
			categories = append(categories, SkillCategory{
				GameSystem:  "midgard",
				Name:        categoryName,
				Description: description,
				Quelle:      "Midgard5",
			})
		}
	}

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			return fmt.Errorf("failed to create skill category %s: %w", category.Name, err)
		}
	}
	return nil
}

// migrateSpellSchools migrates spell school data directly from learningCostsData
func migrateSpellSchools(db *gorm.DB) error {
	// Collect all unique spell schools from learningCostsData.SpellEPPerLE
	schoolDescriptions := map[string]string{
		"Beherrschen": "Beherrschungsmagie",
		"Bewegen":     "Bewegungsmagie",
		"Erkennen":    "Erkenntnismagie",
		"Erschaffen":  "Erschaffungsmagie",
		"Formen":      "Formungsmagie",
		"Verändern":   "Veränderungsmagie",
		"Zerstören":   "Zerstörungsmagie",
		"Wunder":      "Wundermagie",
		"Dweomer":     "Dweomermagie",
		"Lied":        "Liedmagie",
	}

	schoolsFound := make(map[string]bool)
	for _, schools := range learningCostsData.SpellEPPerLE {
		for schoolName := range schools {
			schoolsFound[schoolName] = true
		}
	}

	var schools []SpellSchool
	for schoolName := range schoolsFound {
		if description, exists := schoolDescriptions[schoolName]; exists {
			schools = append(schools, SpellSchool{
				GameSystem:  "midgard",
				Name:        schoolName,
				Description: description,
				Quelle:      "Midgard5",
			})
		}
	}

	for _, school := range schools {
		if err := db.Create(&school).Error; err != nil {
			return fmt.Errorf("failed to create spell school %s: %w", school.Name, err)
		}
	}
	return nil
}

// migrateSkillDifficulties migrates skill difficulty data directly from learningCostsData
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

	// Generate difficulties from learningCostsData.ImprovementCost
	var difficulties []SkillDifficulty
	for categoryName, diffData := range learningCostsData.ImprovementCost {
		categoryID, exists := categoryMap[categoryName]
		if !exists {
			continue
		}

		for difficultyName, data := range diffData {
			difficulties = append(difficulties, SkillDifficulty{
				Name:       difficultyName,
				CategoryID: categoryID,
				LearnCost:  data.LearnCost,
			})
		}
	}

	for _, difficulty := range difficulties {
		if err := db.Create(&difficulty).Error; err != nil {
			return fmt.Errorf("failed to create skill difficulty %s for category %d: %w", difficulty.Name, difficulty.CategoryID, err)
		}
	}
	return nil
}

// migrateClassCategoryEPCosts migrates EP cost data directly from learningCostsData
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

	// Generate EP costs directly from learningCostsData
	var epCosts []ClassCategoryEPCost

	// Mapping from code to class name for lookup
	codeToClassName := map[string]string{
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

	// Iterate through all entries in learningCostsData.EPPerTE
	for classCode, categories := range learningCostsData.EPPerTE {
		className, exists := codeToClassName[classCode]
		if !exists {
			continue // Skip unknown class codes
		}

		classID, exists := classMap[className]
		if !exists {
			continue // Skip if class not found in database
		}

		for categoryName, epValue := range categories {
			categoryID, exists := categoryMap[categoryName]
			if !exists {
				continue // Skip if category not found in database
			}

			epCosts = append(epCosts, ClassCategoryEPCost{
				ClassID:    classID,
				CategoryID: categoryID,
				EPPerTE:    epValue,
			})
		}
	}

	for _, cost := range epCosts {
		if err := db.Create(&cost).Error; err != nil {
			return fmt.Errorf("failed to create class category EP cost: %w", err)
		}
	}
	return nil
}

// migrateClassSpellSchoolCosts migrates spell EP cost data directly from learningCostsData
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

	// Generate spell EP costs directly from learningCostsData.SpellEPPerLE
	var spellCosts []ClassSpellSchoolCost

	// Mapping from code to class name for lookup
	codeToClassName := map[string]string{
		"Dr": "Druide",
		"Hx": "Hexer",
		"Ma": "Magier",
		"PB": "Priester Beschützer",
		"PS": "Priester Streiter",
		"Sc": "Schamane",
		"Ba": "Barde",
	}

	// Iterate through all entries in learningCostsData.SpellEPPerLE
	for classCode, schools := range learningCostsData.SpellEPPerLE {
		className, exists := codeToClassName[classCode]
		if !exists {
			continue // Skip unknown class codes
		}

		classID, exists := classMap[className]
		if !exists {
			continue // Skip if class not found in database
		}

		for schoolName, epValue := range schools {
			// Skip schools with 0 EP cost (not available)
			if epValue == 0 {
				continue
			}

			schoolID, exists := schoolMap[schoolName]
			if !exists {
				continue // Skip if school not found in database
			}

			// For Magier, check if this is a specialization (30 EP cost)
			isSpecialization := (classCode == "Ma" && epValue == 30)

			spellCosts = append(spellCosts, ClassSpellSchoolCost{
				ClassID:          classID,
				SchoolID:         schoolID,
				EPPerLE:          epValue,
				IsSpecialization: isSpecialization,
			})
		}
	}

	for _, cost := range spellCosts {
		if err := db.Create(&cost).Error; err != nil {
			return fmt.Errorf("failed to create class spell school cost: %w", err)
		}
	}
	return nil
}

// migrateSkillDefinitions migrates individual skill definitions directly from learningCostsData
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

	// Generate skills directly from learningCostsData.ImprovementCost
	var skills []SkillDefinitionNew
	for categoryName, difficulties := range learningCostsData.ImprovementCost {
		categoryID, exists := categoryMap[categoryName]
		if !exists {
			continue
		}

		for difficultyName, data := range difficulties {
			difficultyID, exists := difficultyMap[categoryName][difficultyName]
			if !exists {
				continue
			}

			for _, skillName := range data.Skills {
				skills = append(skills, SkillDefinitionNew{
					GameSystem:   "midgard",
					Name:         skillName,
					Quelle:       "Midgard5",
					CategoryID:   categoryID,
					DifficultyID: difficultyID,
				})
			}
		}
	}

	for _, skill := range skills {
		if err := db.Create(&skill).Error; err != nil {
			return fmt.Errorf("failed to create skill definition %s: %w", skill.Name, err)
		}
	}
	return nil
}

// migrateSkillImprovementCosts migrates skill improvement cost data directly from learningCostsData
func migrateSkillImprovementCosts(db *gorm.DB) error {
	// Get difficulty IDs
	var difficulties []SkillDifficulty
	if err := db.Find(&difficulties).Error; err != nil {
		return fmt.Errorf("failed to retrieve skill difficulties: %w", err)
	}

	// Create improvement costs from learningCostsData.ImprovementCost
	var costs []SkillImprovementCost
	for _, diff := range difficulties {
		// Find the corresponding difficulty data in learningCostsData
		var diffData *DifficultyData
		for _, catDifficulties := range learningCostsData.ImprovementCost {
			for diffName, data := range catDifficulties {
				if diffName == diff.Name {
					diffData = &data
					break
				}
			}
			if diffData != nil {
				break
			}
		}

		if diffData == nil {
			continue
		}

		// Use actual train costs from learningCostsData if available
		if len(diffData.TrainCosts) > 0 {
			for level, teCost := range diffData.TrainCosts {
				costs = append(costs, SkillImprovementCost{
					DifficultyID: diff.ID,
					CurrentLevel: level,
					TECost:       teCost,
				})
			}
		} else {
			// Fallback to simple formula for levels without explicit train costs
			for level := 1; level <= 20; level++ {
				teCost := level * diff.LearnCost
				costs = append(costs, SkillImprovementCost{
					DifficultyID: diff.ID,
					CurrentLevel: level,
					TECost:       teCost,
				})
			}
		}
	}

	for _, cost := range costs {
		if err := db.Create(&cost).Error; err != nil {
			return fmt.Errorf("failed to create skill improvement cost: %w", err)
		}
	}
	return nil
}

// migrateSpellLevelCosts migrates spell level cost data directly from learningCostsData
func migrateSpellLevelCosts(db *gorm.DB) error {
	// Generate spell level costs directly from learningCostsData.SpellLEPerLevel
	var costs []SpellLevelCost

	for level, leRequired := range learningCostsData.SpellLEPerLevel {
		costs = append(costs, SpellLevelCost{
			SpellLevel: level,
			LERequired: leRequired,
		})
	}

	for _, cost := range costs {
		if err := db.Create(&cost).Error; err != nil {
			return fmt.Errorf("failed to create spell level cost for level %d: %w", cost.SpellLevel, err)
		}
	}
	return nil
}
