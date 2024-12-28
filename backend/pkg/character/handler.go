package character

import (
	"net/http"
	"strconv"

	"github.com/Bardioc26/bamort/pkg/database"
	"github.com/Bardioc26/bamort/pkg/models"
	"github.com/labstack/echo/v4"
)

// ListCharacters liefert eine Liste der Charaktere des eingeloggten Benutzers (als Beispiel).
func ListCharacters(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	// Angenommen es gibt in Character noch eine user_id-FK (dann müsstest du dein Modell anpassen)
	var characters []models.Character
	if err := database.DB.Where("user_id = ?", userID).Find(&characters).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, characters)
}

// GetCharacterByID ...
func GetCharacterByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var character models.Character
	if err := database.DB.Preload("Fertigkeiten").Preload("Zauber").
		First(&character, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "character not found"})
	}
	return c.JSON(http.StatusOK, character)
}

// CreateCharacter ...
func CreateCharacter(c echo.Context) error {
	var ch models.Character
	if err := c.Bind(&ch); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid data"})
	}
	// userID aus Token
	// ch.UserID = c.Get("user_id").(uint) // falls du die Spalte hast
	if err := database.DB.Create(&ch).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, ch)
}

// UpdateCharacter ...
func UpdateCharacter(c echo.Context) error {
	// ...
	return c.JSON(http.StatusOK, map[string]string{"message": "update not implemented yet"})
}

// DeleteCharacter ...
func DeleteCharacter(c echo.Context) error {
	// ...
	return c.JSON(http.StatusOK, map[string]string{"message": "delete not implemented yet"})
}
