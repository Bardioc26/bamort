package importer

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

// DetectionCache stores cached format detection results
type DetectionCache struct {
	signature  string
	adapterID  string
	confidence float64
	cachedAt   time.Time
}

// Detector handles smart format detection with optimizations
type Detector struct {
	registry *AdapterRegistry
	cache    map[string]*DetectionCache
	cacheMu  sync.RWMutex
	cacheTTL time.Duration
}

// NewDetector creates a new detector
func NewDetector(registry *AdapterRegistry) *Detector {
	return &Detector{
		registry: registry,
		cache:    make(map[string]*DetectionCache),
		cacheTTL: 5 * time.Minute, // Default TTL
	}
}

// DetectFormat implements smart format detection with short-circuit optimization
// Priority:
// 1. User-specified adapter (if provided)
// 2. Extension match (if single match found)
// 3. Signature cache (SHA256 of first 1KB)
// 4. Full adapter detection fan-out
func (d *Detector) DetectFormat(data []byte, filename string, specifiedAdapterID string) (string, float64, error) {
	// Step 1: User-specified adapter (highest priority)
	if specifiedAdapterID != "" {
		adapter := d.registry.Get(specifiedAdapterID)
		if adapter == nil {
			return "", 0, fmt.Errorf("specified adapter not found: %s", specifiedAdapterID)
		}
		if !adapter.Healthy {
			return "", 0, fmt.Errorf("specified adapter is unhealthy: %s", specifiedAdapterID)
		}
		return specifiedAdapterID, 1.0, nil
	}

	// Step 2: Extension match (short-circuit if single match)
	ext := getFileExtension(filename)
	if ext != "" {
		matches := d.getAdaptersByExtension(ext)
		if len(matches) == 1 {
			// Single match - short-circuit!
			return matches[0].ID, 1.0, nil
		}
		// Multiple matches - continue to full detection
	}

	// Step 3: Signature cache
	if cachedAdapterID, cachedConfidence := d.getCachedDetection(data); cachedAdapterID != "" {
		return cachedAdapterID, cachedConfidence, nil
	}

	// Step 4: Full detection fan-out to all healthy adapters
	adapterID, confidence, err := d.registry.Detect(data, filename)
	if err != nil {
		return "", 0, err
	}

	// Cache the successful detection
	d.cacheDetection(data, adapterID, confidence)

	return adapterID, confidence, nil
}

// getAdaptersByExtension returns all healthy adapters that support the given extension
func (d *Detector) getAdaptersByExtension(ext string) []*AdapterMetadata {
	healthy := d.registry.GetHealthy()
	matches := make([]*AdapterMetadata, 0)

	ext = strings.ToLower(ext)
	for _, adapter := range healthy {
		if !adapter.SupportsCapability("detect") {
			continue
		}
		for _, supportedExt := range adapter.SupportedExtensions {
			if strings.ToLower(supportedExt) == ext {
				matches = append(matches, adapter)
				break
			}
		}
	}

	return matches
}

// getCachedDetection retrieves a cached detection result if available and not expired
func (d *Detector) getCachedDetection(data []byte) (string, float64) {
	signature := d.computeSignature(data)

	d.cacheMu.RLock()
	defer d.cacheMu.RUnlock()

	cached, exists := d.cache[signature]
	if !exists {
		return "", 0
	}

	// Check if cache entry is expired
	if time.Since(cached.cachedAt) > d.cacheTTL {
		return "", 0
	}

	return cached.adapterID, cached.confidence
}

// cacheDetection stores a detection result in the cache
func (d *Detector) cacheDetection(data []byte, adapterID string, confidence float64) {
	signature := d.computeSignature(data)

	d.cacheMu.Lock()
	defer d.cacheMu.Unlock()

	d.cache[signature] = &DetectionCache{
		signature:  signature,
		adapterID:  adapterID,
		confidence: confidence,
		cachedAt:   time.Now(),
	}
}

// computeSignature computes SHA256 hash of first 1KB of data
func (d *Detector) computeSignature(data []byte) string {
	// Use first 1KB or full data if smaller
	size := 1024
	if len(data) < size {
		size = len(data)
	}

	hash := sha256.Sum256(data[:size])
	return hex.EncodeToString(hash[:])
}

// getFileExtension extracts the file extension from a filename (case-insensitive)
func getFileExtension(filename string) string {
	if filename == "" {
		return ""
	}

	// Find last dot
	lastDot := strings.LastIndex(filename, ".")
	if lastDot == -1 || lastDot == len(filename)-1 {
		return ""
	}

	return strings.ToLower(filename[lastDot:])
}
