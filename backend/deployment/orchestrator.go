package deployment

import (
	"archive/tar"
	"bamort/config"
	"bamort/deployment/backup"
	"bamort/deployment/migrations"
	"bamort/deployment/version"
	"bamort/gsmaster"
	"bamort/logger"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	// Export all master data (system data, rules, equipment, etc.)
	logger.Info("Exporting master data...")
	err := gsmaster.ExportAll(exportDir)
	if err != nil {
		return nil, fmt.Errorf("master data export failed: %w", err)
	}

	pkg.ExportPath = exportDir
	logger.Info("✓ Master data exported to %s", exportDir)

	// Create tar.gz archive
	logger.Info("Creating deployment package archive...")
	tarballName := fmt.Sprintf("deployment_package_%s_%s.tar.gz",
		config.GetVersion(),
		time.Now().Format("20060102-150405"))
	tarballPath := filepath.Join(filepath.Dir(exportDir), tarballName)

	err = createTarGz(exportDir, tarballPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create tar.gz archive: %w", err)
	}

	pkg.TarballPath = tarballPath
	logger.Info("✓ Package archive created: %s", tarballPath)

	logger.Info("═══════════════════════════════════════════════════")
	logger.Info("Deployment Package Ready")
	logger.Info("Export Directory: %s", exportDir)
	logger.Info("Archive:          %s", tarballPath)
	logger.Info("═══════════════════════════════════════════════════")

	return pkg, nil
}

// DeploymentPackage contains information about a deployment package
type DeploymentPackage struct {
	Version     string
	Timestamp   time.Time
	ExportPath  string
	TarballPath string
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
	logger.Info("Step 2/5: Exporting current master data state...")
	exportDir := "./tmp"
	err = gsmaster.ExportAll(exportDir)
	if err != nil {
		report.Warnings = append(report.Warnings, fmt.Sprintf("Current state export failed: %v", err))
		logger.Warn("Could not export current state: %v", err)
	} else {
		logger.Info("✓ Current state exported to: %s", exportDir)
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
		logger.Info("Step 5/5: Importing master data from %s...", importFilePath)
		err := gsmaster.ImportAll(importFilePath)
		if err != nil {
			report.Errors = append(report.Errors, fmt.Sprintf("Master data import failed: %v", err))
			return report, fmt.Errorf("master data import failed: %w", err)
		}
		logger.Info("✓ Master data imported successfully")
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

// createTarGz creates a tar.gz archive from a directory
func createTarGz(sourceDir, targetPath string) error {
	// Create the tar.gz file
	outFile, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create tar.gz file: %w", err)
	}
	defer outFile.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(outFile)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Get the base name for the archive
	baseName := filepath.Base(sourceDir)

	// Walk the directory tree
	err = filepath.Walk(sourceDir, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return fmt.Errorf("failed to create tar header: %w", err)
		}

		// Update the name to be relative to the source dir
		relPath, err := filepath.Rel(sourceDir, file)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}
		header.Name = filepath.Join(baseName, relPath)

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("failed to write tar header: %w", err)
		}

		// If it's a file, write its content
		if !fi.IsDir() {
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer f.Close()

			if _, err := io.Copy(tarWriter, f); err != nil {
				return fmt.Errorf("failed to write file content: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	return nil
}
