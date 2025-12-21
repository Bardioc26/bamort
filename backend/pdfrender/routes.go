package pdfrender

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers protected PDF routes
func RegisterRoutes(r *gin.RouterGroup) {
	pdfGrp := r.Group("/pdf")

	// List available templates (protected)
	pdfGrp.GET("/templates", ListTemplates)

	// Export character to PDF (protected)
	pdfGrp.GET("/export/:id", ExportCharacterToPDF)

	// Cleanup old PDF files (protected)
	pdfGrp.POST("/cleanup", CleanupExportTemp)
}

// RegisterPublicRoutes registers public PDF routes (no authentication required)
func RegisterPublicRoutes(r *gin.Engine) {
	// Get PDF file from xporttemp (public - for direct browser access)
	r.GET("/api/pdf/file/:filename", GetPDFFile)
}
