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
		name           string
		source         *Source
		aeonLevel      string
		expectedQuote  string
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
			aeonLevel:      "chaos",
			expectedQuote:  "Test quote",
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
			aeonLevel:      "chaos",
			expectedQuote:  "Fallback quote",
			expectFallback: true,
		},
		{
			name: "empty quotes - default fallback",
			source: &Source{
				Name:   "Test Source",
				Icon:   "📜",
				Quotes: map[string][]Quote{},
			},
			aeonLevel:      "chaos",
			expectedQuote:  "Silence is also wisdom.",
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
			aeonLevel:      "chaos",
			expectedQuote:  "Any quote",
			expectFallback: true,
		},
		{
			name: "multiple quotes - random selection",
			source: &Source{
				Name: "Test Source",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"chaos": {
						{Quote: "Quote 1", Source: "Test", Encouragement: "Test"},
						{Quote: "Quote 2", Source: "Test", Encouragement: "Test"},
						{Quote: "Quote 3", Source: "Test", Encouragement: "Test"},
					},
				},
			},
			aeonLevel:      "chaos",
			expectedQuote:  "", // Will check that it's one of the quotes
			expectFallback: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quote := tt.source.GetQuote(tt.aeonLevel)
			if quote == nil {
				t.Fatal("GetQuote returned nil")
			}
			if tt.expectedQuote == "" {
				// For random selection test, verify it's one of the expected quotes
				validQuotes := []string{"Quote 1", "Quote 2", "Quote 3"}
				found := false
				for _, vq := range validQuotes {
					if quote.Quote == vq {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetQuote() = %q, expected one of %v", quote.Quote, validQuotes)
				}
			} else if quote.Quote != tt.expectedQuote {
				t.Errorf("GetQuote() = %q, want %q", quote.Quote, tt.expectedQuote)
			}
		})
	}
}

func TestSource_GetQuote_RandomSelection(t *testing.T) {
	// Test that random selection works and is date-seeded (consistent within same day)
	source := &Source{
		Name: "Test Source",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {
				{Quote: "Quote 1", Source: "Test", Encouragement: "Test"},
				{Quote: "Quote 2", Source: "Test", Encouragement: "Test"},
				{Quote: "Quote 3", Source: "Test", Encouragement: "Test"},
			},
		},
	}

	// Get quote multiple times - should return same quote (date-seeded)
	quote1 := source.GetQuote("chaos")
	quote2 := source.GetQuote("chaos")
	quote3 := source.GetQuote("chaos")

	if quote1 == nil || quote2 == nil || quote3 == nil {
		t.Fatal("GetQuote returned nil")
	}

	// All quotes should be the same (date-seeded consistency)
	if quote1.Quote != quote2.Quote || quote2.Quote != quote3.Quote {
		t.Errorf("Date-seeded random selection not consistent: got %q, %q, %q", quote1.Quote, quote2.Quote, quote3.Quote)
	}

	// Verify it's one of the valid quotes
	validQuotes := []string{"Quote 1", "Quote 2", "Quote 3"}
	found := false
	for _, vq := range validQuotes {
		if quote1.Quote == vq {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("GetQuote() = %q, expected one of %v", quote1.Quote, validQuotes)
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
