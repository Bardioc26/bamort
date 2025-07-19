package gsmaster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLerntabellenForHexer(t *testing.T) {
	// Testfall für Menschenkenntnis als Hexer - Lernkosten
	t.Run("Lernkosten für Menschenkenntnis als Hexer", func(t *testing.T) {
		// Verwende die exportierte Funktion CalculateDetailedSkillLearningCost
		result, err := CalculateDetailedSkillLearningCost("Menschenkenntnis", "Hexer")
		assert.NoError(t, err, "CalculateDetailedSkillLearningCost sollte keinen Fehler zurückgeben")
		assert.NotNil(t, result, "Ergebnis sollte nicht nil sein")

		// Ausgabe der Ergebnisse
		fmt.Printf("Lernkosten für Menschenkenntnis als Hexer:\n")
		fmt.Printf("Lerneinheiten (LE): %d\n", result.LE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", result.Ep)
		fmt.Printf("Geldkosten (GS): %d\n", result.Money)

		// Überprüfung der Werte basierend auf aktueller Implementierung
		// Menschenkenntnis ist Sozial/schwer: 4 LE
		// Hexer EP-Kosten für Sozial: 20 EP/TE
		// Money-Kosten: 20 GS/LE
		assert.Equal(t, 4, result.LE, "LE-Kosten sollten 4 sein")
		assert.Equal(t, 80, result.Ep, "EP-Kosten sollten 80 sein (4 LE * 20 EP)")
		assert.Equal(t, 80, result.Money, "Geldkosten sollten 80 GS sein (4 LE * 20 GS)")
	})

	// Testfall für Menschenkenntnis als Hexer - Verbesserungskosten
	t.Run("Verbesserungskosten für Menschenkenntnis als Hexer", func(t *testing.T) {
		// Verwende die exportierte Funktion CalculateDetailedSkillImprovementCost
		currentLevel := 10
		result, err := CalculateDetailedSkillImprovementCost("Menschenkenntnis", "Hexer", currentLevel)
		assert.NoError(t, err, "CalculateDetailedSkillImprovementCost sollte keinen Fehler zurückgeben")
		assert.NotNil(t, result, "Ergebnis sollte nicht nil sein")

		// Ausgabe der Ergebnisse
		fmt.Printf("\nVerbesserungskosten für Menschenkenntnis als Hexer (von %d auf %d):\n", currentLevel, currentLevel+1)
		fmt.Printf("Zielstufe: %d\n", result.Stufe)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", result.Ep)
		fmt.Printf("Geldkosten (GS): %d\n", result.Money)

		// Überprüfung der Werte basierend auf aktueller Implementierung
		// Result.Stufe ist die Zielstufe (11), nicht die benötigten TE
		assert.Equal(t, 11, result.Stufe, "Zielstufe sollte 11 sein (von 10 auf 11)")
		assert.Equal(t, 80, result.Ep, "EP-Kosten sollten 80 sein")
		assert.Equal(t, 80, result.Money, "Geldkosten sollten 80 GS sein")
	})
}
