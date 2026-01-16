package deployment

import (
	"bamort/config"
	"bamort/deployment/backup"
	"bamort/deployment/migrations"
	"bamort/deployment/version"
	"bamort/logger"
	"bamort/transfer"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DeploymentOrchestrator coordinates the full deployment process
type DeploymentOrchestrator struct {
	DB *gorm.DB
}

// DeploymentReport contains the results of a deployment
type DeploymentReport struct {
	Success          bool
	StartTime        time.Time
	EndTime          time.Time
	Duration         time.Duration
	BackupCreated    bool
	BackupPath       string
	MigrationsRun    int
	ValidationPassed bool
	Errors           []string
	Warnings         []string
}

// NewOrchestrator creates a new deployment orchestrator
func NewOrchestrator(db *gorm.DB) *DeploymentOrchestrator {
	return &DeploymentOrchestrator{
		DB: db,
	}
}

// Deploy executes the full deployment workflow
func (o *DeploymentOrchestrator) Deploy() (*DeploymentReport, error) {
	report := &DeploymentReport{
		StartTime: time.Now(),
	}

	logger.Info("═══════════════════════════════════════════════════")
	logger.Info("Starting Deployment Process")
	logger.Info("═══════════════════════════════════════════════════")

	// Step 1: Create backup
	logger.Info("Step 1/4: Creating pre-deployment backup...")
	backupPath, err := o.createBackup()
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Backup failed: %v", err))
		return report, fmt.Errorf("backup failed: %w", err)
	}
	report.BackupCreated = true
	report.BackupPath = backupPath
	logger.Info("✓ Backup created: %s", backupPath)

	// Step 2: Check version compatibility
	logger.Info("Step 2/4: Checking version compatibility...")
	if err := o.checkCompatibility(); err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Compatibility check failed: %v", err))
		return report, fmt.Errorf("compatibility check failed: %w", err)
	}
	logger.Info("✓ Version compatibility verified")

	// Step 3: Apply migrations
	logger.Info("Step 3/4: Applying database migrations...")
	migrationsRun, err := o.applyMigrations()
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Migration failed: %v", err))
		logger.Error("Migration failed, attempting rollback...")
		// Rollback would happen here in production
		return report, fmt.Errorf("migration failed: %w", err)
	}
	report.MigrationsRun = migrationsRun
	if migrationsRun > 0 {
		logger.Info("✓ Applied %d migration(s)", migrationsRun)
	} else {
		logger.Info("✓ No migrations to apply")
	}

	// Step 4: Validate deployment
	logger.Info("Step 4/4: Validating deployment...")
	if err := o.validateDeployment(); err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Validation failed: %v", err))
		return report, fmt.Errorf("validation failed: %w", err)
	}
	report.ValidationPassed = true
	logger.Info("✓ Deployment validated successfully")

	report.Success = true
	report.EndTime = time.Now()
	report.Duration = report.EndTime.Sub(report.StartTime)

	logger.Info("═══════════════════════════════════════════════════")
	logger.Info("Deployment Completed Successfully")
	logger.Info("Duration: %v", report.Duration)
	logger.Info("═══════════════════════════════════════════════════")

	return report, nil
}

// createBackup creates a pre-deployment backup
func (o *DeploymentOrchestrator) createBackup() (string, error) {
	// Get current version for backup metadata
	runner := migrations.NewMigrationRunner(o.DB)
	currentVer, migNum, err := runner.GetCurrentVersion()
	if err != nil {
		currentVer = "unknown"
		migNum = 0
	}

	// Create backup using backup service
	backupService := backup.NewBackupService()
	metadata, err := backupService.CreateJSONBackup(currentVer, migNum)
	if err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	return metadata.FilePath, nil
}

// checkCompatibility verifies version compatibility
func (o *DeploymentOrchestrator) checkCompatibility() error {
	runner := migrations.NewMigrationRunner(o.DB)
	currentVer, _, err := runner.GetCurrentVersion()
	if err != nil {
		// If version table doesn't exist, this might be a fresh install
		currentVer = ""
	}

	compat := version.CheckCompatibility(currentVer)

	if !compat.Compatible && !compat.MigrationNeeded {
		return fmt.Errorf("version incompatible: %s", compat.Reason)
	}

	return nil
}

// applyMigrations applies pending database migrations
func (o *DeploymentOrchestrator) applyMigrations() (int, error) {
	runner := migrations.NewMigrationRunner(o.DB)
	runner.Verbose = true

	// Get pending migrations
	pending, err := runner.GetPendingMigrations()
	if err != nil {
		return 0, fmt.Errorf("failed to get pending migrations: %w", err)
	}

	if len(pending) == 0 {
		return 0, nil
	}

	// Apply all pending migrations
	results, err := runner.ApplyAll()
	if err != nil {
		return 0, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Count successful migrations
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	return successCount, nil
}

// validateDeployment validates the database after deployment
func (o *DeploymentOrchestrator) validateDeployment() error {
	// Check that version was updated
	runner := migrations.NewMigrationRunner(o.DB)
	currentVer, _, err := runner.GetCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get version after migration: %w", err)
	}

	// Verify version matches required version
	if currentVer != version.GetRequiredDBVersion() {
		return fmt.Errorf("version mismatch after deployment: expected %s, got %s",
			version.GetRequiredDBVersion(), currentVer)
	}

	// Basic sanity check: verify we can query the database
	var count int64
	if err := o.DB.Table("schema_version").Count(&count).Error; err != nil {
		return fmt.Errorf("database sanity check failed: %w", err)
	}

	return nil
}

// PrepareDeploymentPackage creates an export of all system and master data
func (o *DeploymentOrchestrator) PrepareDeploymentPackage(exportDir string) (*DeploymentPackage, error) {
	logger.Info("═══════════════════════════════════════════════════")
	logger.Info("Preparing Deployment Package")
	logger.Info("═══════════════════════════════════════════════════")

	pkg := &DeploymentPackage{
		Version:   config.GetVersion(),
		Timestamp: time.Now(),
	}

	// Export full database (all data, all tables)
	logger.Info("Exporting complete database...")
	result, err := transfer.ExportDatabase(exportDir)
	if err != nil {
		return nil, fmt.Errorf("database export failed: %w", err)
	}

	pkg.ExportPath = result.FilePath
	pkg.RecordCount = result.RecordCount
	logger.Info("✓ Exported %d records to %s", result.RecordCount, result.Filename)

	logger.Info("═══════════════════════════════════════════════════")
	logger.Info("Deployment Package Ready")
	logger.Info("File: %s", result.FilePath)
	logger.Info("Records: %d", result.RecordCount)
	logger.Info("═══════════════════════════════════════════════════")

	return pkg, nil
}

// DeploymentPackage contains information about a deployment package
type DeploymentPackage struct {
	Version     string
	Timestamp   time.Time
	ExportPath  string
	RecordCount int
}

// FullDeploymentWithImport performs a complete deployment including data import
func (o *DeploymentOrchestrator) FullDeploymentWithImport(importFilePath string) (*DeploymentReport, error) {
	report := &DeploymentReport{
		StartTime: time.Now(),
	}

	logger.Info("═══════════════════════════════════════════════════")
	logger.Info("Starting Full Deployment With Data Import")
	logger.Info("═══════════════════════════════════════════════════")

	// Step 1: Create backup of current state
	logger.Info("Step 1/5: Creating pre-deployment backup...")
	backupPath, err := o.createBackup()
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Backup failed: %v", err))
		return report, fmt.Errorf("backup failed: %w", err)
	}
	report.BackupCreated = true
	report.BackupPath = backupPath
	logger.Info("✓ Backup created: %s", backupPath)

	// Step 2: Export current state (before migration)
	logger.Info("Step 2/5: Exporting current database state...")
	exportDir := "./export_temp"
	exportResult, err := transfer.ExportDatabase(exportDir)
	if err != nil {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Current state export failed: %v", err))
		logger.Warn("Could not export current state: %v", err)
	} else {
		logger.Info("✓ Current state exported: %s", exportResult.Filename)
	}

	// Step 3: Check version compatibility
	logger.Info("Step 3/5: Checking version compatibility...")
	if err := o.checkCompatibility(); err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Compatibility check failed: %v", err))
		return report, fmt.Errorf("compatibility check failed: %w", err)
	}
	logger.Info("✓ Version compatibility verified")

	// Step 4: Apply migrations
	logger.Info("Step 4/5: Applying database migrations...")
	migrationsRun, err := o.applyMigrations()
	if err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Migration failed: %v", err))
		logger.Error("Migration failed! Rollback required.")
		return report, fmt.Errorf("migration failed: %w", err)
	}
	report.MigrationsRun = migrationsRun
	if migrationsRun > 0 {
		logger.Info("✓ Applied %d migration(s)", migrationsRun)
	} else {
		logger.Info("✓ No migrations needed")
	}

	// Step 5: Import data if provided
	if importFilePath != "" {
		logger.Info("Step 5/5: Importing data from %s...", importFilePath)
		importResult, err := transfer.ImportDatabase(importFilePath)
		if err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("Data import failed: %v", err))
			return report, fmt.Errorf("data import failed: %w", err)
		}
		logger.Info("✓ Imported %d records", importResult.RecordCount)
	} else {
		logger.Info("Step 5/5: No data import requested")
	}

	// Validate
	logger.Info("Validating deployment...")
	if err := o.validateDeployment(); err != nil {
		report.Errors = append(report.Errors, fmt.Sprintf("Validation failed: %v", err))
		return report, fmt.Errorf("validation failed: %w", err)
	}
	report.ValidationPassed = true
	logger.Info("✓ Deployment validated successfully")

	report.Success = true
	report.EndTime = time.Now()
	report.Duration = report.EndTime.Sub(report.StartTime)

	logger.Info("═══════════════════════════════════════════════════")
	logger.Info("Full Deployment Completed Successfully")
	logger.Info("Duration: %v", report.Duration)
	logger.Info("═══════════════════════════════════════════════════")

	return report, nil
}
