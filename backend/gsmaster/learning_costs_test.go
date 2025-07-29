package gsmaster

//Diese Tests hier sind SCHROTT denn sie testen statisch erzeugt Strukturen und nicht die Abfrage aus der DB wie erhofft
/*
// Test for exported GetAvailableSkillCategories function
func TestGetAvailableSkillCategories(t *testing.T) {
	testCases := []struct {
		skillName     string
		expectedCount int
		description   string
		checkFirst    bool
		firstCategory string
	}{
		{"Menschenkenntnis", 2, "Menschenkenntnis should have two categories (Sozial, Unterwelt)", true, "Sozial"},
		{"Stichwaffen", 1, "Stichwaffen should have one category", true, "Waffen"},
		{"Geländelauf", 1, "Geländelauf should have one category", true, "Körper"},
		{"NonExistentSkill", 1, "Unknown skill should have default category", true, "Alltag"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := GetAvailableSkillCategories(tc.skillName)
			assert.Equal(t, tc.expectedCount, len(result), tc.description)

			if tc.checkFirst && len(result) > 0 {
				assert.Equal(t, tc.firstCategory, result[0].Category,
					fmt.Sprintf("First category for %s should be %s", tc.skillName, tc.firstCategory))
			}
		})
	}
}

// Test for exported GetDefaultCategory function
func TestGetDefaultCategory(t *testing.T) {
	testCases := []struct {
		skillName        string
		expectedCategory string
		description      string
	}{
		{"Menschenkenntnis", "Sozial", "Should return Sozial for Menschenkenntnis"},
		{"Stichwaffen", "Waffen", "Should return Waffen for Stichwaffen"},
		{"Geländelauf", "Körper", "Should return Körper for Geländelauf"},
		{"Sprache", "Alltag", "Should return Alltag for Sprache"},
		{"NonExistentSkill", "Alltag", "Should return default Alltag for unknown skill"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := GetDefaultCategory(tc.skillName)
			assert.Equal(t, tc.expectedCategory, result, tc.description)
		})
	}
}

// Test for exported GetDefaultDifficulty function
func TestGetDefaultDifficulty(t *testing.T) {
	testCases := []struct {
		skillName          string
		expectedDifficulty string
		description        string
	}{
		{"Menschenkenntnis", "schwer", "Should return schwer for Menschenkenntnis"},
		{"Stichwaffen", "leicht", "Should return leicht for Stichwaffen"},
		{"Geländelauf", "leicht", "Should return leicht for Geländelauf"},
		{"Sprache", "normal", "Should return normal for Sprache"},
		{"NonExistentSkill", "normal", "Should return default normal for unknown skill"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result := GetDefaultDifficulty(tc.skillName)
			assert.Equal(t, tc.expectedDifficulty, result, tc.description)
		})
	}
}

func TestCalculateDetailedSkillLearningCostForHexer(t *testing.T) {
	// Testfall, der direkt mit den exportierten Funktionen arbeitet
	t.Run("Lernkosten für Menschenkenntnis als Hexer", func(t *testing.T) {
		// Verwendung der exportierten Funktion CalculateDetailedSkillLearningCost
		result, err := CalculateDetailedSkillLearningCost("Menschenkenntnis", "Hexer")
		assert.NoError(t, err, "CalculateDetailedSkillLearningCost sollte keinen Fehler zurückgeben")
		assert.NotNil(t, result, "Ergebnis sollte nicht nil sein")

		fmt.Printf("Lernkosten für Menschenkenntnis als Hexer:\n")
		fmt.Printf("Lerneinheiten (LE): %d\n", result.LE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", result.Ep)
		fmt.Printf("Geldkosten (GS): %d\n", result.Money)

		// Test der erwarteten Werte basierend auf aktueller Implementierung
		// Menschenkenntnis ist Sozial/schwer: 4 LE
		// Hexer EP-Kosten für Sozial: 20 EP/TE
		// Money-Kosten: 20 GS/LE
		assert.Equal(t, 4, result.LE, "LE-Kosten für Menschenkenntnis (Sozial/schwer) sollten 4 sein")
		assert.Equal(t, 80, result.Ep, "EP-Kosten sollten 80 sein (4 LE * 20 EP)")
		assert.Equal(t, 80, result.Money, "Geldkosten sollten 80 GS sein (4 LE * 20 GS)")
	})

	t.Run("Verbesserungskosten für Menschenkenntnis als Hexer", func(t *testing.T) {
		// Verwendung der exportierten Funktion CalculateDetailedSkillImprovementCost
		currentLevel := 10
		result, err := CalculateDetailedSkillImprovementCost("Menschenkenntnis", "Hexer", currentLevel)
		assert.NoError(t, err, "CalculateDetailedSkillImprovementCost sollte keinen Fehler zurückgeben")
		assert.NotNil(t, result, "Ergebnis sollte nicht nil sein")

		fmt.Printf("\nVerbesserungskosten für Menschenkenntnis als Hexer (von %d auf %d):\n", currentLevel, currentLevel+1)
		fmt.Printf("Zielstufe: %d\n", result.Stufe) // This appears to be the target level
		fmt.Printf("Erfahrungspunkte (EP): %d\n", result.Ep)
		fmt.Printf("Geldkosten (GS): %d\n", result.Money)

		// Test der erwarteten Werte basierend auf aktueller Implementierung
		// Result.Stufe appears to be target level (11), not training units needed
		assert.Equal(t, 11, result.Stufe, "Zielstufe sollte 11 sein (von 10 auf 11)")
		assert.Equal(t, 80, result.Ep, "EP-Kosten sollten 80 sein")
		assert.Equal(t, 80, result.Money, "Geldkosten sollten 80 GS sein")
	})
}
*/
