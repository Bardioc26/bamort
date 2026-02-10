package main

import (
	"bamort/importer"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	if version != "5.x" && version != "" {
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

// TestConvertMoamToBMRT verifies conversion from Moam Import data to BMRT format
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

// testing a real char
func TestCharRealEinskaldir(t *testing.T) {
	data, err := os.ReadFile("testdata/einskaldir.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var moamChar MoamCharacter
	if err := json.Unmarshal(data, &moamChar); err != nil {
		t.Fatalf("Failed to parse Moam JSON: %v", err)
	}

	bmrt, err := toBMRT(&moamChar)
	assert.NoErrorf(t, err, "Conversion to BMRT failed")

	assert.Equalf(t, "Einskaldir", bmrt.Name, "Expected name 'Einskaldir', got '%s'", bmrt.Name)

}

// TestGalaxisKessariusRufusComprehensive tests the complete conversion of the Galaxis character
// including all slices, umlauts (UTF-8), and container relationships
func TestGalaxisKessariusRufusComprehensive(t *testing.T) {
	data, err := os.ReadFile("testdata/galaxis-kessarius-rufus.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var moamChar MoamCharacter
	if err := json.Unmarshal(data, &moamChar); err != nil {
		t.Fatalf("Failed to parse Moam JSON: %v", err)
	}

	bmrt, err := toBMRT(&moamChar)
	assert.NoError(t, err, "Conversion to BMRT failed")

	// Test basic character data with umlauts
	assert.Equal(t, "Galaxis Kessarius Rufus", bmrt.Name, "Character name mismatch")
	assert.Equal(t, 13, bmrt.Grad, "Character grade mismatch")
	assert.Equal(t, "Mensch", bmrt.Rasse, "Character race mismatch")
	assert.Equal(t, "Thaumaturg", bmrt.Typ, "Character type mismatch")
	assert.Equal(t, 19, bmrt.Groesse, "Character size mismatch")
	assert.Equal(t, 10, bmrt.Gewicht, "Character weight mismatch")
	assert.Equal(t, 18, bmrt.Alter, "Character age mismatch")
	assert.Equal(t, "Mittelschicht", bmrt.SocialClass, "Social class mismatch")
	assert.Equal(t, "Jakchos", bmrt.Glaube, "Religion mismatch")

	// Test eigenschaften
	assert.Equal(t, 96, bmrt.Eigenschaften.St, "Stärke (St) mismatch")
	assert.Equal(t, 65, bmrt.Eigenschaften.Gs, "Geschicklichkeit (Gs) mismatch")
	assert.Equal(t, 62, bmrt.Eigenschaften.Gw, "Gewandtheit (Gw) mismatch")
	assert.Equal(t, 87, bmrt.Eigenschaften.Ko, "Konstitution (Ko) mismatch")
	assert.Equal(t, 76, bmrt.Eigenschaften.In, "Intelligenz (In) mismatch")
	assert.Equal(t, 100, bmrt.Eigenschaften.Zt, "Zaubertaltent (Zt) mismatch")

	// Test LP/AP/B
	assert.Equal(t, 19, bmrt.Lp.Max, "LP Max mismatch")
	assert.Equal(t, 19, bmrt.Lp.Value, "LP Value mismatch")
	assert.Equal(t, 61, bmrt.Ap.Max, "AP Max mismatch")
	assert.Equal(t, 61, bmrt.Ap.Value, "AP Value mismatch")
	assert.Equal(t, 28, bmrt.B.Max, "B Max mismatch")

	// Test Fertigkeiten (sample with umlauts)
	assert.Equal(t, 29, len(bmrt.Fertigkeiten), "Fertigkeiten count mismatch")

	// Sample fertigkeiten with umlauts
	var gelaendelauf *importer.Fertigkeit
	var geschaeftssinn *importer.Fertigkeit
	var ueberleben *importer.Fertigkeit
	var ueberlebenGebirge *importer.Fertigkeit

	for i := range bmrt.Fertigkeiten {
		f := &bmrt.Fertigkeiten[i]
		switch f.Name {
		case "Geländelauf":
			gelaendelauf = f
		case "Geschäftssinn":
			geschaeftssinn = f
		case "Überleben":
			ueberleben = f
		case "Überleben:Gebirge":
			ueberlebenGebirge = f
		}
	}

	assert.NotNil(t, gelaendelauf, "Geländelauf not found")
	assert.Equal(t, 12, gelaendelauf.Fertigkeitswert, "Geländelauf value mismatch")
	assert.Equal(t, "KOD5 112", gelaendelauf.Quelle, "Geländelauf source mismatch")

	assert.NotNil(t, geschaeftssinn, "Geschäftssinn not found")
	assert.Equal(t, 10, geschaeftssinn.Fertigkeitswert, "Geschäftssinn value mismatch")
	assert.Equal(t, "KOD5 113", geschaeftssinn.Quelle, "Geschäftssinn source mismatch")

	assert.NotNil(t, ueberleben, "Überleben not found")
	assert.Equal(t, 8, ueberleben.Fertigkeitswert, "Überleben value mismatch")
	assert.Equal(t, "Wald", ueberleben.Beschreibung, "Überleben description mismatch")

	assert.NotNil(t, ueberlebenGebirge, "Überleben:Gebirge not found")
	assert.Equal(t, 5, ueberlebenGebirge.Fertigkeitswert, "Überleben:Gebirge value mismatch")

	// Test Waffenfertigkeiten (all with umlauts)
	assert.Equal(t, 3, len(bmrt.Waffenfertigkeiten), "Waffenfertigkeiten count mismatch")

	var armbrueste *importer.Waffenfertigkeit
	var blasrohre *importer.Waffenfertigkeit

	for i := range bmrt.Waffenfertigkeiten {
		wf := &bmrt.Waffenfertigkeiten[i]
		switch wf.Name {
		case "Armbrüste":
			armbrueste = wf
		case "Blasrohre":
			blasrohre = wf
		}
	}

	assert.NotNil(t, armbrueste, "Armbrüste not found")
	assert.Equal(t, 7, armbrueste.Fertigkeitswert, "Armbrüste value mismatch")

	assert.NotNil(t, blasrohre, "Blasrohre not found")
	assert.Equal(t, 9, blasrohre.Fertigkeitswert, "Blasrohre value mismatch")

	// Test Zauber (sample with umlauts and special characters)
	assert.Equal(t, 34, len(bmrt.Zauber), "Zauber count mismatch")

	var staerke *importer.Zauber
	var vergroessern *importer.Zauber
	var zuendersalz *importer.Zauber
	var loeschsalz *importer.Zauber
	var feuerFinger *importer.Zauber
	var handauflegen *importer.Zauber

	for i := range bmrt.Zauber {
		z := &bmrt.Zauber[i]
		switch z.Name {
		case "Stärke (S)":
			staerke = z
		case "Vergrößern (S)":
			vergroessern = z
		case "Zündersalz":
			zuendersalz = z
		case "Löschsalz":
			loeschsalz = z
		case "Feuerfinger (R)":
			feuerFinger = z
		case "Handauflegen":
			handauflegen = z
		}
	}

	assert.NotNil(t, staerke, "Stärke (S) not found")
	assert.Equal(t, "MYS5 61", staerke.Quelle, "Stärke source mismatch")

	assert.NotNil(t, vergroessern, "Vergrößern (S) not found")
	assert.Equal(t, "MYS5 61", vergroessern.Quelle, "Vergrößern source mismatch")

	assert.NotNil(t, zuendersalz, "Zündersalz not found")
	assert.Equal(t, "MYS5 50", zuendersalz.Quelle, "Zündersalz source mismatch")

	assert.NotNil(t, loeschsalz, "Löschsalz not found")
	assert.Equal(t, "MYS5 48", loeschsalz.Quelle, "Löschsalz source mismatch")

	assert.NotNil(t, feuerFinger, "Feuerfinger (R) not found")
	assert.Equal(t, "MYS5 53", feuerFinger.Quelle, "Feuerfinger source mismatch")

	assert.NotNil(t, handauflegen, "Handauflegen not found")
	assert.Equal(t, "ARK5 143", handauflegen.Quelle, "Handauflegen source mismatch")

	// Test Waffen
	assert.Equal(t, 4, len(bmrt.Waffen), "Waffen count mismatch")

	var armbrust *importer.Waffe
	for i := range bmrt.Waffen {
		if bmrt.Waffen[i].Name == "Armbrust:leicht" {
			armbrust = &bmrt.Waffen[i]
			break
		}
	}
	assert.NotNil(t, armbrust, "Armbrust:leicht not found")
	assert.Equal(t, 2.5, armbrust.Gewicht, "Armbrust weight mismatch")
	assert.Equal(t, 40.0, armbrust.Wert, "Armbrust value mismatch")
	assert.Equal(t, 1, armbrust.Anzahl, "Armbrust count mismatch")

	// Test Ausruestung (with umlauts and container references)
	assert.Equal(t, 11, len(bmrt.Ausruestung), "Ausruestung count mismatch")

	var lederruestung *importer.Ausruestung
	var koecher *importer.Ausruestung
	var seil *importer.Ausruestung

	for i := range bmrt.Ausruestung {
		a := &bmrt.Ausruestung[i]
		switch a.Name {
		case "Lederrüstung":
			lederruestung = a
		case "Köcher, Bolzen":
			koecher = a
		case "Seil":
			seil = a
		}
	}

	assert.NotNil(t, lederruestung, "Lederrüstung not found")
	assert.Equal(t, 13.0, lederruestung.Gewicht, "Lederrüstung weight mismatch")
	assert.Equal(t, 30.0, lederruestung.Wert, "Lederrüstung value mismatch")
	assert.Equal(t, "", lederruestung.BeinhaltetIn, "Lederrüstung should not be in container")

	assert.NotNil(t, koecher, "Köcher not found")
	assert.Equal(t, 0.3, koecher.Gewicht, "Köcher weight mismatch")

	assert.NotNil(t, seil, "Seil not found")
	assert.Equal(t, 0.075, seil.Gewicht, "Seil weight mismatch per unit")
	assert.Equal(t, 20, seil.Anzahl, "Seil count mismatch")
	assert.Equal(t, "moam-container-47422", seil.BeinhaltetIn, "Seil should be in Lederrucksack")

	// Test Behaeltnisse (containers) and ContainedIn relationships
	assert.Equal(t, 3, len(bmrt.Behaeltnisse), "Behaeltnisse count mismatch")

	var lederrucksack *importer.Behaeltniss
	var guerteltasche *importer.Behaeltniss
	var lederhuelle *importer.Behaeltniss

	for i := range bmrt.Behaeltnisse {
		b := &bmrt.Behaeltnisse[i]
		switch b.Name {
		case "Lederrucksack":
			lederrucksack = b
		case "Gürteltasche, Leder":
			guerteltasche = b
		case "Lederhülle":
			lederhuelle = b
		}
	}

	assert.NotNil(t, lederrucksack, "Lederrucksack not found")
	assert.Equal(t, "moam-container-47422", lederrucksack.ID, "Lederrucksack ID mismatch")
	assert.Equal(t, 0.5, lederrucksack.Gewicht, "Lederrucksack weight mismatch")
	assert.Equal(t, 25.0, lederrucksack.Tragkraft, "Lederrucksack capacity mismatch")
	assert.Equal(t, "", lederrucksack.BeinhaltetIn, "Lederrucksack should not be in another container")

	assert.NotNil(t, guerteltasche, "Gürteltasche not found")
	assert.Equal(t, "für 20 Münzen", guerteltasche.Beschreibung, "Gürteltasche description with umlaut mismatch")
	assert.Equal(t, 1.0, guerteltasche.Tragkraft, "Gürteltasche capacity mismatch")

	assert.NotNil(t, lederhuelle, "Lederhülle not found")
	assert.Equal(t, "wasserdicht für Schriftrollen", lederhuelle.Beschreibung, "Lederhülle description with umlaut mismatch")

	// Verify items contained in Lederrucksack (moam-container-47422)
	itemsInBackpack := 0
	for i := range bmrt.Ausruestung {
		if bmrt.Ausruestung[i].BeinhaltetIn == "moam-container-47422" {
			itemsInBackpack++
		}
	}
	assert.Equal(t, 6, itemsInBackpack, "Expected 6 items in Lederrucksack")

	// Test Transportmittel
	assert.Equal(t, 0, len(bmrt.Transportmittel), "Transportmittel should be empty")
	assert.NotNil(t, bmrt.Transportmittel, "Transportmittel should not be nil")

	// Test gestalt with umlauts
	assert.Equal(t, "klein", bmrt.Gestalt.Groesse, "Gestalt Größe mismatch")
	assert.Equal(t, "breit", bmrt.Gestalt.Breite, "Gestalt Breite mismatch")

	// Test merkmale (should preserve UTF-8 empty strings)
	assert.NotNil(t, bmrt.Merkmale, "Merkmale should not be nil")

	// Test erfahrungsschatz
	assert.Equal(t, 3960, bmrt.Erfahrungsschatz.Value, "Erfahrungsschatz mismatch")
}

// TestDorinSchnellhammerComprehensive tests the complete conversion of the Dorin Schnellhammer character
// including all slices, umlauts (UTF-8), specializations, and container relationships
func TestDorinSchnellhammerComprehensive(t *testing.T) {
	data, err := os.ReadFile("testdata/dorin-schnellhammer.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var moamChar MoamCharacter
	if err := json.Unmarshal(data, &moamChar); err != nil {
		t.Fatalf("Failed to parse Moam JSON: %v", err)
	}

	bmrt, err := toBMRT(&moamChar)
	assert.NoError(t, err, "Conversion to BMRT failed")

	// Test basic character data with umlauts
	assert.Equal(t, "Dorin Schnellhammer", bmrt.Name, "Character name mismatch")
	assert.Equal(t, 16, bmrt.Grad, "Character grade mismatch")
	assert.Equal(t, "Zwerg", bmrt.Rasse, "Character race mismatch")
	assert.Equal(t, "Krieger", bmrt.Typ, "Character type mismatch")
	assert.Equal(t, 145, bmrt.Groesse, "Character size mismatch")
	assert.Equal(t, 78, bmrt.Gewicht, "Character weight mismatch")
	assert.Equal(t, 74, bmrt.Alter, "Character age mismatch")
	assert.Equal(t, "Volk", bmrt.SocialClass, "Social class mismatch")
	assert.Equal(t, "Mahal", bmrt.Glaube, "Religion mismatch")
	assert.Equal(t, "beide", bmrt.Hand, "Handedness mismatch")

	// Test Spezialisierung
	assert.Equal(t, 3, len(bmrt.Spezialisierung), "Spezialisierung count mismatch")
	assert.Contains(t, bmrt.Spezialisierung, "Kriegshammer", "Missing Kriegshammer specialization")
	assert.Contains(t, bmrt.Spezialisierung, "Schlachtbeil", "Missing Schlachtbeil specialization")
	assert.Contains(t, bmrt.Spezialisierung, "Streitaxt", "Missing Streitaxt specialization")

	// Test eigenschaften
	assert.Equal(t, 99, bmrt.Eigenschaften.St, "Stärke (St) mismatch")
	assert.Equal(t, 75, bmrt.Eigenschaften.Gs, "Geschicklichkeit (Gs) mismatch")
	assert.Equal(t, 70, bmrt.Eigenschaften.Gw, "Gewandtheit (Gw) mismatch")
	assert.Equal(t, 92, bmrt.Eigenschaften.Ko, "Konstitution (Ko) mismatch")

	// Test LP/AP/B
	assert.Equal(t, 20, bmrt.Lp.Max, "LP Max mismatch")
	assert.Equal(t, 20, bmrt.Lp.Value, "LP Value mismatch")
	assert.Equal(t, 56, bmrt.Ap.Max, "AP Max mismatch")
	assert.Equal(t, 56, bmrt.Ap.Value, "AP Value mismatch")
	assert.Equal(t, 24, bmrt.B.Max, "B Max mismatch")

	// Test Fertigkeiten (sample with umlauts)
	var anfuehren *importer.Fertigkeit
	var beidhaendiger *importer.Fertigkeit
	var gelaendelauf *importer.Fertigkeit
	var kampfVollruestung *importer.Fertigkeit

	for i := range bmrt.Fertigkeiten {
		f := &bmrt.Fertigkeiten[i]
		switch f.Name {
		case "Anführen":
			anfuehren = f
		case "Beidhändiger Kampf":
			beidhaendiger = f
		case "Geländelauf":
			gelaendelauf = f
		case "Kampf in Vollrüstung":
			kampfVollruestung = f
		}
	}

	assert.NotNil(t, anfuehren, "Anführen not found")
	assert.Equal(t, 10, anfuehren.Fertigkeitswert, "Anführen value mismatch")

	assert.NotNil(t, beidhaendiger, "Beidhändiger Kampf not found")
	assert.Equal(t, 15, beidhaendiger.Fertigkeitswert, "Beidhändiger Kampf value mismatch")

	assert.NotNil(t, gelaendelauf, "Geländelauf not found")
	assert.Equal(t, 14, gelaendelauf.Fertigkeitswert, "Geländelauf value mismatch")

	assert.NotNil(t, kampfVollruestung, "Kampf in Vollrüstung not found")
	assert.Equal(t, 15, kampfVollruestung.Fertigkeitswert, "Kampf in Vollrüstung value mismatch")

	// Test Waffenfertigkeiten (with umlauts)
	assert.Equal(t, 7, len(bmrt.Waffenfertigkeiten), "Waffenfertigkeiten count mismatch")

	var armbrueste *importer.Waffenfertigkeit
	var boegen *importer.Waffenfertigkeit

	for i := range bmrt.Waffenfertigkeiten {
		wf := &bmrt.Waffenfertigkeiten[i]
		switch wf.Name {
		case "Armbrüste":
			armbrueste = wf
		case "Bögen":
			boegen = wf
		}
	}

	assert.NotNil(t, armbrueste, "Armbrüste not found")
	assert.Equal(t, 11, armbrueste.Fertigkeitswert, "Armbrüste value mismatch")

	assert.NotNil(t, boegen, "Bögen not found")
	assert.Equal(t, 10, boegen.Fertigkeitswert, "Bögen value mismatch")

	// Test Zauber (should be empty)
	assert.Equal(t, 0, len(bmrt.Zauber), "Zauber should be empty for warrior")
	assert.NotNil(t, bmrt.Zauber, "Zauber should not be nil")

	// Test Waffen
	assert.Greater(t, len(bmrt.Waffen), 0, "Should have weapons")

	var armbrust *importer.Waffe
	for i := range bmrt.Waffen {
		if bmrt.Waffen[i].Name == "Armbrust:leicht" {
			armbrust = &bmrt.Waffen[i]
			break
		}
	}
	assert.NotNil(t, armbrust, "Armbrust:leicht not found")

	// Test Ausruestung (with umlauts)
	assert.Greater(t, len(bmrt.Ausruestung), 0, "Should have equipment")

	var kettenruestung *importer.Ausruestung
	for i := range bmrt.Ausruestung {
		if bmrt.Ausruestung[i].Name == "Kettenrüstung" {
			kettenruestung = &bmrt.Ausruestung[i]
			break
		}
	}

	assert.NotNil(t, kettenruestung, "Kettenrüstung not found")
	assert.Equal(t, 20.0, kettenruestung.Gewicht, "Kettenrüstung weight mismatch")
	assert.Equal(t, 100.0, kettenruestung.Wert, "Kettenrüstung value mismatch")

	// Test Behaeltnisse (containers) with umlauts
	assert.Equal(t, 4, len(bmrt.Behaeltnisse), "Behaeltnisse count mismatch")

	var lederrucksack *importer.Behaeltniss
	var beutel *importer.Behaeltniss

	for i := range bmrt.Behaeltnisse {
		b := &bmrt.Behaeltnisse[i]
		switch b.Name {
		case "Lederrucksack":
			lederrucksack = b
		case "Beutel, Leder":
			beutel = b
		}
	}

	assert.NotNil(t, lederrucksack, "Lederrucksack not found")
	assert.Equal(t, "moam-container-59617", lederrucksack.ID, "Lederrucksack ID mismatch")
	assert.Equal(t, "für 25 kg", lederrucksack.Beschreibung, "Lederrucksack description with umlaut mismatch")
	assert.Equal(t, 25.0, lederrucksack.Tragkraft, "Lederrucksack capacity mismatch")
	assert.Equal(t, "", lederrucksack.BeinhaltetIn, "Lederrucksack should not be in another container")

	assert.NotNil(t, beutel, "Beutel, Leder not found")
	assert.Equal(t, "für 20 Münzen", beutel.Beschreibung, "Beutel description with umlaut mismatch")

	// Verify items contained in Lederrucksack (moam-container-59617)
	itemsInBackpack := 0
	for i := range bmrt.Ausruestung {
		if bmrt.Ausruestung[i].BeinhaltetIn == "moam-container-59617" {
			itemsInBackpack++
		}
	}
	assert.Equal(t, 5, itemsInBackpack, "Expected 5 items in Lederrucksack")

	// Test Transportmittel
	assert.Equal(t, 1, len(bmrt.Transportmittel), "Should have 1 transportmittel")
	assert.Equal(t, "Reitpferd (Transportmittel)", bmrt.Transportmittel[0].Name, "Transportmittel name mismatch")

	// Test gestalt
	assert.Equal(t, "klein", bmrt.Gestalt.Groesse, "Gestalt Größe mismatch")
	assert.Equal(t, "normal", bmrt.Gestalt.Breite, "Gestalt Breite mismatch")

	// Test erfahrungsschatz
	assert.Equal(t, 5586, bmrt.Erfahrungsschatz.Value, "Erfahrungsschatz mismatch")

	// Test Image transfer - Dorin has a base64 image
	assert.NotEmpty(t, bmrt.Image, "Character image should be transferred")
	assert.Contains(t, bmrt.Image, "data:image;base64,", "Image should be base64 encoded")
	assert.Greater(t, len(bmrt.Image), 100, "Image data should be substantial")
}

// TestNicoloSikeriComprehensive tests the complete conversion of the Nicolo Sikeri character
// including all slices, umlauts (UTF-8), priest spells, and container relationships
func TestNicoloSikeriComprehensive(t *testing.T) {
	data, err := os.ReadFile("testdata/nicolo-sikeri.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var moamChar MoamCharacter
	if err := json.Unmarshal(data, &moamChar); err != nil {
		t.Fatalf("Failed to parse Moam JSON: %v", err)
	}

	bmrt, err := toBMRT(&moamChar)
	assert.NoError(t, err, "Conversion to BMRT failed")

	// Test basic character data
	assert.Equal(t, "Nicolo Sikeri", bmrt.Name, "Character name mismatch")
	assert.Equal(t, 7, bmrt.Grad, "Character grade mismatch")
	assert.Equal(t, "Mensch", bmrt.Rasse, "Character race mismatch")
	assert.Equal(t, "Priester, Beschützer", bmrt.Typ, "Character type with umlaut mismatch")
	assert.Equal(t, 172, bmrt.Groesse, "Character size mismatch")
	assert.Equal(t, 80, bmrt.Gewicht, "Character weight mismatch")
	assert.Equal(t, 23, bmrt.Alter, "Character age mismatch")
	assert.Equal(t, "Volk", bmrt.SocialClass, "Social class mismatch")
	assert.Equal(t, "Notuns", bmrt.Glaube, "Religion mismatch")

	// Test Spezialisierung (should be empty)
	assert.Equal(t, 0, len(bmrt.Spezialisierung), "Spezialisierung should be empty")
	assert.NotNil(t, bmrt.Spezialisierung, "Spezialisierung should not be nil")

	// Test eigenschaften
	assert.Equal(t, 95, bmrt.Eigenschaften.St, "Stärke (St) mismatch")
	assert.Equal(t, 99, bmrt.Eigenschaften.Gs, "Geschicklichkeit (Gs) mismatch")
	assert.Equal(t, 89, bmrt.Eigenschaften.Gw, "Gewandtheit (Gw) mismatch")
	assert.Equal(t, 100, bmrt.Eigenschaften.Zt, "Zaubertalent (Zt) mismatch")

	// Test LP/AP/B
	assert.Equal(t, 15, bmrt.Lp.Max, "LP Max mismatch")
	assert.Equal(t, 15, bmrt.Lp.Value, "LP Value mismatch")
	assert.Equal(t, 34, bmrt.Ap.Max, "AP Max mismatch")
	assert.Equal(t, 29, bmrt.Ap.Value, "AP Value mismatch")
	assert.Equal(t, 23, bmrt.B.Max, "B Max mismatch")

	// Test Merkmale with umlauts
	assert.NotNil(t, bmrt.Merkmale, "Merkmale should not be nil")
	assert.Equal(t, "Schwarz", bmrt.Merkmale.Haarfarbe, "Haarfarbe mismatch")
	assert.Equal(t, "grün", bmrt.Merkmale.Augenfarbe, "Augenfarbe with umlaut mismatch")

	// Test Waffenfertigkeiten
	assert.Equal(t, 2, len(bmrt.Waffenfertigkeiten), "Waffenfertigkeiten count mismatch")

	var einhandschlagwaffen *importer.Waffenfertigkeit
	for i := range bmrt.Waffenfertigkeiten {
		if bmrt.Waffenfertigkeiten[i].Name == "Einhandschlagwaffen" {
			einhandschlagwaffen = &bmrt.Waffenfertigkeiten[i]
			break
		}
	}
	assert.NotNil(t, einhandschlagwaffen, "Einhandschlagwaffen not found")
	assert.Equal(t, 9, einhandschlagwaffen.Fertigkeitswert, "Einhandschlagwaffen value mismatch")

	// Test Zauber (priest spells with umlauts)
	assert.Equal(t, 18, len(bmrt.Zauber), "Zauber count mismatch")

	// Check for spell with umlaut 'ö'
	foundBoesen := false
	for i := range bmrt.Zauber {
		if bmrt.Zauber[i].Name == "Austreibung des Bösen" {
			foundBoesen = true
			assert.Equal(t, "ARK5 136", bmrt.Zauber[i].Quelle, "Austreibung des Bösen source mismatch")
			break
		}
	}
	assert.True(t, foundBoesen, "Austreibung des Bösen not found")

	// Check for spell with umlaut 'ä'
	foundStaerke := false
	for i := range bmrt.Zauber {
		if bmrt.Zauber[i].Name == "Stärke" {
			foundStaerke = true
			assert.Equal(t, "ARK5 115", bmrt.Zauber[i].Quelle, "Stärke source mismatch")
			break
		}
	}
	assert.True(t, foundStaerke, "Stärke not found")

	// Test Waffen
	assert.Equal(t, 3, len(bmrt.Waffen), "Waffen count mismatch")

	var kriegshammer *importer.Waffe
	var streitaxt *importer.Waffe

	for i := range bmrt.Waffen {
		w := &bmrt.Waffen[i]
		switch w.Name {
		case "Kriegshammer":
			kriegshammer = w
		case "Streitaxt":
			streitaxt = w
		}
	}

	assert.NotNil(t, kriegshammer, "Kriegshammer not found")
	assert.Equal(t, 2.5, kriegshammer.Gewicht, "Kriegshammer weight mismatch")

	assert.NotNil(t, streitaxt, "Streitaxt not found")
	assert.Equal(t, 2.5, streitaxt.Gewicht, "Streitaxt weight mismatch")

	// Test Behaeltnisse (containers) with umlauts
	assert.Equal(t, 8, len(bmrt.Behaeltnisse), "Behaeltnisse count mismatch")

	var lederrucksack *importer.Behaeltniss
	var guerteltasche *importer.Behaeltniss
	var lederhuelle *importer.Behaeltniss
	var tuschegefaess *importer.Behaeltniss

	for i := range bmrt.Behaeltnisse {
		b := &bmrt.Behaeltnisse[i]
		switch b.Name {
		case "Lederrucksack":
			lederrucksack = b
		case "Gürteltasche, Leder":
			guerteltasche = b
		case "Lederhülle":
			lederhuelle = b
		case "Tuschegefäß":
			tuschegefaess = b
		}
	}

	assert.NotNil(t, lederrucksack, "Lederrucksack not found")
	assert.Equal(t, "moam-container-53059", lederrucksack.ID, "Lederrucksack ID mismatch")
	assert.Equal(t, "für 25 kg", lederrucksack.Beschreibung, "Lederrucksack description with umlaut mismatch")
	assert.Equal(t, 25.0, lederrucksack.Tragkraft, "Lederrucksack capacity mismatch")
	assert.Equal(t, "", lederrucksack.BeinhaltetIn, "Lederrucksack should not be in another container")

	assert.NotNil(t, guerteltasche, "Gürteltasche not found")
	assert.Equal(t, "für 20 Münzen", guerteltasche.Beschreibung, "Gürteltasche description with umlaut mismatch")

	assert.NotNil(t, lederhuelle, "Lederhülle not found")
	assert.Equal(t, "wasserdicht für Schriftrollen", lederhuelle.Beschreibung, "Lederhülle description with umlaut mismatch")

	assert.NotNil(t, tuschegefaess, "Tuschegefäß not found")
	assert.Equal(t, "1/20 Liter", tuschegefaess.Beschreibung, "Tuschegefäß description mismatch")

	// Verify items contained in Lederrucksack (moam-container-53059)
	itemsInBackpack := 0
	for i := range bmrt.Ausruestung {
		if bmrt.Ausruestung[i].BeinhaltetIn == "moam-container-53059" {
			itemsInBackpack++
		}
	}
	assert.Equal(t, 6, itemsInBackpack, "Expected 6 items in Lederrucksack")

	// Test Transportmittel (should be empty)
	assert.Equal(t, 0, len(bmrt.Transportmittel), "Transportmittel should be empty")
	assert.NotNil(t, bmrt.Transportmittel, "Transportmittel should not be nil")

	// Test gestalt
	assert.Equal(t, "mittel", bmrt.Gestalt.Groesse, "Gestalt Größe mismatch")
	assert.Equal(t, "breit", bmrt.Gestalt.Breite, "Gestalt Breite mismatch")

	// Test erfahrungsschatz
	assert.Equal(t, 1300, bmrt.Erfahrungsschatz.Value, "Erfahrungsschatz mismatch")

	// Test Image transfer - Nicolo has a base64 image
	assert.NotEmpty(t, bmrt.Image, "Character image should be transferred")
	assert.Contains(t, bmrt.Image, "data:image;base64,", "Image should be base64 encoded")
	assert.Greater(t, len(bmrt.Image), 100, "Image data should be substantial")
	// Verify image starts with correct PNG header (after base64 prefix)
	assert.Contains(t, bmrt.Image, "iVBORw0KGgo", "Image should contain PNG base64 signature")
}

// TestImageTransfer verifies that character images are correctly transferred during conversion
func TestImageTransfer(t *testing.T) {
	// Test 1: Character WITH image (Dorin Schnellhammer)
	data, err := os.ReadFile("testdata/dorin-schnellhammer.json")
	assert.NoError(t, err, "Failed to read Dorin test data")

	var dorinChar MoamCharacter
	err = json.Unmarshal(data, &dorinChar)
	assert.NoError(t, err, "Failed to parse Dorin JSON")

	dorinBmrt, err := toBMRT(&dorinChar)
	assert.NoError(t, err, "Dorin conversion to BMRT failed")

	// Verify image was transferred
	assert.NotEmpty(t, dorinBmrt.Image, "Dorin image should be transferred")
	assert.Contains(t, dorinBmrt.Image, "data:image;base64,", "Image should have base64 prefix")
	assert.Greater(t, len(dorinBmrt.Image), 1000, "Image should have substantial data")

	// Test 2: Character WITH image (Nicolo Sikeri)
	data, err = os.ReadFile("testdata/nicolo-sikeri.json")
	assert.NoError(t, err, "Failed to read Nicolo test data")

	var nicoloChar MoamCharacter
	err = json.Unmarshal(data, &nicoloChar)
	assert.NoError(t, err, "Failed to parse Nicolo JSON")

	nicoloBmrt, err := toBMRT(&nicoloChar)
	assert.NoError(t, err, "Nicolo conversion to BMRT failed")

	// Verify image was transferred
	assert.NotEmpty(t, nicoloBmrt.Image, "Nicolo image should be transferred")
	assert.Contains(t, nicoloBmrt.Image, "data:image;base64,", "Image should have base64 prefix")

	// Test 3: Character WITHOUT image (Galaxis)
	data, err = os.ReadFile("testdata/galaxis-kessarius-rufus.json")
	assert.NoError(t, err, "Failed to read Galaxis test data")

	var galaxisChar MoamCharacter
	err = json.Unmarshal(data, &galaxisChar)
	assert.NoError(t, err, "Failed to parse Galaxis JSON")

	galaxisBmrt, err := toBMRT(&galaxisChar)
	assert.NoError(t, err, "Galaxis conversion to BMRT failed")

	// Verify no image field (empty string is acceptable)
	assert.Empty(t, galaxisBmrt.Image, "Galaxis should have no image")

	// Test 4: Round-trip conversion preserves image
	convertedMoam, err := fromBMRT(dorinBmrt)
	assert.NoError(t, err, "Failed to convert back to MOAM")

	// Verify image survived round-trip
	assert.Equal(t, dorinBmrt.Image, convertedMoam.Image, "Image should be preserved in round-trip conversion")
	assert.NotEmpty(t, convertedMoam.Image, "Round-trip image should not be empty")
}
