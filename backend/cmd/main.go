package main

import (
	"bamort/character"
	"bamort/database"
	"bamort/gsmaster"
	"bamort/importer"
	"bamort/maintenance"
	"bamort/router"

	"github.com/gin-gonic/gin"
)

// @title Bamort API
// @version 1
// @description This is the API for Bamort
// @host localhost:8180
// @BasePath /
// @schemes http
func main() {
	database.ConnectDatabase()
	//database.DB.AutoMigrate(&models.User{}, &models.Character{}) // Add other models here

	// Migrate Audit-Log table
	if err := character.MigrateAuditLog(); err != nil {
		panic("Failed to migrate audit log table: " + err.Error())
	}

	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	// Register your module routes
	gsmaster.RegisterRoutes(protected)
	character.RegisterRoutes(protected)
	maintenance.RegisterRoutes(protected)
	importer.RegisterRoutes(protected)

	r.Run(":8180") // Start server on port 8080
}
