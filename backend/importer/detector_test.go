package importer

import (
	"testing"
)

func TestDetector_ExtensionMatch(t *testing.T) {
	registry := NewAdapterRegistry()

	// Register adapter with .json extension and a BaseURL
	registry.Register(AdapterMetadata{
		ID:                  "json-adapter",
		Name:                "JSON Adapter",
		SupportedExtensions: []string{".json"},
		Capabilities:        []string{"detect"},
		BaseURL:             "http://localhost:8181",
		Healthy:             true,
	})

	detector := NewDetector(registry)

	// Test extension match - should short-circuit with single match
	data := []byte(`{"name": "test"}`)
	adapterID, confidence, err := detector.DetectFormat(data, "character.json", "")
	if err != nil {
		t.Fatalf("DetectFormat failed: %v", err)
	}

	if adapterID != "json-adapter" {
		t.Errorf("Expected json-adapter, got %s", adapterID)
	}
	if confidence != 1.0 {
		t.Errorf("Expected confidence 1.0 for extension match, got %f", confidence)
	}
}

func TestDetector_ExtensionMatchMultiple(t *testing.T) {
	registry := NewAdapterRegistry()

	// Register two adapters with same extension
	registry.Register(AdapterMetadata{
		ID:                  "adapter-1",
		Name:                "Adapter 1",
		SupportedExtensions: []string{".json"},
		Capabilities:        []string{"detect"},
		Healthy:             true,
		BaseURL:             "http://localhost:8181",
	})
	registry.Register(AdapterMetadata{
		ID:                  "adapter-2",
		Name:                "Adapter 2",
		SupportedExtensions: []string{".json"},
		Capabilities:        []string{"detect"},
		Healthy:             true,
		BaseURL:             "http://localhost:8182",
	})

	detector := NewDetector(registry)

	// With multiple matches, should skip extension match and use full detection
	// Since we don't have a real server, this will fail - that's expected
	data := []byte(`{"name": "test"}`)
	_, _, err := detector.DetectFormat(data, "character.json", "")
	// We expect this to fail because no real servers are running
	if err == nil {
		t.Error("Expected error when no adapters can detect")
	}
}

func TestDetector_SpecifiedAdapter(t *testing.T) {
	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:      "specified-adapter",
		Name:    "Specified",
		BaseURL: "http://localhost:8181",
		Healthy: true,
	})

	detector := NewDetector(registry)

	data := []byte(`{"name": "test"}`)
	adapterID, confidence, err := detector.DetectFormat(data, "test.json", "specified-adapter")
	if err != nil {
		t.Fatalf("DetectFormat failed: %v", err)
	}

	if adapterID != "specified-adapter" {
		t.Errorf("Expected specified-adapter, got %s", adapterID)
	}
	if confidence != 1.0 {
		t.Errorf("Expected confidence 1.0 for specified adapter, got %f", confidence)
	}
}

func TestDetector_SignatureCache(t *testing.T) {
	registry := NewAdapterRegistry()
	detector := NewDetector(registry)

	data := []byte(`{"name": "test character with some data to make it longer than usual"}`)

	// First detection - cache miss
	detector.cacheDetection(data, "test-adapter", 0.95)

	// Second detection - cache hit
	adapterID, confidence := detector.getCachedDetection(data)
	if adapterID != "test-adapter" {
		t.Errorf("Expected cached adapter test-adapter, got %s", adapterID)
	}
	if confidence != 0.95 {
		t.Errorf("Expected cached confidence 0.95, got %f", confidence)
	}
}

func TestDetector_SignatureCacheTTL(t *testing.T) {
	registry := NewAdapterRegistry()
	detector := NewDetector(registry)
	detector.cacheTTL = 0 // Set TTL to 0 for immediate expiration

	data := []byte(`{"name": "test"}`)

	// Cache a detection
	detector.cacheDetection(data, "test-adapter", 0.95)

	// Should be expired immediately
	adapterID, _ := detector.getCachedDetection(data)
	if adapterID != "" {
		t.Errorf("Expected cache miss after TTL expiration, got %s", adapterID)
	}
}

func TestDetector_ComputeSignature(t *testing.T) {
	detector := NewDetector(nil)

	data1 := []byte(`{"name": "test"}`)
	data2 := []byte(`{"name": "test"}`)
	data3 := []byte(`{"name": "different"}`)

	sig1 := detector.computeSignature(data1)
	sig2 := detector.computeSignature(data2)
	sig3 := detector.computeSignature(data3)

	if sig1 != sig2 {
		t.Error("Same data should produce same signature")
	}
	if sig1 == sig3 {
		t.Error("Different data should produce different signature")
	}
}

func TestDetector_ExtensionMatching(t *testing.T) {
	// Test various file extensions
	tests := []struct {
		filename  string
		extension string
	}{
		{"character.json", ".json"},
		{"CHARACTER.JSON", ".json"},
		{"data.xml", ".xml"},
		{"file.txt", ".txt"},
		{"noextension", ""},
		{"multi.part.json", ".json"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			ext := getFileExtension(tt.filename)
			if ext != tt.extension {
				t.Errorf("getFileExtension(%s) = %s, want %s", tt.filename, ext, tt.extension)
			}
		})
	}
}
