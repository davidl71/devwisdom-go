package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/davidl71/devwisdom-go/internal/wisdom"
)

// runBriefing handles the briefing command
func (a *App) runBriefing(args []string) error {
	fs := flag.NewFlagSet("briefing", flag.ExitOnError)
	days := fs.Int("days", 1, "Number of days to include in briefing")
	score := fs.Float64("score", 50.0, "Overall project score (0-100) for consultation mode")
	jsonOutput := fs.Bool("json", false, "Output in JSON format")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Initialize wisdom engine
	engine := wisdom.NewEngine()
	if err := engine.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize wisdom engine: %w", err)
	}

	// Get consultation mode based on provided score
	mode := wisdom.GetConsultationMode(*score)

	// For now, use sample metric scores (can be parameterized or from consultation log later)
	// Note: Full consultation log integration comes in Phase 5
	sampleMetricScores := map[string]float64{
		"security":      40.0,
		"testing":       30.0,
		"documentation": 60.0,
		"completion":    70.0,
		"alignment":     50.0,
	}

	// Find lowest scoring metrics (need most advice)
	type metricScore struct {
		metric string
		score  float64
	}
	metrics := make([]metricScore, 0, len(sampleMetricScores))
	for m, s := range sampleMetricScores {
		metrics = append(metrics, metricScore{metric: m, score: s})
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].score < metrics[j].score
	})

	// Get advisors
	advisors := engine.GetAdvisors()

	// Build briefing with advisor quotes for lowest 3 metrics
	briefingQuotes := make([]map[string]interface{}, 0)
	for i, ms := range metrics {
		if i >= 3 {
			break
		}

		advisorInfo, err := advisors.GetAdvisorForMetric(ms.metric)
		if err != nil {
			continue
		}

		// Get quote from advisor's source
		source, exists := engine.GetSource(advisorInfo.Advisor)
		if !exists {
			continue
		}

		aeonLevel := wisdom.GetAeonLevel(ms.score)
		quote := source.GetQuote(aeonLevel)

		briefingQuotes = append(briefingQuotes, map[string]interface{}{
			"metric":        ms.metric,
			"score":         ms.score,
			"advisor":       advisorInfo.Advisor,
			"advisor_icon":  advisorInfo.Icon,
			"quote":         quote.Quote,
			"source":        quote.Source,
			"encouragement": quote.Encouragement,
		})
	}

	// Output
	if *jsonOutput {
		result := map[string]interface{}{
			"days":          *days,
			"overall_score": *score,
			"consultation_mode": map[string]interface{}{
				"name":        mode.Name,
				"icon":        mode.Icon,
				"frequency":   mode.Frequency,
				"description": mode.Description,
			},
			"advisory_quotes": briefingQuotes,
			"note":            "Consultation log integration coming in Phase 5",
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	// Human-readable output (formatted similar to Python version)
	fmt.Printf("╔══════════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║  🌅 DAILY ADVISOR BRIEFING                                           ║\n")
	fmt.Printf("║  Overall Score: %.1f%% | Mode: %s %-30s ║\n", *score, mode.Icon, strings.ToUpper(mode.Name))
	fmt.Printf("╠══════════════════════════════════════════════════════════════════════╣\n")
	fmt.Println()

	for _, bq := range briefingQuotes {
		metric := bq["metric"].(string)
		score := bq["score"].(float64)
		advisorIcon := bq["advisor_icon"].(string)
		advisor := bq["advisor"].(string)
		quote := bq["quote"].(string)
		encouragement := bq["encouragement"].(string)

		quotePreview := quote
		if len(quotePreview) > 55 {
			quotePreview = quotePreview[:55] + "..."
		}
		encPreview := encouragement
		if len(encPreview) > 55 {
			encPreview = encPreview[:55]
		}

		fmt.Printf("║  %s %s: %.0f%%\n", advisorIcon, strings.ToUpper(metric), score)
		fmt.Printf("║     Advisor: %s\n", advisor)
		fmt.Printf("║     \"%s\"\n", quotePreview)
		fmt.Printf("║     💡 %s\n", encPreview)
		fmt.Println()
	}

	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Note: Consultation log integration coming in Phase 5")

	return nil
}
