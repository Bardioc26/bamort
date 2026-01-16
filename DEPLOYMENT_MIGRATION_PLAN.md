# Deployment and Data Migration Planning Document

**Version:** 1.0  
**Date:** 15. Januar 2026  
**Status:** Planning Phase - No Implementation Yet

---

## Table of Contents
1. [Executive Summary](#executive-summary)
2. [Problem Analysis](#problem-analysis)
3. [Current System Assessment](#current-system-assessment)
4. [Solution Proposals](#solution-proposals)
5. [Recommended Approach](#recommended-approach)
6. [Implementation Plan](#implementation-plan)
7. [Risk Assessment](#risk-assessment)
8. [Rollback Strategy](#rollback-strategy)
9. [Open Questions](#open-questions)

---

## Executive Summary

### Goal
Create a robust, automated deployment procedure that handles:
- Database structure migrations (schema changes, renamed tables, new fields)
- Master data synchronization from development to production
- Preservation of production user data
- Minimal/zero downtime deployments
- Safe rollback capabilities

### Key Challenges
1. **Dual Database Evolution**: Dev and prod databases evolve independently between deployments
2. **ID Mismatch**: Same entities have different primary keys in dev vs prod
3. **GORM Limitations**: GORM's AutoMigrate has issues with renamed tables and index management
4. **Data Types**:
   - System data (users) - production-only, must never be overwritten
   - Master data (gsm_*, learning_*) - updated in dev, must sync to prod
   - User data (char_*, eq_*) - production-only, must be preserved

---

## Problem Analysis

### Current State
Based on codebase analysis:

#### Existing Infrastructure
✅ **Database Migration System**:
- `models.MigrateStructure()` - Uses GORM AutoMigrate
- Separated into domains: gsMaster, character, equipment, skills, learning
- Already handles optional database parameter

✅ **Export/Import Functionality**:
- `gsmaster/export_import.go` - Master data export/import (ID-independent)
- `transfer/database.go` - Full database export/import
- Uses natural keys (name + game_system) for matching
- JSON-based, version-tracked exports

✅ **Docker Infrastructure**:
- Production: `docker-compose.yml` with MariaDB
- Development: `docker-compose.dev.yml` with live-reload
- Health checks for database readiness

#### Identified Gaps
❌ **No Automated Deployment Pipeline**
❌ **No Production Backup Strategy**
❌ **No Schema Comparison/Validation**
❌ **No Data Migration Orchestration**
❌ **No Index Management for Renamed Tables**
❌ **No Deployment UI/Frontend Interface**

### Specific Issues

#### 1. GORM AutoMigrate Limitations
**Problem**: GORM cannot handle:
- Table renames (creates new table instead of renaming)
- Index renames when underlying table/column is renamed
- Constraint modifications
- Data migrations (only schema)

**Evidence from code**:
```go
// backend/models/database.go - Line 9
func MigrateStructure(db ...*gorm.DB) error {
    // ...
    err := targetDB.AutoMigrate(&Skill{}, &WeaponSkill{}, &Spell{}, ...)
    // AutoMigrate only adds columns/tables, doesn't rename or remove
}
```

#### 2. ID-Independent Import System Exists
**Strength**: Already implemented in gsmaster module
```go
// backend/gsmaster/export_import.go
type ExportableSkill struct {
    Name           string `json:"name"`
    GameSystem     string `json:"game_system"` // Natural key
    SourceCode     string `json:"source_code"` // Not SourceID!
    // No ID field - matches by name+system
}
```

#### 3. Multiple Export/Import Systems
**Current systems**:
1. **gsmaster/export_import.go**: Master data only, ID-independent ✅
2. **transfer/database.go**: Full database dump, includes IDs ⚠️
3. **importer/**: VTT format imports

**Issue**: No unified migration strategy

---

## Current System Assessment

### Strengths
1. ✅ **Modular Migration Structure**: Domain-separated migrations are maintainable
2. ✅ **ID-Independent Export**: GSMaster export uses natural keys
3. ✅ **Transaction Support**: Import operations use GORM transactions
4. ✅ **Version Tracking**: Exports include version and timestamp
5. ✅ **Docker Health Checks**: Database readiness is checked
6. ✅ **Test Coverage**: Extensive tests for migrations and imports

### Weaknesses
1. ❌ **No Deployment Orchestration**: Manual process, error-prone
2. ❌ **No Backup Automation**: Must be done manually
3. ❌ **No Schema Validation**: No check for old/orphaned fields
4. ❌ **GORM Limitations**: Cannot handle complex schema changes
5. ❌ **No Migration History**: No tracking of applied migrations
6. ❌ **No Dry-Run Capability**: Cannot preview changes
7. ❌ **No Data Conflict Resolution**: No strategy for dev vs prod data conflicts

### Data Flow Analysis

```
Development Database:
├── System Data (users) - test users only
├── Master Data (gsm_*, learning_*) - actively developed, enriched
└── User Data (char_*, eq_*) - test characters

Production Database:
├── System Data (users) - MUST PRESERVE - real users
├── Master Data (gsm_*, learning_*) - MUST UPDATE from dev
└── User Data (char_*, eq_*) - MUST PRESERVE - real characters
```

---

## Solution Proposals

### Proposal 1: "Pure GORM AutoMigrate" (Current State)
**Approach**: Continue using GORM AutoMigrate, handle edge cases manually

**Implementation**:
```
1. Backup production DB → JSON export
2. Run GORM AutoMigrate (adds new columns/tables)
3. Manually run SQL for renames/drops
4. Import master data via gsmaster import
```

**Pros**:
- ✅ Minimal code changes
- ✅ Leverages existing system
- ✅ Fast implementation

**Cons**:
- ❌ Error-prone (manual SQL steps)
- ❌ No validation of schema correctness
- ❌ GORM still creates duplicate tables on renames
- ❌ No automated rollback
- ❌ Requires database expertise for each deployment

**Risk Level**: 🔴 **HIGH** - Manual steps, no validation

---

### Proposal 2: "Migration Scripts + GORM" (Hybrid)
**Approach**: Add versioned migration scripts for complex changes, keep GORM for simple additions

**Implementation**:
```
backend/migrations/
├── 001_initial_schema.sql
├── 002_rename_gsm_tables.sql
├── 003_add_learning_categories.sql
└── migration_runner.go

Deployment Process:
1. Backup production DB → automated
2. Run pending migration scripts (SQL)
3. Run GORM AutoMigrate (safety net)
4. Import master data updates
5. Validate schema matches expected state
```

**Pros**:
- ✅ Handles complex schema changes (renames, index updates)
- ✅ Version-controlled migrations
- ✅ Testable against SQLite/MariaDB
- ✅ Can be automated
- ✅ Clear audit trail

**Cons**:
- ⚠️ Requires maintaining SQL scripts
- ⚠️ Must write DB-specific SQL (MariaDB vs SQLite)
- ⚠️ Two migration systems (scripts + GORM)

**Risk Level**: 🟡 **MEDIUM** - More complex, but controllable

**Example Implementation**:
```go
// backend/migrations/runner.go
type Migration struct {
    Version     int
    Description string
    UpSQL       string
    DownSQL     string // for rollback
}

func ApplyMigrations(db *gorm.DB) error {
    // Check version table
    // Run pending migrations
    // Update version
}
```

---

### Proposal 3: "State-Based Schema Management" (Advanced)
**Approach**: Define desired schema in code, generate migration from current→desired state

**Tools**: 
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [pressly/goose](https://github.com/pressly/goose)
- [Atlas](https://atlasgo.io/) - schema-as-code

**Implementation**:
```
1. Define models in GORM (current state)
2. Tool compares production DB schema with models
3. Tool generates migration SQL
4. Review + apply migration
5. Import master data
```

**Pros**:
- ✅ Automatic migration generation
- ✅ Handles all schema changes correctly
- ✅ Professional-grade solution
- ✅ Rollback support built-in
- ✅ Schema validation included

**Cons**:
- ❌ New dependency/learning curve
- ❌ Most complex to implement
- ❌ May not handle all GORM edge cases
- ❌ Overkill for current project size

**Risk Level**: 🟢 **LOW** (if implemented correctly) - But 🟡 **MEDIUM** effort

---

### Proposal 4: "Export-Drop-Import" (Nuclear Option)
**Approach**: Full database rebuild from scratch on each deployment

**Implementation**:
```
1. Export production DB → full JSON dump
2. Stop application
3. Drop all tables
4. Run fresh migrations
5. Import master data from dev
6. Import user data from prod backup
7. Validate data integrity
8. Start application
```

**Pros**:
- ✅ Guarantees clean schema
- ✅ No orphaned fields/tables
- ✅ Simple concept
- ✅ Forces good backup hygiene

**Cons**:
- ❌ Requires downtime
- ❌ Slow for large databases
- ❌ High risk if import fails
- ❌ Complicated rollback

**Risk Level**: 🟡 **MEDIUM** - Simple but risky

---

## Recommended Approach

### **Hybrid Strategy: Migration Scripts + Enhanced Export/Import + Validation**

Combines best aspects of Proposals 2 & 4 with pragmatic risk mitigation.

### Architecture

```
Deployment Pipeline:
┌─────────────────────────────────────────────────────────┐
│ 1. PRE-DEPLOYMENT PHASE                                 │
├─────────────────────────────────────────────────────────┤
│ ✓ Backup production DB → timestamped JSON export       │
│ ✓ Backup MariaDB data directory (docker volume)        │
│ ✓ Export dev master data → migration package           │
│ ✓ Run schema validation tests                          │
└─────────────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────────────┐
│ 2. DEPLOYMENT PHASE (Can be frontend-triggered)        │
├─────────────────────────────────────────────────────────┤
│ ✓ Stop backend container                               │
│ ✓ Run migration scripts (if any)                       │
│ ✓ Run GORM AutoMigrate (safety net)                    │
│ ✓ Import master data updates (ID-independent)          │
│ ✓ Run data migrations (if needed)                      │
│ ✓ Start backend container                              │
└─────────────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────────────┐
│ 3. VALIDATION PHASE                                     │
├─────────────────────────────────────────────────────────┤
│ ✓ Health check API endpoint                            │
│ ✓ Schema validation (expected tables/columns exist)    │
│ ✓ Data integrity checks (foreign keys, counts)         │
│ ✓ Smoke tests (create test character, run query)       │
└─────────────────────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────────────────────┐
│ 4. POST-DEPLOYMENT                                      │
├─────────────────────────────────────────────────────────┤
│ ✓ Keep backup for 7 days                               │
│ ✓ Log deployment in audit_log_entries                  │
│ ✓ Monitor error rates                                  │
└─────────────────────────────────────────────────────────┘
```

### Key Components

#### 1. **Migration System** (`backend/migrations/`)
```go
package migrations

type Migration struct {
    Version     int       // Sequential version number
    Description string    // Human-readable description
    UpSQL       []string  // Forward migration SQL statements
    DownSQL     []string  // Rollback SQL statements
    DataMigration func(*gorm.DB) error // Optional data transformation
}

var Migrations = []Migration{
    {
        Version:     1,
        Description: "Add learning_category to spells",
        UpSQL: []string{
            "ALTER TABLE gsm_spells ADD COLUMN learning_category VARCHAR(100)",
            "UPDATE gsm_spells SET learning_category = category WHERE learning_category IS NULL",
        },
        DownSQL: []string{
            "ALTER TABLE gsm_spells DROP COLUMN learning_category",
        },
    },
    // More migrations...
}
```

#### 2. **Backup Service** (`backend/deployment/backup.go`)
```go
package deployment

type BackupService struct {
    BackupDir string
}

func (s *BackupService) CreateBackup() (*BackupResult, error) {
    timestamp := time.Now().Format("20060102_150405")
    
    // 1. Export database to JSON
    jsonBackup := filepath.Join(s.BackupDir, fmt.Sprintf("backup_%s.json", timestamp))
    
    // 2. Create MariaDB dump (if possible)
    sqlBackup := filepath.Join(s.BackupDir, fmt.Sprintf("backup_%s.sql", timestamp))
    
    // 3. Return backup metadata
    return &BackupResult{
        Timestamp:    timestamp,
        JSONPath:     jsonBackup,
        SQLPath:      sqlBackup,
        RecordCounts: recordCounts,
    }, nil
}
```

#### 3. **Master Data Sync** (`backend/deployment/masterdata.go`)
```go
package deployment

// SyncMasterData imports master data from dev export to production
func SyncMasterData(exportDir string) error {
    // Import in correct dependency order:
    // 1. Sources (no dependencies)
    // 2. Character classes, skill categories, difficulties, spell schools
    // 3. Skills, spells (depend on sources, categories)
    // 4. Learning costs (depend on classes, categories, skills)
    
    return gsmaster.ImportAll(exportDir)
}
```

#### 4. **Schema Validator** (`backend/deployment/validator.go`)
```go
package deployment

type SchemaValidator struct {
    DB *gorm.DB
}

func (v *SchemaValidator) Validate() (*ValidationReport, error) {
    // Check expected tables exist
    // Check expected columns exist
    // Check for orphaned tables (old names)
    // Check foreign key constraints
    // Check index definitions
}
```

#### 5. **Deployment Orchestrator** (`backend/deployment/orchestrator.go`)
```go
package deployment

type Orchestrator struct {
    MigrationRunner *MigrationRunner
    BackupService   *BackupService
    Validator       *SchemaValidator
    MasterDataSync  *MasterDataSync
}

func (o *Orchestrator) Deploy(ctx context.Context) (*DeploymentReport, error) {
    // Phase 1: Backup
    backup, err := o.BackupService.CreateBackup()
    
    // Phase 2: Migrations
    migrationResults, err := o.MigrationRunner.ApplyPending()
    
    // Phase 3: Master data
    syncResults, err := o.MasterDataSync.Sync()
    
    // Phase 4: Validation
    validation, err := o.Validator.Validate()
    
    return &DeploymentReport{...}, nil
}
```

#### 6. **API Endpoints** (`backend/deployment/handlers.go`)
```go
// Protected admin-only routes
func RegisterRoutes(r *gin.RouterGroup) {
    deploy := r.Group("/deployment")
    deploy.Use(middleware.RequireAdmin()) // Authentication required!
    
    {
        // Read-only status
        deploy.GET("/status", GetDeploymentStatusHandler)
        deploy.GET("/backups", ListBackupsHandler)
        deploy.GET("/migrations/pending", GetPendingMigrationsHandler)
        
        // Actions
        deploy.POST("/backup/create", CreateBackupHandler)
        deploy.POST("/migrations/apply", ApplyMigrationsHandler)
        deploy.POST("/masterdata/sync", SyncMasterDataHandler)
        deploy.POST("/validate", ValidateSchemaHandler)
        deploy.POST("/deploy", FullDeploymentHandler) // All-in-one
        
        // Rollback
        deploy.POST("/rollback", RollbackHandler)
    }
}
```

#### 7. **Frontend Deployment Page** (Optional)
```vue
<!-- frontend/src/views/DeploymentView.vue -->
<template>
  <div class="deployment-panel">
    <h1>System Deployment</h1>
    
    <section class="backup-section">
      <h2>1. Backup</h2>
      <button @click="createBackup">Create Backup</button>
      <div v-if="backupStatus">
        Last backup: {{ backupStatus.timestamp }}
        Records: {{ backupStatus.recordCount }}
      </div>
    </section>
    
    <section class="migration-section">
      <h2>2. Migrations</h2>
      <div v-for="migration in pendingMigrations" :key="migration.version">
        v{{ migration.version }}: {{ migration.description }}
      </div>
      <button @click="applyMigrations">Apply Migrations</button>
    </section>
    
    <section class="masterdata-section">
      <h2>3. Master Data Sync</h2>
      <button @click="syncMasterData">Sync from Development</button>
    </section>
    
    <section class="deploy-section">
      <h2>Full Deployment</h2>
      <button @click="fullDeploy" class="primary">
        Deploy All
      </button>
    </section>
  </div>
</template>
```

---

## Implementation Plan

### Phase 1: Foundation (Week 1-2)
**Goal**: Set up infrastructure without breaking existing system

#### Tasks:
1. **Create Migration System**
   - [ ] Create `backend/migrations/` package
   - [ ] Implement `Migration` struct and runner
   - [ ] Create migration version table in DB
   - [ ] Write tests for migration runner

2. **Create Backup Service**
   - [ ] Implement automated JSON export
   - [ ] Add MariaDB dump capability (docker exec)
   - [ ] Create backup retention policy (7 days)
   - [ ] Test restore from backup

3. **Add Schema Validator**
   - [ ] Implement table existence checks
   - [ ] Implement column existence checks
   - [ ] Add foreign key validation
   - [ ] Create validation report format

**Deliverable**: Core deployment infrastructure, fully tested

---

### Phase 2: Master Data Sync Enhancement (Week 3)
**Goal**: Ensure master data can be reliably synced

#### Tasks:
1. **Enhance Export System**
   - [ ] Consolidate `gsmaster/export_import.go` as primary system
   - [ ] Add dependency-ordered export/import
   - [ ] Add conflict resolution strategies
   - [ ] Test with real dev→prod scenario

2. **Data Validation**
   - [ ] Add integrity checks post-import
   - [ ] Validate foreign key relationships
   - [ ] Check for orphaned records
   - [ ] Log import statistics

**Deliverable**: Robust master data synchronization

---

### Phase 3: API & Orchestration (Week 4)
**Goal**: Tie everything together with API layer

#### Tasks:
1. **Create Deployment Package**
   - [ ] Implement `deployment/orchestrator.go`
   - [ ] Create deployment report structure
   - [ ] Add comprehensive logging
   - [ ] Implement rollback logic

2. **Add API Endpoints**
   - [ ] Create handlers for each operation
   - [ ] Add admin authentication middleware
   - [ ] Implement operation status tracking
   - [ ] Add error handling and reporting

3. **Testing**
   - [ ] Test complete deployment flow
   - [ ] Test rollback scenarios
   - [ ] Test with realistic data volumes
   - [ ] Load testing

**Deliverable**: Complete backend deployment API

---

### Phase 4: Frontend Interface (Week 5 - Optional)
**Goal**: Add user-friendly deployment interface

#### Tasks:
1. **Create Deployment View**
   - [ ] Build Vue component
   - [ ] Add status indicators
   - [ ] Implement progress tracking
   - [ ] Add confirmation dialogs

2. **Integration**
   - [ ] Connect to backend API
   - [ ] Add real-time updates (websocket or polling)
   - [ ] Implement error display
   - [ ] Add deployment history view

**Deliverable**: Admin deployment interface

---

### Phase 5: Documentation & Training (Week 6)
**Goal**: Ensure deployment process is documented and understood

#### Tasks:
1. **Documentation**
   - [ ] Write deployment runbook
   - [ ] Document rollback procedure
   - [ ] Create troubleshooting guide
   - [ ] Document backup/restore procedures

2. **Automation Scripts**
   - [ ] Create shell script for manual deployment
   - [ ] Add pre-deployment checklist
   - [ ] Create monitoring alerts
   - [ ] Set up automated backups (cron)

**Deliverable**: Complete deployment documentation

---

## Risk Assessment

### High-Risk Areas

#### 1. **Data Loss During Migration**
**Risk**: Migration fails mid-process, data corrupted  
**Mitigation**:
- ✅ Always backup before migration
- ✅ Use database transactions
- ✅ Validate before committing
- ✅ Test migrations on copy of production first
- ✅ Keep backup for 7+ days

**Residual Risk**: 🟢 LOW

#### 2. **Master Data Conflicts**
**Risk**: Dev and prod both have updated same records differently  
**Mitigation**:
- ✅ Use "last write wins" for master data (dev overwrites prod)
- ✅ Master data should only be edited in dev
- ✅ Add audit logging for changes
- ✅ Can implement conflict detection (future)

**Residual Risk**: 🟡 MEDIUM

#### 3. **Schema Migration Failures**
**Risk**: Migration SQL fails, leaves DB in broken state  
**Mitigation**:
- ✅ Test migrations on SQLite first
- ✅ Test on copy of production
- ✅ Implement rollback SQL for each migration
- ✅ Use transactions where possible
- ✅ Validate schema after migration

**Residual Risk**: 🟡 MEDIUM

#### 4. **GORM Index Issues on Table Renames**
**Risk**: GORM creates duplicate indexes after table rename  
**Mitigation**:
- ✅ Handle table renames in migration scripts (before GORM)
- ✅ Explicitly drop old indexes in migration
- ✅ Let GORM recreate indexes on new table
- ✅ Validate index definitions post-migration

**Residual Risk**: 🟢 LOW

#### 5. **Production Downtime**
**Risk**: Deployment takes longer than expected  
**Mitigation**:
- ⚠️ Accept brief downtime window
- ✅ Test deployment duration on staging
- ✅ Schedule during low-traffic periods
- ✅ Can implement blue-green deployment later

**Residual Risk**: 🟡 MEDIUM (acceptable for single-user/small team app)

### Medium-Risk Areas

#### 6. **Foreign Key Constraint Violations**
**Risk**: Master data import violates constraints  
**Mitigation**:
- ✅ Import in dependency order (sources before skills)
- ✅ Validate foreign key integrity post-import
- ✅ Use CASCADE deletes appropriately

**Residual Risk**: 🟢 LOW

#### 7. **Version Skew (Dev too far ahead of Prod)**
**Risk**: Dev database has 10 migrations, prod has 0  
**Mitigation**:
- ✅ Test applying all pending migrations in sequence
- ✅ Can create "catch-up" migration if needed
- ✅ Document deployment frequency requirements

**Residual Risk**: 🟢 LOW

---

## Rollback Strategy

### Automatic Rollback Triggers
- Migration script fails (SQL error)
- Post-migration validation fails
- Master data import fails critically
- Health check fails after deployment

### Rollback Procedure

#### Level 1: Migration Rollback (Fast - 1-2 minutes)
```
1. Run migration DownSQL scripts in reverse order
2. Restore from JSON backup (master data only)
3. Restart backend
4. Validate
```

#### Level 2: Full Database Restore (Slow - 5-15 minutes)
```
1. Stop backend container
2. Restore MariaDB volume from backup
   OR
   Import from JSON backup (full)
3. Start backend container
4. Validate
```

#### Level 3: Container Rollback (Fallback - 2-5 minutes)
```
1. Restore previous backend Docker image
2. Restore database from backup (Level 2)
3. Start containers
```

### Testing Rollback
- [ ] Test each rollback level in development
- [ ] Document expected rollback times
- [ ] Practice rollback procedure
- [ ] Automate as much as possible

---

## Open Questions

### Technical Decisions Needed

1. **Migration SQL: MariaDB-specific or cross-compatible?**
   - Option A: Write separate SQL for MariaDB and SQLite (more work)
   - Option B: MariaDB-only (tests may differ from production)
   - **Recommendation**: B - Tests use SQLite for speed, final integration test on MariaDB

2. **Deployment trigger: API-only or also command-line?**
   - Option A: API only (requires frontend or curl)
   - Option B: Also provide CLI tool
   - **Recommendation**: Both - API for normal use, CLI for emergency/automation

3. **Master data conflict resolution: Automatic or manual?**
   - Option A: Dev always overwrites prod (simple)
   - Option B: Detect conflicts, require manual resolution (safe)
   - **Recommendation**: A initially, add conflict detection in Phase 2

4. **Downtime acceptable?**
   - Option A: Accept 2-5 minutes downtime per deployment
   - Option B: Implement zero-downtime blue-green deployment
   - **Recommendation**: A - Acceptable for current scale, can upgrade later

5. **Backup retention: How long?**
   - Option A: 7 days
   - Option B: 30 days
   - Option C: Keep all backups (with rotation)
   - **Recommendation**: B - 30 days, compress old backups

6. **Who can trigger deployments?**
   - Option A: Admin users only (in app)
   - Option B: Only server SSH access
   - Option C: Both
   - **Recommendation**: C - Frontend for normal, SSH for emergency

### Process Questions

7. **How to handle dev database changes during development?**
   - Should every schema change require a migration script immediately?
   - Or batch changes and create migration before deployment?
   - **Recommendation**: Batch during development, create migration before merge to main

8. **Testing strategy for migrations?**
   - Test on copy of production before actual deployment?
   - Or rely on automated tests?
   - **Recommendation**: Both - automated tests + manual staging deployment

9. **Deployment frequency?**
   - How often are deployments expected?
   - This affects migration strategy (many small vs few large)
   - **Recommendation**: Define target (e.g., weekly/biweekly)

---

## Next Steps

### Immediate Actions (Before Implementation)

1. **Review & Discuss**
   - [ ] Review this document
   - [ ] Discuss solution proposals
   - [ ] Make technical decisions (see Open Questions)
   - [ ] Agree on implementation approach

2. **Prepare Environment**
   - [ ] Set up staging environment (copy of production)
   - [ ] Test current backup/restore process
   - [ ] Document current production state

3. **Prototype**
   - [ ] Create proof-of-concept migration runner
   - [ ] Test backup/restore cycle
   - [ ] Test master data sync dev→staging

### Decision Required
**Choose one of:**
- [ ] **Recommended Approach** (Hybrid Migration Scripts + Enhanced Import)
- [ ] **Proposal 2** (Migration Scripts + GORM)
- [ ] **Proposal 3** (State-Based with Atlas/Goose)
- [ ] **Other** (specify)

---

## Appendix

### A. Data Type Classification

#### System Data (Production-Only, Never Overwrite)
- `users` - User accounts and authentication

#### Master Data (Dev→Prod Sync Required)
- `gsm_skills` - Game system skills
- `gsm_weaponskills` - Weapon skills
- `gsm_spells` - Spell definitions
- `gsm_equipments` - Equipment definitions
- `gsm_weapons` - Weapon definitions
- `gsm_containers` - Container definitions
- `gsm_transportations` - Transportation definitions
- `gsm_believes` - Belief systems
- `gsm_lit_sources` - Source books
- `gsm_character_classes` - Character classes
- `learning_skill_categories` - Skill categories
- `learning_skill_difficulties` - Difficulty levels
- `learning_spell_schools` - Spell schools
- `learning_class_category_ep_costs` - Learning costs
- `learning_class_spell_school_ep_costs` - Spell learning costs
- `learning_spell_level_le_costs` - LE costs
- `learning_skill_category_difficulties` - Relationships
- `learning_skill_improvement_costs` - Improvement costs

#### User Data (Production-Only, Preserve)
- `characters` - Player characters
- `eigenschaften` - Character attributes
- `lps`, `aps`, `bs` - Character stats
- `merkmale` - Character traits
- `erfahrungsschatze`, `bennies`, `vermoegen` - Character progression
- `sk_fertigkeiten`, `sk_waffenfertigkeiten`, `sk_zauber` - Character skills
- `eq_ausruestungen`, `eq_waffen`, `eq_containers` - Character equipment
- `character_creation_sessions` - Character creation state
- `audit_log_entries` - Audit trail

### B. Existing Export/Import Systems

#### 1. gsmaster/export_import.go
**Purpose**: Master data export/import with ID-independence  
**Scope**: Skills, spells, equipment, learning costs  
**Format**: JSON with natural keys  
**Usage**: Development data → Production

#### 2. transfer/database.go
**Purpose**: Full database backup/restore  
**Scope**: All tables  
**Format**: JSON with IDs  
**Usage**: Backup/restore, database migration

#### 3. importer/ (VTT format)
**Purpose**: Import characters from external formats  
**Scope**: Characters with all related data  
**Format**: VTT JSON  
**Usage**: External data import

### C. GORM Migration Behavior

**What AutoMigrate DOES**:
- ✅ Create missing tables
- ✅ Add missing columns
- ✅ Create indexes
- ✅ Create foreign keys
- ✅ Update column types (limited)

**What AutoMigrate DOES NOT**:
- ❌ Rename tables
- ❌ Rename columns
- ❌ Drop columns
- ❌ Drop tables
- ❌ Modify constraints
- ❌ Data migrations
- ❌ Complex index changes

### D. References

**Code Locations**:
- Migration entry point: `backend/models/database.go`
- Main migration runner: `backend/maintenance/handlers.go`
- Master data export: `backend/gsmaster/export_import.go`
- Full DB export: `backend/transfer/database.go`
- Docker config: `docker/docker-compose.yml`

**Documentation**:
- Export/Import: `backend/doc/EXPORT_IMPORT.md`
- Transfer module: `backend/transfer/README.md`
- Data transfer: `backend/doc/DATA_TRANSFER.md`

---

**End of Planning Document**

*This document should be reviewed and updated based on implementation experience and changing requirements.*
