package models

import "time"

type User struct {
	UserID       uint      `gorm:"primaryKey"`
	Username     string    `gorm:"unique" json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `gorm:"unique" json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
