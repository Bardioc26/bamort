package gsmaster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeidhändigerKampfFürPS(t *testing.T) {
	// Testfall für das Erlernen von "Beidhändiger Kampf" durch einen Priester Streiter (PS)
	t.Run("Lernkosten für Beidhändiger Kampf als PS", func(t *testing.T) {
		// Bekannte Werte für Beidhändiger Kampf:
		// - Kategorie: "Kampf"
		// - Schwierigkeit: "sehr_schwer"
		category := "Kampf"
		difficulty := "sehr_schwer"

		// Direkt aus den learningCosts-Tabellen die Werte ablesen
		baseLE, ok := learningCosts.BaseLearnCost[category][difficulty]
		assert.True(t, ok, "Kategorie und Schwierigkeit sollten in der Tabelle existieren")

		epPerTE, ok := learningCosts.EPPerTE["PS"][category]
		assert.True(t, ok, "PS sollte EP-Kosten für diese Kategorie haben")

		// Berechnen der Kosten basierend auf den Tabellenwerten
		expectedLE := baseLE
		expectedEP := baseLE * (epPerTE * 3)
		expectedMoney := baseLE * 200

		// Ausgabe der Ergebnisse
		fmt.Printf("Lernkosten für Beidhändiger Kampf als PS:\n")
		fmt.Printf("Lerneinheiten (LE): %d\n", expectedLE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", expectedEP)
		fmt.Printf("Geldkosten (GS): %d\n", expectedMoney)

		// Überprüfung der Werte
		assert.Equal(t, 10, expectedLE, "LE-Kosten sollten 10 sein")
		assert.Equal(t, 10*(30*3), expectedEP, "EP-Kosten sollten 10 LE * (30 EP * 3) sein")
		assert.Equal(t, 10*200, expectedMoney, "Geldkosten sollten 10 LE * 200 GS sein")
	})

	// Testfall für das Verbessern von "Beidhändiger Kampf" durch einen Priester Streiter (PS)
	// von Stufe 5 auf 6
	t.Run("Verbesserungskosten für Beidhändiger Kampf als PS von 5 auf 6", func(t *testing.T) {
		// Bekannte Werte für Beidhändiger Kampf:
		// - Kategorie: "Kampf"
		// - Schwierigkeit: "sehr_schwer"
		category := "Kampf"
		difficulty := "sehr_schwer"
		currentLevel := 5
		nextLevel := 6

		// Aus der Tabelle die Werte abrufen
		neededTE, ok := learningCosts.ImprovementCost[category][difficulty][nextLevel]
		assert.True(t, ok, "Es sollte Kosten für diese Verbesserung geben")

		epPerTE, ok := learningCosts.EPPerTE["PS"][category]
		assert.True(t, ok, "PS sollte EP-Kosten für diese Kategorie haben")

		// Kosten berechnen
		expectedEP := neededTE * epPerTE
		expectedMoney := neededTE * 20

		// Ausgabe der Ergebnisse
		fmt.Printf("\nVerbesserungskosten für Beidhändiger Kampf als PS (von %d auf %d):\n", currentLevel, nextLevel)
		fmt.Printf("Trainingseinheiten (TE): %d\n", neededTE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", expectedEP)
		fmt.Printf("Geldkosten (GS): %d\n", expectedMoney)

		// Überprüfung der Werte
		assert.Equal(t, 2, neededTE, "TE-Kosten sollten 2 sein")
		assert.Equal(t, 2*30, expectedEP, "EP-Kosten sollten 2 TE * 30 EP sein")
		assert.Equal(t, 2*20, expectedMoney, "Geldkosten sollten 2 TE * 20 GS sein")
	})

	// Testfall für das Verbessern von "Beidhändiger Kampf" durch einen Priester Streiter (PS)
	// von Stufe 6 auf 7
	t.Run("Verbesserungskosten für Beidhändiger Kampf als PS von 6 auf 7", func(t *testing.T) {
		// Bekannte Werte für Beidhändiger Kampf:
		// - Kategorie: "Kampf"
		// - Schwierigkeit: "sehr_schwer"
		category := "Kampf"
		difficulty := "sehr_schwer"
		currentLevel := 6
		nextLevel := 7

		// Aus der Tabelle die Werte abrufen
		neededTE, ok := learningCosts.ImprovementCost[category][difficulty][nextLevel]
		assert.True(t, ok, "Es sollte Kosten für diese Verbesserung geben")

		epPerTE, ok := learningCosts.EPPerTE["PS"][category]
		assert.True(t, ok, "PS sollte EP-Kosten für diese Kategorie haben")

		// Kosten berechnen
		expectedEP := neededTE * epPerTE
		expectedMoney := neededTE * 20

		// Ausgabe der Ergebnisse
		fmt.Printf("\nVerbesserungskosten für Beidhändiger Kampf als PS (von %d auf %d):\n", currentLevel, nextLevel)
		fmt.Printf("Trainingseinheiten (TE): %d\n", neededTE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", expectedEP)
		fmt.Printf("Geldkosten (GS): %d\n", expectedMoney)

		// Überprüfung der Werte
		assert.Equal(t, 5, neededTE, "TE-Kosten sollten 5 sein")
		assert.Equal(t, 5*30, expectedEP, "EP-Kosten sollten 5 TE * 30 EP sein")
		assert.Equal(t, 5*20, expectedMoney, "Geldkosten sollten 5 TE * 20 GS sein")
	})

	// Testfall für das Verbessern von "Beidhändiger Kampf" durch einen Priester Streiter (PS)
	// von Stufe 7 auf 8
	t.Run("Verbesserungskosten für Beidhändiger Kampf als PS von 7 auf 8", func(t *testing.T) {
		// Bekannte Werte für Beidhändiger Kampf:
		// - Kategorie: "Kampf"
		// - Schwierigkeit: "sehr_schwer"
		category := "Kampf"
		difficulty := "sehr_schwer"
		currentLevel := 7
		nextLevel := 8

		// Aus der Tabelle die Werte abrufen
		neededTE, ok := learningCosts.ImprovementCost[category][difficulty][nextLevel]
		assert.True(t, ok, "Es sollte Kosten für diese Verbesserung geben")

		epPerTE, ok := learningCosts.EPPerTE["PS"][category]
		assert.True(t, ok, "PS sollte EP-Kosten für diese Kategorie haben")

		// Kosten berechnen
		expectedEP := neededTE * epPerTE
		expectedMoney := neededTE * 20

		// Ausgabe der Ergebnisse
		fmt.Printf("\nVerbesserungskosten für Beidhändiger Kampf als PS (von %d auf %d):\n", currentLevel, nextLevel)
		fmt.Printf("Trainingseinheiten (TE): %d\n", neededTE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", expectedEP)
		fmt.Printf("Geldkosten (GS): %d\n", expectedMoney)

		// Überprüfung der Werte
		assert.Equal(t, 10, neededTE, "TE-Kosten sollten 10 sein")
		assert.Equal(t, 10*30, expectedEP, "EP-Kosten sollten 10 TE * 30 EP sein")
		assert.Equal(t, 10*20, expectedMoney, "Geldkosten sollten 10 TE * 20 GS sein")
	})
}
