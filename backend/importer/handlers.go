package importer

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Global registry instance (initialized on startup)
var globalRegistry *AdapterRegistry

// InitializeRegistry initializes the global adapter registry
func InitializeRegistry(registry *AdapterRegistry) {
	globalRegistry = registry
}

// DetectHandler handles format detection for uploaded files
// POST /api/import/detect
// Rate limit: 10 requests/minute per user
func DetectHandler(c *gin.Context) {
	// Accept multipart file upload
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Validate JSON depth if JSON file
	if err := ValidateJSONDepth(data, 100); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}

	// Detect format using global registry
	if globalRegistry == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Import service not initialized"})
		return
	}

	adapterID, confidence, err := globalRegistry.Detect(data, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Detection failed: %v", err)})
		return
	}

	// Get adapter metadata for suggested name
	adapter := globalRegistry.Get(adapterID)
	response := gin.H{
		"adapter_id": adapterID,
		"confidence": confidence,
	}
	if adapter != nil {
		response["suggested_adapter_name"] = adapter.Name
	}

	c.JSON(http.StatusOK, response)
}

// ImportHandler handles character import from external formats
// POST /api/import/import
// Rate limit: 5 requests/minute per user
func ImportHandler(c *gin.Context) {
	// Get user ID from context
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Accept file and optional adapter_id
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Validate JSON depth
	if err := ValidateJSONDepth(data, 100); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}

	// Get adapter ID from form or detect
	adapterID := c.PostForm("adapter_id")
	if adapterID == "" {
		// Auto-detect
		if globalRegistry == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Import service not initialized"})
			return
		}
		detectedID, confidence, err := globalRegistry.Detect(data, header.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Detection failed: %v", err)})
			return
		}
		if confidence < 0.7 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      "Could not reliably detect format",
				"confidence": confidence,
				"hint":       "Please specify adapter_id explicitly",
			})
			return
		}
		adapterID = detectedID
	}

	// Import via adapter
	charImport, err := globalRegistry.Import(adapterID, data)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("Import failed: %v", err)})
		return
	}

	// Validate the imported character
	validator := NewValidator()
	validationResult := validator.ValidateCharacter(charImport)
	if !validationResult.Valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Character validation failed",
			"errors": validationResult.Errors,
		})
		return
	}

	// Collect warnings but don't block
	warnings := validationResult.Warnings

	// Convert BMRTCharacter to CharacterImport for import logic
	charData := &charImport.CharacterImport

	// Import character with transaction safety
	result, err := ImportCharacter(charData, userID, adapterID, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create character: %v", err)})
		return
	}

	// Merge validation warnings into result
	result.Warnings = append(result.Warnings, warnings...)

	c.JSON(http.StatusOK, result)
}

// ListAdaptersHandler returns all registered adapters
// GET /api/import/adapters
func ListAdaptersHandler(c *gin.Context) {
	if globalRegistry == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Import service not initialized"})
		return
	}

	adapters := globalRegistry.GetHealthy()
	c.JSON(http.StatusOK, gin.H{
		"adapters": adapters,
		"count":    len(adapters),
	})
}

// ImportHistoryHandler returns user's import history with pagination
// GET /api/import/history?page=1&per_page=20
func ImportHistoryHandler(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	// Query import history
	var histories []ImportHistory
	var total int64

	db := database.DB.Model(&ImportHistory{}).Where("user_id = ?", userID)
	db.Count(&total)

	err := db.Order("imported_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&histories).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch import history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"histories": histories,
		"total":     total,
		"page":      page,
		"per_page":  perPage,
		"pages":     (total + int64(perPage) - 1) / int64(perPage),
	})
}

// ImportDetailsHandler returns detailed information about a specific import
// GET /api/import/history/:id
func ImportDetailsHandler(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	importID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid import ID"})
		return
	}

	// Fetch import history with ownership check
	var history ImportHistory
	err = database.DB.Where("id = ? AND user_id = ?", importID, userID).
		First(&history).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Import not found"})
		return
	}

	// Fetch related master data imports
	var masterDataImports []MasterDataImport
	database.DB.Where("import_history_id = ?", importID).
		Find(&masterDataImports)

	c.JSON(http.StatusOK, gin.H{
		"history":             history,
		"master_data_imports": masterDataImports,
	})
}

// ExportHandler exports a character to an external format
// POST /api/import/export/:id?adapter_id=foundry-vtt-v1
func ExportHandler(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	charID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid character ID"})
		return
	}

	// Load character with ownership check
	var char models.Char
	err = database.DB.Where("id = ? AND user_id = ?", charID, userID).
		First(&char).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found or access denied"})
		return
	}

	// Determine adapter ID (override or original)
	adapterID := c.Query("adapter_id")
	if adapterID == "" {
		// Try to get original adapter from import history
		var history ImportHistory
		err = database.DB.Where("character_id = ?", charID).
			Order("imported_at DESC").
			First(&history).Error

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No adapter specified and character has no import history",
				"hint":  "Specify adapter_id query parameter",
			})
			return
		}
		adapterID = history.AdapterID
	}

	// Check adapter exists and is healthy
	if globalRegistry == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Import service not initialized"})
		return
	}

	adapter := globalRegistry.Get(adapterID)
	if adapter == nil {
		// Suggest available adapters
		available := globalRegistry.GetHealthy()
		availableIDs := make([]string, len(available))
		for i, a := range available {
			availableIDs[i] = a.ID
		}

		c.JSON(http.StatusConflict, gin.H{
			"error":              fmt.Sprintf("Adapter '%s' not available", adapterID),
			"available_adapters": availableIDs,
		})
		return
	}

	if !adapter.Healthy {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf("Adapter '%s' is currently unhealthy", adapterID),
		})
		return
	}

	// Convert Char to CharacterImport
	charImport, err := ConvertCharToImport(&char)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to convert character: %v", err)})
		return
	}

	// Export via adapter
	exportedData, err := globalRegistry.Export(adapterID, charImport)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("Export failed: %v", err)})
		return
	}

	// Return file download
	filename := fmt.Sprintf("%s_%s.json", char.Name, adapterID)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Data(http.StatusOK, "application/json", exportedData)
}

// respondWithError is a helper to send error responses
func respondWithError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
