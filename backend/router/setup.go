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
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // Cache preflight for 12 hours
	}))
}
