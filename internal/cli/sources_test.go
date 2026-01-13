package cli

import (
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
			check:   validateJSONOutput,
		},
	}

	for _, tt := range tests {
		runCLITest(t, tt.name, func() error {
			return app.runSources(tt.args)
		}, tt.wantErr, tt.check)
	}
}
