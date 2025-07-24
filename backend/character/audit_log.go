package character

import (
	"bamort/database"
	"time"
)

// AuditLogEntry repräsentiert einen Eintrag im Audit-Log
type AuditLogEntry struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CharacterID uint      `json:"character_id" gorm:"not null"`
	FieldName   string    `json:"field_name" gorm:"not null"` // "experience_points", "gold", "silver", "copper"
	OldValue    int       `json:"old_value"`
	NewValue    int       `json:"new_value"`
	Difference  int       `json:"difference"` // NewValue - OldValue
	Reason      string    `json:"reason"`     // "manual", "skill_learning", "skill_improvement", "spell_learning", etc.
	UserID      uint      `json:"user_id,omitempty"`
	Timestamp   time.Time `json:"timestamp" gorm:"autoCreateTime"`
	Notes       string    `json:"notes,omitempty"` // Zusätzliche Informationen
}

// AuditLogReason definiert Standard-Gründe für Änderungen
type AuditLogReason string

const (
	ReasonManual           AuditLogReason = "manual"
	ReasonSkillLearning    AuditLogReason = "skill_learning"
	ReasonSkillImprovement AuditLogReason = "skill_improvement"
	ReasonSpellLearning    AuditLogReason = "spell_learning"
	ReasonSpellImprovement AuditLogReason = "spell_improvement"
	ReasonEquipment        AuditLogReason = "equipment"
	ReasonReward           AuditLogReason = "reward"
	ReasonCorrection       AuditLogReason = "correction"
	ReasonImport           AuditLogReason = "import"
)

// CreateAuditLogEntry erstellt einen neuen Audit-Log-Eintrag
func CreateAuditLogEntry(characterID uint, fieldName string, oldValue, newValue int, reason AuditLogReason, userID uint, notes string) error {
	entry := AuditLogEntry{
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
func GetAuditLogForCharacter(characterID uint) ([]AuditLogEntry, error) {
	var entries []AuditLogEntry
	err := database.DB.Where("character_id = ?", characterID).
		Order("timestamp DESC").
		Find(&entries).Error
	return entries, err
}

// GetAuditLogForField holt Audit-Log-Einträge für ein spezifisches Feld
func GetAuditLogForField(characterID uint, fieldName string) ([]AuditLogEntry, error) {
	var entries []AuditLogEntry
	err := database.DB.Where("character_id = ? AND field_name = ?", characterID, fieldName).
		Order("timestamp DESC").
		Find(&entries).Error
	return entries, err
}

// MigrateAuditLog führt die Migration für die Audit-Log-Tabelle durch
func MigrateAuditLog() error {
	return database.DB.AutoMigrate(&AuditLogEntry{})
}
