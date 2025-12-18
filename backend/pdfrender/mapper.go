package pdfrender

import (
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
		Name:       char.Name,
		Type:       char.Typ,
		Grade:      char.Grad,
		Age:        char.Alter,
		Height:     char.Groesse,
		Weight:     char.Gewicht,
		Gender:     char.Gender,
		Homeland:   char.Herkunft,
		Religion:   char.Glaube,
		Stand:      char.SocialClass,
		IconBase64: "", // Will be set later if image exists
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
		case "pA":
			attrs.PA = e.Value
		case "Wk":
			attrs.Wk = e.Value
		}
	}

	attrs.B = char.B.Value

	return attrs
}

// mapDerivedValues extracts derived values like LP, AP, etc.
func mapDerivedValues(char *models.Char) DerivedValueSet {
	return DerivedValueSet{
		LPMax:     char.Lp.Max,
		LPAktuell: char.Lp.Value,
		APMax:     char.Ap.Max,
		APAktuell: char.Ap.Value,
	}
}

// mapSkills converts character skills to SkillViewModel
func mapSkills(char *models.Char) []SkillViewModel {
	skills := make([]SkillViewModel, 0, len(char.Fertigkeiten))

	for _, skill := range char.Fertigkeiten {
		skills = append(skills, SkillViewModel{
			Name:           skill.Name,
			Category:       skill.Category,
			Value:          skill.Fertigkeitswert,
			Bonus:          skill.Bonus,
			PracticePoints: skill.Pp,
			IsLearned:      skill.Fertigkeitswert > 0,
		})
	}

	return skills
}

// mapWeapons converts character weapon skills to WeaponViewModel
func mapWeapons(char *models.Char) []WeaponViewModel {
	weapons := make([]WeaponViewModel, 0, len(char.Waffenfertigkeiten))

	for _, weapon := range char.Waffenfertigkeiten {
		weapons = append(weapons, WeaponViewModel{
			Name:  weapon.Name,
			Value: weapon.Fertigkeitswert,
		})
	}

	return weapons
}

// mapSpells converts character spells to SpellViewModel
func mapSpells(char *models.Char) []SpellViewModel {
	spells := make([]SpellViewModel, 0, len(char.Zauber))

	for _, spell := range char.Zauber {
		spells = append(spells, SpellViewModel{
			Name: spell.Name,
		})
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
