package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers user-related routes
func RegisterRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	{
		// Protected routes - require authentication
		userGroup.GET("/profile", GetUserProfile)
		userGroup.PUT("/email", UpdateEmail)
		userGroup.PUT("/password", UpdatePassword)
	}
}
