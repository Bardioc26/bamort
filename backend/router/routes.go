package router

import (
	"bamort/user"

	"github.com/gin-gonic/gin"
)

func BaseRouterGrp(r *gin.Engine) *gin.RouterGroup {
	// Routes
	r.POST("/register", user.RegisterUser)
	r.POST("/login", user.LoginUser)
	protected := r.Group("/api")
	protected.Use(user.AuthMiddleware())
	return protected
}
