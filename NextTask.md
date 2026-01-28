I want to remove the string "midgard" as a selection criterium for GameSystem (game_system) from the code.
It would be best if I could have a table for Game Systems and select the right records by this ID.
That means that we can show a list of available games systems in the frontend when we select or edit a char, skill, spell, weapon or other equipment. that must be handled to when exporting and importing data.
When exporting we must changethe export format and increase the format version. When importing we should add some kind of compatibility layer or preprocessing to ensure the previous format is imported right.
Create migration scripts for the database structure and data (I think it would be a good idea to add new fields transfer the data and only after that create foreign keys, if this is possible)

Use GORM for strutural migrations whenever possible. Update the migrateAll functions if needed
If you change function parameters provide a wrapper to keep old behaviour for older consumers.
Implemnt backend changes first
Make a concise and comprehensible plan for implementing this changes. Check the Code to see if I have missed szenarios we have to address.
Create a planning document.
If possible plan it step by step package by package.

Follow the TDD and KISS principles



add db version table schema_version -done
add api enpoint for migration info, start - unsure
for each step add migration methods
Create GameSystem struct                - done
Add GameSystemId to GSMData one by one update tests, export and import 
    - believes
add GameSystemID to Char update tests, export and import
add GameSystemID to Sources  update tests, export and import

