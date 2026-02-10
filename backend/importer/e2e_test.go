package importer

import (
	"bamort/database"
	"bamort/models"
	"bamort/user"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestE2E_CompleteImportWorkflow tests the full user workflow from file upload to database
func TestE2E_CompleteImportWorkflow(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB(true)
	db := database.DB

	// Clean up before test to ensure fresh state
	db.Exec("DELETE FROM import_histories")
	db.Exec("DELETE FROM master_data_imports")
	db.Exec("DELETE FROM chars")
	db.Exec("DELETE FROM users")

	// Run migrations
	err := models.MigrateStructure(db)
	require.NoError(t, err)
	err = MigrateStructure(db)
	require.NoError(t, err)

	// Create test user
	testUser := &user.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}
	err = db.Create(testUser).Error
	require.NoError(t, err)

	// Mock adapter server
	mockAdapter := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/metadata":
			metadata := map[string]interface{}{
				"id":                   "test-adapter-v1",
				"name":                 "Test Adapter",
				"version":              "1.0",
				"bmrt_versions":        []string{"1.0"},
				"supported_extensions": []string{".json"},
				"capabilities":         []string{"import", "export", "detect"},
			}
			json.NewEncoder(w).Encode(metadata)
		case "/detect":
			response := map[string]interface{}{
				"confidence": 0.95,
				"version":    "1.0",
			}
			json.NewEncoder(w).Encode(response)
		case "/import":
			// Return minimal BMRT format
			bmrt := CharacterImport{
				Name:  "E2E Test Character",
				Grad:  1,
				Rasse: "Mensch",
				Typ:   "Krieger",
				Alter: 25,
				Eigenschaften: Eigenschaften{
					St: 80, Gs: 75, Gw: 70, Ko: 85,
					In: 65, Zt: 60, Pa: 55, Au: 70, Wk: 60,
				},
				Lp: Lp{Max: 12, Value: 12},
				Ap: Ap{Max: 20, Value: 20},
			}
			json.NewEncoder(w).Encode(bmrt)
		case "/export":
			// Return simple JSON
			response := map[string]string{"name": "E2E Test Character"}
			json.NewEncoder(w).Encode(response)
		default:
			http.NotFound(w, r)
		}
	}))
	defer mockAdapter.Close()

	// Register mock adapter
	registry := NewAdapterRegistry()
	err = registry.Register(AdapterMetadata{
		ID:                  "test-adapter-v1",
		Name:                "Test Adapter",
		Version:             "1.0",
		BmrtVersions:        []string{"1.0"},
		SupportedExtensions: []string{".json"},
		BaseURL:             mockAdapter.URL,
		Capabilities:        []string{"import", "export", "detect"},
		Healthy:             true,
		LastCheckedAt:       time.Now(),
	})
	require.NoError(t, err)

	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", testUser.UserID)
		c.Next()
	})

	// Register routes with our mock registry
	registerRoutesWithRegistry(protected, registry)

	// STEP 1: Upload file and detect format
	t.Run("DetectFormat", func(t *testing.T) {
		body, contentType := createMultipartFile(t, "file", "test.json", []byte(`{"test": "data"}`))

		req := httptest.NewRequest("POST", "/api/import/detect", body)
		req.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "test-adapter-v1", response["adapter_id"])
		assert.Greater(t, response["confidence"], 0.9)
	})

	// STEP 2: Import character
	var importID uint
	var characterID uint
	t.Run("ImportCharacter", func(t *testing.T) {
		body, contentType := createMultipartFile(t, "file", "test.json", []byte(`{"test": "data"}`))

		req := httptest.NewRequest("POST", "/api/import/import", body)
		req.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response ImportResult
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotZero(t, response.CharacterID)
		assert.NotZero(t, response.ImportID)
		assert.Equal(t, "test-adapter-v1", response.AdapterID)
		assert.Equal(t, "success", response.Status)

		importID = response.ImportID
		characterID = response.CharacterID
	})

	// STEP 3: Verify character in database
	t.Run("VerifyCharacterInDB", func(t *testing.T) {
		var char models.Char
		err := db.Where("id = ?", characterID).First(&char).Error
		require.NoError(t, err)

		assert.Equal(t, "E2E Test Character", char.Name)
		assert.Equal(t, 1, char.Grad)
		assert.Equal(t, testUser.UserID, char.UserID)
		assert.NotNil(t, char.ImportedFromAdapter)
		assert.Equal(t, "test-adapter-v1", *char.ImportedFromAdapter)
		assert.NotNil(t, char.ImportedAt)
	})

	// STEP 4: Verify import history
	t.Run("VerifyImportHistory", func(t *testing.T) {
		var history ImportHistory
		err := db.Where("id = ?", importID).First(&history).Error
		require.NoError(t, err)

		assert.Equal(t, testUser.UserID, history.UserID)
		require.NotNil(t, history.CharacterID)
		assert.Equal(t, characterID, *history.CharacterID)
		assert.Equal(t, "test-adapter-v1", history.AdapterID)
		assert.Equal(t, "success", history.Status)
		assert.NotEmpty(t, history.SourceSnapshot)
		assert.Equal(t, "1.0", history.BmrtVersion)
	})

	// STEP 5: List import history
	t.Run("ListImportHistory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/import/history", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Histories []ImportHistory `json:"histories"`
			Total     int64           `json:"total"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, int64(1), response.Total)
		assert.Len(t, response.Histories, 1)
		assert.Equal(t, importID, response.Histories[0].ID)
	})

	// STEP 6: Export character back to original format
	t.Run("ExportCharacter", func(t *testing.T) {
		req := httptest.NewRequest("POST", fmt.Sprintf("/api/import/export/%d", characterID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")

		var exported map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &exported)
		require.NoError(t, err)
		assert.Equal(t, "E2E Test Character", exported["name"])
	})

	// STEP 7: List available adapters
	t.Run("ListAdapters", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/import/adapters", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Adapters []AdapterMetadata `json:"adapters"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Len(t, response.Adapters, 1)
		assert.Equal(t, "test-adapter-v1", response.Adapters[0].ID)
	})

	t.Cleanup(func() {
		db.Exec("DELETE FROM import_histories")
		db.Exec("DELETE FROM master_data_imports")
		db.Exec("DELETE FROM chars")
		db.Exec("DELETE FROM users")
	})
}

// TestE2E_ImportWithMasterDataReconciliation tests master data creation during import
func TestE2E_ImportWithMasterDataReconciliation(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB(true)
	db := database.DB

	err := models.MigrateStructure(db)
	require.NoError(t, err)
	err = MigrateStructure(db)
	require.NoError(t, err)

	// Create Midgard5 game system for testing (required for master data reconciliation)
	gameSystem := &models.GameSystem{
		Code:        "Midgard5",
		Name:        "Midgard 5",
		Description: "Midgard 5th Edition",
	}
	err = db.Create(gameSystem).Error
	require.NoError(t, err)

	testUser := &user.User{
		Username:     "testuser2",
		Email:        "test2@example.com",
		PasswordHash: "hashedpassword",
	}
	err = db.Create(testUser).Error
	require.NoError(t, err)

	// Mock adapter with skills and spells
	mockAdapter := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/metadata":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                   "test-adapter-v2",
				"bmrt_versions":        []string{"1.0"},
				"supported_extensions": []string{".json"},
				"capabilities":         []string{"import", "detect"},
			})
		case "/detect":
			json.NewEncoder(w).Encode(map[string]interface{}{"confidence": 0.95})
		case "/import":
			bmrt := CharacterImport{
				Name:          "Character with Skills",
				Grad:          2,
				Rasse:         "Elf",
				Typ:           "Magier",
				Alter:         100,
				Eigenschaften: Eigenschaften{St: 80, Gs: 75, Gw: 70, Ko: 85, In: 65, Zt: 60, Pa: 55, Au: 70, Wk: 60},
				Lp:            Lp{Max: 10, Value: 10},
				Ap:            Ap{Max: 30, Value: 30},
				Fertigkeiten: []Fertigkeit{
					{ImportBase: ImportBase{Name: "Custom Skill 1"}, Fertigkeitswert: 12},
					{ImportBase: ImportBase{Name: "Custom Skill 2"}, Fertigkeitswert: 15},
				},
				Zauber: []Zauber{
					{ImportBase: ImportBase{Name: "Custom Spell"}, Bonus: 10},
				},
			}
			json.NewEncoder(w).Encode(bmrt)
		}
	}))
	defer mockAdapter.Close()

	registry := NewAdapterRegistry()
	err = registry.Register(AdapterMetadata{
		ID:                  "test-adapter-v2",
		BaseURL:             mockAdapter.URL,
		BmrtVersions:        []string{"1.0"},
		SupportedExtensions: []string{".json"},
		Capabilities:        []string{"import", "detect"},
		Healthy:             true,
	})
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", testUser.UserID)
		c.Next()
	})
	registerRoutesWithRegistry(protected, registry)

	// Import character
	body, contentType := createMultipartFile(t, "file", "test.json", []byte(`{}`))
	req := httptest.NewRequest("POST", "/api/import/import", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Logf("Error response: %s", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response ImportResult
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verify master data imports were logged
	var masterDataCount int64
	db.Model(&MasterDataImport{}).Where("import_history_id = ?", response.ImportID).Count(&masterDataCount)
	assert.Greater(t, masterDataCount, int64(0), "Should have logged master data imports")

	t.Cleanup(func() {
		db.Exec("DELETE FROM import_histories")
		db.Exec("DELETE FROM master_data_imports")
		db.Exec("DELETE FROM chars")
		db.Exec("DELETE FROM users")
	})
}

// TestE2E_ImportFailureRollback tests that failed imports rollback correctly
func TestE2E_ImportFailureRollback(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB(true)
	db := database.DB

	err := models.MigrateStructure(db)
	require.NoError(t, err)
	err = MigrateStructure(db)
	require.NoError(t, err)

	testUser := &user.User{
		Username:     "testuser3",
		Email:        "test3@example.com",
		PasswordHash: "hashedpassword",
	}
	err = db.Create(testUser).Error
	require.NoError(t, err)

	// Mock adapter that returns invalid data
	mockAdapter := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/metadata":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":                   "test-adapter-v3",
				"bmrt_versions":        []string{"1.0"},
				"supported_extensions": []string{".json"},
				"capabilities":         []string{"import"},
			})
		case "/import":
			// Return invalid BMRT (missing required fields)
			bmrt := CharacterImport{
				Name: "", // Invalid: empty name
				Grad: 0,
			}
			json.NewEncoder(w).Encode(bmrt)
		}
	}))
	defer mockAdapter.Close()

	registry := NewAdapterRegistry()
	err = registry.Register(AdapterMetadata{
		ID:           "test-adapter-v3",
		BaseURL:      mockAdapter.URL,
		BmrtVersions: []string{"1.0"},
		Capabilities: []string{"import"},
		Healthy:      true,
	})
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", testUser.UserID)
		c.Next()
	})
	registerRoutesWithRegistry(protected, registry)

	// Attempt import
	body, contentType := createMultipartFile(t, "file", "test.json", []byte(`{}`))
	req := httptest.NewRequest("POST", "/api/import/import", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should fail with validation error
	assert.NotEqual(t, http.StatusOK, w.Code)

	// Verify no character was created
	var charCount int64
	db.Model(&models.Char{}).Where("user_id = ?", testUser.UserID).Count(&charCount)
	assert.Equal(t, int64(0), charCount, "No character should be created on failed import")

	t.Cleanup(func() {
		db.Exec("DELETE FROM import_histories")
		db.Exec("DELETE FROM chars")
		db.Exec("DELETE FROM users")
	})
}

// TestE2E_UnhealthyAdapterHandling tests graceful handling of unavailable adapters
func TestE2E_UnhealthyAdapterHandling(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB(true)
	db := database.DB

	err := models.MigrateStructure(db)
	require.NoError(t, err)
	err = MigrateStructure(db)
	require.NoError(t, err)

	registry := NewAdapterRegistry()
	err = registry.Register(AdapterMetadata{
		ID:           "unhealthy-adapter",
		BaseURL:      "http://localhost:9999", // Non-existent
		BmrtVersions: []string{"1.0"},
		Capabilities: []string{"import"},
		Healthy:      false, // Mark as unhealthy
	})
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	protected := router.Group("/api")
	protected.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})
	registerRoutesWithRegistry(protected, registry)

	// Attempt to use unhealthy adapter
	body, contentType := createMultipartFile(t, "file", "test.json", []byte(`{}`))
	req := httptest.NewRequest("POST", "/api/import/import?adapter_id=unhealthy-adapter", body)
	req.Header.Set("Content-Type", contentType)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "unavailable")
}

// Helper function to create multipart file
func createMultipartFile(t *testing.T, fieldName, filename string, content []byte) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	require.NoError(t, err)

	_, err = io.Copy(part, bytes.NewReader(content))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	return body, writer.FormDataContentType()
}

// Helper function to register routes with custom registry
func registerRoutesWithRegistry(r *gin.RouterGroup, registry *AdapterRegistry) {
	globalRegistry = registry // Set global registry
	RegisterRoutes(r)
}

// TestE2E_RoundTripExportImport tests that exported data can be reimported
func TestE2E_RoundTripExportImport(t *testing.T) {
	t.Skip("Requires full adapter implementation with export support")
	// This test would:
	// 1. Import a character
	// 2. Export it to original format
	// 3. Import the exported file again
	// 4. Verify both characters are identical
}

// TestE2E_ConcurrentImports tests multiple simultaneous imports
func TestE2E_ConcurrentImports(t *testing.T) {
	t.Skip("Stress test - run separately")
	// This test would simulate multiple users importing at the same time
	// to verify thread safety and database transaction handling
}

// TestE2E_LargeFileImport tests import of large character files
func TestE2E_LargeFileImport(t *testing.T) {
	t.Skip("Performance test - run separately")
	// This test would import a character with hundreds of skills/spells
	// to verify performance and memory handling
}
