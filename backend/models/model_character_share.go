package models

import (
	"bamort/database"
	"fmt"
)

type CharShare struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	CharacterID uint   `gorm:"index" json:"character_id"` // ID of the character being shared
	UserID      uint   `gorm:"index" json:"user_id"`      // ID of the user with whom the character is shared
	Permission  string `json:"permission"`                // Permission level (e.g., "read", "write")
}

func (object *CharShare) TableName() string {
	dbPrefix := "char"
	return dbPrefix + "_" + "shares"
}

func (object *CharShare) FirstByChar(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid character ID")
	}
	return database.DB.First(object, "character_id = ?", id).Error
}

func (object *CharShare) FirstByUser(id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid user ID")
	}
	return database.DB.First(object, "user_id = ?", id).Error
}
