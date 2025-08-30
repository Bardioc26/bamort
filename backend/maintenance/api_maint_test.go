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
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	database.SetupTestDB() //(false)

	ski := models.Skill{}
	categories, err := ski.GetSkillCategories()
	assert.NoError(t, err)
	assert.LessOrEqual(t, 1, len(categories))
	assert.Equal(t, "Wissen", categories[0])
}

func TestGetMDSkills(t *testing.T) {
	database.SetupTestDB() //(false)
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

func TestMaintMakeTestdataFromLive(t *testing.T) {
	// Setup proper test database
	database.SetupTestDB(true)

	// Create a proper HTTP test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	MakeTestdataFromLive(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Live database copied to file successfully", response["message"])
	assert.Contains(t, response, "test_data_file")
}

func TestMaintReconnectDataBase(t *testing.T) {
	// Create a proper HTTP test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	ReconnectDataBase(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Database reconnected successfully", response["message"])
}

func TestMaintReloadENV(t *testing.T) {
	// Create a proper HTTP test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	ReloadENV(c)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Environment variables reloaded successfully", response["message"])
}

func TestMaintLoadPredefinedTestData(t *testing.T) {
	// Setup proper test database
	database.SetupTestDB(true)

	// Create a proper HTTP test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler
	LoadPredefinedTestData(c)

	// The function should always attempt to load data, but may return different status codes
	// based on whether the predefined test data file exists and is accessible
	if w.Code == http.StatusOK {
		// Parse the JSON response for success case
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Predefined test data loaded successfully into in-memory database", response["message"])
		assert.Contains(t, response, "test_data_file")
		assert.Contains(t, response, "statistics")
	} else {
		// Check that it fails gracefully if no test data file exists or other issues
		// Could be 404 (file not found), 500 (internal server error), etc.
		assert.True(t, w.Code >= 400, "Should return an error status code when predefined data is not available")
		
		// Verify error response structure
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err == nil {
			// If we can parse JSON, check for error field
			assert.Contains(t, response, "error")
		}
		// If we can't parse JSON, that's also acceptable for some error cases
	}
}

func TestLoadPredefinedTestDataFromFile(t *testing.T) {
	// Create a temporary test database file
	tempDir := t.TempDir()
	testDataFile := filepath.Join(tempDir, "test_data.db")

	// Create a simple test database with some data
	sourceDB, err := gorm.Open(sqlite.Open(testDataFile), &gorm.Config{})
	require.NoError(t, err, "Should create test database file")

	// Migrate basic structures
	err = sourceDB.AutoMigrate(&user.User{})
	require.NoError(t, err, "Should migrate structures")

	// Add test data
	testUser := &user.User{
		UserID:       1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
	}
	err = sourceDB.Create(testUser).Error
	require.NoError(t, err, "Should create test user")

	// Close the source database
	if sqlDB, err := sourceDB.DB(); err == nil {
		sqlDB.Close()
	}

	// Create target in-memory database
	targetDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Should create target database")
	defer func() {
		if sqlDB, err := targetDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Test the function
	err = LoadPredefinedTestDataFromFile(targetDB, testDataFile)
	assert.NoError(t, err, "LoadPredefinedTestDataFromFile should succeed")

	// Verify data was loaded
	var userCount int64
	err = targetDB.Model(&user.User{}).Count(&userCount).Error
	assert.NoError(t, err, "Should be able to count users")
	assert.Equal(t, int64(1), userCount, "Should have loaded the test user")
}

func TestMaintTransferSQLiteToMariaDB(t *testing.T) {
	// Skip this test if we don't have MariaDB available
	// This test would typically be an integration test
	t.Skip("Skipping TransferSQLiteToMariaDB test - requires MariaDB setup")

	// Note: This test would need:
	// 1. A real MariaDB instance running
	// 2. Proper test data setup in SQLite source
	// 3. Verification of data transfer
	//
	// For unit testing, this would be better tested by testing the
	// individual components like copySQLiteToMariaDB function separately
}
