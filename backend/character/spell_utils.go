package character

import (
	"bamort/database"
	"bamort/gsmaster"
)

// getSpellCategory ermittelt die Zaubergruppe für einen gegebenen Zaubernamen
// Wenn es sich um einen Zauber handelt, wird die Kategorie zurückgegeben
// Andernfalls wird der ursprüngliche Name zurückgegeben
func getSpellCategory(name string) string {
	var spell gsmaster.Spell
	if err := database.DB.Where("name = ?", name).First(&spell).Error; err != nil {
		// Kein Zauber gefunden, ursprünglichen Namen verwenden
		return name
	}

	// Zauber gefunden, Kategorie direkt zurückgeben
	// Die Kategorien sind bereits die vollen Namen wie "Beherrschen", "Verändern", etc.
	return spell.Category
}
