# Normalisierte Datenbankstruktur für Lernkosten

## Übersicht

Die `LearningCostsTable2` Strukt### 5. **Erweiterte Migration**
- Erstellt Standardquellen automatisch (KOD, ARK, MYS, UNB)
- Verknüpft Charakterklassen mit KOD
- Verknüpft Fertigkeiten mit KOD (Standardquelle)
- Verknüpft Zauber mit ARK (Standardquelle)
- Verknüpft Basis-Zauberschulen mit KOD, erweiterte mit ARK
- "Unbekannt"-Kategorie wird mit UNB-Quelle verknüpfturde in eine normalisierte Datenbankstruktur transformiert, die es ermöglicht, die Lernkosten-Daten effizient in der Datenbank zu verwalten.

## Neue Tabellen

### 1. Basistabellen

#### `learning_character_classes`
Speichert alle Charakterklassen mit Code und vollständigem Namen.
- `id` (PK)
- `code` (String, 3 Zeichen, z.B. "Hx", "Ma", "Kr")
- `name` (String, z.B. "Hexer", "Magier", "Krieger")
- `description` (Optional)
- `game_system` (Standard: "midgard")

#### `learning_skill_categories`
Speichert alle Fertigkeitskategorien.
- `id` (PK)
- `name` (String, z.B. "Alltag", "Kampf", "Waffen")
- `description` (Optional)
- `game_system` (Standard: "midgard")

#### `learning_skill_difficulties`
Speichert alle Schwierigkeitsgrade.
- `id` (PK)
- `name` (String, z.B. "leicht", "normal", "schwer", "sehr schwer")
- `description` (Optional)
- `game_system` (Standard: "midgard")

#### `learning_spell_schools`
Speichert alle Zauberschulen.
- `id` (PK)
- `name` (String, z.B. "Beherrschen", "Bewegen", "Erkennen")
- `description` (Optional)
- `game_system` (Standard: "midgard")

### 2. Kostentabellen

#### `learning_class_category_ep_costs`
EP-Kosten für 1 Trainingseinheit (TE) pro Charakterklasse und Fertigkeitskategorie.
- `id` (PK)
- `character_class_id` (FK zu learning_character_classes)
- `skill_category_id` (FK zu learning_skill_categories)
- `ep_per_te` (Integer, EP-Kosten pro TE)

#### `learning_class_spell_school_ep_costs`
EP-Kosten für 1 Lerneinheit (LE) für Zauber pro Charakterklasse und Zauberschule.
- `id` (PK)
- `character_class_id` (FK zu learning_character_classes)
- `spell_school_id` (FK zu learning_spell_schools)
- `ep_per_le` (Integer, EP-Kosten pro LE)

#### `learning_spell_level_le_costs`
LE-Kosten pro Zauber-Stufe.
- `id` (PK)
- `level` (Integer, Zauber-Stufe 1-12)
- `le_required` (Integer, benötigte Lerneinheiten)
- `game_system` (Standard: "midgard")

### 3. Verknüpfungstabellen

#### `learning_skill_category_difficulties`
Definiert die Schwierigkeit einer Fertigkeit in einer bestimmten Kategorie.
Eine Fertigkeit kann in mehreren Kategorien mit unterschiedlichen Schwierigkeiten existieren.
- `id` (PK)
- `skill_id` (FK zu lookup_lists/skills)
- `skill_category_id` (FK zu learning_skill_categories)
- `skill_difficulty_id` (FK zu learning_skill_difficulties)
- `learn_cost` (Integer, LE-Kosten für das Erlernen)

#### `learning_skill_improvement_costs`
TE-Kosten für Verbesserungen basierend auf Kategorie, Schwierigkeit und aktuellem Wert.
- `id` (PK)
- `skill_category_difficulty_id` (FK zu learning_skill_category_difficulties)
- `current_level` (Integer, aktueller Fertigkeitswert)
- `te_required` (Integer, benötigte Trainingseinheiten)

## Vorteile der neuen Struktur

### 1. Normalisierung
- Eliminierung von Redundanzen
- Konsistente Datenstruktur
- Einfache Wartung und Updates

### 2. Flexibilität
- Fertigkeiten können in mehreren Kategorien mit unterschiedlichen Schwierigkeiten existieren
- Neue Charakterklassen, Kategorien oder Schwierigkeiten können einfach hinzugefügt werden
- Lernkosten können individuell angepasst werden

### 3. Performance
- Effiziente Abfragen durch indexierte Fremdschlüssel
- Optimierte Joins für komplexe Kostenberechnungen

### 4. Erweiterbarkeit
- Einfache Integration neuer Spielsysteme über `game_system`
- Möglichkeit zur Versionierung von Lernkosten
- Unterstützung für charakterspezifische Modifikatoren

## Migration und Integration

### Migrationsfunktionen

1. **`CreateLearningCostsTables()`**
   - Erstellt alle neuen Tabellen
   - Definiert Indizes und Fremdschlüssel

2. **`MigrateLearningCostsToDatabase()`**
   - Überträgt alle Daten aus `learningCostsData` in die Datenbank
   - Erstellt Basisdaten für alle Charakterklassen, Kategorien, etc.

3. **`EnhanceLearningDataWithExistingTables()`**
   - Verknüpft vorhandene Fertigkeiten und Zauber mit den neuen Tabellen
   - Erstellt fehlende Verknüpfungen mit Standardwerten

4. **`InitializeLearningCostsSystem()`**
   - Hauptfunktion für die komplette Migration
   - Führt alle Schritte in der richtigen Reihenfolge aus

### Query-Funktionen

1. **`GetEPPerTEForClassAndCategory()`**
   - Holt EP-Kosten für eine Klasse und Kategorie

2. **`GetSkillCategoryAndDifficulty()`**
   - Findet die beste Kategorie für eine Fertigkeit (niedrigste EP-Kosten)

3. **`GetSpellLearningInfo()`**
   - Holt alle Informationen für das Erlernen eines Zaubers

4. **`GetImprovementCost()`**
   - Holt Verbesserungskosten für eine spezifische Situation

## Verwendung

```go
// System initialisieren (einmalig)
if err := InitializeLearningCostsSystem(); err != nil {
    log.Fatal("Failed to initialize learning costs system:", err)
}

// Kostenberechnung für Fertigkeit
skillInfo, err := GetSkillCategoryAndDifficulty("Klettern", "Kr")
if err != nil {
    log.Error("Failed to get skill info:", err)
}

// Zauber-Lernkosten ermitteln
spellInfo, err := GetSpellLearningInfo("Feuerball", "Ma")
if err != nil {
    log.Error("Failed to get spell info:", err)
}
```

## Kompatibilität

Die neue Struktur ist vollständig kompatibel mit dem bestehenden System:
- Alle bestehenden Funktionen wie `CalcSkillLernCost()` können weiterhin verwendet werden
- Die neuen Query-Funktionen ersetzen schrittweise die statischen Maps
- Keine Änderungen an bestehenden APIs erforderlich
