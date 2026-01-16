package system

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers protected system routes with the Gin router
func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	system := r.Group("/system")
	{
		system.GET("/health", HealthHandler(db))
		system.GET("/version", VersionHandler(db))
	}
}

// RegisterPublicRoutes registers public system routes (no authentication required)
func RegisterPublicRoutes(r *gin.Engine, db *gorm.DB) {
	system := r.Group("/api/system")
	{
		system.GET("/health", HealthHandler(db))
		system.GET("/version", VersionHandler(db))
	}
}
