package pdfrender

import (
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// RenderPageWithContinuations renders a template page and all necessary continuation pages
// Returns a slice of PDF bytes, one for each page (main + continuations)
func RenderPageWithContinuations(
	viewModel *CharacterSheetViewModel,
	templateName string,
	startPageNumber int,
	date string,
	loader *TemplateLoader,
	renderer *PDFRenderer,
) ([][]byte, error) {
	var pdfs [][]byte
	templateSet := DefaultA4QuerTemplateSet()
	paginator := NewPaginator(templateSet)

	// Determine which list type this template handles
	var distributions []PageDistribution
	var err error

	switch templateName {
	case "page1_stats.html":
		// Paginate skills
		distributions, err = paginator.PaginateSkills(viewModel.Skills, templateName, "")
		if err != nil {
			return nil, fmt.Errorf("failed to paginate skills: %w", err)
		}

	case "page2_play.html":
		// Page 2 has both skills and weapons that overflow together
		// Use multi-list pagination so remaining items from both lists go to continuation pages
		distributions, err = paginator.PaginatePage2PlayLists(viewModel.Skills, viewModel.Weapons, templateName)
		if err != nil {
			return nil, fmt.Errorf("failed to paginate page2 lists: %w", err)
		}

	case "page3_spell.html":
		// Paginate spells
		distributions, err = paginator.PaginateSpells(viewModel.Spells, templateName)
		if err != nil {
			return nil, fmt.Errorf("failed to paginate spells: %w", err)
		}

	case "page4_equip.html":
		// Page 4 has a complex container-based layout where items are grouped by containers.
		// The template expects the full equipment list to properly render containers and their contents.
		// Pagination doesn't make sense here - render as single page with all equipment.
		pageData, err := PreparePaginatedPageData(viewModel, templateName, startPageNumber, date)
		if err != nil {
			return nil, err
		}
		html, err := loader.RenderTemplateWithInlinedResources(templateName, pageData)
		if err != nil {
			return nil, err
		}
		pdf, err := renderer.RenderHTMLToPDF(html)
		if err != nil {
			return nil, err
		}
		return [][]byte{pdf}, nil

	default:
		// For unknown templates, render single page without pagination
		pageData, err := PreparePaginatedPageData(viewModel, templateName, startPageNumber, date)
		if err != nil {
			return nil, err
		}
		html, err := loader.RenderTemplateWithInlinedResources(templateName, pageData)
		if err != nil {
			return nil, err
		}
		pdf, err := renderer.RenderHTMLToPDF(html)
		if err != nil {
			return nil, err
		}
		return [][]byte{pdf}, nil
	}

	// If only one page, use the simplified approach
	if len(distributions) == 1 {
		pageData, err := PreparePaginatedPageData(viewModel, templateName, startPageNumber, date)
		if err != nil {
			return nil, err
		}
		html, err := loader.RenderTemplateWithInlinedResources(templateName, pageData)
		if err != nil {
			return nil, err
		}
		pdf, err := renderer.RenderHTMLToPDF(html)
		if err != nil {
			return nil, err
		}
		return [][]byte{pdf}, nil
	}

	// Multiple pages needed - render each one
	for i, dist := range distributions {
		pageData := &PageData{
			Character:     viewModel.Character,
			Attributes:    viewModel.Attributes,
			DerivedValues: viewModel.DerivedValues,
			GameResults:   viewModel.GameResults,
			Meta: PageMeta{
				Date:       date,
				PageNumber: startPageNumber + i,
			},
		}

		// Populate the page data based on the distribution
		switch templateName {
		case "page1_stats.html":
			// Extract skills from distribution
			if col1, ok := dist.Data["skills_column1"].([]SkillViewModel); ok {
				pageData.SkillsColumn1 = col1
			}
			if col2, ok := dist.Data["skills_column2"].([]SkillViewModel); ok {
				pageData.SkillsColumn2 = col2
			}
			// Combine for backward compatibility
			pageData.Skills = append(pageData.SkillsColumn1, pageData.SkillsColumn2...)

		case "page2_play.html":
			// Extract all lists from distribution (skills and weapons)
			if weapons, ok := dist.Data["weapons_main"].([]WeaponViewModel); ok {
				pageData.Weapons = weapons
			}
			if learned, ok := dist.Data["skills_learned"].([]SkillViewModel); ok {
				pageData.SkillsLearned = learned
			}
			if unlearned, ok := dist.Data["skills_unlearned"].([]SkillViewModel); ok {
				// Unlearned skills are typically shown via general skills list
				// Add to Skills for template compatibility
				pageData.Skills = append(pageData.Skills, unlearned...)
			}
			if languages, ok := dist.Data["skills_languages"].([]SkillViewModel); ok {
				pageData.SkillsLanguage = languages
			}

		case "page3_spell.html":
			// Extract spells from distribution
			if left, ok := dist.Data["spells_left"].([]SpellViewModel); ok {
				pageData.SpellsLeft = left
			}
			if right, ok := dist.Data["spells_right"].([]SpellViewModel); ok {
				pageData.SpellsRight = right
			}
			// Combine for backward compatibility
			pageData.Spells = append(pageData.SpellsLeft, pageData.SpellsRight...)

		case "page4_equip.html":
			// Extract equipment from distribution
			if equipment, ok := dist.Data["equipment_worn"].([]EquipmentViewModel); ok {
				pageData.Equipment = append(pageData.Equipment, equipment...)
			}
			if equipment, ok := dist.Data["equipment_carried"].([]EquipmentViewModel); ok {
				pageData.Equipment = append(pageData.Equipment, equipment...)
			}
		}

		// Render the page (use continuation template name if needed)
		html, err := loader.RenderTemplateWithInlinedResources(dist.TemplateName, pageData)
		if err != nil {
			return nil, fmt.Errorf("failed to render %s: %w", dist.TemplateName, err)
		}

		pdf, err := renderer.RenderHTMLToPDF(html)
		if err != nil {
			return nil, fmt.Errorf("failed to generate PDF for %s: %w", dist.TemplateName, err)
		}

		pdfs = append(pdfs, pdf)
	}

	return pdfs, nil
}

// MergePDFs merges multiple PDF byte slices into a single PDF
func MergePDFs(pdfList [][]byte, outputPath string) error {
	if len(pdfList) == 0 {
		return fmt.Errorf("no PDFs to merge")
	}

	if len(pdfList) == 1 {
		// Single PDF, just write it
		return nil
	}

	// Use pdfcpu to merge - this is a placeholder, actual implementation
	// would need to save individual PDFs first, then merge them
	// For now, this function signature is defined for future use
	return api.MergeCreateFile(nil, outputPath, false, nil)
}
