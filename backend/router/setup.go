package router

import (
	"bamort/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupGin(r *gin.Engine) {
	// Build allowed origins list from configuration
	allowedOrigins := []string{
		config.Cfg.FrontendURL,
		"http://localhost:5173",    // Development frontend
		"https://bamort.trokan.de", // Production frontend
	}

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your frontend's URL
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // Cache preflight for 12 hours
	}))
}
