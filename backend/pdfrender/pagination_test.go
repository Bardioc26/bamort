package pdfrender

import (
	"fmt"
	"testing"
)

// TestPaginateMultiList_SingleListType tests pagination with a single list type (skills only)
func TestPaginateMultiList_SingleListType(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create skills that will fit on one page (less than total capacity)
	skills := make([]SkillViewModel, 20)
	for i := 0; i < 20; i++ {
		skills[i] = SkillViewModel{Name: fmt.Sprintf("Skill%d", i+1)}
	}

	dataMap := map[string]interface{}{
		"skills": skills,
	}

	// Act
	distributions, err := paginator.PaginateMultiList(dataMap, "page_1.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(distributions) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(distributions))
	}

	// Verify data is distributed correctly
	page := distributions[0]
	if page.TemplateName != "page_1.html" {
		t.Errorf("Expected template 'page_1.html', got '%s'", page.TemplateName)
	}

	// Check that skills are distributed to appropriate blocks
	// Note: Skills are distributed to multiple columns, so we count each block
	totalSkills := 0
	skillBlocks := []string{"skills_column1", "skills_column2"}
	for _, blockName := range skillBlocks {
		if data, exists := page.Data[blockName]; exists {
			if skillsList, ok := data.([]SkillViewModel); ok {
				totalSkills += len(skillsList)
			}
		}
	}

	if totalSkills != 20 {
		t.Errorf("Expected 20 total skills distributed across columns, got %d", totalSkills)
	}
}

// TestPaginateMultiList_MultipleListTypes tests pagination with multiple list types
func TestPaginateMultiList_MultipleListTypes(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create skills and weapons for page_2.html
	skills := make([]SkillViewModel, 30)
	for i := 0; i < 30; i++ {
		skills[i] = SkillViewModel{
			Name:      fmt.Sprintf("Skill%d", i+1),
			IsLearned: i < 15, // First 15 are learned
			Category:  "Kampf",
		}
	}

	weapons := make([]WeaponViewModel, 10)
	for i := 0; i < 10; i++ {
		weapons[i] = WeaponViewModel{
			Name:  fmt.Sprintf("Weapon%d", i+1),
			Value: 10 + i,
		}
	}

	dataMap := map[string]interface{}{
		"skills":  skills,
		"weapons": weapons,
	}

	// Act
	distributions, err := paginator.PaginateMultiList(dataMap, "page_2.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(distributions) == 0 {
		t.Fatal("Expected at least 1 page")
	}

	// Verify all data is distributed across pages
	totalSkills := 0
	totalWeapons := 0

	for _, dist := range distributions {
		for _, data := range dist.Data {
			switch v := data.(type) {
			case []SkillViewModel:
				totalSkills += len(v)
			case []WeaponViewModel:
				totalWeapons += len(v)
			}
		}
	}

	if totalSkills != 30 {
		t.Errorf("Expected 30 total skills, got %d", totalSkills)
	}
	if totalWeapons != 10 {
		t.Errorf("Expected 10 total weapons, got %d", totalWeapons)
	}
}

// TestPaginateMultiList_WithOverflow tests continuation pages are created
func TestPaginateMultiList_WithOverflow(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create many skills to force continuation pages
	skills := make([]SkillViewModel, 100)
	for i := 0; i < 100; i++ {
		skills[i] = SkillViewModel{Name: fmt.Sprintf("Skill%d", i+1)}
	}

	dataMap := map[string]interface{}{
		"skills": skills,
	}

	// Act
	distributions, err := paginator.PaginateMultiList(dataMap, "page_1.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(distributions) <= 1 {
		t.Fatalf("Expected multiple pages for 100 skills, got %d", len(distributions))
	}

	// Verify continuation template naming
	if distributions[0].TemplateName != "page_1.html" {
		t.Errorf("First page should use base template, got '%s'", distributions[0].TemplateName)
	}

	for i := 1; i < len(distributions); i++ {
		expectedName := "page_1.2.html"
		if distributions[i].TemplateName != expectedName {
			t.Errorf("Page %d should use continuation template '%s', got '%s'",
				i+1, expectedName, distributions[i].TemplateName)
		}
	}

	// Verify all skills are distributed
	totalSkills := 0
	for _, dist := range distributions {
		for blockName, data := range dist.Data {
			// Only count skill blocks to avoid counting the same skills multiple times
			if skillsList, ok := data.([]SkillViewModel); ok &&
				(blockName == "skills_column1" || blockName == "skills_column2" ||
					blockName == "skills_column3" || blockName == "skills_column4") {
				totalSkills += len(skillsList)
			}
		}
	}

	if totalSkills != 100 {
		t.Errorf("Expected 100 total skills distributed, got %d", totalSkills)
	}
}

// TestPaginateMultiList_EmptyLists tests handling of empty data
func TestPaginateMultiList_EmptyLists(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	dataMap := map[string]interface{}{
		"skills": []SkillViewModel{},
	}

	// Act
	distributions, err := paginator.PaginateMultiList(dataMap, "page_1.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(distributions) != 0 {
		t.Errorf("Expected 0 pages for empty data, got %d", len(distributions))
	}
}

// TestPaginateMultiList_WithFilters tests filtering by learned/unlearned/language
func TestPaginateMultiList_WithFilters(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create mixed skills
	skills := []SkillViewModel{
		{Name: "Learned1", IsLearned: true, Category: "Kampf"},
		{Name: "Learned2", IsLearned: true, Category: "Körper"},
		{Name: "Unlearned1", IsLearned: false, Category: "Kampf"},
		{Name: "Unlearned2", IsLearned: false, Category: "Social"},
		{Name: "Language1", Category: "Sprache"},
		{Name: "Language2", Category: "Sprache"},
	}

	dataMap := map[string]interface{}{
		"skills": skills,
	}

	// Act
	distributions, err := paginator.PaginateMultiList(dataMap, "page_2.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(distributions) == 0 {
		t.Fatal("Expected at least 1 page")
	}

	// Verify filtering worked - check that learned, unlearned, and language skills are in separate blocks
	page := distributions[0]

	if learnedData, ok := page.Data["skills_learned"]; ok {
		learned := learnedData.([]SkillViewModel)
		for _, skill := range learned {
			if !skill.IsLearned {
				t.Errorf("Found unlearned skill '%s' in learned block", skill.Name)
			}
		}
	}

	if languageData, ok := page.Data["skills_languages"]; ok {
		languages := languageData.([]SkillViewModel)
		for _, skill := range languages {
			if skill.Category != "Sprache" {
				t.Errorf("Found non-language skill '%s' in language block", skill.Name)
			}
		}
	}
}

// TestPaginateMultiList_UnknownTemplate tests error handling for unknown templates
func TestPaginateMultiList_UnknownTemplate(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	dataMap := map[string]interface{}{
		"skills": []SkillViewModel{{Name: "Test"}},
	}

	// Act
	distributions, err := paginator.PaginateMultiList(dataMap, "nonexistent_template.html")

	// Assert
	if err == nil {
		t.Fatal("Expected error for unknown template, got nil")
	}

	if distributions != nil {
		t.Error("Expected nil distributions for error case")
	}
}

// TestPaginateMultiList_PageNumbering tests that page numbers are sequential
func TestPaginateMultiList_PageNumbering(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Create enough skills for 3 pages
	skills := make([]SkillViewModel, 120)
	for i := 0; i < 120; i++ {
		skills[i] = SkillViewModel{Name: fmt.Sprintf("Skill%d", i+1)}
	}

	dataMap := map[string]interface{}{
		"skills": skills,
	}

	// Act
	distributions, err := paginator.PaginateMultiList(dataMap, "page_1.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify page numbers are sequential
	for i, dist := range distributions {
		expectedPageNum := i + 1
		if dist.PageNumber != expectedPageNum {
			t.Errorf("Page %d has incorrect PageNumber: expected %d, got %d",
				i+1, expectedPageNum, dist.PageNumber)
		}
	}
}

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

	// Get column capacities from template
	var col1Capacity int
	for _, tmpl := range templateSet.Templates {
		if tmpl.Metadata.Name == "page_1.html" {
			for _, block := range tmpl.Metadata.Blocks {
				if block.Name == "skills_column1" {
					col1Capacity = block.MaxItems
					break
				}
			}
			break
		}
	}

	// Create skills that fit within first column only
	numSkills := col1Capacity
	skills := make([]SkillViewModel, numSkills)
	for i := 0; i < numSkills; i++ {
		skills[i] = SkillViewModel{Name: "Skill" + string(rune('A'+i))}
	}

	// Act
	pages, err := paginator.PaginateSkills(skills, "page_1.html", "")

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

	// Column 1 should have all skills (exactly matching capacity)
	col1Data, ok := page.Data["skills_column1"].([]SkillViewModel)
	if !ok {
		t.Fatal("skills_column1 data not found or wrong type")
	}
	if len(col1Data) != numSkills {
		t.Errorf("Expected %d skills in column 1 (capacity), got %d", numSkills, len(col1Data))
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

	// Get expected capacity from template
	var page1Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_1.html" {
			page1Template = &templateSet.Templates[i]
			break
		}
	}
	col1Block := GetBlockByName(page1Template.Metadata.Blocks, "skills_column1")
	col2Block := GetBlockByName(page1Template.Metadata.Blocks, "skills_column2")
	totalCapacity := col1Block.MaxItems + col2Block.MaxItems

	// Create enough skills to use both columns but fit on one page
	// Use totalCapacity - 1 to test partial fill of second column
	numSkills := col1Block.MaxItems + 2
	if numSkills > totalCapacity {
		numSkills = totalCapacity
	}

	skills := make([]SkillViewModel, numSkills)
	for i := 0; i < numSkills; i++ {
		skills[i] = SkillViewModel{Name: "Skill" + string(rune(i))}
	}

	// Act
	pages, err := paginator.PaginateSkills(skills, "page_1.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should fit on one page
	expectedPages := (numSkills + totalCapacity - 1) / totalCapacity
	if len(pages) != expectedPages {
		t.Fatalf("Expected %d page, got %d", expectedPages, len(pages))
	}

	page := pages[0]

	// Column 1 should have max capacity skills from template
	col1Data := page.Data["skills_column1"].([]SkillViewModel)
	if len(col1Data) != col1Block.MaxItems {
		t.Errorf("Expected %d skills in column 1 (from template), got %d", col1Block.MaxItems, len(col1Data))
	}

	// Column 2 should have remaining skills
	col2Data := page.Data["skills_column2"].([]SkillViewModel)
	expectedCol2 := numSkills - col1Block.MaxItems
	if len(col2Data) != expectedCol2 {
		t.Errorf("Expected %d skills in column 2 (%d total - %d in col1), got %d", expectedCol2, numSkills, col1Block.MaxItems, len(col2Data))
	}
}

func TestPaginateSkills_MultiPage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Get column capacities from template
	var page1Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_1.html" {
			page1Template = &templateSet.Templates[i]
			break
		}
	}
	col1Block := GetBlockByName(page1Template.Metadata.Blocks, "skills_column1")
	col2Block := GetBlockByName(page1Template.Metadata.Blocks, "skills_column2")
	totalCapacity := col1Block.MaxItems + col2Block.MaxItems

	// Create enough skills to span exactly 2 pages
	numSkills := totalCapacity + 1
	skills := make([]SkillViewModel, numSkills)
	for i := 0; i < numSkills; i++ {
		skills[i] = SkillViewModel{Name: "Skill" + string(rune(i))}
	}

	// Act
	pages, err := paginator.PaginateSkills(skills, "page_1.html", "")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedPages := 2
	if len(pages) != expectedPages {
		t.Fatalf("Expected %d pages, got %d", expectedPages, len(pages))
	}

	// Page 1 should have full capacity (col1 + col2)
	page1 := pages[0]
	col1Page1 := page1.Data["skills_column1"].([]SkillViewModel)
	col2Page1 := page1.Data["skills_column2"].([]SkillViewModel)
	if len(col1Page1) != col1Block.MaxItems {
		t.Errorf("Page 1 col1: expected %d skills (template capacity), got %d", col1Block.MaxItems, len(col1Page1))
	}
	if len(col2Page1) != col2Block.MaxItems {
		t.Errorf("Page 1 col2: expected %d skills (template capacity), got %d", col2Block.MaxItems, len(col2Page1))
	}

	// Page 2 should have remaining skills (just 1 skill)
	page2 := pages[1]
	col1Page2 := page2.Data["skills_column1"].([]SkillViewModel)
	col2Page2 := page2.Data["skills_column2"].([]SkillViewModel)
	remainingSkills := numSkills - totalCapacity
	if len(col1Page2) != remainingSkills {
		t.Errorf("Page 2 col1: expected %d skills (remaining), got %d", remainingSkills, len(col1Page2))
	}
	if len(col2Page2) != 0 {
		t.Errorf("Page 2 col2: expected 0 skills, got %d", len(col2Page2))
	}
}

func TestPaginateSpells_TwoColumns(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Get spell column capacities from template
	var page3Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_3.html" {
			page3Template = &templateSet.Templates[i]
			break
		}
	}
	var leftCapacity, rightCapacity int
	for i := range page3Template.Metadata.Blocks {
		if page3Template.Metadata.Blocks[i].Name == "spells_column1" {
			leftCapacity = page3Template.Metadata.Blocks[i].MaxItems
		} else if page3Template.Metadata.Blocks[i].Name == "spells_column2" {
			rightCapacity = page3Template.Metadata.Blocks[i].MaxItems
		}
	}
	totalCapacity := leftCapacity + rightCapacity

	// Create spells that fit within total capacity
	testCount := totalCapacity
	if testCount > 0 {
		spells := make([]SpellViewModel, testCount)
		for i := 0; i < testCount; i++ {
			spells[i] = SpellViewModel{Name: "Spell" + string(rune('A'+i))}
		}

		// Act
		pages, err := paginator.PaginateSpells(spells, "page_3.html")

		// Assert
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(pages) != 1 {
			t.Fatalf("Expected 1 page (capacity %d from template), got %d", totalCapacity, len(pages))
		}

		page := pages[0]

		// Left column should be filled first
		leftData := page.Data["spells_column1"].([]SpellViewModel)
		expectedLeft := leftCapacity
		if testCount < leftCapacity {
			expectedLeft = testCount
		}
		if len(leftData) != expectedLeft {
			t.Errorf("Expected %d spells in left column (template capacity), got %d", expectedLeft, len(leftData))
		}

		// Right column gets remainder
		rightData := page.Data["spells_column2"].([]SpellViewModel)
		expectedRight := testCount - expectedLeft
		if len(rightData) != expectedRight {
			t.Errorf("Expected %d spells in right column (remaining), got %d", expectedRight, len(rightData))
		}
	}
}

func TestPaginateSpells_MultiPage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Get spell column capacities from template
	var page3Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_3.html" {
			page3Template = &templateSet.Templates[i]
			break
		}
	}
	var leftCapacity, rightCapacity int
	for i := range page3Template.Metadata.Blocks {
		if page3Template.Metadata.Blocks[i].Name == "spells_column1" {
			leftCapacity = page3Template.Metadata.Blocks[i].MaxItems
		} else if page3Template.Metadata.Blocks[i].Name == "spells_column2" {
			rightCapacity = page3Template.Metadata.Blocks[i].MaxItems
		}
	}

	// Create 30 spells - should span multiple pages based on template capacity
	spells := make([]SpellViewModel, 30)
	for i := 0; i < 30; i++ {
		spells[i] = SpellViewModel{Name: "Spell" + string(rune(i))}
	}

	// Act
	pages, err := paginator.PaginateSpells(spells, "page_3.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// With capacity from template, 30 spells require 2 pages
	totalCapacity := leftCapacity + rightCapacity
	expectedPages := (30 + totalCapacity - 1) / totalCapacity // ceiling division
	if len(pages) != expectedPages {
		t.Fatalf("Expected %d pages (capacity %d from template), got %d", expectedPages, totalCapacity, len(pages))
	}

	// Page 1 should have full capacity (left + right)
	page1 := pages[0]
	leftPage1 := page1.Data["spells_column1"].([]SpellViewModel)
	rightPage1 := page1.Data["spells_column2"].([]SpellViewModel)
	if len(leftPage1) != leftCapacity {
		t.Errorf("Page 1 left: expected %d spells (template capacity), got %d", leftCapacity, len(leftPage1))
	}
	if len(rightPage1) != rightCapacity {
		t.Errorf("Page 1 right: expected %d spells (template capacity), got %d", rightCapacity, len(rightPage1))
	}
}

func TestPaginateWeapons_SinglePage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Get weapon capacity from template
	var weaponCapacity int
	for _, tmpl := range templateSet.Templates {
		if tmpl.Metadata.Name == "page_2.html" {
			for _, block := range tmpl.Metadata.Blocks {
				if block.ListType == "weapons" {
					weaponCapacity = block.MaxItems
					break
				}
			}
			break
		}
	}

	// Create fewer weapons than capacity
	numWeapons := weaponCapacity - 2
	if numWeapons < 1 {
		numWeapons = 1
	}
	weapons := make([]WeaponViewModel, numWeapons)
	for i := 0; i < numWeapons; i++ {
		weapons[i] = WeaponViewModel{Name: "Weapon" + string(rune('A'+i))}
	}

	// Act
	pages, err := paginator.PaginateWeapons(weapons, "page_2.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedPages := 1
	if len(pages) != expectedPages {
		t.Fatalf("Expected %d page (capacity %d, items %d), got %d", expectedPages, weaponCapacity, numWeapons, len(pages))
	}

	page := pages[0]
	weaponsData := page.Data["weapons_main"].([]WeaponViewModel)
	if len(weaponsData) != numWeapons {
		t.Errorf("Expected %d weapons, got %d", numWeapons, len(weaponsData))
	}
}

func TestPaginateWeapons_MultiPage(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Get weapon capacity from template
	var page2Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_2.html" {
			page2Template = &templateSet.Templates[i]
			break
		}
	}
	if page2Template == nil {
		t.Fatal("page_2.html template not found")
	}
	var weaponsBlock *BlockMetadata
	for i := range page2Template.Metadata.Blocks {
		if page2Template.Metadata.Blocks[i].Name == "weapons_main" {
			weaponsBlock = &page2Template.Metadata.Blocks[i]
			break
		}
	}
	if weaponsBlock == nil {
		t.Fatal("weapons_main block not found")
	}
	weaponCapacity := weaponsBlock.MaxItems

	// Create 50 weapons - should span 3 pages based on template capacity
	weapons := make([]WeaponViewModel, 50)
	for i := 0; i < 50; i++ {
		weapons[i] = WeaponViewModel{Name: "Weapon" + string(rune(i))}
	}

	// Act
	pages, err := paginator.PaginateWeapons(weapons, "page_2.html")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedPages := (50 + weaponCapacity - 1) / weaponCapacity // ceiling division
	if len(pages) != expectedPages {
		t.Fatalf("Expected %d pages (%d capacity from template), got %d", expectedPages, weaponCapacity, len(pages))
	}

	// First pages should have weaponCapacity weapons each
	for i := 0; i < expectedPages-1; i++ {
		pageWeapons := pages[i].Data["weapons_main"].([]WeaponViewModel)
		if len(pageWeapons) != weaponCapacity {
			t.Errorf("Page %d: expected %d weapons (template capacity), got %d", i+1, weaponCapacity, len(pageWeapons))
		}
	}

	// Last page should have remaining weapons
	remainingWeapons := 50 - ((expectedPages - 1) * weaponCapacity)
	lastPageWeapons := pages[expectedPages-1].Data["weapons_main"].([]WeaponViewModel)
	if len(lastPageWeapons) != remainingWeapons {
		t.Errorf("Last page: expected %d weapons (remaining), got %d", remainingWeapons, len(lastPageWeapons))
	}
}

func TestCalculatePagesNeeded(t *testing.T) {
	// Arrange
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Get capacities from templates
	var page1Template, page2Template, page3Template *TemplateWithMeta
	for i := range templateSet.Templates {
		switch templateSet.Templates[i].Metadata.Name {
		case "page_1.html":
			page1Template = &templateSet.Templates[i]
		case "page_2.html":
			page2Template = &templateSet.Templates[i]
		case "page_3.html":
			page3Template = &templateSet.Templates[i]
		}
	}

	// Get skill capacity (col1 + col2)
	var skillCol1, skillCol2 int
	for i := range page1Template.Metadata.Blocks {
		if page1Template.Metadata.Blocks[i].Name == "skills_column1" {
			skillCol1 = page1Template.Metadata.Blocks[i].MaxItems
		} else if page1Template.Metadata.Blocks[i].Name == "skills_column2" {
			skillCol2 = page1Template.Metadata.Blocks[i].MaxItems
		}
	}
	skillCapacity := skillCol1 + skillCol2

	// Get weapon capacity
	var weaponCapacity int
	for i := range page2Template.Metadata.Blocks {
		if page2Template.Metadata.Blocks[i].Name == "weapons_main" {
			weaponCapacity = page2Template.Metadata.Blocks[i].MaxItems
			break
		}
	}

	// Get spell capacity (col1 + col2)
	var spellCol1, spellCol2 int
	for i := range page3Template.Metadata.Blocks {
		if page3Template.Metadata.Blocks[i].Name == "spells_column1" {
			spellCol1 = page3Template.Metadata.Blocks[i].MaxItems
		} else if page3Template.Metadata.Blocks[i].Name == "spells_column2" {
			spellCol2 = page3Template.Metadata.Blocks[i].MaxItems
		}
	}
	spellCapacity := spellCol1 + spellCol2

	testCases := []struct {
		name          string
		templateName  string
		listType      string
		itemCount     int
		expectedPages int
	}{
		{"10 skills on page1", "page_1.html", "skills", 10, 1},
		{fmt.Sprintf("%d skills on page1", skillCapacity), "page_1.html", "skills", skillCapacity, 1},
		{fmt.Sprintf("%d skills on page1", skillCapacity+1), "page_1.html", "skills", skillCapacity + 1, 2},
		{"100 skills on page1", "page_1.html", "skills", 100, (100 + skillCapacity - 1) / skillCapacity},  // Dynamic calculation
		{"10 weapons on page2", "page_2.html", "weapons", 10, (10 + weaponCapacity - 1) / weaponCapacity}, // Dynamic calculation
		{fmt.Sprintf("%d weapons on page2", weaponCapacity), "page_2.html", "weapons", weaponCapacity, 1},
		{fmt.Sprintf("%d weapons on page2", weaponCapacity+1), "page_2.html", "weapons", weaponCapacity + 1, 2},
		{"10 spells on page3", "page_3.html", "spells", 10, (10 + spellCapacity - 1) / spellCapacity}, // Dynamic calculation
		{fmt.Sprintf("%d spells on page3", spellCapacity), "page_3.html", "spells", spellCapacity, 1},
		{fmt.Sprintf("%d spells on page3", spellCapacity+1), "page_3.html", "spells", spellCapacity + 1, 2},
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
	pages, err := paginator.PaginateSkills(skills, "page_1.html", "")

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
