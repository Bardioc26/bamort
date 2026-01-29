package pdfrender

import (
	"bamort/database"
	"bamort/models"
	"testing"
)

func TestMapCharacterToViewModel_BasicInfo(t *testing.T) {
	// Arrange
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Bjarnfinnur Haberdson",
		},
		Rasse:       "Mensch",
		Typ:         "Krieger",
		Grad:        5,
		Alter:       30,
		Gender:      "männlich",
		Groesse:     180,
		Gewicht:     85,
		Herkunft:    "Clanngadarn",
		Glaube:      "Druide",
		SocialClass: "Mittelschicht",
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if vm.Character.Name != "Bjarnfinnur Haberdson" {
		t.Errorf("Expected name 'Bjarnfinnur Haberdson', got '%s'", vm.Character.Name)
	}

	if vm.Character.Type != "Krieger" {
		t.Errorf("Expected type 'Krieger', got '%s'", vm.Character.Type)
	}

	if vm.Character.Grade != 5 {
		t.Errorf("Expected grade 5, got %d", vm.Character.Grade)
	}

	if vm.Character.Herkunft != "Clanngadarn" {
		t.Errorf("Expected Herkunft 'Clanngadarn', got '%s'", vm.Character.Herkunft)
	}

	if vm.Character.Glaube != "Druide" {
		t.Errorf("Expected Glaube 'Druide', got '%s'", vm.Character.Glaube)
	}

	if vm.Character.SocialClass != "Mittelschicht" {
		t.Errorf("Expected SocialClass 'Mittelschicht', got '%s'", vm.Character.SocialClass)
	}
}

func TestMapCharacterToViewModel_Attributes(t *testing.T) {
	// Arrange
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Character",
		},
		Eigenschaften: []models.Eigenschaft{
			{Name: "St", Value: 79},
			{Name: "Gs", Value: 65},
			{Name: "Gw", Value: 70},
		},
		B: models.B{Max: 10},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if vm.Attributes.St != 79 {
		t.Errorf("Expected St=79, got %d", vm.Attributes.St)
	}
	if vm.Attributes.Gs != 65 {
		t.Errorf("Expected Gs=65, got %d", vm.Attributes.Gs)
	}
	if vm.Attributes.B != 10 {
		t.Errorf("Expected B=10, got %d", vm.Attributes.B)
	}
}

func TestMapCharacterToViewModel_LPandAP(t *testing.T) {
	// Arrange
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Character",
		},
		Lp: models.Lp{Max: 20, Value: 18},
		Ap: models.Ap{Max: 30, Value: 25},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if vm.DerivedValues.LPMax != 20 {
		t.Errorf("Expected LP Max=20, got %d", vm.DerivedValues.LPMax)
	}

	if vm.DerivedValues.LPAktuell != 18 {
		t.Errorf("Expected LP Aktuell=18, got %d", vm.DerivedValues.LPAktuell)
	}

	if vm.DerivedValues.APMax != 30 {
		t.Errorf("Expected AP Max=30, got %d", vm.DerivedValues.APMax)
	}

	if vm.DerivedValues.APAktuell != 25 {
		t.Errorf("Expected AP Aktuell=25, got %d", vm.DerivedValues.APAktuell)
	}
}

func TestMapCharacterToViewModel_Skills(t *testing.T) {
	// Arrange
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Character",
		},
		Fertigkeiten: []models.SkFertigkeit{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Schwimmen",
					},
				},
				Fertigkeitswert: 10,
				Bonus:           2,
				Pp:              5,
				Category:        "Körper",
			},
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Klettern",
					},
				},
				Fertigkeitswert: 8,
				Bonus:           1,
				Pp:              3,
				Category:        "Körper",
			},
		},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(vm.Skills) != 2 {
		t.Fatalf("Expected 2 skills, got %d", len(vm.Skills))
	}

	skill := vm.Skills[0]
	if skill.Name != "Schwimmen" {
		t.Errorf("Expected skill name 'Schwimmen', got '%s'", skill.Name)
	}
	if skill.Value != 10 {
		t.Errorf("Expected skill value 10, got %d", skill.Value)
	}
	if skill.Bonus != 2 {
		t.Errorf("Expected skill bonus 2, got %d", skill.Bonus)
	}
	if skill.PracticePoints != 5 {
		t.Errorf("Expected practice points 5, got %d", skill.PracticePoints)
	}
	if skill.Category != "Körper" {
		t.Errorf("Expected category 'Körper', got '%s'", skill.Category)
	}
}

func TestMapCharacterToViewModel_Weapons(t *testing.T) {
	// Setup test database for weapon lookup
	database.SetupTestDB()

	// Create test weapon in gsm_weapons
	database.DB.Where("name = ?", "Langschwert").Delete(&models.Weapon{})
	testWeapon := &models.Weapon{
		Equipment: models.Equipment{
			GameSystemId: 1,
			Name:         "Langschwert",
		},
		SkillRequired: "Schwerter",
		Damage:        "1W6",
	}
	_ = testWeapon.Create()

	// Arrange
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Character",
		},
		Waffenfertigkeiten: []models.SkWaffenfertigkeit{
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{
							Name: "Schwerter",
						},
					},
					Fertigkeitswert: 12,
					Bonus:           3,
					Category:        "Kampf",
				},
			},
		},
		Waffen: []models.EqWaffe{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Langschwert",
					},
				},
				Anb:  0,
				Schb: 0,
			},
		},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Expect 2 weapons: Raufen (always first) + Langschwert
	if len(vm.Weapons) != 2 {
		t.Fatalf("Expected 2 weapons (Raufen + Langschwert), got %d", len(vm.Weapons))
	}

	// Check first weapon is Raufen
	if vm.Weapons[0].Name != "Raufen" {
		t.Errorf("Expected first weapon to be 'Raufen', got '%s'", vm.Weapons[0].Name)
	}

	// Check second weapon is Langschwert
	weapon := vm.Weapons[1]
	if weapon.Name != "Langschwert" {
		t.Errorf("Expected weapon name 'Langschwert', got '%s'", weapon.Name)
	}
	if weapon.Value != 12 {
		t.Errorf("Expected weapon value 12 (skill value), got %d", weapon.Value)
	}
}

func TestMapCharacterToViewModel_Spells(t *testing.T) {
	// Arrange
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Character",
		},
		Zauber: []models.SkZauber{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Macht über die belebte Natur",
					},
				},
				Bonus:  2,
				Quelle: "Arkanum",
			},
		},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(vm.Spells) != 1 {
		t.Fatalf("Expected 1 spell, got %d", len(vm.Spells))
	}

	spell := vm.Spells[0]
	if spell.Name != "Macht über die belebte Natur" {
		t.Errorf("Expected spell name 'Macht über die belebte Natur', got '%s'", spell.Name)
	}
}

func TestMapCharacterToViewModel_Equipment(t *testing.T) {
	// Arrange
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Character",
		},
		Ausruestung: []models.EqAusruestung{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Rucksack",
					},
				},
				Anzahl:       1,
				Gewicht:      2.5,
				Wert:         10,
				BeinhaltetIn: "Am Körper",
			},
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Seil (20m)",
					},
				},
				Anzahl:       1,
				Gewicht:      5.0,
				Wert:         15,
				BeinhaltetIn: "Rucksack",
			},
		},
		Behaeltnisse: []models.EqContainer{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Lederbeutel",
					},
				},
				Gewicht:   0.5,
				Wert:      5,
				Tragkraft: 10,
			},
		},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(vm.Equipment) != 3 {
		t.Fatalf("Expected 3 equipment items, got %d", len(vm.Equipment))
	}

	// Check first item
	item := vm.Equipment[0]
	if item.Name != "Rucksack" {
		t.Errorf("Expected equipment name 'Rucksack', got '%s'", item.Name)
	}
	if item.Quantity != 1 {
		t.Errorf("Expected quantity 1, got %d", item.Quantity)
	}
	if item.Weight != 2.5 {
		t.Errorf("Expected weight 2.5, got %f", item.Weight)
	}
	if item.Location != "Am Körper" {
		t.Errorf("Expected location 'Am Körper', got '%s'", item.Location)
	}

	// Check container
	container := vm.Equipment[2]
	if container.Name != "Lederbeutel" {
		t.Errorf("Expected container name 'Lederbeutel', got '%s'", container.Name)
	}
	if !container.IsContainer {
		t.Error("Expected IsContainer to be true")
	}
}

func TestMapCharacterToViewModel_GameResults(t *testing.T) {
	// Arrange - no game results in domain model yet
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Character",
		},
	}

	// Act
	vm, err := MapCharacterToViewModel(char)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should be empty slice, not nil
	if vm.GameResults == nil {
		t.Error("Expected GameResults to be empty slice, got nil")
	}

	if len(vm.GameResults) != 0 {
		t.Errorf("Expected 0 game results, got %d", len(vm.GameResults))
	}
}
