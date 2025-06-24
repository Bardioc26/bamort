package database

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var migrationDone bool
var isTestDb bool

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(opts ...bool) {
	isTestDb = true
	if len(opts) > 0 {
		isTestDb = opts[0]
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
				panic("failed to connect to the test database")
			}
			//*/
			migrationDone = true
		}
		DB = db
	}
	/*
		if !migrationDone {
			err := MigrateStructure()
			if err != nil {
				panic("failed to MigrateStructure")
			}
			migrationDone = true
		}
	*/
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
