package main

import (
	"bamort/importer"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// MoamCharacter represents the Moam VTT JSON character format.
// This structure mirrors importer.CharacterImport since Moam VTT
// uses a compatible format. We use composition to add Moam-specific metadata.
type MoamCharacter struct {
	importer.CharacterImport
	// Stand is Moam-specific (social status), not in BMRT
	Stand string `json:"stand,omitempty"`
}

// AdapterMetadata represents this adapter's capabilities and version information
type AdapterMetadata struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Version             string   `json:"version"`
	BmrtVersions        []string `json:"bmrt_versions"`
	SupportedExtensions []string `json:"supported_extensions"`
	SupportedVersions   []string `json:"supported_game_versions"`
	Capabilities        []string `json:"capabilities"`
	Description         string   `json:"description"`
}

// DetectResponse contains the confidence score and detected version
type DetectResponse struct {
	Confidence float64 `json:"confidence"`
	Version    string  `json:"version,omitempty"`
}

// ErrorResponse is returned on error
type ErrorResponse struct {
	Error string `json:"error"`
}

// main initializes the Gin router and starts the HTTP server
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8181"
	}

	// Use release mode in production
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Register routes
	r.GET("/metadata", metadataHandler)
	r.POST("/detect", detectHandler)
	r.POST("/import", importHandler)
	r.POST("/export", exportHandler)
	r.GET("/health", healthHandler)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Moam VTT Adapter starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// metadataHandler returns adapter capabilities and version information
func metadataHandler(c *gin.Context) {
	metadata := AdapterMetadata{
		ID:                  "moam-vtt-v1",
		Name:                "Moam VTT Character",
		Version:             "1.0",
		BmrtVersions:        []string{"1.0"},
		SupportedExtensions: []string{".json"},
		SupportedVersions:   []string{"5.x"},
		Capabilities:        []string{"import", "export", "detect"},
		Description:         "Adapter for importing and exporting characters in the Moam VTT JSON format Supports Characters created for Midgard Version 5.",
	}

	c.JSON(http.StatusOK, metadata)
}

// detectHandler analyzes the provided data and returns confidence score
func detectHandler(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to read request body"})
		return
	}

	confidence, version := detectMoamFormat(data)
	c.JSON(http.StatusOK, DetectResponse{
		Confidence: confidence,
		Version:    version,
	})
}

// importHandler converts Moam VTT JSON to BMRT format
func importHandler(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to read request body"})
		return
	}

	var moamChar MoamCharacter
	if err := json.Unmarshal(data, &moamChar); err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Error: fmt.Sprintf("Invalid Moam JSON format: %v", err),
		})
		return
	}

	bmrt, err := toBMRT(&moamChar)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Error: fmt.Sprintf("Conversion to BMRT failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, bmrt)
}

// exportHandler converts BMRT format back to Moam VTT JSON
func exportHandler(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Failed to read request body"})
		return
	}

	var bmrt importer.CharacterImport
	if err := json.Unmarshal(data, &bmrt); err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Error: fmt.Sprintf("Invalid BMRT format: %v", err),
		})
		return
	}

	moamChar, err := fromBMRT(&bmrt)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
			Error: fmt.Sprintf("Conversion from BMRT failed: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, moamChar)
}

// healthHandler provides a simple health check endpoint
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}

// detectMoamFormat analyzes the JSON data and returns a confidence score (0.0 to 1.0)
// and detected Moam version if applicable.
//
// Detection strategy:
// 1. Must be valid JSON (confidence += 0.2)
// 2. Must have "id" field starting with "moam-" (confidence += 0.3)
// 3. Must have standard character fields (name, eigenschaften, etc.) (confidence += 0.3)
// 4. Has fertigkeiten/waffenfertigkeiten arrays (confidence += 0.2)
//
// Returns: (confidence, version)
func detectMoamFormat(data []byte) (float64, string) {
	var confidence float64
	var version string

	// Try to parse as JSON
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return 0.0, ""
	}
	confidence += 0.2

	// Check for Moam-specific ID format
	if id, ok := raw["id"].(string); ok {
		if strings.HasPrefix(id, "moam-character-") {
			confidence += 0.3
			version = "5.x" // Default version, could be refined
		}
	}

	// Check for standard character structure
	requiredFields := []string{"name", "eigenschaften", "grad"}
	for _, field := range requiredFields {
		if _, ok := raw[field]; ok {
			confidence += 0.1
		}
	}

	// Check for typical Moam collections
	collections := []string{"fertigkeiten", "waffenfertigkeiten", "waffen"}
	collectionsFound := 0
	for _, field := range collections {
		if _, ok := raw[field]; ok {
			collectionsFound++
		}
	}
	if collectionsFound > 0 {
		confidence += 0.2
	}

	return confidence, version
}

// toBMRT converts a MoamCharacter to importer.CharacterImport (BMRT format)
//
// Since Moam VTT format is already very close to BMRT format, this is mostly
// a direct mapping. The MoamCharacter embeds CharacterImport, so we can
// return the embedded struct directly after any Moam-specific processing.
func toBMRT(moam *MoamCharacter) (*importer.CharacterImport, error) {
	if moam == nil {
		return nil, fmt.Errorf("moam character is nil")
	}

	// Since MoamCharacter embeds CharacterImport, we can use it directly
	bmrt := &moam.CharacterImport

	// Ensure all slices are initialized (never nil)
	if bmrt.Fertigkeiten == nil {
		bmrt.Fertigkeiten = []importer.Fertigkeit{}
	}
	if bmrt.Waffenfertigkeiten == nil {
		bmrt.Waffenfertigkeiten = []importer.Waffenfertigkeit{}
	}
	if bmrt.Zauber == nil {
		bmrt.Zauber = []importer.Zauber{}
	}
	if bmrt.Waffen == nil {
		bmrt.Waffen = []importer.Waffe{}
	}
	if bmrt.Ausruestung == nil {
		bmrt.Ausruestung = []importer.Ausruestung{}
	}
	if bmrt.Behaeltnisse == nil {
		bmrt.Behaeltnisse = []importer.Behaeltniss{}
	}
	if bmrt.Transportmittel == nil {
		bmrt.Transportmittel = []importer.Transportation{}
	}
	if bmrt.Spezialisierung == nil {
		bmrt.Spezialisierung = []string{}
	}
	if moam.Stand == "" {
		bmrt.SocialClass = "Mittelschicht" // Default social class if not specified
	} else {
		bmrt.SocialClass = moam.Stand
	}

	// Note: Moam-specific fields like "Stand" are dropped during conversion
	// If we needed to preserve them, we would store them in Extensions field
	// (to be added in importer/bmrt.go wrapper)

	return bmrt, nil
}

// fromBMRT converts importer.CharacterImport (BMRT format) to MoamCharacter
//
// This is primarily for export functionality. Since the formats are compatible,
// this is a straightforward conversion.
func fromBMRT(bmrt *importer.CharacterImport) (*MoamCharacter, error) {
	if bmrt == nil {
		return nil, fmt.Errorf("bmrt character is nil")
	}

	// Create MoamCharacter and copy BMRT data
	moam := &MoamCharacter{
		CharacterImport: *bmrt,
	}
	moam.Stand = bmrt.SocialClass // Map social class back to Moam's "Stand"

	// Could set Moam-specific defaults here if needed
	// For example: moam.Stand = "Nicht festgelegt"

	return moam, nil
}
