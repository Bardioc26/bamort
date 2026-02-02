package character

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bamort/database"
	"bamort/models"
	"bamort/testutils"
	"bamort/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestOwnershipGuardsBlockMutations(t *testing.T) {
	testutils.SetupTestEnvironment(t)
	gin.SetMode(gin.TestMode)

	database.SetupTestDB(true, true)
	t.Cleanup(database.ResetTestDB)

	require.NoError(t, models.MigrateStructure())

	t.Run("UpdateCharacter blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 101)
		char := createCharacterOwnedBy(t, owner.UserID)
		originalName := char.Name

		ctx, w := buildJSONContext(t, http.MethodPut, map[string]any{"name": "New Name"}, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		UpdateCharacter(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		reloaded := reloadCharacter(t, char.ID)
		require.Equal(t, originalName, reloaded.Name)
	})

	t.Run("DeleteCharacter blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 102)
		char := createCharacterOwnedBy(t, owner.UserID)

		ctx, w := buildJSONContext(t, http.MethodDelete, nil, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		DeleteCharacter(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		var exists models.Char
		err := database.DB.First(&exists, char.ID).Error
		require.NoError(t, err)
	})

	t.Run("UpdateCharacterExperience blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 103)
		char := createCharacterOwnedBy(t, owner.UserID)
		seedExperience(t, char, 25)

		ctx, w := buildJSONContext(t, http.MethodPost, map[string]any{"experience_points": 50}, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		UpdateCharacterExperience(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		reloaded := reloadCharacterWithPreloads(t, char.ID)
		require.Equal(t, 25, reloaded.Erfahrungsschatz.EP)
	})

	t.Run("UpdateCharacterWealth blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 104)
		char := createCharacterOwnedBy(t, owner.UserID)
		seedWealth(t, char, 5, 4, 3)

		ctx, w := buildJSONContext(t, http.MethodPost, map[string]any{"goldstücke": 50}, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		UpdateCharacterWealth(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		reloaded := reloadCharacterWithPreloads(t, char.ID)
		require.Equal(t, 5, reloaded.Vermoegen.Goldstuecke)
		require.Equal(t, 4, reloaded.Vermoegen.Silberstuecke)
		require.Equal(t, 3, reloaded.Vermoegen.Kupferstuecke)
	})

	t.Run("LearnSkill blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 105)
		char := createCharacterOwnedBy(t, owner.UserID)
		before := countSkills(t, char.ID)

		ctx, w := buildJSONContext(t, http.MethodPost, map[string]any{"name": "Test Skill", "target_level": 1, "type": "skill"}, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		LearnSkill(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		after := countSkills(t, char.ID)
		require.Equal(t, before, after)
	})

	t.Run("ImproveSkill blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 106)
		char := createCharacterOwnedBy(t, owner.UserID)
		seedSkill(t, char, "Athletik", 1, 0)
		before := countSkills(t, char.ID)

		payload := map[string]any{
			"char_id":       char.ID,
			"name":          "Athletik",
			"current_level": 1,
			"target_level":  2,
			"type":          "skill",
			"action":        "improve",
			"reward":        "default",
		}
		ctx, w := buildJSONContext(t, http.MethodPost, payload, char.UserID+1, nil)

		ImproveSkill(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		after := countSkills(t, char.ID)
		require.Equal(t, before, after)
	})

	t.Run("LearnSpell blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 107)
		char := createCharacterOwnedBy(t, owner.UserID)
		before := countSpells(t, char.ID)

		ctx, w := buildJSONContext(t, http.MethodPost, map[string]any{"name": "Test Spell"}, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		LearnSpell(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		after := countSpells(t, char.ID)
		require.Equal(t, before, after)
	})

	t.Run("UpdatePracticePoints blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 108)
		char := createCharacterOwnedBy(t, owner.UserID)
		seedSkill(t, char, "Menschenkenntnis", 1, 2)
		before := fetchSkillPp(t, char.ID, "Menschenkenntnis")

		payload := []map[string]any{{"skill_name": "Menschenkenntnis", "amount": 0}}
		ctx, w := buildJSONContext(t, http.MethodPost, payload, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		UpdatePracticePoints(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		after := fetchSkillPp(t, char.ID, "Menschenkenntnis")
		require.Equal(t, before, after)
	})

	t.Run("AddPracticePoint blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 109)
		char := createCharacterOwnedBy(t, owner.UserID)
		seedSkill(t, char, "Athletik", 1, 1)
		before := fetchSkillPp(t, char.ID, "Athletik")

		payload := map[string]any{"skill_name": "Athletik", "amount": 3}
		ctx, w := buildJSONContext(t, http.MethodPost, payload, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		AddPracticePoint(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		after := fetchSkillPp(t, char.ID, "Athletik")
		require.Equal(t, before, after)
	})

	t.Run("UsePracticePoint blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 110)
		char := createCharacterOwnedBy(t, owner.UserID)
		seedSkill(t, char, "Athletik", 1, 2)
		before := fetchSkillPp(t, char.ID, "Athletik")

		payload := map[string]any{"skill_name": "Athletik", "amount": 1}
		ctx, w := buildJSONContext(t, http.MethodPost, payload, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		UsePracticePoint(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		after := fetchSkillPp(t, char.ID, "Athletik")
		require.Equal(t, before, after)
	})

	t.Run("UpdateCharacterShares blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 111)
		shareTarget := ensureUserExists(t, 112)
		char := createCharacterOwnedBy(t, owner.UserID)

		payload := map[string]any{"user_ids": []uint{shareTarget.UserID}}
		ctx, w := buildJSONContext(t, http.MethodPost, payload, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		UpdateCharacterShares(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		var shares []models.CharShare
		err := database.DB.Where("character_id = ?", char.ID).Find(&shares).Error
		require.NoError(t, err)
		require.Empty(t, shares)
	})

	t.Run("UpdateCharacterImage blocks non-owner", func(t *testing.T) {
		owner := ensureUserExists(t, 113)
		char := createCharacterOwnedBy(t, owner.UserID)
		char.Image = "initial.png"
		require.NoError(t, database.DB.Save(&char).Error)

		ctx, w := buildJSONContext(t, http.MethodPost, map[string]any{"image": "new.png"}, char.UserID+1, map[string]string{"id": fmt.Sprint(char.ID)})

		UpdateCharacterImage(ctx)

		require.Equal(t, http.StatusForbidden, w.Code)
		reloaded := reloadCharacter(t, char.ID)
		require.Equal(t, "initial.png", reloaded.Image)
	})
}

func ensureUserExists(t *testing.T, id uint) user.User {
	var existing user.User
	err := database.DB.First(&existing, "user_id = ?", id).Error
	if err == nil {
		return existing
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		require.NoError(t, err)
	}

	newUser := user.User{
		UserID:   id,
		Username: fmt.Sprintf("user-%d-%d", id, time.Now().UnixNano()),
		Email:    fmt.Sprintf("user-%d@example.com", id),
		Role:     user.RoleStandardUser,
	}
	require.NoError(t, database.DB.Create(&newUser).Error)
	return newUser
}

func createCharacterOwnedBy(t *testing.T, ownerID uint) models.Char {
	char := models.Char{
		BamortBase: models.BamortBase{Name: fmt.Sprintf("Char-%d", time.Now().UnixNano())},
		UserID:     ownerID,
		Typ:        "Krieger",
		Rasse:      "Mensch",
		Grad:       1,
	}
	require.NoError(t, database.DB.Create(&char).Error)
	return char
}

func seedExperience(t *testing.T, char models.Char, ep int) {
	exp := models.Erfahrungsschatz{
		BamortCharTrait: models.BamortCharTrait{CharacterID: char.ID, UserID: char.UserID},
		EP:              ep,
	}
	require.NoError(t, database.DB.Create(&exp).Error)
}

func seedWealth(t *testing.T, char models.Char, gold, silver, copper int) {
	wealth := models.Vermoegen{
		BamortCharTrait: models.BamortCharTrait{CharacterID: char.ID, UserID: char.UserID},
		Goldstuecke:     gold,
		Silberstuecke:   silver,
		Kupferstuecke:   copper,
	}
	require.NoError(t, database.DB.Create(&wealth).Error)
}

func seedSkill(t *testing.T, char models.Char, name string, level, pp int) {
	skill := models.SkFertigkeit{
		BamortCharTrait: models.BamortCharTrait{
			BamortBase:  models.BamortBase{Name: name},
			CharacterID: char.ID,
			UserID:      char.UserID,
		},
		Fertigkeitswert: level,
		Pp:              pp,
		Improvable:      true,
		Category:        "Test",
	}
	require.NoError(t, database.DB.Create(&skill).Error)
}

func countSkills(t *testing.T, charID uint) int {
	var skills []models.SkFertigkeit
	require.NoError(t, database.DB.Where("character_id = ?", charID).Find(&skills).Error)
	return len(skills)
}

func countSpells(t *testing.T, charID uint) int {
	var spells []models.SkZauber
	require.NoError(t, database.DB.Where("character_id = ?", charID).Find(&spells).Error)
	return len(spells)
}

func fetchSkillPp(t *testing.T, charID uint, name string) int {
	var skill models.SkFertigkeit
	err := database.DB.Where("character_id = ? AND name = ?", charID, name).First(&skill).Error
	require.NoError(t, err)
	return skill.Pp
}

func reloadCharacter(t *testing.T, id uint) models.Char {
	var char models.Char
	require.NoError(t, database.DB.First(&char, id).Error)
	return char
}

func reloadCharacterWithPreloads(t *testing.T, id uint) models.Char {
	var char models.Char
	require.NoError(t, database.DB.Preload("Erfahrungsschatz").Preload("Vermoegen").First(&char, id).Error)
	return char
}

func buildJSONContext(t *testing.T, method string, body any, userID uint, params map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}

	req, err := http.NewRequest(method, "/", &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Set("userID", userID)

	for k, v := range params {
		ctx.Params = append(ctx.Params, gin.Param{Key: k, Value: v})
	}

	return ctx, w
}
