package main

import (
	"bamort/database"
	"bamort/gsmaster"
	"flag"
	"log"
	"os"
)

func main() {
	var initLearning = flag.Bool("init-learning", false, "Initialize learning costs system")
	var validateLearning = flag.Bool("validate-learning", false, "Validate learning costs data")
	var summaryLearning = flag.Bool("summary-learning", false, "Show learning costs summary")
	flag.Parse()

	// Datenbank verbinden
	database.ConnectDatabase()
	if database.DB == nil {
		log.Fatal("Failed to connect to database")
	}

	if *initLearning {
		log.Println("Starting learning costs system initialization...")
		if err := gsmaster.InitializeLearningCostsSystem(); err != nil {
			log.Fatalf("Failed to initialize learning costs system: %v", err)
		}
		log.Println("Learning costs system initialized successfully!")
		os.Exit(0)
	}

	if *validateLearning {
		log.Println("Validating learning costs data...")
		if err := gsmaster.ValidateLearningCostsData(); err != nil {
			log.Fatalf("Validation failed: %v", err)
		}
		log.Println("Validation completed successfully!")
		os.Exit(0)
	}

	if *summaryLearning {
		log.Println("Getting learning costs summary...")
		summary, err := gsmaster.GetLearningCostsSummary()
		if err != nil {
			log.Fatalf("Failed to get summary: %v", err)
		}

		log.Println("Learning Costs Summary:")
		for table, count := range summary {
			log.Printf("  %s: %v", table, count)
		}
		os.Exit(0)
	}

	log.Println("Usage:")
	log.Println("  -init-learning     Initialize learning costs system")
	log.Println("  -validate-learning Validate learning costs data")
	log.Println("  -summary-learning  Show learning costs summary")
}
