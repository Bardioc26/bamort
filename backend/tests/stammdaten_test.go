package tests

import (
	"bamort/database"
	"bamort/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initTestDB4Lookup() *gorm.DB {
	db := SetupTestDB()
	db.AutoMigrate(
		&models.LookupSkill{},          //needed for stammdaten.CheckFertigkeit
		&models.LookupSpell{},          //needed for stammdaten.CheckZauber
		&models.LookupWaeponSkill{},    //needed for stammdaten.CheckWaffenFertigkeit
		&models.LookupEquipment{},      //needed for stammdaten.Check...
		&models.LookupContainer{},      //needed for stammdaten.Check...
		&models.LookupTransportation{}, //needed for stammdaten.Check...
	)
	return db
}

func TestCreateLookupSkill(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Lookup()
	database.DB = testDB // Assign test DB to global DB
	stamm := models.LookupSkill{}
	stamm.System = "Midgard"
	stamm.Name = "Lesen"
	stamm.Beschreibung = "Lesen und Schreiben"
	stamm.Quelle = "kod-4711"
	stamm.Initialkeitswert = 8
	stamm.Bonuseigenschaft = "In"
	err := stamm.Create()
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
}

func TestFindLookupSkill(t *testing.T) {
	// Setup test database
	TestCreateLookupSkill(t)
	stamm := models.LookupSkill{}
	stamm.Name = "Lesen"

	err := stamm.First("Lesen")
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, "In", stamm.Bonuseigenschaft)
}

func TestCreateLookupWaeponSkill(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Lookup()
	database.DB = testDB // Assign test DB to global DB
	stamm := models.LookupWaeponSkill{}
	stamm.System = "Midgard"
	stamm.Name = "Stichwaffen"
	stamm.Beschreibung = "Für Dolche und Ochsenzungen"
	stamm.Quelle = "kod-4712"
	stamm.Initialkeitswert = 5
	stamm.Bonuseigenschaft = "Gs"
	err := stamm.Create()
	assert.NoError(t, err, "expexted Error does not exist in WaffenFertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
}

func TestFindLookupWaeponSkill(t *testing.T) {
	// Setup test database
	TestCreateLookupWaeponSkill(t)
	stamm := models.LookupWaeponSkill{}
	stamm.Name = "Lesen"

	err := stamm.First("Stichwaffen")
	assert.NoError(t, err, "expexted Error does not exist in WaffenFertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, "Stichwaffen", stamm.Name)
	assert.Equal(t, "Für Dolche und Ochsenzungen", stamm.Beschreibung)
	assert.Equal(t, "Gs", stamm.Bonuseigenschaft)
}

func TestCreateLookupSpell(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Lookup()
	database.DB = testDB // Assign test DB to global DB
	stamm := models.LookupSpell{}
	stamm.System = "Midgard"
	stamm.Name = "Unsichtbarkeit"
	stamm.Beschreibung = "werde unsichtbar"
	stamm.Quelle = "Ark-4711"
	stamm.Stufe = 1
	stamm.AP = 1
	stamm.Reichweite = 0
	stamm.Wirkungsziel = "Zauberer"
	stamm.Bonus = 0
	err := stamm.Create()
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, "Unsichtbarkeit", stamm.Name)
	assert.Equal(t, "werde unsichtbar", stamm.Beschreibung)
	assert.Equal(t, 1, stamm.Stufe)
}

func TestFindLookupSpell(t *testing.T) {
	// Setup test database
	TestCreateLookupSpell(t)
	stamm := models.LookupSpell{}
	stamm.Name = "lesen"

	err := stamm.First("Unsichtbarkeit")
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, "Unsichtbarkeit", stamm.Name)
	assert.Equal(t, "werde unsichtbar", stamm.Beschreibung)
	assert.Equal(t, 1, stamm.AP)
	assert.Equal(t, 0, stamm.Reichweite)
	assert.Equal(t, "Zauberer", stamm.Wirkungsziel)
	assert.Equal(t, 0, stamm.Bonus)
}

func TestCreateLookupEquipment(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Lookup()
	database.DB = testDB // Assign test DB to global DB
	stamm := models.LookupEquipment{}
	stamm.System = "Midgard"
	stamm.Name = "Decke"
	stamm.Beschreibung = "zum zudecken"
	stamm.Quelle = "kod-4713"
	stamm.Gewicht = 0.2
	stamm.Wert = 300
	err := stamm.Create()
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
}

func TestFindLookupEquipment(t *testing.T) {
	// Setup test database
	TestCreateLookupEquipment(t)
	stamm := models.LookupEquipment{}
	stamm.Name = "Lesen"

	err := stamm.First("Decke")
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, 0.2, stamm.Gewicht)
}

func TestCreateLookupContainer(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Lookup()
	database.DB = testDB // Assign test DB to global DB
	stamm := models.LookupContainer{}
	stamm.System = "Midgard"
	stamm.Name = "Topf"
	stamm.Beschreibung = "zum kochen"
	stamm.Quelle = "kod-4714"
	stamm.Gewicht = 0.6
	stamm.Wert = 300
	stamm.Volumen = 12.2
	err := stamm.Create()
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
}

func TestFindLookupContainer(t *testing.T) {
	// Setup test database
	TestCreateLookupContainer(t)
	stamm := models.LookupContainer{}
	stamm.Name = "Lesen"

	err := stamm.First("Topf")
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, 0.6, stamm.Gewicht)
	assert.Equal(t, 12.2, stamm.Volumen)
}

func TestCreateLookupTransportation(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Lookup()
	database.DB = testDB // Assign test DB to global DB
	stamm := models.LookupTransportation{}
	stamm.System = "Midgard"
	stamm.Name = "Topf"
	stamm.Beschreibung = "zum kochen"
	stamm.Quelle = "kod-4714"
	stamm.Gewicht = 0.6
	stamm.Wert = 300
	stamm.Volumen = 12.5
	err := stamm.Create()
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
}

func TestFindLookupTransportation(t *testing.T) {
	// Setup test database
	TestCreateLookupTransportation(t)
	stamm := models.LookupTransportation{}
	stamm.Name = "Lesen"

	err := stamm.First("Topf")
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, 0.6, stamm.Gewicht)
	assert.Equal(t, 12.5, stamm.Volumen)
}
