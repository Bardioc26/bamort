package pdfrender

import (
	"os"
	"testing"
)

func TestLoadTemplateSetFromFiles(t *testing.T) {
	// Test loading template set from actual files
	templateSet, err := LoadTemplateSetFromFiles("../templates/Default_A4_Quer")
	if err != nil {
		t.Fatalf("Failed to load template set: %v", err)
	}

	// Verify we have templates
	if len(templateSet.Templates) == 0 {
		t.Fatal("Expected templates, got none")
	}

	// Find page_1.html and verify its metadata matches the HTML comments
	var page1 *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_1.html" {
			page1 = &templateSet.Templates[i]
			break
		}
	}

	if page1 == nil {
		t.Fatal("page_1.html not found in template set")
	}

	// Check that blocks were parsed from HTML
	if len(page1.Metadata.Blocks) == 0 {
		t.Error("Expected blocks in page1 metadata")
	}

	// Verify skills_column1 block - read expected value directly from template file
	templateContent, err := os.ReadFile("../templates/Default_A4_Quer/page_1.html")
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	expectedBlocks := ParseTemplateMetadata(string(templateContent))
	var expectedSkillsCol1 *BlockMetadata
	for i := range expectedBlocks {
		if expectedBlocks[i].Name == "skills_column1" {
			expectedSkillsCol1 = &expectedBlocks[i]
			break
		}
	}

	if expectedSkillsCol1 == nil {
		t.Fatal("skills_column1 block not found in template file")
	}

	var skillsCol1 *BlockMetadata
	for i := range page1.Metadata.Blocks {
		if page1.Metadata.Blocks[i].Name == "skills_column1" {
			skillsCol1 = &page1.Metadata.Blocks[i]
			break
		}
	}

	if skillsCol1 == nil {
		t.Error("skills_column1 block not found")
	} else {
		// Should match the MAX value in the template comment
		if skillsCol1.MaxItems != expectedSkillsCol1.MaxItems {
			t.Errorf("Expected skills_column1 MaxItems %d (from template), got %d", expectedSkillsCol1.MaxItems, skillsCol1.MaxItems)
		}
		if skillsCol1.ListType != "skills" {
			t.Errorf("Expected ListType 'skills', got '%s'", skillsCol1.ListType)
		}
	}
}

func TestDefaultA4QuerTemplateSet_LoadsFromFiles(t *testing.T) {
	// Test that DefaultA4QuerTemplateSet now loads from actual files
	templateSet := DefaultA4QuerTemplateSet()

	if len(templateSet.Templates) == 0 {
		t.Fatal("Expected templates, got none")
	}

	// Verify metadata comes from template files, not hardcoded
	// Read expected value directly from template file
	templateContent, err := os.ReadFile("../templates/Default_A4_Quer/page_3.html")
	if err != nil {
		t.Fatalf("Failed to read template file: %v", err)
	}

	expectedBlocks := ParseTemplateMetadata(string(templateContent))
	var expectedSpellsLeft *BlockMetadata
	for i := range expectedBlocks {
		if expectedBlocks[i].Name == "spells_left" {
			expectedSpellsLeft = &expectedBlocks[i]
			break
		}
	}

	if expectedSpellsLeft == nil {
		t.Fatal("spells_left block not found in template file")
	}

	var page3 *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page_3.html" {
			page3 = &templateSet.Templates[i]
			break
		}
	}

	if page3 == nil {
		t.Fatal("page_3.html not found")
	}

	var spellsLeft *BlockMetadata
	for i := range page3.Metadata.Blocks {
		if page3.Metadata.Blocks[i].Name == "spells_left" {
			spellsLeft = &page3.Metadata.Blocks[i]
			break
		}
	}

	if spellsLeft == nil {
		t.Error("spells_left block not found")
	} else {
		// Should match the value from the template file
		if spellsLeft.MaxItems != expectedSpellsLeft.MaxItems {
			t.Errorf("Expected spells_left MaxItems %d (from template file), got %d", expectedSpellsLeft.MaxItems, spellsLeft.MaxItems)
		}
	}
}
