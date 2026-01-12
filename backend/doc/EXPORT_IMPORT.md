# Master Data Export/Import

## Overview

The export/import mechanism allows exporting all master data from the `model_gsmaster` and `model_learning_costs` modules to JSON files and importing them back. The exported data is ID-independent, using natural keys (name + game_system) to identify records, making it suitable for:

- Migrating data between environments
- Version controlling master data
- Manually editing game data
- Sharing/distributing game systems

## Supported Entities

### From `model_learning_costs.go`:
- **Sources** (`gsm_lit_sources`) - Game books and source materials
- **SkillCategories** (`learning_skill_categories`) - Skill classification categories
- **SkillDifficulties** (`learning_skill_difficulties`) - Difficulty levels
- **SkillCategoryDifficulties** (`learning_skill_category_difficulties`) - Relationship between skills, categories, and difficulties with learning costs

### From `model_gsmaster.go`:
- **Skills** (`gsm_skills`) - Character skills
- **Spells** (`gsm_spells`) - Magic spells

## Excluded Entities

The following are NOT exported/imported:
- `AuditLogEntry` - Audit logs (transient data)
- `BamortBase` - Base character data (character-specific)
- `BamortCharTrait` - Character traits (character-specific)
- `Magisch` - Character magic data (character-specific)
- `LookupList` - System lookup lists (system config)

## File Format

All data is exported to JSON files with indentation for easy reading and editing:

```json
[
  {
    "name": "Skill Name",
    "game_system": "midgard",
    "source_code": "KOD",
    ...
  }
]
```

## Natural Keys

Instead of database IDs, the following natural keys are used to identify records:

| Entity | Natural Key |
|--------|------------|
| Source | `code` |
| Skill | `name` + `game_system` |
| Spell | `name` + `game_system` |
| SkillCategory | `name` + `game_system` |
| SkillDifficulty | `name` + `game_system` |
| SkillCategoryDifficulty | `skill_name` + `skill_system` + `category_name` + `category_system` + `difficulty_name` + `difficulty_system` |

## Usage

### Export

```go
import "bamort/gsmaster"

// Export specific entity types
err := gsmaster.ExportSources("/path/to/output")
err := gsmaster.ExportSkills("/path/to/output")
err := gsmaster.ExportSpells("/path/to/output")
err := gsmaster.ExportSkillCategories("/path/to/output")
err := gsmaster.ExportSkillDifficulties("/path/to/output")
err := gsmaster.ExportSkillCategoryDifficulties("/path/to/output")

// Export all master data at once
err := gsmaster.ExportAll("/path/to/output")
```

### Import

```go
import "bamort/gsmaster"

// Import specific entity types
err := gsmaster.ImportSources("/path/to/input")
err := gsmaster.ImportSkills("/path/to/input")
err := gsmaster.ImportSpells("/path/to/input")
err := gsmaster.ImportSkillCategories("/path/to/input")
err := gsmaster.ImportSkillDifficulties("/path/to/input")
err := gsmaster.ImportSkillCategoryDifficulties("/path/to/input")

// Import all master data at once
err := gsmaster.ImportAll("/path/to/input")
```

## Import Behavior

The import mechanism follows an "upsert" pattern:

1. **Check if record exists** using natural keys
2. If **not found**: Create new record
3. If **found**: Update existing record with imported values

This allows for:
- Importing new data
- Updating existing data
- Safe re-import of previously exported data

## Dependency Order

When using `ExportAll()` and `ImportAll()`, entities are processed in dependency order:

**Export/Import Order:**
1. Sources (no dependencies)
2. SkillCategories (depends on Sources)
3. SkillDifficulties (no dependencies)
4. Skills (depends on Sources)
5. SkillCategoryDifficulties (depends on Skills, Categories, Difficulties)
6. Spells (depends on Sources)

## File Names

Each entity type is exported to its own file:

| Entity | Filename |
|--------|----------|
| Sources | `sources.json` |
| Skills | `skills.json` |
| Spells | `spells.json` |
| SkillCategories | `skill_categories.json` |
| SkillDifficulties | `skill_difficulties.json` |
| SkillCategoryDifficulties | `skill_category_difficulties.json` |

## Example Workflow

### Exporting Data

```go
package main

import (
    "bamort/gsmaster"
    "log"
)

func main() {
    outputDir := "./exported_data"
    
    if err := gsmaster.ExportAll(outputDir); err != nil {
        log.Fatalf("Export failed: %v", err)
    }
    
    log.Println("All master data exported to", outputDir)
}
```

### Editing Exported Data

```bash
# Edit the JSON files manually
vim exported_data/skills.json
```

### Importing Modified Data

```go
package main

import (
    "bamort/gsmaster"
    "bamort/database"
    "log"
)

func main() {
    // Initialize database connection
    database.InitDB()
    
    inputDir := "./exported_data"
    
    if err := gsmaster.ImportAll(inputDir); err != nil {
        log.Fatalf("Import failed: %v", err)
    }
    
    log.Println("All master data imported from", inputDir)
}
```

## Error Handling

All export/import functions return errors that should be checked:

- **Export errors**: Usually file system issues (permissions, disk space)
- **Import errors**: Can be:
  - File not found
  - Invalid JSON format
  - Missing dependencies (e.g., referenced source doesn't exist)
  - Database constraint violations

## Testing

The export/import mechanism is fully tested with TDD. See `gsmaster/export_import_test.go` for comprehensive test coverage including:

- Export creates valid JSON files
- Import creates new records
- Import updates existing records
- Relationships are correctly restored using natural keys
- Full export/import cycle works correctly

## Notes

- **No Handlers/Routes**: This functionality is intentionally not exposed as API endpoints. Use it programmatically or via CLI tools.
- **ID Independence**: Exported data does not contain database IDs, making it portable across different database instances.
- **Idempotent**: Import can be run multiple times safely - it will update existing records rather than creating duplicates.
- **Transaction Safety**: Each import operation should ideally be wrapped in a database transaction for atomicity.
