package gsmaster

/*

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"path/filepath"

	"gorm.io/gorm"
)

// ExportSkillImprovementCosts exports all skill improvement costs to a JSON file
func ExportSkillImprovementCosts(outputDir string) error {
	var costs []models.SkillImprovementCost2
	if err := database.DB.Find(&costs).Error; err != nil {
		return fmt.Errorf("failed to fetch skill improvement costs: %w", err)
	}

	exportable := make([]ExportableSkillImprovementCost, 0, len(costs))
	for _, cost := range costs {
		// Resolve category and difficulty names
		var lc models.SkillCategory
		if err := database.DB.First(&lc, cost.CategoryID).Error; err != nil {
			continue
		}
		var ld models.SkillDifficulty
		if err := database.DB.First(&ld, cost.DifficultyID).Error; err != nil {
			continue
		}
		// Find skill via learning_skill_category_difficulties
		var scd models.SkillCategoryDifficulty
		if err := database.DB.Where("skill_category = ? AND skill_difficulty = ?", lc.Name, ld.Name).First(&scd).Error; err != nil {
			continue
		}
		var skill models.Skill
		if err := database.DB.First(&skill, scd.SkillID).Error; err != nil {
			continue
		}

		exportable = append(exportable, ExportableSkillImprovementCost{
			SkillName:        skill.Name,
			SkillSystem:      skill.GameSystem,
			CategoryName:     lc.Name,
			CategorySystem:   lc.GameSystem,
			DifficultyName:   ld.Name,
			DifficultySystem: ld.GameSystem,
			CurrentLevel:     cost.CurrentLevel,
			TERequired:       cost.TERequired,
		})
	}

	return writeJSON(filepath.Join(outputDir, "skill_improvement_costs.json"), exportable)
}

// ImportSkillImprovementCosts imports skill improvement costs from a JSON file
func ImportSkillImprovementCosts(inputDir string) error {
	var exportable []ExportableSkillImprovementCost
	if err := readJSON(filepath.Join(inputDir, "skill_improvement_costs.json"), &exportable); err != nil {
		return err
	}

	// Build lookup maps using helpers
	skillMap := buildSkillMap()
	categoryMap := buildCategoryMap()
	difficultyMap := buildDifficultyMap()

	for _, exp := range exportable {
		// Find skill ID
		skillID, ok := skillMap[exp.SkillSystem][exp.SkillName]
		if !ok {
			return fmt.Errorf("skill not found: %s (%s)", exp.SkillName, exp.SkillSystem)
		}

		// Find category ID
		categoryID, ok := categoryMap[exp.CategorySystem][exp.CategoryName]
		if !ok {
			return fmt.Errorf("category not found: %s (%s)", exp.CategoryName, exp.CategorySystem)
		}

		// Find difficulty ID
		difficultyID, ok := difficultyMap[exp.DifficultySystem][exp.DifficultyName]
		if !ok {
			return fmt.Errorf("difficulty not found: %s (%s)", exp.DifficultyName, exp.DifficultySystem)
		}

		// Find SkillCategoryDifficulty
		var scd models.SkillCategoryDifficulty
		if err := database.DB.Where("skill_id = ? AND skill_category_id = ? AND skill_difficulty_id = ?",
			skillID, categoryID, difficultyID).First(&scd).Error; err != nil {
			return fmt.Errorf("skill category difficulty not found for %s/%s/%s: %w",
				exp.SkillName, exp.CategoryName, exp.DifficultyName, err)
		}

		// Find or create SkillImprovementCost2 using category/difficulty IDs
		var cost models.SkillImprovementCost2
		result := database.DB.Where("skill_category_id = ? AND skill_difficulty_id = ? AND current_level = ?",
			categoryID, difficultyID, exp.CurrentLevel).First(&cost)

		if result.Error == gorm.ErrRecordNotFound {
			cost = models.SkillImprovementCost2{
				CategoryID:   categoryID,
				DifficultyID: difficultyID,
				CurrentLevel: exp.CurrentLevel,
				TERequired:   exp.TERequired,
			}
			if err := database.DB.Create(&cost).Error; err != nil {
				return fmt.Errorf("failed to create skill improvement cost: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query skill improvement cost: %w", result.Error)
		} else {
			cost.TERequired = exp.TERequired
			if err := database.DB.Save(&cost).Error; err != nil {
				return fmt.Errorf("failed to update skill improvement cost: %w", err)
			}
		}
	}

	return nil
}

// ExportWeaponSkills exports all weapon skills to a JSON file
func ExportWeaponSkills(outputDir string) error {
	var skills []models.WeaponSkill
	if err := database.DB.Find(&skills).Error; err != nil {
		return fmt.Errorf("failed to fetch weapon skills: %w", err)
	}

	sourceMap := buildSourceMap()

	exportable := make([]ExportableWeaponSkill, len(skills))
	for i, skill := range skills {
		exportable[i] = ExportableWeaponSkill{
			Name:             skill.Name,
			GameSystem:       skill.GameSystem,
			Beschreibung:     skill.Beschreibung,
			SourceCode:       sourceMap[skill.SourceID],
			PageNumber:       skill.PageNumber,
			Initialwert:      skill.Initialwert,
			BasisWert:        skill.BasisWert,
			Bonuseigenschaft: skill.Bonuseigenschaft,
			Improvable:       skill.Improvable,
			InnateSkill:      skill.InnateSkill,
			Category:         skill.Category,
			Difficulty:       skill.Difficulty,
		}
	}

	return writeJSON(filepath.Join(outputDir, "weapon_skills.json"), exportable)
}

// ImportWeaponSkills imports weapon skills from a JSON file
func ImportWeaponSkills(inputDir string) error {
	var exportable []ExportableWeaponSkill
	if err := readJSON(filepath.Join(inputDir, "weapon_skills.json"), &exportable); err != nil {
		return err
	}

	sourceMap := buildSourceMapReverse()

	for _, exp := range exportable {
		var skill models.WeaponSkill
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&skill)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			skill = models.WeaponSkill{
				Skill: models.Skill{
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
				},
			}
			if err := database.DB.Create(&skill).Error; err != nil {
				return fmt.Errorf("failed to create weapon skill %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query weapon skill %s: %w", exp.Name, result.Error)
		} else {
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
				return fmt.Errorf("failed to update weapon skill %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}

// ExportEquipment exports all equipment to a JSON file
func ExportEquipment(outputDir string) error {
	var equipment []models.Equipment
	if err := database.DB.Find(&equipment).Error; err != nil {
		return fmt.Errorf("failed to fetch equipment: %w", err)
	}

	sourceMap := buildSourceMap()

	exportable := make([]ExportableEquipment, len(equipment))
	for i, eq := range equipment {
		exportable[i] = ExportableEquipment{
			Name:         eq.Name,
			GameSystem:   eq.GameSystem,
			Beschreibung: eq.Beschreibung,
			SourceCode:   sourceMap[eq.SourceID],
			PageNumber:   eq.PageNumber,
			Gewicht:      eq.Gewicht,
			Wert:         eq.Wert,
			PersonalItem: eq.PersonalItem,
		}
	}

	return writeJSON(filepath.Join(outputDir, "equipment.json"), exportable)
}

// ImportEquipment imports equipment from a JSON file
func ImportEquipment(inputDir string) error {
	var exportable []ExportableEquipment
	if err := readJSON(filepath.Join(inputDir, "equipment.json"), &exportable); err != nil {
		return err
	}

	sourceMap := buildSourceMapReverse()

	for _, exp := range exportable {
		var eq models.Equipment
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&eq)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			eq = models.Equipment{
				Name:         exp.Name,
				GameSystem:   exp.GameSystem,
				Beschreibung: exp.Beschreibung,
				SourceID:     sourceID,
				PageNumber:   exp.PageNumber,
				Gewicht:      exp.Gewicht,
				Wert:         exp.Wert,
				PersonalItem: exp.PersonalItem,
			}
			if err := database.DB.Create(&eq).Error; err != nil {
				return fmt.Errorf("failed to create equipment %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query equipment %s: %w", exp.Name, result.Error)
		} else {
			eq.Beschreibung = exp.Beschreibung
			eq.SourceID = sourceID
			eq.PageNumber = exp.PageNumber
			eq.Gewicht = exp.Gewicht
			eq.Wert = exp.Wert
			eq.PersonalItem = exp.PersonalItem

			if err := database.DB.Save(&eq).Error; err != nil {
				return fmt.Errorf("failed to update equipment %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}

// ExportWeapons exports all weapons to a JSON file
func ExportWeapons(outputDir string) error {
	var weapons []models.Weapon
	if err := database.DB.Find(&weapons).Error; err != nil {
		return fmt.Errorf("failed to fetch weapons: %w", err)
	}

	sourceMap := buildSourceMap()

	exportable := make([]ExportableWeapon, len(weapons))
	for i, weapon := range weapons {
		exportable[i] = ExportableWeapon{
			Name:          weapon.Name,
			GameSystem:    weapon.GameSystem,
			Beschreibung:  weapon.Beschreibung,
			SourceCode:    sourceMap[weapon.SourceID],
			PageNumber:    weapon.PageNumber,
			Gewicht:       weapon.Gewicht,
			Wert:          weapon.Wert,
			PersonalItem:  weapon.PersonalItem,
			SkillRequired: weapon.SkillRequired,
			Damage:        weapon.Damage,
			RangeNear:     weapon.RangeNear,
			RangeMiddle:   weapon.RangeMiddle,
			RangeFar:      weapon.RangeFar,
		}
	}

	return writeJSON(filepath.Join(outputDir, "weapons.json"), exportable)
}

// ImportWeapons imports weapons from a JSON file
func ImportWeapons(inputDir string) error {
	var exportable []ExportableWeapon
	if err := readJSON(filepath.Join(inputDir, "weapons.json"), &exportable); err != nil {
		return err
	}

	sourceMap := buildSourceMapReverse()

	for _, exp := range exportable {
		var weapon models.Weapon
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&weapon)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			weapon = models.Weapon{
				Equipment: models.Equipment{
					Name:         exp.Name,
					GameSystem:   exp.GameSystem,
					Beschreibung: exp.Beschreibung,
					SourceID:     sourceID,
					PageNumber:   exp.PageNumber,
					Gewicht:      exp.Gewicht,
					Wert:         exp.Wert,
					PersonalItem: exp.PersonalItem,
				},
				SkillRequired: exp.SkillRequired,
				Damage:        exp.Damage,
				RangeNear:     exp.RangeNear,
				RangeMiddle:   exp.RangeMiddle,
				RangeFar:      exp.RangeFar,
			}
			if err := database.DB.Create(&weapon).Error; err != nil {
				return fmt.Errorf("failed to create weapon %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query weapon %s: %w", exp.Name, result.Error)
		} else {
			weapon.Beschreibung = exp.Beschreibung
			weapon.SourceID = sourceID
			weapon.PageNumber = exp.PageNumber
			weapon.Gewicht = exp.Gewicht
			weapon.Wert = exp.Wert
			weapon.PersonalItem = exp.PersonalItem
			weapon.SkillRequired = exp.SkillRequired
			weapon.Damage = exp.Damage
			weapon.RangeNear = exp.RangeNear
			weapon.RangeMiddle = exp.RangeMiddle
			weapon.RangeFar = exp.RangeFar

			if err := database.DB.Save(&weapon).Error; err != nil {
				return fmt.Errorf("failed to update weapon %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}

// ExportContainers exports all containers to a JSON file
func ExportContainers(outputDir string) error {
	var containers []models.Container
	if err := database.DB.Find(&containers).Error; err != nil {
		return fmt.Errorf("failed to fetch containers: %w", err)
	}

	sourceMap := buildSourceMap()

	exportable := make([]ExportableContainer, len(containers))
	for i, container := range containers {
		exportable[i] = ExportableContainer{
			Name:         container.Name,
			GameSystem:   container.GameSystem,
			Beschreibung: container.Beschreibung,
			SourceCode:   sourceMap[container.SourceID],
			PageNumber:   container.PageNumber,
			Gewicht:      container.Gewicht,
			Wert:         container.Wert,
			PersonalItem: container.PersonalItem,
			Tragkraft:    container.Tragkraft,
			Volumen:      container.Volumen,
		}
	}

	return writeJSON(filepath.Join(outputDir, "containers.json"), exportable)
}

// ImportContainers imports containers from a JSON file
func ImportContainers(inputDir string) error {
	var exportable []ExportableContainer
	if err := readJSON(filepath.Join(inputDir, "containers.json"), &exportable); err != nil {
		return err
	}

	sourceMap := buildSourceMapReverse()

	for _, exp := range exportable {
		var container models.Container
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&container)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			container = models.Container{
				Equipment: models.Equipment{
					Name:         exp.Name,
					GameSystem:   exp.GameSystem,
					Beschreibung: exp.Beschreibung,
					SourceID:     sourceID,
					PageNumber:   exp.PageNumber,
					Gewicht:      exp.Gewicht,
					Wert:         exp.Wert,
					PersonalItem: exp.PersonalItem,
				},
				Tragkraft: exp.Tragkraft,
				Volumen:   exp.Volumen,
			}
			if err := database.DB.Create(&container).Error; err != nil {
				return fmt.Errorf("failed to create container %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query container %s: %w", exp.Name, result.Error)
		} else {
			container.Beschreibung = exp.Beschreibung
			container.SourceID = sourceID
			container.PageNumber = exp.PageNumber
			container.Gewicht = exp.Gewicht
			container.Wert = exp.Wert
			container.PersonalItem = exp.PersonalItem
			container.Tragkraft = exp.Tragkraft
			container.Volumen = exp.Volumen

			if err := database.DB.Save(&container).Error; err != nil {
				return fmt.Errorf("failed to update container %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}

// ExportTransportation exports all transportation to a JSON file
func ExportTransportation(outputDir string) error {
	var transportation []models.Transportation
	if err := database.DB.Find(&transportation).Error; err != nil {
		return fmt.Errorf("failed to fetch transportation: %w", err)
	}

	sourceMap := buildSourceMap()

	exportable := make([]ExportableTransportation, len(transportation))
	for i, trans := range transportation {
		exportable[i] = ExportableTransportation{
			Name:         trans.Name,
			GameSystem:   trans.GameSystem,
			Beschreibung: trans.Beschreibung,
			SourceCode:   sourceMap[trans.SourceID],
			PageNumber:   trans.PageNumber,
			Gewicht:      trans.Gewicht,
			Wert:         trans.Wert,
			PersonalItem: trans.PersonalItem,
			Tragkraft:    trans.Tragkraft,
			Volumen:      trans.Volumen,
		}
	}

	return writeJSON(filepath.Join(outputDir, "transportation.json"), exportable)
}

// ImportTransportation imports transportation from a JSON file
func ImportTransportation(inputDir string) error {
	var exportable []ExportableTransportation
	if err := readJSON(filepath.Join(inputDir, "transportation.json"), &exportable); err != nil {
		return err
	}

	sourceMap := buildSourceMapReverse()

	for _, exp := range exportable {
		var trans models.Transportation
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&trans)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			trans = models.Transportation{
				Container: models.Container{
					Equipment: models.Equipment{
						Name:         exp.Name,
						GameSystem:   exp.GameSystem,
						Beschreibung: exp.Beschreibung,
						SourceID:     sourceID,
						PageNumber:   exp.PageNumber,
						Gewicht:      exp.Gewicht,
						Wert:         exp.Wert,
						PersonalItem: exp.PersonalItem,
					},
					Tragkraft: exp.Tragkraft,
					Volumen:   exp.Volumen,
				},
			}
			if err := database.DB.Create(&trans).Error; err != nil {
				return fmt.Errorf("failed to create transportation %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query transportation %s: %w", exp.Name, result.Error)
		} else {
			trans.Beschreibung = exp.Beschreibung
			trans.SourceID = sourceID
			trans.PageNumber = exp.PageNumber
			trans.Gewicht = exp.Gewicht
			trans.Wert = exp.Wert
			trans.PersonalItem = exp.PersonalItem
			trans.Tragkraft = exp.Tragkraft
			trans.Volumen = exp.Volumen

			if err := database.DB.Save(&trans).Error; err != nil {
				return fmt.Errorf("failed to update transportation %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}

// ExportBelieves exports all beliefs to a JSON file
func ExportBelieves(outputDir string) error {
	var believes []models.Believe
	if err := database.DB.Find(&believes).Error; err != nil {
		return fmt.Errorf("failed to fetch believes: %w", err)
	}

	sourceMap := buildSourceMap()

	exportable := make([]ExportableBelieve, len(believes))
	for i, believe := range believes {
		exportable[i] = ExportableBelieve{
			Name:         believe.Name,
			GameSystem:   believe.GameSystem,
			Beschreibung: believe.Beschreibung,
			SourceCode:   sourceMap[believe.SourceID],
			PageNumber:   believe.PageNumber,
		}
	}

	return writeJSON(filepath.Join(outputDir, "believes.json"), exportable)
}

// ImportBelieves imports believes from a JSON file
func ImportBelieves(inputDir string) error {
	var exportable []ExportableBelieve
	if err := readJSON(filepath.Join(inputDir, "believes.json"), &exportable); err != nil {
		return err
	}

	sourceMap := buildSourceMapReverse()

	for _, exp := range exportable {
		var believe models.Believe
		result := database.DB.Where("name = ? AND game_system = ?", exp.Name, exp.GameSystem).First(&believe)

		sourceID := sourceMap[exp.SourceCode]

		if result.Error == gorm.ErrRecordNotFound {
			believe = models.Believe{
				Name:         exp.Name,
				GameSystem:   exp.GameSystem,
				Beschreibung: exp.Beschreibung,
				SourceID:     sourceID,
				PageNumber:   exp.PageNumber,
			}
			if err := database.DB.Create(&believe).Error; err != nil {
				return fmt.Errorf("failed to create believe %s: %w", exp.Name, err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to query believe %s: %w", exp.Name, result.Error)
		} else {
			believe.Beschreibung = exp.Beschreibung
			believe.SourceID = sourceID
			believe.PageNumber = exp.PageNumber

			if err := database.DB.Save(&believe).Error; err != nil {
				return fmt.Errorf("failed to update believe %s: %w", exp.Name, err)
			}
		}
	}

	return nil
}
*/
