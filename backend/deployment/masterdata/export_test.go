package masterdata

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadExportFile(t *testing.T) {
	// Create temp file with test data
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_export.json")

	exportData := &ExportData{
		ExportVersion:  "1.0",
		BackendVersion: "0.4.0",
		Timestamp:      time.Now(),
		GameSystem:     "midgard",
		Data: map[string]interface{}{
			"skills": []interface{}{
				map[string]interface{}{
					"name":       "Schwimmen",
					"category":   "Körper",
					"difficulty": "leicht",
				},
			},
		},
	}

	err := WriteExportFile(testFile, exportData)
	require.NoError(t, err)

	// Read it back
	readData, err := ReadExportFile(testFile)
	assert.NoError(t, err)
	assert.NotNil(t, readData)
	assert.Equal(t, "1.0", readData.ExportVersion)
	assert.Equal(t, "0.4.0", readData.BackendVersion)
	assert.Equal(t, "midgard", readData.GameSystem)
	assert.NotNil(t, readData.Data)
}

func TestReadExportFile_NoVersion(t *testing.T) {
	// Create file without version (old format)
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "old_export.json")

	// Old format without export_version field
	oldFormat := `{
		"skills": [
			{"name": "Schwimmen", "category": "Körper"}
		]
	}`

	err := os.WriteFile(testFile, []byte(oldFormat), 0644)
	require.NoError(t, err)

	// Should default to version 1.0
	data, err := ReadExportFile(testFile)
	assert.NoError(t, err)
	assert.Equal(t, "1.0", data.ExportVersion)
}

func TestWriteExportFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "write_test.json")

	exportData := &ExportData{
		ExportVersion:  "1.0",
		BackendVersion: "0.4.0",
		Timestamp:      time.Now(),
		GameSystem:     "midgard",
		Data:           map[string]interface{}{},
	}

	err := WriteExportFile(testFile, exportData)
	assert.NoError(t, err)

	// Verify file exists and is valid JSON
	fileData, err := os.ReadFile(testFile)
	assert.NoError(t, err)
	assert.Contains(t, string(fileData), `"export_version": "1.0"`)
	assert.Contains(t, string(fileData), `"game_system": "midgard"`)
}

func TestTransformToCurrentVersion_AlreadyCurrent(t *testing.T) {
	data := &ExportData{
		ExportVersion: CurrentExportVersion,
		Data:          map[string]interface{}{},
	}

	transformed, err := TransformToCurrentVersion(data)
	assert.NoError(t, err)
	assert.Equal(t, CurrentExportVersion, transformed.ExportVersion)
}

func TestRegisterTransformer(t *testing.T) {
	// Save original registry
	original := transformerRegistry
	t.Cleanup(func() {
		transformerRegistry = original
	})

	// Clear registry for test
	transformerRegistry = []ImportTransformer{}

	// Create mock transformer
	mockTransformer := &mockTransformer{
		canTransform:  true,
		targetVersion: "2.0",
	}

	RegisterTransformer(mockTransformer)

	assert.Len(t, transformerRegistry, 1)
}

// Mock transformer for testing
type mockTransformer struct {
	canTransform  bool
	targetVersion string
}

func (m *mockTransformer) CanTransform(version string) bool {
	return m.canTransform
}

func (m *mockTransformer) Transform(data *ExportData) (*ExportData, error) {
	data.ExportVersion = m.targetVersion
	return data, nil
}

func (m *mockTransformer) TargetVersion() string {
	return m.targetVersion
}
