package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// runSources handles the sources command
func (a *App) runSources(args []string) error {
	fs := flag.NewFlagSet("sources", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output in JSON format")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Initialize wisdom engine
	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize wisdom engine: %w", err)
	}

	// Get sources
	sourceIDs := engine.ListSources()
	sort.Strings(sourceIDs)

	// Build source details
	sources := make([]map[string]interface{}, 0, len(sourceIDs))
	for _, id := range sourceIDs {
		source, exists := engine.GetSource(id)
		if !exists {
			continue
		}
		sourceInfo := map[string]interface{}{
			"id":   id,
			"name": source.Name,
			"icon": source.Icon,
		}
		if source.Description != "" {
			sourceInfo["description"] = source.Description
		}
		sources = append(sources, sourceInfo)
	}

	// Output
	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(sources)
	}

	// Human-readable output
	fmt.Printf("Available Wisdom Sources (%d):\n\n", len(sources))
	for _, src := range sources {
		icon := src["icon"].(string)
		name := src["name"].(string)
		id := src["id"].(string)

		if icon != "" {
			fmt.Printf("%s %s", icon, name)
		} else {
			fmt.Printf("%s", name)
		}
		if id != name {
			fmt.Printf(" (%s)", id)
		}
		fmt.Println()

		if desc, ok := src["description"].(string); ok && desc != "" {
			fmt.Printf("  %s\n", desc)
		}
		fmt.Println()
	}

	return nil
}
