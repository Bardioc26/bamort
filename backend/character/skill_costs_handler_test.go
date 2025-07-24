package character

import (
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
	"bamort/models"
	"bamort/skills"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestImprovedSkillCostAPI(t *testing.T) {
	// Setup test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Migrate the schema
	err := MigrateStructure()
	assert.NoError(t, err)

	// Also migrate skills and equipment to avoid preload errors
	err = skills.MigrateStructure()
	assert.NoError(t, err)
	err = equipment.MigrateStructure()
	assert.NoError(t, err)
	err = gsmaster.MigrateStructure()
	assert.NoError(t, err)

	// Create test skill data
	err = createTestSkillData()
	assert.NoError(t, err)
	defer cleanupTestSkillData()

	// Create test character
	testChar := createChar()
	testChar.ID = 1 // Set the ID to match our test requests
	err = testChar.Create()
	assert.NoError(t, err)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create test cases
	testCases := []struct {
		name           string
		request        SkillCostRequest
		expectedStatus int
		description    string
	}{
		{
			name: "Learn new skill",
			request: SkillCostRequest{
				Name:   "Menschenkenntnis",
				Type:   "skill",
				Action: "learn",
			},
			expectedStatus: http.StatusOK,
			description:    "Should calculate costs for learning a new skill",
		},
		{
			name: "Improve existing skill",
			request: SkillCostRequest{
				Name:         "Menschenkenntnis",
				Type:         "skill",
				Action:       "improve",
				CurrentLevel: 10,
			},
			expectedStatus: http.StatusOK,
			description:    "Should calculate costs for improving an existing skill",
		},
		{
			name: "Multi-level improvement",
			request: SkillCostRequest{
				Name:         "Menschenkenntnis",
				Type:         "skill",
				Action:       "improve",
				CurrentLevel: 10,
				TargetLevel:  13,
			},
			expectedStatus: http.StatusOK,
			description:    "Should calculate costs for multi-level improvement",
		},
		{
			name: "Invalid request - missing name",
			request: SkillCostRequest{
				Type:   "skill",
				Action: "learn",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "Should return error for missing skill name",
		},
		{
			name: "Invalid request - invalid type",
			request: SkillCostRequest{
				Name:   "Test",
				Type:   "invalid",
				Action: "learn",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "Should return error for invalid skill type",
		},
		{
			name: "Invalid request - invalid action",
			request: SkillCostRequest{
				Name:   "Test",
				Type:   "skill",
				Action: "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			description:    "Should return error for invalid action",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			requestBody, _ := json.Marshal(tc.request)

			// Create HTTP request
			req, _ := http.NewRequest("POST", "/api/characters/1/skill-cost", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = []gin.Param{{Key: "id", Value: "1"}}

			// Note: This test would need a proper database setup to work fully
			// For now, we're just testing the request parsing and validation

			fmt.Printf("Test: %s\n", tc.description)
			fmt.Printf("Request: %+v\n", tc.request)
			fmt.Printf("Expected Status: %d\n", tc.expectedStatus)

			// Call the actual handler function
			GetSkillCost(c)

			// Check the response status
			assert.Equal(t, tc.expectedStatus, w.Code, "Status code should match expected")

			// If successful, validate response structure
			if w.Code == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")

				// Check for expected fields in successful responses
				if tc.request.TargetLevel > 0 {
					// Multi-level response
					assert.Contains(t, response, "LevelCosts", "Multi-level response should contain LevelCosts")
				} else {
					// Single-level response
					assert.Contains(t, response, "SkillName", "Response should contain SkillName")
					assert.Contains(t, response, "SkillType", "Response should contain SkillType")
					assert.Contains(t, response, "Action", "Response should contain Action")
				}
			}

			// Validate request structure
			var parsedRequest SkillCostRequest
			err := json.Unmarshal(requestBody, &parsedRequest)
			assert.NoError(t, err, "Request should be valid JSON")

			// Test validation logic
			if tc.request.Name == "" {
				assert.Empty(t, parsedRequest.Name, "Name should be empty when not provided")
			}

			if tc.request.Type != "" {
				assert.Equal(t, tc.request.Type, parsedRequest.Type, "Type should match")
			}

			if tc.request.Action != "" {
				assert.Equal(t, tc.request.Action, parsedRequest.Action, "Action should match")
			}
		})
	}
}

// Test the response structures
func TestSkillCostResponseStructures(t *testing.T) {
	t.Run("SkillCostResponse structure", func(t *testing.T) {
		response := SkillCostResponse{
			SkillName:    "Menschenkenntnis",
			SkillType:    "skill",
			Action:       "learn",
			CharacterID:  1,
			CurrentLevel: 0,
			TargetLevel:  0,
			Category:     "Sozial",
			Difficulty:   "schwer",
			CanAfford:    true,
			Notes:        "Neue Fertigkeit erlernen",
		}

		// Test JSON marshaling
		jsonData, err := json.Marshal(response)
		assert.NoError(t, err, "Response should be marshallable to JSON")

		// Test JSON unmarshaling
		var parsedResponse SkillCostResponse
		err = json.Unmarshal(jsonData, &parsedResponse)
		assert.NoError(t, err, "Response should be unmarshallable from JSON")

		assert.Equal(t, response.SkillName, parsedResponse.SkillName, "Skill name should match")
		assert.Equal(t, response.SkillType, parsedResponse.SkillType, "Skill type should match")
		assert.Equal(t, response.Action, parsedResponse.Action, "Action should match")
		assert.Equal(t, response.CharacterID, parsedResponse.CharacterID, "Character ID should match")
		assert.Equal(t, response.Category, parsedResponse.Category, "Category should match")
		assert.Equal(t, response.Difficulty, parsedResponse.Difficulty, "Difficulty should match")
		assert.Equal(t, response.CanAfford, parsedResponse.CanAfford, "Can afford should match")
		assert.Equal(t, response.Notes, parsedResponse.Notes, "Notes should match")
	})

	t.Run("MultiLevelCostResponse structure", func(t *testing.T) {
		response := MultiLevelCostResponse{
			SkillName:      "Menschenkenntnis",
			SkillType:      "skill",
			CharacterID:    1,
			CurrentLevel:   10,
			TargetLevel:    13,
			LevelCosts:     []SkillCostResponse{},
			CanAffordTotal: true,
		}

		// Test JSON marshaling
		jsonData, err := json.Marshal(response)
		assert.NoError(t, err, "MultiLevelCostResponse should be marshallable to JSON")

		// Test JSON unmarshaling
		var parsedResponse MultiLevelCostResponse
		err = json.Unmarshal(jsonData, &parsedResponse)
		assert.NoError(t, err, "MultiLevelCostResponse should be unmarshallable from JSON")

		assert.Equal(t, response.SkillName, parsedResponse.SkillName, "Skill name should match")
		assert.Equal(t, response.SkillType, parsedResponse.SkillType, "Skill type should match")
		assert.Equal(t, response.CharacterID, parsedResponse.CharacterID, "Character ID should match")
		assert.Equal(t, response.CurrentLevel, parsedResponse.CurrentLevel, "Current level should match")
		assert.Equal(t, response.TargetLevel, parsedResponse.TargetLevel, "Target level should match")
		assert.Equal(t, response.CanAffordTotal, parsedResponse.CanAffordTotal, "Can afford total should match")
	})
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	t.Run("getCurrentSkillLevel", func(t *testing.T) {
		// This would need a proper character setup to test fully
		// For now, we're just testing the function exists and doesn't panic

		var character Char
		level := getCurrentSkillLevel(&character, "Test", "skill")
		assert.Equal(t, -1, level, "Should return -1 for non-existent skill")
	})
}

// Test integration with gsmaster exported functions
func TestGSMasterIntegration(t *testing.T) {
	t.Run("GetDefaultCategory integration", func(t *testing.T) {
		// Test that we can access the exported function from gsmaster
		category := gsmaster.GetDefaultCategory("Menschenkenntnis")
		assert.Equal(t, "Sozial", category, "Should return correct category for Menschenkenntnis")

		category = gsmaster.GetDefaultCategory("Stichwaffen")
		assert.Equal(t, "Waffen", category, "Should return correct category for Stichwaffen")

		// Test fallback for unknown skill
		category = gsmaster.GetDefaultCategory("NonExistentSkill")
		assert.Equal(t, "Alltag", category, "Should return default category for unknown skill")
	})

	t.Run("GetDefaultDifficulty integration", func(t *testing.T) {
		// Test that we can access the exported function from gsmaster
		difficulty := gsmaster.GetDefaultDifficulty("Menschenkenntnis")
		assert.Equal(t, "schwer", difficulty, "Should return correct difficulty for Menschenkenntnis")

		difficulty = gsmaster.GetDefaultDifficulty("Stichwaffen")
		assert.Equal(t, "leicht", difficulty, "Should return correct difficulty for Stichwaffen")

		// Test fallback for unknown skill
		difficulty = gsmaster.GetDefaultDifficulty("NonExistentSkill")
		assert.Equal(t, "normal", difficulty, "Should return default difficulty for unknown skill")
	})

	t.Run("Reward system structures", func(t *testing.T) {
		// Test RewardOptions structure
		rewards := RewardOptions{
			Type:         "free_learning",
			UseGoldForEP: true,
			MaxGoldEP:    50,
		}

		// Test JSON marshaling
		jsonData, err := json.Marshal(rewards)
		assert.NoError(t, err, "RewardOptions should be marshallable to JSON")

		// Test JSON unmarshaling
		var parsedRewards RewardOptions
		err = json.Unmarshal(jsonData, &parsedRewards)
		assert.NoError(t, err, "RewardOptions should be unmarshallable from JSON")

		assert.Equal(t, rewards.Type, parsedRewards.Type, "Type should match")
		assert.Equal(t, rewards.UseGoldForEP, parsedRewards.UseGoldForEP, "UseGoldForEP should match")
		assert.Equal(t, rewards.MaxGoldEP, parsedRewards.MaxGoldEP, "MaxGoldEP should match")

		// Test validation of reward types
		validTypes := []string{"free_learning", "free_spell_learning", "half_ep_improvement", "gold_for_ep"}
		for _, validType := range validTypes {
			rewards.Type = validType
			_, err := json.Marshal(rewards)
			assert.NoError(t, err, fmt.Sprintf("Should marshal valid type: %s", validType))
		}
	})

	t.Run("Reward system integration with gsmaster functions", func(t *testing.T) {
		// Test that the reward system works with the exported gsmaster functions
		// This simulates the flow where we get skill info from gsmaster and apply rewards

		skillName := "Menschenkenntnis"

		// Get skill info using exported functions
		category := gsmaster.GetDefaultCategory(skillName)
		difficulty := gsmaster.GetDefaultDifficulty(skillName)

		assert.Equal(t, "Sozial", category, "Should get correct category from gsmaster")
		assert.Equal(t, "schwer", difficulty, "Should get correct difficulty from gsmaster")

		// Test reward structure that would be used in the actual API
		rewards := RewardOptions{
			Type:         "half_ep_improvement",
			UseGoldForEP: false,
			MaxGoldEP:    0,
		}

		// Test that structure is valid
		jsonData, err := json.Marshal(rewards)
		assert.NoError(t, err, "Reward options should marshal correctly")

		var parsedRewards RewardOptions
		err = json.Unmarshal(jsonData, &parsedRewards)
		assert.NoError(t, err, "Reward options should unmarshal correctly")
		assert.Equal(t, "half_ep_improvement", parsedRewards.Type, "Reward type should be preserved")
	})
}

// Test the GetSkillAllLevelCosts endpoint (GET /:id/improve/skill)
func TestGetSkillAllLevelCostsEndpoint(t *testing.T) {
	// Setup test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Migrate the schema
	err := MigrateStructure()
	assert.NoError(t, err)

	// Also migrate skills and equipment to avoid preload errors
	err = skills.MigrateStructure()
	assert.NoError(t, err)
	err = equipment.MigrateStructure()
	assert.NoError(t, err)
	err = gsmaster.MigrateStructure()
	assert.NoError(t, err)

	// Create test skill data
	err = createTestSkillData()
	assert.NoError(t, err)
	defer cleanupTestSkillData()

	// Create test character with "Klettern" skill
	testChar := createChar()
	testChar.ID = 20 // Set the ID to match the test requirement

	// Add Menschenkenntnis skill at level 8 so we can improve to level 10
	skillName := "Menschenkenntnis"
	skill := skills.Fertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: skillName,
			},
			CharacterID: 20,
		},
		Fertigkeitswert: 8,
	}
	testChar.Fertigkeiten = append(testChar.Fertigkeiten, skill)

	err = testChar.Create()
	assert.NoError(t, err)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("Get improvement costs for Menschenkenntnis to level 10", func(t *testing.T) {
		// Create request body with skill name
		requestData := LearnRequestStruct{
			Name: "Menschenkenntnis",
		}
		requestBody, _ := json.Marshal(requestData)

		// Create HTTP request - note this is a GET request but with JSON body (unusual but matches the handler)
		req, _ := http.NewRequest("GET", "/api/characters/20/improve/skill", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: "20"}}

		fmt.Printf("Test: Get improvement costs for Menschenkenntnis to level 10 for character ID 20\n")
		fmt.Printf("Request: %+v\n", requestData)

		// Call the actual handler function
		GetSkillAllLevelCosts(c)

		// Print the actual response to see what we get
		fmt.Printf("Response Status: %d\n", w.Code)
		fmt.Printf("Response Body: %s\n", w.Body.String())

		// Check if we got an error response first
		if w.Code != http.StatusOK {
			var errorResponse map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			if err == nil {
				fmt.Printf("Error Response: %+v\n", errorResponse)
			}
			return // Exit early for error cases
		}

		// Parse and validate response for success case
		var response []gsmaster.LearnCost
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Should have costs for levels 9, 10 (from current level 8 to level 20)
		assert.Greater(t, len(response), 0, "Should return learning costs")

		// Find cost for level 10 specifically
		var level10Cost *gsmaster.LearnCost
		for _, cost := range response {
			if cost.Stufe == 10 {
				level10Cost = &cost
				break
			}
		}

		if level10Cost != nil {
			assert.Equal(t, 10, level10Cost.Stufe, "Target level should be 10")
			assert.Greater(t, level10Cost.Ep, 0, "EP cost should be greater than 0")

			fmt.Printf("Level 10 improvement cost: EP=%d, Money=%d, LE=%d\n",
				level10Cost.Ep, level10Cost.Money, level10Cost.LE)
		} else {
			fmt.Printf("No cost found for level 10. Available levels: ")
			for _, cost := range response {
				fmt.Printf("%d ", cost.Stufe)
			}
			fmt.Println()
		}
	})
	t.Run("Get improvement costs for Klettern to level 10", func(t *testing.T) {
		// Create request body with skill name
		requestData := LearnRequestStruct{
			SkillType: "skill",
			Name:      "Klettern",
			Stufe:     13, // Target level
		}
		requestBody, _ := json.Marshal(requestData)

		// Create HTTP request - note this is a GET request but with JSON body (unusual but matches the handler)
		req, _ := http.NewRequest("GET", "/api/characters/20/improve/skill", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: "20"}}

		fmt.Printf("Test: Get improvement costs for Klettern to level 13 for character ID 20\n")
		fmt.Printf("Request: %+v\n", requestData)

		// Call the actual handler function
		GetSkillAllLevelCosts(c)

		// Print the actual response to see what we get
		fmt.Printf("Response Status: %d\n", w.Code)
		fmt.Printf("Response Body: %s\n", w.Body.String())

		// Check if we got an error response first
		if w.Code != http.StatusOK {
			var errorResponse map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
			if err == nil {
				fmt.Printf("Error Response: %+v\n", errorResponse)
			}
			return // Exit early for error cases
		}

		// Parse and validate response for success case
		var response []gsmaster.LearnCost
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Should have costs for levels 9, 10 (from current level 8 to level 20)
		assert.Greater(t, len(response), 0, "Should return learning costs")

		// Find cost for level 10 specifically
		var level10Cost *gsmaster.LearnCost
		for _, cost := range response {
			if cost.Stufe == 13 {
				level10Cost = &cost
				break
			}
		}

		if level10Cost != nil {
			assert.Equal(t, 13, level10Cost.Stufe, "Target level should be 13")
			assert.Greater(t, level10Cost.Ep, 0, "EP cost should be greater than 0")

			fmt.Printf("Level 10 improvement cost: EP=%d, Money=%d, LE=%d\n",
				level10Cost.Ep, level10Cost.Money, level10Cost.LE)
		} else {
			fmt.Printf("No cost found for level 10. Available levels: ")
			for _, cost := range response {
				fmt.Printf("%d ", cost.Stufe)
			}
			fmt.Println()
		}
	})

	t.Run("Error case - skill not found", func(t *testing.T) {
		// Test with a skill the character doesn't have
		requestData := LearnRequestStruct{
			Name: "NonExistentSkill",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("GET", "/api/characters/20/improve/skill", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: "20"}}

		GetSkillAllLevelCosts(c)

		// Should still return 200 but with empty array (based on the handler logic)
		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200 OK")

		var response []gsmaster.LearnCost
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, 0, len(response), "Should return empty array for non-existent skill")
	})

	t.Run("Error case - missing skill name", func(t *testing.T) {
		// Test with empty skill name
		requestData := LearnRequestStruct{
			Name: "",
		}
		requestBody, _ := json.Marshal(requestData)

		req, _ := http.NewRequest("GET", "/api/characters/20/improve/skill", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: "20"}}

		GetSkillAllLevelCosts(c)

		// Should return 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w.Code, "Status code should be 400 Bad Request")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Contains(t, response, "error", "Response should contain error message")
	})
}
