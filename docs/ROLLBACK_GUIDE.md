# Rollback Guide

**Version:** 1.0  
**Last Updated:** 16. Januar 2026  
**Emergency Contact:** System Administrator

---

## Table of Contents
1. [When to Rollback](#when-to-rollback)
2. [Rollback Options](#rollback-options)
3. [Rollback Procedures](#rollback-procedures)
4. [Emergency Rollback](#emergency-rollback)
5. [Time Estimates](#time-estimates)
6. [Post-Rollback Verification](#post-rollback-verification)

---

## When to Rollback

### Immediate Rollback Triggers

Roll back IMMEDIATELY if:

- ❌ **Critical Migration Failure** - Migration cannot complete or database is corrupted
- ❌ **Data Loss Detected** - Characters, skills, or other critical data missing
- ❌ **System Unusable** - Backend crashes repeatedly or won't start
- ❌ **Security Vulnerability** - New deployment introduces security issue

### Consider Rollback If:

- ⚠️ **Non-Critical Errors** - Minor features broken but system functional
- ⚠️ **Performance Degradation** - System noticeably slower after update
- ⚠️ **User Reports** - Multiple users reporting same issue

### DO NOT Rollback For:

- ✓ **Minor UI Bugs** - CSS issues, minor display problems
- ✓ **Non-Blocking Errors** - Errors that don't affect core functionality
- ✓ **Expected Warnings** - Warnings documented in release notes

---

## Rollback Options

Bamort provides three rollback methods:

| Method | Speed | Scope | Data Loss | When to Use |
|--------|-------|-------|-----------|-------------|
| **Migration Rollback** | Fast (1-5 min) | Database only | None | Migration failed but data intact |
| **JSON Restore** | Medium (5-15 min) | All data | Changes since backup | Complete rollback needed |
| **Full System Rollback** | Slow (10-30 min) | Everything | Changes since backup | Catastrophic failure |

---

## Rollback Procedures

### Option 1: Migration Rollback (Preferred)

**Use When:** Recent migrations caused issues but data is intact

**Prerequisites:**
- Migrations have `DownSQL` defined (rollback scripts)
- No data corruption

**Steps:**

#### 1. Check Current State

```bash
# View migration history
docker exec bamort-backend /app/deploy migrations history
```

**Output:**
```
Migration History:
  #5: Create equipment_cache table (applied 5 min ago) ✓
  #4: Add learning_category column (applied 7 min ago) ✓
  #3: Update skill_indices (applied 2 days ago) ✓
  #2: Add user_preferences (applied 5 days ago) ✓
  #1: Create version_tables (applied 30 days ago) ✓
```

#### 2. Identify Problem Migration

Determine which migration(s) to rollback:
- If last migration failed: Rollback 1 step
- If system broken after multiple migrations: Rollback to last known good state

#### 3. Execute Rollback

```bash
# Rollback last migration
docker exec bamort-backend /app/deploy migrations rollback --steps 1

# Or rollback to specific version
docker exec bamort-backend /app/deploy migrations rollback --to-version 0.4.0
```

**Expected Output:**
```
Rolling back 1 migration(s)...

Rolling back migration #5: Create equipment_cache table
Executing: DROP TABLE IF EXISTS equipment_cache;
✓ Migration #5 rolled back (executed in 120ms)

Rollback completed successfully.
Database version: 0.5.0 → 0.4.0
```

#### 4. Verify Rollback

```bash
# Check version
docker exec bamort-backend /app/deploy status

# Check system health
curl http://localhost:8180/api/system/health | jq
```

#### 5. Restart Services

```bash
# Restart backend to clear caches
docker-compose -f docker/docker-compose.yml restart backend

# Test functionality
curl http://localhost:8180/api/system/health
```

**Time Estimate:** 1-5 minutes

---

### Option 2: JSON Restore

**Use When:** Need to restore data to pre-deployment state

**Prerequisites:**
- Backup created before deployment
- Backup file accessible

**Steps:**

#### 1. List Available Backups

```bash
docker exec bamort-backend /app/deploy backup list
```

**Output:**
```
Available Backups:
  backup_20260116_220000_v0.4.0_m3.json (2.4 MB, 10 minutes ago)
  backup_20260115_180000_v0.4.0_m3.json (2.3 MB, 1 day ago)
  backup_20260114_120000_v0.4.0_m3.json (2.2 MB, 2 days ago)
```

#### 2. Stop Backend (Prevent Data Changes)

```bash
docker-compose -f docker/docker-compose.yml stop backend
```

#### 3. Restore Backup

```bash
docker exec bamort-mariadb /app/deploy backup restore \
  --file /app/backups/backup_20260116_220000_v0.4.0_m3.json
```

**Expected Output:**
```
Restoring from backup: backup_20260116_220000_v0.4.0_m3.json
Backup version: 0.4.0
Backup date: 2026-01-16 22:00:00

WARNING: This will DELETE all current data!
Type 'CONFIRM' to proceed: CONFIRM

Dropping existing tables...
✓ Tables dropped

Restoring data...
✓ Users: 5 records restored
✓ Characters: 23 records restored
✓ Skills: 245 records restored
✓ Spells: 189 records restored
✓ Equipment: 156 records restored

Restore completed successfully.
Database restored to version: 0.4.0
Total records restored: 618
```

#### 4. Restart Backend

```bash
# Ensure backend version matches restored DB version
# If needed, rollback Docker image to previous version

docker-compose -f docker/docker-compose.yml start backend

# Verify
docker logs bamort-backend --tail=50
```

#### 5. Verify Restore

```bash
# Check version compatibility
curl http://localhost:8180/api/system/health | jq

# Check data
# Login to frontend and verify characters exist
```

**Time Estimate:** 5-15 minutes (depends on data size)

---

### Option 3: Full System Rollback

**Use When:** Complete system failure, nothing works

**Prerequisites:**
- Access to Docker host
- Previous Docker images available
- Backup available

**Steps:**

#### 1. Stop All Services

```bash
cd /data/dev/bamort/docker
./stop-prd.sh
```

#### 2. Backup Current State (If Possible)

```bash
# Create emergency backup of current state
docker-compose -f docker-compose.yml start mariadb
sleep 5
docker exec bamort-backend /app/deploy backup create --emergency
docker-compose -f docker-compose.yml stop mariadb
```

#### 3. Restore Database Volume

```bash
# Option A: Restore from volume backup
docker run --rm \
  -v bamort-db:/data \
  -v $(pwd)/backups:/backup \
  alpine sh -c "cd /data && tar -xzf /backup/mariadb_backup_20260116.tar.gz"

# Option B: Recreate volume and import JSON backup
docker volume rm bamort-db
docker volume create bamort-db
# Then start mariadb and import JSON (see Option 2)
```

#### 4. Rollback Docker Images

```bash
# Check available images
docker images | grep bamort

# Tag previous version as latest (if needed)
docker tag bamort-backend:0.4.0 bamort-backend:latest
docker tag bamort-frontend:0.4.0 bamort-frontend:latest
```

#### 5. Start Services

```bash
./start-prd.sh

# Monitor startup
docker-compose -f docker-compose.yml logs --follow
```

#### 6. Verify System

```bash
# Check all containers running
docker ps | grep bamort

# Check health
curl http://localhost:8180/api/system/health | jq

# Check frontend
open http://localhost:5173
```

**Time Estimate:** 10-30 minutes (depends on backup size)

---

## Emergency Rollback

### 🔴 Emergency Procedure (System Down)

If system is completely broken and users cannot access:

#### Quick Rollback (5 minutes):

```bash
# 1. Stop everything
cd /data/dev/bamort/docker
./stop-prd.sh

# 2. Revert to last known good version
git checkout <previous-version-tag>  # e.g., v0.4.0

# 3. Restore database backup
docker-compose -f docker-compose.yml start mariadb
sleep 10
docker exec bamort-mariadb mysql -u root -p<password> -e "DROP DATABASE bamort; CREATE DATABASE bamort;"
docker exec -i bamort-mariadb mysql -u bamort -p<password> bamort < backups/latest.sql

# 4. Restart services
./start-prd.sh

# 5. Verify
curl http://localhost:8180/api/system/health
open http://localhost:5173
```

#### Rollback Decision Tree

```
Deployment failed
      ↓
   Can backend start?
   /              \
 YES              NO
  ↓                ↓
Migrations     Check Docker
 failed?        logs
  ↓                ↓
YES              Fix and
  ↓             restart
Rollback          ↓
Migration      Still
  ↓            failing?
Test             ↓
  ↓            YES
Working?         ↓
  ↓         Full System
 NO         Rollback
  ↓
JSON
Restore
```

---

## Time Estimates

### Rollback Time by Method

| Rollback Method | Minimum | Typical | Maximum |
|-----------------|---------|---------|---------|
| **Migration Rollback (1-3 migrations)** | 30 sec | 2 min | 5 min |
| **JSON Restore (Small DB < 10MB)** | 2 min | 5 min | 10 min |
| **JSON Restore (Large DB > 100MB)** | 5 min | 15 min | 30 min |
| **Full System Rollback** | 10 min | 20 min | 45 min |
| **Emergency Quick Rollback** | 3 min | 5 min | 10 min |

### Rollback Risk by Complexity

| Complexity | Risk | Recovery Time if Fails |
|------------|------|------------------------|
| **Single Migration** | Low | 5 minutes (re-apply) |
| **Multiple Migrations** | Medium | 15 minutes (JSON restore) |
| **Full System** | High | 30-60 minutes (rebuild) |

---

## Post-Rollback Verification

After any rollback, perform these checks:

### 1. System Health Check

```bash
# Check health endpoint
curl http://localhost:8180/api/system/health | jq

# Expected output:
{
  "status": "ok",
  "backend_version": "0.4.0",
  "required_db_version": "0.4.0",
  "actual_db_version": "0.4.0",
  "migrations_pending": false,
  "compatible": true
}
```

### 2. Version Verification

```bash
# Check versions match
docker exec bamort-backend /app/deploy status

# Expected output:
Backend Version:      0.4.0
Database Version:     0.4.0
Migration Number:     3
Migrations Pending:   0
Compatible:           Yes
```

### 3. Data Integrity Check

```bash
# Check record counts
docker exec bamort-mariadb mysql -u bamort -p<password> bamort -e "
  SELECT 'Users' as table_name, COUNT(*) as count FROM users
  UNION SELECT 'Characters', COUNT(*) FROM char_chars
  UNION SELECT 'Skills', COUNT(*) FROM gsm_skills
  UNION SELECT 'Spells', COUNT(*) FROM gsm_spells;
"
```

### 4. Functional Testing

Manual tests:

- [ ] Can login with existing user
- [ ] Can view character list
- [ ] Can view character details
- [ ] Can edit character
- [ ] Can create new skill/spell
- [ ] Can export character PDF
- [ ] No error messages in console

### 5. Log Review

```bash
# Check for errors in backend
docker logs bamort-backend --tail=100 | grep ERROR

# Check for errors in frontend
docker logs bamort-frontend --tail=100 | grep ERROR

# Should be no critical errors
```

---

## Common Rollback Scenarios

### Scenario 1: Last Migration Failed

**Situation:** Applied 3 migrations, 3rd one failed

**Solution:**
```bash
# Rollback the failed migration
docker exec bamort-backend /app/deploy migrations rollback --steps 1

# Verify system works
curl http://localhost:8180/api/system/health

# Fix migration script
# Re-apply when fixed
```

**Time:** 2-3 minutes

---

### Scenario 2: System Broken After Multiple Migrations

**Situation:** Applied 5 migrations, system now broken

**Solution:**
```bash
# Rollback to last known good version
docker exec bamort-backend /app/deploy migrations rollback --to-version 0.4.0

# If rollback fails, use JSON restore
docker exec bamort-backend /app/deploy backup restore \
  --file /app/backups/backup_<timestamp>_v0.4.0.json
```

**Time:** 5-15 minutes

---

### Scenario 3: Data Corruption Detected

**Situation:** Characters missing or corrupted after deployment

**Solution:**
```bash
# IMMEDIATELY stop backend
docker-compose -f docker/docker-compose.yml stop backend

# Restore from backup
docker exec bamort-mariadb /app/deploy backup restore \
  --file /app/backups/backup_<timestamp>_v0.4.0.json

# Rollback backend version
git checkout v0.4.0
docker-compose -f docker/docker-compose.yml up -d
```

**Time:** 10-20 minutes

---

### Scenario 4: Frontend Broken, Backend OK

**Situation:** Backend works but frontend has critical bug

**Solution:**
```bash
# Rollback frontend only
docker-compose -f docker/docker-compose.yml stop frontend
docker tag bamort-frontend:0.4.0 bamort-frontend:latest
docker-compose -f docker/docker-compose.yml start frontend
```

**Time:** 1-2 minutes

---

## Post-Rollback Actions

After successful rollback:

1. **Document the Issue**
   - What went wrong?
   - What was the error message?
   - Which migration/code caused it?

2. **Notify Users**
   - "System restored to previous version"
   - Explain any data changes
   - Estimated time for fix

3. **Fix the Problem**
   - Fix migration script
   - Fix code bug
   - Test on development environment

4. **Plan Re-Deployment**
   - Schedule new deployment window
   - Test thoroughly before re-deploying
   - Prepare better rollback plan

5. **Review Backup Strategy**
   - Ensure backups are working
   - Verify backup retention policy
   - Test restore procedure

---

## Rollback Checklist

Print this checklist for emergency use:

```
□ System failure confirmed
□ Backup identified (version: _______)
□ Stop affected services
□ Create emergency backup (if possible)
□ Execute rollback procedure
□ Verify system health
□ Verify version compatibility
□ Test core functionality
□ Check for data loss
□ Restart all services
□ Monitor for 30 minutes
□ Document incident
□ Notify users
□ Plan fix
```

---

## Getting Help

If rollback fails or you're unsure:

1. **Check Documentation**
   - [DEPLOYMENT_RUNBOOK.md](DEPLOYMENT_RUNBOOK.md)
   - [TROUBLESHOOTING.md](TROUBLESHOOTING.md)

2. **Check Logs**
   ```bash
   docker logs bamort-backend --tail=200 > backend_error.log
   docker logs bamort-mariadb --tail=200 > db_error.log
   ```

3. **Contact Support**
   - Email: admin@bamort.local
   - Emergency: [Phone Number]

4. **Community**
   - GitHub Issues: https://github.com/Bardioc26/bamort/issues
   - Documentation: https://github.com/Bardioc26/bamort/docs

---

**Remember:** A successful rollback is better than a broken system. When in doubt, rollback!

**Last Updated:** 16. Januar 2026  
**Version:** 1.0
