package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config enthält alle Anwendungskonfigurationen
type Config struct {
	// Server Konfiguration
	ServerPort string

	// Database Konfiguration
	DatabaseURL  string
	DatabaseType string

	// Logging Konfiguration
	DebugMode bool
	LogLevel  string

	// Environment
	Environment string

	DevTesting string // "yes" or "no", used to determine if we are in a test environment

	// PDF Templates
	TemplatesDir  string // Directory where PDF templates are stored
	ExportTempDir string // Directory for temporary PDF exports
}

// Cfg ist die globale Konfigurationsvariable
// Sie wird beim Programmstart automatisch geladen
var Cfg *Config

// init lädt die Konfiguration einmal beim Programmstart
func init() {
	Cfg = LoadConfig()
}

// defaultConfig gibt die Standard-Konfiguration zurück
func defaultConfig() *Config {
	return &Config{
		ServerPort:    "8180",
		DatabaseURL:   "",
		DatabaseType:  "mysql",
		DebugMode:     false,
		LogLevel:      "INFO",
		Environment:   "production",
		DevTesting:    "no",          // Default to "no", can be overridden in tests
		TemplatesDir:  "./templates", // Default templates directory
		ExportTempDir: "./xporttemp", // Default export temp directory
	}
}

// LoadConfig lädt die Konfiguration aus Umgebungsvariablen
func LoadConfig() *Config {
	// Lade .env-Datei falls vorhanden
	loadEnvFile()

	config := defaultConfig()

	// Debug: Zeige geladene Umgebungsvariablen
	fmt.Printf("DEBUG LoadConfig - ENVIRONMENT aus ENV: '%s'\n", os.Getenv("ENVIRONMENT"))
	fmt.Printf("DEBUG LoadConfig - TESTING aus ENV: '%s'\n", os.Getenv("DEVTESTING"))
	fmt.Printf("DEBUG LoadConfig - DATABASE_TYPE aus ENV: '%s'\n", os.Getenv("DATABASE_TYPE"))

	// Server Port
	if port := os.Getenv("PORT"); port != "" {
		config.ServerPort = port
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.ServerPort = port
	}

	// Database
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		config.DatabaseURL = dbURL
	}
	if dbType := os.Getenv("DATABASE_TYPE"); dbType != "" {
		config.DatabaseType = strings.ToLower(dbType)
	}

	// Debug Mode
	if debug := os.Getenv("DEBUG"); debug != "" {
		config.DebugMode = strings.ToLower(debug) == "true" || debug == "1"
	}

	// Log Level
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.LogLevel = strings.ToUpper(logLevel)
	}

	// Environment
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		config.Environment = strings.ToLower(env)
	}
	if env := os.Getenv("GO_ENV"); env != "" {
		config.Environment = strings.ToLower(env)
	}

	// Automatisch Debug-Modus für Development-Environment aktivieren
	if config.Environment == "development" || config.Environment == "dev" {
		config.DebugMode = true
		if config.LogLevel == "INFO" {
			config.LogLevel = "DEBUG"
		}
	}
	// Testing  in Development
	if testing := os.Getenv("DEVTESTING"); testing != "" {
		config.DevTesting = strings.ToLower(testing)
		fmt.Printf("DEBUG LoadConfig - DEVTESTING gefunden: '%s' -> DevTesting: '%s'\n", testing, config.DevTesting)
	} else {
		config.DevTesting = "no" // Default to "no"
		fmt.Printf("DEBUG LoadConfig - DEVTESTING nicht gefunden, setze DevTesting auf 'no'\n")
	}

	// Templates Directory
	if templatesDir := os.Getenv("TEMPLATES_DIR"); templatesDir != "" {
		config.TemplatesDir = templatesDir
	}

	// Export Temp Directory
	if exportTempDir := os.Getenv("EXPORT_TEMP_DIR"); exportTempDir != "" {
		config.ExportTempDir = exportTempDir
	}

	fmt.Printf("DEBUG LoadConfig - Finale Config: Environment='%s', DevTesting='%s', DatabaseType='%s'\n",
		config.Environment, config.DevTesting, config.DatabaseType)

	return config
}

// loadEnvFile lädt eine .env-Datei falls vorhanden
func loadEnvFile() {
	envFiles := []string{".env", ".env.local"}

	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err == nil {
			fmt.Printf("DEBUG loadEnvFile - Lade .env-Datei: %s\n", envFile)
			loadEnvFileContent(envFile)
		} else {
			fmt.Printf("DEBUG loadEnvFile - .env-Datei nicht gefunden: %s\n", envFile)
		}
	}
}

// loadEnvFileContent lädt den Inhalt einer .env-Datei
func loadEnvFileContent(filename string) {
	fmt.Printf("DEBUG loadEnvFileContent - Öffne Datei: %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("DEBUG loadEnvFileContent - Fehler beim Öffnen von %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Überspringe leere Zeilen und Kommentare
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Teile die Zeile in Key=Value auf
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Behandle Kommentare am Ende der Zeile (nur wenn nicht in Anführungszeichen)
		if !strings.HasPrefix(value, `"`) && !strings.HasPrefix(value, `'`) {
			// Suche nach Kommentar am Ende der Zeile
			if commentPos := strings.Index(value, "#"); commentPos > 0 {
				// Entferne Kommentar und Leerzeichen davor
				value = strings.TrimSpace(value[:commentPos])
			}
		}

		// Entferne Anführungszeichen falls vorhanden
		value = strings.Trim(value, `"'`)

		fmt.Printf("DEBUG loadEnvFileContent - Zeile %d: %s='%s' (nach Kommentar-Behandlung)\n", lineNum, key, value)

		// Setze die Umgebungsvariable nur, wenn sie noch nicht gesetzt ist
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
			fmt.Printf("DEBUG loadEnvFileContent - Setze ENV %s='%s'\n", key, value)
		} else {
			fmt.Printf("DEBUG loadEnvFileContent - ENV %s bereits gesetzt, überspringe\n", key)
		}
	}
	fmt.Printf("DEBUG loadEnvFileContent - Datei %s vollständig verarbeitet (%d Zeilen)\n", filename, lineNum)
}

// IsDevelopment prüft, ob die Anwendung im Development-Modus läuft
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development" || c.Environment == "dev"
}

// IsProduction prüft, ob die Anwendung im Production-Modus läuft
func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// GetServerAddress gibt die vollständige Server-Adresse zurück
func (c *Config) GetServerAddress() string {
	return ":" + c.ServerPort
}

// GetBoolEnv ist eine Hilfsfunktion zum Laden von Boolean-Umgebungsvariablen
func GetBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
		return strings.ToLower(value) == "true" || value == "1"
	}
	return defaultValue
}

// GetIntEnv ist eine Hilfsfunktion zum Laden von Integer-Umgebungsvariablen
func GetIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
