## COMPLETED

* ✅ Page 2 Weapons Table now shows the following information:
  - For each Weapon a character has:
    - Its name
    - ✅ the Fertigkeitswert (EW): Waffenfertigkeit.Fertigkeitswert + Character.AngriffBonus + Weapon.Anb (if equipped)
    - ❌ TODO: the Schaden (Damage) including character's SchadenBonus and weapon's Schadensbonus
    - ❌ TODO: If it is a ranged weapon, the ranges for near, medium and far

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

## TODO (Remaining)

* currently the template fetched for rendering is set to Default_A4_Quer
* page 2 Waffen Schaden müssen noch die Boni der Waffe und der Schadensbonus des Characters eingerechnet werden.
  Dazu muss vom (1W6)(+1) der zweite teil extrahiert werden, falls nicht vorhanden als +0 annehmen. auf diesen zweiten Teil wird der Bonus aufgerechnet