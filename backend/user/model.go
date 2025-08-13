package user

import (
	"bamort/database"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID             uint       `gorm:"primaryKey" json:"id"`
	Username           string     `gorm:"unique" json:"username"`
	PasswordHash       string     `json:"password"`
	Email              string     `gorm:"unique" json:"email"`
	ResetPwHash        *string    `gorm:"index" json:"-"` // Hash für Password-Reset (wird nicht serialisiert)
	ResetPwHashExpires *time.Time `json:"-"`              // Ablaufzeit für Password-Reset-Hash
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
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

// FindByEmail findet einen User anhand der E-Mail-Adresse
func (object *User) FindByEmail(email string) error {
	err := database.DB.First(&object, "email = ?", email).Error
	return err
}

// FindByResetHash findet einen User anhand des Reset-Hashes
func (object *User) FindByResetHash(resetHash string) error {
	err := database.DB.First(&object, "reset_pw_hash = ? AND reset_pw_hash_expires > ?", resetHash, time.Now()).Error
	return err
}

// SetPasswordResetHash setzt den Reset-Hash und die Ablaufzeit (14 Tage)
func (object *User) SetPasswordResetHash(resetHash string) error {
	expiryTime := time.Now().Add(14 * 24 * time.Hour) // 14 Tage gültig
	object.ResetPwHash = &resetHash
	object.ResetPwHashExpires = &expiryTime
	return object.Save()
}

// ClearPasswordResetHash entfernt den Reset-Hash
func (object *User) ClearPasswordResetHash() error {
	object.ResetPwHash = nil
	object.ResetPwHashExpires = nil
	return object.Save()
}

// IsResetHashValid prüft ob der Reset-Hash gültig und nicht abgelaufen ist
func (object *User) IsResetHashValid(resetHash string) bool {
	if object.ResetPwHash == nil || object.ResetPwHashExpires == nil {
		return false
	}
	return *object.ResetPwHash == resetHash && time.Now().Before(*object.ResetPwHashExpires)
}
