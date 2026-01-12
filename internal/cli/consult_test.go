package cli

import (
	"bytes"
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
			check:   validateJSONOutput,
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
				_, err := buf.ReadFrom(r)
				if err != nil {
					t.Errorf("buf.ReadFrom failed: %v", err)
				}
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
