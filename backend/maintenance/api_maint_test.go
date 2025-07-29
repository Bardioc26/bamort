package maintenance

import (
	"bamort/character"
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"bamort/router"
	"bamort/user"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMaintSetupCheck(t *testing.T) {
	// Setup proper test database
	database.SetupTestDB(true)

	// Create a proper HTTP test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	SetupCheck(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Setup Check OK", response["message"])
}

func TestGetMasterData(t *testing.T) {
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
		Skills       []models.Skill       `json:"skills"`
		Weaponskills []models.WeaponSkill `json:"weaponskills"`
		Spell        []models.Spell       `json:"spells"`
		Equipment    []models.Equipment   `json:"equipment"`
		Weapons      []models.Weapon      `json:"weapons"`
	}
	var dta dtaStruct
	err := json.Unmarshal(respRecorder.Body.Bytes(), &dta)
	assert.NoError(t, err)
}

func TestGetMDSkillCategories(t *testing.T) {
	database.SetupTestDB(false)

	ski := models.Skill{}
	categories, err := ski.GetSkillCategories()
	assert.NoError(t, err)
	assert.LessOrEqual(t, 1, len(categories))
	assert.Equal(t, "Wissen", categories[0])
}

func TestGetMDSkills(t *testing.T) {
	database.SetupTestDB(false)
	// Initialize a Gin router
	r := gin.Default()
	router.SetupGin(r)

	// Routes
	protected := router.BaseRouterGrp(r)

	// Character routes
	character.RegisterRoutes(protected)
	gsmaster.RegisterRoutes(protected)
	protected.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "Test OK"})
	})
	u := user.User{}
	u.FirstId(1)

	// Create a test HTTP request
	req, _ := http.NewRequest("GET", "/api/characters/20", nil)
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
	assert.Equal(t, 20, int(listOfCharacter.ID)) // Check the simulated ID
	assert.Equal(t, "Krieger", listOfCharacter.Typ)
	assert.Equal(t, 3, listOfCharacter.Grad)
	//assert.Equal(t, "test", listOfCharacter.Owner)
	//assert.Equal(t, false, listOfCharacter.Public)

}
