package database

import (
	"fmt"
	"log"

	"github.com/Bardioc26/bamort/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB ist eine globale Variable für den Datenbankzugriff (kann man aber auch in einer struct verpacken)
var DB *gorm.DB

func ConnectDB(user, password, host string, port int, dbname string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Globale DB-Variable setzen
	DB = db
	log.Println("Connected to MariaDB successfully!")

	return nil
}

func Migrate() error {
	// Reihenfolge abhängig von Foreign Keys
	err := DB.AutoMigrate(
		&models.User{},
		&models.Character{},
		&models.Fertigkeit{},
		&models.Zauber{},
		&models.Waffenfertigkeit{},
		&models.Waffe{},
		&models.Merkmale{},
		&models.Lp{},
		// usw. alle deine Models
	)
	if err != nil {
		return err
	}

	log.Println("DB Migration completed.")
	return nil
}
