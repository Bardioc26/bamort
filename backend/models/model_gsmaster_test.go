package models

import (
	"bamort/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGSMasterTestDB(t *testing.T) {
	database.SetupTestDB()

	// Migrate structures
	err := MigrateStructure()
	require.NoError(t, err, "Failed to migrate database structure")
}

func createTestGSMSkill(name string) *Skill {
	return &Skill{
		GameSystem:       "midgard",
		Name:             name,
		Beschreibung:     "Test skill description",
		Category:         "Körper",
		Difficulty:       "Normal",
		Initialwert:      5,
		Bonuseigenschaft: "Gs",
		Improvable:       true,
		InnateSkill:      false,
	}
}

func createTestWeaponSkill(name string) *WeaponSkill {
	return &WeaponSkill{
		Skill: Skill{
			GameSystem:       "midgard",
			Name:             name,
			Beschreibung:     "Test weapon skill description",
			Category:         "Kampf",
			Difficulty:       "Normal",
			Initialwert:      5,
			Bonuseigenschaft: "Gs",
			Improvable:       true,
			InnateSkill:      false,
		},
	}
}

func createTestSpell(name string) *Spell {
	return &Spell{
		GameSystem:       "midgard",
		Name:             name,
		Beschreibung:     "Test spell description",
		Bonus:            0,
		Stufe:            1,
		AP:               "1",
		Art:              "Gestenzauber",
		Zauberdauer:      "10 sec",
		Reichweite:       "10 m",
		Wirkungsziel:     "Person",
		Wirkungsbereich:  "1 Person",
		Wirkungsdauer:    "10 min",
		Ursprung:         "elementar",
		Category:         "Zerstören",
		LearningCategory: "Spruch",
	}
}

func createTestEquipment(name string) *Equipment {
	return &Equipment{
		GameSystem:   "midgard",
		Name:         name,
		Beschreibung: "Test equipment description",
		Gewicht:      1.5,
		Wert:         10.0,
		PersonalItem: false,
	}
}

func createTestWeapon(name string) *Weapon {
	return &Weapon{
		Equipment: Equipment{
			GameSystem:   "midgard",
			Name:         name,
			Beschreibung: "Test weapon description",
			Gewicht:      2.0,
			Wert:         50.0,
			PersonalItem: false,
		},
		SkillRequired: "Einhandschwerter",
		Damage:        "1W6+2",
	}
}

func createTestContainer(name string) *Container {
	return &Container{
		Equipment: Equipment{
			GameSystem:   "midgard",
			Name:         name,
			Beschreibung: "Test container description",
			Gewicht:      0.5,
			Wert:         5.0,
			PersonalItem: false,
		},
		Tragkraft: 10.0,
		Volumen:   20.0,
	}
}

func createTestBelieve(name string) *Believe {
	return &Believe{
		GameSystem:   "midgard",
		Name:         name,
		Beschreibung: "Test believe description",
		SourceID:     1, // Use active source KOD
	}
}

// =============================================================================
// Tests for BamortBase struct
// =============================================================================

func TestBamortBase_StructFields(t *testing.T) {
	base := BamortBase{
		ID:   1,
		Name: "TestBase",
	}

	assert.Equal(t, uint(1), base.ID)
	assert.Equal(t, "TestBase", base.Name)
}

// =============================================================================
// Tests for BamortCharTrait struct
// =============================================================================

func TestBamortCharTrait_StructFields(t *testing.T) {
	trait := BamortCharTrait{
		BamortBase: BamortBase{
			ID:   1,
			Name: "TestTrait",
		},
		CharacterID: 10,
		UserID:      5,
	}

	assert.Equal(t, uint(1), trait.ID)
	assert.Equal(t, "TestTrait", trait.Name)
	assert.Equal(t, uint(10), trait.CharacterID)
	assert.Equal(t, uint(5), trait.UserID)
}

// =============================================================================
// Tests for Skill struct
// =============================================================================

func TestSkill_TableName(t *testing.T) {
	skill := Skill{}
	expected := "gsm_skills"
	actual := skill.TableName()
	assert.Equal(t, expected, actual)
}

func TestSkill_Create(t *testing.T) {
	setupGSMasterTestDB(t)

	skill := createTestGSMSkill("TestCreateSkill")
	err := skill.Create()

	assert.NoError(t, err)
	assert.NotZero(t, skill.ID)
	assert.Equal(t, "midgard", skill.GameSystem)
}

func TestSkill_First_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test skill
	testSkill := createTestGSMSkill("TestFirstSkill")
	err := testSkill.Create()
	require.NoError(t, err)

	// Test First method
	foundSkill := &Skill{}
	err = foundSkill.First("TestFirstSkill")

	assert.NoError(t, err)
	assert.Equal(t, "TestFirstSkill", foundSkill.Name)
	assert.Equal(t, "midgard", foundSkill.GameSystem)
}

func TestSkill_First_NotFound(t *testing.T) {
	setupGSMasterTestDB(t)

	skill := &Skill{}
	err := skill.First("NonExistentSkill")

	assert.Error(t, err)
}

func TestSkill_FirstId_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test skill
	testSkill := createTestGSMSkill("TestFirstIdSkill")
	err := testSkill.Create()
	require.NoError(t, err)

	// Test FirstId method
	foundSkill := &Skill{}
	err = foundSkill.FirstId(testSkill.ID)

	assert.NoError(t, err)
	assert.Equal(t, testSkill.ID, foundSkill.ID)
	assert.Equal(t, "TestFirstIdSkill", foundSkill.Name)
}

func TestSkill_Save(t *testing.T) {
	setupGSMasterTestDB(t)

	skill := createTestGSMSkill("TestSaveSkill")
	err := skill.Create()
	require.NoError(t, err)

	// Modify and save
	skill.Beschreibung = "Updated description"
	err = skill.Save()

	assert.NoError(t, err)

	// Verify the update
	foundSkill := &Skill{}
	err = foundSkill.FirstId(skill.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", foundSkill.Beschreibung)
}

func TestSkill_Delete(t *testing.T) {
	setupGSMasterTestDB(t)

	skill := createTestGSMSkill("TestDeleteSkill")
	err := skill.Create()
	require.NoError(t, err)

	// Delete the skill
	err = skill.Delete()
	assert.NoError(t, err)

	// Verify it's deleted
	foundSkill := &Skill{}
	err = foundSkill.FirstId(skill.ID)
	assert.Error(t, err)
}

func TestSkill_Select(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create test skills with different categories
	skill1 := createTestGSMSkill("TestSelectSkill1")
	skill1.Category = "Körper"
	err := skill1.Create()
	require.NoError(t, err)

	skill2 := createTestGSMSkill("TestSelectSkill2")
	skill2.Category = "Körper"
	err = skill2.Create()
	require.NoError(t, err)

	skill3 := createTestGSMSkill("TestSelectSkill3")
	skill3.Category = "Wissen"
	err = skill3.Create()
	require.NoError(t, err)

	// Test Select method
	skills, err := skill1.Select("category", "Körper")

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(skills), 2) // At least our two test skills

	// Verify all returned skills have the correct category
	for _, skill := range skills {
		assert.Equal(t, "Körper", skill.Category)
	}
}

func TestSkill_GetSkillCategories(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create test skills with different categories
	skill1 := createTestGSMSkill("TestCategorySkill1")
	skill1.Category = "TestCategory1"
	err := skill1.Create()
	require.NoError(t, err)

	skill2 := createTestGSMSkill("TestCategorySkill2")
	skill2.Category = "TestCategory2"
	err = skill2.Create()
	require.NoError(t, err)

	// Test GetSkillCategories
	categories, err := skill1.GetSkillCategories()

	assert.NoError(t, err)
	assert.Contains(t, categories, "TestCategory1")
	assert.Contains(t, categories, "TestCategory2")
}

func TestSelectSkills_All(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create test skills
	skill1 := createTestGSMSkill("TestSelectAllSkill1")
	err := skill1.Create()
	require.NoError(t, err)

	skill2 := createTestGSMSkill("TestSelectAllSkill2")
	err = skill2.Create()
	require.NoError(t, err)

	// Test SelectSkills without filter
	skills, err := SelectSkills()

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(skills), 2)
}

func TestSelectSkills_WithFilter(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create test skills
	skill1 := createTestGSMSkill("TestFilterSkill1")
	skill1.Category = "FilterTest"
	err := skill1.Create()
	require.NoError(t, err)

	skill2 := createTestGSMSkill("TestFilterSkill2")
	skill2.Category = "FilterTest"
	err = skill2.Create()
	require.NoError(t, err)

	// Test SelectSkills with filter
	skills, err := SelectSkills("category", "FilterTest")

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(skills), 2)

	// Verify all returned skills have the correct category
	for _, skill := range skills {
		assert.Equal(t, "FilterTest", skill.Category)
	}
}

// =============================================================================
// Tests for WeaponSkill struct
// =============================================================================

func TestWeaponSkill_TableName(t *testing.T) {
	weaponSkill := WeaponSkill{}
	expected := "gsm_weaponskills"
	actual := weaponSkill.TableName()
	assert.Equal(t, expected, actual)
}

func TestWeaponSkill_Create(t *testing.T) {
	setupGSMasterTestDB(t)

	weaponSkill := createTestWeaponSkill("TestCreateWeaponSkill")
	err := weaponSkill.Create()

	assert.NoError(t, err)
	assert.NotZero(t, weaponSkill.ID)
	assert.Equal(t, "midgard", weaponSkill.GameSystem)
}

func TestWeaponSkill_First_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test weapon skill
	testWeaponSkill := createTestWeaponSkill("TestFirstWeaponSkill")
	err := testWeaponSkill.Create()
	require.NoError(t, err)

	// Test First method
	foundWeaponSkill := &WeaponSkill{}
	err = foundWeaponSkill.First("TestFirstWeaponSkill")

	assert.NoError(t, err)
	assert.Equal(t, "TestFirstWeaponSkill", foundWeaponSkill.Name)
	assert.Equal(t, "midgard", foundWeaponSkill.GameSystem)
}

func TestWeaponSkill_First_NotFound(t *testing.T) {
	setupGSMasterTestDB(t)

	weaponSkill := &WeaponSkill{}
	err := weaponSkill.First("NonExistentWeaponSkill")

	assert.Error(t, err)
}

func TestWeaponSkill_FirstId_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test weapon skill
	testWeaponSkill := createTestWeaponSkill("TestFirstIdWeaponSkill")
	err := testWeaponSkill.Create()
	require.NoError(t, err)

	// Test FirstId method
	foundWeaponSkill := &WeaponSkill{}
	err = foundWeaponSkill.FirstId(testWeaponSkill.ID)

	assert.NoError(t, err)
	assert.Equal(t, testWeaponSkill.ID, foundWeaponSkill.ID)
	assert.Equal(t, "TestFirstIdWeaponSkill", foundWeaponSkill.Name)
}

func TestWeaponSkill_Save(t *testing.T) {
	setupGSMasterTestDB(t)

	weaponSkill := createTestWeaponSkill("TestSaveWeaponSkill")
	err := weaponSkill.Create()
	require.NoError(t, err)

	// Modify and save
	weaponSkill.Beschreibung = "Updated weapon skill description"
	err = weaponSkill.Save()

	assert.NoError(t, err)

	// Verify the update
	foundWeaponSkill := &WeaponSkill{}
	err = foundWeaponSkill.FirstId(weaponSkill.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated weapon skill description", foundWeaponSkill.Beschreibung)
}

// =============================================================================
// Tests for Spell struct
// =============================================================================

func TestSpell_TableName(t *testing.T) {
	spell := Spell{}
	expected := "gsm_spells"
	actual := spell.TableName()
	assert.Equal(t, expected, actual)
}

func TestSpell_Create(t *testing.T) {
	setupGSMasterTestDB(t)

	spell := createTestSpell("TestCreateSpell")
	err := spell.Create()

	assert.NoError(t, err)
	assert.NotZero(t, spell.ID)
	assert.Equal(t, "midgard", spell.GameSystem)
}

func TestSpell_First_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test spell
	testSpell := createTestSpell("TestFirstSpell")
	err := testSpell.Create()
	require.NoError(t, err)

	// Test First method
	foundSpell := &Spell{}
	err = foundSpell.First("TestFirstSpell")

	assert.NoError(t, err)
	assert.Equal(t, "TestFirstSpell", foundSpell.Name)
	assert.Equal(t, "midgard", foundSpell.GameSystem)
}

func TestSpell_First_NotFound(t *testing.T) {
	setupGSMasterTestDB(t)

	spell := &Spell{}
	err := spell.First("NonExistentSpell")

	assert.Error(t, err)
}

func TestSpell_FirstId_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test spell
	testSpell := createTestSpell("TestFirstIdSpell")
	err := testSpell.Create()
	require.NoError(t, err)

	// Test FirstId method
	foundSpell := &Spell{}
	err = foundSpell.FirstId(testSpell.ID)

	assert.NoError(t, err)
	assert.Equal(t, testSpell.ID, foundSpell.ID)
	assert.Equal(t, "TestFirstIdSpell", foundSpell.Name)
}

func TestSpell_Save(t *testing.T) {
	setupGSMasterTestDB(t)

	spell := createTestSpell("TestSaveSpell")
	err := spell.Create()
	require.NoError(t, err)

	// Modify and save
	spell.Beschreibung = "Updated spell description"
	err = spell.Save()

	assert.NoError(t, err)

	// Verify the update
	foundSpell := &Spell{}
	err = foundSpell.FirstId(spell.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated spell description", foundSpell.Beschreibung)
}

func TestSpell_GetSpellCategories(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create test spells with different categories
	spell1 := createTestSpell("TestSpellCategory1")
	spell1.Category = "TestSpellCat1"
	err := spell1.Create()
	require.NoError(t, err)

	spell2 := createTestSpell("TestSpellCategory2")
	spell2.Category = "TestSpellCat2"
	err = spell2.Create()
	require.NoError(t, err)

	// Test GetSpellCategories
	categories, err := spell1.GetSpellCategories()

	assert.NoError(t, err)
	assert.Contains(t, categories, "TestSpellCat1")
	assert.Contains(t, categories, "TestSpellCat2")
}

func TestSelectSpells_All(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create test spells
	spell1 := createTestSpell("TestSelectAllSpell1")
	err := spell1.Create()
	require.NoError(t, err)

	spell2 := createTestSpell("TestSelectAllSpell2")
	err = spell2.Create()
	require.NoError(t, err)

	// Test SelectSpells without filter
	spells, err := SelectSpells()

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(spells), 2)
}

func TestSelectSpells_WithFilter(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create test spells
	spell1 := createTestSpell("TestFilterSpell1")
	spell1.Category = "FilterSpellTest"
	err := spell1.Create()
	require.NoError(t, err)

	spell2 := createTestSpell("TestFilterSpell2")
	spell2.Category = "FilterSpellTest"
	err = spell2.Create()
	require.NoError(t, err)

	// Test SelectSpells with filter
	spells, err := SelectSpells("category", "FilterSpellTest")

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(spells), 2)

	// Verify all returned spells have the correct category
	for _, spell := range spells {
		assert.Equal(t, "FilterSpellTest", spell.Category)
	}
}

// =============================================================================
// Tests for Equipment struct
// =============================================================================

func TestEquipment_TableName(t *testing.T) {
	equipment := Equipment{}
	expected := "gsm_equipments"
	actual := equipment.TableName()
	assert.Equal(t, expected, actual)
}

func TestEquipment_Create(t *testing.T) {
	setupGSMasterTestDB(t)

	equipment := createTestEquipment("TestCreateEquipment")
	err := equipment.Create()

	assert.NoError(t, err)
	assert.NotZero(t, equipment.ID)
	assert.Equal(t, "midgard", equipment.GameSystem)
}

func TestEquipment_First_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test equipment
	testEquipment := createTestEquipment("TestFirstEquipment")
	err := testEquipment.Create()
	require.NoError(t, err)

	// Test First method
	foundEquipment := &Equipment{}
	err = foundEquipment.First("TestFirstEquipment")

	assert.NoError(t, err)
	assert.Equal(t, "TestFirstEquipment", foundEquipment.Name)
	assert.Equal(t, "midgard", foundEquipment.GameSystem)
}

func TestEquipment_First_NotFound(t *testing.T) {
	setupGSMasterTestDB(t)

	equipment := &Equipment{}
	err := equipment.First("NonExistentEquipment")

	assert.Error(t, err)
}

func TestEquipment_FirstId_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test equipment
	testEquipment := createTestEquipment("TestFirstIdEquipment")
	err := testEquipment.Create()
	require.NoError(t, err)

	// Test FirstId method
	foundEquipment := &Equipment{}
	err = foundEquipment.FirstId(testEquipment.ID)

	assert.NoError(t, err)
	assert.Equal(t, testEquipment.ID, foundEquipment.ID)
	assert.Equal(t, "TestFirstIdEquipment", foundEquipment.Name)
}

func TestEquipment_Save(t *testing.T) {
	setupGSMasterTestDB(t)

	equipment := createTestEquipment("TestSaveEquipment")
	err := equipment.Create()
	require.NoError(t, err)

	// Modify and save
	equipment.Beschreibung = "Updated equipment description"
	err = equipment.Save()

	assert.NoError(t, err)

	// Verify the update
	foundEquipment := &Equipment{}
	err = foundEquipment.FirstId(equipment.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated equipment description", foundEquipment.Beschreibung)
}

// =============================================================================
// Tests for Weapon struct
// =============================================================================

func TestWeapon_TableName(t *testing.T) {
	weapon := Weapon{}
	expected := "gsm_weapons"
	actual := weapon.TableName()
	assert.Equal(t, expected, actual)
}

func TestWeapon_Create(t *testing.T) {
	setupGSMasterTestDB(t)

	weapon := createTestWeapon("TestCreateWeapon")
	err := weapon.Create()

	assert.NoError(t, err)
	assert.NotZero(t, weapon.ID)
	assert.Equal(t, "midgard", weapon.GameSystem)
}

func TestWeapon_First_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test weapon
	testWeapon := createTestWeapon("TestFirstWeapon")
	err := testWeapon.Create()
	require.NoError(t, err)

	// Test First method
	foundWeapon := &Weapon{}
	err = foundWeapon.First("TestFirstWeapon")

	assert.NoError(t, err)
	assert.Equal(t, "TestFirstWeapon", foundWeapon.Name)
	assert.Equal(t, "midgard", foundWeapon.GameSystem)
}

func TestWeapon_First_NotFound(t *testing.T) {
	setupGSMasterTestDB(t)

	weapon := &Weapon{}
	err := weapon.First("NonExistentWeapon")

	assert.Error(t, err)
}

func TestWeapon_FirstId_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test weapon
	testWeapon := createTestWeapon("TestFirstIdWeapon")
	err := testWeapon.Create()
	require.NoError(t, err)

	// Test FirstId method
	foundWeapon := &Weapon{}
	err = foundWeapon.FirstId(testWeapon.ID)

	assert.NoError(t, err)
	assert.Equal(t, testWeapon.ID, foundWeapon.ID)
	assert.Equal(t, "TestFirstIdWeapon", foundWeapon.Name)
}

func TestWeapon_Save(t *testing.T) {
	setupGSMasterTestDB(t)

	weapon := createTestWeapon("TestSaveWeapon")
	err := weapon.Create()
	require.NoError(t, err)

	// Modify and save
	weapon.Damage = "2W6+3"
	err = weapon.Save()

	assert.NoError(t, err)

	// Verify the update
	foundWeapon := &Weapon{}
	err = foundWeapon.FirstId(weapon.ID)
	require.NoError(t, err)
	assert.Equal(t, "2W6+3", foundWeapon.Damage)
}

func TestWeapon_RangedWeaponRanges(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a ranged weapon with ranges
	weapon := &Weapon{
		Equipment: Equipment{
			GameSystem:   "midgard",
			Name:         "TestBogen",
			Beschreibung: "Test ranged weapon",
			Gewicht:      1.5,
			Wert:         100.0,
		},
		SkillRequired: "Bogen",
		Damage:        "1W6",
		RangeNear:     10,
		RangeMiddle:   30,
		RangeFar:      100,
	}

	err := weapon.Create()
	require.NoError(t, err)

	// Verify the weapon was created with ranges
	foundWeapon := &Weapon{}
	err = foundWeapon.FirstId(weapon.ID)
	require.NoError(t, err)
	assert.Equal(t, 10, foundWeapon.RangeNear)
	assert.Equal(t, 30, foundWeapon.RangeMiddle)
	assert.Equal(t, 100, foundWeapon.RangeFar)
}

func TestWeapon_IsRanged(t *testing.T) {
	setupGSMasterTestDB(t)

	// Test ranged weapon (has at least one range > 0)
	rangedWeapon := &Weapon{
		Equipment: Equipment{
			GameSystem: "midgard",
			Name:       "TestArmbrust",
		},
		SkillRequired: "Armbrust",
		Damage:        "2W6",
		RangeNear:     15,
		RangeMiddle:   50,
		RangeFar:      150,
	}
	err := rangedWeapon.Create()
	require.NoError(t, err)

	assert.True(t, rangedWeapon.IsRanged(), "Weapon with ranges should be ranged")

	// Test melee weapon (all ranges are 0)
	meleeWeapon := &Weapon{
		Equipment: Equipment{
			GameSystem: "midgard",
			Name:       "TestSchwert",
		},
		SkillRequired: "Einhandschwerter",
		Damage:        "1W6+2",
		RangeNear:     0,
		RangeMiddle:   0,
		RangeFar:      0,
	}
	err = meleeWeapon.Create()
	require.NoError(t, err)

	assert.False(t, meleeWeapon.IsRanged(), "Weapon with no ranges should not be ranged")
}

// =============================================================================
// Tests for Container struct
// =============================================================================

func TestContainer_TableName(t *testing.T) {
	container := Container{}
	expected := "gsm_containers"
	actual := container.TableName()
	assert.Equal(t, expected, actual)
}

func TestContainer_Create(t *testing.T) {
	setupGSMasterTestDB(t)

	container := createTestContainer("TestCreateContainer")
	err := container.Create()

	assert.NoError(t, err)
	assert.NotZero(t, container.ID)
	assert.Equal(t, "midgard", container.GameSystem)
}

func TestContainer_First_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test container
	testContainer := createTestContainer("TestFirstContainer")
	err := testContainer.Create()
	require.NoError(t, err)

	// Test First method
	foundContainer := &Container{}
	err = foundContainer.First("TestFirstContainer")

	assert.NoError(t, err)
	assert.Equal(t, "TestFirstContainer", foundContainer.Name)
	assert.Equal(t, "midgard", foundContainer.GameSystem)
}

func TestContainer_First_NotFound(t *testing.T) {
	setupGSMasterTestDB(t)

	container := &Container{}
	err := container.First("NonExistentContainer")

	assert.Error(t, err)
}

func TestContainer_FirstId_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test container
	testContainer := createTestContainer("TestFirstIdContainer")
	err := testContainer.Create()
	require.NoError(t, err)

	// Test FirstId method
	foundContainer := &Container{}
	err = foundContainer.FirstId(testContainer.ID)

	assert.NoError(t, err)
	assert.Equal(t, testContainer.ID, foundContainer.ID)
	assert.Equal(t, "TestFirstIdContainer", foundContainer.Name)
}

func TestContainer_Save(t *testing.T) {
	setupGSMasterTestDB(t)

	container := createTestContainer("TestSaveContainer")
	err := container.Create()
	require.NoError(t, err)

	// Modify and save
	container.Tragkraft = 25.0
	err = container.Save()

	assert.NoError(t, err)

	// Verify the update
	foundContainer := &Container{}
	err = foundContainer.FirstId(container.ID)
	require.NoError(t, err)
	assert.Equal(t, 25.0, foundContainer.Tragkraft)
}

// =============================================================================
// Tests for Transportation struct
// =============================================================================

func TestTransportation_TableName(t *testing.T) {
	transportation := Transportation{}
	expected := "gsm_transportations"
	actual := transportation.TableName()
	assert.Equal(t, expected, actual)
}

func TestTransportation_Create(t *testing.T) {
	setupGSMasterTestDB(t)

	transportation := &Transportation{
		Container: Container{
			Equipment: Equipment{
				GameSystem:   "midgard",
				Name:         "TestCreateTransportation",
				Beschreibung: "Test transportation description",
				Gewicht:      100.0,
				Wert:         500.0,
				PersonalItem: false,
			},
			Tragkraft: 200.0,
			Volumen:   500.0,
		},
	}
	err := transportation.Create()

	assert.NoError(t, err)
	assert.NotZero(t, transportation.ID)
	assert.Equal(t, "midgard", transportation.GameSystem)
}

func TestTransportation_First_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test transportation
	testTransportation := &Transportation{
		Container: Container{
			Equipment: Equipment{
				GameSystem:   "midgard",
				Name:         "TestFirstTransportation",
				Beschreibung: "Test transportation description",
				Gewicht:      100.0,
				Wert:         500.0,
				PersonalItem: false,
			},
			Tragkraft: 200.0,
			Volumen:   500.0,
		},
	}
	err := testTransportation.Create()
	require.NoError(t, err)

	// Test First method
	foundTransportation := &Transportation{}
	err = foundTransportation.First("TestFirstTransportation")

	assert.NoError(t, err)
	assert.Equal(t, "TestFirstTransportation", foundTransportation.Name)
	assert.Equal(t, "midgard", foundTransportation.GameSystem)
}

func TestTransportation_FirstId_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test transportation
	testTransportation := &Transportation{
		Container: Container{
			Equipment: Equipment{
				GameSystem:   "midgard",
				Name:         "TestFirstIdTransportation",
				Beschreibung: "Test transportation description",
				Gewicht:      100.0,
				Wert:         500.0,
				PersonalItem: false,
			},
			Tragkraft: 200.0,
			Volumen:   500.0,
		},
	}
	err := testTransportation.Create()
	require.NoError(t, err)

	// Test FirstId method
	foundTransportation := &Transportation{}
	err = foundTransportation.FirstId(testTransportation.ID)

	assert.NoError(t, err)
	assert.Equal(t, testTransportation.ID, foundTransportation.ID)
	assert.Equal(t, "TestFirstIdTransportation", foundTransportation.Name)
}

func TestTransportation_Save(t *testing.T) {
	setupGSMasterTestDB(t)

	transportation := &Transportation{
		Container: Container{
			Equipment: Equipment{
				GameSystem:   "midgard",
				Name:         "TestSaveTransportation",
				Beschreibung: "Test transportation description",
				Gewicht:      100.0,
				Wert:         500.0,
				PersonalItem: false,
			},
			Tragkraft: 200.0,
			Volumen:   500.0,
		},
	}
	err := transportation.Create()
	require.NoError(t, err)

	// Modify and save
	transportation.Tragkraft = 300.0
	err = transportation.Save()

	assert.NoError(t, err)

	// Verify the update
	foundTransportation := &Transportation{}
	err = foundTransportation.FirstId(transportation.ID)
	require.NoError(t, err)
	assert.Equal(t, 300.0, foundTransportation.Tragkraft)
}

// =============================================================================
// Tests for Believe struct
// =============================================================================

func TestBelieve_TableName(t *testing.T) {
	believe := Believe{}
	expected := "gsm_believes"
	actual := believe.TableName()
	assert.Equal(t, expected, actual)
}

func TestBelieve_Create(t *testing.T) {
	setupGSMasterTestDB(t)

	believe := createTestBelieve("TestCreateBelieve")
	err := believe.Create()

	assert.NoError(t, err)
	assert.NotZero(t, believe.ID)
	assert.Equal(t, "midgard", believe.GameSystem)
}

func TestBelieve_First_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test believe
	testBelieve := createTestBelieve("TestFirstBelieve")
	err := testBelieve.Create()
	require.NoError(t, err)

	// Test First method
	foundBelieve := &Believe{}
	err = foundBelieve.First("TestFirstBelieve")

	assert.NoError(t, err)
	assert.Equal(t, "TestFirstBelieve", foundBelieve.Name)
	assert.Equal(t, "midgard", foundBelieve.GameSystem)
}

func TestBelieve_FirstId_Success(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create a test believe
	testBelieve := createTestBelieve("TestFirstIdBelieve")
	err := testBelieve.Create()
	require.NoError(t, err)

	// Test FirstId method
	foundBelieve := &Believe{}
	err = foundBelieve.FirstId(testBelieve.ID)

	assert.NoError(t, err)
	assert.Equal(t, testBelieve.ID, foundBelieve.ID)
	assert.Equal(t, "TestFirstIdBelieve", foundBelieve.Name)
}

func TestBelieve_Save(t *testing.T) {
	setupGSMasterTestDB(t)

	believe := createTestBelieve("TestSaveBelieve")
	err := believe.Create()
	require.NoError(t, err)

	// Modify and save
	believe.Beschreibung = "Updated believe description"
	err = believe.Save()

	assert.NoError(t, err)

	// Verify the update
	foundBelieve := &Believe{}
	err = foundBelieve.FirstId(believe.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated believe description", foundBelieve.Beschreibung)
}

// =============================================================================
// Tests for Magisch struct
// =============================================================================

func TestMagisch_StructFields(t *testing.T) {
	magisch := Magisch{
		IstMagisch:  true,
		Abw:         5,
		Ausgebrannt: false,
	}

	assert.True(t, magisch.IstMagisch)
	assert.Equal(t, 5, magisch.Abw)
	assert.False(t, magisch.Ausgebrannt)
}

// =============================================================================
// Tests for LookupList struct
// =============================================================================

func TestLookupList_StructFields(t *testing.T) {
	lookup := LookupList{
		ID:           1,
		GameSystem:   "midgard",
		Name:         "TestLookup",
		Beschreibung: "Test lookup description",
		Quelle:       "Test source",
		SourceID:     10,
		PageNumber:   42,
	}

	assert.Equal(t, uint(1), lookup.ID)
	assert.Equal(t, "midgard", lookup.GameSystem)
	assert.Equal(t, "TestLookup", lookup.Name)
	assert.Equal(t, "Test lookup description", lookup.Beschreibung)
	assert.Equal(t, "Test source", lookup.Quelle)
	assert.Equal(t, uint(10), lookup.SourceID)
	assert.Equal(t, 42, lookup.PageNumber)
}

// Note: LookupList methods are commented out in the source code, so no functional tests

// =============================================================================
// Additional missing tests for existing structs
// =============================================================================

func TestTransportation_First_NotFound(t *testing.T) {
	setupGSMasterTestDB(t)

	transportation := &Transportation{}
	err := transportation.First("NonExistentTransportation")

	assert.Error(t, err)
}

func TestBelieve_First_NotFound(t *testing.T) {
	setupGSMasterTestDB(t)

	believe := &Believe{}
	err := believe.First("NonExistentBelieve")

	assert.Error(t, err)
}

// =============================================================================
// Tests for Global Functions
// =============================================================================

func TestGetBelievesByActiveSources(t *testing.T) {
	setupGSMasterTestDB(t)

	// Create test believes
	believe1 := createTestBelieve("ActiveSourceBelieve1")
	err := believe1.Create()
	require.NoError(t, err)

	believe2 := createTestBelieve("ActiveSourceBelieve2")
	err = believe2.Create()
	require.NoError(t, err)

	// Test GetBelievesByActiveSources
	believes, err := GetBelievesByActiveSources("midgard")

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(believes), 2)

	// Verify believes are ordered by name
	if len(believes) > 1 {
		for i := 1; i < len(believes); i++ {
			assert.LessOrEqual(t, believes[i-1].Name, believes[i].Name, "Believes should be ordered by name")
		}
	}
}

// =============================================================================
// Additional Edge Cases
// =============================================================================

func TestFirstId_EdgeCases(t *testing.T) {
	setupGSMasterTestDB(t)

	t.Run("WeaponSkill FirstId with non-existent ID", func(t *testing.T) {
		weaponSkill := &WeaponSkill{}
		err := weaponSkill.FirstId(99999)
		assert.Error(t, err)
	})

	t.Run("Spell FirstId with non-existent ID", func(t *testing.T) {
		spell := &Spell{}
		err := spell.FirstId(99999)
		assert.Error(t, err)
	})

	t.Run("Equipment FirstId with non-existent ID", func(t *testing.T) {
		equipment := &Equipment{}
		err := equipment.FirstId(99999)
		assert.Error(t, err)
	})

	t.Run("Weapon FirstId with non-existent ID", func(t *testing.T) {
		weapon := &Weapon{}
		err := weapon.FirstId(99999)
		assert.Error(t, err)
	})

	t.Run("Container FirstId with non-existent ID", func(t *testing.T) {
		container := &Container{}
		err := container.FirstId(99999)
		assert.Error(t, err)
	})

	t.Run("Transportation FirstId with non-existent ID", func(t *testing.T) {
		transportation := &Transportation{}
		err := transportation.FirstId(99999)
		assert.Error(t, err)
	})

	t.Run("Believe FirstId with non-existent ID", func(t *testing.T) {
		believe := &Believe{}
		err := believe.FirstId(99999)
		assert.Error(t, err)
	})
}

func TestFirst_EmptyName_EdgeCases(t *testing.T) {
	setupGSMasterTestDB(t)

	t.Run("WeaponSkill First with empty name", func(t *testing.T) {
		weaponSkill := &WeaponSkill{}
		err := weaponSkill.First("")
		assert.Error(t, err)
	})

	t.Run("Spell First with empty name", func(t *testing.T) {
		spell := &Spell{}
		err := spell.First("")
		assert.Error(t, err)
	})

	t.Run("Equipment First with empty name", func(t *testing.T) {
		equipment := &Equipment{}
		err := equipment.First("")
		assert.Error(t, err)
	})

	t.Run("Weapon First with empty name", func(t *testing.T) {
		weapon := &Weapon{}
		err := weapon.First("")
		assert.Error(t, err)
	})

	t.Run("Container First with empty name", func(t *testing.T) {
		container := &Container{}
		err := container.First("")
		assert.Error(t, err)
	})

	t.Run("Transportation First with empty name", func(t *testing.T) {
		transportation := &Transportation{}
		err := transportation.First("")
		assert.Error(t, err)
	})

	t.Run("Believe First with empty name", func(t *testing.T) {
		believe := &Believe{}
		err := believe.First("")
		assert.Error(t, err)
	})
}

// =============================================================================
// Tests for Struct Field Validation
// =============================================================================

func TestSkill_StructFieldValidation(t *testing.T) {
	skill := Skill{
		ID:               1,
		GameSystem:       "midgard",
		Name:             "TestSkill",
		Beschreibung:     "Test description",
		Quelle:           "Test source",
		SourceID:         10,
		PageNumber:       42,
		Initialwert:      5,
		Bonuseigenschaft: "Gs",
		Improvable:       true,
		InnateSkill:      false,
		Category:         "Körper",
		Difficulty:       "Normal",
	}

	assert.Equal(t, uint(1), skill.ID)
	assert.Equal(t, "midgard", skill.GameSystem)
	assert.Equal(t, "TestSkill", skill.Name)
	assert.Equal(t, "Test description", skill.Beschreibung)
	assert.Equal(t, "Test source", skill.Quelle)
	assert.Equal(t, uint(10), skill.SourceID)
	assert.Equal(t, 42, skill.PageNumber)
	assert.Equal(t, 5, skill.Initialwert)
	assert.Equal(t, "Gs", skill.Bonuseigenschaft)
	assert.True(t, skill.Improvable)
	assert.False(t, skill.InnateSkill)
	assert.Equal(t, "Körper", skill.Category)
	assert.Equal(t, "Normal", skill.Difficulty)
}

func TestWeaponSkill_StructFieldValidation(t *testing.T) {
	weaponSkill := WeaponSkill{
		Skill: Skill{
			ID:               1,
			GameSystem:       "midgard",
			Name:             "TestWeaponSkill",
			Beschreibung:     "Test weapon skill description",
			Category:         "Kampf",
			Difficulty:       "Normal",
			Initialwert:      5,
			Bonuseigenschaft: "Gs",
			Improvable:       true,
			InnateSkill:      false,
		},
	}

	assert.Equal(t, uint(1), weaponSkill.ID)
	assert.Equal(t, "midgard", weaponSkill.GameSystem)
	assert.Equal(t, "TestWeaponSkill", weaponSkill.Name)
	assert.Equal(t, "Test weapon skill description", weaponSkill.Beschreibung)
	assert.Equal(t, "Kampf", weaponSkill.Category)
	assert.Equal(t, "Normal", weaponSkill.Difficulty)
	assert.Equal(t, 5, weaponSkill.Initialwert)
	assert.Equal(t, "Gs", weaponSkill.Bonuseigenschaft)
	assert.True(t, weaponSkill.Improvable)
	assert.False(t, weaponSkill.InnateSkill)
}

func TestSpell_StructFieldValidation(t *testing.T) {
	spell := Spell{
		ID:               1,
		GameSystem:       "midgard",
		Name:             "TestSpell",
		Beschreibung:     "Test spell description",
		Quelle:           "Test source",
		SourceID:         10,
		PageNumber:       42,
		Bonus:            2,
		Stufe:            3,
		AP:               "2",
		Art:              "Gestenzauber",
		Zauberdauer:      "20 sec",
		Reichweite:       "15 m",
		Wirkungsziel:     "Person",
		Wirkungsbereich:  "2 Personen",
		Wirkungsdauer:    "15 min",
		Ursprung:         "elementar",
		Category:         "Zerstören",
		LearningCategory: "Spruch",
	}

	assert.Equal(t, uint(1), spell.ID)
	assert.Equal(t, "midgard", spell.GameSystem)
	assert.Equal(t, "TestSpell", spell.Name)
	assert.Equal(t, "Test spell description", spell.Beschreibung)
	assert.Equal(t, "Test source", spell.Quelle)
	assert.Equal(t, uint(10), spell.SourceID)
	assert.Equal(t, 42, spell.PageNumber)
	assert.Equal(t, 2, spell.Bonus)
	assert.Equal(t, 3, spell.Stufe)
	assert.Equal(t, "2", spell.AP)
	assert.Equal(t, "Gestenzauber", spell.Art)
	assert.Equal(t, "20 sec", spell.Zauberdauer)
	assert.Equal(t, "15 m", spell.Reichweite)
	assert.Equal(t, "Person", spell.Wirkungsziel)
	assert.Equal(t, "2 Personen", spell.Wirkungsbereich)
	assert.Equal(t, "15 min", spell.Wirkungsdauer)
	assert.Equal(t, "elementar", spell.Ursprung)
	assert.Equal(t, "Zerstören", spell.Category)
	assert.Equal(t, "Spruch", spell.LearningCategory)
}

func TestEquipment_StructFieldValidation(t *testing.T) {
	equipment := Equipment{
		ID:           1,
		GameSystem:   "midgard",
		Name:         "TestEquipment",
		Beschreibung: "Test equipment description",
		Quelle:       "Test source",
		SourceID:     10,
		PageNumber:   42,
		Gewicht:      2.5,
		Wert:         15.0,
		PersonalItem: true,
	}

	assert.Equal(t, uint(1), equipment.ID)
	assert.Equal(t, "midgard", equipment.GameSystem)
	assert.Equal(t, "TestEquipment", equipment.Name)
	assert.Equal(t, "Test equipment description", equipment.Beschreibung)
	assert.Equal(t, "Test source", equipment.Quelle)
	assert.Equal(t, uint(10), equipment.SourceID)
	assert.Equal(t, 42, equipment.PageNumber)
	assert.Equal(t, 2.5, equipment.Gewicht)
	assert.Equal(t, 15.0, equipment.Wert)
	assert.True(t, equipment.PersonalItem)
}

func TestWeapon_StructFieldValidation(t *testing.T) {
	weapon := Weapon{
		Equipment: Equipment{
			ID:           1,
			GameSystem:   "midgard",
			Name:         "TestWeapon",
			Beschreibung: "Test weapon description",
			Gewicht:      3.0,
			Wert:         75.0,
			PersonalItem: false,
		},
		SkillRequired: "Einhandschwerter",
		Damage:        "1W8+3",
	}

	assert.Equal(t, uint(1), weapon.ID)
	assert.Equal(t, "midgard", weapon.GameSystem)
	assert.Equal(t, "TestWeapon", weapon.Name)
	assert.Equal(t, "Test weapon description", weapon.Beschreibung)
	assert.Equal(t, 3.0, weapon.Gewicht)
	assert.Equal(t, 75.0, weapon.Wert)
	assert.False(t, weapon.PersonalItem)
	assert.Equal(t, "Einhandschwerter", weapon.SkillRequired)
	assert.Equal(t, "1W8+3", weapon.Damage)
}

func TestContainer_StructFieldValidation(t *testing.T) {
	container := Container{
		Equipment: Equipment{
			ID:           1,
			GameSystem:   "midgard",
			Name:         "TestContainer",
			Beschreibung: "Test container description",
			Gewicht:      1.0,
			Wert:         8.0,
			PersonalItem: false,
		},
		Tragkraft: 15.0,
		Volumen:   30.0,
	}

	assert.Equal(t, uint(1), container.ID)
	assert.Equal(t, "midgard", container.GameSystem)
	assert.Equal(t, "TestContainer", container.Name)
	assert.Equal(t, "Test container description", container.Beschreibung)
	assert.Equal(t, 1.0, container.Gewicht)
	assert.Equal(t, 8.0, container.Wert)
	assert.False(t, container.PersonalItem)
	assert.Equal(t, 15.0, container.Tragkraft)
	assert.Equal(t, 30.0, container.Volumen)
}

func TestTransportation_StructFieldValidation(t *testing.T) {
	transportation := Transportation{
		Container: Container{
			Equipment: Equipment{
				ID:           1,
				GameSystem:   "midgard",
				Name:         "TestTransportation",
				Beschreibung: "Test transportation description",
				Gewicht:      150.0,
				Wert:         750.0,
				PersonalItem: false,
			},
			Tragkraft: 300.0,
			Volumen:   600.0,
		},
	}

	assert.Equal(t, uint(1), transportation.ID)
	assert.Equal(t, "midgard", transportation.GameSystem)
	assert.Equal(t, "TestTransportation", transportation.Name)
	assert.Equal(t, "Test transportation description", transportation.Beschreibung)
	assert.Equal(t, 150.0, transportation.Gewicht)
	assert.Equal(t, 750.0, transportation.Wert)
	assert.False(t, transportation.PersonalItem)
	assert.Equal(t, 300.0, transportation.Tragkraft)
	assert.Equal(t, 600.0, transportation.Volumen)
}

func TestBelieve_StructFieldValidation(t *testing.T) {
	believe := Believe{
		ID:           1,
		GameSystem:   "midgard",
		Name:         "TestBelieve",
		Beschreibung: "Test believe description",
		Quelle:       "Test source",
		SourceID:     10,
		PageNumber:   42,
	}

	assert.Equal(t, uint(1), believe.ID)
	assert.Equal(t, "midgard", believe.GameSystem)
	assert.Equal(t, "TestBelieve", believe.Name)
	assert.Equal(t, "Test believe description", believe.Beschreibung)
	assert.Equal(t, "Test source", believe.Quelle)
	assert.Equal(t, uint(10), believe.SourceID)
	assert.Equal(t, 42, believe.PageNumber)
}

// =============================================================================
// Additional Integration Tests
// =============================================================================

func TestGSMasterStructures_WithDatabase(t *testing.T) {
	setupGSMasterTestDB(t)

	t.Run("Skill Database Integration", func(t *testing.T) {
		skill := createTestGSMSkill("IntegrationTestSkill")
		err := skill.Create()
		require.NoError(t, err)

		// Test retrieval
		foundSkill := &Skill{}
		err = foundSkill.First("IntegrationTestSkill")
		assert.NoError(t, err)
		assert.Equal(t, "IntegrationTestSkill", foundSkill.Name)

		// Test update
		foundSkill.Category = "UpdatedCategory"
		err = foundSkill.Save()
		assert.NoError(t, err)

		// Verify update
		reFoundSkill := &Skill{}
		err = reFoundSkill.FirstId(foundSkill.ID)
		assert.NoError(t, err)
		assert.Equal(t, "UpdatedCategory", reFoundSkill.Category)
	})

	t.Run("Spell Database Integration", func(t *testing.T) {
		spell := createTestSpell("IntegrationTestSpell")
		err := spell.Create()
		require.NoError(t, err)

		// Test retrieval
		foundSpell := &Spell{}
		err = foundSpell.First("IntegrationTestSpell")
		assert.NoError(t, err)
		assert.Equal(t, "IntegrationTestSpell", foundSpell.Name)

		// Test update
		foundSpell.Stufe = 5
		err = foundSpell.Save()
		assert.NoError(t, err)

		// Verify update
		reFoundSpell := &Spell{}
		err = reFoundSpell.FirstId(foundSpell.ID)
		assert.NoError(t, err)
		assert.Equal(t, 5, reFoundSpell.Stufe)
	})

	t.Run("Equipment Database Integration", func(t *testing.T) {
		equipment := createTestEquipment("IntegrationTestEquipment")
		err := equipment.Create()
		require.NoError(t, err)

		// Test retrieval
		foundEquipment := &Equipment{}
		err = foundEquipment.First("IntegrationTestEquipment")
		assert.NoError(t, err)
		assert.Equal(t, "IntegrationTestEquipment", foundEquipment.Name)

		// Test update
		foundEquipment.Wert = 25.0
		err = foundEquipment.Save()
		assert.NoError(t, err)

		// Verify update
		reFoundEquipment := &Equipment{}
		err = reFoundEquipment.FirstId(foundEquipment.ID)
		assert.NoError(t, err)
		assert.Equal(t, 25.0, reFoundEquipment.Wert)
	})
}

func TestTableNames_GSMaster_Consistency(t *testing.T) {
	// Test that all table names follow the expected pattern
	skill := Skill{}
	weaponSkill := WeaponSkill{}
	spell := Spell{}
	equipment := Equipment{}
	weapon := Weapon{}
	container := Container{}
	transportation := Transportation{}
	believe := Believe{}

	assert.Equal(t, "gsm_skills", skill.TableName())
	assert.Equal(t, "gsm_weaponskills", weaponSkill.TableName())
	assert.Equal(t, "gsm_spells", spell.TableName())
	assert.Equal(t, "gsm_equipments", equipment.TableName())
	assert.Equal(t, "gsm_weapons", weapon.TableName())
	assert.Equal(t, "gsm_containers", container.TableName())
	assert.Equal(t, "gsm_transportations", transportation.TableName())
	assert.Equal(t, "gsm_believes", believe.TableName())

	// All table names should start with "gsm_"
	assert.Contains(t, skill.TableName(), "gsm_")
	assert.Contains(t, weaponSkill.TableName(), "gsm_")
	assert.Contains(t, spell.TableName(), "gsm_")
	assert.Contains(t, equipment.TableName(), "gsm_")
	assert.Contains(t, weapon.TableName(), "gsm_")
	assert.Contains(t, container.TableName(), "gsm_")
	assert.Contains(t, transportation.TableName(), "gsm_")
	assert.Contains(t, believe.TableName(), "gsm_")
}

// =============================================================================
// Additional Edge Cases and Error Handling
// =============================================================================

func TestGSMaster_EdgeCases(t *testing.T) {
	setupGSMasterTestDB(t)

	t.Run("First with empty name", func(t *testing.T) {
		skill := &Skill{}
		err := skill.First("")
		assert.Error(t, err)
	})

	t.Run("FirstId with non-existent ID", func(t *testing.T) {
		skill := &Skill{}
		err := skill.FirstId(99999)
		assert.Error(t, err)
	})

	t.Run("Delete non-existent record", func(t *testing.T) {
		skill := &Skill{
			ID: 99999,
		}
		err := skill.Delete()
		assert.Error(t, err)
	})
}

// =============================================================================
// Tests for LearnCost struct
// =============================================================================

func TestLearnCost_StructFields(t *testing.T) {
	learnCost := LearnCost{
		Stufe: 3,
		LE:    10,
		TE:    5,
		Ep:    100,
		Money: 50,
		PP:    2,
	}

	assert.Equal(t, 3, learnCost.Stufe)
	assert.Equal(t, 10, learnCost.LE)
	assert.Equal(t, 5, learnCost.TE)
	assert.Equal(t, 100, learnCost.Ep)
	assert.Equal(t, 50, learnCost.Money)
	assert.Equal(t, 2, learnCost.PP)
}

// =============================================================================
// Additional Benchmark Tests
// =============================================================================

func BenchmarkSkill_Create(b *testing.B) {
	database.SetupTestDB()
	err := MigrateStructure()
	if err != nil {
		b.Fatal("Failed to migrate structure:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skill := createTestGSMSkill("BenchmarkSkill")
		skill.Create()
	}
}

func BenchmarkSkill_First(b *testing.B) {
	database.SetupTestDB()
	err := MigrateStructure()
	if err != nil {
		b.Fatal("Failed to migrate structure:", err)
	}

	// Create test skill
	testSkill := createTestGSMSkill("BenchmarkFirstSkill")
	err = testSkill.Create()
	if err != nil {
		b.Fatal("Failed to create test skill:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skill := &Skill{}
		skill.First("BenchmarkFirstSkill")
	}
}

func BenchmarkSpell_Create(b *testing.B) {
	database.SetupTestDB()
	err := MigrateStructure()
	if err != nil {
		b.Fatal("Failed to migrate structure:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		spell := createTestSpell("BenchmarkSpell")
		spell.Create()
	}
}

func BenchmarkEquipment_Create(b *testing.B) {
	database.SetupTestDB()
	err := MigrateStructure()
	if err != nil {
		b.Fatal("Failed to migrate structure:", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		equipment := createTestEquipment("BenchmarkEquipment")
		equipment.Create()
	}
}
