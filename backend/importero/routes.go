package importero

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup) {
	charGrp := r.Group("/importer")

	// Import routes
	charGrp.POST("/upload", UploadFiles)
	charGrp.POST("/spells/csv", ImportSpellCSVHandler)

	// Export routes
	exportGrp := charGrp.Group("/export")
	exportGrp.GET("/vtt/:id", ExportCharacterVTTHandler)
	exportGrp.GET("/vtt/:id/file", ExportCharacterVTTFileHandler)
	exportGrp.GET("/csv/:id", ExportCharacterCSVHandler)
	exportGrp.GET("/spells/csv", ExportSpellsCSVHandler)
}
