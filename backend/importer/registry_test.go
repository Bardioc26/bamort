package importer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdapterRegistry_Register(t *testing.T) {
	registry := NewAdapterRegistry()

	meta := AdapterMetadata{
		ID:                  "test-adapter-v1",
		Name:                "Test Adapter",
		Version:             "1.0",
		BmrtVersions:        []string{"1.0"},
		SupportedExtensions: []string{".json"},
		BaseURL:             "http://localhost:8181",
		Capabilities:        []string{"import", "export", "detect"},
		Healthy:             true,
	}

	err := registry.Register(meta)
	if err != nil {
		t.Fatalf("Failed to register adapter: %v", err)
	}

	// Verify adapter is registered
	adapter := registry.Get("test-adapter-v1")
	if adapter == nil {
		t.Fatal("Adapter should be registered")
	}
	if adapter.ID != "test-adapter-v1" {
		t.Errorf("Adapter ID = %v, want %v", adapter.ID, "test-adapter-v1")
	}
	if adapter.Name != "Test Adapter" {
		t.Errorf("Adapter Name = %v, want %v", adapter.Name, "Test Adapter")
	}
}

func TestAdapterRegistry_RegisterDuplicate(t *testing.T) {
	registry := NewAdapterRegistry()

	meta := AdapterMetadata{
		ID:           "test-adapter-v1",
		Name:         "Test Adapter",
		Version:      "1.0",
		BmrtVersions: []string{"1.0"},
		BaseURL:      "http://localhost:8181",
		Healthy:      true,
	}

	err := registry.Register(meta)
	if err != nil {
		t.Fatalf("Failed to register adapter first time: %v", err)
	}

	// Try to register again - should update, not error
	meta.Version = "1.1"
	err = registry.Register(meta)
	if err != nil {
		t.Fatalf("Failed to update adapter: %v", err)
	}

	// Verify it was updated
	adapter := registry.Get("test-adapter-v1")
	if adapter.Version != "1.1" {
		t.Errorf("Adapter version = %v, want %v", adapter.Version, "1.1")
	}
}

func TestAdapterRegistry_GetHealthy(t *testing.T) {
	registry := NewAdapterRegistry()

	// Register healthy adapter
	registry.Register(AdapterMetadata{
		ID:           "healthy-adapter",
		Name:         "Healthy",
		BmrtVersions: []string{"1.0"},
		BaseURL:      "http://localhost:8181",
		Healthy:      true,
	})

	// Register unhealthy adapter
	registry.Register(AdapterMetadata{
		ID:           "unhealthy-adapter",
		Name:         "Unhealthy",
		BmrtVersions: []string{"1.0"},
		BaseURL:      "http://localhost:8182",
		Healthy:      false,
	})

	healthy := registry.GetHealthy()
	if len(healthy) != 1 {
		t.Errorf("Expected 1 healthy adapter, got %d", len(healthy))
	}
	if len(healthy) > 0 && healthy[0].ID != "healthy-adapter" {
		t.Errorf("Healthy adapter ID = %v, want %v", healthy[0].ID, "healthy-adapter")
	}
}

func TestAdapterRegistry_GetAll(t *testing.T) {
	registry := NewAdapterRegistry()

	registry.Register(AdapterMetadata{
		ID:           "adapter-1",
		Name:         "Adapter 1",
		BmrtVersions: []string{"1.0"},
		BaseURL:      "http://localhost:8181",
		Healthy:      true,
	})

	registry.Register(AdapterMetadata{
		ID:           "adapter-2",
		Name:         "Adapter 2",
		BmrtVersions: []string{"1.0"},
		BaseURL:      "http://localhost:8182",
		Healthy:      true,
	})

	all := registry.GetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 adapters, got %d", len(all))
	}
}

func TestAdapterRegistry_Import(t *testing.T) {
	// Create a mock adapter server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/import" {
			t.Errorf("Expected path /import, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Return a valid BMRT character
		resp := BMRTCharacter{
			BmrtVersion: "1.0",
			Metadata: SourceMetadata{
				SourceFormat: "test-format",
				AdapterID:    "test-adapter-v1",
				ImportedAt:   time.Now(),
			},
		}
		resp.Name = "Test Character"

		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:           "test-adapter-v1",
		Name:         "Test Adapter",
		BmrtVersions: []string{"1.0"},
		BaseURL:      server.URL,
		Capabilities: []string{"import"},
		Healthy:      true,
	})

	data := []byte(`{"name": "Test Character"}`)
	char, err := registry.Import("test-adapter-v1", data)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	if char.Name != "Test Character" {
		t.Errorf("Character name = %v, want %v", char.Name, "Test Character")
	}
}

func TestAdapterRegistry_ImportNotFound(t *testing.T) {
	registry := NewAdapterRegistry()

	data := []byte(`{"name": "Test"}`)
	_, err := registry.Import("non-existent-adapter", data)
	if err == nil {
		t.Error("Expected error for non-existent adapter")
	}
}

func TestAdapterRegistry_ImportUnhealthy(t *testing.T) {
	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:      "unhealthy-adapter",
		Name:    "Unhealthy",
		BaseURL: "http://localhost:8181",
		Healthy: false,
	})

	data := []byte(`{"name": "Test"}`)
	_, err := registry.Import("unhealthy-adapter", data)
	if err == nil {
		t.Error("Expected error for unhealthy adapter")
	}
}

func TestAdapterRegistry_Detect(t *testing.T) {
	// Create a mock adapter server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/detect" {
			t.Errorf("Expected path /detect, got %s", r.URL.Path)
		}

		resp := map[string]interface{}{
			"confidence": 0.95,
			"version":    "1.0",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:           "test-adapter-v1",
		Name:         "Test Adapter",
		BmrtVersions: []string{"1.0"},
		BaseURL:      server.URL,
		Capabilities: []string{"detect"},
		Healthy:      true,
	})

	data := []byte(`{"test": "data"}`)
	adapterID, confidence, err := registry.Detect(data, "test.json")
	if err != nil {
		t.Fatalf("Detect failed: %v", err)
	}

	if adapterID != "test-adapter-v1" {
		t.Errorf("Adapter ID = %v, want %v", adapterID, "test-adapter-v1")
	}
	if confidence < 0.9 {
		t.Errorf("Confidence = %v, want >= 0.9", confidence)
	}
}

func TestAdapterRegistry_DetectNoMatch(t *testing.T) {
	// Create a mock adapter server that returns low confidence
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"confidence": 0.2,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	registry := NewAdapterRegistry()
	registry.Register(AdapterMetadata{
		ID:           "test-adapter-v1",
		Name:         "Test Adapter",
		BmrtVersions: []string{"1.0"},
		BaseURL:      server.URL,
		Capabilities: []string{"detect"},
		Healthy:      true,
	})

	data := []byte(`{"test": "data"}`)
	_, _, err := registry.Detect(data, "test.json")
	if err == nil {
		t.Error("Expected error for low confidence detection")
	}
}

func TestAdapterMetadata_SupportsCapability(t *testing.T) {
	meta := AdapterMetadata{
		Capabilities: []string{"import", "export"},
	}

	if !meta.SupportsCapability("import") {
		t.Error("Should support import capability")
	}
	if !meta.SupportsCapability("export") {
		t.Error("Should support export capability")
	}
	if meta.SupportsCapability("detect") {
		t.Error("Should not support detect capability")
	}
}
