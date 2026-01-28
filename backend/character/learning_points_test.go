package character

import (
	"bamort/database"
	"bamort/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetCharacterClassLearningPoints(t *testing.T) {
	// Setup test database
	database.SetupTestDB()

	// Migrate the new structures
	if err := models.MigrateStructure(database.DB); err != nil {
		t.Fatalf("Failed to migrate structures: %v", err)
	}

	/*
		// Populate test data (character classes and learning points)
		if err := models.PopulateClassLearningPointsData(); err != nil {
			t.Logf("Warning: Failed to populate learning points data: %v", err)
			// Continue anyway - some tests may still work
		}
	*/

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/characters/classes/learning-points", GetCharacterClassLearningPoints)

	tests := []struct {
		name            string
		classParam      string
		standParam      string
		expectedStatus  int
		expectError     bool
		expectedClass   string
		checkWeapons    bool
		expectedWeapons int
		checkSpells     bool
		expectedSpells  int
	}{
		{
			name:            "Valid Spitzbube class Mittelschicht stand",
			classParam:      "Spitzbube",
			standParam:      "Mittelschicht",
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedClass:   "Spitzbube",
			checkWeapons:    true,
			expectedWeapons: 20,
			checkSpells:     true,
			expectedSpells:  0,
		},
		{
			name:            "Valid Hexer class without stand",
			classParam:      "Hexer",
			standParam:      "",
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedClass:   "Hexer",
			checkWeapons:    true,
			expectedWeapons: 2,
			checkSpells:     true,
			expectedSpells:  6,
		},
		{
			name:            "Valid Hexer class with Volk stand",
			classParam:      "Hexer",
			standParam:      "Volk",
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedClass:   "Hexer",
			checkWeapons:    true,
			expectedWeapons: 2,
			checkSpells:     true,
			expectedSpells:  6,
		},
		{
			name:            "Valid Krieger class with Adel stand",
			classParam:      "Krieger",
			standParam:      "Adel",
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedClass:   "Krieger",
			checkWeapons:    true,
			expectedWeapons: 36,
			checkSpells:     false,
		},
		{
			name:            "Valid Magier class",
			classParam:      "Magier",
			standParam:      "",
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedClass:   "Magier",
			checkWeapons:    true,
			expectedWeapons: 2,
			checkSpells:     true,
			expectedSpells:  7,
		},
		{
			name:            "Valid Spitzbube class",
			classParam:      "Spitzbube",
			standParam:      "",
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedClass:   "Spitzbube",
			checkWeapons:    true,
			expectedWeapons: 20,
			checkSpells:     false,
		},
		{
			name:            "Valid Waldläufer class",
			classParam:      "Waldläufer",
			standParam:      "",
			expectedStatus:  http.StatusOK,
			expectError:     false,
			expectedClass:   "Waldläufer",
			checkWeapons:    true,
			expectedWeapons: 20,
			checkSpells:     false,
		},
		{
			name:           "Invalid class should return error",
			classParam:     "InvalidClass",
			standParam:     "",
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:           "Missing class parameter should return error",
			classParam:     "",
			standParam:     "",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Valid class with invalid stand should still work",
			classParam:     "Hexer",
			standParam:     "InvalidStand",
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedClass:  "Hexer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build request URL
			url := "/api/characters/classes/learning-points"
			if tt.classParam != "" || tt.standParam != "" {
				url += "?"
				if tt.classParam != "" {
					url += "class=" + tt.classParam
				}
				if tt.standParam != "" {
					if tt.classParam != "" {
						url += "&"
					}
					url += "stand=" + tt.standParam
				}
			}

			// Create request
			req, err := http.NewRequest("GET", url, nil)
			assert.NoError(t, err)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectError {
				// For error cases, check that we have an error message
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			} else {
				// For success cases, check the response structure
				var response LearningPointsData
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check basic fields
				assert.Equal(t, tt.expectedClass, response.ClassName)
				assert.NotEmpty(t, response.ClassCode)
				assert.NotNil(t, response.LearningPoints)
				assert.NotNil(t, response.TypicalSkills)

				// Check weapon points if specified (now in LearningPoints["Waffen"])
				if tt.checkWeapons {
					assert.Equal(t, tt.expectedWeapons, response.LearningPoints["Waffen"], "Weapon learning points should match")
				}

				// Check spell points if specified
				if tt.checkSpells {
					assert.Equal(t, tt.expectedSpells, response.SpellPoints)
				}

				// Check that learning points are not empty
				assert.NotEmpty(t, response.LearningPoints)

				// Check that typical skills are not empty
				assert.NotEmpty(t, response.TypicalSkills)

			}
		})
	}
}

func TestGetLearningPointsForClass(t *testing.T) {
	// Setup test database
	database.SetupTestDB()

	// Migrate the new structures
	if err := models.MigrateStructure(database.DB); err != nil {
		t.Fatalf("Failed to migrate structures: %v", err)
	}

	/*
		// Populate test data
		if err := models.PopulateClassLearningPointsData(); err != nil {
			t.Logf("Warning: Failed to populate learning points data: %v", err)
		}
	*/

	tests := []struct {
		name          string
		className     string
		stand         string
		expectError   bool
		expectedClass string
		expectedCode  string
		checkPoints   map[string]int
		checkStand    map[string]int
	}{
		{
			name:          "Hexer class data",
			className:     "Hexer",
			stand:         "",
			expectError:   false,
			expectedClass: "Hexer",
			expectedCode:  "Hx",
			checkPoints: map[string]int{
				"Alltag": 3,
				"Sozial": 2,
				"Wissen": 2,
			},
		},
		{
			name:          "Hexer with Volk stand",
			className:     "Hexer",
			stand:         "Volk",
			expectError:   false,
			expectedClass: "Hexer",
			expectedCode:  "Hx",
			checkPoints: map[string]int{
				"Alltag": 5, // Base 3 + Volk bonus 2 = 5
				"Sozial": 2,
				"Wissen": 2,
			},
			checkStand: map[string]int{
				"Alltag": 2,
			},
		},
		{
			name:          "Krieger class data",
			className:     "Krieger",
			stand:         "",
			expectError:   false,
			expectedClass: "Krieger",
			expectedCode:  "Kr",
			checkPoints: map[string]int{
				"Alltag": 2,
				"Kampf":  3,
				"Körper": 1,
			},
		},
		{
			name:          "Krieger with Adel stand",
			className:     "Krieger",
			stand:         "Adel",
			expectError:   false,
			expectedClass: "Krieger",
			expectedCode:  "Kr",
			checkPoints: map[string]int{
				"Alltag": 2,
				"Kampf":  3,
				"Körper": 1,
				"Sozial": 2, // Stand bonus adds this new category
			},
			checkStand: map[string]int{
				"Sozial": 2,
			},
		},
		{
			name:        "Invalid class should return error",
			className:   "InvalidClass",
			stand:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := getLearningPointsForClass(tt.className, tt.stand)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, data)

				// Check basic properties
				assert.Equal(t, tt.expectedClass, data.ClassName)
				assert.Equal(t, tt.expectedCode, data.ClassCode)

				// Check learning points
				if tt.checkPoints != nil {
					for category, expectedPoints := range tt.checkPoints {
						actualPoints, exists := data.LearningPoints[category]
						assert.True(t, exists, "Category %s should exist", category)
						assert.Equal(t, expectedPoints, actualPoints, "Points for category %s", category)
					}
				}

				// Check that we have some typical skills
				assert.NotEmpty(t, data.TypicalSkills)

				// Validate typical skills structure
				for _, skill := range data.TypicalSkills {
					assert.NotEmpty(t, skill.Name)
					assert.NotEmpty(t, skill.Attribute)
					assert.GreaterOrEqual(t, skill.Bonus, 0)
				}
			}
		})
	}
}

func TestGetStandBonusPoints(t *testing.T) {
	database.SetupTestDB()
	defer database.ResetTestDB()

	tests := []struct {
		name     string
		stand    string
		expected map[string]int
	}{
		{
			name:     "Unfreie stand",
			stand:    "Unfreie",
			expected: map[string]int{"Halbwelt": 2},
		},
		{
			name:     "Volk stand",
			stand:    "Volk",
			expected: map[string]int{"Alltag": 2},
		},
		{
			name:     "Mittelschicht stand",
			stand:    "Mittelschicht",
			expected: map[string]int{"Wissen": 2},
		},
		{
			name:     "Adel stand",
			stand:    "Adel",
			expected: map[string]int{"Sozial": 2},
		},
		{
			name:     "Invalid stand",
			stand:    "Invalid",
			expected: map[string]int{},
		},
		{
			name:     "Empty stand",
			stand:    "",
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStandBonusPoints(tt.stand)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test all character classes to ensure they're properly defined
func TestAllCharacterClassesAreDefined(t *testing.T) {
	// Setup test database
	database.SetupTestDB()

	// Migrate the new structures
	if err := models.MigrateStructure(database.DB); err != nil {
		t.Fatalf("Failed to migrate structures: %v", err)
	}

	/*
		// Populate test data
		if err := models.PopulateClassLearningPointsData(); err != nil {
			t.Logf("Warning: Failed to populate learning points data: %v", err)
		}
	*/

	expectedClasses := []string{
		"Assassine", "Barbar", "Glücksritter", "Händler", "Krieger", "Spitzbube", "Waldläufer",
		"Barde", "Ordenskrieger", "Druide", "Hexer", "Magier", "Priester Beschützer", "Priester Streiter", "Schamane",
	}

	for _, className := range expectedClasses {
		t.Run("Class_"+className, func(t *testing.T) {
			data, err := getLearningPointsForClass(className, "")
			assert.NoError(t, err, "Class %s should be defined", className)
			assert.NotNil(t, data, "Class %s should return data", className)
			assert.Equal(t, className, data.ClassName)
			assert.NotEmpty(t, data.ClassCode)
			assert.NotEmpty(t, data.LearningPoints)
			assert.GreaterOrEqual(t, data.WeaponPoints, 0)
		})
	}
}
