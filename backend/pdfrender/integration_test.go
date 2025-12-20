package pdfrender

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// TestIntegration_FullPDFGeneration tests the complete workflow:
// Character -> ViewModel -> Template -> HTML -> PDF
func TestIntegration_FullPDFGeneration(t *testing.T) {
	// Arrange - Create a test character
	char := &models.Char{
		BamortBase: models.BamortBase{
			Name: "Bjarnfinnur Haberdson",
		},
		Typ:         "Krieger",
		Grad:        5,
		Alter:       35,
		Groesse:     180,
		Gewicht:     85,
		Gender:      "m",
		SocialClass: "Frei",
		Glaube:      "Apshai",
		Herkunft:    "Erainn",
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
		Lp: models.Lp{
			Value: 20,
			Max:   25,
		},
		Ap: models.Ap{
			Value: 15,
			Max:   20,
		},
		B: models.B{
			Value: 15,
		},
		Vermoegen: models.Vermoegen{
			Goldstuecke:   100,
			Silberstuecke: 50,
			Kupferstuecke: 25,
		},
		Fertigkeiten: []models.SkFertigkeit{
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{Name: "Schwimmen"},
				},
				Fertigkeitswert: 12,
				Pp:              5,
			},
			{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{Name: "Klettern"},
				},
				Fertigkeitswert: 10,
				Pp:              3,
			},
		},
		Waffenfertigkeiten: []models.SkWaffenfertigkeit{
			{
				SkFertigkeit: models.SkFertigkeit{
					BamortCharTrait: models.BamortCharTrait{
						BamortBase: models.BamortBase{Name: "Langschwert"},
					},
					Fertigkeitswert: 14,
					Pp:              8,
				},
			},
		},
	}

	// Step 1: Map to ViewModel
	viewModel, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("Failed to map character to view model: %v", err)
	}
	if viewModel.Character.Name != "Bjarnfinnur Haberdson" {
		t.Fatalf("ViewModel mapping failed: expected name 'Bjarnfinnur Haberdson', got '%s'", viewModel.Character.Name)
	}

	// Step 2: Load templates
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	if err = loader.LoadTemplates(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	// Step 3: Prepare paginated data and render template to HTML
	pageData, err := PreparePaginatedPageData(viewModel, "page1_stats.html", 1, "18.12.2025")
	if err != nil {
		t.Fatalf("Failed to prepare paginated data: %v", err)
	}

	html, err := loader.RenderTemplate("page1_stats.html", pageData)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Verify HTML contains expected data
	if !strings.Contains(html, "Bjarnfinnur Haberdson") {
		t.Error("HTML does not contain character name")
	}
	if !strings.Contains(html, "Schwimmen") {
		t.Error("HTML does not contain skill 'Schwimmen'")
	}

	// Step 4: Convert HTML to PDF
	renderer := NewPDFRenderer()
	pdfBytes, err := renderer.RenderHTMLToPDF(html)
	if err != nil {
		t.Fatalf("Failed to render PDF: %v", err)
	}

	// Verify PDF was created
	if len(pdfBytes) == 0 {
		t.Fatal("PDF bytes are empty")
	}

	if string(pdfBytes[0:4]) != "%PDF" {
		t.Error("Output does not appear to be a PDF")
	}

	// PDF should be at least 10KB for a page with content
	if len(pdfBytes) < 10000 {
		t.Errorf("PDF seems too small (%d bytes), might be missing content", len(pdfBytes))
	}

	t.Logf("Successfully generated PDF of %d bytes", len(pdfBytes))
}

// TestIntegration_TemplateMetadata verifies that metadata parsing works with actual templates
func TestIntegration_TemplateMetadata(t *testing.T) {
	// Arrange
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	err := loader.LoadTemplates()
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	// Load template set to get expected values from actual template files
	templateSet := DefaultA4QuerTemplateSet()

	// Act & Assert - Check each template has metadata
	testCases := []struct {
		template      string
		expectedBlock string
	}{
		{"page1_stats.html", "skills_column1"},
		{"page2_play.html", "skills_learned"},
		{"page3_spell.html", "spells_left"},
		{"page3_spell.html", "spells_right"},
		{"page4_equip.html", "equipment_worn"},
	}

	for _, tc := range testCases {
		metadata := loader.GetTemplateMetadata(tc.template)
		if len(metadata) == 0 {
			t.Errorf("Template %s has no metadata", tc.template)
			continue
		}

		block := GetBlockByName(metadata, tc.expectedBlock)
		if block == nil {
			t.Errorf("Template %s missing block '%s'", tc.template, tc.expectedBlock)
			continue
		}

		// Get expected value from template set
		var expectedMax int
		for i := range templateSet.Templates {
			if templateSet.Templates[i].Metadata.Name == tc.template {
				for j := range templateSet.Templates[i].Metadata.Blocks {
					if templateSet.Templates[i].Metadata.Blocks[j].Name == tc.expectedBlock {
						expectedMax = templateSet.Templates[i].Metadata.Blocks[j].MaxItems
						break
					}
				}
				break
			}
		}

		if block.MaxItems != expectedMax {
			t.Errorf("Template %s block %s: expected max %d (from template), got %d",
				tc.template, tc.expectedBlock, expectedMax, block.MaxItems)
		}
	}
}

// TestIntegration_PaginationWithPDF tests pagination integrated with PDF generation
func TestIntegration_PaginationWithPDF(t *testing.T) {
	// Arrange - Create 100 skills to force pagination
	skills := make([]SkillViewModel, 100)
	for i := 0; i < 100; i++ {
		skills[i] = SkillViewModel{
			Name:           "Skill" + string(rune(i)),
			Value:          10 + i%20,
			PracticePoints: i % 10,
		}
	}

	// Create paginator and distribute skills
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	pages, err := paginator.PaginateSkills(skills, "page1_stats.html", "")
	if err != nil {
		t.Fatalf("Failed to paginate skills: %v", err)
	}

	// Calculate expected pages based on actual capacity
	var skillsCapacity int
	for _, tmpl := range templateSet.Templates {
		if tmpl.Metadata.Name == "page1_stats.html" {
			for _, block := range tmpl.Metadata.Blocks {
				if block.ListType == "skills" {
					skillsCapacity += block.MaxItems
				}
			}
			break
		}
	}
	expectedPages := (100 + skillsCapacity - 1) / skillsCapacity

	if len(pages) != expectedPages {
		t.Fatalf("Expected %d pages for 100 skills (capacity %d), got %d", expectedPages, skillsCapacity, len(pages))
	}

	t.Logf("Successfully paginated 100 skills into %d pages (capacity %d per page)", len(pages), skillsCapacity)

	// Load templates
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	if err = loader.LoadTemplates(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	// Render first page
	pageData := &PageData{
		Character: CharacterInfo{
			Name:   "Test Warrior",
			Player: "Test Player",
			Type:   "Krieger",
			Grade:  5,
		},
		Attributes: AttributeValues{
			St: 90, Gs: 80, Gw: 70, Ko: 85,
			In: 60, Zt: 55, Au: 75, PA: 65, Wk: 50, B: 15,
		},
		DerivedValues: DerivedValueSet{
			LPMax: 25, LPAktuell: 20,
			APMax: 20, APAktuell: 15,
		},
		Meta: PageMeta{
			Date:       "18.12.2025",
			PageNumber: 1,
		},
	}

	// Add paginated skills for page 1 - now with proper column split
	pageData.SkillsColumn1 = pages[0].Data["skills_column1"].([]SkillViewModel)
	pageData.SkillsColumn2 = pages[0].Data["skills_column2"].([]SkillViewModel)
	pageData.Skills = append(pageData.SkillsColumn1, pageData.SkillsColumn2...) // Keep for logging

	html, err := loader.RenderTemplate("page1_stats.html", pageData)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Verify HTML contains skills
	if !strings.Contains(html, "Skill") {
		t.Error("HTML does not contain skills")
	}

	// Generate PDF
	renderer := NewPDFRenderer()
	pdfBytes, err := renderer.RenderHTMLToPDF(html)
	if err != nil {
		t.Fatalf("Failed to render PDF: %v", err)
	}

	if len(pdfBytes) == 0 {
		t.Fatal("PDF bytes are empty")
	}

	t.Logf("Successfully generated page 1 PDF with %d skills, size: %d bytes", len(pageData.Skills), len(pdfBytes))

	// Verify second page has remaining skills (if there are multiple pages)
	if len(pages) > 1 {
		col1Page2 := pages[1].Data["skills_column1"].([]SkillViewModel)
		col2Page2 := pages[1].Data["skills_column2"].([]SkillViewModel)
		totalPage2 := len(col1Page2) + len(col2Page2)

		expectedPage2 := skillsCapacity
		if 100-skillsCapacity < skillsCapacity {
			expectedPage2 = 100 - skillsCapacity
		}

		if totalPage2 != expectedPage2 {
			t.Errorf("Expected %d skills on page 2, got %d", expectedPage2, totalPage2)
		}

		t.Logf("Page 2 has %d skills distributed across columns", totalPage2)
	}

	t.Logf("Page 2 would have %d skills distributed across columns", 100-skillsCapacity)
}

// TestIntegration_MultiPageSpellList tests spell pagination across multiple pages
func TestIntegration_MultiPageSpellList(t *testing.T) {
	// Get expected capacity from template
	templateSet := DefaultA4QuerTemplateSet()
	var page3Template *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page3_spell.html" {
			page3Template = &templateSet.Templates[i]
			break
		}
	}

	spellsLeftBlock := GetBlockByName(page3Template.Metadata.Blocks, "spells_left")
	spellsRightBlock := GetBlockByName(page3Template.Metadata.Blocks, "spells_right")
	expectedSpellCapacity := spellsLeftBlock.MaxItems + spellsRightBlock.MaxItems

	// Arrange - Create 30 spells
	spells := make([]SpellViewModel, 30)
	for i := 0; i < 30; i++ {
		spells[i] = SpellViewModel{
			Name:          "Zauber Nr. " + string(rune('A'+i%26)),
			AP:            "5",
			Stufe:         1,
			Wirkungsdauer: "1 Minute",
			Zauberdauer:   "1 sec",
		}
	}

	// Create paginator
	paginator := NewPaginator(templateSet)

	// Paginate spells
	pages, err := paginator.PaginateSpells(spells, "page3_spell.html")
	if err != nil {
		t.Fatalf("Failed to paginate spells: %v", err)
	}

	// Calculate expected pages
	expectedPages := (30 + expectedSpellCapacity - 1) / expectedSpellCapacity // Ceiling division

	// Verify distribution
	if len(pages) != expectedPages {
		t.Fatalf("Expected %d pages for 30 spells (capacity %d), got %d", expectedPages, expectedSpellCapacity, len(pages))
	}

	// Page 1 should have min(30, capacity) spells
	leftPage1 := pages[0].Data["spells_left"].([]SpellViewModel)
	rightPage1 := pages[0].Data["spells_right"].([]SpellViewModel)
	totalPage1 := len(leftPage1) + len(rightPage1)

	expectedPage1 := expectedSpellCapacity
	if 30 < expectedSpellCapacity {
		expectedPage1 = 30
	}

	if totalPage1 != expectedPage1 {
		t.Errorf("Expected %d spells on page 1 (capacity %d+%d), got %d", expectedPage1, spellsLeftBlock.MaxItems, spellsRightBlock.MaxItems, totalPage1)
	}

	t.Logf("Successfully distributed 30 spells: Page 1 has %d (left %d, right %d)", totalPage1, len(leftPage1), len(rightPage1))
}

// TestIntegration_CompleteWorkflow demonstrates the full workflow from character to multi-page PDF
func TestIntegration_CompleteWorkflow(t *testing.T) {
	// Step 1: Create a character with lots of data to force pagination
	char := &models.Char{
		BamortBase: models.BamortBase{
			Name: "Complete Test Character",
		},
		Typ:   "Krieger",
		Grad:  10,
		Alter: 45,
		Eigenschaften: []models.Eigenschaft{
			{Name: "St", Value: 95},
			{Name: "Gs", Value: 85},
			{Name: "Gw", Value: 80},
			{Name: "Ko", Value: 90},
			{Name: "In", Value: 70},
			{Name: "Zt", Value: 60},
			{Name: "Au", Value: 75},
			{Name: "pA", Value: 80},
			{Name: "Wk", Value: 65},
		},
		Lp: models.Lp{Value: 45, Max: 50},
		Ap: models.Ap{Value: 30, Max: 35},
		B:  models.B{Value: 20},
	}

	// Add many skills to force pagination
	char.Fertigkeiten = make([]models.SkFertigkeit, 70)
	for i := 0; i < 70; i++ {
		char.Fertigkeiten[i] = models.SkFertigkeit{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: "Fertigkeit " + string(rune('A'+i%26)) + string(rune('0'+i/26))},
			},
			Fertigkeitswert: 10 + i%15,
			Pp:              i % 8,
		}
	}

	// Add weapons
	char.Waffenfertigkeiten = make([]models.SkWaffenfertigkeit, 5)
	for i := 0; i < 5; i++ {
		char.Waffenfertigkeiten[i] = models.SkWaffenfertigkeit{
			SkFertigkeit: models.SkFertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{Name: "Waffe " + string(rune('A'+i))},
				},
				Fertigkeitswert: 12 + i*2,
				Pp:              i,
			},
		}
	}

	// Step 2: Map to view model
	viewModel, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("Failed to map character: %v", err)
	}

	t.Logf("Mapped character with %d skills", len(viewModel.Skills))

	// Step 3: Initialize pagination system
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Step 4: Paginate skills
	skillPages, err := paginator.PaginateSkills(viewModel.Skills, "page1_stats.html", "")
	if err != nil {
		t.Fatalf("Failed to paginate skills: %v", err)
	}

	t.Logf("Skills distributed across %d pages", len(skillPages))

	// Step 5: Load templates
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	if err = loader.LoadTemplates(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	// Step 6: Render and generate PDFs for each page
	renderer := NewPDFRenderer()
	totalPDFSize := 0

	for _, page := range skillPages {
		// Extract paginated data
		col1 := page.Data["skills_column1"].([]SkillViewModel)
		col2 := page.Data["skills_column2"].([]SkillViewModel)

		pageData := &PageData{
			Character:     viewModel.Character,
			Attributes:    viewModel.Attributes,
			DerivedValues: viewModel.DerivedValues,
			Skills:        append(col1, col2...),
			Weapons:       viewModel.Weapons,
			Meta: PageMeta{
				Date:       "18.12.2025",
				PageNumber: page.PageNumber,
			},
		}

		// Render HTML
		html, err := loader.RenderTemplate(page.TemplateName, pageData)
		if err != nil {
			t.Fatalf("Failed to render page %d: %v", page.PageNumber, err)
		}

		// Generate PDF
		pdfBytes, err := renderer.RenderHTMLToPDF(html)
		if err != nil {
			t.Fatalf("Failed to generate PDF for page %d: %v", page.PageNumber, err)
		}

		// Verify PDF
		if len(pdfBytes) == 0 {
			t.Errorf("Page %d: PDF is empty", page.PageNumber)
		}
		if string(pdfBytes[0:4]) != "%PDF" {
			t.Errorf("Page %d: Invalid PDF format", page.PageNumber)
		}

		totalPDFSize += len(pdfBytes)
		t.Logf("Page %d: Generated %d bytes PDF with %d skills",
			page.PageNumber, len(pdfBytes), len(pageData.Skills))
	}

	// Summary
	t.Logf("✓ Complete workflow successful:")
	t.Logf("  - Character mapped: %s (Grade %d)", viewModel.Character.Name, viewModel.Character.Grade)
	t.Logf("  - Total skills: %d", len(viewModel.Skills))
	t.Logf("  - Pages generated: %d", len(skillPages))
	t.Logf("  - Total PDF size: %d bytes", totalPDFSize)
}

// TestVisualInspection_AllPages generates all 4 page types and saves them to disk
// Run with: go test -v ./pdfrender/... -run TestVisualInspection
func TestVisualInspection_AllPages(t *testing.T) {
	database.SetupTestDB()

	// Load character Fanjo Vetrani with ID 18 from test database
	char := &models.Char{}
	err := char.FirstID("18")
	if err != nil {
		t.Fatalf("Failed to load character with ID 18 (Fanjo Vetrani): %v", err)
	}

	// Verify we loaded the correct character
	if char.Name == "" {
		t.Fatalf("Character loaded but has empty name")
	}
	t.Logf("Loaded character: %s (ID: %d)", char.Name, char.ID)

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

	// Generate all pages with continuations if needed
	allPDFs := [][]byte{}
	var filePaths []string
	outputDir := "/tmp/bamort_pdf_test"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Page 1: Stats page with skills (may have continuations)
	t.Log("Generating Page 1: Stats...")
	page1PDFs, err := RenderPageWithContinuations(viewModel, "page1_stats.html", 1, "18.12.2025", loader, renderer)
	if err != nil {
		t.Fatalf("Failed to generate page1: %v", err)
	}
	t.Logf("  Generated %d PDF(s) for page 1", len(page1PDFs))

	for i, pdf := range page1PDFs {
		allPDFs = append(allPDFs, pdf)
		var filename string
		if i == 0 {
			filename = "page1_stats.pdf"
		} else {
			filename = fmt.Sprintf("page1_stats_continuation_%d.pdf", i)
		}
		path := filepath.Join(outputDir, filename)
		if err := os.WriteFile(path, pdf, 0644); err != nil {
			t.Errorf("Failed to write %s: %v", filename, err)
			continue
		}
		filePaths = append(filePaths, path)
		t.Logf("  ✓ Saved %s (%d bytes)", filename, len(pdf))
	}

	// Page 2: Play/Adventure page with weapons (may have continuations)
	t.Log("Generating Page 2: Play...")
	page2PDFs, err := RenderPageWithContinuations(viewModel, "page2_play.html", 2, "18.12.2025", loader, renderer)
	if err != nil {
		t.Fatalf("Failed to generate page2: %v", err)
	}
	t.Logf("  Generated %d PDF(s) for page 2", len(page2PDFs))

	for i, pdf := range page2PDFs {
		allPDFs = append(allPDFs, pdf)
		var filename string
		if i == 0 {
			filename = "page2_play.pdf"
		} else {
			filename = fmt.Sprintf("page2_play_continuation_%d.pdf", i)
		}
		path := filepath.Join(outputDir, filename)
		if err := os.WriteFile(path, pdf, 0644); err != nil {
			t.Errorf("Failed to write %s: %v", filename, err)
			continue
		}
		filePaths = append(filePaths, path)
		t.Logf("  ✓ Saved %s (%d bytes)", filename, len(pdf))
	}

	// Page 3: Spells page (may have continuations)
	t.Log("Generating Page 3: Spells...")
	page3PDFs, err := RenderPageWithContinuations(viewModel, "page3_spell.html", 3, "18.12.2025", loader, renderer)
	if err != nil {
		t.Fatalf("Failed to generate page3: %v", err)
	}
	t.Logf("  Generated %d PDF(s) for page 3", len(page3PDFs))

	for i, pdf := range page3PDFs {
		allPDFs = append(allPDFs, pdf)
		var filename string
		if i == 0 {
			filename = "page3_spell.pdf"
		} else {
			filename = fmt.Sprintf("page3_spell_continuation_%d.pdf", i)
		}
		path := filepath.Join(outputDir, filename)
		if err := os.WriteFile(path, pdf, 0644); err != nil {
			t.Errorf("Failed to write %s: %v", filename, err)
			continue
		}
		filePaths = append(filePaths, path)
		t.Logf("  ✓ Saved %s (%d bytes)", filename, len(pdf))
	}

	// Page 4: Equipment page (may have continuations)
	t.Log("Generating Page 4: Equipment...")
	page4PDFs, err := RenderPageWithContinuations(viewModel, "page4_equip.html", 4, "18.12.2025", loader, renderer)
	if err != nil {
		t.Fatalf("Failed to generate page4: %v", err)
	}
	t.Logf("  Generated %d PDF(s) for page 4", len(page4PDFs))

	for i, pdf := range page4PDFs {
		allPDFs = append(allPDFs, pdf)
		var filename string
		if i == 0 {
			filename = "page4_equip.pdf"
		} else {
			filename = fmt.Sprintf("page4_equip_continuation_%d.pdf", i)
		}
		path := filepath.Join(outputDir, filename)
		if err := os.WriteFile(path, pdf, 0644); err != nil {
			t.Errorf("Failed to write %s: %v", filename, err)
			continue
		}
		filePaths = append(filePaths, path)
		t.Logf("  ✓ Saved %s (%d bytes)", filename, len(pdf))
	}

	// Merge all PDFs into a single file
	combinedPath := outputDir + "/character_sheet_complete.pdf"
	if err := api.MergeCreateFile(filePaths, combinedPath, false, nil); err != nil {
		t.Fatalf("Failed to merge PDFs: %v", err)
	}

	// Get size of combined PDF
	combinedInfo, err := os.Stat(combinedPath)
	if err != nil {
		t.Fatalf("Failed to stat combined PDF: %v", err)
	}

	t.Logf("\n✓ Combined all pages into: %s (%d bytes)", combinedPath, combinedInfo.Size())

	// Summary
	t.Logf("\n✓ All pages generated successfully!")
	t.Logf("  Character: %s (Grade %d)", viewModel.Character.Name, viewModel.Character.Grade)
	t.Logf("  Skills: %d", len(viewModel.Skills))
	t.Logf("  Weapons: %d", len(viewModel.Weapons))
	t.Logf("  Spells: %d", len(viewModel.Spells))
	t.Logf("  Equipment: %d items", len(viewModel.Equipment))
	t.Logf("\n  Output directory: %s", outputDir)
	t.Logf("  Total PDFs generated: %d (including continuations)", len(allPDFs))
	t.Logf("  Combined PDF: %d bytes", combinedInfo.Size())
}
