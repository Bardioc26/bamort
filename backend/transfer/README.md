# Transfer Package
The transfer package handles exporting and importing of characters and complete database content.

## Features

### Character Export/Import
- Export single characters with all related data (skills, spells, equipment, etc.)
- Import characters from JSON with automatic GSMaster data reconciliation
- Download characters as JSON files

### Database Export/Import
- Export complete database to JSON file
- Import complete database from JSON file
- All database tables included (users, characters, skills, equipment, learning data, etc.)
- Stored in `./backend/export_temp` directory

## API Endpoints

### Character Operations
- `GET /api/transfer/export/:id` - Export character as JSON (API response)
- `GET /api/transfer/download/:id` - Download character as JSON file
- `POST /api/transfer/import` - Import character from JSON

### Database Operations
- `POST /api/transfer/database/export` - Export complete database
- `POST /api/transfer/database/import` - Import complete database

## Usage Examples

### Export Database
```bash
curl -X POST http://localhost:8180/api/transfer/database/export \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:
```json
{
  "message": "Database exported successfully",
  "filename": "database_export_20260101_120000.json",
  "filepath": "./backend/export_temp/database_export_20260101_120000.json",
  "record_count": 1234,
  "timestamp": "2026-01-01T12:00:00Z"
}
```

### Import Database
```bash
curl -X POST http://localhost:8180/api/transfer/database/import \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "filepath": "./backend/export_temp/database_export_20260101_120000.json"
  }'
```

Response:
```json
{
  "message": "Database imported successfully",
  "record_count": 1234,
  "timestamp": "2026-01-01T12:00:01Z"
}
```

## Technical Details

### Database Export Format
The export includes all tables:
- Users
- Characters (with all relations: Eigenschaften, Lps, Aps, etc.)
- Skills (SkFertigkeiten, SkWaffenfertigkeiten, SkZauber)
- Equipment (EqAusruestungen, EqWaffen, EqContainers)
- GSMaster data (Skills, WeaponSkills, Spells, Equipment, Weapons, etc.)
- Learning data (Sources, CharacterClasses, SkillCategories, etc.)
- Audit log entries

### Import Behavior
- Uses `Save()` for upsert logic (updates existing records, creates new ones)
- Maintains referential integrity
- Wrapped in transaction (all-or-nothing)

## Testing
All functionality is fully tested with TDD approach:
```bash
go test -v ./transfer/
```