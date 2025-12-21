package pdfrender

import (
	"fmt"
	"os"
)

// TemplateMetadata contains information about a template's capacity and requirements
type TemplateMetadata struct {
	Name        string // Template name (e.g., "page1_stats.html")
	PageType    string // "stats", "play", "spell", "equip"
	Description string
	Blocks      []BlockMetadata // List blocks and their capacities
}

// BlockMetadata defines a list block within a template and how many items it can hold
type BlockMetadata struct {
	Name     string // Logical name (e.g., "skills_column1", "weapons_main")
	ListType string // "skills", "weapons", "spells", "equipment", "magicItems", "gameResults"
	MaxItems int    // Maximum number of items this block can display
	Filter   string // Optional filter criteria (e.g., "learned", "unlearned", "languages")
	Column   int    // Column number for multi-column layouts (0 if single column)
}

// TemplateSet contains all templates and their metadata for a format
type TemplateSet struct {
	Name        string // e.g., "Default_A4_Quer"
	Description string
	Templates   []TemplateWithMeta
}

// TemplateWithMeta combines a template with its metadata
type TemplateWithMeta struct {
	Metadata TemplateMetadata
	Path     string // File path to the template
}

// LoadTemplateSetFromFiles loads template metadata by parsing actual template files
func LoadTemplateSetFromFiles(templateDir string) (TemplateSet, error) {
	templateSet := TemplateSet{
		Name:        "Default_A4_Quer",
		Description: "Standard A4 Querformat Charakterbogen",
		Templates:   []TemplateWithMeta{},
	}

	// Define the template files to load
	templateFiles := []struct {
		filename    string
		pageType    string
		description string
	}{
		{"page1_stats.html", "stats", "Statistikseite mit Grundwerten"},
		{"page2_play.html", "play", "Spielbogen mit gelernten Fertigkeiten und Waffen"},
		{"page3_spell.html", "spell", "Zauberseite mit Zauberliste"},
		{"page4_equip.html", "equip", "Ausrüstungsseite"},
	}

	// Load each template file and parse its metadata
	for _, tmplFile := range templateFiles {
		filePath := templateDir + "/" + tmplFile.filename

		// Read template content
		content, err := os.ReadFile(filePath)
		if err != nil {
			return templateSet, fmt.Errorf("failed to read template %s: %w", tmplFile.filename, err)
		}

		// Parse metadata from HTML comments
		blocks := ParseTemplateMetadata(string(content))

		templateSet.Templates = append(templateSet.Templates, TemplateWithMeta{
			Metadata: TemplateMetadata{
				Name:        tmplFile.filename,
				PageType:    tmplFile.pageType,
				Description: tmplFile.description,
				Blocks:      blocks,
			},
			Path: filePath,
		})
	}

	return templateSet, nil
}

// DefaultA4QuerTemplateSet returns the template set for A4 Querformat
// Now loads from actual template files instead of hardcoded values
func DefaultA4QuerTemplateSet() TemplateSet {
	// Try to load from files
	templateSet, err := LoadTemplateSetFromFiles("backend/templates/Default_A4_Quer")
	if err != nil {
		// Fallback to relative path from test directory
		templateSet, err = LoadTemplateSetFromFiles("../templates/Default_A4_Quer")
		if err != nil {
			// Last fallback: return hardcoded defaults
			return getHardcodedTemplateSet()
		}
	}
	return templateSet
}

// getHardcodedTemplateSet returns hardcoded fallback values
func getHardcodedTemplateSet() TemplateSet {
	return TemplateSet{
		Name:        "Default_A4_Quer",
		Description: "Standard A4 Querformat Charakterbogen",
		Templates: []TemplateWithMeta{
			{
				Metadata: TemplateMetadata{
					Name:        "page1_stats.html",
					PageType:    "stats",
					Description: "Statistikseite mit Grundwerten",
					Blocks: []BlockMetadata{
						{
							Name:     "skills_column1",
							ListType: "skills",
							MaxItems: 29,
							Column:   1,
						},
						{
							Name:     "skills_column2",
							ListType: "skills",
							MaxItems: 29,
							Column:   2,
						},
					},
				},
				Path: "templates/Default_A4_Quer/page1_stats.html",
			},
			{
				Metadata: TemplateMetadata{
					Name:        "page2_play.html",
					PageType:    "play",
					Description: "Spielbogen mit gelernten Fertigkeiten und Waffen",
					Blocks: []BlockMetadata{
						{
							Name:     "skills_learned",
							ListType: "skills",
							MaxItems: 24,
							Filter:   "learned",
						},
						{
							Name:     "skills_unlearned",
							ListType: "skills",
							MaxItems: 15,
							Filter:   "unlearned",
						},
						{
							Name:     "skills_languages",
							ListType: "skills",
							MaxItems: 11,
							Filter:   "languages",
						},
						{
							Name:     "weapons_main",
							ListType: "weapons",
							MaxItems: 22,
						},
					},
				},
				Path: "templates/Default_A4_Quer/page2_play.html",
			},
			{
				Metadata: TemplateMetadata{
					Name:        "page3_spell.html",
					PageType:    "spell",
					Description: "Zauberseite mit Zauberliste",
					Blocks: []BlockMetadata{
						{
							Name:     "spells_left",
							ListType: "spells",
							MaxItems: 15,
							Column:   1,
						},
						{
							Name:     "spells_right",
							ListType: "spells",
							MaxItems: 10,
							Column:   2,
						},
						{
							Name:     "magic_items",
							ListType: "magicItems",
							MaxItems: 5,
						},
					},
				},
				Path: "templates/Default_A4_Quer/page3_spell.html",
			},
			{
				Metadata: TemplateMetadata{
					Name:        "page4_equip.html",
					PageType:    "equip",
					Description: "Ausrüstungsseite",
					Blocks: []BlockMetadata{
						{
							Name:     "equipment_sections",
							ListType: "equipment",
							MaxItems: 20,
						},
						{
							Name:     "game_results",
							ListType: "gameResults",
							MaxItems: 10,
						},
					},
				},
				Path: "templates/Default_A4_Quer/page4_equip.html",
			},
		},
	}
}

// GetBlockMetadata returns the metadata for a specific block in a template
func (tm *TemplateMetadata) GetBlockMetadata(blockName string) *BlockMetadata {
	for i := range tm.Blocks {
		if tm.Blocks[i].Name == blockName {
			return &tm.Blocks[i]
		}
	}
	return nil
}

// GetMaxItems returns the maximum items for a specific list type in this template
func (tm *TemplateMetadata) GetMaxItems(listType string, filter string) int {
	total := 0
	for _, block := range tm.Blocks {
		if block.ListType == listType {
			if filter == "" || block.Filter == filter {
				total += block.MaxItems
			}
		}
	}
	return total
}
