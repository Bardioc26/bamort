package router

import (
	"bamort/database"
	"bamort/user"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	BaseRouterGrp(r)
	return r
}

func setupTestUserForRouter(t *testing.T) user.User {
	database.SetupTestDB()

	// Migrate User structure to ensure reset password fields exist
	err := user.MigrateStructure()
	require.NoError(t, err, "Failed to migrate user structure")

	// Generate unique user for each test
	randomSuffix := rand.Intn(100000)
	testUser := user.User{
		Username:     fmt.Sprintf("routetest_user_%d", randomSuffix),
		PasswordHash: "testpassword123",
		Email:        fmt.Sprintf("routetest.%d@example.com", randomSuffix),
	}

	// Hash password like in RegisterUser
	hashedPassword := md5.Sum([]byte(testUser.PasswordHash))
	testUser.PasswordHash = hex.EncodeToString(hashedPassword[:])

	err = testUser.Create()
	require.NoError(t, err, "Failed to create test user")

	return testUser
}

func TestPasswordResetRoutes_Integration(t *testing.T) {
	testUser := setupTestUserForRouter(t)
	router := setupTestRouter()

	t.Run("Complete Password Reset Flow", func(t *testing.T) {
		// Step 1: Request password reset
		requestData := map[string]interface{}{
			"email":        testUser.Email,
			"redirect_url": "http://localhost:3000",
		}
		jsonData, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "Falls ein Account mit dieser E-Mail-Adresse existiert")

		// Step 2: Get the reset token from database
		var dbUser user.User
		err = dbUser.FindByEmail(testUser.Email)
		require.NoError(t, err)
		require.NotNil(t, dbUser.ResetPwHash, "Reset hash should be set")

		resetToken := *dbUser.ResetPwHash

		// Step 3: Validate the reset token
		req, _ = http.NewRequest("GET", "/password-reset/validate/"+resetToken, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.True(t, response["valid"].(bool))
		assert.Equal(t, testUser.Username, response["username"])

		// Step 4: Reset the password
		resetData := map[string]interface{}{
			"token":        resetToken,
			"new_password": "new_secure_password_123",
		}
		jsonData, _ = json.Marshal(resetData)

		req, _ = http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Passwort erfolgreich zurückgesetzt", response["message"])

		// Step 5: Verify password was changed and reset token was cleared
		var refreshedUser user.User
		err = refreshedUser.FindByEmail(testUser.Email)
		require.NoError(t, err)

		expectedHash := md5.Sum([]byte("new_secure_password_123"))
		expectedHashString := hex.EncodeToString(expectedHash[:])
		assert.Equal(t, expectedHashString, refreshedUser.PasswordHash, "Password should have changed")
		assert.Nil(t, refreshedUser.ResetPwHash, "Reset hash should be cleared")
		assert.Nil(t, refreshedUser.ResetPwHashExpires, "Reset expiry should be cleared")

		// Step 6: Verify old token is no longer valid
		req, _ = http.NewRequest("GET", "/password-reset/validate/"+resetToken, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPasswordResetRequestRoute(t *testing.T) {
	router := setupTestRouter()

	t.Run("POST /password-reset/request - Success", func(t *testing.T) {
		testUser := setupTestUserForRouter(t)

		requestData := map[string]interface{}{
			"email":        testUser.Email,
			"redirect_url": "http://localhost:3000",
		}
		jsonData, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})

	t.Run("POST /password-reset/request - Non-existent email", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

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
	})

	t.Run("POST /password-reset/request - Invalid JSON", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /password-reset/request - Missing email", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		requestData := map[string]interface{}{
			"redirect_url": "http://localhost:3000",
		}
		jsonData, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /password-reset/request - Invalid email format", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		requestData := map[string]interface{}{
			"email":        "not-an-email",
			"redirect_url": "http://localhost:3000",
		}
		jsonData, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPasswordResetValidateRoute(t *testing.T) {
	router := setupTestRouter()

	t.Run("GET /password-reset/validate/:token - Success", func(t *testing.T) {
		testUser := setupTestUserForRouter(t)

		// Set a reset hash for the user
		resetHash := "valid_test_token_123456789"
		err := testUser.SetPasswordResetHash(resetHash)
		require.NoError(t, err)

		req, _ := http.NewRequest("GET", "/password-reset/validate/"+resetHash, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.True(t, response["valid"].(bool))
		assert.Equal(t, testUser.Username, response["username"])
		assert.NotNil(t, response["expires"])
	})

	t.Run("GET /password-reset/validate/:token - Invalid token", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		req, _ := http.NewRequest("GET", "/password-reset/validate/invalid_token_123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GET /password-reset/validate/:token - Empty token", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		req, _ := http.NewRequest("GET", "/password-reset/validate/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 404 because the route doesn't match
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPasswordResetResetRoute(t *testing.T) {
	router := setupTestRouter()

	t.Run("POST /password-reset/reset - Success", func(t *testing.T) {
		testUser := setupTestUserForRouter(t)
		originalPassword := testUser.PasswordHash

		// Set a reset hash for the user
		resetHash := "valid_reset_token_for_password_change"
		err := testUser.SetPasswordResetHash(resetHash)
		require.NoError(t, err)

		resetData := map[string]interface{}{
			"token":        resetHash,
			"new_password": "new_secure_password_456",
		}
		jsonData, _ := json.Marshal(resetData)

		req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Passwort erfolgreich zurückgesetzt", response["message"])

		// Verify password was changed
		var dbUser user.User
		err = dbUser.FindByEmail(testUser.Email)
		require.NoError(t, err)

		assert.NotEqual(t, originalPassword, dbUser.PasswordHash, "Password should have changed")
		assert.Nil(t, dbUser.ResetPwHash, "Reset hash should be cleared")
	})

	t.Run("POST /password-reset/reset - Invalid token", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		resetData := map[string]interface{}{
			"token":        "invalid_token_123",
			"new_password": "new_secure_password_456",
		}
		jsonData, _ := json.Marshal(resetData)

		req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /password-reset/reset - Short password", func(t *testing.T) {
		testUser := setupTestUserForRouter(t)

		resetHash := "valid_token_short_password"
		err := testUser.SetPasswordResetHash(resetHash)
		require.NoError(t, err)

		resetData := map[string]interface{}{
			"token":        resetHash,
			"new_password": "123", // Too short
		}
		jsonData, _ := json.Marshal(resetData)

		req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /password-reset/reset - Missing token", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		resetData := map[string]interface{}{
			"new_password": "new_secure_password_456",
		}
		jsonData, _ := json.Marshal(resetData)

		req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /password-reset/reset - Missing password", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		resetData := map[string]interface{}{
			"token": "some_token_123",
		}
		jsonData, _ := json.Marshal(resetData)

		req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("POST /password-reset/reset - Invalid JSON", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		req, _ := http.NewRequest("POST", "/password-reset/reset", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestPasswordResetRoutes_HTTPMethods(t *testing.T) {
	router := setupTestRouter()
	setupTestUserForRouter(t) // Setup database

	t.Run("Wrong HTTP methods should return 404 or 405", func(t *testing.T) {
		// Test wrong methods for password-reset/request (should be POST)
		req, _ := http.NewRequest("GET", "/password-reset/request", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		req, _ = http.NewRequest("PUT", "/password-reset/request", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Test wrong methods for password-reset/validate/:token (should be GET)
		req, _ = http.NewRequest("POST", "/password-reset/validate/some_token", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Test wrong methods for password-reset/reset (should be POST)
		req, _ = http.NewRequest("GET", "/password-reset/reset", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestPasswordResetRoutes_Security(t *testing.T) {
	router := setupTestRouter()

	t.Run("Routes should not require authentication", func(t *testing.T) {
		testUser := setupTestUserForRouter(t)

		// Test that password reset routes don't require Authorization header
		requestData := map[string]interface{}{
			"email":        testUser.Email,
			"redirect_url": "http://localhost:3000",
		}
		jsonData, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		// Deliberately NOT setting Authorization header

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should work without auth
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Email enumeration protection", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		// Test with non-existent email
		requestData := map[string]interface{}{
			"email":        "definitely.does.not.exist@example.com",
			"redirect_url": "http://localhost:3000",
		}
		jsonData, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/password-reset/request", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return same success message to prevent email enumeration
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "Falls ein Account mit dieser E-Mail-Adresse existiert")
	})
}

func TestPasswordResetRoutes_CORS(t *testing.T) {
	router := setupTestRouter()
	setupTestUserForRouter(t) // Setup database

	t.Run("Routes should handle CORS preflight requests", func(t *testing.T) {
		// Test OPTIONS request for CORS preflight
		req, _ := http.NewRequest("OPTIONS", "/password-reset/request", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should handle OPTIONS request appropriately
		// Note: Actual CORS headers would be set by middleware, not tested here
		// This just ensures the routes don't break with OPTIONS
		assert.NotEqual(t, http.StatusInternalServerError, w.Code)
	})
}
