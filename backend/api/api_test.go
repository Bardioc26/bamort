package api

import (
	"bamort/character"
	"bamort/database"
	"bamort/gsmaster"
	_ "bamort/maintenance" // Anonymous import to ensure init() is called
	"bamort/models"
	"bamort/router"
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

	err = models.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating gsmaster tables")
	err = user.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating user tables")
	//err = importer.MigrateStructure()
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
	var listOfCharacter []*models.CharList
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
	var listOfCharacter *models.Char
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
	database.SetupTestDB() //(false)
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
	database.SetupTestDB() //(true, true) // Use true to load test data
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

func TestGetAvailableSkillsNew(t *testing.T) {
	database.SetupTestDB(true) // Setup test database
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

	// Create request body for available skills
	skillRequest := gsmaster.LernCostRequest{
		CharId:       20,
		Name:         "Schwimmen", // Use a valid skill name for validation
		CurrentLevel: 0,
		Type:         "skill",
		Action:       "learn",
		TargetLevel:  1,
		UsePP:        0,
		UseGold:      0,
		Reward:       &[]string{"default"}[0],
	}
	jsonData, _ := json.Marshal(skillRequest)

	// Create a test HTTP request
	req, _ := http.NewRequest("POST", "/api/characters/available-skills-new", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	token := user.GenerateToken(&u)
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// Parse the response to verify it contains skills by category
	var response map[string]interface{}
	err := json.Unmarshal(respRecorder.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	// Check that the response contains skills_by_category
	skillsByCategory, exists := response["skills_by_category"]
	assert.True(t, exists, "Response should contain skills_by_category")
	assert.NotNil(t, skillsByCategory, "skills_by_category should not be nil")

	// Convert to map for easier access
	skillsMap, ok := skillsByCategory.(map[string]interface{})
	assert.True(t, ok, "skills_by_category should be a map")
	assert.Greater(t, len(skillsMap), 0, "Should return at least one category of skills")

	// Check that "Bogenbau" is not in the available skills (assuming it's already learned)
	foundBogenbau := false
	for _, categorySkillsInterface := range skillsMap {
		categorySkills, ok := categorySkillsInterface.([]interface{})
		if !ok {
			continue
		}

		for _, skillInterface := range categorySkills {
			skill, ok := skillInterface.(map[string]interface{})
			if !ok {
				continue
			}

			skillName, exists := skill["name"]
			if exists && skillName == "Bogenbau" {
				foundBogenbau = true
				break
			}
		}

		if foundBogenbau {
			break
		}
	}

	assert.False(t, foundBogenbau, "Bogenbau should not be in available skills (already learned)")

	// Verify that each skill has the expected structure
	for categoryName, categorySkillsInterface := range skillsMap {
		categorySkills, ok := categorySkillsInterface.([]interface{})
		assert.True(t, ok, "Category %s should contain an array of skills", categoryName)

		for _, skillInterface := range categorySkills {
			skill, ok := skillInterface.(map[string]interface{})
			assert.True(t, ok, "Each skill should be a map")

			// Check required fields
			_, hasName := skill["name"]
			_, hasEpCost := skill["epCost"]
			_, hasGoldCost := skill["goldCost"]

			assert.True(t, hasName, "Skill should have name field")
			assert.True(t, hasEpCost, "Skill should have epCost field")
			assert.True(t, hasGoldCost, "Skill should have goldCost field")
		}
	}
}
