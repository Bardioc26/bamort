package importer

import (
	"bamort/database"
	"time"

	"gorm.io/gorm"
)

// ImportHistory tracks all character import attempts with source data and mapping information
type ImportHistory struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	CharacterID     *uint     `gorm:"index" json:"character_id,omitempty"`          // NULL if import failed
	AdapterID       string    `gorm:"type:varchar(100);not null" json:"adapter_id"` // e.g., "foundry-vtt-v1"
	SourceFormat    string    `gorm:"type:varchar(50)" json:"source_format"`        // e.g., "foundry-vtt"
	SourceFilename  string    `json:"source_filename"`
	SourceSnapshot  []byte    `gorm:"type:MEDIUMBLOB" json:"-"`             // Original file (gzip compressed)
	MappingSnapshot []byte    `gorm:"type:JSON" json:"-"`                   // Adapter->BMRT mappings
	BmrtVersion     string    `gorm:"type:varchar(10)" json:"bmrt_version"` // e.g., "1.0"
	ImportedAt      time.Time `json:"imported_at"`
	Status          string    `gorm:"type:varchar(20)" json:"status"` // "in_progress", "success", "partial", "failed"
	ErrorLog        string    `gorm:"type:TEXT" json:"error_log,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// MasterDataImport tracks created/matched master data items during import
type MasterDataImport struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ImportHistoryID uint      `gorm:"not null;index" json:"import_history_id"`
	ItemType        string    `gorm:"type:varchar(20)" json:"item_type"` // "skill", "spell", "weapon", "equipment"
	ItemID          uint      `gorm:"not null" json:"item_id"`
	ExternalName    string    `json:"external_name"`
	MatchType       string    `gorm:"type:varchar(20)" json:"match_type"` // "exact", "created_personal"
	CreatedAt       time.Time `json:"created_at"`
}

// TableName specifies the table name for ImportHistory
func (ImportHistory) TableName() string {
	return "import_histories"
}

// TableName specifies the table name for MasterDataImport
func (MasterDataImport) TableName() string {
	return "master_data_imports"
}

// MigrateStructure runs migrations for importer-specific tables
func MigrateStructure(db ...*gorm.DB) error {
	// Use provided DB or default to database.DB
	var targetDB *gorm.DB
	if len(db) > 0 && db[0] != nil {
		targetDB = db[0]
	} else {
		targetDB = database.DB
	}

	err := targetDB.AutoMigrate(
		&ImportHistory{},
		&MasterDataImport{},
	)
	if err != nil {
		return err
	}

	return nil
}
