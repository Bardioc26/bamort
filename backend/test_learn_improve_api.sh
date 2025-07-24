#!/bin/bash

# Test-Script für die neuen Learn/Improve API-Endpunkte mit Audit-Log

BASE_URL="http://localhost:8080/api/characters"
CHAR_ID="18"  # Test-Charakter ID
TOKEN="Bearer dc7a780.1:bba7f4daabda117f2a2c14263"

echo "Testing Learn/Improve API with Audit-Log..."
echo "=============================================="

# Test 1: Fertigkeit lernen
echo -e "\n1. Testing LearnSkill endpoint..."
curl -s -X POST "$BASE_URL/$CHAR_ID/learn-skill" \
  -H "Content-Type: application/json" \
  -H "Authorization: $TOKEN" \
  -d '{
    "name": "Bootfahren",
    "notes": "Test: Fertigkeit lernen über API",
    "use_pp": 0
  }' | jq .

# Test 2: Fertigkeit verbessern
echo -e "\n2. Testing ImproveSkill endpoint..."
curl -s -X POST "$BASE_URL/$CHAR_ID/improve-skill" \
  -H "Content-Type: application/json" \
  -H "Authorization: $TOKEN" \
  -d '{
    "name": "Schwimmen",
    "current_level": 8,
    "notes": "Test: Fertigkeit verbessern über API",
    "use_pp": 1
  }' | jq .

# Test 3: Zauber lernen
echo -e "\n3. Testing LearnSpell endpoint..."
curl -s -X POST "$BASE_URL/$CHAR_ID/learn-spell" \
  -H "Content-Type: application/json" \
  -H "Authorization: $TOKEN" \
  -d '{
    "name": "Angst",
    "notes": "Test: Zauber lernen über API"
  }' | jq .

# Test 4: Zauber verbessern
echo -e "\n4. Testing ImproveSpell endpoint..."
curl -s -X POST "$BASE_URL/$CHAR_ID/improve-spell" \
  -H "Content-Type: application/json" \
  -H "Authorization: $TOKEN" \
  -d '{
    "name": "Licht",
    "current_level": 3,
    "notes": "Test: Zauber verbessern über API"
  }' | jq .

# Test 5: Audit-Log prüfen
echo -e "\n5. Checking Audit-Log after learn/improve operations..."
curl -s -X GET "$BASE_URL/$CHAR_ID/audit-log" \
  -H "Authorization: $TOKEN" | jq '.entries[] | select(.reason | contains("skill") or contains("spell")) | {timestamp, field_name, old_value, new_value, difference, reason, notes}'

# Test 6: Audit-Log Statistiken
echo -e "\n6. Checking Audit-Log Statistics..."
curl -s -X GET "$BASE_URL/$CHAR_ID/audit-log/stats" \
  -H "Authorization: $TOKEN" | jq .

echo -e "\nTesting completed!"
