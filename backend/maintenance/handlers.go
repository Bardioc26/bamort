package maintenance

import (
	"bamort/config"
	"bamort/database"
	"bamort/gamesystem"
	"bamort/logger"
	"bamort/models"
	"bamort/user"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Constants for test data management
var (
	testDataDir    = database.TestDataDir
	preparedTestDB = database.PreparedTestDB
)

func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// migrateAllStructures migrates all database structures to the provided database
func migrateAllStructures(db *gorm.DB) error {
	logger.Debug("Starte Migration aller Datenbankstrukturen...")

	// Migrate all structures in the correct order
	logger.Debug("Migriere Datenbankstrukturen...")
	if err := database.MigrateStructure(db); err != nil {
		logger.Error("Fehler beim Migrieren der Datenbankstrukturen: %s", err.Error())
		return fmt.Errorf("failed to migrate database structures: %w", err)
	}

	// Migrate all structures in the correct order
	logger.Debug("Migriere GameSystemstrukturen...")
	if err := gamesystem.MigrateStructure(db); err != nil {
		logger.Error("Fehler beim Migrieren der GameSystemstrukturen: %s", err.Error())
		return fmt.Errorf("failed to migrate game system structures: %w", err)
	}

	logger.Debug("Migriere Benutzerstrukturen...")
	if err := user.MigrateStructure(db); err != nil {
		logger.Error("Fehler beim Migrieren der Benutzerstrukturen: %s", err.Error())
		return fmt.Errorf("failed to migrate user structures: %w", err)
	}

	logger.Debug("Migriere GSMaster-Strukturen...")
	if err := models.MigrateStructure(db); err != nil {
		logger.Error("Fehler beim Migrieren der GSMaster-Strukturen: %s", err.Error())
		return fmt.Errorf("failed to migrate gsmaster structures: %w", err)
	}

	/*if err := importer.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate importer structures: %w", err)
	}*/

	logger.Info("Migration aller Datenbankstrukturen erfolgreich abgeschlossen")
	return nil
}

func migrateDataIfNeeded(db *gorm.DB) error {
	logger.Debug("Starte Datenmigration falls erforderlich...")

	err := database.MigrateDataIfNeeded(db)
	if err != nil {
		logger.Error("Fehler beim Migrieren der Datenbankdaten: %s", err.Error())
		return fmt.Errorf("failed to migrate database data: %w", err)
	}
	err = gamesystem.MigrateDataIfNeeded(db)
	if err != nil {
		logger.Error("Fehler beim Migrieren der GameSystem-Daten: %s", err.Error())
		return fmt.Errorf("failed to migrate game system data: %w", err)
	}
	err = models.MigrateDataIfNeeded(db)
	if err != nil {
		logger.Error("Fehler beim Migrieren der Models-Daten: %s", err.Error())
		return fmt.Errorf("failed to migrate models data: %w", err)
	}

	// Kopiere categorie nach learning_category für Spells, wenn learning_category leer ist
	logger.Debug("Migriere Spell Learning Categories...")
	err = migrateSpellLearningCategories(db)
	if err != nil {
		logger.Error("Fehler beim Migrieren der Spell Learning Categories: %s", err.Error())
		return fmt.Errorf("failed to migrate spell learning categories: %w", err)
	}

	logger.Info("Datenmigration erfolgreich abgeschlossen")
	return nil
}

// migrateSpellLearningCategories kopiert categorie-Werte in learning_category wenn diese leer sind
func migrateSpellLearningCategories(db *gorm.DB) error {
	logger.Debug("Starte Migration der Spell Learning Categories...")

	// SQL-Statement um categorie nach learning_category zu kopieren, wo learning_category leer oder NULL ist
	sql := `
		UPDATE gsm_spells 
		SET learning_category = category 
		WHERE (learning_category IS NULL OR learning_category = '') 
		AND category IS NOT NULL 
		AND category != ''
	`

	logger.Debug("Führe SQL-Update aus: %s", strings.ReplaceAll(sql, "\n", " "))
	result := db.Exec(sql)
	if result.Error != nil {
		logger.Error("Fehler beim SQL-Update der Spell Learning Categories: %s", result.Error.Error())
		return fmt.Errorf("failed to update spell learning categories: %w", result.Error)
	}

	// Log der Anzahl der aktualisierten Datensätze
	if result.RowsAffected > 0 {
		logger.Info("Updated %d spell records with learning_category from categorie", result.RowsAffected)
		fmt.Printf("Updated %d spell records with learning_category from categorie\n", result.RowsAffected)
	} else {
		logger.Debug("Keine Spell-Datensätze benötigten ein Update der learning_category")
	}

	return nil
}

func MakeTestdataFromLive(c *gin.Context) {
	logger.Info("Starte Testdaten-Erstellung aus Live-Datenbank...")

	liveDB := database.ConnectDatabase()
	if liveDB == nil {
		logger.Error("Fehler beim Verbinden mit der Live-Datenbank")
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to live database")
		return
	}
	logger.Debug("Erfolgreich mit Live-Datenbank verbunden")

	// Live-Datenbank in SQLite-Datei kopieren
	backupFile := preparedTestDB
	logger.Info("Kopiere Live-Datenbank nach: %s", backupFile)
	err := copyLiveDatabaseToFile(liveDB, backupFile)
	if err != nil {
		logger.Error("Fehler beim Kopieren der Datenbank: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to copy database: %v", err))
		return
	}

	logger.Info("Live-Datenbank erfolgreich in Datei kopiert: %s", backupFile)
	c.JSON(http.StatusOK, gin.H{
		"message":        "Live database copied to file successfully",
		"test_data_file": backupFile,
	})
}

// CopyLiveDatabaseToFile kopiert die MariaDB-Datenbank in eine SQLite-Datei (exported for testing)
func CopyLiveDatabaseToFile(liveDB *gorm.DB, targetFile string) error {
	return copyLiveDatabaseToFile(liveDB, targetFile)
}

// copyLiveDatabaseToFile kopiert die MariaDB-Datenbank in eine SQLite-Datei
func copyLiveDatabaseToFile(liveDB *gorm.DB, targetFile string) error {
	logger.Debug("Starte Kopiervorgang von Live-DB nach SQLite-Datei: %s", targetFile)

	// Verzeichnis erstellen falls es nicht existiert
	dir := filepath.Dir(targetFile)
	logger.Debug("Erstelle Zielverzeichnis falls erforderlich: %s", dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Fehler beim Erstellen des Verzeichnisses %s: %s", dir, err.Error())
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Backup der existierenden Datei erstellen
	if _, err := os.Stat(targetFile); err == nil {
		backupFile := targetFile + ".backup"
		logger.Debug("Existierende Datei gefunden, erstelle Backup: %s", backupFile)
		os.Remove(backupFile) // Alte Backup entfernen
		if err := os.Rename(targetFile, backupFile); err != nil {
			logger.Error("Fehler beim Erstellen des Backups %s: %s", backupFile, err.Error())
			return fmt.Errorf("failed to backup existing file: %w", err)
		}
		logger.Debug("Backup erfolgreich erstellt")
	}

	// SQLite-Zieldatenbank erstellen
	logger.Debug("Erstelle neue SQLite-Zieldatenbank: %s", targetFile)
	targetDB, err := gorm.Open(sqlite.Open(targetFile), &gorm.Config{})
	if err != nil {
		logger.Error("Fehler beim Erstellen der SQLite-Zieldatenbank: %s", err.Error())
		return fmt.Errorf("failed to create target SQLite database: %w", err)
	}
	defer func() {
		if sqlDB, err := targetDB.DB(); err == nil {
			logger.Debug("Schließe SQLite-Datenbankverbindung")
			sqlDB.Close()
		}
	}()

	// Strukturen in SQLite-DB migrieren
	logger.Debug("Migriere Strukturen in SQLite-Datenbank...")
	if err := migrateAllStructures(targetDB); err != nil {
		logger.Error("Fehler beim Migrieren der Strukturen in SQLite: %s", err.Error())
		return fmt.Errorf("failed to migrate structures to SQLite: %w", err)
	}

	// Daten von MariaDB zu SQLite kopieren
	logger.Info("Kopiere Daten von MariaDB zu SQLite...")
	if err := copyMariaDBToSQLite(liveDB, targetDB); err != nil {
		logger.Error("Fehler beim Kopieren der Daten von MariaDB zu SQLite: %s", err.Error())
		return fmt.Errorf("failed to copy data from MariaDB to SQLite: %w", err)
	}

	logger.Info("Kopiervorgang erfolgreich abgeschlossen")
	return nil
}

// copyMariaDBToSQLite kopiert alle Daten von MariaDB zu SQLite
func copyMariaDBToSQLite(mariaDB, sqliteDB *gorm.DB) error {
	logger.Debug("Starte Kopiervorgang aller Daten von MariaDB zu SQLite...")

	// Vollständige Liste aller Strukturen mit GORM-Tags in der richtigen Reihenfolge
	// (Basis-Tabellen zuerst wegen Foreign Key-Abhängigkeiten)
	tables := []interface{}{
		// Basis-Strukturen (keine Abhängigkeiten)
		&database.MigrationHistory{},
		&database.SchemaVersion{},
		&user.User{},

		// Game System - Basis
		&models.GameSystem{},

		// Learning Costs System - Basis
		&models.Source{},
		&models.CharacterClass{},
		&models.SkillCategory{},
		&models.SkillDifficulty{},
		&models.SpellSchool{},
		&models.MiscLookup{},

		// Learning Costs System - Abhängige Tabellen
		&models.ClassCategoryEPCost{},
		&models.ClassSpellSchoolEPCost{},
		&models.SpellLevelLECost{},
		&models.SkillCategoryDifficulty{},
		&models.WeaponSkillCategoryDifficulty{},
		&models.SkillImprovementCost2{},
		&models.ClassCategoryLearningPoints{},
		&models.ClassSpellPoints{},
		&models.ClassTypicalSkill{},
		&models.ClassTypicalSpell{},

		// GSMaster Basis-Daten
		//&models.LookupList{}, // Basis für Skills, Spells, Equipment
		&models.Skill{},
		&models.WeaponSkill{},
		&models.Spell{},
		&models.Equipment{},
		&models.Weapon{},
		&models.Container{},
		&models.Transportation{},
		&models.Believe{},

		// Charaktere (Basis)
		&models.Char{},

		// Charakter-Eigenschaften (abhängig von Char)
		&models.Eigenschaft{},
		&models.Lp{},
		&models.Ap{},
		&models.B{},
		&models.Merkmale{},
		&models.Erfahrungsschatz{},
		&models.Bennies{},
		&models.Vermoegen{},

		// Charakter-Skills (abhängig von Char und Skills)
		&models.SkFertigkeit{},
		&models.SkWaffenfertigkeit{},
		&models.SkAngeboreneFertigkeit{},
		&models.SkZauber{},

		// Charakter-Equipment (abhängig von Char und Equipment)
		&models.EqAusruestung{},
		&models.EqWaffe{},
		&models.EqContainer{},

		// Character Creation Sessions (abhängig von Char)
		&models.CharacterCreationSession{},

		// Audit Logging (abhängig von Char)
		&models.AuditLogEntry{},

		// View-Strukturen ohne eigene Tabellen werden nicht kopiert:
		// SkillLearningInfo, SpellLearningInfo, CharList, FeChar, etc.
	}

	logger.Info("Kopiere Daten für %d Tabellen...", len(tables))
	for i, model := range tables {
		logger.Debug("Kopiere Tabelle %d/%d: %T", i+1, len(tables), model)
		if err := copyTableData(mariaDB, sqliteDB, model); err != nil {
			logger.Error("Fehler beim Kopieren der Tabellendaten für %T: %s", model, err.Error())
			return fmt.Errorf("failed to copy table data for %T: %w", model, err)
		}
	}

	logger.Info("Alle Tabellendaten erfolgreich kopiert")
	return nil
}

// copyTableData kopiert alle Daten einer Tabelle von MariaDB zu SQLite
func copyTableData(sourceDB, targetDB *gorm.DB, model interface{}) error {
	tableName := fmt.Sprintf("%T", model)
	logger.Debug("Starte Kopiervorgang für Tabelle: %s", tableName)

	// Anzahl der Datensätze prüfen
	var count int64
	err := sourceDB.Model(model).Count(&count).Error
	if err != nil {
		// If table doesn't exist, skip silently (useful for testing with partial schemas)
		if isTableNotExistError(err) {
			logger.Debug("Tabelle %s existiert nicht in der Quelle, überspringe", tableName)
			return nil
		}
		logger.Error("Fehler beim Zählen der Datensätze für %s: %s", tableName, err.Error())
		return err
	}

	if count == 0 {
		logger.Debug("Tabelle %s ist leer, keine Daten zu kopieren", tableName)
		return nil // Keine Daten zu kopieren
	}

	logger.Debug("Kopiere %d Datensätze für Tabelle %s", count, tableName)

	// Daten in Blöcken kopieren (für große Tabellen)
	batchSize := 100
	totalBatches := (int(count) + batchSize - 1) / batchSize

	// Get the element type for creating slice of records
	modelType := reflect.TypeOf(model).Elem()

	for offset := 0; offset < int(count); offset += batchSize {
		batchNum := (offset / batchSize) + 1
		logger.Debug("Kopiere Batch %d/%d für %s (Offset: %d, Limit: %d)", batchNum, totalBatches, tableName, offset, batchSize)

		// Create a slice of the model type using reflection
		sliceType := reflect.SliceOf(modelType)
		recordsValue := reflect.MakeSlice(sliceType, 0, batchSize)
		recordsPtr := reflect.New(sliceType)
		recordsPtr.Elem().Set(recordsValue)

		// Batch aus MariaDB lesen (use proper struct type instead of map)
		if err := sourceDB.Model(model).Offset(offset).Limit(batchSize).Find(recordsPtr.Interface()).Error; err != nil {
			logger.Error("Fehler beim Lesen von Batch %d für %s: %s", batchNum, tableName, err.Error())
			return err
		}

		// Get the records for iteration
		recordsVal := recordsPtr.Elem()
		if recordsVal.Len() == 0 {
			logger.Debug("Keine weiteren Datensätze für %s", tableName)
			break
		}

		// Batch in SQLite einfügen
		// Use Save() with SkipHooks to preserve raw values and avoid callbacks that rely on global DB state
		db := targetDB.Session(&gorm.Session{SkipHooks: true})
		for i := 0; i < recordsVal.Len(); i++ {
			record := recordsVal.Index(i).Addr().Interface()
			if err := db.Save(record).Error; err != nil {
				logger.Error("Fehler beim Speichern von Datensatz in Batch %d für %s: %s", batchNum, tableName, err.Error())
				return err
			}
		}

		logger.Debug("Batch %d/%d für %s erfolgreich kopiert (%d Datensätze)", batchNum, totalBatches, tableName, recordsVal.Len())
	}

	logger.Info("Tabelle %s erfolgreich kopiert (%d Datensätze total)", tableName, count)
	return nil
}

// isTableNotExistError checks if the error indicates a table doesn't exist
func isTableNotExistError(err error) bool {
	errorMsg := err.Error()
	return strings.Contains(errorMsg, "no such table") ||
		strings.Contains(errorMsg, "doesn't exist") ||
		strings.Contains(errorMsg, "Table") && strings.Contains(errorMsg, "doesn't exist")
}

// LoadPredefinedTestDataFromFile loads predefined test data from a specific file into the provided database
func LoadPredefinedTestDataFromFile(targetDB *gorm.DB, dataFile string) error {
	logger.Debug("Lade vordefinierte Testdaten aus Datei: %s", dataFile)

	// Check if file exists
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		logger.Error("Vordefinierte Testdaten-Datei nicht gefunden: %s", dataFile)
		return fmt.Errorf("predefined test data file not found: %s", dataFile)
	}
	logger.Debug("Testdaten-Datei existiert: %s", dataFile)

	// Migrate structures to target DB
	logger.Debug("Migriere Strukturen in Zieldatenbank...")
	err := migrateAllStructures(targetDB)
	if err != nil {
		logger.Error("Fehler beim Migrieren der Strukturen: %s", err.Error())
		return fmt.Errorf("failed to migrate structures: %w", err)
	}

	// Copy data from file database to target database
	logger.Info("Kopiere Testdaten in Zieldatenbank...")
	err = copyDataFromFileToMemory(dataFile, targetDB)
	if err != nil {
		logger.Error("Fehler beim Kopieren der Testdaten: %s", err.Error())
		return fmt.Errorf("failed to copy test data to database: %w", err)
	}

	logger.Info("Vordefinierte Testdaten erfolgreich geladen")
	return nil
}

// LoadPredefinedTestData creates a new in-memory test database and loads predefined test data into it (HTTP handler)
// Todo I think this don't need to be a handler, but can be called directly
func LoadPredefinedTestData(c *gin.Context) {

	testDataFile := filepath.Join(testDataDir, "predefined_test_data.db")

	// Check if file exists
	if _, err := os.Stat(testDataFile); os.IsNotExist(err) {
		respondWithError(c, http.StatusNotFound, "Predefined test data file not found. Run MakeTestdataFromLive first.")
		return
	}

	// Create new in-memory test database using SetupTestDB
	database.SetupTestDB(true)

	// Load test data using the predefined test data file (includes migrations)
	err := LoadPredefinedTestDataFromFile(database.DB, preparedTestDB)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to load test data: "+err.Error())
		return
	}

	// Get statistics about the loaded data
	stats, err := getTestDataStatistics(database.DB)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get test data statistics: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Predefined test data loaded successfully into in-memory database",
		"test_data_file": testDataFile,
		"statistics":     stats,
	})
}

// copyDataFromFileToMemory copies data from a SQLite file to an in-memory database
func copyDataFromFileToMemory(sourceFile string, targetDB *gorm.DB) error {
	logger.Debug("Kopiere Daten von SQLite-Datei in Memory-Datenbank: %s", sourceFile)

	// Copy all tables using ATTACH and INSERT
	attachSQL := fmt.Sprintf("ATTACH DATABASE '%s' AS source", sourceFile)
	logger.Debug("Hänge Quell-Datenbank an: %s", attachSQL)
	if err := targetDB.Exec(attachSQL).Error; err != nil {
		logger.Error("Fehler beim Anhängen der Quell-Datenbank: %s", err.Error())
		return fmt.Errorf("failed to attach source database: %w", err)
	}

	// Get list of tables from source database
	logger.Debug("Ermittle Tabellenliste aus Quell-Datenbank...")
	var tables []string
	if err := targetDB.Raw("SELECT name FROM source.sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables).Error; err != nil {
		logger.Error("Fehler beim Ermitteln der Tabellenliste: %s", err.Error())
		return fmt.Errorf("failed to get table list: %w", err)
	}
	logger.Info("Gefundene Tabellen zum Kopieren: %d (%v)", len(tables), tables)

	// Copy each table
	for i, table := range tables {
		logger.Debug("Kopiere Tabelle %d/%d: %s", i+1, len(tables), table)
		copySQL := fmt.Sprintf("INSERT OR REPLACE INTO main.%s SELECT * FROM source.%s", table, table)
		if err := targetDB.Exec(copySQL).Error; err != nil {
			logger.Error("Fehler beim Kopieren der Tabelle %s: %s", table, err.Error())
			return fmt.Errorf("failed to copy table %s: %w", table, err)
		}
		logger.Debug("Tabelle %s erfolgreich kopiert", table)
	}

	// Detach the source database
	logger.Debug("Löse Quell-Datenbank-Verbindung...")
	if err := targetDB.Exec("DETACH DATABASE source").Error; err != nil {
		logger.Error("Fehler beim Lösen der Quell-Datenbank-Verbindung: %s", err.Error())
		return fmt.Errorf("failed to detach source database: %w", err)
	}

	logger.Info("Daten erfolgreich von Datei in Memory-Datenbank kopiert")
	return nil
}

// getTestDataStatistics returns statistics about the test database
func getTestDataStatistics(db *gorm.DB) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Count records in each table
	tables := map[string]interface{}{
		"users":                     &user.User{},
		"characters":                &models.Char{},
		"gsmaster_skills":           &models.Skill{},
		"gsmaster_spells":           &models.Spell{},
		"gsmaster_equipment":        &models.Equipment{},
		"skills_fertigkeiten":       &models.SkFertigkeit{},
		"skills_waffenfertigkeiten": &models.SkWaffenfertigkeit{},
		"skills_zauber":             &models.SkZauber{},
		"equipment_ausruestung":     &models.EqAusruestung{},
		"equipment_waffen":          &models.EqWaffe{},
	}

	for name, model := range tables {
		var count int64
		if err := db.Model(model).Count(&count).Error; err != nil {
			return stats, fmt.Errorf("failed to count %s: %w", name, err)
		}
		stats[name] = count
	}

	return stats, nil
}
func setupCheck(c *gin.Context, db *gorm.DB) {
	logger.Debug("Führe Strukturmigration durch...")
	err := migrateAllStructures(db)
	if err != nil {
		logger.Error("Fehler bei der Strukturmigration: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Debug("Führe Datenmigration durch...")
	err = migrateDataIfNeeded(db)
	if err != nil {
		logger.Error("Fehler bei der Datenmigration: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to migrate data: " + err.Error()})
		return
	}

	logger.Info("Setup-Check erfolgreich abgeschlossen")
	c.JSON(http.StatusOK, gin.H{"message": "Setup Check OK"})

}

func SetupCheck(c *gin.Context) {
	logger.Info("Starte Setup-Check...")

	db := database.ConnectDatabase()
	if db == nil {
		logger.Error("Fehler beim Verbinden mit der Datenbank für Setup-Check")
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to DataBase")
		return
	}
	logger.Debug("Erfolgreich mit Datenbank für Setup-Check verbunden")
	setupCheck(c, db)
}

func SetupCheckDev(c *gin.Context) {
	logger.Info("Starte Setup-Check... PreparedTestDB")

	// Use the prepared test database for development setup check
	db, dberr := gorm.Open(sqlite.Open(database.PreparedTestDB), &gorm.Config{})
	if dberr != nil {
		logger.Error("SetupTestDB: Fehler beim Verbinden mit der Test-Datenbank: %s", dberr.Error())
		panic("failed to connect to the test database: " + dberr.Error())
	}
	database.DB = db

	logger.Debug("Erfolgreich mit Datenbank für Setup-Check verbunden")
	setupCheck(c, db)
}

/*
// PopulateClassLearningPoints populates the class learning points tables from hardcoded data
func PopulateClassLearningPoints(c *gin.Context) {
	logger.Info("Starte Population der Class Learning Points Daten...")

	err := models.PopulateClassLearningPointsData()
	if err != nil {
		logger.Error("Fehler beim Populieren der Class Learning Points: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to populate class learning points: "+err.Error())
		return
	}

	logger.Info("Class Learning Points erfolgreich populiert")
	c.JSON(http.StatusOK, gin.H{"message": "Class learning points data populated successfully"})
}
*/

func ReconnectDataBase(c *gin.Context) {
	logger.Info("Führe Datenbank-Reconnect durch...")

	db := database.ConnectDatabase()
	if db == nil {
		logger.Error("Fehler beim Reconnect zur Datenbank")
		respondWithError(c, http.StatusInternalServerError, "Failed to reconnect to DataBase")
		return
	}

	logger.Info("Datenbank-Reconnect erfolgreich")
	c.JSON(http.StatusOK, gin.H{"message": "Database reconnected successfully"})
}

func ReloadENV(c *gin.Context) {
	logger.Info("Starte Reload der Umgebungsvariablen...")

	// Reload the environment variables
	config.LoadConfig()
	c.JSON(http.StatusOK, gin.H{"message": "Environment variables reloaded successfully"})
}

// TransferSQLiteToMariaDB transfers data from SQLite test database to MariaDB
func TransferSQLiteToMariaDB(c *gin.Context) {
	logger.Info("Starte Datenübertragung von SQLite zu MariaDB...")

	// Path to the SQLite source database
	sourceFile := preparedTestDB

	// Check if source file exists
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		logger.Error("SQLite-Quelldatei nicht gefunden: %s", sourceFile)
		respondWithError(c, http.StatusNotFound, "SQLite source file not found: "+sourceFile)
		return
	}

	logger.Debug("SQLite-Quelldatei gefunden: %s", sourceFile)

	// Connect to SQLite source database
	logger.Debug("Verbinde mit SQLite-Quelldatenbank...")
	sourceDB, err := gorm.Open(sqlite.Open(sourceFile), &gorm.Config{})
	if err != nil {
		logger.Error("Fehler beim Verbinden mit SQLite-Datenbank: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to SQLite source: "+err.Error())
		return
	}
	defer func() {
		if sqlDB, err := sourceDB.DB(); err == nil {
			logger.Debug("Schließe SQLite-Datenbankverbindung")
			sqlDB.Close()
		}
	}()
	logger.Debug("SQLite-Verbindung erfolgreich")

	// Connect to MariaDB target using the configured connection string
	logger.Debug("Verbinde mit MariaDB-Zieldatenbank...")

	// Temporarily override config to ensure MariaDB connection
	originalType := config.Cfg.DatabaseType
	originalURL := config.Cfg.DatabaseURL
	originalEnv := config.Cfg.Environment

	// Force MariaDB connection parameters
	config.Cfg.DatabaseType = "mysql"
	config.Cfg.DatabaseURL = "bamort:bG4)efozrc@tcp(mariadb:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"
	config.Cfg.Environment = "production" // Ensure we don't get test DB

	targetDB := database.ConnectDatabaseOrig() // Use original connection method to avoid test DB

	// Restore original config
	config.Cfg.DatabaseType = originalType
	config.Cfg.DatabaseURL = originalURL
	config.Cfg.Environment = originalEnv

	if targetDB == nil {
		logger.Error("Fehler beim Verbinden mit MariaDB-Zieldatenbank")
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to MariaDB target")
		return
	}
	logger.Debug("MariaDB-Verbindung erfolgreich")

	// Migrate all structures to MariaDB first
	logger.Debug("Migriere Strukturen in MariaDB-Datenbank...")
	if err := migrateAllStructures(targetDB); err != nil {
		logger.Error("Fehler beim Migrieren der Strukturen in MariaDB: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to migrate structures to MariaDB: "+err.Error())
		return
	}
	logger.Debug("Strukturen erfolgreich migriert")

	// Clear existing data in MariaDB (optional - be careful!)
	clearExisting := c.Query("clear")
	if clearExisting == "true" {
		logger.Info("Lösche bestehende Daten in MariaDB...")
		if err := clearMariaDBData(targetDB); err != nil {
			logger.Error("Fehler beim Löschen bestehender Daten: %s", err.Error())
			respondWithError(c, http.StatusInternalServerError, "Failed to clear existing data: "+err.Error())
			return
		}
		logger.Debug("Bestehende Daten gelöscht")
	}

	// Copy data from SQLite to MariaDB
	logger.Info("Kopiere Daten von SQLite zu MariaDB...")
	if err := copySQLiteToMariaDB(sourceDB, targetDB); err != nil {
		logger.Error("Fehler beim Kopieren der Daten von SQLite zu MariaDB: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to copy data from SQLite to MariaDB: "+err.Error())
		return
	}

	// Get statistics about the transferred data
	stats, err := getTestDataStatistics(targetDB)
	if err != nil {
		logger.Error("Fehler beim Abrufen der Datenstatistiken: %s", err.Error())
		respondWithError(c, http.StatusInternalServerError, "Failed to get data statistics: "+err.Error())
		return
	}

	logger.Info("Datenübertragung von SQLite zu MariaDB erfolgreich abgeschlossen")
	c.JSON(http.StatusOK, gin.H{
		"message":     "Data transfer from SQLite to MariaDB completed successfully",
		"source_file": sourceFile,
		"target":      "mariadb:3306/bamort",
		"statistics":  stats,
	})
}

// copySQLiteToMariaDB copies all data from SQLite to MariaDB
func copySQLiteToMariaDB(sqliteDB, mariaDB *gorm.DB) error {
	logger.Debug("Starte Kopiervorgang aller Daten von SQLite zu MariaDB...")

	// Disable foreign key checks temporarily to avoid constraint issues
	logger.Debug("Deaktiviere Foreign Key Checks...")
	if err := mariaDB.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		logger.Warn("Warnung: Konnte Foreign Key Checks nicht deaktivieren: %s", err.Error())
	}

	// Re-enable foreign key checks at the end
	defer func() {
		logger.Debug("Aktiviere Foreign Key Checks wieder...")
		if err := mariaDB.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
			logger.Warn("Warnung: Konnte Foreign Key Checks nicht reaktivieren: %s", err.Error())
		}
	}()

	// Same table order as copyMariaDBToSQLite but in reverse direction
	tables := []interface{}{
		// Basis-Strukturen (keine Abhängigkeiten)
		&user.User{},

		// Learning Costs System - Basis
		&models.Source{},
		&models.CharacterClass{},
		&models.SkillCategory{},
		&models.SkillDifficulty{},
		&models.SpellSchool{},

		// GSMaster Basis-Daten (müssen vor den abhängigen Learning Cost Tabellen kommen)
		&models.Skill{},
		&models.WeaponSkill{},
		&models.Spell{},
		&models.Equipment{},
		&models.Weapon{},
		&models.Container{},
		&models.Transportation{},
		&models.Believe{},

		// Learning Costs System - Abhängige Tabellen (nach Skills/Spells)
		&models.ClassCategoryEPCost{},
		&models.ClassSpellSchoolEPCost{},
		&models.SpellLevelLECost{},
		&models.SkillCategoryDifficulty{}, // Jetzt nach Skills
		&models.WeaponSkillCategoryDifficulty{},
		&models.SkillImprovementCost2{},
		&models.ClassCategoryLearningPoints{},
		&models.ClassSpellPoints{},
		&models.ClassTypicalSkill{},
		&models.ClassTypicalSpell{},

		// Charaktere (Basis)
		&models.Char{},

		// Charakter-Eigenschaften (abhängig von Char)
		&models.Eigenschaft{},
		&models.Lp{},
		&models.Ap{},
		&models.B{},
		&models.Merkmale{},
		&models.Erfahrungsschatz{},
		&models.Bennies{},
		&models.Vermoegen{},

		// Charakter-Skills (abhängig von Char und Skills)
		&models.SkFertigkeit{},
		&models.SkWaffenfertigkeit{},
		&models.SkAngeboreneFertigkeit{},
		&models.SkZauber{},

		// Charakter-Equipment (abhängig von Char und Equipment)
		&models.EqAusruestung{},
		&models.EqWaffe{},
		&models.EqContainer{},

		// Character Creation Sessions (abhängig von Char)
		&models.CharacterCreationSession{},

		// Audit Logging (abhängig von Char)
		&models.AuditLogEntry{},
	}

	logger.Info("Kopiere Daten für %d Tabellen von SQLite zu MariaDB...", len(tables))
	for i, model := range tables {
		logger.Debug("Kopiere Tabelle %d/%d: %T", i+1, len(tables), model)
		if err := copyTableDataReverse(sqliteDB, mariaDB, model); err != nil {
			logger.Error("Fehler beim Kopieren der Tabellendaten für %T: %s", model, err.Error())
			return fmt.Errorf("failed to copy table data for %T: %w", model, err)
		}
	}

	logger.Info("Alle Tabellendaten erfolgreich von SQLite zu MariaDB kopiert")
	return nil
}

// copyTableDataReverse copies all data from source to target database
func copyTableDataReverse(sourceDB, targetDB *gorm.DB, model interface{}) error {
	tableName := fmt.Sprintf("%T", model)
	logger.Debug("Starte Kopiervorgang für Tabelle: %s", tableName)

	// Count records in source
	var count int64
	err := sourceDB.Model(model).Count(&count).Error
	if err != nil {
		if isTableNotExistError(err) {
			logger.Debug("Tabelle %s existiert nicht in der Quelle, überspringe", tableName)
			return nil
		}
		logger.Error("Fehler beim Zählen der Datensätze für %s: %s", tableName, err.Error())
		return err
	}

	if count == 0 {
		logger.Debug("Tabelle %s ist leer, keine Daten zu kopieren", tableName)
		return nil
	}

	logger.Debug("Kopiere %d Datensätze für Tabelle %s", count, tableName)

	// Copy data in batches
	batchSize := 100
	totalBatches := (int(count) + batchSize - 1) / batchSize

	for batch := 0; batch < totalBatches; batch++ {
		offset := batch * batchSize
		logger.Debug("Verarbeite Batch %d/%d für Tabelle %s (Offset: %d)", batch+1, totalBatches, tableName, offset)

		// Create slice to hold batch data and read from source
		var records interface{}

		// Read batch from source
		switch model.(type) {
		case *user.User:
			var batch []user.User
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Source:
			var batch []models.Source
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.CharacterClass:
			var batch []models.CharacterClass
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SkillCategory:
			var batch []models.SkillCategory
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SkillDifficulty:
			var batch []models.SkillDifficulty
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SpellSchool:
			var batch []models.SpellSchool
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.ClassCategoryEPCost:
			var batch []models.ClassCategoryEPCost
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.ClassSpellSchoolEPCost:
			var batch []models.ClassSpellSchoolEPCost
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SpellLevelLECost:
			var batch []models.SpellLevelLECost
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SkillCategoryDifficulty:
			var batch []models.SkillCategoryDifficulty
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.WeaponSkillCategoryDifficulty:
			var batch []models.WeaponSkillCategoryDifficulty
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SkillImprovementCost2:
			var batch []models.SkillImprovementCost2
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.ClassCategoryLearningPoints:
			var batch []models.ClassCategoryLearningPoints
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.ClassSpellPoints:
			var batch []models.ClassSpellPoints
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.ClassTypicalSkill:
			var batch []models.ClassTypicalSkill
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.ClassTypicalSpell:
			var batch []models.ClassTypicalSpell
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Skill:
			var batch []models.Skill
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.WeaponSkill:
			var batch []models.WeaponSkill
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Spell:
			var batch []models.Spell
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Equipment:
			var batch []models.Equipment
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Weapon:
			var batch []models.Weapon
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Container:
			var batch []models.Container
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Transportation:
			var batch []models.Transportation
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Believe:
			var batch []models.Believe
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Char:
			var batch []models.Char
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Eigenschaft:
			var batch []models.Eigenschaft
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Lp:
			var batch []models.Lp
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Ap:
			var batch []models.Ap
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.B:
			var batch []models.B
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Merkmale:
			var batch []models.Merkmale
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Erfahrungsschatz:
			var batch []models.Erfahrungsschatz
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Bennies:
			var batch []models.Bennies
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.Vermoegen:
			var batch []models.Vermoegen
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SkFertigkeit:
			var batch []models.SkFertigkeit
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SkWaffenfertigkeit:
			var batch []models.SkWaffenfertigkeit
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SkAngeboreneFertigkeit:
			var batch []models.SkAngeboreneFertigkeit
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.SkZauber:
			var batch []models.SkZauber
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.EqAusruestung:
			var batch []models.EqAusruestung
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.EqWaffe:
			var batch []models.EqWaffe
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.EqContainer:
			var batch []models.EqContainer
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.CharacterCreationSession:
			var batch []models.CharacterCreationSession
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		case *models.AuditLogEntry:
			var batch []models.AuditLogEntry
			if err := sourceDB.Limit(batchSize).Offset(offset).Find(&batch).Error; err != nil {
				return fmt.Errorf("failed to read batch from source: %w", err)
			}
			records = batch
		default:
			return fmt.Errorf("unsupported model type: %T", model)
		}

		// Insert batch into target database using CreateInBatches for better performance
		if err := targetDB.Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches(records, batchSize).Error; err != nil {
			logger.Error("Fehler beim Einfügen des Batches für Tabelle %s: %s", tableName, err.Error())
			return fmt.Errorf("failed to insert batch for table %s: %w", tableName, err)
		}

		logger.Debug("Batch %d/%d für Tabelle %s erfolgreich verarbeitet", batch+1, totalBatches, tableName)
	}

	logger.Debug("Kopiervorgang für Tabelle %s abgeschlossen", tableName)
	return nil
}

// clearMariaDBData clears all data from MariaDB tables (use with caution!)
func clearMariaDBData(db *gorm.DB) error {
	logger.Debug("Lösche alle Daten aus MariaDB-Tabellen...")

	// Clear tables in reverse order due to foreign key constraints
	// (reverse of the insertion order in copySQLiteToMariaDB)
	tables := []interface{}{
		// Audit Logging und Character Creation Sessions (abhängig von Char) - zuerst löschen
		&models.AuditLogEntry{},
		&models.CharacterCreationSession{},

		// Charakter-Equipment (abhängig von Char und Equipment)
		&models.EqContainer{},
		&models.EqWaffe{},
		&models.EqAusruestung{},

		// Charakter-Skills (abhängig von Char und Skills)
		&models.SkZauber{},
		&models.SkAngeboreneFertigkeit{},
		&models.SkWaffenfertigkeit{},
		&models.SkFertigkeit{},

		// Charakter-Eigenschaften (abhängig von Char)
		&models.Vermoegen{},
		&models.Bennies{},
		&models.Erfahrungsschatz{},
		&models.Merkmale{},
		&models.B{},
		&models.Ap{},
		&models.Lp{},
		&models.Eigenschaft{},

		// Charaktere (Basis)
		&models.Char{},

		// Learning Costs System - Abhängige Tabellen (vor Skills/Spells löschen)
		&models.SkillImprovementCost2{},
		&models.WeaponSkillCategoryDifficulty{},
		&models.SkillCategoryDifficulty{},
		&models.SpellLevelLECost{},
		&models.ClassSpellSchoolEPCost{},
		&models.ClassCategoryEPCost{},
		&models.ClassTypicalSpell{},
		&models.ClassTypicalSkill{},
		&models.ClassSpellPoints{},
		&models.ClassCategoryLearningPoints{},

		// GSMaster Basis-Daten
		&models.Believe{},
		&models.Transportation{},
		&models.Container{},
		&models.Weapon{},
		&models.Equipment{},
		&models.Spell{},
		&models.WeaponSkill{},
		&models.Skill{},

		// Learning Costs System - Basis
		&models.SpellSchool{},
		&models.SkillDifficulty{},
		&models.SkillCategory{},
		&models.CharacterClass{},
		&models.Source{},

		// Basis-Strukturen (keine Abhängigkeiten) - zuletzt löschen
		&user.User{},
	}

	for _, model := range tables {
		tableName := fmt.Sprintf("%T", model)
		logger.Debug("Lösche Daten aus Tabelle: %s", tableName)

		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model).Error; err != nil {
			// Continue with other tables even if one fails
			logger.Warn("Warnung beim Löschen der Tabelle %s: %s", tableName, err.Error())
		}
	}

	logger.Debug("Alle Tabellendaten gelöscht")
	return nil
}
