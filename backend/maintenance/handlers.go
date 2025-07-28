package maintenance

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
	"bamort/importer"
	"bamort/models"
	"bamort/skills"
	"bamort/user"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Constants for test data management
var (
	testDataDir    = database.TestDataDir
	preparedTestDB = database.PreparedTestDB
)

// init function to register the test data loader and migration callback
// This callback mechanism is necessary to avoid circular imports between database and maintenance packages
func init() {
	database.SetTestDataLoader(func(targetDB *gorm.DB) error {
		return LoadPredefinedTestDataFromFile(targetDB, preparedTestDB)
	})
	database.SetMigrationCallback(migrateAllStructures)
}

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
	if err := character.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate character structures: %w", err)
	}
	if err := gsmaster.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate gsmaster structures: %w", err)
	}
	if err := equipment.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate equipment structures: %w", err)
	}
	if err := skills.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate skills structures: %w", err)
	}
	if err := importer.MigrateStructure(db); err != nil {
		return fmt.Errorf("failed to migrate importer structures: %w", err)
	}
	return nil
}

func MakeTestdataFromLive(c *gin.Context) {
	db := database.ConnectDatabase()
	if db == nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to DataBase")
		return
	}
	// Setup test database
	var testDb *gorm.DB
	var err error

	testDb, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to connect to test database testDb: "+err.Error())
		return
	}

	// Step 1: Migrate all structures to test database
	err = migrateAllStructures(testDb)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to migrate structures to test DB: "+err.Error())
		return
	}

	// Step 2: Copy all data from live database to test database
	copyStats, err := copyAllDataToTestDB(db, testDb)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to copy data to test DB: "+err.Error())
		return
	}

	// Step 3: Save test database to file for reuse
	err = saveTestDatabaseToFile(testDb, preparedTestDB)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to save test DB to file: "+err.Error())
		return
	}

	// databases are set up and data copied

	// cleanup test database
	sqlDB, err := testDb.DB()
	if err == nil {
		sqlDB.Close()
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Test data creation from live DB completed successfully",
		"statistics":     copyStats,
		"test_data_file": preparedTestDB,
		"file_size_info": "Check file system for actual size",
	})
}

// copyAllDataToTestDB copies all data from the live database to the test database
// TODO can we make it more dynamic? Maybe use reflection to get all models?
func copyAllDataToTestDB(liveDB, testDB *gorm.DB) (map[string]int, error) {
	stats := make(map[string]int)

	// Define all the model types that need to be copied
	// Order matters due to foreign key constraints!

	// Step 1: Copy base/lookup data first (no foreign keys)
	count, err := copyTableDataWithCount(liveDB, testDB, &user.User{})
	if err != nil {
		return stats, err
	}
	stats["users"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &models.Skill{})
	if err != nil {
		return stats, err
	}
	stats["gsmaster_skills"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &models.WeaponSkill{})
	if err != nil {
		return stats, err
	}
	stats["gsmaster_weapon_skills"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &models.Spell{})
	if err != nil {
		return stats, err
	}
	stats["gsmaster_spells"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &models.Equipment{})
	if err != nil {
		return stats, err
	}
	stats["gsmaster_equipment"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &models.Weapon{})
	if err != nil {
		return stats, err
	}
	stats["gsmaster_weapons"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &models.Container{})
	if err != nil {
		return stats, err
	}
	stats["gsmaster_containers"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &models.Transportation{})
	if err != nil {
		return stats, err
	}
	stats["gsmaster_transportation"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &gsmaster.Believe{})
	if err != nil {
		return stats, err
	}
	stats["gsmaster_believes"] = count

	// Step 2: Copy character data (depends on nothing)
	count, err = copyTableDataWithCount(liveDB, testDB, &character.Char{})
	if err != nil {
		return stats, err
	}
	stats["characters"] = count

	// Step 3: Copy character-dependent data
	count, err = copyTableDataWithCount(liveDB, testDB, &character.Eigenschaft{})
	if err != nil {
		return stats, err
	}
	stats["character_eigenschaften"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &character.Lp{})
	if err != nil {
		return stats, err
	}
	stats["character_lp"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &character.Ap{})
	if err != nil {
		return stats, err
	}
	stats["character_ap"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &character.B{})
	if err != nil {
		return stats, err
	}
	stats["character_b"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &character.Merkmale{})
	if err != nil {
		return stats, err
	}
	stats["character_merkmale"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &character.Erfahrungsschatz{})
	if err != nil {
		return stats, err
	}
	stats["character_erfahrungsschatz"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &character.Bennies{})
	if err != nil {
		return stats, err
	}
	stats["character_bennies"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &character.Vermoegen{})
	if err != nil {
		return stats, err
	}
	stats["character_vermoegen"] = count

	// Step 4: Copy skills (depends on characters)
	count, err = copyTableDataWithCount(liveDB, testDB, &skills.Fertigkeit{})
	if err != nil {
		return stats, err
	}
	stats["skills_fertigkeiten"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &skills.Waffenfertigkeit{})
	if err != nil {
		return stats, err
	}
	stats["skills_waffenfertigkeiten"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &skills.Zauber{})
	if err != nil {
		return stats, err
	}
	stats["skills_zauber"] = count

	// Step 5: Copy equipment (depends on characters)
	count, err = copyTableDataWithCount(liveDB, testDB, &equipment.Ausruestung{})
	if err != nil {
		return stats, err
	}
	stats["equipment_ausruestung"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &equipment.Waffe{})
	if err != nil {
		return stats, err
	}
	stats["equipment_waffen"] = count

	count, err = copyTableDataWithCount(liveDB, testDB, &equipment.Container{})
	if err != nil {
		return stats, err
	}
	stats["equipment_containers"] = count

	return stats, nil
}

// copyTableDataWithCount copies all records from one table to another and returns the count
func copyTableDataWithCount(liveDB, testDB *gorm.DB, model interface{}) (int, error) {
	// Get all records from live database using the actual model type
	// This prevents GORM from modifying our data during Create operations
	var count int64
	if err := liveDB.Model(model).Count(&count).Error; err != nil {
		return 0, err
	}

	// If no records, skip (this is normal for some tables)
	if count == 0 {
		return 0, nil
	}

	// Process records in batches to avoid memory issues with large tables
	batchSize := 100
	totalProcessed := 0

	for offset := 0; offset < int(count); offset += batchSize {
		// Get a fresh batch of records for each iteration
		var records []map[string]interface{}
		if err := liveDB.Model(model).Offset(offset).Limit(batchSize).Find(&records).Error; err != nil {
			return totalProcessed, err
		}

		if len(records) == 0 {
			break
		}

		// Insert the batch into test database
		if err := testDB.Model(model).Create(&records).Error; err != nil {
			// If batch insert fails, try individual inserts to identify problematic records
			for _, record := range records {
				if err := testDB.Model(model).Create(&record).Error; err != nil {
					return totalProcessed, err
				}
				totalProcessed++
			}
		} else {
			totalProcessed += len(records)
		}
	}

	return totalProcessed, nil
}

// saveTestDatabaseToFile saves the in-memory test database to a file
func saveTestDatabaseToFile(testDB *gorm.DB, filename string) error {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Handle existing file by creating a backup
	backupFilename := filename + ".backup"
	fileExists := false

	// Check if the target file exists
	if _, err := os.Stat(filename); err == nil {
		fileExists = true
		// Remove any existing backup file first
		if _, err := os.Stat(backupFilename); err == nil {
			if err := os.Remove(backupFilename); err != nil {
				return fmt.Errorf("failed to remove existing backup file %s: %w", backupFilename, err)
			}
		}
		// Rename existing file to backup
		if err := os.Rename(filename, backupFilename); err != nil {
			return fmt.Errorf("failed to create backup of existing file %s: %w", filename, err)
		}
	}

	// For SQLite in-memory databases, we need to backup to a file
	// We'll use SQLite's backup API through a raw SQL command
	backupSQL := fmt.Sprintf("VACUUM INTO '%s'", filename)
	if err := testDB.Exec(backupSQL).Error; err != nil {
		// If VACUUM INTO fails and we created a backup, restore it
		if fileExists {
			if restoreErr := os.Rename(backupFilename, filename); restoreErr != nil {
				return fmt.Errorf("failed to backup database to file: %w (and failed to restore backup: %v)", err, restoreErr)
			}
		}
		return fmt.Errorf("failed to backup database to file: %w", err)
	}

	// VACUUM INTO succeeded, remove the backup file if it exists
	if fileExists {
		if err := os.Remove(backupFilename); err != nil {
			// Log the error but don't fail the operation since the main task succeeded
			fmt.Printf("Warning: failed to remove backup file %s: %v\n", backupFilename, err)
		}
	}

	return nil
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
		"characters":                &character.Char{},
		"gsmaster_skills":           &models.Skill{},
		"gsmaster_spells":           &models.Spell{},
		"gsmaster_equipment":        &models.Equipment{},
		"skills_fertigkeiten":       &skills.Fertigkeit{},
		"skills_waffenfertigkeiten": &skills.Waffenfertigkeit{},
		"skills_zauber":             &skills.Zauber{},
		"equipment_ausruestung":     &equipment.Ausruestung{},
		"equipment_waffen":          &equipment.Waffe{},
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
