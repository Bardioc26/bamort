package character

import (
	"bamort/database"
	"bamort/equipment"
	"bamort/models"
	"bamort/skills"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPracticePointsAPI(t *testing.T) {
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
	err = models.MigrateStructure()
	assert.NoError(t, err)

	// Create test skill data
	err = createTestSkillData()
	assert.NoError(t, err)
	defer cleanupTestSkillData()

	// Create a test character
	character := &Char{
		BamortBase: models.BamortBase{
			Name: "Test Character",
		},
		Rasse: "Human",
		Typ:   "Hx",
	}
	err = character.Create()
	assert.NoError(t, err)

	// Add a test skill to the character
	testSkill := &skills.Fertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Menschenkenntnis",
			},
			CharacterID: character.ID,
		},
		Pp: 0,
	}

	// Save the skill to database
	result := database.GetDB().Create(testSkill)
	assert.NoError(t, result.Error)

	// Add a test spell to the character
	testSpell := &skills.Fertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Macht über das Selbst",
			},
			CharacterID: character.ID,
		},
		Pp: 0,
	}

	// Save the spell to database
	result = database.GetDB().Create(testSpell)
	assert.NoError(t, result.Error)

	// Add the "Beherrschen" magic school as a skill since spell PP go to the magic school
	beherrschenSkill := &skills.Fertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Beherrschen",
			},
			CharacterID: character.ID,
		},
		Pp: 0,
	}

	// Save the magic school skill to database
	result = database.GetDB().Create(beherrschenSkill)
	assert.NoError(t, result.Error)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Register routes
	api := router.Group("/api")
	RegisterRoutes(api)

	t.Run("GetPracticePoints", func(t *testing.T) {
		// Test getting practice points for a character with no PP
		req := httptest.NewRequest(http.MethodGet, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var pp []PracticePointResponse
		err := json.Unmarshal(w.Body.Bytes(), &pp)
		assert.NoError(t, err)
		assert.Empty(t, pp) // Should be empty initially
	})

	t.Run("AddPracticePoint", func(t *testing.T) {
		// Add practice points to a specific skill
		request := map[string]interface{}{
			"skill_name": "Menschenkenntnis",
			"amount":     3,
		}
		jsonData, _ := json.Marshal(request)

		req := httptest.NewRequest(http.MethodPost, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points/add", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Debug: print response body if test fails
		if w.Code != http.StatusOK {
			t.Logf("Response body: %s", w.Body.String())
		}

		var pp []PracticePointResponse
		err := json.Unmarshal(w.Body.Bytes(), &pp)
		assert.NoError(t, err)
		assert.Len(t, pp, 1)
		assert.Equal(t, "Menschenkenntnis", pp[0].SkillName)
		assert.Equal(t, 3, pp[0].Amount)
	})

	t.Run("UsePracticePoint", func(t *testing.T) {
		// Reset skill to 0 PP first to ensure clean test
		testSkill.Pp = 0
		database.GetDB().Save(testSkill)

		// First add practice points
		addRequest := map[string]interface{}{
			"skill_name": "Menschenkenntnis",
			"amount":     3,
		}
		jsonData, _ := json.Marshal(addRequest)

		req := httptest.NewRequest(http.MethodPost, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points/add", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify that practice points were actually added by checking the response
		var addedPP []PracticePointResponse
		addErr := json.Unmarshal(w.Body.Bytes(), &addedPP)
		assert.NoError(t, addErr)
		assert.Len(t, addedPP, 1)
		assert.Equal(t, "Menschenkenntnis", addedPP[0].SkillName)
		assert.Equal(t, 3, addedPP[0].Amount)

		// Then use one practice point
		useRequest := map[string]interface{}{
			"skill_name": "Menschenkenntnis",
			"amount":     1,
		}
		jsonData, _ = json.Marshal(useRequest)

		req = httptest.NewRequest(http.MethodPost, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points/use", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Debug: print response body if test fails
		if w.Code != http.StatusOK {
			t.Logf("Response body: %s", w.Body.String())
		}

		var pp []PracticePointResponse
		err := json.Unmarshal(w.Body.Bytes(), &pp)
		assert.NoError(t, err)
		assert.Len(t, pp, 1)
		assert.Equal(t, "Menschenkenntnis", pp[0].SkillName)
		assert.Equal(t, 2, pp[0].Amount) // Should be reduced by 1
	})

	t.Run("SkillCostWithPP", func(t *testing.T) {
		// First check current PP status
		req := httptest.NewRequest(http.MethodGet, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var currentPP []PracticePointResponse
		json.Unmarshal(w.Body.Bytes(), &currentPP)
		t.Logf("Current PP before skill cost test: %+v", currentPP)

		// Find Menschenkenntnis PP
		var humanKnowledgePP int
		for _, pp := range currentPP {
			if pp.SkillName == "Menschenkenntnis" {
				humanKnowledgePP = pp.Amount
				break
			}
		}

		// Test skill cost calculation with practice points
		request := map[string]interface{}{
			"name":          "Menschenkenntnis",
			"type":          "skill",
			"action":        "improve",
			"current_level": 10,
			"use_pp":        1,
		}
		jsonData, _ := json.Marshal(request)

		req = httptest.NewRequest(http.MethodPost, "/api/characters/"+strconv.Itoa(int(character.ID))+"/skill-cost", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response SkillCostResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify PP information is included
		assert.Equal(t, 1, response.PPUsed)
		assert.Equal(t, humanKnowledgePP, response.PPAvailable)      // Should match current available PP
		assert.Greater(t, response.OriginalCost, response.FinalCost) // Final cost should be lower
	})

	t.Run("SpellCostWithPP", func(t *testing.T) {
		// Add PP for spell - should go to the "Beherrschen" magic school, not the specific spell
		request := map[string]interface{}{
			"skill_name": "Macht über das Selbst",
			"amount":     2,
		}
		jsonData, _ := json.Marshal(request)

		req := httptest.NewRequest(http.MethodPost, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points/add", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify that PP were added to "Beherrschen" skill, not "Macht über das Selbst"
		var ppResponse []PracticePointResponse
		err := json.Unmarshal(w.Body.Bytes(), &ppResponse)
		assert.NoError(t, err)

		// Should have PP on "Beherrschen", not on the specific spell
		found := false
		for _, pp := range ppResponse {
			if pp.SkillName == "Beherrschen" && pp.Amount == 2 {
				found = true
				break
			}
		}
		assert.True(t, found, "Practice points should be added to 'Beherrschen' magic school, not the specific spell")

		// Test spell learning with practice points
		spellRequest := map[string]interface{}{
			"name":   "Macht über das Selbst",
			"type":   "spell",
			"action": "learn",
			"use_pp": 1,
		}
		jsonData, _ = json.Marshal(spellRequest)

		req = httptest.NewRequest(http.MethodPost, "/api/characters/"+strconv.Itoa(int(character.ID))+"/skill-cost", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response SkillCostResponse
		unmarshalErr := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, unmarshalErr)

		// Verify PP information is included for spells
		assert.Equal(t, 1, response.PPUsed)
		assert.Equal(t, 2, response.PPAvailable)
	})
}
