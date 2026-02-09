package importer

import (
	"bamort/database"
	"bamort/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setupReconcilerTest initializes the test database and runs migrations
func setupReconcilerTest() {
	database.SetupTestDB()
	// Run migrations to ensure PersonalItem field exists in gsm tables
	models.MigrateStructure(database.DB)
	// Run importer migrations to create ImportHistory and MasterDataImport tables
	MigrateStructure(database.DB)
}

func TestReconcileSkill_ExactMatch(t *testing.T) {
	setupReconcilerTest()

	// Create a master skill
	gs := models.GetGameSystem(0, "midgard")

	masterSkill := &models.Skill{
		Name:         "TestExactMatchSkill",
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
		Beschreibung: "Test skill for exact match",
		Initialwert:  12,
		SourceID:     1,
	}
	err := database.DB.Create(masterSkill).Error
	assert.NoError(t, err)

	// Import a skill with same name
	importSkill := Fertigkeit{
		ImportBase:      ImportBase{Name: "TestExactMatchSkill"},
		Fertigkeitswert: 15,
	}

	// Reconcile
	result, matchType, err := ReconcileSkill(importSkill, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "exact", matchType)
	assert.NotNil(t, result)
	assert.Equal(t, masterSkill.ID, result.ID)
	assert.Equal(t, "TestExactMatchSkill", result.Name)
}

func TestReconcileSkill_CreatePersonal(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Import a skill that doesn't exist
	importSkill := Fertigkeit{
		ImportBase:      ImportBase{Name: "Unbekannte Fertigkeit"},
		Beschreibung:    "Eine neue Fertigkeit",
		Fertigkeitswert: 10,
	}

	// Reconcile
	result, matchType, err := ReconcileSkill(importSkill, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "created_personal", matchType)
	assert.NotNil(t, result)
	assert.Equal(t, "Unbekannte Fertigkeit", result.Name)
	assert.True(t, result.PersonalItem)
	assert.Equal(t, gs.ID, result.GameSystemId)

	// Verify it was created in database
	var dbSkill models.Skill
	err = database.DB.Where("name = ? AND game_system = ?", "Unbekannte Fertigkeit", gs.Name).First(&dbSkill).Error
	assert.NoError(t, err)
	assert.True(t, dbSkill.PersonalItem)
}

func TestReconcileSkill_LogsToMasterDataImport(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")
	importHistoryID := uint(123)

	// Import a skill
	importSkill := Fertigkeit{
		ImportBase:   ImportBase{Name: "Test Fertigkeit"},
		Beschreibung: "Test",
	}

	// Reconcile
	result, matchType, err := ReconcileSkillWithHistory(importSkill, importHistoryID, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "created_personal", matchType)

	// Verify log entry was created
	var logEntry MasterDataImport
	err = database.DB.Where("import_history_id = ? AND item_type = ? AND item_id = ?",
		importHistoryID, "skill", result.ID).First(&logEntry).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test Fertigkeit", logEntry.ExternalName)
	assert.Equal(t, "created_personal", logEntry.MatchType)
}

func TestReconcileWeaponSkill_ExactMatch(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Create a master weapon skill
	masterWS := &models.WeaponSkill{
		Skill: models.Skill{
			Name:         "Langschwert",
			GameSystem:   gs.Name,
			GameSystemId: gs.ID,
			SourceID:     1,
		},
	}

	err := database.DB.Create(masterWS).Error
	assert.NoError(t, err)

	// Import weapon skill
	importWS := Waffenfertigkeit{
		ImportBase:      ImportBase{Name: "Langschwert"},
		Fertigkeitswert: 10,
	}

	// Reconcile
	result, matchType, err := ReconcileWeaponSkill(importWS, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "exact", matchType)
	assert.Equal(t, masterWS.ID, result.ID)
}

func TestReconcileWeaponSkill_CreatePersonal(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Import unknown weapon skill
	importWS := Waffenfertigkeit{
		ImportBase:      ImportBase{Name: "Magisches Schwert"},
		Beschreibung:    "Eine besondere Waffe",
		Fertigkeitswert: 8,
	}

	// Reconcile
	result, matchType, err := ReconcileWeaponSkill(importWS, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "created_personal", matchType)
	assert.True(t, result.PersonalItem)
	assert.Equal(t, "Magisches Schwert", result.Skill.Name)
}

func TestReconcileSpell_ExactMatch(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Create master spell
	masterSpell := &models.Spell{
		Name:         "Feuerball",
		GameSystem:   gs.Name,
		GameSystemId: gs.ID,
		SourceID:     1,
	}
	err := database.DB.Create(masterSpell).Error
	assert.NoError(t, err)

	// Import spell
	importSpell := Zauber{
		ImportBase: ImportBase{Name: "Feuerball"},
	}

	// Reconcile
	result, matchType, err := ReconcileSpell(importSpell, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "exact", matchType)
	assert.Equal(t, masterSpell.ID, result.ID)
}

func TestReconcileSpell_CreatePersonal(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Import unknown spell
	importSpell := Zauber{
		ImportBase:   ImportBase{Name: "Drachenruf"},
		Beschreibung: "Ruft einen Drachen",
	}

	// Reconcile
	result, matchType, err := ReconcileSpell(importSpell, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "created_personal", matchType)
	assert.True(t, result.PersonalItem)
}

func TestReconcileWeapon_ExactMatch(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Create master weapon
	masterWeapon := &models.Weapon{
		Equipment: models.Equipment{
			Name:         "TestExactMatchWeapon",
			GameSystem:   gs.Name,
			GameSystemId: gs.ID,
			SourceID:     1,
		},
	}
	err := database.DB.Create(masterWeapon).Error
	assert.NoError(t, err)

	// Import weapon
	importWeapon := Waffe{
		ImportBase: ImportBase{Name: "TestExactMatchWeapon"},
		Schb:       10,
	}

	// Reconcile
	result, matchType, err := ReconcileWeapon(importWeapon, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "exact", matchType)
	assert.Equal(t, masterWeapon.ID, result.ID)
}

func TestReconcileWeapon_CreatePersonal(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Import unknown weapon
	importWeapon := Waffe{
		ImportBase:   ImportBase{Name: "Excalibur"},
		Beschreibung: "Das legendäre Schwert",
		Schb:         20,
	}

	// Reconcile
	result, matchType, err := ReconcileWeapon(importWeapon, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "created_personal", matchType)
	assert.True(t, result.PersonalItem)
}

func TestReconcileEquipment_CreatePersonal(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Import equipment
	importEquip := Ausruestung{
		ImportBase:   ImportBase{Name: "Magischer Umhang"},
		Beschreibung: "Ein besonderer Umhang",
		Gewicht:      0.5,
	}

	// Reconcile
	result, matchType, err := ReconcileEquipment(importEquip, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "created_personal", matchType)
	assert.True(t, result.PersonalItem)
	assert.Equal(t, "Magischer Umhang", result.Name)
}

func TestReconcileContainer_CreatePersonal(t *testing.T) {
	setupReconcilerTest()

	gs := models.GetGameSystem(0, "midgard")

	// Import container
	importContainer := Behaeltniss{
		ImportBase:   ImportBase{Name: "Magische Tasche"},
		Beschreibung: "Eine verzauberte Tasche",
		Tragkraft:    50.0,
	}

	// Reconcile
	result, matchType, err := ReconcileContainer(importContainer, 1, gs.Name)
	assert.NoError(t, err)
	assert.Equal(t, "created_personal", matchType)
	assert.True(t, result.PersonalItem)
}
