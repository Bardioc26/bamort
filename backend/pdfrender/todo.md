* ✓ Template meta data must NOT be hard coded (func DefaultA4QuerTemplateSet ) but must be read from theTemplate itself.
    OR if default values must be defined they must be overwritten from the statements fount in the template
    **COMPLETED**: DefaultA4QuerTemplateSet() now calls LoadTemplateSetFromFiles() which parses HTML comments.
    Falls back to hardcoded values if files can't be read.
* ✓ fill table with empty lines up to the my max value
    **COMPLETED**: Added FillToCapacity() function, integrated into PreparePaginatedPageData().
    All lists now filled to MAX capacity from template metadata.

* ✓ the tests for pagination must take the max values to test against from the template or from the metadata already updated by loading from template
    **COMPLETED**: All tests updated to read MAX values from templates dynamically using GetBlockCapacity().
    Tests no longer hardcode expected values.

* ✓ paginating does not work in page 2 for 
   <!-- BLOCK: skills_learned, TYPE: skills, MAX: 18, FILTER: learned -->
   <!-- BLOCK: skills_unlearned, TYPE: skills, MAX: 15, FILTER: unlearned -->
    **COMPLETED**: PreparePaginatedPageData() now reads capacities from template metadata using GetBlockCapacity().
    Page2 correctly uses MAX:18 for skills_learned, MAX:15 for skills_unlearned, MAX:5 for skills_languages, MAX:24 for weapons_main.

* ✓ FillToCapacity() does not work 
   in page 2 for <!-- BLOCK: skills_languages, TYPE: skills, MAX: 5, FILTER: language -->
   in page 3 for <!-- BLOCK: magic_items, TYPE: magicItems, MAX: 8 -->
    **COMPLETED**: Fixed PreparePaginatedPageData() to read correct MAX values from templates.
    skills_languages now uses MAX:5 (not 11), magic_items now uses MAX:8 (not 5).
    All blocks properly filled to capacity with empty rows.

* weapons_main list currently uses Waffenfertigkeiten (weapon skills).
  NOTE: Equipment.Weapons (EqWaffe) contains physical weapons with metadata like Abwb/Schb.
  Waffenfertigkeiten already provides EW (Fertigkeitswert) which is the skill value needed for the character sheet.
  Current implementation is correct - weapons_main shows weapon skills with their EW values.