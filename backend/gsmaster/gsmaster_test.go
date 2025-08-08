package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestExportGSMasterdata(t *testing.T) {
	/*
		database.SetupTestDB()
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
	*/
}

func TestImportGSMasterdata(t *testing.T) {
	database.SetupTestDB()
	err := Import("../testdata/gsmaster_exported_gsdata.json")
	assert.NoError(t, err, "expexted no Error during Export if gsmaster data")
}

func TestMigrateStructure(t *testing.T) {
	database.SetupTestDB(true) // Use in-memory SQLite for testing
	err := models.MigrateStructure()
	assert.NoError(t, err, "expected no Error during MigrateStructure")
}
