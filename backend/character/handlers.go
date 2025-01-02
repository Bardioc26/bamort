package character

import (
	"bamort/database"
	"bamort/skills"

	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

/*
Character Handlers

Add CRUD operations for characters:
*/

func ListCharacters(c *gin.Context) {
	var characters []Char
	if err := database.DB.Find(&characters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}
	c.JSON(http.StatusOK, characters)
}

func CreateCharacter(c *gin.Context) {
	var character Char
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&character).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
		return
	}

	c.JSON(http.StatusCreated, character)
}
func GetCharacter(c *gin.Context) {
	id := c.Param("id")
	var character Char
	err := character.FirstID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character"})
		return
	}
	c.JSON(http.StatusOK, character)
}
func UpdateCharacter(c *gin.Context) {
	var character Char
	/*
		if err := c.ShouldBindJSON(&character.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.DB.Create(&character).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
			return
		}
	*/
	c.JSON(http.StatusCreated, character)
}
func DeleteCharacter(c *gin.Context) {
	id := c.Param("id")
	var character Char
	err := character.FirstID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve character"})
		return
	}
	err = character.Delete()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete character"})
		return
	}
	/*
		if err := c.ShouldBindJSON(&character.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.DB.Create(&character).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
			return
		}
	*/
	c.JSON(http.StatusCreated, character)
}

// Upload files
func UploadFiles(c *gin.Context) {
	// Get files from the request
	file_vtt, err1 := c.FormFile("file_vtt")
	file_csv, err2 := c.FormFile("file_csv")

	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_vtt is required"})
		return
	}
	if !isValidFileType(file_vtt.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File1 must be a .csv or .json file"})
		return
	}
	// Validate file2 if provided
	if file_csv != nil && !isValidFileType(file_csv.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File2 must be a .csv or .json file"})
		return
	}

	// Save File 1
	err := c.SaveUploadedFile(file_vtt, fmt.Sprintf("./uploads/%s", file_vtt.Filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file_vtt"})
		return
	}

	// Save File 2 if provided
	if err2 == nil {
		err := c.SaveUploadedFile(file_csv, fmt.Sprintf("./uploads/%s", file_csv.Filename))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file_csv"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Files uploaded successfully"})

	// Open and parse JSON
	var character Char
	filePath := fmt.Sprintf("./uploads/%s", file_vtt.Filename)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	if err := json.Unmarshal(fileContent, &character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON structure"})
		return
	}

	// Save character data to the database
	if err := SaveCharacterToDB(&character); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save character to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Character imported successfully",
		"character": character,
	})
}

func isValidFileType(filename string) bool {
	allowedExtensions := []string{".csv", ".json"}
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// Add Fertigkeit by putting it directly to the DB
func AddFertigkeit(charID uint, fertigkeit *skills.Fertigkeit) error {
	// Set the foreign key for the new Eigenschaft
	fertigkeit.CharacterID = charID

	// Save the new Eigenschaft to the database
	if err := database.DB.Create(&fertigkeit).Error; err != nil {
		return fmt.Errorf("failed to add Eigenschaft: %w", err)
	}
	return nil
}

// Append the new Fertigkeit to the slice of the characters property
//character.Fertigkeiten = append(character.Fertigkeiten, fertigkeit)
