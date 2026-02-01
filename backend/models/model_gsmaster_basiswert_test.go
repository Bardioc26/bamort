package models

import (
	"testing"
)

func TestSkillBasisWert(t *testing.T) {
	// Test that BasisWert defaults to 0 in master skill data
	skill := Skill{
		Name:         "TestSkill",
		GameSystemId: 1,
		Initialwert:  5,
	}

	// BasisWert should default to 0
	if skill.BasisWert != 0 {
		t.Errorf("Expected BasisWert to default to 0, got %d", skill.BasisWert)
	}
}

func TestSkillBasisWertSet(t *testing.T) {
	// Test that BasisWert can be set in master skill data
	skill := Skill{
		Name:         "TestSkill",
		GameSystemId: 1,
		Initialwert:  5,
		BasisWert:    3,
	}

	if skill.BasisWert != 3 {
		t.Errorf("Expected BasisWert to be 3, got %d", skill.BasisWert)
	}
}

func TestWeaponSkillBasisWert(t *testing.T) {
	// Test that BasisWert works for WeaponSkill (inherited from Skill)
	weaponSkill := WeaponSkill{
		Skill: Skill{
			Name:         "TestWeaponSkill",
			GameSystemId: 1,
			Initialwert:  5,
			BasisWert:    2,
		},
	}

	if weaponSkill.BasisWert != 2 {
		t.Errorf("Expected BasisWert to be 2, got %d", weaponSkill.BasisWert)
	}
}
