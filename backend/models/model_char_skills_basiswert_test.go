package models

import (
	"testing"
)

func TestSkFertigkeitBasisWert(t *testing.T) {
	// Test that BasisWert defaults to 0
	skill := SkFertigkeit{
		Fertigkeitswert: 10,
	}
	skill.Name = "TestSkill"

	// BasisWert should default to 0
	if skill.BasisWert != 0 {
		t.Errorf("Expected BasisWert to default to 0, got %d", skill.BasisWert)
	}
}

func TestSkFertigkeitBasisWertSet(t *testing.T) {
	// Test that BasisWert can be set
	skill := SkFertigkeit{
		Fertigkeitswert: 10,
		BasisWert:       5,
	}
	skill.Name = "TestSkill"

	if skill.BasisWert != 5 {
		t.Errorf("Expected BasisWert to be 5, got %d", skill.BasisWert)
	}
}

func TestSkWaffenfertigkeitBasisWert(t *testing.T) {
	// Test that BasisWert works for SkWaffenfertigkeit (inherited from SkFertigkeit)
	weaponSkill := SkWaffenfertigkeit{
		SkFertigkeit: SkFertigkeit{
			Fertigkeitswert: 8,
			BasisWert:       3,
		},
	}
	weaponSkill.Name = "TestWeaponSkill"

	if weaponSkill.BasisWert != 3 {
		t.Errorf("Expected BasisWert to be 3, got %d", weaponSkill.BasisWert)
	}
}

func TestSkAngeboreneFertigkeitBasisWert(t *testing.T) {
	// Test that BasisWert works for SkAngeboreneFertigkeit (inherited from SkFertigkeit)
	innateSkill := SkAngeboreneFertigkeit{
		SkFertigkeit: SkFertigkeit{
			Fertigkeitswert: 12,
			BasisWert:       7,
		},
	}
	innateSkill.Name = "TestInnateSkill"

	if innateSkill.BasisWert != 7 {
		t.Errorf("Expected BasisWert to be 7, got %d", innateSkill.BasisWert)
	}
}
