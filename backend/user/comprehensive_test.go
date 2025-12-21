package user

import (
	"bamort/database"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestUserWithReset creates a test user for password reset functionality tests
func setupTestUserWithReset(t *testing.T) User {
	database.SetupTestDB()

	// Migrate User structure to ensure reset password fields exist
	err := MigrateStructure()
	require.NoError(t, err, "Failed to migrate user structure")

	// Generate unique email for each test to avoid conflicts
	randomSuffix := rand.Intn(100000)
	user := User{
		Username:     fmt.Sprintf("testuser_reset_%d", randomSuffix),
		PasswordHash: "testpassword123",
		Email:        fmt.Sprintf("test.reset.%d@example.com", randomSuffix),
	}

	// Hash password like in RegisterUser
	hashedPassword := md5.Sum([]byte(user.PasswordHash))
	user.PasswordHash = hex.EncodeToString(hashedPassword[:])

	err = user.Create()
	require.NoError(t, err, "Failed to create test user")

	return user
}

// setupUserModelTestEnvironment sets up the test environment for user model tests
func setupUserModelTestEnvironment(t *testing.T) *gorm.DB {
	// Save original state
	originalDB := database.DB

	// Create a test database
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "user_model_test.db")
	testDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	require.NoError(t, err, "Should be able to create test database")

	// Migrate the User table
	err = testDB.AutoMigrate(&User{})
	require.NoError(t, err, "Should be able to migrate User table")

	// Set global DB to test DB
	database.DB = testDB

	// Cleanup function
	t.Cleanup(func() {
		// Restore original state
		database.DB = originalDB
	})

	return testDB
}

// =============================================================================
// User Model Method Tests - Core CRUD Operations
// =============================================================================

func TestUser_Create(t *testing.T) {
	testDB := setupUserModelTestEnvironment(t)

	user := &User{
		Username:     "testuser",
		PasswordHash: "hashedpassword123",
		Email:        "test@example.com",
	}

	// Test successful creation
	err := user.Create()
	assert.NoError(t, err, "Create should succeed with valid user data")
	assert.NotZero(t, user.UserID, "UserID should be set after creation")
	assert.NotZero(t, user.CreatedAt, "CreatedAt should be set after creation")
	assert.NotZero(t, user.UpdatedAt, "UpdatedAt should be set after creation")

	// Verify user was saved to database
	var retrievedUser User
	err = testDB.First(&retrievedUser, user.UserID).Error
	assert.NoError(t, err, "Should be able to retrieve created user")
	assert.Equal(t, user.Username, retrievedUser.Username, "Username should match")
	assert.Equal(t, user.Email, retrievedUser.Email, "Email should match")
	assert.Equal(t, user.PasswordHash, retrievedUser.PasswordHash, "PasswordHash should match")
}

func TestUser_Create_DuplicateConstraints(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create first user
	user1 := &User{
		Username:     "uniqueuser",
		PasswordHash: "hash1",
		Email:        "unique@example.com",
	}
	err := user1.Create()
	require.NoError(t, err, "First user creation should succeed")

	// Try to create second user with same username
	user2 := &User{
		Username:     "uniqueuser", // Same username
		PasswordHash: "hash2",
		Email:        "different@example.com",
	}
	err = user2.Create()
	assert.Error(t, err, "Should fail to create user with duplicate username")

	// Try to create third user with same email
	user3 := &User{
		Username:     "differentuser",
		PasswordHash: "hash3",
		Email:        "unique@example.com", // Same email
	}
	err = user3.Create()
	assert.Error(t, err, "Should fail to create user with duplicate email")
}

func TestUser_Create_WithNilDatabase(t *testing.T) {
	// Save original state
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	// Set database to nil
	database.DB = nil

	user := &User{
		Username:     "nildbuser",
		PasswordHash: "hash",
		Email:        "nildb@example.com",
	}

	// Should error when database is nil
	err := user.Create()
	assert.Error(t, err, "Should error when database is nil")
}

func TestUser_First(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create a test user
	originalUser := &User{
		Username:     "findmeuser",
		PasswordHash: "findhash",
		Email:        "findme@example.com",
	}
	err := originalUser.Create()
	require.NoError(t, err, "Should be able to create test user")

	// Test finding existing user
	var foundUser User
	err = foundUser.First("findmeuser")
	assert.NoError(t, err, "Should find existing user")
	assert.Equal(t, originalUser.UserID, foundUser.UserID, "Found user should have same ID")
	assert.Equal(t, originalUser.Username, foundUser.Username, "Found user should have same username")
	assert.Equal(t, originalUser.Email, foundUser.Email, "Found user should have same email")

	// Test finding non-existent user
	var notFoundUser User
	err = notFoundUser.First("nonexistentuser")
	assert.Error(t, err, "Should return error for non-existent user")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Should return RecordNotFound error")

	// Test with empty username
	var emptyUser User
	err = emptyUser.First("")
	assert.Error(t, err, "Should return error for empty username")
}

func TestUser_First_WithNilDatabase(t *testing.T) {
	// Save original state
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	// Set database to nil
	database.DB = nil

	var user User
	err := user.First("anyuser")
	assert.Error(t, err, "Should error when database is nil")
}

func TestUser_FirstId(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create a test user
	originalUser := &User{
		Username:     "findbyiduser",
		PasswordHash: "idhash",
		Email:        "findbyid@example.com",
	}
	err := originalUser.Create()
	require.NoError(t, err, "Should be able to create test user")

	// Test finding existing user by ID
	var foundUser User
	err = foundUser.FirstId(originalUser.UserID)
	assert.NoError(t, err, "Should find existing user by ID")
	assert.Equal(t, originalUser.UserID, foundUser.UserID, "Found user should have same ID")
	assert.Equal(t, originalUser.Username, foundUser.Username, "Found user should have same username")
	assert.Equal(t, originalUser.Email, foundUser.Email, "Found user should have same email")
	assert.Equal(t, originalUser.PasswordHash, foundUser.PasswordHash, "Found user should have same password hash")

	// Test finding non-existent user by ID
	var notFoundUser User
	err = notFoundUser.FirstId(99999) // Non-existent ID
	assert.Error(t, err, "Should return error for non-existent user ID")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Should return RecordNotFound error")

	// Test with ID 0 (should be invalid)
	var zeroIdUser User
	err = zeroIdUser.FirstId(0)
	assert.Error(t, err, "Should return error for ID 0")

	// Test with very large ID
	var largeIdUser User
	err = largeIdUser.FirstId(4294967295) // Max uint32
	assert.Error(t, err, "Should return error for very large non-existent ID")
}

func TestUser_FirstId_WithNilDatabase(t *testing.T) {
	// Save original state
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	// Set database to nil
	database.DB = nil

	var user User
	err := user.FirstId(1)
	assert.Error(t, err, "Should error when database is nil")
}

func TestUser_Save(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create a test user
	user := &User{
		Username:     "saveuser",
		PasswordHash: "savehash",
		Email:        "save@example.com",
	}
	err := user.Create()
	require.NoError(t, err, "Should be able to create test user")

	originalID := user.UserID
	originalCreatedAt := user.CreatedAt

	// Modify the user
	user.Username = "modifieduser"
	user.Email = "modified@example.com"
	user.PasswordHash = "modifiedhash"

	// Test saving changes
	err = user.Save()
	assert.NoError(t, err, "Save should succeed")

	// Verify changes were saved to database
	var retrievedUser User
	err = retrievedUser.FirstId(user.UserID)
	assert.NoError(t, err, "Should be able to retrieve saved user")
	assert.Equal(t, "modifieduser", retrievedUser.Username, "Username should be updated")
	assert.Equal(t, "modified@example.com", retrievedUser.Email, "Email should be updated")
	assert.Equal(t, "modifiedhash", retrievedUser.PasswordHash, "PasswordHash should be updated")

	// Verify ID and CreatedAt remain unchanged
	assert.Equal(t, originalID, retrievedUser.UserID, "UserID should remain unchanged")
	assert.Equal(t, originalCreatedAt.Unix(), retrievedUser.CreatedAt.Unix(), "CreatedAt should remain unchanged")

	// Verify UpdatedAt was changed
	assert.True(t, retrievedUser.UpdatedAt.After(retrievedUser.CreatedAt), "UpdatedAt should be after CreatedAt")
}

func TestUser_Save_UniqueConstraintViolation(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create two test users
	user1 := &User{
		Username:     "user1",
		PasswordHash: "hash1",
		Email:        "user1@example.com",
	}
	err := user1.Create()
	require.NoError(t, err, "Should be able to create first user")

	user2 := &User{
		Username:     "user2",
		PasswordHash: "hash2",
		Email:        "user2@example.com",
	}
	err = user2.Create()
	require.NoError(t, err, "Should be able to create second user")

	// Try to update user2 to have same username as user1
	user2.Username = "user1"
	err = user2.Save()
	assert.Error(t, err, "Save should fail due to unique constraint violation on username")

	// Try to update user2 to have same email as user1
	user2.Username = "user2" // Reset username
	user2.Email = "user1@example.com"
	err = user2.Save()
	assert.Error(t, err, "Save should fail due to unique constraint violation on email")
}

func TestUser_Save_WithNilDatabase(t *testing.T) {
	// Save original state
	originalDB := database.DB
	defer func() {
		database.DB = originalDB
	}()

	// Set database to nil
	database.DB = nil

	user := &User{
		UserID:       1,
		Username:     "nildbuser",
		PasswordHash: "hash",
		Email:        "nildb@example.com",
	}

	// Should error when database is nil
	err := user.Save()
	assert.Error(t, err, "Should error when database is nil")
}

func TestUser_Save_NewRecord(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create user without using Create() (no ID set)
	user := &User{
		Username:     "newrecorduser",
		PasswordHash: "newrecordhash",
		Email:        "newrecord@example.com",
	}

	// Save should work like Create for new records
	err := user.Save()
	assert.NoError(t, err, "Save should succeed for new record")
	assert.NotZero(t, user.UserID, "UserID should be set after save")
	assert.NotZero(t, user.CreatedAt, "CreatedAt should be set after save")
	assert.NotZero(t, user.UpdatedAt, "UpdatedAt should be set after save")

	// Verify record was saved to database
	var retrievedUser User
	err = retrievedUser.FirstId(user.UserID)
	assert.NoError(t, err, "Should be able to retrieve saved user")
	assert.Equal(t, user.Username, retrievedUser.Username, "Username should match")
	assert.Equal(t, user.Email, retrievedUser.Email, "Email should match")
}

func TestUser_Save_PartialUpdate(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create a test user
	user := &User{
		Username:     "partialuser",
		PasswordHash: "partialhash",
		Email:        "partial@example.com",
	}
	err := user.Create()
	require.NoError(t, err, "Should be able to create test user")

	// Update only username
	originalEmail := user.Email
	originalHash := user.PasswordHash
	user.Username = "updatedpartialuser"

	err = user.Save()
	assert.NoError(t, err, "Partial save should succeed")

	// Verify only username was updated
	var retrievedUser User
	err = retrievedUser.FirstId(user.UserID)
	assert.NoError(t, err, "Should be able to retrieve saved user")
	assert.Equal(t, "updatedpartialuser", retrievedUser.Username, "Username should be updated")
	assert.Equal(t, originalEmail, retrievedUser.Email, "Email should remain unchanged")
	assert.Equal(t, originalHash, retrievedUser.PasswordHash, "PasswordHash should remain unchanged")
}

// =============================================================================
// Password Reset Method Tests
// =============================================================================

func TestUser_FindByEmail(t *testing.T) {
	user := setupTestUserWithReset(t)

	// Test finding existing email
	var foundUser User
	err := foundUser.FindByEmail(user.Email)
	assert.NoError(t, err, "Should find user by existing email")
	assert.Equal(t, user.UserID, foundUser.UserID, "Found user should have correct ID")
	assert.Equal(t, user.Username, foundUser.Username, "Found user should have correct username")

	// Test finding non-existent email
	var notFoundUser User
	err = notFoundUser.FindByEmail("nonexistent@example.com")
	assert.Error(t, err, "Should return error for non-existent email")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Should return RecordNotFound error")

	// Test with empty email
	var emptyEmailUser User
	err = emptyEmailUser.FindByEmail("")
	assert.Error(t, err, "Should return error for empty email")
}

func TestUser_FindByResetHash(t *testing.T) {
	user := setupTestUserWithReset(t)

	resetHash := "find_by_hash_test_123"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err, "Should be able to set reset hash")

	// Test finding valid hash
	var foundUser User
	err = foundUser.FindByResetHash(resetHash)
	assert.NoError(t, err, "Should find user by valid reset hash")
	assert.Equal(t, user.UserID, foundUser.UserID, "Found user should have correct ID")
	assert.Equal(t, user.Email, foundUser.Email, "Found user should have correct email")

	// Test finding invalid hash
	var notFoundUser User
	err = notFoundUser.FindByResetHash("invalid_hash")
	assert.Error(t, err, "Should return error for invalid reset hash")

	// Test finding expired hash
	expiredTime := time.Now().Add(-1 * time.Hour)
	user.ResetPwHashExpires = &expiredTime
	err = user.Save()
	require.NoError(t, err, "Should be able to save expired hash")

	var expiredUser User
	err = expiredUser.FindByResetHash(resetHash)
	assert.Error(t, err, "Should not find expired reset hash")

	// Test with empty hash
	var emptyHashUser User
	err = emptyHashUser.FindByResetHash("")
	assert.Error(t, err, "Should return error for empty hash")
}

func TestUser_SetPasswordResetHash(t *testing.T) {
	user := setupTestUserWithReset(t)

	resetHash := "test_hash_123456789"
	beforeTime := time.Now()

	err := user.SetPasswordResetHash(resetHash)
	afterTime := time.Now()

	assert.NoError(t, err, "Should successfully set reset hash")
	assert.NotNil(t, user.ResetPwHash, "ResetPwHash should be set")
	assert.Equal(t, resetHash, *user.ResetPwHash, "ResetPwHash should match provided value")
	assert.NotNil(t, user.ResetPwHashExpires, "ResetPwHashExpires should be set")

	// Verify expiry time is approximately 14 days from now
	expectedMinExpiry := beforeTime.Add(14 * 24 * time.Hour)
	expectedMaxExpiry := afterTime.Add(14 * 24 * time.Hour)
	assert.True(t, user.ResetPwHashExpires.After(expectedMinExpiry), "Expiry should be after minimum expected time")
	assert.True(t, user.ResetPwHashExpires.Before(expectedMaxExpiry), "Expiry should be before maximum expected time")

	// Test setting multiple times (should overwrite)
	newResetHash := "new_reset_hash_456"
	err = user.SetPasswordResetHash(newResetHash)
	assert.NoError(t, err, "Should successfully overwrite reset hash")
	assert.Equal(t, newResetHash, *user.ResetPwHash, "ResetPwHash should be updated to new value")
}

func TestUser_ClearPasswordResetHash(t *testing.T) {
	user := setupTestUserWithReset(t)

	// First set a reset hash
	resetHash := "test_hash_to_clear"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err, "Should be able to set reset hash")
	require.NotNil(t, user.ResetPwHash, "Reset hash should be set before clearing")

	// Then clear it
	err = user.ClearPasswordResetHash()
	assert.NoError(t, err, "Should successfully clear reset hash")
	assert.Nil(t, user.ResetPwHash, "ResetPwHash should be nil after clearing")
	assert.Nil(t, user.ResetPwHashExpires, "ResetPwHashExpires should be nil after clearing")

	// Verify changes were saved to database
	var retrievedUser User
	err = retrievedUser.FirstId(user.UserID)
	assert.NoError(t, err, "Should be able to retrieve user after clearing")
	assert.Nil(t, retrievedUser.ResetPwHash, "Retrieved user should have nil reset hash")
	assert.Nil(t, retrievedUser.ResetPwHashExpires, "Retrieved user should have nil expiry")

	// Test clearing when already cleared (should not error)
	err = user.ClearPasswordResetHash()
	assert.NoError(t, err, "Should not error when clearing already cleared hash")
}

func TestUser_IsResetHashValid(t *testing.T) {
	user := setupTestUserWithReset(t)

	resetHash := "valid_test_hash_123"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err, "Should be able to set reset hash")

	// Test valid hash
	assert.True(t, user.IsResetHashValid(resetHash), "Should be valid with correct hash and future expiry")

	// Test invalid hash
	assert.False(t, user.IsResetHashValid("wrong_hash"), "Should be invalid with wrong hash")

	// Test expired hash
	expiredTime := time.Now().Add(-1 * time.Hour)
	user.ResetPwHashExpires = &expiredTime
	assert.False(t, user.IsResetHashValid(resetHash), "Should be invalid with expired hash")

	// Test nil hash
	user.ResetPwHash = nil
	user.ResetPwHashExpires = nil
	assert.False(t, user.IsResetHashValid(resetHash), "Should be invalid when hash is nil")

	// Test only hash set (no expiry)
	user.ResetPwHash = &resetHash
	user.ResetPwHashExpires = nil
	assert.False(t, user.IsResetHashValid(resetHash), "Should be invalid when expiry is nil")

	// Test only expiry set (no hash)
	user.ResetPwHash = nil
	futureTime := time.Now().Add(1 * time.Hour)
	user.ResetPwHashExpires = &futureTime
	assert.False(t, user.IsResetHashValid(resetHash), "Should be invalid when hash is nil")

	// Test with empty string hash
	user.ResetPwHash = &resetHash
	user.ResetPwHashExpires = &futureTime
	assert.False(t, user.IsResetHashValid(""), "Should be invalid with empty hash string")
}

// =============================================================================
// Legacy Tests from original user_test.go (Fixed for isolated testing)
// =============================================================================

func TestRegisterUser(t *testing.T) {
	setupUserModelTestEnvironment(t)

	usr := User{
		Username:     "bebe",
		PasswordHash: "osiris",
		Email:        "frank@wuenscheonline.de",
	}

	hashedPassword := md5.Sum([]byte(usr.PasswordHash))
	usr.PasswordHash = hex.EncodeToString(hashedPassword[:])
	err := usr.Create()
	assert.NoError(t, err, "no error expected when creating record")

	usr2 := User{
		Username:     "bubnu",
		PasswordHash: "osiris",
		Email:        "spacer@wuenscheonline.de",
	}
	hashedPassword = md5.Sum([]byte(usr2.PasswordHash))
	usr2.PasswordHash = hex.EncodeToString(hashedPassword[:])
	err = usr2.Create()
	assert.NoError(t, err, "no error expected when creating record")
}

func TestLoginUser(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create test user first
	usr := User{
		Username:     "logintest",
		PasswordHash: "osiris",
		Email:        "login@test.com",
	}
	hashedPassword := md5.Sum([]byte(usr.PasswordHash))
	usr.PasswordHash = hex.EncodeToString(hashedPassword[:])
	err := usr.Create()
	require.NoError(t, err, "Should create test user")

	// Test login functionality
	var foundUser User
	input := struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		HashedPassword string
	}{
		Username: "logintest",
		Password: "osiris",
	}
	err = foundUser.First(input.Username)
	assert.NoError(t, err, "no error expected when finding record")

	hashedPassword = md5.Sum([]byte(input.Password))
	input.HashedPassword = hex.EncodeToString(hashedPassword[:])
	assert.Equal(t, input.HashedPassword, foundUser.PasswordHash, "Password hashes should match")
}

func TestHashing(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create users similar to TestRegisterUser
	usr := User{
		Username:     "hashtest1",
		PasswordHash: "osiris",
		Email:        "hash1@test.com",
	}
	hashedPassword := md5.Sum([]byte(usr.PasswordHash))
	usr.PasswordHash = hex.EncodeToString(hashedPassword[:])
	err := usr.Create()
	require.NoError(t, err, "Should create first user")

	usr2 := User{
		Username:     "hashtest2",
		PasswordHash: "osiris",
		Email:        "hash2@test.com",
	}
	hashedPassword = md5.Sum([]byte(usr2.PasswordHash))
	usr2.PasswordHash = hex.EncodeToString(hashedPassword[:])
	err = usr2.Create()
	require.NoError(t, err, "Should create second user")

	// Test first user hash
	var foundUser1 User
	input1 := struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		HashedPassword string
	}{
		Username: "hashtest1",
		Password: "osiris",
	}
	err = foundUser1.First(input1.Username)
	assert.NoError(t, err, "no error expected when finding record")

	hashedPassword = md5.Sum([]byte(input1.Password))
	input1.HashedPassword = hex.EncodeToString(hashedPassword[:])
	assert.Equal(t, input1.HashedPassword, foundUser1.PasswordHash)

	// Test second user hash
	var foundUser2 User
	input2 := struct {
		Username       string `json:"username"`
		Password       string `json:"password"`
		HashedPassword string
	}{
		Username: "hashtest2",
		Password: "osiris",
	}
	err = foundUser2.First(input2.Username)
	assert.NoError(t, err, "no error expected when finding record")

	hashedPassword = md5.Sum([]byte(input2.Password))
	input2.HashedPassword = hex.EncodeToString(hashedPassword[:])
	assert.Equal(t, input2.HashedPassword, foundUser2.PasswordHash)
}

// =============================================================================
// Edge Cases and Error Handling Tests
// =============================================================================

func TestUser_EdgeCases(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Test with empty strings
	user := &User{
		Username:     "",
		PasswordHash: "",
		Email:        "",
	}
	err := user.Save()
	assert.NoError(t, err, "Should save user with empty strings")
	assert.NotZero(t, user.UserID, "UserID should be set")

	// Test FirstId with the created user
	var foundUser User
	err = foundUser.FirstId(user.UserID)
	assert.NoError(t, err, "Should find user with empty strings")
	assert.Equal(t, "", foundUser.Username, "Username should be empty")
	assert.Equal(t, "", foundUser.Email, "Email should be empty")

	// Test with very long strings
	longString := strings.Repeat("a", 500)
	user2 := &User{
		Username:     "longuser",
		PasswordHash: longString,
		Email:        "long@example.com",
	}
	err = user2.Create()
	assert.NoError(t, err, "Should create user with long password hash")
}

func TestUser_ConcurrentAccess(t *testing.T) {
	setupUserModelTestEnvironment(t)

	// Create a test user
	user := &User{
		Username:     "concurrentuser",
		PasswordHash: "hash",
		Email:        "concurrent@example.com",
	}
	err := user.Create()
	require.NoError(t, err, "Should be able to create test user")

	// Test concurrent modifications
	done := make(chan bool, 2)

	// Goroutine 1: Update username
	go func() {
		defer func() { done <- true }()
		var u1 User
		u1.FirstId(user.UserID)
		u1.Username = "concurrent1"
		u1.Save()
	}()

	// Goroutine 2: Update email
	go func() {
		defer func() { done <- true }()
		var u2 User
		u2.FirstId(user.UserID)
		u2.Email = "concurrent2@example.com"
		u2.Save()
	}()

	// Wait for both to complete
	<-done
	<-done

	// Verify final state (one of the updates should have won)
	var finalUser User
	err = finalUser.FirstId(user.UserID)
	assert.NoError(t, err, "Should be able to retrieve user after concurrent updates")
	// We can't predict exact final state due to race conditions, but should not crash
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkUser_Create(b *testing.B) {
	// Setup test database
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "benchmark_create.db")
	testDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}

	err = testDB.AutoMigrate(&User{})
	if err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		user := &User{
			Username:     fmt.Sprintf("benchuser%d", i),
			PasswordHash: "benchhash",
			Email:        fmt.Sprintf("bench%d@example.com", i),
		}
		err := user.Create()
		if err != nil {
			b.Fatalf("Create failed: %v", err)
		}
	}
}

func BenchmarkUser_First(b *testing.B) {
	// Setup test database with users
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "benchmark_first.db")
	testDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}

	err = testDB.AutoMigrate(&User{})
	if err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	// Create test users
	for i := 0; i < 100; i++ {
		user := &User{
			Username:     fmt.Sprintf("finduser%d", i),
			PasswordHash: "findhash",
			Email:        fmt.Sprintf("find%d@example.com", i),
		}
		user.Create()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var user User
		username := fmt.Sprintf("finduser%d", i%100)
		err := user.First(username)
		if err != nil {
			b.Fatalf("First failed: %v", err)
		}
	}
}

func BenchmarkUser_Save(b *testing.B) {
	// Setup test database with a user
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "benchmark_save.db")
	testDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}

	err = testDB.AutoMigrate(&User{})
	if err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	// Create a test user
	user := &User{
		Username:     "saveuser",
		PasswordHash: "savehash",
		Email:        "save@example.com",
	}
	user.Create()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		user.Username = fmt.Sprintf("saveuser%d", i)
		err := user.Save()
		if err != nil {
			b.Fatalf("Save failed: %v", err)
		}
	}
}

func BenchmarkUser_FirstId(b *testing.B) {
	// Setup test database with users
	tempDir := b.TempDir()
	dbPath := filepath.Join(tempDir, "benchmark_firstid.db")
	testDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		b.Fatalf("Failed to create test database: %v", err)
	}

	err = testDB.AutoMigrate(&User{})
	if err != nil {
		b.Fatalf("Failed to migrate: %v", err)
	}

	originalDB := database.DB
	database.DB = testDB
	defer func() { database.DB = originalDB }()

	// Create test users and collect their IDs
	var userIDs []uint
	for i := 0; i < 100; i++ {
		user := &User{
			Username:     fmt.Sprintf("finduser%d", i),
			PasswordHash: "findhash",
			Email:        fmt.Sprintf("find%d@example.com", i),
		}
		user.Create()
		userIDs = append(userIDs, user.UserID)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var user User
		userID := userIDs[i%len(userIDs)]
		err := user.FirstId(userID)
		if err != nil {
			b.Fatalf("FirstId failed: %v", err)
		}
	}
}
