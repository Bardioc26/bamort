package character

import (
	"bamort/database"
	"bamort/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDebugPracticePoints(t *testing.T) {
	// Setup test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Migrate the schema
	err := models.MigrateStructure()
	assert.NoError(t, err)

	// Also migrate skills and equipment to avoid errors
	/*
		err = skills.MigrateStructure()
		assert.NoError(t, err)
		err = equipment.MigrateStructure()
		assert.NoError(t, err)*/

	// Create a test character manually using GORM
	character := &models.Char{
		BamortBase: models.BamortBase{
			Name: "Test Character",
		},
		Rasse: "Human",
		Typ:   "Hx",
	}
	err = database.DB.Create(character).Error
	assert.NoError(t, err)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register routes
	api := router.Group("/api")
	RegisterRoutes(api)

	// Test the endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Logf("Status: %d", w.Code)
	t.Logf("Body: %s", w.Body.String())

	// Check if the character exists in database
	var testChar models.Char
	err = database.DB.First(&testChar, character.ID).Error
	t.Logf("Character exists in DB: %v, ID: %d, Error: %v", err == nil, character.ID, err)
}
