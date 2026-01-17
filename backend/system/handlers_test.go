package system

import (
	"bamort/config"
	"bamort/database"
	"bamort/deployment/migrations"
	"bamort/deployment/version"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// setupTestEnvironment sets up test environment variables
func setupTestEnvironment(t *testing.T) {
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original != "" {
			os.Setenv("ENVIRONMENT", original)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})
}

// setupTestDBWithVersion creates a test database and optionally sets a version
func setupTestDBWithVersion(t *testing.T, dbVersion string) *gorm.DB {
	setupTestEnvironment(t)
	database.SetupTestDB()
	db := database.DB

	// Create version tables
	err := db.AutoMigrate(&migrations.SchemaVersion{}, &migrations.MigrationHistory{})
	assert.NoError(t, err)

	// Clean any existing data
	db.Exec("DELETE FROM schema_version")
	db.Exec("DELETE FROM migration_history")

	// Insert version if provided
	if dbVersion != "" {
		versionRecord := map[string]interface{}{
			"version":          dbVersion,
			"migration_number": 1,
			"applied_at":       time.Now(),
			"backend_version":  config.GetVersion(),
			"description":      "Test version",
		}
		err = db.Table("schema_version").Create(versionRecord).Error
		assert.NoError(t, err)
	}

	return db
}

func TestHealthHandler_Compatible(t *testing.T) {
	// Setup: DB version matches required version
	db := setupTestDBWithVersion(t, version.RequiredDBVersion)

	// Create Gin context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterPublicRoutes(router, db)

	// Make request
	req, _ := http.NewRequest("GET", "/api/system/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, "ok", result["status"])
	assert.Equal(t, version.RequiredDBVersion, result["required_db_version"])
	assert.Equal(t, config.GetVersion(), result["actual_backend_version"])
	assert.Equal(t, version.RequiredDBVersion, result["db_version"])
	assert.Equal(t, false, result["migrations_pending"])
	assert.Equal(t, float64(0), result["pending_count"]) // JSON numbers are float64
	assert.Equal(t, true, result["compatible"])
	assert.NotNil(t, result["timestamp"])
}

func TestHealthHandler_MigrationPending(t *testing.T) {
	// Setup: Old DB version with migration number 0 (pre-migration system)
	oldVersion := "0.3.0"
	setupTestEnvironment(t)
	database.SetupTestDB()
	db := database.DB

	// Create version tables
	err := db.AutoMigrate(&migrations.SchemaVersion{}, &migrations.MigrationHistory{})
	assert.NoError(t, err)

	// Clean any existing data
	db.Exec("DELETE FROM schema_version")
	db.Exec("DELETE FROM migration_history")

	// Insert old version with migration_number 0 to simulate pending migrations
	versionRecord := map[string]interface{}{
		"version":          oldVersion,
		"migration_number": 0,
		"applied_at":       time.Now(),
		"backend_version":  "0.3.0",
		"description":      "Old version",
	}
	err = db.Table("schema_version").Create(versionRecord).Error
	assert.NoError(t, err)

	// Create Gin context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterPublicRoutes(router, db)

	// Make request
	req, _ := http.NewRequest("GET", "/api/system/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, "ok", result["status"])
	assert.Equal(t, oldVersion, result["db_version"])
	assert.Equal(t, true, result["migrations_pending"])
	assert.Greater(t, result["pending_count"], float64(0))
	assert.Equal(t, false, result["compatible"])
}

func TestHealthHandler_NoVersion(t *testing.T) {
	// Setup: DB without version (new installation)
	db := setupTestDBWithVersion(t, "")

	// Create Gin context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterPublicRoutes(router, db)

	// Make request
	req, _ := http.NewRequest("GET", "/api/system/health", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, "ok", result["status"])
	assert.Equal(t, "", result["db_version"])
	assert.Equal(t, true, result["migrations_pending"])
	assert.Equal(t, false, result["compatible"])
}

func TestVersionHandler_Success(t *testing.T) {
	// Setup
	db := setupTestDBWithVersion(t, version.RequiredDBVersion)

	// Create Gin context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterPublicRoutes(router, db)

	// Make request
	req, _ := http.NewRequest("GET", "/api/system/version", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Check backend section
	backend, ok := result["backend"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, config.GetVersion(), backend["version"])

	// Check database section
	database, ok := result["database"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, version.RequiredDBVersion, database["version"])
	assert.Equal(t, float64(1), database["migration_number"]) // JSON numbers are float64
	// last_migration can be nil or a valid time - both are acceptable
}

func TestVersionHandler_NoDBVersion(t *testing.T) {
	// Setup: DB without version
	db := setupTestDBWithVersion(t, "")

	// Create Gin context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterPublicRoutes(router, db)

	// Make request
	req, _ := http.NewRequest("GET", "/api/system/version", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	// Check database section
	database, ok := result["database"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "", database["version"])
	assert.Equal(t, float64(0), database["migration_number"])
	assert.Nil(t, database["last_migration"])
}
