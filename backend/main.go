package main

import (
	"bamort/character"
	"bamort/database"
	"bamort/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()
	//database.DB.AutoMigrate(&models.User{}, &models.Character{}) // Add other models here

	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your frontend's URL
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routes
	r.POST("/register", user.RegisterUser)
	r.POST("/login", user.LoginUser)
	protected := r.Group("/api")
	protected.Use(user.AuthMiddleware())
	protected.GET("/characters", character.GetCharacters)
	protected.POST("/characters", character.CreateCharacter)
	protected.POST("/ausruestung", character.CreateAusruestung)
	protected.GET("/ausruestung/:character_id", character.GetAusruestung)
	protected.PUT("/ausruestung/:ausruestung_id", character.UpdateAusruestung)
	protected.DELETE("/ausruestung/:ausruestung_id", character.DeleteAusruestung)
	protected.POST("/upload", character.UploadFiles)
	protected.GET("/setupcheck", database.SetupCheck)

	r.Run(":8180") // Start server on port 8080
}
