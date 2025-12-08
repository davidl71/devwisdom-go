package main

import (
	"fmt"
	"log"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

func main() {
	// Initialize engine
	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}

	// Create a custom project source
	config := &wisdom.SourceConfig{
		ID:          "my_project",
		Name:        "My Project Wisdom",
		Icon:        "🚀",
		Description: "Custom wisdom quotes for my project",
		Quotes: map[string][]wisdom.Quote{
			"chaos": {
				{
					Quote:        "When everything breaks, remember: you built this.",
					Source:       "Project Wisdom",
					Encouragement: "You can fix it.",
				},
				{
					Quote:        "The best debugging happens at 2 AM.",
					Source:       "Project Wisdom",
					Encouragement: "Take a break, come back fresh.",
				},
			},
			"lower_aeons": {
				{
					Quote:        "Progress is progress, even if it's slow.",
					Source:       "Project Wisdom",
					Encouragement: "Keep moving forward.",
				},
			},
			"middle_aeons": {
				{
					Quote:        "We're getting there, one commit at a time.",
					Source:       "Project Wisdom",
					Encouragement: "Consistency wins.",
				},
			},
			"upper_aeons": {
				{
					Quote:        "The architecture is solid, the tests are passing.",
					Source:       "Project Wisdom",
					Encouragement: "You've built something good.",
				},
			},
			"treasury": {
				{
					Quote:        "Everything is working perfectly!",
					Source:       "Project Wisdom",
					Encouragement: "Enjoy the moment.",
				},
			},
		},
	}

	// Add and save to project directory
	if err := engine.AddProjectSource(config); err != nil {
		log.Fatalf("Failed to add project source: %v", err)
	}

	fmt.Println("✅ Project source added successfully!")
	fmt.Printf("📁 Saved to: %s\n", engine.GetLoader().GetProjectSourcesPath())

	// Test the new source
	quote, err := engine.GetWisdom(75.0, "my_project")
	if err != nil {
		log.Fatalf("Failed to get wisdom: %v", err)
	}

	fmt.Println("\n📜 Test quote from your source:")
	fmt.Printf("   %s\n", quote.Quote)
	fmt.Printf("   — %s\n", quote.Source)
	fmt.Printf("   💡 %s\n", quote.Encouragement)
}

