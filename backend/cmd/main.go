package main

import (
	"log"
	"net/http"

	"github.com/Bardioc26/bamort/pkg/auth"
	"github.com/Bardioc26/bamort/pkg/character"
	"github.com/Bardioc26/bamort/pkg/database"
	"github.com/Bardioc26/bamort/pkg/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// DB verbinden
	err := database.ConnectDB("root", "1234", "localhost", 3306, "rollenspiel_db")
	if err != nil {
		log.Fatal(err)
	}
	// Migration
	if err := database.Migrate(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// Öffentliche Routen
	e.POST("/register", user.Register)
	e.POST("/login", user.Login)
	e.POST("/reset-password", user.ResetPassword)

	// Geschützte Routen (JWT-Middleware)
	protected := e.Group("/api")
	protected.Use(auth.JWTMiddleware)
	{
		// Character
		protected.GET("/characters", character.ListCharacters)
		protected.GET("/characters/:id", character.GetCharacterByID)
		protected.POST("/characters", character.CreateCharacter)
		protected.PUT("/characters/:id", character.UpdateCharacter)
		protected.DELETE("/characters/:id", character.DeleteCharacter)
	}

	e.Logger.Fatal(e.Start(":8080"))
}
