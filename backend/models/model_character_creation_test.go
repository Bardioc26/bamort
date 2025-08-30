package models

import (
	"bamort/database"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCharacterCreationTestDB(t *testing.T) {
	database.SetupTestDB()

	// Migrate structures
	err := MigrateStructure()
	require.NoError(t, err, "Failed to migrate database structure")

	// Clean up any existing character creation sessions
	cleanupCharacterCreationTestData(t)
}

func cleanupCharacterCreationTestData(t *testing.T) {
	// Delete all character creation sessions to ensure clean state
	err := database.DB.Exec("DELETE FROM char_char_creation_session").Error
	require.NoError(t, err, "Failed to clean up character creation sessions")
}

func createTestCharacterCreationSession(sessionID string) *CharacterCreationSession {
	return &CharacterCreationSession{
		ID:         sessionID,
		UserID:     1,
		Name:       "Test Character",
		Geschlecht: "männlich",
		Rasse:      "Mensch",
		Typ:        "Krieger",
		Herkunft:   "Alba",
		Stand:      "Mittelschicht",
		Glaube:     "Xan",
		Attributes: AttributesData{
			ST: 15,
			GS: 12,
			GW: 14,
			KO: 16,
			IN: 13,
			ZT: 11,
			AU: 10,
		},
		DerivedValues: DerivedValuesData{
			PA:                    10,
			WK:                    13,
			LPMax:                 20,
			APMax:                 30,
			BMax:                  15,
			ResistenzKoerper:      16,
			ResistenzGeist:        13,
			ResistenzBonusKoerper: 2,
			ResistenzBonusGeist:   1,
			Abwehr:                12,
			AbwehrBonus:           0,
			AusdauerBonus:         2,
			AngriffsBonus:         1,
			Zaubern:               11,
			ZauberBonus:           0,
			Raufen:                13,
			SchadensBonus:         2,
			SG:                    2,
			GG:                    1,
			GP:                    3,
		},
		Skills: CharacterCreationSkills{
			{Name: "Athletik", Level: 8, Category: "Körper", Cost: 3},
			{Name: "Schwimmen", Level: 6, Category: "Körper", Cost: 2},
		},
		Spells: CharacterCreationSpells{
			{Name: "Licht", Cost: 1},
			{Name: "Handauflegen", Cost: 2},
		},
		SkillPoints: SkillPointsData{
			"Körper": 10,
			"Geist":  15,
			"Kampf":  8,
		},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		CurrentStep: 3,
	}
}

func createTestAttributesData() AttributesData {
	return AttributesData{
		ST: 15,
		GS: 12,
		GW: 14,
		KO: 16,
		IN: 13,
		ZT: 11,
		AU: 10,
	}
}

func createTestDerivedValuesData() DerivedValuesData {
	return DerivedValuesData{
		PA:                    10,
		WK:                    13,
		LPMax:                 20,
		APMax:                 30,
		BMax:                  15,
		ResistenzKoerper:      16,
		ResistenzGeist:        13,
		ResistenzBonusKoerper: 2,
		ResistenzBonusGeist:   1,
		Abwehr:                12,
		AbwehrBonus:           0,
		AusdauerBonus:         2,
		AngriffsBonus:         1,
		Zaubern:               11,
		ZauberBonus:           0,
		Raufen:                13,
		SchadensBonus:         2,
		SG:                    2,
		GG:                    1,
		GP:                    3,
	}
}

// =============================================================================
// Tests for CharacterCreationSession struct
// =============================================================================

func TestCharacterCreationSession_TableName(t *testing.T) {
	session := CharacterCreationSession{}
	expected := "char_char_creation_session"
	actual := session.TableName()
	assert.Equal(t, expected, actual)
}

// =============================================================================
// Tests for AttributesData JSON handling
// =============================================================================

func TestAttributesData_Value_Success(t *testing.T) {
	attrs := createTestAttributesData()

	value, err := attrs.Value()

	assert.NoError(t, err, "Value should succeed")
	assert.IsType(t, []byte{}, value, "Value should return byte slice")

	// Verify it's valid JSON
	var result AttributesData
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err, "Should be valid JSON")
	assert.Equal(t, attrs.ST, result.ST, "ST should match")
	assert.Equal(t, attrs.GS, result.GS, "GS should match")
}

func TestAttributesData_Scan_Success(t *testing.T) {
	original := createTestAttributesData()
	jsonData, err := json.Marshal(original)
	require.NoError(t, err, "Marshal should succeed")

	var attrs AttributesData
	err = attrs.Scan(jsonData)

	assert.NoError(t, err, "Scan should succeed")
	assert.Equal(t, original.ST, attrs.ST, "ST should match")
	assert.Equal(t, original.GS, attrs.GS, "GS should match")
	assert.Equal(t, original.GW, attrs.GW, "GW should match")
	assert.Equal(t, original.KO, attrs.KO, "KO should match")
}

func TestAttributesData_Scan_Nil(t *testing.T) {
	var attrs AttributesData
	err := attrs.Scan(nil)

	assert.NoError(t, err, "Scan with nil should succeed")
}

func TestAttributesData_Scan_InvalidData(t *testing.T) {
	var attrs AttributesData
	err := attrs.Scan("invalid json")

	assert.NoError(t, err, "Scan with invalid data should not error but may have zero values")
}

// =============================================================================
// Tests for DerivedValuesData JSON handling
// =============================================================================

func TestDerivedValuesData_Value_Success(t *testing.T) {
	derived := createTestDerivedValuesData()

	value, err := derived.Value()

	assert.NoError(t, err, "Value should succeed")
	assert.IsType(t, []byte{}, value, "Value should return byte slice")

	// Verify it's valid JSON
	var result DerivedValuesData
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err, "Should be valid JSON")
	assert.Equal(t, derived.PA, result.PA, "PA should match")
	assert.Equal(t, derived.WK, result.WK, "WK should match")
}

func TestDerivedValuesData_Scan_Success(t *testing.T) {
	original := createTestDerivedValuesData()
	jsonData, err := json.Marshal(original)
	require.NoError(t, err, "Marshal should succeed")

	var derived DerivedValuesData
	err = derived.Scan(jsonData)

	assert.NoError(t, err, "Scan should succeed")
	assert.Equal(t, original.PA, derived.PA, "PA should match")
	assert.Equal(t, original.WK, derived.WK, "WK should match")
	assert.Equal(t, original.LPMax, derived.LPMax, "LPMax should match")
}

func TestDerivedValuesData_Scan_Nil(t *testing.T) {
	var derived DerivedValuesData
	err := derived.Scan(nil)

	assert.NoError(t, err, "Scan with nil should succeed")
}

// =============================================================================
// Tests for SkillPointsData JSON handling
// =============================================================================

func TestSkillPointsData_Value_Success(t *testing.T) {
	skillPoints := SkillPointsData{
		"Körper": 10,
		"Geist":  15,
		"Kampf":  8,
	}

	value, err := skillPoints.Value()

	assert.NoError(t, err, "Value should succeed")
	assert.IsType(t, []byte{}, value, "Value should return byte slice")

	// Verify it's valid JSON
	var result SkillPointsData
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err, "Should be valid JSON")
	assert.Equal(t, 10, result["Körper"], "Körper points should match")
	assert.Equal(t, 15, result["Geist"], "Geist points should match")
}

func TestSkillPointsData_Scan_Success(t *testing.T) {
	original := SkillPointsData{
		"Körper": 10,
		"Geist":  15,
		"Kampf":  8,
	}
	jsonData, err := json.Marshal(original)
	require.NoError(t, err, "Marshal should succeed")

	var skillPoints SkillPointsData
	err = skillPoints.Scan(jsonData)

	assert.NoError(t, err, "Scan should succeed")
	assert.Equal(t, 10, skillPoints["Körper"], "Körper points should match")
	assert.Equal(t, 15, skillPoints["Geist"], "Geist points should match")
	assert.Equal(t, 8, skillPoints["Kampf"], "Kampf points should match")
}

func TestSkillPointsData_Scan_Nil(t *testing.T) {
	var skillPoints SkillPointsData
	err := skillPoints.Scan(nil)

	assert.NoError(t, err, "Scan with nil should succeed")
}

// =============================================================================
// Tests for CharacterCreationSkills JSON handling
// =============================================================================

func TestCharacterCreationSkills_Value_Success(t *testing.T) {
	skills := CharacterCreationSkills{
		{Name: "Athletik", Level: 8, Category: "Körper", Cost: 3},
		{Name: "Schwimmen", Level: 6, Category: "Körper", Cost: 2},
	}

	value, err := skills.Value()

	assert.NoError(t, err, "Value should succeed")
	assert.IsType(t, []byte{}, value, "Value should return byte slice")

	// Verify it's valid JSON
	var result CharacterCreationSkills
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err, "Should be valid JSON")
	assert.Len(t, result, 2, "Should have 2 skills")
	assert.Equal(t, "Athletik", result[0].Name, "First skill name should match")
}

func TestCharacterCreationSkills_Value_Empty(t *testing.T) {
	skills := CharacterCreationSkills{}

	value, err := skills.Value()

	assert.NoError(t, err, "Value should succeed")
	assert.Equal(t, "[]", value, "Empty skills should return empty JSON array")
}

func TestCharacterCreationSkills_Scan_Success(t *testing.T) {
	original := CharacterCreationSkills{
		{Name: "Athletik", Level: 8, Category: "Körper", Cost: 3},
		{Name: "Schwimmen", Level: 6, Category: "Körper", Cost: 2},
	}
	jsonData, err := json.Marshal(original)
	require.NoError(t, err, "Marshal should succeed")

	var skills CharacterCreationSkills
	err = skills.Scan(jsonData)

	assert.NoError(t, err, "Scan should succeed")
	assert.Len(t, skills, 2, "Should have 2 skills")
	assert.Equal(t, "Athletik", skills[0].Name, "First skill name should match")
	assert.Equal(t, 8, skills[0].Level, "First skill level should match")
}

func TestCharacterCreationSkills_Scan_Nil(t *testing.T) {
	var skills CharacterCreationSkills
	err := skills.Scan(nil)

	assert.NoError(t, err, "Scan with nil should succeed")
	assert.Empty(t, skills, "Should be empty slice")
}

func TestCharacterCreationSkills_Scan_InvalidData(t *testing.T) {
	var skills CharacterCreationSkills
	err := skills.Scan("invalid json")

	assert.NoError(t, err, "Scan with invalid data should not error")
	assert.Empty(t, skills, "Should be empty slice")
}

// =============================================================================
// Tests for CharacterCreationSpells JSON handling
// =============================================================================

func TestCharacterCreationSpells_Value_Success(t *testing.T) {
	spells := CharacterCreationSpells{
		{Name: "Licht", Cost: 1},
		{Name: "Handauflegen", Cost: 2},
	}

	value, err := spells.Value()

	assert.NoError(t, err, "Value should succeed")
	assert.IsType(t, []byte{}, value, "Value should return byte slice")

	// Verify it's valid JSON
	var result CharacterCreationSpells
	err = json.Unmarshal(value.([]byte), &result)
	assert.NoError(t, err, "Should be valid JSON")
	assert.Len(t, result, 2, "Should have 2 spells")
	assert.Equal(t, "Licht", result[0].Name, "First spell name should match")
}

func TestCharacterCreationSpells_Value_Empty(t *testing.T) {
	spells := CharacterCreationSpells{}

	value, err := spells.Value()

	assert.NoError(t, err, "Value should succeed")
	assert.Equal(t, "[]", value, "Empty spells should return empty JSON array")
}

func TestCharacterCreationSpells_Scan_Success(t *testing.T) {
	original := CharacterCreationSpells{
		{Name: "Licht", Cost: 1},
		{Name: "Handauflegen", Cost: 2},
	}
	jsonData, err := json.Marshal(original)
	require.NoError(t, err, "Marshal should succeed")

	var spells CharacterCreationSpells
	err = spells.Scan(jsonData)

	assert.NoError(t, err, "Scan should succeed")
	assert.Len(t, spells, 2, "Should have 2 spells")
	assert.Equal(t, "Licht", spells[0].Name, "First spell name should match")
	assert.Equal(t, 1, spells[0].Cost, "First spell cost should match")
}

func TestCharacterCreationSpells_Scan_Nil(t *testing.T) {
	var spells CharacterCreationSpells
	err := spells.Scan(nil)

	assert.NoError(t, err, "Scan with nil should succeed")
	assert.Empty(t, spells, "Should be empty slice")
}

func TestCharacterCreationSpells_Scan_InvalidData(t *testing.T) {
	var spells CharacterCreationSpells
	err := spells.Scan("invalid json")

	assert.NoError(t, err, "Scan with invalid data should not error")
	assert.Empty(t, spells, "Should be empty slice")
}

// =============================================================================
// Tests for CleanupExpiredSessions function
// =============================================================================

func TestCleanupExpiredSessions_Success(t *testing.T) {
	setupCharacterCreationTestDB(t)

	// Create test sessions - some expired, some not
	expiredSession1 := createTestCharacterCreationSession("expired-1")
	expiredSession1.ExpiresAt = time.Now().Add(-1 * time.Hour) // Expired 1 hour ago

	expiredSession2 := createTestCharacterCreationSession("expired-2")
	expiredSession2.ExpiresAt = time.Now().Add(-2 * time.Hour) // Expired 2 hours ago

	validSession := createTestCharacterCreationSession("valid-1")
	validSession.ExpiresAt = time.Now().Add(1 * time.Hour) // Expires in 1 hour

	// Save sessions to database
	err := database.DB.Create(expiredSession1).Error
	require.NoError(t, err, "Should create expired session 1")
	err = database.DB.Create(expiredSession2).Error
	require.NoError(t, err, "Should create expired session 2")
	err = database.DB.Create(validSession).Error
	require.NoError(t, err, "Should create valid session")

	// Cleanup expired sessions
	err = CleanupExpiredSessions(database.DB)
	assert.NoError(t, err, "CleanupExpiredSessions should succeed")

	// Verify expired sessions are deleted
	var count int64
	database.DB.Model(&CharacterCreationSession{}).Where("id IN ?", []string{"expired-1", "expired-2"}).Count(&count)
	assert.Equal(t, int64(0), count, "Expired sessions should be deleted")

	// Verify valid session still exists
	database.DB.Model(&CharacterCreationSession{}).Where("id = ?", "valid-1").Count(&count)
	assert.Equal(t, int64(1), count, "Valid session should still exist")
}

func TestCleanupExpiredSessions_NoExpiredSessions(t *testing.T) {
	setupCharacterCreationTestDB(t)

	// Create only valid sessions
	validSession1 := createTestCharacterCreationSession("valid-1")
	validSession1.ExpiresAt = time.Now().Add(1 * time.Hour)

	validSession2 := createTestCharacterCreationSession("valid-2")
	validSession2.ExpiresAt = time.Now().Add(2 * time.Hour)

	// Save sessions to database
	err := database.DB.Create(validSession1).Error
	require.NoError(t, err, "Should create valid session 1")
	err = database.DB.Create(validSession2).Error
	require.NoError(t, err, "Should create valid session 2")

	// Cleanup expired sessions
	err = CleanupExpiredSessions(database.DB)
	assert.NoError(t, err, "CleanupExpiredSessions should succeed")

	// Verify all sessions still exist
	var count int64
	database.DB.Model(&CharacterCreationSession{}).Count(&count)
	assert.Equal(t, int64(2), count, "All valid sessions should still exist")
}

// =============================================================================
// Tests for GetUserSessions function
// =============================================================================

func TestGetUserSessions_Success(t *testing.T) {
	setupCharacterCreationTestDB(t)

	// Create sessions for different users
	user1Session1 := createTestCharacterCreationSession("user1-session1")
	user1Session1.UserID = 1
	user1Session1.ExpiresAt = time.Now().Add(1 * time.Hour)

	user1Session2 := createTestCharacterCreationSession("user1-session2")
	user1Session2.UserID = 1
	user1Session2.ExpiresAt = time.Now().Add(2 * time.Hour)

	user2Session := createTestCharacterCreationSession("user2-session1")
	user2Session.UserID = 2
	user2Session.ExpiresAt = time.Now().Add(1 * time.Hour)

	expiredUser1Session := createTestCharacterCreationSession("user1-expired")
	expiredUser1Session.UserID = 1
	expiredUser1Session.ExpiresAt = time.Now().Add(-1 * time.Hour)

	// Save sessions to database
	sessions := []*CharacterCreationSession{user1Session1, user1Session2, user2Session, expiredUser1Session}
	for _, session := range sessions {
		err := database.DB.Create(session).Error
		require.NoError(t, err, "Should create session")
	}

	// Get sessions for user 1
	userSessions, err := GetUserSessions(database.DB, 1)

	assert.NoError(t, err, "GetUserSessions should succeed")
	assert.Len(t, userSessions, 2, "Should return 2 valid sessions for user 1")

	// Verify returned sessions belong to user 1 and are not expired
	for _, session := range userSessions {
		assert.Equal(t, uint(1), session.UserID, "All sessions should belong to user 1")
		assert.True(t, session.ExpiresAt.After(time.Now()), "All sessions should not be expired")
	}
}

func TestGetUserSessions_NoSessions(t *testing.T) {
	setupCharacterCreationTestDB(t)

	// Get sessions for user that has no sessions
	userSessions, err := GetUserSessions(database.DB, 999)

	assert.NoError(t, err, "GetUserSessions should succeed")
	assert.Empty(t, userSessions, "Should return empty slice for user with no sessions")
}

func TestGetUserSessions_OnlyExpiredSessions(t *testing.T) {
	setupCharacterCreationTestDB(t)

	// Create only expired sessions for user 1
	expiredSession1 := createTestCharacterCreationSession("expired-1")
	expiredSession1.UserID = 1
	expiredSession1.ExpiresAt = time.Now().Add(-1 * time.Hour)

	expiredSession2 := createTestCharacterCreationSession("expired-2")
	expiredSession2.UserID = 1
	expiredSession2.ExpiresAt = time.Now().Add(-2 * time.Hour)

	// Save sessions to database
	err := database.DB.Create(expiredSession1).Error
	require.NoError(t, err, "Should create expired session 1")
	err = database.DB.Create(expiredSession2).Error
	require.NoError(t, err, "Should create expired session 2")

	// Get sessions for user 1
	userSessions, err := GetUserSessions(database.DB, 1)

	assert.NoError(t, err, "GetUserSessions should succeed")
	assert.Empty(t, userSessions, "Should return empty slice when all sessions are expired")
}

// =============================================================================
// Integration tests for complete character creation workflow
// =============================================================================

func TestCharacterCreationSession_CompleteWorkflow(t *testing.T) {
	setupCharacterCreationTestDB(t)

	// Create a character creation session
	session := createTestCharacterCreationSession("workflow-test")

	// Save session to database
	err := database.DB.Create(session).Error
	require.NoError(t, err, "Should create character creation session")

	// Verify the session was saved with all JSON fields properly serialized
	var savedSession CharacterCreationSession
	err = database.DB.Where("id = ?", "workflow-test").First(&savedSession).Error
	require.NoError(t, err, "Should find saved session")

	// Verify basic fields
	assert.Equal(t, session.Name, savedSession.Name, "Name should match")
	assert.Equal(t, session.Rasse, savedSession.Rasse, "Race should match")
	assert.Equal(t, session.Typ, savedSession.Typ, "Type should match")

	// Verify JSON fields were properly saved and loaded
	assert.Equal(t, session.Attributes.ST, savedSession.Attributes.ST, "Attributes ST should match")
	assert.Equal(t, session.DerivedValues.PA, savedSession.DerivedValues.PA, "Derived values PA should match")
	assert.Len(t, savedSession.Skills, 2, "Should have 2 skills")
	assert.Equal(t, session.Skills[0].Name, savedSession.Skills[0].Name, "First skill name should match")
	assert.Len(t, savedSession.Spells, 2, "Should have 2 spells")
	assert.Equal(t, session.Spells[0].Name, savedSession.Spells[0].Name, "First spell name should match")
	assert.Equal(t, session.SkillPoints["Körper"], savedSession.SkillPoints["Körper"], "Skill points should match")
}

func TestCharacterCreationSession_ProgressionSteps(t *testing.T) {
	setupCharacterCreationTestDB(t)

	// Test progression through creation steps
	session := createTestCharacterCreationSession("progression-test")
	session.CurrentStep = 1 // Start at basic info

	// Save initial session
	err := database.DB.Create(session).Error
	require.NoError(t, err, "Should create character creation session")

	// Simulate progression to step 2 (attributes)
	session.CurrentStep = 2
	session.Attributes = AttributesData{ST: 14, GS: 13, GW: 12, KO: 15, IN: 11, ZT: 10, AU: 9}
	err = database.DB.Save(session).Error
	require.NoError(t, err, "Should update session to step 2")

	// Simulate progression to step 3 (derived values)
	session.CurrentStep = 3
	session.DerivedValues = createTestDerivedValuesData()
	err = database.DB.Save(session).Error
	require.NoError(t, err, "Should update session to step 3")

	// Simulate progression to step 4 (skills)
	session.CurrentStep = 4
	session.Skills = CharacterCreationSkills{
		{Name: "Athletik", Level: 8, Category: "Körper", Cost: 3},
	}
	err = database.DB.Save(session).Error
	require.NoError(t, err, "Should update session to step 4")

	// Verify final state
	var finalSession CharacterCreationSession
	err = database.DB.Where("id = ?", "progression-test").First(&finalSession).Error
	require.NoError(t, err, "Should find final session")

	assert.Equal(t, 4, finalSession.CurrentStep, "Should be at step 4")
	assert.Equal(t, 14, finalSession.Attributes.ST, "Attributes should be saved")
	assert.NotZero(t, finalSession.DerivedValues.PA, "Derived values should be saved")
	assert.Len(t, finalSession.Skills, 1, "Skills should be saved")
}

func TestCharacterCreationSession_EdgeCases(t *testing.T) {
	setupCharacterCreationTestDB(t)

	// Test session with minimal data
	minimalSession := &CharacterCreationSession{
		ID:          "minimal-test",
		UserID:      1,
		Name:        "Minimal Character",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		CurrentStep: 1,
	}

	err := database.DB.Create(minimalSession).Error
	assert.NoError(t, err, "Should create minimal session")

	// Test session with special characters
	specialSession := createTestCharacterCreationSession("special-test")
	specialSession.Name = "Ärger McÇharacter ñ"
	specialSession.Herkunft = "Ålba with spëcial chars"

	err = database.DB.Create(specialSession).Error
	assert.NoError(t, err, "Should create session with special characters")

	// Test session with extreme values
	extremeSession := createTestCharacterCreationSession("extreme-test")
	extremeSession.Attributes.ST = 99
	extremeSession.DerivedValues.LPMax = 999
	extremeSession.CurrentStep = 10

	err = database.DB.Create(extremeSession).Error
	assert.NoError(t, err, "Should create session with extreme values")
}
