package database

import (
	"bamort/logger"
	"io"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var isTestDb bool
var testdbTempDir string

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	logger.Debug("copyFile: Kopiere Datei von %s nach %s", src, dst)

	sourceFile, err := os.Open(src)
	if err != nil {
		logger.Error("copyFile: Fehler beim Öffnen der Quelldatei %s: %s", src, err.Error())
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		logger.Error("copyFile: Fehler beim Erstellen der Zieldatei %s: %s", dst, err.Error())
		return err
	}
	defer destFile.Close()

	copied, err := io.Copy(destFile, sourceFile)
	if err != nil {
		logger.Error("copyFile: Fehler beim Kopieren der Datei: %s", err.Error())
		return err
	}

	logger.Debug("copyFile: Erfolgreich %d Bytes kopiert von %s nach %s", copied, src, dst)
	return err
}

// The Database for testing is created from the live database whenever needed and stored in the path defined in database.PreparedTestDB
// to use the test database make a temporary copy of it and then open this new copy as testing database
// This allows to have a clean database for each test run without affecting the live database
// However SetupTestDB can still open the live database if required by setting isTestDb to false
// SetupTestDB creates an in-memory SQLite database for testing
// Parameters:
// - opts[0]: isTestDb (bool) - whether to use precopied SQLite (true) or persistent (Live) MariaDB (false)
func SetupTestDB(opts ...bool) {
	logger.Debug("SetupTestDB aufgerufen")

	isTestDb = true

	if len(opts) > 0 {
		isTestDb = opts[0]
		logger.Debug("SetupTestDB: isTestDb Parameter überschrieben auf %t", isTestDb)
	}

	logger.Debug("SetupTestDB: Verwende Test-Datenbank: %t", isTestDb)

	if DB == nil {
		logger.Debug("SetupTestDB: DB ist nil, erstelle neue Datenbankverbindung")
		var db *gorm.DB

		if isTestDb {
			logger.Info("SetupTestDB: Erstelle SQLite Test-Datenbank")

			testdbTempDir, err := os.MkdirTemp("", "bamort-test-")
			if err != nil {
				logger.Error("SetupTestDB: Fehler beim Erstellen des temporären Verzeichnisses: %s", err.Error())
				panic("failed to create temporary directory: " + err.Error())
			}
			logger.Debug("SetupTestDB: Temporäres Verzeichnis erstellt: %s", testdbTempDir)

			targetFile := filepath.Join(testdbTempDir, "test_backup.db")
			logger.Debug("SetupTestDB: Ziel-Datei: %s", targetFile)
			logger.Debug("SetupTestDB: Quelle-Datei: %s", PreparedTestDB)

			err = copyFile(PreparedTestDB, targetFile)
			if err != nil {
				logger.Error("SetupTestDB: Fehler beim Kopieren der Test-Datenbank: %s", err.Error())
				panic("failed to copy prepared test database: " + err.Error())
			}
			logger.Info("SetupTestDB: Test-Datenbank erfolgreich kopiert")

			db, err = gorm.Open(sqlite.Open(targetFile), &gorm.Config{})
			if err != nil {
				logger.Error("SetupTestDB: Fehler beim Verbinden mit der Test-Datenbank: %s", err.Error())
				panic("failed to connect to the test database: " + err.Error())
			}
			logger.Info("SetupTestDB: Erfolgreich mit SQLite Test-Datenbank verbunden")
			//defer os.RemoveAll(testdbTempDir)
		} else {
			logger.Info("SetupTestDB: Verwende Live-Datenbank (MariaDB)")
			//* //testing with persistent MariaDB
			db = ConnectDatabase()
			if db == nil {
				logger.Error("SetupTestDB: Fehler beim Verbinden mit der Live-Datenbank")
				panic("failed to connect to the live database")
			}
			logger.Info("SetupTestDB: Erfolgreich mit Live-Datenbank verbunden")
		}
		DB = db
		logger.Info("SetupTestDB: Datenbankverbindung erfolgreich eingerichtet")
	} else {
		logger.Debug("SetupTestDB: DB bereits initialisiert, überspringe Setup")
	}
}
func ResetTestDB() {
	logger.Debug("ResetTestDB aufgerufen")

	if isTestDb {
		logger.Debug("ResetTestDB: Verwende Test-Datenbank, führe Cleanup durch")

		// Check if DB is not nil before trying to use it
		if DB != nil {
			sqlDB, err := DB.DB()
			if err == nil {
				logger.Debug("ResetTestDB: Schließe Datenbankverbindung")
				sqlDB.Close()
			} else {
				logger.Error("ResetTestDB: Fehler beim Abrufen der SQL-Datenbank: %s", err.Error())
			}
		} else {
			logger.Debug("ResetTestDB: DB ist bereits nil, überspringe Verbindungsschließung")
		}

		// Always set DB to nil and clean up temp directory
		DB = nil
		logger.Debug("ResetTestDB: DB auf nil gesetzt")

		if testdbTempDir != "" {
			logger.Debug("ResetTestDB: Lösche temporäres Verzeichnis: %s", testdbTempDir)
			err := os.RemoveAll(testdbTempDir)
			if err != nil {
				logger.Error("ResetTestDB: Fehler beim Löschen des temporären Verzeichnisses: %s", err.Error())
			} else {
				logger.Info("ResetTestDB: Temporäres Verzeichnis erfolgreich gelöscht")
			}
			testdbTempDir = ""
		}
	} else {
		logger.Debug("ResetTestDB: Verwende Live-Datenbank, überspringe Cleanup")
	}

	logger.Info("ResetTestDB: Cleanup abgeschlossen")
}
