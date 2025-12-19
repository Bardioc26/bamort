package pdfrender

import (
	"testing"
)

func TestPreparePaginatedPageData_Page1Stats(t *testing.T) {
	// Create test view model with many skills to test pagination
	viewModel := &CharacterSheetViewModel{
		Skills: make([]SkillViewModel, 50), // 50 skills should exceed column capacities
	}

	// Fill with test data
	for i := range viewModel.Skills {
		viewModel.Skills[i] = SkillViewModel{
			Name:  "Test Skill",
			Value: 10,
		}
	}

	pageData, err := PreparePaginatedPageData(viewModel, "page1_stats.html", 1, "2024-01-01")
	if err != nil {
		t.Fatalf("PreparePaginatedPageData failed: %v", err)
	}

	// Verify columns are populated
	if len(pageData.SkillsColumn1) == 0 {
		t.Error("SkillsColumn1 is empty")
	}

	if len(pageData.SkillsColumn2) == 0 {
		t.Error("SkillsColumn2 is empty")
	}

	// Get capacities from template
	templateSet := DefaultA4QuerTemplateSet()
	var page1Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page1_stats.html" {
			page1Template = &templateSet.Templates[i]
			break
		}
	}
	var col1MaxItems, col2MaxItems int
	for i := range page1Template.Metadata.Blocks {
		if page1Template.Metadata.Blocks[i].Name == "skills_column1" {
			col1MaxItems = page1Template.Metadata.Blocks[i].MaxItems
		} else if page1Template.Metadata.Blocks[i].Name == "skills_column2" {
			col2MaxItems = page1Template.Metadata.Blocks[i].MaxItems
		}
	}

	// Check capacities (from template)
	if len(pageData.SkillsColumn1) > col1MaxItems {
		t.Errorf("SkillsColumn1 exceeds capacity: got %d, max %d (from template)", len(pageData.SkillsColumn1), col1MaxItems)
	}

	if len(pageData.SkillsColumn2) > col2MaxItems {
		t.Errorf("SkillsColumn2 exceeds capacity: got %d, max %d (from template)", len(pageData.SkillsColumn2), col2MaxItems)
	}

	// Verify skills are split correctly
	totalPaginated := len(pageData.SkillsColumn1) + len(pageData.SkillsColumn2)
	expectedTotal := col1MaxItems + col2MaxItems
	if totalPaginated > expectedTotal {
		t.Errorf("Total paginated skills exceeds capacity: got %d, max %d (from template)", totalPaginated, expectedTotal)
	}

	t.Logf("Column1: %d skills, Column2: %d skills (total: %d)",
		len(pageData.SkillsColumn1), len(pageData.SkillsColumn2), totalPaginated)
}

func TestSplitSkillsForColumns(t *testing.T) {
	// Get actual column capacities from template
	templateSet := DefaultA4QuerTemplateSet()
	var page1Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page1_stats.html" {
			page1Template = &templateSet.Templates[i]
			break
		}
	}
	var col1MaxItems, col2MaxItems int
	for i := range page1Template.Metadata.Blocks {
		if page1Template.Metadata.Blocks[i].Name == "skills_column1" {
			col1MaxItems = page1Template.Metadata.Blocks[i].MaxItems
		} else if page1Template.Metadata.Blocks[i].Name == "skills_column2" {
			col2MaxItems = page1Template.Metadata.Blocks[i].MaxItems
		}
	}

	tests := []struct {
		name     string
		skills   int
		col1Max  int
		col2Max  int
		wantCol1 int
		wantCol2 int
	}{
		{
			name:     "few skills - only column 1",
			skills:   10,
			col1Max:  col1MaxItems,
			col2Max:  col2MaxItems,
			wantCol1: 10,
			wantCol2: 0,
		},
		{
			name:     "exactly column 1 capacity",
			skills:   col1MaxItems,
			col1Max:  col1MaxItems,
			col2Max:  col2MaxItems,
			wantCol1: col1MaxItems,
			wantCol2: 0,
		},
		{
			name:     "overflow to column 2",
			skills:   col1MaxItems + 11,
			col1Max:  col1MaxItems,
			col2Max:  col2MaxItems,
			wantCol1: col1MaxItems,
			wantCol2: 11,
		},
		{
			name:     "both columns full",
			skills:   col1MaxItems + col2MaxItems,
			col1Max:  col1MaxItems,
			col2Max:  col2MaxItems,
			wantCol1: col1MaxItems,
			wantCol2: col2MaxItems,
		},
		{
			name:     "more than both columns - truncate",
			skills:   col1MaxItems + col2MaxItems + 12,
			col1Max:  col1MaxItems,
			col2Max:  col2MaxItems,
			wantCol1: col1MaxItems,
			wantCol2: col2MaxItems,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test skills
			skills := make([]SkillViewModel, tt.skills)
			for i := range skills {
				skills[i] = SkillViewModel{
					Name:  "Test Skill",
					Value: 10,
				}
			}

			col1, col2 := SplitSkillsForColumns(skills, tt.col1Max, tt.col2Max)

			if len(col1) != tt.wantCol1 {
				t.Errorf("Column1: got %d skills, want %d", len(col1), tt.wantCol1)
			}

			if len(col2) != tt.wantCol2 {
				t.Errorf("Column2: got %d skills, want %d", len(col2), tt.wantCol2)
			}
		})
	}
}

func TestPreparePaginatedPageData_Page2Play(t *testing.T) {
	// Create 40 weapons to test capacity limiting
	viewModel := &CharacterSheetViewModel{
		Weapons: make([]WeaponViewModel, 40),
	}
	for i := range viewModel.Weapons {
		viewModel.Weapons[i] = WeaponViewModel{
			Name: "Test Weapon",
		}
	}

	pageData, err := PreparePaginatedPageData(viewModel, "page2_play.html", 2, "2024-01-01")
	if err != nil {
		t.Fatalf("PreparePaginatedPageData failed: %v", err)
	}

	// Page 2 should have weapons limited to 30
	if len(pageData.Weapons) > 30 {
		t.Errorf("Weapons exceed capacity: got %d, max 30", len(pageData.Weapons))
	}

	t.Logf("Page2: %d weapons", len(pageData.Weapons))
}

func TestPreparePaginatedPageData_Page3Spell(t *testing.T) {
	// Get capacities from template
	templateSet := DefaultA4QuerTemplateSet()
	leftCap := GetBlockCapacity(&templateSet, "page3_spell.html", "spells_left")
	rightCap := GetBlockCapacity(&templateSet, "page3_spell.html", "spells_right")
	magicItemsCap := GetBlockCapacity(&templateSet, "page3_spell.html", "magic_items")

	// Create test data exceeding capacities
	viewModel := &CharacterSheetViewModel{
		Spells:     make([]SpellViewModel, 50),
		MagicItems: make([]MagicItemViewModel, 20),
	}
	for i := range viewModel.Spells {
		viewModel.Spells[i] = SpellViewModel{
			Name: "Test Spell",
		}
	}
	for i := range viewModel.MagicItems {
		viewModel.MagicItems[i] = MagicItemViewModel{
			Name: "Test Item",
		}
	}

	pageData, err := PreparePaginatedPageData(viewModel, "page3_spell.html", 3, "2024-01-01")
	if err != nil {
		t.Fatalf("PreparePaginatedPageData failed: %v", err)
	}

	// Verify capacities match template
	if len(pageData.SpellsLeft) != leftCap {
		t.Errorf("SpellsLeft should be filled to %d (from template), got %d", leftCap, len(pageData.SpellsLeft))
	}

	if len(pageData.SpellsRight) != rightCap {
		t.Errorf("SpellsRight should be filled to %d (from template), got %d", rightCap, len(pageData.SpellsRight))
	}

	if len(pageData.MagicItems) != magicItemsCap {
		t.Errorf("MagicItems should be filled to %d (from template), got %d", magicItemsCap, len(pageData.MagicItems))
	}

	t.Logf("Page3: left=%d, right=%d (total=%d), magic_items=%d (from template)",
		len(pageData.SpellsLeft), len(pageData.SpellsRight),
		len(pageData.SpellsLeft)+len(pageData.SpellsRight), len(pageData.MagicItems))
}

func TestPreparePaginatedPageData_Page4Equipment(t *testing.T) {
	// Get capacity from template
	templateSet := DefaultA4QuerTemplateSet()
	equipmentCap := GetBlockCapacity(&templateSet, "page4_equip.html", "equipment_worn")

	// Create test data exceeding capacity
	viewModel := &CharacterSheetViewModel{
		Equipment: make([]EquipmentViewModel, 50),
	}
	for i := range viewModel.Equipment {
		viewModel.Equipment[i] = EquipmentViewModel{
			Name: "Test Equipment",
		}
	}

	pageData, err := PreparePaginatedPageData(viewModel, "page4_equip.html", 4, "2024-01-01")
	if err != nil {
		t.Fatalf("PreparePaginatedPageData failed: %v", err)
	}

	// Verify capacity matches template
	if len(pageData.Equipment) != equipmentCap {
		t.Errorf("Equipment should be filled to %d (from template), got %d", equipmentCap, len(pageData.Equipment))
	}

	t.Logf("Page4: %d equipment items (from template)", len(pageData.Equipment))
}
