package config

import (
	"bufio"
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

	Testing string // "yes" or "no", used to determine if we are in a test environment
}

// defaultConfig gibt die Standard-Konfiguration zurück
func defaultConfig() *Config {
	return &Config{
		ServerPort:   "8180",
		DatabaseURL:  "",
		DatabaseType: "mysql",
		DebugMode:    false,
		LogLevel:     "INFO",
		Environment:  "production",
		Testing:      "no", // Default to "no", can be overridden in tests
	}
}

// LoadConfig lädt die Konfiguration aus Umgebungsvariablen
func LoadConfig() *Config {
	// Lade .env-Datei falls vorhanden
	loadEnvFile()

	config := defaultConfig()

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
	if testing := os.Getenv("TESTING"); testing != "" {
		config.Testing = strings.ToLower(testing)
	} else {
		config.Testing = "no" // Default to "no"
	}

	return config
}

// loadEnvFile lädt eine .env-Datei falls vorhanden
func loadEnvFile() {
	envFiles := []string{".env", ".env.local"}

	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err == nil {
			loadEnvFileContent(envFile)
		}
	}
}

// loadEnvFileContent lädt den Inhalt einer .env-Datei
func loadEnvFileContent(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
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

		// Entferne Anführungszeichen falls vorhanden
		value = strings.Trim(value, `"'`)

		// Setze die Umgebungsvariable nur, wenn sie noch nicht gesetzt ist
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
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
