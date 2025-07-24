package character

import (
	"bamort/database"
	"bamort/skills"
)

// determineSkillType ermittelt automatisch den Typ einer Fertigkeit anhand des Namens
func determineSkillType(skillName string) string {
	// Prüfe ob es eine normale Fertigkeit ist
	var fertigkeit skills.Fertigkeit
	if err := database.DB.Where("name = ?", skillName).First(&fertigkeit).Error; err == nil {
		return "fertigkeit"
	}

	// Prüfe ob es eine Waffenfertigkeit ist
	var waffenfertigkeit skills.Waffenfertigkeit
	if err := database.DB.Where("name = ?", skillName).First(&waffenfertigkeit).Error; err == nil {
		return "waffenfertigkeit"
	}

	// Prüfe ob es ein Zauber ist
	var zauber skills.Zauber
	if err := database.DB.Where("name = ?", skillName).First(&zauber).Error; err == nil {
		return "zauber"
	}

	// Fallback: unbekannt
	return ""
}

// skillExists prüft, ob eine Fertigkeit mit dem gegebenen Namen existiert
func skillExists(skillName string) bool {
	return determineSkillType(skillName) != ""
}
