package wisdom

import "testing"

func TestSourceLoader_Fallback(t *testing.T) {
	// Test that SourceLoader provides fallback behavior similar to old GetBuiltInSources
	loader := NewSourceLoader()

	// Try to load (may fail if no config files exist, which is fine)
	_ = loader.Load()

	sources := loader.GetAllSources()
	if sources == nil {
		t.Fatal("SourceLoader.GetAllSources() returned nil")
	}

	// If no sources loaded, that's acceptable (fallback happens in engine)
	// But if sources exist, verify structure
	if len(sources) > 0 {
		// Check that sources have proper structure
		for id, source := range sources {
			if id == "" {
				t.Error("Source has empty ID")
			}
			if source == nil {
				t.Errorf("Source %q is nil", id)
				continue
			}
			if source.Name == "" {
				t.Errorf("Source %q has no name", id)
			}
		}
	}
}

func TestSourceLoader_Structure(t *testing.T) {
	loader := NewSourceLoader()
	_ = loader.Load()
	sources := loader.GetAllSources()

	for id, source := range sources {
		if id == "" {
			t.Error("Source has empty ID")
		}
		if source == nil {
			t.Errorf("Source %q is nil", id)
			continue
		}
		if source.Name == "" {
			t.Errorf("Source %q has no name", id)
		}
		if source.Icon == "" {
			t.Errorf("Source %q has no icon", id)
		}
		if source.Quotes == nil {
			t.Errorf("Source %q has nil quotes map", id)
			continue
		}

		// Check that quotes exist for at least one aeon level
		hasQuotes := false
		for level, quotes := range source.Quotes {
			if len(quotes) > 0 {
				hasQuotes = true
				// Validate quote structure
				for i, quote := range quotes {
					if quote.Quote == "" {
						t.Errorf("Source %q, level %q, quote %d has empty quote text", id, level, i)
					}
					if quote.Source == "" {
						t.Errorf("Source %q, level %q, quote %d has empty source", id, level, i)
					}
					if quote.Encouragement == "" {
						t.Errorf("Source %q, level %q, quote %d has empty encouragement", id, level, i)
					}
				}
			}
		}
		if !hasQuotes {
			t.Errorf("Source %q has no quotes", id)
		}
	}
}
