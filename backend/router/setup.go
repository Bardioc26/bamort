package router

import (
	"bamort/importer"
	"bamort/maintenance"
	"bamort/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupGin(r *gin.Engine) {
	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your frontend's URL
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}
func BaseRouterGrp(r *gin.Engine) *gin.RouterGroup {
	// Routes
	r.POST("/register", user.RegisterUser)
	r.POST("/login", user.LoginUser)
	protected := r.Group("/api")
	protected.Use(user.AuthMiddleware())

	protected.POST("/upload", importer.UploadFiles)
	protected.GET("/setupcheck", maintenance.SetupCheck)
	return protected
}
