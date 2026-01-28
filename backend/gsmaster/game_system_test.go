package gsmaster

import (
	"testing"

	"bamort/database"

	"github.com/stretchr/testify/assert"
)

func TestGetGameSystem(t *testing.T) {
	// Initialize test DB and migrations
	database.SetupTestDB(true)
	defer database.ResetTestDB()
	t.Run("GetGameSystem", func(t *testing.T) {
		gs := GetGameSystem(1, "")
		assert.Equal(t, "M5", gs.Code)
		gs = GetGameSystem(0, "midgard")
		assert.Equal(t, "M5", gs.Code)
		gs = GetGameSystem(1, "M5")
		assert.Equal(t, "M5", gs.Code)
	})
}
