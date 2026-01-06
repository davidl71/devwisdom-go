package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// runConsult handles the consult command
func (a *App) runConsult(args []string) error {
	fs := flag.NewFlagSet("consult", flag.ExitOnError)
	metric := fs.String("metric", "", "Metric name (e.g., security, testing)")
	tool := fs.String("tool", "", "Tool name (e.g., project_scorecard)")
	stage := fs.String("stage", "", "Stage name (e.g., daily_checkin)")
	score := fs.Float64("score", 50.0, "Project score (0-100)")
	jsonOutput := fs.Bool("json", false, "Output in JSON format")
	quiet := fs.Bool("quiet", false, "Output only the quote text")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validate that at least one of metric, tool, or stage is provided
	if *metric == "" && *tool == "" && *stage == "" {
		return fmt.Errorf("must provide at least one of --metric, --tool, or --stage: use 'devwisdom consult --help' for usage examples")
	}

	// Initialize wisdom engine
	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize wisdom engine (check sources.json configuration): %w", err)
	}

	// Determine advisor
	var advisorInfo *wisdom.AdvisorInfo
	var err error
	advisorRegistry := engine.GetAdvisors()

	if *metric != "" {
		advisorInfo, err = advisorRegistry.GetAdvisorForMetric(*metric)
		if err != nil {
			return fmt.Errorf("no advisor found for metric %q: %w. Use 'devwisdom advisors' to see available metrics", *metric, err)
		}
	} else if *tool != "" {
		advisorInfo, err = advisorRegistry.GetAdvisorForTool(*tool)
		if err != nil {
			return fmt.Errorf("no advisor found for tool %q: %w. Use 'devwisdom advisors' to see available tools", *tool, err)
		}
	} else if *stage != "" {
		advisorInfo, err = advisorRegistry.GetAdvisorForStage(*stage)
		if err != nil {
			return fmt.Errorf("no advisor found for stage %q: %w. Use 'devwisdom advisors' to see available stages", *stage, err)
		}
	}

	// Get quote from advisor's source
	aeonLevel := wisdom.GetAeonLevel(*score)

	// Try to get quote from advisor's source
	source, exists := engine.GetSource(advisorInfo.Advisor)
	if !exists {
		// Fallback: try to get quote from any source
		sources := engine.ListSources()
		if len(sources) == 0 {
			return fmt.Errorf("no wisdom sources available: ensure sources.json exists and contains valid source definitions. Use 'devwisdom sources' to verify")
		}
		var quote *wisdom.Quote
		quote, err = engine.GetWisdom(*score, sources[0])
		if err != nil {
			return fmt.Errorf("failed to get wisdom quote (source: %q, score: %.1f): %w", sources[0], *score, err)
		}

		// Output
		if *quiet {
			fmt.Println(quote.Quote)
			return nil
		}

		if *jsonOutput {
			result := map[string]interface{}{
				"advisor":       advisorInfo.Advisor,
				"advisor_icon":  advisorInfo.Icon,
				"rationale":     advisorInfo.Rationale,
				"quote":         quote.Quote,
				"source":        quote.Source,
				"encouragement": quote.Encouragement,
			}
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			return encoder.Encode(result)
		}

		// Human-readable output
		if advisorInfo.Icon != "" {
			fmt.Printf("%s ", advisorInfo.Icon)
		}
		fmt.Printf("Advisor: %s\n", advisorInfo.Advisor)
		if advisorInfo.Rationale != "" {
			fmt.Printf("Rationale: %s\n", advisorInfo.Rationale)
		}
		fmt.Printf("\n\"%s\"\n", quote.Quote)
		if quote.Encouragement != "" {
			fmt.Printf("— %s\n", quote.Encouragement)
		}
		return nil
	}

	// Get quote from advisor's source
	quote := source.GetQuote(aeonLevel)

	// Output
	if *quiet {
		fmt.Println(quote.Quote)
		return nil
	}

	if *jsonOutput {
		result := map[string]interface{}{
			"advisor":       advisorInfo.Advisor,
			"advisor_icon":  advisorInfo.Icon,
			"rationale":     advisorInfo.Rationale,
			"quote":         quote.Quote,
			"source":        quote.Source,
			"encouragement": quote.Encouragement,
		}
		if *metric != "" {
			result["metric"] = *metric
		}
		if *tool != "" {
			result["tool"] = *tool
		}
		if *stage != "" {
			result["stage"] = *stage
		}
		result["score"] = *score
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	// Human-readable output
	if advisorInfo.Icon != "" {
		fmt.Printf("%s ", advisorInfo.Icon)
	}
	fmt.Printf("Advisor: %s\n", advisorInfo.Advisor)
	if advisorInfo.Rationale != "" {
		fmt.Printf("Rationale: %s\n", advisorInfo.Rationale)
	}
	if advisorInfo.HelpsWith != "" {
		fmt.Printf("Helps with: %s\n", advisorInfo.HelpsWith)
	}
	fmt.Printf("\n\"%s\"\n", quote.Quote)
	if quote.Encouragement != "" {
		fmt.Printf("— %s\n", quote.Encouragement)
	}
	if quote.Source != "" {
		fmt.Printf("Source: %s\n", quote.Source)
	}

	return nil
}
