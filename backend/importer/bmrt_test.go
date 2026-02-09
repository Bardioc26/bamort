package importer

import (
	"encoding/json"
	"testing"
	"time"
)

func TestValidateBMRTVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    bool
	}{
		{
			name:    "valid current version",
			version: "1.0",
			want:    true,
		},
		{
			name:    "invalid future version",
			version: "2.0",
			want:    false,
		},
		{
			name:    "invalid old version",
			version: "0.9",
			want:    false,
		},
		{
			name:    "empty version",
			version: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateBMRTVersion(tt.version); got != tt.want {
				t.Errorf("ValidateBMRTVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBMRTCharacter(t *testing.T) {
	char := CharacterImport{
		Name: "Test Character",
		Typ:  "Krieger",
	}
	adapterID := "foundry-vtt-v1"
	sourceFormat := "foundry-vtt"

	before := time.Now()
	bmrt := NewBMRTCharacter(char, adapterID, sourceFormat)
	after := time.Now()

	if bmrt.Name != "Test Character" {
		t.Errorf("Character name = %v, want %v", bmrt.Name, "Test Character")
	}
	if bmrt.Typ != "Krieger" {
		t.Errorf("Typ = %v, want %v", bmrt.Typ, "Krieger")
	}
	if bmrt.BmrtVersion != CurrentBMRTVersion {
		t.Errorf("BmrtVersion = %v, want %v", bmrt.BmrtVersion, CurrentBMRTVersion)
	}
	if bmrt.Metadata.AdapterID != adapterID {
		t.Errorf("Metadata.AdapterID = %v, want %v", bmrt.Metadata.AdapterID, adapterID)
	}
	if bmrt.Metadata.SourceFormat != sourceFormat {
		t.Errorf("Metadata.SourceFormat = %v, want %v", bmrt.Metadata.SourceFormat, sourceFormat)
	}
	if bmrt.Metadata.ImportedAt.Before(before) || bmrt.Metadata.ImportedAt.After(after) {
		t.Errorf("Metadata.ImportedAt not within expected time range")
	}
	if bmrt.Extensions == nil {
		t.Error("Extensions map should be initialized")
	}
}

func TestBMRTCharacterJSONSerialization(t *testing.T) {
	char := CharacterImport{
		Name: "Test Character",
		Typ:  "Krieger",
	}
	bmrt := NewBMRTCharacter(char, "test-adapter", "test-format")

	// Add extension data
	extensionData := map[string]interface{}{
		"original_id": "abc123",
		"version":     "11.x",
	}
	extensionJSON, _ := json.Marshal(extensionData)
	bmrt.Extensions["foundry"] = extensionJSON

	// Serialize
	data, err := json.Marshal(bmrt)
	if err != nil {
		t.Fatalf("Failed to marshal BMRTCharacter: %v", err)
	}

	// Deserialize
	var decoded BMRTCharacter
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal BMRTCharacter: %v", err)
	}

	// Verify
	if decoded.Name != bmrt.Name {
		t.Errorf("Decoded name = %v, want %v", decoded.Name, bmrt.Name)
	}
	if decoded.BmrtVersion != bmrt.BmrtVersion {
		t.Errorf("Decoded BmrtVersion = %v, want %v", decoded.BmrtVersion, bmrt.BmrtVersion)
	}
	if decoded.Metadata.AdapterID != bmrt.Metadata.AdapterID {
		t.Errorf("Decoded AdapterID = %v, want %v", decoded.Metadata.AdapterID, bmrt.Metadata.AdapterID)
	}
	if len(decoded.Extensions) != 1 {
		t.Errorf("Decoded Extensions length = %v, want 1", len(decoded.Extensions))
	}
}
