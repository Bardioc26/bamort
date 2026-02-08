package transfero

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ImportCharacter imports a character from export data
func ImportCharacter(exportData *CharacterExport, userID uint) (uint, error) {
	var importedCharID uint

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Import GSM master data first
		if err := importGSMData(tx, exportData); err != nil {
			return fmt.Errorf("failed to import GSM data: %w", err)
		}

		// Import learning data
		if err := importLearningData(tx, exportData); err != nil {
			return fmt.Errorf("failed to import learning data: %w", err)
		}

		// Import character
		char := exportData.Character
		char.ID = 0 // Reset ID for new character
		char.UserID = userID

		// Reset all related IDs
		if char.Lp.ID != 0 {
			char.Lp.ID = 0
		}
		if char.Ap.ID != 0 {
			char.Ap.ID = 0
		}
		if char.B.ID != 0 {
			char.B.ID = 0
		}
		if char.Merkmale.ID != 0 {
			char.Merkmale.ID = 0
		}
		if char.Bennies.ID != 0 {
			char.Bennies.ID = 0
		}
		if char.Vermoegen.ID != 0 {
			char.Vermoegen.ID = 0
		}
		if char.Erfahrungsschatz.ID != 0 {
			char.Erfahrungsschatz.ID = 0
		}

		// Reset skill IDs
		for i := range char.Eigenschaften {
			char.Eigenschaften[i].ID = 0
			char.Eigenschaften[i].UserID = userID
		}
		for i := range char.Fertigkeiten {
			char.Fertigkeiten[i].ID = 0
			char.Fertigkeiten[i].UserID = userID
		}
		for i := range char.Waffenfertigkeiten {
			char.Waffenfertigkeiten[i].ID = 0
			char.Waffenfertigkeiten[i].UserID = userID
		}
		for i := range char.Zauber {
			char.Zauber[i].ID = 0
			char.Zauber[i].UserID = userID
		}

		// Reset equipment IDs
		for i := range char.Waffen {
			char.Waffen[i].ID = 0
			char.Waffen[i].UserID = userID
		}
		for i := range char.Behaeltnisse {
			char.Behaeltnisse[i].ID = 0
			char.Behaeltnisse[i].UserID = userID
		}
		for i := range char.Transportmittel {
			char.Transportmittel[i].ID = 0
			char.Transportmittel[i].UserID = userID
		}
		for i := range char.Ausruestung {
			char.Ausruestung[i].ID = 0
			char.Ausruestung[i].UserID = userID
		}

		// Create character
		if err := tx.Create(&char).Error; err != nil {
			return fmt.Errorf("failed to create character: %w", err)
		}

		importedCharID = char.ID

		// Import audit log entries
		if len(exportData.AuditLogEntries) > 0 {
			for i := range exportData.AuditLogEntries {
				exportData.AuditLogEntries[i].ID = 0
				exportData.AuditLogEntries[i].CharacterID = importedCharID
			}
			if err := tx.Create(&exportData.AuditLogEntries).Error; err != nil {
				return fmt.Errorf("failed to import audit log: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return importedCharID, nil
}

// importGSMData imports or updates GSM master data
func importGSMData(tx *gorm.DB, exportData *CharacterExport) error {
	// Import skills
	for _, skill := range exportData.GSMSkills {
		if err := importOrUpdateSkill(tx, &skill); err != nil {
			return err
		}
	}

	// Import weapon skills
	for _, weaponSkill := range exportData.GSMWeaponSkills {
		if err := importOrUpdateWeaponSkill(tx, &weaponSkill); err != nil {
			return err
		}
	}

	// Import spells
	for _, spell := range exportData.GSMSpells {
		if err := importOrUpdateSpell(tx, &spell); err != nil {
			return err
		}
	}

	// Import weapons
	for _, weapon := range exportData.GSMWeapons {
		if err := importOrUpdateWeapon(tx, &weapon); err != nil {
			return err
		}
	}

	// Import equipment
	for _, equipment := range exportData.GSMEquipment {
		if err := importOrUpdateEquipment(tx, &equipment); err != nil {
			return err
		}
	}

	// Import containers
	for _, container := range exportData.GSMContainers {
		if err := importOrUpdateContainer(tx, &container); err != nil {
			return err
		}
	}

	return nil
}

// importOrUpdateSkill imports or updates a skill based on name
func importOrUpdateSkill(tx *gorm.DB, skill *models.Skill) error {
	// Set default source_id if 0
	if skill.SourceID == 0 {
		skill.SourceID = 1
	}

	// Ensure game system fields are populated so queries match existing records
	gs := models.GetGameSystem(skill.GameSystemId, skill.GameSystem)
	skill.GameSystem = gs.Name
	skill.GameSystemId = gs.ID

	var existing models.Skill
	err := tx.Where("name = ? AND game_system = ?", skill.Name, skill.GameSystem).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Create new skill
		skill.ID = 0
		return tx.Create(skill).Error
	} else if err != nil {
		return err
	}

	// Update if existing has empty fields
	updates := make(map[string]interface{})
	if existing.Beschreibung == "" && skill.Beschreibung != "" {
		updates["beschreibung"] = skill.Beschreibung
	}
	if existing.Category == "" && skill.Category != "" {
		updates["category"] = skill.Category
	}
	if existing.Difficulty == "" && skill.Difficulty != "" {
		updates["difficulty"] = skill.Difficulty
	}
	if existing.Bonuseigenschaft == "" && skill.Bonuseigenschaft != "" {
		updates["bonuseigenschaft"] = skill.Bonuseigenschaft
	}
	if existing.SourceID == 0 && skill.SourceID != 0 {
		updates["source_id"] = skill.SourceID
	}
	if skill.PageNumber != 0 && existing.PageNumber == 0 {
		updates["page_number"] = skill.PageNumber
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}

// importOrUpdateWeaponSkill imports or updates a weapon skill
func importOrUpdateWeaponSkill(tx *gorm.DB, weaponSkill *models.WeaponSkill) error {
	if weaponSkill.SourceID == 0 {
		weaponSkill.SourceID = 1
	}

	var existing models.WeaponSkill
	err := tx.Where("name = ? AND game_system = ?", weaponSkill.Name, weaponSkill.GameSystem).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		weaponSkill.ID = 0
		return tx.Create(weaponSkill).Error
	} else if err != nil {
		return err
	}

	// Update empty fields
	updates := make(map[string]interface{})
	if existing.Beschreibung == "" && weaponSkill.Beschreibung != "" {
		updates["beschreibung"] = weaponSkill.Beschreibung
	}
	if existing.Category == "" && weaponSkill.Category != "" {
		updates["category"] = weaponSkill.Category
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}

// importOrUpdateSpell imports or updates a spell
func importOrUpdateSpell(tx *gorm.DB, spell *models.Spell) error {
	// Set default source_id if 0 (spells get source_id 2)
	if spell.SourceID == 0 {
		spell.SourceID = 2
	}

	var existing models.Spell
	err := tx.Where("name = ? AND game_system = ?", spell.Name, spell.GameSystem).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		spell.ID = 0
		return tx.Create(spell).Error
	} else if err != nil {
		return err
	}

	// Update empty fields
	updates := make(map[string]interface{})
	if existing.Beschreibung == "" && spell.Beschreibung != "" {
		updates["beschreibung"] = spell.Beschreibung
	}
	if existing.Category == "" && spell.Category != "" {
		updates["category"] = spell.Category
	}
	if existing.LearningCategory == "" && spell.LearningCategory != "" {
		updates["learning_category"] = spell.LearningCategory
	}
	if existing.Ursprung == "" && spell.Ursprung != "" {
		updates["ursprung"] = spell.Ursprung
	}
	if existing.SourceID == 0 && spell.SourceID != 0 {
		updates["source_id"] = spell.SourceID
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}

// importOrUpdateWeapon imports or updates a weapon
func importOrUpdateWeapon(tx *gorm.DB, weapon *models.Weapon) error {
	if weapon.SourceID == 0 {
		weapon.SourceID = 1
	}

	var existing models.Weapon
	err := tx.Where("name = ? AND game_system = ?", weapon.Name, weapon.GameSystem).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		weapon.ID = 0
		return tx.Create(weapon).Error
	} else if err != nil {
		return err
	}

	// Update empty fields
	updates := make(map[string]interface{})
	if existing.Beschreibung == "" && weapon.Beschreibung != "" {
		updates["beschreibung"] = weapon.Beschreibung
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}

// importOrUpdateEquipment imports or updates equipment
func importOrUpdateEquipment(tx *gorm.DB, equipment *models.Equipment) error {
	if equipment.SourceID == 0 {
		equipment.SourceID = 1
	}

	var existing models.Equipment
	err := tx.Where("name = ? AND game_system = ?", equipment.Name, equipment.GameSystem).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		equipment.ID = 0
		return tx.Create(equipment).Error
	} else if err != nil {
		return err
	}

	// Update empty fields
	updates := make(map[string]interface{})
	if existing.Beschreibung == "" && equipment.Beschreibung != "" {
		updates["beschreibung"] = equipment.Beschreibung
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}

// importOrUpdateContainer imports or updates a container
func importOrUpdateContainer(tx *gorm.DB, container *models.Container) error {
	if container.SourceID == 0 {
		container.SourceID = 1
	}

	var existing models.Container
	err := tx.Where("name = ? AND game_system = ?", container.Name, container.GameSystem).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		container.ID = 0
		return tx.Create(container).Error
	} else if err != nil {
		return err
	}

	// Update empty fields
	updates := make(map[string]interface{})
	if existing.Beschreibung == "" && container.Beschreibung != "" {
		updates["beschreibung"] = container.Beschreibung
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}

// importLearningData imports learning-related master data
func importLearningData(tx *gorm.DB, exportData *CharacterExport) error {
	// Import sources
	for _, source := range exportData.LearningData.Sources {
		if err := importOrUpdateSource(tx, &source); err != nil {
			return err
		}
	}

	// Import character classes
	for _, cc := range exportData.LearningData.CharacterClasses {
		if err := importOrUpdateCharacterClass(tx, &cc); err != nil {
			return err
		}
	}

	// Import skill categories
	for _, sc := range exportData.LearningData.SkillCategories {
		if err := importOrUpdateSkillCategory(tx, &sc); err != nil {
			return err
		}
	}

	// Import skill difficulties
	for _, sd := range exportData.LearningData.SkillDifficulties {
		if err := importOrUpdateSkillDifficulty(tx, &sd); err != nil {
			return err
		}
	}

	// Import spell schools
	for _, ss := range exportData.LearningData.SpellSchools {
		if err := importOrUpdateSpellSchool(tx, &ss); err != nil {
			return err
		}
	}

	// More complex tables - skip if already exist (identified by combination of fields)
	// These don't need updates as they're typically static cost tables

	return nil
}

func importOrUpdateSource(tx *gorm.DB, source *models.Source) error {
	var existing models.Source
	err := tx.Where("code = ?", source.Code).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		source.ID = 0
		return tx.Create(source).Error
	} else if err != nil {
		return err
	}

	// Update empty fields
	updates := make(map[string]interface{})
	if existing.FullName == "" && source.FullName != "" {
		updates["full_name"] = source.FullName
	}
	if existing.Description == "" && source.Description != "" {
		updates["description"] = source.Description
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}

func importOrUpdateCharacterClass(tx *gorm.DB, cc *models.CharacterClass) error {
	var existing models.CharacterClass
	err := tx.Where("code = ?", strings.TrimSpace(cc.Code)).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		cc.ID = 0
		return tx.Create(cc).Error
	} else if err != nil {
		return err
	}

	// Update empty fields
	updates := make(map[string]interface{})
	if existing.Description == "" && cc.Description != "" {
		updates["description"] = cc.Description
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}

func importOrUpdateSkillCategory(tx *gorm.DB, sc *models.SkillCategory) error {
	var existing models.SkillCategory
	err := tx.Where("name = ?", sc.Name).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		sc.ID = 0
		return tx.Create(sc).Error
	}

	return err
}

func importOrUpdateSkillDifficulty(tx *gorm.DB, sd *models.SkillDifficulty) error {
	var existing models.SkillDifficulty
	err := tx.Where("name = ?", sd.Name).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		sd.ID = 0
		return tx.Create(sd).Error
	}

	return err
}

func importOrUpdateSpellSchool(tx *gorm.DB, ss *models.SpellSchool) error {
	var existing models.SpellSchool
	err := tx.Where("name = ?", ss.Name).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		ss.ID = 0
		return tx.Create(ss).Error
	} else if err != nil {
		return err
	}

	// Update empty fields
	updates := make(map[string]interface{})
	if existing.Description == "" && ss.Description != "" {
		updates["description"] = ss.Description
	}

	if len(updates) > 0 {
		return tx.Model(&existing).Updates(updates).Error
	}

	return nil
}
