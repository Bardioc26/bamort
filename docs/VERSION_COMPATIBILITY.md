# Version Compatibility Reference

**Version:** 1.0  
**Last Updated:** 16. Januar 2026

---

## Overview

This document defines the compatibility requirements between Bamort backend and database versions.

### Version Strategy

**Rule:** Backend version and database version must match exactly.

Each backend version defines exactly which database version it requires via the `RequiredDBVersion` constant in `backend/deployment/version/version.go`.

---

## Version Matrix

| Backend Version | Required DB Version | Migration Count | Release Date | Status |
|-----------------|---------------------|-----------------|--------------|--------|
| 0.5.0 | 0.5.0 | 5 | TBD | Planned |
| 0.4.0 | 0.4.0 | 3 | 2026-01-15 | Current |
| 0.3.0 | 0.3.0 | 2 | 2025-12-20 | Deprecated |
| 0.2.0 | 0.2.0 | 1 | 2025-12-01 | Deprecated |
| 0.1.x | 0.1.0 | 0 | 2025-11-15 | Legacy |

---

## Compatibility Rules

### ✅ Compatible Combinations

```
Backend 0.4.0 + Database 0.4.0 = ✅ Compatible
Backend 0.5.0 + Database 0.5.0 = ✅ Compatible
```

### ❌ Incompatible Combinations

```
Backend 0.5.0 + Database 0.4.0 = ❌ Database too old (migration needed)
Backend 0.4.0 + Database 0.5.0 = ❌ Database too new (backend too old)
Backend 0.5.0 + Database 0.3.0 = ❌ Cannot skip versions
```

---

## Migration Paths

### Sequential Upgrade Required

You cannot skip versions. Must upgrade sequentially:

```
0.1.0 → 0.2.0 → 0.3.0 → 0.4.0 → 0.5.0
```

**Example:** To upgrade from 0.2.0 to 0.5.0:

```bash
# Step 1: Upgrade to 0.3.0
docker exec bamort-backend /app/deploy migrations apply --to-version 0.3.0

# Step 2: Upgrade to 0.4.0
docker exec bamort-backend /app/deploy migrations apply --to-version 0.4.0

# Step 3: Upgrade to 0.5.0
docker exec bamort-backend /app/deploy migrations apply --to-version 0.5.0
```

### Direct Upgrade (Only One Version Apart)

You can directly upgrade one version:

```bash
# From 0.4.0 to 0.5.0 (OK - one version)
docker exec bamort-backend /app/deploy migrations apply --all
```

---

## Version Details

### Version 0.5.0 (Planned)

**Required DB Version:** 0.5.0  
**Migrations:** 5 total  
**New Migrations:** 4, 5

**Database Changes:**
- Migration #4: Add `learning_category` column to spells
- Migration #5: Create `equipment_cache` table for performance

**Breaking Changes:**
- None

**Upgrade Path:**
- From 0.4.0: Direct upgrade (apply migrations 4-5)
- From 0.3.0: Upgrade to 0.4.0 first

---

### Version 0.4.0 (Current)

**Required DB Version:** 0.4.0  
**Migrations:** 3 total  
**New Migrations:** 3

**Database Changes:**
- Migration #1: Create `schema_version` and `migration_history` tables
- Migration #2: Add `user_preferences` table
- Migration #3: Update skill indices

**Breaking Changes:**
- None

**Upgrade Path:**
- From 0.3.0: Direct upgrade (apply migrations 1-3)
- From 0.2.0: Upgrade to 0.3.0 first

---

### Version 0.3.0 (Deprecated)

**Required DB Version:** 0.3.0  
**Migrations:** 2 total

**Status:** Deprecated - upgrade to 0.4.0 recommended

**Upgrade Path:**
- Must upgrade to 0.4.0

---

### Version 0.2.0 (Deprecated)

**Required DB Version:** 0.2.0  
**Migrations:** 1 total

**Status:** Deprecated - upgrade to 0.4.0 required

**Upgrade Path:**
- Upgrade to 0.3.0, then to 0.4.0

---

### Version 0.1.x (Legacy)

**Required DB Version:** 0.1.0  
**Migrations:** 0 (pre-migration system)

**Status:** Legacy - no longer supported

**Upgrade Path:**
- Must perform fresh installation or manual migration

---

## Checking Compatibility

### Command Line

```bash
# Check current versions
docker exec bamort-backend /app/deploy status
```

**Output:**
```
Backend Version:      0.4.0
Required DB Version:  0.4.0
Actual DB Version:    0.4.0
Migration Number:     3
Migrations Pending:   0
Compatible:           Yes
```

### API Endpoint

```bash
curl http://localhost:8180/api/system/health | jq
```

**Output:**
```json
{
  "status": "ok",
  "backend_version": "0.4.0",
  "required_db_version": "0.4.0",
  "actual_db_version": "0.4.0",
  "migrations_pending": false,
  "compatible": true
}
```

### Frontend Warning Banner

When `compatible: false`, frontend shows:

```
⚠️ Database migration required. Please contact administrator.
Backend: v0.5.0 | Database: v0.4.0
```

---

## Upgrade Procedures

### Before Upgrading

1. **Check Current Version**
   ```bash
   docker exec bamort-backend /app/deploy status
   ```

2. **Create Backup**
   ```bash
   docker exec bamort-backend /app/deploy backup create
   ```

3. **Review Release Notes**
   - Check `CHANGELOG.md`
   - Review migration scripts
   - Note breaking changes

### Standard Upgrade (One Version)

```bash
# 1. Stop frontend
docker-compose -f docker/docker-compose.yml stop frontend

# 2. Pull new backend
docker-compose -f docker/docker-compose.yml pull backend

# 3. Start backend (auto-runs migrations)
docker-compose -f docker/docker-compose.yml up -d backend

# 4. Verify
docker exec bamort-backend /app/deploy status

# 5. Start frontend
docker-compose -f docker/docker-compose.yml start frontend
```

### Multi-Version Upgrade

```bash
# Example: 0.2.0 → 0.4.0

# Step 1: Upgrade to 0.3.0
docker pull bamort/backend:0.3.0
docker tag bamort/backend:0.3.0 bamort-backend:latest
docker-compose -f docker/docker-compose.yml up -d backend
docker exec bamort-backend /app/deploy status

# Step 2: Upgrade to 0.4.0
docker pull bamort/backend:0.4.0
docker tag bamort/backend:0.4.0 bamort-backend:latest
docker-compose -f docker/docker-compose.yml up -d backend
docker exec bamort-backend /app/deploy status
```

---

## Downgrade Procedures

### One Version Rollback

```bash
# Example: 0.5.0 → 0.4.0

# 1. Rollback migrations
docker exec bamort-backend /app/deploy migrations rollback --to-version 0.4.0

# 2. Downgrade backend image
docker tag bamort-backend:0.4.0 bamort-backend:latest
docker-compose -f docker/docker-compose.yml restart backend

# 3. Verify
docker exec bamort-backend /app/deploy status
```

### Multi-Version Rollback

Not recommended - use backup restore instead:

```bash
docker exec bamort-backend /app/deploy backup restore \
  --file /app/backups/backup_<timestamp>_v0.4.0.json
```

---

## Migration Details by Version

### Migrations in 0.5.0

**Migration #4: Add learning_category**
- **Purpose:** Separate spell categories for learning costs
- **Rollback:** Safe - removes column
- **Time:** < 1 second

**Migration #5: Create equipment_cache**
- **Purpose:** Performance optimization for equipment queries
- **Rollback:** Safe - drops table
- **Time:** < 1 second

### Migrations in 0.4.0

**Migration #1: Create version tables**
- **Purpose:** Initialize version tracking system
- **Rollback:** NOT SAFE - removes version tracking
- **Time:** < 1 second

**Migration #2: Add user_preferences**
- **Purpose:** Store user UI preferences
- **Rollback:** Safe - drops table
- **Time:** < 1 second

**Migration #3: Update skill indices**
- **Purpose:** Performance improvement for skill queries
- **Rollback:** Safe - drops indices
- **Time:** 1-3 seconds

---

## FAQ

### Q: Can I skip versions during upgrade?

**A:** No. You must upgrade sequentially: 0.2.0 → 0.3.0 → 0.4.0 → 0.5.0

### Q: Can I run newer backend with older database?

**A:** No. Backend will refuse to start if database version doesn't match required version.

### Q: Can I run older backend with newer database?

**A:** No. This is not supported and will cause errors.

### Q: How do I check if migration is needed?

**A:** Run `docker exec bamort-backend /app/deploy status` or check `/api/system/health`

### Q: What if migration fails halfway?

**A:** Migrations run in transactions. If migration fails, it rolls back automatically. Use `migrations rollback` to revert completed migrations.

### Q: Can I manually change database version?

**A:** Not recommended. Use the deployment tools to ensure consistency.

### Q: How long do migrations take?

**A:** Most migrations complete in < 5 seconds. Large data migrations may take minutes.

### Q: Do I need downtime for migrations?

**A:** Yes. Stop frontend during migration to prevent user access.

---

## Version Update Checklist

When releasing new version:

```
□ Update RequiredDBVersion constant in version.go
□ Create migration scripts (UpSQL and DownSQL)
□ Add migration to AllMigrations slice
□ Test migration on development database
□ Test rollback on development database
□ Update this VERSION_COMPATIBILITY.md
□ Update CHANGELOG.md
□ Tag release in git
□ Build and tag Docker images
□ Test full upgrade path
□ Document breaking changes
```

---

## Support Matrix

| Version | Status | Support | Updates |
|---------|--------|---------|---------|
| 0.5.0 | Planned | TBD | TBD |
| 0.4.0 | Current | Full | Security + Features |
| 0.3.0 | Deprecated | Security only | Critical fixes only |
| 0.2.0 | Deprecated | None | Upgrade required |
| 0.1.x | Legacy | None | Not supported |

---

**Recommendation:** Always run latest stable version (currently 0.4.0)

**Last Updated:** 16. Januar 2026  
**Version:** 1.0
