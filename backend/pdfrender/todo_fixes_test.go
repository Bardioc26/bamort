package pdfrender

import (
	"bamort/database"
	"bamort/models"
	"os"
	"testing"
)

// TestPaginationUsesTemplateMetadata verifies tests use actual template MAX values
func TestPaginationUsesTemplateMetadata(t *testing.T) {
	// Read expected values directly from template file
	templateContent, err := os.ReadFile("../templates/Default_A4_Quer/page_2.html")
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	expectedBlocks := ParseTemplateMetadata(string(templateContent))
	expectedSkillsLearned := GetBlockByName(expectedBlocks, "skills_learned")
	expectedSkillsLanguages := GetBlockByName(expectedBlocks, "skills_languages")
	expectedWeaponsMain := GetBlockByName(expectedBlocks, "weapons_main")

	// Load template set from actual files
	templateSet := DefaultA4QuerTemplateSet()

	// Find page2
	var page2 *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_2.html" {
			page2 = &templateSet.Templates[i]
			break
		}
	}

	if page2 == nil {
		t.Fatal("page_2.html not found")
	}

	// Verify blocks exist and have correct MAX from template
	skillsLearned := GetBlockByName(page2.Metadata.Blocks, "skills_learned")
	if skillsLearned == nil {
		t.Fatal("skills_learned block not found")
	}
	if expectedSkillsLearned != nil && skillsLearned.MaxItems != expectedSkillsLearned.MaxItems {
		t.Errorf("skills_learned: expected MAX %d from template, got %d", expectedSkillsLearned.MaxItems, skillsLearned.MaxItems)
	}

	skillsLanguages := GetBlockByName(page2.Metadata.Blocks, "skills_languages")
	if skillsLanguages == nil {
		t.Fatal("skills_languages block not found")
	}
	if expectedSkillsLanguages != nil && skillsLanguages.MaxItems != expectedSkillsLanguages.MaxItems {
		t.Errorf("skills_languages: expected MAX %d from template, got %d", expectedSkillsLanguages.MaxItems, skillsLanguages.MaxItems)
	}

	weaponsMain := GetBlockByName(page2.Metadata.Blocks, "weapons_main")
	if weaponsMain == nil {
		t.Fatal("weapons_main block not found")
	}
	if expectedWeaponsMain != nil && weaponsMain.MaxItems != expectedWeaponsMain.MaxItems {
		t.Errorf("weapons_main: expected MAX %d from template, got %d", expectedWeaponsMain.MaxItems, weaponsMain.MaxItems)
	}
}

func TestPage2PaginationWithCorrectCapacities(t *testing.T) {
	// Read expected values directly from template file
	templateContent, err := os.ReadFile("../templates/Default_A4_Quer/page_2.html")
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	expectedBlocks := ParseTemplateMetadata(string(templateContent))
	expectedSkillsLearned := GetBlockByName(expectedBlocks, "skills_learned")
	expectedSkillsLanguages := GetBlockByName(expectedBlocks, "skills_languages")
	expectedWeaponsMain := GetBlockByName(expectedBlocks, "weapons_main")

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

	pageData, err := PreparePaginatedPageData(viewModel, "page_2.html", 2, "2024-01-01")
	if err != nil {
		t.Fatalf("Failed to prepare page data: %v", err)
	}

	// Verify capacities match template values
	if expectedSkillsLearned != nil && len(pageData.SkillsLearned) != expectedSkillsLearned.MaxItems {
		t.Errorf("SkillsLearned should be filled to %d (from template), got %d", expectedSkillsLearned.MaxItems, len(pageData.SkillsLearned))
	}

	if expectedSkillsLanguages != nil && len(pageData.SkillsLanguage) != expectedSkillsLanguages.MaxItems {
		t.Errorf("SkillsLanguage should be filled to %d (from template), got %d", expectedSkillsLanguages.MaxItems, len(pageData.SkillsLanguage))
	}

	if expectedWeaponsMain != nil && len(pageData.Weapons) != expectedWeaponsMain.MaxItems {
		t.Errorf("Weapons should be filled to %d (from template), got %d", expectedWeaponsMain.MaxItems, len(pageData.Weapons))
	}
}

func TestPage3MagicItemsCapacity(t *testing.T) {
	// Read expected values directly from template file
	templateContent, err := os.ReadFile("../templates/Default_A4_Quer/page_3.html")
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	expectedBlocks := ParseTemplateMetadata(string(templateContent))
	expectedMagicItems := GetBlockByName(expectedBlocks, "magic_items")

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

	pageData, err := PreparePaginatedPageData(viewModel, "page_3.html", 3, "2024-01-01")
	if err != nil {
		t.Fatalf("Failed to prepare page data: %v", err)
	}

	// Verify capacity matches template
	if expectedMagicItems != nil && len(pageData.MagicItems) != expectedMagicItems.MaxItems {
		t.Errorf("MagicItems should be filled to %d (from template), got %d", expectedMagicItems.MaxItems, len(pageData.MagicItems))
	}
}

func TestWeaponsWithEW(t *testing.T) {
	// Setup test database for weapon lookup
	database.SetupTestDB()

	// Create test weapon in gsm_weapons
	database.DB.Where("name = ?", "Schwert").Delete(&models.Weapon{})
	testWeapon := &models.Weapon{
		Equipment: models.Equipment{
			GameSystemId: 1,
			Name:         "Schwert",
		},
		SkillRequired: "Schwerter",
		Damage:        "1W6",
	}
	_ = testWeapon.Create()

	// Test that equipped weapons use Waffenfertigkeiten with correct EW
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
							Name: "Schwerter",
						},
					},
					Fertigkeitswert: 15,
					Category:        "Kampf",
				},
			},
		},
		Waffen: []models.EqWaffe{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{
						Name: "Schwert",
					},
				},
				Anb:  0,
				Schb: 0,
			},
		},
	}

	viewModel, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("Failed to map character: %v", err)
	}

	// Weapons should contain equipped weapons with EW from skill
	// Note: Raufen is always added as first weapon
	if len(viewModel.Weapons) < 2 {
		t.Fatalf("Expected at least 2 weapons (Raufen + Schwert), got %d", len(viewModel.Weapons))
	}

	// Find the Schwert weapon (not Raufen)
	var weapon *WeaponViewModel
	for i := range viewModel.Weapons {
		if viewModel.Weapons[i].Name == "Schwert" {
			weapon = &viewModel.Weapons[i]
			break
		}
	}

	if weapon == nil {
		t.Fatal("Schwert not found in weapons list")
	}

	if weapon.Name != "Schwert" {
		t.Errorf("Expected weapon name 'Schwert', got '%s'", weapon.Name)
	}
	if weapon.Value != 15 {
		t.Errorf("Expected weapon EW 15 (from skill), got %d", weapon.Value)
	}
}
