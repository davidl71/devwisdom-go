package cli

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestRunAdvisors(t *testing.T) {
	app := NewApp("test")

	// Test JSON output
	t.Run("JSON output", func(t *testing.T) {
		// Capture output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := app.runAdvisors([]string{"--json"})
		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("runAdvisors() error = %v", err)
		}

		// Read output
		var output map[string]interface{}
		decoder := json.NewDecoder(r)
		if err := decoder.Decode(&output); err != nil {
			t.Fatalf("Failed to decode JSON output: %v", err)
		}

		// Verify structure
		if _, ok := output["metric_advisors"]; !ok {
			t.Error("JSON output missing metric_advisors")
		}
		if _, ok := output["tool_advisors"]; !ok {
			t.Error("JSON output missing tool_advisors")
		}
		if _, ok := output["stage_advisors"]; !ok {
			t.Error("JSON output missing stage_advisors")
		}

		// Verify metrics
		metrics, ok := output["metric_advisors"].([]interface{})
		if !ok {
			t.Fatal("metric_advisors is not an array")
		}
		if len(metrics) < 10 {
			t.Errorf("Expected at least 10 metric advisors, got %d", len(metrics))
		}

		// Verify first metric has required fields
		if len(metrics) > 0 {
			firstMetric, ok := metrics[0].(map[string]interface{})
			if !ok {
				t.Fatal("First metric is not an object")
			}
			if _, ok := firstMetric["metric"]; !ok {
				t.Error("Metric missing 'metric' field")
			}
			if _, ok := firstMetric["advisor"]; !ok {
				t.Error("Metric missing 'advisor' field")
			}
		}
	})

	// Test human-readable output contains expected sections
	t.Run("Human-readable output", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := app.runAdvisors([]string{})
		w.Close()
		os.Stdout = oldStdout

		if err != nil {
			t.Fatalf("runAdvisors() error = %v", err)
		}

		// Read output
		output := make([]byte, 4096)
		n, _ := r.Read(output)
		outputStr := string(output[:n])

		// Verify sections exist
		if !strings.Contains(outputStr, "Metric Advisors") {
			t.Error("Output missing 'Metric Advisors' section")
		}
		if !strings.Contains(outputStr, "Tool Advisors") {
			t.Error("Output missing 'Tool Advisors' section")
		}
		if !strings.Contains(outputStr, "Stage Advisors") {
			t.Error("Output missing 'Stage Advisors' section")
		}

		// Verify some known advisors are mentioned
		if !strings.Contains(outputStr, "security") && !strings.Contains(outputStr, "bofh") {
			t.Error("Output missing known metric advisor (security/bofh)")
		}
		if !strings.Contains(outputStr, "project_scorecard") && !strings.Contains(outputStr, "pistis_sophia") {
			t.Error("Output missing known tool advisor")
		}
		if !strings.Contains(outputStr, "daily_checkin") && !strings.Contains(outputStr, "pistis_sophia") {
			t.Error("Output missing known stage advisor")
		}
	})
}

func TestRunAdvisors_EmptyLists(t *testing.T) {
	// This test would require mocking the engine, which is complex
	// For now, we assume engine initialization works (tested elsewhere)
	// This is more of an integration test
	t.Skip("Requires engine mocking - integration test")
}
