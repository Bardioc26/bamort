package pdfrender

import (
	"bamort/database"
	"bamort/models"
	"testing"
)

// TestMapWeapons_WithEWCalculation tests that weapon EW is correctly calculated
// EW = Waffenfertigkeit.Fertigkeitswert + Character.AngriffBonus + Weapon.Anb
func TestMapWeapons_WithEWCalculation(t *testing.T) {
	// Setup test database
	database.SetupTestDB()

	// Create weapons in gsmaster with skill requirements
	testSword := &models.Weapon{
		Equipment: models.Equipment{
			GameSystem: "midgard",
			Name:       "Langschwert",
		},
		SkillRequired: "Schwerter",
		Damage:        "1W6",
	}
	_ = testSword.Create()

	testBow := &models.Weapon{
		Equipment: models.Equipment{
			GameSystem: "midgard",
			Name:       "Langbogen",
		},
		SkillRequired: "Bogen",
		Damage:        "1W6",
	}
	_ = testBow.Create()

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
							Name: "Schwerter",
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
				Anb:  2, // Attack bonus of weapon
				Schb: 1, // Damage bonus of weapon
			},
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Langbogen",
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

	// Expect 3 weapons: Raufen + Langschwert + Langbogen
	if len(vm.Weapons) != 3 {
		t.Fatalf("Expected 3 weapons (Raufen + 2 equipped), got %d", len(vm.Weapons))
	}

	// Check first weapon is Raufen
	if vm.Weapons[0].Name != "Raufen" {
		t.Errorf("Expected first weapon to be 'Raufen', got '%s'", vm.Weapons[0].Name)
	}

	// Test Langschwert (index 1)
	sword := vm.Weapons[1]
	if sword.Name != "Langschwert" {
		t.Errorf("Expected weapon name 'Langschwert', got '%s'", sword.Name)
	}

	// EW should be: skill(12) + char_attack_bonus + weapon_bonus(2)
	// Character attack bonus is derived from St/Gw, typically calculated in DerivedValues
	// For now, we expect at least: 12 + 2 = 14 (we'll enhance with char bonus later)
	if sword.Value < 14 {
		t.Errorf("Expected sword EW >= 14 (skill 12 + weapon bonus 2), got %d", sword.Value)
	}

	// Test Langbogen (index 2)
	bow := vm.Weapons[2]
	if bow.Name != "Langbogen" {
		t.Errorf("Expected weapon name 'Langbogen', got '%s'", bow.Name)
	}

	if bow.Value < 10 {
		t.Errorf("Expected bow EW >= 10 (skill 10 + no weapon bonus), got %d", bow.Value)
	}
}

// TestMapWeapons_WithDamageCalculation tests that weapon damage is correctly calculated
// Damage = Base Weapon Damage + Character SchadenBonus + Weapon Schb
func TestMapWeapons_WithDamageCalculation(t *testing.T) {
	// Setup test database to load weapon data
	database.SetupTestDB()

	// Create test weapon in gsmaster if it doesn't exist
	testWeapon := &models.Weapon{
		Equipment: models.Equipment{
			GameSystem: "midgard",
			Name:       "Langschwert",
		},
		SkillRequired: "Schwerter",
		Damage:        "1W6",
	}
	_ = testWeapon.Create() // Ignore error if already exists

	// Arrange - Create a character with weapons that have damage values
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Fighter",
		},
		// Weapon skills
		Waffenfertigkeiten: []models.SkWaffenfertigkeit{
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{
							Name: "Schwerter",
						},
					},
					Fertigkeitswert: 12,
				},
			},
		},
		// Equipped weapons with damage bonuses
		Waffen: []models.EqWaffe{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Langschwert",
					},
				},
				Schb: 2, // Weapon's damage bonus
			},
		},
		// Character attributes that influence damage bonus
		Eigenschaften: []models.Eigenschaft{
			{Name: "St", Value: 95}, // High strength gives SchadenBonus = +1
			{Name: "Gw", Value: 80},
		},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Expect 2 weapons: Raufen + Langschwert
	if len(vm.Weapons) != 2 {
		t.Fatalf("Expected 2 weapons (Raufen + Langschwert), got %d", len(vm.Weapons))
	}

	// Check Langschwert (index 1, after Raufen)
	sword := vm.Weapons[1]

	// Damage should be in format like "1W6+3" where the bonus combines
	// Character's SchadenBonus (from St=95 -> +1) + Weapon's Schb (2) = +3
	if sword.Damage == "" {
		t.Error("Expected weapon to have damage value, got empty string")
	}

	// The damage string should be "1W6+3"
	// Base damage "1W6" + total bonus (+3)
	expectedDamage := "1W6+3"
	if sword.Damage != expectedDamage {
		t.Errorf("Expected damage '%s', got '%s'", expectedDamage, sword.Damage)
	}

	t.Logf("✓ Weapon damage correctly calculated: %s (SchadenBonus %d + Weapon Schb 2 = +3)",
		sword.Damage, vm.DerivedValues.SchadenBonus)
}

// TestMapWeapons_WithRangedWeaponRanges tests that ranged weapons show their ranges
func TestMapWeapons_WithRangedWeaponRanges(t *testing.T) {
	// Setup test database
	database.SetupTestDB()

	// Delete existing test weapons to ensure clean state
	database.DB.Where("name IN (?, ?)", "Langbogen", "Kurzschwert").Delete(&models.Weapon{})

	// Create a ranged weapon in gsmaster with range values
	testBow := &models.Weapon{
		Equipment: models.Equipment{
			GameSystem: "midgard",
			Name:       "Langbogen",
		},
		SkillRequired: "Bogen",
		Damage:        "1W6",
		RangeNear:     10,
		RangeMiddle:   30,
		RangeFar:      100,
	}
	err := testBow.Create()
	if err != nil {
		t.Fatalf("Failed to create Langbogen: %v", err)
	}

	// Verify the weapon was created
	verifyWeapon := &models.Weapon{}
	verifyErr := verifyWeapon.First("Langbogen")
	if verifyErr != nil || verifyWeapon.ID == 0 {
		t.Fatalf("Failed to create/find Langbogen in database: %v", verifyErr)
	}
	if verifyWeapon.RangeNear != 10 || verifyWeapon.RangeMiddle != 30 || verifyWeapon.RangeFar != 100 {
		t.Fatalf("Langbogen ranges not set correctly: got %d/%d/%d, expected 10/30/100",
			verifyWeapon.RangeNear, verifyWeapon.RangeMiddle, verifyWeapon.RangeFar)
	}

	// Create a melee weapon without ranges
	testSword := &models.Weapon{
		Equipment: models.Equipment{
			GameSystem: "midgard",
			Name:       "Kurzschwert",
		},
		SkillRequired: "Schwerter",
		Damage:        "1W6+1",
		RangeNear:     0,
		RangeMiddle:   0,
		RangeFar:      0,
	}
	_ = testSword.Create()

	// Arrange - Character with both ranged and melee weapons
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Archer",
		},
		Waffenfertigkeiten: []models.SkWaffenfertigkeit{
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{
							Name: "Bogen",
						},
					},
					Fertigkeitswert: 14,
				},
			},
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{
							Name: "Schwerter",
						},
					},
					Fertigkeitswert: 10,
				},
			},
		},
		Waffen: []models.EqWaffe{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Langbogen",
					},
				},
			},
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Kurzschwert",
					},
				},
			},
		},
		Eigenschaften: []models.Eigenschaft{
			{Name: "St", Value: 75},
			{Name: "Gs", Value: 85},
		},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Expect 3 weapons: Raufen + Langbogen + Kurzschwert
	if len(vm.Weapons) != 3 {
		t.Fatalf("Expected 3 weapons (Raufen + 2 equipped), got %d", len(vm.Weapons))
	}

	// Check first weapon is Raufen
	if vm.Weapons[0].Name != "Raufen" {
		t.Errorf("Expected first weapon to be 'Raufen', got '%s'", vm.Weapons[0].Name)
	}

	// Find the bow and sword
	var bow *WeaponViewModel
	var sword *WeaponViewModel
	for i := range vm.Weapons {
		if vm.Weapons[i].Name == "Langbogen" {
			bow = &vm.Weapons[i]
		} else if vm.Weapons[i].Name == "Kurzschwert" {
			sword = &vm.Weapons[i]
		}
	}

	// Test ranged weapon
	if bow == nil {
		t.Fatal("Langbogen not found in weapons")
	}

	if !bow.IsRanged {
		t.Error("Expected Langbogen to be marked as ranged weapon")
	}

	// Range should be formatted as "10/30/100" (Nah/Mittel/Fern)
	expectedRange := "10/30/100"
	if bow.Range != expectedRange {
		t.Errorf("Expected bow range '%s', got '%s'", expectedRange, bow.Range)
	}

	t.Logf("✓ Ranged weapon ranges: %s", bow.Range)

	// Test melee weapon
	if sword == nil {
		t.Fatal("Kurzschwert not found in weapons")
	}

	if sword.IsRanged {
		t.Error("Expected Kurzschwert not to be marked as ranged weapon")
	}

	if sword.Range != "" {
		t.Errorf("Expected melee weapon to have empty range, got '%s'", sword.Range)
	}
}
