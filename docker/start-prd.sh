#!/bin/bash

echo "🚀 Starting Bamort Production Environment..."

# Prüfe ob Docker läuft
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker ist nicht gestartet. Bitte starte Docker zuerst."
    exit 1
fi

# Gehe ins Docker-Verzeichnis
cd "$(dirname "$0")"

echo "📦 Building and starting production containers..."
# Build before stopping existing containers
docker-compose -f docker-compose.yml build

# Stoppe vorhandene Container
docker-compose -f docker-compose.yml down

# Baue und starte die Container
docker-compose -f docker-compose.yml up --build -d

echo "✅ Production environment started."