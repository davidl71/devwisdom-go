package wisdom

import (
	"testing"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Fatal("NewEngine returned nil")
	}
	if engine.sources == nil {
		t.Error("Engine sources map is nil")
	}
	if engine.loader == nil {
		t.Error("Engine loader is nil")
	}
	if engine.advisors == nil {
		t.Error("Engine advisors is nil")
	}
}

func TestEngine_Initialize(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	if !engine.initialized {
		t.Error("Engine not marked as initialized")
	}
}

func TestEngine_Initialize_Twice(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("First Initialize failed: %v", err)
	}
	// Second initialize should not error
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Second Initialize failed: %v", err)
	}
}

func TestEngine_GetWisdom_NotInitialized(t *testing.T) {
	engine := NewEngine()
	_, err := engine.GetWisdom(50.0, "test")
	if err == nil {
		t.Error("GetWisdom should fail when engine not initialized")
	}
}

func TestEngine_GetWisdom_UnknownSource(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	_, err := engine.GetWisdom(50.0, "nonexistent")
	if err == nil {
		t.Error("GetWisdom should fail for unknown source")
	}
}

func TestEngine_GetWisdom_Success(t *testing.T) {
	engine := NewEngine()
	
	// Add a test source directly
	testSource := &Source{
		Name: "Test Source",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"middle_aeons": {
				{Quote: "Test quote", Source: "Test", Encouragement: "Test"},
			},
		},
	}
	engine.sources["test"] = testSource
	engine.initialized = true

	quote, err := engine.GetWisdom(60.0, "test")
	if err != nil {
		t.Fatalf("GetWisdom failed: %v", err)
	}
	if quote == nil {
		t.Fatal("GetWisdom returned nil quote")
	}
	if quote.Quote != "Test quote" {
		t.Errorf("GetWisdom quote = %q, want %q", quote.Quote, "Test quote")
	}
}

func TestEngine_ListSources(t *testing.T) {
	engine := NewEngine()
	
	// Add sources directly
	engine.sources["source1"] = &Source{Name: "Source 1", Icon: "📜", Quotes: make(map[string][]Quote)}
	engine.sources["source2"] = &Source{Name: "Source 2", Icon: "📜", Quotes: make(map[string][]Quote)}
	engine.loader = nil // Disable loader to use internal sources
	engine.initialized = true

	sources := engine.ListSources()
	if len(sources) != 2 {
		t.Errorf("ListSources returned %d sources, want 2", len(sources))
	}
}

func TestEngine_GetSource(t *testing.T) {
	engine := NewEngine()
	testSource := &Source{
		Name: "Test Source",
		Icon: "📜",
		Quotes: make(map[string][]Quote),
	}
	engine.sources["test"] = testSource
	engine.loader = nil // Disable loader to use internal sources
	engine.initialized = true

	source, found := engine.GetSource("test")
	if !found {
		t.Fatal("GetSource returned false for existing source")
	}
	if source.Name != "Test Source" {
		t.Errorf("GetSource name = %q, want %q", source.Name, "Test Source")
	}

	_, found = engine.GetSource("nonexistent")
	if found {
		t.Error("GetSource returned true for nonexistent source")
	}
}

func TestEngine_ReloadSources(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Reload should not error even if no sources loaded
	if err := engine.ReloadSources(); err != nil {
		t.Fatalf("ReloadSources failed: %v", err)
	}
}

func TestEngine_GetLoader(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	loader := engine.GetLoader()
	if loader == nil {
		t.Error("GetLoader returned nil")
	}
}

func TestEngine_GetWisdom_AeonLevels(t *testing.T) {
	engine := NewEngine()
	
	// Create source with quotes for all aeon levels
	testSource := &Source{
		Name: "Test Source",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos":        {{Quote: "Chaos quote", Source: "Test", Encouragement: "Test"}},
			"lower_aeons":  {{Quote: "Lower quote", Source: "Test", Encouragement: "Test"}},
			"middle_aeons": {{Quote: "Middle quote", Source: "Test", Encouragement: "Test"}},
			"upper_aeons":  {{Quote: "Upper quote", Source: "Test", Encouragement: "Test"}},
			"treasury":     {{Quote: "Treasury quote", Source: "Test", Encouragement: "Test"}},
		},
	}
	engine.sources["test"] = testSource
	engine.initialized = true

	tests := []struct {
		score    float64
		expected string
	}{
		{10.0, "Chaos quote"},
		{40.0, "Lower quote"},
		{60.0, "Middle quote"},
		{80.0, "Upper quote"},
		{90.0, "Treasury quote"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			quote, err := engine.GetWisdom(tt.score, "test")
			if err != nil {
				t.Fatalf("GetWisdom failed: %v", err)
			}
			if quote.Quote != tt.expected {
				t.Errorf("GetWisdom(%.1f) quote = %q, want %q", tt.score, quote.Quote, tt.expected)
			}
		})
	}
}

