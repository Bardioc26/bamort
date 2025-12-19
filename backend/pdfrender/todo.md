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

## TODO (Remaining)

* Page 2 Weapons - Damage calculation:
  - Calculate and display weapon damage including:
    - Base weapon damage (from models.Weapon)
    - Character's SchadenBonus
    - Weapon's Schadensbonus (Schb from EqWaffe)
  - Format: e.g., "1W6+3" where +3 = SchadenBonus + Schb

* Page 2 Weapons - Ranged weapon ranges:
  - For ranged weapons (Bogen, Armbrust, Wurfwaffe), show:
    - Range for "Nah" (near)
    - Range for "Mittel" (medium)  
    - Range for "Fern" (far)
  - This data should come from models.Weapon.Range or similar field
.

* add 3 field to gsm_weapons that holds the ranges near, middle, far all values are measured in meters so integer seems to be a good datatype. If at least 1 of 3 values is > 0 the weapon is treated as a ranged weapon
* add a field to gsm_weapons that holds the damage the weapon creates Damage is notated as 2W6+3
* EqWaffe already contains the values for the bunus values we need Attack (Anb), Damage (Schb) und Defence (Abwb)

* currently the template fetched for rendering is set to Default_A4_Quer