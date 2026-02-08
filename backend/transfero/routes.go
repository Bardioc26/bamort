package transfero

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the transfer routes
func RegisterRoutes(r *gin.RouterGroup) {
	router := r.Group("/transfer")
	{
		// Export character as JSON (for API consumption)
		router.GET("/export/:id", ExportCharacterHandler)

		// Download character as JSON file
		router.GET("/download/:id", DownloadCharacterHandler)

		// Import character from JSON
		router.POST("/import", ImportCharacterHandler)

		// Full database export/import
		router.POST("/database/export", ExportDatabaseHandler)
		router.POST("/database/import", ImportDatabaseHandler)

		// methods for new importer handling
		router.POST("/vtt-import", dummyproc)
		router.POST("/vtt-export", dummyproc)
		router.POST("/csv-import", dummyproc)
		router.POST("/csv-export", dummyproc)
		router.POST("/moam-import", dummyproc)
		router.POST("/moam-export", dummyproc)
	}
}
