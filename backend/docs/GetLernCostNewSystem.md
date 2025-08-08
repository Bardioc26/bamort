# GetLernCostNewSystem - Neue Datenbank-basierte Lernkosten

## Übersicht

Die Funktion `GetLernCostNewSystem` ist eine neue Implementierung des Lernkosten-Systems, die anstelle der hardkodierten `learningCostsData` eine vollständige Datenbank-Lösung verwendet.

## Funktionen

### Alte vs Neue Implementation

| Aspect | GetLernCost (Alt) | GetLernCostNewSystem (Neu) |
|--------|------------------|---------------------------|
| Datenquelle | `learningCostsData` (hardkodiert) | Datenbank-Tabellen (`learning_*`) |
| Flexibilität | Statisch | Dynamisch konfigurierbar |
| Verwaltung | Code-Änderungen nötig | Admin-Interface möglich |
| Performance | Schnell (im Speicher) | Leicht langsamer (DB-Queries) |
| Skalierbarkeit | Begrenzt | Unbegrenzt |

### API-Kompatibilität

Beide Funktionen verwenden die gleiche Request/Response-Struktur:

```go
// Request (identisch)
type LernCostRequest struct {
    CharId       uint    `json:"char_id"`
    Name         string  `json:"name"`
    CurrentLevel int     `json:"current_level"`
    Type         string  `json:"type"`
    Action       string  `json:"action"`
    UsePP        int     `json:"use_pp"`
    UseGold      int     `json:"use_gold"`
    Reward       *string `json:"reward"`
}

// Response (identisch)
type SkillCostResultNew struct {
    CharacterID    string `json:"character_id"`
    CharacterClass string `json:"character_class"`
    SkillName      string `json:"skill_name"`
    Category       string `json:"category"`
    Difficulty     string `json:"difficulty"`
    TargetLevel    int    `json:"target_level"`
    EP             int    `json:"ep"`
    GoldCost       int    `json:"gold_cost"`
    // ... weitere Felder
}
```

## Endpoints

```
POST /api/characters/lerncost      # Altes System (learningCostsData)
POST /api/characters/lerncost-new  # Neues System (Datenbank)
```

## Setup und Initialisierung

### 1. Lernkosten-System initialisieren

```bash
# CLI-Tool verwenden
cd backend
go run cmd/learning_costs_cli/main.go -init-learning

# Oder HTTP-Endpoint
curl -X POST http://localhost:8080/api/maintenance/initialize-learning-costs
```

### 2. Validierung

```bash
# CLI-Tool
go run cmd/learning_costs_cli/main.go -validate-learning

# Zusammenfassung anzeigen
go run cmd/learning_costs_cli/main.go -summary-learning
```

## Verwendung

### Beispiel-Request

```bash
curl -X POST http://localhost:8080/api/characters/lerncost-new \
  -H "Content-Type: application/json" \
  -d '{
    "char_id": 20,
    "name": "Athletik",
    "current_level": 9,
    "type": "skill",
    "action": "improve",
    "use_pp": 0,
    "use_gold": 0,
    "reward": "default"
  }'
```

### Beispiel-Response

```json
[
  {
    "character_id": "20",
    "character_class": "Kr",
    "skill_name": "Athletik",
    "category": "Körper",
    "difficulty": "normal",
    "target_level": 10,
    "ep": 40,
    "gold_cost": 40,
    "le": 0,
    "pp_used": 0,
    "gold_used": 0
  }
  // ... weitere Level bis 18
]
```

## Vorteile des neuen Systems

### 1. **Flexibilität**
- Lernkosten können zur Laufzeit geändert werden
- Neue Charakterklassen, Fertigkeiten und Schwierigkeitsgrade einfach hinzufügbar
- Verschiedene Regelbücher/Quellen können aktiviert/deaktiviert werden

### 2. **Verwaltbarkeit**
- Admin-Interface möglich
- Datenexport und -import
- Versioning und Änderungshistorie

### 3. **Skalierbarkeit**
- Unterstützt beliebig viele Fertigkeiten und Klassen
- Optimierte Datenbankabfragen
- Caching-Möglichkeiten

### 4. **Datenintegrität**
- Foreign Key Constraints
- Validierung auf Datenbankebene
- Transaktionale Konsistenz

## Migration vom alten System

### Schritt 1: Parallel-Betrieb
Beide Systeme können parallel laufen. Clients können zwischen den Endpoints wählen:
- `/lerncost` für das alte System
- `/lerncost-new` für das neue System

### Schritt 2: Testing und Validierung
```go
// Test beide Systeme und vergleiche Ergebnisse
func compareSystems(request LernCostRequest) {
    oldResponse := callOldSystem(request)
    newResponse := callNewSystem(request)
    
    // Vergleiche Ergebnisse und identifiziere Unterschiede
    compareResults(oldResponse, newResponse)
}
```

### Schritt 3: Umstellung
Nach erfolgreicher Validierung kann der alte Endpoint durch den neuen ersetzt werden.

## Fehlerbehebung

### Häufige Fehler

1. **"Fertigkeit nicht gefunden"**
   - Lernkosten-System nicht initialisiert
   - Fertigkeit existiert nicht in der Datenbank
   - Lösung: `InitializeLearningCostsSystem()` ausführen

2. **"Charakterklasse nicht verfügbar"**
   - Klasse nicht in `learning_character_classes` Tabelle
   - Lösung: Klasse zur Datenbank hinzufügen

3. **"EP-Kosten nicht gefunden"**
   - Fehlende Einträge in `learning_class_category_ep_costs`
   - Lösung: Daten für die spezifische Klasse/Kategorie-Kombination hinzufügen

### Debug-Informationen

```bash
# Prüfe Datenbank-Status
go run cmd/learning_costs_cli/main.go -summary-learning

# Validiere Daten
go run cmd/learning_costs_cli/main.go -validate-learning
```

## Technische Details

### Datenbank-Schema

Das neue System verwendet folgende Tabellen:
- `learning_sources` - Regelbücher/Quellen
- `learning_character_classes` - Charakterklassen
- `learning_skill_categories` - Fertigkeitskategorien
- `learning_skill_difficulties` - Schwierigkeitsgrade
- `learning_class_category_ep_costs` - EP-Kosten pro Klasse/Kategorie
- `learning_skill_category_difficulties` - Fertigkeiten-Zuordnungen
- `learning_skill_improvement_costs` - Verbesserungskosten
- `learning_spell_schools` - Zauberschulen
- `learning_class_spell_school_ep_costs` - Zauber-EP-Kosten
- `learning_spell_level_le_costs` - LE-Kosten pro Zaubergrad

### Performance-Optimierungen

- Indizierte Spalten für häufige Abfragen
- JOIN-optimierte Queries
- Möglichkeit für Query-Caching
- Prepared Statements
