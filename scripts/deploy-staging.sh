#!/bin/bash
# Staging Deployment Script
# Usage: ./deploy-staging.sh [version]

set -e  # Exit on error

VERSION="${1:-latest}"

echo "================================"
echo "Bamort Staging Deployment"
echo "Version: $VERSION"
echo "================================"
echo ""

# Configuration
DOCKER_COMPOSE_FILE="docker/docker-compose.yml"
PROJECT_ROOT="/data/dev/bamort"
BACKUP_DIR="$PROJECT_ROOT/backups"

cd "$PROJECT_ROOT"

# Pre-deployment checks
echo "→ Running pre-deployment checks..."

# Check disk space
AVAILABLE_SPACE=$(df -h . | awk 'NR==2 {print $4}')
echo "  Available disk space: $AVAILABLE_SPACE"

# Create backup
echo ""
echo "→ Creating backup..."
mkdir -p "$BACKUP_DIR"
BACKUP_FILE="$BACKUP_DIR/pre-deploy-$(date +%Y%m%d-%H%M%S).sql"

if docker ps | grep -q bamort-mariadb; then
    docker exec bamort-mariadb mysqldump -u bamort -p\${MARIADB_PASSWORD} bamort > "$BACKUP_FILE" 2>/dev/null || {
        echo "⚠️  Backup failed - continuing without backup"
        rm -f "$BACKUP_FILE"
    }
    if [ -f "$BACKUP_FILE" ]; then
        echo "✓ Backup created: $BACKUP_FILE"
    fi
else
    echo "⚠️  MariaDB not running - skipping backup"
fi

# Pull latest changes
echo ""
echo "→ Pulling latest changes..."
git fetch origin
if [ "$VERSION" != "latest" ]; then
    git checkout "$VERSION"
fi
git pull

# Pull Docker images
echo ""
echo "→ Pulling Docker images..."
docker-compose -f "$DOCKER_COMPOSE_FILE" pull

# Stop frontend (prevent user access during migration)
echo ""
echo "→ Stopping frontend..."
docker-compose -f "$DOCKER_COMPOSE_FILE" stop frontend

# Restart backend (applies migrations)
echo ""
echo "→ Updating backend..."
docker-compose -f "$DOCKER_COMPOSE_FILE" up -d backend

# Wait for backend to be ready
echo "→ Waiting for backend to start..."
for i in {1..30}; do
    if curl -f -s http://localhost:8182/api/system/health > /dev/null 2>&1; then
        echo "✓ Backend is ready"
        break
    fi
    echo "  Waiting... ($i/30)"
    sleep 2
done

# Check for pending migrations
echo ""
echo "→ Checking migration status..."
HEALTH_CHECK=$(curl -s http://localhost:8182/api/system/health 2>/dev/null || echo "{}")
MIGRATIONS_PENDING=$(echo "$HEALTH_CHECK" | grep -o '"migrations_pending":[^,}]*' | cut -d':' -f2 || echo "unknown")

if [ "$MIGRATIONS_PENDING" = "true" ]; then
    echo "⚠️  Migrations pending - please run migrations manually"
    echo "  docker exec bamort-backend /app/deploy migrations apply --all"
elif [ "$MIGRATIONS_PENDING" = "false" ]; then
    echo "✓ No pending migrations"
fi

# Restart frontend
echo ""
echo "→ Starting frontend..."
docker-compose -f "$DOCKER_COMPOSE_FILE" up -d frontend

# Final health check
echo ""
echo "→ Final health check..."
sleep 3

if curl -f -s http://localhost:8182/api/system/health > /dev/null 2>&1; then
    echo "✓ Backend is healthy"
    curl -s http://localhost:8182/api/system/health | python3 -m json.tool 2>/dev/null || echo "  (health data available at /api/system/health)"
else
    echo "❌ Backend health check failed"
    echo ""
    echo "Recent logs:"
    docker logs bamort-backend --tail=30
    exit 1
fi

# Show status
echo ""
echo "================================"
echo "Deployment Complete!"
echo "================================"
echo ""
echo "Backend:  http://localhost:8182"
echo "Frontend: http://localhost:5174"
echo ""
echo "Backup saved to: $BACKUP_FILE"
echo ""
echo "Monitor logs:"
echo "  docker-compose -f $DOCKER_COMPOSE_FILE logs --follow"
echo ""
