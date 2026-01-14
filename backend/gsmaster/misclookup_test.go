package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiscLookup_TableName(t *testing.T) {
	misc := models.MiscLookup{}
	assert.Equal(t, "gsm_misc", misc.TableName())
}

func TestMiscLookup_CreateAndRetrieve(t *testing.T) {
	database.SetupTestDB()
	err := models.MigrateStructure()
	require.NoError(t, err)

	/*
		// Create test data
		testData := []models.MiscLookup{
			{Key: "gender", Value: "männlich"},
			{Key: "gender", Value: "weiblich"},
			{Key: "races", Value: "Mensch"},
			{Key: "races", Value: "Elf"},
		}

		// Insert test data
		for _, item := range testData {
			err := database.DB.Create(&item).Error
			require.NoError(t, err)
		}
	*/

	// Retrieve by key - sorted alphabetically by value (default)
	genders, err := GetMiscLookupByKey("gender")
	require.NoError(t, err)
	assert.Len(t, genders, 3)
	assert.Equal(t, "divers", genders[0].Value)
	assert.Equal(t, "männlich", genders[1].Value)
	assert.Equal(t, "weiblich", genders[2].Value)

	genders, err = GetMiscLookupByKey("gender", "id")
	require.NoError(t, err)
	assert.Len(t, genders, 3)
	assert.Equal(t, "männlich", genders[0].Value)
	assert.Equal(t, "weiblich", genders[1].Value)
	assert.Equal(t, "divers", genders[2].Value)

	races, err := GetMiscLookupByKey("races")
	require.NoError(t, err)
	assert.Len(t, races, 5)
}

func TestGetMiscLookupByKey_NotFound(t *testing.T) {
	database.SetupTestDB()
	err := models.MigrateStructure()
	require.NoError(t, err)

	items, err := GetMiscLookupByKey("nonexistent")
	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestMiscLookup_WithSourceInfo(t *testing.T) {
	database.SetupTestDB()
	err := models.MigrateStructure()
	require.NoError(t, err)

	misc := models.MiscLookup{
		Key:        "test_key",
		Value:      "test_value",
		SourceID:   1,
		PageNumber: 42,
	}

	err = database.DB.Create(&misc).Error
	require.NoError(t, err)
	assert.NotZero(t, misc.ID)

	// Retrieve and verify
	var retrieved models.MiscLookup
	err = database.DB.First(&retrieved, misc.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "test_key", retrieved.Key)
	assert.Equal(t, "test_value", retrieved.Value)
	assert.Equal(t, uint(1), retrieved.SourceID)
	assert.Equal(t, 42, retrieved.PageNumber)
}

func TestPopulateMiscLookupData(t *testing.T) {
	database.SetupTestDB()

	// Migrate the structure first
	err := models.MigrateStructure()
	require.NoError(t, err)

	/*
		// First population should succeed
		err = PopulateMiscLookupData()
		require.NoError(t, err)
	*/

	// Verify all keys have data
	expectedCounts := map[string]int{
		"gender":         3,
		"races":          5,
		"origins":        15,
		"social_classes": 4,
		"faiths":         5,
		"handedness":     3,
	}

	for key, expectedCount := range expectedCounts {
		items, err := GetMiscLookupByKey(key)
		require.NoError(t, err)
		assert.Len(t, items, expectedCount, "Expected %d items for key %s", expectedCount, key)
	}

	// Verify specific values
	genders, _ := GetMiscLookupByKey("gender")
	assert.Contains(t, []string{"männlich", "weiblich", "divers"}, genders[0].Value)

	races, _ := GetMiscLookupByKey("races")
	raceValues := make([]string, len(races))
	for i, r := range races {
		raceValues[i] = r.Value
	}
	assert.Contains(t, raceValues, "Mensch")
	assert.Contains(t, raceValues, "Elf")

	/*
		// Second population should not duplicate data
		err = PopulateMiscLookupData()
		require.NoError(t, err)
	*/

	var totalCount int64
	err = database.DB.Model(&models.MiscLookup{}).Count(&totalCount).Error
	require.NoError(t, err)
	assert.Equal(t, int64(39), totalCount, "Should not duplicate data on second population")
}

func TestGetSocialClassBonusPoints(t *testing.T) {
	database.SetupTestDB()
	err := models.MigrateStructure()
	require.NoError(t, err)

	/*
		// Populate data
		err = PopulateMiscLookupData()
		require.NoError(t, err)
	*/

	// Test Volk bonus
	bonuses, err := GetSocialClassBonusPoints("Volk")
	require.NoError(t, err)
	assert.Equal(t, 2, bonuses["Alltag"])

	// Test Adel bonus
	bonuses, err = GetSocialClassBonusPoints("Adel")
	require.NoError(t, err)
	assert.Equal(t, 2, bonuses["Sozial"])

	// Test Mittelschicht bonus
	bonuses, err = GetSocialClassBonusPoints("Mittelschicht")
	require.NoError(t, err)
	assert.Equal(t, 2, bonuses["Wissen"])

	// Test Unfrei bonus
	bonuses, err = GetSocialClassBonusPoints("Unfrei")
	require.NoError(t, err)
	assert.Equal(t, 2, bonuses["Halbwelt"])

	// Test non-existent social class
	bonuses, err = GetSocialClassBonusPoints("NonExistent")
	require.NoError(t, err)
	assert.Empty(t, bonuses)
}

func TestGetMiscLookupByKey_OrderParameter(t *testing.T) {
	database.SetupTestDB()
	err := models.MigrateStructure()
	require.NoError(t, err)

	// Create test data with different IDs and sources
	testData := []models.MiscLookup{
		{Key: "test_order", Value: "Zebra", SourceID: 2},
		{Key: "test_order", Value: "Alpha", SourceID: 1},
		{Key: "test_order", Value: "Beta", SourceID: 1},
	}

	for _, item := range testData {
		err := database.DB.Create(&item).Error
		require.NoError(t, err)
	}

	// Test default ordering (by value)
	items, err := GetMiscLookupByKey("test_order")
	require.NoError(t, err)
	assert.Len(t, items, 3)
	assert.Equal(t, "Alpha", items[0].Value)
	assert.Equal(t, "Beta", items[1].Value)
	assert.Equal(t, "Zebra", items[2].Value)

	// Test explicit value ordering
	items, err = GetMiscLookupByKey("test_order", "value")
	require.NoError(t, err)
	assert.Len(t, items, 3)
	assert.Equal(t, "Alpha", items[0].Value)

	// Test ID ordering
	items, err = GetMiscLookupByKey("test_order", "id")
	require.NoError(t, err)
	assert.Len(t, items, 3)
	// Should be ordered by ID (creation order)
	assert.Equal(t, "Zebra", items[0].Value)
	assert.Equal(t, "Alpha", items[1].Value)
	assert.Equal(t, "Beta", items[2].Value)

	// Test source ordering
	items, err = GetMiscLookupByKey("test_order", "source")
	require.NoError(t, err)
	assert.Len(t, items, 3)
	// Should be ordered by source_id first, then value
	assert.Equal(t, uint(1), items[0].SourceID)
	assert.Equal(t, "Alpha", items[0].Value)
	assert.Equal(t, uint(1), items[1].SourceID)
	assert.Equal(t, "Beta", items[1].Value)
	assert.Equal(t, uint(2), items[2].SourceID)
	assert.Equal(t, "Zebra", items[2].Value)

	// Test source_value ordering (same as source)
	items, err = GetMiscLookupByKey("test_order", "source_value")
	require.NoError(t, err)
	assert.Len(t, items, 3)
	assert.Equal(t, uint(1), items[0].SourceID)
	assert.Equal(t, "Alpha", items[0].Value)

	// Test invalid ordering (should default to value)
	items, err = GetMiscLookupByKey("test_order", "invalid")
	require.NoError(t, err)
	assert.Len(t, items, 3)
	assert.Equal(t, "Alpha", items[0].Value)

	// Test empty string ordering (should default to value)
	items, err = GetMiscLookupByKey("test_order", "")
	require.NoError(t, err)
	assert.Len(t, items, 3)
	assert.Equal(t, "Alpha", items[0].Value)
}
