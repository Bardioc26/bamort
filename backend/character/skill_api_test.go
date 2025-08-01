package character

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"bamort/database"
	"bamort/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAvailableSkills(t *testing.T) {
	// Setup test database with real data
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("Get available skills for existing character - default reward type", func(t *testing.T) {
		// Get a character ID from the test data
		var testChar models.Char
		err := database.DB.Preload("Fertigkeiten").Preload("Erfahrungsschatz").Preload("Vermoegen").First(&testChar).Error
		assert.NoError(t, err, "Should find a test character")

		// Gib dem Charakter genug EP und Gold für Tests
		testChar.Erfahrungsschatz.EP = 1000
		testChar.Vermoegen.Goldstücke = 1000
		database.DB.Save(&testChar.Erfahrungsschatz)
		database.DB.Save(&testChar.Vermoegen)

		characterID := fmt.Sprintf("%d", testChar.ID)

		// Create HTTP request with default reward type
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/characters/%s/available-skills?reward_type=default", characterID), nil)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: characterID}}

		fmt.Printf("Test: Get available skills for character ID %s (default reward)\n", characterID)
		fmt.Printf("Character EP: %d, Gold: %d\n", testChar.Erfahrungsschatz.EP, testChar.Vermoegen.Goldstücke)

		// Call the handler function
		GetAvailableSkillsOld(c)

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
		assert.Contains(t, response, "skills_by_category", "Response should contain skills_by_category field")

		// Check if skills_by_category is an object
		skillsByCategory, ok := response["skills_by_category"].(map[string]interface{})
		assert.True(t, ok, "skills_by_category should be an object")

		fmt.Printf("Found %d skill categories\n", len(skillsByCategory))

		// Check each category
		totalSkills := 0
		for categoryName, categorySkills := range skillsByCategory {
			fmt.Printf("Category: %s\n", categoryName)

			skills, ok := categorySkills.([]interface{})
			assert.True(t, ok, fmt.Sprintf("Category %s should contain an array of skills", categoryName))

			totalSkills += len(skills)
			fmt.Printf("  Skills in category: %d\n", len(skills))

			// Check structure of first skill in category if exists
			if len(skills) > 0 {
				firstSkill, ok := skills[0].(map[string]interface{})
				assert.True(t, ok, "First skill should be an object")

				// Check required fields
				assert.Contains(t, firstSkill, "name", "Skill should have name field")
				assert.Contains(t, firstSkill, "epCost", "Skill should have epCost field")
				assert.Contains(t, firstSkill, "goldCost", "Skill should have goldCost field")

				skillName, ok := firstSkill["name"].(string)
				assert.True(t, ok, "Skill name should be a string")
				assert.NotEmpty(t, skillName, "Skill name should not be empty")

				epCost, ok := firstSkill["epCost"].(float64)
				assert.True(t, ok, "EP cost should be a number")
				assert.GreaterOrEqual(t, epCost, float64(0), "EP cost should be non-negative")

				goldCost, ok := firstSkill["goldCost"].(float64)
				assert.True(t, ok, "Gold cost should be a number")
				assert.GreaterOrEqual(t, goldCost, float64(0), "Gold cost should be non-negative")

				fmt.Printf("  First skill: %s (EP: %.0f, Gold: %.0f)\n",
					skillName, epCost, goldCost)
			}
		}

		fmt.Printf("Total available skills: %d\n", totalSkills)
	})

	t.Run("Get available skills for existing character - noGold reward type", func(t *testing.T) {
		// Get a character ID from the test data
		var testChar models.Char
		err := database.DB.Preload("Fertigkeiten").Preload("Erfahrungsschatz").Preload("Vermoegen").First(&testChar).Error
		assert.NoError(t, err, "Should find a test character")

		// Gib dem Charakter genug EP und Gold für Tests
		testChar.Erfahrungsschatz.EP = 1000
		testChar.Vermoegen.Goldstücke = 1000
		database.DB.Save(&testChar.Erfahrungsschatz)
		database.DB.Save(&testChar.Vermoegen)

		characterID := fmt.Sprintf("%d", testChar.ID)

		// Create HTTP request with noGold reward type
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/characters/%s/available-skills?reward_type=noGold", characterID), nil)
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: characterID}}

		fmt.Printf("Test: Get available skills for character ID %s (noGold reward)\n", characterID)

		// Call the handler function
		GetAvailableSkillsOld(c)

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
		assert.Contains(t, response, "skills_by_category", "Response should contain skills_by_category field")

		skillsByCategory, ok := response["skills_by_category"].(map[string]interface{})
		assert.True(t, ok, "skills_by_category should be an object")

		// Check that skills have goldCost = 0 for noGold reward type
		for categoryName, categorySkills := range skillsByCategory {
			skills, ok := categorySkills.([]interface{})
			assert.True(t, ok, fmt.Sprintf("Category %s should contain an array of skills", categoryName))

			if len(skills) > 0 {
				firstSkill, ok := skills[0].(map[string]interface{})
				assert.True(t, ok, "First skill should be an object")

				goldCost, ok := firstSkill["goldCost"].(float64)
				assert.True(t, ok, "Gold cost should be a number")
				assert.Equal(t, float64(0), goldCost, "Gold cost should be 0 for noGold reward type")

				fmt.Printf("Category %s - First skill gold cost: %.0f (should be 0)\n", categoryName, goldCost)
			}
		}
	})

	t.Run("Error case - character not found", func(t *testing.T) {
		// Test with non-existent character ID
		req, _ := http.NewRequest("GET", "/api/characters/99999/available-skills", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: "99999"}}

		GetAvailableSkillsOld(c)

		// Should return 404
		assert.Equal(t, http.StatusNotFound, w.Code, "Status code should be 404 Not Found")

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")
		assert.Contains(t, response, "error", "Response should contain error message")
	})

	t.Run("Check that learned skills are excluded", func(t *testing.T) {
		// Get a character with some skills
		var testChar models.Char
		err := database.DB.Preload("Fertigkeiten").First(&testChar).Error
		assert.NoError(t, err, "Should find a test character")

		characterID := fmt.Sprintf("%d", testChar.ID)

		// Get list of learned skills
		learnedSkillNames := make(map[string]bool)
		for _, skill := range testChar.Fertigkeiten {
			learnedSkillNames[skill.Name] = true
		}

		// Make request
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/characters/%s/available-skills", characterID), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: characterID}}

		GetAvailableSkillsOld(c)

		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200 OK")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "Response should be valid JSON")

		skillsByCategory, ok := response["skills_by_category"].(map[string]interface{})
		assert.True(t, ok, "skills_by_category should be an object")

		// Check that no learned skills appear in available skills
		for categoryName, categorySkills := range skillsByCategory {
			skills, ok := categorySkills.([]interface{})
			assert.True(t, ok, fmt.Sprintf("Category %s should contain an array of skills", categoryName))

			for _, skillInterface := range skills {
				skill, ok := skillInterface.(map[string]interface{})
				assert.True(t, ok, "Skill should be an object")

				skillName, ok := skill["name"].(string)
				assert.True(t, ok, "Skill name should be a string")

				// Assert that this skill is not in the learned skills list
				assert.False(t, learnedSkillNames[skillName],
					fmt.Sprintf("Learned skill '%s' should not appear in available skills", skillName))
			}
		}

		fmt.Printf("Verified that %d learned skills are excluded from available skills\n", len(testChar.Fertigkeiten))
	})
}
