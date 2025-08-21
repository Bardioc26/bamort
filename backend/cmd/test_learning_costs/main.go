package main

import (
	"bamort/database"
	"bamort/gsmaster"
	"bamort/models"
	"fmt"
	"strings"
)

func main() {
	// Initialize config and database
	database.SetupTestDB(true, true)
	defer database.ResetTestDB()

	// Test some skills for Magier character class
	testSkills := []string{"Schreiben", "Athletik", "Lesen von Zauberformeln", "Zaubern"}
	characterClass := "Ma" // Magier abbreviation

	fmt.Printf("Testing learning costs for character class: %s\n", characterClass)
	fmt.Println(strings.Repeat("=", 50))

	for _, skillName := range testSkills {
		// Test the same logic as our handler
		fmt.Printf("\n--- Testing skill: %s ---\n", skillName)

		// Try gsmaster.FindBestCategoryForSkillLearningOld
		bestCategory, difficulty, err := gsmaster.FindBestCategoryForSkillLearningOld(skillName, characterClass)
		if err == nil && bestCategory != "" {
			fmt.Printf("✅ FindBestCategory: Category=%s, Difficulty=%s\n", bestCategory, difficulty)
		} else {
			fmt.Printf("❌ FindBestCategory: ERROR: %v\n", err)
		}

		// Try models.GetSkillCategoryAndDifficultyNewSystem
		skillLearningInfo, err := models.GetSkillCategoryAndDifficultyNewSystem(skillName, characterClass)
		if err == nil {
			fmt.Printf("✅ GetSkillCategory: Category=%s, LearnCost=%d\n",
				skillLearningInfo.CategoryName, skillLearningInfo.LearnCost)
		} else {
			fmt.Printf("❌ GetSkillCategory: ERROR: %v\n", err)
		}
	}

	// Also test with empty character class
	fmt.Println("\nTesting with empty character class:")
	for _, skillName := range testSkills {
		skillLearningInfo, err := models.GetSkillCategoryAndDifficultyNewSystem(skillName, "")
		if err != nil {
			fmt.Printf("❌ Skill: %s (empty class) - ERROR: %v\n", skillName, err)
			continue
		}

		fmt.Printf("✅ Skill: %s (empty class) - Category: %s, LearnCost: %d\n",
			skillName,
			skillLearningInfo.CategoryName,
			skillLearningInfo.LearnCost)
	}
}
