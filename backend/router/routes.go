package router

import (
	"bamort/user"

	"github.com/gin-gonic/gin"
)

func BaseRouterGrp(r *gin.Engine) *gin.RouterGroup {
	// Routes
	r.POST("/register", user.RegisterUser)
	r.POST("/login", user.LoginUser)

	// Password Reset Routes (unprotected)
	r.POST("/password-reset/request", user.RequestPasswordReset)
	r.GET("/password-reset/validate/:token", user.ValidateResetToken)
	r.POST("/password-reset/reset", user.ResetPassword)

	protected := r.Group("/api")
	protected.Use(user.AuthMiddleware())
	return protected
}
