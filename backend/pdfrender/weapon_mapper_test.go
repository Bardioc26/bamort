package pdfrender

import (
	"bamort/models"
	"testing"
)

// TestMapWeapons_WithEWCalculation tests that weapon EW is correctly calculated
// EW = Waffenfertigkeit.Fertigkeitswert + Character.AngriffBonus + Weapon.Anb
func TestMapWeapons_WithEWCalculation(t *testing.T) {
	// Arrange - Create a character with weapon skills and equipped weapons
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Fighter",
		},
		// Weapon skills the character has learned
		Waffenfertigkeiten: []models.SkWaffenfertigkeit{
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{
							Name: "Langschwert",
						},
					},
					Fertigkeitswert: 12, // Base skill value
				},
			},
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{
							Name: "Bogen",
						},
					},
					Fertigkeitswert: 10,
				},
			},
		},
		// Actual equipped weapons with bonuses
		Waffen: []models.EqWaffe{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Langschwert",
					},
				},
				Anb:  2,  // Attack bonus of weapon
				Schb: 1,  // Damage bonus of weapon
			},
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Bogen",
					},
				},
				Anb:  0,
				Schb: 0,
			},
		},
		// Character attributes that influence combat
		Eigenschaften: []models.Eigenschaft{
			{Name: "St", Value: 90}, // Strength influences attack bonus
			{Name: "Gw", Value: 80}, // Dexterity influences attack bonus
		},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(vm.Weapons) != 2 {
		t.Fatalf("Expected 2 weapons, got %d", len(vm.Weapons))
	}

	// Test Langschwert
	sword := vm.Weapons[0]
	if sword.Name != "Langschwert" {
		t.Errorf("Expected weapon name 'Langschwert', got '%s'", sword.Name)
	}
	
	// EW should be: skill(12) + char_attack_bonus + weapon_bonus(2)
	// Character attack bonus is derived from St/Gw, typically calculated in DerivedValues
	// For now, we expect at least: 12 + 2 = 14 (we'll enhance with char bonus later)
	if sword.Value < 14 {
		t.Errorf("Expected sword EW >= 14 (skill 12 + weapon bonus 2), got %d", sword.Value)
	}

	// Test Bogen
	bow := vm.Weapons[1]
	if bow.Name != "Bogen" {
		t.Errorf("Expected weapon name 'Bogen', got '%s'", bow.Name)
	}
	
	if bow.Value < 10 {
		t.Errorf("Expected bow EW >= 10 (skill 10 + no weapon bonus), got %d", bow.Value)
	}
}

// TestMapWeapons_WithDamageCalculation tests that weapon damage is correctly calculated
func TestMapWeapons_WithDamageCalculation(t *testing.T) {
	// This test will be implemented to verify damage calculation
	// For now, we mark it as a placeholder for TDD
	t.Skip("TODO: Implement damage calculation test")
}

// TestMapWeapons_WithRangedWeaponRanges tests that ranged weapons show their ranges
func TestMapWeapons_WithRangedWeaponRanges(t *testing.T) {
	// This test will verify that ranged weapons (Bogen, Armbrust) show ranges
	t.Skip("TODO: Implement ranged weapon ranges test")
}
