package database

import (
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

// The Database for testing is created from the live database whenever needed and stored in the path defined in database.PreparedTestDB
// to use the test database make a temporary copy of it and then open this new copy as testing database
// This allows to have a clean database for each test run without affecting the live database
// However SetupTestDB can still open the live database if required by setting isTestDb to false
// SetupTestDB creates an in-memory SQLite database for testing
// Parameters:
// - opts[0]: isTestDb (bool) - whether to use precopied SQLite (true) or persistent (Live) MariaDB (false)
func SetupTestDB(opts ...bool) {
	isTestDb = true

	if len(opts) > 0 {
		isTestDb = opts[0]
	}

	if DB == nil {
		var db *gorm.DB
		if isTestDb {
			testdbTempDir, err := os.MkdirTemp("", "bamort-test-")
			if err != nil {
				panic("failed to create temporary directory: " + err.Error())
			}
			targetFile := filepath.Join(testdbTempDir, "test_backup.db")
			err = copyFile(PreparedTestDB, targetFile)
			if err != nil {
				panic("failed to copy prepared test database: " + err.Error())
			}
			db, err = gorm.Open(sqlite.Open(targetFile), &gorm.Config{})
			if err != nil {
				panic("failed to connect to the test database: " + err.Error())
			}
			//defer os.RemoveAll(testdbTempDir)
		} else {
			//* //testing with persistent MariaDB
			db = ConnectDatabase()
			if db == nil {
				panic("failed to connect to the live database")
			}
		}
		DB = db
	}

}
func ResetTestDB() {
	if isTestDb {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			DB = nil
			os.RemoveAll(testdbTempDir)
		}
	}
}
