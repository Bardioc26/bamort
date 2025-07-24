package importer

import (
	"fmt"
	"net/http"
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
