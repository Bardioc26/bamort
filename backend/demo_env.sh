#!/bin/bash

echo "=== Bamort .env Konfiguration Demo ==="
echo ""

# Erstelle verschiedene .env-Dateien für Demo
echo "# Development .env" > .env.dev.demo
echo "ENVIRONMENT=development" >> .env.dev.demo
echo "DEBUG=true" >> .env.dev.demo
echo "LOG_LEVEL=DEBUG" >> .env.dev.demo
echo "PORT=8180" >> .env.dev.demo
echo "" >> .env.dev.demo

echo "# Production .env" > .env.prod.demo
echo "ENVIRONMENT=production" >> .env.prod.demo
echo "DEBUG=false" >> .env.prod.demo
echo "LOG_LEVEL=INFO" >> .env.prod.demo
echo "PORT=8080" >> .env.prod.demo
echo "" >> .env.prod.demo

echo "1. Development-Konfiguration (.env.dev.demo):"
cat .env.dev.demo
echo ""

echo "2. Production-Konfiguration (.env.prod.demo):"
cat .env.prod.demo
echo ""

echo "3. Test mit Development-Konfiguration:"
cp .env.dev.demo .env
echo "Building and testing with development config..."
go run cmd/logging_demo.go 2>/dev/null | head -10
echo ""

echo "4. Test mit Production-Konfiguration:"
cp .env.prod.demo .env
echo "Building and testing with production config..."
go run cmd/logging_demo.go 2>/dev/null | head -10
echo ""

echo "5. Test mit Umgebungsvariablen-Override:"
echo "Setze DEBUG=false als Umgebungsvariable (sollte .env überschreiben)"
DEBUG=false go run cmd/logging_demo.go 2>/dev/null | head -10
echo ""

# Cleanup
rm -f .env.dev.demo .env.prod.demo

echo "=== Demo beendet ==="
echo ""
echo "Tipp: Kopiere .env.example nach .env und passe sie an deine Bedürfnisse an:"
echo "cp .env.example .env"
