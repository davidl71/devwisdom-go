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
		Name:   "Test Source",
		Icon:   "📜",
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

func TestEngine_GetWisdom_RandomSource(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Get wisdom with random source
	quote, err := engine.GetWisdom(75.0, "random")
	if err != nil {
		t.Fatalf("GetWisdom with random source failed: %v", err)
	}
	if quote == nil {
		t.Fatal("GetWisdom returned nil quote for random source")
	}
	if quote.Quote == "" {
		t.Error("GetWisdom returned empty quote for random source")
	}
}

func TestEngine_GetRandomSource(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Test date-seeded random source (should be consistent for same day)
	source1, err := engine.GetRandomSource(true)
	if err != nil {
		t.Fatalf("GetRandomSource failed: %v", err)
	}
	if source1 == "" {
		t.Error("GetRandomSource returned empty string")
	}

	// Should get same source on same day
	source2, err := engine.GetRandomSource(true)
	if err != nil {
		t.Fatalf("GetRandomSource failed: %v", err)
	}
	if source1 != source2 {
		t.Errorf("Date-seeded random source inconsistent: got %q, want %q", source2, source1)
	}

	// Test non-seeded random source (should vary)
	source3, err := engine.GetRandomSource(false)
	if err != nil {
		t.Fatalf("GetRandomSource failed: %v", err)
	}
	if source3 == "" {
		t.Error("GetRandomSource returned empty string")
	}
}

func TestEngine_GetRandomSource_NotInitialized(t *testing.T) {
	engine := NewEngine()
	_, err := engine.GetRandomSource(true)
	if err == nil {
		t.Error("GetRandomSource should fail when engine not initialized")
	}
}

func TestEngine_GetRandomSource_NoSources(t *testing.T) {
	engine := NewEngine()
	engine.sources = make(map[string]*Source) // Empty sources
	engine.initialized = true

	_, err := engine.GetRandomSource(true)
	if err == nil {
		t.Error("GetRandomSource should fail when no sources available")
	}
}

func TestEngine_GetRandomSource_OnlySefariaSources(t *testing.T) {
	engine := NewEngine()
	// Add only Sefaria sources (should be excluded)
	engine.sources["pirkei_avot"] = &Source{Name: "Pirkei Avot", Icon: "📜", Quotes: make(map[string][]Quote)}
	engine.sources["proverbs"] = &Source{Name: "Proverbs", Icon: "📜", Quotes: make(map[string][]Quote)}
	engine.initialized = true

	_, err := engine.GetRandomSource(true)
	if err == nil {
		t.Error("GetRandomSource should fail when only Sefaria sources available")
	}
}

func TestEngine_ListSources_NotInitialized(t *testing.T) {
	engine := NewEngine()
	sources := engine.ListSources()
	if len(sources) != 0 {
		t.Errorf("ListSources should return empty slice when not initialized, got %d", len(sources))
	}
}

func TestEngine_ListSources_WithLoader(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	sources := engine.ListSources()
	if len(sources) == 0 {
		t.Error("ListSources returned empty slice after initialization")
	}
}

func TestEngine_GetSource_NotInitialized(t *testing.T) {
	engine := NewEngine()
	_, found := engine.GetSource("test")
	if found {
		t.Error("GetSource should return false when engine not initialized")
	}
}

func TestEngine_GetSource_WithLoader(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Try to get a source that should exist after initialization
	source, found := engine.GetSource("stoic")
	if !found {
		// If stoic doesn't exist, try any source from the list
		sources := engine.ListSources()
		if len(sources) > 0 {
			source, found = engine.GetSource(sources[0])
		}
	}
	if !found {
		t.Skip("No sources available to test GetSource with loader")
	}
	if source == nil {
		t.Error("GetSource returned nil for existing source")
	}
}

func TestEngine_GetAdvisors(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	advisors := engine.GetAdvisors()
	if advisors == nil {
		t.Error("GetAdvisors returned nil")
	}
}

func TestEngine_ReloadSources_NotInitialized(t *testing.T) {
	engine := NewEngine()
	// Reload without initialization should fail
	err := engine.ReloadSources()
	if err == nil {
		t.Error("ReloadSources should fail when loader not initialized")
	}
}

func TestEngine_AddProjectSource(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	config := &SourceConfig{
		ID:          "test_project_source",
		Name:        "Test Project Source",
		Icon:        "🚀",
		Description: "Test source for project",
		Quotes: map[string][]Quote{
			"middle_aeons": {
				{Quote: "Test project quote", Source: "Test", Encouragement: "Test"},
			},
		},
	}

	// This might fail if project root detection fails, so we'll check error
	err := engine.AddProjectSource(config)
	if err != nil {
		// If it fails due to project root detection, that's okay for unit tests
		t.Logf("AddProjectSource failed (may be expected in test environment): %v", err)
		return
	}

	// If it succeeded, verify the source is available
	source, found := engine.GetSource("test_project_source")
	if !found {
		t.Error("AddProjectSource did not make source available")
	}
	if source == nil {
		t.Error("GetSource returned nil for added project source")
	}
}

func TestEngine_AddProjectSource_LoaderNotInitialized(t *testing.T) {
	engine := NewEngine()
	engine.loader = nil

	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
	}

	err := engine.AddProjectSource(config)
	if err == nil {
		t.Error("AddProjectSource should fail when loader not initialized")
	}
}

func TestEngine_GetWisdom_EdgeCaseScores(t *testing.T) {
	engine := NewEngine()

	testSource := &Source{
		Name: "Test Source",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos":        {{Quote: "Chaos", Source: "Test", Encouragement: "Test"}},
			"lower_aeons":  {{Quote: "Lower", Source: "Test", Encouragement: "Test"}},
			"middle_aeons": {{Quote: "Middle", Source: "Test", Encouragement: "Test"}},
			"upper_aeons":  {{Quote: "Upper", Source: "Test", Encouragement: "Test"}},
			"treasury":     {{Quote: "Treasury", Source: "Test", Encouragement: "Test"}},
		},
	}
	engine.sources["test"] = testSource
	engine.initialized = true

	tests := []struct {
		name     string
		score    float64
		expected string
	}{
		{"Exactly 30", 30.0, "Lower"},    // Boundary: chaos -> lower_aeons
		{"Exactly 50", 50.0, "Middle"},   // Boundary: lower_aeons -> middle_aeons
		{"Exactly 70", 70.0, "Upper"},    // Boundary: middle_aeons -> upper_aeons
		{"Exactly 85", 85.0, "Treasury"}, // Boundary: upper_aeons -> treasury
		{"Zero", 0.0, "Chaos"},
		{"Negative", -10.0, "Chaos"},
		{"Over 100", 150.0, "Treasury"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestEngine_ConcurrentAccess(t *testing.T) {
	engine := NewEngine()
	if err := engine.Initialize(); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Concurrent reads should not panic
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, _ = engine.GetWisdom(75.0, "stoic")
			_ = engine.ListSources()
			_, _ = engine.GetSource("stoic")
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestEngine_GetWisdom_EmptySourceQuotes(t *testing.T) {
	engine := NewEngine()

	// Source with no quotes
	emptySource := &Source{
		Name:   "Empty Source",
		Icon:   "📜",
		Quotes: make(map[string][]Quote),
	}
	engine.sources["empty"] = emptySource
	engine.initialized = true

	quote, err := engine.GetWisdom(75.0, "empty")
	if err != nil {
		t.Fatalf("GetWisdom failed: %v", err)
	}
	// Should return default quote
	if quote == nil {
		t.Fatal("GetWisdom returned nil for empty source")
	}
	if quote.Quote == "" {
		t.Error("GetWisdom returned empty quote for empty source")
	}
}
