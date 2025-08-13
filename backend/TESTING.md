# Test Environment Configuration

## Überblick

Alle Tests in diesem Projekt sind so konfiguriert, dass sie `ENVIRONMENT=test` verwenden, um sicherzustellen, dass:

1. Die richtige Datenbank-Konfiguration verwendet wird (Test-DB statt Production-DB)
2. Logger-Ausgaben reduziert sind 
3. Tests isoliert von der Produktionsumgebung laufen

## Test-Umgebung einrichten

### Option 1: Lokale setupTestEnvironment Funktion (empfohlen)

Jeder Test sollte eine lokale `setupTestEnvironment` Funktion verwenden:

```go
package mypackage

import (
    "os"
    "testing"
)

// setupTestEnvironment setzt ENVIRONMENT=test für Tests
func setupTestEnvironment(t *testing.T) {
    original := os.Getenv("ENVIRONMENT")
    os.Setenv("ENVIRONMENT", "test")
    t.Cleanup(func() {
        if original != "" {
            os.Setenv("ENVIRONMENT", original)
        } else {
            os.Unsetenv("ENVIRONMENT")
        }
    })
}

func TestMyFunction(t *testing.T) {
    setupTestEnvironment(t)
    
    // Ihr Test-Code hier...
}
```

### Option 2: testutils Package (für komplexere Setups)

Für komplexere Test-Setups verwenden Sie das `testutils` Package:

```go
package mypackage

import (
    "bamort/testutils"
    "testing"
)

func TestMyFunction(t *testing.T) {
    testutils.SetupTestEnvironment(t)
    
    // Oder für spezifische Konfiguration:
    testutils.SetupTestEnvironmentWithConfig(t, map[string]string{
        "DEBUG": "false",
        "LOG_LEVEL": "ERROR",
    })
    
    // Ihr Test-Code hier...
}
```

## Aktualisierte Test-Dateien

Folgende Test-Dateien wurden bereits aktualisiert:

- ✅ `config/config_test.go` - Alle Tests setzen ENVIRONMENT=test
- ✅ `logger/logger_test.go` - Alle Tests setzen ENVIRONMENT=test  
- ✅ `character/character_test.go` - setupTestEnvironment hinzugefügt
- ✅ `character/handlers_test.go` - Test-Umgebung konfiguriert

## Database-Tests

Die `database.ConnectDatabase()` Funktion erkennt automatisch `ENVIRONMENT=test` und verwendet die Test-Datenbank:

```go
func ConnectDatabase() *gorm.DB {
    cfg := config.LoadConfig()
    
    if cfg.Environment == "test" {
        logger.Debug("Test-Umgebung erkannt, verwende Test-Datenbank")
        SetupTestDB()
    } else {
        logger.Debug("Verwende konfigurierte Datenbank (%s)", cfg.DatabaseType)
        return ConnectDatabaseOrig()
    }
    
    return DB
}
```

## Best Practices

1. **Immer setupTestEnvironment() aufrufen**: Jeder Test sollte die Umgebung konfigurieren
2. **Cleanup automatisch**: Verwenden Sie `t.Cleanup()` um ursprüngliche Werte wiederherzustellen
3. **Test-Isolation**: Tests sollten sich nicht gegenseitig beeinflussen
4. **Keine Produktionsdaten**: Tests verwenden immer Test-Datenbank

## Verifikation

Um zu prüfen, ob alle Tests korrekt konfiguriert sind:

```bash
# Script ausführen
./verify_test_environment.sh

# Einzelne Packages testen
go test ./config -v
go test ./logger -v
go test ./character -v
```

## Migration bestehender Tests

Für bestehende Tests, die noch nicht aktualisiert wurden:

1. Import für `os` hinzufügen
2. `setupTestEnvironment` Funktion hinzufügen  
3. `setupTestEnvironment(t)` am Anfang jeder Test-Funktion aufrufen

Beispiel:

```go
// Vor der Migration
func TestMyFunction(t *testing.T) {
    // Test-Code...
}

// Nach der Migration  
func TestMyFunction(t *testing.T) {
    setupTestEnvironment(t)
    // Test-Code...
}
```

## Troubleshooting

**Problem**: Test verwendet Produktions-Datenbank
**Lösung**: `setupTestEnvironment(t)` am Anfang des Tests aufrufen

**Problem**: Umgebungsvariablen beeinflussen sich zwischen Tests
**Lösung**: `t.Cleanup()` verwenden für automatisches Aufräumen

**Problem**: Import-Zyklen bei testutils
**Lösung**: Lokale `setupTestEnvironment` Funktion in jedem Package verwenden
