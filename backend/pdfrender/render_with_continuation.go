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

	// Build data map from view model
	dataMap := map[string]interface{}{
		"skills":     viewModel.Skills,
		"weapons":    viewModel.Weapons,
		"spells":     viewModel.Spells,
		"equipment":  viewModel.Equipment,
		"magicItems": viewModel.MagicItems,
	}

	// Use unified pagination for all templates
	distributions, err := paginator.PaginateMultiList(dataMap, templateName)
	if err != nil {
		return nil, fmt.Errorf("failed to paginate: %w", err)
	}

	// If no distributions (empty data), render single empty page
	if len(distributions) == 0 {
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

	// Render each distributed page
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

		// Populate page data from distribution
		populatePageDataFromDistribution(pageData, dist)

		// Render the page
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

// populatePageDataFromDistribution populates PageData from a distribution
// This replaces the hardcoded switch statements for each template type
func populatePageDataFromDistribution(pageData *PageData, dist PageDistribution) {
	// Populate data based on block names in distribution
	for blockName, data := range dist.Data {
		switch blockName {
		// Skills blocks
		case "skills_column1":
			if skills, ok := data.([]SkillViewModel); ok {
				pageData.SkillsColumn1 = skills
				pageData.Skills = append(pageData.Skills, skills...)
			}
		case "skills_column2":
			if skills, ok := data.([]SkillViewModel); ok {
				pageData.SkillsColumn2 = skills
				pageData.Skills = append(pageData.Skills, skills...)
			}
		case "skills_column3":
			if skills, ok := data.([]SkillViewModel); ok {
				pageData.SkillsColumn3 = skills
				pageData.Skills = append(pageData.Skills, skills...)
			}
		case "skills_column4":
			if skills, ok := data.([]SkillViewModel); ok {
				pageData.SkillsColumn4 = skills
				pageData.Skills = append(pageData.Skills, skills...)
			}
		case "skills_learned":
			if skills, ok := data.([]SkillViewModel); ok {
				pageData.SkillsLearned = skills
			}
		case "skills_unlearned":
			if skills, ok := data.([]SkillViewModel); ok {
				// Add to general Skills list for template compatibility
				pageData.Skills = append(pageData.Skills, skills...)
			}
		case "skills_languages":
			if skills, ok := data.([]SkillViewModel); ok {
				pageData.SkillsLanguage = skills
			}

		// Weapons blocks
		case "weapons_main":
			if weapons, ok := data.([]WeaponViewModel); ok {
				pageData.Weapons = weapons
			}

		// Spells blocks
		case "spells_column1":
			if spells, ok := data.([]SpellViewModel); ok {
				pageData.SpellsLeft = spells
				pageData.Spells = append(pageData.Spells, spells...)
			}
		case "spells_column2":
			if spells, ok := data.([]SpellViewModel); ok {
				pageData.SpellsRight = spells
				pageData.Spells = append(pageData.Spells, spells...)
			}

		// Equipment blocks
		case "equipment_worn":
			if equipment, ok := data.([]EquipmentViewModel); ok {
				pageData.Equipment = append(pageData.Equipment, equipment...)
			}
		case "equipment_carried":
			if equipment, ok := data.([]EquipmentViewModel); ok {
				pageData.Equipment = append(pageData.Equipment, equipment...)
			}

		// Magic items
		case "magic_items":
			if items, ok := data.([]MagicItemViewModel); ok {
				pageData.MagicItems = items
			}
		}
	}
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
