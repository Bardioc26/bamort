package logger

import (
	"os"
	"testing"
)

// setupTestEnvironment setzt ENVIRONMENT=test für Tests
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

func TestLogLevels(t *testing.T) {
	setupTestEnvironment(t)
	// Test String-Representation der Log-Levels
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.level.String())
		}
	}
}

func TestDebugModeFromEnv(t *testing.T) {
	setupTestEnvironment(t)

	// Test verschiedene Umgebungsvariablen-Werte
	tests := []struct {
		envValue string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"True", true},
		{"1", true},
		{"false", false},
		{"FALSE", false},
		{"False", false},
		{"0", false},
		{"", false},
		{"invalid", false},
	}

	for _, test := range tests {
		os.Setenv("DEBUG", test.envValue)
		result := getDebugModeFromEnv()
		if result != test.expected {
			t.Errorf("For DEBUG=%s, expected %v, got %v", test.envValue, test.expected, result)
		}
	}

	// Cleanup
	os.Unsetenv("DEBUG")
}

func TestMinLogLevelFromEnv(t *testing.T) {
	setupTestEnvironment(t)

	// Test verschiedene LOG_LEVEL Werte
	tests := []struct {
		envValue string
		expected LogLevel
	}{
		{"DEBUG", DEBUG},
		{"debug", DEBUG},
		{"INFO", INFO},
		{"info", INFO},
		{"WARN", WARN},
		{"warn", WARN},
		{"ERROR", ERROR},
		{"error", ERROR},
		{"invalid", INFO}, // Default zu INFO
		{"", INFO},        // Default zu INFO
	}

	// DEBUG-Modus deaktivieren für diese Tests
	os.Setenv("DEBUG", "false")

	for _, test := range tests {
		os.Setenv("LOG_LEVEL", test.envValue)
		result := getMinLogLevelFromEnv()
		if result != test.expected {
			t.Errorf("For LOG_LEVEL=%s, expected %v, got %v", test.envValue, test.expected, result)
		}
	}

	// Test: Wenn DEBUG=true und kein LOG_LEVEL gesetzt, sollte DEBUG returned werden
	os.Setenv("DEBUG", "true")
	os.Unsetenv("LOG_LEVEL")
	result := getMinLogLevelFromEnv()
	if result != DEBUG {
		t.Errorf("When DEBUG=true and no LOG_LEVEL set, expected DEBUG, got %v", result)
	}

	// Cleanup
	os.Unsetenv("DEBUG")
	os.Unsetenv("LOG_LEVEL")
}

func TestSetDebugMode(t *testing.T) {
	setupTestEnvironment(t)

	// Test Debug-Modus aktivieren
	SetDebugMode(true)
	if !IsDebugEnabled() {
		t.Error("Debug mode should be enabled")
	}

	// Test Debug-Modus deaktivieren
	SetDebugMode(false)
	if IsDebugEnabled() {
		t.Error("Debug mode should be disabled")
	}
}

func TestSetMinLogLevel(t *testing.T) {
	setupTestEnvironment(t)

	// Test verschiedene Log-Level setzen
	levels := []LogLevel{DEBUG, INFO, WARN, ERROR}

	for _, level := range levels {
		SetMinLogLevel(level)
		if defaultLogger.minLevel != level {
			t.Errorf("Expected min log level %v, got %v", level, defaultLogger.minLevel)
		}
	}
}
