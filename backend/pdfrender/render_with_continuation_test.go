package pdfrender

import (
	"bamort/models"
	"testing"
)

// TestRenderWithContinuations_SkillsOverflow tests that continuation pages are rendered
// when skills exceed template capacity
func TestRenderWithContinuations_SkillsOverflow(t *testing.T) {
	// Arrange - Create character with many skills to force continuation
	char := &models.Char{
		BamortBase: models.BamortBase{
			Name: "Test Character With Many Skills",
		},
		Typ:   "Krieger",
		Grad:  5,
		Alter: 30,
		Eigenschaften: []models.Eigenschaft{
			{Name: "St", Value: 90},
			{Name: "Gs", Value: 80},
			{Name: "Gw", Value: 70},
			{Name: "Ko", Value: 85},
			{Name: "In", Value: 60},
			{Name: "Zt", Value: 55},
			{Name: "Au", Value: 75},
			{Name: "pA", Value: 65},
			{Name: "Wk", Value: 50},
		},
		Lp: models.Lp{Value: 20, Max: 25},
		Ap: models.Ap{Value: 15, Max: 20},
		B:  models.B{Value: 15},
	}

	// Add 50 skills to force multiple pages
	char.Fertigkeiten = make([]models.SkFertigkeit, 50)
	for i := 0; i < 50; i++ {
		char.Fertigkeiten[i] = models.SkFertigkeit{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: "Skill " + string(rune('A'+i%26))},
			},
			Fertigkeitswert: 10 + i%15,
			Pp:              i % 8,
		}
	}

	// Map to view model
	viewModel, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("Failed to map character: %v", err)
	}

	// Load templates
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	if err = loader.LoadTemplates(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	renderer := NewPDFRenderer()

	// Act - Render page1 with continuations
	pdfs, err := RenderPageWithContinuations(
		viewModel,
		"page_1.html",
		1,
		"20.12.2025",
		loader,
		renderer,
	)

	// Assert
	if err != nil {
		t.Fatalf("Failed to render with continuations: %v", err)
	}

	// Get template capacity to calculate expected pages
	templateSet := DefaultA4QuerTemplateSet()
	var skillsCapacity int
	for _, tmpl := range templateSet.Templates {
		if tmpl.Metadata.Name == "page_1.html" {
			for _, block := range tmpl.Metadata.Blocks {
				if block.ListType == "skills" {
					skillsCapacity += block.MaxItems
				}
			}
			break
		}
	}

	expectedPages := (50 + skillsCapacity - 1) / skillsCapacity
	if len(pdfs) != expectedPages {
		t.Errorf("Expected %d PDFs (pages) for 50 skills with capacity %d, got %d",
			expectedPages, skillsCapacity, len(pdfs))
	}

	// Verify each PDF is valid
	for i, pdf := range pdfs {
		if len(pdf) == 0 {
			t.Errorf("PDF %d is empty", i+1)
		}
		if string(pdf[0:4]) != "%PDF" {
			t.Errorf("PDF %d does not start with PDF marker", i+1)
		}
	}

	t.Logf("✓ Successfully generated %d continuation pages for 50 skills (capacity: %d)",
		len(pdfs), skillsCapacity)
}

// TestRenderWithContinuations_NoOverflow tests that single page is rendered
// when items fit within capacity
func TestRenderWithContinuations_NoOverflow(t *testing.T) {
	// Arrange - Create character with few skills
	char := &models.Char{
		BamortBase: models.BamortBase{
			Name: "Test Character",
		},
		Typ:   "Krieger",
		Grad:  5,
		Alter: 30,
		Eigenschaften: []models.Eigenschaft{
			{Name: "St", Value: 90},
			{Name: "Gs", Value: 80},
			{Name: "Gw", Value: 70},
			{Name: "Ko", Value: 85},
			{Name: "In", Value: 60},
			{Name: "Zt", Value: 55},
			{Name: "Au", Value: 75},
			{Name: "pA", Value: 65},
			{Name: "Wk", Value: 50},
		},
		Lp: models.Lp{Value: 20, Max: 25},
		Ap: models.Ap{Value: 15, Max: 20},
		B:  models.B{Value: 15},
	}

	// Add only 3 skills - should fit on one page
	char.Fertigkeiten = make([]models.SkFertigkeit, 3)
	for i := 0; i < 3; i++ {
		char.Fertigkeiten[i] = models.SkFertigkeit{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: "Skill " + string(rune('A'+i))},
			},
			Fertigkeitswert: 10,
			Pp:              5,
		}
	}

	// Map to view model
	viewModel, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("Failed to map character: %v", err)
	}

	// Load templates
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	if err = loader.LoadTemplates(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	renderer := NewPDFRenderer()

	// Act
	pdfs, err := RenderPageWithContinuations(
		viewModel,
		"page_1.html",
		1,
		"20.12.2025",
		loader,
		renderer,
	)

	// Assert
	if err != nil {
		t.Fatalf("Failed to render: %v", err)
	}

	if len(pdfs) != 1 {
		t.Errorf("Expected 1 PDF for 3 skills, got %d", len(pdfs))
	}

	t.Logf("✓ Correctly generated single page for 3 skills")
}
