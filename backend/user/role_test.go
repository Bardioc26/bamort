package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bamort/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupRoleTestEnvironment sets up the test environment for role tests
func setupRoleTestEnvironment(t *testing.T) {
	setupTestEnvironment(t)
	database.SetupTestDB()
	err := MigrateStructure()
	require.NoError(t, err, "Should migrate user structure")
	gin.SetMode(gin.TestMode)
}

// TestUserRoleDefaults tests that new users get standard role by default
func TestUserRoleDefaults(t *testing.T) {
	setupRoleTestEnvironment(t)

	user := &User{
		Username:     "roletest_user",
		PasswordHash: "hashedpw",
		Email:        "roletest@example.com",
	}

	err := user.Create()
	require.NoError(t, err, "Should create user")
	assert.Equal(t, RoleStandardUser, user.Role, "New users should have standard role")
}

// TestUserDisplayNameDefaultsToUsername ensures a user's display name falls back to the username
func TestUserDisplayNameDefaultsToUsername(t *testing.T) {
	setupRoleTestEnvironment(t)

	user := &User{
		Username:     "display_user",
		PasswordHash: "hashedpw",
		Email:        "display@example.com",
	}

	err := user.Create()
	require.NoError(t, err, "Should create user")
	assert.Equal(t, user.Username, user.DisplayName, "DisplayName should default to username")
}

// TestRoleValidation tests role validation
func TestRoleValidation(t *testing.T) {
	assert.True(t, ValidateRole(RoleStandardUser), "standard should be valid")
	assert.True(t, ValidateRole(RoleMaintainer), "maintainer should be valid")
	assert.True(t, ValidateRole(RoleAdmin), "admin should be valid")
	assert.False(t, ValidateRole("invalid"), "invalid should not be valid")
	assert.False(t, ValidateRole(""), "empty should not be valid")
}

// TestRoleHierarchy tests role hierarchy methods
func TestRoleHierarchy(t *testing.T) {
	standardUser := &User{Role: RoleStandardUser}
	maintainer := &User{Role: RoleMaintainer}
	admin := &User{Role: RoleAdmin}

	// Test IsStandardUser
	assert.True(t, standardUser.IsStandardUser(), "standard user should pass IsStandardUser")
	assert.True(t, maintainer.IsStandardUser(), "maintainer should pass IsStandardUser")
	assert.True(t, admin.IsStandardUser(), "admin should pass IsStandardUser")

	// Test IsMaintainer
	assert.False(t, standardUser.IsMaintainer(), "standard user should not pass IsMaintainer")
	assert.True(t, maintainer.IsMaintainer(), "maintainer should pass IsMaintainer")
	assert.True(t, admin.IsMaintainer(), "admin should pass IsMaintainer")

	// Test IsAdmin
	assert.False(t, standardUser.IsAdmin(), "standard user should not pass IsAdmin")
	assert.False(t, maintainer.IsAdmin(), "maintainer should not pass IsAdmin")
	assert.True(t, admin.IsAdmin(), "admin should pass IsAdmin")
}

// TestListUsers tests listing all users (admin only)
func TestListUsers(t *testing.T) {
	setupRoleTestEnvironment(t)

	// Create test users
	admin := &User{Username: "listadmin", PasswordHash: "hash", Email: "listadmin@test.com", Role: RoleAdmin, DisplayName: "List Admin"}
	require.NoError(t, admin.Create())

	standardUser := &User{Username: "listuser", PasswordHash: "hash", Email: "listuser@test.com", Role: RoleStandardUser}
	require.NoError(t, standardUser.Create())

	// Test admin access
	router := gin.Default()
	router.GET("/users", func(c *gin.Context) {
		c.Set("user", admin)
	}, RequireAdmin(), ListUsers)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Admin should access list")

	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	var adminEntry map[string]interface{}
	var standardEntry map[string]interface{}
	for _, entry := range response {
		switch entry["username"] {
		case admin.Username:
			adminEntry = entry
		case standardUser.Username:
			standardEntry = entry
		}
	}

	require.NotNil(t, adminEntry, "admin should be present in response")
	require.NotNil(t, standardEntry, "standard user should be present in response")

	assert.Equal(t, admin.DisplayName, adminEntry["display_name"], "admin display name should be returned")
	assert.Equal(t, standardUser.Username, standardEntry["display_name"], "standard user should fall back to username when display name is empty")

	// Test standard user access (should fail)
	router2 := gin.Default()
	router2.GET("/users", func(c *gin.Context) {
		c.Set("user", standardUser)
	}, RequireAdmin(), ListUsers)

	req2, _ := http.NewRequest("GET", "/users", nil)
	w2 := httptest.NewRecorder()
	router2.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusForbidden, w2.Code, "Standard user should not access list")
}

// TestUpdateUserRole tests updating user role (admin only)
func TestUpdateUserRole(t *testing.T) {
	setupRoleTestEnvironment(t)

	// Create admin and target user
	admin := &User{Username: "roleadmin", PasswordHash: "hash", Email: "roleadmin@test.com", Role: RoleAdmin}
	require.NoError(t, admin.Create())

	targetUser := &User{Username: "roletarget", PasswordHash: "hash", Email: "roletarget@test.com", Role: RoleStandardUser}
	require.NoError(t, targetUser.Create())

	// Test valid role update
	router := gin.Default()
	router.PUT("/users/:id/role", func(c *gin.Context) {
		c.Set("user", admin)
	}, RequireAdmin(), UpdateUserRole)

	updateData := map[string]string{"role": RoleMaintainer}
	jsonData, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d/role", targetUser.UserID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Admin should update role")

	// Verify role was updated
	var updatedUser User
	require.NoError(t, updatedUser.FirstId(targetUser.UserID))
	assert.Equal(t, RoleMaintainer, updatedUser.Role, "Role should be updated to maintainer")

	// Test invalid role update
	invalidData := map[string]string{"role": "invalid"}
	jsonData2, _ := json.Marshal(invalidData)

	req2, _ := http.NewRequest("PUT", fmt.Sprintf("/users/%d/role", targetUser.UserID), bytes.NewBuffer(jsonData2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code, "Should reject invalid role")
}

// TestDeleteUser tests deleting a user (admin only)
func TestDeleteUser(t *testing.T) {
	setupRoleTestEnvironment(t)

	// Create admin and target user
	admin := &User{Username: "deladmin", PasswordHash: "hash", Email: "deladmin@test.com", Role: RoleAdmin}
	require.NoError(t, admin.Create())

	targetUser := &User{Username: "deltarget", PasswordHash: "hash", Email: "deltarget@test.com", Role: RoleStandardUser}
	require.NoError(t, targetUser.Create())

	// Test deletion
	router := gin.Default()
	router.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("user", admin)
	}, RequireAdmin(), DeleteUser)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", targetUser.UserID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Admin should delete user")

	// Verify user was deleted
	var deletedUser User
	err := deletedUser.FirstId(targetUser.UserID)
	assert.Error(t, err, "User should not exist after deletion")
}

// TestMaintainerPermissions tests maintainer-specific permissions
func TestMaintainerPermissions(t *testing.T) {
	setupRoleTestEnvironment(t)

	standardUser := &User{Username: "permuser", PasswordHash: "hash", Email: "permuser@test.com", Role: RoleStandardUser}
	require.NoError(t, standardUser.Create())

	maintainer := &User{Username: "permmaintainer", PasswordHash: "hash", Email: "permmaintainer@test.com", Role: RoleMaintainer}
	require.NoError(t, maintainer.Create())

	admin := &User{Username: "permadmin", PasswordHash: "hash", Email: "permadmin@test.com", Role: RoleAdmin}
	require.NoError(t, admin.Create())

	router := gin.Default()
	router.GET("/maintainer-only", func(c *gin.Context) {
		c.Set("user", standardUser)
	}, RequireMaintainer(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Standard user should fail
	req, _ := http.NewRequest("GET", "/maintainer-only", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code, "Standard user should not access")

	// Maintainer should succeed
	router2 := gin.Default()
	router2.GET("/maintainer-only", func(c *gin.Context) {
		c.Set("user", maintainer)
	}, RequireMaintainer(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req2, _ := http.NewRequest("GET", "/maintainer-only", nil)
	w2 := httptest.NewRecorder()
	router2.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code, "Maintainer should access")

	// Admin should also succeed
	router3 := gin.Default()
	router3.GET("/maintainer-only", func(c *gin.Context) {
		c.Set("user", admin)
	}, RequireMaintainer(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req3, _ := http.NewRequest("GET", "/maintainer-only", nil)
	w3 := httptest.NewRecorder()
	router3.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code, "Admin should access maintainer endpoints")
}
