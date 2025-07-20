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

	// Zauber gefunden, Kategorie zurückgeben
	// Die Kategorien sind typischerweise kurze Formen wie "Beherr", "Veränd", etc.
	// Wir müssen diese zu den vollen Zaubergruppen-Namen mappen
	switch spell.Category {
	case "Beherr":
		return "Beherrschen"
	case "Veränd":
		return "Verändern"
	case "Erken":
		return "Erkennen"
	case "Erschaf":
		return "Erschaffen"
	case "Zerstör":
		return "Zerstören"
	case "Bewegen":
		return "Bewegen"
	default:
		// Fallback: Kategorie direkt verwenden
		return spell.Category
	}
}
