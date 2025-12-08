package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// runBriefing handles the briefing command
func (a *App) runBriefing(args []string) error {
	fs := flag.NewFlagSet("briefing", flag.ExitOnError)
	days := fs.Int("days", 1, "Number of days to include in briefing")
	jsonOutput := fs.Bool("json", false, "Output in JSON format")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Initialize wisdom engine
	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize wisdom engine: %w", err)
	}

	// Get daily quote (for now, just get one quote)
	// TODO: Implement proper daily briefing with multiple quotes and consultation log
	sources := engine.ListSources()
	if len(sources) == 0 {
		return fmt.Errorf("no wisdom sources available")
	}
	// Use first available source (TODO: implement daily random source selector)
	quote, err := engine.GetWisdom(50.0, sources[0])
	if err != nil {
		return fmt.Errorf("failed to get wisdom: %w", err)
	}

	// Output
	if *jsonOutput {
		result := map[string]interface{}{
			"days":          *days,
			"quote":         quote.Quote,
			"source":        quote.Source,
			"encouragement": quote.Encouragement,
			"note":          "Full briefing with consultation log coming soon",
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	// Human-readable output
	fmt.Printf("Daily Briefing (%d day(s))\n\n", *days)
	if quote.WisdomIcon != "" {
		fmt.Printf("%s ", quote.WisdomIcon)
	}
	fmt.Printf("\"%s\"\n", quote.Quote)
	if quote.Encouragement != "" {
		fmt.Printf("— %s\n", quote.Encouragement)
	}
	if quote.Source != "" {
		fmt.Printf("Source: %s\n", quote.Source)
	}
	fmt.Println()
	fmt.Println("Note: Full briefing with consultation log coming soon")

	return nil
}
