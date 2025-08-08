#!/bin/bash

# Example usage of the /:id/improve/skill endpoint
# This endpoint gets improvement costs for all levels of a skill for a specific character

# Example 1: Get improvement costs for Menschenkenntnis for character ID 20
echo "Getting improvement costs for Menschenkenntnis (character ID 20):"
curl -X GET http://localhost:8180/api/characters/20/improve/skill \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Klettern"
  }'

echo -e "\n\nNote: This is a GET request with JSON body (unusual API design)"
echo "The endpoint returns an array of models.LearnCost objects with fields:"
echo "- stufe: target level"
echo "- ep: experience points cost"
echo "- money: money cost"
echo "- le: life energy cost"
