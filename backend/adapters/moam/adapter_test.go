package main

import (
	"bamort/importer"
	"encoding/json"
	"os"
	"testing"
)

// TestDetectMoamFormat verifies that Moam VTT format can be detected
func TestDetectMoamFormat(t *testing.T) {
	data, err := os.ReadFile("testdata/moam_character.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	confidence, version := detectMoamFormat(data)

	// Should have high confidence for valid Moam format
	if confidence < 0.9 {
		t.Errorf("Expected confidence >= 0.9, got %f", confidence)
	}

	// Should detect version (can be empty for generic detection)
	if version != "10.x" && version != "" {
		t.Logf("Detected version: %s", version)
	}
}

// TestDetectNonMoamFormat verifies that non-Moam format returns low confidence
func TestDetectNonMoamFormat(t *testing.T) {
	invalidData := []byte(`{"random": "data", "not": "moam"}`)

	confidence, _ := detectMoamFormat(invalidData)

	// Should have low confidence for invalid format
	if confidence > 0.5 {
		t.Errorf("Expected confidence < 0.5 for invalid data, got %f", confidence)
	}
}

// TestConvertMoamToBMRT verifies conversion from Moam to BMRT format
func TestConvertMoamToBMRT(t *testing.T) {
	data, err := os.ReadFile("testdata/moam_character.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var moamChar MoamCharacter
	if err := json.Unmarshal(data, &moamChar); err != nil {
		t.Fatalf("Failed to parse Moam JSON: %v", err)
	}

	bmrt, err := toBMRT(&moamChar)
	if err != nil {
		t.Fatalf("Conversion to BMRT failed: %v", err)
	}

	// Verify basic character data
	if bmrt.Name != "Test Character" {
		t.Errorf("Expected name 'Test Character', got '%s'", bmrt.Name)
	}

	if bmrt.Grad != 1 {
		t.Errorf("Expected grad 1, got %d", bmrt.Grad)
	}

	if bmrt.Rasse != "Mensch" {
		t.Errorf("Expected race 'Mensch', got '%s'", bmrt.Rasse)
	}

	// Verify eigenschaften (stats)
	if bmrt.Eigenschaften.St != 80 {
		t.Errorf("Expected St=80, got %d", bmrt.Eigenschaften.St)
	}

	if bmrt.Eigenschaften.Gw != 70 {
		t.Errorf("Expected Gw=70, got %d", bmrt.Eigenschaften.Gw)
	}

	// Verify LP/AP
	if bmrt.Lp.Max != 12 || bmrt.Lp.Value != 12 {
		t.Errorf("Expected LP max=12, value=12, got max=%d, value=%d", bmrt.Lp.Max, bmrt.Lp.Value)
	}

	if bmrt.Ap.Max != 20 || bmrt.Ap.Value != 20 {
		t.Errorf("Expected AP max=20, value=20, got max=%d, value=%d", bmrt.Ap.Max, bmrt.Ap.Value)
	}

	// Verify fertigkeiten
	if len(bmrt.Fertigkeiten) != 2 {
		t.Errorf("Expected 2 fertigkeiten, got %d", len(bmrt.Fertigkeiten))
	}

	// Verify waffenfertigkeiten
	if len(bmrt.Waffenfertigkeiten) != 1 {
		t.Errorf("Expected 1 waffenfertigkeit, got %d", len(bmrt.Waffenfertigkeiten))
	}

	// Verify equipment
	if len(bmrt.Waffen) != 1 {
		t.Errorf("Expected 1 weapon, got %d", len(bmrt.Waffen))
	}

	if len(bmrt.Ausruestung) != 1 {
		t.Errorf("Expected 1 armor, got %d", len(bmrt.Ausruestung))
	}

	if len(bmrt.Behaeltnisse) != 1 {
		t.Errorf("Expected 1 container, got %d", len(bmrt.Behaeltnisse))
	}
}

// TestConvertBMRTToMoam verifies round-trip conversion (BMRT back to Moam)
func TestConvertBMRTToMoam(t *testing.T) {
	data, err := os.ReadFile("testdata/moam_character.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var originalMoam MoamCharacter
	if err := json.Unmarshal(data, &originalMoam); err != nil {
		t.Fatalf("Failed to parse Moam JSON: %v", err)
	}

	// Convert to BMRT
	bmrt, err := toBMRT(&originalMoam)
	if err != nil {
		t.Fatalf("Conversion to BMRT failed: %v", err)
	}

	// Convert back to Moam
	convertedMoam, err := fromBMRT(bmrt)
	if err != nil {
		t.Fatalf("Conversion from BMRT failed: %v", err)
	}

	// Verify structural equality of key fields
	if convertedMoam.Name != originalMoam.Name {
		t.Errorf("Name mismatch: expected '%s', got '%s'", originalMoam.Name, convertedMoam.Name)
	}

	if convertedMoam.Grad != originalMoam.Grad {
		t.Errorf("Grad mismatch: expected %d, got %d", originalMoam.Grad, convertedMoam.Grad)
	}

	if convertedMoam.Rasse != originalMoam.Rasse {
		t.Errorf("Rasse mismatch: expected '%s', got '%s'", originalMoam.Rasse, convertedMoam.Rasse)
	}

	// Verify stats
	if convertedMoam.Eigenschaften.St != originalMoam.Eigenschaften.St {
		t.Errorf("St mismatch: expected %d, got %d", originalMoam.Eigenschaften.St, convertedMoam.Eigenschaften.St)
	}

	// Verify collection counts
	if len(convertedMoam.Fertigkeiten) != len(originalMoam.Fertigkeiten) {
		t.Errorf("Fertigkeiten count mismatch: expected %d, got %d",
			len(originalMoam.Fertigkeiten), len(convertedMoam.Fertigkeiten))
	}

	if len(convertedMoam.Waffenfertigkeiten) != len(originalMoam.Waffenfertigkeiten) {
		t.Errorf("Waffenfertigkeiten count mismatch: expected %d, got %d",
			len(originalMoam.Waffenfertigkeiten), len(convertedMoam.Waffenfertigkeiten))
	}
}

// TestInvalidJSON verifies that invalid JSON is handled gracefully
func TestInvalidJSON(t *testing.T) {
	invalidData := []byte(`{invalid json}`)

	var moamChar MoamCharacter
	err := json.Unmarshal(invalidData, &moamChar)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

// TestEmptyCharacterConversion verifies handling of minimal character data
func TestEmptyCharacterConversion(t *testing.T) {
	moamChar := MoamCharacter{
		CharacterImport: importer.CharacterImport{
			ID:   "test-empty",
			Name: "Empty Character",
		},
	}

	bmrt, err := toBMRT(&moamChar)
	if err != nil {
		t.Fatalf("Conversion should not fail for minimal data: %v", err)
	}

	if bmrt.Name != "Empty Character" {
		t.Errorf("Expected name 'Empty Character', got '%s'", bmrt.Name)
	}

	// Should have empty collections, not nil
	if bmrt.Fertigkeiten == nil {
		t.Error("Fertigkeiten should be empty slice, not nil")
	}
}

// TestMagischFieldConversion verifies magical item fields are preserved
func TestMagischFieldConversion(t *testing.T) {
	moamChar := MoamCharacter{
		CharacterImport: importer.CharacterImport{
			ID:   "test-magic",
			Name: "Magic Test",
			Waffen: []importer.Waffe{
				{
					ImportBase: importer.ImportBase{
						ID:   "magic-sword-1",
						Name: "Verzaubertes Schwert",
					},
					Schb:    3,
					Gewicht: 2.0,
					Magisch: importer.Magisch{
						IstMagisch:  true,
						Abw:         2,
						Ausgebrannt: false,
					},
				},
			},
		},
	}

	bmrt, err := toBMRT(&moamChar)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	if len(bmrt.Waffen) != 1 {
		t.Fatalf("Expected 1 weapon, got %d", len(bmrt.Waffen))
	}

	weapon := bmrt.Waffen[0]
	if !weapon.Magisch.IstMagisch {
		t.Error("Expected weapon to be magical")
	}

	if weapon.Magisch.Abw != 2 {
		t.Errorf("Expected Abw=2, got %d", weapon.Magisch.Abw)
	}
}

// TestContainerHierarchy verifies beinhaltet_in relationships are preserved
func TestContainerHierarchy(t *testing.T) {
	moamChar := MoamCharacter{
		CharacterImport: importer.CharacterImport{
			ID:   "test-containers",
			Name: "Container Test",
			Behaeltnisse: []importer.Behaeltniss{
				{
					ImportBase: importer.ImportBase{
						ID:   "container-1",
						Name: "Rucksack",
					},
					BeinhaltetIn: "",
					Tragkraft:    20.0,
				},
			},
			Waffen: []importer.Waffe{
				{
					ImportBase: importer.ImportBase{
						ID:   "weapon-1",
						Name: "Schwert",
					},
					BeinhaltetIn: "container-1",
				},
			},
		},
	}

	bmrt, err := toBMRT(&moamChar)
	if err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	// Verify container exists
	if len(bmrt.Behaeltnisse) != 1 {
		t.Fatalf("Expected 1 container, got %d", len(bmrt.Behaeltnisse))
	}

	// Verify weapon references container
	if len(bmrt.Waffen) != 1 {
		t.Fatalf("Expected 1 weapon, got %d", len(bmrt.Waffen))
	}

	if bmrt.Waffen[0].BeinhaltetIn != "container-1" {
		t.Errorf("Expected weapon in container-1, got '%s'", bmrt.Waffen[0].BeinhaltetIn)
	}
}
