package database

import (
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
			//* //testin with persitant MariaDB
			dsn := "bamort:bG4)efozrc@tcp(192.168.0.5:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"
			//dsn := "root:26Osiris-Mar@tcp(192.168.0.5:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"
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
