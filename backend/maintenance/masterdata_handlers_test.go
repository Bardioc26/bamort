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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createGameSystem(t *testing.T, code string) models.GameSystem {
	t.Helper()
	gs := models.GameSystem{Code: code, Name: "Game System " + code, Description: "desc", IsActive: true}
	require.NoError(t, database.DB.Create(&gs).Error)
	return gs
}

func TestListGameSystems(t *testing.T) {
	token, router, _ := setupMaintenanceTest(t)
	created := createGameSystem(t, "TSTGS")

	req, err := http.NewRequest(http.MethodGet, "/api/maintenance/game-systems", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var payload struct {
		GameSystems []models.GameSystem `json:"game_systems"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))

	var found bool
	for _, gs := range payload.GameSystems {
		if gs.ID == created.ID {
			found = true
			assert.Equal(t, created.Code, gs.Code)
			assert.Equal(t, created.Name, gs.Name)
		}
	}
	assert.True(t, found, "expected created game system in response")
}

func TestUpdateGameSystem(t *testing.T) {
	token, router, _ := setupMaintenanceTest(t)
	gs := createGameSystem(t, "UPDGS")

	body := map[string]interface{}{
		"name":        "Updated GS",
		"description": "Updated desc",
		"is_active":   false,
	}
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/maintenance/game-systems/%d", gs.ID), bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var updated models.GameSystem
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &updated))
	assert.Equal(t, gs.ID, updated.ID)
	assert.Equal(t, "Updated GS", updated.Name)
	assert.Equal(t, "Updated desc", updated.Description)
	assert.False(t, updated.IsActive)
}

func createLitSource(t *testing.T, gs models.GameSystem, code string) models.Source {
	t.Helper()
	src := models.Source{
		Code:         code,
		Name:         "Source " + code,
		FullName:     "Full " + code,
		Edition:      "1",
		Publisher:    "Pub",
		PublishYear:  2025,
		Description:  "Desc",
		IsCore:       false,
		IsActive:     true,
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
	}
	require.NoError(t, database.DB.Create(&src).Error)
	return src
}

func TestListLitSources(t *testing.T) {
	token, router, gs := setupMaintenanceTest(t)
	src := createLitSource(t, *gs, "SRC01")

	req, _ := http.NewRequest(http.MethodGet, "/api/maintenance/gsm-lit-sources", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var payload struct {
		Sources []models.Source `json:"sources"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))

	var found bool
	for _, s := range payload.Sources {
		if s.ID == src.ID {
			found = true
			assert.Equal(t, src.Code, s.Code)
			assert.Equal(t, src.Name, s.Name)
		}
	}
	assert.True(t, found)
}

func TestUpdateLitSource(t *testing.T) {
	token, router, gs := setupMaintenanceTest(t)
	src := createLitSource(t, *gs, "SRCUPD")

	body := map[string]interface{}{
		"name":         "Updated Source",
		"full_name":    "Updated Full",
		"edition":      "2",
		"publisher":    "NewPub",
		"publish_year": 2026,
		"description":  "New Desc",
		"is_active":    false,
		"is_core":      true,
	}
	payload, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/maintenance/gsm-lit-sources/%d", src.ID), bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var updated models.Source
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &updated))
	assert.Equal(t, "Updated Source", updated.Name)
	assert.Equal(t, "Updated Full", updated.FullName)
	assert.Equal(t, "2", updated.Edition)
	assert.Equal(t, "NewPub", updated.Publisher)
	assert.Equal(t, 2026, updated.PublishYear)
	assert.False(t, updated.IsActive)
	assert.True(t, updated.IsCore)
}

func createMisc(t *testing.T, gs models.GameSystem, key, value string) models.MiscLookup {
	t.Helper()
	m := models.MiscLookup{
		Key:          key,
		Value:        value,
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
	}
	require.NoError(t, database.DB.Create(&m).Error)
	return m
}

func TestListMisc(t *testing.T) {
	token, router, gs := setupMaintenanceTest(t)
	item := createMisc(t, *gs, "origin", "North")

	req, _ := http.NewRequest(http.MethodGet, "/api/maintenance/gsm-misc", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var payload struct {
		Items []models.MiscLookup `json:"misc"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))

	var found bool
	for _, it := range payload.Items {
		if it.ID == item.ID {
			found = true
			assert.Equal(t, "origin", it.Key)
			assert.Equal(t, "North", it.Value)
		}
	}
	assert.True(t, found)
}

func TestUpdateMisc(t *testing.T) {
	token, router, gs := setupMaintenanceTest(t)
	item := createMisc(t, *gs, "race", "Human")

	body := map[string]interface{}{
		"value": "Elf",
	}
	payload, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/maintenance/gsm-misc/%d", item.ID), bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var updated models.MiscLookup
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &updated))
	assert.Equal(t, "Elf", updated.Value)
}

func createSkillImprovementCost(t *testing.T, gs models.GameSystem) models.SkillImprovementCost {
	t.Helper()
	cost := models.SkillImprovementCost{
		CurrentLevel: 5,
		TERequired:   2,
		CategoryID:   1,
		DifficultyID: 1,
	}
	require.NoError(t, database.DB.Create(&cost).Error)
	return cost
}

func TestListSkillImprovementCost(t *testing.T) {
	token, router, gs := setupMaintenanceTest(t)
	created := createSkillImprovementCost(t, *gs)

	req, _ := http.NewRequest(http.MethodGet, "/api/maintenance/skill-improvement-cost2", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var payload struct {
		Costs []models.SkillImprovementCost `json:"costs"`
	}
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))

	var found bool
	for _, c := range payload.Costs {
		if c.ID == created.ID {
			found = true
			assert.Equal(t, created.CurrentLevel, c.CurrentLevel)
			assert.Equal(t, created.TERequired, c.TERequired)
		}
	}
	assert.True(t, found)
}

func TestUpdateSkillImprovementCost(t *testing.T) {
	token, router, gs := setupMaintenanceTest(t)
	created := createSkillImprovementCost(t, *gs)

	body := map[string]interface{}{
		"te_required":   5,
		"current_level": 6,
	}
	payload, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/maintenance/skill-improvement-cost2/%d", created.ID), bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var updated models.SkillImprovementCost
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &updated))
	assert.Equal(t, 5, updated.TERequired)
	assert.Equal(t, 6, updated.CurrentLevel)
}
