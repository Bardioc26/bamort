package pdfrender

import (
	"strings"
	"testing"
)

func TestLoadTemplate_Success(t *testing.T) {
	// Arrange
	loader := NewTemplateLoader("../templates/Default_A4_Quer")

	// Act
	err := loader.LoadTemplates()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that templates were loaded
	if loader.templates == nil {
		t.Error("Expected templates to be loaded, got nil")
	}
}

func TestRenderTemplate_BasicData(t *testing.T) {
	// Arrange
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	err := loader.LoadTemplates()
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	data := &PageData{
		Character: CharacterInfo{
			Name:  "Test Character",
			Grade: 5,
		},
		Attributes: AttributeValues{
			St: 80,
			Gs: 70,
		},
		Meta: PageMeta{
			Date: "18.12.2025",
		},
	}

	// Act
	html, err := loader.RenderTemplate("page1_stats.html", data)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if html == "" {
		t.Error("Expected non-empty HTML, got empty string")
	}

	// Check that template variables were replaced
	if !strings.Contains(html, "Test Character") {
		t.Error("Expected HTML to contain 'Test Character'")
	}
	if !strings.Contains(html, "18.12.2025") {
		t.Error("Expected HTML to contain date '18.12.2025'")
	}
}

func TestGetTemplateMetadata(t *testing.T) {
	// Arrange
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	err := loader.LoadTemplates()
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	// Act
	metadata := loader.GetTemplateMetadata("page3_spell.html")

	// Assert
	if len(metadata) == 0 {
		t.Fatal("Expected metadata blocks, got none")
	}

	// Check for spells_left block (template says MAX: 15)
	leftBlock := GetBlockByName(metadata, "spells_left")
	if leftBlock == nil {
		t.Error("Expected to find 'spells_left' block")
	} else {
		if leftBlock.MaxItems != 15 {
			t.Errorf("Expected spells_left max 15 (from template), got %d", leftBlock.MaxItems)
		}
	}

	// Check for spells_right block (template says MAX: 10)
	rightBlock := GetBlockByName(metadata, "spells_right")
	if rightBlock == nil {
		t.Error("Expected to find 'spells_right' block")
	} else {
		if rightBlock.MaxItems != 10 {
			t.Errorf("Expected spells_right max 10 (from template), got %d", rightBlock.MaxItems)
		}
	}
}

func TestRenderTemplate_WithSkills(t *testing.T) {
	// Arrange
	loader := NewTemplateLoader("../templates/Default_A4_Quer")
	err := loader.LoadTemplates()
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	data := &PageData{
		Character: CharacterInfo{
			Name: "Test",
		},
		SkillsColumn1: []SkillViewModel{
			{Name: "Schwimmen", Value: 10, PracticePoints: 2},
		},
		SkillsColumn2: []SkillViewModel{
			{Name: "Klettern", Value: 8, PracticePoints: 3},
		},
		Meta: PageMeta{
			Date: "18.12.2025",
		},
	}

	// Act
	html, err := loader.RenderTemplate("page1_stats.html", data)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that skills were rendered
	if !strings.Contains(html, "Schwimmen") {
		t.Error("Expected HTML to contain 'Schwimmen'")
	}
	if !strings.Contains(html, "Klettern") {
		t.Error("Expected HTML to contain 'Klettern'")
	}
}

func TestLoadTemplate_InvalidPath(t *testing.T) {
	// Arrange
	loader := NewTemplateLoader("/invalid/path")

	// Act
	err := loader.LoadTemplates()

	// Assert
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}
