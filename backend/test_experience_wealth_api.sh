#!/bin/bash

# Test-Script für die neuen Experience/Wealth API-Endpunkte

CHARACTER_ID=15
BASE_URL="http://localhost:8180/api/characters"
LOGIN_URL="http://localhost:8180/login"

echo "=== Testing Experience and Wealth Update API ==="
echo

# Zuerst einloggen und Token holen
echo "0. Logging in to get token..."
TOKEN_RESPONSE=$(curl -s -X POST "${LOGIN_URL}" \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser", "password": "testpass"}')

TOKEN=$(echo $TOKEN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "Login failed. Trying to run tests without authentication..."
  echo "Response: $TOKEN_RESPONSE"
  echo
else
  echo "Login successful. Token obtained."
  echo "Token: ${TOKEN:0:20}..."
  echo
fi

# Auth header für alle folgenden Requests
AUTH_HEADER=""
if [ ! -z "$TOKEN" ]; then
  AUTH_HEADER="-H \"Authorization: Bearer $TOKEN\""
fi

echo "=== Testing Experience and Wealth Update API ==="
echo

# Test 1: Update Experience Points with Reason
echo "1. Testing Experience Update with Reason..."
if [ ! -z "$TOKEN" ]; then
  curl -X PUT "${BASE_URL}/${CHARACTER_ID}/experience" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"experience_points": 150, "reason": "manual", "notes": "Test manual EP increase"}' \
    -w "\nStatus: %{http_code}\n\n"
else
  curl -X PUT "${BASE_URL}/${CHARACTER_ID}/experience" \
    -H "Content-Type: application/json" \
    -d '{"experience_points": 150, "reason": "manual", "notes": "Test manual EP increase"}' \
    -w "\nStatus: %{http_code}\n\n"
fi

# Test 2: Update Gold with Reason
echo "2. Testing Gold Update with Reason..."
if [ ! -z "$TOKEN" ]; then
  curl -X PUT "${BASE_URL}/${CHARACTER_ID}/wealth" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"goldstücke": 250, "reason": "equipment", "notes": "Sold equipment for gold"}' \
    -w "\nStatus: %{http_code}\n\n"
else
  curl -X PUT "${BASE_URL}/${CHARACTER_ID}/wealth" \
    -H "Content-Type: application/json" \
    -d '{"goldstücke": 250, "reason": "equipment", "notes": "Sold equipment for gold"}' \
    -w "\nStatus: %{http_code}\n\n"
fi

# Test 3: Simulate skill learning (EP reduction)
echo "3. Testing Skill Learning (EP reduction)..."
if [ ! -z "$TOKEN" ]; then
  curl -X PUT "${BASE_URL}/${CHARACTER_ID}/experience" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"experience_points": 100, "reason": "skill_learning", "notes": "Learned new sword fighting skill"}' \
    -w "\nStatus: %{http_code}\n\n"
else
  curl -X PUT "${BASE_URL}/${CHARACTER_ID}/experience" \
    -H "Content-Type: application/json" \
    -d '{"experience_points": 100, "reason": "skill_learning", "notes": "Learned new sword fighting skill"}' \
    -w "\nStatus: %{http_code}\n\n"
fi

# Test 4: Get Audit Log
echo "4. Testing Get Audit Log..."
if [ ! -z "$TOKEN" ]; then
  curl -X GET "${BASE_URL}/${CHARACTER_ID}/audit-log" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -w "\nStatus: %{http_code}\n\n"
else
  curl -X GET "${BASE_URL}/${CHARACTER_ID}/audit-log" \
    -H "Content-Type: application/json" \
    -w "\nStatus: %{http_code}\n\n"
fi

# Test 5: Get Audit Log Stats
echo "5. Testing Get Audit Log Stats..."
if [ ! -z "$TOKEN" ]; then
  curl -X GET "${BASE_URL}/${CHARACTER_ID}/audit-log/stats" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -w "\nStatus: %{http_code}\n\n"
else
  curl -X GET "${BASE_URL}/${CHARACTER_ID}/audit-log/stats" \
    -H "Content-Type: application/json" \
    -w "\nStatus: %{http_code}\n\n"
fi

# Test 6: Get Experience-only Audit Log
echo "6. Testing Get Experience-only Audit Log..."
if [ ! -z "$TOKEN" ]; then
  curl -X GET "${BASE_URL}/${CHARACTER_ID}/audit-log?field=experience_points" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -w "\nStatus: %{http_code}\n\n"
else
  curl -X GET "${BASE_URL}/${CHARACTER_ID}/audit-log?field=experience_points" \
    -H "Content-Type: application/json" \
    -w "\nStatus: %{http_code}\n\n"
fi

echo "=== Tests completed ==="
