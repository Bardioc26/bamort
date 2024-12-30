package tests

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() *gorm.DB {
	//*
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the test database")
	}
	//*/
	/* //testin with persitant MariaDB
	dsn := "bamort:bG4)efozrc@tcp(192.168.0.5:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:26Osiris-Mar@tcp(192.168.0.5:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the test database")
	}
	//*/

	return db
}
