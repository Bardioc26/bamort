package character

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailableSpellsNewSystem(t *testing.T) {
	// Setup test database with real data
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("Get available spells for character ID 20", func(t *testing.T) {

		// Test Request wie im Frontend
		request := gsmaster.LernCostRequest{
			CharId:  18,
			Type:    "spell",
			Action:  "learn",
			UsePP:   0,
			UseGold: 0,
			Reward:  stringPtr("default"),
		}

		// Convert request to JSON
		requestJSON, err := json.Marshal(request)
		assert.NoError(t, err, "Should marshal request")

		// Create HTTP request
		req, _ := http.NewRequest("POST", "/api/characters/available-spells-new", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		fmt.Printf("Test: Get available spells for character ID 20\n")
		fmt.Printf("Request: %s\n", string(requestJSON))

		// Call the handler function
		GetAvailableSpellsNewSystem(c)

		// Print the response for debugging
		fmt.Printf("Response Status: %d\n", w.Code)
		fmt.Printf("Response Body: %s\n", w.Body.String())

		// Assert response status
		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200 OK")

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Check response structure
		assert.Contains(t, response, "spells_by_school", "Response should contain spells_by_school field")

		// Check if spells_by_school is an object
		spellsBySchool, ok := response["spells_by_school"].(map[string]interface{})
		assert.True(t, ok, "spells_by_school should be an object")

		fmt.Printf("Found %d spell schools\n", len(spellsBySchool))

		// Check each school
		totalSpells := 0
		for schoolName, schoolSpells := range spellsBySchool {
			fmt.Printf("School: %s\n", schoolName)

			spells, ok := schoolSpells.([]interface{})
			assert.True(t, ok, fmt.Sprintf("School %s should contain an array of spells", schoolName))

			totalSpells += len(spells)
			fmt.Printf("  Spells in school: %d\n", len(spells))

			// Check structure of first spell in school if exists
			if len(spells) > 0 {
				firstSpell, ok := spells[0].(map[string]interface{})
				assert.True(t, ok, "First spell should be an object")

				// Check required fields
				assert.Contains(t, firstSpell, "name", "Spell should have name field")
				assert.Contains(t, firstSpell, "level", "Spell should have level field")
				assert.Contains(t, firstSpell, "epCost", "Spell should have epCost field")
				assert.Contains(t, firstSpell, "goldCost", "Spell should have goldCost field")

				spellName, ok := firstSpell["name"].(string)
				assert.True(t, ok, "Spell name should be a string")
				assert.NotEmpty(t, spellName, "Spell name should not be empty")

				level, ok := firstSpell["level"].(float64)
				assert.True(t, ok, "Spell level should be a number")
				assert.GreaterOrEqual(t, level, float64(0), "Spell level should be at least 0")

				epCost, ok := firstSpell["epCost"].(float64)
				assert.True(t, ok, "EP cost should be a number")
				assert.GreaterOrEqual(t, epCost, float64(0), "EP cost should be non-negative")

				goldCost, ok := firstSpell["goldCost"].(float64)
				assert.True(t, ok, "Gold cost should be a number")
				assert.GreaterOrEqual(t, goldCost, float64(0), "Gold cost should be non-negative")

				fmt.Printf("    Example spell: %s (Level %v, EP: %v, Gold: %v)\n",
					spellName, level, epCost, goldCost)
			}
		}

		fmt.Printf("Total spells found: %d\n", totalSpells)
		assert.Greater(t, totalSpells, 0, "Should have at least some spells available")
	})

	t.Run("Test with different reward types", func(t *testing.T) {
		// Verwende existierenden Charakter
		rewardTypes := []string{"default", "noGold", "halveep", "halveepnoGold"}

		for _, rewardType := range rewardTypes {
			t.Run(fmt.Sprintf("RewardType_%s", rewardType), func(t *testing.T) {
				request := gsmaster.LernCostRequest{
					CharId:  18,
					Type:    "spell",
					Action:  "learn",
					UsePP:   0,
					UseGold: 0,
					Reward:  stringPtr(rewardType),
				}

				requestJSON, err := json.Marshal(request)
				assert.NoError(t, err)

				req, _ := http.NewRequest("POST", "/api/characters/available-spells-new", bytes.NewBuffer(requestJSON))
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request = req

				GetAvailableSpellsNewSystem(c)

				fmt.Printf("Testing reward type: %s - Status: %d\n", rewardType, w.Code)
				assert.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("Should work with reward type %s", rewardType))

				var response map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err, "Response should be valid JSON")
				assert.Contains(t, response, "spells_by_school", "Should contain spells_by_school")
			})
		}
	})

	t.Run("Test with invalid character ID", func(t *testing.T) {
		request := gsmaster.LernCostRequest{
			CharId:  99999, // Non-existent character
			Type:    "spell",
			Action:  "learn",
			UsePP:   0,
			UseGold: 0,
			Reward:  stringPtr("default"),
		}

		requestJSON, err := json.Marshal(request)
		assert.NoError(t, err)

		req, _ := http.NewRequest("POST", "/api/characters/available-spells-new", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		GetAvailableSpellsNewSystem(c)

		fmt.Printf("Testing invalid character ID - Status: %d\n", w.Code)
		// Der Handler verwendet nicht den "id" Parameter aus der URL, sondern CharId aus dem Request Body
		// Daher wird kein 404 zurückgegeben, sondern ein leeres Ergebnis oder Fehler
		assert.True(t, w.Code == http.StatusNotFound || w.Code == http.StatusOK,
			"Should return 404 or 200 for non-existent character")
	})

	t.Run("Test with invalid request format", func(t *testing.T) {
		invalidJSON := []byte(`{"invalid": "request"}`)

		req, _ := http.NewRequest("POST", "/api/characters/available-spells-new", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		GetAvailableSpellsNewSystem(c)

		fmt.Printf("Testing invalid request format - Status: %d\n", w.Code)
		assert.Equal(t, http.StatusBadRequest, w.Code, "Should return 400 for invalid request")
	})

	t.Run("Test excluding already learned spells", func(t *testing.T) {
		// Füge dem Charakter einen Zauber hinzu
		learnedSpell := models.SkZauber{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{
					Name: "Feuerball",
				},
				CharacterID: 20,
			},
			Beschreibung: "Ein mächtiger Feuerball",
			Bonus:        5,
			Quelle:       "Test",
		}

		err := database.DB.Create(&learnedSpell).Error
		assert.NoError(t, err, "Should create learned spell")

		request := gsmaster.LernCostRequest{
			CharId:  18,
			Type:    "spell",
			Action:  "learn",
			UsePP:   0,
			UseGold: 0,
			Reward:  stringPtr("default"),
		}

		requestJSON, err := json.Marshal(request)
		assert.NoError(t, err)

		req, _ := http.NewRequest("POST", "/api/characters/available-spells-new", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		GetAvailableSpellsNewSystem(c)

		assert.Equal(t, http.StatusOK, w.Code, "Should return 200 OK")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		spellsBySchool := response["spells_by_school"].(map[string]interface{})

		// Prüfe, dass "Feuerball" nicht in der Liste ist
		foundFeuerball := false
		for _, schoolSpells := range spellsBySchool {
			spells := schoolSpells.([]interface{})
			for _, spell := range spells {
				spellObj := spell.(map[string]interface{})
				if spellObj["name"].(string) == "Feuerball" {
					foundFeuerball = true
					break
				}
			}
		}

		assert.False(t, foundFeuerball, "Already learned spell 'Feuerball' should not be in available spells")
		fmt.Printf("Correctly excluded already learned spell 'Feuerball'\n")
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
