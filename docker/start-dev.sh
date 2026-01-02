#!/bin/bash

echo "🚀 Starting Bamort Development Environment..."

# Prüfe ob Docker läuft
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker ist nicht gestartet. Bitte starte Docker zuerst."
    exit 1
fi

# Gehe ins Docker-Verzeichnis
cd "$(dirname "$0")"

# Get current git commit
export GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
echo "📝 Git Commit: $GIT_COMMIT"

echo "📦 Building and starting development containers..."

# Stoppe vorhandene Container
docker-compose -f docker-compose.dev.yml down

# Baue und starte die Container
COMPOSE_PROJECT_NAME=bamort docker-compose -f docker-compose.dev.yml up --build -d

echo "✅ Development environment started."
