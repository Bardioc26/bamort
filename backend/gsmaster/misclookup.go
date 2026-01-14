package gsmaster

import (
	"bamort/database"
	"bamort/models"
)

// GetMiscLookupByKey retrieves all values for a given key
func GetMiscLookupByKey(key string) ([]models.MiscLookup, error) {
	var items []models.MiscLookup
	err := database.DB.Where("`key` = ?", key).Order("value ASC").Find(&items).Error
	return items, err
}

/*
// PopulateMiscLookupData populates initial misc lookup data if table is empty
func PopulateMiscLookupData() error {
	// Check if data already exists
	var count int64
	if err := database.DB.Model(&models.MiscLookup{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // Data already exists
	}

	// Define initial data
	initialData := []struct {
		key    string
		values []string
	}{
		{
			key:    "gender",
			values: []string{"männlich", "weiblich", "divers"},
		},
		{
			key:    "races",
			values: []string{"Mensch", "Elf", "Zwerg", "Gnom", "Halbling"},
		},
		{
			key: "origins",
			values: []string{
				"Albai", "Aran", "Chryseia", "Clanngadarn", "Erainn",
				"Eschar", "Fuardain", "Ikengabecken", "KanThaiPan", "Küstenstaaten",
				"Moravod", "Nahuatlan", "Rawindra", "Twyneddin", "Valian",
			},
		},
		{
			key:    "social_classes",
			values: []string{"Volk", "Mittelschicht", "Adel"},
		},
		{
			key:    "faiths",
			values: []string{"Keine", "Nathir", "Deis Albai", "Mahal", "Druide"},
		},
		{
			key:    "handedness",
			values: []string{"rechts", "links", "beidhändig"},
		},
	}

	// Insert data
	for _, item := range initialData {
		for _, value := range item.values {
			misc := models.MiscLookup{
				Key:   item.key,
				Value: value,
			}
			if err := database.DB.Create(&misc).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
*/
