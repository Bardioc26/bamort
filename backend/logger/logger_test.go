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

func TestLoggingFunctions(t *testing.T) {
	setupTestEnvironment(t)

	// Capture output for testing
	// Since logger writes to os.Stdout, we need to capture that output
	// For simplicity, we'll test that the functions don't panic and work with different log levels

	tests := []struct {
		name    string
		logFunc func(string, ...interface{})
		level   LogLevel
		message string
		args    []interface{}
		enabled bool
	}{
		{"Debug", Debug, DEBUG, "Debug message: %s", []interface{}{"test"}, true},
		{"Info", Info, INFO, "Info message: %s", []interface{}{"test"}, true},
		{"Warn", Warn, WARN, "Warn message: %s", []interface{}{"test"}, true},
		{"Error", Error, ERROR, "Error message: %s", []interface{}{"test"}, true},
		{"Debugf", Debugf, DEBUG, "Debugf message: %s", []interface{}{"test"}, true},
		{"Infof", Infof, INFO, "Infof message: %s", []interface{}{"test"}, true},
		{"Warnf", Warnf, WARN, "Warnf message: %s", []interface{}{"test"}, true},
		{"Errorf", Errorf, ERROR, "Errorf message: %s", []interface{}{"test"}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set appropriate log level and debug mode
			SetDebugMode(true)
			SetMinLogLevel(DEBUG)

			// Test that the function doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s function panicked: %v", test.name, r)
				}
			}()

			// Call the logging function
			test.logFunc(test.message, test.args...)

			// The test passes if no panic occurred
			// In a more sophisticated test, we could capture stdout and verify the output format
		})
	}
}

func TestLoggingWithDifferentLevels(t *testing.T) {
	setupTestEnvironment(t)

	// Test that messages are filtered based on log level
	tests := []struct {
		minLevel    LogLevel
		debugMode   bool
		description string
	}{
		{ERROR, false, "ERROR level - only errors should log"},
		{WARN, false, "WARN level - warnings and errors should log"},
		{INFO, false, "INFO level - info, warnings and errors should log"},
		{DEBUG, true, "DEBUG level with debug mode - all should log"},
		{DEBUG, false, "DEBUG level without debug mode - debug should not log"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			// Configure logger
			SetDebugMode(test.debugMode)
			SetMinLogLevel(test.minLevel)

			// Test that functions don't panic with this configuration
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Logging function panicked with %s: %v", test.description, r)
				}
			}()

			// Call all logging functions
			Debug("Debug message")
			Info("Info message")
			Warn("Warn message")
			Error("Error message")
			Debugf("Debug message %s", "formatted")
			Infof("Info message %s", "formatted")
			Warnf("Warn message %s", "formatted")
			Errorf("Error message %s", "formatted")

			// Test passes if no panics occurred
		})
	}
}

func TestLoggingWithNoArgs(t *testing.T) {
	setupTestEnvironment(t)

	// Test logging functions with no format arguments
	SetDebugMode(true)
	SetMinLogLevel(DEBUG)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logging function panicked with no args: %v", r)
		}
	}()

	// Test all functions with simple messages (no format placeholders)
	Debug("Simple debug message")
	Info("Simple info message")
	Warn("Simple warn message")
	Error("Simple error message")
	Debugf("Simple debugf message")
	Infof("Simple infof message")
	Warnf("Simple warnf message")
	Errorf("Simple errorf message")
}

func TestLoggingWithMultipleArgs(t *testing.T) {
	setupTestEnvironment(t)

	// Test logging functions with multiple format arguments
	SetDebugMode(true)
	SetMinLogLevel(DEBUG)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logging function panicked with multiple args: %v", r)
		}
	}()

	// Test all functions with multiple format arguments
	Debug("Debug: %s %d %v", "test", 42, true)
	Info("Info: %s %d %v", "test", 42, true)
	Warn("Warn: %s %d %v", "test", 42, true)
	Error("Error: %s %d %v", "test", 42, true)
	Debugf("Debugf: %s %d %v", "test", 42, true)
	Infof("Infof: %s %d %v", "test", 42, true)
	Warnf("Warnf: %s %d %v", "test", 42, true)
	Errorf("Errorf: %s %d %v", "test", 42, true)
}

func TestDebugModeFiltering(t *testing.T) {
	setupTestEnvironment(t)

	// Test that Debug messages are actually filtered when debug mode is disabled
	SetMinLogLevel(DEBUG) // Allow debug level
	SetDebugMode(false)   // But disable debug mode

	// Verify IsDebugEnabled returns false
	if IsDebugEnabled() {
		t.Error("IsDebugEnabled should return false when debug mode is disabled")
	}

	// Test that debug functions don't panic even when filtered
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Debug function panicked when filtered: %v", r)
		}
	}()

	Debug("This debug message should be filtered")
	Debugf("This debug message %s should be filtered", "also")

	// Now enable debug mode and verify IsDebugEnabled returns true
	SetDebugMode(true)
	if !IsDebugEnabled() {
		t.Error("IsDebugEnabled should return true when debug mode is enabled and level is DEBUG")
	}

	// These should not panic
	Debug("This debug message should now log")
	Debugf("This debug message %s should now log", "also")
}

func TestLogLevelFiltering(t *testing.T) {
	setupTestEnvironment(t)

	// Test that messages below the minimum level are filtered
	SetDebugMode(false)  // Disable debug mode for this test
	SetMinLogLevel(WARN) // Only warnings and errors should log

	// These should not panic even though they're filtered
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logging function panicked when filtered: %v", r)
		}
	}()

	Debug("This should be filtered") // Below WARN level
	Info("This should be filtered")  // Below WARN level
	Warn("This should log")          // At WARN level
	Error("This should log")         // Above WARN level

	Debugf("This should be filtered") // Below WARN level
	Infof("This should be filtered")  // Below WARN level
	Warnf("This should log")          // At WARN level
	Errorf("This should log")         // Above WARN level
}
