package maintenance

import (
	"bamort/character"
	"bamort/database"
	"bamort/equipment"
	"bamort/models"
	"bamort/skills"
	"bamort/user"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// cleanupTestEnvironment creates a clean test environment
func cleanupTestEnvironment(t *testing.T) {
	// Clean up any existing test files
	localTestDataFile := filepath.Join(database.TestDataDir, "predefined_test_data.db")
	if err := os.RemoveAll(filepath.Dir(localTestDataFile)); err != nil {
		t.Logf("Warning: Could not clean test directory: %v", err)
	}

	// Reset any existing database connections
	database.ResetTestDB()
}

// createTestDataInLiveDB creates some test data in a live database for testing
func createTestDataInLiveDB(t *testing.T, liveDB *gorm.DB) {
	// Create test user
	testUser := &user.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	err := liveDB.Create(testUser).Error
	require.NoError(t, err)

	// Create test character
	testChar := &character.Char{
		BamortBase: models.BamortBase{
			Name: "Test Character",
		},
		Rasse: "Human",
		Typ:   "Warrior",
	}
	err = liveDB.Create(testChar).Error
	require.NoError(t, err)

	// Create test skill
	testSkill := &models.Skill{
		LookupList: models.LookupList{
			Name:         "Test Skill",
			Beschreibung: "A test skill",
		},
		Category:   "Combat",
		Difficulty: "1",
	}
	err = liveDB.Create(testSkill).Error
	require.NoError(t, err)

	// Create character skill
	testCharSkill := &skills.Fertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Test Skill",
			},
			CharacterID: testChar.ID,
		},
		Fertigkeitswert: 5,
	}
	err = liveDB.Create(testCharSkill).Error
	require.NoError(t, err)
}

// checks if copiing data from one DB to another works
// it uses 2 in Memory SQLITE databases
func TestMakeTestdataFromLiveRewrite(t *testing.T) {
	cleanupTestEnvironment(t)

	// Create a temporary live database with test data
	liveDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate structures to live DB
	err = database.MigrateStructure(liveDB)
	require.NoError(t, err)
	err = user.MigrateStructure(liveDB)
	require.NoError(t, err)
	err = character.MigrateStructure(liveDB)
	require.NoError(t, err)
	err = models.MigrateStructure(liveDB)
	require.NoError(t, err)
	err = skills.MigrateStructure(liveDB)
	require.NoError(t, err)
	err = equipment.MigrateStructure(liveDB)
	require.NoError(t, err)

	// Create test data in live DB
	createTestDataInLiveDB(t, liveDB)

	// We'll test the copyAllDataToTestDB function directly since we can't easily mock database.ConnectDatabase
	// Create a test database to copy to
	testDb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate structures to test database
	err = database.MigrateStructure(testDb)
	require.NoError(t, err)
	err = user.MigrateStructure(testDb)
	require.NoError(t, err)
	err = character.MigrateStructure(testDb)
	require.NoError(t, err)
	err = models.MigrateStructure(testDb)
	require.NoError(t, err)
	err = skills.MigrateStructure(testDb)
	require.NoError(t, err)
	err = equipment.MigrateStructure(testDb)
	require.NoError(t, err)

	// Test the copyAllDataToTestDB function
	copyStats, err := copyAllDataToTestDB(liveDB, testDb)
	require.NoError(t, err)

	// Verify statistics
	assert.Greater(t, copyStats["users"], 0, "Should have copied users")
	assert.Greater(t, copyStats["characters"], 0, "Should have copied characters")
	assert.Greater(t, copyStats["gsmaster_skills"], 0, "Should have copied skills")

	// Test saving to file
	localTestDataFile := filepath.Join(database.TestDataDir, "predefined_test_data.db")
	err = saveTestDatabaseToFile(testDb, localTestDataFile)
	require.NoError(t, err)

	// Check that the test data file was created
	_, err = os.Stat(localTestDataFile)
	assert.NoError(t, err, "Test data file should exist")

	// Verify the file contains data by loading it back
	fileDB, err := gorm.Open(sqlite.Open(localTestDataFile), &gorm.Config{})
	require.NoError(t, err)

	var userCount int64
	err = fileDB.Model(&user.User{}).Count(&userCount).Error
	require.NoError(t, err)
	assert.Greater(t, userCount, int64(0), "File should contain users")

	sqlDB, _ := fileDB.DB()
	sqlDB.Close()
}

func TestLoadPredefinedTestDataFromFile(t *testing.T) {
	cleanupTestEnvironment(t)

	// Create the directory if it doesn't exist
	tmpDBFile := filepath.Join(database.TestDataDir, "test_source.db")

	dir := filepath.Dir(tmpDBFile)
	var mErr error
	if err := os.MkdirAll(dir, 0755); err != nil {
		mErr = fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	require.NoError(t, mErr)

	// First, we need to create a test data file
	// Create a temporary source database with test data
	sourceDB, err := gorm.Open(sqlite.Open(tmpDBFile), &gorm.Config{})
	require.NoError(t, err)

	// Migrate structures to source DB
	err = database.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = user.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = character.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = models.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = skills.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = equipment.MigrateStructure(sourceDB)
	require.NoError(t, err)

	// Create test data in source DB
	createTestDataInLiveDB(t, sourceDB)

	// Save source DB as predefined test data file
	backupSQL := "VACUUM INTO '" + filepath.Join(database.TestDataDir, "predefined_test_data.db") + "'"

	err = sourceDB.Exec(backupSQL).Error
	require.NoError(t, err)

	// Close source DB
	sqlDB, _ := sourceDB.DB()
	sqlDB.Close()

	// Remove temporary source file
	defer os.Remove("testdata/test_source.db")

	// Now test LoadPredefinedTestDataFromFile
	// Create a new target database
	targetDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate structures to target DB
	err = database.MigrateStructure(targetDB)
	require.NoError(t, err)
	err = user.MigrateStructure(targetDB)
	require.NoError(t, err)
	err = character.MigrateStructure(targetDB)
	require.NoError(t, err)
	err = models.MigrateStructure(targetDB)
	require.NoError(t, err)
	err = skills.MigrateStructure(targetDB)
	require.NoError(t, err)
	err = equipment.MigrateStructure(targetDB)
	require.NoError(t, err)

	// Load test data from the source file
	err = LoadPredefinedTestDataFromFile(targetDB, "testdata/test_source.db")
	require.NoError(t, err)

	// Verify data was loaded
	var userCount int64
	err = targetDB.Model(&user.User{}).Count(&userCount).Error
	require.NoError(t, err)
	assert.Greater(t, userCount, int64(0), "Should have loaded users")

	var charCount int64
	err = targetDB.Model(&character.Char{}).Count(&charCount).Error
	require.NoError(t, err)
	assert.Greater(t, charCount, int64(0), "Should have loaded characters")

	var skillCount int64
	err = targetDB.Model(&models.Skill{}).Count(&skillCount).Error
	require.NoError(t, err)
	assert.Greater(t, skillCount, int64(0), "Should have loaded skills")

	// Verify specific data
	var loadedUser user.User
	err = targetDB.Where("username = ?", "testuser").First(&loadedUser).Error
	require.NoError(t, err)
	assert.Equal(t, "testuser", loadedUser.Username)
	assert.Equal(t, "test@example.com", loadedUser.Email)

	var loadedChar character.Char
	err = targetDB.Where("name = ?", "Test Character").First(&loadedChar).Error
	require.NoError(t, err)
	assert.Equal(t, "Test Character", loadedChar.Name)
	assert.Equal(t, "Human", loadedChar.Rasse)
}

func TestLoadPredefinedTestDataFromFile_FileNotFound(t *testing.T) {
	cleanupTestEnvironment(t)

	// Create target database
	targetDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Try to load non-existent test data
	err = LoadPredefinedTestDataFromFile(targetDB, "nonexistent_file.db")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "predefined test data file not found")
}

func TestSetupTestDBWithTestData(t *testing.T) {
	cleanupTestEnvironment(t)

	// Create directory if it doesn't exist
	tmpDBFile := filepath.Join(database.TestDataDir, "test_source.db")

	dir := filepath.Dir(tmpDBFile)
	var mErr error
	if err := os.MkdirAll(dir, 0755); err != nil {
		mErr = fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	require.NoError(t, mErr)

	// First create a predefined test data file
	sourceDB, err := gorm.Open(sqlite.Open("testdata/test_source.db"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate structures to live DB
	err = database.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = user.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = character.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = models.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = skills.MigrateStructure(sourceDB)
	require.NoError(t, err)
	err = equipment.MigrateStructure(sourceDB)
	require.NoError(t, err)

	createTestDataInLiveDB(t, sourceDB)

	// Save as predefined test data
	backupSQL := "VACUUM INTO '" + filepath.Join(database.TestDataDir, "predefined_test_data.db") + "'"
	err = sourceDB.Exec(backupSQL).Error
	require.NoError(t, err)

	sqlDB, _ := sourceDB.DB()
	sqlDB.Close()
	defer os.Remove("testdata/test_source.db")

	// Reset database state
	database.ResetTestDB()

	// Test SetupTestDB with test data loading
	database.SetupTestDB(true, true)

	// Verify that database.DB is available and has data
	require.NotNil(t, database.DB)

	// Check that data was loaded
	var userCount int64
	err = database.DB.Model(&user.User{}).Count(&userCount).Error
	require.NoError(t, err)
	assert.Greater(t, userCount, int64(0), "Should have loaded test users")

	// Clean up
	database.ResetTestDB()
}

func TestCopyDataFromFileToMemory(t *testing.T) {
	cleanupTestEnvironment(t)
	// Create the directory if it doesn't exist
	dir := filepath.Dir("testdata/test_source.db")
	if err := os.MkdirAll(dir, 0755); err != nil {
		require.NoError(t, err)
	}

	// Create source file database
	sourceDB, err := gorm.Open(sqlite.Open("testdata/test_source.db"), &gorm.Config{})
	require.NoError(t, err)

	//err = os.MkdirAll("testdata", 0755)
	//require.NoError(t, err)

	// Setup source database
	err = user.MigrateStructure(sourceDB)
	require.NoError(t, err)

	testUser := &user.User{
		Username:     "copytest",
		Email:        "copy@test.com",
		PasswordHash: "password",
	}
	err = sourceDB.Create(testUser).Error
	require.NoError(t, err)

	// Close source DB
	sqlDB, _ := sourceDB.DB()
	sqlDB.Close()

	// Create target in-memory database
	targetDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate target database
	err = user.MigrateStructure(targetDB)
	require.NoError(t, err)

	// Copy data
	err = copyDataFromFileToMemory("testdata/test_source.db", targetDB)
	require.NoError(t, err)

	// Verify data was copied
	var copiedUser user.User
	err = targetDB.Where("username = ?", "copytest").First(&copiedUser).Error
	require.NoError(t, err)
	assert.Equal(t, "copytest", copiedUser.Username)
	assert.Equal(t, "copy@test.com", copiedUser.Email)

	// Clean up
	defer os.Remove("testdata/test_source.db")
}

// Cleanup function to run after tests
func TestMain(m *testing.M) {
	// Setup
	gin.SetMode(gin.TestMode)

	// Run tests
	code := m.Run()

	// Cleanup
	os.RemoveAll("testdata")

	os.Exit(code)
}
