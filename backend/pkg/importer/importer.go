package importer

import (
	"github.com/Bardioc26/bamort/pkg/database"
	"github.com/Bardioc26/bamort/pkg/models"
)

// ImportCharacter importiert einen externen Character (z.B. JSON-Format) und transformiert ihn
func ImportCharacter(ext models.CharacterImport) (*models.Character, error) {
	// 1) Ggf. prüfen, ob wir ext.ImportID schon kennen
	// 2) Falls nein, neu anlegen, oder updaten
	// 3) Stammdaten (Fertigkeiten, Waffen etc.) checken und anlegen, falls neu

	// Dummy-Implementierung
	char := models.Character{
		Name: ext.Name,
		// ...
	}

	// Save
	if err := database.DB.Create(&char).Error; err != nil {
		return nil, err
	}

	return &char, nil
}

// Additional Helper-Funktionen, z.B. für das Mapping von ext.Waffen → models.Waffe
