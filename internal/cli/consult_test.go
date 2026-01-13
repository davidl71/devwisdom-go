package cli

import (
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
		runCLITest(t, tt.name, func() error {
			return app.runConsult(tt.args)
		}, tt.wantErr, tt.check)
	}
}
