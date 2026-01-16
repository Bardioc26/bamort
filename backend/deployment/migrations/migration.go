package migrations

import (
	"bamort/logger"

	"gorm.io/gorm"
)

// Migration represents a single database migration
type Migration struct {
	Number      int                  // Sequential migration number
	Version     string               // Target version (e.g., "0.5.0")
	Description string               // Human-readable description
	UpSQL       []string             // Forward migration SQL statements
	DownSQL     []string             // Rollback SQL statements
	DataFunc    func(*gorm.DB) error // Optional data migration function
	Critical    bool                 // If true, stops on error; if false, warns
}

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

// createSchemaVersionTables creates the schema_version and migration_history tables using GORM
func createSchemaVersionTables(db *gorm.DB) error {
	logger.Debug("Creating schema_version and migration_history tables using GORM")

	// Use GORM AutoMigrate for database-agnostic table creation
	if err := db.AutoMigrate(&SchemaVersion{}, &MigrationHistory{}); err != nil {
		return err
	}

	logger.Debug("Schema version tables created successfully")
	return nil
}

// AllMigrations contains all migrations in sequential order
var AllMigrations = []Migration{
	{
		Number:      1,
		Version:     "0.4.0",
		Description: "Initial schema version tracking",
		DataFunc:    createSchemaVersionTables,
		DownSQL: []string{
			"DROP TABLE IF EXISTS migration_history",
			"DROP TABLE IF EXISTS schema_version",
		},
		Critical: true,
	},
}

// GetMigrationByNumber returns a migration by its number
func GetMigrationByNumber(number int) *Migration {
	for _, m := range AllMigrations {
		if m.Number == number {
			return &m
		}
	}
	return nil
}

// GetLatestMigration returns the latest migration
func GetLatestMigration() *Migration {
	if len(AllMigrations) == 0 {
		return nil
	}
	return &AllMigrations[len(AllMigrations)-1]
}
