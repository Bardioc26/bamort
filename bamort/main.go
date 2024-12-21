package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	ConnectDatabase()
	DB.AutoMigrate(&User{}, &Character{}, &Eigenschaften{}) // Add other models here

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
	r.POST("/register", RegisterUser)
	r.POST("/login", LoginUser)
	r.GET("/characters", GetCharacters)
	r.POST("/characters", CreateCharacter)
	r.POST("/ausruestung", CreateAusruestung)
	r.GET("/ausruestung/:character_id", GetAusruestung)
	r.PUT("/ausruestung/:ausruestung_id", UpdateAusruestung)
	r.DELETE("/ausruestung/:ausruestung_id", DeleteAusruestung)

	r.Run(":8180") // Start server on port 8080
}
