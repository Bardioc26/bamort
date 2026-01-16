#!/bin/bash
# Production Deployment Script
# Usage: ./deploy-production.sh <version> [deployment-package.tar.gz]
#
# Example: ./deploy-production.sh v0.5.0
# Example: ./deploy-production.sh v0.5.0 deployment_package_0.5.0.tar.gz

set -e  # Exit on error

if [ -z "$1" ]; then
    echo "❌ Version required"
    echo "Usage: ./deploy-production.sh <version> [deployment-package]"
    echo "Example: ./deploy-production.sh v0.5.0"
    echo "Example: ./deploy-production.sh v0.5.0 deployment_package_0.5.0.tar.gz"
    exit 1
fi

VERSION="$1"
DEPLOYMENT_PACKAGE="$2"

echo "================================"
echo "Bamort PRODUCTION Deployment"
echo "Version: $VERSION"
echo "================================"
echo ""
echo "⚠️  WARNING: This will deploy to PRODUCTION"
echo ""
read -p "Type 'DEPLOY' to continue: " CONFIRM

if [ "$CONFIRM" != "DEPLOY" ]; then
    echo "Deployment cancelled"
    exit 0
fi

# Configuration
DOCKER_COMPOSE_FILE="docker/docker-compose.yml"
PROJECT_ROOT="/data/dev/bamort"
BACKUP_DIR="$PROJECT_ROOT/backups"
LOG_FILE="$PROJECT_ROOT/logs/deploy-$(date +%Y%m%d-%H%M%S).log"

mkdir -p "$PROJECT_ROOT/logs"

# Log function
log() {
    echo "$1" | tee -a "$LOG_FILE"
}

log ""
log "=== Deployment started at $(date) ==="
log ""

cd "$PROJECT_ROOT"

# Pre-deployment checks
log "→ Running pre-deployment checks..."

# Check disk space
AVAILABLE_GB=$(df -BG . | awk 'NR==2 {print $4}' | tr -d 'G')
if [ "$AVAILABLE_GB" -lt 2 ]; then
    log "❌ Insufficient disk space (${AVAILABLE_GB}GB available, 2GB required)"
    exit 1
fi
log "✓ Disk space: ${AVAILABLE_GB}GB available"

# Check Docker running
if ! docker ps > /dev/null 2>&1; then
    log "❌ Docker is not running"
    exit 1
fi
log "✓ Docker is running"

# Create backup
log ""
log "→ Creating backup..."
mkdir -p "$BACKUP_DIR"
BACKUP_FILE="$BACKUP_DIR/pre-deploy-$VERSION-$(date +%Y%m%d-%H%M%S).sql"

if docker ps | grep -q bamort-mariadb; then
    docker exec bamort-mariadb mysqldump -u bamort -p\${MARIADB_PASSWORD} bamort > "$BACKUP_FILE" 2>/dev/null && {
        BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
        log "✓ Backup created: $BACKUP_FILE ($BACKUP_SIZE)"
    } || {
        log "❌ Backup failed - aborting deployment"
        exit 1
    }
else
    log "❌ MariaDB not running - aborting deployment"
    exit 1
fi

# Pull specific version
log ""
log "→ Pulling version $VERSION..."
git fetch origin
git checkout "$VERSION" || {
    log "❌ Failed to checkout version $VERSION"
    exit 1
}
log "✓ Checked out version $VERSION"

# Pull Docker images
log ""
log "→ Pulling Docker images for $VERSION..."
docker-compose -f "$DOCKER_COMPOSE_FILE" build || {
    log "❌ Failed to build Docker images"
    exit 1
}
log "✓ Docker images built"

# Stop frontend
log ""
log "→ Stopping frontend to prevent user access..."
docker-compose -f "$DOCKER_COMPOSE_FILE" stop frontend
log "✓ Frontend stopped"

# Run database migrations
log ""
log "→ Running database migrations..."
docker exec bamort-backend /app/deploy migrate --verbose || {
    log "❌ Database migration failed"
    log "Rolling back..."
    docker-compose -f "$DOCKER_COMPOSE_FILE" down
    git checkout main
    docker-compose -f "$DOCKER_COMPOSE_FILE" up -d
    log ""
    log "To restore database backup:"
    log "  ./scripts/rollback.sh $BACKUP_FILE"
    exit 1
}
log "✓ Migrations completed successfully"

# Import master data if deployment package provided
if [ -n "$DEPLOYMENT_PACKAGE" ] && [ -f "$DEPLOYMENT_PACKAGE" ]; then
    log ""
    log "→ Extracting deployment package..."
    TEMP_DIR=$(mktemp -d)
    tar -xzf "$DEPLOYMENT_PACKAGE" -C "$TEMP_DIR" || {
        log "❌ Failed to extract deployment package"
        rm -rf "$TEMP_DIR"
        exit 1
    }
    log "✓ Package extracted"
    
    log ""
    log "→ Importing master data..."
    docker cp "$TEMP_DIR/masterdata" bamort-backend:/app/masterdata || {
        log "❌ Failed to copy master data"
        rm -rf "$TEMP_DIR"
        exit 1
    }
    
    docker exec bamort-backend /app/deploy import-masterdata --path /app/masterdata --verbose || {
        log "⚠️  Master data import had issues (check logs)"
    }
    log "✓ Master data import completed"
    
    rm -rf "$TEMP_DIR"
elif [ -n "$DEPLOYMENT_PACKAGE" ]; then
    log "⚠️  Deployment package not found: $DEPLOYMENT_PACKAGE"
else
    log "ℹ️  No deployment package provided, skipping master data import"
fi

# Update backend (restart to ensure clean state)
log ""
log "→ Restarting backend..."
docker-compose -f "$DOCKER_COMPOSE_FILE" restart backend

# Wait for backend
log "→ Waiting for backend to start (max 60s)..."
for i in {1..30}; do
    if curl -f -s http://localhost:8182/api/system/health > /dev/null 2>&1; then
        log "✓ Backend is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        log "❌ Backend failed to start within 60 seconds"
        log "Rolling back..."
        docker-compose -f "$DOCKER_COMPOSE_FILE" down
        git checkout main
        docker-compose -f "$DOCKER_COMPOSE_FILE" up -d
        exit 1
    fi
    sleep 4
done

# Check migrations
log ""
log "→ Checking system health..."
HEALTH_JSON=$(curl -s http://localhost:8182/api/system/health 2>/dev/null || echo "{}")
log "$HEALTH_JSON"

COMPATIBLE=$(echo "$HEALTH_JSON" | grep -o '"compatible":[^,}]*' | cut -d':' -f2 | tr -d ' ')
if [ "$COMPATIBLE" = "false" ]; then
    log "❌ Version incompatibility detected"
    log "Please check migrations and database version"
    log ""
    log "To rollback:"
    log "  ./scripts/rollback.sh $BACKUP_FILE"
    exit 1
fi
log "✓ System is compatible"

# Start frontend
log ""
log "→ Starting frontend..."
docker-compose -f "$DOCKER_COMPOSE_FILE" up -d frontend
sleep 3

if curl -f -s http://localhost:5174 > /dev/null 2>&1; then
    log "✓ Frontend is accessible"
else
    log "⚠️  Frontend may still be starting"
fi

# Final validation
log ""
log "→ Final validation..."

# Check all services running
SERVICES_OK=true
for service in backend frontend mariadb; do
    if docker-compose -f "$DOCKER_COMPOSE_FILE" ps | grep -q "$service.*Up"; then
        log "✓ $service is running"
    else
        log "❌ $service is not running"
        SERVICES_OK=false
    fi
done

if [ "$SERVICES_OK" = "false" ]; then
    log "❌ Some services are not running"
    exit 1
fi

# Deployment summary
log ""
log "================================"
log "Deployment Successful!"
log "================================"
log ""
log "Version deployed: $VERSION"
log "Backup location:  $BACKUP_FILE"
log "Log file:         $LOG_FILE"
log ""
log "Services:"
log "  Backend:  http://localhost:8182"
log "  Frontend: http://localhost:5174"
log ""
log "Next steps:"
log "  1. Monitor logs for errors"
log "  2. Test core functionality"
log "  3. Notify users of update"
log "  4. Monitor system for 24 hours"
log ""
log "To rollback:"
log "  ./scripts/rollback.sh $BACKUP_FILE"
log ""
log "=== Deployment completed at $(date) ==="
log ""
