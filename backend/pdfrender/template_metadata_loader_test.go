package pdfrender

import (
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

	// Find page1_stats.html and verify its metadata matches the HTML comments
	var page1 *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page1_stats.html" {
			page1 = &templateSet.Templates[i]
			break
		}
	}

	if page1 == nil {
		t.Fatal("page1_stats.html not found in template set")
	}

	// Check that blocks were parsed from HTML
	if len(page1.Metadata.Blocks) == 0 {
		t.Error("Expected blocks in page1 metadata")
	}

	// Verify skills_column1 block
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
		if skillsCol1.MaxItems != 29 {
			t.Errorf("Expected skills_column1 MaxItems 29 (from template), got %d", skillsCol1.MaxItems)
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
	// Check page3_spell.html spells_left should be 20 (from template)
	var page3 *TemplateWithMeta
	for i := range templateSet.Templates {
		if templateSet.Templates[i].Metadata.Name == "page3_spell.html" {
			page3 = &templateSet.Templates[i]
			break
		}
	}

	if page3 == nil {
		t.Fatal("page3_spell.html not found")
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
		// Should be 26 from the template file (<!-- BLOCK: spells_left, TYPE: spells, MAX: 26 -->)
		if spellsLeft.MaxItems != 26 {
			t.Errorf("Expected spells_left MaxItems 26 (from template file), got %d", spellsLeft.MaxItems)
		}
	}
}
