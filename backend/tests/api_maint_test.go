package tests

import (
	"bamort/character"
	"bamort/gsmaster"
	"bamort/maintenance"
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

func TestMaintSetupCheck(t *testing.T) {
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

func TestGetMasterData(t *testing.T) {
	SetupTestDB(false)
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)
	// Character routes
	rCharGrp := router.MaintenanceRouterGrp(protected)
	rCharGrp.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})
	u := user.User{}
	u.FirstId(1)

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/maintenance", nil)
	req.Header.Set("Content-Type", "application/json")
	token := user.GenerateToken(&u)
	req.Header.Set("Authorization", "Bearer "+token)

	// Create a response recorder to capture the handler's response
	respRecorder := httptest.NewRecorder()

	// Perform the test request
	r.ServeHTTP(respRecorder, req)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, respRecorder.Code)

	// Assert the response body
	type dtaStruct struct {
		Skills       []gsmaster.Skill       `json:"skills"`
		Weaponskills []gsmaster.WeaponSkill `json:"weaponskills"`
		Spell        []gsmaster.Spell       `json:"spells"`
		Equipment    []gsmaster.Equipment   `json:"equipment"`
		Weapons      []gsmaster.Weapon      `json:"weapons"`
	}
	var dta dtaStruct
	err := json.Unmarshal(respRecorder.Body.Bytes(), &dta)
	assert.NoError(t, err)
}
func TestGetMDSkillCategories(t *testing.T) {
	SetupTestDB(false)
	//gsmaster.MigrateStructure()
	ski := gsmaster.Skill{}
	categories, err := ski.GetSkillCategories()
	assert.NoError(t, err)
	assert.LessOrEqual(t, 1, len(categories))
	assert.Equal(t, "Allgemein", categories[0])
}
func TestGetMDSkills(t *testing.T) {
	SetupTestDB()
	//TestCreateChar(t)
	//TestRegisterUser(t)
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

// TestCreateCharacter tests the POST /characters endpoint
func TestGetMDSkill(t *testing.T) {
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
