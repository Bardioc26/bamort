# BaMoRT Deployment Guide

This directory contains the deployment tool for BaMoRT. This guide explains the complete deployment workflow from preparing a new release to deploying it on the target system.

## Table of Contents

- [Overview](#overview)
- [Deployment Tool Commands](#deployment-tool-commands)
- [Pre-Deployment Checklist](#pre-deployment-checklist)
- [Step-by-Step Deployment Process](#step-by-step-deployment-process)
  - [1. Development System: Prepare Release](#1-development-system-prepare-release)
  - [2. Development System: Version Management](#2-development-system-version-management)
  - [3. Development System: Create Deployment Package](#3-development-system-create-deployment-package)
  - [4. Git: Commit and Tag](#4-git-commit-and-tag)
  - [5. Target System: Pre-Deployment](#5-target-system-pre-deployment)
  - [6. Target System: Deployment](#6-target-system-deployment)
  - [7. Target System: Post-Deployment](#7-target-system-post-deployment)
- [Rollback Procedure](#rollback-procedure)
- [Troubleshooting](#troubleshooting)

## Overview

The BaMoRT deployment process uses a multi-phase approach:

1. **Version Update**: Set new version numbers across the codebase
2. **Package Preparation**: Export master data and system configurations
3. **Git Management**: Commit, tag, and push to repository
4. **Target Deployment**: Pull code, migrate database, import data, restart services
5. **Validation**: Verify database schema, service health, and data integrity

The deployment tool (`backend/cmd/deploy/main.go`) provides several commands to manage this process safely.

## Deployment Tool Commands

Build the deployment tool first:

```bash
export BASEDIR=$(pwd)
cd $BASEDIR/backend
go build -o deploy cmd/deploy/main.go
```

Available commands:

| Command | Description | Usage |
|---------|-------------|-------|
| `version` | Show version information | `./deploy version` |
| `status` | Show current DB version and pending migrations | `./deploy status` |
| `prepare [dir]` | Create deployment package with master data | `./deploy prepare [export_dir]` |
| `deploy [dir]` | Run full deployment (backup → migrate → import → validate) | `./deploy deploy [import_dir]` |
| `validate` | Validate database schema and data integrity | `./deploy validate` |
| `help` | Show help message | `./deploy help` |

## Pre-Deployment Checklist

Before starting a deployment, ensure:

- [ ] All features are tested locally
- [ ] All tests pass: `cd backend && go test ./...`
- [ ] Frontend builds without errors
- [ ] Database migrations are tested locally
- [ ] Breaking changes are documented
- [ ] You have access to the target system
- [ ] Target system has sufficient disk space (>2GB free) mostly for docker images
- [ ] Docker is running on target system
- [ ] You know the new version number (semantic versioning: MAJOR.MINOR.PATCH)

## Step-by-Step Deployment Process

### 1. Development System: Prepare Release

Ensure your development environment is clean and up-to-date:

```bash
# Navigate to project root
cd $BASEDIR

# Ensure you're on the main branch
git checkout main
git pull origin main

# Verify all changes are committed
git status

# Run tests
cd backend && go test ./...
cd ../frontend && npm run test
```

### 2. Development System: Version Management

Update the version number across all components using the automated script:

```bash
cd $BASEDIR

# Update both backend and frontend to same version
./scripts/update-version.sh 0.1.38

# Or update to different versions
./scripts/update-version.sh 0.1.38 0.2.0

# Or use auto-commit mode (sets version AND commits)
./scripts/update-version.sh 0.1.38 auto
```

This script updates:
- `backend/config/version.go` - Backend application version
- `frontend/src/version.js` - Frontend application version
- `frontend/package.json` - NPM package version
- `backend/VERSION.md` - Backend version documentation
- `frontend/VERSION.md` - Frontend version documentation

**Manual version update** (if not using script):

Edit `backend/config/version.go`:
```go
const Version = "0.1.38"
```

Edit `frontend/src/version.js`:
```js
export const VERSION = '0.1.38'
```

Edit `frontend/package.json`:
```json
{
  "version": "0.1.38"
}
```

### 3. Development System: Create Deployment Package

Create a deployment package containing all master data and system configurations:

```bash
cd ./backend

# Build the deployment tool
go build -o deploy cmd/deploy/main.go

# Check current database status
./deploy status

# Create deployment package (exports to ./tmp by default)
./deploy prepare

# Or specify custom export directory
./deploy prepare /path/to/export_dir
```

The deployment package includes:
- Master data (skills, spells, equipment definitions, etc.)
- System configurations
- Database structure metadata
- Version information

**Important**: The deployment package does NOT include user data (characters, user accounts). User data remains on the target system and is migrated during deployment.

Package files:
- Export directory: `backend/tmp/`
- Archive: `backend/deployment_package_<version>_<timestamp>.tar.gz`

**Transfer the archive file** to the target system for deployment.

### 4. Git: Commit and Tag

Commit all version changes and create a git tag:

```bash
cd $BASEDIR

# If using auto-commit mode, skip this section
# Otherwise, commit manually:

# Add version files
git add backend/config/version.go
git add frontend/src/version.js
git add frontend/package.json
git add backend/VERSION.md
git add frontend/VERSION.md

# Commit with descriptive message
git commit -m "Bump version to 0.1.38"

# Create annotated tag
git tag -a v0.1.38 -m "Release version 0.1.38

Features:
- Feature 1 description
- Feature 2 description

Bug Fixes:
- Fix 1 description
- Fix 2 description

Breaking Changes:
- None (or list breaking changes)
"

# Push commits and tags
git push origin main
git push origin v0.1.38
```

**Tag naming convention**:
- Standard release: `v0.1.38`
- Backend-specific: `backend-v0.1.38` (if versioned separately)
- Frontend-specific: `frontend-v0.1.38` (if versioned separately)

### 5. Target System: Pre-Deployment

SSH into the target system and prepare for deployment:

```bash
# SSH to production server
ssh user@production-server

# Navigate to project directory
cd $BASEDIR  # Or production install location

# Check Docker status
docker ps

# Check disk space (need >2GB free)
df -h .

# Check current running version
curl http://localhost:8182/api/system/health | jq .

# Verify database connectivity
docker exec bamort-backend /app/deploy status
```

**Pre-deployment checks**:
- [ ] All services are running (`docker ps`)
- [ ] Sufficient disk space available (`df -h`)
- [ ] Database is accessible
- [ ] No pending migrations or issues

### 6. Target System: Deployment

Run the deployment script on the target system:

```bash
cd $BASEDIR

# Option 1: Deploy with migrations and master data import
# (Recommended for version upgrades with new game content)
./scripts/deploy-production.sh v0.1.38 deployment_package_0.1.38_20260118-120000.tar.gz

# Option 2: Deploy with migrations only (no master data changes)
# (Use for bug fixes or feature updates without game content changes)
./scripts/deploy-production.sh v0.1.38
```

**Deployment will prompt for confirmation**:
```
⚠️  WARNING: This will deploy to PRODUCTION
Type 'DEPLOY' to continue:
```

**The script performs these steps automatically**:

1. **Pre-flight checks**
   - Verify disk space (minimum 2GB required)
   - Verify Docker is running
   - Verify MariaDB is accessible

2. **Backup current database**
   - Creates timestamped backup in `backups/` directory
   - Format: `pre-deploy-v0.1.38-20260117-143022.sql`
   - Skips backup on fresh installation (when tables don't exist)
   - Aborts deployment if backup fails on existing installation

3. **Checkout version**
   - Fetches from git origin
   - Checks out the specified tag (e.g., `v0.1.38`)

4. **Build Docker images**
   - Builds new backend and frontend containers
   - Uses production Dockerfiles

5. **Stop frontend**
   - Stops frontend container to prevent user access during migration
   - Backend remains running

6. **Extract deployment package** (if provided)
   - Extracts deployment package to temporary directory
   - Copies master data to backend container
   - Prepares import directory path

7. **Run deployment command**
   - Executes `deploy deploy [importDir]` in backend container
   - Creates backup of current database state (skipped on fresh install)
   - Exports current master data (skipped on fresh install)
   - Checks version compatibility
   - Applies pending database migrations
   - Imports master data (if package provided)
   - Validates database schema
   - Automatically rolls back on failure

8. **Restart backend**
   - Restarts backend container with new code
   - Ensures clean state

9. **Health checks**
   - Waits for backend to start (max 120 seconds)
   - Verifies API endpoint responds
   - Checks version compatibility
   - Validates database schema

10. **Start frontend**
    - Starts frontend container
    - Verifies frontend accessibility

11. **Final validation**
    - Verifies all services are running
    - Reports deployment status
    - Cleans up temporary files

**Fresh Installation**: On first deployment to an empty database, backup and export steps are automatically skipped. This is expected behavior and not an error.

**Deployment log**: Saved to `logs/deploy-YYYYMMDD-HHMMSS.log`

### 7. Target System: Post-Deployment

After successful deployment, perform these validation steps:

```bash
# Check service status
docker ps

# View logs
docker logs bamort-backend --tail=100
docker logs bamort-frontend --tail=100

# Check system health
curl http://localhost:8182/api/system/health | jq .

# Verify database version
docker exec bamort-backend /app/deploy status

# Validate database schema
docker exec bamort-backend /app/deploy validate
```

**Manual Testing**:
1. Open the application in a browser
2. Login with test account
3. Navigate through main features
4. Create/edit a character
5. Generate a PDF export
6. Check responsive behavior on mobile


## Rollback Procedure

If deployment fails or critical issues are discovered:

### Automatic Rollback

The deployment script automatically rolls back if:
- Database backup fails
- Git checkout fails
- Docker build fails
- Database migration fails
- Backend fails to start within 60 seconds
- Version incompatibility detected

Automatic rollback performs:
1. Stops all containers
2. Checks out previous version (main branch)
3. Restarts containers
4. Displays rollback instructions for database

### Manual Rollback

If you need to rollback manually:

```bash
cd $BASEDIR

# Stop all services
docker-compose -f docker/docker-compose.yml down

# Restore database backup
BACKUP_FILE="backups/pre-deploy-v0.1.38-20260117-143022.sql"
cat "$BACKUP_FILE" | docker exec -i bamort-mariadb mysql -u bamort -p bamort

# Checkout previous version
git checkout v0.1.37  # Or main branch

# Rebuild and restart services
docker-compose -f docker/docker-compose.yml build
docker-compose -f docker/docker-compose.yml up -d

# Verify services
docker ps
curl http://localhost:8182/api/system/health | jq .
```

### Post-Rollback

After rollback:
1. Identify root cause of deployment failure
2. Fix issues in development environment
3. Test thoroughly
4. Increment version number
5. Retry deployment process

## Troubleshooting

### Problem: Version mismatch after deployment

**Symptoms**: Health check shows `"compatible": false`

**Solution**:
```bash
# Check versions
docker exec bamort-backend /app/deploy status

# Check for pending migrations
docker exec bamort-backend /app/deploy validate

# Run migrations manually if needed
docker exec bamort-backend /app/deploy deploy
```

### Problem: Backend won't start

**Symptoms**: Backend container exits immediately

**Solution**:
```bash
# Check logs
docker logs bamort-backend

# Common causes:
# 1. Database connection issues - check DATABASE_URL env var
# 2. Missing environment variables - check .env file
# 3. Port conflicts - check if port 8180 is in use

# Check database connectivity
docker exec bamort-backend sh -c 'nc -zv mariadb 3306'
```

### Problem: Frontend shows old version

**Symptoms**: UI displays previous version number

**Solution**:
```bash
# Clear browser cache
# Or force reload: Ctrl+Shift+R (Linux/Windows) or Cmd+Shift+R (Mac)

# Rebuild frontend container
docker-compose -f docker/docker-compose.yml build frontend
docker-compose -f docker/docker-compose.yml up -d frontend
```

### Problem: Migration fails

**Symptoms**: Migration error during deployment

**Solution**:
```bash
# Check migration status
docker exec bamort-backend /app/deploy status

# Check specific migration file
# Migrations located in: backend/deployment/migrations/

# Test migration locally first (without import)
cd $BASEDIR/backend
go build -o deploy cmd/deploy/main.go
./deploy deploy

# Or test with import
./deploy prepare ./test_export
./deploy deploy ./test_export

# Fix migration code if needed
# Rollback production deployment
# Test fixed migration locally
# Redeploy
```

### Problem: Deployment package import fails

**Symptoms**: Master data import errors

**Solution**:
```bash
# Check deployment package contents
tar -tzf deployment_package_0.1.38.tar.gz

# Verify package was created correctly
cd $BASEDIR/backend
go build -o deploy cmd/deploy/main.go
./deploy prepare ./test_export
ls -lh ./test_export/

# Import manually if needed
docker cp ./test_export bamort-backend:/tmp/import_data
docker exec bamort-backend /app/deploy deploy /tmp/import_data

# Or just run migrations without import
docker exec bamort-backend /app/deploy deploy
```

### Problem: Insufficient disk space

**Symptoms**: Deployment fails at pre-flight check

**Solution**:
```bash
# Check disk usage
df -h

# Clean up old Docker images
docker system prune -a

# Clean up old backups (keep last 10)
cd $BASEDIR/backups
ls -lt *.sql | tail -n +11 | awk '{print $NF}' | xargs rm

# Clean up old logs
cd $BASEDIR/logs
find . -name "deploy-*.log" -mtime +30 -delete
```

### Problem: Docker daemon not running

**Symptoms**: `Cannot connect to Docker daemon`

**Solution**:
```bash
# Start Docker service
sudo systemctl start docker

# Enable Docker on boot
sudo systemctl enable docker

# Verify Docker status
sudo systemctl status docker
```

---

## Quick Reference

### Development System Commands
```bash
# Update version
./scripts/update-version.sh 0.1.38 auto

# Create deployment package (includes tar.gz archive)
cd backend && go build -o deploy cmd/deploy/main.go
./deploy prepare
# Transfer the generated .tar.gz file to target system
```

### Git Commands
```bash
git tag -a v0.1.38 -m "Release v0.1.38"
git push origin main
git push origin v0.1.38
```

### Target System Commands
```bash
# Deploy without master data import (migrations only)
./scripts/deploy-production.sh v0.1.38

# Deploy with master data import
./scripts/deploy-production.sh v0.1.38 deployment_package_0.1.38.tar.gz

# Or run deployment tool directly
docker exec bamort-backend /app/deploy deploy              # Migrations only
docker exec bamort-backend /app/deploy deploy /import/dir  # With import

# Rollback
./scripts/rollback.sh backups/pre-deploy-v0.1.38-TIMESTAMP.sql

# Status check
docker exec bamort-backend /app/deploy status
docker exec bamort-backend /app/deploy validate
```

---

**Last Updated**: 2026-01-17  
**Version**: 1.0  
**Maintainer**: BaMoRT Development Team
