# Praxispunkte (PP) API - Fertigkeitsspezifisch (Vereinfacht)

## Übersicht

Praxispunkte sind in Midgard **fertigkeitsspezifisch** vergeben. Jede Fertigkeit und jeder Zauber kann individuelle Praxispunkte haben, die nur für die Verbesserung oder das Lernen dieser spezifischen Fertigkeit verwendet werden können.

**Vereinfachung**: Der SkillType wird automatisch anhand des Fertigkeitsnamens ermittelt, sodass nur der Name angegeben werden muss.

## Regeln

- **1 TE für 1 PP** beim Verbessern von Fertigkeiten
- **1 LE für 1 PP** beim Lernen von Zaubern aus der Gruppe
- Praxispunkte sind **fertigkeitsspezifisch** - sie können nur für die Fertigkeit verwendet werden, für die sie erworben wurden
- Praxispunkte für "Menschenkenntnis" können nur für "Menschenkenntnis" verwendet werden
- Praxispunkte für "Bannbalken" können nur für "Bannbalken" verwendet werden

## Datenstruktur

```go
type Praxispunkt struct {
    CharacterID uint   `json:"character_id"`
    SkillName   string `json:"skill_name"` // Name der spezifischen Fertigkeit
    Anzahl      int    `json:"anzahl"`     // Anzahl der verfügbaren PP
}
```

## API-Endpunkte

### 1. Praxispunkte abrufen
```http
GET /api/characters/{id}/practice-points
```

**Response:**
```json
[
    {
        "character_id": 1,
        "skill_name": "Menschenkenntnis",
        "anzahl": 5
    },
    {
        "character_id": 1,
        "skill_name": "Bannbalken",
        "anzahl": 3
    }
]
```

### 2. Praxispunkte hinzufügen
```http
POST /api/characters/{id}/practice-points/add
```

**Request Body:**
```json
{
    "skill_name": "Menschenkenntnis",
    "anzahl": 2
}
```

### 3. Praxispunkte verwenden
```http
POST /api/characters/{id}/practice-points/use
```

**Request Body:**
```json
{
    "skill_name": "Menschenkenntnis",
    "anzahl": 1
}
```

### 4. Aktualisierte Praxispunkte setzen
```http
PUT /api/characters/{id}/practice-points
```

**Request Body:**
```json
[
    {
        "skill_name": "Menschenkenntnis",
        "anzahl": 3
    },
    {
        "skill_name": "Bannbalken", 
        "anzahl": 2
    }
]
```

## Integration in Kostenberechnung

Die PP-Funktionalität ist in die bestehende `/api/characters/{id}/skill-cost` API integriert:

```http
POST /api/characters/{id}/skill-cost
```

**Request Body mit PP:**
```json
{
    "name": "Menschenkenntnis",
    "type": "skill",
    "action": "improve",
    "current_level": 10,
    "use_pp": 2
}
```

**Response:**
```json
{
    "skill_name": "Menschenkenntnis",
    "skill_type": "skill",
    "action": "improve",
    "character_id": 1,
    "current_level": 10,
    "target_level": 11,
    "ep": 8,
    "le": 0,
    "money": 0,
    "pp_used": 2,
    "pp_available": 5,
    "pp_reduction": 2,
    "original_cost": 10,
    "final_cost": 8,
    "can_afford": true
}
```

## Beispiele

### Praxispunkte für Fertigkeit vergeben
```bash
curl -X POST http://localhost:8180/api/characters/1/practice-points/add \
  -H "Content-Type: application/json" \
  -d '{
    "skill_name": "Menschenkenntnis",
    "anzahl": 3
  }'
```

### Praxispunkte für Zauber vergeben
```bash
curl -X POST http://localhost:8180/api/characters/1/practice-points/add \
  -H "Content-Type: application/json" \
  -d '{
    "skill_name": "Bannbalken",
    "anzahl": 2
  }'
```

### Kostenberechnung mit PP für Fertigkeitsverbesserung
```bash
curl -X POST http://localhost:8180/api/characters/1/skill-cost \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Menschenkenntnis",
    "type": "skill",
    "action": "improve",
    "current_level": 10,
    "use_pp": 2
  }'
```

### Kostenberechnung mit PP für Zauber lernen
```bash
curl -X POST http://localhost:8180/api/characters/1/skill-cost \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bannbalken",
    "type": "spell", 
    "action": "learn",
    "use_pp": 1
  }'
```

## Technische Details

### Automatische Skill-Type-Erkennung

Das System erkennt automatisch den Typ einer Fertigkeit anhand des Namens:

```go
func determineSkillType(skillName string) string {
    // Prüft automatisch in den Tabellen:
    // - skill_skills (Fertigkeiten)
    // - skill_weaponskills (Waffenfertigkeiten)
    // - skill_spells (Zauber)
}
```

### Validierung

In der Produktionsumgebung kann optional validiert werden, ob eine Fertigkeit existiert. In der Test-Umgebung ist dies deaktiviert für Flexibilität.
