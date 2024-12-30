package user

import (
	"time"
)

type User struct {
	UserID       uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"unique" json:"name"`
	PasswordHash string    `json:"password"`
	Email        string    `gorm:"unique" json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
