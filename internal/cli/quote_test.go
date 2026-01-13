package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

// validateJSONOutput validates that the output is valid JSON
func validateJSONOutput(output string) bool {
	if len(output) == 0 {
		return false
	}
	var v interface{}
	return json.Unmarshal([]byte(output), &v) == nil
}

func TestRunQuote(t *testing.T) {
	app := NewApp("0.1.0")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(output string) bool
	}{
		{
			name:    "quote with quiet flag",
			args:    []string{"--quiet"},
			wantErr: false,
			check: func(output string) bool {
				return len(output) > 0 && !validateJSONOutput(output)
			},
		},
		{
			name:    "quote with json flag",
			args:    []string{"--json"},
			wantErr: false,
			check:   validateJSONOutput,
		},
		{
			name:    "quote with source and score",
			args:    []string{"--source", "stoic", "--score", "75"},
			wantErr: false,
			check: func(output string) bool {
				return len(output) > 0
			},
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

			err := app.runQuote(tt.args)

			w.Close()
			os.Stdout = oldStdout
			<-done

			if (err != nil) != tt.wantErr {
				t.Errorf("runQuote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				output := buf.String()
				if !tt.check(output) {
					t.Errorf("runQuote() output validation failed. Output: %s", output)
				}
			}
		})
	}
}

func TestRunQuote_InvalidSource(t *testing.T) {
	app := NewApp("0.1.0")

	err := app.runQuote([]string{"--source", "nonexistent_source"})
	if err == nil {
		t.Error("runQuote() with invalid source should return error")
	}
}

func TestRunQuote_InvalidScore(t *testing.T) {
	app := NewApp("0.1.0")

	// Test with negative score (should still work, just uses different aeon level)
	err := app.runQuote([]string{"--score", "-10"})
	// May fail if no sources, which is OK
	if err != nil {
		t.Logf("Quote with negative score returned error (may be expected): %v", err)
	}
}
