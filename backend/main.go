package main

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
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
	// Character routes
	rCharGrp := protected.Group("/characters")
	rCharGrp.GET("", character.ListCharacters)
	rCharGrp.POST("", character.CreateCharacter)
	rCharGrp.GET("/:id", character.GetCharacter)
	rCharGrp.PUT("/:id", character.UpdateCharacter)
	rCharGrp.DELETE("/:id", character.DeleteCharacter)
	//rCharGrp.GET("/{id}/skills", character.ListSkills)                      //	List skills for a character
	//rCharGrp.GET("/{id}/skills/{id}", character.GetSkill)                      //	get a skill for a character
	//rCharGrp.PUT("/{id}/skills/{id}", character.UpdateSkill)                      //	Update skill for a character
	//rCharGrp.POST("/{id}/skills", character.AddSkill)                       //Add a skill to a character
	//rCharGrp.DELETE("/{id}/skills/{id}", character.DeleteSkill)                       //ADEletedd a skill to a character

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
	protected.POST("/upload", character.UploadFiles)
	protected.GET("/setupcheck", database.SetupCheck)

	r.Run(":8180") // Start server on port 8080
}
