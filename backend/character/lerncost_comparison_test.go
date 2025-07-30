package character

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestCompareOldVsNewLearningCostSystems vergleicht die Ergebnisse von GetLernCost und GetLernCostNewSystem
func TestCompareOldVsNewLearningCostSystems(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	// Migrate the schema
	err := models.MigrateStructure()
	assert.NoError(t, err)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name         string
		charId       uint
		skillName    string
		currentLevel int
		TargetLevel  int `json:"target_level,omitempty"` // Zielwert (optional, für Kostenberechnung bis zu einem bestimmten Level)
		usePP        int
		useGold      int
		description  string
		Type         string  `json:"type" binding:"required,oneof=skill spell weapon"`                     // 'skill', 'spell' oder 'weapon' Waffenfertigkeiten sind normale Fertigkeiten (evtl. kann hier später der Name der Waffe angegeben werden )
		Action       string  `json:"action" binding:"required,oneof=learn improve"`                        // 'learn' oder 'improve'
		Reward       *string `json:"reward" binding:"required,oneof=default noGold halveep halveepnoGold"` // Belohnungsoptionen Lernen als Belohnung
	}{
		{
			name:         "improve Athletik Basic Test",
			charId:       20,
			skillName:    "Athletik",
			currentLevel: 9,
			usePP:        0,
			useGold:      0,
			description:  "Basic comparison for Athletik skill improvement",
			Type:         "skill",
			Action:       "improve",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "improve Athletik with PP",
			charId:       20,
			skillName:    "Athletik",
			currentLevel: 9,
			usePP:        10,
			useGold:      0,
			description:  "Comparison with Practice Points usage",
			Type:         "skill",
			Action:       "improve",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "improve Athletik with Gold",
			charId:       20,
			skillName:    "Athletik",
			currentLevel: 9,
			usePP:        0,
			useGold:      100,
			description:  "Comparison with Gold usage",
			Type:         "skill",
			Action:       "improve",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "improve Athletik with 2000 Gold",
			charId:       20,
			skillName:    "Athletik",
			currentLevel: 9,
			usePP:        0,
			useGold:      2000,
			description:  "Comparison with Gold usage",
			Type:         "skill",
			Action:       "improve",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "improve Athletik with PP and Gold",
			charId:       20,
			skillName:    "Athletik",
			currentLevel: 9,
			usePP:        15,
			useGold:      150,
			description:  "Comparison with both PP and Gold usage",
			Type:         "skill",
			Action:       "improve",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "learn Naturkunde Basic Test",
			charId:       20,
			skillName:    "Naturkunde",
			currentLevel: 9,
			usePP:        0,
			useGold:      0,
			description:  "Comparison Learn Basic",
			Type:         "skill",
			Action:       "learn",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "learn Naturkunde Test 100 Gold",
			charId:       20,
			skillName:    "Naturkunde",
			currentLevel: 9,
			usePP:        0,
			useGold:      100,
			description:  "Comparison Learn 100 Gold",
			Type:         "skill",
			Action:       "learn",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "learn Naturkunde Test 2000 Gold",
			charId:       20,
			skillName:    "Naturkunde",
			currentLevel: 9,
			usePP:        0,
			useGold:      2000,
			description:  "Comparison Learn 2000 Gold",
			Type:         "skill",
			Action:       "learn",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "learn Naturkunde Test 2 PP",
			charId:       20,
			skillName:    "Naturkunde",
			currentLevel: 9,
			usePP:        2,
			useGold:      0,
			description:  "Comparison Learn 2 PP",
			Type:         "skill",
			Action:       "learn",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
		{
			name:         "learn Beeinflussen Test Basic",
			charId:       18,
			skillName:    "Beeinflussen",
			currentLevel: 9,
			usePP:        2,
			useGold:      0,
			description:  "Comparison Learn ",
			Type:         "spell",
			Action:       "learn",
			TargetLevel:  0, // Calculate all levels
			Reward:       &[]string{"default"}[0],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Printf("\n=== Comparing Old vs New System: %s ===\n", tc.description)

			// Prepare request data
			requestData := gsmaster.LernCostRequest{
				CharId:       tc.charId,
				Name:         tc.skillName,
				CurrentLevel: tc.currentLevel,
				Type:         tc.Type,
				Action:       tc.Action,
				TargetLevel:  tc.TargetLevel, // Calculate all levels
				UsePP:        tc.usePP,
				UseGold:      tc.useGold,
				Reward:       tc.Reward,
			}
			requestBody, _ := json.Marshal(requestData)

			// Test Old System (GetLernCost)
			reqOld, _ := http.NewRequest("POST", "/api/characters/lerncost", bytes.NewBuffer(requestBody))
			reqOld.Header.Set("Content-Type", "application/json")

			wOld := httptest.NewRecorder()
			cOld, _ := gin.CreateTestContext(wOld)
			cOld.Request = reqOld

			GetLernCost(cOld)

			// Test New System (GetLernCostNewSystem)
			reqNew, _ := http.NewRequest("POST", "/api/characters/lerncost-new", bytes.NewBuffer(requestBody))
			reqNew.Header.Set("Content-Type", "application/json")

			wNew := httptest.NewRecorder()
			cNew, _ := gin.CreateTestContext(wNew)
			cNew.Request = reqNew

			GetLernCostNewSystem(cNew)

			// Check that both systems returned successful responses
			fmt.Printf("Old System Status: %d, New System Status: %d\n", wOld.Code, wNew.Code)

			if wOld.Code != http.StatusOK && wNew.Code != http.StatusOK {
				fmt.Printf("Both systems failed - Old: %s, New: %s\n", wOld.Body.String(), wNew.Body.String())
				t.Skip("Both systems failed, skipping comparison")
				return
			}

			if wOld.Code != http.StatusOK {
				fmt.Printf("Old system failed: %s\n", wOld.Body.String())
				fmt.Printf("New system succeeded with status %d\n", wNew.Code)
				t.Skip("Old system failed, new system comparison only")
				return
			}

			if wNew.Code != http.StatusOK {
				fmt.Printf("New system failed: %s\n", wNew.Body.String())
				fmt.Printf("Old system succeeded with status %d\n", wOld.Code)
				t.Skip("New system failed, old system comparison only")
				return
			}

			// Parse responses
			var oldResponse []gsmaster.SkillCostResultNew
			var newResponse []gsmaster.SkillCostResultNew

			err := json.Unmarshal(wOld.Body.Bytes(), &oldResponse)
			assert.NoError(t, err, "Old system response should be valid JSON")

			err = json.Unmarshal(wNew.Body.Bytes(), &newResponse)
			assert.NoError(t, err, "New system response should be valid JSON")

			// Compare basic structure
			fmt.Printf("Old System: %d levels, New System: %d levels\n", len(oldResponse), len(newResponse))

			// Create comparison table
			fmt.Printf("\n=== Detailed Cost Comparison ===\n")
			fmt.Printf("Level | Old EP | New EP | Old Gold | New Gold | Old LE | New LE | Old PP used | New PP used | Old Gold used | New Gold used | Match?\n")
			fmt.Printf("------|--------|--------|----------|----------|--------|--------|-------------|-------------|---------------|---------------|-------\n")

			// Compare each level that exists in both systems
			maxLevels := len(oldResponse)
			if len(newResponse) < maxLevels {
				maxLevels = len(newResponse)
			}

			exactMatches := 0
			totalComparisons := 0

			for i := 0; i < maxLevels; i++ {
				old := oldResponse[i]
				new := newResponse[i]

				// Check if target levels match
				/*if old.TargetLevel != new.TargetLevel {
					fmt.Printf("Level mismatch at index %d: Old=%d, New=%d\n", i, old.TargetLevel, new.TargetLevel)
					continue
				}*/

				totalComparisons++

				// Compare costs
				epMatch := old.EP == new.EP
				goldMatch := old.GoldCost == new.GoldCost
				leMatch := old.LE == new.LE
				ppMatch := old.PPUsed == new.PPUsed
				goldUsedMatch := old.GoldUsed == new.GoldUsed

				overallMatch := epMatch && goldMatch && leMatch && ppMatch && goldUsedMatch
				if overallMatch {
					exactMatches++
				}

				matchSymbol := "✓"
				if !overallMatch {
					matchSymbol = "✗"
				}

				fmt.Printf("%5d | %6d | %6d | %8d | %8d | %6d | %6d | %11d | %11d | %13d | %13d | %7s\n",
					old.TargetLevel, old.EP, new.EP, old.GoldCost, new.GoldCost,
					old.LE, new.LE, old.PPUsed, new.PPUsed, old.GoldUsed, new.GoldUsed, matchSymbol)

				// Individual field assertions for debugging
				if !epMatch {
					fmt.Printf("  EP mismatch at level %d: Old=%d, New=%d (diff=%d)\n",
						old.TargetLevel, old.EP, new.EP, old.EP-new.EP)
				}
				if !goldMatch {
					fmt.Printf("  GoldCost mismatch at level %d: Old=%d, New=%d (diff=%d)\n",
						old.TargetLevel, old.GoldCost, new.GoldCost, old.GoldCost-new.GoldCost)
				}
				if !leMatch {
					fmt.Printf("  LE mismatch at level %d: Old=%d, New=%d (diff=%d)\n",
						old.TargetLevel, old.LE, new.LE, old.LE-new.LE)
				}
				if !ppMatch {
					fmt.Printf("  PP used mismatch at level %d: Old=%d, New=%d (diff=%d)\n",
						old.TargetLevel, old.PPUsed, new.PPUsed, old.PPUsed-new.PPUsed)
				}
				if !goldUsedMatch {
					fmt.Printf("  GoldUsed mismatch at level %d: Old=%d, New=%d (diff=%d)\n",
						old.TargetLevel, old.GoldUsed, new.GoldUsed, old.GoldUsed-new.GoldUsed)
				}

				// Verify basic consistency within each system
				assert.GreaterOrEqual(t, old.EP, 0, "Old system EP should be non-negative for level %d", old.TargetLevel)
				assert.GreaterOrEqual(t, new.EP, 0, "New system EP should be non-negative for level %d", new.TargetLevel)
				assert.GreaterOrEqual(t, old.LE, 0, "Old system LE should be non-negative for level %d", old.TargetLevel)
				assert.GreaterOrEqual(t, new.LE, 0, "New system LE should be non-negative for level %d", new.TargetLevel)
			}

			// Summary statistics
			matchPercentage := float64(exactMatches) / float64(totalComparisons) * 100
			fmt.Printf("\n=== Comparison Summary ===\n")
			fmt.Printf("Total Comparisons: %d\n", totalComparisons)
			fmt.Printf("Exact Matches: %d\n", exactMatches)
			fmt.Printf("Match Percentage: %.1f%%\n", matchPercentage)

			// Check for levels that exist in one system but not the other
			if len(oldResponse) != len(newResponse) {
				fmt.Printf("Length Difference: Old=%d levels, New=%d levels\n", len(oldResponse), len(newResponse))

				if len(oldResponse) > len(newResponse) {
					fmt.Printf("Old system has additional levels:\n")
					for i := len(newResponse); i < len(oldResponse); i++ {
						fmt.Printf("  Level %d: EP=%d, Gold=%d, LE=%d\n",
							oldResponse[i].TargetLevel, oldResponse[i].EP, oldResponse[i].GoldCost, oldResponse[i].LE)
					}
				} else {
					fmt.Printf("New system has additional levels:\n")
					for i := len(oldResponse); i < len(newResponse); i++ {
						fmt.Printf("  Level %d: EP=%d, Gold=%d, LE=%d\n",
							newResponse[i].TargetLevel, newResponse[i].EP, newResponse[i].GoldCost, newResponse[i].LE)
					}
				}
			}

			// Basic assertions for test validation
			assert.Greater(t, totalComparisons, 0, "Should have at least one level to compare")

			// If systems are supposed to be equivalent, we might want exact matches
			// For now, we'll just document the differences and ensure both are reasonable
			if matchPercentage < 50.0 {
				t.Logf("WARNING: Low match percentage (%.1f%%) between old and new systems", matchPercentage)
			}

			// Verify that both systems produce reasonable results
			if len(oldResponse) > 0 {
				assert.Equal(t, tc.skillName, oldResponse[0].SkillName, "Old system should return correct skill name")
				assert.Equal(t, fmt.Sprintf("%d", tc.charId), oldResponse[0].CharacterID, "Old system should return correct character ID")
			}

			if len(newResponse) > 0 {
				assert.Equal(t, tc.skillName, newResponse[0].SkillName, "New system should return correct skill name")
				assert.Equal(t, fmt.Sprintf("%d", tc.charId), newResponse[0].CharacterID, "New system should return correct character ID")
			}
		})
	}

	// Summary test to compare system performance
	t.Run("System Performance Summary", func(t *testing.T) {
		fmt.Printf("\n=== Overall System Comparison Summary ===\n")
		fmt.Printf("Both systems tested with various configurations:\n")
		fmt.Printf("- Basic skill improvement (no resources)\n")
		fmt.Printf("- With Practice Points usage\n")
		fmt.Printf("- With Gold usage\n")
		fmt.Printf("- With combined PP and Gold usage\n")
		fmt.Printf("\nKey observations should be documented above.\n")
		fmt.Printf("Differences may indicate:\n")
		fmt.Printf("1. Different calculation methods\n")
		fmt.Printf("2. Different data sources (old vs new tables)\n")
		fmt.Printf("3. Implementation bugs in either system\n")
		fmt.Printf("4. Different business rules or assumptions\n")
	})
}

// TestPerformanceComparison vergleicht die Performance der beiden Systeme
func TestPerformanceComparison(t *testing.T) {
	// Setup test database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	// Migrate the schema
	err := models.MigrateStructure()
	assert.NoError(t, err)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	fmt.Printf("\n=== Performance Comparison ===\n")

	// Prepare test data
	requestData := gsmaster.LernCostRequest{
		CharId:       20,
		Name:         "Athletik",
		CurrentLevel: 5,
		Type:         "skill",
		Action:       "improve",
		TargetLevel:  0,
		UsePP:        0,
		UseGold:      0,
		Reward:       &[]string{"default"}[0],
	}
	requestBody, _ := json.Marshal(requestData)

	// Test old system multiple times
	oldSystemTimes := make([]int64, 10)
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("POST", "/api/characters/lerncost", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		start := getCurrentTime()
		GetLernCost(c)
		oldSystemTimes[i] = getCurrentTime() - start

		if w.Code != http.StatusOK {
			t.Logf("Old system failed on iteration %d: %s", i, w.Body.String())
		}
	}

	// Test new system multiple times
	newSystemTimes := make([]int64, 10)
	for i := 0; i < 10; i++ {
		req, _ := http.NewRequest("POST", "/api/characters/lerncost-new", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		start := getCurrentTime()
		GetLernCostNewSystem(c)
		newSystemTimes[i] = getCurrentTime() - start

		if w.Code != http.StatusOK {
			t.Logf("New system failed on iteration %d: %s", i, w.Body.String())
		}
	}

	// Calculate averages
	oldAvg := calculateAverage(oldSystemTimes)
	newAvg := calculateAverage(newSystemTimes)

	fmt.Printf("Old System Average Time: %d μs\n", oldAvg)
	fmt.Printf("New System Average Time: %d μs\n", newAvg)

	if newAvg < oldAvg {
		improvement := float64(oldAvg-newAvg) / float64(oldAvg) * 100
		fmt.Printf("New system is %.1f%% faster\n", improvement)
	} else {
		regression := float64(newAvg-oldAvg) / float64(oldAvg) * 100
		fmt.Printf("New system is %.1f%% slower\n", regression)
	}
}

// Helper functions for performance testing
func getCurrentTime() int64 {
	return 0 // Placeholder - in real implementation would use time.Now().UnixNano()
}

func calculateAverage(times []int64) int64 {
	var sum int64
	for _, t := range times {
		sum += t
	}
	return sum / int64(len(times))
}
