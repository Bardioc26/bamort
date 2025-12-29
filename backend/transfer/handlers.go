package transfer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExportCharacterHandler handles character export requests
func ExportCharacterHandler(c *gin.Context) {
	// Get character ID from URL parameter
	charIDStr := c.Param("id")
	charID, err := strconv.ParseUint(charIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	// Export character
	exportData, err := ExportCharacter(uint(charID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to export character: %v", err)})
		return
	}

	// Return as JSON
	c.JSON(http.StatusOK, exportData)
}

// DownloadCharacterHandler exports character as downloadable JSON file
func DownloadCharacterHandler(c *gin.Context) {
	// Get character ID from URL parameter
	charIDStr := c.Param("id")
	charID, err := strconv.ParseUint(charIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	// Export character
	exportData, err := ExportCharacter(uint(charID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to export character: %v", err)})
		return
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("character_%s_export.json", exportData.Character.Name)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", jsonData)
}

// ImportCharacterHandler handles character import requests
func ImportCharacterHandler(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse import data from request body
	var importData CharacterExport
	if err := c.ShouldBindJSON(&importData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}

	// Import character
	charID, err := ImportCharacter(&importData, userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to import character: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Character imported successfully",
		"character_id": charID,
	})
}
