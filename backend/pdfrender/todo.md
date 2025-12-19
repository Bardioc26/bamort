* weapons_main list currently uses Waffenfertigkeiten (weapon skills).
  NOTE: Equipment.Weapons (EqWaffe) contains physical weapons with metadata like Abwb/Schb.
  Waffenfertigkeiten already provides EW (Fertigkeitswert) which is the skill value needed for the character sheet.
  Current implementation is correct - weapons_main shows weapon skills with their EW values.
  * The implementation is NOT correct, because the eapons_main shows weapon skills. But it should show the weapons list from the equipment