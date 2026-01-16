package masterdata

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// CurrentExportVersion is the current version of the export format
const CurrentExportVersion = "1.0"

// ExportData represents a versioned master data export
type ExportData struct {
	ExportVersion  string                 `json:"export_version"`
	BackendVersion string                 `json:"backend_version"`
	Timestamp      time.Time              `json:"timestamp"`
	GameSystem     string                 `json:"game_system"`
	Data           map[string]interface{} `json:"data"`
}

// ReadExportFile reads and parses an export file
func ReadExportFile(filePath string) (*ExportData, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var export ExportData
	if err := json.Unmarshal(data, &export); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// If no version specified, assume 1.0 (old format)
	if export.ExportVersion == "" {
		export.ExportVersion = "1.0"
	}

	return &export, nil
}

// WriteExportFile writes export data to a JSON file
func WriteExportFile(filePath string, export *ExportData) error {
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
