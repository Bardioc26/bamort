package pdfrender

import (
	"strings"
	"testing"
)

func TestFillToCapacity_Skills(t *testing.T) {
	// Test filling skills list to capacity
	skills := []SkillViewModel{
		{Name: "Skill 1", Value: 10},
		{Name: "Skill 2", Value: 12},
	}

	filled := FillToCapacity(skills, 5)

	if len(filled) != 5 {
		t.Errorf("Expected 5 items after filling, got %d", len(filled))
	}

	// First 2 should be original skills
	if filled[0].Name != "Skill 1" {
		t.Error("First item should be original")
	}

	// Last 3 should be empty
	if filled[2].Name != "" {
		t.Error("Filled items should have empty Name")
	}
	if filled[4].Value != 0 {
		t.Error("Filled items should have zero Value")
	}
}

func TestFillToCapacity_LessThanCapacity(t *testing.T) {
	// If already at or over capacity, should not add more
	skills := []SkillViewModel{
		{Name: "Skill 1"},
		{Name: "Skill 2"},
		{Name: "Skill 3"},
	}

	filled := FillToCapacity(skills, 2)

	// Should keep original 3, not truncate
	if len(filled) != 3 {
		t.Errorf("Expected 3 items (original), got %d", len(filled))
	}
}

func TestTemplateWithEmptyRows(t *testing.T) {
	// Integration test: render template with filled rows
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	err := loader.LoadTemplates()
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	// Create data with few skills
	skills := []SkillViewModel{
		{Name: "Schwimmen", Value: 10},
		{Name: "Klettern", Value: 8},
	}

	// Fill to column capacity (29)
	filledCol1 := FillToCapacity(skills, 29)

	pageData := &PageData{
		Character: CharacterInfo{
			Name: "Test Character",
		},
		SkillsColumn1: filledCol1,
		SkillsColumn2: FillToCapacity([]SkillViewModel{}, 29),
		Meta: PageMeta{
			Date: "19.12.2025",
		},
	}

	html, err := loader.RenderTemplate("page_1.html", pageData)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Get expected skill capacity from template
	templateSet := DefaultA4QuerTemplateSet()
	var page1Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_1.html" {
			page1Template = &templateSet.Templates[i]
			break
		}
	}
	if page1Template == nil {
		t.Fatal("page_1.html template not found")
	}
	var col1Capacity int
	for i := range page1Template.Metadata.Blocks {
		if page1Template.Metadata.Blocks[i].Name == "skills_column1" {
			col1Capacity = page1Template.Metadata.Blocks[i].MaxItems
			break
		}
	}

	// Count the number of <tr> tags in skills table
	// Should have col1Capacity rows (2 filled + remaining empty)
	trCount := strings.Count(html, "<tr><td>")
	if trCount < col1Capacity {
		t.Errorf("Expected at least %d skill rows in HTML (from template), got %d", col1Capacity, trCount)
	}
}
