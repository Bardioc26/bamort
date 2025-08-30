package character

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

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

		var character models.Char
		level := getCurrentSkillLevel(&character, "Test", "skill")
		assert.Equal(t, -1, level, "Should return -1 for non-existent skill")
	})
}

// Test GetLernCost endpoint specifically with gsmaster.LernCostRequest structure

// Test GetLernCost endpoint specifically with gsmaster.LernCostRequest structure
func TestGetLernCostEndpointNewSystem(t *testing.T) {
	// Setup test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Migrate the schema
	err := models.MigrateStructure()
	assert.NoError(t, err)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("GetLernCost with Athletik for Krieger character", func(t *testing.T) {
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
		req, _ := http.NewRequest("POST", "/api/characters/lerncost", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: "20"}}

		fmt.Printf("Test: GetLernCost for Athletik improvement for Krieger character ID 20\n")
		fmt.Printf("Request: CharId=%d, SkillName=%s, CurrentLevel=%d, TargetLevel=%d\n",
			requestData.CharId, requestData.Name, requestData.CurrentLevel, requestData.TargetLevel)

		// Call the actual handler function
		GetLernCostNewSystem(c)

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
			assert.Fail(t, "Expected successful response but got error: %s", w.Body.String())
			return
		}

		// Parse and validate response for success case
		var response []gsmaster.SkillCostResultNew
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON array of SkillCostResultNew")

		// Should have costs for levels 10, 11, 12, ... up to 18 (from current level 9)
		assert.Greater(t, len(response), 0, "Should return learning costs for multiple levels")
		assert.LessOrEqual(t, len(response), 9, "Should not return more than 9 levels (10-18)")

		// Validate the first entry (level 10)
		if len(response) > 0 {
			firstResult := response[0]
			assert.Equal(t, "20", firstResult.CharacterID, "Character ID should match")
			assert.Equal(t, "Athletik", firstResult.SkillName, "Skill name should match")
			assert.Equal(t, 10, firstResult.TargetLevel, "First target level should be 10")

			// Character class should be "Kr" (abbreviation for "Krieger")
			assert.Equal(t, "Kr", firstResult.CharacterClass, "Character class should be abbreviated to 'Kr'")

			// Should have valid costs
			assert.Greater(t, firstResult.EP, 0, "EP cost should be greater than 0")
			assert.GreaterOrEqual(t, firstResult.GoldCost, 0, "Gold cost should be 0 or greater")
			assert.Equal(t, firstResult.EP, 20, "EP cost should be 20")
			assert.Equal(t, firstResult.GoldCost, 40, "Gold cost should be 40")

			fmt.Printf("Level %d cost: EP=%d, GoldCost=%d, LE=%d\n", firstResult.TargetLevel,
				firstResult.EP, firstResult.GoldCost, firstResult.LE)
			fmt.Printf("Category=%s, Difficulty=%s\n",
				firstResult.Category, firstResult.Difficulty)
		}

		// Find cost for level 12 specifically to test mid-range
		var level12Cost *gsmaster.SkillCostResultNew
		for i := range response {
			if response[i].TargetLevel == 12 {
				level12Cost = &response[i]
				break
			}
		}

		if level12Cost != nil {
			assert.Equal(t, 12, level12Cost.TargetLevel, "Target level should be 12")
			assert.Greater(t, level12Cost.EP, 0, "EP cost should be greater than 0 for level 12")

			fmt.Printf("Level 12 cost: EP=%d, GoldCost=%d, LE=%d\n",
				level12Cost.EP, level12Cost.GoldCost, level12Cost.LE)
		} else {
			fmt.Printf("No cost found for level 12. Available levels: ")
			for _, cost := range response {
				fmt.Printf("%d ", cost.TargetLevel)
			}
			fmt.Println()
		}

		// Verify all target levels are sequential and start from current level + 1
		expectedLevel := 10 // Current level 9 + 1
		for _, cost := range response {
			assert.Equal(t, expectedLevel, cost.TargetLevel,
				"Target levels should be sequential starting from %d", expectedLevel)
			assert.Equal(t, "Athletik", cost.SkillName, "All entries should have correct skill name")
			assert.Equal(t, "Kr", cost.CharacterClass, "All entries should have correct character class")
			expectedLevel++
		}
	})

}
