package testutils

import (
	"os"
	"testing"
)

// SetupTestEnvironment konfiguriert die Umgebung für Tests
// Diese Funktion setzt ENVIRONMENT=test und stellt sicher, dass nach dem Test
// die ursprüngliche Umgebung wiederhergestellt wird
func SetupTestEnvironment(t *testing.T) {
	// Sicherstellen, dass ENVIRONMENT auf "test" gesetzt ist
	originalEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")

	// Cleanup-Funktion registrieren
	t.Cleanup(func() {
		if originalEnv != "" {
			os.Setenv("ENVIRONMENT", originalEnv)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})
}

// SetupTestEnvironmentWithConfig konfiguriert die Test-Umgebung mit spezifischen Werten
func SetupTestEnvironmentWithConfig(t *testing.T, envVars map[string]string) {
	// Ursprüngliche Werte sichern
	originalVars := make(map[string]string)
	for key := range envVars {
		originalVars[key] = os.Getenv(key)
	}

	// Sicherstellen, dass ENVIRONMENT auf "test" gesetzt ist
	envVars["ENVIRONMENT"] = "test"

	// Test-Umgebungsvariablen setzen
	for key, value := range envVars {
		os.Setenv(key, value)
	}

	// Cleanup-Funktion registrieren
	t.Cleanup(func() {
		for key, originalValue := range originalVars {
			if originalValue != "" {
				os.Setenv(key, originalValue)
			} else {
				os.Unsetenv(key)
			}
		}
	})
}

// EnsureTestEnvironment ist eine einfache Prüfung, ob die Test-Umgebung korrekt ist
// Kann in Tests verwendet werden um sicherzustellen, dass ENVIRONMENT=test gesetzt ist
func EnsureTestEnvironment(t *testing.T) {
	if os.Getenv("ENVIRONMENT") != "test" {
		t.Errorf("ENVIRONMENT sollte 'test' sein, ist aber '%s'. Vergessen Sie SetupTestEnvironment() aufzurufen?",
			os.Getenv("ENVIRONMENT"))
	}
}
