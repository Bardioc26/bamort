package character

import (
	"bamort/database"
	"bamort/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCharacterImage(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	router := gin.Default()
	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})
	RegisterRoutes(protected)

	// Get existing character
	var char models.Char
	err := char.FirstID("18")
	assert.NoError(t, err, "Test character 18 should exist")

	// Prepare image data
	imageData := "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="

	requestBody := map[string]string{
		"image": imageData,
	}
	jsonData, _ := json.Marshal(requestBody)

	// Update character image
	req, _ := http.NewRequest("PUT", "/api/characters/18/image", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Should successfully update image")

	// Verify image was saved
	var updatedChar models.Char
	err = updatedChar.FirstID("18")
	assert.NoError(t, err)
	assert.Equal(t, imageData, updatedChar.Image, "Image should be updated in database")
}

func TestUpdateCharacterImageInvalidID(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	router := gin.Default()
	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})
	RegisterRoutes(protected)

	requestBody := map[string]string{
		"image": "data:image/png;base64,test",
	}
	jsonData, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("PUT", "/api/characters/99999/image", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent character")
}

func TestUpdateCharacterImageInvalidData(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()

	router := gin.Default()
	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})
	RegisterRoutes(protected)

	req, _ := http.NewRequest("PUT", "/api/characters/18/image", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 for invalid JSON")
}
