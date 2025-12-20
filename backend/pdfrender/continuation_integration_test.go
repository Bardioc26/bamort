package pdfrender

import (
	"bamort/models"
	"os"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// TestIntegration_ContinuationPages_ActualFiles tests that continuation pages
// are actually generated and saved as separate PDF files
func TestIntegration_ContinuationPages_ActualFiles(t *testing.T) {
	// Arrange - Create character with 50 skills to force multiple pages
	char := &models.Char{
		BamortBase: models.BamortBase{
			Name: "Test Character With Many Skills",
		},
		Typ:   "Krieger",
		Grad:  10,
		Alter: 35,
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
		Lp: models.Lp{Value: 30, Max: 35},
		Ap: models.Ap{Value: 25, Max: 30},
		B:  models.B{Value: 15},
	}

	// Add 50 skills to force multiple continuation pages
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

	t.Logf("Character has %d skills", len(viewModel.Skills))

	// Load templates
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	if err = loader.LoadTemplates(); err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	renderer := NewPDFRenderer()

	// Act - Render page1 with continuations
	pdfs, err := RenderPageWithContinuations(
		viewModel,
		"page1_stats.html",
		1,
		"20.12.2025",
		loader,
		renderer,
	)

	// Assert
	if err != nil {
		t.Fatalf("Failed to render with continuations: %v", err)
	}

	if len(pdfs) < 2 {
		t.Fatalf("Expected at least 2 PDFs (main + continuation), got %d", len(pdfs))
	}

	t.Logf("Generated %d PDF pages (1 main + %d continuations)", len(pdfs), len(pdfs)-1)

	// Save all PDFs to disk
	outputDir := "/tmp/bamort_continuation_test"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	var filePaths []string
	for i, pdf := range pdfs {
		var filename string
		if i == 0 {
			filename = "page1_stats.pdf"
		} else {
			filename = "page1_stats_continuation_" + string(rune('0'+i)) + ".pdf"
		}

		path := outputDir + "/" + filename
		if err := os.WriteFile(path, pdf, 0644); err != nil {
			t.Errorf("Failed to write %s: %v", filename, err)
			continue
		}
		filePaths = append(filePaths, path)
		t.Logf("✓ Saved %s (%d bytes)", path, len(pdf))

		// Verify PDF is valid
		if string(pdf[0:4]) != "%PDF" {
			t.Errorf("%s does not start with PDF marker", filename)
		}
	}

	// Merge all PDFs into a single file
	combinedPath := outputDir + "/page1_stats_combined.pdf"
	if err := api.MergeCreateFile(filePaths, combinedPath, false, nil); err != nil {
		t.Fatalf("Failed to merge PDFs: %v", err)
	}

	// Get size of combined PDF
	combinedInfo, err := os.Stat(combinedPath)
	if err != nil {
		t.Fatalf("Failed to stat combined PDF: %v", err)
	}

	t.Logf("\n✓ Combined all %d pages into: %s (%d bytes)", len(pdfs), combinedPath, combinedInfo.Size())
	t.Logf("  Output directory: %s", outputDir)
	t.Logf("\n✅ CONTINUATION PAGES SUCCESSFULLY GENERATED AND SAVED!")
}
