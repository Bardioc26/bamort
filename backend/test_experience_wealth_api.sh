#!/bin/bash

# Test-Script für die neuen Experience/Wealth API-Endpunkte

CHARACTER_ID=1
BASE_URL="http://localhost:8080/api/characters"

echo "=== Testing Experience and Wealth Update API ==="
echo

# Test 1: Update Experience Points
echo "1. Testing Experience Update..."
curl -X PUT "${BASE_URL}/${CHARACTER_ID}/experience" \
  -H "Content-Type: application/json" \
  -d '{"experience_points": 150}' \
  -w "\nStatus: %{http_code}\n\n"

# Test 2: Update Gold
echo "2. Testing Gold Update..."
curl -X PUT "${BASE_URL}/${CHARACTER_ID}/wealth" \
  -H "Content-Type: application/json" \
  -d '{"goldstücke": 250}' \
  -w "\nStatus: %{http_code}\n\n"

# Test 3: Update multiple wealth values
echo "3. Testing Multiple Wealth Update..."
curl -X PUT "${BASE_URL}/${CHARACTER_ID}/wealth" \
  -H "Content-Type: application/json" \
  -d '{"goldstücke": 100, "silberstücke": 50, "kupferstücke": 200}' \
  -w "\nStatus: %{http_code}\n\n"

# Test 4: Get Experience and Wealth
echo "4. Testing Get Experience and Wealth..."
curl -X GET "${BASE_URL}/${CHARACTER_ID}/experience-wealth" \
  -H "Content-Type: application/json" \
  -w "\nStatus: %{http_code}\n\n"

# Test 5: Invalid Experience (negative)
echo "5. Testing Invalid Experience (negative)..."
curl -X PUT "${BASE_URL}/${CHARACTER_ID}/experience" \
  -H "Content-Type: application/json" \
  -d '{"experience_points": -10}' \
  -w "\nStatus: %{http_code}\n\n"

echo "=== Tests completed ==="
