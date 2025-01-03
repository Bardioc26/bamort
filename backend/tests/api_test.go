package tests

import (
	"bamort/character"
	"bamort/maintenance"
	"bamort/router"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock character creation handler
func MockCreateCharacter(c *gin.Context) {
	var character Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simulate saving the character and returning success
	character.ID = 1 // Simulated ID from the database
	c.JSON(http.StatusCreated, character)
}

// Character struct for testing
type Character struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Race string `json:"race"`
}

func TestSetupCheck(t *testing.T) {
	//r := gin.Default()
	c := gin.Context{}
	maintenance.SetupCheck(&c)
	assert.Empty(t, nil, "expected NIL to be empty")
	/*
		SetupTestDB()
		TestCreateChar(t)
		// Initialize a Gin router
		r := gin.Default()
		router.SetupGin(r)

		// Routes
		protected := router.BaseRouterGrp(r)
		// Character routes
		rCharGrp := router.CharRouterGrp(protected)
		rCharGrp.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
		})

		// Create a test HTTP request
		req, _ := http.NewRequest("GET", "/api/characters", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer ${token}")

		// Create a response recorder to capture the handler's response
		respRecorder := httptest.NewRecorder()

		// Perform the test request
		r.ServeHTTP(respRecorder, req)

		// Assert the response status code
		assert.Equal(t, http.StatusOK, respRecorder.Code)

		// Assert the response body
		var listOfCharacter []*character.CharList
		err := json.Unmarshal(respRecorder.Body.Bytes(), &listOfCharacter)
		assert.NoError(t, err)
		assert.Equal(t, "Harsk Hammerhuter, Zen", listOfCharacter[0].Name)
		assert.Equal(t, "Zwerg", listOfCharacter[0].Rasse)
		assert.Equal(t, 1, int(listOfCharacter[0].ID)) // Check the simulated ID
		assert.Equal(t, "Krieger", listOfCharacter[0].Typ)
		assert.Equal(t, 3, listOfCharacter[0].Grad)
		assert.Equal(t, "test", listOfCharacter[0].Owner)
		assert.Equal(t, false, listOfCharacter[0].Public)
	*/
}
func TestListCharacters(t *testing.T) {
	SetupTestDB()
	TestCreateChar(t)
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	// Character routes
	rCharGrp := router.CharRouterGrp(protected)
	rCharGrp.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/characters", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer ${token}")

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// Assert the response body
	var listOfCharacter []*character.CharList
	err := json.Unmarshal(respRecorder.Body.Bytes(), &listOfCharacter)
	assert.NoError(t, err)
	assert.Equal(t, "Harsk Hammerhuter, Zen", listOfCharacter[0].Name)
	assert.Equal(t, "Zwerg", listOfCharacter[0].Rasse)
	assert.Equal(t, 1, int(listOfCharacter[0].ID)) // Check the simulated ID
	assert.Equal(t, "Krieger", listOfCharacter[0].Typ)
	assert.Equal(t, 3, listOfCharacter[0].Grad)
	assert.Equal(t, "test", listOfCharacter[0].Owner)
	assert.Equal(t, false, listOfCharacter[0].Public)

}

// TestCreateCharacter tests the POST /characters endpoint
func TestCreateCharacter(t *testing.T) {
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	// Character routes
	rCharGrp := router.CharRouterGrp(protected)
	rCharGrp.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})
	// Define the test case input
	testCharacter := Character{
		Name: "Aragorn",
		Race: "Human",
	}
	jsonData, _ := json.Marshal(testCharacter)

	// Create a test HTTP request
	req, _ := http.NewRequest("POST", "/characters", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusCreated, respRecorder.Code)

	// Assert the response body
	var createdCharacter Character
	err := json.Unmarshal(respRecorder.Body.Bytes(), &createdCharacter)
	assert.NoError(t, err)
	assert.Equal(t, "Aragorn", createdCharacter.Name)
	assert.Equal(t, "Human", createdCharacter.Race)
	assert.Equal(t, 1, createdCharacter.ID) // Check the simulated ID
}
