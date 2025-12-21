package pdfrender

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseTemplateMetadata extracts block metadata from HTML comments in template content
// Format: <!-- BLOCK: name, TYPE: type, MAX: number, FILTER: filter -->
func ParseTemplateMetadata(templateContent string) []BlockMetadata {
	blocks := []BlockMetadata{}

	// Regex to match: <!-- BLOCK: name, TYPE: type, MAX: number, FILTER: filter -->
	re := regexp.MustCompile(`<!--\s*BLOCK:\s*([^,]+),\s*TYPE:\s*([^,]+),\s*MAX:\s*(\d+)(?:,\s*FILTER:\s*([^-]+))?\s*-->`)

	matches := re.FindAllStringSubmatch(templateContent, -1)
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		name := strings.TrimSpace(match[1])
		listType := strings.TrimSpace(match[2])
		maxItems, _ := strconv.Atoi(strings.TrimSpace(match[3]))

		filter := ""
		if len(match) > 4 && match[4] != "" {
			filter = strings.TrimSpace(match[4])
		}

		blocks = append(blocks, BlockMetadata{
			Name:     name,
			ListType: listType,
			MaxItems: maxItems,
			Filter:   filter,
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
