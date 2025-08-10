package config

import (
	"os"
	"testing"
)

// setupTestEnvironment setzt ENVIRONMENT=test für Tests und stellt es nach dem Test wieder her
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

func TestLoadEnvFile(t *testing.T) {
	setupTestEnvironment(t)
	// Test-Datei erstellen
	envContent := `# Test .env file
DEBUG=true
LOG_LEVEL=DEBUG
PORT=9999
# Comment line should be ignored

ENVIRONMENT=test
DATABASE_URL="postgresql://test:test@localhost:5432/test"
QUOTED_VALUE='single quotes'
`

	// Temporäre .env-Datei erstellen
	err := os.WriteFile(".env.test", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen der Test-.env-Datei: %v", err)
	}
	defer os.Remove(".env.test")

	// Ursprüngliche Umgebungsvariablen sichern
	originalDebug := os.Getenv("DEBUG")
	originalLogLevel := os.Getenv("LOG_LEVEL")
	originalPort := os.Getenv("PORT")
	originalEnv := os.Getenv("ENVIRONMENT")
	originalDB := os.Getenv("DATABASE_URL")
	originalQuoted := os.Getenv("QUOTED_VALUE")

	// Umgebungsvariablen zurücksetzen
	os.Unsetenv("DEBUG")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("QUOTED_VALUE")

	// Test-Datei laden
	loadEnvFileContent(".env.test")

	// Tests
	tests := []struct {
		key      string
		expected string
	}{
		{"DEBUG", "true"},
		{"LOG_LEVEL", "DEBUG"},
		{"PORT", "9999"},
		{"ENVIRONMENT", "test"},
		{"DATABASE_URL", "postgresql://test:test@localhost:5432/test"},
		{"QUOTED_VALUE", "single quotes"},
	}

	for _, test := range tests {
		if value := os.Getenv(test.key); value != test.expected {
			t.Errorf("Für %s: erwartet '%s', erhalten '%s'", test.key, test.expected, value)
		}
	}

	// Ursprüngliche Werte wiederherstellen
	if originalDebug != "" {
		os.Setenv("DEBUG", originalDebug)
	} else {
		os.Unsetenv("DEBUG")
	}
	if originalLogLevel != "" {
		os.Setenv("LOG_LEVEL", originalLogLevel)
	} else {
		os.Unsetenv("LOG_LEVEL")
	}
	if originalPort != "" {
		os.Setenv("PORT", originalPort)
	} else {
		os.Unsetenv("PORT")
	}
	if originalEnv != "" {
		os.Setenv("ENVIRONMENT", originalEnv)
	} else {
		os.Unsetenv("ENVIRONMENT")
	}
	if originalDB != "" {
		os.Setenv("DATABASE_URL", originalDB)
	} else {
		os.Unsetenv("DATABASE_URL")
	}
	if originalQuoted != "" {
		os.Setenv("QUOTED_VALUE", originalQuoted)
	} else {
		os.Unsetenv("QUOTED_VALUE")
	}
}

func TestEnvVariablesPrecedence(t *testing.T) {
	setupTestEnvironment(t)

	// Test, dass bereits gesetzte Umgebungsvariablen Vorrang haben
	envContent := `DEBUG=false
LOG_LEVEL=ERROR`

	// Temporäre .env-Datei erstellen
	err := os.WriteFile(".env.precedence", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen der Test-.env-Datei: %v", err)
	}
	defer os.Remove(".env.precedence")

	// Umgebungsvariable vorher setzen
	os.Setenv("DEBUG", "true")
	os.Setenv("LOG_LEVEL", "INFO")

	// .env-Datei laden
	loadEnvFileContent(".env.precedence")

	// Tests - bereits gesetzte Werte sollten nicht überschrieben werden
	if debug := os.Getenv("DEBUG"); debug != "true" {
		t.Errorf("DEBUG sollte 'true' bleiben, aber ist '%s'", debug)
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "INFO" {
		t.Errorf("LOG_LEVEL sollte 'INFO' bleiben, aber ist '%s'", logLevel)
	}

	// Cleanup
	os.Unsetenv("DEBUG")
	os.Unsetenv("LOG_LEVEL")
}

func TestLoadConfigWithEnvFile(t *testing.T) {
	setupTestEnvironment(t)

	// Test-Konfiguration mit .env-Datei
	envContent := `ENVIRONMENT=development
DEBUG=true
LOG_LEVEL=DEBUG
PORT=7777
DATABASE_URL=test://localhost/testdb`

	// Temporäre .env-Datei erstellen
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Fehler beim Erstellen der .env-Datei: %v", err)
	}
	defer os.Remove(".env")

	// Alle relevanten Umgebungsvariablen zurücksetzen
	originalVars := map[string]string{
		"ENVIRONMENT":  os.Getenv("ENVIRONMENT"),
		"DEBUG":        os.Getenv("DEBUG"),
		"LOG_LEVEL":    os.Getenv("LOG_LEVEL"),
		"PORT":         os.Getenv("PORT"),
		"DATABASE_URL": os.Getenv("DATABASE_URL"),
	}

	for key := range originalVars {
		os.Unsetenv(key)
	}

	// Konfiguration laden
	config := LoadConfig()

	// Tests
	if config.Environment != "development" {
		t.Errorf("Environment: erwartet 'development', erhalten '%s'", config.Environment)
	}

	if !config.DebugMode {
		t.Error("DebugMode sollte true sein")
	}

	if config.LogLevel != "DEBUG" {
		t.Errorf("LogLevel: erwartet 'DEBUG', erhalten '%s'", config.LogLevel)
	}

	if config.ServerPort != "7777" {
		t.Errorf("ServerPort: erwartet '7777', erhalten '%s'", config.ServerPort)
	}

	if config.DatabaseURL != "test://localhost/testdb" {
		t.Errorf("DatabaseURL: erwartet 'test://localhost/testdb', erhalten '%s'", config.DatabaseURL)
	}

	// Ursprüngliche Werte wiederherstellen
	for key, value := range originalVars {
		if value != "" {
			os.Setenv(key, value)
		} else {
			os.Unsetenv(key)
		}
	}
}
