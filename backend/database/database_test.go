package database

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Test model for migration testing
type TestModel struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"size:100;not null"`
	Age  int    `gorm:"default:0"`
}

// Another test model for migration testing
type AnotherTestModel struct {
	ID          uint   `gorm:"primarykey"`
	Description string `gorm:"size:200"`
	Active      bool   `gorm:"default:true"`
}

// setupTestDB creates a test database for migration testing
func setupTestDB(t *testing.T) *gorm.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_migrate.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err, "Should be able to create test database")

	return db
}

// MockDB creates a mock database that can be configured to return errors
type MockDB struct {
	*gorm.DB
	shouldError bool
	errorMsg    string
}

func (m *MockDB) AutoMigrate(dst ...interface{}) error {
	if m.shouldError {
		return errors.New(m.errorMsg)
	}
	return m.DB.AutoMigrate(dst...)
}

func TestMigrateStructure_WithProvidedDB(t *testing.T) {
	// Create a test database
	testDB := setupTestDB(t)

	// Test migration with provided database
	err := MigrateStructure(testDB)
	assert.NoError(t, err, "MigrateStructure should succeed with valid provided database")

	// Verify that the database connection is working
	sqlDB, err := testDB.DB()
	assert.NoError(t, err, "Should be able to get underlying sql.DB")
	err = sqlDB.Ping()
	assert.NoError(t, err, "Should be able to ping the database")
}

func TestMigrateStructure_WithGlobalDB(t *testing.T) {
	// Save original global DB
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set up global test database
	DB = setupTestDB(t)

	// Test migration without providing database (should use global DB)
	err := MigrateStructure()
	assert.NoError(t, err, "MigrateStructure should succeed with global database")

	// Verify that the global database connection is working
	sqlDB, err := DB.DB()
	assert.NoError(t, err, "Should be able to get underlying sql.DB from global DB")
	err = sqlDB.Ping()
	assert.NoError(t, err, "Should be able to ping the global database")
}

func TestMigrateStructure_WithNilProvidedDB(t *testing.T) {
	// Save original global DB
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set up global test database
	DB = setupTestDB(t)

	// Test migration with nil provided database (should fall back to global DB)
	err := MigrateStructure(nil)
	assert.NoError(t, err, "MigrateStructure should succeed when nil is provided and global DB exists")
}

func TestMigrateStructure_WithMultipleDBParams(t *testing.T) {
	// Create two test databases
	testDB1 := setupTestDB(t)
	testDB2 := setupTestDB(t)

	// Test that it uses the first provided database
	err := MigrateStructure(testDB1, testDB2)
	assert.NoError(t, err, "MigrateStructure should succeed with multiple DB parameters")

	// The function should have used testDB1 (the first parameter)
	// We can't directly verify this, but the test should pass if the first DB is valid
}

func TestMigrateStructure_WithActualModels(t *testing.T) {
	// Create a test database
	testDB := setupTestDB(t)

	// Register models for migration
	err := testDB.AutoMigrate(&TestModel{}, &AnotherTestModel{})
	require.NoError(t, err, "Should be able to register test models")

	// Test migration (AutoMigrate is called without models, which is valid in GORM)
	err = MigrateStructure(testDB)
	assert.NoError(t, err, "MigrateStructure should succeed with registered models")

	// Verify tables exist by checking if we can create records
	testRecord := TestModel{Name: "Test", Age: 25}
	err = testDB.Create(&testRecord).Error
	assert.NoError(t, err, "Should be able to create test record after migration")

	anotherRecord := AnotherTestModel{Description: "Test Description", Active: true}
	err = testDB.Create(&anotherRecord).Error
	assert.NoError(t, err, "Should be able to create another test record after migration")
}

func TestMigrateStructure_ErrorHandling(t *testing.T) {
	// Save original global DB
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set global DB to nil to force an error scenario
	DB = nil

	// Test migration without providing database when global DB is nil
	// This should cause a panic when trying to call AutoMigrate on nil
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic when trying to call methods on nil database
			assert.NotNil(t, r, "Should panic when DB is nil")
		}
	}()

	// This should panic because targetDB will be nil and we can't call AutoMigrate on nil
	MigrateStructure()

	// If we reach here without panic, the test should fail
	t.Fatal("Expected panic when calling MigrateStructure with nil DB")
}

func TestMigrateStructure_EmptyMigration(t *testing.T) {
	// Test that migration works even when there are no models to migrate
	testDB := setupTestDB(t)

	// Call AutoMigrate with no models (empty migration)
	err := testDB.AutoMigrate()
	assert.NoError(t, err, "AutoMigrate should succeed with no models")

	// Test our MigrateStructure function
	err = MigrateStructure(testDB)
	assert.NoError(t, err, "MigrateStructure should succeed with empty migration")
}

func TestMigrateStructure_DatabaseTypes(t *testing.T) {
	// Test with SQLite (in-memory)
	db1, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Should be able to create in-memory SQLite database")

	err = MigrateStructure(db1)
	assert.NoError(t, err, "MigrateStructure should work with in-memory SQLite")

	// Test with SQLite (file-based)
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "file_test.db")
	db2, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err, "Should be able to create file-based SQLite database")

	err = MigrateStructure(db2)
	assert.NoError(t, err, "MigrateStructure should work with file-based SQLite")
}

func TestMigrateStructure_ParameterVariations(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) []*gorm.DB
		expectError bool
		description string
	}{
		{
			name: "No parameters",
			setupFunc: func(t *testing.T) []*gorm.DB {
				// Set up global DB
				originalDB := DB
				t.Cleanup(func() { DB = originalDB })
				DB = setupTestDB(t)
				return []*gorm.DB{}
			},
			expectError: false,
			description: "Should use global DB when no parameters provided",
		},
		{
			name: "One valid parameter",
			setupFunc: func(t *testing.T) []*gorm.DB {
				return []*gorm.DB{setupTestDB(t)}
			},
			expectError: false,
			description: "Should use provided DB when one parameter given",
		},
		{
			name: "Multiple parameters",
			setupFunc: func(t *testing.T) []*gorm.DB {
				return []*gorm.DB{setupTestDB(t), setupTestDB(t)}
			},
			expectError: false,
			description: "Should use first DB when multiple parameters given",
		},
		{
			name: "Nil first parameter with global DB",
			setupFunc: func(t *testing.T) []*gorm.DB {
				originalDB := DB
				t.Cleanup(func() { DB = originalDB })
				DB = setupTestDB(t)
				return []*gorm.DB{nil}
			},
			expectError: false,
			description: "Should fall back to global DB when first parameter is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbs := tt.setupFunc(t)

			err := MigrateStructure(dbs...)

			if tt.expectError {
				assert.Error(t, err, tt.description)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkMigrateStructure(b *testing.B) {
	// Set up test database
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "benchmark.db")
	testDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := MigrateStructure(testDB)
		if err != nil {
			b.Fatalf("MigrateStructure failed: %v", err)
		}
	}
}

func BenchmarkMigrateStructure_GlobalDB(b *testing.B) {
	// Save original global DB
	originalDB := DB
	defer func() {
		DB = originalDB
	}()

	// Set up global test database
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "benchmark_global.db")
	testDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}
	DB = testDB

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := MigrateStructure()
		if err != nil {
			b.Fatalf("MigrateStructure failed: %v", err)
		}
	}
}
