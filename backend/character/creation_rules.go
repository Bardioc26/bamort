package character

import "fmt"

// SpecialAbility repräsentiert eine besondere Fähigkeit mit Name und Modifier
type SpecialAbility struct {
	Name     string `json:"name"`
	Modifier int    `json:"modifier"`
	Special  string `json:"special,omitempty"` // Für spezielle Fälle wie Berserkergang
}

// GetSpecialAbilityByRoll gibt die besondere Fähigkeit basierend auf einem W100-Wurf zurück
// W%	Besondere Fähigkeit
// 01-05:	Sehen–2
// 06-10:	Hören–2
// 11-15:	Riechen–2
// 16-20:	Sechster Sinn+2
// 21-30:	Sehen+2
// 31-40:	Hören+2
// 41-50:	Riechen+2
// 51-55:	Nachtsicht+2
// 56-60:	Gute Reﬂexe+9
// 61-65:	Richtungssinn+12
// 66-70:	Robustheit+9
// 71-75:	Schmerzunempﬁndlichkeit+9
// 76-80:	Trinken+12
// 81-85:	Wachgabe+6
// 86-90:	Wahrnehmung+8
// 91-95:	Einprägen+4
// 96-99:	Berserkergang+(18–Wk/5)
// 100:	freie Wahl und zweiter Wurf
func GetSpecialAbilityByRoll(roll int) (*SpecialAbility, error) {
	if roll < 1 || roll > 100 {
		return nil, fmt.Errorf("ungültiger Würfelwurf: %d (muss zwischen 1 und 100 liegen)", roll)
	}

	switch {
	case roll >= 1 && roll <= 5:
		return &SpecialAbility{Name: "Sehen", Modifier: -2}, nil
	case roll >= 6 && roll <= 10:
		return &SpecialAbility{Name: "Hören", Modifier: -2}, nil
	case roll >= 11 && roll <= 15:
		return &SpecialAbility{Name: "Riechen", Modifier: -2}, nil
	case roll >= 16 && roll <= 20:
		return &SpecialAbility{Name: "Sechster Sinn", Modifier: 2}, nil
	case roll >= 21 && roll <= 30:
		return &SpecialAbility{Name: "Sehen", Modifier: 2}, nil
	case roll >= 31 && roll <= 40:
		return &SpecialAbility{Name: "Hören", Modifier: 2}, nil
	case roll >= 41 && roll <= 50:
		return &SpecialAbility{Name: "Riechen", Modifier: 2}, nil
	case roll >= 51 && roll <= 55:
		return &SpecialAbility{Name: "Nachtsicht", Modifier: 2}, nil
	case roll >= 56 && roll <= 60:
		return &SpecialAbility{Name: "Gute Reflexe", Modifier: 9}, nil
	case roll >= 61 && roll <= 65:
		return &SpecialAbility{Name: "Richtungssinn", Modifier: 12}, nil
	case roll >= 66 && roll <= 70:
		return &SpecialAbility{Name: "Robustheit", Modifier: 9}, nil
	case roll >= 71 && roll <= 75:
		return &SpecialAbility{Name: "Schmerzunempfindlichkeit", Modifier: 9}, nil
	case roll >= 76 && roll <= 80:
		return &SpecialAbility{Name: "Trinken", Modifier: 12}, nil
	case roll >= 81 && roll <= 85:
		return &SpecialAbility{Name: "Wachgabe", Modifier: 6}, nil
	case roll >= 86 && roll <= 90:
		return &SpecialAbility{Name: "Wahrnehmung", Modifier: 8}, nil
	case roll >= 91 && roll <= 95:
		return &SpecialAbility{Name: "Einprägen", Modifier: 4}, nil
	case roll >= 96 && roll <= 99:
		return &SpecialAbility{
			Name:     "Berserkergang",
			Modifier: 0,
			Special:  "+(18–Wk/5)",
		}, nil
	case roll == 100:
		return &SpecialAbility{
			Name:     "Freie Wahl",
			Modifier: 0,
			Special:  "freie Wahl und zweiter Wurf",
		}, nil
	default:
		return nil, fmt.Errorf("unerwarteter Würfelwurf: %d", roll)
	}
}

// GetAllSpecialAbilities gibt alle möglichen besonderen Fähigkeiten mit ihren Wahrscheinlichkeiten zurück
func GetAllSpecialAbilities() []struct {
	RangeStart int             `json:"range_start"`
	RangeEnd   int             `json:"range_end"`
	Ability    *SpecialAbility `json:"ability"`
} {
	return []struct {
		RangeStart int             `json:"range_start"`
		RangeEnd   int             `json:"range_end"`
		Ability    *SpecialAbility `json:"ability"`
	}{
		{1, 5, &SpecialAbility{Name: "Sehen", Modifier: -2}},
		{6, 10, &SpecialAbility{Name: "Hören", Modifier: -2}},
		{11, 15, &SpecialAbility{Name: "Riechen", Modifier: -2}},
		{16, 20, &SpecialAbility{Name: "Sechster Sinn", Modifier: 2}},
		{21, 30, &SpecialAbility{Name: "Sehen", Modifier: 2}},
		{31, 40, &SpecialAbility{Name: "Hören", Modifier: 2}},
		{41, 50, &SpecialAbility{Name: "Riechen", Modifier: 2}},
		{51, 55, &SpecialAbility{Name: "Nachtsicht", Modifier: 2}},
		{56, 60, &SpecialAbility{Name: "Gute Reflexe", Modifier: 9}},
		{61, 65, &SpecialAbility{Name: "Richtungssinn", Modifier: 12}},
		{66, 70, &SpecialAbility{Name: "Robustheit", Modifier: 9}},
		{71, 75, &SpecialAbility{Name: "Schmerzunempfindlichkeit", Modifier: 9}},
		{76, 80, &SpecialAbility{Name: "Trinken", Modifier: 12}},
		{81, 85, &SpecialAbility{Name: "Wachgabe", Modifier: 6}},
		{86, 90, &SpecialAbility{Name: "Wahrnehmung", Modifier: 8}},
		{91, 95, &SpecialAbility{Name: "Einprägen", Modifier: 4}},
		{96, 99, &SpecialAbility{Name: "Berserkergang", Modifier: 0, Special: "+(18–Wk/5)"}},
		{100, 100, &SpecialAbility{Name: "Freie Wahl", Modifier: 0, Special: "freie Wahl und zweiter Wurf"}},
	}
}
