package tests

import (
	"bamort/database"
	"bamort/gsmaster"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func initTestDB4Lookup() *gorm.DB {
	db := SetupTestDB()
	db.AutoMigrate(
		&gsmaster.Skill{},          //needed for stammdaten.CheckFertigkeit
		&gsmaster.Spell{},          //needed for stammdaten.CheckZauber
		&gsmaster.WeaponSkill{},    //needed for stammdaten.CheckWaffenFertigkeit
		&gsmaster.Equipment{},      //needed for stammdaten.Check...
		&gsmaster.Container{},      //needed for stammdaten.Check...
		&gsmaster.Transportation{}, //needed for stammdaten.Check...
		&gsmaster.Weapon{},
		&gsmaster.Believe{},
	)
	return db
}

// GenerateFilename generates a filename based on the prefix and the current date/time
func generateFilename(prefix string, extension string) string {
	// Get the current date and time
	now := time.Now()

	// Format the date and time as "YYYY-MM-DD_HH-MM-SS"
	//timestamp = now.Format("2006-01-02_15-04-05")
	timestamp := now.Format("20060102_150405")

	// Combine the prefix and the timestamp to form the filename
	return fmt.Sprintf("%s_%s.%s", prefix, timestamp, extension)
}

/*s
func TestCreateLookupSkill(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Lookup()
	database.DB = testDB // Assign test DB to global DB
	stamm := gsmaster.Skill{}
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
	stamm := gsmaster.Skill{}
	stamm.Name = "Lesen"

	err := stamm.First("Lesen")
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, "In", stamm.Bonuseigenschaft)
}

func TestCreateLookupWeaponSkill(t *testing.T) {
	// Setup test database
	testDB := initTestDB4Lookup()
	database.DB = testDB // Assign test DB to global DB
	stamm := gsmaster.WeaponSkill{}
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

func TestFindLookupWeaponSkill(t *testing.T) {
	// Setup test database
	TestCreateLookupWeaponSkill(t)
	stamm := gsmaster.WeaponSkill{}
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
	stamm := gsmaster.Spell{}
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
	stamm := gsmaster.Spell{}
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
	stamm := gsmaster.Equipment{}
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
	stamm := gsmaster.Equipment{}
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
	stamm := gsmaster.Container{}
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
	stamm := gsmaster.Container{}
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
	stamm := gsmaster.Transportation{}
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
	stamm := gsmaster.Transportation{}
	stamm.Name = "Lesen"

	err := stamm.First("Topf")
	assert.NoError(t, err, "expexted Error does not exist in Fertigkeit Stammdaten")
	assert.GreaterOrEqual(t, 1, int(stamm.ID), "exepets an ID to be present")
	assert.Equal(t, "midgard", stamm.System)
	assert.Equal(t, 0.6, stamm.Gewicht)
	assert.Equal(t, 12.5, stamm.Volumen)
}
*/

func TestExportGSMasterdata(t *testing.T) {
	// Setup test database
	if database.DB == nil {
		testDB := initTestDB4Import()
		database.DB = testDB // Assign test DB to global DB
	}
	//testDB := initTestDB4Lookup()
	//database.DB = testDB // Assign test DB to global DB
	TestImportSkill2GSMaster(t)
	TestImportWeaponSkill2GSMaster(t)
	TestImportSpell2GSMaster(t)
	TestImportWeapon2GSMaster(t)
	TestImportContainer2GSMaster(t)
	TestImportTransportation2GSMaster(t)
	TestImportEquipment2GSMaster(t)
	TestImportBelieve2GSMaster(t)
	err := gsmaster.Export(generateFilename("../testdata/gsmaster_", "json"))
	assert.NoError(t, err, "expexted no Error during Export if gsmaster data")
}

func TestImportGSMasterdata(t *testing.T) {
	// Setup test database
	if database.DB == nil {
		testDB := initTestDB4Import()
		database.DB = testDB // Assign test DB to global DB
	}
	//testDB := initTestDB4Lookup()
	//database.DB = testDB // Assign test DB to global DB
	err := gsmaster.Import("../testdata/gsmaster_exported_gsdata.json")
	assert.NoError(t, err, "expexted no Error during Export if gsmaster data")
}
