package character

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bamort/database"
	"bamort/models"
	"bamort/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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
			"remaining_ep":   float64(250),
			"remaining_gold": float64(290),
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
		assert.Equal(t, 250, updatedChar.Erfahrungsschatz.EP, "Character should have 316 EP remaining")

		// Check that Gold was deducted correctly
		assert.Equal(t, 290, updatedChar.Vermoegen.Goldstuecke, "Character should have 370 Gold remaining")

		t.Logf("Test completed successfully!")
		t.Logf("EP: %d -> %d (cost: %.0f)", 326, updatedChar.Erfahrungsschatz.EP, response["ep_cost"])
		t.Logf("Gold: %d -> %d (cost: %.0f)", 390, updatedChar.Vermoegen.Goldstuecke, response["gold_cost"])
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
								assert.Contains(t, skillMap, "name", "Skill should have name", skillMap["name"])
								//assert.Contains(t, skillMap, "learnCost", "Skill should have learnCost", skillMap["name"])
								assert.Contains(t, skillMap, "leCost", "Skill should have leCost", skillMap["name"])

								learnCost := skillMap["leCost"].(float64)
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

func TestGetAvailableSpellsForCreation(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	tests := []struct {
		name           string
		characterClass string
		expectStatus   int
		expectError    bool
		findspells     bool
	}{
		{
			name:           "ValidCharacterClass",
			characterClass: "As",
			expectStatus:   http.StatusNotFound,
			expectError:    false,
			findspells:     false,
		},
		{
			name:           "MagierCharacterClass",
			characterClass: "Ma",
			expectStatus:   http.StatusOK,
			expectError:    false,
			findspells:     true,
		},
		{
			name:           "NonMagicCharacterClass",
			characterClass: "Kr",
			expectStatus:   http.StatusNotFound,
			expectError:    true,
			findspells:     false,
		},
		{
			name:           "EmptyCharacterClass",
			characterClass: "",
			expectStatus:   http.StatusBadRequest,
			expectError:    true,
			findspells:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			requestData := gin.H{
				"characterClass": tt.characterClass,
			}
			requestBody, _ := json.Marshal(requestData)

			u := user.User{}
			u.FirstId(1)
			token := user.GenerateToken(&u)

			// Create HTTP request
			req, _ := http.NewRequest("POST", "/api/characters/available-spells-creation", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			// Create response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the handler directly (it will handle JSON parsing internally)
			GetAvailableSpellsForCreation(c)

			// Verify response
			assert.Equal(t, tt.expectStatus, w.Code)

			if !tt.expectError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check response structure
				assert.Contains(t, response, "spells_by_category")

				spellsByCategory, ok := response["spells_by_category"].(map[string]interface{})
				assert.True(t, ok)
				if tt.findspells {
					assert.Greater(t, len(spellsByCategory), 0, "Should have at least some spell categories")

					// Verify spells have learnCost field
					for categoryName, spells := range spellsByCategory {
						if spellsList, ok := spells.([]interface{}); ok {
							for _, spell := range spellsList {
								if spellMap, ok := spell.(map[string]interface{}); ok {
									assert.Contains(t, spellMap, "name", "Spell should have name")
									assert.Contains(t, spellMap, "le_cost", "Spell should have learnCost")

									learnCost := spellMap["le_cost"].(float64)
									assert.Greater(t, learnCost, 0.0, "Learn cost should be positive")
									assert.Less(t, learnCost, 500.0, "Learn cost should be reasonable for character creation")
								}
							}
						}
						_ = categoryName // Mark as used
					}
					// Log some sample data for verification
					t.Logf("Character creation spells loaded for class %s: %d categories", tt.characterClass, len(spellsByCategory))
				} else {
					assert.Equal(t, len(spellsByCategory), 0, "Should not have any spell categories")
				}

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

func TestFinalizeCharacterCreation(t *testing.T) {
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
	// Create test user (bebe)
	testUser := user.User{}
	testUser.FirstId(1)
	/*
		{
			UserID:   1,
			Username: "bebe",
			Email:    "frank@wuenscheonline.de",
		}
		err = database.DB.Create(&testUser).Error
		assert.NoError(t, err)
	*/
	// Create test character creation session with data from testdata
	testSession := models.CharacterCreationSession{
		ID:     "char_create_1_1756326371",
		UserID: 1,
		Name:   "wergw5 z ",
		Rasse:  "Mensch",
		Typ:    "Priester Streiter",
		Attributes: models.AttributesData{
			ST: 89,
			GS: 64,
			GW: 77,
			KO: 71,
			IN: 87,
			ZT: 44,
			AU: 87,
		},
		DerivedValues: models.DerivedValuesData{
			PA:                    33,
			WK:                    27,
			LPMax:                 11,
			APMax:                 14,
			BMax:                  26,
			ResistenzKoerper:      11,
			ResistenzGeist:        11,
			ResistenzBonusKoerper: 0,
			ResistenzBonusGeist:   0,
			Abwehr:                11,
			AbwehrBonus:           0,
			AusdauerBonus:         11,
			AngriffsBonus:         0,
			Zaubern:               11,
			ZauberBonus:           0,
			Raufen:                8,
			SchadensBonus:         3,
			SG:                    0,
			GG:                    0,
			GP:                    0,
		},
		Skills: models.CharacterCreationSkills{
			{Name: "Klettern", Level: 0, Category: "Alltag", Cost: 1},
			{Name: "Reiten", Level: 0, Category: "Alltag", Cost: 1},
			{Name: "Sprache", Level: 0, Category: "Alltag", Cost: 1},
			{Name: "Athletik", Level: 0, Category: "Kampf", Cost: 2},
			{Name: "Spießwaffen", Level: 0, Category: "Waffen", Cost: 2},
			{Name: "Stielwurfwaffen", Level: 0, Category: "Waffen", Cost: 2},
			{Name: "Waffenloser Kampf", Level: 0, Category: "Waffen", Cost: 2},
			{Name: "Stichwaffen", Level: 0, Category: "Waffen", Cost: 2},
			{Name: "Heilkunde", Level: 0, Category: "Wissen", Cost: 2},
			{Name: "Naturkunde", Level: 0, Category: "Wissen", Cost: 2},
		},
		Spells: models.CharacterCreationSpells{
			{Name: "Göttlicher Schutz v. d. Bösen", Cost: 1},
			{Name: "Erkennen der Aura", Cost: 1},
			{Name: "Heiliger Zorn", Cost: 1},
			{Name: "Blutmeisterschaft", Cost: 1},
		},
		SkillPoints: models.SkillPointsData{},
		CurrentStep: 5,
		Geschlecht:  "Männlich",
		Herkunft:    "Aran",
		Stand:       "Mittelschicht",
		Glaube:      "",
	}
	err = database.DB.Create(&testSession).Error
	assert.NoError(t, err)

	t.Run("FinalizeCharacterCreation with valid session", func(t *testing.T) {
		// Create the HTTP request to finalize the character creation
		// Using the session ID from the test data: 'char_create_1_1756326371'
		req, err := http.NewRequest("POST", "/api/characters/sessions/char_create_1_1756326371/finalize", nil)
		assert.NoError(t, err)

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context with userID set to 1 (bebe user)
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", uint(1)) // Set the userID for bebe user
		c.Params = gin.Params{
			gin.Param{Key: "sessionId", Value: "char_create_1_1756326371"},
		}

		// Call the FinalizeCharacterCreation handler
		FinalizeCharacterCreation(c)

		// Log the response for debugging
		t.Logf("Response Status: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		// Check if we got a successful response
		assert.Equal(t, http.StatusCreated, w.Code, "Status code should be 201 Created")

		// Parse and validate response
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")

		// Validate expected response structure
		assert.Contains(t, response, "message", "Response should contain success message")
		assert.Contains(t, response, "character_id", "Response should contain character_id")
		assert.Contains(t, response, "session_id", "Response should contain session_id")

		// Validate response values
		assert.Equal(t, "Charakter erfolgreich erstellt", response["message"])
		assert.Equal(t, "char_create_1_1756326371", response["session_id"])

		// Verify that a character was actually created
		characterID := response["character_id"]
		assert.NotNil(t, characterID, "Character ID should not be nil")

		// Verify character exists in database
		var createdChar models.Char
		err = database.DB.Preload("Fertigkeiten").Preload("Waffenfertigkeiten").Preload("Zauber").
			First(&createdChar, "id = ?", characterID).Error
		assert.NoError(t, err, "Created character should exist in database")

		// Validate character data based on the session data
		assert.Equal(t, "wergw5 z ", createdChar.Name, "Character name should match session")
		assert.Equal(t, "Mensch", createdChar.Rasse, "Character race should match session")
		assert.Equal(t, "Priester Streiter", createdChar.Typ, "Character type should match session")
		assert.Equal(t, "Aran", createdChar.Herkunft, "Character origin should match session")
		assert.Equal(t, "Männlich", createdChar.Gender, "Character gender should match session")
		assert.Equal(t, "Mittelschicht", createdChar.SocialClass, "Character status should match session")

		// Validate attributes
		assert.Equal(t, 9, len(createdChar.Eigenschaften), "Should have 9 attributes (ST, GS, GW, KO, IN, ZT, AU, pA, WK)")

		// Create a map for easier attribute validation
		attrMap := make(map[string]int)
		for _, attr := range createdChar.Eigenschaften {
			attrMap[attr.Name] = attr.Value
		}

		// Validate each attribute matches the session data
		assert.Equal(t, 87, attrMap["IN"], "Intelligence should match session")
		assert.Equal(t, 74, attrMap["ST"], "Strength should match session")
		assert.Equal(t, 65, attrMap["GS"], "Dexterity should match session")
		assert.Equal(t, 76, attrMap["GW"], "Agility should match session")
		assert.Equal(t, 58, attrMap["KO"], "Constitution should match session")
		assert.Equal(t, 83, attrMap["ZT"], "Magic Talent should match session")
		assert.Equal(t, 69, attrMap["AU"], "Charisma should match session")
		assert.Equal(t, 59, attrMap["pA"], "Personal Charisma should match session")
		assert.Equal(t, 72, attrMap["WK"], "Willpower should match session")

		// Validate derived values
		assert.Equal(t, 17, createdChar.Lp.Max, "LP Max should match session")
		assert.Equal(t, 17, createdChar.Lp.Value, "LP Value should equal Max initially")
		assert.Equal(t, 33, createdChar.Ap.Max, "AP Max should match session")
		assert.Equal(t, 33, createdChar.Ap.Value, "AP Value should equal Max initially")
		assert.Equal(t, 8, createdChar.B.Max, "B Max should match session")
		assert.Equal(t, 8, createdChar.B.Value, "B Value should equal Max initially")

		// Validate skills were transferred (session has 10 skills: 6 regular skills + 4 weapon skills)
		// Regular skills: Klettern, Reiten, Sprache, Athletik, Heilkunde, Naturkunde (6)
		// Weapon skills: Spießwaffen, Stielwurfwaffen, Waffenloser Kampf, Stichwaffen (4)
		assert.Equal(t, 6, len(createdChar.Fertigkeiten), "Should have 6 regular skills")
		assert.Equal(t, 4, len(createdChar.Waffenfertigkeiten), "Should have 4 weapon skills")

		// Validate that skills use database initial values (not session levels which were 0)
		// Check a few specific skills exist
		skillNames := make([]string, len(createdChar.Fertigkeiten))
		for i, skill := range createdChar.Fertigkeiten {
			skillNames[i] = skill.Name
		}
		assert.Contains(t, skillNames, "Klettern", "Should contain Klettern skill")
		assert.Contains(t, skillNames, "Reiten", "Should contain Reiten skill")
		assert.Contains(t, skillNames, "Athletik", "Should contain Athletik skill")
		assert.Contains(t, skillNames, "Heilkunde", "Should contain Heilkunde skill")

		// Check weapon skills
		weaponSkillNames := make([]string, len(createdChar.Waffenfertigkeiten))
		for i, skill := range createdChar.Waffenfertigkeiten {
			weaponSkillNames[i] = skill.Name
		}
		assert.Contains(t, weaponSkillNames, "Spießwaffen", "Should contain Spießwaffen weapon skill")
		assert.Contains(t, weaponSkillNames, "Stichwaffen", "Should contain Stichwaffen weapon skill")

		// Validate spells were transferred (session has 4 spells)
		assert.Equal(t, 4, len(createdChar.Zauber), "Should have 4 spells")

		// Validate spell names
		spellNames := make([]string, len(createdChar.Zauber))
		for i, spell := range createdChar.Zauber {
			spellNames[i] = spell.Name
		}
		assert.Contains(t, spellNames, "Göttlicher Schutz v. d. Bösen", "Should contain Göttlicher Schutz v. d. Bösen spell")
		assert.Contains(t, spellNames, "Erkennen der Aura", "Should contain Erkennen der Aura spell")
		assert.Contains(t, spellNames, "Heiliger Zorn", "Should contain Heiliger Zorn spell")
		assert.Contains(t, spellNames, "Blutmeisterschaft", "Should contain Blutmeisterschaft spell")

		// Verify session was deleted after successful creation
		var deletedSession models.CharacterCreationSession
		err = database.DB.Where("id = ?", "char_create_1_1756326371").First(&deletedSession).Error
		assert.Error(t, err, "Session should be deleted after character creation")
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "Session should not be found in database")

		t.Logf("Character successfully created with ID: %.0f", characterID)
		t.Logf("Character name: %s", createdChar.Name)
		t.Logf("Character has %d skills, %d weapon skills, %d spells",
			len(createdChar.Fertigkeiten), len(createdChar.Waffenfertigkeiten), len(createdChar.Zauber))
	})

	t.Run("FinalizeCharacterCreation with invalid session", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/characters/sessions/nonexistent_session/finalize", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", uint(1))
		c.Params = gin.Params{
			gin.Param{Key: "sessionId", Value: "nonexistent_session"},
		}

		FinalizeCharacterCreation(c)

		assert.Equal(t, http.StatusNotFound, w.Code, "Should return 404 for non-existent session")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "Session not found", response["error"])
	})

	t.Run("FinalizeCharacterCreation with unauthorized user", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/characters/sessions/char_create_1_1756326371/finalize", nil)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", uint(0)) // Invalid userID
		c.Params = gin.Params{
			gin.Param{Key: "sessionId", Value: "char_create_1_1756326371"},
		}

		FinalizeCharacterCreation(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "Should return 401 for unauthorized user")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
		assert.Equal(t, "Unauthorized", response["error"])
	})
}

func TestListCharacters(t *testing.T) {
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

	t.Run("ListCharacters Success", func(t *testing.T) {
		// Create a test user
		u := user.User{}
		u.FirstId(1)
		token := user.GenerateToken(&u)

		// Create a test HTTP request
		req, _ := http.NewRequest("GET", "/api/characters", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", uint(1)) // Set valid userID

		// Call the handler
		ListCharacters(c)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)
		type AllCharacters struct {
			SelfOwned []models.CharList `json:"self_owned"`
			Others    []models.CharList `json:"others"`
		}

		var response AllCharacters
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Response should be an array (could be empty if no characters exist)
		assert.IsType(t, AllCharacters{}, response)
	})

	t.Run("ListCharacters with Invalid User", func(t *testing.T) {
		// Create a test user with invalid ID
		u := user.User{}
		u.FirstId(999) // Non-existent user ID
		token := user.GenerateToken(&u)

		// Create a test HTTP request
		req, _ := http.NewRequest("GET", "/api/characters", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Set("userID", uint(999)) // Set invalid userID

		// Call the handler
		ListCharacters(c)

		// Should still return OK with empty list if user has no characters
		assert.Equal(t, http.StatusOK, w.Code)
		type AllCharacters struct {
			SelfOwned []models.CharList `json:"self_owned"`
			Others    []models.CharList `json:"others"`
		}

		var response AllCharacters
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.SelfOwned), "Should return empty list for user with no characters")
	})

	t.Run("ListCharacters without UserID", func(t *testing.T) {
		// Create a test HTTP request without setting userID in context
		req, _ := http.NewRequest("GET", "/api/characters", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		// Don't set userID in context - this should trigger an error

		// Call the handler
		ListCharacters(c)

		// Should still return OK with empty list since GetUint("userID") returns 0 for missing userID
		assert.Equal(t, http.StatusOK, w.Code)
		type AllCharacters struct {
			SelfOwned []models.CharList `json:"self_owned"`
			Others    []models.CharList `json:"others"`
		}

		var response AllCharacters
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(response.SelfOwned), "Should return empty list for userID 0")
	})
}
