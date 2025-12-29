package pdfrender

import (
	"bamort/config"
	"bamort/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// TemplateInfo represents information about an available export template
type TemplateInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListTemplates returns a list of available export templates
func ListTemplates(c *gin.Context) {
	templatesDir := config.Cfg.TemplatesDir

	// Read template directories
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read templates directory"})
		return
	}

	var templates []TemplateInfo
	for _, entry := range entries {
		if entry.IsDir() {
			templates = append(templates, TemplateInfo{
				ID:          entry.Name(),
				Name:        entry.Name(),
				Description: "PDF Export Template: " + entry.Name(),
			})
		}
	}

	c.JSON(http.StatusOK, templates)
}

// ExportCharacterToPDF exports a character to PDF and saves it to xporttemp directory
// Query params:
//   - template: template ID to use (default: "Default_A4_Quer")
//   - showUserName: whether to show user name (default: false)
//
// Returns JSON with filename: {"filename": "CharacterName_20231225_143045.pdf"}
func ExportCharacterToPDF(c *gin.Context) {
	// Get character ID
	charID := c.Param("id")
	if charID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Character ID is required"})
		return
	}

	// Load character
	char := &models.Char{}
	if err := char.FirstID(charID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	// Get template parameter (default to Default_A4_Quer)
	templateID := c.DefaultQuery("template", "Default_A4_Quer")

	// Map character to view model
	viewModel, err := MapCharacterToViewModel(char)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to map character: " + err.Error()})
		return
	}

	// Load templates
	templateDir := filepath.Join(config.Cfg.TemplatesDir, templateID)
	loader := NewTemplateLoader(templateDir)
	if err := loader.LoadTemplates(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load templates: " + err.Error()})
		return
	}

	renderer := NewPDFRenderer()
	currentDate := time.Now().Format("02.01.2006")

	// Generate all pages with continuations if needed
	var allPDFs [][]byte

	// Page 1: Stats
	page1PDFs, err := RenderPageWithContinuations(viewModel, "page_1.html", 1, currentDate, loader, renderer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render page 1: " + err.Error()})
		return
	}
	//allPDFs = append(allPDFs, page1PDFs...)
	for _, pdf := range page1PDFs {
		allPDFs = append(allPDFs, pdf)
	}

	// Page 2: Play
	page2PDFs, err := RenderPageWithContinuations(viewModel, "page_2.html", 2, currentDate, loader, renderer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render page 2: " + err.Error()})
		return
	}
	//allPDFs = append(allPDFs, page2PDFs...)
	for _, pdf := range page2PDFs {
		allPDFs = append(allPDFs, pdf)
	}

	// Page 3: Spells
	page3PDFs, err := RenderPageWithContinuations(viewModel, "page_3.html", 3, currentDate, loader, renderer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render page 3: " + err.Error()})
		return
	}
	//allPDFs = append(allPDFs, page3PDFs...)
	for _, pdf := range page3PDFs {
		allPDFs = append(allPDFs, pdf)
	}

	// Page 4: Equipment
	page4PDFs, err := RenderPageWithContinuations(viewModel, "page_4.html", 4, currentDate, loader, renderer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render page 4: " + err.Error()})
		return
	}
	//allPDFs = append(allPDFs, page4PDFs...)
	for _, pdf := range page4PDFs {
		allPDFs = append(allPDFs, pdf)
	}

	// Merge PDFs if needed
	var finalPDF []byte
	if len(allPDFs) == 1 {
		finalPDF = allPDFs[0]
	} else {
		// Merge multiple PDFs
		tmpDir := fmt.Sprintf("/tmp/bamort_pdf_export_%d", time.Now().UnixNano())
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory"})
			return
		}
		defer os.RemoveAll(tmpDir)

		// Save individual PDFs
		var filePaths []string
		for i, pdf := range allPDFs {
			filename := fmt.Sprintf("%s/page_%d.pdf", tmpDir, i)
			if err := os.WriteFile(filename, pdf, 0644); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write temporary PDF"})
				return
			}
			filePaths = append(filePaths, filename)
		}

		// Merge PDFs
		combinedPath := fmt.Sprintf("%s/combined_%d.pdf", tmpDir, time.Now().UnixNano())
		if err := api.MergeCreateFile(filePaths, combinedPath, false, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to merge PDFs: " + err.Error()})
			return
		}

		// Read combined PDF
		finalPDF, err = os.ReadFile(combinedPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read combined PDF"})
			return
		}
	}

	// Ensure export temp directory exists
	if err := EnsureExportTempDir(config.Cfg.ExportTempDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create export directory"})
		return
	}

	// Generate filename
	filename := GenerateExportFilename(char.Name, time.Now())
	filePath := filepath.Join(config.Cfg.ExportTempDir, filename)

	// Save PDF to file
	if err := os.WriteFile(filePath, finalPDF, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save PDF file: " + err.Error()})
		return
	}

	// Return filename
	c.JSON(http.StatusOK, gin.H{"filename": filename})
}

// GetPDFFile serves a PDF file from the xporttemp directory
func GetPDFFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Prevent path traversal attacks - only allow base filename
	if filepath.Base(filename) != filename {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	// Only allow .pdf files
	if filepath.Ext(filename) != ".pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}

	// Construct full path
	filePath := filepath.Join(config.Cfg.ExportTempDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Serve the file
	c.File(filePath)
}

// CleanupExportTemp removes PDF files older than 7 days from xporttemp directory
func CleanupExportTemp(c *gin.Context) {
	// Clean up files older than 7 days
	maxAge := 7 * 24 * time.Hour
	count, err := CleanupOldFiles(config.Cfg.ExportTempDir, maxAge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup files: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": count,
		"message": fmt.Sprintf("Deleted %d files older than 7 days", count),
	})
}
