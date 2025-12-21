## COMPLETED

* ✅ Page 2 Weapons Table now shows the following information:
  - For each Weapon a character has:
    - Its name
    - ✅ the Fertigkeitswert (EW): Waffenfertigkeit.Fertigkeitswert + Character.AngriffBonus + Weapon.Anb (if equipped)
    - ✅ the Schaden (Damage) including character's SchadenBonus and weapon's Schadensbonus
    - ✅ If it is a ranged weapon, the ranges for near, medium and far

  Implementation notes:
  - ✅ AngriffBonus and SchadenBonus are now calculated in DerivedValueSet from character attributes
  - ✅ mapWeapons() function now:
    - Calculates EW = Waffenfertigkeit.Fertigkeitswert + Character.AngriffBonus + Weapon.Anb
    - Matches Waffenfertigkeiten with equipped Waffen by name
    - Adds weapon attack bonus (Anb) if weapon is equipped
  - ✅ Added test TestMapWeapons_WithEWCalculation to verify correct EW calculation
  - ✅ Used TDD approach: wrote failing test first, then implemented solution

* ✅ Template MaxItems expectations are now dynamic:
  - Tests now read MaxItems values directly from template HTML comments
  - Updated TestLoadTemplateSetFromFiles to parse templates dynamically
  - Updated TestDefaultA4QuerTemplateSet_LoadsFromFiles to use dynamic expectations
  - Updated TestPaginationUsesTemplateMetadata to read from template files
  - Tests will automatically adapt when template capacities change

* ✅ Weapon model enhancements:
  - ✅ Added RangeNear, RangeMiddle, RangeFar integer fields to gsm_weapons table
  - ✅ Added IsRanged() method that returns true if at least one range value > 0
  - ✅ Damage field already exists in gsm_weapons table
  - ✅ EqWaffe already contains bonus values: Anb (Attack), Schb (Damage), Abwb (Defense)
  - ✅ All tests pass with new fields

* ✅ Page 2 Weapons Table - Complete implementation:
  - ✅ Changed TestVisualInspection_AllPages to load character Fanjo Vetrani (ID 18) from test database
  - ✅ Damage calculation implemented:
    - Calculates total damage: Base Weapon Damage + Character SchadenBonus + Weapon Schb
    - Format: "1W6+3" where +3 = SchadenBonus + Schb
    - Implemented calculateWeaponDamage() function
    - Added TestMapWeapons_WithDamageCalculation test
  - ✅ Ranged weapon ranges implemented:
    - Shows ranges for ranged weapons (Bogen, Armbrust, etc.)
    - Format: "Nah/Mittel/Fern" (e.g., "10/30/100")
    - Implemented calculateWeaponRange() function
    - Added TestMapWeapons_WithRangedWeaponRanges test
    - Marks weapons as ranged using IsRanged field
  - ✅ All tests pass

* ✅ Continuation pages for overflow items:
  - When items exceed template capacity, continuation pages are automatically created
  - Continuation pages follow naming pattern: page1.2_stats.html, page1.3_stats.html, etc.
  - Template loader automatically falls back to base template for continuation pages
  - No physical continuation template files needed - reuses base template structure
  - **NEW: RenderPageWithContinuations() function generates actual PDF files**
  - Each continuation page is rendered as a separate PDF
  - PDFs can be merged into a single combined file
  - Implemented using TDD:
    - Created comprehensive tests in continuation_test.go
    - Added GenerateContinuationTemplateName() function
    - Added ExtractBaseTemplateName() function
    - Updated paginateList() to generate continuation template names
    - Updated RenderTemplate() to handle continuation template fallback
    - **Created RenderPageWithContinuations() to actually render multiple PDFs**
    - **Created integration test that saves real PDF files to disk**
  - All existing tests updated to work with dynamic template capacities
  - Fully tested and working end-to-end
  - **VERIFIED: 5 continuation pages generated for 50 skills, saved to /tmp/bamort_continuation_test/**


* ✅ 1. create API endpoint for listing available export templates
  * Endpoint: GET /api/pdf/templates
  * Returns: JSON array of TemplateInfo objects [{id, name, description}]
  * Test: TestListTemplates passes
  * Configuration: Uses config.Cfg.TemplatesDir (default: "./templates")
* ✅ 2. create API endpoint for exporting character to PDF. 
  * Endpoint: GET /api/pdf/export/:id?template=xxx&showUserName=true
  * Endpoint takes parameter "template", "show user name". 
  * Return combined PDF file for download or display in Browser
  * Renders all 4 pages (stats, play, spells, equipment) with continuation pages
  * Merges all PDFs into single combined file
  * Returns PDF with proper headers: Content-Type: application/pdf, Content-Disposition
  * Tests: TestExportCharacterToPDF, TestExportCharacterToPDF_WithTemplate, TestExportCharacterToPDF_CharacterNotFound all pass
  * Configuration: Uses config.Cfg.TemplatesDir for template path resolution
  * Status: ✅ Deployed and running in Docker container, verified with logs
* ✅ 3. create exporting function in Frontend
  * ✅ The UI element to start the export function is to the left side of the character's name (CharacterDetails.vue)
  * ✅ Template selection dropdown implemented - auto-selects first template
  * ✅ Export button with loading state (disabled while exporting)
  * ✅ PDF opens in new browser tab using window.open()
  * ✅ Translations added for German and English
  * ✅ API integration: Fetches templates on component load, calls export endpoint with selected template
  * ✅ Error handling with user-friendly alerts
  * Status: ✅ Deployed with HMR, ready for testing

## TODO (Remaining)
* 1. create a directory xporttemp in the backend
* 2. save the PDF to file during the ExportCharacterToPDF call and return the filename charname+timestamp.pdf. Make shure the filename contains no spaces or special chars that might disturb the download
* 3. create an API endpoint to load the file from xporttemp
* 4. create a maintenance endpoint to clean up the xporttemp directory. remove all files that are older than 7 days
* 5. change the frontend to get the PDF from the new API endpoint

### Later
* continuation of lists does not work as expected but good enough for a first shot
  * generalize handling so that only on set of functions can handle ALL kinds of templates. Needs massive refactoring

* currently the template fetched for rendering is set to Default_A4_Quer
* remove inline css as far as possible
* make pdf download popup an own view