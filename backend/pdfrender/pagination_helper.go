package pdfrender

// PreparePaginatedPageData prepares data for rendering a template page with proper pagination
// It takes the full view model and returns PageData with lists split according to template capacity
func PreparePaginatedPageData(viewModel *CharacterSheetViewModel, templateName string, pageNumber int, date string) (*PageData, error) {
	// Get template metadata to determine capacities
	templateSet := DefaultA4QuerTemplateSet()

	pageData := &PageData{
		Character:     viewModel.Character,
		Attributes:    viewModel.Attributes,
		DerivedValues: viewModel.DerivedValues,
		GameResults:   viewModel.GameResults,
		Meta: PageMeta{
			Date:       date,
			PageNumber: pageNumber,
		},
	}

	// For page1_stats.html - paginate skills across two columns
	if templateName == "page1_stats.html" {
		// Get the template metadata
		var template *TemplateMetadata
		for _, tmpl := range templateSet.Templates {
			if tmpl.Metadata.Name == templateName {
				template = &tmpl.Metadata
				break
			}
		}

		if template != nil {
			// Get skill blocks (should be skills_column1 and skills_column2)
			var skillBlocks []BlockMetadata
			for _, block := range template.Blocks {
				if block.ListType == "skills" {
					skillBlocks = append(skillBlocks, block)
				}
			}

			if len(skillBlocks) >= 2 {
				// Calculate how to split skills across columns
				col1Capacity := skillBlocks[0].MaxItems
				col2Capacity := skillBlocks[1].MaxItems

				// Debug logging
				// fmt.Printf("DEBUG: col1Capacity=%d, col2Capacity=%d\n", col1Capacity, col2Capacity)

				col1Skills, col2Skills := SplitSkillsForColumns(viewModel.Skills, col1Capacity, col2Capacity)
				pageData.SkillsColumn1 = col1Skills
				pageData.SkillsColumn2 = col2Skills
				pageData.Skills = viewModel.Skills // Keep for backward compatibility
			} else {
				pageData.Skills = viewModel.Skills
			}
		} else {
			pageData.Skills = viewModel.Skills
		}
	} else if templateName == "page2_play.html" {
		// Limit weapons according to capacity (30)
		pageData.Weapons = viewModel.Weapons
		if len(pageData.Weapons) > 30 {
			pageData.Weapons = pageData.Weapons[:30]
		}
		pageData.Skills = viewModel.Skills
	} else if templateName == "page3_spell.html" {
		// Limit spells according to capacity (24 total: 12+12)
		pageData.Spells = viewModel.Spells
		if len(pageData.Spells) > 24 {
			pageData.Spells = pageData.Spells[:24]
		}
		pageData.MagicItems = viewModel.MagicItems
		if len(pageData.MagicItems) > 5 {
			pageData.MagicItems = pageData.MagicItems[:5]
		}
	} else if templateName == "page4_equip.html" {
		pageData.Equipment = viewModel.Equipment
		if len(pageData.Equipment) > 20 {
			pageData.Equipment = pageData.Equipment[:20]
		}
	}

	return pageData, nil
}

// SplitSkillsForColumns splits skills into two separate lists for two-column layout
// Returns (column1Skills, column2Skills)
func SplitSkillsForColumns(skills []SkillViewModel, col1Max, col2Max int) ([]SkillViewModel, []SkillViewModel) {
	col1 := skills
	if len(col1) > col1Max {
		col1 = col1[:col1Max]
	}

	col2 := []SkillViewModel{}
	if len(skills) > col1Max {
		remaining := skills[col1Max:]
		if len(remaining) > col2Max {
			col2 = remaining[:col2Max]
		} else {
			col2 = remaining
		}
	}

	return col1, col2
}
