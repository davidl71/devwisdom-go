package main

import (
	"fmt"
	"log"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// Example 1: Basic wisdom quote retrieval
func exampleBasicQuote() {
	fmt.Println("=== Example 1: Basic Wisdom Quote ===")

	// Initialize the wisdom engine
	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	// Get a quote with a specific score and source
	quote, err := engine.GetWisdom(75.0, "stoic")
	if err != nil {
		log.Fatalf("Failed to get wisdom: %v", err)
	}

	fmt.Printf("Quote: %s\n", quote.Quote)
	fmt.Printf("Source: %s\n", quote.Source)
	fmt.Printf("Encouragement: %s\n", quote.Encouragement)
	fmt.Println()
}

// Example 2: Random source selection
func exampleRandomSource() {
	fmt.Println("=== Example 2: Random Source Selection ===")

	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	// Get a date-seeded random source (same source all day)
	randomSource, err := engine.GetRandomSource(true)
	if err != nil {
		log.Fatalf("Failed to get random source: %v", err)
	}

	fmt.Printf("Today's random source: %s\n", randomSource)

	// Get quote from random source
	quote, err := engine.GetWisdom(50.0, randomSource)
	if err != nil {
		log.Fatalf("Failed to get wisdom: %v", err)
	}

	fmt.Printf("Quote: %s\n", quote.Quote)
	fmt.Println()
}

// Example 3: Advisor consultation
func exampleAdvisorConsultation() {
	fmt.Println("=== Example 3: Advisor Consultation ===")

	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	// Get advisor registry
	advisors := engine.GetAdvisors()

	// Consult advisor for a metric
	advisorInfo, err := advisors.GetAdvisorForMetric("security")
	if err != nil {
		log.Fatalf("Failed to get advisor: %v", err)
	}

	fmt.Printf("Advisor: %s (%s)\n", advisorInfo.Advisor, advisorInfo.Icon)
	fmt.Printf("Rationale: %s\n", advisorInfo.Rationale)

	// Get quote from advisor's source
	score := 40.0
	quote, err := engine.GetWisdom(score, advisorInfo.Advisor)
	if err != nil {
		log.Fatalf("Failed to get wisdom: %v", err)
	}

	fmt.Printf("Quote: %s\n", quote.Quote)
	fmt.Printf("Encouragement: %s\n", quote.Encouragement)
	fmt.Println()
}

// Example 4: Listing sources and advisors
func exampleListResources() {
	fmt.Println("=== Example 4: Listing Sources and Advisors ===")

	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	// List all sources
	sources := engine.ListSources()
	fmt.Printf("Available sources (%d):\n", len(sources))
	for i, sourceID := range sources {
		if i >= 5 { // Show first 5
			fmt.Printf("  ... and %d more\n", len(sources)-5)
			break
		}
		source, exists := engine.GetSource(sourceID)
		if exists {
			fmt.Printf("  %s %s (%s)\n", source.Icon, source.Name, sourceID)
		}
	}
	fmt.Println()

	// List advisors
	advisors := engine.GetAdvisors()

	// Metric advisors
	metricAdvisors := advisors.GetAllMetricAdvisors()
	fmt.Printf("Metric advisors (%d):\n", len(metricAdvisors))
	count := 0
	for metric, info := range metricAdvisors {
		if count >= 5 {
			fmt.Printf("  ... and %d more\n", len(metricAdvisors)-5)
			break
		}
		fmt.Printf("  %s → %s\n", metric, info.Advisor)
		count++
	}
	fmt.Println()
}

// Example 5: Working with aeon levels
func exampleAeonLevels() {
	fmt.Println("=== Example 5: Aeon Levels ===")

	scores := []float64{15.0, 40.0, 60.0, 80.0, 90.0}

	for _, score := range scores {
		aeonLevel := wisdom.GetAeonLevel(score)
		mode := wisdom.GetConsultationMode(score)
		fmt.Printf("Score: %.1f → Aeon Level: %s, Mode: %s %s\n",
			score, aeonLevel, mode.Icon, mode.Name)
	}
	fmt.Println()
}

// Example 6: Custom source configuration
func exampleCustomSource() {
	fmt.Println("=== Example 6: Custom Source Configuration ===")

	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	// Create a custom source
	config := &wisdom.SourceConfig{
		ID:          "team_wisdom",
		Name:        "Team Wisdom",
		Icon:        "👥",
		Description: "Wisdom quotes for our development team",
		Quotes: map[string][]wisdom.Quote{
			"middle_aeons": {
				{
					Quote:         "Code reviews make us better.",
					Source:        "Team Wisdom",
					Encouragement: "Every review is a learning opportunity.",
				},
			},
			"upper_aeons": {
				{
					Quote:         "We ship quality code.",
					Source:        "Team Wisdom",
					Encouragement: "Excellence is our standard.",
				},
			},
		},
	}

	// Add the source
	if err := engine.AddProjectSource(config); err != nil {
		log.Fatalf("Failed to add source: %v", err)
	}

	fmt.Println("✅ Custom source added successfully!")

	// Use the new source
	quote, err := engine.GetWisdom(75.0, "team_wisdom")
	if err != nil {
		log.Fatalf("Failed to get wisdom: %v", err)
	}

	fmt.Printf("Quote from custom source: %s\n", quote.Quote)
	fmt.Println()
}

// Example 7: Consultation mode and frequency
func exampleConsultationMode() {
	fmt.Println("=== Example 7: Consultation Mode ===")

	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	scores := []float64{25.0, 45.0, 65.0, 85.0}

	for _, score := range scores {
		mode := wisdom.GetConsultationMode(score)
		fmt.Printf("Score: %.1f → Mode: %s %s\n", score, mode.Icon, mode.Name)
		fmt.Printf("  Frequency: %s\n", mode.Frequency)
		fmt.Printf("  Description: %s\n", mode.Description)
		fmt.Println()
	}
}

func main() {
	fmt.Println("devwisdom-go Programmatic API Examples")
	fmt.Println("=====================================")
	fmt.Println()

	exampleBasicQuote()
	exampleRandomSource()
	exampleAdvisorConsultation()
	exampleListResources()
	exampleAeonLevels()
	exampleCustomSource()
	exampleConsultationMode()

	fmt.Println("✅ All examples completed!")
}
