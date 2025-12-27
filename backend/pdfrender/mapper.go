package pdfrender

import (
	"fmt"

	"bamort/character"
	"bamort/database"
	"bamort/models"
)

// MapCharacterToViewModel converts a domain Character to a CharacterSheetViewModel
func MapCharacterToViewModel(char *models.Char) (*CharacterSheetViewModel, error) {
	vm := &CharacterSheetViewModel{
		Skills:      make([]SkillViewModel, 0),
		Weapons:     make([]WeaponViewModel, 0),
		Spells:      make([]SpellViewModel, 0),
		MagicItems:  make([]MagicItemViewModel, 0),
		Equipment:   make([]EquipmentViewModel, 0),
		GameResults: make([]GameResultViewModel, 0),
	}

	// Map basic character info
	vm.Character = CharacterInfo{
		Name:        char.Name,
		Type:        char.Typ,
		Grade:       char.Grad,
		Age:         char.Alter,
		Height:      char.Groesse,
		Weight:      char.Gewicht,
		Gender:      char.Gender,
		Hand:        char.Hand,
		Herkunft:    char.Herkunft,
		Glaube:      char.Glaube,
		SocialClass: char.SocialClass,
		IconBase64:  "", // Will be set later if image exists
		Vermoegen: WealthInfo{
			Goldstuecke:   char.Vermoegen.Goldstuecke,
			Silberstuecke: char.Vermoegen.Silberstuecke,
			Kupferstuecke: char.Vermoegen.Kupferstuecke,
		},
	}

	// Map attributes
	vm.Attributes = mapAttributes(char)

	// Map derived values
	vm.DerivedValues = mapDerivedValues(char)

	// Map skills
	vm.Skills = mapSkills(char)

	// Map weapons
	vm.Weapons = mapWeapons(char)

	// Map spells
	vm.Spells = mapSpells(char)

	// Map equipment
	vm.Equipment = mapEquipment(char)

	return vm, nil
}

// mapAttributes extracts attribute values from Eigenschaften slice
func mapAttributes(char *models.Char) AttributeValues {
	attrs := AttributeValues{}

	for _, e := range char.Eigenschaften {
		switch e.Name {
		case "St":
			attrs.St = e.Value
		case "Gs":
			attrs.Gs = e.Value
		case "Gw":
			attrs.Gw = e.Value
		case "Ko":
			attrs.Ko = e.Value
		case "In":
			attrs.In = e.Value
		case "Zt":
			attrs.Zt = e.Value
		case "Au":
			attrs.Au = e.Value
		case "PA":
			attrs.PA = e.Value
		case "Wk":
			attrs.Wk = e.Value
		}
	}

	attrs.B = char.B.Max

	return attrs
}

// mapDerivedValues extracts derived values like LP, AP, etc.
func mapDerivedValues(char *models.Char) DerivedValueSet {
	// Get attributes for bonus calculations
	attrs := mapAttributes(char)
	// Get Bonus values from character logic
	bonusValues := character.CalculateStaticFieldsLogic(character.CalculateStaticFieldsRequest{
		St:    attrs.St,
		Gs:    attrs.Gs,
		Gw:    attrs.Gw,
		Ko:    attrs.Ko,
		In:    attrs.In,
		Zt:    attrs.Zt,
		Au:    attrs.Au,
		Rasse: char.Rasse,
		Typ:   char.Typ,
		Grad:  char.Grad,
	})

	return DerivedValueSet{
		LPMax:                 char.Lp.Max,
		LPAktuell:             char.Lp.Value,
		APMax:                 char.Ap.Max,
		APAktuell:             char.Ap.Value,
		AusdauerBonus:         bonusValues.AusdauerBonus,
		SchadenBonus:          bonusValues.SchadensBonus,
		AngriffBonus:          bonusValues.AngriffsBonus,
		Abwehr:                bonusValues.Abwehr,
		AbwehrBonus:           bonusValues.AbwehrBonus,
		Zaubern:               bonusValues.Zaubern,
		ZauberBonus:           bonusValues.ZauberBonus,
		ResistenzKoerper:      bonusValues.ResistenzKoerper,
		ResistenzBonusKoerper: bonusValues.ResistenzBonusKoerper,
		ResistenzGeist:        bonusValues.ResistenzGeist,
		ResistenzBonusGeist:   bonusValues.ResistenzBonusGeist,
		Raufen:                bonusValues.Raufen,
		GG:                    char.B.Max,
		SG:                    char.B.Max / 2,
		Sehen:                 0, // TODO: Add to character model
		Horen:                 0, // TODO: Add to character model
		Riechen:               0, // TODO: Add to character model
		Sechster:              0, // TODO: Add to character model
		RKMalus:               0,
		RKSave:                2,
		RKShort:               "LR",
		RKName:                "Lederrüstung",
	}
}

// mapSkills converts character skills to SkillViewModel
func mapSkills(char *models.Char) []SkillViewModel {
	skills := make([]SkillViewModel, 0, len(char.Fertigkeiten))

	for _, skill := range char.Fertigkeiten {
		skl := SkillViewModel{
			Name:           skill.Name,
			Category:       skill.Category,
			Value:          skill.Fertigkeitswert,
			Bonus:          skill.Bonus,
			PracticePoints: skill.Pp,
			IsLearned:      skill.Fertigkeitswert > 0,
			Bemerkung:      skill.Bemerkung,
		}
		if skl.Name == "Schreiben" || skl.Name == "Sprache" {
			skl.Category = "Sprache"
		}
		skills = append(skills, skl)
	}

	return skills
}

// mapWeapons converts equipped weapons to WeaponViewModel
// EW = Waffenfertigkeit.Fertigkeitswert + Character.AngriffBonus + Weapon.Anb
func mapWeapons(char *models.Char) []WeaponViewModel {
	weapons := make([]WeaponViewModel, 0, len(char.Waffen))

	// Calculate character's bonuses using character logic
	attrs := mapAttributes(char)
	bonusValues := character.CalculateStaticFieldsLogic(character.CalculateStaticFieldsRequest{
		St:    attrs.St,
		Gs:    attrs.Gs,
		Gw:    attrs.Gw,
		Ko:    attrs.Ko,
		In:    attrs.In,
		Zt:    attrs.Zt,
		Au:    attrs.Au,
		Rasse: char.Rasse,
		Typ:   char.Typ,
		Grad:  char.Grad,
	})

	// Create a map of weapon skills for quick lookup
	weaponSkills := make(map[string]int)
	for _, skill := range char.Waffenfertigkeiten {
		weaponSkills[skill.Name] = skill.Fertigkeitswert
	}

	// Iterate over equipped weapons
	for _, equippedWeapon := range char.Waffen {
		vm := WeaponViewModel{
			Name: equippedWeapon.Name,
		}

		// Load weapon from gsm_weapons to get base stats and required skill
		baseWeapon := &models.Weapon{}
		err := baseWeapon.First(equippedWeapon.Name)

		if err == nil && baseWeapon.ID > 0 {
			// Calculate attack value: skill + character bonus + weapon bonus
			skillValue := 0
			if baseWeapon.SkillRequired != "" {
				skillValue = weaponSkills[baseWeapon.SkillRequired]
			}
			vm.Value = skillValue + bonusValues.AngriffsBonus + equippedWeapon.Anb

			// Calculate damage: Base weapon damage + character bonus + weapon damage bonus
			vm.Damage = calculateWeaponDamageWithBase(baseWeapon, bonusValues.SchadensBonus, equippedWeapon.Schb)

			// Add range information for ranged weapons
			if baseWeapon.IsRanged() {
				vm.Range = fmt.Sprintf("%d/%d/%d",
					baseWeapon.RangeNear,
					baseWeapon.RangeMiddle,
					baseWeapon.RangeFar)
				vm.IsRanged = true
			}
		} else {
			// Weapon not found in gsm_weapons, use basic info
			vm.Value = bonusValues.AngriffsBonus + equippedWeapon.Anb
		}

		weapons = append(weapons, vm)
	}

	return weapons
}

// mapSpells converts character spells to SpellViewModel
func mapSpells(char *models.Char) []SpellViewModel {
	spells := make([]SpellViewModel, 0, len(char.Zauber))

	for _, charSpell := range char.Zauber {
		vm := SpellViewModel{
			Name:  charSpell.Name,
			Bonus: charSpell.Bonus,
		}

		// Try to load spell details from gsmaster (only if database is available)
		// In test environments without DB, we skip this enrichment
		if database.DB != nil {
			masterSpell := &models.Spell{}
			err := masterSpell.First(charSpell.Name)

			// If master spell found, add all details
			if err == nil && masterSpell.ID > 0 {
				vm.Stufe = masterSpell.Stufe
				vm.AP = masterSpell.AP
				vm.Art = masterSpell.Art
				vm.Zauberdauer = masterSpell.Zauberdauer
				vm.Reichweite = masterSpell.Reichweite
				vm.Wirkungsziel = masterSpell.Wirkungsziel
				vm.Wirkungsbereich = masterSpell.Wirkungsbereich
				vm.Wirkungsdauer = masterSpell.Wirkungsdauer
				vm.Ursprung = masterSpell.Ursprung
				vm.Category = masterSpell.Category
				vm.LearningCategory = masterSpell.LearningCategory
				vm.Beschreibung = masterSpell.Beschreibung

				// If description is empty, use source code and page number
				if vm.Beschreibung == "" && masterSpell.SourceID > 0 && masterSpell.PageNumber > 0 {
					source := &models.Source{}
					if err := database.DB.First(source, masterSpell.SourceID).Error; err == nil {
						vm.Beschreibung = source.Code + " S." + fmt.Sprintf("%d", masterSpell.PageNumber)
					}
				}
			}
		}
		// If spell details not found or DB not available, just use name and bonus from character

		spells = append(spells, vm)
	}

	return spells
}

// mapEquipment converts character equipment to EquipmentViewModel
func mapEquipment(char *models.Char) []EquipmentViewModel {
	equipment := make([]EquipmentViewModel, 0, len(char.Ausruestung)+len(char.Behaeltnisse))

	// Add regular equipment
	for _, item := range char.Ausruestung {
		equipment = append(equipment, EquipmentViewModel{
			Name:        item.Name,
			Quantity:    item.Anzahl,
			Weight:      item.Gewicht,
			TotalWeight: item.Gewicht * float64(item.Anzahl),
			Value:       int(item.Wert),
			Location:    item.BeinhaltetIn,
			Container:   item.BeinhaltetIn,
			IsWorn:      item.BeinhaltetIn == "Am Körper",
			IsContainer: false,
		})
	}

	// Add containers
	for _, container := range char.Behaeltnisse {
		equipment = append(equipment, EquipmentViewModel{
			Name:        container.Name,
			Quantity:    1,
			Weight:      container.Gewicht,
			TotalWeight: container.Gewicht,
			Value:       int(container.Wert),
			IsContainer: true,
		})
	}

	return equipment
}

// calculateAttributeBonus calculates attribute bonus based on value
// Same logic as in character/derived_values_calculator.go
func calculateAttributeBonus(value int) int {
	if value >= 1 && value <= 5 {
		return -2
	} else if value >= 6 && value <= 20 {
		return -1
	} else if value >= 21 && value <= 80 {
		return 0
	} else if value >= 81 && value <= 95 {
		return 1
	} else if value >= 96 && value <= 100 {
		return 2
	}
	return 0
}

// calculateWeaponDamageWithBase calculates the total damage string using an already-loaded weapon
// Format: BaseDamage+TotalBonus, e.g., "1W6+3"
// TotalBonus = Character's SchadenBonus + Weapon's Schadensbonus (Schb)
func calculateWeaponDamageWithBase(baseWeapon *models.Weapon, schadenBonus int, weaponSchb int) string {
	baseDamage := baseWeapon.Damage
	if baseDamage == "" {
		return ""
	}

	// Calculate total damage bonus
	totalBonus := schadenBonus + weaponSchb

	// Format the damage string
	if totalBonus > 0 {
		return fmt.Sprintf("%s+%d", baseDamage, totalBonus)
	} else if totalBonus < 0 {
		return fmt.Sprintf("%s%d", baseDamage, totalBonus)
	}

	return baseDamage
}

// calculateWeaponDamage calculates the total damage string for a weapon
// Format: BaseDamage+TotalBonus, e.g., "1W6+3"
// TotalBonus = Character's SchadenBonus + Weapon's Schadensbonus (Schb)
func calculateWeaponDamage(weaponName string, schadenBonus int, weaponSchb int) string {
	// Try to load weapon details from gsmaster to get base damage
	var baseDamage string

	if database.DB != nil {
		masterWeapon := &models.Weapon{}
		err := masterWeapon.First(weaponName)
		if err == nil && masterWeapon.ID > 0 && masterWeapon.Damage != "" {
			return calculateWeaponDamageWithBase(masterWeapon, schadenBonus, weaponSchb)
		}
	}

	return baseDamage
}

// calculateWeaponRange returns the range string for a weapon and whether it's ranged
// Format: "Nah/Mittel/Fern", e.g., "10/30/100"
// Returns (rangeString, isRanged)
func calculateWeaponRange(weaponName string) (string, bool) {
	// Try to load weapon details from gsmaster to get ranges
	if database.DB != nil {
		masterWeapon := &models.Weapon{}
		err := masterWeapon.First(weaponName)
		if err == nil && masterWeapon.ID > 0 {
			// Check if weapon is ranged (at least one range value > 0)
			if masterWeapon.IsRanged() {
				// Format as "Nah/Mittel/Fern"
				rangeStr := fmt.Sprintf("%d/%d/%d",
					masterWeapon.RangeNear,
					masterWeapon.RangeMiddle,
					masterWeapon.RangeFar)
				return rangeStr, true
			}
		}
	}

	return "", false
}
