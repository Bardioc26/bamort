package user

import (
	"bamort/database"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Role constants
const (
	RoleStandardUser = "standard"
	RoleMaintainer   = "maintainer"
	RoleAdmin        = "admin"
)

type User struct {
	UserID             uint       `gorm:"primaryKey" json:"id"`
	Username           string     `gorm:"unique" json:"username"`
	PasswordHash       string     `json:"password"`
	Email              string     `gorm:"unique" json:"email"`
	Role               string     `gorm:"default:standard" json:"role"`
	ResetPwHash        *string    `gorm:"index" json:"-"` // Hash für Password-Reset (wird nicht serialisiert)
	ResetPwHashExpires *time.Time `json:"-"`              // Ablaufzeit für Password-Reset-Hash
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (object *User) Create() error {
	if database.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

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
	if database.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	err := database.DB.First(&object, "username = ?", value).Error
	if err != nil {
		// User found
		return err
	}
	return nil
}

func (object *User) FirstId(value uint) error {
	if database.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	err := database.DB.First(&object, "user_id = ?", value).Error
	if err != nil {
		// User found
		return err
	}
	return nil
}

func (object *User) Save() error {
	if database.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	err := database.DB.Save(&object).Error
	if err != nil {
		// User found
		return err
	}
	return nil
}

// FindByEmail findet einen User anhand der E-Mail-Adresse
func (object *User) FindByEmail(email string) error {
	if database.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	err := database.DB.First(&object, "email = ?", email).Error
	return err
}

// FindByResetHash findet einen User anhand des Reset-Hashes
func (object *User) FindByResetHash(resetHash string) error {
	if database.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

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

// HasRole checks if the user has the specified role
func (u *User) HasRole(role string) bool {
	return u.Role == role
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsMaintainer checks if the user is a maintainer or higher
func (u *User) IsMaintainer() bool {
	return u.Role == RoleMaintainer || u.Role == RoleAdmin
}

// IsStandardUser checks if the user is a standard user or higher
func (u *User) IsStandardUser() bool {
	return u.Role == RoleStandardUser || u.IsMaintainer()
}

// ValidateRole checks if the given role is valid
func ValidateRole(role string) bool {
	return role == RoleStandardUser || role == RoleMaintainer || role == RoleAdmin
}
