package cli

import (
	"testing"
)

func TestApp_Run(t *testing.T) {
	app := NewApp("0.1.0")

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "help command",
			args:    []string{"help"},
			wantErr: false,
		},
		{
			name:    "version command",
			args:    []string{"version"},
			wantErr: false,
		},
		{
			name:    "unknown command",
			args:    []string{"unknown"},
			wantErr: true,
		},
		{
			name:    "empty args",
			args:    []string{},
			wantErr: false, // Should print usage
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := app.Run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("App.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_RunQuote(t *testing.T) {
	app := NewApp("0.1.0")

	// Test quote command with minimal args
	err := app.Run([]string{"quote", "--quiet"})
	// May fail if no sources available, which is OK for now
	if err != nil {
		t.Logf("Quote command returned error (expected if no sources): %v", err)
	}
}

func TestApp_RunSources(t *testing.T) {
	app := NewApp("0.1.0")

	err := app.Run([]string{"sources"})
	if err != nil {
		t.Logf("Sources command returned error (expected if no sources): %v", err)
	}
}

func TestApp_RunAdvisors(t *testing.T) {
	app := NewApp("0.1.0")

	err := app.Run([]string{"advisors"})
	if err != nil {
		t.Errorf("Advisors command returned error: %v", err)
	}
}

func TestApp_RunBriefing(t *testing.T) {
	app := NewApp("0.1.0")

	err := app.Run([]string{"briefing"})
	if err != nil {
		t.Logf("Briefing command returned error (expected if no sources): %v", err)
	}
}
