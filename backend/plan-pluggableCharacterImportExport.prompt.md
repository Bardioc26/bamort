# Plan: Pluggable Character Import/Export with Microservice Adapters

This plan extends the existing `importero` package into a full-featured, extensible import/export system using Docker-based adapter microservices. The canonical `CharacterImport` format becomes the system-wide interchange format (BMRT-Format), and new external formats (starting with Foundry VTT) are handled by isolated adapter services. New master data is automatically flagged as personal items (house rules).

**Key Decisions**:
- Microservice architecture for adapters (Docker containers)
- Auto-flag imported master data as personal items
- Foundry VTT JSON as first format
- Backend-only implementation (no Vue components)
- Keep [transfero/](transfero/) untouched (BaMoRT-to-BaMoRT transfers)
- Extend [importero/](importero/) as the adapter orchestration layer

## 1. Core Infrastructure (Backend)

### 1.1 Formalize BMRT-Format
- Document [importero/model.go](importero/model.go) `CharacterImport` as the canonical interchange format
- Add JSON schema validation using `github.com/xeipuuv/gojsonschema`
- Add `BmrtVersion` field to `CharacterImport` (start at "1.0")
- Add `SourceMetadata` struct to track original format, adapter ID, import timestamp
- Update existing VTT adapter to populate these fields

### 1.2 Database Migrations
Add new tables to [models/model_character.go](models/model_character.go):

```go
type ImportHistory struct {
    ID              uint   `gorm:"primaryKey"`
    UserID          uint   `gorm:"not null;index"`
    CharacterID     uint   `gorm:"index"`
    AdapterID       string `gorm:"type:varchar(100);not null"` // "foundry-vtt-v1"
    SourceFormat    string `gorm:"type:varchar(50)"`           // "foundry-vtt"
    SourceFilename  string
    SourceSnapshot  []byte `gorm:"type:MEDIUMBLOB"`            // Original file
    MappingSnapshot []byte `gorm:"type:JSON"`                  // Adapter->BMRT mappings
    ImportedAt      time.Time
    Status          string `gorm:"type:varchar(20)"` // "success", "partial", "failed"
    ErrorLog        string `gorm:"type:TEXT"`
}

type MasterDataImport struct {
    ID              uint `gorm:"primaryKey"`
    ImportHistoryID uint `gorm:"not null;index"`
    ItemType        string `gorm:"type:varchar(20)"` // "skill", "spell", "weapon", "equipment"
    ItemID          uint   `gorm:"not null"`
    ExternalName    string
    MatchType       string `gorm:"type:varchar(20)"` // "exact", "created_personal"
    CreatedAt       time.Time
}
```

Add to [models/database.go](models/database.go) `MigrateStructure()` function
Add to [models/model_character.go](models/model_character.go) migration function

### 1.3 Adapter Service Registry
Create [importero/registry.go](importero/registry.go):

```go
type AdapterMetadata struct {
    ID                  string   // "foundry-vtt-v1"
    Name                string   // "Foundry VTT Character"
    Version             string   // "1.0"
    SupportedExtensions []string // [".json"]
    BaseURL             string   // "http://adapter-foundry:8181"
    Capabilities        []string // ["import", "export", "detect"]
}

type AdapterRegistry struct {
    adapters map[string]*AdapterMetadata
    mu       sync.RWMutex
}

func (r *AdapterRegistry) Register(meta AdapterMetadata) error
func (r *AdapterRegistry) Detect(data []byte) (string, float64, error) // Returns adapter ID + confidence
func (r *AdapterRegistry) Import(adapterID string, data []byte) (*CharacterImport, error)
func (r *AdapterRegistry) Export(adapterID string, char *CharacterImport) ([]byte, error)
```

Load adapters from config on startup ([importero/routes.go](importero/routes.go)):
- Environment variable `IMPORT_ADAPTERS` (JSON array of adapter configs)
- Ping each adapter's `/metadata` endpoint to register
- Cache metadata in memory

### 1.4 Format Detection
Create [importero/detector.go](importero/detector.go):

```go
func DetectFormat(data []byte, filename string) (adapterID string, confidence float64, err error) {
    // Call all registered adapters' POST /detect endpoints in parallel
    // Return highest confidence match
    // Fallback to filename extension matching
}
```

### 1.5 Validation Framework
Create [importero/validator.go](importero/validator.go):

```go
type ValidationResult struct {
    Valid    bool
    Errors   []ValidationError
    Warnings []ValidationWarning
}

type ValidationRule interface {
    Validate(char *CharacterImport) ValidationResult
}

// Rules:
// - RequiredFieldsRule (name, gameSystem must exist)
// - StatsRangeRule (stats 0-100 for Midgard)
// - ReferentialIntegrityRule (skills reference valid categories)
```

Register system-specific rules by `GameSystem` field
Never block import on warnings (log only)

### 1.6 Master Data Reconciliation
Enhance [importero/importer.go](importero/importer.go) existing `CheckSkill()`, `CheckSpell()` functions:

```go
func ReconcileSkill(skill Fertigkeit, importHistoryID uint) (*models.Skill, string, error) {
    // 1. Exact match by (Name + GameSystem) → "exact"
    // 2. Not found → Create with PersonalItem=true → "created_personal"
    // 3. Log to MasterDataImport table
}
```

Apply to all types: skills, weapon skills, spells, equipment, weapons, containers
Set `PersonalItem = true` for all created master data
Chain user's `UserID` to created items via `CreatedByUserID` (add field to GSM models)

## 2. API Endpoints (Backend)

Add to [importero/routes.go](importero/routes.go):

```go
func RegisterRoutes(r *gin.RouterGroup) {
    importer := r.Group("/importer")
    
    // Existing endpoints remain
    
    // NEW endpoints:
    importer.POST("/detect", DetectHandler)           // Upload file, returns detected format
    importer.POST("/import", ImportHandler)           // Upload + import with adapter
    importer.GET("/adapters", ListAdaptersHandler)    // List registered adapters
    importer.GET("/history", ImportHistoryHandler)    // User's import history
    importer.GET("/history/:id", ImportDetailsHandler) // Details + errors
    importer.POST("/export/:id", ExportHandler)        // Export char to original format
}
```

**Handler Implementations** in [importero/handlers.go](importero/handlers.go):

**DetectHandler**:
- Accept multipart file upload
- Save to `./uploads/detect_<uuid>`
- Call `DetectFormat()`
- Return `{adapter_id, confidence, suggested_adapter_name}`
- Clean up temp file

**ImportHandler**:
- Accept `file` + optional `adapter_id` (from detect)
- Save original file to `./uploads/import_<uuid>`
- If no `adapter_id`, call `DetectFormat()`
- Call `registry.Import(adapterID, fileData)`
- Validate result with `validator.Validate()`
- Create `models.Char` via existing `CreateCharacterFromImport()` helper (new function)
- Reconcile all master data, log to `MasterDataImport`
- Save original file to `ImportHistory.SourceSnapshot`
- Return `{character_id, warnings, created_items: {skills: 3, spells: 1}}`

**ListAdaptersHandler**:
- Return `registry.GetAll()` metadata

**ImportHistoryHandler**:
- Query `ImportHistory` filtered by `userID`
- Return paginated list

**ExportHandler**:
- Load `Char` by ID (check ownership)
- Load `ImportHistory` to get original `AdapterID`
- Convert `Char` back to `CharacterImport` (reverse of import)
- Call `registry.Export(adapterID, charImport)`
- Return file download with `Content-Disposition: attachment`

## 3. Adapter Service Protocol

### 3.1 Adapter HTTP API Contract
All adapter services must implement:

**GET `/metadata`**
```json
{
  "id": "foundry-vtt-v1",
  "name": "Foundry VTT Character",
  "version": "1.0",
  "supported_extensions": [".json"],
  "capabilities": ["import", "export", "detect"]
}
```

**POST `/detect`**
- Body: raw file bytes
- Response: `{"confidence": 0.95, "version": "10.x"}`

**POST `/import`**
- Body: raw file bytes
- Response: `CharacterImport` JSON (BMRT-Format)

**POST `/export`**  
- Body: `CharacterImport` JSON
- Response: original format file bytes

### 3.2 Error Handling
- 400 Bad Request: malformed input
- 422 Unprocessable Entity: valid format but conversion failed
- 500 Internal Server Error: adapter crash

All adapter calls have 30-second timeout
Retry logic: 3 attempts with exponential backoff

## 4. Foundry VTT Adapter Service (First Implementation)

### 4.1 Docker Service
Create `docker/Dockerfile.adapter-foundry`:
```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY backend/adapters/foundry/ .
RUN go build -o adapter-foundry .

FROM alpine:latest
COPY --from=builder /app/adapter-foundry /adapter-foundry
EXPOSE 8181
CMD ["/adapter-foundry"]
```

### 4.2 Service Code
Create [backend/adapters/foundry/main.go](backend/adapters/foundry/main.go):

```go
package main

import (
    "github.com/gin-gonic/gin"
    "bamort/importero"  // Import BMRT-Format types
)

type FoundryCharacter struct {
    Name   string `json:"name"`
    System struct {
        Abilities map[string]struct {
            Value int `json:"value"`
        } `json:"abilities"`
        // ... Foundry schema
    } `json:"system"`
}

func metadata(c *gin.Context) {
    c.JSON(200, gin.H{
        "id": "foundry-vtt-v1",
        "name": "Foundry VTT Character",
        "version": "1.0",
        "supported_extensions": []string{".json"},
        "capabilities": []string{"import", "export", "detect"},
    })
}

func detect(c *gin.Context) {
    data, _ := c.GetRawData()
    // Parse JSON, check for Foundry-specific fields
    // Return confidence 0.0-1.0
}

func importChar(c *gin.Context) {
    var foundry FoundryCharacter
    c.BindJSON(&foundry)
    
    // Convert to importero.CharacterImport
    bmrt := toBMRT(foundry)
    c.JSON(200, bmrt)
}

func exportChar(c *gin.Context) {
    var bmrt importero.CharacterImport
    c.BindJSON(&bmrt)
    
    // Convert back to Foundry format
    foundry := fromBMRT(bmrt)
    c.JSON(200, foundry)
}
```

### 4.3 Conversion Logic
- Map Foundry abilities → BMRT stats (St, Gw, In...)
- Map Foundry items → BMRT equipment
- Map Foundry features → BMRT skills
- Preserve unmapped fields in `CharacterImport.Extensions["foundry"]` (add Extensions map to model)

### 4.4 Docker Compose Integration
Add to [docker/docker-compose.dev.yml](docker/docker-compose.dev.yml):

```yaml
adapter-foundry:
  build:
    context: ../
    dockerfile: docker/Dockerfile.adapter-foundry
  container_name: bamort-adapter-foundry-dev
  ports:
    - "8181:8181"
  networks:
    - bamort-network
  environment:
    - PORT=8181
  restart: unless-stopped
```

Update backend environment to register adapter:
```yaml
bamort-backend-dev:
  environment:
    - IMPORT_ADAPTERS=[{"id":"foundry-vtt-v1","base_url":"http://adapter-foundry:8181"}]
```

## 5. Testing Strategy

### 5.1 Unit Tests
Create [importero/registry_test.go](importero/registry_test.go):
- Test adapter registration
- Test detection with multiple adapters
- Mock HTTP responses using `httptest`

Create [importero/validator_test.go](importero/validator_test.go):
- Test each validation rule
- Test warning vs error distinction

### 5.2 Integration Tests
Create [importero/import_integration_test.go](importero/import_integration_test.go):
- Use `testutils.SetupTestDB()`
- Test full import flow with mock adapter
- Verify `ImportHistory` created
- Verify personal items flagged
- Test character creation

### 5.3 Adapter Tests
Create [backend/adapters/foundry/adapter_test.go](backend/adapters/foundry/adapter_test.go):
- Golden file tests: `testdata/foundry_character.json` → BMRT → compare
- Round-trip tests: Foundry → BMRT → Foundry (structural equality)
- Detection tests with sample files

### 5.4 End-to-End Tests
Create [backend/api/import_e2e_test.go](backend/api/import_e2e_test.go):
- Start real adapter service in Docker
- Upload Foundry character via API
- Verify character created
- Verify export produces valid Foundry JSON
- Use `docker-compose -f docker/docker-compose.test.yml` with test services

## 6. Documentation

### 6.1 Update Files
- [backend/PlanNewFeature.md](backend/PlanNewFeature.md) → Mark as "Implemented, see IMPORT_EXPORT_GUIDE.md"
- Create `backend/IMPORT_EXPORT_GUIDE.md` with architecture overview
- Create `backend/adapters/ADAPTER_DEVELOPMENT.md` with adapter creation guide
- Update [docker/SERVICES_REFERENCE.md](docker/SERVICES_REFERENCE.md) with adapter services

### 6.2 API Documentation
Add OpenAPI/Swagger annotations to handlers (use `swaggo/swag`)
Generate docs with `swag init`

## 7. Deployment Considerations

### 7.1 Production Configuration
- Adapter URLs from environment variables
- Health checks for adapter services
- Graceful degradation if adapter unavailable (return error, don't crash)
- Rate limiting on import endpoints (prevent abuse)

### 7.2 Monitoring
- Log all import attempts (success/failure) with `logger` package
- Metrics: imports per adapter, detection accuracy, errors by adapter
- Alert on adapter unavailability

### 7.3 File Cleanup
- Cron job to delete old uploads (>30 days)
- `ImportHistory.SourceSnapshot` compressed with gzip
- Configurable retention policy

## 8. Future Extensibility

### 8.1 Adding New Adapters
1. Create adapter service in `backend/adapters/<format>/`
2. Add Dockerfile
3. Add to `docker-compose.dev.yml`
4. Register in backend env vars
5. Deploy container
6. No backend code changes required

### 8.2 Master Data Approval Workflow (Future)
- Add `MasterDataPending` table
- Admin UI in Vue to approve/reject
- Change reconciliation to create pending records instead of auto-creating

### 8.3 Fuzzy Matching (Future)
- Add `github.com/texttheater/golang-levenshtein` for string distance
- Configurable threshold (e.g., distance < 3)
- Return suggestions to user for manual mapping

## Verification

### Step-by-Step Testing
1. Start dev environment: `cd docker && ./start-dev.sh`
2. Verify adapter container running: `docker ps | grep bamort-adapter-foundry`
3. Check adapter metadata: `curl http://localhost:8181/metadata`
4. Run backend tests: `cd backend && go test ./importero/... -v`
5. Run adapter tests: `go test ./adapters/foundry/... -v`
6. Upload test character: `curl -F "file=@testdata/foundry_sample.json" http://localhost:8180/api/importer/import -H "Authorization: Bearer <token>"`
7. Verify character created in database via phpMyAdmin
8. Check `ImportHistory` table populated
9. Export character: `curl http://localhost:8180/api/importer/export/1 -H "Authorization: Bearer <token>" -o exported.json`
10. Compare original vs exported (structural equivalence)

### Database Verification
```sql
SELECT * FROM import_histories ORDER BY imported_at DESC LIMIT 10;
SELECT * FROM master_data_imports WHERE item_type='skill';
SELECT * FROM skills WHERE personal_item = true;
```

## Key Decisions

- **Microservice vs Monolith**: Chose microservices for adapters despite added complexity, enables language-agnostic adapters and crash isolation
- **Master Data Handling**: Auto-flag as personal items (no approval workflow) to avoid blocking imports
- **Format Priority**: Foundry VTT first, enables testing of full architecture before adding more formats
- **Frontend Scope**: Backend-only to establish stable API before UI/UX work
- **BMRT-Format**: Use existing `CharacterImport` rather than create new structure, reduces refactoring
- **transfero Separation**: Keep untouched, serves different purpose (BaMoRT-to-BaMoRT lossless transfer)

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1-2)
- Database migrations
- Adapter registry
- Format detection
- Validation framework
- Master data reconciliation updates

### Phase 2: API Endpoints (Week 2-3)
- Implement all handlers
- Error handling
- File management
- Testing infrastructure

### Phase 3: Foundry Adapter (Week 3-4)
- Docker service setup
- Conversion logic
- Round-trip testing
- Integration with backend

### Phase 4: Testing & Documentation (Week 4-5)
- Comprehensive test suite
- Documentation updates
- E2E testing
- Performance testing

### Phase 5: Deployment & Monitoring (Week 5-6)
- Production configuration
- Monitoring setup
- File cleanup jobs
- Security hardening

## Success Criteria

- [ ] Foundry VTT characters import successfully
- [ ] Round-trip export produces valid Foundry JSON
- [ ] Personal items flagged automatically
- [ ] ImportHistory tracks all imports
- [ ] Adapters run in isolated Docker containers
- [ ] 90%+ test coverage on new code
- [ ] API documentation complete
- [ ] Adding new adapter requires no backend changes
- [ ] Zero data loss on import/export cycle
- [ ] Performance: <5s for typical character import
