# Pagination and Template Metadata Fixes

## Summary
Fixed all issues from todo.md related to template metadata, pagination capacities, and empty row filling.

## Changes Made

### 1. Dynamic Capacity Loading (✓ COMPLETED)
**Problem**: PreparePaginatedPageData hardcoded capacity values that didn't match template MAX values.

**Solution**: 
- Added `GetBlockCapacity()` helper function to read MAX from template metadata
- Updated PreparePaginatedPageData() to dynamically load capacities for all pages
- Removed all hardcoded capacity values (24, 11, 5, 30, 20, 10, etc.)

**Files Modified**:
- `pagination_helper.go`: Added GetBlockCapacity(), updated all page handlers

### 2. Page2 Pagination Fixed (✓ COMPLETED)
**Problem**: 
- skills_learned used capacity 24 instead of template's MAX:18
- skills_languages used capacity 11 instead of template's MAX:5
- weapons_main capacity was inconsistent

**Solution**:
- Page2 now reads correct capacities from template:
  - skills_learned: MAX:18 (was 24)
  - skills_unlearned: MAX:15 (correct)
  - skills_languages: MAX:5 (was 11)
  - weapons_main: MAX:24 (was 30)

**Template Values** (from page2_play.html):
```html
<!-- BLOCK: skills_learned, TYPE: skills, MAX: 18, FILTER: learned -->
<!-- BLOCK: skills_unlearned, TYPE: skills, MAX: 15, FILTER: unlearned -->
<!-- BLOCK: skills_languages, TYPE: skills, MAX: 5, FILTER: language -->
<!-- BLOCK: weapons_main, TYPE: weapons, MAX: 24 -->
```

### 3. Page3 Magic Items Fixed (✓ COMPLETED)
**Problem**: magic_items used capacity 5 instead of template's MAX:8

**Solution**: Updated to read from template (MAX:8)

**Template Values** (from page3_spell.html):
```html
<!-- BLOCK: spells_left, TYPE: spells, MAX: 26 -->
<!-- BLOCK: spells_right, TYPE: spells, MAX: 15 -->
<!-- BLOCK: magic_items, TYPE: magicItems, MAX: 8 -->
```

### 4. Page4 Equipment Fixed (✓ COMPLETED)
**Problem**: Used wrong block name ("equipment" instead of "equipment_worn")

**Solution**: Updated to use correct block name from template

**Template Values** (from page4_equip.html):
```html
<!-- BLOCK: equipment_worn, TYPE: equipment, MAX: 10, FILTER: worn -->
```

### 5. Tests Updated (✓ COMPLETED)
**Problem**: Tests hardcoded expected MAX values instead of reading from templates

**Solution**: Updated all tests to dynamically read capacities:
- `TestPaginationUsesTemplateMetadata`: Verifies template parsing works
- `TestPage2PaginationWithCorrectCapacities`: Uses GetBlockCapacity()
- `TestPage3MagicItemsCapacity`: Uses GetBlockCapacity()
- `TestPreparePaginatedPageData_Page3Spell`: Updated expectations
- `TestPreparePaginatedPageData_Page4Equipment`: Fixed block name
- `TestCalculatePagesNeeded`: Updated to match template capacities (24 not 30, 41 not 30)
- `TestPaginateSpells_MultiPage`: Updated to 26+15=41
- `TestPaginateWeapons_MultiPage`: Updated to 24 capacity (3 pages for 50 weapons)
- `TestIntegration_TemplateMetadata`: Updated expected values
- `templates_test.go`: Updated TestGetTemplateMetadata expectations
- `template_metadata_loader_test.go`: Updated expected values

**Files Modified**:
- `todo_fixes_test.go`: New test file with TDD tests
- `pagination_helper_test.go`: Updated expectations
- `pagination_test.go`: Updated test cases
- `integration_test.go`: Updated expected block MAX values
- `templates_test.go`: Updated to expect 26/15 not 20/10
- `template_metadata_loader_test.go`: Updated to expect 26 not 20

## Test Results
```
70 tests passing
0 tests failing

All pdfrender tests pass:
- 13 pagination tests
- 8 template metadata tests
- 6 fill_capacity tests
- 9 mapper tests
- 17 integration tests
- 4 new todo_fixes tests
- 13 other tests
```

## Visual Verification
Generated test PDFs confirm:
- ✓ Page1: Skills split correctly across 2 columns (29+29)
- ✓ Page2: 18 learned skills, 5 language skills, 24 weapons
- ✓ Page3: 26 left spells, 15 right spells, 8 magic items
- ✓ Page4: 10 equipment items
- ✓ All empty rows render correctly
- ✓ Combined PDF merges all pages successfully

## Weapons Implementation Note
The current implementation correctly uses `Waffenfertigkeiten` (weapon skills) for the weapons_main list. Each weapon skill already contains:
- Name: Weapon name (e.g., "Schwert", "Bogen")
- Value (EW): Fertigkeitswert (skill effectiveness value)
- This is the correct data for character sheet display

The `Equipment.Weapons` (EqWaffe) contains physical weapon objects with different metadata (Abwb, Schb, weight, etc.) which is not needed for the weapons_main table on page2.

## Architecture Improvements
1. **Separation of Concerns**: Template metadata is now the single source of truth
2. **Maintainability**: Adding/changing template capacities only requires HTML comment updates
3. **Type Safety**: GetBlockCapacity() provides consistent interface
4. **Testability**: Tests verify against actual templates, not hardcoded values
5. **DRY Principle**: No duplication between template definitions and code

## Future Enhancements
All current issues resolved. System ready for:
- Multi-page pagination (already working for >58 skills, >24 weapons, >41 spells)
- Additional template blocks (just add HTML comments)
- Different template sets (already supported via LoadTemplateSetFromFiles)
