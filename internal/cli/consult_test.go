package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestRunConsult(t *testing.T) {
	app := NewApp("0.1.0")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(output string) bool
	}{
		{
			name:    "consult with metric",
			args:    []string{"--metric", "security", "--score", "40"},
			wantErr: false,
			check: func(output string) bool {
				return strings.Contains(output, "Advisor") || strings.Contains(output, "bofh")
			},
		},
		{
			name:    "consult with tool",
			args:    []string{"--tool", "project_scorecard", "--score", "60"},
			wantErr: false,
			check: func(output string) bool {
				return len(output) > 0
			},
		},
		{
			name:    "consult with stage",
			args:    []string{"--stage", "daily_checkin", "--score", "50"},
			wantErr: false,
			check: func(output string) bool {
				return len(output) > 0
			},
		},
		{
			name:    "consult with json flag",
			args:    []string{"--metric", "security", "--json"},
			wantErr: false,
			check: func(output string) bool {
				// Extract JSON from output (may have warnings before it)
				lines := strings.Split(output, "\n")
				var jsonLines []string
				inJSON := false
				braceCount := 0
				bracketCount := 0

				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if trimmed == "" {
						if inJSON {
							jsonLines = append(jsonLines, "")
						}
						continue
					}
					// Skip warning lines
					if strings.HasPrefix(trimmed, "Warning:") {
						continue
					}
					// Start collecting JSON when we see { or [
					if !inJSON && (trimmed[0] == '{' || trimmed[0] == '[') {
						inJSON = true
					}

					if inJSON {
						jsonLines = append(jsonLines, trimmed)
						// Count braces and brackets to know when JSON ends
						for _, char := range trimmed {
							if char == '{' {
								braceCount++
							} else if char == '}' {
								braceCount--
							} else if char == '[' {
								bracketCount++
							} else if char == ']' {
								bracketCount--
							}
						}
						// JSON is complete when all braces/brackets are closed
						if braceCount == 0 && bracketCount == 0 && len(jsonLines) > 0 {
							break
						}
					}
				}
				if len(jsonLines) > 0 {
					jsonStr := strings.Join(jsonLines, "\n")
					return json.Valid([]byte(jsonStr))
				}
				// Try validating entire output as fallback
				return json.Valid([]byte(output))
			},
		},
		{
			name:    "consult with quiet flag",
			args:    []string{"--metric", "testing", "--quiet"},
			wantErr: false,
			check: func(output string) bool {
				// Quiet mode should only output quote text
				return len(output) > 0 && !strings.Contains(output, "Advisor")
			},
		},
		{
			name:    "consult without metric/tool/stage",
			args:    []string{"--score", "50"},
			wantErr: true, // Should fail - need at least one
		},
		{
			name:    "consult with invalid metric",
			args:    []string{"--metric", "nonexistent"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture output using a pipe
			r, w, _ := os.Pipe()
			oldStdout := os.Stdout
			os.Stdout = w

			var buf bytes.Buffer
			done := make(chan bool)
			go func() {
				buf.ReadFrom(r)
				done <- true
			}()

			err := app.runConsult(tt.args)

			w.Close()
			os.Stdout = oldStdout
			<-done

			if (err != nil) != tt.wantErr {
				t.Errorf("runConsult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				output := buf.String()
				if !tt.check(output) {
					t.Errorf("runConsult() output validation failed. Output: %s", output)
				}
			}
		})
	}
}
