#!/bin/bash

# Script to transfer data from SQLite test database to MariaDB using the maintenance endpoint
# Usage: ./transfer_sqlite_to_mariadb.sh [clear]

BACKEND_URL="http://localhost:8180"
ENDPOINT="/api/maintenance/transfer-sqlite-to-mariadb"

echo "=== Bamort Data Transfer: SQLite to MariaDB ==="
echo ""

# Check if clear parameter is provided
CLEAR_PARAM=""
if [ "$1" == "clear" ]; then
    echo "⚠️  WARNING: This will clear all existing data in MariaDB before transfer!"
    read -p "Are you sure you want to continue? (y/N): " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        echo "Transfer cancelled."
        exit 0
    fi
    CLEAR_PARAM="?clear=true"
    echo "✅ Proceeding with data clearing..."
else
    echo "ℹ️  Existing data in MariaDB will be preserved (use 'clear' parameter to clear first)"
fi

echo ""
echo "🔄 Starting data transfer..."
echo "📍 Source: ./testdata/prepared_test_data.db (SQLite)"
echo "📍 Target: mariadb:3306/bamort (MariaDB)"
echo ""

# Make the API request
echo "🚀 Calling transfer endpoint..."
response=$(curl -s -w "\n%{http_code}" -X POST "${BACKEND_URL}${ENDPOINT}${CLEAR_PARAM}")

# Extract HTTP status code (last line)
http_code=$(echo "$response" | tail -n1)
# Extract response body (all lines except last)
response_body=$(echo "$response" | head -n -1)

echo ""
if [ "$http_code" = "200" ]; then
    echo "✅ Transfer completed successfully!"
    echo ""
    echo "📊 Response:"
    echo "$response_body" | jq '.' 2>/dev/null || echo "$response_body"
else
    echo "❌ Transfer failed with HTTP status: $http_code"
    echo ""
    echo "📋 Error details:"
    echo "$response_body" | jq '.' 2>/dev/null || echo "$response_body"
    exit 1
fi

echo ""
echo "🎉 Data transfer complete!"
echo ""
echo "💡 Tips:"
echo "  - Check the statistics above to verify the transfer"
echo "  - Use the backend logs for detailed information"
echo "  - Update your .env file to use MariaDB if needed"
