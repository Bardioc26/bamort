package install

import (
	"bamort/config"
	"bamort/deployment/masterdata"
	"bamort/deployment/migrations"
	"bamort/logger"
	"bamort/models"
	"bamort/user"
	"crypto/md5"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// NewInstallation handles fresh database installation
type NewInstallation struct {
	DB              *gorm.DB
	MasterDataPath  string
	CreateAdminUser bool
	AdminUsername   string
	AdminPassword   string
	GameSystem      string
}

// InstallationResult contains the result of the installation
type InstallationResult struct {
	Success        bool
	Version        string
	TablesCreated  int
	AdminCreated   bool
	MasterDataOK   bool
	ExecutionTime  time.Duration
	Errors         []string
}

// NewInstaller creates a new installation instance
func NewInstaller(db *gorm.DB) *NewInstallation {
	return &NewInstallation{
		DB:              db,
		MasterDataPath:  "./masterdata",
		CreateAdminUser: false,
		GameSystem:      "midgard",
	}
}

// Initialize performs a fresh installation
func (n *NewInstallation) Initialize() (*InstallationResult, error) {
	startTime := time.Now()
	result := &InstallationResult{
		Version: config.GetVersion(),
	}

	logger.Info("Initializing new Bamort installation...")
	logger.Info("Backend version: %s", result.Version)

	// Step 1: Create database schema using GORM
	logger.Info("Step 1/4: Creating database schema...")
	if err := n.createDatabaseSchema(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Schema creation failed: %v", err))
		return result, fmt.Errorf("schema creation failed: %w", err)
	}
	logger.Info("✓ Database schema created successfully")

	// Step 2: Initialize version tracking
	logger.Info("Step 2/4: Initializing version tracking...")
	if err := n.initializeVersionTracking(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Version tracking failed: %v", err))
		return result, fmt.Errorf("version tracking failed: %w", err)
	}
	logger.Info("✓ Version tracking initialized (DB version: %s)", config.GetVersion())

	// Step 3: Import master data
	logger.Info("Step 3/4: Importing master data from %s...", n.MasterDataPath)
	if err := n.importMasterData(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Master data import failed: %v", err))
		return result, fmt.Errorf("master data import failed: %w", err)
	}
	result.MasterDataOK = true
	logger.Info("✓ Master data imported successfully")

	// Step 4: Create admin user if requested
	if n.CreateAdminUser {
		logger.Info("Step 4/4: Creating admin user '%s'...", n.AdminUsername)
		if err := n.createAdmin(); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Admin creation failed: %v", err))
			return result, fmt.Errorf("admin creation failed: %w", err)
		}
		result.AdminCreated = true
		logger.Info("✓ Admin user created successfully")
	} else {
		logger.Info("Step 4/4: Skipping admin user creation (not requested)")
	}

	result.Success = true
	result.ExecutionTime = time.Since(startTime)

	logger.Info("═══════════════════════════════════════════════════")
	logger.Info("Installation completed successfully!")
	logger.Info("Version: %s", result.Version)
	logger.Info("Execution time: %v", result.ExecutionTime)
	logger.Info("═══════════════════════════════════════════════════")

	return result, nil
}

// createDatabaseSchema creates all tables using GORM AutoMigrate
func (n *NewInstallation) createDatabaseSchema() error {
	logger.Debug("Running GORM AutoMigrate for all models...")

	if err := models.MigrateStructure(n.DB); err != nil {
		return fmt.Errorf("GORM AutoMigrate failed: %w", err)
	}

	logger.Debug("All tables created successfully")
	return nil
}

// initializeVersionTracking creates version tables and records initial version
func (n *NewInstallation) initializeVersionTracking() error {
	// Get the first migration (creates version tables)
	if len(migrations.AllMigrations) == 0 {
		return fmt.Errorf("no migrations available")
	}

	firstMigration := migrations.AllMigrations[0]
	logger.Debug("Applying initial migration: %s", firstMigration.Description)

	// Create tables using the migration's DataFunc (GORM-based)
	if firstMigration.DataFunc != nil {
		if err := firstMigration.DataFunc(n.DB); err != nil {
			return fmt.Errorf("failed to create version tables: %w", err)
		}
	} else {
		// Fallback: execute SQL if no DataFunc
		for _, sql := range firstMigration.UpSQL {
			if err := n.DB.Exec(sql).Error; err != nil {
				return fmt.Errorf("failed to execute SQL: %w", err)
			}
		}
	}

	// Record initial version (all migrations are considered "pre-applied")
	latestMigration := migrations.GetLatestMigration()
	if latestMigration == nil {
		return fmt.Errorf("no migrations available")
	}

	version := map[string]interface{}{
		"version":          latestMigration.Version,
		"migration_number": latestMigration.Number,
		"applied_at":       time.Now(),
		"backend_version":  config.GetVersion(),
		"description":      "Initial installation",
	}

	if err := n.DB.Table("schema_version").Create(version).Error; err != nil {
		return fmt.Errorf("failed to record version: %w", err)
	}

	// Record migration history for all migrations (as pre-applied)
	for _, m := range migrations.AllMigrations {
		history := map[string]interface{}{
			"migration_number":   m.Number,
			"version":            m.Version,
			"description":        m.Description,
			"applied_at":         time.Now(),
			"applied_by":         "installer",
			"execution_time_ms":  0,
			"success":            true,
			"rollback_available": len(m.DownSQL) > 0,
		}

		if err := n.DB.Table("migration_history").Create(history).Error; err != nil {
			return fmt.Errorf("failed to record migration history: %w", err)
		}
	}

	logger.Debug("Version tracking initialized with version %s (migration %d)",
		latestMigration.Version, latestMigration.Number)

	return nil
}

// importMasterData imports all master data using MasterDataSync
func (n *NewInstallation) importMasterData() error {
	sync := masterdata.NewMasterDataSync(n.DB, n.MasterDataPath)
	sync.Verbose = true

	if err := sync.SyncAll(); err != nil {
		return err
	}

	return nil
}

// createAdmin creates the admin user
func (n *NewInstallation) createAdmin() error {
	if n.AdminUsername == "" {
		return fmt.Errorf("admin username not specified")
	}

	if n.AdminPassword == "" {
		return fmt.Errorf("admin password not specified")
	}

	// Check if user already exists
	var existing user.User
	if err := n.DB.Where("username = ?", n.AdminUsername).First(&existing).Error; err == nil {
		logger.Warn("Admin user '%s' already exists, skipping creation", n.AdminUsername)
		return nil
	}

	// Create new admin user with MD5 password hash (matching user/handlers.go)
	admin := &user.User{
		Username: n.AdminUsername,
		Email:    n.AdminUsername + "@localhost",
		Role:     user.RoleAdmin,
	}

	// Hash password using MD5 (same as Register handler)
	hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(n.AdminPassword)))
	admin.PasswordHash = hashedPassword

	if err := n.DB.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	logger.Debug("Admin user '%s' created with ID %d", admin.Username, admin.UserID)
	return nil
}
