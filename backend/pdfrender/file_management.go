package pdfrender

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// SanitizeFilename removes or replaces characters that are problematic in filenames
func SanitizeFilename(name string) string {
	// Replace umlauts
	replacements := map[string]string{
		"ä": "ae", "Ä": "Ae",
		"ö": "oe", "Ö": "Oe",
		"ü": "ue", "Ü": "Ue",
		"ß": "ss",
	}

	result := name
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	// Replace problematic characters with underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9_.-]+`)
	result = reg.ReplaceAllString(result, "_")

	// Remove consecutive underscores
	reg = regexp.MustCompile(`_+`)
	result = reg.ReplaceAllString(result, "_")

	// Trim underscores from start and end
	result = strings.Trim(result, "_")

	return result
}

// GenerateExportFilename creates a filename from character name and timestamp
func GenerateExportFilename(charName string, timestamp time.Time) string {
	sanitizedName := SanitizeFilename(charName)
	timeStr := timestamp.Format("20060102_150405")
	return sanitizedName + "_" + timeStr + ".pdf"
}

// EnsureExportTempDir creates the export temp directory if it doesn't exist
func EnsureExportTempDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// CleanupOldFiles removes files older than maxAge from the directory
// Returns the number of files deleted
func CleanupOldFiles(dir string, maxAge time.Duration) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	count := 0
	cutoff := time.Now().Add(-maxAge)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filePath); err == nil {
				count++
			}
		}
	}

	return count, nil
}
