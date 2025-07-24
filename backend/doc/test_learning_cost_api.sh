#!/bin/bash

# Konsolidierte Character Skill Cost API Examples
# Diese Skripts zeigen, wie die einheitliche API verwendet werden kann

BASE_URL="http://localhost:8080/api/characters"
CHAR_ID="1" # Beispiel Charakter-ID

echo "=== Konsolidierte Character Skill Cost API Examples ==="
echo

# 1. Fertigkeit lernen
echo "1. Fertigkeit lernen:"
curl -s -X POST "$BASE_URL/$CHAR_ID/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Stehlen",
    "type": "skill",
    "action": "learn"
  }' | jq .

echo
echo "2. Fertigkeit verbessern (nächste Stufe):"
curl -s -X POST "$BASE_URL/$CHAR_ID/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Menschenkenntnis", 
    "type": "skill",
    "action": "improve",
    "current_level": 14
  }' | jq .

echo
echo "3. Multi-Level Verbesserung:"
curl -s -X POST "$BASE_URL/$CHAR_ID/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Menschenkenntnis",
    "type": "skill", 
    "action": "improve",
    "current_level": 14,
    "target_level": 16
  }' | jq .

echo
echo "4. Mit Praxispunkten:"
curl -s -X POST "$BASE_URL/$CHAR_ID/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Menschenkenntnis",
    "type": "skill",
    "action": "improve", 
    "current_level": 14,
    "use_pp": 1
  }' | jq .

echo
echo "5. Zauber lernen:"
curl -s -X POST "$BASE_URL/$CHAR_ID/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Macht über das Selbst",
    "type": "spell",
    "action": "learn"
  }' | jq .

echo
echo "6. Legacy: Nächste Verbesserungsstufe:"
curl -s -X GET "$BASE_URL/$CHAR_ID/improve" | jq .

echo
echo "7. System-Info: Charakterklassen:"
curl -s -X GET "$BASE_URL/character-classes" | jq .

echo
echo "8. System-Info: Fertigkeitskategorien:"
curl -s -X GET "$BASE_URL/skill-categories" | jq .

echo
echo "=== Konsolidierte API-Beispiele abgeschlossen ==="
