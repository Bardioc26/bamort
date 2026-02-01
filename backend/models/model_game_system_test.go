package models

import (
	"testing"

	"bamort/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGameSystem_Methods(t *testing.T) {
	// Initialize test DB and migrations
	database.SetupTestDB(true)
	defer database.ResetTestDB()

	// Ensure default exists (use FirstOrCreate to avoid unique constraint)
	defaultGS := GameSystem{}
	database.DB.Where(GameSystem{Code: "M5"}).FirstOrCreate(&defaultGS, GameSystem{Code: "M5", Name: "M5", IsActive: true})

	t.Run("FirstByCode returns matching record", func(t *testing.T) {
		gs := GameSystem{Code: "TESTCODE", Name: "Test System", IsActive: true}
		err := database.DB.Create(&gs).Error
		require.NoError(t, err)

		var found GameSystem
		err = found.FirstByCode("TESTCODE")
		assert.NoError(t, err)
		assert.Equal(t, "TESTCODE", found.Code)
		assert.Equal(t, "Test System", found.Name)
	})

	t.Run("GetDefault returns M5", func(t *testing.T) {
		var found GameSystem
		err := found.GetDefault()
		assert.NoError(t, err)
		assert.Equal(t, "M5", found.Code)
	})

	t.Run("FirstByName behaviour", func(t *testing.T) {
		// existing name
		gs := GameSystem{Code: "CUSTOM", Name: "CustomSys", IsActive: true}
		err := database.DB.Create(&gs).Error
		require.NoError(t, err)

		var byName GameSystem
		err = byName.FirstByName("CustomSys")
		assert.NoError(t, err)
		assert.Equal(t, "CustomSys", byName.Name)

		// empty name should fallback to default
		byName = GameSystem{}
		err = byName.FirstByName("")
		assert.NoError(t, err)
		assert.Equal(t, "M5", byName.Code)

		// empty name should fallback to default
		byName = GameSystem{}
		err = byName.FirstByName("midgard")
		assert.NoError(t, err)
		assert.Equal(t, "M5", byName.Code)

		// non-existent name should return an error
		var notFound GameSystem
		err = notFound.FirstByName("NoSuchSystem")
		assert.Error(t, err)
	})

	t.Run("Get By Id", func(t *testing.T) {
		var found GameSystem
		err := found.FirstByID(1)
		assert.NoError(t, err)
		assert.Equal(t, "M5", found.Code)
	})
}
