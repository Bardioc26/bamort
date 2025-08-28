package importer

import (
	"bamort/database"
	"bamort/models"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestImportCsv2Spell(t *testing.T) {
	// Clear source cache to ensure clean test state
	ClearSourceCache()

	// Setup test database
	database.SetupTestDB(true, false) // Use in-memory SQLite, no test data loading
	defer database.ResetTestDB()
	models.MigrateStructure()
	/*
		// Create test source data
		testSources := []models.Source{
			{Code: "ARK", Name: "Arkanum", GameSystem: "midgard"},
			{Code: "MYS", Name: "Mysterium", GameSystem: "midgard"},
			{Code: "KOD", Name: "Kodex", GameSystem: "midgard"},
		}
		for _, source := range testSources {
			source.Create()
		}
	*/
	t.Run("Import Zauber-Arkanum.csv", func(t *testing.T) {
		// Test if file exists
		csvPath := "/data/dev/bamort/backend/doc/Zauber-Arkanum.csv"
		if _, err := os.Stat(csvPath); os.IsNotExist(err) {
			t.Skipf("CSV file %s not found, skipping test", csvPath)
			return
		}

		// Test import
		err := ImportCsv2Spell(csvPath)
		assert.NoError(t, err, "Import should succeed")

		// Verify some spells were imported
		var count int64
		database.DB.Model(&models.Spell{}).Count(&count)
		assert.Greater(t, count, int64(0), "Should have imported some spells")

		// Test specific spell with more detailed debugging
		var spell models.Spell
		err = spell.First("Angst")
		if err != nil {
			// Let's check what spells actually exist
			var allSpells []models.Spell
			database.DB.Limit(10).Find(&allSpells)
			t.Logf("Found %d spells in database, first 10:", count)
			for i, s := range allSpells {
				t.Logf("Spell %d: %s (Category: %s, Level: %d)", i+1, s.Name, s.Category, s.Stufe)
			}
		}
		assert.NoError(t, err, "Should find 'Angst' spell")
		assert.Equal(t, "Beherrschen", spell.Category, "Spell category should match")
		assert.Equal(t, 2, spell.Stufe, "Spell level should be 2")

		// Test that source_id was set correctly (ARK should have been looked up)
		if spell.SourceID != 0 {
			var source models.Source
			err = database.DB.First(&source, spell.SourceID).Error
			assert.NoError(t, err, "Should find source by ID")
			assert.Equal(t, "ARK", source.Code, "Source code should be ARK")
		}
	})

	t.Run("Import Zauber-Mysterium.csv", func(t *testing.T) {
		// Test if file exists
		csvPath := "/data/dev/bamort/backend/doc/Zauber-Mysterium.csv"
		if _, err := os.Stat(csvPath); os.IsNotExist(err) {
			t.Skipf("CSV file %s not found, skipping test", csvPath)
			return
		}

		// Get current spell count
		var countBefore int64
		database.DB.Model(&models.Spell{}).Count(&countBefore)

		// Test import
		err := ImportCsv2Spell(csvPath)
		assert.NoError(t, err, "Import should succeed")

		// Verify more spells were imported
		var countAfter int64
		database.DB.Model(&models.Spell{}).Count(&countAfter)
		assert.Greater(t, countAfter, countBefore, "Should have imported additional spells")
	})

	t.Run("Test update existing spell", func(t *testing.T) {
		// Create a test spell first
		testSpell := models.Spell{
			GameSystem:   "midgard",
			Name:         "Test Zauber",
			Beschreibung: "Original description",
			Stufe:        1,
		}
		err := testSpell.Create()
		assert.NoError(t, err, "Should create test spell")

		// Create temporary CSV with updated data
		csvContent := `name,Beschreibung,stufe
Test Zauber,Updated description,2`

		tmpFile, err := os.CreateTemp("", "test_spell_*.csv")
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(csvContent)
		assert.NoError(t, err, "Should write to temp file")
		tmpFile.Close()

		// Import the updated data
		err = ImportCsv2Spell(tmpFile.Name())
		assert.NoError(t, err, "Import should succeed")

		// Verify the spell was updated
		var updatedSpell models.Spell
		err = updatedSpell.First("Test Zauber")
		assert.NoError(t, err, "Should find updated spell")
		assert.Equal(t, "Updated description", updatedSpell.Beschreibung, "Description should be updated")
		assert.Equal(t, 2, updatedSpell.Stufe, "Level should be updated")
	})

	t.Run("Test source lookup function", func(t *testing.T) {
		// Clear cache to ensure fresh lookups
		ClearSourceCache()

		// Test successful lookup of existing source
		sourceID, err := lookupSourceID("ARK")
		assert.NoError(t, err, "Should find ARK source")
		assert.Greater(t, sourceID, uint(0), "Source ID should be greater than 0")

		// Test auto-creation of non-existent source
		newSourceID, err := lookupSourceID("NEWCODE")
		assert.NoError(t, err, "Should auto-create NEWCODE source")
		assert.Greater(t, newSourceID, uint(0), "New source ID should be greater than 0")

		// Verify the source was actually created in the database
		var newSource models.Source
		err = database.DB.First(&newSource, newSourceID).Error
		assert.NoError(t, err, "Should find newly created source")
		assert.Equal(t, "NEWCODE", newSource.Code, "Source code should match")
		assert.Equal(t, "NEWCODE", newSource.Name, "Source name should default to code")
		assert.Equal(t, "midgard", newSource.GameSystem, "Game system should be midgard")
		assert.True(t, newSource.IsActive, "New source should be active")

		// Test that the second lookup uses cache (should return same ID)
		cachedSourceID, err := lookupSourceID("NEWCODE")
		assert.NoError(t, err, "Should find cached source")
		assert.Equal(t, newSourceID, cachedSourceID, "Cached lookup should return same ID")

		// Test empty source code
		_, err = lookupSourceID("")
		assert.Error(t, err, "Should return error for empty source code")
	})

	t.Run("Test auto-creation of sources during CSV import", func(t *testing.T) {
		// Clear cache and check initial source count
		ClearSourceCache()
		var initialSourceCount int64
		database.DB.Model(&models.Source{}).Count(&initialSourceCount)

		// Create temporary CSV with a new source code
		csvContent := `name,quelle,stufe,category
Test Spell,TESTSRC,1,Beherrschen`

		tmpFile, err := os.CreateTemp("", "test_auto_source_*.csv")
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(csvContent)
		assert.NoError(t, err, "Should write to temp file")
		tmpFile.Close()

		// Import the CSV
		err = ImportCsv2Spell(tmpFile.Name())
		assert.NoError(t, err, "Import should succeed")

		// Verify new source was created
		var finalSourceCount int64
		database.DB.Model(&models.Source{}).Count(&finalSourceCount)
		assert.Greater(t, finalSourceCount, initialSourceCount, "Should have created new source")

		// Verify the spell was imported with correct source
		var importedSpell models.Spell
		err = importedSpell.First("Test Spell")
		assert.NoError(t, err, "Should find imported spell")
		assert.Equal(t, "TESTSRC", importedSpell.Quelle, "Spell quelle should match")
		assert.Greater(t, importedSpell.SourceID, uint(0), "Spell should have source ID")

		// Verify the source details
		var createdSource models.Source
		err = database.DB.First(&createdSource, importedSpell.SourceID).Error
		assert.NoError(t, err, "Should find created source")
		assert.Equal(t, "TESTSRC", createdSource.Code, "Source code should match")
	})
}

func TestImportSpellCSVHandler(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, false)
	defer database.ResetTestDB()
	models.MigrateStructure()

	// Create test source
	testSource := models.Source{Code: "ARK", Name: "Arkanum", GameSystem: "midgard"}
	testSource.Create()

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("Successful CSV import via HTTP", func(t *testing.T) {
		// Create temporary test CSV file
		csvContent := `name,quelle,stufe,category
Test Spell HTTP,ARK,3,Beherrschen`

		tmpFile, err := os.CreateTemp("", "test_http_*.csv")
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(csvContent)
		assert.NoError(t, err, "Should write to temp file")
		tmpFile.Close()

		// Create test request
		router := gin.New()
		router.POST("/test", ImportSpellCSVHandler)

		req, err := http.NewRequest("POST", "/test?file="+tmpFile.Name(), nil)
		assert.NoError(t, err, "Should create request")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse JSON response")
		assert.True(t, response["success"].(bool), "Should indicate success")
		assert.Contains(t, response, "total_spells", "Should contain spell count")

		// Verify spell was imported
		var spell models.Spell
		err = spell.First("Test Spell HTTP")
		assert.NoError(t, err, "Should find imported spell")
		assert.Equal(t, 3, spell.Stufe, "Spell level should be 3")
	})

	t.Run("Missing file parameter", func(t *testing.T) {
		router := gin.New()
		router.POST("/test", ImportSpellCSVHandler)

		req, err := http.NewRequest("POST", "/test", nil)
		assert.NoError(t, err, "Should create request")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 Bad Request")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse JSON response")
		assert.Equal(t, "Missing file parameter", response["error"], "Should return correct error")
	})

	t.Run("File not found", func(t *testing.T) {
		router := gin.New()
		router.POST("/test", ImportSpellCSVHandler)

		req, err := http.NewRequest("POST", "/test?file=/nonexistent/file.csv", nil)
		assert.NoError(t, err, "Should create request")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 Bad Request")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse JSON response")
		assert.Equal(t, "File not found", response["error"], "Should return correct error")
	})

	t.Run("Invalid file type", func(t *testing.T) {
		// Create temporary non-CSV file
		tmpFile, err := os.CreateTemp("", "test_*.txt")
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile.Name())
		tmpFile.Close()

		router := gin.New()
		router.POST("/test", ImportSpellCSVHandler)

		req, err := http.NewRequest("POST", "/test?file="+tmpFile.Name(), nil)
		assert.NoError(t, err, "Should create request")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 Bad Request")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse JSON response")
		assert.Equal(t, "Invalid file type", response["error"], "Should return correct error")
	})

	t.Run("File upload with multipart form", func(t *testing.T) {
		// Setup test database
		database.SetupTestDB(true, false)
		models.MigrateStructure()

		// Create test CSV content
		csvContent := `name,beschreibung,quelle,stufe,ap
Test Spell Upload,Test description,ARK,1,2`

		// Create temporary file
		tmpFile, err := os.CreateTemp("", "test_spell_upload_*.csv")
		assert.NoError(t, err, "Should create temp file")
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(csvContent)
		assert.NoError(t, err, "Should write CSV content")
		tmpFile.Close()

		// Create multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add file field
		file, err := os.Open(tmpFile.Name())
		assert.NoError(t, err, "Should open temp file")
		defer file.Close()

		part, err := writer.CreateFormFile("file", "test_spells.csv")
		assert.NoError(t, err, "Should create form file")

		_, err = io.Copy(part, file)
		assert.NoError(t, err, "Should copy file content")

		err = writer.Close()
		assert.NoError(t, err, "Should close writer")

		// Create request
		router := gin.New()
		router.POST("/test", ImportSpellCSVHandler)

		req, err := http.NewRequest("POST", "/test", body)
		assert.NoError(t, err, "Should create request")
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Should parse JSON response")
		assert.True(t, response["success"].(bool), "Should be successful")
		assert.Contains(t, response["message"], "imported successfully", "Should contain success message")

		// Verify spell was imported
		var spell models.Spell
		err = database.DB.Where("name = ?", "Test Spell Upload").First(&spell).Error
		assert.NoError(t, err, "Should find imported spell")
		assert.Equal(t, "Test description", spell.Beschreibung, "Should have correct description")
	})
}
