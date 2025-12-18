package pdfrender

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
