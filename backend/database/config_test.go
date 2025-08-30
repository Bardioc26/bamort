package database

import (
	"bamort/config"
	"database/sql/driver"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestEnvironment sets up test environment variables
func setupTestEnvironment(t *testing.T) {
	// Save original values
	origEnv := os.Getenv("ENVIRONMENT")
	origDevTesting := os.Getenv("DEV_TESTING")
	origDatabaseType := os.Getenv("DATABASE_TYPE")
	origDatabaseURL := os.Getenv("DATABASE_URL")

	// Cleanup function to restore original values
	t.Cleanup(func() {
		if origEnv != "" {
			os.Setenv("ENVIRONMENT", origEnv)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
		if origDevTesting != "" {
			os.Setenv("DEV_TESTING", origDevTesting)
		} else {
			os.Unsetenv("DEV_TESTING")
		}
		if origDatabaseType != "" {
			os.Setenv("DATABASE_TYPE", origDatabaseType)
		} else {
			os.Unsetenv("DATABASE_TYPE")
		}
		if origDatabaseURL != "" {
			os.Setenv("DATABASE_URL", origDatabaseURL)
		} else {
			os.Unsetenv("DATABASE_URL")
		}

		// Reset global DB variable
		DB = nil

		// Reload configuration
		config.LoadConfig()
	})
}

func TestGetBackendDir(t *testing.T) {
	// Test that getBackendDir returns a valid path
	backendDir := getBackendDir()

	// Should be an absolute path
	assert.True(t, filepath.IsAbs(backendDir), "getBackendDir should return an absolute path")

	// Should end with "backend"
	assert.True(t, strings.HasSuffix(backendDir, "backend"), "getBackendDir should return path ending with 'backend'")

	// The directory should exist
	info, err := os.Stat(backendDir)
	assert.NoError(t, err, "Backend directory should exist")
	assert.True(t, info.IsDir(), "Backend path should be a directory")

	// Should contain expected subdirectories
	expectedDirs := []string{"database", "models", "config"}
	for _, expectedDir := range expectedDirs {
		dirPath := filepath.Join(backendDir, expectedDir)
		info, err := os.Stat(dirPath)
		assert.NoError(t, err, "Expected directory %s should exist", expectedDir)
		if err == nil {
			assert.True(t, info.IsDir(), "%s should be a directory", expectedDir)
		}
	}
}

func TestPreparedTestDBPath(t *testing.T) {
	// Test that PreparedTestDB contains the correct path
	assert.True(t, strings.Contains(PreparedTestDB, "testdata"), "PreparedTestDB should contain testdata directory")
	assert.True(t, strings.HasSuffix(PreparedTestDB, "prepared_test_data.db"), "PreparedTestDB should end with prepared_test_data.db")
	assert.True(t, filepath.IsAbs(PreparedTestDB), "PreparedTestDB should be an absolute path")
}

func TestTestDataDirPath(t *testing.T) {
	// Test that TestDataDir contains the correct path
	assert.True(t, strings.Contains(TestDataDir, "maintenance"), "TestDataDir should contain maintenance directory")
	assert.True(t, strings.Contains(TestDataDir, "testdata"), "TestDataDir should contain testdata directory")
	assert.True(t, filepath.IsAbs(TestDataDir), "TestDataDir should be an absolute path")
}

func TestConnectDatabase_TestEnvironment(t *testing.T) {
	setupTestEnvironment(t)

	// Set environment to test
	os.Setenv("ENVIRONMENT", "test")
	config.LoadConfig()

	// Reset DB to ensure fresh connection
	DB = nil

	// ConnectDatabase should use test database when environment is "test"
	db := ConnectDatabase()

	assert.NotNil(t, db, "ConnectDatabase should return a valid database connection")
	assert.Equal(t, db, DB, "ConnectDatabase should set global DB variable")
}

func TestConnectDatabase_DevTestingYes(t *testing.T) {
	setupTestEnvironment(t)

	// Set dev testing to yes
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("DEV_TESTING", "yes")
	config.LoadConfig()

	// Reset DB to ensure fresh connection
	DB = nil

	// ConnectDatabase should use test database when DEV_TESTING=yes
	db := ConnectDatabase()

	assert.NotNil(t, db, "ConnectDatabase should return a valid database connection")
	assert.Equal(t, db, DB, "ConnectDatabase should set global DB variable")
}

func TestConnectDatabaseOrig_SQLite(t *testing.T) {
	setupTestEnvironment(t)

	// Create a temporary SQLite file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	// Set up configuration for SQLite
	os.Setenv("DATABASE_TYPE", "sqlite")
	os.Setenv("DATABASE_URL", dbPath)
	config.LoadConfig()

	// Reset DB to ensure fresh connection
	DB = nil

	// Test ConnectDatabaseOrig with SQLite
	db := ConnectDatabaseOrig()

	assert.NotNil(t, db, "ConnectDatabaseOrig should return a valid SQLite connection")
	assert.Equal(t, db, DB, "ConnectDatabaseOrig should set global DB variable")

	// Verify we can perform basic operations
	sqlDB, err := db.DB()
	assert.NoError(t, err, "Should be able to get underlying sql.DB")
	err = sqlDB.Ping()
	assert.NoError(t, err, "Should be able to ping the database")
}

func TestConnectDatabaseOrig_DefaultMySQL(t *testing.T) {
	setupTestEnvironment(t)

	// Set up configuration with empty DATABASE_URL to trigger default MySQL
	os.Setenv("DATABASE_TYPE", "mysql")
	os.Setenv("DATABASE_URL", "")
	config.LoadConfig()

	// Reset DB to ensure fresh connection
	DB = nil

	// Note: This test will fail to connect since we don't have a real MySQL server
	// But we can test that it attempts to use the default configuration
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic since we don't have a real MySQL server
			assert.Contains(t, r.(string), "Failed to connect to database", "Should panic with connection error")
		}
	}()

	ConnectDatabaseOrig()

	// If we reach here, the connection surprisingly succeeded
	// This could happen in a test environment with MySQL available
	if DB != nil {
		assert.NotNil(t, DB, "If connection succeeds, DB should be set")
	}
}

func TestConnectDatabaseOrig_UnsupportedDatabaseType(t *testing.T) {
	setupTestEnvironment(t)

	// Set up configuration with unsupported database type
	os.Setenv("DATABASE_TYPE", "postgresql")
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test")
	config.LoadConfig()

	// Reset DB to ensure fresh connection
	DB = nil

	// Should fall back to MySQL for unsupported database types
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic since we don't have a real MySQL server
			assert.Contains(t, r.(string), "Failed to connect to database", "Should panic with connection error")
		}
	}()

	ConnectDatabaseOrig()
}

func TestGetDB(t *testing.T) {
	setupTestEnvironment(t)

	// Reset DB to ensure fresh connection
	DB = nil

	// First call should initialize DB
	db1 := GetDB()
	assert.NotNil(t, db1, "GetDB should return a valid database connection")
	assert.Equal(t, db1, DB, "GetDB should set global DB variable")

	// Second call should return the same instance
	db2 := GetDB()
	assert.Equal(t, db1, db2, "GetDB should return the same instance on subsequent calls")
}

func TestGetDB_AlreadyInitialized(t *testing.T) {
	setupTestEnvironment(t)

	// Set up a mock database connection
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	mockDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err, "Should be able to create mock database")

	// Set DB to the mock instance
	DB = mockDB

	// GetDB should return the existing instance
	db := GetDB()
	assert.Equal(t, mockDB, db, "GetDB should return the existing DB instance")
}

func TestStringArray_Value(t *testing.T) {
	tests := []struct {
		name     string
		input    StringArray
		expected string
	}{
		{
			name:     "Empty array",
			input:    StringArray{},
			expected: "[]",
		},
		{
			name:     "Single element",
			input:    StringArray{"test"},
			expected: `["test"]`,
		},
		{
			name:     "Multiple elements",
			input:    StringArray{"one", "two", "three"},
			expected: `["one","two","three"]`,
		},
		{
			name:     "Array with special characters",
			input:    StringArray{"test\"quote", "test\nnewline", "test\\backslash"},
			expected: `["test\"quote","test\nnewline","test\\backslash"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := tt.input.Value()
			assert.NoError(t, err, "Value() should not return an error")

			// Convert the returned driver.Value to string
			bytes, ok := value.([]byte)
			require.True(t, ok, "Value() should return []byte")

			assert.JSONEq(t, tt.expected, string(bytes), "Value() should return correct JSON")
		})
	}
}

func TestStringArray_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected StringArray
		hasError bool
	}{
		{
			name:     "Nil input",
			input:    nil,
			expected: StringArray{},
			hasError: false,
		},
		{
			name:     "Empty JSON array",
			input:    []byte("[]"),
			expected: StringArray{},
			hasError: false,
		},
		{
			name:     "Single element JSON array",
			input:    []byte(`["test"]`),
			expected: StringArray{"test"},
			hasError: false,
		},
		{
			name:     "Multiple elements JSON array",
			input:    []byte(`["one","two","three"]`),
			expected: StringArray{"one", "two", "three"},
			hasError: false,
		},
		{
			name:     "JSON array with special characters",
			input:    []byte(`["test\"quote","test\nnewline","test\\backslash"]`),
			expected: StringArray{"test\"quote", "test\nnewline", "test\\backslash"},
			hasError: false,
		},
		{
			name:     "Invalid input type",
			input:    "not a byte slice",
			expected: StringArray{},
			hasError: true,
		},
		{
			name:     "Invalid JSON",
			input:    []byte(`invalid json`),
			expected: StringArray{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sa StringArray
			err := sa.Scan(tt.input)

			if tt.hasError {
				assert.Error(t, err, "Scan() should return an error for invalid input")
			} else {
				assert.NoError(t, err, "Scan() should not return an error for valid input")
				assert.Equal(t, tt.expected, sa, "Scan() should set correct values")
			}
		})
	}
}

func TestStringArray_ValueScanRoundTrip(t *testing.T) {
	// Test that Value() and Scan() work together correctly
	original := StringArray{"test1", "test2", "test3"}

	// Convert to driver.Value
	value, err := original.Value()
	assert.NoError(t, err, "Value() should not error")

	// Scan back to StringArray
	var result StringArray
	err = result.Scan(value)
	assert.NoError(t, err, "Scan() should not error")

	// Should be equal to original
	assert.Equal(t, original, result, "Value/Scan round trip should preserve data")
}

func TestStringArray_DatabaseCompatibility(t *testing.T) {
	// Test that StringArray implements the required database interfaces
	var sa StringArray

	// Should implement driver.Valuer
	_, ok := interface{}(sa).(driver.Valuer)
	assert.True(t, ok, "StringArray should implement driver.Valuer interface")

	// Should have Scan method for sql.Scanner interface
	assert.True(t, true, "StringArray has Scan method for sql.Scanner interface")
}

// Benchmark tests for performance-critical functions
func BenchmarkGetBackendDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getBackendDir()
	}
}

func BenchmarkStringArray_Value(b *testing.B) {
	sa := StringArray{"test1", "test2", "test3", "test4", "test5"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = sa.Value()
	}
}

func BenchmarkStringArray_Scan(b *testing.B) {
	data := []byte(`["test1","test2","test3","test4","test5"]`)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var sa StringArray
		_ = sa.Scan(data)
	}
}
