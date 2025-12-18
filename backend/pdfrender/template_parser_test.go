package pdfrender

import "testing"

func TestParseTemplateMetadata(t *testing.T) {
	// Arrange
	templateContent := `
<!DOCTYPE html>
<html>
<body>
    <!-- BLOCK: spells_left, TYPE: spells, MAX: 12 -->
    <table>
        {{range .Spells}}
        {{end}}
    </table>
    
    <!-- BLOCK: spells_right, TYPE: spells, MAX: 10 -->
    <table>
        {{range .Spells}}
        {{end}}
    </table>
    
    <!-- BLOCK: magic_items, TYPE: magicItems, MAX: 5 -->
    <table>
        {{range .MagicItems}}
        {{end}}
    </table>
</body>
</html>
`

	// Act
	blocks := ParseTemplateMetadata(templateContent)

	// Assert
	if len(blocks) != 3 {
		t.Fatalf("Expected 3 blocks, got %d", len(blocks))
	}

	// Check first block
	if blocks[0].Name != "spells_left" {
		t.Errorf("Expected name 'spells_left', got '%s'", blocks[0].Name)
	}
	if blocks[0].ListType != "spells" {
		t.Errorf("Expected type 'spells', got '%s'", blocks[0].ListType)
	}
	if blocks[0].MaxItems != 12 {
		t.Errorf("Expected max 12, got %d", blocks[0].MaxItems)
	}

	// Check second block
	if blocks[1].Name != "spells_right" {
		t.Errorf("Expected name 'spells_right', got '%s'", blocks[1].Name)
	}
	if blocks[1].MaxItems != 10 {
		t.Errorf("Expected max 10, got %d", blocks[1].MaxItems)
	}

	// Check third block
	if blocks[2].Name != "magic_items" {
		t.Errorf("Expected name 'magic_items', got '%s'", blocks[2].Name)
	}
	if blocks[2].ListType != "magicItems" {
		t.Errorf("Expected type 'magicItems', got '%s'", blocks[2].ListType)
	}
	if blocks[2].MaxItems != 5 {
		t.Errorf("Expected max 5, got %d", blocks[2].MaxItems)
	}
}

func TestParseTemplateMetadata_WithFilter(t *testing.T) {
	// Arrange
	templateContent := `
<!-- BLOCK: skills_learned, TYPE: skills, MAX: 24, FILTER: learned -->
<table>{{range .Skills}}{{end}}</table>
`

	// Act
	blocks := ParseTemplateMetadata(templateContent)

	// Assert
	if len(blocks) != 1 {
		t.Fatalf("Expected 1 block, got %d", len(blocks))
	}

	if blocks[0].Name != "skills_learned" {
		t.Errorf("Expected name 'skills_learned', got '%s'", blocks[0].Name)
	}
	if blocks[0].Filter != "learned" {
		t.Errorf("Expected filter 'learned', got '%s'", blocks[0].Filter)
	}
	if blocks[0].MaxItems != 24 {
		t.Errorf("Expected max 24, got %d", blocks[0].MaxItems)
	}
}

func TestGetBlockByName(t *testing.T) {
	// Arrange
	blocks := []BlockMetadata{
		{Name: "skills_left", ListType: "skills", MaxItems: 32},
		{Name: "skills_right", ListType: "skills", MaxItems: 32},
		{Name: "weapons_main", ListType: "weapons", MaxItems: 30},
	}

	// Act
	block := GetBlockByName(blocks, "weapons_main")

	// Assert
	if block == nil {
		t.Fatal("Expected block, got nil")
	}
	if block.Name != "weapons_main" {
		t.Errorf("Expected 'weapons_main', got '%s'", block.Name)
	}
	if block.MaxItems != 30 {
		t.Errorf("Expected max 30, got %d", block.MaxItems)
	}
}

func TestGetBlocksByType(t *testing.T) {
	// Arrange
	blocks := []BlockMetadata{
		{Name: "skills_left", ListType: "skills", MaxItems: 32},
		{Name: "skills_right", ListType: "skills", MaxItems: 32},
		{Name: "weapons_main", ListType: "weapons", MaxItems: 30},
	}

	// Act
	skillBlocks := GetBlocksByType(blocks, "skills")

	// Assert
	if len(skillBlocks) != 2 {
		t.Fatalf("Expected 2 skill blocks, got %d", len(skillBlocks))
	}
	if skillBlocks[0].ListType != "skills" {
		t.Errorf("Expected type 'skills', got '%s'", skillBlocks[0].ListType)
	}
	if skillBlocks[1].ListType != "skills" {
		t.Errorf("Expected type 'skills', got '%s'", skillBlocks[1].ListType)
	}
}
