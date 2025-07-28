package character

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bamort/database"
	"bamort/models"
	"bamort/skills"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestImproveSkillUpdatesLevel(t *testing.T) {
	// Setup test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Migrate the schema
	err := MigrateStructure()
	assert.NoError(t, err)

	// Also migrate skills and equipment to avoid preload errors
	err = skills.MigrateStructure()
	assert.NoError(t, err)
	err = models.MigrateStructure()
	assert.NoError(t, err)

	// Try to migrate equipment if it exists
	if equipmentDB := database.DB.Exec("CREATE TABLE IF NOT EXISTS equi_equipments (id INTEGER PRIMARY KEY, character_id INTEGER, name TEXT)"); equipmentDB.Error != nil {
		t.Logf("Warning: Could not create equipment table: %v", equipmentDB.Error)
	}

	// Create audit log table for tests
	database.DB.Exec(`CREATE TABLE IF NOT EXISTS audit_log_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		character_id INTEGER,
		field_name TEXT,
		old_value INTEGER,
		new_value INTEGER,
		difference INTEGER,
		reason TEXT,
		user_id INTEGER,
		timestamp DATETIME,
		notes TEXT
	)`)

	// Create container table to avoid preload errors
	database.DB.Exec("CREATE TABLE IF NOT EXISTS equi_containers (id INTEGER PRIMARY KEY, character_id INTEGER)")
	database.DB.Exec("CREATE TABLE IF NOT EXISTS equi_weapons (id INTEGER PRIMARY KEY, character_id INTEGER)")
	database.DB.Exec("CREATE TABLE IF NOT EXISTS equi_transportmittel (id INTEGER PRIMARY KEY, character_id INTEGER)")
	database.DB.Exec("CREATE TABLE IF NOT EXISTS equi_equipments (id INTEGER PRIMARY KEY, character_id INTEGER)")

	// Create test character with ID 20
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
			EP: 326, // Starting EP
		},
		Vermoegen: Vermoegen{
			BamortCharTrait: models.BamortCharTrait{
				CharacterID: 20,
			},
			Goldstücke: 390, // Starting Gold
		},
	}

	// Add Athletik skill at level 9
	athletikSkill := skills.Fertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Athletik",
			},
			CharacterID: 20,
		},
		Fertigkeitswert: 9,
		Improvable:      true,
	}
	testChar.Fertigkeiten = append(testChar.Fertigkeiten, athletikSkill)

	err = testChar.Create()
	assert.NoError(t, err)

	// Verify character was created correctly
	var verifyChar Char
	err = database.DB.Preload("Fertigkeiten").Preload("Erfahrungsschatz").Preload("Vermoegen").First(&verifyChar, 20).Error
	assert.NoError(t, err)
	t.Logf("Character created with ID: %d, EP: %d, Skills: %d", verifyChar.ID, verifyChar.Erfahrungsschatz.EP, len(verifyChar.Fertigkeiten))

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	t.Run("ImproveSkill updates skill level from 9 to 10", func(t *testing.T) {
		// Create request body
		requestData := map[string]interface{}{
			"char_id":       20,
			"name":          "Athletik",
			"current_level": 9,
			"target_level":  10,
			"type":          "skill",
			"action":        "improve",
			"reward":        "default",
			"use_pp":        1,
			"use_gold":      0,
			"notes":         "Fertigkeit Athletik von 9 auf 10 verbessert (1 Level)",
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

		// Verify that the skill level was actually updated in the database
		var updatedChar Char
		err = database.DB.Preload("Fertigkeiten").First(&updatedChar, 20).Error
		assert.NoError(t, err)

		// Find the Athletik skill and check its level
		skillFound := false
		for _, skill := range updatedChar.Fertigkeiten {
			if skill.Name == "Athletik" {
				assert.Equal(t, 10, skill.Fertigkeitswert, "Athletik skill should be level 10 after improvement")
				skillFound = true
				t.Logf("Found Athletik skill with level: %d", skill.Fertigkeitswert)
				break
			}
		}
		assert.True(t, skillFound, "Athletik skill should be found in character's skills")

		t.Logf("Test completed successfully!")
		t.Logf("Athletik skill successfully updated from level 9 to level 10")
	})

	t.Run("LearnSkill creates new skill at level 1", func(t *testing.T) {
		// Create request body for learning a new skill
		requestData := map[string]interface{}{
			"name":          "Schwimmen",
			"current_level": 0,
			"target_level":  1,
			"type":          "skill",
			"action":        "learn",
			"reward":        "default",
			"notes":         "Neue Fertigkeit Schwimmen gelernt",
		}
		requestBody, _ := json.Marshal(requestData)

		// Create HTTP request
		req, _ := http.NewRequest("POST", "/api/characters/20/learn-skill", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		// Create response recorder
		w := httptest.NewRecorder()

		// Create Gin context with character ID parameter
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: "20"}}

		// Call the actual handler function
		LearnSkill(c)

		// Check if we got a successful response
		assert.Equal(t, http.StatusOK, w.Code, "Status code should be 200 OK")

		// Verify that the new skill was created in the database
		var updatedChar Char
		err = updatedChar.FirstID("20")
		assert.NoError(t, err)

		// Find the Schwimmen skill and check its level
		skillFound := false
		for _, skill := range updatedChar.Fertigkeiten {
			if skill.Name == "Schwimmen" {
				assert.Equal(t, 1, skill.Fertigkeitswert, "Schwimmen skill should be level 1 after learning")
				skillFound = true
				break
			}
		}
		assert.True(t, skillFound, "Schwimmen skill should be found in character's skills")

		t.Logf("New skill 'Schwimmen' successfully created at level 1")
	})
}
