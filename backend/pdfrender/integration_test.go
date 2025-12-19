package pdfrender

import (
	"bamort/models"
	"os"
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

	// Act & Assert - Check each template has metadata
	testCases := []struct {
		template      string
		expectedBlock string
		expectedMax   int
	}{
		{"page1_stats.html", "skills_column1", 29},
		{"page2_play.html", "skills_learned", 18}, // From template: MAX: 18
		{"page3_spell.html", "spells_left", 26},   // From template: MAX: 26
		{"page3_spell.html", "spells_right", 15},  // From template: MAX: 15
		{"page4_equip.html", "equipment_worn", 10},
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

		if block.MaxItems != tc.expectedMax {
			t.Errorf("Template %s block %s: expected max %d, got %d",
				tc.template, tc.expectedBlock, tc.expectedMax, block.MaxItems)
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

	// Should create 2 pages (64 + 36 skills)
	if len(pages) != 2 {
		t.Fatalf("Expected 2 pages, got %d", len(pages))
	}

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

	// Verify second page has remaining skills (94 total - 58 from page 1 = 36 remaining)
	// But with 29+29 capacity, it will be 29+13 = 42 on page 2
	col1Page2 := pages[1].Data["skills_column1"].([]SkillViewModel)
	col2Page2 := pages[1].Data["skills_column2"].([]SkillViewModel)
	totalPage2 := len(col1Page2) + len(col2Page2)

	if totalPage2 != 42 { // 100 total - 58 from page 1 = 42 remaining
		t.Errorf("Expected 42 skills on page 2, got %d", totalPage2)
	}

	t.Logf("Page 2 would have %d skills distributed across columns", totalPage2)
}

// TestIntegration_MultiPageSpellList tests spell pagination across multiple pages
func TestIntegration_MultiPageSpellList(t *testing.T) {
	// Arrange - Create 30 spells (will need 2 pages with 24 capacity each)
	spells := make([]SpellViewModel, 30)
	for i := 0; i < 30; i++ {
		spells[i] = SpellViewModel{
			Name:     "Zauber Nr. " + string(rune('A'+i%26)),
			AP:       5,
			Category: 1,
			Duration: "1 Minute",
			CastTime: "1 sec",
		}
	}

	// Create paginator
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Paginate spells
	pages, err := paginator.PaginateSpells(spells, "page3_spell.html")
	if err != nil {
		t.Fatalf("Failed to paginate spells: %v", err)
	}

	// With 30 capacity (20+10), should create 1 page for 30 spells

	// Verify distribution
	// With 20+10 capacity, 30 spells should fit on 1 page
	if len(pages) != 1 {
		t.Fatalf("Expected 1 page for 30 spells, got %d", len(pages))
	}

	// Page 1: 20 (left) + 10 (right) = 30 spells
	leftPage1 := pages[0].Data["spells_left"].([]SpellViewModel)
	rightPage1 := pages[0].Data["spells_right"].([]SpellViewModel)
	totalPage1 := len(leftPage1) + len(rightPage1)

	if totalPage1 != 30 {
		t.Errorf("Expected 30 spells on page 1 (20+10), got %d", totalPage1)
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
	// Create a rich character with data for all page types
	char := &models.Char{
		BamortBase: models.BamortBase{
			Name: "Integration Test",
		},
		Typ:         "Krieger",
		Grad:        8,
		Alter:       42,
		Groesse:     185,
		Gewicht:     92,
		Gender:      "m",
		SocialClass: "Frei",
		Glaube:      "Apshai",
		Herkunft:    "Erainn",
		Eigenschaften: []models.Eigenschaft{
			{Name: "St", Value: 95},
			{Name: "Gs", Value: 85},
			{Name: "Gw", Value: 80},
			{Name: "Ko", Value: 90},
			{Name: "In", Value: 75},
			{Name: "Zt", Value: 70},
			{Name: "Au", Value: 80},
			{Name: "pA", Value: 85},
			{Name: "Wk", Value: 70},
		},
		Lp: models.Lp{Value: 42, Max: 48},
		Ap: models.Ap{Value: 28, Max: 32},
		B:  models.B{Value: 18},
	}

	// Add skills
	skillNames := []string{
		"Schwimmen", "Klettern", "Reiten", "Laufen", "Springen",
		"Balancieren", "Schleichen", "Sich Verstecken", "Singen",
		"Tanzen", "Musizieren", "Malen", "Kochen", "Erste Hilfe",
		"Himmelskunde", "Pflanzenkunde", "Tierkunde", "Geografie",
		"Geschichte", "Lesen/Schreiben", "Rechnen", "Schätzen",
		"Heilkunde", "Giftmischen", "Alchimie", "Schmieden",
		"Lederarbeiten", "Holzbearbeitung", "Steinmetzkunst",
		"Schlösser öffnen", "Fallen entschärfen", "Taschendiebstahl",
	}
	char.Fertigkeiten = make([]models.SkFertigkeit, len(skillNames))
	for i, name := range skillNames {
		char.Fertigkeiten[i] = models.SkFertigkeit{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: name},
			},
			Fertigkeitswert: 8 + i%12,
			Pp:              i % 6,
		}
	}

	// Add weapons
	weaponNames := []string{
		"Langschwert", "Kurzschwert", "Kriegshammer", "Streitaxt",
		"Speer", "Langbogen", "Armbrust", "Dolch", "Schild",
	}
	char.Waffenfertigkeiten = make([]models.SkWaffenfertigkeit, len(weaponNames))
	for i, name := range weaponNames {
		char.Waffenfertigkeiten[i] = models.SkWaffenfertigkeit{
			SkFertigkeit: models.SkFertigkeit{
				BamortCharTrait: models.BamortCharTrait{
					BamortBase: models.BamortBase{Name: name},
				},
				Fertigkeitswert: 12 + i*2,
				Pp:              i,
			},
		}
	}

	// Add spells
	spellNames := []string{
		"Macht über die belebte Natur", "Macht über das Selbst",
		"Erkennen von Gift", "Heilen von Wunden", "Bannen von Zauberwerk",
		"Schutz vor Dämonen", "Macht über Unbelebtes", "Angst",
		"Unsichtbarkeit", "Feuerlanze", "Eisstrahl", "Blitz",
		"Verwandlung", "Teleportation", "Hellsicht",
	}
	char.Zauber = make([]models.SkZauber, len(spellNames))
	for i, name := range spellNames {
		char.Zauber[i] = models.SkZauber{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: name},
			},
			Bonus: i % 3,
		}
	}

	// Add equipment
	equipmentNames := []string{
		"Rüstung (Leder)", "Helm", "Stiefel", "Umhang", "Rucksack",
		"Seil (20m)", "Fackel (5x)", "Öllampe", "Zunderbüchse",
		"Wasserschlauch", "Proviant (7 Tage)", "Schlafsack",
		"Zelt", "Kochgeschirr", "Werkzeug",
	}
	char.Ausruestung = make([]models.EqAusruestung, len(equipmentNames))
	for i, name := range equipmentNames {
		char.Ausruestung[i] = models.EqAusruestung{
			BamortCharTrait: models.BamortCharTrait{
				BamortBase: models.BamortBase{Name: name},
			},
			Anzahl:  1 + i%3,
			Gewicht: 0.5 + float64(i%10)*0.5,
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

	// Page 1: Stats page with skills
	t.Log("Generating Page 1: Stats...")
	page1Data, err := PreparePaginatedPageData(viewModel, "page1_stats.html", 1, "18.12.2025")
	if err != nil {
		t.Fatalf("Failed to prepare page1 data: %v", err)
	}

	html1, err := loader.RenderTemplateWithInlinedResources("page1_stats.html", page1Data)
	if err != nil {
		t.Fatalf("Failed to render page1: %v", err)
	}

	pdf1, err := renderer.RenderHTMLToPDF(html1)
	if err != nil {
		t.Fatalf("Failed to generate PDF for page1: %v", err)
	}

	// Page 2: Play/Adventure page with weapons
	t.Log("Generating Page 2: Play...")
	page2Data, err := PreparePaginatedPageData(viewModel, "page2_play.html", 2, "18.12.2025")
	if err != nil {
		t.Fatalf("Failed to prepare page2 data: %v", err)
	}

	html2, err := loader.RenderTemplateWithInlinedResources("page2_play.html", page2Data)
	if err != nil {
		t.Fatalf("Failed to render page2: %v", err)
	}

	pdf2, err := renderer.RenderHTMLToPDF(html2)
	if err != nil {
		t.Fatalf("Failed to generate PDF for page2: %v", err)
	}

	// Page 3: Spells page
	t.Log("Generating Page 3: Spells...")
	page3Data, err := PreparePaginatedPageData(viewModel, "page3_spell.html", 3, "18.12.2025")
	if err != nil {
		t.Fatalf("Failed to prepare page3 data: %v", err)
	}

	html3, err := loader.RenderTemplateWithInlinedResources("page3_spell.html", page3Data)
	if err != nil {
		t.Fatalf("Failed to render page3: %v", err)
	}

	pdf3, err := renderer.RenderHTMLToPDF(html3)
	if err != nil {
		t.Fatalf("Failed to generate PDF for page3: %v", err)
	}

	// Page 4: Equipment page
	t.Log("Generating Page 4: Equipment...")
	page4Data, err := PreparePaginatedPageData(viewModel, "page4_equip.html", 4, "18.12.2025")
	if err != nil {
		t.Fatalf("Failed to prepare page4 data: %v", err)
	}

	html4, err := loader.RenderTemplateWithInlinedResources("page4_equip.html", page4Data)
	if err != nil {
		t.Fatalf("Failed to render page4: %v", err)
	}

	pdf4, err := renderer.RenderHTMLToPDF(html4)
	if err != nil {
		t.Fatalf("Failed to generate PDF for page4: %v", err)
	}

	// Save all PDFs to disk
	outputDir := "/tmp/bamort_pdf_test"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	files := []struct {
		name string
		data []byte
	}{
		{"page1_stats.pdf", pdf1},
		{"page2_play.pdf", pdf2},
		{"page3_spell.pdf", pdf3},
		{"page4_equip.pdf", pdf4},
	}

	var filePaths []string
	for _, file := range files {
		path := outputDir + "/" + file.name
		if err := os.WriteFile(path, file.data, 0644); err != nil {
			t.Errorf("Failed to write %s: %v", file.name, err)
			continue
		}
		filePaths = append(filePaths, path)
		t.Logf("✓ Saved %s (%d bytes)", path, len(file.data))
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
	t.Logf("\n✓ All 4 pages generated successfully!")
	t.Logf("  Character: %s (Grade %d)", viewModel.Character.Name, viewModel.Character.Grade)
	t.Logf("  Skills: %d", len(viewModel.Skills))
	t.Logf("  Weapons: %d", len(viewModel.Weapons))
	t.Logf("  Spells: %d", len(viewModel.Spells))
	t.Logf("  Equipment: %d items", len(viewModel.Equipment))
	t.Logf("\n  Output directory: %s", outputDir)
	t.Logf("  Individual PDFs: %d bytes", len(pdf1)+len(pdf2)+len(pdf3)+len(pdf4))
	t.Logf("  Combined PDF: %d bytes", combinedInfo.Size())
}
