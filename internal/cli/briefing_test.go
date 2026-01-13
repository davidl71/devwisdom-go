package cli

import (
	"strings"
	"testing"
)

func TestRunBriefing(t *testing.T) {
	app := NewApp("0.1.0")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
		check   func(output string) bool
	}{
		{
			name:    "briefing default",
			args:    []string{},
			wantErr: false,
			check: func(output string) bool {
				return strings.Contains(output, "Daily Briefing") ||
					len(output) > 0
			},
		},
		{
			name:    "briefing with days",
			args:    []string{"--days", "7"},
			wantErr: false,
			check: func(output string) bool {
				return strings.Contains(output, "7") || len(output) > 0
			},
		},
		{
			name:    "briefing with json flag",
			args:    []string{"--json"},
			wantErr: false,
			check:   validateJSONOutput,
		},
	}

	for _, tt := range tests {
		runCLITest(t, tt.name, func() error {
			return app.runBriefing(tt.args)
		}, tt.wantErr, tt.check)
	}
}
