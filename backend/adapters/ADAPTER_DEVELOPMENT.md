# Adapter Development Guide

## Overview

This guide explains how to create a new adapter microservice for importing/exporting characters from external formats into BaMoRT's BMRT format.

## Adapter Architecture

Each adapter is a standalone HTTP service that:
1. Receives raw file data from BaMoRT backend
2. Converts to/from BMRT format (BaMoRT's canonical interchange format)
3. Returns converted data or error information

### Benefits of Microservice Approach

- **Language Agnostic**: Write adapters in any language (Go, Python, Node.js, etc.)
- **Crash Isolation**: Adapter failures don't crash main backend
- **Independent Deployment**: Update adapters without backend changes
- **Easy Testing**: Test adapters independently with sample files

## Prerequisites

- Docker for containerization
- Understanding of target format (e.g., Foundry VTT JSON schema)
- Access to sample files in target format

## Adapter Contract

All adapters MUST implement 4 HTTP endpoints:

### 1. GET `/metadata`

Returns adapter capabilities and version information.

**Response Schema:**
```json
{
  "id": "string",                    // Unique ID (e.g., "foundry-vtt-v1")
  "name": "string",                  // Human-readable name
  "version": "string",               // Adapter version (semantic versioning)
  "bmrt_versions": ["string"],       // Supported BMRT versions (e.g., ["1.0"])
  "supported_extensions": ["string"], // File extensions (e.g., [".json"])
  "supported_game_versions": ["string"], // Optional: external format versions
  "capabilities": ["string"]         // ["import", "export", "detect"]
}
```

**Example:**
```json
{
  "id": "foundry-vtt-v1",
  "name": "Foundry VTT Character",
  "version": "1.0.2",
  "bmrt_versions": ["1.0"],
  "supported_extensions": [".json"],
  "supported_game_versions": ["10.x", "11.x", "12.x"],
  "capabilities": ["import", "export", "detect"]
}
```

### 2. POST `/detect`

Determines if uploaded file matches this adapter's format.

**Request:**
- Content-Type: `application/octet-stream`
- Body: Raw file bytes

**Response Schema:**
```json
{
  "confidence": 0.95,  // Float 0.0-1.0 (threshold: 0.7 for positive match)
  "version": "10.x"    // Optional: detected version of external format
}
```

**Detection Logic:**

```go
func detect(data []byte) (confidence float64, version string) {
    // 1. Parse JSON
    var obj map[string]interface{}
    if err := json.Unmarshal(data, &obj); err != nil {
        return 0.0, ""
    }
    
    confidence := 0.0
    
    // 2. Check required fields
    if _, ok := obj["system"]; ok {
        confidence += 0.3
    }
    if abilities, ok := obj["system"].(map[string]interface{})["abilities"]; ok {
        confidence += 0.3
    }
    
    // 3. Check signature fields unique to format
    if foundryVersion, ok := obj["system"].(map[string]interface{})["version"]; ok {
        confidence += 0.4
        version = detectVersion(foundryVersion)
    }
    
    return confidence, version
}
```

**Performance:** Must respond within 2 seconds (backend timeout)

### 3. POST `/import`

Converts external format to BMRT format.

**Request:**
- Content-Type: `application/octet-stream`
- Body: Raw file bytes (same as uploaded by user)

**Response:**
- Content-Type: `application/json`
- Body: BMRT CharacterImport JSON

**BMRT Format** (based on `importer.CharacterImport`):

```json
{
  "name": "Character Name",
  "grad": 1,
  "game_system": "Midgard5",
  "stats": {
    "st": 80,
    "gs": 75,
    "gw": 70,
    "ko": 85,
    "in": 65,
    "zt": 60,
    "pa": 55,
    "au": 70,
    "wk": 60
  },
  "herkunft": {
    "rasse": "Mensch",
    "typ": "Krieger",
    "stand": "Bürger"
  },
  "basics": {
    "lp": 12,
    "ap": 20,
    "alter": 25,
    "groesse": 180,
    "gewicht": 75,
    "geschlecht": "m",
    "hand": "rechts",
    "glaube": "keine"
  },
  "skills": [
    {
      "name": "Langschwert",
      "wert": 10,
      "kategorie": "Kampf"
    }
  ],
  "spells": [
    {
      "name": "Feuerball",
      "wert": 8
    }
  ],
  "equipment": [],
  "weapons": [],
  "waffen": []
}
```

**Error Handling:**
- 400 Bad Request: Malformed input (not valid file)
- 422 Unprocessable Entity: Valid file but conversion failed
- 500 Internal Server Error: Adapter crash/unexpected error

**Performance:** Must respond within 30 seconds (backend timeout)

### 4. POST `/export`

Converts BMRT format back to external format.

**Request:**
- Content-Type: `application/json`
- Body: BMRT CharacterImport JSON

**Response:**
- Content-Type: `application/json` (or format-specific)
- Body: External format file bytes

**Note:** Export is best-effort. Some BMRT fields may not have equivalents in external format.

## Step-by-Step: Creating a New Adapter

### Step 1: Project Setup

```bash
mkdir -p backend/adapters/myformat
cd backend/adapters/myformat
go mod init bamort-adapter-myformat

# Or for Python:
# python -m venv venv
# source venv/bin/activate
# pip install flask
```

### Step 2: Implement Adapter Server

**Go Example:**

```go
package main

import (
    "encoding/json"
    "io"
    "net/http"
    "github.com/gin-gonic/gin"
    "bamort/importer"
)

type MyFormatChar struct {
    Name   string                 `json:"name"`
    Level  int                    `json:"level"`
    Attrs  map[string]int         `json:"attributes"`
    Items  []MyFormatItem         `json:"items"`
}

func main() {
    r := gin.Default()
    
    r.GET("/metadata", metadataHandler)
    r.POST("/detect", detectHandler)
    r.POST("/import", importHandler)
    r.POST("/export", exportHandler)
    
    r.Run(":8182")
}

func metadataHandler(c *gin.Context) {
    c.JSON(200, gin.H{
        "id":                    "myformat-v1",
        "name":                  "My Format Adapter",
        "version":               "1.0",
        "bmrt_versions":         []string{"1.0"},
        "supported_extensions":  []string{".myformat"},
        "capabilities":          []string{"import", "export", "detect"},
    })
}

func detectHandler(c *gin.Context) {
    data, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    var myChar MyFormatChar
    if err := json.Unmarshal(data, &myChar); err != nil {
        c.JSON(200, gin.H{"confidence": 0.0})
        return
    }
    
    confidence := calculateConfidence(myChar)
    c.JSON(200, gin.H{"confidence": confidence, "version": "1.0"})
}

func importHandler(c *gin.Context) {
    data, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    
    var myChar MyFormatChar
    if err := json.Unmarshal(data, &myChar); err != nil {
        c.JSON(422, gin.H{"error": "invalid format"})
        return
    }
    
    // Convert to BMRT
    bmrt := convertToBMRT(myChar)
    c.JSON(200, bmrt)
}

func exportHandler(c *gin.Context) {
    var bmrt importer.CharacterImport
    if err := c.ShouldBindJSON(&bmrt); err != nil {
        c.JSON(400, gin.H{"error": "invalid BMRT format"})
        return
    }
    
    // Convert from BMRT
    myChar := convertFromBMRT(bmrt)
    c.JSON(200, myChar)
}

func calculateConfidence(char MyFormatChar) float64 {
    confidence := 0.0
    
    // Check required fields
    if char.Name != "" {
        confidence += 0.3
    }
    if char.Level > 0 {
        confidence += 0.3
    }
    if len(char.Attrs) > 0 {
        confidence += 0.4
    }
    
    return confidence
}

func convertToBMRT(myChar MyFormatChar) importer.CharacterImport {
    return importer.CharacterImport{
        Name:       myChar.Name,
        Grad:       uint(myChar.Level),
        GameSystem: "Midgard5",
        Stats: importer.Stats{
            St: myChar.Attrs["strength"],
            Gs: myChar.Attrs["dexterity"],
            Gw: myChar.Attrs["constitution"],
            // ... map other stats
        },
        // ... map other fields
    }
}

func convertFromBMRT(bmrt importer.CharacterImport) MyFormatChar {
    return MyFormatChar{
        Name:  bmrt.Name,
        Level: int(bmrt.Grad),
        Attrs: map[string]int{
            "strength":     bmrt.Stats.St,
            "dexterity":    bmrt.Stats.Gs,
            "constitution": bmrt.Stats.Gw,
            // ... map other stats
        },
        // ... map other fields
    }
}
```

### Step 3: Create Dockerfile

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o adapter-myformat .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/adapter-myformat /adapter-myformat
EXPOSE 8182
CMD ["/adapter-myformat"]
```

### Step 4: Add Docker Compose Service

Edit `docker/docker-compose.dev.yml`:

```yaml
services:
  adapter-myformat:
    build:
      context: ../
      dockerfile: docker/Dockerfile.adapter-myformat
    container_name: bamort-adapter-myformat-dev
    ports:
      - "8182:8182"
    networks:
      - bamort-network
    environment:
      - PORT=8182
    restart: unless-stopped
```

### Step 5: Register Adapter

Edit backend environment in `docker/docker-compose.dev.yml`:

```yaml
bamort-backend-dev:
  environment:
    - IMPORT_ADAPTERS=[
        {"id":"moam-vtt-v1","base_url":"http://adapter-moam:8181"},
        {"id":"myformat-v1","base_url":"http://adapter-myformat:8182"}
      ]
```

### Step 6: Create Test Data

Create `backend/adapters/myformat/testdata/sample.myformat`:

```json
{
  "name": "Test Character",
  "level": 3,
  "attributes": {
    "strength": 80,
    "dexterity": 75,
    "constitution": 85
  },
  "items": []
}
```

### Step 7: Write Tests

Create `backend/adapters/myformat/adapter_test.go`:

```go
package main

import (
    "encoding/json"
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestDetectMyFormat(t *testing.T) {
    data, err := os.ReadFile("testdata/sample.myformat")
    require.NoError(t, err)
    
    var char MyFormatChar
    err = json.Unmarshal(data, &char)
    require.NoError(t, err)
    
    confidence := calculateConfidence(char)
    assert.GreaterOrEqual(t, confidence, 0.7)
}

func TestConvertToBMRT(t *testing.T) {
    myChar := MyFormatChar{
        Name:  "Test",
        Level: 1,
        Attrs: map[string]int{"strength": 80},
    }
    
    bmrt := convertToBMRT(myChar)
    
    assert.Equal(t, "Test", bmrt.Name)
    assert.Equal(t, uint(1), bmrt.Grad)
    assert.Equal(t, 80, bmrt.Stats.St)
}

func TestRoundTrip(t *testing.T) {
    // Original -> BMRT -> Original
    original := MyFormatChar{
        Name:  "Round Trip Test",
        Level: 2,
        Attrs: map[string]int{"strength": 75},
    }
    
    bmrt := convertToBMRT(original)
    result := convertFromBMRT(bmrt)
    
    assert.Equal(t, original.Name, result.Name)
    assert.Equal(t, original.Level, result.Level)
}
```

Run tests:
```bash
go test -v
```

### Step 8: Build and Test

```bash
# Build adapter
docker build -t bamort-adapter-myformat -f docker/Dockerfile.adapter-myformat .

# Run adapter standalone
docker run -p 8182:8182 bamort-adapter-myformat

# Test metadata endpoint
curl http://localhost:8182/metadata

# Test with sample file
curl -X POST http://localhost:8182/import \
  -H "Content-Type: application/octet-stream" \
  --data-binary @testdata/sample.myformat
```

### Step 9: Integration Testing

Start full stack:
```bash
cd docker
./start-dev.sh
```

Test via BaMoRT API:
```bash
# Get token
TOKEN=$(curl -X POST http://localhost:8180/api/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}' | jq -r .token)

# Import via BaMoRT
curl -X POST http://localhost:8180/api/import/import \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@testdata/sample.myformat"
```

## Best Practices

### 1. Version Detection

Always detect and report external format version:

```go
func detectVersion(obj map[string]interface{}) string {
    if v, ok := obj["schema_version"].(string); ok {
        return v
    }
    
    // Fallback: heuristic version detection
    if _, ok := obj["new_field_v2"]; ok {
        return "2.x"
    }
    return "1.x"
}
```

### 2. Graceful Degradation

Handle missing optional fields:

```go
func convertToBMRT(char MyFormatChar) importer.CharacterImport {
    bmrt := importer.CharacterImport{
        Name:       char.Name,
        GameSystem: "Midgard5",
    }
    
    // Optional fields with fallbacks
    if char.Level > 0 {
        bmrt.Grad = uint(char.Level)
    } else {
        bmrt.Grad = 1 // Default
    }
    
    return bmrt
}
```

### 3. Preserve Unmapped Data

Store extra fields in Extensions:

```go
// Extensions field in BMRT wrapper
bmrt := importer.BMRTCharacter{
    CharacterImport: baseImport,
    Extensions: map[string]json.RawMessage{
        "myformat": rawExtensionData,
    },
}
```

### 4. Logging

Log all conversions for debugging:

```go
import "log"

func importHandler(c *gin.Context) {
    log.Printf("[IMPORT] Starting conversion for adapter myformat-v1")
    
    // ... conversion logic ...
    
    log.Printf("[IMPORT] Success: converted character '%s'", bmrt.Name)
}
```

### 5. Error Messages

Provide helpful error messages:

```go
if char.Name == "" {
    c.JSON(422, gin.H{
        "error": "Character name is required",
        "field": "name",
        "help": "Set the 'name' field in your character JSON"
    })
    return
}
```

## Testing Checklist

- [ ] Unit tests for detection logic
- [ ] Unit tests for BMRT conversion
- [ ] Round-trip tests (import → export → import)
- [ ] Test with real sample files
- [ ] Test with malformed input
- [ ] Test with missing optional fields
- [ ] Performance test (< 30s for import)
- [ ] Integration test with BaMoRT backend

## Deployment

### Development

```bash
cd docker
./start-dev.sh
```

### Production

1. Build production image:
```bash
docker build -t bamort-adapter-myformat:1.0 -f docker/Dockerfile.adapter-myformat .
```

2. Update production compose:
```yaml
# docker/docker-compose.yml
adapter-myformat:
  image: bamort-adapter-myformat:1.0
  restart: unless-stopped
  networks:
    - bamort-network
```

3. Deploy:
```bash
cd docker
./stop-prd.sh
./start-prd.sh
```

## Troubleshooting

### Adapter Not Detected

Check logs:
```bash
docker logs bamort-adapter-myformat-dev
```

Verify metadata endpoint:
```bash
curl http://localhost:8182/metadata
```

### Import Fails

Test adapter directly:
```bash
curl -X POST http://localhost:8182/import \
  --data-binary @testdata/sample.myformat \
  -v
```

Check backend logs:
```bash
docker logs bamort-backend-dev | grep myformat
```

### Low Detection Confidence

Adjust confidence calculation:
```go
func calculateConfidence(char MyFormatChar) float64 {
    // Add debug logging
    log.Printf("Calculating confidence for: %+v", char)
    
    confidence := 0.0
    // ... increase weights for signature fields
    
    log.Printf("Final confidence: %f", confidence)
    return confidence
}
```

## Examples

See reference implementations:
- [Moam VTT Adapter](../adapters/moam/) - Full-featured adapter
- [Simple CSV Adapter](../adapters/csv/) - Minimal example (future)

## Support

For questions or issues:
1. Check [TROUBLESHOOTING.md](../TROUBLESHOOTING.md)
2. Review existing adapter implementations
3. Open GitHub issue with adapter logs
