package user

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bamort/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnvironment sets up the test environment
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

// setupHandlerTestEnvironment sets up the test environment for handler tests
func setupHandlerTestEnvironment(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()
	err := MigrateStructure()
	require.NoError(t, err, "Should migrate user structure")
	gin.SetMode(gin.TestMode)
}

// createTestUser creates a test user and returns it
func createTestUser(t *testing.T, username, password, email string) *User {
	user := &User{
		Username:     username,
		PasswordHash: password,
		Email:        email,
	}

	hashedPassword := md5.Sum([]byte(password))
	user.PasswordHash = hex.EncodeToString(hashedPassword[:])

	err := user.Create()
	require.NoError(t, err, "Should create test user")

	return user
}

// TestGetUserProfile tests the GetUserProfile handler
func TestGetUserProfile(t *testing.T) {
	setupHandlerTestEnvironment(t)

	t.Run("Success - Get user profile", func(t *testing.T) {
		// Create test user
		testUser := createTestUser(t, "profileuser", "password123", "profile@test.com")

		// Create HTTP request
		req, _ := http.NewRequest("GET", "/api/user/profile", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set user ID in context (simulating AuthMiddleware)
		c.Set("userID", testUser.UserID)

		// Call handler
		GetUserProfile(c)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, testUser.Username, response["username"])
		assert.Equal(t, testUser.Email, response["email"])
		assert.Equal(t, float64(testUser.UserID), response["id"])
	})

	t.Run("Failure - No user ID in context", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user/profile", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Don't set userID in context

		GetUserProfile(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized", response["error"])
	})

	t.Run("Failure - User not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/user/profile", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set non-existent user ID
		c.Set("userID", uint(99999))

		GetUserProfile(c)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response["error"])
	})
}

// TestUpdateEmail tests the UpdateEmail handler
func TestUpdateEmail(t *testing.T) {
	setupHandlerTestEnvironment(t)

	t.Run("Success - Update email", func(t *testing.T) {
		// Create test user
		testUser := createTestUser(t, "emailuser", "password123", "old@test.com")

		// Create request body
		requestData := map[string]interface{}{
			"email": "new@test.com",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/email", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdateEmail(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Email updated successfully", response["message"])
		assert.Equal(t, "new@test.com", response["email"])

		// Verify email was actually updated in database
		var updatedUser User
		err = updatedUser.FirstId(testUser.UserID)
		assert.NoError(t, err)
		assert.Equal(t, "new@test.com", updatedUser.Email)
	})

	t.Run("Failure - Email already in use", func(t *testing.T) {
		// Create two test users
		testUser1 := createTestUser(t, "emailuser1", "password123", "user1@test.com")
		createTestUser(t, "emailuser2", "password123", "user2@test.com")

		// Try to change user1's email to user2's email
		requestData := map[string]interface{}{
			"email": "user2@test.com",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/email", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser1.UserID)

		UpdateEmail(c)

		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Email already in use", response["error"])
	})

	t.Run("Failure - Invalid email format", func(t *testing.T) {
		testUser := createTestUser(t, "emailuser3", "password123", "valid@test.com")

		requestData := map[string]interface{}{
			"email": "invalid-email",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/email", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdateEmail(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failure - No user ID in context", func(t *testing.T) {
		requestData := map[string]interface{}{
			"email": "new@test.com",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/email", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		UpdateEmail(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Success - Update to same email (idempotent)", func(t *testing.T) {
		testUser := createTestUser(t, "emailuser4", "password123", "same@test.com")

		requestData := map[string]interface{}{
			"email": "same@test.com",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/email", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdateEmail(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestUpdatePassword tests the UpdatePassword handler
func TestUpdatePassword(t *testing.T) {
	setupHandlerTestEnvironment(t)

	t.Run("Success - Update password", func(t *testing.T) {
		// Create test user
		testUser := createTestUser(t, "passworduser", "oldpassword123", "password@test.com")

		// Create request body
		requestData := map[string]interface{}{
			"current_password": "oldpassword123",
			"new_password":     "newpassword456",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdatePassword(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Password updated successfully", response["message"])

		// Verify password was actually updated in database
		var updatedUser User
		err = updatedUser.FirstId(testUser.UserID)
		assert.NoError(t, err)

		// Hash the new password and check it matches
		hashedNewPassword := md5.Sum([]byte("newpassword456"))
		expectedHash := hex.EncodeToString(hashedNewPassword[:])
		assert.Equal(t, expectedHash, updatedUser.PasswordHash)
	})

	t.Run("Failure - Incorrect current password", func(t *testing.T) {
		testUser := createTestUser(t, "passworduser2", "correctpassword", "password2@test.com")

		requestData := map[string]interface{}{
			"current_password": "wrongpassword",
			"new_password":     "newpassword456",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdatePassword(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Current password is incorrect", response["error"])
	})

	t.Run("Failure - New password too short", func(t *testing.T) {
		testUser := createTestUser(t, "passworduser3", "oldpassword", "password3@test.com")

		requestData := map[string]interface{}{
			"current_password": "oldpassword",
			"new_password":     "short",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdatePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failure - Missing current password", func(t *testing.T) {
		testUser := createTestUser(t, "passworduser4", "oldpassword", "password4@test.com")

		requestData := map[string]interface{}{
			"new_password": "newpassword456",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdatePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failure - Missing new password", func(t *testing.T) {
		testUser := createTestUser(t, "passworduser5", "oldpassword", "password5@test.com")

		requestData := map[string]interface{}{
			"current_password": "oldpassword",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdatePassword(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Failure - No user ID in context", func(t *testing.T) {
		requestData := map[string]interface{}{
			"current_password": "oldpassword",
			"new_password":     "newpassword456",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		UpdatePassword(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Failure - User not found", func(t *testing.T) {
		requestData := map[string]interface{}{
			"current_password": "oldpassword",
			"new_password":     "newpassword456",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", uint(99999))

		UpdatePassword(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Success - Password can be changed multiple times", func(t *testing.T) {
		testUser := createTestUser(t, fmt.Sprintf("passworduser6_%d", os.Getpid()), "password1", "password6@test.com")

		// First password change
		requestData1 := map[string]interface{}{
			"current_password": "password1",
			"new_password":     "password2",
		}
		requestBody1, _ := json.Marshal(requestData1)
		req1, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody1))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		c1, _ := gin.CreateTestContext(w1)
		c1.Request = req1
		c1.Set("userID", testUser.UserID)

		UpdatePassword(c1)
		assert.Equal(t, http.StatusOK, w1.Code)

		// Second password change
		requestData2 := map[string]interface{}{
			"current_password": "password2",
			"new_password":     "password3",
		}
		requestBody2, _ := json.Marshal(requestData2)
		req2, _ := http.NewRequest("PUT", "/api/user/password", bytes.NewBuffer(requestBody2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		c2, _ := gin.CreateTestContext(w2)
		c2.Request = req2
		c2.Set("userID", testUser.UserID)

		UpdatePassword(c2)
		assert.Equal(t, http.StatusOK, w2.Code)

		// Verify final password
		var updatedUser User
		err := updatedUser.FirstId(testUser.UserID)
		assert.NoError(t, err)
		hashedPassword3 := md5.Sum([]byte("password3"))
		expectedHash := hex.EncodeToString(hashedPassword3[:])
		assert.Equal(t, expectedHash, updatedUser.PasswordHash)
	})
}

// TestUpdateLanguage tests the UpdateLanguage handler
func TestUpdateLanguage(t *testing.T) {
	setupHandlerTestEnvironment(t)

	t.Run("Success - Update language to en", func(t *testing.T) {
		testUser := createTestUser(t, "languser1", "password123", "languser1@test.com")

		requestData := map[string]interface{}{
			"language": "en",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/language", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdateLanguage(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Language updated successfully", response["message"])
		assert.Equal(t, "en", response["language"])

		// Verify language was actually updated in database
		var updatedUser User
		err = updatedUser.FirstId(testUser.UserID)
		assert.NoError(t, err)
		assert.Equal(t, "en", updatedUser.PreferredLanguage)
	})

	t.Run("Success - Update language to de", func(t *testing.T) {
		testUser := createTestUser(t, "languser2", "password123", "languser2@test.com")

		requestData := map[string]interface{}{
			"language": "de",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/language", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdateLanguage(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "de", response["language"])
	})

	t.Run("Failure - Invalid language", func(t *testing.T) {
		testUser := createTestUser(t, "languser3", "password123", "languser3@test.com")

		requestData := map[string]interface{}{
			"language": "fr",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/language", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdateLanguage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid language. Supported languages: de, en", response["error"])
	})

	t.Run("Failure - No user ID in context", func(t *testing.T) {
		requestData := map[string]interface{}{
			"language": "en",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/language", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		UpdateLanguage(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Unauthorized", response["error"])
	})

	t.Run("Failure - Empty language", func(t *testing.T) {
		testUser := createTestUser(t, "languser4", "password123", "languser4@test.com")

		requestData := map[string]interface{}{
			"language": "",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("PUT", "/api/user/language", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		UpdateLanguage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Language is required", response["error"])
	})
}

// TestGetUserProfileWithLanguage tests that GetUserProfile returns the language field
func TestGetUserProfileWithLanguage(t *testing.T) {
	setupHandlerTestEnvironment(t)

	t.Run("Success - Profile includes preferred language", func(t *testing.T) {
		testUser := createTestUser(t, "langprofileuser", "password123", "langprofile@test.com")

		// Set language to en
		testUser.PreferredLanguage = "en"
		err := testUser.Save()
		require.NoError(t, err)

		req, _ := http.NewRequest("GET", "/api/user/profile", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		GetUserProfile(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "en", response["preferred_language"])
	})

	t.Run("Success - New user has default language de", func(t *testing.T) {
		testUser := createTestUser(t, "newlanguser", "password123", "newlang@test.com")

		req, _ := http.NewRequest("GET", "/api/user/profile", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", testUser.UserID)

		GetUserProfile(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check that preferred_language is present and defaults to "de"
		lang, ok := response["preferred_language"]
		assert.True(t, ok, "preferred_language should be present in response")
		if lang == "" || lang == nil {
			// If empty, GORM default should be "de"
			assert.Equal(t, "de", testUser.PreferredLanguage)
		} else {
			assert.Equal(t, "de", lang)
		}
	})
}
