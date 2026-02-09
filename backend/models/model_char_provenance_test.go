package models

import (
	"bamort/database"
	"testing"
	"time"
)

// TestCharProvenanceFields verifies that the provenance fields (ImportedFromAdapter and ImportedAt) are correctly defined in the Char struct and can be set and retrieved without issues.
// This test ensures that the character model can track the source of imported characters, which is crucial for data provenance and debugging import issues.
// logical it belongs to the character model tests since it is testing fields on the Char struct, but it is specifically focused on the provenance tracking fields that were recently added to support the import system. This test ensures that these new fields are properly integrated and functional within the character model.

func TestCharProvenanceFields(t *testing.T) {
	setupTestDB(t)
	db := database.DB

	// Run migrations
	err := characterMigrateStructure(db)
	if err != nil {
		t.Fatalf("Failed to migrate character structure: %v", err)
	}

	// Check that provenance fields exist
	// This is verified by the migration not failing and the struct compiling

	// Create a test character with provenance
	adapterID := "foundry-vtt-v1"
	importedAt := time.Now()
	char := Char{
		BamortBase: BamortBase{
			Name: "Test Import Character",
		},
		GameSystem:          "midgard",
		UserID:              uint(1),
		Typ:                 "Krieger",
		ImportedFromAdapter: &adapterID,
		ImportedAt:          &importedAt,
	}

	// Create the character
	result := db.Create(&char)
	if result.Error != nil {
		t.Fatalf("Failed to create character with provenance: %v", result.Error)
	}

	// Fetch it back
	var fetched Char
	result = db.First(&fetched, char.ID)
	if result.Error != nil {
		t.Fatalf("Failed to fetch character: %v", result.Error)
	}

	// Verify provenance fields
	if fetched.ImportedFromAdapter == nil || *fetched.ImportedFromAdapter != adapterID {
		t.Errorf("ImportedFromAdapter = %v, want %v", fetched.ImportedFromAdapter, adapterID)
	}
	if fetched.ImportedAt == nil {
		t.Error("ImportedAt should not be nil")
	}
}
