# Bamort Logging System

## Überblick

Das Bamort Logging System bietet konfigurierbare Debug- und Standard-Logs, die über Umgebungsvariablen oder .env-Dateien gesteuert werden können.

## Konfiguration

### .env-Dateien

Das System lädt automatisch `.env` und `.env.local` Dateien aus dem Projektverzeichnis, falls vorhanden.

**Priorität der Konfigurationsquellen:**
1. Bereits gesetzte Umgebungsvariablen (höchste Priorität)
2. .env.local Datei
3. .env Datei
4. Standard-Werte (niedrigste Priorität)

**Beispiel .env-Datei:**
```bash
# Bamort Server Konfiguration
ENVIRONMENT=development
DEBUG=true
LOG_LEVEL=DEBUG
PORT=8180
DATABASE_URL=postgresql://user:pass@localhost:5432/bamort
```

### Umgebungsvariablen

### Debug-Modus
```bash
# Debug-Modus aktivieren
DEBUG=true

# oder
DEBUG=1
```

### Log-Level
```bash
# Verfügbare Log-Level: DEBUG, INFO, WARN, ERROR
LOG_LEVEL=DEBUG
LOG_LEVEL=INFO
LOG_LEVEL=WARN
LOG_LEVEL=ERROR
```

### Environment
```bash
# Umgebung definieren
ENVIRONMENT=development  # Aktiviert automatisch Debug-Modus
ENVIRONMENT=production   # Deaktiviert Debug-Modus standardmäßig

# oder
GO_ENV=development
GO_ENV=production
```

### Server-Port
```bash
# Server-Port konfigurieren
PORT=8180
# oder
SERVER_PORT=8180
```

## Verwendung im Code

### Import
```go
import "bamort/logger"
```

### Logging-Funktionen
```go
// Debug-Messages (nur sichtbar wenn Debug-Modus aktiviert)
logger.Debug("Debug-Information: %s", variable)
logger.Debugf("Debug mit Printf-Syntax: %d", number)

// Info-Messages
logger.Info("Server gestartet auf Port %s", port)
logger.Infof("Benutzer %s angemeldet", username)

// Warning-Messages
logger.Warn("Warnung: %s", warningMessage)
logger.Warnf("Performance-Warnung: %dms", duration)

// Error-Messages
logger.Error("Fehler beim Laden der Datenbank: %s", err.Error())
logger.Errorf("HTTP-Fehler %d: %s", statusCode, message)
```

### Programmgesteuerte Konfiguration
```go
import (
    "bamort/logger"
    "bamort/config"
)

// Konfiguration laden
cfg := config.LoadConfig()

// Debug-Modus setzen
logger.SetDebugMode(true)

// Minimales Log-Level setzen
logger.SetMinLogLevel(logger.DEBUG)

// Debug-Status prüfen
if logger.IsDebugEnabled() {
    // Nur ausführen wenn Debug aktiv
}
```

## Beispiel-Konfigurationen

### Development
```bash
export ENVIRONMENT=development
export DEBUG=true
export LOG_LEVEL=DEBUG
export PORT=8180
```

### Production
```bash
export ENVIRONMENT=production
export DEBUG=false
export LOG_LEVEL=INFO
export PORT=8180
```

### Testing
```bash
export ENVIRONMENT=test
export DEBUG=true
export LOG_LEVEL=WARN
export PORT=8181
```

## Log-Format

```
[2025-08-10 15:04:05] DEBUG: Debug-Nachricht hier
[2025-08-10 15:04:05] INFO: Info-Nachricht hier
[2025-08-10 15:04:05] WARN: Warning-Nachricht hier
[2025-08-10 15:04:05] ERROR: Error-Nachricht hier
```

## Docker-Konfiguration

### Dockerfile
```dockerfile
# Für Development
ENV ENVIRONMENT=development
ENV DEBUG=true
ENV LOG_LEVEL=DEBUG

# Für Production
ENV ENVIRONMENT=production
ENV DEBUG=false
ENV LOG_LEVEL=INFO
```

### docker-compose.yml
```yaml
version: '3.8'
services:
  bamort-api:
    environment:
      - ENVIRONMENT=development
      - DEBUG=true
      - LOG_LEVEL=DEBUG
      - PORT=8180
```

## Best Practices

1. **Debug-Messages**: Verwenden Sie Debug-Messages für detaillierte Informationen, die nur während der Entwicklung relevant sind.

2. **Info-Messages**: Verwenden Sie Info-Messages für wichtige Ereignisse wie Server-Start, Benutzeraktionen, etc.

3. **Warning-Messages**: Verwenden Sie Warn-Messages für potentielle Probleme, die nicht zum Absturz führen.

4. **Error-Messages**: Verwenden Sie Error-Messages für tatsächliche Fehler und Exceptions.

5. **Performance**: Debug-Messages werden automatisch gefiltert, wenn der Debug-Modus deaktiviert ist, daher keine Performance-Einbußen in der Produktion.

6. **Sensitive Daten**: Niemals Passwörter, Token oder andere sensitive Daten loggen, besonders nicht in höheren Log-Levels.
