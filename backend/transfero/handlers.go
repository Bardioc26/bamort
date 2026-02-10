package transfero

import (
	"bamort/config"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func dummyproc(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "This endpoint is a placeholder for future VTT import functionality."})
}

// ExportCharacterHandler handles character export requests
// ExportCharacterHandler godoc
// @Summary Export character as JSON
// @Description Exports a character as JSON for API consumption
// @Tags Import/Export
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {object} object "Character data in JSON format"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "Character not found"
// @Security BearerAuth
// @Router /api/transfer/export/{id} [get]
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
// DownloadCharacterHandler godoc
// @Summary Download character as JSON file
// @Description Downloads a character as a JSON file attachment
// @Tags Import/Export
// @Produce application/json
// @Param id path int true "Character ID"
// @Success 200 {file} binary "JSON file download"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Access denied"
// @Failure 404 {object} map[string]string "Character not found"
// @Security BearerAuth
// @Router /api/transfer/download/{id} [get]
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
// ImportCharacterHandler godoc
// @Summary Import character from JSON
// @Description Imports a character from JSON data
// @Tags Import/Export
// @Accept json
// @Produce json
// @Param character body object true "Character JSON data"
// @Success 201 {object} object "Imported character"
// @Failure 400 {object} map[string]string "Invalid JSON data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /api/transfer/import [post]
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

// ExportDatabaseHandler handles full database export requests
// ExportDatabaseHandler godoc
// @Summary Export full database
// @Description Exports the entire database as JSON (admin only)
// @Tags Import/Export
// @Produce json
// @Success 200 {object} object "Database export data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - admin access required"
// @Security BearerAuth
// @Router /api/transfer/database/export [post]
func ExportDatabaseHandler(c *gin.Context) {
	// Use configured export directory
	exportDir := config.Cfg.ExportTempDir
	if exportDir == "" {
		exportDir = "./backend/export_temp"
	}

	// Export database
	result, err := ExportDatabase(exportDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to export database: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Database exported successfully",
		"filename":     result.Filename,
		"filepath":     result.FilePath,
		"record_count": result.RecordCount,
		"timestamp":    result.Timestamp,
	})
}

// ImportDatabaseHandler handles full database import requests
// ImportDatabaseHandler godoc
// @Summary Import full database
// @Description Imports an entire database from JSON (admin only)
// @Tags Import/Export
// @Accept json
// @Produce json
// @Param database body object true "Database JSON data"
// @Success 200 {object} object "Import result"
// @Failure 400 {object} map[string]string "Invalid JSON data"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden - admin access required"
// @Security BearerAuth
// @Router /api/transfer/database/import [post]
func ImportDatabaseHandler(c *gin.Context) {
	// Parse request body with filepath
	var req struct {
		FilePath string `json:"filepath" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	// Import database
	result, err := ImportDatabase(req.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to import database: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Database imported successfully",
		"record_count": result.RecordCount,
		"timestamp":    result.Timestamp,
	})
}
