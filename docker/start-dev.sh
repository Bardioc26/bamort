#!/bin/bash

echo "🚀 Starting Bamort Development Environment..."

# Prüfe ob Docker läuft
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker ist nicht gestartet. Bitte starte Docker zuerst."
    exit 1
fi

# Gehe ins Docker-Verzeichnis
cd "$(dirname "$0")"

echo "📦 Building and starting development containers..."

# Stoppe vorhandene Container
docker-compose -f docker-compose.dev.yml down

# Baue und starte die Container
docker-compose -f docker-compose.dev.yml up --build

echo "✅ Development environment stopped."
