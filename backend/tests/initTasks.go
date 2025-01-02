package tests

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
	"bamort/importer"
	"bamort/skills"
	"bamort/user"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var migrationDone bool

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

func MigrateStructure() error {
	err := database.MigrateStructure()
	if err != nil {
		return err
	}
	err = character.MigrateStructure()
	if err != nil {
		return err
	}
	err = equipment.MigrateStructure()
	if err != nil {
		return err
	}
	err = gsmaster.MigrateStructure()
	if err != nil {
		return err
	}
	err = importer.MigrateStructure()
	if err != nil {
		return err
	}
	err = skills.MigrateStructure()
	if err != nil {
		return err
	}
	err = user.MigrateStructure()
	if err != nil {
		return err
	}
	migrationDone = true

	return nil
}
