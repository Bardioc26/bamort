package models

import (
	"bamort/database"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEquipmentTestDB(t *testing.T) {
	database.SetupTestDB()

	// Migrate structures
	err := MigrateStructure()
	require.NoError(t, err, "Failed to migrate database structure")

	// Clean up any existing test data
	cleanupEquipmentTestData(t)
}

func cleanupEquipmentTestData(t *testing.T) {
	// Delete all equipment data to ensure clean state
	// Delete from actual equipment tables that exist
	err := database.DB.Exec("DELETE FROM equi_containers").Error
	require.NoError(t, err, "Failed to clean up equipment containers")

	err = database.DB.Exec("DELETE FROM equi_weapons").Error
	require.NoError(t, err, "Failed to clean up equipment weapons")

	err = database.DB.Exec("DELETE FROM equi_equipments").Error
	require.NoError(t, err, "Failed to clean up equipment items")
}

func createTestAusruestung(name string) *EqAusruestung {
	return &EqAusruestung{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: name,
			},
			CharacterID: 1,
			UserID:      1,
		},
		Magisch: Magisch{
			IstMagisch:  false,
			Abw:         0,
			Ausgebrannt: false,
		},
		Beschreibung:  "Test equipment description",
		Anzahl:        1,
		BeinhaltetIn:  "",
		ContainedIn:   0,
		ContainerType: "",
		Bonus:         0,
		Gewicht:       1.5,
		Wert:          10.0,
	}
}

func createTestWaffe(name string) *EqWaffe {
	return &EqWaffe{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: name,
			},
			CharacterID: 1,
			UserID:      1,
		},
		Magisch: Magisch{
			IstMagisch:  false,
			Abw:         0,
			Ausgebrannt: false,
		},
		Beschreibung:            "Test weapon description",
		Abwb:                    2,
		Anb:                     1,
		Anzahl:                  1,
		BeinhaltetIn:            "",
		ContainedIn:             0,
		ContainerType:           "",
		Gewicht:                 2.0,
		NameFuerSpezialisierung: "Test Weapon Specialization",
		Schb:                    6,
		Wert:                    25.0,
	}
}

func createTestEqContainer(name string) *EqContainer {
	return &EqContainer{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: name,
			},
			CharacterID: 1,
			UserID:      1,
		},
		Magisch: Magisch{
			IstMagisch:  false,
			Abw:         0,
			Ausgebrannt: false,
		},
		Beschreibung:     "Test container description",
		BeinhaltetIn:     "",
		ContainedIn:      0,
		IsTransportation: false,
		Gewicht:          0.5,
		Wert:             5.0,
		Tragkraft:        10.0,
		Volumen:          15.0,
		ExtID:            "", // Empty by default, to be set by individual tests
	}
}

// =============================================================================
// Tests for EqAusruestung struct
// =============================================================================

func TestEqAusruestung_TableName(t *testing.T) {
	equipment := EqAusruestung{}
	expected := "equi_equipments"
	actual := equipment.TableName()
	assert.Equal(t, expected, actual)
}

func TestEqAusruestung_Save_Success(t *testing.T) {
	setupEquipmentTestDB(t)

	testEquipment := createTestAusruestung("Test Equipment")

	err := testEquipment.Save()

	assert.NoError(t, err, "Save should succeed")
	assert.Greater(t, testEquipment.ID, uint(0), "ID should be set after save")
}

func TestEqAusruestung_Save_Update(t *testing.T) {
	setupEquipmentTestDB(t)

	testEquipment := createTestAusruestung("Test Equipment Update")

	// Initial save
	err := testEquipment.Save()
	require.NoError(t, err, "Initial save should succeed")

	originalID := testEquipment.ID

	// Update the equipment
	testEquipment.Beschreibung = "Updated description"
	testEquipment.Wert = 20.0

	err = testEquipment.Save()

	assert.NoError(t, err, "Update save should succeed")
	assert.Equal(t, originalID, testEquipment.ID, "ID should remain the same after update")
	assert.Equal(t, "Updated description", testEquipment.Beschreibung)
	assert.Equal(t, 20.0, testEquipment.Wert)
}

func TestEqAusruestung_LinkContainer_NoContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	testEquipment := createTestAusruestung("Test Equipment No Container")
	testEquipment.BeinhaltetIn = "" // No container

	err := testEquipment.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed with no container")
	assert.Equal(t, uint(0), testEquipment.ContainedIn, "ContainedIn should remain 0")
}

func TestEqAusruestung_LinkContainer_WithExistingContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	// Create and save a container first
	testContainer := createTestEqContainer("Test Container")
	testContainer.ExtID = "equipment-container-001" // Unique ExtID for this test
	err := testContainer.Save()
	require.NoError(t, err, "Container save should succeed")

	// Create equipment that references the container
	testEquipment := createTestAusruestung("Test Equipment With Container")
	testEquipment.BeinhaltetIn = testContainer.ExtID

	err = testEquipment.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed")
	assert.Equal(t, testContainer.ID, testEquipment.ContainedIn, "ContainedIn should be set to container ID")
}

func TestEqAusruestung_LinkContainer_NonExistentContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	testEquipment := createTestAusruestung("Test Equipment Non-Existent Container")
	testEquipment.BeinhaltetIn = "non-existent-container"

	err := testEquipment.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed even with non-existent container")
	assert.Equal(t, uint(0), testEquipment.ContainedIn, "ContainedIn should remain 0")
}

// =============================================================================
// Tests for EqWaffe struct
// =============================================================================

func TestEqWaffe_TableName(t *testing.T) {
	weapon := EqWaffe{}
	expected := "equi_weapons"
	actual := weapon.TableName()
	assert.Equal(t, expected, actual)
}

func TestEqWaffe_Save_Success(t *testing.T) {
	setupEquipmentTestDB(t)

	testWeapon := createTestWaffe("Test Weapon")

	err := testWeapon.Save()

	assert.NoError(t, err, "Save should succeed")
	assert.Greater(t, testWeapon.ID, uint(0), "ID should be set after save")
}

func TestEqWaffe_Save_Update(t *testing.T) {
	setupEquipmentTestDB(t)

	testWeapon := createTestWaffe("Test Weapon Update")

	// Initial save
	err := testWeapon.Save()
	require.NoError(t, err, "Initial save should succeed")

	originalID := testWeapon.ID

	// Update the weapon
	testWeapon.Beschreibung = "Updated weapon description"
	testWeapon.Schb = 8
	testWeapon.Wert = 50.0

	err = testWeapon.Save()

	assert.NoError(t, err, "Update save should succeed")
	assert.Equal(t, originalID, testWeapon.ID, "ID should remain the same after update")
	assert.Equal(t, "Updated weapon description", testWeapon.Beschreibung)
	assert.Equal(t, 8, testWeapon.Schb)
	assert.Equal(t, 50.0, testWeapon.Wert)
}

func TestEqWaffe_LinkContainer_NoContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	testWeapon := createTestWaffe("Test Weapon No Container")
	testWeapon.BeinhaltetIn = "" // No container

	err := testWeapon.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed with no container")
	assert.Equal(t, uint(0), testWeapon.ContainedIn, "ContainedIn should remain 0")
}

func TestEqWaffe_LinkContainer_WithExistingContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	// Create and save a container first
	testContainer := createTestEqContainer("Test Container for Weapon")
	testContainer.ExtID = "weapon-container-001" // Unique ExtID for this test
	err := testContainer.Save()
	require.NoError(t, err, "Container save should succeed")

	// Create weapon that references the container
	testWeapon := createTestWaffe("Test Weapon With Container")
	testWeapon.BeinhaltetIn = testContainer.ExtID

	err = testWeapon.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed")
	assert.Equal(t, testContainer.ID, testWeapon.ContainedIn, "ContainedIn should be set to container ID")
}

func TestEqWaffe_LinkContainer_NonExistentContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	testWeapon := createTestWaffe("Test Weapon Non-Existent Container")
	testWeapon.BeinhaltetIn = "non-existent-weapon-container"

	err := testWeapon.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed even with non-existent container")
	assert.Equal(t, uint(0), testWeapon.ContainedIn, "ContainedIn should remain 0")
}

// =============================================================================
// Tests for EqContainer struct
// =============================================================================

func TestEqContainer_TableName(t *testing.T) {
	container := EqContainer{}
	expected := "equi_containers"
	actual := container.TableName()
	assert.Equal(t, expected, actual)
}

func TestEqContainer_Save_Success(t *testing.T) {
	setupEquipmentTestDB(t)

	testContainer := createTestEqContainer("Test Container")

	err := testContainer.Save()

	assert.NoError(t, err, "Save should succeed")
	assert.Greater(t, testContainer.ID, uint(0), "ID should be set after save")
}

func TestEqContainer_Save_Update(t *testing.T) {
	setupEquipmentTestDB(t)

	testContainer := createTestEqContainer("Test Container Update")

	// Initial save
	err := testContainer.Save()
	require.NoError(t, err, "Initial save should succeed")

	originalID := testContainer.ID

	// Update the container
	testContainer.Beschreibung = "Updated container description"
	testContainer.Tragkraft = 20.0
	testContainer.Volumen = 25.0

	err = testContainer.Save()

	assert.NoError(t, err, "Update save should succeed")
	assert.Equal(t, originalID, testContainer.ID, "ID should remain the same after update")
	assert.Equal(t, "Updated container description", testContainer.Beschreibung)
	assert.Equal(t, 20.0, testContainer.Tragkraft)
	assert.Equal(t, 25.0, testContainer.Volumen)
}

func TestEqContainer_FirstExtId_Success(t *testing.T) {
	setupEquipmentTestDB(t)

	testContainer := createTestEqContainer("Test Container Find")
	testContainer.ExtID = "find-container-001" // Unique ExtID for this test
	err := testContainer.Save()
	require.NoError(t, err, "Container save should succeed")

	// Create a new container instance to test FirstExtId
	findContainer := EqContainer{}
	err = findContainer.FirstExtId(testContainer.ExtID)

	assert.NoError(t, err, "FirstExtId should succeed")
	assert.Equal(t, testContainer.ID, findContainer.ID, "Found container should have same ID")
	assert.Equal(t, testContainer.ExtID, findContainer.ExtID, "Found container should have same ExtID")
	assert.Equal(t, testContainer.Name, findContainer.Name, "Found container should have same name")
}

func TestEqContainer_FirstExtId_NotFound(t *testing.T) {
	setupEquipmentTestDB(t)

	findContainer := EqContainer{}
	err := findContainer.FirstExtId("non-existent-ext-id")

	assert.Error(t, err, "FirstExtId should return error for non-existent ExtID")
}

func TestEqContainer_LinkContainer_NoContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	testContainer := createTestEqContainer("Test Container No Parent")
	testContainer.BeinhaltetIn = "" // No parent container

	err := testContainer.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed with no parent container")
	assert.Equal(t, uint(0), testContainer.ContainedIn, "ContainedIn should remain 0")
}

func TestEqContainer_LinkContainer_WithExistingContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	// Create and save a parent container first
	parentContainer := createTestEqContainer("Parent Container")
	parentContainer.ExtID = "parent-container-001"
	err := parentContainer.Save()
	require.NoError(t, err, "Parent container save should succeed")

	// Create child container that references the parent container
	childContainer := createTestEqContainer("Child Container")
	childContainer.ExtID = "child-container-001"
	childContainer.BeinhaltetIn = parentContainer.ExtID

	err = childContainer.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed")
	assert.Equal(t, parentContainer.ID, childContainer.ContainedIn, "ContainedIn should be set to parent container ID")
}

func TestEqContainer_LinkContainer_NonExistentContainer(t *testing.T) {
	setupEquipmentTestDB(t)

	testContainer := createTestEqContainer("Test Container Non-Existent Parent")
	testContainer.BeinhaltetIn = "non-existent-parent-container"

	err := testContainer.LinkContainer()

	assert.NoError(t, err, "LinkContainer should succeed even with non-existent parent container")
	assert.Equal(t, uint(0), testContainer.ContainedIn, "ContainedIn should remain 0")
}

// =============================================================================
// Integration tests for equipment linking workflow
// =============================================================================

func TestEquipmentLinking_CompleteWorkflow(t *testing.T) {
	setupEquipmentTestDB(t)

	// Create a hierarchy: Backpack -> Sword Sheath -> Sword
	// And: Backpack -> Potion

	// 1. Create and save backpack (top-level container)
	backpack := createTestEqContainer("Backpack")
	backpack.ExtID = "backpack-001"
	err := backpack.Save()
	require.NoError(t, err, "Backpack save should succeed")

	// 2. Create and save sword sheath (container within backpack)
	sheath := createTestEqContainer("Sword Sheath")
	sheath.ExtID = "sheath-001"
	sheath.BeinhaltetIn = backpack.ExtID
	err = sheath.Save()
	require.NoError(t, err, "Sheath save should succeed")

	err = sheath.LinkContainer()
	require.NoError(t, err, "Sheath link should succeed")

	// 3. Create sword (weapon within sheath)
	sword := createTestWaffe("Magic Sword")
	sword.BeinhaltetIn = sheath.ExtID
	err = sword.Save()
	require.NoError(t, err, "Sword save should succeed")

	err = sword.LinkContainer()
	require.NoError(t, err, "Sword link should succeed")

	// 4. Create potion (equipment within backpack)
	potion := createTestAusruestung("Healing Potion")
	potion.BeinhaltetIn = backpack.ExtID
	err = potion.Save()
	require.NoError(t, err, "Potion save should succeed")

	err = potion.LinkContainer()
	require.NoError(t, err, "Potion link should succeed")

	// Verify the hierarchy
	assert.Equal(t, uint(0), backpack.ContainedIn, "Backpack should not be in any container")
	assert.Equal(t, backpack.ID, sheath.ContainedIn, "Sheath should be in backpack")
	assert.Equal(t, sheath.ID, sword.ContainedIn, "Sword should be in sheath")
	assert.Equal(t, backpack.ID, potion.ContainedIn, "Potion should be in backpack")
}

func TestEquipmentLinking_CircularReference(t *testing.T) {
	setupEquipmentTestDB(t)

	// Test that the code handles potential circular references gracefully
	// (though the current implementation doesn't prevent them)

	container1 := createTestEqContainer("Container 1")
	container1.ExtID = "container-1"
	err := container1.Save()
	require.NoError(t, err, "Container 1 save should succeed")

	container2 := createTestEqContainer("Container 2")
	container2.ExtID = "container-2"
	container2.BeinhaltetIn = container1.ExtID
	err = container2.Save()
	require.NoError(t, err, "Container 2 save should succeed")

	err = container2.LinkContainer()
	assert.NoError(t, err, "Container 2 link should succeed")

	// This would create a circular reference if we tried to make container1 contain container2
	// But we won't actually do it as it would be invalid
	assert.Equal(t, container1.ID, container2.ContainedIn, "Container 2 should be in Container 1")
}

// =============================================================================
// Tests for struct field validation and data integrity
// =============================================================================

func TestEquipmentStructs_FieldValidation(t *testing.T) {
	setupEquipmentTestDB(t)

	// Test EqAusruestung with various field values
	equipment := createTestAusruestung("Field Validation Test")
	equipment.Magisch.IstMagisch = true
	equipment.Magisch.Abw = 5
	equipment.Magisch.Ausgebrannt = true
	equipment.Anzahl = 3
	equipment.Bonus = 2

	err := equipment.Save()
	assert.NoError(t, err, "Equipment with magic properties should save successfully")

	// Test EqWaffe with combat stats
	weapon := createTestWaffe("Combat Weapon Test")
	weapon.Abwb = 3
	weapon.Anb = 2
	weapon.Schb = 10
	weapon.NameFuerSpezialisierung = "Special Combat Weapon"

	err = weapon.Save()
	assert.NoError(t, err, "Weapon with combat stats should save successfully")

	// Test EqContainer with capacity limits
	container := createTestEqContainer("Large Container Test")
	container.IsTransportation = true
	container.Tragkraft = 100.0
	container.Volumen = 200.0

	err = container.Save()
	assert.NoError(t, err, "Container with large capacity should save successfully")
}

func TestEquipmentStructs_ZeroValues(t *testing.T) {
	setupEquipmentTestDB(t)

	// Test that equipment can be saved with zero/empty values
	equipment := &EqAusruestung{
		BamortCharTrait: BamortCharTrait{
			BamortBase: BamortBase{
				Name: "Minimal Equipment",
			},
			CharacterID: 1,
			UserID:      1,
		},
		// All other fields will be zero values
	}

	err := equipment.Save()
	assert.NoError(t, err, "Equipment with minimal fields should save successfully")
	assert.Greater(t, equipment.ID, uint(0), "ID should be set after save")
}

// =============================================================================
// Performance and edge case tests
// =============================================================================

func TestEquipment_MultipleItems(t *testing.T) {
	setupEquipmentTestDB(t)

	// Test creating multiple items of each type
	const itemCount = 5

	// Create multiple containers
	containerIDs := make([]uint, itemCount)
	for i := 0; i < itemCount; i++ {
		container := createTestEqContainer(fmt.Sprintf("Container %d", i+1))
		container.ExtID = fmt.Sprintf("container-%03d", i+1)
		err := container.Save()
		require.NoError(t, err, "Container %d save should succeed", i+1)
		containerIDs[i] = container.ID
	}

	// Create multiple equipment items
	for i := 0; i < itemCount; i++ {
		equipment := createTestAusruestung(fmt.Sprintf("Equipment %d", i+1))
		if i > 0 {
			// Link some equipment to containers
			equipment.BeinhaltetIn = fmt.Sprintf("container-%03d", (i%2)+1)
		}
		err := equipment.Save()
		require.NoError(t, err, "Equipment %d save should succeed", i+1)

		if equipment.BeinhaltetIn != "" {
			err = equipment.LinkContainer()
			assert.NoError(t, err, "Equipment %d link should succeed", i+1)
		}
	}

	// Create multiple weapons
	for i := 0; i < itemCount; i++ {
		weapon := createTestWaffe(fmt.Sprintf("Weapon %d", i+1))
		if i > 1 {
			// Link some weapons to containers
			weapon.BeinhaltetIn = fmt.Sprintf("container-%03d", ((i-1)%2)+1)
		}
		err := weapon.Save()
		require.NoError(t, err, "Weapon %d save should succeed", i+1)

		if weapon.BeinhaltetIn != "" {
			err = weapon.LinkContainer()
			assert.NoError(t, err, "Weapon %d link should succeed", i+1)
		}
	}

	// Verify all containers were created
	for i, containerID := range containerIDs {
		assert.Greater(t, containerID, uint(0), "Container %d should have valid ID", i+1)
	}
}

func TestEquipment_EdgeCases(t *testing.T) {
	setupEquipmentTestDB(t)

	// Test equipment with very long strings
	equipment := createTestAusruestung("Equipment with Long Description")
	equipment.Beschreibung = string(make([]byte, 1000)) // Very long description
	for i := range equipment.Beschreibung {
		equipment.Beschreibung = equipment.Beschreibung[:i] + "A" + equipment.Beschreibung[i+1:]
	}

	err := equipment.Save()
	assert.NoError(t, err, "Equipment with long description should save successfully")

	// Test equipment with special characters
	weaponSpecial := createTestWaffe("Weapon with Special Chars äöü€")
	weaponSpecial.Beschreibung = "Special description with éñ symbols!"

	err = weaponSpecial.Save()
	assert.NoError(t, err, "Weapon with special characters should save successfully")

	// Test container with extreme values
	container := createTestEqContainer("Extreme Container")
	container.Gewicht = 999999.99
	container.Wert = 0.01
	container.Tragkraft = 0.0
	container.Volumen = 1000000.0

	err = container.Save()
	assert.NoError(t, err, "Container with extreme values should save successfully")
}
