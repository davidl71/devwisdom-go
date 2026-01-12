package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunSources(t *testing.T) {
	app := NewApp("0.1.0")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(output string) bool
	}{
		{
			name:    "sources default output",
			args:    []string{},
			wantErr: false,
			check: func(output string) bool {
				return strings.Contains(output, "Available Wisdom Sources") ||
					strings.Contains(output, "sources available")
			},
		},
		{
			name:    "sources with json flag",
			args:    []string{"--json"},
			wantErr: false,
			check: validateJSONOutput,
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

			err := app.runSources(tt.args)

			w.Close()
			os.Stdout = oldStdout
			<-done

			if (err != nil) != tt.wantErr {
				t.Errorf("runSources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				output := buf.String()
				if !tt.check(output) {
					t.Errorf("runSources() output validation failed. Output: %s", output)
				}
			}
		})
	}
}
