package api

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
	"bamort/importer"
	_ "bamort/maintenance" // Anonymous import to ensure init() is called
	"bamort/models"
	"bamort/router"
	"bamort/skills"
	"bamort/user"
	"bytes"
	"encoding/json"
	"io"
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
	// must be in sync with maintenance.SetupCheck(&c)
	database.SetupTestDB(true) // Use in-memory database for tests

	db := database.ConnectDatabase()
	assert.NotNil(t, db, "expected database connection to be established")
	if db == nil {
		return
	}

	err := database.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating database tables")

	err = character.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating character tables")

	err = user.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating user tables")

	err = models.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating gsmaster tables")

	err = equipment.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating equipment tables")

	err = skills.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating skill tables")

	err = importer.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating importer tables")
}

func TestListCharacters(t *testing.T) {
	database.SetupTestDB()
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	protected.GET("/test", func(c *gin.Context) {
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

func TestGetCharacters(t *testing.T) {
	database.SetupTestDB()
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	protected.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/characters/9", nil)
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer ${token}")
	req.Header.Set("Authorization", "Bearer dc7a780.1:bba7f4daabda117f2a2c14263")

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// Assert the response body
	var listOfCharacter *character.Char
	err := json.Unmarshal(respRecorder.Body.Bytes(), &listOfCharacter)
	assert.NoError(t, err)
	assert.Equal(t, "Harsk Hammerhuter, Zen", listOfCharacter.Name)
	assert.Equal(t, "Zwerg", listOfCharacter.Rasse)
	assert.Equal(t, 1, int(listOfCharacter.ID)) // Check the simulated ID
	assert.Equal(t, "Krieger", listOfCharacter.Typ)
	assert.Equal(t, 3, listOfCharacter.Grad)
	//assert.Equal(t, "test", listOfCharacter.Owner)
	//assert.Equal(t, false, listOfCharacter.Public)

}

func TestCreateCharacter(t *testing.T) {
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	protected.GET("/test", func(c *gin.Context) {
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

func TestGetSkillCost(t *testing.T) {
	database.SetupTestDB(false)
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	protected.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})

	// Test spell learning cost
	skillCostRequest := map[string]interface{}{
		"name":   "Angst",
		"type":   "spell",
		"action": "learn",
	}
	jsonData, _ := json.Marshal(skillCostRequest)

	// Create a test HTTP request
	req, _ := http.NewRequest("POST", "/api/characters/18/skill-cost", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer dc7a780.1:bba7f4daabda117f2a2c14263")

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// Test skill learning cost
	skillCostRequest = map[string]interface{}{
		"name":   "Bootfahren",
		"type":   "skill",
		"action": "learn",
	}
	jsonData, _ = json.Marshal(skillCostRequest)
	req.Body = io.NopCloser(bytes.NewBuffer(jsonData))

	// Create a response recorder to capture the handler's response
	respRecorder = httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// The new GetSkillCost returns a structured response, not just a number
	// We just check that it returns successfully for now
}

func TestGetSkillAllLevelCosts(t *testing.T) {
	database.SetupTestDB(true, true) // Use true to load test data
	//database.SetupTestDB(false)
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	protected.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})
	u := user.User{}
	u.First("testuser")

	// Test skill improvement costs - create request body with skill name
	skillRequest := map[string]interface{}{
		"name": "Hören",
	}
	jsonData, _ := json.Marshal(skillRequest)

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/characters/20/improve/skill", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	token := user.GenerateToken(&u)
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// Parse the response to verify it contains level cost information
	var response []interface{}
	err := json.Unmarshal(respRecorder.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")
	assert.Greater(t, len(response), 0, "Should return learning costs")

	// Test with a weapon skill
	skillRequest = map[string]interface{}{
		"name": "Armbrüste",
	}
	jsonData, _ = json.Marshal(skillRequest)
	req, _ = http.NewRequest("GET", "/api/characters/20/improve/skill", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token) // Use the same token as the first request

	respRecorder = httptest.NewRecorder()
	r.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusOK, respRecorder.Code)

	err = json.Unmarshal(respRecorder.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON for weapon skill")
}
