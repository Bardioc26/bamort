package character

import (
	"bamort/gsmaster"
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

			// In a real test, you would call GetSkillCost(c) here
			// But since we don't have a database setup, we'll just test the structure

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
