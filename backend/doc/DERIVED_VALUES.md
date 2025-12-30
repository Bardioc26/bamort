# Derived Values Implementation

## Summary

Static derived values (Resistenz Körper, Resistenz Geist, Abwehr, Zaubern, Raufen) are now stored in the database as part of the `Char` model. These values can increase when a new grade is reached.

Bonus values are calculated on-demand from the character's attributes using the `CalculateBonuses()` method.

## Database Fields Added to `Char`

- `ResistenzKoerper` (int) - Resistenz Körper base value
- `ResistenzGeist` (int) - Resistenz Geist base value
- `Abwehr` (int) - Abwehr base value
- `Zaubern` (int) - Zaubern base value
- `Raufen` (int) - Raufen base value

## Calculated Bonuses (Not Stored)

The following bonuses are calculated from attributes and are NOT stored in the database:

- `AusdauerBonus` = Ko/10 + St/20
- `SchadensBonus` = St/20 + Gs/30 - 3
- `AngriffsBonus` = attribute bonus from Gs
- `AbwehrBonus` = attribute bonus from Gw
- `ZauberBonus` = attribute bonus from Zt
- `ResistenzBonusKoerper` = attribute/race bonus from Ko
- `ResistenzBonusGeist` = attribute/race bonus from In

## Usage Example

```go
// Load character from database
var char models.Char
err := char.FirstID("123")

// Get static derived values (stored in DB)
resistenzKoerper := char.ResistenzKoerper
abwehr := char.Abwehr

// Calculate bonuses from attributes (on-demand)
bonuses := char.CalculateBonuses()
abwehrBonus := bonuses.AbwehrBonus
zauberBonus := bonuses.ZauberBonus

// Total values
totalAbwehr := char.Abwehr + bonuses.AbwehrBonus
totalZaubern := char.Zaubern + bonuses.ZauberBonus
```

## When Grade Increases

When a character reaches a new grade, the static derived values may increase. To update them:

```go
// After grade increase, recalculate and update static values
char.Grad++
char.ResistenzKoerper++ // or apply the appropriate grade bonus
char.Abwehr++
// Save updated character
database.DB.Save(&char)
```

Bonuses will automatically reflect any attribute changes without needing to update the database, as they are calculated on-demand from the `Eigenschaften` values.
