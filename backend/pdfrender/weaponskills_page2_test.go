package pdfrender

import (
	"testing"

	"bamort/models"
)

func TestWeaponSkillsAppearInSkillsList(t *testing.T) {
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Fighter",
		},
		Fertigkeiten: []models.SkFertigkeit{
			{
				BamortCharTrait: models.BamortCharTrait{BamortBase: models.BamortBase{Name: "Klettern"}},
				Fertigkeitswert: 10,
				Category:        "Körper",
			},
		},
		Waffenfertigkeiten: []models.SkWaffenfertigkeit{
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{BamortBase: models.BamortBase{Name: "Langschwert"}},
					Fertigkeitswert: 12,
				},
			},
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{BamortBase: models.BamortBase{Name: "Dolch"}},
					Fertigkeitswert: 8,
				},
			},
		},
	}

	vm, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("MapCharacterToViewModel failed: %v", err)
	}

	// Should have 3 skills: 1 regular + 2 weapon skills
	expectedSkillCount := 3
	if len(vm.Skills) != expectedSkillCount {
		t.Errorf("Expected %d skills, got %d", expectedSkillCount, len(vm.Skills))
		for i, s := range vm.Skills {
			t.Logf("  Skill %d: %s (Category: %s)", i, s.Name, s.Category)
		}
	}

	// Check that weapon skills are present
	foundLangschwert := false
	foundDolch := false
	for _, skill := range vm.Skills {
		if skill.Name == "Langschwert" && skill.Category == "Waffenfertigkeit" {
			foundLangschwert = true
			if skill.Value != 12 {
				t.Errorf("Langschwert value: expected 12, got %d", skill.Value)
			}
		}
		if skill.Name == "Dolch" && skill.Category == "Waffenfertigkeit" {
			foundDolch = true
			if skill.Value != 8 {
				t.Errorf("Dolch value: expected 8, got %d", skill.Value)
			}
		}
	}

	if !foundLangschwert {
		t.Error("Langschwert weapon skill not found in Skills list")
	}
	if !foundDolch {
		t.Error("Dolch weapon skill not found in Skills list")
	}
}

func TestWeaponSkillsAppearOnPage2(t *testing.T) {
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Fighter",
		},
		Fertigkeiten: []models.SkFertigkeit{
			{
				BamortCharTrait: models.BamortCharTrait{BamortBase: models.BamortBase{Name: "Klettern"}},
				Fertigkeitswert: 10,
				Category:        "Körper",
			},
		},
		Waffenfertigkeiten: []models.SkWaffenfertigkeit{
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{BamortBase: models.BamortBase{Name: "Langschwert"}},
					Fertigkeitswert: 12,
				},
			},
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{BamortBase: models.BamortBase{Name: "Dolch"}},
					Fertigkeitswert: 8,
				},
			},
		},
	}

	vm, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("MapCharacterToViewModel failed: %v", err)
	}

	pageData, err := PreparePaginatedPageData(vm, "page_2.html", 2, "2024-01-01")
	if err != nil {
		t.Fatalf("PreparePaginatedPageData failed: %v", err)
	}

	// SkillsLearned should contain both regular skills and weapon skills
	// Expected: Klettern + Langschwert + Dolch = 3 learned skills
	expectedLearnedCount := 3
	if len(pageData.SkillsLearned) < expectedLearnedCount {
		t.Errorf("Expected at least %d learned skills on page 2, got %d", expectedLearnedCount, len(pageData.SkillsLearned))
		t.Log("SkillsLearned contents:")
		for i, s := range pageData.SkillsLearned {
			t.Logf("  %d: %s (Category: %s, Value: %d)", i, s.Name, s.Category, s.Value)
		}
	}

	// Verify weapon skills are in SkillsLearned
	foundLangschwert := false
	foundDolch := false
	for _, skill := range pageData.SkillsLearned {
		if skill.Name == "Langschwert" {
			foundLangschwert = true
		}
		if skill.Name == "Dolch" {
			foundDolch = true
		}
	}

	if !foundLangschwert {
		t.Error("Langschwert weapon skill not found in page 2 SkillsLearned")
	}
	if !foundDolch {
		t.Error("Dolch weapon skill not found in page 2 SkillsLearned")
	}
}
