# Character Creation System

## Übersicht

Das Character Creation System ermöglicht es Benutzern, neue Charaktere in einem mehrstufigen Prozess zu erstellen. Der Fortschritt wird nach jedem Schritt gespeichert und steht dem Benutzer für 14 Tage zur Verfügung.

## Implementierte Features

### ✅ Backend (Go/Gin)

**Neue API-Endpunkte in `/backend/character/routes.go`:**
- `GET /api/characters/create-sessions` - Aktive Sessions für Benutzer auflisten
- `POST /api/characters/create-session` - Neue Session erstellen
- `GET /api/characters/create-session/:sessionId` - Session-Daten abrufen  
- `PUT /api/characters/create-session/:sessionId/basic` - Grundinformationen
- `PUT /api/characters/create-session/:sessionId/attributes` - Grundwerte
- `PUT /api/characters/create-session/:sessionId/derived` - Abgeleitete Werte
- `PUT /api/characters/create-session/:sessionId/skills` - Fertigkeiten
- `PUT /api/characters/create-session/:sessionId/spells` - Zauber
- `POST /api/characters/create-session/:sessionId/finalize` - Abschließen
- `DELETE /api/characters/create-session/:sessionId` - Session löschen
- `GET /api/characters/races` - Verfügbare Rassen
- `GET /api/characters/classes` - Verfügbare Klassen
- `GET /api/characters/origins` - Verfügbare Herkünfte
- `GET /api/characters/beliefs?q=searchterm` - Glaube-Suche
- `GET /api/characters/skill-categories` - Kategorien mit Lernpunkten

**Handler in `/backend/character/handlers.go`:**
- Session-Management mit 14-Tage Gültigkeit
- Reference Data APIs mit Dummy-Daten
- Step-by-Step Validierung und Speicherung

**Datenmodell in `/backend/models/model_character_creation.go`:**
- CharacterCreationSession mit JSON-Feldern
- Automatische Session-Bereinigung
- Type-safe Datenstrukturen

### ✅ Frontend (Vue 3)

**Button in CharacterList.vue:**
- "Create New Character" Button mit Styling
- Session-Erstellung und Navigation
- **Aktive Sessions-Anzeige**: Grid mit Session-Karten
- **Session-Details**: Name, Fortschritt, Rasse/Klasse, Daten
- **Continue/Delete**: Fortsetzen oder Löschen von Drafts

**Neue Komponenten:**
- `CharacterCreation.vue` - Hauptcontainer mit Progress-Indicator
- `CharacterCreation/CharacterBasicInfo.vue` - Schritt 1: Grunddaten
- `CharacterCreation/CharacterAttributes.vue` - Schritt 2: Grundwerte
- `CharacterCreation/CharacterDerivedValues.vue` - Schritt 3: Abgeleitete Werte
- `CharacterCreation/CharacterSkills.vue` - Schritt 4: Fertigkeiten
- `CharacterCreation/CharacterSpells.vue` - Schritt 5: Zauber

**Router-Integration:**
- Route `/character/create/:sessionId` hinzugefügt
- Props-basierte Parameter-Übergabe

## Funktionsdetails

### Schritt 1: Grundinformationen
- **Name**: Text-Input mit Validierung (2-50 Zeichen)
- **Rasse**: Dropdown aus API-Daten
- **Klasse**: Dropdown aus API-Daten  
- **Herkunft**: Dropdown aus API-Daten
- **Glaube**: Suchfeld mit dynamischem Autocomplete (min. 2 Zeichen)

### Schritt 2: Grundwerte
- **9 Attribute**: ST, GS, GW, KO, IN, ZT, AU, PA, WK (1-100)
- **Live-Berechnung**: Gesamtpunkte und Durchschnitt
- **Responsive Grid**: Automatisches Layout für verschiedene Bildschirmgrößen

### Schritt 3: Abgeleitete Werte
- **LP/AP/B**: Automatische Berechnung aus Grundwerten
- **Bennies**: SG, GG, GP basierend auf Klasse
- **Recalculate-Button**: Reset auf berechnete Werte
- **Manuelle Anpassung**: Überschreibung der Berechnungen möglich

### Schritt 4: Fertigkeiten
- **Kategorien-Ansicht**: Grid mit Lernpunkt-Anzeige und Progress-Bars
- **Skill-Listen**: Dynamische Listen pro Kategorie aus Datenbank
- **Punkt-Tracking**: Echtzeit-Update der verbleibenden Lernpunkte
- **Selected-Overview**: Zusammenfassung aller gewählten Skills
- **Navigation**: "Next: Spells" für nächsten Schritt

### Schritt 5: Zauber
- **Zauber-Liste**: Separate Kategorie für magische Fähigkeiten
- **Zauber-Punkte**: Eigenständiges Punktesystem für Magie
- **Visual Feedback**: Kann/kann nicht lernen Indikationen
- **Finalisierung**: "Create Character" Button für Abschluss

## Session-Management

### Persistierung
- **Session-ID**: UUID-basiert für Eindeutigkeit
- **Auto-Save**: Nach jedem Schritt automatische Speicherung
- **14-Tage Gültigkeit**: Automatische Bereinigung alter Sessions
- **Step-Tracking**: Fortschritts-Verfolgung über alle Schritte

### Benutzerführung
- **Progress-Indicator**: Visuelle Fortschrittsanzeige (5 Schritte) - **CLICKABLE!**
- **Navigation**: Vor/Zurück zwischen Schritten über klickbare Progress-Steps
- **Validierung**: Schritt-Blockierung bei fehlenden Daten
- **Draft-Management**: Session-Liste mit Continue/Delete-Optionen
- **Session-Übersicht**: Aktive Drafts auf Character-Liste-Seite
- **Datenerhaltung**: Alle eingegebenen Daten bleiben beim Wechseln erhalten

## Dummy-Implementierungen

**Backend Placeholders:**
- Reference Data als hardcodierte Arrays
- Session-Speicherung in Memory (TODO: Redis/DB)
- User-ID als Konstante (TODO: Auth-Context)
- Skill/Spell-Kosten als Fallback-Werte

**Frontend Fallbacks:**
- API-Fehler-Behandlung mit Sample-Daten
- Loading-States für bessere UX
- Offline-fähige Formulare

## Nächste Schritte

### 🎯 Produktive Implementierung
1. **Session-Storage**: Redis oder DB-Tabelle
2. **Reference Data**: Echte Datenbank-Anbindung
3. **Auth-Integration**: User-ID aus JWT/Session
4. **Validation**: Server-seitige Eingabe-Prüfung

### 🔧 Optimierungen  
1. **Performance**: Reference Data Caching
2. **UX**: Loading-Spinner und Transitions
3. **Error Handling**: Robuste Fehlerbehandlung
4. **Tests**: Unit- und E2E-Tests

Die Implementierung ist funktional vollständig und kann sofort getestet werden. Alle Dummy-Implementierungen sind klar markiert für schrittweise Produktionsreife.
