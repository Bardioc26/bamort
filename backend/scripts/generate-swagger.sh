#!/bin/bash
# generate-swagger.sh - Generate Swagger documentation for BaMoRT API

set -e

echo "===================================="
echo "Generating Swagger Documentation"
echo "===================================="
echo ""

# Change to backend directory
cd "$(dirname "$0")/.."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate Swagger docs
echo "Generating Swagger specification..."
swag init -g cmd/main.go -o docs/swagger

echo ""
echo "✓ Swagger documentation generated successfully!"
echo ""
echo "Generated files:"
echo "  - docs/swagger/docs.go"
echo "  - docs/swagger/swagger.json"
echo "  - docs/swagger/swagger.yaml"
echo ""
echo "To view the documentation:"
echo "  1. Start the backend server: cd docker && ./start-dev.sh"
echo "  2. Open browser to: http://localhost:8180/swagger/index.html"
echo ""
echo "===================================="
