# Bamort Development Environment

Diese Docker-Konfiguration ermöglicht es, die Bamort-Anwendung in einer Entwicklungsumgebung zu starten, in der Änderungen am Code sofort verfügbar sind.

## Voraussetzungen

- Docker und Docker Compose installiert
- Ports 8080 (Backend) und 5173 (Frontend) verfügbar

## Starten der Entwicklungsumgebung

```bash
# Aus dem docker/ Verzeichnis
./start-dev.sh

# Oder manuell
docker-compose -f docker-compose.dev.yml up --build
```

## Stoppen der Entwicklungsumgebung

```bash
# Aus dem docker/ Verzeichnis
./stop-dev.sh

# Oder manuell
docker-compose -f docker-compose.dev.yml down
```

## Verfügbare Services

### Backend (Go mit Air Live-Reloading)
- **URL**: http://localhost:8080
- **API-Dokumentation**: http://localhost:8080/api (falls verfügbar)
- **Live-Reloading**: ✅ Änderungen an Go-Dateien lösen automatisch einen Neustart aus
- **Volume**: Backend-Code wird vom lokalen Verzeichnis gemountet

### Frontend (Vue.js mit Vite Dev Server)
- **URL**: http://localhost:5173
- **Live-Reloading**: ✅ Änderungen an Vue-Dateien werden sofort übernommen
- **Hot Module Replacement**: ✅ Aktiviert für schnelle Entwicklung
- **Volume**: Frontend-Code wird vom lokalen Verzeichnis gemountet

## Entwicklung

### Backend-Entwicklung
- Alle Änderungen an `.go` Dateien im `backend/` Verzeichnis werden automatisch erkannt
- Der Server wird automatisch neu gestartet bei Änderungen
- Logs sind in der Docker-Konsole sichtbar

### Frontend-Entwicklung
- Alle Änderungen an Vue-Komponenten, CSS und anderen Assets werden sofort übernommen
- Vite's Hot Module Replacement sorgt für schnelle Updates ohne vollständigen Page-Reload
- Build-Fehler werden im Browser angezeigt

## Debugging

### Backend-Logs anzeigen
```bash
docker-compose -f docker-compose.dev.yml logs backend-dev
```

### Frontend-Logs anzeigen
```bash
docker-compose -f docker-compose.dev.yml logs frontend-dev
```

### In Container einsteigen
```bash
# Backend
docker exec -it bamort-backend-dev sh

# Frontend
docker exec -it bamort-frontend-dev sh
```

## Konfiguration

### Backend
- Environment: `development`
- Air-Konfiguration: `backend/.air.toml`
- Ausgeschlossene Dateien: Tests, Uploads, temporäre Dateien

### Frontend
- Environment: `development`
- Vite dev server mit `--host 0.0.0.0` für Docker-Zugriff
- API URL: `http://localhost:8080` (konfigurierbar über VITE_API_URL)

## Troubleshooting

### Port bereits in Verwendung
Stelle sicher, dass die Ports 8080 und 5173 nicht von anderen Anwendungen verwendet werden.

### Node_modules Probleme
Falls es Probleme mit node_modules gibt:
```bash
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up --build
```

### Go Modules Cache
Der Go-Module-Cache wird in einem Volume gespeichert, um Downloads zu beschleunigen.
