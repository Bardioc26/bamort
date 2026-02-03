package character

import (
	"bamort/database"
	"bamort/models"
	"bamort/user"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupShareHandlerTestEnvironment(t *testing.T) {
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original != "" {
			os.Setenv("ENVIRONMENT", original)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})

	database.SetupTestDB(true, true)
	t.Cleanup(database.ResetTestDB)

	err := user.MigrateStructure()
	require.NoError(t, err, "Should migrate user structure")

	err = models.MigrateStructure()
	require.NoError(t, err, "Should migrate models structure")

	gin.SetMode(gin.TestMode)
}

func createHashedUser(t *testing.T, username, password, email, displayName string) *user.User {
	u := &user.User{
		Username:     username,
		PasswordHash: password,
		Email:        email,
		DisplayName:  displayName,
	}

	hashed := md5.Sum([]byte(password))
	u.PasswordHash = hex.EncodeToString(hashed[:])

	err := u.Create()
	require.NoError(t, err, "Should create user")
	return u
}

func TestGetAvailableUsersForSharingReturnsDisplayNames(t *testing.T) {
	setupShareHandlerTestEnvironment(t)

	owner := createHashedUser(t, "owneruser", "ownerpass", "owner@example.com", "")
	sharedUser := createHashedUser(t, "shareduser", "sharedpass", "shared@example.com", "Shared Display")

	char := models.Char{
		BamortBase: models.BamortBase{Name: "Shared Character"},
		UserID:     owner.UserID,
	}
	err := database.DB.Create(&char).Error
	require.NoError(t, err, "Should create character")

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/characters/%d/available-users", char.ID), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: fmt.Sprintf("%d", char.ID)}}
	c.Set("userID", owner.UserID)

	GetAvailableUsersForSharing(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Len(t, response, 1, "Only the non-owner user should be returned")

	entry := response[0]
	assert.Equal(t, float64(sharedUser.UserID), entry["user_id"])
	assert.Equal(t, sharedUser.DisplayName, entry["display_name"])
}
