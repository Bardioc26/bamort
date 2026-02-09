package importer

import (
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all import/export API endpoints
// Following the plan: POST /detect, POST /import, GET /adapters, GET /history, etc.
func RegisterRoutes(r *gin.RouterGroup) {
	// Rate limiters per endpoint as specified in the plan
	detectLimiter := NewRateLimiter(10, time.Minute) // 10/min
	importLimiter := NewRateLimiter(5, time.Minute)  // 5/min
	exportLimiter := NewRateLimiter(20, time.Minute) // 20/min

	// File size limit (10MB as per plan)
	maxFileSize := int64(10 << 20)

	importer := r.Group("/import")
	importer.Use(ValidateFileSizeMiddleware(maxFileSize))

	// Format detection endpoint
	importer.POST("/detect", detectLimiter.Middleware(), DetectHandler)

	// Import character from external format
	importer.POST("/import", importLimiter.Middleware(), ImportHandler)

	// List registered adapters
	importer.GET("/adapters", ListAdaptersHandler)

	// Get user's import history
	importer.GET("/history", ImportHistoryHandler)

	// Get details for specific import
	importer.GET("/history/:id", ImportDetailsHandler)

	// Export character to external format
	importer.POST("/export/:id", exportLimiter.Middleware(), ExportHandler)
}
