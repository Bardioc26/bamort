# Transfer Package

The transfer package provides character export and import functionality for the Bamort application.

## Features

- **Complete Character Export**: Exports a character with all related data including:
  - Basic character information (race, type, grade, attributes, etc.)
  - Skills (Fertigkeiten) and Weapon Skills (Waffenfertigkeiten)
  - Spells (Zauber)
  - Equipment (Waffen, Ausrüstung, Behältnisse, Transportmittel)
  - Learning data (all learning_* tables)
  - Audit log entries
  - GSM master data (skills, spells, weapons, equipment definitions)
  - **Security**: Sensitive user data (password, updated_at, reset tokens) is removed from exports

- **Smart Import**: 
  - Identifies existing GSM data by name (not ID) to avoid duplicates
  - Updates incomplete GSM records with missing information
  - Sets default source_id values (1 for skills/equipment, 2 for spells)
  - Preserves audit log history
  - Creates new character with fresh IDs

- **JSON Format**: All data is exported/imported as JSON for portability

## API Endpoints

All endpoints are under `/api/transfer`:

### GET /api/transfer/export/:id
Exports a character as JSON data (for API consumption)

**Response**: CharacterExport JSON object

### GET /api/transfer/download/:id
Downloads a character as a JSON file

**Response**: JSON file with `Content-Disposition: attachment` header

### POST /api/transfer/import
Imports a character from JSON data

**Request Body**: CharacterExport JSON object
**Response**: 
```json
{
  "message": "Character imported successfully",
  "character_id": 123
}
```

## Usage Examples

### Export a Character
```bash
curl http://localhost:8180/api/transfer/export/18 \
  -H "Authorization: Bearer XXXXXXXXXXX" \
  > character_export.json
```

### Download a Character
```bash
wget http://localhost:8180/api/transfer/download/18
```

### Import a Character
```bash
curl -X POST http://localhost:8180/api/transfer/import \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer XXXXXXXXXXX" \
  -d @character_export.json
```

## Testing

The package includes comprehensive tests:

- **Export Tests** (7 tests):
  - Basic export functionality
  - Skills inclusion
  - Spells inclusion
  - Equipment inclusion
  - Learning data inclusion
  - Audit log inclusion
  - Error handling for non-existent characters

- **Import Tests** (6 tests):
  - Basic import functionality
  - Handling existing GSM data
  - Updating incomplete GSM data
  - Default source_id assignment
  - Audit log import
  - JSON round-trip testing

- **Handler/API Tests** (5 tests):
  - Export endpoint
  - Download endpoint
  - Import endpoint
  - Error handling

**Total: 18 tests, all passing**

Run tests with:
```bash
cd backend
go test -v ./transfer/
```

## Implementation Notes

### TDD Approach
The package was developed using Test-Driven Development (TDD):
1. Tests written first
2. Implementation follows to make tests pass
3. Refactoring as needed

### KISS Principle
The implementation follows the "Keep It Simple, Stupid" principle:
- Simple, clear function names
- Each function does one thing well
- No over-engineering
- Straightforward error handling

### GSM Data Handling
- Skills, weapons, equipment, and spells are identified by **name**, not ID
- Prevents duplicate creation of master data
- Updates existing records only if they have missing information
- Default source_id: 1 for general data, 2 for spells

### Source ID Rules
When importing, if `source_id` is 0:
- **Spells**: Set to 2
- **All other data**: Set to 1

## Files

- `exporter.go` - Character export functionality
- `exporter_test.go` - Export tests
- `importer.go` - Character import functionality
- `importer_test.go` - Import tests
- `handlers.go` - HTTP handlers for API endpoints
- `handlers_test.go` - API handler tests
- `routes.go` - Route registration
- `README.md` - This file

## Future Enhancements

Possible improvements:
- Batch export/import of multiple characters
- Export filtering (e.g., export without audit log)
- Import validation and conflict resolution options
- Export format versioning
- Compression for large exports
