package pdfrender

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseTemplateMetadata extracts block metadata from HTML comments in template content
// Format: <!-- BLOCK: name, TYPE: type, MAX: number, FILTER: filter, NOEMPTY -->
func ParseTemplateMetadata(templateContent string) []BlockMetadata {
	blocks := []BlockMetadata{}

	// Simple regex to find all BLOCK comments
	re := regexp.MustCompile(`<!--\s*BLOCK:([^>]+)-->`)

	matches := re.FindAllStringSubmatch(templateContent, -1)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		// Parse the content of the BLOCK comment
		blockContent := match[1]

		// Extract BLOCK name
		nameRe := regexp.MustCompile(`BLOCK:\s*([^,]+)`)
		if !nameRe.MatchString(match[0]) {
			continue
		}
		nameMatch := nameRe.FindStringSubmatch(match[0])
		name := strings.TrimSpace(nameMatch[1])

		// Extract TYPE
		typeRe := regexp.MustCompile(`TYPE:\s*([^,]+)`)
		typeMatch := typeRe.FindStringSubmatch(blockContent)
		if typeMatch == nil {
			continue
		}
		listType := strings.TrimSpace(typeMatch[1])

		// Extract MAX
		maxRe := regexp.MustCompile(`MAX:\s*(\d+)`)
		maxMatch := maxRe.FindStringSubmatch(blockContent)
		if maxMatch == nil {
			continue
		}
		maxItems, _ := strconv.Atoi(strings.TrimSpace(maxMatch[1]))

		// Extract FILTER (optional)
		filter := ""
		filterRe := regexp.MustCompile(`FILTER:\s*([^,]+)`)
		filterMatch := filterRe.FindStringSubmatch(blockContent)
		if filterMatch != nil {
			filter = strings.TrimSpace(filterMatch[1])
		}

		// Check for NOEMPTY flag (optional)
		noEmpty := strings.Contains(blockContent, "NOEMPTY")

		blocks = append(blocks, BlockMetadata{
			Name:     name,
			ListType: listType,
			MaxItems: maxItems,
			Filter:   filter,
			NoEmpty:  noEmpty,
		})
	}

	return blocks
}

// GetBlockByName returns the block metadata for a specific block name
func GetBlockByName(blocks []BlockMetadata, name string) *BlockMetadata {
	for i := range blocks {
		if blocks[i].Name == name {
			return &blocks[i]
		}
	}
	return nil
}

// GetBlocksByType returns all blocks of a specific list type
func GetBlocksByType(blocks []BlockMetadata, listType string) []BlockMetadata {
	result := []BlockMetadata{}
	for _, block := range blocks {
		if block.ListType == listType {
			result = append(result, block)
		}
	}
	return result
}
