# Praxispunkte (PP) Funktionalität

## Übersicht

Die Praxispunkte-Funktionalität ermöglicht es Spielern, ihre Charaktere effizienter zu entwickeln, indem sie gesammelte Praxispunkte nutzen können, um die Kosten für das Erlernen und Verbessern von Fertigkeiten und Zaubern zu reduzieren.

## Regeln

### Ausbildungsmöglichkeiten mit Praxispunkten (PP)

1. **1 TE für 1 PP** beim Verbessern der Fertigkeit, für die der PP vergeben wurde
2. **1 LE für 1 PP** beim Lernen von Zaubern aus der Gruppe, für die der PP vergeben wurde

## API-Endpunkte

### Praxispunkte abrufen
```
GET /api/characters/{id}/practice-points
```
Gibt alle verfügbaren Praxispunkte für einen Charakter zurück.

### Praxispunkte hinzufügen
```
POST /api/characters/{id}/practice-points/add
```
Fügt Praxispunkte zu einer Kategorie hinzu.

Request Body:
```json
{
  "kategorie": "Sozial",
  "anzahl": 3
}
```

### Praxispunkte verwenden
```
POST /api/characters/{id}/practice-points/use
```
Verbraucht Praxispunkte aus einer Kategorie.

Request Body:
```json
{
  "kategorie": "Sozial", 
  "anzahl": 1
}
```

### Kostenberechnung mit Praxispunkten
```
POST /api/characters/{id}/skill-cost
```
Berechnet die Kosten für das Erlernen oder Verbessern einer Fertigkeit/Zauber und berücksichtigt dabei die Verwendung von Praxispunkten.

Request Body:
```json
{
  "name": "Menschenkenntnis",
  "type": "skill",
  "action": "improve",
  "current_level": 10,
  "use_pp": 2
}
```

Response:
```json
{
  "stufe": 11,
  "le": 0,
  "ep": 180,
  "money": 200,
  "skill_name": "Menschenkenntnis",
  "skill_type": "skill",
  "action": "improve",
  "character_id": 1,
  "current_level": 10,
  "target_level": 11,
  "category": "Sozial",
  "difficulty": "schwer",
  "can_afford": true,
  "notes": "Verbesserung von 10 auf 11. Kosten für Hx. Verwendung von 2 Praxispunkten",
  "pp_used": 2,
  "pp_available": 5,
  "pp_reduction": 2,
  "original_cost": 200,
  "final_cost": 180
}
```

## Datenstruktur

### Praxispunkt
```go
type Praxispunkt struct {
    CharacterID uint   `json:"character_id"`
    SkillName   string `json:"skill_name"` // Name der spezifischen Fertigkeit
    Anzahl      int    `json:"anzahl"`     // Anzahl der verfügbaren PP
}
```

## Kategorien

### Fertigkeiten
- Alltag
- Freiland  
- Halbwelt
- Kampf
- Körper
- Sozial
- Unterwelt
- Waffen
- Wissen

### Zauber
- Beherr (Beherrschen)
- Beweg (Bewegen)
- Erkenn (Erkennen)
- Erschaff (Erschaffen)
- Formen
- Veränd (Verändern)
- Zerstören
- Wunder
- Dweom
- Lied

## Verwendung

1. **Fertigkeiten verbessern**: Wenn ein Charakter eine Fertigkeit verbessert, können verfügbare Praxispunkte der entsprechenden Kategorie verwendet werden, um die benötigten Trainingseinheiten (TE) zu reduzieren.

2. **Zauber lernen**: Beim Erlernen neuer Zauber können Praxispunkte der entsprechenden Zauberschule verwendet werden, um die benötigten Lerneinheiten (LE) zu reduzieren.

3. **Automatische Berechnung**: Das System berechnet automatisch die reduzierten Kosten basierend auf den verfügbaren und verwendeten Praxispunkten.

## Implementierungsdetails

- Praxispunkte werden pro Kategorie verwaltet
- Die API verhindert die Verwendung von mehr PP als verfügbar sind
- Originalkosten und reduzierte Kosten werden beide in der Response angezeigt
- Geldkosten bleiben unverändert, nur EP/LE-Kosten werden reduziert
- Multi-Level-Verbesserungen verteilen PP intelligent auf die verschiedenen Level
