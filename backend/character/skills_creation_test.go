package character

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"bamort/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetAllSkillsWithLearningCosts tests the helper function directly
func TestGetAllSkillsWithLearningCosts(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	tests := []struct {
		name           string
		characterClass string
		expectError    bool
		expectSkills   bool
	}{
		{
			name:           "ValidCharacterClassAs",
			characterClass: "As",
			expectError:    false,
			expectSkills:   true,
		},
		{
			name:           "ValidCharacterClassMagier",
			characterClass: "Magier",
			expectError:    false,
			expectSkills:   true,
		},
		{
			name:           "EmptyCharacterClass",
			characterClass: "",
			expectError:    false,
			expectSkills:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skillsByCategory, err := GetAllSkillsWithLearningCosts(tt.characterClass)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err, "GetAllSkillsWithLearningCosts should not return error")

			if tt.expectSkills {
				assert.Greater(t, len(skillsByCategory), 0, "Should return at least some skill categories")

				// Verify structure of returned data
				for categoryName, skills := range skillsByCategory {
					assert.NotEmpty(t, categoryName, "Category name should not be empty")

					skillsArray := skills
					assert.Greater(t, len(skillsArray), 0, "Category should have at least one skill")

					// Check first skill structure
					if len(skillsArray) > 0 {
						skill := skillsArray[0]
						assert.Contains(t, skill, "name", "Skill should have name")
						assert.Contains(t, skill, "learnCost", "Skill should have learnCost")

						name, nameOk := skill["name"].(string)
						assert.True(t, nameOk, "Skill name should be string")
						assert.NotEmpty(t, name, "Skill name should not be empty")

						learnCost, costOk := skill["learnCost"].(int)
						assert.True(t, costOk, "Learn cost should be int")
						assert.Greater(t, learnCost, 0, "Learn cost should be positive")
						assert.LessOrEqual(t, learnCost, 500, "Learn cost should be reasonable")
					}
				}

				t.Logf("Character class %s: Found %d skill categories", tt.characterClass, len(skillsByCategory))

				// Log some sample skills for verification
				count := 0
				for categoryName, skills := range skillsByCategory {
					skillsArray := skills
					for _, skill := range skillsArray {
						skillMap := skill
						t.Logf("  %s -> %s: %v LP", categoryName, skillMap["name"], skillMap["learnCost"])
						count++
						if count >= 5 { // Only log first 5 skills
							break
						}
					}
					if count >= 5 {
						break
					}
				}
			}
		})
	}
}

// TestGetAvailableSkillsForCreationHandler tests the HTTP handler directly
func TestHandlerGetAvailableSkillsForCreation(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	tests := []struct {
		name         string
		requestBody  interface{}
		expectStatus int
		expectError  bool
		validateFunc func(t *testing.T, response map[string]interface{})
	}{
		{
			name: "ValidMagierRequest",
			requestBody: gin.H{
				"characterClass": "Magier",
			},
			expectStatus: http.StatusOK,
			expectError:  false,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "skills_by_category")
				skillsByCategory := response["skills_by_category"].(map[string]interface{})
				assert.Greater(t, len(skillsByCategory), 0, "Should have skill categories")
			},
		},
		{
			name: "ValidAsRequest",
			requestBody: gin.H{
				"characterClass": "As",
			},
			expectStatus: http.StatusOK,
			expectError:  false,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "skills_by_category")
				skillsByCategory := response["skills_by_category"].(map[string]interface{})
				assert.Greater(t, len(skillsByCategory), 0, "Should have skill categories")

				// Verify that skills have reasonable learning costs
				hasReasonableCosts := false
				for _, skills := range skillsByCategory {
					if skillsList, ok := skills.([]interface{}); ok {
						for _, skill := range skillsList {
							if skillMap, ok := skill.(map[string]interface{}); ok {
								if learnCost, exists := skillMap["leCost"]; exists {
									if cost, ok := learnCost.(float64); ok {
										assert.Greater(t, cost, 0.0, "Learn cost should be positive")
										assert.LessOrEqual(t, cost, 500.0, "Learn cost should be reasonable")
										hasReasonableCosts = true
									}
								}
							}
						}
					}
				}
				assert.True(t, hasReasonableCosts, "Should find skills with reasonable costs")
			},
		},
		{
			name: "EmptyCharacterClass",
			requestBody: gin.H{
				"characterClass": "",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  true,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "MissingCharacterClass",
			requestBody: gin.H{
				"someOtherField": "value",
			},
			expectStatus: http.StatusBadRequest,
			expectError:  true,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "error")
			},
		},
		{
			name:         "InvalidJSON",
			requestBody:  "invalid json string",
			expectStatus: http.StatusBadRequest,
			expectError:  true,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				assert.Contains(t, response, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request body
			var requestBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				requestBody = []byte(str)
			} else {
				requestBody, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			// Create HTTP request
			req, err := http.NewRequest("POST", "/api/characters/available-skills-creation", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the handler
			GetAvailableSkillsForCreation(c)

			// Verify response status
			assert.Equal(t, tt.expectStatus, w.Code, "HTTP status should match expected")

			// Parse response body
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "Response should be valid JSON")

			// Run custom validation if provided
			if tt.validateFunc != nil {
				tt.validateFunc(t, response)
			}

			// Log response for debugging
			if !tt.expectError {
				t.Logf("Response for %s: %v", tt.name, response)
			}
		})
	}
}

// BenchmarkGetAllSkillsWithLearningCosts benchmarks the skills loading function
func BenchmarkGetAllSkillsWithLearningCosts(b *testing.B) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	characterClass := "Magier"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GetAllSkillsWithLearningCosts(characterClass)
		if err != nil {
			b.Fatalf("GetAllSkillsWithLearningCosts failed: %v", err)
		}
	}
}

// TestSkillsCreationEndpointIntegration tests the full integration
func TestSkillsCreationEndpointIntegration(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	// Test different character classes
	characterClasses := []string{"As", "Magier", "Krieger", "Spitzbube"}

	for _, class := range characterClasses {
		t.Run("CharacterClass_"+class, func(t *testing.T) {
			// Create request
			requestData := gin.H{
				"characterClass": class,
			}
			requestBody, err := json.Marshal(requestData)
			assert.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest("POST", "/api/characters/available-skills-creation", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the handler
			GetAvailableSkillsForCreation(c)

			// Verify response
			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Verify structure
			assert.Contains(t, response, "skills_by_category")
			skillsByCategory := response["skills_by_category"].(map[string]interface{})
			assert.Greater(t, len(skillsByCategory), 0, "Should have at least one skill category")

			// Count total skills available
			totalSkills := 0
			for categoryName, skills := range skillsByCategory {
				if skillsList, ok := skills.([]interface{}); ok {
					totalSkills += len(skillsList)
					t.Logf("Category %s has %d skills", categoryName, len(skillsList))
				}
			}

			assert.Greater(t, totalSkills, 0, "Should have at least some skills available")
			t.Logf("Character class %s: %d total skills in %d categories", class, totalSkills, len(skillsByCategory))
		})
	}
}
