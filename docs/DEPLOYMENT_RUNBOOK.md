# Deployment Runbook

**Version:** 1.0  
**Last Updated:** 16. Januar 2026  
**Target System:** BaMoRT 

---

## Table of Contents
1. [Overview](#overview)
2. [Pre-Deployment Checklist](#pre-deployment-checklist)
3. [Deployment Procedures](#deployment-procedures)
4. [Post-Deployment Validation](#post-deployment-validation)
5. [Common Scenarios](#common-scenarios)
6. [Troubleshooting](#troubleshooting)

---

## Overview

### What is a Deployment?

A deployment updates the Bamort system with:
- New backend code (Go)
- New frontend code (Vue.js)
- Database migrations (schema changes)
- Master data updates (skills, spells, etc.)

### Deployment Frequency

- **Expected**: Once per week to once per month
- **Duration**: 20-50 minutes including backup
- **Downtime**: System unavailable during deployment

### Deployment Components

```
Deployment Package Contains:
├── Backend Binary (`deploy`)
├── Database Migrations (SQL scripts)
├── Master Data (JSON exports)
├── Frontend Build (static files)
└── Deployment Metadata (version info)
```

---

## Pre-Deployment Checklist

###Before You Start

- [ ] **Backup Current System**
  ```bash
  docker exec bamort-backend /app/deploy backup create
  ```

- [ ] **Verify Docker Containers Running**
  ```bash
  docker ps | grep bamort
  # Should show: backend, frontend, mariadb, phpmyadmin
  ```

- [ ] **Check Current Version**
  ```bash
  docker exec bamort-backend /app/deploy status
  ```

- [ ] **Check Disk Space**
  ```bash
  df -h
  # Ensure at least 2GB free
  ```

- [ ] **Notify Users** (if applicable)
  - Send maintenance notification
  - Set maintenance window (e.g., 22:00-23:00)

- [ ] **Review Changes**
  - Check `CHANGELOG.md`
  - Review migration scripts in `backend/deployment/migrations/`
  - Review new features/fixes

### Required Information

- [ ] **Current Version**: ______________
- [ ] **Target Version**: ______________
- [ ] **Migration Count**: ______________
- [ ] **Deployment Package**: ______________
- [ ] **Backup Location**: ______________

---

## Deployment Procedures

### Procedure 1: Standard Deployment (Update Existing System)

**When to Use:** Updating an existing Bamort installation

**Steps:**

#### 1. Create Pre-Deployment Backup

```bash
# Create full backup
docker exec bamort-backend /app/deploy backup create

# Verify backup created
docker exec bamort-backend /app/deploy backup list
```

**Expected Output:**
```
Backup created: backup_20260116_220000_v0.4.0_m3.json
Size: 2.4 MB
```

#### 2. Stop Frontend (Prevent User Access)

```bash
docker-compose -f docker/docker-compose.yml stop frontend
```

**Verification:**
```bash
curl http://localhost:5173
# Should fail or show "connection refused"
```

#### 3. Apply Database Migrations

```bash
# Check pending migrations
docker exec bamort-backend /app/deploy migrations pending

# Apply all pending migrations
docker exec bamort-backend /app/deploy migrations apply --all
```

**Expected Output:**
```
Found 2 pending migrations:
  Migration 4: Add learning_category column
  Migration 5: Create equipment_cache table

Applying migrations...
✓ Migration 4 completed (executed in 450ms)
✓ Migration 5 completed (executed in 320ms)

All migrations applied successfully.
Database version: 0.4.0 → 0.5.0
```

**If Errors Occur:** See [Rollback Guide](ROLLBACK_GUIDE.md)

#### 4. Import Master Data

```bash
# Check what will be imported
docker exec bamort-backend /app/deploy masterdata import \
  --source /app/masterdata \
  --dry-run

# Import master data
docker exec bamort-backend /app/deploy masterdata import \
  --source /app/masterdata
```

**Expected Output:**
```
Importing master data from /app/masterdata...
✓ Sources: 12 records (2 new, 10 updated)
✓ Character Classes: 15 records (0 new, 0 updated)
✓ Skills: 245 records (5 new, 8 updated)
✓ Spells: 189 records (3 new, 2 updated)
✓ Learning Costs: 1,234 records (45 new, 12 updated)

Master data import completed successfully.
```

#### 5. Restart Backend

```bash
# Pull latest backend image
docker-compose -f docker/docker-compose.yml pull backend

# Restart backend
docker-compose -f docker/docker-compose.yml restart backend

# Wait for backend to start (check logs)
docker logs bamort-backend --tail=50 --follow
# Press Ctrl+C when you see "Server started on :8180"
```

#### 6. Restart Frontend

```bash
# Pull latest frontend image
docker-compose -f docker/docker-compose.yml pull frontend

# Start frontend
docker-compose -f docker/docker-compose.yml start frontend
```

#### 7. Verify Deployment

See [Post-Deployment Validation](#post-deployment-validation) section below.

---

### Procedure 2: Fresh Installation (New System)

**When to Use:** Setting up Bamort for the first time

**Steps:**

#### 1. Prepare Environment

```bash
# Clone repository
git clone https://github.com/Bardioc26/bamort.git
cd bamort

# Copy environment file
cp .env.example .env

# Edit configuration
nano .env
# Set:
#   DATABASE_PASSWORD=<secure-password>
#   JWT_SECRET=<random-secret>
```

#### 2. Start Docker Containers

```bash
# Start all services
cd docker
./start-prd.sh

# Verify all containers running
docker ps | grep bamort
```

#### 3. Initialize Database

```bash
# Run initialization
docker exec bamort-backend /app/deploy init \
  --masterdata /app/masterdata \
  --create-admin \
  --admin-user admin \
  --admin-password <secure-password>
```

**Expected Output:**
```
Initializing new Bamort installation...
Backend version: 0.5.0

Step 1/4: Creating database schema...
✓ Database schema created successfully

Step 2/4: Initializing version tracking...
✓ Version tracking initialized (DB version: 0.5.0)

Step 3/4: Importing master data...
✓ Master data imported successfully

Step 4/4: Creating admin user...
✓ Admin user 'admin' created successfully

═══════════════════════════════════════════
Installation completed successfully!
Version: 0.5.0
Tables created: 42
Admin created: Yes
Master data: Imported
Execution time: 8.5s
═══════════════════════════════════════════
```

#### 4. Verify Installation

```bash
# Check system health
curl http://localhost:8180/api/system/health | jq

# Access frontend
open http://localhost:5173

# Login with admin credentials
```

---

### Procedure 3: Rollback Deployment

See [ROLLBACK_GUIDE.md](ROLLBACK_GUIDE.md) for detailed rollback procedures.

---

## Post-Deployment Validation

### Automated Checks

Run the validation script:

```bash
docker exec bamort-backend /app/deploy validate
```

**Expected Output:**
```
Running post-deployment validation...

Database Checks:
✓ All expected tables exist
✓ All expected columns exist
✓ No orphaned tables
✓ Foreign key constraints valid
✓ Indexes created

Version Checks:
✓ Backend version: 0.5.0
✓ Database version: 0.5.0
✓ Versions compatible

Data Integrity:
✓ Master data tables populated
✓ Sources: 12 records
✓ Skills: 245 records
✓ Spells: 189 records

All validation checks passed ✓
```

### Manual Verification

#### 1. Check System Health

```bash
curl http://localhost:8180/api/system/health | jq
```

**Expected:**
```json
{
  "status": "ok",
  "backend_version": "0.5.0",
  "required_db_version": "0.5.0",
  "actual_db_version": "0.5.0",
  "migrations_pending": false,
  "compatible": true,
  "timestamp": "2026-01-16T22:30:00Z"
}
```

#### 2. Check Frontend

- [ ] Open http://localhost:5173
- [ ] No warning banner displayed
- [ ] Login works
- [ ] Dashboard loads
- [ ] Character list loads
- [ ] Can view character details

#### 3. Check Backend Logs

```bash
docker logs bamort-backend --tail=100
```

**Look for:**
- [ ] No ERROR messages
- [ ] Server started successfully
- [ ] Database connection OK
- [ ] No migration errors

#### 4. Test Core Functionality

- [ ] **Create Character**: Can create new character
- [ ] **Edit Character**: Can edit existing character
- [ ] **Skills**: Skills load and display correctly
- [ ] **Spells**: Spells load and display correctly
- [ ] **Learning Costs**: Learning cost calculations work
- [ ] **Export**: Can export character to PDF

#### 5. Check Database

```bash
# Connect to database
docker exec -it bamort-mariadb mysql -u bamort -p bamort

# Check version
SELECT * FROM schema_version ORDER BY id DESC LIMIT 1;

# Check migration history
SELECT * FROM migration_history ORDER BY migration_number DESC LIMIT 5;

# Exit
exit
```

---

## Common Scenarios

### Scenario 1: Minor Update (No Migrations)

When deployment only updates code without database changes:

```bash
# No migrations needed, just restart services
docker-compose -f docker/docker-compose.yml pull
docker-compose -f docker/docker-compose.yml restart
```

### Scenario 2: Major Update (Multiple Migrations)

When deployment includes 5+ migrations:

```bash
# Create backup first
docker exec bamort-backend /app/deploy backup create

# Apply migrations with verbose output
docker exec bamort-backend /app/deploy migrations apply --all --verbose

# Verify each migration
docker exec bamort-backend /app/deploy migrations history
```

### Scenario 3: Hotfix Deployment

Emergency fix during business hours:

```bash
# 1. Create backup (fast JSON backup)
docker exec bamort-backend /app/deploy backup create

# 2. Pull latest code
docker-compose -f docker/docker-compose.yml pull backend

# 3. Restart backend only
docker-compose -f docker/docker-compose.yml restart backend

# 4. Monitor logs
docker logs bamort-backend --tail=100 --follow

# 5. Verify health endpoint
curl http://localhost:8180/api/system/health
```

### Scenario 4: Master Data Update Only

When only updating skills/spells/costs:

```bash
# Import new master data (dev exports to production)
docker cp ./masterdata bamort-backend:/tmp/masterdata
docker exec bamort-backend /app/deploy masterdata import \
  --source /tmp/masterdata
```

---

## Troubleshooting

### Issue: Migration Fails

**Symptom:** Migration error during `migrations apply`

**Solution:**
1. Check error message in logs
2. If safe, rollback last migration:
   ```bash
   docker exec bamort-backend /app/deploy migrations rollback --steps 1
   ```
3. Fix migration script if needed
4. Re-apply migration

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for detailed solutions.

### Issue: Frontend Shows Warning Banner

**Symptom:** Yellow warning banner: "Database migration required"

**Cause:** Database version doesn't match backend version

**Solution:**
```bash
# Check versions
docker exec bamort-backend /app/deploy status

# Apply pending migrations
docker exec bamort-backend /app/deploy migrations apply --all
```

### Issue: Backend Won't Start

**Symptom:** Backend container crashes or restarts repeatedly

**Solution:**
```bash
# Check logs
docker logs bamort-backend --tail=100

# Common causes:
# - Database connection failed → Check DATABASE_PASSWORD in .env
# - Migration failed → Rollback to previous version
# - Port conflict → Check if port 8180 is available
```

### Issue: Cannot Access Frontend

**Symptom:** http://localhost:5173 not reachable

**Solution:**
```bash
# Check frontend container
docker ps | grep frontend

# If not running, start it
docker-compose -f docker/docker-compose.yml start frontend

# Check frontend logs
docker logs bamort-frontend --tail=50
```

---

## Deployment Timeline

Typical deployment timeline:

| Step | Duration | Description |
|------|----------|-------------|
| 1. Backup | 2-5 min | Create database backup |
| 2. Stop Frontend | 10 sec | Prevent user access |
| 3. Migrations | 5-15 min | Apply database changes |
| 4. Master Data | 3-8 min | Import updated data |
| 5. Restart Backend | 1-2 min | Load new code |
| 6. Restart Frontend | 30 sec | Load new UI |
| 7. Validation | 5-10 min | Verify deployment |
| **Total** | **20-45 min** | Complete deployment |

---

## Next Steps

After successful deployment:

1. **Monitor System** for 24 hours
   - Check error logs daily
   - Monitor system health endpoint
   - Watch for user-reported issues

2. **Update Documentation**
   - Update `CHANGELOG.md`
   - Document any manual changes
   - Update version in README

3. **Cleanup Old Backups**
   ```bash
   docker exec bamort-backend /app/deploy backup cleanup --keep 30
   ```

4. **Notify Users**
   - Send "System Updated" notification
   - Highlight new features
   - Mention any breaking changes

---

**For Emergency Rollback:** See [ROLLBACK_GUIDE.md](ROLLBACK_GUIDE.md)  
**For Common Issues:** See [TROUBLESHOOTING.md](TROUBLESHOOTING.md)  
**For Version Info:** See [VERSION_COMPATIBILITY.md](VERSION_COMPATIBILITY.md)
