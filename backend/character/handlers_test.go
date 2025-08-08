package character

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bamort/database"
	"bamort/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestImproveSkillHandler(t *testing.T) {
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
