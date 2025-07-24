package gsmaster

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeidhändigerKampfFürPS(t *testing.T) {
	// Testfall für das Erlernen von "Beidhändiger Kampf" durch einen Priester Streiter (PS)
	t.Run("Lernkosten für Beidhändiger Kampf als PS", func(t *testing.T) {
		// Verwende exportierte Funktion zum Berechnen der Lernkosten
		result, err := CalculateDetailedSkillLearningCost("Beidhändiger Kampf", "PS")
		assert.NoError(t, err, "Es sollte keinen Fehler beim Berechnen der Lernkosten geben")
		assert.True(t, result.LE > 0, "LE-Kosten sollten größer als 0 sein")

		// Ausgabe der Ergebnisse
		fmt.Printf("Lernkosten für Beidhändiger Kampf als PS:\n")
		fmt.Printf("Lerneinheiten (LE): %d\n", result.LE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", result.Ep)
		fmt.Printf("Geldkosten (GS): %d\n", result.Money)

		// Überprüfung der Werte basierend auf der aktuellen Implementierung
		assert.Equal(t, 2, result.LE, "LE-Kosten sollten 2 sein")
		assert.Equal(t, 40, result.Ep, "EP-Kosten sollten 40 sein")
		assert.Equal(t, 40, result.Money, "Geldkosten sollten 40 GS sein")
	})

	// Testfall für das Verbessern von "Beidhändiger Kampf" durch einen Priester Streiter (PS)
	// von Stufe 5 auf 6
	t.Run("Verbesserungskosten für Beidhändiger Kampf als PS von 5 auf 6", func(t *testing.T) {
		// Verwende exportierte Funktion zum Berechnen der Verbesserungskosten
		result, err := CalculateDetailedSkillImprovementCost("Beidhändiger Kampf", "PS", 5)
		assert.NoError(t, err, "Es sollte keinen Fehler beim Berechnen der Verbesserungskosten geben")
		assert.True(t, result.LE > 0, "LE-Kosten sollten größer als 0 sein")

		// Ausgabe der Ergebnisse
		fmt.Printf("\nVerbesserungskosten für Beidhändiger Kampf als PS (von 5 auf 6):\n")
		fmt.Printf("Trainingseinheiten (TE): %d\n", result.LE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", result.Ep)
		fmt.Printf("Geldkosten (GS): %d\n", result.Money)

		// Überprüfung der Werte basierend auf der aktuellen vereinfachten Implementierung
		assert.Equal(t, 1, result.LE, "TE-Kosten sollten 1 sein (vereinfachte Implementierung)")
		assert.Equal(t, 40, result.Ep, "EP-Kosten sollten 40 sein")
		assert.Equal(t, 40, result.Money, "Geldkosten sollten 40 GS sein")
	})

	// Testfall für das Verbessern von "Beidhändiger Kampf" durch einen Priester Streiter (PS)
	// von Stufe 6 auf 7
	t.Run("Verbesserungskosten für Beidhändiger Kampf als PS von 6 auf 7", func(t *testing.T) {
		// Verwende exportierte Funktion zum Berechnen der Verbesserungskosten
		result, err := CalculateDetailedSkillImprovementCost("Beidhändiger Kampf", "PS", 6)
		assert.NoError(t, err, "Es sollte keinen Fehler beim Berechnen der Verbesserungskosten geben")
		assert.True(t, result.LE > 0, "LE-Kosten sollten größer als 0 sein")

		// Ausgabe der Ergebnisse
		fmt.Printf("\nVerbesserungskosten für Beidhändiger Kampf als PS (von 6 auf 7):\n")
		fmt.Printf("Trainingseinheiten (TE): %d\n", result.LE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", result.Ep)
		fmt.Printf("Geldkosten (GS): %d\n", result.Money)

		// Überprüfung der Werte basierend auf der aktuellen vereinfachten Implementierung
		assert.Equal(t, 1, result.LE, "TE-Kosten sollten 1 sein (vereinfachte Implementierung)")
		assert.Equal(t, 40, result.Ep, "EP-Kosten sollten 40 sein")
		assert.Equal(t, 40, result.Money, "Geldkosten sollten 40 GS sein")
	})

	// Testfall für das Verbessern von "Beidhändiger Kampf" durch einen Priester Streiter (PS)
	// von Stufe 7 auf 8
	t.Run("Verbesserungskosten für Beidhändiger Kampf als PS von 7 auf 8", func(t *testing.T) {
		// Verwende exportierte Funktion zum Berechnen der Verbesserungskosten
		result, err := CalculateDetailedSkillImprovementCost("Beidhändiger Kampf", "PS", 7)
		assert.NoError(t, err, "Es sollte keinen Fehler beim Berechnen der Verbesserungskosten geben")
		assert.True(t, result.LE > 0, "LE-Kosten sollten größer als 0 sein")

		// Ausgabe der Ergebnisse
		fmt.Printf("\nVerbesserungskosten für Beidhändiger Kampf als PS (von 7 auf 8):\n")
		fmt.Printf("Trainingseinheiten (TE): %d\n", result.LE)
		fmt.Printf("Erfahrungspunkte (EP): %d\n", result.Ep)
		fmt.Printf("Geldkosten (GS): %d\n", result.Money)

		// Überprüfung der Werte basierend auf der aktuellen vereinfachten Implementierung
		assert.Equal(t, 1, result.LE, "TE-Kosten sollten 1 sein (vereinfachte Implementierung)")
		assert.Equal(t, 40, result.Ep, "EP-Kosten sollten 40 sein")
		assert.Equal(t, 40, result.Money, "Geldkosten sollten 40 GS sein")
	})
}
