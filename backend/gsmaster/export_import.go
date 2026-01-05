package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

// ExportableCategoryDifficulty represents a category/difficulty combination for a skill
type ExportableCategoryDifficulty struct {
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
	LearnCost  int    `json:"learn_cost"`
}

// ExportableSkill represents a skill without database IDs for export
type ExportableSkill struct {
	Name                   string                         `json:"name"`
	GameSystem             string                         `json:"game_system"`
	Beschreibung           string                         `json:"beschreibung"`
	SourceCode             string                         `json:"source_code"` // Instead of SourceID
	PageNumber             int                            `json:"page_number"`
	Initialwert            int                            `json:"initialwert"`
	BasisWert              int                            `json:"basiswert"`
	Bonuseigenschaft       string                         `json:"bonuseigenschaft"`
	Improvable             bool                           `json:"improvable"`
	InnateSkill            bool                           `json:"innate_skill"`
	Category               string                         `json:"category"`                // Deprecated: use CategoriesDifficulties
	Difficulty             string                         `json:"difficulty"`              // Deprecated: use CategoriesDifficulties
	CategoriesDifficulties []ExportableCategoryDifficulty `json:"categories_difficulties"` // All category/difficulty combinations
}

// ExportableSource represents a source for export
type ExportableSource struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Edition     string `json:"edition"`
	Publisher   string `json:"publisher"`
	PublishYear int    `json:"publish_year"`
	Description string `json:"description"`
	IsCore      bool   `json:"is_core"`
	IsActive    bool   `json:"is_active"`
	GameSystem  string `json:"game_system"`
}

// ExportableSkillCategory represents a skill category for export
type ExportableSkillCategory struct {
	Name       string `json:"name"`
	GameSystem string `json:"game_system"`
	SourceCode string `json:"source_code"`
}

// ExportableSkillDifficulty represents a skill difficulty for export
type ExportableSkillDifficulty struct {
	Name       string `json:"name"`
	GameSystem string `json:"game_system"`
}

// ExportableSkillCategoryDifficulty represents the relationship for export
type ExportableSkillCategoryDifficulty struct {
	SkillName        string `json:"skill_name"`
	SkillSystem      string `json:"skill_system"`
	CategoryName     string `json:"category_name"`
	CategorySystem   string `json:"category_system"`
	DifficultyName   string `json:"difficulty_name"`
	DifficultySystem string `json:"difficulty_system"`
	LearnCost        int    `json:"learn_cost"`
}

// ExportableSpell represents a spell for export
type ExportableSpell struct {
	Name             string `json:"name"`
	GameSystem       string `json:"game_system"`
	Beschreibung     string `json:"beschreibung"`
	SourceCode       string `json:"source_code"`
	PageNumber       int    `json:"page_number"`
	Bonus            int    `json:"bonus"`
	Stufe            int    `json:"level"`
	AP               string `json:"ap"`
	Art              string `json:"art"`
	Zauberdauer      string `json:"zauberdauer"`
	Reichweite       string `json:"reichweite"`
	Wirkungsziel     string `json:"wirkungsziel"`
	Wirkungsbereich  string `json:"wirkungsbereich"`
	Wirkungsdauer    string `json:"wirkungsdauer"`
	Ursprung         string `json:"ursprung"`
	Category         string `json:"category"`
	LearningCategory string `json:"learning_category"`
}

// ExportableCharacterClass represents a character class for export
type ExportableCharacterClass struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SourceCode  string `json:"source_code"`
	GameSystem  string `json:"game_system"`
}

// ExportableSpellSchool represents a spell school for export
type ExportableSpellSchool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SourceCode  string `json:"source_code"`
	GameSystem  string `json:"game_system"`
}

// ExportableClassCategoryEPCost represents class-category EP costs for export
type ExportableClassCategoryEPCost struct {
	CharacterClassCode string `json:"character_class_code"`
	SkillCategoryName  string `json:"skill_category_name"`
	EPPerTE            int    `json:"ep_per_te"`
}

// ExportableClassSpellSchoolEPCost represents class-spell school EP costs for export
type ExportableClassSpellSchoolEPCost struct {
	CharacterClassCode string `json:"character_class_code"`
	SpellSchoolName    string `json:"spell_school_name"`
	EPPerLE            int    `json:"ep_per_le"`
}

// ExportableSpellLevelLECost represents spell level LE costs for export
type ExportableSpellLevelLECost struct {
	Level      int    `json:"level"`
	LERequired int    `json:"le_required"`
	GameSystem string `json:"game_system"`
}

// ExportableSkillImprovementCost represents skill improvement costs for export
type ExportableSkillImprovementCost struct {
	SkillName        string `json:"skill_name"`
	SkillSystem      string `json:"skill_system"`
	CategoryName     string `json:"category_name"`
	CategorySystem   string `json:"category_system"`
	DifficultyName   string `json:"difficulty_name"`
	DifficultySystem string `json:"difficulty_system"`
	CurrentLevel     int    `json:"current_level"`
	TERequired       int    `json:"te_required"`
}

// ExportableWeaponSkill represents a weapon skill for export
type ExportableWeaponSkill struct {
	Name             string `json:"name"`
	GameSystem       string `json:"game_system"`
	Beschreibung     string `json:"beschreibung"`
	SourceCode       string `json:"source_code"`
	PageNumber       int    `json:"page_number"`
	Initialwert      int    `json:"initialwert"`
	BasisWert        int    `json:"basiswert"`
	Bonuseigenschaft string `json:"bonuseigenschaft"`
	Improvable       bool   `json:"improvable"`
	InnateSkill      bool   `json:"innate_skill"`
	Category         string `json:"category"`
	Difficulty       string `json:"difficulty"`
}

// ExportableEquipment represents equipment for export
type ExportableEquipment struct {
	Name         string  `json:"name"`
	GameSystem   string  `json:"game_system"`
	Beschreibung string  `json:"beschreibung"`
	SourceCode   string  `json:"source_code"`
	PageNumber   int     `json:"page_number"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
	PersonalItem bool    `json:"personal_item"`
}

// ExportableWeapon represents a weapon for export
type ExportableWeapon struct {
	Name          string  `json:"name"`
	GameSystem    string  `json:"game_system"`
	Beschreibung  string  `json:"beschreibung"`
	SourceCode    string  `json:"source_code"`
	PageNumber    int     `json:"page_number"`
	Gewicht       float64 `json:"gewicht"`
	Wert          float64 `json:"wert"`
	PersonalItem  bool    `json:"personal_item"`
	SkillRequired string  `json:"skill_required"`
	Damage        string  `json:"damage"`
	RangeNear     int     `json:"range_near"`
	RangeMiddle   int     `json:"range_middle"`
	RangeFar      int     `json:"range_far"`
}

// ExportableContainer represents a container for export
type ExportableContainer struct {
	Name         string  `json:"name"`
	GameSystem   string  `json:"game_system"`
	Beschreibung string  `json:"beschreibung"`
	SourceCode   string  `json:"source_code"`
	PageNumber   int     `json:"page_number"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
	PersonalItem bool    `json:"personal_item"`
	Tragkraft    float64 `json:"tragkraft"`
	Volumen      float64 `json:"volumen"`
}

// ExportableTransportation represents transportation for export
type ExportableTransportation struct {
	Name         string  `json:"name"`
	GameSystem   string  `json:"game_system"`
	Beschreibung string  `json:"beschreibung"`
	SourceCode   string  `json:"source_code"`
	PageNumber   int     `json:"page_number"`
	Gewicht      float64 `json:"gewicht"`
	Wert         float64 `json:"wert"`
	PersonalItem bool    `json:"personal_item"`
	Tragkraft    float64 `json:"tragkraft"`
	Volumen      float64 `json:"volumen"`
}

// ExportableBelieve represents a belief system for export
type ExportableBelieve struct {
	Name         string `json:"name"`
	GameSystem   string `json:"game_system"`
	Beschreibung string `json:"beschreibung"`
	SourceCode   string `json:"source_code"`
	PageNumber   int    `json:"page_number"`
}

// ExportSkills exports all skills to a JSON file
func ExportSkills(outputDir string) error {
	var skills []models.Skill
	if err := database.DB.Find(&skills).Error; err != nil {
		return fmt.Errorf("failed to fetch skills: %w", err)
	}

	// Get all sources for mapping
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[uint]string)
	for _, s := range sources {
		sourceMap[s.ID] = s.Code
	}

	// Get all skill category difficulties
	var scds []models.SkillCategoryDifficulty
	database.DB.Preload("SkillCategory").Preload("SkillDifficulty").Find(&scds)

	// Build map of skill_id -> []category/difficulty combinations
	scdMap := make(map[uint][]ExportableCategoryDifficulty)
	for _, scd := range scds {
		scdMap[scd.SkillID] = append(scdMap[scd.SkillID], ExportableCategoryDifficulty{
			Category:   scd.SkillCategory.Name,
			Difficulty: scd.SkillDifficulty.Name,
			LearnCost:  scd.LearnCost,
		})
	}

	// Convert to exportable format
	exportable := make([]ExportableSkill, len(skills))
	for i, skill := range skills {
		exportable[i] = ExportableSkill{
			Name:                   skill.Name,
			GameSystem:             skill.GameSystem,
			Beschreibung:           skill.Beschreibung,
			SourceCode:             sourceMap[skill.SourceID],
			PageNumber:             skill.PageNumber,
			Initialwert:            skill.Initialwert,
			BasisWert:              skill.BasisWert,
			Bonuseigenschaft:       skill.Bonuseigenschaft,
			Improvable:             skill.Improvable,
			InnateSkill:            skill.InnateSkill,
			Category:               skill.Category,
			Difficulty:             skill.Difficulty,
			CategoriesDifficulties: scdMap[skill.ID],
		}
	}

	return writeJSON(filepath.Join(outputDir, "skills.json"), exportable)
}

// ImportSkills imports skills from a JSON file
func ImportSkills(inputDir string) error {
	var exportable []ExportableSkill
	if err := readJSON(filepath.Join(inputDir, "skills.json"), &exportable); err != nil {
		return err
	}

	// Get all sources for mapping
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[string]uint)
	for _, s := range sources {
		sourceMap[s.Code] = s.ID
	}

	// Get all categories for mapping
	var categories []models.SkillCategory
	database.DB.Find(&categories)
	categoryMap := make(map[string]map[string]uint) // game_system -> name -> id
	for _, c := range categories {
		if categoryMap[c.GameSystem] == nil {
			categoryMap[c.GameSystem] = make(map[string]uint)
		}
		categoryMap[c.GameSystem][c.Name] = c.ID
	}

	// Get all difficulties for mapping
	var difficulties []models.SkillDifficulty
	database.DB.Find(&difficulties)
	difficultyMap := make(map[string]map[string]uint) // game_system -> name -> id
	for _, d := range difficulties {
		if difficultyMap[d.GameSystem] == nil {
			difficultyMap[d.GameSystem] = make(map[string]uint)
		}
		difficultyMap[d.GameSystem][d.Name] = d.ID
	}

	for _, exp := range exportable {
		var skill models.Skill
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&skill)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			// Create new skill
			skill = models.Skill{
				Name:             exp.Name,
				GameSystem:       exp.GameSystem,
				Beschreibung:     exp.Beschreibung,
				SourceID:         sourceID,
				PageNumber:       exp.PageNumber,
				Initialwert:      exp.Initialwert,
				BasisWert:        exp.BasisWert,
				Bonuseigenschaft: exp.Bonuseigenschaft,
				Improvable:       exp.Improvable,
				InnateSkill:      exp.InnateSkill,
				Category:         exp.Category,
				Difficulty:       exp.Difficulty,
			}
			if err := database.DB.Create(&skill).Error; err != nil {
				return fmt.Errorf("failed to create skill %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query skill %s: %w", exp.Name, result.Error)
		} else {
			// Update existing skill
			skill.Beschreibung = exp.Beschreibung
			skill.SourceID = sourceID
			skill.PageNumber = exp.PageNumber
			skill.Initialwert = exp.Initialwert
			skill.BasisWert = exp.BasisWert
			skill.Bonuseigenschaft = exp.Bonuseigenschaft
			skill.Improvable = exp.Improvable
			skill.InnateSkill = exp.InnateSkill
			skill.Category = exp.Category
			skill.Difficulty = exp.Difficulty

			if err := database.DB.Save(&skill).Error; err != nil {
				return fmt.Errorf("failed to update skill %s: %w", exp.Name, err)
			}
		}

		// Import category/difficulty combinations if present
		if len(exp.CategoriesDifficulties) > 0 {
			// Delete existing relationships for this skill
			database.DB.Where("skill_id = ?", skill.ID).Delete(&models.SkillCategoryDifficulty{})

			// Create new relationships
			for _, cd := range exp.CategoriesDifficulties {
				categoryID := categoryMap[exp.GameSystem][cd.Category]
				difficultyID := difficultyMap[exp.GameSystem][cd.Difficulty]

				if categoryID == 0 || difficultyID == 0 {
					continue // Skip if category or difficulty not found
				}

				scd := models.SkillCategoryDifficulty{
					SkillID:           skill.ID,
					SkillCategoryID:   categoryID,
					SkillDifficultyID: difficultyID,
					LearnCost:         cd.LearnCost,
					SCategory:         cd.Category,
					SDifficulty:       cd.Difficulty,
				}
				if err := database.DB.Create(&scd).Error; err != nil {
					return fmt.Errorf("failed to create skill category difficulty for %s: %w", exp.Name, err)
				}
			}
		}
	}

	return nil
}

// ExportSources exports all sources to a JSON file
func ExportSources(outputDir string) error {
	var sources []models.Source
	if err := database.DB.Find(&sources).Error; err != nil {
		return fmt.Errorf("failed to fetch sources: %w", err)
	}

	exportable := make([]ExportableSource, len(sources))
	for i, s := range sources {
		exportable[i] = ExportableSource{
			Code:        s.Code,
			Name:        s.Name,
			FullName:    s.FullName,
			Edition:     s.Edition,
			Publisher:   s.Publisher,
			PublishYear: s.PublishYear,
			Description: s.Description,
			IsCore:      s.IsCore,
			IsActive:    s.IsActive,
			GameSystem:  s.GameSystem,
		}
	}

	return writeJSON(filepath.Join(outputDir, "sources.json"), exportable)
}

// ImportSources imports sources from a JSON file
func ImportSources(inputDir string) error {
	var exportable []ExportableSource
	if err := readJSON(filepath.Join(inputDir, "sources.json"), &exportable); err != nil {
		return err
	}

	for _, exp := range exportable {
		var source models.Source
		result := database.DB.Where("code = ?", exp.Code).First(&source)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new source
			source = models.Source{
				Code:        exp.Code,
				Name:        exp.Name,
				FullName:    exp.FullName,
				Edition:     exp.Edition,
				Publisher:   exp.Publisher,
				PublishYear: exp.PublishYear,
				Description: exp.Description,
				IsCore:      exp.IsCore,
				IsActive:    exp.IsActive,
				GameSystem:  exp.GameSystem,
			}
			if err := database.DB.Create(&source).Error; err != nil {
				return fmt.Errorf("failed to create source %s: %w", exp.Code, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query source %s: %w", exp.Code, result.Error)
		} else {
			// Update existing source
			source.Name = exp.Name
			source.FullName = exp.FullName
			source.Edition = exp.Edition
			source.Publisher = exp.Publisher
			source.PublishYear = exp.PublishYear
			source.Description = exp.Description
			source.IsCore = exp.IsCore
			source.IsActive = exp.IsActive
			source.GameSystem = exp.GameSystem

			if err := database.DB.Save(&source).Error; err != nil {
				return fmt.Errorf("failed to update source %s: %w", exp.Code, err)
			}
		}
	}

	return nil
}

// ExportSkillCategoryDifficulties exports skill-category-difficulty relationships
func ExportSkillCategoryDifficulties(outputDir string) error {
	var scds []models.SkillCategoryDifficulty
	if err := database.DB.Find(&scds).Error; err != nil {
		return fmt.Errorf("failed to fetch skill category difficulties: %w", err)
	}

	exportable := make([]ExportableSkillCategoryDifficulty, len(scds))
	for i, scd := range scds {
		var skill models.Skill
		var category models.SkillCategory
		var difficulty models.SkillDifficulty

		database.DB.First(&skill, scd.SkillID)
		database.DB.First(&category, scd.SkillCategoryID)
		database.DB.First(&difficulty, scd.SkillDifficultyID)

		exportable[i] = ExportableSkillCategoryDifficulty{
			SkillName:        skill.Name,
			SkillSystem:      skill.GameSystem,
			CategoryName:     category.Name,
			CategorySystem:   category.GameSystem,
			DifficultyName:   difficulty.Name,
			DifficultySystem: difficulty.GameSystem,
			LearnCost:        scd.LearnCost,
		}
	}

	return writeJSON(filepath.Join(outputDir, "skill_category_difficulties.json"), exportable)
}

// ImportSkillCategoryDifficulties imports skill-category-difficulty relationships
func ImportSkillCategoryDifficulties(inputDir string) error {
	var exportable []ExportableSkillCategoryDifficulty
	if err := readJSON(filepath.Join(inputDir, "skill_category_difficulties.json"), &exportable); err != nil {
		return err
	}

	for _, exp := range exportable {
		// Find the skill
		var skill models.Skill
		if err := database.DB.Where("name = ? AND game_system = ?", exp.SkillName, exp.SkillSystem).First(&skill).Error; err != nil {
			return fmt.Errorf("skill not found: %s/%s", exp.SkillName, exp.SkillSystem)
		}

		// Find the category
		var category models.SkillCategory
		if err := database.DB.Where("name = ? AND game_system = ?", exp.CategoryName, exp.CategorySystem).First(&category).Error; err != nil {
			return fmt.Errorf("category not found: %s/%s", exp.CategoryName, exp.CategorySystem)
		}

		// Find the difficulty
		var difficulty models.SkillDifficulty
		if err := database.DB.Where("name = ? AND game_system = ?", exp.DifficultyName, exp.DifficultySystem).First(&difficulty).Error; err != nil {
			return fmt.Errorf("difficulty not found: %s/%s", exp.DifficultyName, exp.DifficultySystem)
		}

		// Check if relationship exists
		var scd models.SkillCategoryDifficulty
		result := database.DB.Where("skill_id = ? AND skill_category_id = ? AND skill_difficulty_id = ?",
			skill.ID, category.ID, difficulty.ID).First(&scd)

		if result.Error == gorm.ErrRecordNotFound {
			// Create new relationship
			scd = models.SkillCategoryDifficulty{
				SkillID:           skill.ID,
				SkillCategoryID:   category.ID,
				SkillDifficultyID: difficulty.ID,
				LearnCost:         exp.LearnCost,
				SCategory:         category.Name,
				SDifficulty:       difficulty.Name,
			}
			if err := database.DB.Create(&scd).Error; err != nil {
				return fmt.Errorf("failed to create relationship: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query relationship: %w", result.Error)
		} else {
			// Update existing relationship
			scd.LearnCost = exp.LearnCost
			scd.SCategory = category.Name
			scd.SDifficulty = difficulty.Name

			if err := database.DB.Save(&scd).Error; err != nil {
				return fmt.Errorf("failed to update relationship: %w", err)
			}
		}
	}

	return nil
}

// Helper functions for JSON I/O
func writeJSON(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON to %s: %w", filename, err)
	}

	return nil
}

func readJSON(filename string, data interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode JSON from %s: %w", filename, err)
	}

	return nil
}

// ExportSkillCategories exports all skill categories to a JSON file
func ExportSkillCategories(outputDir string) error {
	var categories []models.SkillCategory
	if err := database.DB.Find(&categories).Error; err != nil {
		return fmt.Errorf("failed to fetch skill categories: %w", err)
	}

	// Get source map
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[uint]string)
	for _, s := range sources {
		sourceMap[s.ID] = s.Code
	}

	exportable := make([]ExportableSkillCategory, len(categories))
	for i, cat := range categories {
		exportable[i] = ExportableSkillCategory{
			Name:       cat.Name,
			GameSystem: cat.GameSystem,
			SourceCode: sourceMap[cat.SourceID],
		}
	}

	return writeJSON(filepath.Join(outputDir, "skill_categories.json"), exportable)
}

// ImportSkillCategories imports skill categories from a JSON file
func ImportSkillCategories(inputDir string) error {
	var exportable []ExportableSkillCategory
	if err := readJSON(filepath.Join(inputDir, "skill_categories.json"), &exportable); err != nil {
		return err
	}

	// Get source map
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[string]uint)
	for _, s := range sources {
		sourceMap[s.Code] = s.ID
	}

	for _, exp := range exportable {
		var category models.SkillCategory
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&category)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			category = models.SkillCategory{
				Name:       exp.Name,
				GameSystem: exp.GameSystem,
				SourceID:   sourceID,
			}
			if err := database.DB.Create(&category).Error; err != nil {
				return fmt.Errorf("failed to create category %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query category %s: %w", exp.Name, result.Error)
		} else {
			category.SourceID = sourceID
			if err := database.DB.Save(&category).Error; err != nil {
				return fmt.Errorf("failed to update category %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}

// ExportSkillDifficulties exports all skill difficulties to a JSON file
func ExportSkillDifficulties(outputDir string) error {
	var difficulties []models.SkillDifficulty
	if err := database.DB.Find(&difficulties).Error; err != nil {
		return fmt.Errorf("failed to fetch skill difficulties: %w", err)
	}

	exportable := make([]ExportableSkillDifficulty, len(difficulties))
	for i, diff := range difficulties {
		exportable[i] = ExportableSkillDifficulty{
			Name:       diff.Name,
			GameSystem: diff.GameSystem,
		}
	}

	return writeJSON(filepath.Join(outputDir, "skill_difficulties.json"), exportable)
}

// ImportSkillDifficulties imports skill difficulties from a JSON file
func ImportSkillDifficulties(inputDir string) error {
	var exportable []ExportableSkillDifficulty
	if err := readJSON(filepath.Join(inputDir, "skill_difficulties.json"), &exportable); err != nil {
		return err
	}

	for _, exp := range exportable {
		var difficulty models.SkillDifficulty
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&difficulty)

		if result.Error == gorm.ErrRecordNotFound {
			difficulty = models.SkillDifficulty{
				Name:       exp.Name,
				GameSystem: exp.GameSystem,
			}
			if err := database.DB.Create(&difficulty).Error; err != nil {
				return fmt.Errorf("failed to create difficulty %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query difficulty %s: %w", exp.Name, result.Error)
		}
		// No update needed for difficulties (only name and game_system)
	}

	return nil
}

// ExportSpells exports all spells to a JSON file
func ExportSpells(outputDir string) error {
	var spells []models.Spell
	if err := database.DB.Find(&spells).Error; err != nil {
		return fmt.Errorf("failed to fetch spells: %w", err)
	}

	// Get source map
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[uint]string)
	for _, s := range sources {
		sourceMap[s.ID] = s.Code
	}

	exportable := make([]ExportableSpell, len(spells))
	for i, spell := range spells {
		exportable[i] = ExportableSpell{
			Name:             spell.Name,
			GameSystem:       spell.GameSystem,
			Beschreibung:     spell.Beschreibung,
			SourceCode:       sourceMap[spell.SourceID],
			PageNumber:       spell.PageNumber,
			Bonus:            spell.Bonus,
			Stufe:            spell.Stufe,
			AP:               spell.AP,
			Art:              spell.Art,
			Zauberdauer:      spell.Zauberdauer,
			Reichweite:       spell.Reichweite,
			Wirkungsziel:     spell.Wirkungsziel,
			Wirkungsbereich:  spell.Wirkungsbereich,
			Wirkungsdauer:    spell.Wirkungsdauer,
			Ursprung:         spell.Ursprung,
			Category:         spell.Category,
			LearningCategory: spell.LearningCategory,
		}
	}

	return writeJSON(filepath.Join(outputDir, "spells.json"), exportable)
}

// ImportSpells imports spells from a JSON file
func ImportSpells(inputDir string) error {
	var exportable []ExportableSpell
	if err := readJSON(filepath.Join(inputDir, "spells.json"), &exportable); err != nil {
		return err
	}

	// Get source map
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[string]uint)
	for _, s := range sources {
		sourceMap[s.Code] = s.ID
	}

	for _, exp := range exportable {
		var spell models.Spell
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&spell)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			spell = models.Spell{
				Name:             exp.Name,
				GameSystem:       exp.GameSystem,
				Beschreibung:     exp.Beschreibung,
				SourceID:         sourceID,
				PageNumber:       exp.PageNumber,
				Bonus:            exp.Bonus,
				Stufe:            exp.Stufe,
				AP:               exp.AP,
				Art:              exp.Art,
				Zauberdauer:      exp.Zauberdauer,
				Reichweite:       exp.Reichweite,
				Wirkungsziel:     exp.Wirkungsziel,
				Wirkungsbereich:  exp.Wirkungsbereich,
				Wirkungsdauer:    exp.Wirkungsdauer,
				Ursprung:         exp.Ursprung,
				Category:         exp.Category,
				LearningCategory: exp.LearningCategory,
			}
			if err := database.DB.Create(&spell).Error; err != nil {
				return fmt.Errorf("failed to create spell %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query spell %s: %w", exp.Name, result.Error)
		} else {
			// Update existing spell
			spell.Beschreibung = exp.Beschreibung
			spell.SourceID = sourceID
			spell.PageNumber = exp.PageNumber
			spell.Bonus = exp.Bonus
			spell.Stufe = exp.Stufe
			spell.AP = exp.AP
			spell.Art = exp.Art
			spell.Zauberdauer = exp.Zauberdauer
			spell.Reichweite = exp.Reichweite
			spell.Wirkungsziel = exp.Wirkungsziel
			spell.Wirkungsbereich = exp.Wirkungsbereich
			spell.Wirkungsdauer = exp.Wirkungsdauer
			spell.Ursprung = exp.Ursprung
			spell.Category = exp.Category
			spell.LearningCategory = exp.LearningCategory

			if err := database.DB.Save(&spell).Error; err != nil {
				return fmt.Errorf("failed to update spell %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}

// ExportCharacterClasses exports all character classes to a JSON file
func ExportCharacterClasses(outputDir string) error {
	var classes []models.CharacterClass
	if err := database.DB.Find(&classes).Error; err != nil {
		return fmt.Errorf("failed to fetch character classes: %w", err)
	}

	// Get source map
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[uint]string)
	for _, s := range sources {
		sourceMap[s.ID] = s.Code
	}

	exportable := make([]ExportableCharacterClass, len(classes))
	for i, class := range classes {
		exportable[i] = ExportableCharacterClass{
			Code:        class.Code,
			Name:        class.Name,
			Description: class.Description,
			SourceCode:  sourceMap[class.SourceID],
			GameSystem:  class.GameSystem,
		}
	}

	return writeJSON(filepath.Join(outputDir, "character_classes.json"), exportable)
}

// ImportCharacterClasses imports character classes from a JSON file
func ImportCharacterClasses(inputDir string) error {
	var exportable []ExportableCharacterClass
	if err := readJSON(filepath.Join(inputDir, "character_classes.json"), &exportable); err != nil {
		return err
	}

	// Get source map
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[string]uint)
	for _, s := range sources {
		sourceMap[s.Code] = s.ID
	}

	for _, exp := range exportable {
		var class models.CharacterClass
		result := database.DB.Where("code = ?", exp.Code).First(&class)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			class = models.CharacterClass{
				Code:        exp.Code,
				Name:        exp.Name,
				Description: exp.Description,
				SourceID:    sourceID,
				GameSystem:  exp.GameSystem,
			}
			if err := database.DB.Create(&class).Error; err != nil {
				return fmt.Errorf("failed to create character class %s: %w", exp.Code, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query character class %s: %w", exp.Code, result.Error)
		} else {
			class.Name = exp.Name
			class.Description = exp.Description
			class.SourceID = sourceID
			class.GameSystem = exp.GameSystem

			if err := database.DB.Save(&class).Error; err != nil {
				return fmt.Errorf("failed to update character class %s: %w", exp.Code, err)
			}
		}
	}

	return nil
}

// ExportSpellSchools exports all spell schools to a JSON file
func ExportSpellSchools(outputDir string) error {
	var schools []models.SpellSchool
	if err := database.DB.Find(&schools).Error; err != nil {
		return fmt.Errorf("failed to fetch spell schools: %w", err)
	}

	// Get source map
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[uint]string)
	for _, s := range sources {
		sourceMap[s.ID] = s.Code
	}

	exportable := make([]ExportableSpellSchool, len(schools))
	for i, school := range schools {
		exportable[i] = ExportableSpellSchool{
			Name:        school.Name,
			Description: school.Description,
			SourceCode:  sourceMap[school.SourceID],
			GameSystem:  school.GameSystem,
		}
	}

	return writeJSON(filepath.Join(outputDir, "spell_schools.json"), exportable)
}

// ImportSpellSchools imports spell schools from a JSON file
func ImportSpellSchools(inputDir string) error {
	var exportable []ExportableSpellSchool
	if err := readJSON(filepath.Join(inputDir, "spell_schools.json"), &exportable); err != nil {
		return err
	}

	// Get source map
	var sources []models.Source
	database.DB.Find(&sources)
	sourceMap := make(map[string]uint)
	for _, s := range sources {
		sourceMap[s.Code] = s.ID
	}

	for _, exp := range exportable {
		var school models.SpellSchool
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&school)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			school = models.SpellSchool{
				Name:        exp.Name,
				Description: exp.Description,
				SourceID:    sourceID,
				GameSystem:  exp.GameSystem,
			}
			if err := database.DB.Create(&school).Error; err != nil {
				return fmt.Errorf("failed to create spell school %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query spell school %s: %w", exp.Name, result.Error)
		} else {
			school.Description = exp.Description
			school.SourceID = sourceID
			school.GameSystem = exp.GameSystem

			if err := database.DB.Save(&school).Error; err != nil {
				return fmt.Errorf("failed to update spell school %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}

// Export/Import functions for ClassCategoryEPCost, ClassSpellSchoolEPCost, SpellLevelLECost
// SkillImprovementCost, WeaponSkill, Equipment, Weapon, Container, Transportation, Believe

// ExportClassCategoryEPCosts exports class-category EP costs
func ExportClassCategoryEPCosts(outputDir string) error {
	var costs []models.ClassCategoryEPCost
	if err := database.DB.Find(&costs).Error; err != nil {
		return fmt.Errorf("failed to fetch class category EP costs: %w", err)
	}

	var classes []models.CharacterClass
	database.DB.Find(&classes)
	classMap := make(map[uint]string)
	for _, c := range classes {
		classMap[c.ID] = c.Code
	}

	var categories []models.SkillCategory
	database.DB.Find(&categories)
	categoryMap := make(map[uint]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	exportable := make([]ExportableClassCategoryEPCost, len(costs))
	for i, cost := range costs {
		exportable[i] = ExportableClassCategoryEPCost{
			CharacterClassCode: classMap[cost.CharacterClassID],
			SkillCategoryName:  categoryMap[cost.SkillCategoryID],
			EPPerTE:            cost.EPPerTE,
		}
	}

	return writeJSON(filepath.Join(outputDir, "class_category_ep_costs.json"), exportable)
}

// ImportClassCategoryEPCosts imports class-category EP costs
func ImportClassCategoryEPCosts(inputDir string) error {
	var exportable []ExportableClassCategoryEPCost
	if err := readJSON(filepath.Join(inputDir, "class_category_ep_costs.json"), &exportable); err != nil {
		return err
	}

	var classes []models.CharacterClass
	database.DB.Find(&classes)
	classMap := make(map[string]uint)
	for _, c := range classes {
		classMap[c.Code] = c.ID
	}

	var categories []models.SkillCategory
	database.DB.Find(&categories)
	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	for _, exp := range exportable {
		classID := classMap[exp.CharacterClassCode]
		categoryID := categoryMap[exp.SkillCategoryName]

		var cost models.ClassCategoryEPCost
		result := database.DB.Where("character_class_id = ? AND skill_category_id = ?", classID, categoryID).First(&cost)

		if result.Error == gorm.ErrRecordNotFound {
			cost = models.ClassCategoryEPCost{
				CharacterClassID: classID,
				SkillCategoryID:  categoryID,
				EPPerTE:          exp.EPPerTE,
				CCLass:           exp.CharacterClassCode,
				SCategory:        exp.SkillCategoryName,
			}
			if err := database.DB.Create(&cost).Error; err != nil {
				return fmt.Errorf("failed to create class category EP cost: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query class category EP cost: %w", result.Error)
		} else {
			cost.EPPerTE = exp.EPPerTE
			cost.CCLass = exp.CharacterClassCode
			cost.SCategory = exp.SkillCategoryName
			if err := database.DB.Save(&cost).Error; err != nil {
				return fmt.Errorf("failed to update class category EP cost: %w", err)
			}
		}
	}

	return nil
}

// ExportClassSpellSchoolEPCosts exports class-spell school EP costs
func ExportClassSpellSchoolEPCosts(outputDir string) error {
	var costs []models.ClassSpellSchoolEPCost
	if err := database.DB.Find(&costs).Error; err != nil {
		return fmt.Errorf("failed to fetch class spell school EP costs: %w", err)
	}

	var classes []models.CharacterClass
	database.DB.Find(&classes)
	classMap := make(map[uint]string)
	for _, c := range classes {
		classMap[c.ID] = c.Code
	}

	var schools []models.SpellSchool
	database.DB.Find(&schools)
	schoolMap := make(map[uint]string)
	for _, s := range schools {
		schoolMap[s.ID] = s.Name
	}

	exportable := make([]ExportableClassSpellSchoolEPCost, len(costs))
	for i, cost := range costs {
		exportable[i] = ExportableClassSpellSchoolEPCost{
			CharacterClassCode: classMap[cost.CharacterClassID],
			SpellSchoolName:    schoolMap[cost.SpellSchoolID],
			EPPerLE:            cost.EPPerLE,
		}
	}

	return writeJSON(filepath.Join(outputDir, "class_spell_school_ep_costs.json"), exportable)
}

// ImportClassSpellSchoolEPCosts imports class-spell school EP costs
func ImportClassSpellSchoolEPCosts(inputDir string) error {
	var exportable []ExportableClassSpellSchoolEPCost
	if err := readJSON(filepath.Join(inputDir, "class_spell_school_ep_costs.json"), &exportable); err != nil {
		return err
	}

	var classes []models.CharacterClass
	database.DB.Find(&classes)
	classMap := make(map[string]uint)
	for _, c := range classes {
		classMap[c.Code] = c.ID
	}

	var schools []models.SpellSchool
	database.DB.Find(&schools)
	schoolMap := make(map[string]uint)
	for _, s := range schools {
		schoolMap[s.Name] = s.ID
	}

	for _, exp := range exportable {
		classID := classMap[exp.CharacterClassCode]
		schoolID := schoolMap[exp.SpellSchoolName]

		var cost models.ClassSpellSchoolEPCost
		result := database.DB.Where("character_class_id = ? AND spell_school_id = ?", classID, schoolID).First(&cost)

		if result.Error == gorm.ErrRecordNotFound {
			cost = models.ClassSpellSchoolEPCost{
				CharacterClassID: classID,
				SpellSchoolID:    schoolID,
				EPPerLE:          exp.EPPerLE,
				CCLass:           exp.CharacterClassCode,
				SCategory:        exp.SpellSchoolName,
			}
			if err := database.DB.Create(&cost).Error; err != nil {
				return fmt.Errorf("failed to create class spell school EP cost: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query class spell school EP cost: %w", result.Error)
		} else {
			cost.EPPerLE = exp.EPPerLE
			cost.CCLass = exp.CharacterClassCode
			cost.SCategory = exp.SpellSchoolName
			if err := database.DB.Save(&cost).Error; err != nil {
				return fmt.Errorf("failed to update class spell school EP cost: %w", err)
			}
		}
	}

	return nil
}

// ExportSpellLevelLECosts exports spell level LE costs
func ExportSpellLevelLECosts(outputDir string) error {
	var costs []models.SpellLevelLECost
	if err := database.DB.Find(&costs).Error; err != nil {
		return fmt.Errorf("failed to fetch spell level LE costs: %w", err)
	}

	exportable := make([]ExportableSpellLevelLECost, len(costs))
	for i, cost := range costs {
		exportable[i] = ExportableSpellLevelLECost{
			Level:      cost.Level,
			LERequired: cost.LERequired,
			GameSystem: cost.GameSystem,
		}
	}

	return writeJSON(filepath.Join(outputDir, "spell_level_le_costs.json"), exportable)
}

// ImportSpellLevelLECosts imports spell level LE costs
func ImportSpellLevelLECosts(inputDir string) error {
	var exportable []ExportableSpellLevelLECost
	if err := readJSON(filepath.Join(inputDir, "spell_level_le_costs.json"), &exportable); err != nil {
		return err
	}

	for _, exp := range exportable {
		var cost models.SpellLevelLECost
		result := database.DB.Where("level = ? AND game_system = ?", exp.Level, exp.GameSystem).First(&cost)

		if result.Error == gorm.ErrRecordNotFound {
			cost = models.SpellLevelLECost{
				Level:      exp.Level,
				LERequired: exp.LERequired,
				GameSystem: exp.GameSystem,
			}
			if err := database.DB.Create(&cost).Error; err != nil {
				return fmt.Errorf("failed to create spell level LE cost: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query spell level LE cost: %w", result.Error)
		} else {
			cost.LERequired = exp.LERequired
			if err := database.DB.Save(&cost).Error; err != nil {
				return fmt.Errorf("failed to update spell level LE cost: %w", err)
			}
		}
	}

	return nil
}

// Note: SkillImprovementCost, WeaponSkill, Equipment, Weapon, Container, Transportation, Believe
// export/import functions follow similar patterns - implement as needed

// ExportAll exports all master data to the specified directory
func ExportAll(outputDir string) error {
	// Export in dependency order
	if err := ExportSources(outputDir); err != nil {
		return err
	}
	if err := ExportCharacterClasses(outputDir); err != nil {
		return err
	}
	if err := ExportSkillCategories(outputDir); err != nil {
		return err
	}
	if err := ExportSkillDifficulties(outputDir); err != nil {
		return err
	}
	if err := ExportSpellSchools(outputDir); err != nil {
		return err
	}
	if err := ExportSkills(outputDir); err != nil {
		return err
	}
	if err := ExportSkillCategoryDifficulties(outputDir); err != nil {
		return err
	}
	if err := ExportSpells(outputDir); err != nil {
		return err
	}
	if err := ExportClassCategoryEPCosts(outputDir); err != nil {
		return err
	}
	if err := ExportClassSpellSchoolEPCosts(outputDir); err != nil {
		return err
	}
	if err := ExportSpellLevelLECosts(outputDir); err != nil {
		return err
	}
	if err := ExportSkillImprovementCosts(outputDir); err != nil {
		return err
	}
	if err := ExportWeaponSkills(outputDir); err != nil {
		return err
	}
	if err := ExportEquipment(outputDir); err != nil {
		return err
	}
	if err := ExportWeapons(outputDir); err != nil {
		return err
	}
	if err := ExportContainers(outputDir); err != nil {
		return err
	}
	if err := ExportTransportation(outputDir); err != nil {
		return err
	}
	if err := ExportBelieves(outputDir); err != nil {
		return err
	}

	return nil
}

// ImportAll imports all master data from the specified directory
func ImportAll(inputDir string) error {
	// Import in dependency order
	if err := ImportSources(inputDir); err != nil {
		return err
	}
	if err := ImportCharacterClasses(inputDir); err != nil {
		return err
	}
	if err := ImportSkillCategories(inputDir); err != nil {
		return err
	}
	if err := ImportSkillDifficulties(inputDir); err != nil {
		return err
	}
	if err := ImportSpellSchools(inputDir); err != nil {
		return err
	}
	if err := ImportSkills(inputDir); err != nil {
		return err
	}
	if err := ImportSkillCategoryDifficulties(inputDir); err != nil {
		return err
	}
	if err := ImportSpells(inputDir); err != nil {
		return err
	}
	if err := ImportClassCategoryEPCosts(inputDir); err != nil {
		return err
	}
	if err := ImportClassSpellSchoolEPCosts(inputDir); err != nil {
		return err
	}
	if err := ImportSpellLevelLECosts(inputDir); err != nil {
		return err
	}
	if err := ImportSkillImprovementCosts(inputDir); err != nil {
		return err
	}
	if err := ImportWeaponSkills(inputDir); err != nil {
		return err
	}
	if err := ImportEquipment(inputDir); err != nil {
		return err
	}
	if err := ImportWeapons(inputDir); err != nil {
		return err
	}
	if err := ImportContainers(inputDir); err != nil {
		return err
	}
	if err := ImportTransportation(inputDir); err != nil {
		return err
	}
	if err := ImportBelieves(inputDir); err != nil {
		return err
	}

	return nil
}
