package pdfrender

import (
	"fmt"
	"strings"
)

// GenerateContinuationTemplateName creates a continuation template name
// Example: "page_1.html" + pageNum 2 -> "page_1.2.html"
// Note: All continuation pages (2, 3, 4, ...) use the same template name: page_1.2.html
func GenerateContinuationTemplateName(originalTemplate string, pageNum int) string {
	if pageNum == 1 {
		return originalTemplate
	}

	// All continuation pages use .2 template (page_1.2, page_2.2, etc.)
	// NOT page_1.3, page_1.4, etc.

	// New format: "page_1.html" -> "page_1.2.html"
	// Pattern: page_N.html where N is a number
	ext := ".html"
	base := strings.TrimSuffix(originalTemplate, ext)

	// Check if it's already a continuation (has .2 in it)
	if strings.Contains(base, ".2") {
		return originalTemplate
	}

	// Append .2 before .html
	return fmt.Sprintf("%s.2%s", base, ext)
}

// ExtractBaseTemplateName extracts the base template name from a continuation template
// Example: "page_1.2.html" -> "page_1.html"
func ExtractBaseTemplateName(templateName string) string {
	// New format: "page_1.2.html" or "page_1.10.html" -> "page_1.html"
	// Pattern: ends with .N.html where N is any number
	ext := ".html"
	if !strings.HasSuffix(templateName, ext) {
		return templateName
	}

	// Remove .html
	base := strings.TrimSuffix(templateName, ext)

	// Find the last dot
	lastDotIdx := strings.LastIndex(base, ".")
	if lastDotIdx == -1 {
		return templateName // No dot found, not a continuation
	}

	// Check if everything after the last dot is a number
	numPart := base[lastDotIdx+1:]
	if len(numPart) == 0 {
		return templateName
	}

	for _, c := range numPart {
		if c < '0' || c > '9' {
			return templateName // Not a number, not a continuation template
		}
	}

	// It's a continuation template, return the base name
	return base[:lastDotIdx] + ext
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

// PaginateMultiList is a unified pagination function that handles multiple list types
// It replaces PaginateSkills, PaginateSpells, and PaginatePage2PlayLists
// dataMap keys: "skills", "weapons", "spells", "equipment", "magicItems"
func (p *Paginator) PaginateMultiList(dataMap map[string]interface{}, templateName string) ([]PageDistribution, error) {
	template := p.findTemplate(templateName)
	if template == nil {
		return nil, fmt.Errorf("template not found: %s", templateName)
	}

	// Build filtered lists for each unique list type + filter combination
	type listTracker struct {
		items      interface{}
		currentIdx int
		totalCount int
	}

	// Track by "listType:filter" to avoid duplicates
	listTrackers := make(map[string]*listTracker)

	// Scan both base template and continuation template to find all unique listType:filter combinations
	templatesToScan := []*TemplateMetadata{template}
	continuationName := GenerateContinuationTemplateName(templateName, 2)
	if contTemplate := p.findTemplate(continuationName); contTemplate != nil {
		templatesToScan = append(templatesToScan, contTemplate)
	}

	// First pass: create filtered lists for each unique listType+filter combination
	for _, tmpl := range templatesToScan {
		for _, block := range tmpl.Blocks {
			// Get the source data based on block's ListType
			var sourceData interface{}
			switch block.ListType {
			case "skills":
				sourceData = dataMap["skills"]
			case "weapons":
				sourceData = dataMap["weapons"]
			case "spells":
				sourceData = dataMap["spells"]
			case "equipment":
				sourceData = dataMap["equipment"]
			case "magicItems":
				sourceData = dataMap["magicItems"]
			default:
				continue
			}

			if sourceData == nil {
				continue
			}

			// Create unique key for this list type + filter combination
			trackerKey := block.ListType
			if block.Filter != "" {
				trackerKey += ":" + block.Filter
			}

			// Only create tracker once per unique combination
			if _, exists := listTrackers[trackerKey]; !exists {
				// Apply filter if specified
				filteredItems := p.applyFilter(sourceData, block.Filter)
				itemCount := p.getItemCount(filteredItems)

				if itemCount > 0 {
					listTrackers[trackerKey] = &listTracker{
						items:      filteredItems,
						currentIdx: 0,
						totalCount: itemCount,
					}
				}
			}
		}
	}

	// If all lists are empty, return empty result
	if len(listTrackers) == 0 {
		return []PageDistribution{}, nil
	}

	// Generate pages until all items are distributed
	distributions := []PageDistribution{}
	pageNum := 1

	for {
		// Check if there are any remaining items
		hasRemainingItems := false
		for _, tracker := range listTrackers {
			if tracker.currentIdx < tracker.totalCount {
				hasRemainingItems = true
				break
			}
		}

		if !hasRemainingItems {
			break
		}

		// Determine template name for this page
		pageTemplateName := GenerateContinuationTemplateName(templateName, pageNum)

		// Load the correct template metadata for this page
		// Fall back to base template if continuation template doesn't exist
		currentTemplate := p.findTemplate(pageTemplateName)
		if currentTemplate == nil {
			currentTemplate = template
		}

		// Create page data
		pageData := make(map[string]interface{})

		// Distribute items to each block for this page
		for _, block := range currentTemplate.Blocks {
			// Get tracker for this block's list type + filter
			trackerKey := block.ListType
			if block.Filter != "" {
				trackerKey += ":" + block.Filter
			}

			tracker, exists := listTrackers[trackerKey]
			if !exists {
				// Skip blocks with NOEMPTY flag when they have no data
				if block.NoEmpty {
					continue
				}
				// Block has no data, fill with empty items up to MAX
				pageData[block.Name] = p.createEmptySliceWithCapacity(block.ListType, block.MaxItems)
				continue
			}

			// Calculate how many items to take for this block
			itemsToTake := block.MaxItems
			remaining := tracker.totalCount - tracker.currentIdx
			if itemsToTake > remaining {
				itemsToTake = remaining
			}

			// Skip NOEMPTY blocks when no items remain
			if block.NoEmpty && itemsToTake == 0 {
				continue
			}
			// For regular blocks with no items remaining, fill with empty rows
			if itemsToTake == 0 {
				pageData[block.Name] = p.createEmptySliceWithCapacity(block.ListType, block.MaxItems)
				continue
			}
			// Extract slice for this block
			blockItems := p.extractSlice(tracker.items, tracker.currentIdx, itemsToTake)

			// Always fill to capacity (both normal and NOEMPTY blocks get filled when they have items)
			blockItems = p.fillSliceToCapacity(blockItems, block.MaxItems)

			pageData[block.Name] = blockItems
			tracker.currentIdx += itemsToTake
		}

		// Template name was already determined at the start of the loop
		distributions = append(distributions, PageDistribution{
			TemplateName: pageTemplateName,
			PageNumber:   pageNum,
			Data:         pageData,
		})

		pageNum++
	}

	return distributions, nil
}

// applyFilter filters a list based on filter criteria
func (p *Paginator) applyFilter(items interface{}, filter string) interface{} {
	if filter == "" {
		return items
	}

	switch v := items.(type) {
	case []SkillViewModel:
		filtered := []SkillViewModel{}
		for _, skill := range v {
			include := false
			switch filter {
			case "learned":
				include = skill.IsLearned && skill.Category != "Sprache"
			case "unlearned":
				include = !skill.IsLearned && skill.Category != "Sprache"
			case "language", "languages":
				include = skill.Category == "Sprache"
			default:
				include = true
			}
			if include {
				filtered = append(filtered, skill)
			}
		}
		return filtered
	default:
		// No filtering for other types
		return items
	}
}

// getItemCount returns the count of items in a list
func (p *Paginator) getItemCount(items interface{}) int {
	switch v := items.(type) {
	case []SkillViewModel:
		return len(v)
	case []WeaponViewModel:
		return len(v)
	case []SpellViewModel:
		return len(v)
	case []EquipmentViewModel:
		return len(v)
	case []MagicItemViewModel:
		return len(v)
	default:
		return 0
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
	case []MagicItemViewModel:
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
	case "magicItems":
		return []MagicItemViewModel{}
	default:
		return []interface{}{}
	}
}

// createEmptySliceWithCapacity creates an empty slice filled to capacity with zero values
func (p *Paginator) createEmptySliceWithCapacity(listType string, capacity int) interface{} {
	switch listType {
	case "skills":
		return FillToCapacity([]SkillViewModel{}, capacity)
	case "weapons":
		return FillToCapacity([]WeaponViewModel{}, capacity)
	case "spells":
		return FillToCapacity([]SpellViewModel{}, capacity)
	case "equipment":
		return FillToCapacity([]EquipmentViewModel{}, capacity)
	case "magicItems":
		return FillToCapacity([]MagicItemViewModel{}, capacity)
	default:
		return []interface{}{}
	}
}

// fillSliceToCapacity fills an existing slice to the specified capacity
func (p *Paginator) fillSliceToCapacity(items interface{}, capacity int) interface{} {
	switch v := items.(type) {
	case []SkillViewModel:
		return FillToCapacity(v, capacity)
	case []WeaponViewModel:
		return FillToCapacity(v, capacity)
	case []SpellViewModel:
		return FillToCapacity(v, capacity)
	case []EquipmentViewModel:
		return FillToCapacity(v, capacity)
	case []MagicItemViewModel:
		return FillToCapacity(v, capacity)
	default:
		return items
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
