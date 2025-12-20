package pdfrender

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	pdfGrp := r.Group("/pdf")

	// List available templates
	pdfGrp.GET("/templates", ListTemplates)

	// Export character to PDF
	pdfGrp.GET("/export/:id", ExportCharacterToPDF)
}
