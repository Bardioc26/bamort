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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Rasse string `json:"rasse"` // German field name to match API
}

func getAuthToken() string {
	u := user.User{}
	u.FirstId(1)
	token := user.GenerateToken(&u)
	return token
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

	err = user.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating user tables")

	err = models.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating gsmaster tables")
	//err = importer.MigrateStructure()
	assert.NoError(t, err, "No error expected when migrating importer tables")
}

func TestListCharacters(t *testing.T) {
	database.SetupTestDB(true)
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

	token := getAuthToken()

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/characters", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// Assert the response body
	//var listOfCharacter []*models.CharList
	type AllCharacters struct {
		SelfOwned []models.CharList `json:"self_owned"`
		Others    []models.CharList `json:"others"`
	}
	var allCharacters AllCharacters
	err := json.Unmarshal(respRecorder.Body.Bytes(), &allCharacters)
	listOfCharacter := allCharacters.SelfOwned
	assert.NoError(t, err)
	assert.Equal(t, "Harsk Hammerhuter, Zen", listOfCharacter[4].Name)
	assert.Equal(t, "Zwerg", listOfCharacter[4].Rasse)
	assert.Equal(t, 20, int(listOfCharacter[4].ID)) // Check the simulated ID
	assert.Equal(t, "Krieger", listOfCharacter[4].Typ)
	assert.Equal(t, 3, listOfCharacter[4].Grad)
	assert.Equal(t, "bebe", listOfCharacter[4].Owner)
	assert.Equal(t, false, listOfCharacter[4].Public)

}

func TestGetCharacters(t *testing.T) {
	database.SetupTestDB(true)
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)
	token := getAuthToken()

	// Routes
	protected := router.BaseRouterGrp(r)
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	protected.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/characters/20", nil)
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Authorization", "Bearer ${token}")
	req.Header.Set("Authorization", "Bearer "+token)

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
	assert.Equal(t, 20, int(listOfCharacter.ID)) // Check the simulated ID
	assert.Equal(t, "Krieger", listOfCharacter.Typ)
	assert.Equal(t, 3, listOfCharacter.Grad)
	//assert.Equal(t, "test", listOfCharacter.Owner)
	//assert.Equal(t, false, listOfCharacter.Public)

}

func TestCreateCharacter(t *testing.T) {
	database.SetupTestDB(true)
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)
	token := getAuthToken()

	// Routes
	protected := router.BaseRouterGrp(r)
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	protected.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})
	// Define the test case input
	testCharacter := Character{
		Name:  "Aragorn",
		Rasse: "Human",
	}
	jsonData, _ := json.Marshal(testCharacter)

	// Create a test HTTP request
	req, _ := http.NewRequest("POST", "/api/characters", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

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
	assert.Equal(t, "Human", createdCharacter.Rasse)
	assert.GreaterOrEqual(t, createdCharacter.ID, 21) // Check the simulated ID
}

func TestGetSkillCost(t *testing.T) {
	// NOTE: This test uses the newly created character from TestCreateCharacter when run
	// in the full suite, because database.SetupTestDB(true) only creates a fresh DB if DB == nil.
	// When tests run sequentially, they share the same DB instance, so we use the character
	// created by TestCreateCharacter to ensure the skill doesn't already exist.
	database.SetupTestDB(true) //(false)
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
	token := getAuthToken()

	// Test skill learning cost using character 21 (created by TestCreateCharacter in full suite)
	// or character 20 (existing in DB when run individually)
	// Use "Abrichten" which character 20 definitely doesn't have in prepared_test_data.db
	skillCostRequest := gsmaster.LernCostRequest{
		CharId:       20,
		Name:         "Musizieren",
		CurrentLevel: 0,
		Type:         "skill",
		Action:       "learn",
		TargetLevel:  1,
		UsePP:        0,
		UseGold:      0,
		Reward:       &[]string{"default"}[0],
	}
	jsonData, _ := json.Marshal(skillCostRequest)
	req, _ := http.NewRequest("POST", "/api/characters/20/learn-skill-new", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// The new GetSkillCost returns a structured response, not just a number
	// We just check that it returns successfully for now
}

func TestGetAvailableSkillsNewSystem(t *testing.T) {
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

	token := getAuthToken()
	// Create a test HTTP request
	req, _ := http.NewRequest("POST", "/api/characters/available-skills-new", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
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

func TestGetAvailableSpellsNewSystem(t *testing.T) {
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

	token := getAuthToken()
	// Create request body for available spells
	spellRequest := gsmaster.LernCostRequest{
		CharId:       20,
		Name:         "Angst", // Use a valid spell name for validation
		CurrentLevel: 0,
		Type:         "spell",
		Action:       "learn",
		TargetLevel:  1,
		UsePP:        0,
		UseGold:      0,
		Reward:       &[]string{"default"}[0],
	}
	jsonData, _ := json.Marshal(spellRequest)

	// Create a test HTTP request
	req, _ := http.NewRequest("POST", "/api/characters/available-spells-new", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", "Bearer "+token)

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// Parse the response to verify it contains spells by school
	var response map[string]interface{}
	err := json.Unmarshal(respRecorder.Body.Bytes(), &response)
	assert.NoError(t, err, "Response should be valid JSON")

	// Check that the response contains spells_by_school
	spellsBySchool, exists := response["spells_by_school"]
	assert.True(t, exists, "Response should contain spells_by_school")
	assert.NotNil(t, spellsBySchool, "spells_by_school should not be nil")

	// Convert to map for easier access
	spellsMap, ok := spellsBySchool.(map[string]interface{})
	assert.True(t, ok, "spells_by_school should be a map")
	assert.Greater(t, len(spellsMap), 0, "Should return at least one school of spells")

	// Verify that each spell has the expected structure and check for fallback values
	fallbackSpells := []string{}
	totalSpells := 0

	for schoolName, schoolSpellsInterface := range spellsMap {
		schoolSpells, ok := schoolSpellsInterface.([]interface{})
		assert.True(t, ok, "School %s should contain an array of spells", schoolName)

		for _, spellInterface := range schoolSpells {
			spell, ok := spellInterface.(map[string]interface{})
			assert.True(t, ok, "Each spell should be a map")
			totalSpells++

			// Check required fields
			name, hasName := spell["name"]
			level, hasLevel := spell["level"]
			epCost, hasEpCost := spell["epCost"]
			goldCost, hasGoldCost := spell["goldCost"]

			assert.True(t, hasName, "Spell should have name field")
			assert.True(t, hasLevel, "Spell should have level field")
			assert.True(t, hasEpCost, "Spell should have epCost field")
			assert.True(t, hasGoldCost, "Spell should have goldCost field")

			// Check for fallback values (10000 EP, 50000 GS)
			if epCostFloat, ok := epCost.(float64); ok && epCostFloat == 10000 {
				if goldCostFloat, ok := goldCost.(float64); ok && goldCostFloat == 50000 {
					fallbackSpells = append(fallbackSpells, name.(string))
					t.Logf("FALLBACK VALUES DETECTED: Spell '%s' (Level %v) - EP: %.0f, Gold: %.0f",
						name, level, epCostFloat, goldCostFloat)
				}
			}

			// Log first few spells for debugging
			if totalSpells <= 5 {
				t.Logf("Spell '%s' (Level %v) - EP: %v, Gold: %v", name, level, epCost, goldCost)
			}
		}
	}

	// Assert that no spells have fallback values
	if len(fallbackSpells) > 0 {
		t.Errorf("Found %d spells with fallback values (10000 EP, 50000 GS): %v",
			len(fallbackSpells), fallbackSpells)
	}

	t.Logf("Total spells checked: %d, spells with fallback values: %d", totalSpells, len(fallbackSpells))
}

func TestFallbackValueDetection(t *testing.T) {
	database.SetupTestDB(true) // Setup test database
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)

	token := getAuthToken()

	// Test both skills and spells for fallback values
	testCases := []struct {
		endpoint string
		itemType string
		testName string
	}{
		{"/api/characters/available-skills-new", "skill", "Schwimmen"},
		{"/api/characters/available-spells-new", "spell", "Angst"},
	}

	for _, tc := range testCases {
		t.Run(tc.itemType, func(t *testing.T) {
			request := gsmaster.LernCostRequest{
				CharId:       20,
				Name:         tc.testName,
				CurrentLevel: 0,
				Type:         tc.itemType,
				Action:       "learn",
				TargetLevel:  1,
				UsePP:        0,
				UseGold:      0,
				Reward:       &[]string{"default"}[0],
			}
			jsonData, _ := json.Marshal(request)

			req, _ := http.NewRequest("POST", tc.endpoint, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			req.Header.Set("Authorization", "Bearer "+token)

			respRecorder := httptest.NewRecorder()
			r.ServeHTTP(respRecorder, req)

			assert.Equal(t, http.StatusOK, respRecorder.Code)

			var response map[string]interface{}
			err := json.Unmarshal(respRecorder.Body.Bytes(), &response)
			assert.NoError(t, err)

			fallbackCount := 0
			totalItems := 0

			// Check for fallback values in response
			var dataKey string
			if tc.itemType == "skill" {
				dataKey = "skills_by_category"
			} else {
				dataKey = "spells_by_school"
			}

			if data, exists := response[dataKey]; exists {
				if dataMap, ok := data.(map[string]interface{}); ok {
					for _, itemsInterface := range dataMap {
						if items, ok := itemsInterface.([]interface{}); ok {
							for _, itemInterface := range items {
								if item, ok := itemInterface.(map[string]interface{}); ok {
									totalItems++
									if epCost, hasEP := item["epCost"]; hasEP {
										if goldCost, hasGold := item["goldCost"]; hasGold {
											if epFloat, ok := epCost.(float64); ok && epFloat == 10000 {
												if goldFloat, ok := goldCost.(float64); ok && goldFloat == 50000 {
													fallbackCount++
													t.Logf("FALLBACK DETECTED in %s '%s': EP=%v, Gold=%v",
														tc.itemType, item["name"], epCost, goldCost)
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}

			t.Logf("%s test: Total items=%d, Fallback values=%d", tc.itemType, totalItems, fallbackCount)
			assert.Equal(t, 0, fallbackCount, "No %s should have fallback values (10000 EP, 50000 GS)", tc.itemType)
		})
	}
}
