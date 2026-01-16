#!/bin/bash
# Development Deployment Script
# Usage: ./deploy-dev.sh

set -e  # Exit on error

echo "================================"
echo "Bamort Development Deployment"
echo "================================"
echo ""

# Configuration
DOCKER_COMPOSE_FILE="docker/docker-compose.dev.yml"
PROJECT_ROOT="/data/dev/bamort"

cd "$PROJECT_ROOT"

# Pre-deployment checks
echo "→ Running pre-deployment checks..."

if ! docker --version > /dev/null 2>&1; then
    echo "❌ Docker not installed"
    exit 1
fi

if ! docker-compose --version > /dev/null 2>&1; then
    echo "❌ Docker Compose not installed"
    exit 1
fi

echo "✓ Docker check passed"

# Check if containers are running
echo ""
echo "→ Checking container status..."
if docker ps | grep -q bamort-backend-dev; then
    echo "✓ Backend container running"
else
    echo "⚠️  Backend container not running"
fi

if docker ps | grep -q bamort-frontend-dev; then
    echo "✓ Frontend container running"
else
    echo "⚠️  Frontend container not running"
fi

# Pull latest changes
echo ""
echo "→ Pulling latest changes..."
git pull origin main || echo "⚠️  Git pull failed (continuing anyway)"

# Rebuild containers
echo ""
echo "→ Rebuilding containers..."
docker-compose -f "$DOCKER_COMPOSE_FILE" build

# Restart services
echo ""
echo "→ Restarting services..."
docker-compose -f "$DOCKER_COMPOSE_FILE" up -d

# Wait for services to start
echo ""
echo "→ Waiting for services to start..."
sleep 5

# Check health
echo ""
echo "→ Checking system health..."
if curl -f -s http://localhost:8180/api/system/health > /dev/null 2>&1; then
    echo "✓ Backend is healthy"
else
    echo "❌ Backend health check failed"
    docker logs bamort-backend-dev --tail=20
    exit 1
fi

if curl -f -s http://localhost:5173 > /dev/null 2>&1; then
    echo "✓ Frontend is accessible"
else
    echo "⚠️  Frontend not accessible (may still be starting)"
fi

# Show status
echo ""
echo "================================"
echo "Deployment Complete!"
echo "================================"
echo ""
echo "Backend:  http://localhost:8180"
echo "Frontend: http://localhost:5173"
echo ""
echo "View logs:"
echo "  docker logs bamort-backend-dev --follow"
echo "  docker logs bamort-frontend-dev --follow"
echo ""
