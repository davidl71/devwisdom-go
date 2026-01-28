package wisdom

import (
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine() returned nil")
	}
}

func TestGetAeonLevel(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		score    float64
	}{
		{"negative score", -10, "chaos"},
		{"zero score", 0, "chaos"},
		{"low score", 15, "chaos"},
		{"boundary 30", 30, "lower_aeons"},
		{"middle low", 40, "lower_aeons"},
		{"boundary 50", 50, "middle_aeons"},
		{"middle", 60, "middle_aeons"},
		{"boundary 70", 70, "upper_aeons"},
		{"high", 80, "upper_aeons"},
		{"boundary 85", 85, "treasury"},
		{"very high", 95, "treasury"},
		{"over 100", 150, "treasury"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAeonLevel(tt.score)
			if result != tt.expected {
				t.Errorf("GetAeonLevel(%v) = %q, want %q", tt.score, result, tt.expected)
			}
		})
	}
}

func TestGetConsultationMode(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		score    float64
	}{
		{"negative score", -10, "chaos"},
		{"zero score", 0, "chaos"},
		{"low score", 15, "chaos"},
		{"boundary 30", 30, "building"},
		{"middle low", 40, "building"},
		{"boundary 60", 60, "maturing"},
		{"middle", 70, "maturing"},
		{"boundary 80", 80, "mastery"},
		{"high", 95, "mastery"},
		{"over 100", 150, "mastery"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetConsultationMode(tt.score)
			if result.Name != tt.expected {
				t.Errorf("GetConsultationMode(%v).Name = %q, want %q", tt.score, result.Name, tt.expected)
			}
		})
	}
}

func TestGetConsultationMode_Config(t *testing.T) {
	// Test that GetConsultationMode returns valid configs.
	scores := []float64{-10, 0, 15, 30, 40, 50, 60, 70, 80, 90, 100, 150}

	for _, score := range scores {
		t.Run("score_"+string(rune(score)), func(t *testing.T) {
			config := GetConsultationMode(score)
			if config.Name == "" {
				t.Errorf("GetConsultationMode(%v) returned config with empty Name", score)
			}
			if config.MinScore < 0 || config.MinScore > 100 {
				t.Errorf("GetConsultationMode(%v) returned config with invalid MinScore: %v", score, config.MinScore)
			}
			if config.MaxScore < 0 || config.MaxScore > 100 {
				t.Errorf("GetConsultationMode(%v) returned config with invalid MaxScore: %v", score, config.MaxScore)
			}
			if config.MinScore > config.MaxScore {
				t.Errorf("GetConsultationMode(%v) returned config with MinScore > MaxScore", score)
			}
		})
	}
}
