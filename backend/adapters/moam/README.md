# Moam VTT Adapter

This adapter service provides import/export functionality for Moam VTT character formats into BaMoRT's BMRT format (based on `importer.CharacterImport`).

## Overview

The Moam VTT adapter is a microservice that:
- Detects Moam VTT JSON character files
- Converts Moam VTT format to BMRT (BaMoRT interchange format)
- Exports BMRT format back to Moam VTT format
- Supports Moam game versions: 10.x, 11.x, 12.x

## Architecture

### Format Compatibility

Moam VTT format is structurally very similar to BMRT format, as both share a common origin. The adapter primarily:
- Validates the Moam-specific ID format (`moam-character-*`)
- Ensures all collections are initialized (never nil)
- Preserves container hierarchies (`beinhaltet_in` relationships)
- Maintains magical item properties

### Conversion Strategy

Since `MoamCharacter` embeds `importer.CharacterImport`, the conversion is mostly a direct pass-through with initialization of empty collections. Moam-specific fields (like `Stand` for social status) are currently dropped during conversion but could be preserved in a future `Extensions` field.

## API Endpoints

### GET /metadata
Returns adapter capabilities and version information.

**Response:**
```json
{
  "id": "moam-vtt-v1",
  "name": "Moam VTT Character",
  "version": "1.0",
  "bmrt_versions": ["1.0"],
  "supported_extensions": [".json"],
  "supported_game_versions": ["10.x", "11.x", "12.x"],
  "capabilities": ["import", "export", "detect"]
}
```

### POST /detect
Analyzes JSON data and returns confidence score (0.0 to 1.0).

**Request:** Raw JSON file bytes

**Response:**
```json
{
  "confidence": 0.95,
  "version": "10.x"
}
```

**Detection Logic:**
1. Valid JSON structure (+0.2 confidence)
2. ID starts with `moam-character-` (+0.3 confidence)
3. Required fields present (name, eigenschaften, grad) (+0.3 confidence)
4. Expected collections exist (fertigkeiten, waffenfertigkeiten, waffen) (+0.2 confidence)

### POST /import
Converts Moam VTT JSON to BMRT format.

**Request:** Raw Moam VTT JSON bytes

**Response:** `CharacterImport` JSON (BMRT format)

**Error Codes:**
- `400`: Invalid request body
- `422`: Invalid Moam JSON or conversion failed

### POST /export
Converts BMRT format to Moam VTT JSON.

**Request:** `CharacterImport` JSON (BMRT format)

**Response:** Moam VTT JSON bytes

**Error Codes:**
- `400`: Invalid request body
- `422`: Invalid BMRT format or conversion failed

### GET /health
Health check endpoint for container orchestration.

**Response:**
```json
{
  "status": "healthy"
}
```

## Development

### Running Tests

```bash
cd /data/dev/bamort/backend/adapters/moam
go test -v
```

**Test Coverage:**
- Format detection (valid and invalid)
- Moam to BMRT conversion
- BMRT to Moam round-trip conversion
- Empty character handling
- Magical item preservation
- Container hierarchy preservation

### Running Locally

```bash
# Build the adapter
go build -o adapter-moam .

# Run with default port (8181)
./adapter-moam

# Run with custom port
PORT=9000 ./adapter-moam
```

### Docker Development

The adapter runs in a Docker container with Air for live-reload during development:

```bash
cd /data/dev/bamort/docker
docker-compose -f docker-compose.dev.yml up -d adapter-moam-dev

# View logs
docker logs -f bamort-adapter-moam-dev

# Test endpoints
curl http://localhost:8181/metadata
curl -X POST http://localhost:8181/detect -d @testdata/moam_character.json
```

## Test Data

Sample test character: `testdata/moam_character.json`

Contains a minimal Moam VTT character with:
- Basic stats (eigenschaften)
- Skills (fertigkeiten)
- Weapon skills (waffenfertigkeiten)
- Equipment, weapons, containers
- LP/AP/B values

## Implementation Notes

### Go Version
Requires Go 1.24+ (matches backend requirements)

### Dependencies
- `github.com/gin-gonic/gin` - HTTP framework
- `bamort/importer` - BMRT format types

### Code Style
Follows [Go instructions](.github/instructions/go.instructions.md):
- Idiomatic Go code
- Clear error handling
- Comprehensive documentation
- Test-driven development

### TDD Approach
This adapter was developed following TDD principles:
1. Write failing tests first
2. Implement minimal code to pass tests
3. Refactor while keeping tests green
4. Verify with integration tests

### KISS Principle
Implementation follows Keep It Simple principles:
- Direct struct embedding (MoamCharacter embeds CharacterImport)
- No complex transformations (formats are compatible)
- Clear, readable code over clever optimizations
- Explicit error messages

## Integration with BaMoRT

The adapter is registered with the main BaMoRT backend via the `IMPORT_ADAPTERS` environment variable:

```yaml
environment:
  - IMPORT_ADAPTERS=[{"id":"moam-vtt-v1","base_url":"http://adapter-moam-dev:8181"}]
```

The backend's `importer` package:
1. Discovers the adapter on startup
2. Uses it for format detection
3. Calls import/export endpoints as needed
4. Handles health checks and failover

## Future Enhancements

### Planned Features
1. **Extensions Support**: Preserve Moam-specific fields (like `Stand`) in BMRT `Extensions` field
2. **Version Detection**: Better detection of specific Moam versions (10.x vs 11.x vs 12.x)
3. **Validation**: Deeper validation of Moam-specific business rules
4. **Error Details**: More detailed error messages with field-level feedback

### Extensibility
Adding support for new Moam versions:
1. Update `SupportedVersions` in metadata
2. Add version-specific detection logic in `detectMoamFormat()`
3. Add version-specific conversion logic if needed
4. Add test data for new version

## Troubleshooting

### Adapter Won't Start
- Check Go version: `go version` (must be 1.24+)
- Verify dependencies: `go mod tidy`
- Check port availability: `lsof -i :8181`

### Low Confidence Detection
- Verify JSON structure matches Moam format
- Check for `moam-character-` prefix in ID
- Ensure required fields are present

### Conversion Failures
- Check Moam JSON is valid: `jq . < file.json`
- Verify all required fields exist
- Check logs for specific error messages

### Health Check Failures
- Ensure container can bind to port 8181
- Check Air is running: `docker logs bamort-adapter-moam-dev`
- Verify no firewall blocking health check endpoint

## License

Same license as BaMoRT main project.
