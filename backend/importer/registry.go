package importer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// AdapterMetadata contains information about a registered adapter service
type AdapterMetadata struct {
	ID                  string    `json:"id"`                   // e.g., "foundry-vtt-v1"
	Name                string    `json:"name"`                 // e.g., "Foundry VTT Character"
	Version             string    `json:"version"`              // e.g., "1.0"
	BmrtVersions        []string  `json:"bmrt_versions"`        // Supported BMRT versions
	SupportedExtensions []string  `json:"supported_extensions"` // e.g., [".json"]
	BaseURL             string    `json:"base_url"`             // e.g., "http://adapter-foundry:8181"
	Capabilities        []string  `json:"capabilities"`         // e.g., ["import", "export", "detect"]
	Healthy             bool      `json:"healthy"`              // Runtime health status
	LastCheckedAt       time.Time `json:"last_checked_at"`
	LastError           string    `json:"last_error,omitempty"`
}

// SupportsCapability checks if the adapter supports a specific capability
func (m *AdapterMetadata) SupportsCapability(capability string) bool {
	for _, cap := range m.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

// AdapterRegistry manages registered adapter services
type AdapterRegistry struct {
	adapters          map[string]*AdapterMetadata
	mu                sync.RWMutex
	client            *http.Client
	stopHealthChecker chan struct{}
}

// NewAdapterRegistry creates a new adapter registry
func NewAdapterRegistry() *AdapterRegistry {
	return &AdapterRegistry{
		adapters: make(map[string]*AdapterMetadata),
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Disable redirects for security
			},
		},
		stopHealthChecker: make(chan struct{}),
	}
}

// Register registers or updates an adapter
func (r *AdapterRegistry) Register(meta AdapterMetadata) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Validate required fields
	if meta.ID == "" {
		return fmt.Errorf("adapter ID is required")
	}
	if meta.BaseURL == "" {
		return fmt.Errorf("adapter base URL is required")
	}

	// Register/update the adapter
	r.adapters[meta.ID] = &meta
	return nil
}

// StartBackgroundHealthChecker starts a background goroutine that periodically checks adapter health
// Runs every 30 seconds as specified in the plan
func (r *AdapterRegistry) StartBackgroundHealthChecker() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				_ = r.HealthCheck() // Ignore errors, they're logged in adapter metadata
			case <-r.stopHealthChecker:
				ticker.Stop()
				return
			}
		}
	}()
}

// StopBackgroundHealthChecker stops the background health checker
func (r *AdapterRegistry) StopBackgroundHealthChecker() {
	close(r.stopHealthChecker)
}

// Get retrieves an adapter by ID
func (r *AdapterRegistry) Get(id string) *AdapterMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.adapters[id]
}

// GetAll returns all registered adapters
func (r *AdapterRegistry) GetAll() []*AdapterMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*AdapterMetadata, 0, len(r.adapters))
	for _, adapter := range r.adapters {
		result = append(result, adapter)
	}
	return result
}

// GetHealthy returns only healthy adapters
func (r *AdapterRegistry) GetHealthy() []*AdapterMetadata {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*AdapterMetadata, 0)
	for _, adapter := range r.adapters {
		if adapter.Healthy {
			result = append(result, adapter)
		}
	}
	return result
}

// Import calls an adapter's import endpoint
func (r *AdapterRegistry) Import(adapterID string, data []byte) (*BMRTCharacter, error) {
	adapter := r.Get(adapterID)
	if adapter == nil {
		return nil, fmt.Errorf("adapter not found: %s", adapterID)
	}

	if !adapter.Healthy {
		return nil, fmt.Errorf("adapter is unhealthy: %s", adapterID)
	}

	if !adapter.SupportsCapability("import") {
		return nil, fmt.Errorf("adapter does not support import: %s", adapterID)
	}

	url := adapter.BaseURL + "/import"
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("adapter request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("adapter returned status %d: %s", resp.StatusCode, string(body))
	}

	var char BMRTCharacter
	if err := json.NewDecoder(resp.Body).Decode(&char); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &char, nil
}

// Export calls an adapter's export endpoint
func (r *AdapterRegistry) Export(adapterID string, char *CharacterImport) ([]byte, error) {
	adapter := r.Get(adapterID)
	if adapter == nil {
		return nil, fmt.Errorf("adapter not found: %s", adapterID)
	}

	if !adapter.Healthy {
		return nil, fmt.Errorf("adapter is unhealthy: %s", adapterID)
	}

	if !adapter.SupportsCapability("export") {
		return nil, fmt.Errorf("adapter does not support export: %s", adapterID)
	}

	// Marshal character to JSON
	charData, err := json.Marshal(char)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal character: %w", err)
	}

	url := adapter.BaseURL + "/export"
	req, err := http.NewRequest("POST", url, bytes.NewReader(charData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("adapter request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("adapter returned status %d: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}

// DetectResponse represents the response from an adapter's detect endpoint
type DetectResponse struct {
	Confidence float64 `json:"confidence"`
	Version    string  `json:"version,omitempty"`
}

// Detect calls all healthy adapters' detect endpoints and returns the best match
func (r *AdapterRegistry) Detect(data []byte, filename string) (string, float64, error) {
	healthy := r.GetHealthy()
	if len(healthy) == 0 {
		return "", 0, fmt.Errorf("no healthy adapters available")
	}

	// Create a short timeout client for detection
	detectClient := &http.Client{
		Timeout: 2 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	bestAdapterID := ""
	bestConfidence := 0.0

	for _, adapter := range healthy {
		if !adapter.SupportsCapability("detect") {
			continue
		}

		url := adapter.BaseURL + "/detect"
		req, err := http.NewRequest("POST", url, bytes.NewReader(data))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/octet-stream")

		resp, err := detectClient.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var detectResp DetectResponse
			if err := json.NewDecoder(resp.Body).Decode(&detectResp); err == nil {
				if detectResp.Confidence > bestConfidence {
					bestConfidence = detectResp.Confidence
					bestAdapterID = adapter.ID
				}
			}
		}
		resp.Body.Close()
	}

	// Require minimum confidence threshold
	if bestConfidence < 0.7 {
		return "", bestConfidence, fmt.Errorf("no adapter reached confidence threshold (best: %.2f)", bestConfidence)
	}

	return bestAdapterID, bestConfidence, nil
}

// HealthCheck performs health checks on all adapters
func (r *AdapterRegistry) HealthCheck() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, adapter := range r.adapters {
		// Try to ping the metadata endpoint
		url := adapter.BaseURL + "/metadata"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			adapter.Healthy = false
			adapter.LastError = err.Error()
			adapter.LastCheckedAt = time.Now()
			continue
		}

		resp, err := r.client.Do(req)
		if err != nil {
			adapter.Healthy = false
			adapter.LastError = err.Error()
			adapter.LastCheckedAt = time.Now()
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			adapter.Healthy = true
			adapter.LastError = ""
			adapter.LastCheckedAt = time.Now()
		} else {
			adapter.Healthy = false
			adapter.LastError = fmt.Sprintf("status code: %d", resp.StatusCode)
			adapter.LastCheckedAt = time.Now()
		}

		r.adapters[id] = adapter
	}

	return nil
}
