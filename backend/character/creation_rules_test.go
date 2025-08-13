package character

import (
	"testing"
)

func TestGetSpecialAbilityByRoll(t *testing.T) {
	testCases := []struct {
		roll     int
		expected string
		modifier int
	}{
		{1, "Sehen", -2},
		{5, "Sehen", -2},
		{6, "Hören", -2},
		{16, "Sechster Sinn", 2},
		{25, "Sehen", 2},
		{55, "Nachtsicht", 2},
		{60, "Gute Reflexe", 9},
		{65, "Richtungssinn", 12},
		{70, "Robustheit", 9},
		{75, "Schmerzunempfindlichkeit", 9},
		{80, "Trinken", 12},
		{85, "Wachgabe", 6},
		{90, "Wahrnehmung", 8},
		{95, "Einprägen", 4},
		{96, "Berserkergang", 0},
		{99, "Berserkergang", 0},
		{100, "Freie Wahl", 0},
	}

	for _, tc := range testCases {
		ability, err := GetSpecialAbilityByRoll(tc.roll)
		if err != nil {
			t.Errorf("Unerwarteter Fehler für Wurf %d: %v", tc.roll, err)
			continue
		}
		if ability.Name != tc.expected {
			t.Errorf("Für Wurf %d: erwartet '%s', erhalten '%s'", tc.roll, tc.expected, ability.Name)
		}
		if ability.Modifier != tc.modifier {
			t.Errorf("Für Wurf %d: erwarteter Modifier %d, erhalten %d", tc.roll, tc.modifier, ability.Modifier)
		}
	}
}

func TestGetSpecialAbilityByRollInvalidInput(t *testing.T) {
	invalidRolls := []int{0, -1, 101, 1000}

	for _, roll := range invalidRolls {
		_, err := GetSpecialAbilityByRoll(roll)
		if err == nil {
			t.Errorf("Erwartete Fehler für ungültigen Wurf %d, aber keinen erhalten", roll)
		}
	}
}

func TestGetAllSpecialAbilities(t *testing.T) {
	abilities := GetAllSpecialAbilities()

	if len(abilities) != 18 {
		t.Errorf("Erwartete 18 besondere Fähigkeiten, erhalten %d", len(abilities))
	}

	// Teste, dass alle Bereiche lückenlos von 1-100 abgedeckt sind
	for i, ability := range abilities {
		if i == 0 && ability.RangeStart != 1 {
			t.Errorf("Erster Bereich sollte bei 1 starten, startet bei %d", ability.RangeStart)
		}
		if i == len(abilities)-1 && ability.RangeEnd != 100 {
			t.Errorf("Letzter Bereich sollte bei 100 enden, endet bei %d", ability.RangeEnd)
		}
	}
}
