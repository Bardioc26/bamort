# BaMoRT

Bamort ist ein Ersatzwerkzeug für MOAM mit dem Schwerpunkt auf Charakterverwaltung für Rollenspiele. Das Projekt bietet ein einfaches, erweiterbares System zur Charaktererstellung, zum Erlernen neuer Fertigkeiten und zur Charakterentwicklung. Langfristiges Ziel ist es, die Regeln für Charaktererstellung und -weiterentwicklung aus dem Code zu extrahieren und ein austauschbares Regelwerk‑Framework bereitzustellen.

## Ziele
- Ein modernes, wartbares Ersatzwerkzeug für MOAM mit Fokus auf Charaktere bereitstellen.
- Primäre Schwerpunkte:
  - Charaktererstellung (Attribute, abgeleitete Werte, Startfertigkeiten).
  - Charakterentwicklung (neue Fertigkeiten erlernen, bestehende Fertigkeiten verbessern).
- Langfristig: Regeln für Generierung und Fortschritt im Code abbilden und ein Framework anbieten, mit dem verschiedene Regelwerke implementiert und gewechselt werden können.

## Architekturüberblick
Monorepo mit zwei Hauptteilen: Backend und Frontend.

- Backend
  - Sprache: Go (1.25+)
  - Frameworks: Gin (HTTP), GORM (ORM)
  - Zuständigkeiten: API, Geschäftslogik, Regelimplementierung, PDF‑Export
  - Typische Struktur: Module pro Domäne (character, pdfrender, equipment, ...)
  - Einstiegspunkt: cmd/main.go

- Frontend
  - Framework: Vue 3 + Vite
  - State: Pinia Stores
  - i18n: Lokalisierungsdateien unter src/locales/
  - Zuständigkeiten: UI für Charaktere (Erstellung, Bearbeitung, Entwicklung), Admin‑Werkzeuge

- Daten & DevOps
  - Datenbank: MariaDB für Entwicklung/Produktion; SQLite für Tests
  - Container: Docker Compose mit separaten Dev‑Containern für Backend, Frontend und DB
  - PDF‑Rendering: chromedp im pdfrender‑Modul (Chromium im Container erforderlich)

## Kurze Entwicklerhinweise
- Geschützte API‑Routen liegen unter `/api` und erfordern JWT‑Auth.
- Vorlagen für Charakterexport liegen in `templates/` und sind per ENV‑Var TEMPLATES_DIR konfigurierbar.
- Dev‑Docker‑Compose für Live‑Reload verwenden:
  - Backend: `bamort-backend-dev` (Air, Port 8180)
  - Frontend: `bamort-frontend-dev` (Vite, Port 5173)
  - DB: `bamort-mariadb-dev` (MariaDB)

## Tests
- Backend‑Tests nutzen `testutils.SetupTestDB()` und die Testumgebung (`ENVIRONMENT=test`).
- NIEMALS `main()` in Testdateien hinzufügen; Tests müssen auf `_test.go` enden.
- Beispielbefehle:
  - cd backend && go test -v ./character/
  - cd backend && go test -v ./pdfrender/ -run TestExportCharacterToPDF

## Plan für erweiterbare Regelwerke
- Kernidee: Regeln für Generierung und Fortschritt aus unstrukturierter Logik in modulare, versionierte Regelwerk‑Pakete überführen.
- Regelwerke sollen:
  - Validierungs‑ und Generierungsschritte für Charaktere definieren.
  - Hooks für Lern‑ und Fortschrittsalgorithmen von Fertigkeiten bereitstellen.
  - Laufzeit‑wahl ermöglichen (Konfiguration oder pro Kampagne).
- Anfangs werden kanonische Regeln hartkodiert; spätere Refactorings extrahieren Schnittstellen und liefern Beispiel‑Regelwerke.

## Mitwirkung
- KISS und TDD: zuerst fehlschlagende Tests schreiben, dann minimale Implementierung.
- Neue UI‑Strings in beiden `src/locales/de` und `src/locales/en` ergänzen.
- Änderungen modular pro Domänenmodul halten (Backend‑Modulkonvention beachten).

## Wo anfangen
- backend/character — Charakter‑Handler und Routen
- backend/models  — Datenmodelle
- backend/pdfrender — PDF‑Templates und Rendering‑Logik
- frontend/src/views — Haupt‑UI für Charakterarbeiten

## Lizenz
Dieses Projekt ist unter der **PolyForm Noncommercial License 1.0.0** lizenziert (siehe `LICENSE`).

Kurzfassung:
- Nicht-kommerzielle Nutzung ist erlaubt.
- Für kommerzielle Nutzung ist eine separate kommerzielle Lizenz nötig; siehe `LICENSE-COMMERCIAL.md`.

Commercial licensing: contact **(Bamort Admin / Bamort@trokan.de)**.

# BaMoRT for english readers

Bamort is a replacement tool for MOAM focused on character management for role‑playing games. The project aims to provide a simple, extensible system for character creation, learning new skills, and character development, with a long‑term goal of extracting the rules for character generation and advancement from code and exposing a pluggable ruleset framework.

## Goals
- Provide a modern, maintainable replacement for MOAM focused on characters.
- First-class focus on:
  - Character creation (attributes, derived values, starting skills).
  - Character development (learning new skills, improving existing skills).
- Long-term: encode generation & progression rules in code and offer a framework so different rulesets can be implemented and swapped.

## Architecture overview
Monorepo with two primary parts: backend and frontend.

- Backend
  - Language: Go (1.25+)
  - Frameworks: Gin (HTTP), GORM (ORM)
  - Responsibilities: API, business logic, rules implementation, PDF export
  - Typical layout: modules per domain (character, pdfrender, equipment, ...)
  - Entry point: cmd/main.go

- Frontend
  - Framework: Vue 3 + Vite
  - State: Pinia stores
  - i18n: locale files under src/locales/
  - Responsibilities: character UI (creation, edit, development), admin tools

- Data & Devops
  - Database: MariaDB for dev/prod; SQLite used for tests
  - Containers: Docker Compose with separate dev containers for backend, frontend, and DB
  - PDF rendering: chromedp in the pdfrender module (requires Chromium in container)

## Quick dev notes
- Protected API routes live under `/api` and require JWT auth.
- Templates for character export are stored under `templates/` and the directory can be configured via env var TEMPLATES_DIR.
- Use the included docker-compose dev setup for live reload:
  - Backend: `bamort-backend-dev` (Air, port 8180)
  - Frontend: `bamort-frontend-dev` (Vite, port 5173)
  - DB: `bamort-mariadb-dev` (MariaDB)

## Testing
- Backend tests use `testutils.SetupTestDB()` and a test environment (`ENVIRONMENT=test`).
- NEVER add a `main()` to test files; tests must end with `_test.go`.
- Example commands:
  - cd backend && go test -v ./character/
  - cd backend && go test -v ./pdfrender/ -run TestExportCharacterToPDF

## Extensible ruleset plan
- Core idea: move rules for generation and progression from ad‑hoc code into modular, versioned ruleset packages.
- Rulesets should:
  - Define validation and generation steps for a character.
  - Expose hooks for skill learning and progression algorithms.
  - Be selectable at runtime (config or per‑campaign).
- Early work will hardcode canonical rules; later refactors will extract interfaces and provide example alternate rulesets.

## Contributing
- Follow KISS and TDD: write failing tests first, implement minimal change.
- Add translations to both `src/locales/de` and `src/locales/en` for new UI strings.
- Keep changes modular per domain module (see backend module pattern).

## Where to look first
- backend/character — character handlers and routes
- backend/models  — data models
- backend/pdfrender — PDF template system and rendering logic handlers and routes
- frontend/src/views — primary UI entry points for character workflows

## License
This project is licensed under the **PolyForm Noncommercial License 1.0.0** (see `LICENSE`).

Summary:
- Noncommercial use is permitted.
- Commercial use requires a separate commercial license; see `LICENSE-COMMERCIAL.md`.

Commercial licensing: contact **(Bamort Admin / Bamort@trokan.de)**.