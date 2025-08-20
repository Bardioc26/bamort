package importer

import (
	"bamort/models"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// sourceIDCache stores cached source ID lookups to avoid repeated database queries
var sourceIDCache = make(map[string]uint)

// ClearSourceCache clears the cached source ID lookups
// Useful for testing or when working with different databases
func ClearSourceCache() {
	sourceIDCache = make(map[string]uint)
}

// lookupSourceID retrieves the source ID from the database based on the source code
// Uses caching to avoid repeated database calls
// If the source doesn't exist, it creates a new one automatically
func lookupSourceID(sourceCode string) (uint, error) {
	if sourceCode == "" {
		return 0, fmt.Errorf("source code is empty")
	}

	// Check cache first
	if sourceID, cached := sourceIDCache[sourceCode]; cached {
		return sourceID, nil
	}

	// Not in cache, look up from database
	var source models.Source
	if err := source.FirstByCode(sourceCode); err != nil {
		// Source not found, create it automatically
		newSource := models.Source{
			Code:       sourceCode,
			Name:       sourceCode, // Use code as name initially
			GameSystem: "midgard",  // Default game system
			IsActive:   true,       // Set as active by default
		}

		if createErr := newSource.Create(); createErr != nil {
			return 0, fmt.Errorf("source with code '%s' not found and could not be created: %w", sourceCode, createErr)
		}

		// Cache the newly created source
		sourceIDCache[sourceCode] = newSource.ID
		return newSource.ID, nil
	}

	// Cache the result for future use
	sourceIDCache[sourceCode] = source.ID
	return source.ID, nil
}

func ImportChar(char CharacterImport) (*models.Char, error) {
	return nil, fmt.Errorf("char could not be imported %s", "Weil Wegen Kommt noch")
}

func CheckSkill(fertigkeit *Fertigkeit, autocreate bool) (*models.Skill, error) {
	stammF := models.Skill{}
	//err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, fertigkeit.Name).Error
	err := stammF.First(fertigkeit.Name)
	if err == nil {
		// Fertigkeit found
		return &stammF, nil
	}
	if !autocreate {
		return nil, fmt.Errorf("does not exist in Fertigkeit importer")
	}
	stammF.GameSystem = "midgard"
	stammF.Name = fertigkeit.Name
	if stammF.Name != "Sprache" {
		stammF.Beschreibung = fertigkeit.Beschreibung
	}
	if fertigkeit.Fertigkeitswert < 12 {
		stammF.Initialwert = 5
	} else {
		stammF.Initialwert = 12
	}
	stammF.Bonuseigenschaft = "keine"
	stammF.Quelle = fertigkeit.Quelle
	//fmt.Println(stammF)

	err = stammF.Create()
	if err != nil {
		// Fertigkeit found
		return nil, err
	}

	//err = database.DB.First(&stammF, "system=? AND name = ?", gameSystem, fertigkeit.Name).Error
	err = stammF.First(fertigkeit.Name)
	if err != nil {
		// Fertigkeit found
		return nil, err
	}
	return &stammF, nil
}

func CheckSpell(zauber *Zauber, autocreate bool) (*models.Spell, error) {
	stammF := models.Spell{}

	//err := database.DB.First(&stammF, "system=? AND name = ?", gameSystem, zauber.Name).Error
	err := stammF.First(zauber.Name)
	if err == nil {
		// zauber found
		return &stammF, nil
	}
	if !autocreate {
		return nil, fmt.Errorf("does not exist in zauber importer")
	}
	stammF.GameSystem = "midgard"
	stammF.Name = zauber.Name
	stammF.Beschreibung = zauber.Beschreibung
	stammF.AP = "1"
	stammF.Stufe = 1
	stammF.Wirkungsziel = "Zauberer"
	stammF.Reichweite = "15 m"

	stammF.Quelle = zauber.Quelle
	//fmt.Println(stammF)
	err = stammF.Create()
	if err != nil {
		// spell found
		return nil, err
	}

	//err = database.DB.First(&stammF, "system=? AND name = ?", gameSystem, zauber.Name).Error
	err = stammF.First(zauber.Name)
	if err != nil {
		// spell found
		return nil, err
	}
	return &stammF, nil
}

func ImportCsv2Spell(filepath string) error {
	// Open the CSV file
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filepath, err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV file: %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}

	// Get headers from first row
	headers := records[0]

	// Create a map for field name mapping to handle case variations
	fieldMap := make(map[string]int)
	for i, header := range headers {
		// Normalize header names for mapping
		normalizedHeader := strings.ToLower(strings.TrimSpace(header))
		fieldMap[normalizedHeader] = i
	}

	// Process each record (skip header row)
	for i, record := range records[1:] {
		if len(record) == 0 {
			continue // Skip empty rows
		}

		// Create spell struct
		spell := models.Spell{
			GameSystem: "midgard", // Default value
		}

		// Map CSV fields to struct fields
		if idx, exists := fieldMap["name"]; exists && idx < len(record) {
			spell.Name = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["beschreibung"]; exists && idx < len(record) {
			spell.Beschreibung = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["quelle"]; exists && idx < len(record) {
			quelleCode := strings.ToUpper(strings.TrimSpace(record[idx]))
			spell.Quelle = quelleCode

			// Look up source_id from database based on quelle code using cached lookup
			if quelleCode != "" {
				if sourceID, err := lookupSourceID(quelleCode); err == nil {
					spell.SourceID = sourceID
				}
				// If source lookup fails, keep the existing source_id from CSV if present
			}
		}
		if idx, exists := fieldMap["game_system"]; exists && idx < len(record) && strings.TrimSpace(record[idx]) != "" {
			spell.GameSystem = strings.ToLower(strings.TrimSpace(record[idx]))
		}
		/*
			if idx, exists := fieldMap["source_id"]; exists && idx < len(record) {
				if sourceID, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
					spell.SourceID = uint(sourceID)
				}
			}
		*/
		if idx, exists := fieldMap["page_number"]; exists && idx < len(record) {
			if pageNum, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
				spell.PageNumber = pageNum
			}
		}
		if idx, exists := fieldMap["bonus"]; exists && idx < len(record) {
			if bonus, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
				spell.Bonus = bonus
			}
		}
		if idx, exists := fieldMap["stufe"]; exists && idx < len(record) {
			if level, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
				spell.Stufe = level
			}
		}
		if idx, exists := fieldMap["ap"]; exists && idx < len(record) {
			spell.AP = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["art"]; exists && idx < len(record) {
			spell.Art = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["zauberdauer"]; exists && idx < len(record) {
			spell.Zauberdauer = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["reichweite"]; exists && idx < len(record) {
			spell.Reichweite = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["wirkungsziel"]; exists && idx < len(record) {
			spell.Wirkungsziel = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["wirkungsbereich"]; exists && idx < len(record) {
			spell.Wirkungsbereich = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["wirkungsdauer"]; exists && idx < len(record) {
			spell.Wirkungsdauer = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["ursprung"]; exists && idx < len(record) {
			spell.Ursprung = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["category"]; exists && idx < len(record) {
			spell.Category = strings.TrimSpace(record[idx])
		}
		if idx, exists := fieldMap["learning_category"]; exists && idx < len(record) {
			spell.LearningCategory = strings.TrimSpace(record[idx])
		}

		// Skip if name is empty
		if spell.Name == "" {
			continue
		}

		// Try to find existing spell by name
		existingSpell := models.Spell{}
		err := existingSpell.First(spell.Name)

		if err == nil {
			// Spell exists, update it
			existingSpell.Beschreibung = spell.Beschreibung
			existingSpell.Quelle = spell.Quelle
			existingSpell.GameSystem = spell.GameSystem
			// Update SourceID if we found one from quelle lookup
			if spell.SourceID != 0 {
				existingSpell.SourceID = spell.SourceID
			}
			if spell.PageNumber != 0 {
				existingSpell.PageNumber = spell.PageNumber
			}
			existingSpell.Bonus = spell.Bonus
			existingSpell.Stufe = spell.Stufe
			existingSpell.AP = spell.AP
			existingSpell.Art = spell.Art
			existingSpell.Zauberdauer = spell.Zauberdauer
			existingSpell.Reichweite = spell.Reichweite
			existingSpell.Wirkungsziel = spell.Wirkungsziel
			existingSpell.Wirkungsbereich = spell.Wirkungsbereich
			existingSpell.Wirkungsdauer = spell.Wirkungsdauer
			existingSpell.Ursprung = spell.Ursprung
			existingSpell.Category = spell.Category
			existingSpell.LearningCategory = spell.LearningCategory

			err = existingSpell.Save()
			if err != nil {
				return fmt.Errorf("error updating spell '%s' at row %d: %w", spell.Name, i+2, err)
			}
		} else {
			// Spell doesn't exist, create new one
			err = spell.Create()
			if err != nil {
				return fmt.Errorf("error creating spell '%s' at row %d: %w", spell.Name, i+2, err)
			}
		}
	}

	return nil
}
