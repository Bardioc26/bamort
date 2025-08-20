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

	character, err3 := ImportVTTJSON(vttFileName)
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
