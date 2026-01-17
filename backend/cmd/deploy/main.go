package main

import (
	"bamort/config"
	"bamort/database"
	"bamort/deployment"
	"bamort/deployment/migrations"
	"bamort/deployment/validator"
	"bamort/deployment/version"
	"fmt"
	"os"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		cmdVersion()
	case "status":
		cmdStatus()
	case "prepare":
		cmdPrepare()
	case "deploy":
		cmdDeploy()
	case "validate":
		cmdValidate()
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("%s✗ Unknown command: %s%s\n", ColorRed, command, ColorReset)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf("\n%s%sBaMoRT Deployment Tool%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("Version: %s\n\n", config.GetVersion())
	fmt.Println("Usage: deploy <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Printf("  %sprepare%s [dir]     Create deployment package (export all master data)\n", ColorGreen, ColorReset)
	fmt.Printf("  %sdeploy%s [dir]      Run full deployment (backup → migrate → import → validate)\n", ColorGreen, ColorReset)
	fmt.Printf("  %svalidate%s          Validate database schema and data integrity\n", ColorGreen, ColorReset)
	fmt.Printf("  %sstatus%s            Show current database version and pending migrations\n", ColorGreen, ColorReset)
	fmt.Printf("  %sversion%s           Show version information\n", ColorGreen, ColorReset)
	fmt.Printf("  %shelp%s              Show this help message\n", ColorGreen, ColorReset)
	fmt.Println("\nArguments:")
	fmt.Printf("  %s[dir]%s             Directory for export/import (default: ./export_temp)\n", ColorCyan, ColorReset)
	fmt.Println("\nExamples:")
	fmt.Println("  deploy prepare              # Create deployment package in ./export_temp")
	fmt.Println("  deploy prepare /path/pkg    # Create deployment package in /path/pkg")
	fmt.Println("  deploy deploy               # Run deployment without importing data")
	fmt.Println("  deploy deploy ./export_temp # Run deployment and import master data")
	fmt.Println("  deploy validate             # Validate database schema")
	fmt.Println("\nDeployment Workflow:")
	fmt.Println("  Source System:  deploy prepare /shared/pkg     # Export master data")
	fmt.Println("  Target System:  deploy deploy /shared/pkg      # Migrate DB + Import data")
	fmt.Println()
}

func cmdVersion() {
	fmt.Printf("\n%s%sBaMoRT Deployment Tool%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("Backend Version: %s%s%s\n", ColorGreen, config.GetVersion(), ColorReset)
	fmt.Printf("Required DB Version: %s%s%s\n", ColorGreen, version.GetRequiredDBVersion(), ColorReset)
	fmt.Println()
}

func cmdStatus() {
	printBanner("Database Status")

	// Connect to database
	database.DB = database.ConnectDatabase()
	if database.DB == nil {
		printError("Failed to connect to database")
		os.Exit(1)
	}

	runner := migrations.NewMigrationRunner(database.DB)

	currentVer, _, err := runner.GetCurrentVersion()
	if err != nil {
		if strings.Contains(err.Error(), "no such table") {
			printWarning("Database not initialized")
			fmt.Printf("\nDatabase appears to be uninitialized.\n\n")
			return
		}
		printError("Failed to get current version: %v", err)
		os.Exit(1)
	}

	fmt.Printf("\n%sCurrent Database Version:%s %s%s%s\n", ColorBold, ColorReset, ColorCyan, currentVer, ColorReset)
	fmt.Printf("%sBackend Version:%s %s%s%s\n", ColorBold, ColorReset, ColorCyan, config.GetVersion(), ColorReset)
	fmt.Printf("%sRequired DB Version:%s %s%s%s\n", ColorBold, ColorReset, ColorCyan, version.GetRequiredDBVersion(), ColorReset)

	compat := version.CheckCompatibility(currentVer)
	fmt.Printf("\n%sCompatibility:%s ", ColorBold, ColorReset)

	if compat.Compatible {
		fmt.Printf("%s✓ Compatible%s\n", ColorGreen, ColorReset)
	} else if compat.MigrationNeeded {
		fmt.Printf("%s⚠ Migration Required%s\n", ColorYellow, ColorReset)
	} else {
		fmt.Printf("%s✗ Version Mismatch%s\n", ColorRed, ColorReset)
	}
	fmt.Printf("  %s\n", compat.Reason)

	pending, _ := runner.GetPendingMigrations()

	if len(pending) > 0 {
		fmt.Printf("\n%sPending Migrations: %d%s\n", ColorYellow, len(pending), ColorReset)
		for _, m := range pending {
			fmt.Printf("  • Migration %d: %s (→ %s)\n", m.Number, m.Description, m.Version)
		}
	} else {
		fmt.Printf("\n%s✓ No pending migrations%s\n", ColorGreen, ColorReset)
	}

	fmt.Println()
}

func printBanner(title string) {
	fmt.Printf("\n%s%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("%s%s %s%s\n", ColorBold, ColorCyan, title, ColorReset)
	fmt.Printf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", ColorCyan, ColorReset)
}

func printError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s✗ "+format+"%s\n", append([]interface{}{ColorRed}, append(args, ColorReset)...)...)
}

func printWarning(format string, args ...interface{}) {
	fmt.Printf("%s⚠ "+format+"%s\n", append([]interface{}{ColorYellow}, append(args, ColorReset)...)...)
}

func printSuccess(format string, args ...interface{}) {
	fmt.Printf("%s✓ "+format+"%s\n", append([]interface{}{ColorGreen}, append(args, ColorReset)...)...)
}

// cmdPrepare creates a deployment package with full database export
func cmdPrepare() {
	printBanner("Prepare Deployment Package")

	// Connect to database
	database.DB = database.ConnectDatabase()
	if database.DB == nil {
		printError("Failed to connect to database")
		os.Exit(1)
	}

	orchestrator := deployment.NewOrchestrator(database.DB)

	exportDir := "./export_temp"
	if len(os.Args) > 2 {
		exportDir = os.Args[2]
	}

	fmt.Printf("\nExporting to: %s%s%s\n", ColorCyan, exportDir, ColorReset)
	fmt.Println("This will create a complete backup of all system and master data...")
	fmt.Println()

	pkg, err := orchestrator.PrepareDeploymentPackage(exportDir)
	if err != nil {
		printError("Failed to prepare deployment package: %v", err)
		os.Exit(1)
	}

	fmt.Println()
	printSuccess("Deployment package created successfully!")
	fmt.Printf("\n%sPackage Details:%s\n", ColorBold, ColorReset)
	fmt.Printf("  Version:      %s\n", pkg.Version)
	fmt.Printf("  Export File:  %s\n", pkg.ExportPath)
	fmt.Printf("  Timestamp:    %s\n", pkg.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Println("This export can be imported on the target system after migration.")
	fmt.Println()
}

// cmdDeploy runs the full deployment workflow
func cmdDeploy() {
	printBanner("Full Deployment")

	// Connect to database
	database.DB = database.ConnectDatabase()
	if database.DB == nil {
		printError("Failed to connect to database")
		os.Exit(1)
	}

	orchestrator := deployment.NewOrchestrator(database.DB)

	// Check if import directory is provided
	importDir := ""
	if len(os.Args) > 2 {
		importDir = os.Args[2]
	}

	if importDir != "" {
		fmt.Println("\nThis will:")
		fmt.Println("  1. Create a backup of the current database")
		fmt.Println("  2. Export current master data state")
		fmt.Println("  3. Check version compatibility")
		fmt.Println("  4. Apply pending migrations")
		fmt.Printf("  5. Import master data from: %s%s%s\n", ColorCyan, importDir, ColorReset)
		fmt.Println("  6. Validate the deployment")
	} else {
		fmt.Println("\nThis will:")
		fmt.Println("  1. Create a backup of the current database")
		fmt.Println("  2. Check version compatibility")
		fmt.Println("  3. Apply pending migrations")
		fmt.Println("  4. Validate the deployment")
		fmt.Println()
		fmt.Printf("%sNOTE:%s No import directory specified. Master data will not be imported.\n", ColorYellow, ColorReset)
	}
	fmt.Println()

	fmt.Printf("%sWARNING:%s This operation will modify the database!\n", ColorYellow, ColorReset)
	fmt.Print("Continue? (yes/no): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" && confirm != "y" {
		fmt.Println("Deployment cancelled.")
		os.Exit(0)
	}

	fmt.Println()

	// Run deployment (with or without import based on importDir)
	report, err := orchestrator.FullDeploymentWithImport(importDir)

	if err != nil {
		printError("Deployment failed: %v", err)
		fmt.Println()
		if report.BackupCreated {
			fmt.Printf("Backup available at: %s\n", report.BackupPath)
		}
		if len(report.Errors) > 0 {
			fmt.Printf("\n%sErrors:%s\n", ColorRed, ColorReset)
			for _, e := range report.Errors {
				fmt.Printf("  • %s\n", e)
			}
		}
		os.Exit(1)
	}

	fmt.Println()
	printSuccess("Deployment completed successfully!")
	fmt.Printf("\n%sDeployment Summary:%s\n", ColorBold, ColorReset)
	fmt.Printf("  Backup:       %s\n", report.BackupPath)
	fmt.Printf("  Migrations:   %d applied\n", report.MigrationsRun)
	if importDir != "" {
		fmt.Printf("  Data Import:  %s✓ Master data imported%s\n", ColorGreen, ColorReset)
	} else {
		fmt.Printf("  Data Import:  %s- Not performed%s\n", ColorYellow, ColorReset)
	}
	fmt.Printf("  Duration:     %v\n", report.Duration)
	fmt.Printf("  Validated:    %s✓%s\n", ColorGreen, ColorReset)
	if len(report.Warnings) > 0 {
		fmt.Printf("\n%sWarnings:%s\n", ColorYellow, ColorReset)
		for _, w := range report.Warnings {
			fmt.Printf("  ⚠ %s\n", w)
		}
	}
	fmt.Println()
}

// cmdValidate validates the database schema and data
func cmdValidate() {
	printBanner("Database Validation")

	// Connect to database
	database.DB = database.ConnectDatabase()
	if database.DB == nil {
		printError("Failed to connect to database")
		os.Exit(1)
	}

	v := validator.NewValidator(database.DB)

	fmt.Println("\nValidating database schema and data integrity...")
	fmt.Println()

	report, err := v.Validate()
	if err != nil {
		printError("Validation failed: %v", err)
		os.Exit(1)
	}

	fmt.Printf("\n%sValidation Results:%s\n", ColorBold, ColorReset)
	fmt.Printf("  Tables Checked: %d\n", report.TablesChecked)
	fmt.Printf("  Tables Valid:   %d\n", report.TablesValid)

	if len(report.Errors) > 0 {
		fmt.Printf("\n%sErrors (%d):%s\n", ColorRed, len(report.Errors), ColorReset)
		for _, e := range report.Errors {
			fmt.Printf("  %s✗%s %s\n", ColorRed, ColorReset, e)
		}
	}

	if len(report.Warnings) > 0 {
		fmt.Printf("\n%sWarnings (%d):%s\n", ColorYellow, len(report.Warnings), ColorReset)
		for _, w := range report.Warnings {
			fmt.Printf("  %s⚠%s %s\n", ColorYellow, ColorReset, w)
		}
	}

	if len(report.MissingTables) > 0 {
		fmt.Printf("\n%sMissing Tables:%s\n", ColorRed, ColorReset)
		for _, t := range report.MissingTables {
			fmt.Printf("  • %s\n", t)
		}
	}

	if len(report.MissingColumns) > 0 {
		fmt.Printf("\n%sMissing Columns:%s\n", ColorRed, ColorReset)
		for table, cols := range report.MissingColumns {
			fmt.Printf("  %s: %v\n", table, cols)
		}
	}

	fmt.Println()
	if report.Success {
		printSuccess("Validation passed!")
	} else {
		printError("Validation failed with %d error(s)", len(report.Errors))
		os.Exit(1)
	}
	fmt.Println()
}
