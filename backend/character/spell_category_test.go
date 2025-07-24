package character

import (
	"bamort/database"
	"bamort/equipment"
	"bamort/gsmaster"
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

func TestSpellCategoryMapping(t *testing.T) {
	// Setup test database
	database.SetupTestDB()
	defer database.ResetTestDB()

	// Migrate the schema
	err := MigrateStructure()
	assert.NoError(t, err)
	err = skills.MigrateStructure()
	assert.NoError(t, err)
	err = equipment.MigrateStructure()
	assert.NoError(t, err)
	err = gsmaster.MigrateStructure()
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

	// Add the "Beherrschen" magic school as a skill
	beherrschenSkill := &skills.Fertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase: models.BamortBase{
				Name: "Beherrschen",
			},
			CharacterID: character.ID,
		},
		Pp: 0,
	}
	result := database.GetDB().Create(beherrschenSkill)
	assert.NoError(t, result.Error)

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	RegisterRoutes(api)

	t.Run("Spell PP goes to magic school", func(t *testing.T) {
		// Add PP for spell "Macht über das Selbst" - should go to "Beherrschen"
		request := map[string]interface{}{
			"skill_name": "Macht über das Selbst",
			"amount":     3,
		}
		jsonData, _ := json.Marshal(request)

		req := httptest.NewRequest(http.MethodPost, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points/add", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response contains PP for "Beherrschen", not "Macht über das Selbst"
		var ppResponse []PracticePointResponse
		err := json.Unmarshal(w.Body.Bytes(), &ppResponse)
		assert.NoError(t, err)

		// Should have exactly one entry for "Beherrschen"
		assert.Len(t, ppResponse, 1)
		assert.Equal(t, "Beherrschen", ppResponse[0].SkillName)
		assert.Equal(t, 3, ppResponse[0].Amount)
	})

	t.Run("Use PP from magic school for spell", func(t *testing.T) {
		// Use PP for spell "Macht über das Selbst" - should use from "Beherrschen"
		request := map[string]interface{}{
			"skill_name": "Macht über das Selbst",
			"amount":     1,
		}
		jsonData, _ := json.Marshal(request)

		req := httptest.NewRequest(http.MethodPost, "/api/characters/"+strconv.Itoa(int(character.ID))+"/practice-points/use", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response contains reduced PP for "Beherrschen"
		var ppResponse []PracticePointResponse
		err := json.Unmarshal(w.Body.Bytes(), &ppResponse)
		assert.NoError(t, err)

		// Should have exactly one entry for "Beherrschen" with 2 PP remaining
		assert.Len(t, ppResponse, 1)
		assert.Equal(t, "Beherrschen", ppResponse[0].SkillName)
		assert.Equal(t, 2, ppResponse[0].Amount)
	})
}
