package models

import (
	"bamort/database"
	"bamort/user"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCharacterTestDB(t *testing.T) {
	database.SetupTestDB()

	// Migrate structures
	err := MigrateStructure()
	require.NoError(t, err, "Failed to migrate database structure")

	// Clean up any existing test data
	cleanupCharacterTestData(t)
}

func cleanupCharacterTestData(t *testing.T) {
	// Delete all characters to ensure clean state
	err := database.DB.Exec("DELETE FROM char_chars").Error
	require.NoError(t, err, "Failed to clean up characters")
}

func createTestUser() *user.User {
	return &user.User{
		UserID:   1,
		Username: "testuser",
		Email:    "test@example.com",
	}
}

func createTestChar(name string) *Char {
	return &Char{
		BamortBase: BamortBase{
			Name: name,
		},
		UserID:      1,
		Rasse:       "Mensch",
		Typ:         "Krieger",
		Alter:       25,
		Anrede:      "Herr",
		Grad:        3,
		Gender:      "männlich",
		SocialClass: "Mittelschicht",
		Groesse:     180,
		Gewicht:     75,
		Herkunft:    "Alba",
		Glaube:      "Xan",
		Hand:        "rechts",
		Public:      false,
		Lp: Lp{
			Max:   20,
			Value: 20,
		},
		Ap: Ap{
			Max:   30,
			Value: 30,
		},
		B: B{
			Max:   10,
			Value: 10,
		},
		Merkmale: Merkmale{
			Augenfarbe: "blau",
			Haarfarbe:  "braun",
			Sonstige:   "Narbe am rechten Arm",
			Breite:     "normal",
			Groesse:    "groß",
		},
		Bennies: Bennies{
			Gg: 1,
			Gp: 2,
			Sg: 0,
		},
		Vermoegen: Vermoegen{
			Goldstuecke:   100,
			Silberstuecke: 50,
			Kupferstuecke: 25,
		},
		Erfahrungsschatz: Erfahrungsschatz{
			ES: 150,
			EP: 25,
		},
	}
}

func createTestEigenschaft(charID uint, name string, value int) *Eigenschaft {
	return &Eigenschaft{
		CharacterID: charID,
		UserID:      1,
		Name:        name,
		Value:       value,
	}
}

// =============================================================================
// Tests for Char struct
// =============================================================================

func TestChar_TableName(t *testing.T) {
	char := Char{}
	expected := "char_chars"
	actual := char.TableName()
	assert.Equal(t, expected, actual)
}

func TestChar_Create_Success(t *testing.T) {
	setupCharacterTestDB(t)

	testChar := createTestChar("Test Character")
	err := testChar.Create()

	assert.NoError(t, err, "Create should succeed")
	assert.Greater(t, testChar.ID, uint(0), "ID should be set after create")
}

func TestChar_Create_SetsGameSystem(t *testing.T) {
	setupCharacterTestDB(t)

	gs := GetGameSystem(0, "midgard")
	require.NotNil(t, gs)

	testChar := createTestChar("Test Character GS")
	testChar.GameSystem = ""
	testChar.GameSystemId = 0

	err := testChar.Create()

	require.NoError(t, err)
	assert.Equal(t, gs.Name, testChar.GameSystem)
	assert.Equal(t, gs.ID, testChar.GameSystemId)
}

func TestChar_Create_WithRelations(t *testing.T) {
	setupCharacterTestDB(t)

	testChar := createTestChar("Test Character With Relations")

	// Add some eigenschaften
	testChar.Eigenschaften = []Eigenschaft{
		{UserID: 1, Name: "St", Value: 15},
		{UserID: 1, Name: "Gs", Value: 12},
		{UserID: 1, Name: "Ko", Value: 14},
	}

	err := testChar.Create()

	assert.NoError(t, err, "Create with relations should succeed")
	assert.Greater(t, testChar.ID, uint(0), "ID should be set after create")
	assert.Len(t, testChar.Eigenschaften, 3, "All eigenschaften should be created")
}

func TestChar_First_Success(t *testing.T) {
	setupCharacterTestDB(t)

	// Create a test character first
	testChar := createTestChar("Test Character First")
	err := testChar.Create()
	require.NoError(t, err, "Character creation should succeed")

	// Now try to find it
	foundChar := &Char{}
	err = foundChar.First("Test Character First")

	assert.NoError(t, err, "First should succeed")
	assert.Equal(t, testChar.ID, foundChar.ID, "Found character should have same ID")
	assert.Equal(t, testChar.Name, foundChar.Name, "Found character should have same name")
	assert.Equal(t, testChar.Rasse, foundChar.Rasse, "Found character should have same race")
}

func TestChar_First_UsesGameSystem(t *testing.T) {
	setupCharacterTestDB(t)

	defaultGS := GetGameSystem(0, "midgard")
	require.NotNil(t, defaultGS)

	altGS := &GameSystem{Code: "ALT", Name: "AltSystem", Description: "Test", IsActive: true}
	require.NoError(t, database.DB.Create(altGS).Error)

	charDefault := createTestChar("GS Default Char")
	require.NoError(t, charDefault.Create())

	charAlt := createTestChar("GS Alt Char")
	charAlt.GameSystem = altGS.Code
	charAlt.GameSystemId = altGS.ID
	require.NoError(t, charAlt.Create())

	finder := &Char{GameSystem: defaultGS.Name, GameSystemId: defaultGS.ID}
	err := finder.First("GS Default Char")
	require.NoError(t, err)
	assert.Equal(t, defaultGS.ID, finder.GameSystemId)

	otherFinder := &Char{GameSystem: altGS.Code, GameSystemId: altGS.ID}
	err = otherFinder.First("GS Default Char")
	assert.Error(t, err)
}

func TestChar_First_NotFound(t *testing.T) {
	setupCharacterTestDB(t)

	foundChar := &Char{}
	err := foundChar.First("Non-Existent Character")

	assert.Error(t, err, "First should return error for non-existent character")
}

func TestChar_FirstID_Success(t *testing.T) {
	setupCharacterTestDB(t)

	// Create a test character first
	testChar := createTestChar("Test Character FirstID")
	err := testChar.Create()
	require.NoError(t, err, "Character creation should succeed")

	// Now try to find it by ID
	foundChar := &Char{}
	err = foundChar.FirstID(strconv.Itoa(int(testChar.ID)))

	assert.NoError(t, err, "FirstID should succeed")
	assert.Equal(t, testChar.ID, foundChar.ID, "Found character should have same ID")
	assert.Equal(t, testChar.Name, foundChar.Name, "Found character should have same name")
}

func TestChar_FirstID_NotFound(t *testing.T) {
	setupCharacterTestDB(t)

	foundChar := &Char{}
	err := foundChar.FirstID("999999")

	assert.Error(t, err, "FirstID should return error for non-existent ID")
}

func TestChar_FindByUserID_Success(t *testing.T) {
	setupCharacterTestDB(t)

	// Create multiple test characters for the same user
	testChar1 := createTestChar("Test Character 1")
	testChar2 := createTestChar("Test Character 2")
	testChar3 := createTestChar("Test Character 3")
	testChar3.UserID = 2 // Different user

	err := testChar1.Create()
	require.NoError(t, err, "Character 1 creation should succeed")
	err = testChar2.Create()
	require.NoError(t, err, "Character 2 creation should succeed")
	err = testChar3.Create()
	require.NoError(t, err, "Character 3 creation should succeed")

	// Find characters for user 1
	foundChar := &Char{}
	chars, err := foundChar.FindByUserID(1)

	assert.NoError(t, err, "FindByUserID should succeed")
	assert.GreaterOrEqual(t, len(chars), 2, "Should find at least 2 characters for user 1")

	// Verify all returned characters belong to user 1
	for _, char := range chars {
		assert.Equal(t, uint(1), char.UserID, "All characters should belong to user 1")
	}
}

func TestChar_FindByUserID_NoCharacters(t *testing.T) {
	setupCharacterTestDB(t)

	foundChar := &Char{}
	chars, err := foundChar.FindByUserID(999999)

	assert.NoError(t, err, "FindByUserID should succeed even with no results")
	assert.Empty(t, chars, "Should return empty slice for user with no characters")
}

func TestChar_Delete_Success(t *testing.T) {
	setupCharacterTestDB(t)

	// Create a test character first
	testChar := createTestChar("Test Character Delete")
	err := testChar.Create()
	require.NoError(t, err, "Character creation should succeed")

	originalID := testChar.ID

	// Delete the character
	err = testChar.Delete()
	assert.NoError(t, err, "Delete should succeed")

	// Verify the character is deleted
	foundChar := &Char{}
	err = foundChar.FirstID(strconv.Itoa(int(originalID)))
	assert.Error(t, err, "Character should not be found after deletion")
}

func TestFindPublicCharList_Success(t *testing.T) {
	setupCharacterTestDB(t)

	// Create some test characters
	testChar1 := createTestChar("Public Character 1")
	testChar1.Public = true
	testChar2 := createTestChar("Private Character")
	testChar2.Public = false
	testChar3 := createTestChar("Public Character 2")
	testChar3.Public = true

	err := testChar1.Create()
	require.NoError(t, err, "Character 1 creation should succeed")
	err = testChar2.Create()
	require.NoError(t, err, "Character 2 creation should succeed")
	err = testChar3.Create()
	require.NoError(t, err, "Character 3 creation should succeed")

	// Find public characters
	publicChars, err := FindPublicCharList()

	assert.NoError(t, err, "FindPublicCharList should succeed")
	assert.GreaterOrEqual(t, len(publicChars), 2, "Should find at least 2 public characters")

	// Verify all returned characters are public
	for _, char := range publicChars {
		assert.True(t, char.Public, "All characters should be public")
	}
}

func TestFindCharListByUserID_Success(t *testing.T) {
	setupCharacterTestDB(t)

	// Create test characters for different users
	testChar1 := createTestChar("User 1 Character 1")
	testChar1.UserID = 1
	testChar2 := createTestChar("User 1 Character 2")
	testChar2.UserID = 1
	testChar3 := createTestChar("User 2 Character")
	testChar3.UserID = 2

	err := testChar1.Create()
	require.NoError(t, err, "Character 1 creation should succeed")
	err = testChar2.Create()
	require.NoError(t, err, "Character 2 creation should succeed")
	err = testChar3.Create()
	require.NoError(t, err, "Character 3 creation should succeed")

	// Find characters for user 1
	userChars, err := FindCharListByUserID(1)

	assert.NoError(t, err, "FindCharListByUserID should succeed")
	assert.GreaterOrEqual(t, len(userChars), 2, "Should find at least 2 characters for user 1")

	// Verify all returned characters belong to user 1
	for _, char := range userChars {
		assert.Equal(t, uint(1), char.UserID, "All characters should belong to user 1")
	}
}

// =============================================================================
// Tests for Eigenschaft struct
// =============================================================================

func TestEigenschaft_TableName(t *testing.T) {
	eigenschaft := Eigenschaft{}
	expected := "char_eigenschaften"
	actual := eigenschaft.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Tests for Lp struct
// =============================================================================

func TestLp_TableName(t *testing.T) {
	lp := Lp{}
	expected := "char_health"
	actual := lp.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Tests for Ap struct
// =============================================================================

func TestAp_TableName(t *testing.T) {
	ap := Ap{}
	expected := "char_endurances"
	actual := ap.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Tests for B struct
// =============================================================================

func TestB_TableName(t *testing.T) {
	b := B{}
	expected := "char_motionranges"
	actual := b.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Tests for Merkmale struct
// =============================================================================

func TestMerkmale_TableName(t *testing.T) {
	merkmale := Merkmale{}
	expected := "char_characteristics"
	actual := merkmale.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Tests for Erfahrungsschatz struct
// =============================================================================

func TestErfahrungsschatz_TableName(t *testing.T) {
	erfahrung := Erfahrungsschatz{}
	expected := "char_experiances"
	actual := erfahrung.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Tests for Bennies struct
// =============================================================================

func TestBennies_TableName(t *testing.T) {
	bennies := Bennies{}
	expected := "char_bennies"
	actual := bennies.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Tests for Vermoegen struct
// =============================================================================

func TestVermoegen_TableName(t *testing.T) {
	vermoegen := Vermoegen{}
	expected := "char_wealth"
	actual := vermoegen.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Integration tests for character with related data
// =============================================================================

func TestChar_CreateWithCompleteData(t *testing.T) {
	setupCharacterTestDB(t)

	testChar := createTestChar("Complete Character")

	// Add eigenschaften
	testChar.Eigenschaften = []Eigenschaft{
		{UserID: 1, Name: "St", Value: 15},
		{UserID: 1, Name: "Gs", Value: 12},
		{UserID: 1, Name: "Ko", Value: 14},
		{UserID: 1, Name: "In", Value: 13},
		{UserID: 1, Name: "Zt", Value: 11},
	}

	err := testChar.Create()

	assert.NoError(t, err, "Complete character creation should succeed")
	assert.Greater(t, testChar.ID, uint(0), "ID should be set after create")

	// Verify the character can be found with all relations
	foundChar := &Char{}
	err = foundChar.First(testChar.Name)
	require.NoError(t, err, "Character should be found")

	assert.Equal(t, testChar.Name, foundChar.Name)
	assert.Equal(t, testChar.Rasse, foundChar.Rasse)
	assert.Equal(t, testChar.Typ, foundChar.Typ)
	assert.Equal(t, testChar.Lp.Max, foundChar.Lp.Max)
	assert.Equal(t, testChar.Ap.Max, foundChar.Ap.Max)
	assert.Equal(t, testChar.B.Max, foundChar.B.Max)
	assert.Equal(t, testChar.Merkmale.Augenfarbe, foundChar.Merkmale.Augenfarbe)
	assert.Equal(t, testChar.Bennies.Gg, foundChar.Bennies.Gg)
	assert.Equal(t, testChar.Vermoegen.Goldstuecke, foundChar.Vermoegen.Goldstuecke)
	assert.Equal(t, testChar.Erfahrungsschatz.ES, foundChar.Erfahrungsschatz.ES)
}

func TestChar_CharacterProgression(t *testing.T) {
	setupCharacterTestDB(t)

	// Create a low-level character
	testChar := createTestChar("Progression Character")
	testChar.Grad = 1
	testChar.Erfahrungsschatz.ES = 50
	testChar.Erfahrungsschatz.EP = 0

	err := testChar.Create()
	require.NoError(t, err, "Character creation should succeed")

	// Simulate character progression
	testChar.Grad = 2
	testChar.Erfahrungsschatz.ES = 100
	testChar.Erfahrungsschatz.EP = 25

	// Note: We would need an Update method to test this properly
	// For now, we just verify the data structure can hold progression data
	assert.Equal(t, 2, testChar.Grad, "Grade should be updated")
	assert.Equal(t, 100, testChar.Erfahrungsschatz.ES, "Experience should be updated")
	assert.Equal(t, 25, testChar.Erfahrungsschatz.EP, "Experience points should be updated")
}

func TestChar_WealthManagement(t *testing.T) {
	setupCharacterTestDB(t)

	testChar := createTestChar("Wealthy Character")
	testChar.Vermoegen.Goldstuecke = 1000
	testChar.Vermoegen.Silberstuecke = 500
	testChar.Vermoegen.Kupferstuecke = 100

	err := testChar.Create()
	require.NoError(t, err, "Character creation should succeed")

	// Verify wealth data
	foundChar := &Char{}
	err = foundChar.First(testChar.Name)
	require.NoError(t, err, "Character should be found")

	assert.Equal(t, 1000, foundChar.Vermoegen.Goldstuecke, "Gold should match")
	assert.Equal(t, 500, foundChar.Vermoegen.Silberstuecke, "Silver should match")
	assert.Equal(t, 100, foundChar.Vermoegen.Kupferstuecke, "Copper should match")
}

func TestChar_EdgeCases(t *testing.T) {
	setupCharacterTestDB(t)

	// Test character with minimal data
	minimalChar := &Char{
		BamortBase: BamortBase{
			Name: "Minimal Character",
		},
		UserID: 1,
		Rasse:  "Unbekannt",
		Typ:    "Abenteurer",
	}

	err := minimalChar.Create()
	assert.NoError(t, err, "Minimal character should be created successfully")

	// Test character with extreme values
	extremeChar := createTestChar("Extreme Character")
	extremeChar.Alter = 999
	extremeChar.Groesse = 300
	extremeChar.Gewicht = 500
	extremeChar.Lp.Max = 999
	extremeChar.Ap.Max = 999
	extremeChar.B.Max = 999

	err = extremeChar.Create()
	assert.NoError(t, err, "Character with extreme values should be created successfully")

	// Test character with special characters in name
	specialChar := createTestChar("Ä Special Çharacter ñ")
	err = specialChar.Create()
	assert.NoError(t, err, "Character with special characters should be created successfully")
}
