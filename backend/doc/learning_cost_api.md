# Character Skill Cost API

Konsolidierte API für die Berechnung von Lern- und Verbesserungskosten für Fertigkeiten und Zauber.

## Hauptendpunkt

### Kostenberechnung für Charakter
```
POST /api/characters/:id/skill-cost
```

Dies ist der **Hauptendpunkt** für alle Kostenberechnungen. Er nutzt das bestehende GSMaster-System und die Praxispunkte-Integration.

**Request Body:**
```json
{
  "name": "Stehlen",
  "type": "skill",
  "action": "learn",
  "current_level": 14,
  "target_level": 16,
  "use_pp": 2
}
```

**Response:**
```json
{
  "skill_name": "Stehlen",
  "skill_type": "skill", 
  "action": "improve",
  "character_id": 123,
  "current_level": 14,
  "target_level": 16,
  "category": "Halbwelt",
  "difficulty": "schwer",
  "stufe": 16,
  "le": 0,
  "ep": 180,
  "money": 180,
  "original_cost": 200,
  "final_cost": 180,
  "pp_used": 1,
  "pp_available": 2,
  "description": "Verbesserung von Stehlen von +14 auf +16 (mit 1 PP)"
}
```

### Parameter

#### Pflichtfelder
- `name`: Name der Fertigkeit oder des Zaubers
- `type`: "skill", "spell" oder "weapon"  
- `action`: "learn" oder "improve"

#### Optionale Felder
- `current_level`: Aktueller Wert (automatisch ermittelt falls nicht angegeben)
- `target_level`: Zielwert (für Multi-Level-Berechnung)
- `use_pp`: Anzahl der zu verwendenden Praxispunkte

## Legacy-Endpunkte

### Nächste Stufe
```
GET /api/characters/:id/improve
```
Zeigt Kosten für die nächste Verbesserungsstufe aller Fertigkeiten.

### Alle Stufen
```
GET /api/characters/:id/improve/skill  
```
Zeigt Kosten für alle möglichen Verbesserungsstufen.

### Einfache Lernkosten
```
GET /api/characters/:id/learn
```
Zeigt grundlegende Lernkosten für neue Fertigkeiten.

## System-Information

### Charakterklassen
```
GET /api/characters/character-classes
```

**Response:**
```json
{
  "character_classes": {
    "Sp": {
      "code": "Sp",
      "name": "Spitzbube", 
      "skill_ep_costs": {
        "Halbwelt": 10,
        "Sozial": 10,
        "Waffen": 20
      },
      "spell_ep_costs": {}
    }
  }
}
```

### Fertigkeitskategorien
```
GET /api/characters/skill-categories
```

**Response:**
```json
{
  "skill_categories": {
    "Halbwelt": {
      "name": "Halbwelt",
      "learn_costs": {
        "leicht": 1,
        "schwer": 2
      },
      "improve_costs": {
        "schwer": {
          "15": 20,
          "16": 50
        }
      }
    }
  }
}
```

## Praxispunkte-Integration

Das System unterstützt automatisch Praxispunkte:

1. **Automatische Ermittlung**: Verfügbare PP werden automatisch geladen
2. **Kostenreduzierung**: 1 PP = 1 LE weniger Kosten
3. **Validierung**: Prüfung ob genügend PP vorhanden sind

### Beispiel mit Praxispunkten
```json
{
  "name": "Menschenkenntnis",
  "type": "skill", 
  "action": "improve",
  "current_level": 14,
  "use_pp": 2
}
```

## Features

### ✅ Automatische Erkennung
- Charakterklasse aus Charakter-Daten
- Aktueller Fertigkeitswert aus Charakter
- Fertigkeitskategorie und Schwierigkeit aus GSMaster-DB

### ✅ Multi-Level-Berechnung  
- Kosten für mehrere Stufen auf einmal
- Berücksichtigung gestaffelter Kosten

### ✅ Praxispunkte-Integration
- Automatische PP-Verwaltung
- Kostenreduzierung durch PP
- PP-Verfügbarkeitsprüfung

### ✅ Vollständige Validierung
- Charakterklassen-Beschränkungen
- Fertigkeits-Verfügbarkeit
- Level-Grenzen

## Fehlerbehandlung

Das System gibt detaillierte Fehlermeldungen zurück:

```json
{
  "error": "Fertigkeit nicht bei diesem Charakter vorhanden oder current_level erforderlich"
}
```

## Migration von alten Endpunkten

**Alt:** Multiple verschiedene Endpunkte
**Neu:** Ein einheitlicher Endpunkt `POST /api/characters/:id/skill-cost`

Der neue Endpunkt kombiniert alle Funktionalitäten und bietet:
- Konsistente Request/Response-Struktur
- Vollständige GSMaster-Integration  
- Automatische Praxispunkte-Berücksichtigung
- Multi-Level-Kostenberechnung
