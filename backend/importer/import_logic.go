package importer

import (
	"bamort/database"
	"bamort/models"
	"bytes"
	"compress/gzip"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ImportResult represents the result of a character import operation
type ImportResult struct {
	CharacterID  uint                `json:"character_id"`
	ImportID     uint                `json:"import_id"`
	AdapterID    string              `json:"adapter_id"`
	Warnings     []ValidationWarning `json:"warnings"`
	CreatedItems map[string]int      `json:"created_items"` // {"skills": 3, "spells": 1}
	Status       string              `json:"status"`
}

// ImportCharacter imports a character with full transaction safety
// This implements the transaction-wrapped import logic from Phase 1
func ImportCharacter(char *CharacterImport, userID uint, adapterID string, originalData []byte) (*ImportResult, error) {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	result := &ImportResult{
		AdapterID:    adapterID,
		CreatedItems: make(map[string]int),
		Status:       "in_progress",
	}

	// 1. Create ImportHistory record (failed status initially)
	history := &ImportHistory{
		UserID:         userID,
		AdapterID:      adapterID,
		SourceFormat:   adapterID, // Can be refined based on adapter metadata
		SourceFilename: fmt.Sprintf("%s_import_%d.json", char.Name, time.Now().Unix()),
		BmrtVersion:    "1.0",
		ImportedAt:     time.Now(),
		Status:         "in_progress",
	}

	// Compress original data
	if originalData != nil {
		compressed, err := compressData(originalData)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to compress source data: %w", err)
		}
		history.SourceSnapshot = compressed
	}

	if err := tx.Create(history).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create import history: %w", err)
	}

	result.ImportID = history.ID

	// 2. Reconcile master data and track created items
	createdCounts := make(map[string]int)

	// Reconcile skills
	for _, skill := range char.Fertigkeiten {
		_, matchType, err := reconcileSkillWithTx(tx, skill, history.ID, char.Typ)
		if err != nil {
			history.Status = "failed"
			history.ErrorLog = fmt.Sprintf("Failed to reconcile skill %s: %v", skill.Name, err)
			tx.Save(history)
			tx.Rollback()
			return nil, fmt.Errorf("failed to reconcile skill: %w", err)
		}
		if matchType == "created_personal" {
			createdCounts["skills"]++
		}
	}

	// Reconcile spells
	for _, spell := range char.Zauber {
		_, matchType, err := reconcileSpellWithTx(tx, spell, history.ID, char.Typ)
		if err != nil {
			history.Status = "failed"
			history.ErrorLog = fmt.Sprintf("Failed to reconcile spell %s: %v", spell.Name, err)
			tx.Save(history)
			tx.Rollback()
			return nil, fmt.Errorf("failed to reconcile spell: %w", err)
		}
		if matchType == "created_personal" {
			createdCounts["spells"]++
		}
	}

	// Reconcile weapon skills
	for _, ws := range char.Waffenfertigkeiten {
		_, matchType, err := reconcileWeaponSkillWithTx(tx, ws, history.ID, char.Typ)
		if err != nil {
			history.Status = "failed"
			history.ErrorLog = fmt.Sprintf("Failed to reconcile weapon skill %s: %v", ws.Name, err)
			tx.Save(history)
			tx.Rollback()
			return nil, fmt.Errorf("failed to reconcile weapon skill: %w", err)
		}
		if matchType == "created_personal" {
			createdCounts["weapon_skills"]++
		}
	}

	// Reconcile weapons
	for _, weapon := range char.Waffen {
		_, matchType, err := reconcileWeaponWithTx(tx, weapon, history.ID, char.Typ)
		if err != nil {
			history.Status = "failed"
			history.ErrorLog = fmt.Sprintf("Failed to reconcile weapon %s: %v", weapon.Name, err)
			tx.Save(history)
			tx.Rollback()
			return nil, fmt.Errorf("failed to reconcile weapon: %w", err)
		}
		if matchType == "created_personal" {
			createdCounts["weapons"]++
		}
	}

	// Reconcile equipment
	for _, eq := range char.Ausruestung {
		_, matchType, err := reconcileEquipmentWithTx(tx, eq, history.ID, char.Typ)
		if err != nil {
			history.Status = "failed"
			history.ErrorLog = fmt.Sprintf("Failed to reconcile equipment %s: %v", eq.Name, err)
			tx.Save(history)
			tx.Rollback()
			return nil, fmt.Errorf("failed to reconcile equipment: %w", err)
		}
		if matchType == "created_personal" {
			createdCounts["equipment"]++
		}
	}

	// Reconcile containers
	for _, container := range char.Behaeltnisse {
		_, matchType, err := reconcileContainerWithTx(tx, container, history.ID, char.Typ)
		if err != nil {
			history.Status = "failed"
			history.ErrorLog = fmt.Sprintf("Failed to reconcile container %s: %v", container.Name, err)
			tx.Save(history)
			tx.Rollback()
			return nil, fmt.Errorf("failed to reconcile container: %w", err)
		}
		if matchType == "created_personal" {
			createdCounts["containers"]++
		}
	}

	result.CreatedItems = createdCounts

	// 3. Create models.Char
	// TODO: Implement CreateCharacterFromImport helper
	// For now, create a minimal character record
	newChar := &models.Char{
		Rasse:               char.Rasse,
		Typ:                 char.Typ,
		Alter:               char.Alter,
		Anrede:              char.Anrede,
		Grad:                char.Grad,
		Groesse:             char.Groesse,
		Gewicht:             char.Gewicht,
		Glaube:              char.Glaube,
		Hand:                char.Hand,
		UserID:              userID,
		ImportedFromAdapter: &adapterID,
		ImportedAt:          &history.ImportedAt,
	}

	// Set the name through BamortBase (inherited)
	newChar.BamortBase.Name = char.Name

	if err := tx.Create(newChar).Error; err != nil {
		history.Status = "failed"
		history.ErrorLog = fmt.Sprintf("Failed to create character: %v", err)
		tx.Save(history)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create character: %w", err)
	}

	// Create attributes as separate records
	attributes := []models.Eigenschaft{
		{CharacterID: newChar.ID, UserID: userID, Name: "St", Value: char.Eigenschaften.St},
		{CharacterID: newChar.ID, UserID: userID, Name: "Gs", Value: char.Eigenschaften.Gs},
		{CharacterID: newChar.ID, UserID: userID, Name: "Gw", Value: char.Eigenschaften.Gw},
		{CharacterID: newChar.ID, UserID: userID, Name: "Ko", Value: char.Eigenschaften.Ko},
		{CharacterID: newChar.ID, UserID: userID, Name: "In", Value: char.Eigenschaften.In},
		{CharacterID: newChar.ID, UserID: userID, Name: "Zt", Value: char.Eigenschaften.Zt},
		{CharacterID: newChar.ID, UserID: userID, Name: "Au", Value: char.Eigenschaften.Au},
		{CharacterID: newChar.ID, UserID: userID, Name: "Pa", Value: char.Eigenschaften.Pa},
		{CharacterID: newChar.ID, UserID: userID, Name: "Wk", Value: char.Eigenschaften.Wk},
	}

	for _, attr := range attributes {
		if err := tx.Create(&attr).Error; err != nil {
			history.Status = "failed"
			history.ErrorLog = fmt.Sprintf("Failed to create attribute %s: %v", attr.Name, err)
			tx.Save(history)
			tx.Rollback()
			return nil, fmt.Errorf("failed to create attribute: %w", err)
		}
	}

	// Create LP, AP, B records
	lp := &models.Lp{
		CharacterID: newChar.ID,
		Max:         char.Lp.Max,
		Value:       char.Lp.Value,
	}
	if err := tx.Create(lp).Error; err != nil {
		history.Status = "failed"
		history.ErrorLog = fmt.Sprintf("Failed to create LP: %v", err)
		tx.Save(history)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create LP: %w", err)
	}

	ap := &models.Ap{
		CharacterID: newChar.ID,
		Max:         char.Ap.Max,
		Value:       char.Ap.Value,
	}
	if err := tx.Create(ap).Error; err != nil {
		history.Status = "failed"
		history.ErrorLog = fmt.Sprintf("Failed to create AP: %v", err)
		tx.Save(history)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create AP: %w", err)
	}

	b := &models.B{
		CharacterID: newChar.ID,
		Max:         char.B.Max,
		Value:       char.B.Value,
	}
	if err := tx.Create(b).Error; err != nil {
		history.Status = "failed"
		history.ErrorLog = fmt.Sprintf("Failed to create B: %v", err)
		tx.Save(history)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create B: %w", err)
	}

	// Create XP record
	xp := &models.Erfahrungsschatz{
		BamortCharTrait: models.BamortCharTrait{
			CharacterID: newChar.ID,
			UserID:      userID,
		},
		EP: char.Erfahrungsschatz.Value,
	}
	if err := tx.Create(xp).Error; err != nil {
		history.Status = "failed"
		history.ErrorLog = fmt.Sprintf("Failed to create XP: %v", err)
		tx.Save(history)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create XP: %w", err)
	}

	// Create Bennies record
	bennies := &models.Bennies{
		BamortCharTrait: models.BamortCharTrait{
			CharacterID: newChar.ID,
			UserID:      userID,
		},
		Gg: char.Bennies.Gg,
		Gp: char.Bennies.Gp,
		Sg: char.Bennies.Sg,
	}
	if err := tx.Create(bennies).Error; err != nil {
		history.Status = "failed"
		history.ErrorLog = fmt.Sprintf("Failed to create bennies: %v", err)
		tx.Save(history)
		tx.Rollback()
		return nil, fmt.Errorf("failed to create bennies: %w", err)
	}

	result.CharacterID = newChar.ID

	// 4. Update ImportHistory (success status)
	history.CharacterID = &newChar.ID
	history.Status = "success"
	if err := tx.Save(history).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update import history: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		// Transaction commit failed - try to keep ImportHistory with failed status
		// This is best-effort since we're outside the transaction now
		database.DB.Model(&ImportHistory{}).
			Where("id = ?", history.ID).
			Updates(map[string]interface{}{
				"status":    "failed",
				"error_log": fmt.Sprintf("Transaction commit failed: %v", err),
			})
		return nil, err
	}

	result.Status = "success"
	return result, nil
}

// compressData compresses data using gzip
func compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	if _, err := gzipWriter.Write(data); err != nil {
		return nil, err
	}

	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// decompressData decompresses gzip data
func decompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Transaction-aware reconciliation helpers that use the provided transaction instead of database.DB

func reconcileSkillWithTx(tx *gorm.DB, skill Fertigkeit, importHistoryID uint, gameSystem string) (*models.Skill, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Skill
	err := tx.Where("name = ? AND game_system = ?", skill.Name, gs.Name).First(&existing).Error

	if err == nil {
		// Exact match found
		if importHistoryID > 0 {
			logMasterDataImportWithTx(tx, importHistoryID, "skill", existing.ID, skill.Name, "exact")
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
		SourceID:         1,
	}

	if err := tx.Create(newSkill).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create skill: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImportWithTx(tx, importHistoryID, "skill", newSkill.ID, skill.Name, "created_personal")
	}

	return newSkill, "created_personal", nil
}

func reconcileSpellWithTx(tx *gorm.DB, spell Zauber, importHistoryID uint, gameSystem string) (*models.Spell, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Spell
	err := tx.Where("name = ? AND game_system = ?", spell.Name, gs.Name).First(&existing).Error

	if err == nil {
		if importHistoryID > 0 {
			logMasterDataImportWithTx(tx, importHistoryID, "spell", existing.ID, spell.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query spell: %w", err)
	}

	newSpell := &models.Spell{
		Name:         spell.Name,
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
		Beschreibung: spell.Beschreibung,
		Quelle:       spell.Quelle,
		PersonalItem: true,
		SourceID:     1,
	}

	if err := tx.Create(newSpell).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create spell: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImportWithTx(tx, importHistoryID, "spell", newSpell.ID, spell.Name, "created_personal")
	}

	return newSpell, "created_personal", nil
}

func reconcileWeaponSkillWithTx(tx *gorm.DB, ws Waffenfertigkeit, importHistoryID uint, gameSystem string) (*models.WeaponSkill, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.WeaponSkill
	err := tx.Where("name = ? AND game_system = ?", ws.Name, gs.Name).First(&existing).Error

	if err == nil {
		if importHistoryID > 0 {
			logMasterDataImportWithTx(tx, importHistoryID, "weapon_skill", existing.ID, ws.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query weapon skill: %w", err)
	}

	newWS := &models.WeaponSkill{
		Skill: models.Skill{
			Name:             ws.Name,
			GameSystem:       gs.Name,
			GameSystemId:     gs.ID,
			Beschreibung:     ws.Beschreibung,
			Quelle:           ws.Quelle,
			Bonuseigenschaft: "check",
			Improvable:       true,
			PersonalItem:     true,
			SourceID:         1,
		},
	}

	if err := tx.Create(newWS).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create weapon skill: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImportWithTx(tx, importHistoryID, "weapon_skill", newWS.ID, ws.Name, "created_personal")
	}

	return newWS, "created_personal", nil
}

func reconcileWeaponWithTx(tx *gorm.DB, weapon Waffe, importHistoryID uint, gameSystem string) (*models.Weapon, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Weapon
	err := tx.Where("name = ? AND game_system = ?", weapon.Name, gs.Name).First(&existing).Error

	if err == nil {
		if importHistoryID > 0 {
			logMasterDataImportWithTx(tx, importHistoryID, "weapon", existing.ID, weapon.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query weapon: %w", err)
	}

	newWeapon := &models.Weapon{
		Equipment: models.Equipment{
			Name:         weapon.Name,
			GameSystem:   gs.Name,
			GameSystemId: gs.ID,
			Beschreibung: weapon.Beschreibung,
			PersonalItem: true,
			SourceID:     1,
		},
	}

	if err := tx.Create(newWeapon).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create weapon: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImportWithTx(tx, importHistoryID, "weapon", newWeapon.ID, weapon.Name, "created_personal")
	}

	return newWeapon, "created_personal", nil
}

func reconcileEquipmentWithTx(tx *gorm.DB, eq Ausruestung, importHistoryID uint, gameSystem string) (*models.Equipment, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Equipment
	err := tx.Where("name = ? AND game_system = ?", eq.Name, gs.Name).First(&existing).Error

	if err == nil {
		if importHistoryID > 0 {
			logMasterDataImportWithTx(tx, importHistoryID, "equipment", existing.ID, eq.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query equipment: %w", err)
	}

	newEq := &models.Equipment{
		Name:         eq.Name,
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
		Beschreibung: eq.Beschreibung,
		PersonalItem: true,
		SourceID:     1,
	}

	if err := tx.Create(newEq).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create equipment: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImportWithTx(tx, importHistoryID, "equipment", newEq.ID, eq.Name, "created_personal")
	}

	return newEq, "created_personal", nil
}

func reconcileContainerWithTx(tx *gorm.DB, container Behaeltniss, importHistoryID uint, gameSystem string) (*models.Container, string, error) {
	gs := models.GetGameSystem(0, gameSystem)

	var existing models.Container
	err := tx.Where("name = ? AND game_system = ?", container.Name, gs.Name).First(&existing).Error

	if err == nil {
		if importHistoryID > 0 {
			logMasterDataImportWithTx(tx, importHistoryID, "container", existing.ID, container.Name, "exact")
		}
		return &existing, "exact", nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, "", fmt.Errorf("failed to query container: %w", err)
	}

	newContainer := &models.Container{
		Equipment: models.Equipment{
			Name:         container.Name,
			GameSystem:   gs.Name,
			GameSystemId: gs.ID,
			Beschreibung: container.Beschreibung,
			PersonalItem: true,
			SourceID:     1,
		},
	}

	if err := tx.Create(newContainer).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create container: %w", err)
	}

	if importHistoryID > 0 {
		logMasterDataImportWithTx(tx, importHistoryID, "container", newContainer.ID, container.Name, "created_personal")
	}

	return newContainer, "created_personal", nil
}

func logMasterDataImportWithTx(tx *gorm.DB, importHistoryID uint, itemType string, itemID uint, externalName string, matchType string) {
	log := &MasterDataImport{
		ImportHistoryID: importHistoryID,
		ItemType:        itemType,
		ItemID:          itemID,
		ExternalName:    externalName,
		MatchType:       matchType,
		CreatedAt:       time.Now(),
	}

	if err := tx.Create(log).Error; err != nil {
		fmt.Printf("Warning: Failed to log master data import: %v\n", err)
	}
}
