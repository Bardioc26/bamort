# Belohnungssystem API

Das Belohnungssystem ermöglicht es, verschiedene Arten von Belohnungen auf Lern- und Verbesserungskosten anzuwenden.

## Belohnungstypen

### 1. `free_learning` - Kostenlose Fertigkeiten
- **Anwendung**: Beim Lernen neuer Fertigkeiten (`action: "learn", type: "skill"`)
- **Effekt**: Die Geldkosten für das Lernen der Fertigkeit entfallen (Gold = 0)
- **EP/LE**: Bleiben unverändert

**Beispiel:**
```json
{
    "name": "Stichwaffen",
    "type": "skill",
    "action": "learn",
    "reward": {
        "type": "free_learning"
    }
}
```

### 2. `free_spell_learning` - Kostenlose Zauber
- **Anwendung**: Beim Lernen neuer Zauber (`action: "learn", type: "spell"`)
- **Effekt**: Die Lehrzeit für das Lernen des Zaubers entfällt (LE = 0)
- **EP/Gold**: Bleiben unverändert

**Beispiel:**
```json
{
    "name": "Licht",
    "type": "spell", 
    "action": "learn",
    "reward": {
        "type": "free_spell_learning"
    }
}
```

### 3. `half_ep_improvement` - Halbe EP für Verbesserungen
- **Anwendung**: Bei Fertigkeitsverbesserungen (`action: "improve"`)
- **Effekt**: Die EP-Kosten werden halbiert
- **LE/Gold**: Bleiben unverändert

**Beispiel:**
```json
{
    "name": "Stichwaffen",
    "type": "skill",
    "action": "improve",
    "current_level": 8,
    "reward": {
        "type": "half_ep_improvement"
    }
}
```

### 4. `gold_for_ep` - Gold statt EP verwenden
- **Anwendung**: Bei allen Aktionen mit EP-Kosten
- **Effekt**: EP werden durch Gold ersetzt (Wechselkurs: 10 GS = 1 EP)
- **Parameter**: 
  - `use_gold_for_ep`: Muss `true` sein
  - `max_gold_ep`: Maximale EP die durch Gold ersetzt werden (optional, Standard: Hälfte der EP-Kosten)

**Beispiel:**
```json
{
    "name": "Schwert",
    "type": "skill",
    "action": "improve",
    "current_level": 10,
    "reward": {
        "type": "gold_for_ep",
        "use_gold_for_ep": true,
        "max_gold_ep": 5
    }
}
```

## Kombination mit Praxispunkten

Belohnungen können mit Praxispunkten kombiniert werden:

```json
{
    "name": "Stichwaffen",
    "type": "skill",
    "action": "improve",
    "current_level": 12,
    "use_pp": 2,
    "reward": {
        "type": "half_ep_improvement"
    }
}
```

**Anwendungsreihenfolge:**
1. Zunächst werden Belohnungen angewendet
2. Dann werden Praxispunkte angewendet

## Response-Struktur

Die API-Response enthält zusätzliche Informationen über die angewendeten Belohnungen:

```json
{
    "ep": 15,
    "le": 0,
    "money": 50,
    "skill_name": "Stichwaffen",
    "action": "improve",
    "can_afford": true,
    
    // Belohnungsdetails
    "reward_applied": "half_ep_improvement",
    "original_cost_struct": {
        "ep": 30,
        "le": 0,
        "money": 50
    },
    "savings": {
        "ep": 15,
        "le": 0,
        "money": 0
    },
    "gold_used_for_ep": 0,
    
    // Praxispunkt-Details
    "pp_used": 2,
    "pp_available": 5,
    "pp_reduction": 2,
    "original_cost": 30,
    "final_cost": 13
}
```

## Felderbeschreibung

- `reward_applied`: Art der angewendeten Belohnung
- `original_cost_struct`: Ursprüngliche Kosten ohne Belohnung
- `savings`: Ersparnisse durch die Belohnung
- `gold_used_for_ep`: Anzahl EP die durch Gold ersetzt wurden
- `original_cost`: EP-Kosten vor PP-Reduktion
- `final_cost`: Finale EP-Kosten nach allen Reduktionen

## Curl-Beispiele

**Wichtig**: Alle API-Aufrufe benötigen einen Authorization Header.

### Kostenlose Fertigkeit lernen
```bash
curl -X POST http://localhost:8180/api/characters/20/skill-cost \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0fec2f7.1:2f9ecb4206b1fc311d91b5530" \
  -d '{
    "name": "Stichwaffen",
    "type": "skill",
    "action": "learn",
    "reward": {
      "type": "free_learning"
    }
  }'
```

### Halbe EP für Verbesserung mit Praxispunkten
```bash
curl -X POST http://localhost:8180/api/characters/20/skill-cost \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0fec2f7.1:2f9ecb4206b1fc311d91b5530" \
  -d '{
    "name": "Stichwaffen",
    "type": "skill",
    "action": "improve",
    "current_level": 12,
    "use_pp": 2,
    "reward": {
      "type": "half_ep_improvement"
    }
  }'
```

### Gold statt EP verwenden
```bash
curl -X POST http://localhost:8180/api/characters/20/skill-cost \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 0fec2f7.1:2f9ecb4206b1fc311d91b5530" \
  -d '{
    "name": "Schwert",
    "type": "skill",
    "action": "improve",
    "current_level": 10,
    "reward": {
      "type": "gold_for_ep",
      "use_gold_for_ep": true,
      "max_gold_ep": 5
    }
  }'
```
