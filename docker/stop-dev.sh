#!/bin/bash

echo "🛑 Stopping Bamort Development Environment..."

# Gehe ins Docker-Verzeichnis
cd "$(dirname "$0")"

# Stoppe und entferne Container
docker-compose -f docker-compose.dev.yml down

echo "✅ Development environment stopped."
