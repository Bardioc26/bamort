package transfer

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers the transfer routes
func RegisterRoutes(r *gin.RouterGroup) {
	transfer := r.Group("/transfer")
	{
		// Export character as JSON (for API consumption)
		transfer.GET("/export/:id", ExportCharacterHandler)

		// Download character as JSON file
		transfer.GET("/download/:id", DownloadCharacterHandler)

		// Import character from JSON
		transfer.POST("/import", ImportCharacterHandler)

		// Full database export/import
		transfer.POST("/database/export", ExportDatabaseHandler)
		transfer.POST("/database/import", ImportDatabaseHandler)
	}
}
