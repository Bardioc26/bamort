package router

import (
	"bamort/character"
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
	// Character routes
	return protected
}
func CharRouterGrp(rt *gin.RouterGroup) *gin.RouterGroup {
	rCharGrp := rt.Group("/characters")
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
	return rCharGrp
}
