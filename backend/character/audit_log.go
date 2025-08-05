package character

import (
	"bamort/database"
	"bamort/models"
)

// AuditLogReason definiert Standard-Gründe für Änderungen
type AuditLogReason string

const (
	ReasonManual           AuditLogReason = "manual"
	ReasonSkillLearning    AuditLogReason = "skill_learning"
	ReasonSkillImprovement AuditLogReason = "skill_improvement"
	ReasonSpellLearning    AuditLogReason = "spell_learning"
	ReasonEquipment        AuditLogReason = "equipment"
	ReasonReward           AuditLogReason = "reward"
	ReasonCorrection       AuditLogReason = "correction"
	ReasonImport           AuditLogReason = "import"
)

// CreateAuditLogEntry erstellt einen neuen Audit-Log-Eintrag
func CreateAuditLogEntry(characterID uint, fieldName string, oldValue, newValue int, reason AuditLogReason, userID uint, notes string) error {
	entry := models.AuditLogEntry{
		CharacterID: characterID,
		FieldName:   fieldName,
		OldValue:    oldValue,
		NewValue:    newValue,
		Difference:  newValue - oldValue,
		Reason:      string(reason),
		UserID:      userID,
		Notes:       notes,
	}

	return database.DB.Create(&entry).Error
}

// GetAuditLogForCharacter holt alle Audit-Log-Einträge für einen Charakter
func GetAuditLogForCharacter(characterID uint) ([]models.AuditLogEntry, error) {
	var entries []models.AuditLogEntry
	err := database.DB.Where("character_id = ?", characterID).
		Order("timestamp DESC").
		Find(&entries).Error
	return entries, err
}

// GetAuditLogForField holt Audit-Log-Einträge für ein spezifisches Feld
func GetAuditLogForField(characterID uint, fieldName string) ([]models.AuditLogEntry, error) {
	var entries []models.AuditLogEntry
	err := database.DB.Where("character_id = ? AND field_name = ?", characterID, fieldName).
		Order("timestamp DESC").
		Find(&entries).Error
	return entries, err
}

// MigrateAuditLog führt die Migration für die Audit-Log-Tabelle durch
func MigrateAuditLog() error {
	return database.DB.AutoMigrate(&models.AuditLogEntry{})
}
