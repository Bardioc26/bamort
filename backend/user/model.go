package user

import (
	"bamort/database"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID       uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"username"`
	PasswordHash string    `json:"password"`
	Email        string    `gorm:"unique" json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (object *User) Create() error {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save the User record
		if err := tx.Create(&object).Error; err != nil {
			return fmt.Errorf("failed to save User: %w", err)
		}
		return nil
	})

	return err
}

func (object *User) First(value string) error {
	err := database.DB.First(&object, "username = ?", value).Error
	if err != nil {
		// User found
		return err
	}
	return nil
}

func (object *User) FirstId(value uint) error {
	err := database.DB.First(&object, "user_id = ?", value).Error
	if err != nil {
		// User found
		return err
	}
	return nil
}

func (object *User) Save() error {
	err := database.DB.Save(&object).Error
	if err != nil {
		// User found
		return err
	}
	return nil
}
