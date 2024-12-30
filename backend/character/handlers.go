package character

import (
	"bamort/database"

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

func GetCharacters(c *gin.Context) {
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
