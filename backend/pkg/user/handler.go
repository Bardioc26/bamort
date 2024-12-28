package user

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/Bardioc26/bamort/pkg/auth"
	"github.com/Bardioc26/bamort/pkg/database"
	"github.com/Bardioc26/bamort/pkg/models"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid body"})
	}

	// Passwort hashen
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "password hash error"})
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
	}

	// in DB ablegen
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "user created"})
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Versuchen, über E-Mail oder Username zu finden
	var user models.User
	if err := database.DB.Where("email = ? OR username = ?", req.EmailOrUsername, req.EmailOrUsername).
		First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}

	// Passwort checken
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "wrong password"})
	}

	// JWT generieren
	token, err := auth.GenerateToken(user.UserID, user.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token generation error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "login successful",
		"token":   token,
	})
}

// Passwort zurücksetzen (nur Skizze)
func ResetPassword(c echo.Context) error {
	// E-Mail validieren, neuen Token generieren, oder Link mailen etc.
	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset link sent."})
}
