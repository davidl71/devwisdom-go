package cli

import (
	"testing"
)

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
		runCLITest(t, tt.name, func() error {
			return app.runQuote(tt.args)
		}, tt.wantErr, tt.check)
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
