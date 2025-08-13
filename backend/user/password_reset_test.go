package user

import (
	"bamort/database"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestUser(t *testing.T) User {
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

func TestRequestPasswordReset_Success(t *testing.T) {
	user := setupTestUser(t)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/password-reset/request", RequestPasswordReset)

	// Test data
	requestData := map[string]interface{}{
		"email":        user.Email,
		"redirect_url": "http://localhost:3000",
	}
	jsonData, _ := json.Marshal(requestData)

	// Create request
	req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Contains(t, response["message"], "Falls ein Account mit dieser E-Mail-Adresse existiert")

	// Check that user has reset hash set in database
	var dbUser User
	err = dbUser.FindByEmail(user.Email)
	require.NoError(t, err)

	assert.NotNil(t, dbUser.ResetPwHash, "Reset hash should be set")
	assert.NotNil(t, dbUser.ResetPwHashExpires, "Reset expiry should be set")
	assert.True(t, dbUser.ResetPwHashExpires.After(time.Now()), "Reset expiry should be in future")
}

func TestRequestPasswordReset_WithoutRedirectURL(t *testing.T) {
	user := setupTestUser(t)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/password-reset/request", RequestPasswordReset)

	// Test data without redirect_url
	requestData := map[string]interface{}{
		"email": user.Email,
	}
	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should still work with fallback URL
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequestPasswordReset_NonExistentEmail(t *testing.T) {
	database.SetupTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/password-reset/request", RequestPasswordReset)

	requestData := map[string]interface{}{
		"email":        "nonexistent@example.com",
		"redirect_url": "http://localhost:3000",
	}
	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return success to prevent email enumeration
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["message"], "Falls ein Account mit dieser E-Mail-Adresse existiert")
}

func TestRequestPasswordReset_InvalidEmail(t *testing.T) {
	database.SetupTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/password-reset/request", RequestPasswordReset)

	requestData := map[string]interface{}{
		"email":        "invalid-email-format",
		"redirect_url": "http://localhost:3000",
	}
	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Should return 400 for invalid email format
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateResetToken_Success(t *testing.T) {
	user := setupTestUser(t)

	// Set reset hash
	resetHash := "test_reset_hash_123456789"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/password-reset/validate/:token", ValidateResetToken)

	req, _ := http.NewRequest("GET", "/password-reset/validate/"+resetHash, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["valid"].(bool))
	assert.Equal(t, user.Username, response["username"])
	assert.NotNil(t, response["expires"])
}

func TestValidateResetToken_InvalidToken(t *testing.T) {
	database.SetupTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/password-reset/validate/:token", ValidateResetToken)

	req, _ := http.NewRequest("GET", "/password-reset/validate/invalid_token", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestValidateResetToken_ExpiredToken(t *testing.T) {
	user := setupTestUser(t)

	// Set expired reset hash
	resetHash := "expired_reset_hash_123456789"
	expiredTime := time.Now().Add(-1 * time.Hour) // 1 hour ago
	user.ResetPwHash = &resetHash
	user.ResetPwHashExpires = &expiredTime
	err := user.Save()
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/password-reset/validate/:token", ValidateResetToken)

	req, _ := http.NewRequest("GET", "/password-reset/validate/"+resetHash, nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_Success(t *testing.T) {
	user := setupTestUser(t)
	originalPassword := user.PasswordHash

	// Set reset hash
	resetHash := "test_reset_hash_for_password_change"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/password-reset/reset", ResetPassword)

	requestData := map[string]interface{}{
		"token":        resetHash,
		"new_password": "new_secure_password123",
	}
	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Passwort erfolgreich zurückgesetzt", response["message"])

	// Verify password was changed and reset hash was cleared
	var dbUser User
	err = dbUser.FindByEmail(user.Email)
	require.NoError(t, err)

	assert.NotEqual(t, originalPassword, dbUser.PasswordHash, "Password should have changed")
	assert.Nil(t, dbUser.ResetPwHash, "Reset hash should be cleared")
	assert.Nil(t, dbUser.ResetPwHashExpires, "Reset expiry should be cleared")

	// Verify new password hash
	expectedHash := md5.Sum([]byte("new_secure_password123"))
	expectedHashString := hex.EncodeToString(expectedHash[:])
	assert.Equal(t, expectedHashString, dbUser.PasswordHash, "New password hash should match")
}

func TestResetPassword_InvalidToken(t *testing.T) {
	database.SetupTestDB()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/password-reset/reset", ResetPassword)

	requestData := map[string]interface{}{
		"token":        "invalid_token",
		"new_password": "new_secure_password123",
	}
	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_ShortPassword(t *testing.T) {
	user := setupTestUser(t)

	resetHash := "test_reset_hash_short_password"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/password-reset/reset", ResetPassword)

	requestData := map[string]interface{}{
		"token":        resetHash,
		"new_password": "123", // Too short
	}
	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	user := setupTestUser(t)

	// Set expired reset hash
	resetHash := "expired_reset_hash_for_reset"
	expiredTime := time.Now().Add(-1 * time.Hour) // 1 hour ago
	user.ResetPwHash = &resetHash
	user.ResetPwHashExpires = &expiredTime
	err := user.Save()
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/password-reset/reset", ResetPassword)

	requestData := map[string]interface{}{
		"token":        resetHash,
		"new_password": "new_secure_password123",
	}
	jsonData, _ := json.Marshal(requestData)

	req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test User model methods
func TestUser_SetPasswordResetHash(t *testing.T) {
	user := setupTestUser(t)

	resetHash := "test_hash_123456789"
	err := user.SetPasswordResetHash(resetHash)

	assert.NoError(t, err)
	assert.NotNil(t, user.ResetPwHash)
	assert.Equal(t, resetHash, *user.ResetPwHash)
	assert.NotNil(t, user.ResetPwHashExpires)
	assert.True(t, user.ResetPwHashExpires.After(time.Now()))
	assert.True(t, user.ResetPwHashExpires.Before(time.Now().Add(15*24*time.Hour))) // Should be ~14 days
}

func TestUser_ClearPasswordResetHash(t *testing.T) {
	user := setupTestUser(t)

	// First set a reset hash
	resetHash := "test_hash_to_clear"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err)
	require.NotNil(t, user.ResetPwHash)

	// Then clear it
	err = user.ClearPasswordResetHash()

	assert.NoError(t, err)
	assert.Nil(t, user.ResetPwHash)
	assert.Nil(t, user.ResetPwHashExpires)
}

func TestUser_IsResetHashValid(t *testing.T) {
	user := setupTestUser(t)

	resetHash := "valid_test_hash_123"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err)

	// Test valid hash
	assert.True(t, user.IsResetHashValid(resetHash))

	// Test invalid hash
	assert.False(t, user.IsResetHashValid("wrong_hash"))

	// Test expired hash
	expiredTime := time.Now().Add(-1 * time.Hour)
	user.ResetPwHashExpires = &expiredTime
	assert.False(t, user.IsResetHashValid(resetHash))

	// Test nil hash
	user.ResetPwHash = nil
	user.ResetPwHashExpires = nil
	assert.False(t, user.IsResetHashValid(resetHash))
}

func TestUser_FindByResetHash(t *testing.T) {
	user := setupTestUser(t)

	resetHash := "find_by_hash_test_123"
	err := user.SetPasswordResetHash(resetHash)
	require.NoError(t, err)

	// Test finding valid hash
	var foundUser User
	err = foundUser.FindByResetHash(resetHash)
	assert.NoError(t, err)
	assert.Equal(t, user.UserID, foundUser.UserID)
	assert.Equal(t, user.Email, foundUser.Email)

	// Test finding invalid hash
	var notFoundUser User
	err = notFoundUser.FindByResetHash("invalid_hash")
	assert.Error(t, err)

	// Test finding expired hash
	expiredTime := time.Now().Add(-1 * time.Hour)
	user.ResetPwHashExpires = &expiredTime
	err = user.Save()
	require.NoError(t, err)

	var expiredUser User
	err = expiredUser.FindByResetHash(resetHash)
	assert.Error(t, err) // Should not find expired token
}

func TestUser_FindByEmail(t *testing.T) {
	user := setupTestUser(t)

	// Test finding existing email
	var foundUser User
	err := foundUser.FindByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.UserID, foundUser.UserID)
	assert.Equal(t, user.Username, foundUser.Username)

	// Test finding non-existent email
	var notFoundUser User
	err = notFoundUser.FindByEmail("nonexistent@example.com")
	assert.Error(t, err)
}
