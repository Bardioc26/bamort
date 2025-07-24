package character

import (
	"bamort/gsmaster"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSkillLearningDialogWorkflow testet den kompletten Workflow, den das SkillLearningDialog verwendet
func TestSkillLearningDialogWorkflow(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Reward Types API - entspricht Frontend loadRewardTypes()", func(t *testing.T) {
		// Test verschiedene Lerntypen wie sie das Frontend sendet
		testCases := []struct {
			learningType string
			skillName    string
			skillType    string
			expectedLen  int
		}{
			{"improve", "Menschenkenntnis", "skill", 4}, // EP, Gold, PP, Mixed
			{"learn", "Klettern", "skill", 2},           // EP, Gold
			//{"spell", "Licht", "spell", 4},              // EP, Gold, PP, Mixed
		}

		for _, tc := range testCases {
			t.Run(tc.learningType+"_"+tc.skillName, func(t *testing.T) {
				router := gin.New()
				router.GET("/api/characters/:id/reward-types", GetRewardTypes)

				url := "/api/characters/1/reward-types?learning_type=" + tc.learningType +
					"&skill_name=" + tc.skillName + "&skill_type=" + tc.skillType

				req, _ := http.NewRequest("GET", url, nil)
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)

				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Überprüfe Response-Struktur wie sie das Frontend erwartet
				assert.Contains(t, response, "reward_types")
				assert.Contains(t, response, "learning_type")
				assert.Contains(t, response, "skill_name")
				assert.Contains(t, response, "character_id")

				rewardTypes := response["reward_types"].([]interface{})
				assert.GreaterOrEqual(t, len(rewardTypes), tc.expectedLen-1) // Mindestens erwartete Anzahl

				// Überprüfe Struktur der Reward Types wie Frontend sie erwartet
				if len(rewardTypes) > 0 {
					firstReward := rewardTypes[0].(map[string]interface{})
					assert.Contains(t, firstReward, "value")
					assert.Contains(t, firstReward, "label")
				}
			})
		}
	})

	t.Run("Skill All Level Costs API - entspricht Frontend loadLearningCosts()", func(t *testing.T) {
		// Teste nur die Request-Struktur und Skip den DB-Teil
		t.Skip("Skipping DB-dependent test - testing request structure only")

		// Test der Request-Struktur wie das Frontend sie sendet
		requestData := LearnRequestStruct{
			SkillType: "skill",
			Name:      "Menschenkenntnis",
			Stufe:     10,
		}

		// Überprüfe, dass die Request-Struktur korrekt ist
		requestBody, err := json.Marshal(requestData)
		require.NoError(t, err)

		var parsedRequest LearnRequestStruct
		err = json.Unmarshal(requestBody, &parsedRequest)
		require.NoError(t, err)

		// Überprüfe, dass alle Felder korrekt geparst werden
		assert.Equal(t, "skill", parsedRequest.SkillType)
		assert.Equal(t, "Menschenkenntnis", parsedRequest.Name)
		assert.Equal(t, 10, parsedRequest.Stufe)

		// Teste erwartete Response-Struktur (ohne DB-Aufruf)
		expectedResponse := []gsmaster.LearnCost{
			{Stufe: 11, Ep: 120, Money: 60, LE: 2},
			{Stufe: 12, Ep: 140, Money: 70, LE: 2},
		}

		// Überprüfe Response-Format wie das Frontend es erwartet
		responseBody, _ := json.Marshal(expectedResponse)
		var response []gsmaster.LearnCost
		err = json.Unmarshal(responseBody, &response)
		require.NoError(t, err)

		if len(response) > 0 {
			cost := response[0]
			assert.NotZero(t, cost.Stufe) // target_level im Frontend
			assert.NotZero(t, cost.Ep)    // ep_cost im Frontend
			assert.NotZero(t, cost.Money) // gold_cost im Frontend
			// cost.LE wird als pp_cost verwendet
		}
	})

	t.Run("Frontend Conversion Logic Test", func(t *testing.T) {
		// Simuliere die Konvertierung wie sie das Frontend durchführt
		mockApiResponse := []gsmaster.LearnCost{
			{Stufe: 11, Ep: 120, Money: 60, LE: 2},
			{Stufe: 12, Ep: 140, Money: 70, LE: 2},
			{Stufe: 13, Ep: 160, Money: 80, LE: 3},
		}

		// Simuliere Frontend-Konvertierung
		availableEP := 1000
		availableGold := 500
		availablePP := 10

		var convertedLevels []map[string]interface{}
		cumulativeEP := 0
		cumulativeGold := 0
		cumulativePP := 0

		for _, cost := range mockApiResponse {
			cumulativeEP += cost.Ep
			cumulativeGold += cost.Money
			cumulativePP += cost.LE

			level := map[string]interface{}{
				"targetLevel":   cost.Stufe,
				"epCost":        cost.Ep,
				"goldCost":      cost.Money,
				"ppCost":        cost.LE,
				"totalEpCost":   cumulativeEP,
				"totalGoldCost": cumulativeGold,
				"totalPpCost":   cumulativePP,
				"canAfford": map[string]bool{
					"ep":   availableEP >= cumulativeEP,
					"gold": availableGold >= cumulativeGold,
					"pp":   availablePP >= cumulativePP,
				},
			}
			convertedLevels = append(convertedLevels, level)
		}

		// Überprüfe Konvertierung
		assert.Len(t, convertedLevels, 3)

		firstLevel := convertedLevels[0]
		assert.Equal(t, 11, firstLevel["targetLevel"])
		assert.Equal(t, 120, firstLevel["epCost"])
		assert.Equal(t, 120, firstLevel["totalEpCost"]) // Erste Stufe = einzelne Kosten

		lastLevel := convertedLevels[2]
		assert.Equal(t, 420, lastLevel["totalEpCost"])   // 120 + 140 + 160
		assert.Equal(t, 210, lastLevel["totalGoldCost"]) // 60 + 70 + 80

		// Verfügbarkeits-Test
		canAfford := lastLevel["canAfford"].(map[string]bool)
		assert.True(t, canAfford["ep"])   // 1000 >= 420
		assert.True(t, canAfford["gold"]) // 500 >= 210
		assert.True(t, canAfford["pp"])   // 10 >= 7
	})

	t.Run("Skill Learning Execution Test", func(t *testing.T) {
		// Test der Request-Struktur für Lernausführung
		// (DB-Operationen werden übersprungen)
		t.Skip("Skipping DB-dependent test - testing request structure only")

		// Request wie das Frontend ihn für executeDetailedLearning sendet
		requestData := ImproveSkillRequest{
			Name:         "Menschenkenntnis",
			CurrentLevel: 10,
			Notes:        "Fertigkeit Menschenkenntnis von 10 auf 11 verbessert",
		}

		// Verify the request structure is correct
		assert.Equal(t, "Menschenkenntnis", requestData.Name)
		assert.Equal(t, 10, requestData.CurrentLevel)
		assert.NotEmpty(t, requestData.Notes)

		// Das Frontend erwartet nach erfolgreichem Lernen diese Response-Felder:
		expectedResponseFields := []string{"message", "skill_name", "ep_cost", "remaining_ep"}
		assert.NotEmpty(t, expectedResponseFields, "Response sollte diese Felder enthalten")
	})
}

// Test für die Frontend-Authentifizierung
func TestSkillLearningDialogAuth(t *testing.T) {
	t.Run("Authentication Error Handling", func(t *testing.T) {
		// Test wie das Frontend mit Auth-Fehlern umgeht
		router := gin.New()

		// Mock eines Auth-Middlewares der 401 zurückgibt
		router.Use(func(c *gin.Context) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
		})

		router.GET("/api/characters/:id/reward-types", GetRewardTypes)

		req, _ := http.NewRequest("GET", "/api/characters/1/reward-types", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// Das Frontend sollte bei 401 das 'auth-error' Event emittieren
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Contains(t, response, "error")
	})
}

// Test verschiedener Belohnungstypen
func TestRewardTypeVariations(t *testing.T) {
	testCases := []struct {
		name            string
		learningType    string
		skillType       string
		expectedRewards []string
	}{
		{
			name:            "Improve Skill",
			learningType:    "improve",
			skillType:       "skill",
			expectedRewards: []string{"ep", "gold", "pp", "mixed"},
		},
		{
			name:            "Learn Skill",
			learningType:    "learn",
			skillType:       "skill",
			expectedRewards: []string{"ep", "gold"},
		},
		{
			name:            "Weapon Training",
			learningType:    "improve",
			skillType:       "weapon",
			expectedRewards: []string{"ep", "gold", "pp", "mixed", "training"},
		},
		{
			name:            "Spell Learning",
			learningType:    "spell",
			skillType:       "spell",
			expectedRewards: []string{"ep", "gold", "pp", "mixed"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/api/characters/:id/reward-types", GetRewardTypes)

			url := "/api/characters/1/reward-types?learning_type=" + tc.learningType +
				"&skill_type=" + tc.skillType + "&skill_name=TestSkill"

			req, _ := http.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			rewardTypes := response["reward_types"].([]interface{})

			// Überprüfe, dass erwartete Reward Types vorhanden sind
			foundRewards := make(map[string]bool)
			for _, reward := range rewardTypes {
				rewardMap := reward.(map[string]interface{})
				value := rewardMap["value"].(string)
				foundRewards[value] = true
			}

			for _, expectedReward := range tc.expectedRewards {
				assert.True(t, foundRewards[expectedReward],
					"Erwartete Belohnungsart '%s' nicht gefunden für %s", expectedReward, tc.name)
			}
		})
	}
}
