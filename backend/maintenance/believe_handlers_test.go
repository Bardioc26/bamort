package maintenance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bamort/database"
	"bamort/models"
	"bamort/router"
	"bamort/testutils"
	"bamort/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMaintenanceTest(t *testing.T) (string, *gin.Engine, *models.GameSystem) {
	t.Helper()

	testutils.SetupTestEnvironment(t)
	database.ResetTestDB()
	t.Cleanup(database.ResetTestDB)
	database.SetupTestDB(true)

	var maintainer user.User
	require.NoError(t, database.DB.First(&maintainer, 1).Error)
	maintainer.Role = user.RoleMaintainer
	require.NoError(t, maintainer.Save())

	token := user.GenerateToken(&maintainer)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	router.SetupGin(r)
	protected := router.BaseRouterGrp(r)
	RegisterRoutes(protected)

	gs := models.GetGameSystem(0, "")
	require.NotNil(t, gs)

	return token, r, gs
}

func createSource(t *testing.T, gs *models.GameSystem, code string) models.Source {
	t.Helper()

	source := models.Source{
		Code:         code,
		Name:         fmt.Sprintf("Source %s", code),
		FullName:     fmt.Sprintf("Source %s", code),
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
		IsActive:     true,
	}
	require.NoError(t, database.DB.Create(&source).Error)
	return source
}

func createBelieve(t *testing.T, gs *models.GameSystem, source models.Source, name string) models.Believe {
	t.Helper()

	believe := models.Believe{
		Name:         name,
		Beschreibung: "Initial description",
		SourceID:     source.ID,
		PageNumber:   7,
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
	}
	require.NoError(t, database.DB.Create(&believe).Error)
	return believe
}

func TestListBelievesReturnsData(t *testing.T) {
	token, router, gs := setupMaintenanceTest(t)
	source := createSource(t, gs, "TSTBEL")
	created := createBelieve(t, gs, source, "Test Believe One")

	req, err := http.NewRequest(http.MethodGet, "/api/maintenance/gsm-believes", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var payload struct {
		Believes []struct {
			ID         uint   `json:"id"`
			Name       string `json:"name"`
			SourceID   uint   `json:"source_id"`
			SourceCode string `json:"source_code"`
		} `json:"believes"`
		Sources []models.Source `json:"sources"`
	}

	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	assert.NotEmpty(t, payload.Sources)

	var found bool
	for _, b := range payload.Believes {
		if b.ID == created.ID {
			found = true
			assert.Equal(t, created.Name, b.Name)
			assert.Equal(t, source.ID, b.SourceID)
			assert.Equal(t, source.Code, b.SourceCode)
		}
	}

	assert.True(t, found, "expected created believe in response")
}

func TestUpdateBelieve(t *testing.T) {
	token, router, gs := setupMaintenanceTest(t)
	sourceOld := createSource(t, gs, "OLD01")
	sourceNew := createSource(t, gs, "NEW01")
	created := createBelieve(t, gs, sourceOld, "Old Believe")

	updateBody := map[string]interface{}{
		"name":         "Updated Believe",
		"beschreibung": "Updated description",
		"source_id":    sourceNew.ID,
		"page_number":  123,
	}

	bodyBytes, err := json.Marshal(updateBody)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/maintenance/gsm-believes/%d", created.ID), bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var updated models.Believe
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &updated))
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Updated Believe", updated.Name)
	assert.Equal(t, "Updated description", updated.Beschreibung)
	assert.Equal(t, sourceNew.ID, updated.SourceID)
	assert.Equal(t, 123, updated.PageNumber)
	assert.Equal(t, gs.Name, updated.GameSystem)
	assert.Equal(t, gs.ID, updated.GameSystemId)

	var stored models.Believe
	require.NoError(t, database.DB.First(&stored, created.ID).Error)
	assert.Equal(t, "Updated Believe", stored.Name)
	assert.Equal(t, sourceNew.ID, stored.SourceID)
	assert.Equal(t, 123, stored.PageNumber)
}
