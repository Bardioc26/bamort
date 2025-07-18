package main

import (
	"bamort/character"
	"bamort/database"
	"bamort/router"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func main() {
	// Setup
	database.InitDB("")
	database.CreateTestTables()

	// Create test character
	testChar := character.Char{
		CharID:   1,
		CharName: "Test Character",
		Praxispunkte: []character.Praxispunkt{
			{
				Kategorie: "Kampf",
				Anzahl:    5,
			},
		},
	}

	// Save to database
	if err := database.DB.Create(&testChar).Error; err != nil {
		fmt.Printf("Fehler beim Erstellen des Test-Charakters: %v\n", err)
		return
	}

	// Setup router
	gin.SetMode(gin.TestMode)
	r := router.SetupRouter()

	// Test GET practice points
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/characters/1/practice-points", nil)
	r.ServeHTTP(w, req)

	fmt.Printf("Status: %d\n", w.Code)
	fmt.Printf("Body: %s\n", w.Body.String())

	// Test with wrong ID
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/characters/999/practice-points", nil)
	r.ServeHTTP(w2, req2)

	fmt.Printf("\nTest with wrong ID:\n")
	fmt.Printf("Status: %d\n", w2.Code)
	fmt.Printf("Body: %s\n", w2.Body.String())
}
