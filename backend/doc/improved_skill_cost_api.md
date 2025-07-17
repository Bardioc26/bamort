# Improved Skill Cost API Documentation

## Overview
The improved `/api/characters/:id/skill-cost` endpoint provides enhanced functionality for calculating learning and improvement costs for skills, weapons, and spells.

## Endpoint
```
POST /api/characters/:id/skill-cost
```

## Request Body

### Basic Request Structure
```json
{
  "name": "string (required)",
  "type": "string (required): skill|spell|weapon",
  "action": "string (required): learn|improve",
  "current_level": "number (optional)",
  "target_level": "number (optional)"
}
```

### Validation Rules
- `name`: Required, skill/spell/weapon name
- `type`: Required, must be one of: `skill`, `spell`, `weapon`
- `action`: Required, must be one of: `learn`, `improve`
- `current_level`: Optional for `learn`, required for `improve` (unless character already has the skill)
- `target_level`: Optional, enables multi-level cost calculation

## Response Formats

### Single Cost Response
```json
{
  "stufe": 11,
  "le": 0,
  "ep": 200,
  "money": 200,
  "skill_name": "Menschenkenntnis",
  "skill_type": "skill",
  "action": "improve",
  "character_id": 1,
  "current_level": 10,
  "target_level": 0,
  "category": "Sozial",
  "difficulty": "schwer",
  "can_afford": true,
  "notes": "Verbesserung von 10 auf 11. Kosten für Hexer"
}
```

### Multi-Level Cost Response
```json
{
  "skill_name": "Menschenkenntnis",
  "skill_type": "skill",
  "character_id": 1,
  "current_level": 10,
  "target_level": 13,
  "level_costs": [
    {
      "stufe": 11,
      "le": 0,
      "ep": 200,
      "money": 200,
      "skill_name": "Menschenkenntnis",
      "skill_type": "skill",
      "action": "improve",
      "character_id": 1,
      "current_level": 10,
      "target_level": 11,
      "category": "Sozial",
      "difficulty": "schwer",
      "can_afford": true
    },
    {
      "stufe": 12,
      "le": 0,
      "ep": 200,
      "money": 200,
      "skill_name": "Menschenkenntnis",
      "skill_type": "skill",
      "action": "improve",
      "character_id": 1,
      "current_level": 11,
      "target_level": 12,
      "category": "Sozial",
      "difficulty": "schwer",
      "can_afford": true
    },
    {
      "stufe": 13,
      "le": 0,
      "ep": 400,
      "money": 400,
      "skill_name": "Menschenkenntnis",
      "skill_type": "skill",
      "action": "improve",
      "character_id": 1,
      "current_level": 12,
      "target_level": 13,
      "category": "Sozial",
      "difficulty": "schwer",
      "can_afford": true
    }
  ],
  "total_cost": {
    "stufe": 13,
    "le": 0,
    "ep": 800,
    "money": 800
  },
  "can_afford_total": true
}
```

## API Improvements

### 1. Enhanced Validation
- **Binding validation**: Uses Gin's binding tags for automatic validation
- **Input sanitization**: Trims whitespace from skill names
- **Type safety**: Validates skill types and actions against allowed values

### 2. Automatic Level Detection
- **Smart level detection**: Automatically detects current skill level from character data
- **Fallback handling**: Gracefully handles missing skill data

### 3. Multi-Level Cost Calculation
- **Range calculations**: Calculate costs for improving from current level to target level
- **Detailed breakdown**: Shows individual costs for each level improvement
- **Total cost summary**: Provides cumulative cost for the entire improvement path

### 4. Enhanced Response Data
- **Skill metadata**: Includes category and difficulty information
- **Affordability check**: Indicates whether character can afford the costs
- **Contextual notes**: Provides helpful information about the cost calculation

### 5. Better Error Handling
- **Detailed error messages**: More specific error descriptions
- **Graceful degradation**: Handles missing data gracefully
- **Input validation**: Comprehensive validation with clear error messages

## Example Usage

### Learn a New Skill
```bash
curl -X POST "http://localhost:8180/api/characters/1/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Menschenkenntnis",
    "type": "skill",
    "action": "learn"
  }'
```

### Improve an Existing Skill
```bash
curl -X POST "http://localhost:8180/api/characters/1/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Menschenkenntnis",
    "type": "skill",
    "action": "improve",
    "current_level": 10
  }'
```

### Multi-Level Improvement
```bash
curl -X POST "http://localhost:8180/api/characters/1/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Menschenkenntnis",
    "type": "skill",
    "action": "improve",
    "current_level": 10,
    "target_level": 15
  }'
```

### Learn a Spell
```bash
curl -X POST "http://localhost:8180/api/characters/1/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Feuerball",
    "type": "spell",
    "action": "learn"
  }'
```

### Improve a Weapon Skill
```bash
curl -X POST "http://localhost:8180/api/characters/1/skill-cost" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Beidhändiger Kampf",
    "type": "weapon",
    "action": "improve",
    "current_level": 5
  }'
```

## Error Responses

### Validation Errors
```json
{
  "error": "Ungültige Anfrageparameter: Key: 'SkillCostRequest.Type' Error:Field validation for 'Type' failed on the 'oneof' tag"
}
```

### Character Not Found
```json
{
  "error": "Charakter nicht gefunden"
}
```

### Skill Not Found
```json
{
  "error": "Fehler bei der Kostenberechnung: unbekannte Fertigkeit"
}
```

### Invalid Level Range
```json
{
  "error": "Fehler bei der Multi-Level-Kostenberechnung"
}
```

## Benefits of the Improvements

1. **Better User Experience**: More detailed and helpful responses
2. **Reduced API Calls**: Multi-level calculation reduces need for multiple requests
3. **Better Error Handling**: Clear, actionable error messages
4. **Automatic Detection**: Reduces manual input requirements
5. **Comprehensive Data**: Includes all relevant information for frontend display
6. **Type Safety**: Prevents invalid requests through validation
7. **Extensibility**: Easy to add new skill types or actions
8. **Performance**: Efficient calculation and response structure
