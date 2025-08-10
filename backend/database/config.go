package database

import (
	"bamort/config"
	"bamort/logger"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"path/filepath"
	"runtime"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// getBackendDir returns the absolute path to the backend directory
func getBackendDir() string {
	_, filename, _, _ := runtime.Caller(0)
	// This file is in backend/database/config.go
	// So we need to go up one level to get to backend/
	return filepath.Dir(filepath.Dir(filename))
}

// Test database configuration paths
var (
	// PreparedTestDB is the path to the prepared test database file
	// This file contains a snapshot of test data that can be loaded into test databases
	// Usage: database.PreparedTestDB to access from any package
	PreparedTestDB = filepath.Join(getBackendDir(), "testdata", "prepared_test_data.db")

	// TestDataDir is the directory for maintenance test data
	// This directory contains temporary test data files during test execution
	// Usage: database.TestDataDir to access from any package
	TestDataDir = filepath.Join(getBackendDir(), "maintenance", "testdata")
)

func ConnectDatabase() *gorm.DB {
	// Konfiguration laden um zu entscheiden, welche Datenbank verwendet werden soll
	cfg := config.LoadConfig()

	// In Test-Umgebung verwende Test-DB, sonst die konfigurierte Datenbank
	if cfg.Environment == "test" || cfg.Testing == "yes" {
		logger.Debug("Test-Umgebung erkannt, verwende Test-Datenbank")
		SetupTestDB()
	} else {
		logger.Debug("Verwende konfigurierte Datenbank (%s)", cfg.DatabaseType)
		return ConnectDatabaseOrig()
	}

	return DB
}
func ConnectDatabaseOrig() *gorm.DB {
	// Konfiguration laden
	cfg := config.LoadConfig()

	logger.Debug("Datenbank-Konfiguration geladen: Type=%s, URL=%s", cfg.DatabaseType, cfg.DatabaseURL)

	// Falls keine URL konfiguriert ist, verwende Standard-MySQL-Konfiguration als Fallback
	dbURL := cfg.DatabaseURL
	if dbURL == "" {
		dbURL = "bamort:bG4)efozrc@tcp(192.168.0.5:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"
		logger.Warn("Keine DATABASE_URL konfiguriert, verwende Standard-MySQL-Konfiguration")
		cfg.DatabaseType = "mysql"
	}

	var database *gorm.DB
	var err error

	// Datenbanktyp-spezifische Verbindung
	switch cfg.DatabaseType {
	case "mysql":
		logger.Debug("Verbinde mit MySQL-Datenbank...")
		database, err = gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	case "sqlite":
		logger.Debug("Verbinde mit SQLite-Datenbank...")
		database, err = gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	default:
		logger.Error("Nicht unterstützter Datenbanktyp: %s. Verwende MySQL als Fallback.", cfg.DatabaseType)
		database, err = gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	}

	if err != nil {
		logger.Error("Fehler beim Verbinden mit der Datenbank (%s): %s", cfg.DatabaseType, err.Error())
		log.Fatal("Failed to connect to database:", err)
	}

	logger.Info("Erfolgreich mit %s-Datenbank verbunden", cfg.DatabaseType)
	DB = database
	return DB
}

func GetDB() *gorm.DB {
	if DB == nil {
		DB = ConnectDatabase()
	}
	return DB
}

// StringArray is a custom type for []string
type StringArray []string

// Value implements the driver.Valuer interface for database storage
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s) // Serialize []string to JSON
}

// Scan implements the sql.Scanner interface for database retrieval
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to convert database value to []byte")
	}

	return json.Unmarshal(bytes, s) // Deserialize JSON to []string
}
