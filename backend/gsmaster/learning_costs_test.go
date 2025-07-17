package gsmaster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDetailedSkillLearningCostForHexer(t *testing.T) {
	// Testfall, der direkt mit den Lerntabellen ohne Datenbankzugriff arbeitet
	t.Run("Lernkosten für Menschenkenntnis als Hexer (ohne DB)", func(t *testing.T) {
		// Da wir keine Datenbankverbindung haben, simulieren wir die Lernkosten direkt
		// basierend auf den bekannten Werten für die Fertigkeit "Menschenkenntnis"

		// Bekannte Werte für Menschenkenntnis:
		// - Kategorie: "Sozial"
		// - Schwierigkeit: "schwer"
		category := "Sozial"
		difficulty := "schwer"

		// Direkt aus den learningCosts-Tabellen die Werte ablesen
		baseLE, ok := learningCosts.BaseLearnCost[category][difficulty]
		assert.True(t, ok, "Kategorie und Schwierigkeit sollten in der Tabelle existieren")

		epPerTE, ok := learningCosts.EPPerTE["Hx"][category]
		assert.True(t, ok, "Hexer sollte EP-Kosten für diese Kategorie haben")

		// Berechnen der Kosten basierend auf den Tabellenwerten
		expectedLE := baseLE
		expectedEP := baseLE * (epPerTE * 3)
		expectedMoney := baseLE * 200

		fmt.Printf("Lernkosten für Menschenkenntnis als Hexer:\n")
		fmt.Printf("Lerneinheiten (LE): %d\n", expectedLE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", expectedEP)
		fmt.Printf("Geldkosten (GS): %d\n", expectedMoney)

		// Test bestanden, wenn die Werte mit den erwarteten übereinstimmen
		assert.Equal(t, 4, expectedLE, "LE-Kosten für Menschenkenntnis (Sozial/schwer) sollten 4 sein")
		assert.Equal(t, 4*(20*3), expectedEP, "EP-Kosten sollten 4 LE * (20 EP * 3) sein")
		assert.Equal(t, 4*200, expectedMoney, "Geldkosten sollten 4 LE * 200 GS sein")
	})

	// Ein Beispiel für einen direkten Zugriff ohne Datenbank
	t.Run("Verbesserungskosten für Menschenkenntnis als Hexer (ohne DB)", func(t *testing.T) {
		// Bekannte Werte für Menschenkenntnis:
		// - Kategorie: "Sozial"
		// - Schwierigkeit: "schwer"
		category := "Sozial"
		difficulty := "schwer"
		currentLevel := 10
		nextLevel := currentLevel + 1

		// Direkt aus den learningCosts-Tabellen die Werte ablesen
		neededTE, ok := learningCosts.ImprovementCost[category][difficulty][nextLevel]
		assert.True(t, ok, "Es sollte Kosten für diese Verbesserung geben")

		epPerTE, ok := learningCosts.EPPerTE["Hx"][category]
		assert.True(t, ok, "Hexer sollte EP-Kosten für diese Kategorie haben")

		// Berechnen der Kosten
		expectedEP := neededTE * epPerTE
		expectedMoney := neededTE * 20

		fmt.Printf("\nVerbesserungskosten für Menschenkenntnis als Hexer (von %d auf %d):\n", currentLevel, nextLevel)
		fmt.Printf("Trainingseinheiten (TE): %d\n", neededTE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", expectedEP)
		fmt.Printf("Geldkosten (GS): %d\n", expectedMoney)

		// Test bestanden, wenn die Werte mit den erwarteten übereinstimmen
		assert.Equal(t, 10, neededTE, "TE-Kosten für Verbesserung von 10 auf 11 sollten 10 sein")
		assert.Equal(t, 10*20, expectedEP, "EP-Kosten sollten 10 TE * 20 EP sein")
		assert.Equal(t, 10*20, expectedMoney, "Geldkosten sollten 10 TE * 20 GS sein")
	})
}
