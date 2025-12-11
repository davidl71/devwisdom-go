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

	// Build comprehensive advisor listing
	advisorList := map[string]interface{}{
		"metric_advisors": make([]map[string]interface{}, 0),
		"tool_advisors":   make([]map[string]interface{}, 0),
		"stage_advisors":  make([]map[string]interface{}, 0),
	}

	// Get all metric advisors
	metricAdvisors := advisors.GetAllMetricAdvisors()
	metricList := make([]map[string]interface{}, 0, len(metricAdvisors))
	metricKeys := make([]string, 0, len(metricAdvisors))
	for k := range metricAdvisors {
		metricKeys = append(metricKeys, k)
	}
	sort.Strings(metricKeys)
	for _, metric := range metricKeys {
		info := metricAdvisors[metric]
		item := map[string]interface{}{
			"metric":    metric,
			"advisor":   info.Advisor,
			"icon":      info.Icon,
			"rationale": info.Rationale,
		}
		if info.HelpsWith != "" {
			item["helps_with"] = info.HelpsWith
		}
		if info.Language != "" {
			item["language"] = info.Language
		}
		metricList = append(metricList, item)
	}
	advisorList["metric_advisors"] = metricList

	// Get all tool advisors
	toolAdvisors := advisors.GetAllToolAdvisors()
	toolList := make([]map[string]interface{}, 0, len(toolAdvisors))
	toolKeys := make([]string, 0, len(toolAdvisors))
	for k := range toolAdvisors {
		toolKeys = append(toolKeys, k)
	}
	sort.Strings(toolKeys)
	for _, tool := range toolKeys {
		info := toolAdvisors[tool]
		item := map[string]interface{}{
			"tool":      tool,
			"advisor":   info.Advisor,
			"rationale": info.Rationale,
		}
		if info.Language != "" {
			item["language"] = info.Language
		}
		toolList = append(toolList, item)
	}
	advisorList["tool_advisors"] = toolList

	// Get all stage advisors
	stageAdvisors := advisors.GetAllStageAdvisors()
	stageList := make([]map[string]interface{}, 0, len(stageAdvisors))
	stageKeys := make([]string, 0, len(stageAdvisors))
	for k := range stageAdvisors {
		stageKeys = append(stageKeys, k)
	}
	sort.Strings(stageKeys)
	for _, stage := range stageKeys {
		info := stageAdvisors[stage]
		item := map[string]interface{}{
			"stage":     stage,
			"advisor":   info.Advisor,
			"icon":      info.Icon,
			"rationale": info.Rationale,
		}
		if info.Language != "" {
			item["language"] = info.Language
		}
		stageList = append(stageList, item)
	}
	advisorList["stage_advisors"] = stageList

	// Output
	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(advisorList)
	}

	// Human-readable output
	fmt.Println("Available Advisor Mappings")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	// Metric Advisors
	fmt.Printf("📊 Metric Advisors (%d):\n\n", len(metricList))
	for _, item := range metricList {
		metric := item["metric"].(string)
		advisor := item["advisor"].(string)
		icon := item["icon"].(string)
		rationale := item["rationale"].(string)

		fmt.Printf("  %s %s → %s\n", icon, metric, advisor)
		if rationale != "" {
			fmt.Printf("    %s\n", rationale)
		}
		if helpsWith, ok := item["helps_with"].(string); ok && helpsWith != "" {
			fmt.Printf("    Helps with: %s\n", helpsWith)
		}
		if lang, ok := item["language"].(string); ok && lang != "" {
			fmt.Printf("    Language: %s\n", lang)
		}
		fmt.Println()
	}

	// Tool Advisors
	fmt.Printf("🔧 Tool Advisors (%d):\n\n", len(toolList))
	for _, item := range toolList {
		tool := item["tool"].(string)
		advisor := item["advisor"].(string)
		rationale := item["rationale"].(string)

		fmt.Printf("  %s → %s\n", tool, advisor)
		if rationale != "" {
			fmt.Printf("    %s\n", rationale)
		}
		if lang, ok := item["language"].(string); ok && lang != "" {
			fmt.Printf("    Language: %s\n", lang)
		}
		fmt.Println()
	}

	// Stage Advisors
	fmt.Printf("🎭 Stage Advisors (%d):\n\n", len(stageList))
	for _, item := range stageList {
		stage := item["stage"].(string)
		advisor := item["advisor"].(string)
		icon := item["icon"].(string)
		rationale := item["rationale"].(string)

		fmt.Printf("  %s %s → %s\n", icon, stage, advisor)
		if rationale != "" {
			fmt.Printf("    %s\n", rationale)
		}
		if lang, ok := item["language"].(string); ok && lang != "" {
			fmt.Printf("    Language: %s\n", lang)
		}
		fmt.Println()
	}

	fmt.Println("Use 'devwisdom consult' to get advisor guidance:")
	fmt.Println("  devwisdom consult --metric security")
	fmt.Println("  devwisdom consult --tool project_scorecard")
	fmt.Println("  devwisdom consult --stage daily_checkin")

	return nil
}
