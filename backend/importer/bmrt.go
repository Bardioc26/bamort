package importer

import (
	"encoding/json"
	"time"
)

// BMRTCharacter wraps CharacterImport with version and metadata tracking.
// This is the canonical interchange format for all adapter conversions.
type BMRTCharacter struct {
	CharacterImport                            // Embedded canonical format
	BmrtVersion     string                     `json:"bmrt_version"`         // e.g., "1.0"
	Extensions      map[string]json.RawMessage `json:"extensions,omitempty"` // Adapter-specific data
	Metadata        SourceMetadata             `json:"_metadata"`
}

// SourceMetadata tracks the origin of imported data for provenance
type SourceMetadata struct {
	SourceFormat string    `json:"source_format"` // e.g., "foundry-vtt"
	AdapterID    string    `json:"adapter_id"`    // e.g., "foundry-vtt-v1"
	ImportedAt   time.Time `json:"imported_at"`
}

// CurrentBMRTVersion is the supported BMRT format version
const CurrentBMRTVersion = "1.0"

// ValidateBMRTVersion checks if the BMRT version is supported
func ValidateBMRTVersion(version string) bool {
	return version == CurrentBMRTVersion
}

// NewBMRTCharacter creates a new BMRT character with metadata
func NewBMRTCharacter(char CharacterImport, adapterID, sourceFormat string) *BMRTCharacter {
	return &BMRTCharacter{
		CharacterImport: char,
		BmrtVersion:     CurrentBMRTVersion,
		Extensions:      make(map[string]json.RawMessage),
		Metadata: SourceMetadata{
			SourceFormat: sourceFormat,
			AdapterID:    adapterID,
			ImportedAt:   time.Now(),
		},
	}
}
