package wisdom

import (
	"fmt"
	"testing"
)

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

func TestGetConsultationMode(t *testing.T) {
	tests := []struct {
		score    float64
		wantMode string
		wantIcon string
		wantFreq string
		boundary bool
	}{
		{0, "chaos", "🔥", "every_action", true},
		{15, "chaos", "🔥", "every_action", false},
		{29.9, "chaos", "🔥", "every_action", true},
		{30, "building", "🏗️", "start_and_review", true},
		{45, "building", "🏗️", "start_and_review", false},
		{59.9, "building", "🏗️", "start_and_review", true},
		{60, "maturing", "🌱", "milestones", true},
		{70, "maturing", "🌱", "milestones", false},
		{79.9, "maturing", "🌱", "milestones", true},
		{80, "mastery", "🎯", "weekly", true},
		{90, "mastery", "🎯", "weekly", false},
		{100, "mastery", "🎯", "weekly", true},
		{150, "mastery", "🎯", "weekly", false},    // Edge case: score > 100
		{-10, "chaos", "🔥", "every_action", true}, // Edge case: negative score (should default to chaos)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("score_%.1f", tt.score), func(t *testing.T) {
			mode := GetConsultationMode(tt.score)
			if mode.Name != tt.wantMode {
				t.Errorf("GetConsultationMode(%.1f) mode = %q, want %q", tt.score, mode.Name, tt.wantMode)
			}
			if mode.Icon != tt.wantIcon {
				t.Errorf("GetConsultationMode(%.1f) icon = %q, want %q", tt.score, mode.Icon, tt.wantIcon)
			}
			if mode.Frequency != tt.wantFreq {
				t.Errorf("GetConsultationMode(%.1f) frequency = %q, want %q", tt.score, mode.Frequency, tt.wantFreq)
			}
			if mode.MinScore > tt.score || tt.score >= mode.MaxScore {
				if !tt.boundary && tt.score >= 0 && tt.score <= 100 {
					t.Errorf("GetConsultationMode(%.1f) score range invalid: %.1f not in [%.1f, %.1f)", tt.score, tt.score, mode.MinScore, mode.MaxScore)
				}
			}
		})
	}
}

func TestGetModeConfig(t *testing.T) {
	tests := []struct {
		mode         SessionMode
		wantAdvisors int
		wantTone     string
	}{
		{SessionModeAgent, 3, "strategic"},
		{SessionModeAsk, 3, "direct"},
		{SessionModeManual, 3, "observational"},
		{SessionMode("UNKNOWN"), 0, ""}, // Should return nil
	}

	for _, tt := range tests {
		t.Run(string(tt.mode), func(t *testing.T) {
			config := GetModeConfig(tt.mode)
			if tt.wantAdvisors == 0 {
				if config != nil {
					t.Errorf("GetModeConfig(%q) = %v, want nil", tt.mode, config)
				}
				return
			}
			if config == nil {
				t.Fatalf("GetModeConfig(%q) = nil, want config", tt.mode)
			}
			if len(config.PreferredAdvisors) != tt.wantAdvisors {
				t.Errorf("GetModeConfig(%q) advisors count = %d, want %d", tt.mode, len(config.PreferredAdvisors), tt.wantAdvisors)
			}
			if config.Tone != tt.wantTone {
				t.Errorf("GetModeConfig(%q) tone = %q, want %q", tt.mode, config.Tone, tt.wantTone)
			}
		})
	}
}

func TestAdjustAdvisorForMode(t *testing.T) {
	availableSources := []string{"art_of_war", "tao_of_programming", "confucius", "stoic", "bofh"}

	tests := []struct {
		name             string
		sessionMode      SessionMode
		consultationType string
		availableSources []string
		wantAdvisor      string
		wantRationale    string
		shouldAdjust     bool
	}{
		{
			name:             "random consultation with AGENT mode",
			sessionMode:      SessionModeAgent,
			consultationType: "random",
			availableSources: availableSources,
			wantAdvisor:      "art_of_war", // First available preferred
			wantRationale:    "Mode-aware selection for AGENT",
			shouldAdjust:     true,
		},
		{
			name:             "random consultation with ASK mode",
			sessionMode:      SessionModeAsk,
			consultationType: "random",
			availableSources: availableSources,
			wantAdvisor:      "confucius", // First available preferred
			wantRationale:    "Mode-aware selection for ASK",
			shouldAdjust:     true,
		},
		{
			name:             "metric consultation should not adjust",
			sessionMode:      SessionModeAgent,
			consultationType: "metric",
			availableSources: availableSources,
			wantAdvisor:      "",
			wantRationale:    "",
			shouldAdjust:     false,
		},
		{
			name:             "no available preferred advisors",
			sessionMode:      SessionModeManual,
			consultationType: "random",
			availableSources: []string{"bofh", "stoic"}, // No preferred advisors available
			wantAdvisor:      "",
			wantRationale:    "",
			shouldAdjust:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			advisor, rationale := AdjustAdvisorForMode(tt.sessionMode, tt.consultationType, tt.availableSources)
			if tt.shouldAdjust {
				if advisor == "" {
					t.Errorf("AdjustAdvisorForMode() advisor = %q, want non-empty", advisor)
				}
				if rationale == "" {
					t.Errorf("AdjustAdvisorForMode() rationale = %q, want non-empty", rationale)
				}
				if advisor != tt.wantAdvisor {
					t.Errorf("AdjustAdvisorForMode() advisor = %q, want %q", advisor, tt.wantAdvisor)
				}
				if rationale != tt.wantRationale {
					t.Errorf("AdjustAdvisorForMode() rationale = %q, want %q", rationale, tt.wantRationale)
				}
			} else {
				if advisor != "" || rationale != "" {
					t.Errorf("AdjustAdvisorForMode() should not adjust, got advisor=%q rationale=%q", advisor, rationale)
				}
			}
		})
	}
}

func TestAdvisorRegistry_GetAllMetricAdvisors(t *testing.T) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	advisors := registry.GetAllMetricAdvisors()
	if advisors == nil {
		t.Fatal("GetAllMetricAdvisors returned nil")
	}
	if len(advisors) == 0 {
		t.Error("GetAllMetricAdvisors returned empty map")
	}

	// Check that known metrics exist
	if _, exists := advisors["security"]; !exists {
		t.Error("GetAllMetricAdvisors missing security advisor")
	}
	if _, exists := advisors["testing"]; !exists {
		t.Error("GetAllMetricAdvisors missing testing advisor")
	}
}

func TestAdvisorRegistry_GetAllToolAdvisors(t *testing.T) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	advisors := registry.GetAllToolAdvisors()
	if advisors == nil {
		t.Fatal("GetAllToolAdvisors returned nil")
	}
	if len(advisors) == 0 {
		t.Error("GetAllToolAdvisors returned empty map")
	}

	// Check that known tools exist
	if _, exists := advisors["project_scorecard"]; !exists {
		t.Error("GetAllToolAdvisors missing project_scorecard advisor")
	}
}

func TestAdvisorRegistry_GetAllStageAdvisors(t *testing.T) {
	registry := NewAdvisorRegistry()
	registry.Initialize()

	advisors := registry.GetAllStageAdvisors()
	if advisors == nil {
		t.Fatal("GetAllStageAdvisors returned nil")
	}
	if len(advisors) == 0 {
		t.Error("GetAllStageAdvisors returned empty map")
	}

	// Check that known stages exist
	if _, exists := advisors["daily_checkin"]; !exists {
		t.Error("GetAllStageAdvisors missing daily_checkin advisor")
	}
}
