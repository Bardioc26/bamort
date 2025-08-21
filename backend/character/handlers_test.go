package character

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bamort/database"
	"bamort/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestImproveSkillHandler(t *testing.T) {
	// Setup test environment
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original != "" {
			os.Setenv("ENVIRONMENT", original)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	err := models.MigrateStructure()
	assert.NoError(t, err)

	// Create test character with ID 20
	/*
		testChar := Char{
			BamortBase: models.BamortBase{
				ID:   20,
				Name: "Test Krieger",
			},
			Typ:   "Krieger",
			Rasse: "Mensch",
			Grad:  1,
			Erfahrungsschatz: Erfahrungsschatz{
				BamortCharTrait: models.BamortCharTrait{
					CharacterID: 20,
				},
				Value: 326, // Starting EP (should end with 316 after spending 10)
			},
			Vermoegen: Vermoegen{
				BamortCharTrait: models.BamortCharTrait{
					CharacterID: 20,
				},
				Goldstücke: 390, // Starting Gold (should end with 370 after spending 20)
			},
		}


		// Add Athletik skill at level 9
		athletikSkill := models.Fertigkeit{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Athletik",
				},
				CharacterID: 20,
			},
			Fertigkeitswert: 9,
		}
		testChar.Fertigkeiten = append(testChar.Fertigkeiten, athletikSkill)

		err = testChar.Create()
		assert.NoError(t, err)
	*/
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("ImproveSkill Athletik from level 9 to 10", func(t *testing.T) {
		// Create request body matching gsmaster.LernCostRequest structure
		requestData := map[string]interface{}{
			"char_id":         20,
			"name":            "Athletik",
			"current_level":   9,
			"target_level":    10,
			"type":            "skill",
			"action":          "improve",
			"reward":          "default",
			"use_pp":          1,
			"use_gold":        0,
			"levels_to_learn": []int{10},
			"notes":           "Fertigkeit Athletik von 9 auf 10 verbessert (1 Level)",
		}
		requestBody, _ := json.Marshal(requestData)

		// Create HTTP request
		req, _ := http.NewRequest("POST", "/api/characters/improve-skill", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the actual handler function
		ImproveSkill(c)

		// Print the actual response to see what we get
		t.Logf("Response Status: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		// Check if we got a successful response
		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200 OK")

		// Parse and validate response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Validate expected response values
		expectedResponse := map[string]interface{}{
			"ep_cost":        float64(10), // JSON numbers are float64
			"from_level":     float64(9),
			"gold_cost":      float64(20),
			"message":        "Fertigkeit erfolgreich verbessert",
			"remaining_ep":   float64(316),
			"remaining_gold": float64(370),
			"skill_name":     "Athletik",
			"to_level":       float64(10),
		}

		// Check each expected field
		for key, expectedValue := range expectedResponse {
			actualValue, exists := response[key]
			assert.True(t, exists, "Response should contain field: %s", key)
			assert.Equal(t, expectedValue, actualValue, "Field %s should match expected value", key)
		}

		// Additional validations
		assert.Contains(t, response, "message", "Response should contain success message")
		assert.Equal(t, "Fertigkeit erfolgreich verbessert", response["message"])

		// Verify character state was updated correctly
		var updatedChar models.Char
		err = updatedChar.FirstID("20")
		assert.NoError(t, err)

		// Check that EP was deducted correctly
		assert.Equal(t, 316, updatedChar.Erfahrungsschatz.EP, "Character should have 316 EP remaining")

		// Check that Gold was deducted correctly
		assert.Equal(t, 370, updatedChar.Vermoegen.Goldstücke, "Character should have 370 Gold remaining")

		t.Logf("Test completed successfully!")
		t.Logf("EP: %d -> %d (cost: %.0f)", 326, updatedChar.Erfahrungsschatz.EP, response["ep_cost"])
		t.Logf("Gold: %d -> %d (cost: %.0f)", 390, updatedChar.Vermoegen.Goldstücke, response["gold_cost"])
	})

	t.Run("ImproveSkill with insufficient EP", func(t *testing.T) {
		// Create character with insufficient EP
		poorChar := models.Char{
			BamortBase: models.BamortBase{
				ID:   21,
				Name: "Poor Test Character",
			},
			Typ:   "Krieger",
			Rasse: "Mensch",
			Grad:  1,
			Erfahrungsschatz: models.Erfahrungsschatz{
				BamortCharTrait: models.BamortCharTrait{
					CharacterID: 21,
				},
				ES: 5, // Insufficient EP
			},
			Vermoegen: models.Vermoegen{
				BamortCharTrait: models.BamortCharTrait{
					CharacterID: 21,
				},
				Goldstücke: 100,
			},
		}

		// Add skill
		skill := models.SkFertigkeit{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Athletik",
				},
				CharacterID: 21,
			},
			Fertigkeitswert: 9,
		}
		poorChar.Fertigkeiten = append(poorChar.Fertigkeiten, skill)

		err = poorChar.Create()
		assert.NoError(t, err)

		requestData := map[string]interface{}{
			"char_id":       21,
			"name":          "Athletik",
			"current_level": 9,
			"target_level":  10,
			"type":          "skill",
			"action":        "improve",
			"reward":        "default",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/api/characters/improve-skill", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		ImproveSkillOld(c)

		assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 for insufficient EP")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Nicht genügend Erfahrungspunkte")
	})

	t.Run("ImproveSkill with nonexistent character", func(t *testing.T) {
		requestData := map[string]interface{}{
			"char_id":       999, // Non-existent character
			"name":          "Athletik",
			"current_level": 9,
			"target_level":  10,
			"type":          "skill",
			"action":        "improve",
			"reward":        "default",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("POST", "/api/characters/improve-skill", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		ImproveSkillOld(c)

		assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent character")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Contains(t, response["error"], "Charakter nicht gefunden")
	})
}

func TestGetAvailableSkillsNewSystem(t *testing.T) {
	// Setup test environment
	original := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "test")
	t.Cleanup(func() {
		if original != "" {
			os.Setenv("ENVIRONMENT", original)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	})

	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	err := models.MigrateStructure()
	assert.NoError(t, err)

	t.Run("GetAvailableSkillsForCharacterCreation", func(t *testing.T) {
		// Test data - the exact request format for character creation
		requestData := map[string]interface{}{
			"CharId":        0,
			"name":          "",
			"current_level": 0,
			"target_level":  1,
			"type":          "skill",
			"action":        "learn",
			"use_pp":        0,
			"use_gold":      0,
			"reward":        "default",
		}

		requestBody, err := json.Marshal(requestData)
		assert.NoError(t, err)

		// Create test request
		req, _ := http.NewRequest("POST", "/api/characters/available-skills-new", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the handler
		GetAvailableSkillsNewSystem(c)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK for character creation request")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Check response structure
		assert.Contains(t, response, "skills_by_category", "Response should contain skills_by_category")

		skillsByCategory, ok := response["skills_by_category"].(map[string]interface{})
		assert.True(t, ok, "skills_by_category should be a map")

		// Verify we have reasonable costs for character creation
		assert.Greater(t, len(skillsByCategory), 0, "Should have at least some skill categories")

		// Check that costs are reasonable for character creation (not the high fallback values)
		for _, skills := range skillsByCategory {
			if skillsList, ok := skills.([]interface{}); ok {
				for _, skill := range skillsList {
					if skillMap, ok := skill.(map[string]interface{}); ok {
						epCost := skillMap["epCost"].(float64)
						goldCost := skillMap["goldCost"].(float64)

						// Verify costs are reasonable for character creation (not fallback values)
						assert.Less(t, epCost, 1000.0, "EP cost should be reasonable for character creation")
						assert.Less(t, goldCost, 1000.0, "Gold cost should be reasonable for character creation")
						assert.Greater(t, epCost, 0.0, "EP cost should be positive")
						assert.Greater(t, goldCost, 0.0, "Gold cost should be positive")
					}
				}
			}
		}
		assert.Greater(t, len(skillsByCategory), 0, "Should return at least some skill categories")

		// Check that each category contains skills with proper structure
		for categoryName, categorySkills := range skillsByCategory {
			assert.NotEmpty(t, categoryName, "Category name should not be empty")

			skillsArray, ok := categorySkills.([]interface{})
			assert.True(t, ok, "Category skills should be an array")

			if len(skillsArray) > 0 {
				// Check first skill structure
				firstSkill, ok := skillsArray[0].(map[string]interface{})
				assert.True(t, ok, "Skill should be a map")

				// Verify skill has required fields
				assert.Contains(t, firstSkill, "name", "Skill should have name field")
				assert.Contains(t, firstSkill, "epCost", "Skill should have epCost field")
				assert.Contains(t, firstSkill, "goldCost", "Skill should have goldCost field")

				// Verify field types
				assert.IsType(t, "", firstSkill["name"], "Name should be string")
				assert.IsType(t, float64(0), firstSkill["epCost"], "epCost should be numeric")
				assert.IsType(t, float64(0), firstSkill["goldCost"], "goldCost should be numeric")
			}
		}
	})

	t.Run("GetAvailableSkillsInvalidRequest", func(t *testing.T) {
		// Test with missing required fields (type, action, reward)
		requestData := map[string]interface{}{
			"CharId":        0, // CharId 0 is valid for character creation
			"name":          "",
			"current_level": 0,
			"target_level":  1,
			// Missing "type", "action", and "reward" - should fail
			"use_pp":   0,
			"use_gold": 0,
		}

		requestBody, err := json.Marshal(requestData)
		assert.NoError(t, err)

		req, _ := http.NewRequest("POST", "/api/characters/available-skills-new", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		GetAvailableSkillsNewSystem(c)

		// Should return 400 Bad Request for missing required fields
		assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 for missing required fields")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.True(t, len(response["error"].(string)) > 0, "Error message should not be empty")
	})

	t.Run("GetAvailableSkillsInvalidRewardType", func(t *testing.T) {
		// Test with invalid reward type
		requestData := map[string]interface{}{
			"CharId":        0,
			"name":          "",
			"current_level": 0,
			"target_level":  1,
			"type":          "skill",
			"action":        "learn",
			"use_pp":        0,
			"use_gold":      0,
			"reward":        "invalid",
		}

		requestBody, err := json.Marshal(requestData)
		assert.NoError(t, err)

		req, _ := http.NewRequest("POST", "/api/characters/available-skills-new", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		GetAvailableSkillsNewSystem(c)

		// Should return 400 Bad Request for invalid reward type
		assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 for invalid reward type")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})
}

// TestGetAvailableSkillsForCreation tests the character creation skills endpoint
func TestGetAvailableSkillsForCreation(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	tests := []struct {
		name           string
		characterClass string
		expectStatus   int
		expectError    bool
	}{
		{
			name:           "ValidCharacterClass",
			characterClass: "As",
			expectStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "MagierCharacterClass",
			characterClass: "Ma",
			expectStatus:   http.StatusOK,
			expectError:    false,
		},
		{
			name:           "EmptyCharacterClass",
			characterClass: "",
			expectStatus:   http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			requestData := gin.H{
				"characterClass": tt.characterClass,
			}
			requestBody, _ := json.Marshal(requestData)

			// Create HTTP request
			req, _ := http.NewRequest("POST", "/api/characters/available-skills-creation", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer test-token")

			// Create response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the handler directly (it will handle JSON parsing internally)
			GetAvailableSkillsForCreation(c)

			// Verify response
			assert.Equal(t, tt.expectStatus, w.Code)

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check response structure
				assert.Contains(t, response, "skills_by_category")

				skillsByCategory, ok := response["skills_by_category"].(map[string]interface{})
				assert.True(t, ok)
				assert.Greater(t, len(skillsByCategory), 0, "Should have at least some skill categories")

				// Verify skills have learnCost field
				hasNonDefaultCost := false
				for categoryName, skills := range skillsByCategory {
					if skillsList, ok := skills.([]interface{}); ok {
						for _, skill := range skillsList {
							if skillMap, ok := skill.(map[string]interface{}); ok {
								assert.Contains(t, skillMap, "name", "Skill should have name")
								assert.Contains(t, skillMap, "learnCost", "Skill should have learnCost")

								learnCost := skillMap["learnCost"].(float64)
								assert.Greater(t, learnCost, 0.0, "Learn cost should be positive")
								assert.Less(t, learnCost, 500.0, "Learn cost should be reasonable for character creation")

								// Check if we have skills with non-default costs
								if learnCost != 50 {
									hasNonDefaultCost = true
								}

								// Log individual skill costs for debugging
								if tt.characterClass == "Ma" {
									t.Logf("Skill: %s, Category: %s, LearnCost: %.0f", skillMap["name"], categoryName, learnCost)
								}
							}
						}
						if tt.characterClass == "Ma" {
							break // Only log first category for Ma
						}
					}
					_ = categoryName // Mark as used
				}

				// For Ma class, we should get some skills with different costs than the default 50
				if tt.characterClass == "Ma" {
					// This is more of an informational test - we want to see what costs we get
					t.Logf("Ma class - has skills with non-default costs: %v", hasNonDefaultCost)
				} // Log some sample data for verification
				t.Logf("Character creation skills loaded for class %s: %d categories", tt.characterClass, len(skillsByCategory))
			} else {
				// For error cases, verify error response
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			}
		})
	}
}
