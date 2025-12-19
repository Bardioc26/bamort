package pdfrender

import (
	"bamort/models"
	"testing"
)

// TestPaginationUsesTemplateMetadata verifies tests use actual template MAX values
func TestPaginationUsesTemplateMetadata(t *testing.T) {
	// Load template set from actual files
	templateSet := DefaultA4QuerTemplateSet()

	// Find page2
	var page2 *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page2_play.html" {
			page2 = &templateSet.Templates[i]
			break
		}
	}

	if page2 == nil {
		t.Fatal("page2_play.html not found")
	}

	// Verify blocks exist and have correct MAX from template
	skillsLearned := GetBlockByName(page2.Metadata.Blocks, "skills_learned")
	if skillsLearned == nil {
		t.Fatal("skills_learned block not found")
	}
	if skillsLearned.MaxItems != 18 {
		t.Errorf("skills_learned: expected MAX 18 from template, got %d", skillsLearned.MaxItems)
	}

	skillsLanguages := GetBlockByName(page2.Metadata.Blocks, "skills_languages")
	if skillsLanguages == nil {
		t.Fatal("skills_languages block not found")
	}
	if skillsLanguages.MaxItems != 5 {
		t.Errorf("skills_languages: expected MAX 5 from template, got %d", skillsLanguages.MaxItems)
	}

	weaponsMain := GetBlockByName(page2.Metadata.Blocks, "weapons_main")
	if weaponsMain == nil {
		t.Fatal("weapons_main block not found")
	}
	if weaponsMain.MaxItems != 24 {
		t.Errorf("weapons_main: expected MAX 24 from template, got %d", weaponsMain.MaxItems)
	}
}

func TestPage2PaginationWithCorrectCapacities(t *testing.T) {
	// Create test data
	viewModel := &CharacterSheetViewModel{
		Skills: []SkillViewModel{
			{Name: "Learned 1", IsLearned: true, Category: "Combat"},
			{Name: "Learned 2", IsLearned: true, Category: "Combat"},
			{Name: "Sprache 1", Category: "Sprache"},
			{Name: "Sprache 2", Category: "Sprache"},
		},
		Weapons: []WeaponViewModel{
			{Name: "Sword", Value: 10},
			{Name: "Bow", Value: 12},
		},
	}

	pageData, err := PreparePaginatedPageData(viewModel, "page2_play.html", 2, "2024-01-01")
	if err != nil {
		t.Fatalf("Failed to prepare page data: %v", err)
	}

	// Verify capacities match template (18, 5, 24)
	if len(pageData.SkillsLearned) != 18 {
		t.Errorf("SkillsLearned should be filled to 18, got %d", len(pageData.SkillsLearned))
	}

	if len(pageData.SkillsLanguage) != 5 {
		t.Errorf("SkillsLanguage should be filled to 5, got %d", len(pageData.SkillsLanguage))
	}

	if len(pageData.Weapons) != 24 {
		t.Errorf("Weapons should be filled to 24, got %d", len(pageData.Weapons))
	}
}

func TestPage3MagicItemsCapacity(t *testing.T) {
	// Create test data with magic items
	viewModel := &CharacterSheetViewModel{
		MagicItems: []MagicItemViewModel{
			{Name: "Wand"},
			{Name: "Ring"},
		},
		Spells: []SpellViewModel{
			{Name: "Fireball"},
		},
	}

	pageData, err := PreparePaginatedPageData(viewModel, "page3_spell.html", 3, "2024-01-01")
	if err != nil {
		t.Fatalf("Failed to prepare page data: %v", err)
	}

	// Template says MAX: 8 for magic_items
	if len(pageData.MagicItems) != 8 {
		t.Errorf("MagicItems should be filled to 8, got %d", len(pageData.MagicItems))
	}
}

func TestWeaponsWithEW(t *testing.T) {
	// Test that weapons use Waffenfertigkeiten with correct EW
	char := &models.Char{
		BamortBase: models.BamortBase{
			ID:   1,
			Name: "Test Fighter",
		},
		Waffenfertigkeiten: []models.SkWaffenfertigkeit{
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{
							Name: "Schwert",
						},
					},
					Fertigkeitswert: 15,
					Category:        "Kampf",
				},
			},
		},
	}

	viewModel, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("Failed to map character: %v", err)
	}

	// Weapons should contain weapon skills with EW
	if len(viewModel.Weapons) == 0 {
		t.Fatal("Expected weapons from Waffenfertigkeiten, got none")
	}

	weapon := viewModel.Weapons[0]
	if weapon.Name != "Schwert" {
		t.Errorf("Expected weapon name 'Schwert', got '%s'", weapon.Name)
	}
	if weapon.Value != 15 {
		t.Errorf("Expected weapon EW 15, got %d", weapon.Value)
	}
}
