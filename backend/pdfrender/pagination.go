package pdfrender

import "fmt"

// SliceList slices a list based on start index and max items
// Returns the sliced list and whether there are more items
func SliceList[T any](fullList []T, startIndex, maxItems int) ([]T, bool) {
	totalCount := len(fullList)
	endIndex := startIndex + maxItems

	if startIndex >= totalCount {
		return []T{}, false
	}

	if endIndex > totalCount {
		endIndex = totalCount
	}

	return fullList[startIndex:endIndex], endIndex < totalCount
}

// PageDistribution represents how data is distributed across pages
type PageDistribution struct {
	TemplateName string                 // Template to use for this page
	PageNumber   int                    // Page number (1-indexed)
	Data         map[string]interface{} // Block name -> data slice
}

// Paginator handles pagination of lists according to template metadata
type Paginator struct {
	templateSet TemplateSet
}

// NewPaginator creates a new paginator with template metadata
func NewPaginator(templateSet TemplateSet) *Paginator {
	return &Paginator{
		templateSet: templateSet,
	}
}

// PaginateSkills splits skills across multiple pages according to template capacity
func (p *Paginator) PaginateSkills(skills []SkillViewModel, templateName string, filter string) ([]PageDistribution, error) {
	template := p.findTemplate(templateName)
	if template == nil {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	blocks := p.getBlocksForType(template, "skills", filter)
	if len(blocks) == 0 {
		return []PageDistribution{}, nil
	}

	return p.paginateList(skills, blocks, templateName, "skills")
}

// PaginateWeapons splits weapons across multiple pages
func (p *Paginator) PaginateWeapons(weapons []WeaponViewModel, templateName string) ([]PageDistribution, error) {
	template := p.findTemplate(templateName)
	if template == nil {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	blocks := p.getBlocksForType(template, "weapons", "")
	if len(blocks) == 0 {
		return []PageDistribution{}, nil
	}

	return p.paginateList(weapons, blocks, templateName, "weapons")
}

// PaginateSpells splits spells across multiple pages and columns
func (p *Paginator) PaginateSpells(spells []SpellViewModel, templateName string) ([]PageDistribution, error) {
	template := p.findTemplate(templateName)
	if template == nil {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	blocks := p.getBlocksForType(template, "spells", "")
	if len(blocks) == 0 {
		return []PageDistribution{}, nil
	}

	return p.paginateList(spells, blocks, templateName, "spells")
}

// PaginateEquipment splits equipment across multiple pages
func (p *Paginator) PaginateEquipment(equipment []EquipmentViewModel, templateName string) ([]PageDistribution, error) {
	template := p.findTemplate(templateName)
	if template == nil {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	blocks := p.getBlocksForType(template, "equipment", "")
	if len(blocks) == 0 {
		return []PageDistribution{}, nil
	}

	return p.paginateList(equipment, blocks, templateName, "equipment")
}

// paginateList is the core pagination algorithm
func (p *Paginator) paginateList(items interface{}, blocks []BlockMetadata, templateName string, listType string) ([]PageDistribution, error) {
	// Convert items to slice length
	itemCount := 0
	switch v := items.(type) {
	case []SkillViewModel:
		itemCount = len(v)
	case []WeaponViewModel:
		itemCount = len(v)
	case []SpellViewModel:
		itemCount = len(v)
	case []EquipmentViewModel:
		itemCount = len(v)
	default:
		return nil, fmt.Errorf("unsupported item type")
	}

	if itemCount == 0 {
		return []PageDistribution{}, nil
	}

	// Calculate total capacity per page
	capacityPerPage := 0
	for _, block := range blocks {
		capacityPerPage += block.MaxItems
	}

	if capacityPerPage == 0 {
		return nil, fmt.Errorf("template has no capacity for list type: %s", listType)
	}

	// Calculate number of pages needed
	pageCount := (itemCount + capacityPerPage - 1) / capacityPerPage

	distributions := make([]PageDistribution, 0, pageCount)
	currentIndex := 0

	for pageNum := 1; pageNum <= pageCount; pageNum++ {
		pageData := make(map[string]interface{})

		// Distribute items across blocks in this page
		for _, block := range blocks {
			if currentIndex >= itemCount {
				// No more items, add empty slice
				pageData[block.Name] = p.createEmptySlice(listType)
				continue
			}

			// Calculate how many items to put in this block
			itemsToTake := block.MaxItems
			if currentIndex+itemsToTake > itemCount {
				itemsToTake = itemCount - currentIndex
			}

			// Extract slice for this block
			blockItems := p.extractSlice(items, currentIndex, itemsToTake)
			pageData[block.Name] = blockItems
			currentIndex += itemsToTake
		}

		distributions = append(distributions, PageDistribution{
			TemplateName: templateName,
			PageNumber:   pageNum,
			Data:         pageData,
		})
	}

	return distributions, nil
}

// findTemplate finds a template by name
func (p *Paginator) findTemplate(templateName string) *TemplateMetadata {
	for _, tmpl := range p.templateSet.Templates {
		if tmpl.Metadata.Name == templateName {
			return &tmpl.Metadata
		}
	}
	return nil
}

// getBlocksForType returns all blocks matching the list type and filter
func (p *Paginator) getBlocksForType(template *TemplateMetadata, listType string, filter string) []BlockMetadata {
	var blocks []BlockMetadata
	for _, block := range template.Blocks {
		if block.ListType == listType {
			if filter == "" || block.Filter == filter {
				blocks = append(blocks, block)
			}
		}
	}
	return blocks
}

// extractSlice extracts a slice of items based on type
func (p *Paginator) extractSlice(items interface{}, start, count int) interface{} {
	switch v := items.(type) {
	case []SkillViewModel:
		end := start + count
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []WeaponViewModel:
		end := start + count
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []SpellViewModel:
		end := start + count
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	case []EquipmentViewModel:
		end := start + count
		if end > len(v) {
			end = len(v)
		}
		return v[start:end]
	}
	return nil
}

// createEmptySlice creates an empty slice of the appropriate type
func (p *Paginator) createEmptySlice(listType string) interface{} {
	switch listType {
	case "skills":
		return []SkillViewModel{}
	case "weapons":
		return []WeaponViewModel{}
	case "spells":
		return []SpellViewModel{}
	case "equipment":
		return []EquipmentViewModel{}
	default:
		return []interface{}{}
	}
}

// CalculatePagesNeeded calculates how many pages are needed for given data
func (p *Paginator) CalculatePagesNeeded(templateName string, listType string, itemCount int) (int, error) {
	template := p.findTemplate(templateName)
	if template == nil {
		return 0, fmt.Errorf("template not found: %s", templateName)
	}

	blocks := p.getBlocksForType(template, listType, "")
	if len(blocks) == 0 {
		return 0, nil
	}

	capacityPerPage := 0
	for _, block := range blocks {
		capacityPerPage += block.MaxItems
	}

	if capacityPerPage == 0 {
		return 0, fmt.Errorf("template has no capacity for list type: %s", listType)
	}

	return (itemCount + capacityPerPage - 1) / capacityPerPage, nil
}
