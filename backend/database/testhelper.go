package database

import (
	"io"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var migrationDone bool
var isTestDb bool

// testDataLoader is a callback function that can be set to load test data
var testDataLoader func(*gorm.DB) error

// migrationCallback is a callback function that can be set to migrate all structures
var migrationCallback func(*gorm.DB) error

// SetTestDataLoader sets the function to load test data
func SetTestDataLoader(loader func(*gorm.DB) error) {
	testDataLoader = loader
}

// SetMigrationCallback sets the function to migrate all structures
func SetMigrationCallback(migrator func(*gorm.DB) error) {
	migrationCallback = migrator
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// SetupTestDB creates an in-memory SQLite database for testing
// Parameters:
// - opts[0]: isTestDb (bool) - whether to use in-memory SQLite (true) or persistent MariaDB (false)
// - opts[1]: loadTestData (bool) - whether to load predefined test data from file
func SetupTestDB(opts ...bool) {
	isTestDb = true
	loadTestData := false

	if len(opts) > 0 {
		isTestDb = opts[0]
	}
	if len(opts) > 1 {
		loadTestData = opts[1]
	}

	if DB == nil {
		var db *gorm.DB
		var err error
		if isTestDb {
			//*
			db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				panic("failed to connect to the test database")
			}
			//*/
		} else {
			//* //testing with persistent MariaDB
			dsn := os.Getenv("TEST_DB_DSN")
			if dsn == "" {
				dsn = "bamort:password@tcp(localhost:3306)/bamort_test?charset=utf8mb4&parseTime=True&loc=Local"
			}
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				//panic("failed to connect to the live database")
				// noch ein versuch mit der lokalen SQLite-Datenbank
				testDataPath := filepath.Join("..", "testdata", "test_data.db")

				// Check if test_data.db exists, if not try to copy from predefined_test_data.db
				if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
					predefinedPath := filepath.Join("..", "testdata", "prepared_test_data.db")
					if _, err := os.Stat(predefinedPath); err == nil {
						// Create directory if it doesn't exist
						if mkdirErr := os.MkdirAll(filepath.Dir(testDataPath), 0755); mkdirErr != nil {
							panic("failed to create testdata directory: " + mkdirErr.Error())
						}
						// Copy predefined_test_data.db to test_data.db
						if copyErr := copyFile(predefinedPath, testDataPath); copyErr != nil {
							panic("failed to copy predefined test data: " + copyErr.Error())
						}
					}
				}

				db, err = gorm.Open(sqlite.Open(testDataPath), &gorm.Config{})
			}
			if err != nil {
				panic("failed to connect to the live database")
			}
			//*/
			migrationDone = true
		}
		DB = db
	}

	// If loadTestData is requested and we have an in-memory database
	if loadTestData && isTestDb && DB != nil {
		// First migrate the structures using callback if available
		if migrationCallback != nil {
			err := migrationCallback(DB)
			if err != nil {
				panic("failed to migrate all structures for test data loading: " + err.Error())
			}
		} else {
			// Fallback to basic migration
			err := MigrateStructure()
			if err != nil {
				panic("failed to MigrateStructure for test data loading: " + err.Error())
			}
		}

		// Load test data using maintenance function
		// We need to import the maintenance package, but to avoid circular imports,
		// we'll create a callback mechanism
		if testDataLoader != nil {
			err := testDataLoader(DB)
			if err != nil {
				panic("failed to load test data: " + err.Error())
			}
		}
		migrationDone = true
	}
}
func ResetTestDB() {
	if isTestDb {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			DB = nil
			migrationDone = false
		}
	}
}
