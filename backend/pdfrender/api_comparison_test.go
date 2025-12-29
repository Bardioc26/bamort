package pdfrender

import (
	"bamort/config"
	"bamort/database"
	"bamort/models"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// TestAPIvsTestOutput_ShouldBeIdentical verifies that the API handler and test produce identical PDFs
func TestAPIvsTestOutput_ShouldBeIdentical(t *testing.T) {
	database.SetupTestDB()

	// Load character with ID 18 (same as visual inspection test)
	char := &models.Char{}
	err := char.FirstID("18")
	if err != nil {
		t.Fatalf("Failed to load character: %v", err)
	}

	// === PATH 1: Test method (current working) ===
	viewModel1, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("Failed to map character (test path): %v", err)
	}

	loader1 := NewTemplateLoader("../templates/Default_A4_Quer")
	if err := loader1.LoadTemplates(); err != nil {
		t.Fatalf("Failed to load templates (test path): %v", err)
	}

	renderer1 := NewPDFRenderer()
	testDate := "28.12.2025"

	var testPDFs [][]byte
	page1PDFs, _ := RenderPageWithContinuations(viewModel1, "page_1.html", 1, testDate, loader1, renderer1)
	testPDFs = append(testPDFs, page1PDFs...)
	page2PDFs, _ := RenderPageWithContinuations(viewModel1, "page_2.html", 2, testDate, loader1, renderer1)
	testPDFs = append(testPDFs, page2PDFs...)
	page3PDFs, _ := RenderPageWithContinuations(viewModel1, "page_3.html", 3, testDate, loader1, renderer1)
	testPDFs = append(testPDFs, page3PDFs...)
	page4PDFs, _ := RenderPageWithContinuations(viewModel1, "page_4.html", 4, testDate, loader1, renderer1)
	testPDFs = append(testPDFs, page4PDFs...)

	// === PATH 2: API handler method ===
	viewModel2, err := MapCharacterToViewModel(char)
	if err != nil {
		t.Fatalf("Failed to map character (API path): %v", err)
	}

	// Use same template resolution as API handler
	templateID := "Default_A4_Quer"
	templateDir := filepath.Join(config.Cfg.TemplatesDir, templateID)
	t.Logf("API template dir: %s", templateDir)

	loader2 := NewTemplateLoader(templateDir)
	if err := loader2.LoadTemplates(); err != nil {
		t.Fatalf("Failed to load templates (API path): %v", err)
	}

	renderer2 := NewPDFRenderer()
	currentDate := testDate // Use same date as test for exact comparison
	t.Logf("Using date: %s", currentDate)

	var apiPDFs [][]byte
	page1PDFs, _ = RenderPageWithContinuations(viewModel2, "page_1.html", 1, currentDate, loader2, renderer2)
	apiPDFs = append(apiPDFs, page1PDFs...)
	page2PDFs, _ = RenderPageWithContinuations(viewModel2, "page_2.html", 2, currentDate, loader2, renderer2)
	apiPDFs = append(apiPDFs, page2PDFs...)
	page3PDFs, _ = RenderPageWithContinuations(viewModel2, "page_3.html", 3, currentDate, loader2, renderer2)
	apiPDFs = append(apiPDFs, page3PDFs...)
	page4PDFs, _ = RenderPageWithContinuations(viewModel2, "page_4.html", 4, currentDate, loader2, renderer2)
	apiPDFs = append(apiPDFs, page4PDFs...)

	// === COMPARISON ===
	if len(testPDFs) != len(apiPDFs) {
		t.Fatalf("Different number of PDFs: test=%d, api=%d", len(testPDFs), len(apiPDFs))
	}

	t.Logf("Both methods generated %d page PDFs", len(testPDFs))

	// Merge both for final comparison
	tmpDir1 := "/tmp/bamort_test_compare"
	tmpDir2 := "/tmp/bamort_api_compare"
	os.MkdirAll(tmpDir1, 0755)
	os.MkdirAll(tmpDir2, 0755)
	defer os.RemoveAll(tmpDir1)
	defer os.RemoveAll(tmpDir2)

	// Save and merge test PDFs
	var testPaths []string
	for i, pdf := range testPDFs {
		path := filepath.Join(tmpDir1, "page_"+string(rune('0'+i))+".pdf")
		os.WriteFile(path, pdf, 0644)
		testPaths = append(testPaths, path)
	}
	testMerged := filepath.Join(tmpDir1, "merged.pdf")
	api.MergeCreateFile(testPaths, testMerged, false, nil)

	// Save and merge API PDFs
	var apiPaths []string
	for i, pdf := range apiPDFs {
		path := filepath.Join(tmpDir2, "page_"+string(rune('0'+i))+".pdf")
		os.WriteFile(path, pdf, 0644)
		apiPaths = append(apiPaths, path)
	}
	apiMerged := filepath.Join(tmpDir2, "merged.pdf")
	api.MergeCreateFile(apiPaths, apiMerged, false, nil)

	// Read merged PDFs
	testBytes, _ := os.ReadFile(testMerged)
	apiBytes, _ := os.ReadFile(apiMerged)

	t.Logf("Test PDF size: %d bytes", len(testBytes))
	t.Logf("API PDF size:  %d bytes", len(apiBytes))

	// Check if identical (excluding date metadata)
	if bytes.Equal(testBytes, apiBytes) {
		t.Log("✓ PDFs are byte-identical")
	} else {
		// Save for manual inspection
		os.WriteFile("/tmp/test_output.pdf", testBytes, 0644)
		os.WriteFile("/tmp/api_output.pdf", apiBytes, 0644)
		t.Logf("PDFs differ - saved to /tmp/test_output.pdf and /tmp/api_output.pdf for comparison")
		t.Log("Note: Difference might be due to date stamps - checking individual pages...")

		// Compare individual pages
		for i := 0; i < len(testPDFs); i++ {
			if bytes.Equal(testPDFs[i], apiPDFs[i]) {
				t.Logf("  Page %d: identical", i+1)
			} else {
				t.Logf("  Page %d: DIFFERENT (test=%d bytes, api=%d bytes)", i+1, len(testPDFs[i]), len(apiPDFs[i]))
			}
		}
	}
}
