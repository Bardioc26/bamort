package importer

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ReconcileSkill reconciles an imported skill with master data.
// Returns the master data skill, match type ("exact" or "created_personal"), and error.
func ReconcileSkill(skill Fertigkeit, userID uint, gameSystem string) (*models.Skill, string, error) {
	return ReconcileSkillWithHistory(skill, 0, gameSystem)
}

// ReconcileSkillWithHistory reconciles a skill and logs to ImportHistory
func ReconcileSkillWithHistory(skill Fertigkeit, importHistoryID uint, gameSystem string) (*models.Skill, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Skill
	err := database.DB.Where("name = ? AND game_system = ?", skill.Name, gs.Name).First(&existing).Error

	if err == nil {
		// Exact match found
		if importHistoryID > 0 {
			logMasterDataImport(importHistoryID, "skill", existing.ID, skill.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query skill: %w", err)
	}

	// Create new personal item
	newSkill := &models.Skill{
		Name:             skill.Name,
		GameSystem:       gs.Name,
		GameSystemId:     gs.ID,
		Beschreibung:     skill.Beschreibung,
		Initialwert:      skill.Fertigkeitswert,
		Quelle:           skill.Quelle,
		Bonuseigenschaft: "check",
		Improvable:       true,
		PersonalItem:     true,
		SourceID:         1, // Default source
	}

	if err := database.DB.Create(newSkill).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create skill: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImport(importHistoryID, "skill", newSkill.ID, skill.Name, "created_personal")
	}

	return newSkill, "created_personal", nil
}

// ReconcileWeaponSkill reconciles an imported weapon skill with master data
func ReconcileWeaponSkill(ws Waffenfertigkeit, userID uint, gameSystem string) (*models.WeaponSkill, string, error) {
	return ReconcileWeaponSkillWithHistory(ws, 0, gameSystem)
}

// ReconcileWeaponSkillWithHistory reconciles a weapon skill and logs to ImportHistory
func ReconcileWeaponSkillWithHistory(ws Waffenfertigkeit, importHistoryID uint, gameSystem string) (*models.WeaponSkill, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.WeaponSkill
	err := database.DB.Where("name = ? AND game_system = ?", ws.Name, gs.Name).First(&existing).Error

	if err == nil {
		// Exact match found
		if importHistoryID > 0 {
			logMasterDataImport(importHistoryID, "weaponskill", existing.ID, ws.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query weapon skill: %w", err)
	}

	// Create new personal item
	newWS := &models.WeaponSkill{
		Skill: models.Skill{
			Name:         ws.Name,
			GameSystem:   gs.Name,
			GameSystemId: gs.ID,
			Beschreibung: ws.Beschreibung,
			Quelle:       ws.Quelle,
			PersonalItem: true,
			SourceID:     1,
		},
	}

	if err := database.DB.Create(newWS).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create weapon skill: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImport(importHistoryID, "weaponskill", newWS.ID, ws.Name, "created_personal")
	}

	return newWS, "created_personal", nil
}

// ReconcileSpell reconciles an imported spell with master data
func ReconcileSpell(spell Zauber, userID uint, gameSystem string) (*models.Spell, string, error) {
	return ReconcileSpellWithHistory(spell, 0, gameSystem)
}

// ReconcileSpellWithHistory reconciles a spell and logs to ImportHistory
func ReconcileSpellWithHistory(spell Zauber, importHistoryID uint, gameSystem string) (*models.Spell, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Spell
	err := database.DB.Where("name = ? AND game_system = ?", spell.Name, gs.Name).First(&existing).Error

	if err == nil {
		// Exact match found
		if importHistoryID > 0 {
			logMasterDataImport(importHistoryID, "spell", existing.ID, spell.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query spell: %w", err)
	}

	// Create new personal item
	newSpell := &models.Spell{
		Name:         spell.Name,
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
		Beschreibung: spell.Beschreibung,
		Quelle:       spell.Quelle,
		PersonalItem: true,
		SourceID:     2, // Default source for spells
	}

	if err := database.DB.Create(newSpell).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create spell: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImport(importHistoryID, "spell", newSpell.ID, spell.Name, "created_personal")
	}

	return newSpell, "created_personal", nil
}

// ReconcileWeapon reconciles an imported weapon with master data
func ReconcileWeapon(weapon Waffe, userID uint, gameSystem string) (*models.Weapon, string, error) {
	return ReconcileWeaponWithHistory(weapon, 0, gameSystem)
}

// ReconcileWeaponWithHistory reconciles a weapon and logs to ImportHistory
func ReconcileWeaponWithHistory(weapon Waffe, importHistoryID uint, gameSystem string) (*models.Weapon, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Weapon
	err := database.DB.Where("name = ? AND game_system = ?", weapon.Name, gs.Name).First(&existing).Error

	if err == nil {
		// Exact match found
		if importHistoryID > 0 {
			logMasterDataImport(importHistoryID, "weapon", existing.ID, weapon.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query weapon: %w", err)
	}

	// Create new personal item
	newWeapon := &models.Weapon{
		Equipment: models.Equipment{
			Name:         weapon.Name,
			GameSystem:   gs.Name,
			GameSystemId: gs.ID,
			Beschreibung: weapon.Beschreibung,
			Gewicht:      weapon.Gewicht,
			Wert:         weapon.Wert,
			PersonalItem: true,
			SourceID:     1,
		},
	}

	if err := database.DB.Create(newWeapon).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create weapon: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImport(importHistoryID, "weapon", newWeapon.ID, weapon.Name, "created_personal")
	}

	return newWeapon, "created_personal", nil
}

// ReconcileEquipment reconciles imported equipment with master data
func ReconcileEquipment(equip Ausruestung, userID uint, gameSystem string) (*models.Equipment, string, error) {
	return ReconcileEquipmentWithHistory(equip, 0, gameSystem)
}

// ReconcileEquipmentWithHistory reconciles equipment and logs to ImportHistory
func ReconcileEquipmentWithHistory(equip Ausruestung, importHistoryID uint, gameSystem string) (*models.Equipment, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Equipment
	err := database.DB.Where("name = ? AND game_system = ?", equip.Name, gs.Name).First(&existing).Error

	if err == nil {
		// Exact match found
		if importHistoryID > 0 {
			logMasterDataImport(importHistoryID, "equipment", existing.ID, equip.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query equipment: %w", err)
	}

	// Create new personal item
	newEquip := &models.Equipment{
		Name:         equip.Name,
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
		Beschreibung: equip.Beschreibung,
		Gewicht:      equip.Gewicht,
		Wert:         equip.Wert,
		PersonalItem: true,
		SourceID:     1,
	}

	if err := database.DB.Create(newEquip).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create equipment: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImport(importHistoryID, "equipment", newEquip.ID, equip.Name, "created_personal")
	}

	return newEquip, "created_personal", nil
}

// ReconcileContainer reconciles an imported container with master data
func ReconcileContainer(container Behaeltniss, userID uint, gameSystem string) (*models.Container, string, error) {
	return ReconcileContainerWithHistory(container, 0, gameSystem)
}

// ReconcileContainerWithHistory reconciles a container and logs to ImportHistory
func ReconcileContainerWithHistory(container Behaeltniss, importHistoryID uint, gameSystem string) (*models.Container, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Container
	err := database.DB.Where("name = ? AND game_system = ?", container.Name, gs.Name).First(&existing).Error

	if err == nil {
		// Exact match found
		if importHistoryID > 0 {
			logMasterDataImport(importHistoryID, "container", existing.ID, container.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query container: %w", err)
	}

	// Create new personal item
	newContainer := &models.Container{
		Equipment: models.Equipment{
			Name:         container.Name,
			GameSystem:   gs.Name,
			GameSystemId: gs.ID,
			Beschreibung: container.Beschreibung,
			Gewicht:      container.Gewicht,
			Wert:         container.Wert,
			PersonalItem: true,
			SourceID:     1,
		},
		Tragkraft: container.Tragkraft,
		Volumen:   container.Volumen,
	}

	if err := database.DB.Create(newContainer).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create container: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImport(importHistoryID, "container", newContainer.ID, container.Name, "created_personal")
	}

	return newContainer, "created_personal", nil
}

// logMasterDataImport creates a log entry in the MasterDataImport table
func logMasterDataImport(importHistoryID uint, itemType string, itemID uint, externalName, matchType string) {
	log := MasterDataImport{
		ImportHistoryID: importHistoryID,
		ItemType:        itemType,
		ItemID:          itemID,
		ExternalName:    externalName,
		MatchType:       matchType,
		CreatedAt:       time.Now(),
	}

	// Best effort logging - don't fail import if logging fails
	if err := database.DB.Create(&log).Error; err != nil {
		// Log error but don't return it
		fmt.Printf("Warning: Failed to log master data import: %v\n", err)
	}
}
