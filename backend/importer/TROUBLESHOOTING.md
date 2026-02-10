# Troubleshooting Guide - BaMoRT Import/Export System

## Quick Diagnosis

### Is the adapter running?
```bash
docker ps | grep adapter
# Should show: bamort-adapter-moam-dev (or your adapter name)
```

### Can you reach the adapter?
```bash
curl http://localhost:8181/metadata
# Should return JSON with adapter info
```

### Is the backend registered with the adapter?
```bash
docker logs bamort-backend-dev | grep "Registered adapter"
# Should show: Registered adapter: moam-vtt-v1
```

### Check recent imports
```sql
SELECT id, adapter_id, status, error_log, imported_at 
FROM import_histories 
ORDER BY imported_at DESC 
LIMIT 10;
```

## Common Issues

### 1. Adapter Not Detected

**Symptoms:**
- GET `/api/import/adapters` returns empty array or missing adapter
- Import fails with "no adapter found"

**Diagnosis:**
```bash
# Check adapter container
docker ps | grep adapter-moam

# Check adapter logs
docker logs bamort-adapter-moam-dev --tail=50

# Check backend logs
docker logs bamort-backend-dev | grep adapter

# Test adapter directly
curl http://localhost:8181/metadata
```

**Solutions:**

**A. Adapter container not running**
```bash
cd docker
docker-compose -f docker-compose.dev.yml up -d adapter-moam
```

**B. Adapter not in backend environment**
```yaml
# docker/docker-compose.dev.yml
bamort-backend-dev:
  environment:
    - IMPORT_ADAPTERS=[{"id":"moam-vtt-v1","base_url":"http://adapter-moam:8181"}]
```

Rebuild and restart:
```bash
docker-compose -f docker-compose.dev.yml down
docker-compose -f docker-compose.dev.yml up -d
```

**C. Network connectivity issue**
```bash
# Test from backend container
docker exec bamort-backend-dev wget -O- http://adapter-moam:8181/metadata
```

If fails, check Docker network:
```bash
docker network inspect bamort-network
# Both backend and adapter should be in same network
```

**D. Adapter marked unhealthy**
```bash
# Check health status
curl http://localhost:8180/api/import/adapters

# If "healthy": false, check adapter logs
docker logs bamort-adapter-moam-dev
```

Restart adapter:
```bash
docker restart bamort-adapter-moam-dev
```

---

### 2. Import Fails with Validation Error

**Symptoms:**
- HTTP 400 or 422 response
- Error message mentions validation

**Diagnosis:**
```bash
# Get detailed error from import history
curl http://localhost:8180/api/import/history/<import_id> \
  -H "Authorization: Bearer <token>"
```

Check database:
```sql
SELECT error_log FROM import_histories WHERE id = <import_id>;
```

**Solutions:**

**A. Missing required fields**
```
Error: "Character name is required"
```

Fix: Ensure source file has required fields:
- `name` (not empty)
- `game_system` (e.g., "Midgard5")

**B. Invalid stat values**
```
Error: "Stat St value -5 is out of range (0-100)"
```

This is a WARNING, not an error. Import should still succeed.
If blocked, check validator configuration.

**C. BMRT version mismatch**
```
Error: "BMRT version 2.0 not supported (expected: 1.0)"
```

Update adapter to support BMRT 1.0:
```go
// In adapter /metadata endpoint
"bmrt_versions": []string{"1.0"}
```

**D. JSON depth limit exceeded**
```
Error: "JSON depth exceeds maximum of 100 levels"
```

File has deeply nested JSON (possible attack or corrupted file).
Simplify structure or adjust limit in security.go (not recommended).

---

### 3. Character Created But Skills/Spells Missing

**Symptoms:**
- Import status: "success"
- Character exists but skills/spells/equipment missing

**Diagnosis:**
```sql
-- Check master data imports
SELECT item_type, external_name, match_type 
FROM master_data_imports 
WHERE import_history_id = <import_id>;

-- Check if skills were created as personal items
SELECT name, personal_item 
FROM skills 
WHERE created_by_user_id = <user_id>
ORDER BY created_at DESC;
```

**Solutions:**

**A. Reconciliation didn't create items**

Check reconciler logic in `reconciler.go`:
- Exact match failed (name mismatch)
- Personal item creation disabled (check code)

Debug:
```go
// Add logging in reconciler.go
log.Printf("Reconciling skill: %s (game: %s)", skill.Name, gameSystem)
```

**B. Master data created but not linked to character**

Check character creation logic in `import_logic.go`:
```go
// Verify skills are being linked
char.Skills = append(char.Skills, reconciledSkill)
```

**C. Game system mismatch**

Skills created for wrong game system:
```sql
SELECT name, game_system FROM skills WHERE created_by_user_id = <user_id>;
```

Fix: Ensure adapter sets correct `GameSystem` in BMRT:
```go
bmrt.GameSystem = "Midgard5"  // Must match GSM game system
```

---

### 4. Import Fails with "Rate Limit Exceeded"

**Symptoms:**
- HTTP 429 response
- Error: "Rate limit exceeded"

**Diagnosis:**
```bash
# Rate limits per user:
# - Detection: 10/minute
# - Import: 5/minute
# - Export: 20/minute
```

**Solutions:**

**A. Wait and retry**
```
Wait 60 seconds and retry request
```

**B. Adjust rate limits (development only)**

Edit `importer/routes.go`:
```go
// Increase limits for testing
detectLimiter := NewRateLimiter(100, time.Minute)
importLimiter := NewRateLimiter(50, time.Minute)
```

Rebuild backend:
```bash
docker-compose -f docker-compose.dev.yml restart bamort-backend-dev
```

---

### 5. Export Returns 409 Conflict

**Symptoms:**
- HTTP 409 response
- Error: "Original adapter unavailable"

**Diagnosis:**
```bash
# Check character's import history
curl http://localhost:8180/api/import/history \
  -H "Authorization: Bearer <token>"

# Check adapter health
curl http://localhost:8180/api/import/adapters
```

**Solutions:**

**A. Original adapter offline**
```bash
# Restart adapter
docker restart bamort-adapter-moam-dev

# Wait for health check (30 seconds)
sleep 30

# Retry export
```

**B. Use alternate adapter**
```bash
# Override adapter in export request
curl -X POST "http://localhost:8180/api/import/export/<char_id>?adapter_id=alternate-adapter-v1" \
  -H "Authorization: Bearer <token>"
```

**C. Character not imported (no original adapter)**

Character was created manually or via old importero system.
Cannot export without specifying adapter:

```bash
# Specify adapter explicitly
curl -X POST "http://localhost:8180/api/import/export/<char_id>?adapter_id=moam-vtt-v1" \
  -H "Authorization: Bearer <token>"
```

---

### 6. Detection Returns Low Confidence

**Symptoms:**
- Detection fails or returns wrong adapter
- Confidence score < 0.7

**Diagnosis:**
```bash
# Test detection directly on adapter
curl -X POST http://localhost:8181/detect \
  --data-binary @yourfile.json

# Check response
# {"confidence": 0.45, "version": ""}
```

**Solutions:**

**A. File format doesn't match adapter**

Try different adapter or create new adapter for this format.

**B. Adapter detection logic too strict**

Adjust confidence calculation in adapter:
```go
func calculateConfidence(char MoamCharacter) float64 {
    confidence := 0.0
    
    // Loosen requirements
    if char.Name != "" {
        confidence += 0.5  // Increased from 0.3
    }
    
    // ... adjust other checks
    
    return confidence
}
```

**C. Multiple adapters match (false positive)**

Check all adapters:
```bash
# List all registered adapters
curl http://localhost:8180/api/import/adapters

# Test detection with each adapter
for adapter in moam-vtt-v1 foundry-vtt-v1; do
    echo "Testing $adapter..."
    curl -X POST http://localhost:8181/detect --data-binary @file.json
done
```

Make detection more specific by adding unique signature checks.

**D. Manually specify adapter**
```bash
# Skip detection, specify adapter directly
curl -X POST "http://localhost:8180/api/import/import?adapter_id=moam-vtt-v1" \
  -H "Authorization: Bearer <token>" \
  -F "file=@yourfile.json"
```

---

### 7. Import Hangs or Times Out

**Symptoms:**
- Request never completes
- 504 Gateway Timeout after 30 seconds

**Diagnosis:**
```bash
# Check adapter logs for stuck processing
docker logs bamort-adapter-moam-dev --follow

# Check backend logs
docker logs bamort-backend-dev --follow

# Monitor database connections
docker exec bamort-mariadb-dev mysql -u bamort -pbamort -e "SHOW PROCESSLIST;"
```

**Solutions:**

**A. Large file processing**

Check file size:
```bash
ls -lh yourfile.json
```

If > 1MB, optimize adapter conversion logic.

**B. Database lock**

Transaction timeout or deadlock:
```sql
-- Check for locked tables
SHOW OPEN TABLES WHERE In_use > 0;

-- Kill stuck queries
SHOW PROCESSLIST;
KILL <process_id>;
```

**C. Infinite loop in adapter**

Add debug logging:
```go
func importHandler(c *gin.Context) {
    log.Println("Starting import...")
    
    // ... conversion logic with progress logs ...
    
    log.Println("Conversion complete")
}
```

Restart adapter:
```bash
docker restart bamort-adapter-moam-dev
```

**D. Network timeout**

Increase timeout in `registry.go`:
```go
httpClient := &http.Client{
    Timeout: 60 * time.Second,  // Increased from 30s
}
```

---

### 8. Compressed Data Corruption

**Symptoms:**
- Cannot decompress source_snapshot
- Error: "gzip: invalid header"

**Diagnosis:**
```sql
-- Check snapshot size
SELECT id, LENGTH(source_snapshot) as size 
FROM import_histories 
WHERE id = <import_id>;
```

```go
// Test decompression
data, err := decompressData(snapshot)
if err != nil {
    log.Printf("Decompression failed: %v", err)
}
```

**Solutions:**

**A. Data not compressed**

Old imports may have uncompressed data:
```go
// Try direct parse first
var bmrt CharacterImport
if err := json.Unmarshal(snapshot, &bmrt); err == nil {
    // Already uncompressed
    return snapshot, nil
}

// Then try decompression
return decompressData(snapshot)
```

**B. Partial write**

Transaction rollback didn't complete:
```sql
-- Check import status
SELECT status, error_log FROM import_histories WHERE id = <import_id>;

-- Delete corrupted import
DELETE FROM import_histories WHERE id = <import_id>;
```

---

### 9. Performance Issues

**Symptoms:**
- Import takes > 5 seconds
- Detection takes > 2 seconds

**Diagnosis:**
```bash
# Run performance tests
cd backend
go test -bench=BenchmarkImportCharacter ./importer/

# Profile import
go test -bench=BenchmarkImportCharacter -cpuprofile=cpu.prof ./importer/
go tool pprof cpu.prof
```

**Solutions:**

**A. Too many skills/spells**

Batch reconciliation:
```go
// Instead of individual reconciliation
for _, skill := range skills {
    ReconcileSkill(skill, importID, gameSystem)
}

// Use bulk insert
db.CreateInBatches(reconciledSkills, 100)
```

**B. No detection cache**

Verify cache is enabled:
```go
// In detector.go
cache := &DetectionCache{
    entries: make(map[string]*CacheEntry),
    ttl:     1 * time.Hour,
}
```

**C. Adapter too slow**

Profile adapter with sample file:
```bash
time curl -X POST http://localhost:8181/import --data-binary @sample.json
```

Optimize conversion logic.

**D. Database indexes missing**

Create indexes:
```sql
CREATE INDEX idx_import_user ON import_histories(user_id);
CREATE INDEX idx_import_char ON import_histories(character_id);
CREATE INDEX idx_masterdata_import ON master_data_imports(import_history_id);
```

---

## Debugging Tools

### Enable Debug Logging

Backend:
```yaml
# docker/docker-compose.dev.yml
bamort-backend-dev:
  environment:
    - LOG_LEVEL=debug
```

Adapter:
```go
// In adapter main.go
gin.SetMode(gin.DebugMode)
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### Database Inspection

```sql
-- Recent imports by user
SELECT u.username, ih.adapter_id, ih.status, ih.imported_at
FROM import_histories ih
JOIN users u ON ih.user_id = u.id
ORDER BY ih.imported_at DESC
LIMIT 20;

-- Failed imports
SELECT id, adapter_id, source_filename, error_log
FROM import_histories
WHERE status = 'failed'
ORDER BY imported_at DESC;

-- Personal items created
SELECT item_type, COUNT(*) as count
FROM master_data_imports
WHERE match_type = 'created_personal'
GROUP BY item_type;

-- Top imported characters
SELECT c.name, c.user_id, ih.adapter_id, ih.imported_at
FROM chars c
JOIN import_histories ih ON c.id = ih.character_id
ORDER BY ih.imported_at DESC
LIMIT 20;
```

### HTTP Request Tracing

```bash
# Verbose curl
curl -v -X POST http://localhost:8180/api/import/import \
  -H "Authorization: Bearer <token>" \
  -F "file=@test.json"

# With timing
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8180/api/import/adapters

# curl-format.txt content:
#     time_namelookup:  %{time_namelookup}\n
#        time_connect:  %{time_connect}\n
#     time_appconnect:  %{time_appconnect}\n
#    time_pretransfer:  %{time_pretransfer}\n
#       time_redirect:  %{time_redirect}\n
#  time_starttransfer:  %{time_starttransfer}\n
#                     ----------\n
#          time_total:  %{time_total}\n
```

### Container Health

```bash
# Check all containers
docker ps -a | grep bamort

# Inspect container
docker inspect bamort-adapter-moam-dev

# Check resource usage
docker stats bamort-adapter-moam-dev

# View full logs
docker logs bamort-adapter-moam-dev --since=1h > adapter-logs.txt
```

## Getting Help

### 1. Gather Information

Before reporting an issue, collect:
- Backend logs: `docker logs bamort-backend-dev`
- Adapter logs: `docker logs bamort-adapter-<name>-dev`
- Sample file (if possible)
- Import history ID and error log from database
- Steps to reproduce

### 2. Check Existing Issues

Search GitHub issues: https://github.com/Bardioc26/bamort/issues

### 3. Create Issue

Include:
- BaMoRT version
- Adapter ID and version
- Error messages
- Relevant logs
- Sample file (if not sensitive)

## Maintenance

### Clean Up Old Imports

```sql
-- Delete imports older than 90 days (keeps character)
DELETE FROM import_histories 
WHERE imported_at < DATE_SUB(NOW(), INTERVAL 90 DAY);

-- Rebuild indexes
OPTIMIZE TABLE import_histories;
OPTIMIZE TABLE master_data_imports;
```

### Monitor Adapter Health

```bash
# Run health check manually
curl http://localhost:8180/api/import/adapters | jq '.adapters[] | {id: .id, healthy: .healthy, last_checked: .last_checked_at}'
```

### Backup Import Data

```bash
# Export import history
docker exec bamort-mariadb-dev mysqldump \
  -u bamort -pbamort bamort \
  import_histories master_data_imports \
  > import_backup_$(date +%Y%m%d).sql
```

## Performance Monitoring

### Metrics to Track

1. **Import Success Rate**: `(successful_imports / total_imports) * 100`
2. **Average Import Time**: Track with performance tests
3. **Adapter Availability**: Track health check failures
4. **Rate Limit Hits**: Log 429 responses

### Query Examples

```sql
-- Import success rate by adapter
SELECT 
    adapter_id,
    COUNT(*) as total,
    SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as successful,
    ROUND(SUM(CASE WHEN status = 'success' THEN 1 ELSE 0 END) / COUNT(*) * 100, 2) as success_rate_pct
FROM import_histories
GROUP BY adapter_id;

-- Average imports per user
SELECT 
    AVG(import_count) as avg_imports_per_user
FROM (
    SELECT user_id, COUNT(*) as import_count
    FROM import_histories
    GROUP BY user_id
) user_imports;

-- Recent error patterns
SELECT 
    LEFT(error_log, 100) as error_prefix,
    COUNT(*) as occurrences
FROM import_histories
WHERE status = 'failed'
AND imported_at > DATE_SUB(NOW(), INTERVAL 7 DAY)
GROUP BY error_prefix
ORDER BY occurrences DESC
LIMIT 10;
```
