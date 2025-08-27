# GetAllSpellsWithLE Funktion

## Überblick

Die `GetAllSpellsWithLE` Funktion wurde implementiert, um alle verfügbaren Zauber mit ihren Lerneinheiten (LE) und Erfahrungspunkten (EP) gruppiert nach Kategorien aus der Datenbank abzurufen.

## Funktionssignatur

```go
func GetAllSpellsWithLE(characterClass string, maxLevel int) (map[string][]gin.H, error)
```

### Parameter

- `characterClass`: Die Charakterklasse (z.B. "Magier", "Hexer", "Druide")
- `maxLevel`: Maximale Zauberstufe, die abgerufen werden soll

### Rückgabe

- `map[string][]gin.H`: Eine Map mit Zauberschulen als Schlüssel und Arrays von Zaubern als Werte
- `error`: Fehlermeldung falls etwas schiefgeht

## Features

### 1. Charakterklassen-Zauberschulen-Mapping

Die Funktion verwendet ein vordefiniertes Mapping, das festlegt, welche Zauberschulen eine Charakterklasse erlernen kann:

```go
"Magier": {
    "Beherrschen": true,
    "Bewegen":     true,
    "Dweomer":     true,
    "Erkennen":    true,
    "Erschaffen":  true,
    "Verändern":   true,
    "Zerstören":   true,
},
"Hexer": {
    "Beherrschen": true,
    "Dweomer":     true,
    "Erkennen":    true,
    "Verändern":   true,
},
```

### 2. Lerneinheiten-Berechnung

Die LE-Kosten werden aus der Datenbank-Tabelle `learning_spell_level_le_costs` abgerufen, die die Zuordnung von Zauberstufen zu LE-Kosten enthält.

**Datenbankstruktur für SpellLevelLECost:**
```go
type SpellLevelLECost struct {
    ID         uint   `gorm:"primaryKey" json:"id"`
    Level      int    `gorm:"uniqueIndex;not null" json:"level"`
    LERequired int    `gorm:"not null" json:"le_required"`
    GameSystem string `gorm:"index;default:midgard" json:"game_system"`
}
```

**Fallback bei fehlenden Datenbankeinträgen:**
Falls keine Einträge in der Datenbank gefunden werden, verwendet die Funktion Standard Midgard-Werte:

```go
Stufe 1:  1 LE    Stufe 7:  8 LE
Stufe 2:  2 LE    Stufe 8:  10 LE
Stufe 3:  3 LE    Stufe 9:  12 LE
Stufe 4:  4 LE    Stufe 10: 15 LE
Stufe 5:  5 LE    Stufe 11: 18 LE
Stufe 6:  6 LE    Stufe 12: 21 LE
```

### 3. EP-Kostenberechnung

EP-Kosten werden basierend auf Charakterklasse und Zauberschule berechnet:

- **Magier**: 10 EP pro LE für alle Schulen
- **Hexer**: 30 EP pro LE für erlaubte Schulen
- **Druide**: 20 EP pro LE für erlaubte Schulen
- **Schamane**: 40 EP pro LE für erlaubte Schulen
- **Priester**: 30 EP pro LE für erlaubte Schulen
- **Barde**: 60 EP pro LE für erlaubte Schulen
- **Ordenskrieger**: 90 EP pro LE für erlaubte Schulen

## Beispiel-Nutzung

```go
// Alle Zauber bis Stufe 3 für einen Hexer
spells, err := GetAllSpellsWithLE("Hexer", 3)
if err != nil {
    log.Fatal(err)
}

// Ausgabe: Zauber gruppiert nach Schulen (Beherrschen, Verändern, etc.)
```

## Beispiel-Ausgabe für Hexer (maxLevel: 3)

```json
{
  "Beherrschen": [
    {
      "description": "Hebt andere Zauber auf",
      "ep_cost": 30,
      "id": 1,
      "le_cost": 1,
      "learning_category": "Beherrschen",
      "level": 1,
      "name": "Bannen von Zauberwerk",
      "school": "Beherrschen"
    },
    {
      "description": "Beruhigt den Geist",
      "ep_cost": 60,
      "id": 2,
      "le_cost": 2,
      "learning_category": "Beherrschen",
      "level": 2,
      "name": "Geistesruhe",
      "school": "Beherrschen"
    }
  ],
  "Verändern": [
    {
      "description": "Heilt körperliche Verletzungen",
      "ep_cost": 30,
      "id": 6,
      "le_cost": 1,
      "learning_category": "Verändern",
      "level": 1,
      "name": "Heilen von Wunden",
      "school": "Verändern"
    }
  ]
}
```

## Datenbankstruktur

Die Funktion erwartet folgende Felder in der `gsm_spells` Tabelle:

- `id`: Eindeutige ID des Zaubers
- `name`: Name des Zaubers
- `beschreibung`: Beschreibung des Zaubers
- `stufe`: Zauberstufe (1-12)
- `category`: Anzeige-Kategorie des Zaubers
- `learning_category`: Interne Lernkategorie für EP-Berechnung
- Zusätzliche Felder: `ap`, `art`, `zauberdauer`, `reichweite`, `wirkungsziel`, `wirkungsbereich`, `wirkungsdauer`

## Beispiel-Testdaten

Siehe `example_spell_data.sql` für vollständige Beispieldaten verschiedener Zauberschulen und Stufen.

## Sortierung

Zauber werden innerhalb jeder Kategorie zuerst nach Stufe, dann alphabetisch nach Name sortiert.

## Unterstützte Charakterklassen

- Magier (Ma) - Vollzugriff auf alle Schulen
- Hexer (Hx) - Beherrschen, Dweomer, Erkennen, Verändern  
- Druide (Dr) - Bewegen, Erkennen, Erschaffen, Verändern
- Schamane (Sc) - Beherrschen, Erkennen, Verändern
- Priester Beschützer (PB) - Dweomer, Erkennen, Verändern
- Priester Streiter (PS) - Dweomer, Erkennen, Verändern
- Barde (Ba) - Beherrschen, Dweomer, Erkennen
- Ordenskrieger (Or) - Dweomer, Erkennen

Charakterklassen, die nicht in der Mapping-Tabelle stehen, erhalten eine leere Antwort (können keine Zauber lernen).
