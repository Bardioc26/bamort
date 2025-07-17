package gsmaster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLerntabellenForHexer(t *testing.T) {
	// Testfall für Menschenkenntnis als Hexer - Lernkosten
	t.Run("Lernkosten für Menschenkenntnis als Hexer", func(t *testing.T) {
		// Bekannte Werte für Menschenkenntnis
		category := "Sozial"
		difficulty := "schwer"

		// Aus der Tabelle die Werte abrufen
		baseLE, ok := learningCosts.BaseLearnCost[category][difficulty]
		assert.True(t, ok, "Kategorie und Schwierigkeit sollten in der Tabelle existieren")

		epPerTE, ok := learningCosts.EPPerTE["Hx"][category]
		assert.True(t, ok, "Hexer sollte EP-Kosten für diese Kategorie haben")

		// Kosten berechnen
		expectedLE := baseLE
		expectedEP := baseLE * (epPerTE * 3)
		expectedMoney := baseLE * 200

		// Ausgabe der Ergebnisse
		fmt.Printf("Lernkosten für Menschenkenntnis als Hexer:\n")
		fmt.Printf("Lerneinheiten (LE): %d\n", expectedLE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", expectedEP)
		fmt.Printf("Geldkosten (GS): %d\n", expectedMoney)

		// Überprüfung der Werte
		assert.Equal(t, 4, expectedLE, "LE-Kosten sollten 4 sein")
		assert.Equal(t, 4*(20*3), expectedEP, "EP-Kosten sollten 240 sein")
		assert.Equal(t, 4*200, expectedMoney, "Geldkosten sollten 800 GS sein")
	})

	// Testfall für Menschenkenntnis als Hexer - Verbesserungskosten
	t.Run("Verbesserungskosten für Menschenkenntnis als Hexer", func(t *testing.T) {
		// Bekannte Werte für Menschenkenntnis
		category := "Sozial"
		difficulty := "schwer"
		currentLevel := 10
		nextLevel := currentLevel + 1

		// Aus der Tabelle die Werte abrufen
		neededTE, ok := learningCosts.ImprovementCost[category][difficulty][nextLevel]
		assert.True(t, ok, "Es sollte Kosten für diese Verbesserung geben")

		epPerTE, ok := learningCosts.EPPerTE["Hx"][category]
		assert.True(t, ok, "Hexer sollte EP-Kosten für diese Kategorie haben")

		// Kosten berechnen
		expectedEP := neededTE * epPerTE
		expectedMoney := neededTE * 20

		// Ausgabe der Ergebnisse
		fmt.Printf("\nVerbesserungskosten für Menschenkenntnis als Hexer (von %d auf %d):\n", currentLevel, nextLevel)
		fmt.Printf("Trainingseinheiten (TE): %d\n", neededTE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", expectedEP)
		fmt.Printf("Geldkosten (GS): %d\n", expectedMoney)

		// Überprüfung der Werte
		assert.Equal(t, 10, neededTE, "TE-Kosten sollten 10 sein")
		assert.Equal(t, 10*20, expectedEP, "EP-Kosten sollten 200 sein")
		assert.Equal(t, 10*20, expectedMoney, "Geldkosten sollten 200 GS sein")
	})
}
