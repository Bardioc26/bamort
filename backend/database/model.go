package database

// SchemaVersion represents the schema_version table
type SchemaVersion struct {
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	Version         string `gorm:"size:20;not null;index"`
	MigrationNumber int    `gorm:"not null;index"`
	AppliedAt       int64  `gorm:"autoCreateTime"`
	BackendVersion  string `gorm:"size:20;not null"`
	Description     string `gorm:"type:text"`
	Checksum        string `gorm:"size:64"`
}

// TableName sets the table name for SchemaVersion
func (SchemaVersion) TableName() string {
	return "schema_version"
}

// MigrationHistory represents the migration_history table
type MigrationHistory struct {
	ID                uint   `gorm:"primaryKey;autoIncrement"`
	MigrationNumber   int    `gorm:"not null;uniqueIndex"`
	Version           string `gorm:"size:20;not null;index"`
	Description       string `gorm:"type:text;not null"`
	AppliedAt         int64  `gorm:"autoCreateTime"`
	AppliedBy         string `gorm:"size:100"`
	ExecutionTimeMs   int64
	Success           bool   `gorm:"default:true"`
	ErrorMessage      string `gorm:"type:text"`
	RollbackAvailable bool   `gorm:"default:true"`
}

// TableName sets the table name for MigrationHistory
func (MigrationHistory) TableName() string {
	return "migration_history"
}
