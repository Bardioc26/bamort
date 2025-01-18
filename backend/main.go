package main

import (
	"bamort/database"
	"bamort/importer"
	"bamort/maintenance"
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
	// Character routes
	//rCharGrp := router.CharRouterGrp(protected)
	router.CharRouterGrp(protected)

	/*
		rCharGrp.GET("/{id}/equipment", equipment.ListAusruestung) //	List equipment for a character
		//rCharGrp.POST("/{id}/equipment", equipment.AddAusruestung)           //	Add equipment to a character
		rCharGrp.PUT("/{id}/equipment/{id}", equipment.UpdateAusruestung)    //	List equipment for a character
		rCharGrp.DELETE("/{id}/equipment/{id}", equipment.DeleteAusruestung) //	List equipment for a character
		//rCharGrp.GET("/{id}/spells", spells.ListAusruestung)                 //	List equipment for a character
		//rCharGrp.POST("/{id}/spells", equipment.AddAusruestung)              //	Add equipment to a character
		rCharGrp.PUT("/{id}/spells/{id}", equipment.UpdateAusruestung)    //	List equipment for a character
		rCharGrp.DELETE("/{id}/spells/{id}", equipment.DeleteAusruestung) //	List equipment for a character
		rCharGrp.POST("/{id}/spells", character.DeleteCharacter)          //	Add a spell to a character
		rCharGrp.GET("/{id}/spells", character.DeleteCharacter)           //	List spells for a character
		// Equipment routes
		protected.POST("/ausruestung", equipment.CreateAusruestung)
		//protected.GET("/ausruestung/:character_id", equipment.GetAusruestung)
		protected.PUT("/ausruestung/:ausruestung_id", equipment.UpdateAusruestung)
		protected.DELETE("/ausruestung/:ausruestung_id", equipment.DeleteAusruestung)
	*/
	router.MaintenanceRouterGrp(protected)
	protected.POST("/upload", importer.UploadFiles)
	protected.GET("/setupcheck", maintenance.SetupCheck)

	r.Run(":8180") // Start server on port 8080
}
