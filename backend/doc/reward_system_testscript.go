package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"bamort/character"
)

func main() {
	fmt.Println("=== Belohnungssystem API Tests ===\n")

	// Warten bis Server gestartet ist
	time.Sleep(2 * time.Second)

	// Base URL für lokale API
	baseURL := "http://localhost:8080/api/characters/1/skill-cost"

	fmt.Println("1. Kostenlose Fertigkeit lernen (free_learning)")
	testRequest1 := character.SkillCostRequest{
		Name:   "Dolch",
		Type:   "skill",
		Action: "learn",
		Reward: &character.RewardOptions{
			Type: "free_learning",
		},
	}
	testRewardRequest(baseURL, testRequest1)

	fmt.Println("\n2. Kostenloser Zauber lernen (free_spell_learning)")
	testRequest2 := character.SkillCostRequest{
		Name:   "Licht",
		Type:   "spell",
		Action: "learn",
		Reward: &character.RewardOptions{
			Type: "free_spell_learning",
		},
	}
	testRewardRequest(baseURL, testRequest2)

	fmt.Println("\n3. Halbe EP für Verbesserung (half_ep_improvement)")
	testRequest3 := character.SkillCostRequest{
		Name:         "Dolch",
		Type:         "skill",
		Action:       "improve",
		CurrentLevel: 8,
		Reward: &character.RewardOptions{
			Type: "half_ep_improvement",
		},
	}
	testRewardRequest(baseURL, testRequest3)

	fmt.Println("\n4. Gold statt EP verwenden (gold_for_ep)")
	testRequest4 := character.SkillCostRequest{
		Name:         "Schwert",
		Type:         "skill",
		Action:       "improve",
		CurrentLevel: 10,
		Reward: &character.RewardOptions{
			Type:         "gold_for_ep",
			UseGoldForEP: true,
			MaxGoldEP:    0, // Automatisch die Hälfte
		},
	}
	testRewardRequest(baseURL, testRequest4)

	fmt.Println("\n5. Kombination: Belohnung + Praxispunkte")
	testRequest5 := character.SkillCostRequest{
		Name:         "Dolch",
		Type:         "skill",
		Action:       "improve",
		CurrentLevel: 12,
		UsePP:        2,
		Reward: &character.RewardOptions{
			Type: "half_ep_improvement",
		},
	}
	testRewardRequest(baseURL, testRequest5)
}

func testRewardRequest(url string, request character.SkillCostRequest) {
	// Convert request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshaling request: %v\n", err)
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d\n", resp.StatusCode)
		return
	}

	// Parse response
	var response character.SkillCostResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("Error parsing response: %v\n", err)
		return
	}

	// Print results
	fmt.Printf("Fertigkeit: %s (%s)\n", response.SkillName, response.Action)
	if response.RewardApplied != "" {
		fmt.Printf("Belohnung: %s\n", response.RewardApplied)
	}

	if response.OriginalCostStruct != nil {
		fmt.Printf("Ursprüngliche Kosten: %d EP, %d LE, %d GS\n",
			response.OriginalCostStruct.Ep,
			response.OriginalCostStruct.LE,
			response.OriginalCostStruct.Money)
	}

	fmt.Printf("Finale Kosten: %d EP, %d LE, %d GS\n",
		response.Ep, response.LE, response.Money)

	if response.Savings != nil && (response.Savings.Ep > 0 || response.Savings.LE > 0 || response.Savings.Money > 0) {
		fmt.Printf("Ersparnisse: %d EP, %d LE, %d GS\n",
			response.Savings.Ep, response.Savings.LE, response.Savings.Money)
	}

	if response.GoldUsedForEP > 0 {
		fmt.Printf("Gold für EP verwendet: %d EP = %d GS\n",
			response.GoldUsedForEP, response.GoldUsedForEP*10)
	}

	if response.PPUsed > 0 {
		fmt.Printf("Praxispunkte verwendet: %d\n", response.PPUsed)
	}

	fmt.Printf("Kann sich leisten: %t\n", response.CanAfford)
}
