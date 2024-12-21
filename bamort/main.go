package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	ConnectDatabase()
	DB.AutoMigrate(&User{}, &Character{}, &Eigenschaften{}) // Add other models here

	r := gin.Default()

	// Routes
	r.POST("/register", RegisterUser)
	r.POST("/login", LoginUser)
	r.GET("/characters", GetCharacters)
	r.POST("/characters", CreateCharacter)

	r.Run(":8080") // Start server on port 8080
}
