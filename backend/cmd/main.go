package main

import (
	"bamort/character"
	"bamort/database"
	"bamort/gsmaster"
	"bamort/router"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()
	//database.DB.AutoMigrate(&models.User{}, &models.Character{}) // Add other models here

	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	// Register your module routes
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)

	r.Run(":8180") // Start server on port 8080
}
