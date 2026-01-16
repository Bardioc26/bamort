package main

import (
	"bamort/config"
	"bamort/database"
	"bamort/deployment/migrations"
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
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Printf("%s✗ Unknown command: %s%s\n", ColorRed, command, ColorReset)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf("\n%s%sBamort Deployment Tool%s\n", ColorBold, ColorCyan, ColorReset)
	fmt.Printf("Version: %s\n\n", config.GetVersion())
	fmt.Println("Usage: deploy <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Printf("  %sstatus%s            Show current database version and pending migrations\n", ColorGreen, ColorReset)
	fmt.Printf("  %sversion%s           Show version information\n", ColorGreen, ColorReset)
	fmt.Printf("  %shelp%s              Show this help message\n", ColorGreen, ColorReset)
	fmt.Println()
}

func cmdVersion() {
	fmt.Printf("\n%s%sBamort Deployment Tool%s\n", ColorBold, ColorCyan, ColorReset)
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
