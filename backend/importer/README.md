# Importer Package

The `importer/` package provides a pluggable character import/export system using Docker-based microservice adapters.

## Quick Links

- **[Complete Guide](../IMPORT_EXPORT_GUIDE.md)** - Full system documentation
- **[Adapter Development](../adapters/ADAPTER_DEVELOPMENT.md)** - Create new adapters
- **[Troubleshooting](./TROUBLESHOOTING.md)** - Common issues and solutions
- **[API Documentation](http://localhost:8180/swagger/index.html)** - Swagger UI (when server running)

## Architecture

This package orchestrates character imports from external formats (e.g., Foundry VTT) using isolated adapter microservices:

```
External Format (Foundry VTT JSON)
  ↓
Adapter Microservice (Docker container)
  ↓
importer.CharacterImport (BMRT-Format - canonical interchange format)
  ↓  
importer/ package handlers (validation, reconciliation)
  ↓
models.Char (BaMoRT database)
```

## Package vs Related Packages

- **`transfero/`** - BaMoRT-to-BaMoRT lossless transfer (existing, untouched)
- **`importero/`** - Legacy format handlers (VTT JSON, CSV) with direct imports (deprecated, untouched)
- **`importer/`** - NEW microservice adapter orchestration layer (self-contained)

## Core Components

### bmrt.go
BMRT (BaMoRT Format) wrapper with source metadata tracking. Uses `CharacterImport` (defined in character.go) as the canonical interchange format.

### registry.go
Adapter service registry with health monitoring, version negotiation, and runtime failover.

### detector.go
Smart format detection with short-circuit optimization (extension match → signature cache → fan-out).

### validator.go
3-phase validation framework:
1. BMRT structural validation (JSON schema, required fields)
2. Game system semantic validation (stat ranges, referential integrity)
3. Adapter-specific validation (format compatibility)

### reconciler.go
Master data reconciliation:
- Exact match by (Name + GameSystem)
- Auto-create with PersonalItem=true flag
- Track in MasterDataImport table

### handlers.go & routes.go
HTTP API endpoints for detection, import, export, and history.

## API Endpoints

- `POST /api/import/detect` - Upload file, detect format
- `POST /api/import/import` - Import character with adapter
- `GET /api/import/adapters` - List registered adapters
- `GET /api/import/history` - User's import history
- `GET /api/import/history/:id` - Import details + errors
- `POST /api/import/export/:id` - Export character to original format

## Database Models

### ImportHistory
Tracks all import attempts with compressed source snapshots, mapping data, and error logs.

### MasterDataImport
Tracks created/matched master data items (skills, spells, equipment).

### Char Extensions
Added fields: `ImportedFromAdapter`, `ImportedAt` for provenance tracking.

## Security

- Rate limiting: 10/min detect, 5/min import, 20/min export (per user)
- File size limit: 10MB max
- JSON depth limit: 100 levels
- SSRF protection: whitelisted adapter URLs only
- No persistent disk storage: files only in DB after import

## Testing

All components follow TDD:
- Unit tests: registry_test.go, validator_test.go
- Integration tests: integration_test.go
- Use `testutils.SetupTestDB()` for database tests

## Development Guidelines

- **TDD**: Write failing test first, then implement
- **KISS**: Simplest solution that works
- **Zero modification** to importero or transfero packages
- **Self-contained**: importer package owns all its types (no importero dependency)
- **Transaction safety**: Full ACID compliance for imports

## Adapter Development

See `backend/adapters/ADAPTER_DEVELOPMENT.md` for creating new adapter microservices.

## Usage Example

```go
// Register adapter (happens at startup)
registry.Register(AdapterMetadata{
    ID: "foundry-vtt-v1",
    BaseURL: "http://adapter-foundry:8181",
})

// Import character
result, err := ImportCharacter(fileData, userID, "")
// Auto-detects format, reconciles master data, creates character
```
