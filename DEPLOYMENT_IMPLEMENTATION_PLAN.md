# Deployment System Implementation Plan

**Version:** 2.0 - Final Specification  
**Date:** 16. Januar 2026  
**Status:** Ready for Implementation  
**Approach:** Hybrid Migration System with CLI Deployment
**Rules:** Always work in the TestDriven Development Approach. Always Keep It Small and Simple

---

## Table of Contents
1. [Executive Summary](#executive-summary)
2. [System Architecture](#system-architecture)
3. [Version Management Strategy](#version-management-strategy)
4. [Database Migration System](#database-migration-system)
5. [Master Data Management](#master-data-management)
6. [Deployment Workflows](#deployment-workflows)
7. [Backward Compatibility](#backward-compatibility)
8. [Implementation Phases](#implementation-phases)
9. [File Structure](#file-structure)
10. [API Specifications](#api-specifications)
11. [Testing Strategy](#testing-strategy)
12. [Rollback Procedures](#rollback-procedures)
13. [Monitoring & Validation](#monitoring--validation)

---

## Executive Summary

### Approved Approach
**Hybrid Migration System** with:
- CLI-based deployment tool (no frontend UI)
- Versioned SQL migration scripts for complex changes
- GORM AutoMigrate as safety net
- Database version tracking with backend compatibility checks
- Frontend warning banner for pending migrations
- ID-independent master data sync
- Backward-compatible import system

### Deployment Frequency
- **Expected**: Once per week to once per month
- **Implication**: Medium-sized migration batches (1-10 migrations per deployment)
- **Downtime**: 20-50 minutes acceptable

### Key Decisions
✅ CLI-only deployment (no frontend deployment UI)  
✅ Frontend shows warning if DB version < backend version  
✅ Database version must match backend version range  
✅ New installations: GORM init + master data import  
✅ Master data import backward compatible with older exports  

---

## System Architecture

### Component Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Deployment Tool                     │
│                   (backend/cmd/deploy)                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ├─────────────────┐
                              │                 │
                              ▼                 ▼
┌──────────────────────────────────┐  ┌─────────────────────┐
│     Migration System             │  │   Backup System     │
│  (backend/deployment/migrations) │  │  (backend/backup)   │
├──────────────────────────────────┤  ├─────────────────────┤
│ - Version tracking               │  │ - JSON export       │
│ - SQL migration runner           │  │ - MariaDB dump      │
│ - GORM AutoMigrate fallback      │  │ - Restore logic     │
│ - Rollback support               │  │ - Retention policy  │
└──────────────────────────────────┘  └─────────────────────┘
                              │
                              ├─────────────────┐
                              │                 │
                              ▼                 ▼
┌──────────────────────────────────┐  ┌─────────────────────┐
│   Master Data Sync               │  │   Validator         │
│  (backend/deployment/masterdata) │  │  (backend/validator)│
├──────────────────────────────────┤  ├─────────────────────┤
│ - Import from dev exports        │  │ - Schema validation │
│ - Version compatibility          │  │ - Data integrity    │
│ - Dependency ordering            │  │ - Health checks     │
│ - Conflict resolution            │  │ - Version checks    │
└──────────────────────────────────┘  └─────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Version Management                         │
│               (backend/deployment/version)                  │
├─────────────────────────────────────────────────────────────┤
│ - DB version table (schema_version)                         │
│ - Backend version (from config.GetVersion())                │
│ - Compatibility matrix                                      │
│ - Migration history tracking                                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  API Health Endpoint                        │
│               GET /api/system/health                        │
├─────────────────────────────────────────────────────────────┤
│ Returns: {                                                  │
│   "backend_version": "0.4.0",                               │
│   "db_version": "0.4.0",                                    │
│   "migrations_pending": false,                              │
│   "compatible": true                                        │
│ }                                                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Frontend Warning Banner                        │
│           (frontend/src/components/SystemAlert)             │
├─────────────────────────────────────────────────────────────┤
│ Polls /api/system/health every 30 seconds                   │
│ Shows warning if migrations_pending = true                  │
│ "Database migration required. Please contact admin."        |
| Has VERY low prio                                           │
└─────────────────────────────────────────────────────────────┘
```

---

## Version Management Strategy

### Version Number Structure

**Backend Version** (existing): `MAJOR.MINOR.PATCH`
- Example: `0.4.0`, `0.5.0`, `1.0.0`
- Defined in: `backend/config/version.go`
- Already implemented: `config.GetVersion()`

**Database Version** (new): Same format as backend
- Example: `0.4.0`, `0.5.0`
- Stored in: New table `schema_version`
- Updated by: Migration system

**Required Database Version** (new): Constant defined in backend
- Example: `const RequiredDBVersion = "0.5.0"`
- Defined in: `backend/deployment/version/version.go`
- Each backend version knows exactly which DB version it requires
- No complex compatibility matrix needed

### Version Table Schema

```sql
CREATE TABLE schema_version (
    id INT PRIMARY KEY AUTO_INCREMENT,
    version VARCHAR(20) NOT NULL,           -- e.g., "0.4.0"
    migration_number INT NOT NULL,          -- e.g., 5 (5th migration applied)
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    backend_version VARCHAR(20) NOT NULL,   -- Backend version that applied this
    description TEXT,                       -- Human-readable description
    checksum VARCHAR(64),                   -- SHA256 of migration SQL (integrity check)
    
    INDEX idx_version (version),
    INDEX idx_migration_number (migration_number)
);

-- Also track individual migration applications
CREATE TABLE migration_history (
    id INT PRIMARY KEY AUTO_INCREMENT,
    migration_number INT NOT NULL UNIQUE,
    version VARCHAR(20) NOT NULL,           -- Target version
    description TEXT NOT NULL,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    applied_by VARCHAR(100),                -- CLI user or "auto"
    execution_time_ms INT,                  -- Performance tracking
    success BOOLEAN DEFAULT TRUE,
    error_message TEXT,
    rollback_available BOOLEAN DEFAULT TRUE,
    
    INDEX idx_migration_number (migration_number),
    INDEX idx_version (version)
);
```

### Compatibility Matrix

**Rule**: Database version must EXACTLY match the required version defined in backend

```go
// backend/deployment/version/version.go
const RequiredDBVersion = "0.5.0"  // Each backend version defines this
```

**Compatibility Examples:**
```
Backend 0.4.0 (requires DB 0.4.0)  ↔  DB 0.4.0  ✅ Compatible (exact match)
Backend 0.5.0 (requires DB 0.5.0)  ↔  DB 0.5.0  ✅ Compatible (exact match)
Backend 0.5.0 (requires DB 0.5.0)  ↔  DB 0.4.0  ❌ Incompatible (DB too old - migration needed)
Backend 0.4.0 (requires DB 0.4.0)  ↔  DB 0.5.0  ❌ Incompatible (DB too new - backend too old)
```

### Update Path Strategy

**Sequential Updates Required:**

If your production DB is multiple versions behind, you must update through intermediate versions:

```
Production DB: 0.3.0
Target:        0.6.0

Update Path:
0.3.0 → 0.4.0 (deploy backend 0.4.0, migrate to DB 0.4.0)
       ↓
0.4.0 → 0.5.0 (deploy backend 0.5.0, migrate to DB 0.5.0)  
       ↓
0.5.0 → 0.6.0 (deploy backend 0.6.0, migrate to DB 0.6.0)
```

**Rule:** You can only migrate to the NEXT database version. Skipping versions is not supported.

**Documentation:** Each release will document its required DB version and update path.

**Implementation:**

```go
// backend/deployment/version/version.go
package version

import (
    "bamort/config"
    "fmt"
)

// RequiredDBVersion defines the exact database version this backend requires
// This must be updated whenever database migrations are added
const RequiredDBVersion = "0.5.0"

type VersionCompatibility struct {
    BackendVersion   string
    RequiredDBVersion string
    ActualDBVersion  string
    Compatible       bool
    MigrationNeeded  bool
    Reason           string
}

func CheckCompatibility(actualDBVersion string) *VersionCompatibility {
    compatible := actualDBVersion == RequiredDBVersion
    migrationNeeded := actualDBVersion != RequiredDBVersion
    
    var reason string
    if compatible {
        reason = "Database version matches required version"
    } else if actualDBVersion < RequiredDBVersion {
        reason = fmt.Sprintf("Database migration required: %s → %s", 
            actualDBVersion, RequiredDBVersion)
    } else {
        reason = fmt.Sprintf("Backend too old for database version. Backend requires %s, database is %s", 
            RequiredDBVersion, actualDBVersion)
    }
    
    return &VersionCompatibility{
        BackendVersion:    config.GetVersion(),
        RequiredDBVersion: RequiredDBVersion,
        ActualDBVersion:   actualDBVersion,
        Compatible:        compatible,
        MigrationNeeded:   migrationNeeded,
        Reason:            reason,
    }
}

// GetRequiredDBVersion returns the database version this backend requires
func GetRequiredDBVersion() string {
    return RequiredDBVersion
}
```

### Version Update Strategy

**During Migration:**
```
1. Backend 0.4.0 deployed (requires DB 0.4.0)
2. Migrations 1-3 applied → DB version becomes 0.4.0
3. Backend checks: RequiredDBVersion (0.4.0) == ActualDBVersion (0.4.0) ✅

Later:
1. Backend 0.5.0 deployed (requires DB 0.5.0) 
2. Migration 4 applied → DB version becomes 0.5.0
3. Backend checks: RequiredDBVersion (0.5.0) == ActualDBVersion (0.5.0) ✅
```

**Version in Migration Definition:**
```go
Migration{
    Number:      4,
    Version:     "0.5.0",  // Target DB version after this migration
    Description: "Add learning_category to spells",
    UpSQL:       []string{...},
}
```

**Backend Version Constant:**
```go
// This is updated in backend/deployment/version/version.go when migrations are added
const RequiredDBVersion = "0.5.0"  // Must match migration target version
```

**Update Process:**
1. Add migrations that target version 0.5.0
2. Update `RequiredDBVersion` constant to "0.5.0"
3. Update `config.GetVersion()` to "0.5.0" 
4. Deploy - migration runs automatically, DB becomes 0.5.0
5. Backend checks compatibility: 0.5.0 == 0.5.0 ✅

---

## Database Migration System

### Migration Structure

```go
// backend/deployment/migrations/migration.go
package migrations

type Migration struct {
    Number      int                        // Sequential migration number
    Version     string                     // Target version (e.g., "0.5.0")
    Description string                     // Human-readable description
    UpSQL       []string                   // Forward migration SQL statements
    DownSQL     []string                   // Rollback SQL statements
    DataFunc    func(*gorm.DB) error       // Optional data migration function
    Critical    bool                       // If true, stops on error; if false, warns
}

// All migrations in order
var AllMigrations = []Migration{
    {
        Number:      1,
        Version:     "0.4.0",
        Description: "Initial schema version tracking",
        UpSQL: []string{
            `CREATE TABLE IF NOT EXISTS schema_version (
                id INT PRIMARY KEY AUTO_INCREMENT,
                version VARCHAR(20) NOT NULL,
                migration_number INT NOT NULL,
                applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                backend_version VARCHAR(20) NOT NULL,
                description TEXT,
                checksum VARCHAR(64)
            )`,
            `CREATE TABLE IF NOT EXISTS migration_history (
                id INT PRIMARY KEY AUTO_INCREMENT,
                migration_number INT NOT NULL UNIQUE,
                version VARCHAR(20) NOT NULL,
                description TEXT NOT NULL,
                applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                applied_by VARCHAR(100),
                execution_time_ms INT,
                success BOOLEAN DEFAULT TRUE,
                error_message TEXT,
                rollback_available BOOLEAN DEFAULT TRUE
            )`,
        },
        DownSQL: []string{
            "DROP TABLE IF EXISTS migration_history",
            "DROP TABLE IF EXISTS schema_version",
        },
        Critical: true,
    },
    {
        Number:      2,
        Version:     "0.4.0",
        Description: "Migrate spell categories to learning_category",
        UpSQL: []string{
            "ALTER TABLE gsm_spells ADD COLUMN IF NOT EXISTS learning_category VARCHAR(100)",
        },
        DownSQL: []string{
            "ALTER TABLE gsm_spells DROP COLUMN IF EXISTS learning_category",
        },
        DataFunc: func(db *gorm.DB) error {
            // Copy categorie to learning_category where empty
            return db.Exec(`
                UPDATE gsm_spells 
                SET learning_category = category 
                WHERE (learning_category IS NULL OR learning_category = '') 
                AND category IS NOT NULL
            `).Error
        },
        Critical: false,
    },
    // Future migrations will be added here
}
```

### Migration Runner

```go
// backend/deployment/migrations/runner.go
package migrations

type MigrationRunner struct {
    DB      *gorm.DB
    DryRun  bool
    Verbose bool
}

type MigrationResult struct {
    Number           int
    Description      string
    Success          bool
    ExecutionTimeMs  int64
    Error            error
    SQLExecuted      []string
}

func (r *MigrationRunner) GetCurrentVersion() (string, int, error) {
    // Query schema_version table for latest
    var version struct {
        Version         string
        MigrationNumber int
    }
    
    err := r.DB.Raw(`
        SELECT version, migration_number 
        FROM schema_version 
        ORDER BY id DESC 
        LIMIT 1
    `).Scan(&version).Error
    
    if err == gorm.ErrRecordNotFound {
        return "", 0, nil // No migrations applied yet
    }
    
    return version.Version, version.MigrationNumber, err
}

func (r *MigrationRunner) GetPendingMigrations() ([]Migration, error) {
    _, currentNumber, err := r.GetCurrentVersion()
    if err != nil {
        return nil, err
    }
    
    var pending []Migration
    for _, m := range AllMigrations {
        if m.Number > currentNumber {
            pending = append(pending, m)
        }
    }
    
    return pending, nil
}

func (r *MigrationRunner) ApplyMigration(m Migration) (*MigrationResult, error) {
    startTime := time.Now()
    result := &MigrationResult{
        Number:      m.Number,
        Description: m.Description,
    }
    
    // Transaction for safety
    err := r.DB.Transaction(func(tx *gorm.DB) error {
        // Execute SQL statements
        for _, sql := range m.UpSQL {
            if r.Verbose {
                log.Printf("Executing: %s", sql)
            }
            
            if r.DryRun {
                log.Printf("[DRY RUN] Would execute: %s", sql)
                result.SQLExecuted = append(result.SQLExecuted, sql)
                continue
            }
            
            if err := tx.Exec(sql).Error; err != nil {
                return fmt.Errorf("SQL failed: %s - Error: %w", sql, err)
            }
            result.SQLExecuted = append(result.SQLExecuted, sql)
        }
        
        // Execute data migration function if exists
        if m.DataFunc != nil && !r.DryRun {
            if r.Verbose {
                log.Printf("Executing data migration function")
            }
            if err := m.DataFunc(tx); err != nil {
                return fmt.Errorf("data migration failed: %w", err)
            }
        }
        
        // Record migration in history
        if !r.DryRun {
            history := map[string]interface{}{
                "migration_number":    m.Number,
                "version":             m.Version,
                "description":         m.Description,
                "applied_at":          time.Now(),
                "applied_by":          "deploy-cli",
                "execution_time_ms":   time.Since(startTime).Milliseconds(),
                "success":             true,
                "rollback_available":  len(m.DownSQL) > 0,
            }
            
            if err := tx.Table("migration_history").Create(history).Error; err != nil {
                return fmt.Errorf("failed to record migration: %w", err)
            }
            
            // Update schema_version
            version := map[string]interface{}{
                "version":          m.Version,
                "migration_number": m.Number,
                "applied_at":       time.Now(),
                "backend_version":  config.GetVersion(),
                "description":      m.Description,
            }
            
            if err := tx.Table("schema_version").Create(version).Error; err != nil {
                return fmt.Errorf("failed to update version: %w", err)
            }
        }
        
        return nil
    })
    
    result.ExecutionTimeMs = time.Since(startTime).Milliseconds()
    
    if err != nil {
        result.Success = false
        result.Error = err
        return result, err
    }
    
    result.Success = true
    return result, nil
}

func (r *MigrationRunner) ApplyAll() ([]*MigrationResult, error) {
    pending, err := r.GetPendingMigrations()
    if err != nil {
        return nil, err
    }
    
    if len(pending) == 0 {
        log.Println("No pending migrations")
        return nil, nil
    }
    
    log.Printf("Found %d pending migrations", len(pending))
    
    var results []*MigrationResult
    for _, migration := range pending {
        log.Printf("Applying migration %d: %s", migration.Number, migration.Description)
        
        result, err := r.ApplyMigration(migration)
        results = append(results, result)
        
        if err != nil {
            if migration.Critical {
                log.Printf("Critical migration failed, stopping: %v", err)
                return results, err
            }
            log.Printf("Warning: Non-critical migration failed: %v", err)
        }
    }
    
    return results, nil
}

func (r *MigrationRunner) Rollback(steps int) error {
    // Get last N migrations from history
    // Execute DownSQL in reverse order
    // Update version table
    // Implementation similar to ApplyMigration but reversed
}
```

### GORM AutoMigrate Integration

```go
// backend/deployment/migrations/gorm_fallback.go
package migrations

func (r *MigrationRunner) RunGORMAutoMigrate() error {
    log.Println("Running GORM AutoMigrate as safety net...")
    
    // Run models.MigrateStructure() after SQL migrations
    // This catches any columns we might have missed
    if err := models.MigrateStructure(r.DB); err != nil {
        return fmt.Errorf("GORM AutoMigrate failed: %w", err)
    }
    
    log.Println("GORM AutoMigrate completed successfully")
    return nil
}
```

---

## Master Data Management

### Export File Versioning

**Export Format with Version:**

```json
{
  "export_version": "1.0",
  "backend_version": "0.5.0",
  "timestamp": "2026-01-16T10:30:00Z",
  "game_system": "midgard",
  "data": {
    "sources": [...],
    "skills": [...],
    "spells": [...],
    "character_classes": [...],
    "learning_costs": [...]
  }
}
```

### Backward Compatibility Strategy

**Version Transformers:**

```go
// backend/deployment/masterdata/import.go
package masterdata

type ExportData struct {
    ExportVersion  string                 `json:"export_version"`
    BackendVersion string                 `json:"backend_version"`
    Timestamp      time.Time              `json:"timestamp"`
    GameSystem     string                 `json:"game_system"`
    Data           map[string]interface{} `json:"data"`
}

type ImportTransformer interface {
    CanTransform(exportVersion string) bool
    Transform(data *ExportData) (*ExportData, error)
}

// Transform old export formats to current format
type V1ToV2Transformer struct{}

func (t *V1ToV2Transformer) CanTransform(exportVersion string) bool {
    return exportVersion == "1.0"
}

func (t *V1ToV2Transformer) Transform(data *ExportData) (*ExportData, error) {
    // Example: Add missing fields with defaults
    // Example: Rename fields
    // Example: Restructure data
    
    log.Printf("Transforming export from v%s to v2.0", data.ExportVersion)
    
    // Transform spells: add learning_category if missing
    if spells, ok := data.Data["spells"].([]interface{}); ok {
        for _, spell := range spells {
            if s, ok := spell.(map[string]interface{}); ok {
                if _, exists := s["learning_category"]; !exists {
                    // Copy from category field
                    if cat, ok := s["category"].(string); ok {
                        s["learning_category"] = cat
                    }
                }
            }
        }
    }
    
    data.ExportVersion = "2.0"
    return data, nil
}

// Registry of transformers
var transformers = []ImportTransformer{
    &V1ToV2Transformer{},
    // Future transformers added here
}

func ImportMasterData(filePath string) error {
    // Read export file
    data, err := readExportFile(filePath)
    if err != nil {
        return err
    }
    
    // Get current export version (from code)
    currentVersion := "2.0" // This would be a constant
    
    // Transform if needed
    if data.ExportVersion != currentVersion {
        log.Printf("Export version %s, current version %s - transforming...", 
            data.ExportVersion, currentVersion)
        
        for _, transformer := range transformers {
            if transformer.CanTransform(data.ExportVersion) {
                data, err = transformer.Transform(data)
                if err != nil {
                    return fmt.Errorf("transformation failed: %w", err)
                }
                break
            }
        }
    }
    
    // Import data using existing gsmaster.Import system
    return importTransformedData(data)
}
```

### Master Data Import Order

**Dependency Graph:**

```
Sources (no dependencies)
    ↓
Character Classes, Skill Categories, Difficulties, Spell Schools
    ↓
Skills, Weapon Skills, Spells, Equipment
    ↓
Learning Costs (depend on classes, categories, skills)
    ↓
Skill-Category-Difficulty relationships
```

**Import Implementation:**

```go
// backend/deployment/masterdata/sync.go
package masterdata

type MasterDataSync struct {
    ImportDir string
    DB        *gorm.DB
    DryRun    bool
}

func (s *MasterDataSync) SyncAll() error {
    log.Println("Starting master data synchronization...")
    
    // Import in dependency order
    steps := []struct {
        Name     string
        ImportFn func() error
    }{
        {"Sources", s.importSources},
        {"Character Classes", s.importCharacterClasses},
        {"Skill Categories", s.importSkillCategories},
        {"Skill Difficulties", s.importSkillDifficulties},
        {"Spell Schools", s.importSpellSchools},
        {"Skills", s.importSkills},
        {"Weapon Skills", s.importWeaponSkills},
        {"Spells", s.importSpells},
        {"Equipment", s.importEquipment},
        {"Learning Costs", s.importLearningCosts},
    }
    
    for _, step := range steps {
        log.Printf("Importing %s...", step.Name)
        if err := step.ImportFn(); err != nil {
            return fmt.Errorf("failed to import %s: %w", step.Name, err)
        }
    }
    
    log.Println("Master data synchronization completed")
    return nil
}

func (s *MasterDataSync) importSources() error {
    // Use existing gsmaster.ImportSources()
    return gsmaster.ImportSources(s.ImportDir)
}

// Similar for other import functions...
```

### New Installation Setup

```go
// backend/deployment/install/initializer.go
package install

type NewInstallation struct {
    DB                *gorm.DB
    MasterDataPath    string
    CreateAdminUser   bool
    AdminUsername     string
    AdminPassword     string
}

func (n *NewInstallation) Initialize() error {
    log.Println("Initializing new installation...")
    
    // 1. Run GORM AutoMigrate to create all tables
    log.Println("Creating database schema...")
    if err := models.MigrateStructure(n.DB); err != nil {
        return fmt.Errorf("schema creation failed: %w", err)
    }
    
    // 2. Initialize version tracking
    log.Println("Initializing version tracking...")
    if err := n.initializeVersionTracking(); err != nil {
        return fmt.Errorf("version tracking failed: %w", err)
    }
    
    // 3. Import master data
    log.Println("Importing master data...")
    sync := &masterdata.MasterDataSync{
        ImportDir: n.MasterDataPath,
        DB:        n.DB,
    }
    if err := sync.SyncAll(); err != nil {
        return fmt.Errorf("master data import failed: %w", err)
    }
    
    // 4. Create admin user if requested
    if n.CreateAdminUser {
        log.Println("Creating admin user...")
        if err := n.createAdmin(); err != nil {
            return fmt.Errorf("admin creation failed: %w", err)
        }
    }
    
    log.Println("Installation completed successfully!")
    return nil
}

func (n *NewInstallation) initializeVersionTracking() error {
    // Create version tables (migration #1)
    migration := migrations.AllMigrations[0] // The schema_version migration
    
    for _, sql := range migration.UpSQL {
        if err := n.DB.Exec(sql).Error; err != nil {
            return err
        }
    }
    
    // Record initial version
    version := map[string]interface{}{
        "version":          config.GetVersion(),
        "migration_number": len(migrations.AllMigrations),
        "applied_at":       time.Now(),
        "backend_version":  config.GetVersion(),
        "description":      "Initial installation",
    }
    
    return n.DB.Table("schema_version").Create(version).Error
}

func (n *NewInstallation) createAdmin() error {
    // Use existing user creation logic
    admin := &user.User{
        Username: n.AdminUsername,
        Email:    n.AdminUsername + "@localhost",
        IsAdmin:  true,
    }
    
    if err := admin.SetPassword(n.AdminPassword); err != nil {
        return err
    }
    
    return n.DB.Create(admin).Error
}
```

---

## Deployment Workflows

### Workflow 1: Update Existing Production

```
┌─────────────────────────────────────────────────────┐
│  1. PRE-DEPLOYMENT (Development Machine)            │
├─────────────────────────────────────────────────────┤
│  $ cd backend                                       │
│  $ go run cmd/deploy/main.go prepare               │
│                                                     │
│  Actions:                                           │
│  - Export master data from dev DB                  │
│  - Create deployment package:                      │
│    deployment_package_0.5.0.tar.gz                 │
│      ├── masterdata/                               │
│      │   ├── sources.json                          │
│      │   ├── skills.json                           │
│      │   └── ...                                   │
│      ├── metadata.json (version, timestamp)        │
│      └── checksums.txt                             │
└─────────────────────────────────────────────────────┘
                        ↓
                  Transfer package
                        ↓
┌─────────────────────────────────────────────────────┐
│  2. DEPLOYMENT (Production Server)                  │
├─────────────────────────────────────────────────────┤
│  $ cd /data/prod/bamort                            │
│  $ ./deploy.sh deployment_package_0.5.0.tar.gz     │
│                                                     │
│  Or manually:                                       │
│  $ docker exec bamort-backend /app/deploy \        │
│      --package /app/deployment_package_0.5.0.tar.gz│
│      --backup                                      │
│      --migrate                                     │
│      --import-masterdata                           │
│      --validate                                    │
│                                                     │
│  Actions:                                           │
│  1. ✓ Backup production DB                        │
│  2. ✓ Check version compatibility                 │
│  3. ✓ Apply migrations                            │
│  4. ✓ Run GORM AutoMigrate                        │
│  5. ✓ Import master data                          │
│  6. ✓ Validate schema & data                      │
│  7. ✓ Update DB version                           │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  3. RESTART & VERIFY                                │
├─────────────────────────────────────────────────────┤
│  $ docker-compose restart backend                  │
│  $ curl http://localhost:8182/api/system/health    │
│                                                     │
│  Response:                                          │
│  {                                                  │
│    "backend_version": "0.5.0",                     │
│    "db_version": "0.5.0",                          │
│    "migrations_pending": false,                    │
│    "compatible": true                              │
│  }                                                  │
└─────────────────────────────────────────────────────┘
```

### Workflow 2: New Installation

```
┌─────────────────────────────────────────────────────┐
│  1. INITIAL SETUP                                   │
├─────────────────────────────────────────────────────┤
│  $ git clone repo                                   │
│  $ cd bamort/docker                                │
│  $ cp .env.example .env                            │
│  $ edit .env (set passwords, etc.)                 │
│  $ ./start-prd.sh                                  │
│                                                     │
│  Containers start, DB is empty                     │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  2. INITIALIZE DATABASE                             │
├─────────────────────────────────────────────────────┤
│  $ docker exec bamort-backend /app/deploy init \   │
│      --masterdata /app/masterdata \                │
│      --admin-user admin \                          │
│      --admin-password <secure-password>            │
│                                                     │
│  Actions:                                           │
│  1. ✓ Create all tables (GORM AutoMigrate)        │
│  2. ✓ Initialize version tracking                 │
│  3. ✓ Import master data                          │
│  4. ✓ Create admin user                           │
│  5. ✓ Validate installation                       │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│  3. READY TO USE                                    │
├─────────────────────────────────────────────────────┤
│  Access frontend at https://bamort.domain.de       │
│  Login with admin credentials                      │
└─────────────────────────────────────────────────────┘
```

### Workflow 3: Rollback

```
┌─────────────────────────────────────────────────────┐
│  ROLLBACK TO PREVIOUS VERSION                       │
├─────────────────────────────────────────────────────┤
│  $ docker exec bamort-backend /app/deploy rollback \│
│      --steps 3                                     │
│                                                     │
│  Or restore from backup:                            │
│  $ docker exec bamort-backend /app/deploy restore \ │
│      --backup /app/backups/backup_20260116.json    │
│                                                     │
│  Actions:                                           │
│  1. ✓ Stop backend                                │
│  2. ✓ Execute DownSQL for last 3 migrations       │
│  3. ✓ Update version table                        │
│  4. ✓ Validate                                     │
│  5. ✓ Restart backend                             │
└─────────────────────────────────────────────────────┘
```

---

## Backward Compatibility

### Export Version History

**Version 1.0** (Current - up to backend 0.4.x)
```json
{
  "skills": [
    {
      "name": "Schwimmen",
      "category": "Körper",
      "difficulty": "leicht"
    }
  ],
  "spells": [
    {
      "name": "Feuerkugel",
      "category": "Zerstören"  // Old field
    }
  ]
}
```

**Version 2.0** (New - backend 0.5.0+)
```json
{
  "export_version": "2.0",
  "backend_version": "0.5.0",
  "data": {
    "skills": [
      {
        "name": "Schwimmen",
        "categories_difficulties": [  // New structure
          {
            "category": "Körper",
            "difficulty": "leicht",
            "learn_cost": 1
          }
        ]
      }
    ],
    "spells": [
      {
        "name": "Feuerkugel",
        "category": "Zerstören",           // Deprecated but kept
        "learning_category": "Zerstören"   // New field
      }
    ]
  }
}
```

### Compatibility Rules

1. **Import system must accept:**
   - Version 1.0 exports (transform to 2.0)
   - Version 2.0 exports (direct import)

2. **Export system always produces:**
   - Latest version (2.0)
   - Includes deprecated fields for backward compatibility

3. **Transformation happens automatically:**
   - Detected by `export_version` field
   - Applied before import
   - Logged for audit

### Breaking Changes Policy

**When to increment export version:**
- Field renamed (provide transformer)
- Field removed (provide default value)
- Structure changed (provide converter)
- Required field added (provide transformer)

**Backward compatibility guarantee:**
- Last 2 major versions supported
- Example: Backend 1.0 supports export versions 1.x and 2.x
- Older versions require manual conversion

---

## Implementation Phases

### Phase 1: Foundation (Week 1) - Core Infrastructure

**Goal:** Build version tracking and migration system

#### Tasks:

**1.1 Version Tracking System**
- [ ] Create `backend/deployment/version/` package
- [ ] Define `RequiredDBVersion` constant
- [ ] Implement simplified `CheckCompatibility()` function (exact match only)
- [ ] Create `schema_version` and `migration_history` tables
- [ ] Implement version getters/setters
- [ ] Write unit tests for version comparison logic

**1.2 Migration Framework**
- [ ] Create `backend/deployment/migrations/` package
- [ ] Implement `Migration` struct
- [ ] Implement `MigrationRunner` with transaction support
- [ ] Add dry-run capability
- [ ] Write first migration (create version tables)
- [ ] Write tests for migration runner

**1.3 Backup Service**
- [ ] Create `backend/deployment/backup/` package
- [ ] Implement JSON export backup
- [ ] Implement MariaDB dump backup (via docker exec)
- [ ] Add backup metadata tracking
- [ ] Implement retention policy (30 days)
- [ ] Test restore procedures

**Deliverables:**
- ✅ Version tracking functional
- ✅ Migration system can apply/rollback
- ✅ Backup/restore tested
- ✅ All code tested with unit tests

**Validation:**
```bash
# Test migration system
cd backend
go test -v ./deployment/migrations/
go test -v ./deployment/version/
go test -v ./deployment/backup/
```

---

### Phase 2: Master Data & Compatibility (Week 2)

**Goal:** Enhance master data import with versioning and backward compatibility

#### Tasks:

**2.1 Export Versioning**
- [ ] Update `gsmaster/export_import.go` to include version metadata
- [ ] Define export version constant (1.0 → 2.0)
- [ ] Add version to all export files
- [ ] Update export functions to include metadata

**2.2 Import Transformers**
- [ ] Create `backend/deployment/masterdata/` package
- [ ] Implement `ImportTransformer` interface
- [ ] Create V1ToV2 transformer
- [ ] Add transformer registry
- [ ] Implement version detection logic

**2.3 Master Data Sync Enhancement**
- [ ] Implement `MasterDataSync` orchestrator
- [ ] Add dependency-ordered import
- [ ] Add conflict resolution (dev overwrites prod)
- [ ] Add import validation
- [ ] Add dry-run mode for imports

**2.4 New Installation Logic**
- [ ] Create `backend/deployment/install/` package
- [ ] Implement `NewInstallation` initializer
- [ ] Add admin user creation
- [ ] Test complete new installation flow

**Deliverables:**
- ✅ Export files versioned
- ✅ Import backward compatible (v1.0 imports work)
- ✅ Master data sync reliable
- ✅ New installation tested

**Validation:**
```bash
# Test with old export file
go test -v ./deployment/masterdata/ -run TestV1Import

# Test new installation
go test -v ./deployment/install/ -run TestNewInstallation
```

---

### Phase 3: CLI Tool (Week 3)

**Goal:** Create command-line deployment tool

#### Tasks:

**3.1 CLI Framework**
- [ ] Create `backend/cmd/deploy/` package
- [ ] Implement command structure using cobra or flag
- [ ] Add commands: `prepare`, `deploy`, `init`, `rollback`, `status`
- [ ] Add logging and progress output
- [ ] Add interactive confirmations

**3.2 Deploy Commands**
```bash
# Status check
deploy status
  - Shows current versions
  - Lists pending migrations
  - Checks compatibility

# Prepare deployment package (dev)
deploy prepare --output deployment_package.tar.gz
  - Exports master data
  - Creates package with metadata

# Full deployment (prod)
deploy deploy --package deployment_package.tar.gz
  - Backup
  - Migrate
  - Import masterdata
  - Validate

# New installation
deploy init --masterdata ./masterdata --admin-user admin

# Rollback
deploy rollback --steps 1
```

**3.3 Deployment Orchestration**
- [ ] Create `backend/deployment/orchestrator.go`
- [ ] Implement full deployment workflow
- [ ] Add progress tracking
- [ ] Add error handling and rollback triggers
- [ ] Add deployment report generation

**3.4 Shell Script Wrapper**
- [ ] Create `docker/deploy.sh` for easy use
- [ ] Handle docker exec calls
- [ ] Add parameter validation
- [ ] Add help text

**Deliverables:**
- ✅ CLI tool builds and runs
- ✅ All commands functional
- ✅ Deployment can be run end-to-end
- ✅ Shell wrapper for convenience

**Validation:**
```bash
# Build CLI
cd backend
go build -o deploy cmd/deploy/main.go

# Test commands
./deploy status
./deploy prepare --dry-run
./deploy deploy --package test.tar.gz --dry-run
```

---

### Phase 4: API Health Endpoint (Week 4)

**Goal:** Add backend API for version checking

#### Tasks:

**4.1 System Module**
- [ ] Create `backend/system/` package
- [ ] Implement health check logic
- [ ] Implement version comparison
- [ ] Add migration status check

**4.2 API Endpoints**
- [ ] Create `backend/system/handlers.go`
- [ ] Implement `GET /api/system/health`
- [ ] Implement `GET /api/system/version`
- [ ] Register routes in `cmd/main.go`
- [ ] Add tests for handlers

**4.3 Response Format**
```json
GET /api/system/health
{
  "status": "ok",
  "required_db_version": "0.5.0",
  "actual_backend_version": "0.5.0",
  "db_version": "0.5.0",
  "migrations_pending": false,
  "pending_count": 0,
  "compatible": true,
  "timestamp": "2026-01-16T10:30:00Z"
}

GET /api/system/version
{
  "backend": {
    "version": "0.5.0",
    "commit": "abc123",
    "build_date": "2026-01-16"
  },
  "database": {
    "version": "0.5.0",
    "migration_number": 5,
    "last_migration": "2026-01-15T18:00:00Z"
  }
}
```

**Deliverables:**
- ✅ Health endpoint functional
- ✅ Returns accurate version info
- ✅ Tested with unit and integration tests

**Validation:**
```bash
# Test endpoint
curl http://localhost:8180/api/system/health | jq
```

---

### Phase 5: Frontend Warning Banner (Week 5)

**Goal:** Add visual warning in frontend when migration pending

#### Tasks:

**5.1 System Alert Component**
- [ ] Create `frontend/src/components/SystemAlert.vue`
- [ ] Implement polling of `/api/system/health`
- [ ] Add warning banner UI
- [ ] Add auto-dismiss when resolved
- [ ] Style according to project CSS

**5.2 Integration**
- [ ] Add to `App.vue` or main layout
- [ ] Configure polling interval (30 seconds)
- [ ] Add translations (DE/EN)
- [ ] Test with mock data

**5.3 UI/UX**
```vue
<!-- Warning shown when migrations_pending = true -->
<div class="system-alert warning">
  <i class="icon-warning"></i>
  <span>
    Datenbank-Migration erforderlich. 
    Bitte kontaktieren Sie den Administrator.
  </span>
  <span class="version-info">
    Backend: v0.5.0 | Datenbank: v0.4.0
  </span>
</div>

<!-- Success shown when compatible = true -->
<div class="system-alert success" v-if="justUpdated">
  <i class="icon-check"></i>
  <span>System aktualisiert auf Version 0.5.0</span>
</div>
```

**Deliverables:**
- ✅ Warning banner functional
- ✅ Polls health endpoint
- ✅ Translations added
- ✅ Styled appropriately

**Validation:**
```bash
# Manually trigger warning
# Set DB version to 0.4.0, backend to 0.5.0
# Check banner appears in frontend
```

---

### Phase 6: Testing & Documentation (Week 6)

**Goal:** Comprehensive testing and documentation

#### Tasks:

**6.1 Integration Testing**
- [ ] Create test scenario: Fresh installation
- [ ] Create test scenario: Update from 0.4.0 to 0.5.0
- [ ] Create test scenario: Rollback
- [ ] Create test scenario: Import old export (v1.0)
- [ ] Test on staging environment (copy of production)
- [ ] Performance testing (large datasets)

**6.2 Documentation**
- [ ] Write deployment runbook
  - Step-by-step instructions
  - Screenshots/examples
  - Common issues & solutions
- [ ] Write rollback procedure
- [ ] Document backup/restore process
- [ ] Create troubleshooting guide
- [ ] Document version update paths (sequential upgrade requirements)
- [ ] Create version compatibility reference (RequiredDBVersion per backend version)

**6.3 Runbook Structure**
```markdown
# Deployment Runbook

## Pre-Deployment Checklist
- [ ] Backend code merged to main
- [ ] Tests passing
- [ ] Deployment package prepared
- [ ] Backup strategy confirmed
- [ ] Maintenance window scheduled

## Deployment Steps
1. Notify users of maintenance
2. Create backup
3. Transfer deployment package
4. Run deployment
5. Validate
6. Monitor for 15 minutes

## Rollback Procedure
If deployment fails:
1. Stop backend
2. Run rollback command
3. Restore backup if needed
4. Validate
5. Notify users

## Common Issues
- Migration fails: [solution]
- Version mismatch: [solution]
- Import fails: [solution]
```

**6.4 Automation Scripts**
- [ ] Create `scripts/deploy-dev.sh` (for development testing)
- [ ] Create `scripts/deploy-staging.sh`
- [ ] Create `scripts/deploy-production.sh`
- [ ] Add pre-deployment checks
- [ ] Add post-deployment validation

**Deliverables:**
- ✅ All scenarios tested successfully
- ✅ Documentation complete
- ✅ Runbook validated on staging
- ✅ Scripts ready for production use

**Validation:**
```bash
# Full deployment test on staging
cd scripts
./deploy-staging.sh deployment_package_0.5.0.tar.gz

# Verify all steps complete
# Check logs
# Test application functionality
```

---

## File Structure

```
backend/
├── cmd/
│   ├── main.go                          # Existing
│   └── deploy/                          # NEW
│       └── main.go                      # CLI deployment tool
│
├── deployment/                          # NEW PACKAGE
│   ├── orchestrator.go                  # Main deployment orchestrator
│   │
│   ├── version/                         # Version management
│   │   ├── version.go                   # Version struct and utilities
│   │   ├── compatibility.go             # Compatibility checking
│   │   └── version_test.go
│   │
│   ├── migrations/                      # Migration system
│   │   ├── migration.go                 # Migration struct
│   │   ├── runner.go                    # Migration runner
│   │   ├── gorm_fallback.go            # GORM AutoMigrate integration
│   │   ├── all_migrations.go           # All migrations registry
│   │   └── migrations_test.go
│   │
│   ├── backup/                          # Backup service
│   │   ├── backup.go                    # Backup creation
│   │   ├── restore.go                   # Restore from backup
│   │   ├── retention.go                 # Cleanup old backups
│   │   └── backup_test.go
│   │
│   ├── masterdata/                      # Master data sync
│   │   ├── import.go                    # Import with versioning
│   │   ├── export.go                    # Export with metadata
│   │   ├── transformers.go              # Version transformers
│   │   ├── sync.go                      # Sync orchestrator
│   │   └── masterdata_test.go
│   │
│   ├── install/                         # New installation
│   │   ├── initializer.go               # Fresh install logic
│   │   └── install_test.go
│   │
│   └── validator/                       # Post-deployment validation
│       ├── schema_validator.go          # Schema validation
│       ├── data_validator.go            # Data integrity checks
│       └── validator_test.go
│
├── system/                              # NEW PACKAGE - System health
│   ├── health.go                        # Health check logic
│   ├── handlers.go                      # API handlers
│   ├── routes.go                        # Route registration
│   └── system_test.go
│
├── gsmaster/                            # EXISTING - Enhanced
│   ├── export_import.go                 # Modified: add version metadata
│   └── ...
│
└── config/
    └── version.go                       # EXISTING - no changes needed

frontend/
├── src/
│   ├── components/
│   │   └── SystemAlert.vue             # NEW - Warning banner
│   │
│   └── utils/
│       └── system-health.js            # NEW - Health check polling
│
└── locales/
    ├── de/
    │   └── system.js                   # NEW - System message translations
    └── en/
        └── system.js                   # NEW

docker/
├── deploy.sh                           # NEW - Deployment wrapper script
└── ...

scripts/
├── deploy-dev.sh                       # NEW - Dev deployment
├── deploy-staging.sh                   # NEW - Staging deployment
├── deploy-production.sh                # NEW - Prod deployment
└── ...

docs/                                   # NEW
├── DEPLOYMENT_RUNBOOK.md               # Deployment procedures
├── ROLLBACK_GUIDE.md                   # Rollback procedures
├── TROUBLESHOOTING.md                  # Common issues
└── VERSION_COMPATIBILITY.md            # Version matrix

backups/                                # NEW - Created by backup service
├── backup_20260116_103000.json
├── backup_20260116_103000.sql
└── ...

masterdata/                             # Production master data exports
├── metadata.json
├── sources.json
├── skills.json
└── ...
```

---

## API Specifications

### System Health Endpoint

**Endpoint:** `GET /api/system/health`  
**Authentication:** None (public)  
**Rate Limit:** 1 req/sec per IP

**Response:**
```json
{
  "status": "ok" | "warning" | "error",
  "required_db_version": "0.5.0",
  "actual_db_version": "0.5.0",
  "migrations_pending": false,
  "pending_count": 0,
  "pending_migrations": [],
  "compatible": true,
  "migration_needed": false,
  "compatibility_reason": "Database version matches required version"
  "compatibility_reason": "Versions compatible",
  "last_migration": {
    "number": 5,
    "description": "Add learning_category",
    "applied_at": "2026-01-15T18:00:00Z"
  },
  "timestamp": "2026-01-16T10:30:00Z"
}
```

**Status Codes:**
- `200 OK` - All good, compatible
- `200 OK` + `"status": "warning"` - Migrations pending but compatible
- `500 Internal Server Error` + `"status": "error"` - Incompatible versions

**Example Responses:**
 - versions match
{
  "status": "ok",
  "backend_version": "0.5.0",
  "required_db_version": "0.5.0",
  "actual_db_version": "0.5.0",
  "migrations_pending": false,
  "compatible": true,
  "migration_needed": false,
  "compatibility_reason": "Database version matches required version"
}

// Migration needed - DB too old
{
  "status": "warning",
  "backend_version": "0.5.0",
  "required_db_version": "0.5.0",
  "actual_db_version": "0.4.0",
  "migrations_pending": true,
  "pending_count": 3,
  "pending_migrations": [
    {"number": 3, "description": "Add field X"},
    {"number": 4, "description": "Rename table Y"},
    {"number": 5, "description": "Create index Z"}
  ],
  "compatible": false,
  "migration_needed": true,
  "compatibility_reason": "Database migration required: 0.4.0 → 0.5.0"
}

// Incompatible - DB too new (backend too old)
{
  "status": "error",
  "backend_version": "0.4.0",
  "required_db_version": "0.4.0",
  "actual_db_version": "0.5.0",
  "migrations_pending": false,
  "compatible": false,
  "migration_needed": false,
  "compatibility_reason": "Backend too old for database version. Backend requires 0.4.0, database is 0.5.0
  "compatible": false,
  "compatibility_reason": "Database version newer than backend - update backend first"
}
```

### System Version Endpoint

**Endpoint:** `GET /api/system/version`  
**Authentication:** None (public)

**Response:**
```json
{
  "backend": {
    "version": "0.5.0",
    "commit_hash": "abc123def456",
    "build_date": "2026-01-16T08:00:00Z",
    "go_version": "1.25.0"
  },
  "database": {
    "version": "0.5.0",
    "migration_number": 5,
    "last_migration": {
      "number": 5,
      "description": "Add learning_category",
      "applied_at": "2026-01-15T18:00:00Z",
      "applied_by": "deploy-cli"
    },
    "total_migrations_available": 5
  },
  "compatibility": {
    "compatible": true,
    "reason": "Versions match"
  }
}
```

---

## Testing Strategy

### Unit Tests

**Coverage Target:** 80%+

**Test Files:**
```
deployment/version/version_test.go
  - TestParseVersion
  - TestCheckCompatibility
  - TestVersionComparison

deployment/migrations/runner_test.go
  - TestGetCurrentVersion
  - TestGetPendingMigrations
  - TestApplyMigration
  - TestApplyMigrationRollback
  - TestTransactionRollback

deployment/backup/backup_test.go
  - TestCreateBackup
  - TestRestoreBackup
  - TestBackupRetention

deployment/masterdata/transformers_test.go
  - TestV1ToV2Transform
  - TestTransformSkills
  - TestTransformSpells

system/health_test.go
  - TestHealthEndpoint
  - TestVersionEndpoint
  - TestCompatibilityCheck
```

### Integration Tests

**Test Scenarios:**

**Scenario 1: Fresh Installation**
```go
func TestFreshInstallation(t *testing.T) {
    // 1. Create empty test database
    // 2. Run NewInstallation.Initialize()
    // 3. Verify all tables created
    // 4. Verify version tracking initialized
    // 5. Verify master data imported
    // 6. Verify admin user created
}
```

**Scenario 2: Update Deployment**
```go
func TestUpdateDeployment(t *testing.T) {
    // 1. Set up DB with version 0.4.0
    // 2. Apply migrations to 0.5.0
    // 3. Import new master data
    // 4. Verify version updated
    // 5. Verify data integrity
}
```

**Scenario 3: Backward Compatible Import**
```go
func TestImportV1Export(t *testing.T) {
    // 1. Load v1.0 export file
    // 2. Import with transformation
    // 3. Verify data imported correctly
    // 4. Verify new fields have defaults
}
```

**Scenario 4: Rollback**
```go
func TestMigrationRollback(t *testing.T) {
    // 1. Apply 3 migrations
    // 2. Rollback 2 migrations
    // 3. Verify schema matches expected state
    // 4. Verify version updated correctly
}
```

### Manual Testing Checklist

**Pre-Release Testing:**
- [ ] Fresh install on clean Docker environment
- [ ] Update from 0.4.0 → 0.5.0 on staging
- [ ] Import old export file (v1.0)
- [ ] Rollback migration
- [ ] Health endpoint returns correct data
- [ ] Frontend warning appears when pending
- [ ] Frontend warning disappears after migration
- [ ] Performance test with 1000+ characters
- [ ] Concurrent migration attempt (should lock)

---

## Rollback Procedures

### Automatic Rollback Triggers

**Migration Runner automatically rolls back if:**
- SQL statement fails in critical migration
- Post-migration validation fails
- Version update fails
- Transaction error occurs

### Manual Rollback Options

**Option 1: Rollback N Migrations**
```bash
docker exec bamort-backend /app/deploy rollback --steps 2

# Executes DownSQL for last 2 migrations in reverse order
# Updates version table
# Updates migration history
```

**Option 2: Rollback to Specific Version**
```bash
docker exec bamort-backend /app/deploy rollback --to-version 0.4.0

# Rolls back all migrations after 0.4.0
# Updates version table to 0.4.0
```

**Option 3: Restore from Backup**
```bash
# List available backups
docker exec bamort-backend /app/deploy backup list

# Restore specific backup
docker exec bamort-backend /app/deploy restore \
    --backup /app/backups/backup_20260116_103000.json

# Full restore (stops backend, restores DB, starts backend)
```

### Emergency Rollback (Full)

```bash
# 1. Stop backend
docker-compose -f docker/docker-compose.yml stop backend

# 2. Restore MariaDB volume from backup
docker-compose -f docker/docker-compose.yml stop mariadb
rm -rf docker/bamort-db/*
tar -xzf backups/mariadb_backup_20260116.tar.gz -C docker/bamort-db/

# 3. Start MariaDB
docker-compose -f docker/docker-compose.yml start mariadb

# 4. Wait for health check
sleep 10

# 5. Start backend
docker-compose -f docker/docker-compose.yml start backend

# 6. Verify
curl http://localhost:8182/api/system/health
```

### Rollback Time Estimates

- **Migration Rollback (1-3 steps):** 30 seconds - 2 minutes
- **JSON Restore:** 2-5 minutes (depends on data size)
- **Full MariaDB Restore:** 5-15 minutes (depends on volume size)
- **Emergency Full Rollback:** 10-20 minutes

---

## Monitoring & Validation

### Post-Deployment Validation Checklist

**Automated Checks (by validator):**
- [ ] All expected tables exist
- [ ] All expected columns exist in each table
- [ ] No orphaned tables (old renamed tables)
- [ ] Foreign key constraints valid
- [ ] Version table updated correctly
- [ ] Migration history recorded
- [ ] Master data counts reasonable (not empty)

**Manual Verification:**
- [ ] Health endpoint returns `compatible: true`
- [ ] Frontend loads without errors
- [ ] Can create new character
- [ ] Can view existing characters
- [ ] Skills/spells searchable
- [ ] Learning cost calculations work
- [ ] No errors in backend logs

### Monitoring Hooks

**Log Files:**
```
backend/logs/deployment_20260116_103000.log
  - All deployment steps
  - Migration SQL executed
  - Import results
  - Validation results
  - Timing information
```

**Metrics to Track:**
- Migration execution time
- Number of records imported
- Backup size
- Database size before/after
- Error count during deployment

### Health Check Integration

**Automated Monitoring (optional):**
```bash
# Add to cron for continuous monitoring
*/5 * * * * curl -s http://localhost:8182/api/system/health | \
    jq '.compatible' | grep -q true || \
    echo "ALERT: Version incompatibility detected" | mail -s "Bamort Alert" admin@example.com
```

**Prometheus Metrics (future enhancement):**
```go
// Expose metrics for monitoring tools
bamort_backend_version{version="0.5.0"}
bamort_db_version{version="0.5.0"}
bamort_migrations_pending 0
bamort_compatibility_status 1
```

---

## Summary & Next Steps

### What We're Building

A **production-ready deployment system** with:

✅ **Versioned migrations** - SQL scripts + GORM safety net  
✅ **Version tracking** - DB version must be compatible with backend  
✅ **Backward compatible imports** - Old export files still work  
✅ **CLI deployment tool** - Easy to use, automatable  
✅ **Frontend warnings** - Users notified if migration needed  
✅ **Rollback support** - Can undo migrations or restore backups  
✅ **New install support** - Fresh setups initialize cleanly  

### Implementation Timeline

- **Week 1:** Version tracking + Migration system + Backup
- **Week 2:** Master data versioning + Import compatibility
- **Week 3:** CLI tool + Deployment orchestration
- **Week 4:** Health API endpoints
- **Week 5:** Frontend warning banner
- **Week 6:** Testing + Documentation

**Total:** 6 weeks to production-ready deployment system

### Approval Checklist

Before starting implementation, confirm:

- [ ] Hybrid approach approved
- [ ] CLI-only deployment approved (no frontend UI)
- [ ] Frontend warning banner approved
- [ ] Version compatibility strategy approved
- [ ] 30-day backup retention approved
- [ ] Implementation timeline acceptable
- [ ] Resource allocation confirmed

### Ready to Start?

Once approved, implementation begins with **Phase 1: Foundation**.

First commit will include:
- `backend/deployment/version/` package
- `backend/deployment/migrations/` package  
- `backend/deployment/backup/` package
- Unit tests for all packages
- First migration (create version tables)

**Next Command:** 
```bash
# Start implementation
git checkout -b feature/deployment-system
# Begin Phase 1 implementation
```

---

**End of Implementation Plan**

*This plan should be reviewed and approved before implementation begins. Updates will be tracked in git commits.*
