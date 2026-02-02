package importer

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// Helper function for error responses
func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// Upload files
func UploadFiles(c *gin.Context) {
	// Get files from the request
	file_vtt, err1 := c.FormFile("file_vtt")
	file_csv, err2 := c.FormFile("file_csv")
	if err1 != nil {
		respondWithError(c, http.StatusBadRequest, "file_vtt is required")
		return
	}
	if !isValidFileType(file_vtt.Filename) {
		respondWithError(c, http.StatusBadRequest, "File1 must be a .csv or .json file")
		return
	}

	vttFileName := fmt.Sprintf("./uploads/%s", file_vtt.Filename)
	csvFileName := "./uploads/default.csv"
	if file_csv != nil {
		csvFileName = fmt.Sprintf("./uploads/%s", file_csv.Filename)
	}

	// Validate file2 if provided
	if file_csv != nil && !isValidFileType(file_csv.Filename) {
		respondWithError(c, http.StatusBadRequest, "File2 must be a .csv or .json file")
		return
	}

	// Save File 1
	err := c.SaveUploadedFile(file_vtt, vttFileName)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to save file_vtt")
		return
	}

	// Save File 2 if provided
	if err2 == nil {
		err := c.SaveUploadedFile(file_csv, csvFileName)
		if err != nil {
			respondWithError(c, http.StatusInternalServerError, "Failed to save file_csv")
			return
		}
	}
	userID := c.GetUint("userID")

	character, err3 := ImportVTTJSON(vttFileName, userID)
	if err3 != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to Import Character from file %s", err3.Error()))
		return
	}
	if character.ID < 1 {
		respondWithError(c, http.StatusInternalServerError, "Failed to Import Character from file ID is < 1")
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

// ImportSpellCSVHandler handles HTTP requests to import spell data from CSV files
// @Summary Import spells from CSV file
// @Description Imports spell data from a CSV file into the database. The CSV file should contain spell information with headers matching the database fields.
// @Tags importer
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file to import"
// @Success 200 {object} map[string]interface{} "Import successful"
// @Failure 400 {object} map[string]interface{} "Bad request - missing file parameter, file not found, or invalid file type"
// @Failure 500 {object} map[string]interface{} "Internal server error - import failed"
// @Router /api/importer/spells/csv [post]
func ImportSpellCSVHandler(c *gin.Context) {
	// Try to get file from multipart form first
	file, err := c.FormFile("file")
	var filePath string

	if err != nil {
		// Fallback to query parameter for backward compatibility
		filePath = c.Query("file")
		if filePath == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Missing file parameter",
				"message": "Please provide a CSV file via multipart upload or file path using the 'file' parameter",
			})
			return
		}

		// Validate file exists and has proper extension
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "File not found",
				"message": fmt.Sprintf("File '%s' does not exist", filePath),
			})
			return
		}
	} else {
		// Handle uploaded file
		// Check file extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".csv" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid file type",
				"message": "Only CSV files are supported",
			})
			return
		}

		// Create uploads directory if it doesn't exist
		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			if err := os.MkdirAll(uploadDir, 0755); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to create upload directory",
					"message": err.Error(),
				})
				return
			}
		}

		// Save the uploaded file
		filePath = filepath.Join(uploadDir, file.Filename)
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to save uploaded file",
				"message": err.Error(),
			})
			return
		}
	}

	// Check file extension for query parameter path
	if file == nil {
		ext := strings.ToLower(filepath.Ext(filePath))
		if ext != ".csv" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid file type",
				"message": "Only CSV files are supported",
			})
			return
		}
	}

	// Clear source cache before import to ensure fresh data
	ClearSourceCache()

	// Perform the import
	err = ImportCsv2Spell(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Import failed",
			"message": err.Error(),
		})
		return
	}

	// Count imported spells for response
	var spellCount int64
	if countErr := database.DB.Model(&models.Spell{}).Count(&spellCount).Error; countErr != nil {
		// If count fails, just report success without count
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Spells imported successfully",
			"file":    filepath.Base(filePath),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"message":      "Spells imported successfully",
		"file":         filepath.Base(filePath),
		"total_spells": spellCount,
	})
}

// ExportCharacterVTTHandler exports a character to VTT JSON format
// @Summary Export character to VTT JSON format
// @Description Exports a character to VTT JSON format for use in other systems
// @Tags importer
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {object} CharacterImport "Export successful"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid character ID"
// @Failure 404 {object} map[string]interface{} "Character not found"
// @Failure 500 {object} map[string]interface{} "Internal server error - export failed"
// @Router /api/importer/export/vtt/{id} [get]
func ExportCharacterVTTHandler(c *gin.Context) {
	// Get character ID from URL parameter
	charID := c.Param("id")
	if charID == "" {
		respondWithError(c, http.StatusBadRequest, "Character ID is required")
		return
	}

	// Load character from database
	var char models.Char
	err := database.DB.Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Zauber").
		Preload("Waffen").
		Preload("Ausruestung").
		Preload("Behaeltnisse").
		Preload("Transportmittel").
		First(&char, charID).Error

	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Export to VTT format
	vttChar, err := ExportCharToVTT(&char)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to export character: %s", err.Error()))
		return
	}

	// Return as JSON
	c.JSON(http.StatusOK, vttChar)
}

// ExportCharacterVTTFileHandler exports a character to VTT JSON file
// @Summary Export character to VTT JSON file
// @Description Exports a character to VTT JSON file and returns it as a download
// @Tags importer
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {file} file "VTT JSON file"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid character ID"
// @Failure 404 {object} map[string]interface{} "Character not found"
// @Failure 500 {object} map[string]interface{} "Internal server error - export failed"
// @Router /api/importer/export/vtt/{id}/file [get]
func ExportCharacterVTTFileHandler(c *gin.Context) {
	// Get character ID from URL parameter
	charID := c.Param("id")
	if charID == "" {
		respondWithError(c, http.StatusBadRequest, "Character ID is required")
		return
	}

	// Load character from database
	var char models.Char
	err := database.DB.Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Zauber").
		Preload("Waffen").
		Preload("Ausruestung").
		Preload("Behaeltnisse").
		Preload("Transportmittel").
		First(&char, charID).Error

	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Create temp file
	tempFile, err := os.CreateTemp("", fmt.Sprintf("vtt_export_%s_*.json", char.Name))
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create temp file")
		return
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Export to file
	err = ExportCharToVTTFile(&char, tempFile.Name())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to export character: %s", err.Error()))
		return
	}

	// Send file as download
	filename := fmt.Sprintf("%s_vtt_export.json", char.Name)
	c.FileAttachment(tempFile.Name(), filename)
}

// ExportSpellsCSVHandler exports spell master data to CSV file
// @Summary Export spells to CSV file
// @Description Exports spell master data to CSV format
// @Tags importer
// @Produce text/csv
// @Param game_system query string false "Game system filter (e.g., 'midgard')"
// @Success 200 {file} file "CSV file"
// @Failure 500 {object} map[string]interface{} "Internal server error - export failed"
// @Router /api/importer/export/spells/csv [get]
func ExportSpellsCSVHandler(c *gin.Context) {
	gameSystem := c.Query("game_system")

	// Load spells from database
	var spells []models.Spell
	query := database.DB
	if gameSystem != "" {
		query = query.Where("game_system = ?", gameSystem)
	}
	err := query.Find(&spells).Error

	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to load spells")
		return
	}

	// Create temp file
	tempFile, err := os.CreateTemp("", "spells_export_*.csv")
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create temp file")
		return
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Export to CSV
	err = ExportSpellsToCSV(spells, tempFile.Name())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to export spells: %s", err.Error()))
		return
	}

	// Send file as download
	filename := "spells_export.csv"
	if gameSystem != "" {
		filename = fmt.Sprintf("spells_%s_export.csv", gameSystem)
	}
	c.FileAttachment(tempFile.Name(), filename)
}

// ExportCharacterCSVHandler exports a character to CSV file
// @Summary Export character to CSV file
// @Description Exports a character to CSV format (MOAM-compatible)
// @Tags importer
// @Produce text/csv
// @Param id path int true "Character ID"
// @Success 200 {file} file "CSV file"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid character ID"
// @Failure 404 {object} map[string]interface{} "Character not found"
// @Failure 500 {object} map[string]interface{} "Internal server error - export failed"
// @Router /api/importer/export/csv/{id} [get]
func ExportCharacterCSVHandler(c *gin.Context) {
	// Get character ID from URL parameter
	charID := c.Param("id")
	if charID == "" {
		respondWithError(c, http.StatusBadRequest, "Character ID is required")
		return
	}

	// Load character from database
	var char models.Char
	err := database.DB.Preload("Eigenschaften").
		Preload("Fertigkeiten").
		Preload("Waffenfertigkeiten").
		Preload("Zauber").
		Preload("Waffen").
		Preload("Ausruestung").
		Preload("Behaeltnisse").
		Preload("Transportmittel").
		First(&char, charID).Error

	if err != nil {
		respondWithError(c, http.StatusNotFound, "Character not found")
		return
	}

	// Create temp file
	tempFile, err := os.CreateTemp("", fmt.Sprintf("csv_export_%s_*.csv", char.Name))
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to create temp file")
		return
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Export to CSV
	err = ExportCharToCSV(&char, tempFile.Name())
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, fmt.Sprintf("Failed to export character: %s", err.Error()))
		return
	}

	// Send file as download
	filename := fmt.Sprintf("%s_export.csv", char.Name)
	c.FileAttachment(tempFile.Name(), filename)
}
