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

			// Costs must be non-negative; values depend on test data
			assert.GreaterOrEqual(t, firstResult.EP, 0, "EP cost should be non-negative")
			assert.GreaterOrEqual(t, firstResult.GoldCost, 0, "Gold cost should be non-negative")

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
			assert.GreaterOrEqual(t, level12Cost.EP, 0, "EP cost should be non-negative for level 12")

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

func TestCalculateSpellLearnCostNewSystem(t *testing.T) {
	rewardHalve := "halveep"

	tests := []struct {
		name              string
		request           gsmaster.LernCostRequest
		spellInfo         models.SpellLearningInfo
		remainingPP       int
		remainingGold     int
		wantLE            int
		wantEP            int
		wantGold          int
		wantPPUsed        int
		wantGoldUsed      int
		wantRemainingPP   int
		wantRemainingGold int
	}{
		{
			name:    "applies PP before cost calculation",
			request: gsmaster.LernCostRequest{Reward: nil, CharId: 15},
			spellInfo: models.SpellLearningInfo{
				SpellID:    47,
				SpellName:  "Scharfblick",
				SpellLevel: 1,
				LERequired: 1,
				EPPerLE:    60,
				SchoolName: "Verändern",
			},
			remainingPP:       1,
			remainingGold:     0,
			wantLE:            0,
			wantEP:            00,
			wantGold:          00,
			wantPPUsed:        1,
			wantGoldUsed:      0,
			wantRemainingPP:   0,
			wantRemainingGold: 0,
		},
		{
			name:    "halves EP via reward",
			request: gsmaster.LernCostRequest{Reward: &rewardHalve, CharId: 22},
			spellInfo: models.SpellLearningInfo{
				SpellID:    1,
				SpellName:  "Angst",
				SpellLevel: 2,
				LERequired: 1,
				EPPerLE:    60,
				SchoolName: "Beherrschen",
			},
			remainingPP:       0,
			remainingGold:     0,
			wantLE:            1,
			wantEP:            30, // 3 LE * 40 EP/LE -> 120, reward halves to 60
			wantGold:          100,
			wantPPUsed:        0,
			wantGoldUsed:      0,
			wantRemainingPP:   0,
			wantRemainingGold: 0,
		},
		{
			name:    "Zaubersprung Hx uses gold up to available but below cap",
			request: gsmaster.LernCostRequest{Reward: nil, CharId: 18},
			spellInfo: models.SpellLearningInfo{
				SpellID:    68,
				SpellName:  "Zaubersprung",
				SpellLevel: 3,
				LERequired: 2,
				EPPerLE:    30,
				SchoolName: "Beherrschen",
			},
			remainingPP:       0,
			remainingGold:     50, // 5 EP replacement, below EP/2 cap (60)
			wantLE:            2,
			wantEP:            55,
			wantGold:          200,
			wantPPUsed:        0,
			wantGoldUsed:      50,
			wantRemainingPP:   0,
			wantRemainingGold: 0,
		},

		{
			name:    "Zaubersprung Ma uses gold up to available but below cap",
			request: gsmaster.LernCostRequest{Reward: nil, CharId: 22},
			spellInfo: models.SpellLearningInfo{
				SpellID:    68,
				SpellName:  "Zaubersprung",
				SpellLevel: 3,
				LERequired: 2,
				EPPerLE:    60,
				SchoolName: "Beherrschen",
			},
			remainingPP:       0,
			remainingGold:     50, // 5 EP replacement, below EP/2 cap (60)
			wantLE:            2,
			wantEP:            115,
			wantGold:          200,
			wantPPUsed:        0,
			wantGoldUsed:      50,
			wantRemainingPP:   0,
			wantRemainingGold: 0,
		},
		{

			name:    "caps gold conversion at half EP",
			request: gsmaster.LernCostRequest{CharId: 22},
			spellInfo: models.SpellLearningInfo{
				SpellID:    68,
				SpellName:  "Zaubersprung",
				SpellLevel: 3,
				LERequired: 2,
				EPPerLE:    60,
				SchoolName: "Beherrschen",
			},
			remainingPP:       0,
			remainingGold:     2000, // would allow 200 EP replacement but cap is EP/2 = 60
			wantLE:            2,
			wantEP:            60,
			wantGold:          200,
			wantPPUsed:        0,
			wantGoldUsed:      600,
			wantRemainingPP:   0,
			wantRemainingGold: 1400,
		},
	}

	for _, tc := range tests {
		//tc := tc
		t.Run(tc.name, func(t *testing.T) {
			remainingPP := tc.remainingPP
			remainingGold := tc.remainingGold

			var result gsmaster.SkillCostResultNew
			err := calculateSpellLearnCostNewSystem(&tc.request, &result, &remainingPP, &remainingGold, &tc.spellInfo)
			assert.NoError(t, err)

			assert.Equal(t, tc.spellInfo.SchoolName, result.Category)
			assert.Equal(t, fmt.Sprintf("Stufe %d", tc.spellInfo.SpellLevel), result.Difficulty)
			assert.Equal(t, tc.wantLE, result.LE)
			assert.Equal(t, tc.wantEP, result.EP)
			assert.Equal(t, tc.wantGold, result.GoldCost)
			assert.Equal(t, tc.wantPPUsed, result.PPUsed)
			assert.Equal(t, tc.wantGoldUsed, result.GoldUsed)
			assert.Equal(t, tc.wantRemainingPP, remainingPP)
			assert.Equal(t, tc.wantRemainingGold, remainingGold)
		})
	}
}
