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
		userGroup.PUT("/language", UpdateLanguage)
	}

	// Admin routes - require admin role
	adminGroup := r.Group("/users")
	adminGroup.Use(RequireAdmin())
	{
		adminGroup.GET("", ListUsers)
		adminGroup.GET("/:id", GetUser)
		adminGroup.PUT("/:id/role", UpdateUserRole)
		adminGroup.PUT("/:id/password", ChangeUserPassword)
		adminGroup.DELETE("/:id", DeleteUser)
	}
}
