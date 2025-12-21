package pdfrender

// FillToCapacity fills a slice to a specified capacity with empty items
// This ensures tables render with the correct number of empty rows
func FillToCapacity[T any](items []T, capacity int) []T {
	if len(items) >= capacity {
		return items
	}

	// Create filled slice with capacity
	filled := make([]T, capacity)

	// Copy existing items
	copy(filled, items)

	// Remaining items are already zero-valued
	return filled
}
