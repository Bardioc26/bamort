package character

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetLernCostNewSystem tests the new database-driven learning cost system
func TestGetLernCostNewSystem(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	// Migrate the schema
	err := models.MigrateStructure()
	assert.NoError(t, err)
	/*
		// Try to initialize the new learning costs system
		// This might fail in read-only test databases, so we handle that gracefully
		err = gsmaster.InitializeLearningCostsSystem()
		hasLearningCosts := (err == nil)

		if !hasLearningCosts {
			t.Logf("Note: Learning costs system not available in test environment: %v", err)
		}
	*/
	hasLearningCosts := true

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("GetLernCostNewSystem functionality test", func(t *testing.T) {
		// Create request body using gsmaster.LernCostRequest structure
		requestData := gsmaster.LernCostRequest{
			CharId:       20,         // CharacterID = 20
			Name:         "Athletik", // SkillName = Athletik
			CurrentLevel: 9,          // CurrentLevel = 9
			Type:         "skill",    // Type = skill
			Action:       "improve",  // Action = improve (since we have current level)
			TargetLevel:  0,          // TargetLevel = 0 (will calculate up to level 18)
			UsePP:        0,          // No practice points used
			UseGold:      0,
			Reward:       &[]string{"default"}[0], // Default reward type
		}
		requestBody, _ := json.Marshal(requestData)

		// Create HTTP request
		req, _ := http.NewRequest("POST", "/api/characters/lerncost-new", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the new function
		GetLernCostNewSystem(c)

		// Check response based on whether learning costs are available
		if hasLearningCosts {
			assert.Equal(t, http.StatusOK, w.Code, "Request should succeed when learning costs are available: %s", w.Body.String())

			var response []gsmaster.SkillCostResultNew
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")

			if len(response) > 0 {
				// Basic structure validation
				assert.Equal(t, "20", response[0].CharacterID, "Character ID should match")
				assert.Equal(t, "Athletik", response[0].SkillName, "Skill name should match")
				assert.Equal(t, 10, response[0].TargetLevel, "First target level should be 10")
				assert.NotEmpty(t, response[0].Category, "Category should be set")
				assert.NotEmpty(t, response[0].Difficulty, "Difficulty should be set")

				t.Logf("New system successfully calculated costs for %d levels", len(response))
				t.Logf("First level (10): EP=%d, Category=%s, Difficulty=%s",
					response[0].EP, response[0].Category, response[0].Difficulty)
			}
		} else {
			// Should return an error indicating learning costs are not available
			assert.Equal(t, http.StatusBadRequest, w.Code, "Should return error when learning costs not available")
			assert.Contains(t, w.Body.String(), "nicht gefunden", "Error message should indicate data not found")

			t.Logf("Expected error response: %s", w.Body.String())
		}
	})

	t.Run("Function structure and API compatibility", func(t *testing.T) {
		// Test that the function has the correct signature and can handle the request
		// Even if learning costs aren't available, the function should parse the request correctly

		requestData := gsmaster.LernCostRequest{
			CharId:       20,
			Name:         "Nonexistent Skill",
			CurrentLevel: 5,
			Type:         "skill",
			Action:       "improve",
			TargetLevel:  0,
			UsePP:        0,
			UseGold:      0,
			Reward:       &[]string{"default"}[0],
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/api/characters/lerncost-new", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		GetLernCostNewSystem(c)

		// Should return an error (either character not found, skill not found, or learning costs not available)
		// but shouldn't panic or have internal server errors
		assert.True(t, w.Code == http.StatusBadRequest || w.Code == http.StatusNotFound,
			"Should return client error, not server error: %d - %s", w.Code, w.Body.String())

		// Ensure response is JSON
		var jsonResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &jsonResponse)
		assert.NoError(t, err, "Response should be valid JSON even on error")

		t.Logf("Function handles error gracefully: %s", w.Body.String())
	})
}
