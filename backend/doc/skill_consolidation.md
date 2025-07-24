# Skill-Kategorie und Schwierigkeit Konsolidierung

## Problemanalyse
1. Skills in der JSON-Datei haben keine Category/Difficulty
2. Fallback-Funktionen sind inkonsistent mit den Lerntabellen
3. Manche Kategorien/Schwierigkeiten existieren nicht in den BaseLearnCost

## Konsolidierte Skill-Definitionen (basierend auf Midgard KOD5)

### Alltag-Fertigkeiten
- **Erste Hilfe**: normal (1 LE)
- **Wahrnehmung**: leicht (1 LE) 
- **Hören**: leicht (1 LE)
- **Sehen**: leicht (1 LE)
- **Nachtsicht**: leicht (1 LE)
- **Riechen**: leicht (1 LE)
- **Sechster Sinn**: leicht (1 LE)

### Freiland-Fertigkeiten  
- **Geländelauf**: normal (2 LE)

### Körper-Fertigkeiten
- **Athletik**: normal (1 LE)
- **Klettern**: schwer (2 LE) 
- **Laufen**: normal (1 LE)
- **Robustheit**: normal (1 LE)

### Kampf-Fertigkeiten
- **Kampf in Vollrüstung**: schwer (10 LE)

### Wissen-Fertigkeiten
- **Landeskunde**: normal (2 LE)
- **Sprache**: normal (2 LE)

### Waffen-Fertigkeiten
- **Stichwaffen**: normal (2 LE)
- **Einhandschlagwaffen**: normal (2 LE)
- **Zweihandschlagwaffen**: normal (2 LE)
- **Spießwaffen**: normal (2 LE)
- **Waffenloser Kampf**: normal (2 LE)
- **Armbrüste**: normal (2 LE)
- **Stielwurfwaffen**: normal (2 LE)
- **Schilde**: leicht (1 LE)

## Fehlende Kategorien in BaseLearnCost
- Wissen: leicht fehlt (sollte 1 LE sein)
- Waffen: sehr_schwer fehlt (wird normalerweise nicht verwendet)

## Empfohlene Korrekturen
1. BaseLearnCost um fehlende Einträge erweitern
2. Skill-Mapping korrigieren
3. Konsistente Kategoriezuordnung
