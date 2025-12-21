package pdfrender

// GetBlockCapacity gets the MAX capacity for a block from template metadata
// Returns 0 if block not found
func GetBlockCapacity(templateSet *TemplateSet, templateName, blockName string) int {
	for _, tmpl := range templateSet.Templates {
		if tmpl.Metadata.Name == templateName {
			block := GetBlockByName(tmpl.Metadata.Blocks, blockName)
			if block != nil {
				return block.MaxItems
			}
		}
	}
	return 0
}

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
	if templateName == "page_1.html" {
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
				// Fill to capacity to ensure empty rows render in template
				pageData.SkillsColumn1 = FillToCapacity(col1Skills, col1Capacity)
				pageData.SkillsColumn2 = FillToCapacity(col2Skills, col2Capacity)
				pageData.Skills = viewModel.Skills // Keep for backward compatibility
			} else {
				pageData.Skills = viewModel.Skills
			}
		} else {
			pageData.Skills = viewModel.Skills
		}
	} else if templateName == "page_2.html" {
		// Get capacities from template
		weaponsCapacity := GetBlockCapacity(&templateSet, templateName, "weapons_main")
		learnedCapacity := GetBlockCapacity(&templateSet, templateName, "skills_learned")
		languageCapacity := GetBlockCapacity(&templateSet, templateName, "skills_languages")

		// Limit and fill weapons to capacity
		weapons := viewModel.Weapons
		if weaponsCapacity > 0 && len(weapons) > weaponsCapacity {
			weapons = weapons[:weaponsCapacity]
		}
		if weaponsCapacity > 0 {
			pageData.Weapons = FillToCapacity(weapons, weaponsCapacity)
		} else {
			pageData.Weapons = weapons
		}

		// Filter skills by category for page2 blocks
		var learnedSkills, languageSkills []SkillViewModel
		for _, skill := range viewModel.Skills {
			if skill.Category == "Sprache" {
				languageSkills = append(languageSkills, skill)
			} else if skill.IsLearned {
				learnedSkills = append(learnedSkills, skill)
			}
		}

		// Apply capacity limits
		if learnedCapacity > 0 && len(learnedSkills) > learnedCapacity {
			learnedSkills = learnedSkills[:learnedCapacity]
		}
		if languageCapacity > 0 && len(languageSkills) > languageCapacity {
			languageSkills = languageSkills[:languageCapacity]
		}

		// Fill to capacity to ensure empty rows render
		if learnedCapacity > 0 {
			pageData.SkillsLearned = FillToCapacity(learnedSkills, learnedCapacity)
		} else {
			pageData.SkillsLearned = learnedSkills
		}
		if languageCapacity > 0 {
			pageData.SkillsLanguage = FillToCapacity(languageSkills, languageCapacity)
		} else {
			pageData.SkillsLanguage = languageSkills
		}
		pageData.Skills = viewModel.Skills // Keep for backward compatibility
	} else if templateName == "page_3.html" {
		// Get capacities from template
		spellsLeftCapacity := GetBlockCapacity(&templateSet, templateName, "spells_left")
		spellsRightCapacity := GetBlockCapacity(&templateSet, templateName, "spells_right")
		magicItemsCapacity := GetBlockCapacity(&templateSet, templateName, "magic_items")

		// Split spells into left and right columns
		leftSpells, rightSpells := SplitSkillsIntoColumns(viewModel.Spells, spellsLeftCapacity, spellsRightCapacity)
		// Fill to capacity to ensure empty rows render
		if spellsLeftCapacity > 0 {
			pageData.SpellsLeft = FillToCapacity(leftSpells, spellsLeftCapacity)
		} else {
			pageData.SpellsLeft = leftSpells
		}
		if spellsRightCapacity > 0 {
			pageData.SpellsRight = FillToCapacity(rightSpells, spellsRightCapacity)
		} else {
			pageData.SpellsRight = rightSpells
		}
		pageData.Spells = viewModel.Spells // Keep for backward compatibility

		// Limit and fill magic items
		magicItems := viewModel.MagicItems
		if magicItemsCapacity > 0 && len(magicItems) > magicItemsCapacity {
			magicItems = magicItems[:magicItemsCapacity]
		}
		if magicItemsCapacity > 0 {
			pageData.MagicItems = FillToCapacity(magicItems, magicItemsCapacity)
		} else {
			pageData.MagicItems = magicItems
		}
	} else if templateName == "page_4.html" {
		// Page 4 needs ALL equipment to properly render containers
		// The template has complex logic showing containers on left, worn items and container sections on right
		// Don't truncate based on capacity - let the template handle all items
		pageData.Equipment = viewModel.Equipment
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

// SplitSkillsIntoColumns splits spells/items into two columns (generic for any slice type)
// Returns (column1, column2)
func SplitSkillsIntoColumns[T any](items []T, col1Max, col2Max int) ([]T, []T) {
	col1 := items
	if len(col1) > col1Max {
		col1 = col1[:col1Max]
	}

	col2 := []T{}
	if len(items) > col1Max {
		remaining := items[col1Max:]
		if len(remaining) > col2Max {
			col2 = remaining[:col2Max]
		} else {
			col2 = remaining
		}
	}

	return col1, col2
}
