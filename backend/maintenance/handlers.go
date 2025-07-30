package maintenance

import (
	"bamort/database"
	"bamort/models"
	"bamort/user"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	// Migrate all structures in the correct order
	if err := database.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate database structures: %w", err)
	}
	if err := user.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate user structures: %w", err)
	}
	if err := models.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate gsmaster structures: %w", err)
	}

	/*if err := importer.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate importer structures: %w", err)
	}*/
	return nil
}

func MakeTestdataFromLive(c *gin.Context) {
	liveDB := database.ConnectDatabase()
	if liveDB == nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to live database")
		return
	}

	// Live-Datenbank in SQLite-Datei kopieren
	backupFile := preparedTestDB
	err := copyLiveDatabaseToFile(liveDB, backupFile)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to copy database: %v", err))
		return
	}

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
	// Verzeichnis erstellen falls es nicht existiert
	dir := filepath.Dir(targetFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Backup der existierenden Datei erstellen
	if _, err := os.Stat(targetFile); err == nil {
		backupFile := targetFile + ".backup"
		os.Remove(backupFile) // Alte Backup entfernen
		if err := os.Rename(targetFile, backupFile); err != nil {
			return fmt.Errorf("failed to backup existing file: %w", err)
		}
	}

	// SQLite-Zieldatenbank erstellen
	targetDB, err := gorm.Open(sqlite.Open(targetFile), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to create target SQLite database: %w", err)
	}
	defer func() {
		if sqlDB, err := targetDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Strukturen in SQLite-DB migrieren
	if err := migrateAllStructures(targetDB); err != nil {
		return fmt.Errorf("failed to migrate structures to SQLite: %w", err)
	}

	// Daten von MariaDB zu SQLite kopieren
	if err := copyMariaDBToSQLite(liveDB, targetDB); err != nil {
		return fmt.Errorf("failed to copy data from MariaDB to SQLite: %w", err)
	}

	return nil
}

// copyMariaDBToSQLite kopiert alle Daten von MariaDB zu SQLite
func copyMariaDBToSQLite(mariaDB, sqliteDB *gorm.DB) error {
	// Vollständige Liste aller Strukturen mit GORM-Tags in der richtigen Reihenfolge
	// (Basis-Tabellen zuerst wegen Foreign Key-Abhängigkeiten)
	tables := []interface{}{
		// Basis-Strukturen (keine Abhängigkeiten)
		&user.User{},

		// Learning Costs System - Basis
		&models.Source{},
		&models.CharacterClass{},
		&models.SkillCategory{},
		&models.SkillDifficulty{},
		&models.SpellSchool{},

		// Learning Costs System - Abhängige Tabellen
		&models.ClassCategoryEPCost{},
		&models.ClassSpellSchoolEPCost{},
		&models.SpellLevelLECost{},
		&models.SkillCategoryDifficulty{},
		&models.SkillImprovementCost{},

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

		// View-Strukturen ohne eigene Tabellen werden nicht kopiert:
		// SkillLearningInfo, SpellLearningInfo, CharList, FeChar, etc.
	}

	for _, model := range tables {
		if err := copyTableData(mariaDB, sqliteDB, model); err != nil {
			return fmt.Errorf("failed to copy table data for %T: %w", model, err)
		}
	}

	return nil
}

// copyTableData kopiert alle Daten einer Tabelle von MariaDB zu SQLite
func copyTableData(sourceDB, targetDB *gorm.DB, model interface{}) error {
	// Anzahl der Datensätze prüfen
	var count int64
	err := sourceDB.Model(model).Count(&count).Error
	if err != nil {
		// If table doesn't exist, skip silently (useful for testing with partial schemas)
		if isTableNotExistError(err) {
			return nil
		}
		return err
	}

	if count == 0 {
		return nil // Keine Daten zu kopieren
	}

	// Daten in Blöcken kopieren (für große Tabellen)
	batchSize := 100
	for offset := 0; offset < int(count); offset += batchSize {
		var records []map[string]interface{}

		// Batch aus MariaDB lesen
		if err := sourceDB.Model(model).Offset(offset).Limit(batchSize).Find(&records).Error; err != nil {
			return err
		}

		if len(records) == 0 {
			break
		}

		// Batch in SQLite einfügen mit Konflikt-Behandlung
		// Verwende Clauses.OnConflict um bestehende Datensätze zu ersetzen
		if err := targetDB.Model(model).Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&records).Error; err != nil {
			return err
		}
	}

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
	// Check if file exists
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return fmt.Errorf("predefined test data file not found: %s", dataFile)
	}

	// Migrate structures to target DB
	err := migrateAllStructures(targetDB)
	if err != nil {
		return fmt.Errorf("failed to migrate structures: %w", err)
	}

	// Copy data from file database to target database
	err = copyDataFromFileToMemory(dataFile, targetDB)
	if err != nil {
		return fmt.Errorf("failed to copy test data to database: %w", err)
	}

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
	// Copy all tables using ATTACH and INSERT
	attachSQL := fmt.Sprintf("ATTACH DATABASE '%s' AS source", sourceFile)
	if err := targetDB.Exec(attachSQL).Error; err != nil {
		return fmt.Errorf("failed to attach source database: %w", err)
	}

	// Get list of tables from source database
	var tables []string
	if err := targetDB.Raw("SELECT name FROM source.sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get table list: %w", err)
	}

	// Copy each table
	for _, table := range tables {
		copySQL := fmt.Sprintf("INSERT OR REPLACE INTO main.%s SELECT * FROM source.%s", table, table)
		if err := targetDB.Exec(copySQL).Error; err != nil {
			return fmt.Errorf("failed to copy table %s: %w", table, err)
		}
	}

	// Detach the source database
	if err := targetDB.Exec("DETACH DATABASE source").Error; err != nil {
		return fmt.Errorf("failed to detach source database: %w", err)
	}

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

func SetupCheck(c *gin.Context) {
	db := database.ConnectDatabase()
	if db == nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to DataBase")
		return
	}
	err := migrateAllStructures(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Setup Check OK"})
}

/*
// InitializeLearningCosts initialisiert das Lernkosten-System
// Wird danach nicht mehr benötigt
func InitializeLearningCosts(c *gin.Context) {
	err := gsmaster.InitializeLearningCostsSystem()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to initialize learning costs: %v", err))
		return
	}

	// Validierung
	if err := gsmaster.ValidateLearningCostsData(); err != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Learning costs initialized but validation failed: %v", err))
		return
	}

	// Zusammenfassung
	summary, err := gsmaster.GetLearningCostsSummary()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to get summary: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Learning costs system initialized successfully",
		"summary": summary,
	})
}
*/
