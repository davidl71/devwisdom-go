package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// runQuote handles the quote command
func (a *App) runQuote(args []string) error {
	fs := flag.NewFlagSet("quote", flag.ExitOnError)
	source := fs.String("source", "", "Wisdom source name (e.g., stoic, tao, pistis_sophia)")
	score := fs.Float64("score", 50.0, "Project score (0-100) for aeon level selection")
	jsonOutput := fs.Bool("json", false, "Output in JSON format")
	quiet := fs.Bool("quiet", false, "Output only the quote text")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Initialize wisdom engine
	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize wisdom engine (check sources.json configuration): %w", err)
	}

	// Determine source
	selectedSource := *source
	if selectedSource == "" {
		// Use date-seeded random source selector for daily consistency
		randomSource, err := engine.GetRandomSource(true)
		if err != nil {
			return fmt.Errorf("failed to get random source (no sources available): %w", err)
		}
		selectedSource = randomSource
	}

	// Get quote - if source is empty string, GetWisdom will fail, so we handle it above
	var quote *wisdom.Quote
	var err error
	if selectedSource != "" {
		quote, err = engine.GetWisdom(*score, selectedSource)
	} else {
		// This shouldn't happen due to check above, but handle it
		return fmt.Errorf("no source specified and no sources available: ensure sources.json exists and contains valid source definitions. Use 'devwisdom sources' to list available sources")
	}
	if err != nil {
		scoreStr := "default"
		if score != nil {
			scoreStr = fmt.Sprintf("%.1f", *score)
		}
		return fmt.Errorf("failed to get wisdom quote (source: %q, score: %s): %w", selectedSource, scoreStr, err)
	}

	// Output based on format
	if *quiet {
		fmt.Println(quote.Quote)
		return nil
	}

	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(quote)
	}

	// Human-readable output
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
	if selectedSource != quote.Source {
		fmt.Printf("Wisdom Source: %s\n", selectedSource)
	}

	return nil
}
