package pdfrender

import (
	"fmt"
	"strings"
)

// GenerateContinuationTemplateName creates a continuation template name
// Example: "page1_stats.html" + pageNum 2 -> "page1.2_stats.html"
func GenerateContinuationTemplateName(originalTemplate string, pageNum int) string {
	if pageNum == 1 {
		return originalTemplate
	}

	// Split template name at first underscore to insert page continuation number
	// Example: "page1_stats.html" -> "page1" + "_stats.html"
	parts := strings.SplitN(originalTemplate, "_", 2)
	if len(parts) != 2 {
		// Fallback: just append .2, .3, etc before extension
		ext := ".html"
		base := strings.TrimSuffix(originalTemplate, ext)
		return fmt.Sprintf("%s.%d%s", base, pageNum, ext)
	}

	// Extract page number and base name
	// "page1" -> "page" + "1"
	baseName := parts[0]
	suffix := parts[1]

	// Insert continuation number: "page1.2_stats.html"
	return fmt.Sprintf("%s.%d_%s", baseName, pageNum, suffix)
}

// ExtractBaseTemplateName extracts the base template name from a continuation template
// Example: "page1.2_stats.html" -> "page1_stats.html"
func ExtractBaseTemplateName(templateName string) string {
	// Check if it's a continuation template (contains .N_ pattern)
	parts := strings.SplitN(templateName, "_", 2)
	if len(parts) != 2 {
		return templateName
	}

	baseName := parts[0]
	suffix := parts[1]

	// Check if baseName contains a dot followed by a number (e.g., "page1.2")
	dotIdx := strings.LastIndex(baseName, ".")
	if dotIdx == -1 {
		return templateName // Not a continuation template
	}

	// Verify the part after the dot is a number
	numPart := baseName[dotIdx+1:]
	if len(numPart) == 0 {
		return templateName
	}

	for _, c := range numPart {
		if c < '0' || c > '9' {
			return templateName // Not a number, not a continuation template
		}
	}

	// It's a continuation template, return the base name
	basePrefix := baseName[:dotIdx]
	return fmt.Sprintf("%s_%s", basePrefix, suffix)
}

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

		// Determine template name - use continuation naming for pages 2+
		pageTemplateName := GenerateContinuationTemplateName(templateName, pageNum)

		distributions = append(distributions, PageDistribution{
			TemplateName: pageTemplateName,
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

// PaginatePage2PlayLists handles pagination for page2_play.html which has both skills and weapons
// Skills and weapons overflow together - if either overflows, create continuation pages with remaining items from both
func (p *Paginator) PaginatePage2PlayLists(skills []SkillViewModel, weapons []WeaponViewModel, templateName string) ([]PageDistribution, error) {
	template := p.findTemplate(templateName)
	if template == nil {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	// Get capacities for each block type
	learnedCap := GetBlockCapacity(&p.templateSet, templateName, "skills_learned")
	unlearnedCap := GetBlockCapacity(&p.templateSet, templateName, "skills_unlearned")
	languageCap := GetBlockCapacity(&p.templateSet, templateName, "skills_languages")
	weaponsCap := GetBlockCapacity(&p.templateSet, templateName, "weapons_main")

	// Filter skills into categories
	var learnedSkills, unlearnedSkills, languageSkills []SkillViewModel
	for _, skill := range skills {
		if skill.Category == "Sprache" {
			languageSkills = append(languageSkills, skill)
		} else if skill.IsLearned {
			learnedSkills = append(learnedSkills, skill)
		} else {
			unlearnedSkills = append(unlearnedSkills, skill)
		}
	}

	// Track current position in each list
	learnedIdx := 0
	unlearnedIdx := 0
	languageIdx := 0
	weaponsIdx := 0

	distributions := []PageDistribution{}
	pageNum := 1

	// Continue creating pages while there are remaining items in any list
	for learnedIdx < len(learnedSkills) || unlearnedIdx < len(unlearnedSkills) ||
		languageIdx < len(languageSkills) || weaponsIdx < len(weapons) {

		pageData := make(map[string]interface{})

		// Add learned skills for this page
		learnedEnd := learnedIdx + learnedCap
		if learnedEnd > len(learnedSkills) {
			learnedEnd = len(learnedSkills)
		}
		pageData["skills_learned"] = learnedSkills[learnedIdx:learnedEnd]
		learnedIdx = learnedEnd

		// Add unlearned skills for this page
		unlearnedEnd := unlearnedIdx + unlearnedCap
		if unlearnedEnd > len(unlearnedSkills) {
			unlearnedEnd = len(unlearnedSkills)
		}
		pageData["skills_unlearned"] = unlearnedSkills[unlearnedIdx:unlearnedEnd]
		unlearnedIdx = unlearnedEnd

		// Add language skills for this page
		languageEnd := languageIdx + languageCap
		if languageEnd > len(languageSkills) {
			languageEnd = len(languageSkills)
		}
		pageData["skills_languages"] = languageSkills[languageIdx:languageEnd]
		languageIdx = languageEnd

		// Add weapons for this page
		weaponsEnd := weaponsIdx + weaponsCap
		if weaponsEnd > len(weapons) {
			weaponsEnd = len(weapons)
		}
		pageData["weapons_main"] = weapons[weaponsIdx:weaponsEnd]
		weaponsIdx = weaponsEnd

		// Create page distribution
		pageTemplateName := GenerateContinuationTemplateName(templateName, pageNum)
		distributions = append(distributions, PageDistribution{
			TemplateName: pageTemplateName,
			PageNumber:   pageNum,
			Data:         pageData,
		})

		pageNum++
	}

	return distributions, nil
}
