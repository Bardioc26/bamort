package importer

import (
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates a test database for importer model testing
func setupTestDB(t *testing.T) *gorm.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_importer.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Run migrations
	err = MigrateStructure(db)
	if err != nil {
		t.Fatalf("Failed to migrate importer structure: %v", err)
	}

	return db
}

func TestImportHistoryModel(t *testing.T) {
	db := setupTestDB(t)

	// Create a test import history
	userID := uint(1)
	charID := uint(18) // Test character
	history := ImportHistory{
		UserID:         userID,
		CharacterID:    &charID,
		AdapterID:      "foundry-vtt-v1",
		SourceFormat:   "foundry-vtt",
		SourceFilename: "test-character.json",
		BmrtVersion:    "1.0",
		ImportedAt:     time.Now(),
		Status:         "success",
	}

	// Create the record
	result := db.Create(&history)
	if result.Error != nil {
		t.Fatalf("Failed to create import history: %v", result.Error)
	}

	// Verify it was created
	if history.ID == 0 {
		t.Error("Import history ID should be set after creation")
	}

	// Fetch it back
	var fetched ImportHistory
	result = db.First(&fetched, history.ID)
	if result.Error != nil {
		t.Fatalf("Failed to fetch import history: %v", result.Error)
	}

	// Verify fields
	if fetched.AdapterID != "foundry-vtt-v1" {
		t.Errorf("AdapterID = %v, want %v", fetched.AdapterID, "foundry-vtt-v1")
	}
	if fetched.Status != "success" {
		t.Errorf("Status = %v, want %v", fetched.Status, "success")
	}
	if fetched.CharacterID == nil || *fetched.CharacterID != charID {
		t.Errorf("CharacterID = %v, want %v", fetched.CharacterID, charID)
	}
}

func TestMasterDataImportModel(t *testing.T) {
	db := setupTestDB(t)

	// Create a test import history first
	history := ImportHistory{
		UserID:       uint(1),
		AdapterID:    "test-adapter",
		SourceFormat: "test",
		ImportedAt:   time.Now(),
		Status:       "success",
	}
	db.Create(&history)

	// Create a test master data import
	masterData := MasterDataImport{
		ImportHistoryID: history.ID,
		ItemType:        "skill",
		ItemID:          uint(1),
		ExternalName:    "Swordfighting",
		MatchType:       "exact",
	}

	// Create the record
	result := db.Create(&masterData)
	if result.Error != nil {
		t.Fatalf("Failed to create master data import: %v", result.Error)
	}

	// Verify it was created
	if masterData.ID == 0 {
		t.Error("Master data import ID should be set after creation")
	}

	// Fetch it back
	var fetched MasterDataImport
	result = db.First(&fetched, masterData.ID)
	if result.Error != nil {
		t.Fatalf("Failed to fetch master data import: %v", result.Error)
	}

	// Verify fields
	if fetched.ItemType != "skill" {
		t.Errorf("ItemType = %v, want %v", fetched.ItemType, "skill")
	}
	if fetched.MatchType != "exact" {
		t.Errorf("MatchType = %v, want %v", fetched.MatchType, "exact")
	}
	if fetched.ImportHistoryID != history.ID {
		t.Errorf("ImportHistoryID = %v, want %v", fetched.ImportHistoryID, history.ID)
	}
}

func TestMigrateStructure(t *testing.T) {
	db := setupTestDB(t)

	// Verify tables exist by trying to query them
	var count int64

	err := db.Model(&ImportHistory{}).Count(&count).Error
	if err != nil {
		t.Errorf("import_histories table should exist: %v", err)
	}

	err = db.Model(&MasterDataImport{}).Count(&count).Error
	if err != nil {
		t.Errorf("master_data_imports table should exist: %v", err)
	}
}

func TestImportHistory_TableName(t *testing.T) {
	history := ImportHistory{}
	if history.TableName() != "import_histories" {
		t.Errorf("TableName() = %v, want 'import_histories'", history.TableName())
	}
}

func TestMasterDataImport_TableName(t *testing.T) {
	masterData := MasterDataImport{}
	if masterData.TableName() != "master_data_imports" {
		t.Errorf("TableName() = %v, want 'master_data_imports'", masterData.TableName())
	}
}
