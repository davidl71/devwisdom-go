package wisdom

import "testing"

func TestNewAdvisorRegistry(t *testing.T) {
	registry := NewAdvisorRegistry()
	if registry == nil {
		t.Fatal("NewAdvisorRegistry returned nil")
	}
	if registry.metricAdvisors == nil {
		t.Error("metricAdvisors map is nil")
	}
	if registry.toolAdvisors == nil {
		t.Error("toolAdvisors map is nil")
	}
	if registry.stageAdvisors == nil {
		t.Error("stageAdvisors map is nil")
	}
}

func TestAdvisorRegistry_Initialize(t *testing.T) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	if !registry.initialized {
		t.Error("Registry not marked as initialized")
	}
}

func TestAdvisorRegistry_Initialize_Twice(t *testing.T) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	// Second initialize should not error
	registry.Initialize()
	if !registry.initialized {
		t.Error("Registry not marked as initialized after second call")
	}
}

func TestAdvisorRegistry_GetAdvisorForMetric(t *testing.T) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	tests := []struct {
		metric  string
		wantErr bool
		checkID string
	}{
		{"security", false, "bofh"},
		{"testing", false, "stoic"},
		{"nonexistent", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.metric, func(t *testing.T) {
			advisor, err := registry.GetAdvisorForMetric(tt.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAdvisorForMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if advisor == nil {
					t.Fatal("GetAdvisorForMetric returned nil advisor")
				}
				if advisor.Advisor != tt.checkID {
					t.Errorf("GetAdvisorForMetric() advisor = %q, want %q", advisor.Advisor, tt.checkID)
				}
			}
		})
	}
}

func TestAdvisorRegistry_GetAdvisorForTool(t *testing.T) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	tests := []struct {
		tool    string
		wantErr bool
		checkID string
	}{
		{"project_scorecard", false, "pistis_sophia"},
		{"nonexistent", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			advisor, err := registry.GetAdvisorForTool(tt.tool)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAdvisorForTool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if advisor == nil {
					t.Fatal("GetAdvisorForTool returned nil advisor")
				}
				if advisor.Advisor != tt.checkID {
					t.Errorf("GetAdvisorForTool() advisor = %q, want %q", advisor.Advisor, tt.checkID)
				}
			}
		})
	}
}

func TestAdvisorRegistry_GetAdvisorForStage(t *testing.T) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	tests := []struct {
		stage   string
		wantErr bool
		checkID string
	}{
		{"daily_checkin", false, "pistis_sophia"},
		{"nonexistent", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.stage, func(t *testing.T) {
			advisor, err := registry.GetAdvisorForStage(tt.stage)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAdvisorForStage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if advisor == nil {
					t.Fatal("GetAdvisorForStage returned nil advisor")
				}
				if advisor.Advisor != tt.checkID {
					t.Errorf("GetAdvisorForStage() advisor = %q, want %q", advisor.Advisor, tt.checkID)
				}
			}
		})
	}
}
