package gsmaster

import (
	"bamort/database"
	"bamort/models"
	"fmt"
	"strings"
)

// GetMiscLookupByKey retrieves all values for a given key
// Optional order parameter can be: "id", "value", "source", "source_value"
// Default is "value" if not specified or invalid
func GetMiscLookupByKey(key string, order ...string) ([]models.MiscLookup, error) {
	var items []models.MiscLookup

	// Determine ordering
	orderBy := "value ASC"
	if len(order) > 0 && order[0] != "" {
		switch order[0] {
		case "id":
			orderBy = "id ASC"
		case "value":
			orderBy = "value ASC"
		case "source":
			orderBy = "source_id ASC, value ASC"
		case "source_value":
			orderBy = "source_id ASC, value ASC"
		default:
			orderBy = "value ASC"
		}
	}

	err := database.DB.Where("`key` = ?", key).Order(orderBy).Find(&items).Error
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
			values: []string{"Unfrei", "Volk", "Mittelschicht", "Adel"},
		},
		{
			key:    "faiths",
			values: []string{"Keine", "Nathir", "Deis Albai", "Mahal", "Druide"},
		},
		{
			key:    "handedness",
			values: []string{"rechts", "links", "beidhändig"},
		},
		{
			key:    "social_class_bonus",
			values: []string{"Unfrei|Halbwelt|2", "Volk|Alltag|2", "Mittelschicht|Wissen|2", "Adel|Sozial|2"},
		},
	}

	// Insert data
	for _, item := range initialData {
		for _, value := range item.values {
			misc := models.MiscLookup{
				Key:      item.key,
				Value:    value,
				SourceID: 1,
			}
			if err := database.DB.Create(&misc).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
*/

// GetSocialClassBonusPoints retrieves bonus learning points for a social class from database
func GetSocialClassBonusPoints(socialClass string) (map[string]int, error) {
	bonuses, err := GetMiscLookupByKey("social_class_bonus")
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, bonus := range bonuses {
		// Parse format: "SocialClass|Category|Points"
		parts := strings.Split(bonus.Value, "|")
		if len(parts) == 3 && parts[0] == socialClass {
			category := parts[1]
			points := 0
			fmt.Sscanf(parts[2], "%d", &points)
			result[category] = points
		}
	}

	return result, nil
}
