package importer

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testChar(t *testing.T, object *CharacterImport) {
	assert.Equal(t, "moam-character-41421", object.ID)
	assert.Equal(t, "Harsk Hammerhuter, Zen", object.Name)
	assert.Equal(t, "Zwerg", object.Rasse)
	assert.Equal(t, "Krieger", object.Typ)
	assert.Equal(t, 39, object.Alter)
	assert.Equal(t, "er", object.Anrede)
	assert.Equal(t, 3, object.Grad)
	assert.Equal(t, 140, object.Groesse)
	assert.Equal(t, 82, object.Gewicht)
	assert.Equal(t, "Torkin", object.Glaube)
	assert.Equal(t, "rechts", object.Hand)
	assert.Equal(t, 17, object.Lp.Max)
	assert.Equal(t, 17, object.Lp.Value)
	assert.Equal(t, 31, object.Ap.Max)
	assert.Equal(t, 31, object.Ap.Value)
	assert.Equal(t, 18, object.B.Max)
	assert.Equal(t, 0, object.B.Value)
	assert.Equal(t, 74, object.Eigenschaften.Au)
	assert.Equal(t, 96, object.Eigenschaften.Gs)
	assert.Equal(t, 70, object.Eigenschaften.Gw)
	assert.Equal(t, 65, object.Eigenschaften.In)
	assert.Equal(t, 85, object.Eigenschaften.Ko)
	assert.Equal(t, 75, object.Eigenschaften.Pa)
	assert.Equal(t, 95, object.Eigenschaften.St)
	assert.Equal(t, 71, object.Eigenschaften.Wk)
	assert.Equal(t, 35, object.Eigenschaften.Zt)
	assert.Equal(t, "blau", object.Merkmale.Augenfarbe)
	assert.Equal(t, "sandfarben", object.Merkmale.Haarfarbe)
	assert.Equal(t, "", object.Merkmale.Sonstige)
	assert.Equal(t, 0, object.Bennies.Gg)
	assert.Equal(t, 1, object.Bennies.Sg)
	assert.Equal(t, 0, object.Bennies.Gp)
	assert.Equal(t, "breit", object.Gestalt.Breite)
	assert.Equal(t, "klein", object.Gestalt.Groesse)
	assert.Equal(t, 325, object.Erfahrungsschatz.Value)
	assert.Equal(t, 3, len(object.Spezialisierung))
	assert.Equal(t, "Kriegshammer", object.Spezialisierung[0])
	assert.Equal(t, "Armbrust:schwer", object.Spezialisierung[1])
	assert.Equal(t, "Stielhammer", object.Spezialisierung[2])
	assert.Contains(t, object.Image, "data:image;base64,")

}

func testSkill(t *testing.T, objects []Fertigkeit) {
	assert.Equal(t, 19, len(objects))
	i := 0
	assert.Equal(t, "moam-ability-horen", objects[i].ID)
	assert.Equal(t, "Hören", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 6, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 99", objects[i].Quelle)
	i = 6
	assert.Equal(t, "moam-ability-759918", objects[i].ID)
	assert.Equal(t, "Athletik", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 9, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 104", objects[i].Quelle)
	i = 16
	assert.Equal(t, "moam-ability-759920", objects[i].ID)
	assert.Equal(t, "Sprache", objects[i].Name)
	assert.Equal(t, "Albisch", objects[i].Beschreibung)
	assert.Equal(t, 8, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 127", objects[i].Quelle)

}

func testWeaponSkill(t *testing.T, objects []Waffenfertigkeit) {
	assert.Equal(t, 8, len(objects))
	i := 0
	assert.Equal(t, "moam-ability-759916", objects[i].ID)
	assert.Equal(t, "Armbrüste", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 8, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 144", objects[i].Quelle)
	i = 2
	assert.Equal(t, "moam-ability-759912", objects[i].ID)
	assert.Equal(t, "Schilde", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 3, objects[i].Fertigkeitswert)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, 0, objects[i].Pp)
	assert.Equal(t, "KOD5 145", objects[i].Quelle)

}

func testSpell(t *testing.T, objects []Zauber) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-spell-134630", objects[i].ID)
	assert.Equal(t, "Angst", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 0, objects[i].Bonus)
	assert.Equal(t, "ARK5 63", objects[i].Quelle)
}

func testWeapon(t *testing.T, objects []Waffe) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-weapon-126819", objects[i].ID)
	assert.Equal(t, "Armbrust:schwer", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 0, objects[i].Abwb)
	assert.Equal(t, 0, objects[i].Anb)
	assert.Equal(t, 1, objects[i].Anzahl)
	assert.Equal(t, "moam-container-47363", objects[i].BeinhaltetIn)
	assert.Equal(t, 5.0, objects[i].Gewicht)
	assert.Equal(t, false, objects[i].Magisch.IstMagisch)
	assert.Equal(t, 0, objects[i].Magisch.Abw)
	assert.Equal(t, false, objects[i].Magisch.Ausgebrannt)
	assert.Equal(t, "Armbrust:schwer", objects[i].NameFuerSpezialisierung)
	assert.Equal(t, 0, objects[i].Schb)
	assert.Equal(t, 40.0, objects[i].Wert)

}

func testEquipment(t *testing.T, objects []Ausruestung) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-armor-48616", objects[i].ID)
	assert.Equal(t, "Lederrüstung", objects[i].Name)
	assert.Equal(t, "", objects[i].Beschreibung)
	assert.Equal(t, 1, objects[i].Anzahl)
	assert.Equal(t, "", objects[i].BeinhaltetIn)
	assert.Equal(t, 13.0, objects[i].Gewicht)
	assert.Equal(t, "", objects[i].BeinhaltetIn)
	assert.Equal(t, false, objects[i].Magisch.IstMagisch)
	assert.Equal(t, 0, objects[i].Magisch.Abw)
	assert.Equal(t, false, objects[i].Magisch.Ausgebrannt)
	assert.Equal(t, 30.0, objects[i].Wert)
	assert.Equal(t, 0, objects[i].Bonus)
}

func testContainer(t *testing.T, objects []Behaeltniss) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-container-47363", objects[i].ID)
	assert.Equal(t, "Lederrucksack", objects[i].Name)
	assert.Equal(t, "für 25 kg", objects[i].Beschreibung)
	assert.Equal(t, 4.0, objects[i].Wert)
	assert.Equal(t, 0.50, objects[i].Gewicht)
	assert.Equal(t, 25.0, objects[i].Volumen)
	assert.Equal(t, 25.0, objects[i].Tragkraft)
	assert.Empty(t, "", objects[i].BeinhaltetIn) //Value in json is null
	assert.Equal(t, false, objects[i].Magisch.IstMagisch)
	assert.Equal(t, 0, objects[i].Magisch.Abw)
	assert.Equal(t, false, objects[i].Magisch.Ausgebrannt)

}

func testTransportation(t *testing.T, objects []Transportation) {
	assert.Equal(t, 1, len(objects))
	i := 0
	assert.Equal(t, "moam-container-47000", objects[i].ID)
	assert.Equal(t, "Karren", objects[i].Name)
	assert.Equal(t, "für 250 kg", objects[i].Beschreibung)
	assert.Equal(t, 14.0, objects[i].Wert)
	assert.Equal(t, 40, objects[i].Gewicht)
	assert.Equal(t, 250.0, objects[i].Tragkraft)
	assert.Empty(t, "", objects[i].BeinhaltetIn) //Value in json is null
	assert.Equal(t, false, objects[i].Magisch.IstMagisch)
	assert.Equal(t, 0, objects[i].Magisch.Abw)
	assert.Equal(t, false, objects[i].Magisch.Ausgebrannt)

}

func TestImportVTTStructure(t *testing.T) {
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	testChar(t, character)
	testSkill(t, character.Fertigkeiten)
	testWeaponSkill(t, character.Waffenfertigkeiten)
	testSpell(t, character.Zauber)
	testEquipment(t, character.Ausruestung)
	testWeapon(t, character.Waffen)
	testContainer(t, character.Behaeltnisse)
	testTransportation(t, character.Transportmittel)
	//fmt.Println(character)
}

func TestImportSkill2GSMaster(t *testing.T) {

	database.SetupTestDB()

	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	skill, erro := TransformImportFertigkeit2GSDMaster(&character.Fertigkeiten[0])

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(skill.ID), 1)
	assert.Equal(t, "Hören", skill.Name)
	assert.Equal(t, "", skill.Beschreibung)
	assert.Equal(t, 6, skill.Initialwert)
	assert.Equal(t, "check", skill.Bonuseigenschaft)
	assert.Equal(t, "KOD5 99", skill.Quelle)
	assert.Equal(t, false, skill.Improvable)
	assert.Equal(t, "midgard", skill.GameSystem)
	//}
	skill2 := models.Skill{}
	erro = skill2.First("Hören")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(skill.ID))

	skill3 := models.Skill{}
	erro = skill3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Hören", skill3.Name)

	assert.Equal(t, skill2.ID, skill3.ID)
	assert.Equal(t, skill2.Name, skill3.Name)
	assert.Equal(t, skill2.Beschreibung, skill3.Beschreibung)
	assert.Equal(t, skill2.Initialwert, skill3.Initialwert)
	assert.Equal(t, skill2.Bonuseigenschaft, skill3.Bonuseigenschaft)
	assert.Equal(t, skill2.Quelle, skill3.Quelle)
	assert.Equal(t, skill2.Improvable, skill3.Improvable)
	assert.Equal(t, skill2.GameSystem, skill3.GameSystem)

	err = CheckFertigkeiten2GSMaster(character.Fertigkeiten)
	assert.NoError(t, err, "Expected no error when checkimg Skills against gsmaster")
}

func TestImportWeaponSkill2GSMaster(t *testing.T) {
	database.SetupTestDB()
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	skill, erro := TransformImportWaffenFertigkeit2GSDMaster(&character.Waffenfertigkeiten[0])

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(skill.ID), 1)
	assert.Equal(t, "Armbrüste", skill.Name)
	assert.Equal(t, "", skill.Beschreibung)
	assert.Equal(t, 8, skill.Initialwert)
	assert.Equal(t, "check", skill.Bonuseigenschaft)
	assert.Equal(t, "KOD5 144", skill.Quelle)
	assert.Equal(t, true, skill.Improvable)
	assert.Equal(t, "midgard", skill.GameSystem)
	//}
	skill2 := models.WeaponSkill{}
	erro = skill2.First("Armbrüste")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(skill.ID))

	skill3 := models.WeaponSkill{}
	erro = skill3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Armbrüste", skill3.Name)

	assert.Equal(t, skill2.ID, skill3.ID)
	assert.Equal(t, skill2.Name, skill3.Name)
	assert.Equal(t, skill2.Beschreibung, skill3.Beschreibung)
	assert.Equal(t, skill2.Initialwert, skill3.Initialwert)
	assert.Equal(t, skill2.Bonuseigenschaft, skill3.Bonuseigenschaft)
	assert.Equal(t, skill2.Quelle, skill3.Quelle)
	assert.Equal(t, skill2.Improvable, skill3.Improvable)
	assert.Equal(t, skill2.GameSystem, skill3.GameSystem)

	err = CheckWaffenFertigkeiten2GSMaster(character.Waffenfertigkeiten)
	assert.NoError(t, err, "Expected no error when checkimg WeaponSkills against gsmaster")
}

func TestImportSpell2GSMaster(t *testing.T) {
	database.SetupTestDB()
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	skill, erro := TransformImportSpell2GSDMaster(&character.Zauber[0])

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(skill.ID), 1)
	assert.Equal(t, "midgard", skill.GameSystem)
	assert.Equal(t, "Angst", skill.Name)
	assert.Equal(t, "", skill.Beschreibung)
	assert.Equal(t, "ARK5 63", skill.Quelle)
	assert.Equal(t, 0, skill.Stufe)
	assert.Equal(t, 0, skill.AP)
	assert.Equal(t, "", skill.Wirkungsziel)
	assert.Equal(t, 0, skill.Reichweite)
	//}
	skill2 := models.Spell{}
	erro = skill2.First("Angst")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(skill.ID))

	skill3 := models.Spell{}
	erro = skill3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Angst", skill3.Name)

	assert.Equal(t, skill2.ID, skill3.ID)
	assert.Equal(t, skill2.GameSystem, skill3.GameSystem)
	assert.Equal(t, skill2.Name, skill3.Name)
	assert.Equal(t, skill2.Beschreibung, skill3.Beschreibung)
	assert.Equal(t, skill2.Quelle, skill3.Quelle)
	assert.Equal(t, skill2.Stufe, skill3.Stufe)
	assert.Equal(t, skill2.AP, skill3.AP)
	assert.Equal(t, skill2.Reichweite, skill3.Reichweite)
	assert.Equal(t, skill2.Wirkungsziel, skill3.Wirkungsziel)

	err = CheckSpells2GSMaster(character.Zauber)
	assert.NoError(t, err, "Expected no error when checkimg Spells against gsmaster")
}

func TestImportWeapon2GSMaster(t *testing.T) {
	database.SetupTestDB()
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	skill, erro := TransformImportWeapon2GSDMaster(&character.Waffen[0])

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(skill.ID), 1)
	assert.Equal(t, "midgard", skill.GameSystem)
	assert.Equal(t, "Armbrust:schwer", skill.Name)
	assert.Equal(t, "", skill.Beschreibung)
	assert.Equal(t, "", skill.Quelle)
	assert.Equal(t, 5.0, skill.Gewicht)
	assert.Equal(t, 40.0, skill.Wert)
	assert.Equal(t, "check", skill.SkillRequired)
	assert.Equal(t, "check", skill.Damage)
	//}
	skill2 := models.Weapon{}
	erro = skill2.First("Armbrust:schwer")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(skill.ID))

	skill3 := models.Weapon{}
	erro = skill3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Armbrust:schwer", skill3.Name)

	assert.Equal(t, skill2.ID, skill3.ID)
	assert.Equal(t, skill2.GameSystem, skill3.GameSystem)
	assert.Equal(t, skill2.Name, skill3.Name)
	assert.Equal(t, skill2.Beschreibung, skill3.Beschreibung)
	assert.Equal(t, skill2.Quelle, skill3.Quelle)
	assert.Equal(t, skill2.Gewicht, skill3.Gewicht)
	assert.Equal(t, skill2.Wert, skill3.Wert)
	assert.Equal(t, skill2.SkillRequired, skill3.SkillRequired)
	assert.Equal(t, skill2.Damage, skill3.Damage)

	err = CheckWeapons2GSMaster(character.Waffen)
	assert.NoError(t, err, "Expected no error when checkimg Weapons against gsmaster")
}

func TestImportContainer2GSMaster(t *testing.T) {
	database.SetupTestDB()
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	container, erro := TransformImportContainer2GSDMaster(&character.Behaeltnisse[0])

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(container.ID), 1)
	assert.Equal(t, "midgard", container.GameSystem)
	assert.Equal(t, "Lederrucksack", container.Name)
	assert.Equal(t, "für 25 kg", container.Beschreibung)
	assert.Equal(t, "", container.Quelle)
	assert.Equal(t, 0.5, container.Gewicht)
	assert.Equal(t, 4.0, container.Wert)
	assert.Equal(t, 25.0, container.Volumen)
	assert.Equal(t, 25.0, container.Tragkraft)
	//}
	container2 := models.Container{}
	erro = container2.First("Lederrucksack")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(container.ID))

	container3 := models.Container{}
	erro = container3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Lederrucksack", container3.Name)

	assert.Equal(t, container2.ID, container3.ID)
	assert.Equal(t, container2.GameSystem, container3.GameSystem)
	assert.Equal(t, container2.Name, container3.Name)
	assert.Equal(t, container2.Beschreibung, container3.Beschreibung)
	assert.Equal(t, container2.Quelle, container3.Quelle)
	assert.Equal(t, container2.Gewicht, container3.Gewicht)
	assert.Equal(t, container2.Wert, container3.Wert)
	assert.Equal(t, container2.Volumen, container3.Volumen)
	assert.Equal(t, container2.Tragkraft, container3.Tragkraft)

	err = CheckContainers2GSMaster(character.Behaeltnisse)
	assert.NoError(t, err, "Expected no error when checkimg Containers against gsmaster")
}

func TestImportTransportation2GSMaster(t *testing.T) {
	database.SetupTestDB()
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	skill, erro := TransformImportTransportation2GSDMaster(&character.Transportmittel[0])

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(skill.ID), 1)
	assert.Equal(t, "midgard", skill.GameSystem)
	assert.Equal(t, "Karren", skill.Name)
	assert.Equal(t, "für 250 kg", skill.Beschreibung)
	assert.Equal(t, "", skill.Quelle)
	assert.Equal(t, 40.0, skill.Gewicht)
	assert.Equal(t, 14.0, skill.Wert)
	assert.Equal(t, 0.0, skill.Volumen)
	assert.Equal(t, 250.0, skill.Tragkraft)
	//}
	skill2 := models.Transportation{}
	erro = skill2.First("Karren")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(skill.ID))

	skill3 := models.Transportation{}
	erro = skill3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Karren", skill3.Name)

	assert.Equal(t, skill2.ID, skill3.ID)
	assert.Equal(t, skill2.GameSystem, skill3.GameSystem)
	assert.Equal(t, skill2.Name, skill3.Name)
	assert.Equal(t, skill2.Beschreibung, skill3.Beschreibung)
	assert.Equal(t, skill2.Quelle, skill3.Quelle)
	assert.Equal(t, skill2.Gewicht, skill3.Gewicht)
	assert.Equal(t, skill2.Wert, skill3.Wert)
	assert.Equal(t, skill2.Volumen, skill3.Volumen)
	assert.Equal(t, skill2.Tragkraft, skill3.Tragkraft)

	err = CheckTransportations2GSMaster(character.Transportmittel)
	assert.NoError(t, err, "Expected no error when checkimg Transüportations against gsmaster")
}

func TestImportEquipment2GSMaster(t *testing.T) {
	database.SetupTestDB()
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	skill, erro := TransformImportEquipment2GSDMaster(&character.Ausruestung[0])

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(skill.ID), 1)
	assert.Equal(t, "midgard", skill.GameSystem)
	assert.Equal(t, "Lederrüstung", skill.Name)
	assert.Equal(t, "", skill.Beschreibung)
	assert.Equal(t, "", skill.Quelle)
	assert.Equal(t, 13.0, skill.Gewicht)
	assert.Equal(t, 30.0, skill.Wert)
	//}
	skill2 := models.Equipment{}
	erro = skill2.First("Lederrüstung")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(skill.ID))

	skill3 := models.Equipment{}
	erro = skill3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Lederrüstung", skill3.Name)

	assert.Equal(t, skill2.ID, skill3.ID)
	assert.Equal(t, skill2.GameSystem, skill3.GameSystem)
	assert.Equal(t, skill2.Name, skill3.Name)
	assert.Equal(t, skill2.Beschreibung, skill3.Beschreibung)
	assert.Equal(t, skill2.Quelle, skill3.Quelle)
	assert.Equal(t, skill2.Gewicht, skill3.Gewicht)
	assert.Equal(t, skill2.Wert, skill3.Wert)

	err = CheckEquipments2GSMaster(character.Ausruestung)
	assert.NoError(t, err, "Expected no error when checkimg Transüportations against gsmaster")
}
func TestImportBelieve2GSMaster(t *testing.T) {
	database.SetupTestDB()
	fileName := fmt.Sprintf("../testdata/%s", "VTT_Import1.json")
	character, err := readImportChar(fileName)
	assert.NoError(t, err, "Expected no error when Unmarshal filecontent")
	//for i := range character.Fertigkeiten {
	skill, erro := TransformImportBelieve2GSDMaster(character.Glaube)

	assert.NoError(t, erro, "Expected no error when Unmarshal filecontent")
	assert.GreaterOrEqual(t, int(skill.ID), 1)
	assert.Equal(t, "midgard", skill.GameSystem)
	assert.Equal(t, "Torkin", skill.Name)
	assert.Equal(t, "", skill.Beschreibung)
	assert.Equal(t, "", skill.Quelle)
	//}
	skill2 := gsmaster.Believe{}
	erro = skill2.First("Torkin")
	assert.NoError(t, erro, "Expected no error when finding Record by name")
	assert.Equal(t, 1, int(skill.ID))

	skill3 := gsmaster.Believe{}
	erro = skill3.FirstId(1)
	assert.NoError(t, erro, "Expected no error when finding Record by ID")
	assert.Equal(t, "Torkin", skill3.Name)

	assert.Equal(t, skill2.ID, skill3.ID)
	assert.Equal(t, skill2.GameSystem, skill3.GameSystem)
	assert.Equal(t, skill2.Name, skill3.Name)
	assert.Equal(t, skill2.Beschreibung, skill3.Beschreibung)
	assert.Equal(t, skill2.Quelle, skill3.Quelle)
	err = CheckBelieve2GSMaster(character)
	assert.NoError(t, err, "Expected no error when checkimg Transüportations against gsmaster")
}
