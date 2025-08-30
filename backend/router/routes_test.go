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

// =============================================================================
// Register Route Tests
// =============================================================================

func TestRegisterRoute(t *testing.T) {
	router := setupTestRouter()

	t.Run("POST /register - Success", func(t *testing.T) {
		database.SetupTestDB()
		err := user.MigrateStructure()
		require.NoError(t, err, "Failed to migrate user structure")

		randomSuffix := rand.Intn(100000)
		userData := map[string]interface{}{
			"username": fmt.Sprintf("testuser_%d", randomSuffix),
			"password": "testpassword123",
			"email":    fmt.Sprintf("test_%d@example.com", randomSuffix),
		}
		jsonData, _ := json.Marshal(userData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "User registered successfully")

		// Verify user was created in database
		var createdUser user.User
		err = createdUser.First(userData["username"].(string))
		assert.NoError(t, err, "User should be created in database")
		assert.Equal(t, userData["username"], createdUser.Username)
		assert.Equal(t, userData["email"], createdUser.Email)
		// Verify password was hashed
		assert.NotEqual(t, userData["password"], createdUser.PasswordHash)
	})

	t.Run("POST /register - Missing Required Fields", func(t *testing.T) {
		// Note: The current implementation allows creation of users with missing fields
		// The database structure allows empty strings for non-primary key fields
		database.SetupTestDB()
		err := user.MigrateStructure()
		require.NoError(t, err, "Failed to migrate user structure")

		testCases := []struct {
			name           string
			data           map[string]interface{}
			expectedStatus int
		}{
			{
				name: "Missing username - should return error",
				data: map[string]interface{}{
					"password": "testpassword123",
					"email":    fmt.Sprintf("missing_username_%d@example.com", rand.Intn(100000)),
				},
				expectedStatus: http.StatusBadRequest, // Now returns error for empty username
			},
			{
				name: "Missing password - should return error",
				data: map[string]interface{}{
					"username": fmt.Sprintf("testuser_nopass_%d", rand.Intn(100000)),
					"email":    fmt.Sprintf("testnopass_%d@example.com", rand.Intn(100000)),
				},
				expectedStatus: http.StatusBadRequest, // Now returns error for empty password
			},
			{
				name: "Missing email - should return error",
				data: map[string]interface{}{
					"username": fmt.Sprintf("testuser_noemail_%d", rand.Intn(100000)),
					"password": "testpassword123",
				},
				expectedStatus: http.StatusBadRequest, // Now returns error for empty email
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				jsonData, _ := json.Marshal(tc.data)

				req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatus, w.Code)

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				if tc.expectedStatus == http.StatusCreated {
					assert.Contains(t, response["message"], "User registered successfully")
				} else {
					assert.Contains(t, response, "error")
				}
			})
		}
	})

	t.Run("POST /register - Duplicate Username", func(t *testing.T) {
		database.SetupTestDB()
		err := user.MigrateStructure()
		require.NoError(t, err, "Failed to migrate user structure")

		// Create first user
		randomSuffix := rand.Intn(100000)
		userData1 := map[string]interface{}{
			"username": fmt.Sprintf("duplicate_user_%d", randomSuffix),
			"password": "testpassword123",
			"email":    fmt.Sprintf("first_%d@example.com", randomSuffix),
		}
		jsonData1, _ := json.Marshal(userData1)

		req1, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		// Try to create second user with same username
		userData2 := map[string]interface{}{
			"username": userData1["username"], // Same username
			"password": "differentpassword",
			"email":    fmt.Sprintf("second_%d@example.com", randomSuffix),
		}
		jsonData2, _ := json.Marshal(userData2)

		req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusInternalServerError, w2.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Failed to create user")
	})

	t.Run("POST /register - Duplicate Email", func(t *testing.T) {
		database.SetupTestDB()
		err := user.MigrateStructure()
		require.NoError(t, err, "Failed to migrate user structure")

		// Create first user
		randomSuffix := rand.Intn(100000)
		userData1 := map[string]interface{}{
			"username": fmt.Sprintf("first_user_%d", randomSuffix),
			"password": "testpassword123",
			"email":    fmt.Sprintf("duplicate_%d@example.com", randomSuffix),
		}
		jsonData1, _ := json.Marshal(userData1)

		req1, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		assert.Equal(t, http.StatusCreated, w1.Code)

		// Try to create second user with same email
		userData2 := map[string]interface{}{
			"username": fmt.Sprintf("second_user_%d", randomSuffix),
			"password": "differentpassword",
			"email":    userData1["email"], // Same email
		}
		jsonData2, _ := json.Marshal(userData2)

		req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")

		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusInternalServerError, w2.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Failed to create user")
	})

	t.Run("POST /register - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
	})
}

// =============================================================================
// Login Route Tests
// =============================================================================

func TestLoginRoute(t *testing.T) {
	router := setupTestRouter()

	t.Run("POST /login - Success", func(t *testing.T) {
		// Setup user first
		testUser := setupTestUserForRouter(t)

		loginData := map[string]interface{}{
			"username": testUser.Username,
			"password": "testpassword123", // Original password before hashing
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Login successful", response["message"])
		assert.Contains(t, response, "token")
		assert.NotEmpty(t, response["token"])

		// Verify token format (should contain user ID)
		token := response["token"].(string)
		assert.Contains(t, token, ".")
		assert.Contains(t, token, ":")
	})

	t.Run("POST /login - Invalid Username", func(t *testing.T) {
		setupTestUserForRouter(t) // Setup database

		loginData := map[string]interface{}{
			"username": "nonexistentuser",
			"password": "testpassword123",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Invalid username")
	})

	t.Run("POST /login - Invalid Password", func(t *testing.T) {
		testUser := setupTestUserForRouter(t)

		loginData := map[string]interface{}{
			"username": testUser.Username,
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Invalid username")
	})

	t.Run("POST /login - Missing Required Fields", func(t *testing.T) {
		// Note: The database may contain users with empty usernames from test data
		// We need to test with usernames that definitely don't exist to get proper errors
		setupTestUserForRouter(t) // Setup database

		testCases := []struct {
			name           string
			data           map[string]interface{}
			expectedStatus int
		}{
			{
				name: "Missing username - should be unauthorized",
				data: map[string]interface{}{
					"password": "testpassword123",
				},
				expectedStatus: http.StatusUnauthorized, // No username provided
			},
			{
				name: "Missing password - should be unauthorized",
				data: map[string]interface{}{
					"username": fmt.Sprintf("nonexistent_%d", rand.Intn(100000)),
				},
				expectedStatus: http.StatusUnauthorized, // No password provided
			},
			{
				name: "Nonexistent username - should be unauthorized",
				data: map[string]interface{}{
					"username": fmt.Sprintf("definitely_nonexistent_%d", rand.Intn(100000)),
					"password": "testpassword123",
				},
				expectedStatus: http.StatusUnauthorized, // Username doesn't exist
			},
			{
				name: "Empty password with nonexistent user - should be unauthorized",
				data: map[string]interface{}{
					"username": fmt.Sprintf("another_nonexistent_%d", rand.Intn(100000)),
					"password": "",
				},
				expectedStatus: http.StatusUnauthorized, // Password mismatch
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				jsonData, _ := json.Marshal(tc.data)

				req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedStatus, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Contains(t, response, "error")
			})
		}
	})

	t.Run("POST /login - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	t.Run("POST /login - Empty Request Body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
	})
}

// =============================================================================
// Integration Tests for Register and Login Flow
// =============================================================================

func TestRegisterLoginFlow(t *testing.T) {
	router := setupTestRouter()

	t.Run("Complete Register and Login Flow", func(t *testing.T) {
		database.SetupTestDB()
		err := user.MigrateStructure()
		require.NoError(t, err, "Failed to migrate user structure")

		randomSuffix := rand.Intn(100000)
		username := fmt.Sprintf("flowtest_user_%d", randomSuffix)
		password := "testpassword123"
		email := fmt.Sprintf("flowtest_%d@example.com", randomSuffix)

		// Step 1: Register user
		registerData := map[string]interface{}{
			"username": username,
			"password": password,
			"email":    email,
		}
		jsonData, _ := json.Marshal(registerData)

		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Step 2: Login with registered credentials
		loginData := map[string]interface{}{
			"username": username,
			"password": password,
		}
		jsonData, _ = json.Marshal(loginData)

		req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Login successful", response["message"])
		assert.Contains(t, response, "token")
		assert.NotEmpty(t, response["token"])

		// Step 3: Verify login fails with wrong password
		wrongLoginData := map[string]interface{}{
			"username": username,
			"password": "wrongpassword",
		}
		jsonData, _ = json.Marshal(wrongLoginData)

		req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
