/*
User Handlers

Add handlers for user registration and login:
*/
package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully:"})
}

func LoginUser(c *gin.Context) {
	var user User
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username. or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Apply middleware to protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Add token validation logic here

		c.Next()
	}
}

/*
Character Handlers

Add CRUD operations for characters:
*/

func GetCharacters(c *gin.Context) {
	var characters []Character
	if err := DB.Find(&characters).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve characters"})
		return
	}
	c.JSON(http.StatusOK, characters)
}

func CreateCharacter(c *gin.Context) {
	var character Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Create(&character).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create character"})
		return
	}

	c.JSON(http.StatusCreated, character)
}

/*
Endpoints for Managing Ausruestung
1. Create Ausruestung

Allows users to add new equipment items for a specific character.
*/

func CreateAusruestung(c *gin.Context) {
	var ausruestung Ausruestung
	if err := c.ShouldBindJSON(&ausruestung); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Create(&ausruestung).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Ausruestung"})
		return
	}

	c.JSON(http.StatusCreated, ausruestung)
}

func GetAusruestung(c *gin.Context) {
	characterID := c.Param("character_id")

	var ausruestung []Ausruestung
	if err := DB.Where("character_id = ?", characterID).Find(&ausruestung).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Ausruestung"})
		return
	}

	c.JSON(http.StatusOK, ausruestung)
}

func UpdateAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")
	var ausruestung Ausruestung

	if err := DB.First(&ausruestung, ausruestungID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ausruestung not found"})
		return
	}

	if err := c.ShouldBindJSON(&ausruestung); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Save(&ausruestung).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Ausruestung"})
		return
	}

	c.JSON(http.StatusOK, ausruestung)
}

func DeleteAusruestung(c *gin.Context) {
	ausruestungID := c.Param("ausruestung_id")
	if err := DB.Delete(&Ausruestung{}, ausruestungID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Ausruestung"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ausruestung deleted successfully"})
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
