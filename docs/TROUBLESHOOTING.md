# Deployment Troubleshooting Guide

**Version:** 1.0  
**Last Updated:** 16. Januar 2026

---

## Quick Diagnosis

```bash
# Run full system diagnosis
docker exec bamort-backend /app/deploy diagnose

# Check system health
curl http://localhost:8180/api/system/health | jq

# View recent logs
docker logs bamort-backend --tail=100
docker logs bamort-mariadb --tail=100
```

---

## Common Issues

### Issue 1: Migration Fails with SQL Error

**Symptoms:**
- Migration command fails
- Error message contains SQL syntax error
- Database left in inconsistent state

**Error Example:**
```
Error applying migration #5: SQL syntax error
near "CREAT TABLE": syntax error
```

**Diagnosis:**
```bash
# Check which migration failed
docker exec bamort-backend /app/deploy migrations history

# Review migration SQL
cat backend/deployment/migrations/all_migrations.go | grep -A 20 "Migration{Number: 5"
```

**Solutions:**

1. **Rollback Failed Migration**
   ```bash
   docker exec bamort-backend /app/deploy migrations rollback --steps 1
   ```

2. **Fix SQL and Re-apply**
   - Fix the SQL in migration file
   - Rebuild backend
   - Re-apply migration

3. **Manual SQL Fix** (if table half-created)
   ```bash
   docker exec -it bamort-mariadb mysql -u bamort -p bamort
   # Manually DROP or FIX table
   # Then rollback migration number in schema_version
   ```

---

### Issue 2: Version Mismatch Error

**Symptoms:**
- Frontend shows yellow warning banner
- Backend logs show version incompatibility
- Some features not working

**Error Example:**
```
Backend version 0.5.0 requires database 0.5.0, but found 0.4.0
```

**Diagnosis:**
```bash
docker exec bamort-backend /app/deploy status
```

**Output:**
```
Backend Version:      0.5.0
Required DB Version:  0.5.0
Actual DB Version:    0.4.0    ← MISMATCH
Migrations Pending:   2
Compatible:           No
```

**Solutions:**

1. **Apply Pending Migrations**
   ```bash
   docker exec bamort-backend /app/deploy migrations pending
   docker exec bamort-backend /app/deploy migrations apply --all
   ```

2. **If Migrations Fail** → Rollback backend to 0.4.0
   ```bash
   docker tag bamort-backend:0.4.0 bamort-backend:latest
   docker-compose -f docker/docker-compose.yml restart backend
   ```

---

### Issue 3: Backend Won't Start

**Symptoms:**
- `docker ps` shows backend constantly restarting
- Cannot access http://localhost:8180
- Frontend cannot connect to API

**Diagnosis:**
```bash
# Check container status
docker ps -a | grep backend

# View crash logs
docker logs bamort-backend --tail=100

# Common error patterns:
# - "connection refused" → Database not ready
# - "migration failed" → Database schema broken
# - "port already in use" → Port conflict
```

**Solutions:**

**A. Database Connection Failed**
```bash
# Check mariadb running
docker ps | grep mariadb

# Check database credentials
docker exec bamort-backend env | grep DATABASE

# Test connection manually
docker exec bamort-mariadb mysql -u bamort -p<password> -e "SELECT 1"
```

**B. Migration on Startup Failed**
```bash
# Disable auto-migration temporarily
docker exec bamort-backend /app/deploy migrations rollback --steps 1

# Restart backend
docker-compose -f docker/docker-compose.yml restart backend
```

**C. Port Conflict**
```bash
# Check what's using port 8180
lsof -i :8180

# Kill conflicting process or change port
```

---

### Issue 4: Master Data Import Fails

**Symptoms:**
- Import command fails with file not found
- Import succeeds but skills/spells missing
- Duplicate key errors during import

**Error Examples:**
```
Error: failed to open file masterdata/skills.json: no such file or directory
Error: duplicate entry 'Heimlichkeit' for key 'name'
```

**Diagnosis:**
```bash
# Check masterdata directory
docker exec bamort-backend ls -la /app/masterdata/

# Check import logs
docker logs bamort-backend | grep "Importing"
```

**Solutions:**

**A. Files Missing**
```bash
# Copy masterdata to container
docker cp ./masterdata bamort-backend:/app/masterdata

# Verify files
docker exec bamort-backend ls /app/masterdata/
```

**B. Duplicate Keys**
```bash
# Use --force flag to overwrite
docker exec bamort-backend /app/deploy masterdata import \
  --source /app/masterdata \
  --force

# Or clean database first
docker exec bamort-mariadb mysql -u bamort -p<password> bamort \
  -e "DELETE FROM gsm_skills"
```

**C. JSON Parse Errors**
```bash
# Validate JSON files
docker exec bamort-backend /app/deploy masterdata validate \
  --source /app/masterdata

# Check file encoding (should be UTF-8)
file ./masterdata/skills.json
```

---

### Issue 5: Frontend Shows 404 for API Calls

**Symptoms:**
- Frontend loads but shows errors
- Browser console shows "404 Not Found" for /api/* calls
- Login fails with network error

**Diagnosis:**
```bash
# Check backend responding
curl http://localhost:8180/api/system/health

# Check frontend API configuration
docker exec bamort-frontend cat /app/dist/.env
```

**Solutions:**

**A. Backend Not Running**
```bash
docker-compose -f docker/docker-compose.yml start backend
```

**B. Wrong API URL in Frontend**
```bash
# Check VITE_API_URL environment variable
docker-compose -f docker/docker-compose.yml restart frontend
```

**C. CORS Issues**
```bash
# Check browser console for CORS errors
# Verify frontend origin in backend CORS config
docker exec bamort-backend env | grep CORS_ORIGINS
```

---

### Issue 6: Backup Creation Fails

**Symptoms:**
- Backup command fails
- Disk full error
- Backup file empty or corrupted

**Error Examples:**
```
Error: no space left on device
Error: backup file size is 0 bytes
```

**Diagnosis:**
```bash
# Check disk space
df -h

# Check backup directory
docker exec bamort-backend ls -lh /app/backups/

# Check backup permissions
docker exec bamort-backend ls -ld /app/backups
```

**Solutions:**

**A. Disk Full**
```bash
# Clean old backups
docker exec bamort-backend /app/deploy backup cleanup --keep 5

# Or manually delete old backups
docker exec bamort-backend rm /app/backups/backup_*.json
```

**B. Permission Denied**
```bash
# Fix permissions
docker exec bamort-backend chmod 777 /app/backups
```

**C. Database Export Fails**
```bash
# Check database connection
docker exec bamort-mariadb mysqldump --help

# Try manual export
docker exec bamort-mariadb mysqldump -u bamort -p<password> bamort > manual_backup.sql
```

---

### Issue 7: Cannot Rollback Migration

**Symptoms:**
- Rollback command fails
- "No DownSQL defined" error
- Table dependencies prevent DROP

**Error Examples:**
```
Error: migration #5 has no rollback script (DownSQL empty)
Error: Cannot drop table 'skills': foreign key constraint fails
```

**Diagnosis:**
```bash
# Check if migration has DownSQL
cat backend/deployment/migrations/all_migrations.go | grep -A 30 "Number: 5"

# Check table dependencies
docker exec bamort-mariadb mysql -u bamort -p<password> bamort \
  -e "SHOW CREATE TABLE skills"
```

**Solutions:**

**A. Missing DownSQL**
```bash
# Must restore from backup
docker exec bamort-backend /app/deploy backup restore \
  --file /app/backups/backup_<timestamp>.json
```

**B. Foreign Key Constraints**
```bash
# Disable FK checks temporarily
docker exec bamort-mariadb mysql -u bamort -p<password> bamort -e "
  SET FOREIGN_KEY_CHECKS=0;
  DROP TABLE IF EXISTS skills;
  SET FOREIGN_KEY_CHECKS=1;
"

# Then rollback migration number manually
```

---

### Issue 8: Container Health Check Failing

**Symptoms:**
- `docker ps` shows (unhealthy) status
- Container keeps restarting
- Services intermittently unavailable

**Diagnosis:**
```bash
# Check health status
docker inspect bamort-backend | jq '.[0].State.Health'

# Check health check command
docker inspect bamort-backend | jq '.[0].Config.Healthcheck'
```

**Solutions:**

**A. Backend Unhealthy**
```bash
# Check if backend actually responding
curl -f http://localhost:8180/api/system/health

# If not, check logs for errors
docker logs bamort-backend --tail=50
```

**B. Database Unhealthy**
```bash
# Check mariadb responding
docker exec bamort-mariadb mysqladmin ping

# If not, restart mariadb
docker-compose -f docker/docker-compose.yml restart mariadb
```

---

## Diagnostic Commands

### System Overview
```bash
# Complete system status
docker-compose -f docker/docker-compose.yml ps

# Resource usage
docker stats --no-stream

# Network connectivity
docker exec bamort-backend ping -c 3 mariadb
```

### Logs Analysis
```bash
# All logs from last hour
docker-compose -f docker/docker-compose.yml logs --since 1h

# Follow live logs
docker-compose -f docker/docker-compose.yml logs --follow

# Search for errors
docker logs bamort-backend 2>&1 | grep -i error | tail -20
```

### Database Inspection
```bash
# Connect to database
docker exec -it bamort-mariadb mysql -u bamort -p bamort

# Check tables
SHOW TABLES;

# Check version
SELECT * FROM schema_version ORDER BY id DESC LIMIT 1;

# Check migration history
SELECT * FROM migration_history ORDER BY migration_number DESC LIMIT 10;

# Exit
exit
```

---

## Error Messages Dictionary

| Error Message | Meaning | Solution |
|---------------|---------|----------|
| `record not found` | Database query returned no results | Normal in some cases, check context |
| `duplicate entry` | Trying to insert duplicate unique key | Use UPDATE or clean table first |
| `foreign key constraint fails` | Cannot delete/update due to FK | Delete child records first or disable FK checks |
| `table already exists` | Migration trying to create existing table | Migration already applied or rollback needed |
| `connection refused` | Cannot connect to database | Check mariadb running and credentials |
| `port already in use` | Port conflict | Kill process using port or change port |
| `no space left on device` | Disk full | Clean old files, logs, backups |
| `permission denied` | File/directory permission issue | Fix permissions with chmod/chown |

---

## When All Else Fails

### Nuclear Option: Complete Reset

⚠️ **WARNING**: This deletes ALL data!

```bash
# 1. Stop everything
cd /data/dev/bamort/docker
./stop-prd.sh

# 2. Remove all volumes
docker volume rm bamort-db
docker volume rm bamort-backend-tmp
docker volume rm bamort-frontend-tmp

# 3. Remove all containers
docker-compose -f docker-compose.yml rm -f

# 4. Start fresh
./start-prd.sh

# 5. Initialize
docker exec bamort-backend /app/deploy init \
  --masterdata /app/masterdata \
  --create-admin \
  --admin-user admin
```

---

## Getting Help

### 1. Gather Information

Before requesting help, collect:

```bash
# System info
docker-compose -f docker/docker-compose.yml ps > system_info.txt
docker version >> system_info.txt
docker-compose version >> system_info.txt

# Logs
docker logs bamort-backend --tail=200 > backend.log
docker logs bamort-mariadb --tail=200 > mariadb.log
docker logs bamort-frontend --tail=200 > frontend.log

# Version info
docker exec bamort-backend /app/deploy status > version_info.txt

# Database schema
docker exec bamort-mariadb mysqldump -u bamort -p --no-data bamort > schema.sql
```

### 2. Check Documentation

- [DEPLOYMENT_RUNBOOK.md](DEPLOYMENT_RUNBOOK.md) - Deployment procedures
- [ROLLBACK_GUIDE.md](ROLLBACK_GUIDE.md) - Rollback procedures
- [VERSION_COMPATIBILITY.md](VERSION_COMPATIBILITY.md) - Version requirements

### 3. Search Issues

- GitHub Issues: https://github.com/Bardioc26/bamort/issues
- Search for error message
- Check closed issues

### 4. Create Issue

If problem persists, create GitHub issue with:
- Error message (full stack trace)
- Steps to reproduce
- System info (from step 1)
- Logs (from step 1)
- What you've tried

---

**Last Updated:** 16. Januar 2026  
**Version:** 1.0
