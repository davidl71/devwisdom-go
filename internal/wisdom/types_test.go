package wisdom

import "testing"

func TestGetAeonLevel(t *testing.T) {
	tests := []struct {
		name     string
		score    float64
		expected string
	}{
		{"chaos - zero", 0.0, "chaos"},
		{"chaos - low", 15.0, "chaos"},
		{"chaos - boundary", 29.9, "chaos"},
		{"lower_aeons - boundary", 30.0, "lower_aeons"},
		{"lower_aeons - middle", 40.0, "lower_aeons"},
		{"lower_aeons - upper", 49.9, "lower_aeons"},
		{"middle_aeons - boundary", 50.0, "middle_aeons"},
		{"middle_aeons - middle", 60.0, "middle_aeons"},
		{"middle_aeons - upper", 69.9, "middle_aeons"},
		{"upper_aeons - boundary", 70.0, "upper_aeons"},
		{"upper_aeons - middle", 77.5, "upper_aeons"},
		{"upper_aeons - upper", 84.9, "upper_aeons"},
		{"treasury - boundary", 85.0, "treasury"},
		{"treasury - high", 90.0, "treasury"},
		{"treasury - perfect", 100.0, "treasury"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAeonLevel(tt.score)
			if result != tt.expected {
				t.Errorf("GetAeonLevel(%.1f) = %s, want %s", tt.score, result, tt.expected)
			}
		})
	}
}

func TestSource_GetQuote(t *testing.T) {
	tests := []struct {
		name          string
		source        *Source
		aeonLevel     string
		expectedQuote string
		expectFallback bool
	}{
		{
			name: "valid aeon level",
			source: &Source{
				Name: "Test Source",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"chaos": {
						{Quote: "Test quote", Source: "Test", Encouragement: "Test"},
					},
				},
			},
			aeonLevel:     "chaos",
			expectedQuote: "Test quote",
			expectFallback: false,
		},
		{
			name: "missing aeon level - fallback",
			source: &Source{
				Name: "Test Source",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"treasury": {
						{Quote: "Fallback quote", Source: "Test", Encouragement: "Test"},
					},
				},
			},
			aeonLevel:     "chaos",
			expectedQuote: "Fallback quote",
			expectFallback: true,
		},
		{
			name: "empty quotes - default fallback",
			source: &Source{
				Name:  "Test Source",
				Icon:  "📜",
				Quotes: map[string][]Quote{},
			},
			aeonLevel:     "chaos",
			expectedQuote: "Silence is also wisdom.",
			expectFallback: true,
		},
		{
			name: "empty aeon level - fallback to any",
			source: &Source{
				Name: "Test Source",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"middle_aeons": {
						{Quote: "Any quote", Source: "Test", Encouragement: "Test"},
					},
				},
			},
			aeonLevel:     "chaos",
			expectedQuote: "Any quote",
			expectFallback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quote := tt.source.GetQuote(tt.aeonLevel)
			if quote == nil {
				t.Fatal("GetQuote returned nil")
			}
			if quote.Quote != tt.expectedQuote {
				t.Errorf("GetQuote() = %q, want %q", quote.Quote, tt.expectedQuote)
			}
		})
	}
}

func TestAeonLevelConstants(t *testing.T) {
	if string(AeonChaos) != "chaos" {
		t.Errorf("AeonChaos = %q, want %q", AeonChaos, "chaos")
	}
	if string(AeonLower) != "lower_aeons" {
		t.Errorf("AeonLower = %q, want %q", AeonLower, "lower_aeons")
	}
	if string(AeonMiddle) != "middle_aeons" {
		t.Errorf("AeonMiddle = %q, want %q", AeonMiddle, "middle_aeons")
	}
	if string(AeonUpper) != "upper_aeons" {
		t.Errorf("AeonUpper = %q, want %q", AeonUpper, "upper_aeons")
	}
	if string(AeonTreasury) != "treasury" {
		t.Errorf("AeonTreasury = %q, want %q", AeonTreasury, "treasury")
	}
}

