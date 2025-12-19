package pdfrender

import (
	"fmt"

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
