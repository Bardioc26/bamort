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

### MariaDB (Datenbankserver)
- **Port**: 3306 (localhost:3306)
- **Datenbank**: bamort
- **Benutzer**: bamort
- **Passwort**: bG4)efozrc (Development)
- **Root-Passwort**: root_password_dev
- **Persistent**: ✅ Daten bleiben bei Container-Neustarts erhalten
- **Health Check**: ✅ Backend wartet auf Datenbankbereitschaft

### phpMyAdmin (Datenbank-Management)
- **URL**: http://localhost:8082
- **Server**: mariadb (automatisch konfiguriert)
- **Login**: root / root_password_dev (oder bamort / bG4)efozrc)
- **Features**: ✅ Web-basiertes Datenbank-Management
- **Auto-Login**: ✅ Vorkonfiguriert für MariaDB-Container

### Backend (Go mit Air Live-Reloading)
- **URL**: http://localhost:8180
- **API-Dokumentation**: http://localhost:8180/api (falls verfügbar)
- **Datenbankverbindung**: Automatisch konfiguriert zu MariaDB-Container
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
# MariaDB
docker exec -it bamort-mariadb-dev mysql -u bamort -p bamort

# phpMyAdmin (Web-Interface)
# Zugriff über Browser: http://localhost:8082

# Backend
docker exec -it bamort-backend-dev sh

# Frontend
docker exec -it bamort-frontend-dev sh
```

## Datenbank-Management

### Web-Interface (phpMyAdmin)
- **URL**: http://localhost:8082
- **Login**: 
  - Root: `root` / `root_password_dev`
  - User: `bamort` / `bG4)efozrc`
- **Features**: Vollständige Datenbankverwaltung über Web-Interface

### Kommandozeilen-Zugriff
```bash
# Mit mysql client im Container
docker exec -it bamort-mariadb-dev mysql -u bamort -p bamort

# Oder als root
docker exec -it bamort-mariadb-dev mysql -u root -p
```

### Backup/Restore
```bash
# Backup erstellen
docker exec bamort-mariadb-dev mysqldump -u bamort -pbG4)efozrc bamort > backup.sql

# Backup einspielen
docker exec -i bamort-mariadb-dev mysql -u bamort -pbG4)efozrc bamort < backup.sql
```

### Datenbank zurücksetzen
```bash
# Volumes löschen (alle Daten gehen verloren!)
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up --build
```

## Konfiguration

### MariaDB
- Version: 11.4
- Charset: utf8mb4
- Collation: utf8mb4_unicode_ci
- Initialisierungsskripte: `docker/init-db/`

### Backend
- Environment: `development`
- Datenbanktyp: mysql (MariaDB)
- Datenbankverbindung: Automatisch konfiguriert
- Air-Konfiguration: `backend/.air.toml`
- Ausgeschlossene Dateien: Tests, Uploads, temporäre Dateien

### Frontend
- Environment: `development`
- Vite dev server mit `--host 0.0.0.0` für Docker-Zugriff
- API URL: `http://localhost:8080` (konfigurierbar über VITE_API_URL)

## Troubleshooting

### Port bereits in Verwendung
Stelle sicher, dass die Ports 3306, 8080, 8082, 8180 und 5173 nicht von anderen Anwendungen verwendet werden.

### Datenbankverbindungsfehler
- Warte nach dem Start 10-15 Sekunden, bis MariaDB vollständig initialisiert ist
- Überprüfe die Logs: `docker-compose -f docker-compose.dev.yml logs mariadb`

### phpMyAdmin nicht erreichbar
- Stelle sicher, dass Port 8082 frei ist
- Warte bis MariaDB vollständig gestartet ist
- Überprüfe die Logs: `docker-compose -f docker-compose.dev.yml logs phpmyadmin`

### Node_modules Probleme
Falls es Probleme mit node_modules gibt:
```bash
docker-compose -f docker-compose.dev.yml down -v
docker-compose -f docker-compose.dev.yml up --build
```

### Go Modules Cache
Der Go-Module-Cache wird in einem Volume gespeichert, um Downloads zu beschleunigen.
