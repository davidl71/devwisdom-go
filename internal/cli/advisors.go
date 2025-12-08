package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// runAdvisors handles the advisors command
func (a *App) runAdvisors(args []string) error {
	fs := flag.NewFlagSet("advisors", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output in JSON format")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Initialize wisdom engine
	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize wisdom engine: %w", err)
	}

	// Get advisors registry
	advisors := engine.GetAdvisors()

	// Build advisor list (we'll need to expose the mappings from AdvisorRegistry)
	// For now, return a placeholder message
	advisorList := []map[string]interface{}{
		{
			"note":  "Advisor mappings are available via consult command",
			"usage": "Use 'devwisdom consult --metric <metric>' or 'devwisdom consult --tool <tool>' or 'devwisdom consult --stage <stage>'",
		},
	}

	// Output
	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(advisorList)
	}

	// Human-readable output
	fmt.Println("Available Advisors:")
	fmt.Println()
	fmt.Println("Advisors are mapped to metrics, tools, and stages.")
	fmt.Println("Use the consult command to get advisor guidance:")
	fmt.Println()
	fmt.Println("  devwisdom consult --metric security")
	fmt.Println("  devwisdom consult --tool project_scorecard")
	fmt.Println("  devwisdom consult --stage daily_checkin")
	fmt.Println()
	fmt.Println("For more information, see: https://github.com/davidl71/devwisdom-go")

	_ = advisors // Suppress unused variable warning
	return nil
}
