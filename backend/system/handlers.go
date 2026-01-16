package system

import (
	"bamort/config"
	"bamort/deployment/migrations"
	"bamort/deployment/version"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status               string    `json:"status"`
	RequiredDBVersion    string    `json:"required_db_version"`
	ActualBackendVersion string    `json:"actual_backend_version"`
	DBVersion            string    `json:"db_version"`
	MigrationsPending    bool      `json:"migrations_pending"`
	PendingCount         int       `json:"pending_count"`
	Compatible           bool      `json:"compatible"`
	Timestamp            time.Time `json:"timestamp"`
}

// VersionResponse represents the version information response
type VersionResponse struct {
	Backend  BackendInfo  `json:"backend"`
	Database DatabaseInfo `json:"database"`
}

// BackendInfo contains backend version information
type BackendInfo struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
}

// DatabaseInfo contains database version information
type DatabaseInfo struct {
	Version         string     `json:"version"`
	MigrationNumber int        `json:"migration_number"`
	LastMigration   *time.Time `json:"last_migration"`
}

// HealthHandler handles GET /api/system/health
func HealthHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current DB version
		runner := migrations.NewMigrationRunner(db)
		dbVersion, _, err := runner.GetCurrentVersion()
		if err != nil {
			// Log error but continue - treat as no version
			dbVersion = ""
		}

		// Get pending migrations
		pending, err := runner.GetPendingMigrations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check pending migrations",
			})
			return
		}

		// Check compatibility
		compat := version.CheckCompatibility(dbVersion)

		response := HealthResponse{
			Status:               "ok",
			RequiredDBVersion:    version.RequiredDBVersion,
			ActualBackendVersion: config.GetVersion(),
			DBVersion:            dbVersion,
			MigrationsPending:    len(pending) > 0,
			PendingCount:         len(pending),
			Compatible:           compat.Compatible,
			Timestamp:            time.Now(),
		}

		c.JSON(http.StatusOK, response)
	}
}

// VersionHandler handles GET /api/system/version
func VersionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get backend version info
		backendInfo := BackendInfo{
			Version: config.GetVersion(),
			Commit:  config.GitCommit,
		}

		// Get database version info
		var dbInfo DatabaseInfo
		var versionRecord struct {
			Version         string
			MigrationNumber int
			AppliedAt       string
		}

		err := db.Raw(`
			SELECT version, migration_number, applied_at 
			FROM schema_version 
			ORDER BY id DESC 
			LIMIT 1
		`).Scan(&versionRecord).Error

		if err == nil && versionRecord.Version != "" {
			// Parse time if available
			var lastMigration *time.Time
			if versionRecord.AppliedAt != "" {
				if parsed, parseErr := time.Parse(time.RFC3339, versionRecord.AppliedAt); parseErr == nil {
					lastMigration = &parsed
				} else if parsed, parseErr := time.Parse("2006-01-02 15:04:05", versionRecord.AppliedAt); parseErr == nil {
					lastMigration = &parsed
				}
			}

			dbInfo = DatabaseInfo{
				Version:         versionRecord.Version,
				MigrationNumber: versionRecord.MigrationNumber,
				LastMigration:   lastMigration,
			}
		} else {
			// No version found - new installation or error
			dbInfo = DatabaseInfo{
				Version:         "",
				MigrationNumber: 0,
				LastMigration:   nil,
			}
		}

		response := VersionResponse{
			Backend:  backendInfo,
			Database: dbInfo,
		}

		c.JSON(http.StatusOK, response)
	}
}
