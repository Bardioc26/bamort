# BaMoRT Import/Export System - Complete Guide

## Table of Contents
1. [System Overview](#system-overview)
2. [Architecture](#architecture)
3. [Getting Started](#getting-started)
4. [API Reference](#api-reference)
5. [Adapter Development](#adapter-development)
6. [Testing](#testing)
7. [Security](#security)
8. [Performance](#performance)
9. [Monitoring](#monitoring)
10. [Troubleshooting](#troubleshooting)

## System Overview

The BaMoRT Import/Export system provides a pluggable, microservice-based architecture for importing characters from external formats (e.g., Foundry VTT, Roll20) and exporting them back to their original formats.

### Key Features

- **Pluggable Adapters**: Add support for new formats without changing core code
- **Microservice Architecture**: Each format adapter runs as an isolated Docker container
- **BMRT Format**: Canonical interchange format based on BaMoRT's internal data model
- **Automatic Reconciliation**: Master data (skills, spells, equipment) automatically matched or created as personal items
- **Full Audit Trail**: Complete import history with source file snapshots and error logs
- **Validation Framework**: 3-phase validation (structural, semantic, adapter-specific)
- **Security First**: Rate limiting, input validation, SSRF protection
- **Test Driven**: 90%+ test coverage, E2E and performance tests

### Package vs Related Packages

BaMoRT has three separate systems for character transfer/import:

| Package | Purpose | Status |
|---------|---------|--------|
| **`transfero/`** | BaMoRT-to-BaMoRT lossless transfer | Existing (untouched) |
| **`importero/`** | Legacy VTT/CSV direct imports | Deprecated (untouched) |
| **`importer/`** | NEW microservice adapter orchestration | Active Development |

## Architecture

### Data Flow

```
External Format (Foundry VTT JSON)
  ↓
Adapter Microservice (Docker container)
  ↓
importer.CharacterImport (BMRT-Format)
  ↓
Validation Framework (3 phases)
  ↓
Master Data Reconciliation
  ↓
models.Char (BaMoRT database)
```

### Components

#### Core Backend (`backend/importer/`)

- **bmrt.go**: BMRT format wrapper with source metadata
- **registry.go**: Adapter service registry with health monitoring
- **detector.go**: Smart format detection with caching
- **validator.go**: 3-phase validation framework
- **reconciler.go**: Master data reconciliation
- **security.go**: Rate limiting, input validation, SSRF protection
- **handlers.go**: HTTP request handlers
- **routes.go**: Route registration
- **models.go**: Database models (ImportHistory, MasterDataImport)
- **import_logic.go**: Core import logic with transaction handling

#### Adapter Services (`backend/adapters/`)

Each adapter is a standalone microservice:

```
adapters/
  ├── moam/          # Moam VTT adapter
  │   ├── main.go    # Adapter HTTP server
  │   ├── adapter_test.go
  │   └── testdata/
  └── foundry/       # Future: Foundry VTT adapter
      └── ...
```

### Database Schema

#### ImportHistory Table

```sql
CREATE TABLE import_histories (
    id              INT PRIMARY KEY AUTO_INCREMENT,
    user_id         INT NOT NULL,
    character_id    INT,
    adapter_id      VARCHAR(100) NOT NULL,
    source_format   VARCHAR(50),
    source_filename VARCHAR(255),
    source_snapshot MEDIUMBLOB,           -- Compressed original file
    mapping_snapshot JSON,                 -- Adapter conversion mappings
    bmrt_version    VARCHAR(10),
    imported_at     DATETIME,
    status          VARCHAR(20),           -- 'in_progress', 'success', 'partial', 'failed'
    error_log       TEXT,
    INDEX idx_user_id (user_id),
    INDEX idx_character_id (character_id)
);
```

#### MasterDataImport Table

```sql
CREATE TABLE master_data_imports (
    id               INT PRIMARY KEY AUTO_INCREMENT,
    import_history_id INT NOT NULL,
    item_type        VARCHAR(20),          -- 'skill', 'spell', 'weapon', 'equipment'
    item_id          INT NOT NULL,
    external_name    VARCHAR(255),
    match_type       VARCHAR(20),          -- 'exact', 'created_personal'
    created_at       DATETIME,
    INDEX idx_import_history_id (import_history_id)
);
```

#### Character Provenance (added to chars table)

```sql
ALTER TABLE chars ADD COLUMN imported_from_adapter VARCHAR(100);
ALTER TABLE chars ADD COLUMN imported_at DATETIME;
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.25+
- MariaDB (or SQLite for testing)
- Backend server running

### Starting Development Environment

```bash
cd /data/dev/bamort/docker
./start-dev.sh
```

This starts:
- `bamort-backend-dev` - Main API server (port 8180)
- `bamort-adapter-moam-dev` - Moam VTT adapter (port 8181)
- `bamort-mariadb-dev` - Database (port 3306)
- `bamort-frontend-dev` - Vue.js frontend (port 5173)

### Registering an Adapter

Adapters are registered via environment variable:

```yaml
# docker/docker-compose.dev.yml
bamort-backend-dev:
  environment:
    - IMPORT_ADAPTERS=[{"id":"moam-vtt-v1","base_url":"http://adapter-moam:8181"}]
```

On startup, the backend:
1. Pings each adapter's `/metadata` endpoint
2. Verifies BMRT version compatibility
3. Registers adapter in memory
4. Starts background health checker (every 30s)

### Basic Usage Example

```bash
# 1. Detect format
curl -X POST http://localhost:8180/api/import/detect \
  -H "Authorization: Bearer <token>" \
  -F "file=@character.json"

# Response:
# {
#   "adapter_id": "moam-vtt-v1",
#   "confidence": 0.95,
#   "suggested_name": "Moam VTT Character"
# }

# 2. Import character
curl -X POST http://localhost:8180/api/import/import \
  -H "Authorization: Bearer <token>" \
  -F "file=@character.json"

# Response:
# {
#   "character_id": 123,
#   "import_id": 456,
#   "adapter_id": "moam-vtt-v1",
#   "status": "success",
#   "warnings": [],
#   "created_items": {
#     "skills": 3,
#     "spells": 1
#   }
# }

# 3. Export character
curl -X POST http://localhost:8180/api/import/export/123 \
  -H "Authorization: Bearer <token>" \
  -o exported.json
```

## API Reference

### Endpoints

All endpoints require JWT authentication and are under `/api/import` prefix.

#### POST `/detect`

Detect format of uploaded file.

**Request:**
- Content-Type: `multipart/form-data`
- Field: `file` (max 10MB)

**Response:**
```json
{
  "adapter_id": "moam-vtt-v1",
  "confidence": 0.95,
  "suggested_name": "Moam VTT Character"
}
```

**Rate Limit:** 10 requests/minute per user

---

#### POST `/import`

Import character from file.

**Request:**
- Content-Type: `multipart/form-data`
- Field: `file` (max 10MB)
- Optional query param: `adapter_id` (override auto-detection)

**Response:**
```json
{
  "character_id": 123,
  "import_id": 456,
  "adapter_id": "moam-vtt-v1",
  "status": "success",
  "warnings": [
    {
      "field": "Stats.St",
      "message": "Stat value 101 exceeds typical range (0-100)",
      "source": "gamesystem"
    }
  ],
  "created_items": {
    "skills": 3,
    "spells": 1,
    "equipment": 2
  }
}
```

**Rate Limit:** 5 requests/minute per user

**Transaction Safety:**
- Full import wrapped in database transaction
- On failure: rollback all changes, keep ImportHistory with status="failed"
- Character, master data, and import history are atomic

---

#### GET `/adapters`

List all registered adapters.

**Response:**
```json
{
  "adapters": [
    {
      "id": "moam-vtt-v1",
      "name": "Moam VTT Character",
      "version": "1.0",
      "bmrt_versions": ["1.0"],
      "supported_extensions": [".json"],
      "capabilities": ["import", "export", "detect"],
      "healthy": true,
      "last_checked_at": "2026-02-10T10:30:00Z"
    }
  ]
}
```

---

#### GET `/history`

Get user's import history.

**Query Params:**
- `page` (default: 1)
- `limit` (default: 20, max: 100)

**Response:**
```json
{
  "imports": [
    {
      "id": 456,
      "character_id": 123,
      "adapter_id": "moam-vtt-v1",
      "source_filename": "character.json",
      "imported_at": "2026-02-10T10:00:00Z",
      "status": "success"
    }
  ],
  "total": 1
}
```

---

#### GET `/history/:id`

Get detailed import history including errors.

**Response:**
```json
{
  "id": 456,
  "character_id": 123,
  "adapter_id": "moam-vtt-v1",
  "source_filename": "character.json",
  "bmrt_version": "1.0",
  "imported_at": "2026-02-10T10:00:00Z",
  "status": "success",
  "error_log": "",
  "master_data_imports": [
    {
      "item_type": "skill",
      "item_id": 789,
      "external_name": "Custom Sword Fighting",
      "match_type": "created_personal"
    }
  ]
}
```

---

#### POST `/export/:id`

Export character to original format.

**Query Params:**
- `adapter_id` (optional: override original adapter)

**Response:**
- Content-Type: `application/json` (or format-specific)
- Content-Disposition: `attachment; filename="character_123_moam.json"`
- Body: Original format file

**Error Handling:**
- 404 Not Found: Character doesn't exist
- 403 Forbidden: User doesn't own character
- 409 Conflict: Original adapter unavailable (includes suggested alternatives)

**Rate Limit:** 20 requests/minute per user

## Adapter Development

See [ADAPTER_DEVELOPMENT.md](../adapters/ADAPTER_DEVELOPMENT.md) for complete guide.

### Adapter HTTP Contract

All adapters must implement 4 endpoints:

1. **GET `/metadata`** - Return adapter capabilities
2. **POST `/detect`** - Return confidence score (0.0-1.0)
3. **POST `/import`** - Convert to BMRT format
4. **POST `/export`** - Convert from BMRT format

### Minimal Adapter Example

```go
package main

import (
    "github.com/gin-gonic/gin"
    "bamort/importer"
)

func main() {
    r := gin.Default()
    
    r.GET("/metadata", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "id": "my-adapter-v1",
            "name": "My Format Adapter",
            "version": "1.0",
            "bmrt_versions": []string{"1.0"},
            "supported_extensions": []string{".myformat"},
            "capabilities": []string{"import", "export", "detect"},
        })
    })
    
    r.POST("/detect", detectHandler)
    r.POST("/import", importHandler)
    r.POST("/export", exportHandler)
    
    r.Run(":8182")
}
```

## Testing

### Running Tests

```bash
# Unit tests
cd backend
go test ./importer/

# E2E tests
go test ./importer/e2e_test.go

# Performance tests
go test ./importer/performance_test.go

# All tests with coverage
./scripts/test-coverage.sh
```

### Coverage Target

**Minimum: 90% code coverage**

Current coverage by component:
- bmrt.go: 95%
- registry.go: 92%
- detector.go: 91%
- validator.go: 94%
- reconciler.go: 88%
- handlers.go: 90%

### Test Categories

1. **Unit Tests** (`*_test.go`): Test individual functions
2. **Integration Tests** (`import_logic_test.go`): Test full workflows with DB
3. **E2E Tests** (`e2e_test.go`): Test HTTP handlers end-to-end
4. **Performance Tests** (`performance_test.go`): Benchmark and performance targets

## Security

### Input Validation

- **File Size Limit**: 10MB max
- **JSON Depth Limit**: 100 levels max
- **Content Type**: Validated multipart/form-data
- **Filename Sanitization**: Prevent directory traversal

### Rate Limiting

Per-user rate limits:
- Detection: 10/minute
- Import: 5/minute
- Export: 20/minute

Implementation: Token bucket algorithm with sliding window

### SSRF Protection

- **Adapter URL Whitelist**: Only registered adapters can be called
- **No Redirects**: HTTP client blocks all redirects
- **Internal Network Block**: Prevent access to 127.0.0.1, 10.x, 192.168.x, etc.
- **Timeout Enforcement**: 2s for detect, 30s for import/export

### Authentication

All endpoints require JWT token in `Authorization: Bearer <token>` header.

User ID extracted from token and used for ownership checks.

## Performance

### Performance Targets

- **Import Time**: < 5 seconds for typical character (20 skills, 5 spells)
- **Detection Time**: < 2 seconds with cache
- **Export Time**: < 3 seconds

### Optimization Techniques

1. **Smart Detection with Short-Circuit**:
   - Extension match → skip detection if only one adapter
   - Signature cache (SHA256 of first 1KB)
   - Parallel API calls with 2s timeout

2. **Data Compression**:
   - Original files gzipped in database (~70% reduction)
   - Decompressed on-demand only

3. **Database Optimization**:
   - Indexes on user_id, character_id, import_history_id
   - Batch inserts for master data imports
   - Single transaction for entire import

4. **Caching**:
   - In-memory adapter registry
   - Detection signature cache (1 hour TTL)
   - Adapter health status cache (30s refresh)

### Running Benchmarks

```bash
cd backend

# Run all benchmarks
go test -bench=. -benchmem ./importer/

# Specific benchmark
go test -bench=BenchmarkImportCharacter -benchmem ./importer/

# With CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./importer/
go tool pprof cpu.prof
```

## Monitoring

### Health Checks

Background health checker runs every 30 seconds:

```go
// Pings each adapter's /metadata endpoint
// Updates Healthy status in registry
// Logs errors for failed adapters
```

Unhealthy adapters:
- Skipped during auto-detection
- Return 503 Service Unavailable if explicitly requested
- Re-checked on next health cycle

### Logging

All import attempts logged with:
- User ID
- Adapter ID
- Source filename
- Status (success/failed)
- Error messages
- Duration

### Metrics to Monitor

- Import success rate by adapter
- Average import time by adapter
- Detection accuracy (user overrides)
- Rate limit hits
- Adapter availability percentage

## Troubleshooting

See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) for detailed troubleshooting guide.

### Common Issues

**1. Adapter Not Detected**

Check adapter is running:
```bash
docker ps | grep adapter-moam
curl http://localhost:8181/metadata
```

**2. Import Fails with Validation Error**

Check error log in ImportHistory:
```sql
SELECT error_log FROM import_histories WHERE id = <import_id>;
```

**3. Character Created But Skills Missing**

Check MasterDataImport table:
```sql
SELECT * FROM master_data_imports WHERE import_history_id = <import_id>;
```

**4. Export Returns 409 Conflict**

Original adapter unavailable. Use adapter override:
```bash
curl -X POST "http://localhost:8180/api/import/export/123?adapter_id=alternate-adapter-v1"
```

## Contributing

1. Follow TDD: Write tests first
2. Follow KISS: Simplest solution that works
3. Maintain 90%+ test coverage
4. Document all public functions
5. Add integration tests for new features

## License

See [LICENSE](../../LICENSE) for license information.
