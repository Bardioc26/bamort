package pdfrender

import (
	"testing"
)

func TestSliceList_Basic(t *testing.T) {
	// Arrange
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Act
	result, hasMore := SliceList(items, 0, 5)

	// Assert
	if len(result) != 5 {
		t.Errorf("Expected 5 items, got %d", len(result))
	}
	if !hasMore {
		t.Error("Expected hasMore to be true")
	}
	if result[0] != 1 || result[4] != 5 {
		t.Error("Unexpected slice content")
	}
}

func TestSliceList_LastPage(t *testing.T) {
	// Arrange
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Act
	result, hasMore := SliceList(items, 5, 10)

	// Assert
	if len(result) != 5 {
		t.Errorf("Expected 5 items, got %d", len(result))
	}
	if hasMore {
		t.Error("Expected hasMore to be false")
	}
}

func TestSliceList_BeyondEnd(t *testing.T) {
	// Arrange
	items := []int{1, 2, 3}

	// Act
	result, hasMore := SliceList(items, 10, 5)

	// Assert
	if len(result) != 0 {
		t.Errorf("Expected 0 items, got %d", len(result))
	}
	if hasMore {
		t.Error("Expected hasMore to be false")
	}
}

func TestPaginateSkills_SinglePage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	skills := make([]SkillViewModel, 10)
	for i := 0; i < 10; i++ {
		skills[i] = SkillViewModel{Name: "Skill" + string(rune('A'+i))}
	}

	// Act - page1_stats has 2 columns with 32 each = 64 total capacity
	pages, err := paginator.PaginateSkills(skills, "page1_stats.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(pages))
	}

	// Check that data is distributed across columns
	page := pages[0]
	if page.PageNumber != 1 {
		t.Errorf("Expected page number 1, got %d", page.PageNumber)
	}

	// Column 1 should have all 10 skills (max 32)
	col1Data, ok := page.Data["skills_column1"].([]SkillViewModel)
	if !ok {
		t.Fatal("skills_column1 data not found or wrong type")
	}
	if len(col1Data) != 10 {
		t.Errorf("Expected 10 skills in column 1, got %d", len(col1Data))
	}

	// Column 2 should be empty (no overflow)
	col2Data, ok := page.Data["skills_column2"].([]SkillViewModel)
	if !ok {
		t.Fatal("skills_column2 data not found or wrong type")
	}
	if len(col2Data) != 0 {
		t.Errorf("Expected 0 skills in column 2, got %d", len(col2Data))
	}
}

func TestPaginateSkills_MultiColumn(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create 40 skills - should fill first column (32) and spill to second (8)
	skills := make([]SkillViewModel, 40)
	for i := 0; i < 40; i++ {
		skills[i] = SkillViewModel{Name: "Skill" + string(rune(i))}
	}

	// Act
	pages, err := paginator.PaginateSkills(skills, "page1_stats.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(pages))
	}

	page := pages[0]

	// Column 1 should have 29 skills
	col1Data := page.Data["skills_column1"].([]SkillViewModel)
	if len(col1Data) != 29 {
		t.Errorf("Expected 29 skills in column 1, got %d", len(col1Data))
	}

	// Column 2 should have 11 skills (40 total - 29 in col1)
	col2Data := page.Data["skills_column2"].([]SkillViewModel)
	if len(col2Data) != 11 {
		t.Errorf("Expected 11 skills in column 2, got %d", len(col2Data))
	}
}

func TestPaginateSkills_MultiPage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create 100 skills - should span 2 pages (64 capacity per page)
	skills := make([]SkillViewModel, 100)
	for i := 0; i < 100; i++ {
		skills[i] = SkillViewModel{Name: "Skill" + string(rune(i))}
	}

	// Act
	pages, err := paginator.PaginateSkills(skills, "page1_stats.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("Expected 2 pages, got %d", len(pages))
	}

	// Page 1 should have 58 skills (29 + 29)
	page1 := pages[0]
	col1Page1 := page1.Data["skills_column1"].([]SkillViewModel)
	col2Page1 := page1.Data["skills_column2"].([]SkillViewModel)
	if len(col1Page1) != 29 {
		t.Errorf("Page 1 col1: expected 29 skills, got %d", len(col1Page1))
	}
	if len(col2Page1) != 29 {
		t.Errorf("Page 1 col2: expected 29 skills, got %d", len(col2Page1))
	}

	// Page 2 should have 42 skills (29 + 13) - 100 total - 58 from page 1 = 42 remaining
	page2 := pages[1]
	col1Page2 := page2.Data["skills_column1"].([]SkillViewModel)
	col2Page2 := page2.Data["skills_column2"].([]SkillViewModel)
	if len(col1Page2) != 29 {
		t.Errorf("Page 2 col1: expected 29 skills, got %d", len(col1Page2))
	}
	if len(col2Page2) != 13 {
		t.Errorf("Page 2 col2: expected 13 skills, got %d", len(col2Page2))
	}
}

func TestPaginateSpells_TwoColumns(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create 15 spells - should fit in first column (12) with 3 in second column (12)
	spells := make([]SpellViewModel, 15)
	for i := 0; i < 15; i++ {
		spells[i] = SpellViewModel{Name: "Spell" + string(rune('A'+i))}
	}

	// Act
	pages, err := paginator.PaginateSpells(spells, "page3_spell.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(pages))
	}

	page := pages[0]

	// Left column should have 15 spells (all fit in 20 capacity)
	leftData := page.Data["spells_left"].([]SpellViewModel)
	if len(leftData) != 15 {
		t.Errorf("Expected 15 spells in left column, got %d", len(leftData))
	}

	// Right column should be empty (15 spells all fit in left)
	rightData := page.Data["spells_right"].([]SpellViewModel)
	if len(rightData) != 0 {
		t.Errorf("Expected 0 spells in right column, got %d", len(rightData))
	}
}

func TestPaginateSpells_MultiPage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create 30 spells - should fit on 1 page (30 capacity = 20 left + 10 right)
	spells := make([]SpellViewModel, 30)
	for i := 0; i < 30; i++ {
		spells[i] = SpellViewModel{Name: "Spell" + string(rune(i))}
	}

	// Act
	pages, err := paginator.PaginateSpells(spells, "page3_spell.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// With capacity of 20+10=30, all 30 spells fit on 1 page
	if len(pages) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(pages))
	}

	// Page 1 should have all 30 spells (20 left + 10 right)
	page1 := pages[0]
	leftPage1 := page1.Data["spells_left"].([]SpellViewModel)
	rightPage1 := page1.Data["spells_right"].([]SpellViewModel)
	if len(leftPage1) != 20 {
		t.Errorf("Page 1 left: expected 20 spells, got %d", len(leftPage1))
	}
	if len(rightPage1) != 10 {
		t.Errorf("Page 1 right: expected 10 spells, got %d", len(rightPage1))
	}
}

func TestPaginateWeapons_SinglePage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	weapons := make([]WeaponViewModel, 10)
	for i := 0; i < 10; i++ {
		weapons[i] = WeaponViewModel{Name: "Weapon" + string(rune('A'+i))}
	}

	// Act - page2_play has weapons_main with MAX:30
	pages, err := paginator.PaginateWeapons(weapons, "page2_play.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(pages))
	}

	page := pages[0]
	weaponsData := page.Data["weapons_main"].([]WeaponViewModel)
	if len(weaponsData) != 10 {
		t.Errorf("Expected 10 weapons, got %d", len(weaponsData))
	}
}

func TestPaginateWeapons_MultiPage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create 50 weapons - should span 2 pages (30 capacity per page)
	weapons := make([]WeaponViewModel, 50)
	for i := 0; i < 50; i++ {
		weapons[i] = WeaponViewModel{Name: "Weapon" + string(rune(i))}
	}

	// Act
	pages, err := paginator.PaginateWeapons(weapons, "page2_play.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 2 {
		t.Fatalf("Expected 2 pages, got %d", len(pages))
	}

	// Page 1 should have 30 weapons
	page1Weapons := pages[0].Data["weapons_main"].([]WeaponViewModel)
	if len(page1Weapons) != 30 {
		t.Errorf("Page 1: expected 30 weapons, got %d", len(page1Weapons))
	}

	// Page 2 should have 20 weapons
	page2Weapons := pages[1].Data["weapons_main"].([]WeaponViewModel)
	if len(page2Weapons) != 20 {
		t.Errorf("Page 2: expected 20 weapons, got %d", len(page2Weapons))
	}
}

func TestCalculatePagesNeeded(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	testCases := []struct {
		name          string
		templateName  string
		listType      string
		itemCount     int
		expectedPages int
	}{
		{"10 skills on page1", "page1_stats.html", "skills", 10, 1},
		{"58 skills on page1", "page1_stats.html", "skills", 58, 1}, // 29+29 = 58 fits on 1 page
		{"59 skills on page1", "page1_stats.html", "skills", 59, 2}, // 59 requires 2 pages
		{"100 skills on page1", "page1_stats.html", "skills", 100, 2},
		{"10 weapons on page2", "page2_play.html", "weapons", 10, 1},
		{"30 weapons on page2", "page2_play.html", "weapons", 30, 1},
		{"31 weapons on page2", "page2_play.html", "weapons", 31, 2},
		{"10 spells on page3", "page3_spell.html", "spells", 10, 1},
		{"30 spells on page3", "page3_spell.html", "spells", 30, 1}, // 20+10 = 30 fits on 1 page
		{"31 spells on page3", "page3_spell.html", "spells", 31, 2}, // 31 requires 2 pages
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			pages, err := paginator.CalculatePagesNeeded(tc.templateName, tc.listType, tc.itemCount)

			// Assert
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if pages != tc.expectedPages {
				t.Errorf("Expected %d pages, got %d", tc.expectedPages, pages)
			}
		})
	}
}

func TestPaginateSkills_EmptyList(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	skills := []SkillViewModel{}

	// Act
	pages, err := paginator.PaginateSkills(skills, "page1_stats.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(pages) != 0 {
		t.Errorf("Expected 0 pages for empty list, got %d", len(pages))
	}
}

func TestPaginateSkills_InvalidTemplate(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	skills := []SkillViewModel{{Name: "Test"}}

	// Act
	_, err := paginator.PaginateSkills(skills, "nonexistent.html", "")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid template, got nil")
	}
}
