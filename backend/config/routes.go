package config

import "github.com/gin-gonic/gin"

// RegisterRoutes registers config-related routes (protected)
func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/version", Versionsinfo)
}

// RegisterPublicRoutes registers public config routes (no auth required)
func RegisterPublicRoutes(r *gin.Engine) {
	// Public version endpoint - no authentication required
	public := r.Group("/api/public")
	public.GET("/version", Versionsinfo)
}
