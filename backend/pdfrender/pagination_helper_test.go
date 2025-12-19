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

	// Check capacities (MAX: 29 each)
	if len(pageData.SkillsColumn1) > 29 {
		t.Errorf("SkillsColumn1 exceeds capacity: got %d, max 29", len(pageData.SkillsColumn1))
	}

	if len(pageData.SkillsColumn2) > 29 {
		t.Errorf("SkillsColumn2 exceeds capacity: got %d, max 29", len(pageData.SkillsColumn2))
	}

	// Verify skills are split correctly
	totalPaginated := len(pageData.SkillsColumn1) + len(pageData.SkillsColumn2)
	expectedTotal := 58 // 29 + 29
	if totalPaginated > expectedTotal {
		t.Errorf("Total paginated skills exceeds capacity: got %d, max %d", totalPaginated, expectedTotal)
	}

	t.Logf("Column1: %d skills, Column2: %d skills (total: %d)",
		len(pageData.SkillsColumn1), len(pageData.SkillsColumn2), totalPaginated)
}

func TestSplitSkillsForColumns(t *testing.T) {
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
			col1Max:  29,
			col2Max:  29,
			wantCol1: 10,
			wantCol2: 0,
		},
		{
			name:     "exactly column 1 capacity",
			skills:   29,
			col1Max:  29,
			col2Max:  29,
			wantCol1: 29,
			wantCol2: 0,
		},
		{
			name:     "overflow to column 2",
			skills:   40,
			col1Max:  29,
			col2Max:  29,
			wantCol1: 29,
			wantCol2: 11,
		},
		{
			name:     "both columns full",
			skills:   58,
			col1Max:  29,
			col2Max:  29,
			wantCol1: 29,
			wantCol2: 29,
		},
		{
			name:     "more than both columns - truncate",
			skills:   70,
			col1Max:  29,
			col2Max:  29,
			wantCol1: 29,
			wantCol2: 29,
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
	// Create 30 spells and 10 magic items to test capacity
	viewModel := &CharacterSheetViewModel{
		Spells:     make([]SpellViewModel, 30),
		MagicItems: make([]MagicItemViewModel, 10),
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

	// Page 3 should have spells limited to 24 (12+12)
	if len(pageData.Spells) > 24 {
		t.Errorf("Spells exceed capacity: got %d, max 24", len(pageData.Spells))
	}

	// Magic items limited to 5
	if len(pageData.MagicItems) > 5 {
		t.Errorf("MagicItems exceed capacity: got %d, max 5", len(pageData.MagicItems))
	}

	t.Logf("Page3: %d spells, %d magic items", len(pageData.Spells), len(pageData.MagicItems))
}

func TestPreparePaginatedPageData_Page4Equipment(t *testing.T) {
	// Create 30 equipment items to test capacity
	viewModel := &CharacterSheetViewModel{
		Equipment: make([]EquipmentViewModel, 30),
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

	// Page 4 should have equipment limited to 20
	if len(pageData.Equipment) > 20 {
		t.Errorf("Equipment exceeds capacity: got %d, max 20", len(pageData.Equipment))
	}

	t.Logf("Page4: %d equipment items", len(pageData.Equipment))
}
