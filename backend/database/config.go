package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	dsn := "bamort:bG4)efozrc@tcp(192.168.0.5:3306)/bamort?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	DB = database
	return DB
}

func GetDB() *gorm.DB {
	if DB == nil {
		DB = ConnectDatabase()
	}
	return DB
}

// StringArray is a custom type for []string
type StringArray []string

// Value implements the driver.Valuer interface for database storage
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s) // Serialize []string to JSON
}

// Scan implements the sql.Scanner interface for database retrieval
func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to convert database value to []byte")
	}

	return json.Unmarshal(bytes, s) // Deserialize JSON to []string
}
