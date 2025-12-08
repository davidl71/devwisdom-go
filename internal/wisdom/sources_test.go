package wisdom

import "testing"

func TestGetBuiltInSources(t *testing.T) {
	sources := GetBuiltInSources()
	if sources == nil {
		t.Fatal("GetBuiltInSources returned nil")
	}
	if len(sources) == 0 {
		t.Error("GetBuiltInSources returned empty map")
	}
	
	// Check that at least bofh exists (our fallback)
	bofh, exists := sources["bofh"]
	if !exists {
		t.Error("GetBuiltInSources missing bofh source")
	}
	if bofh == nil {
		t.Error("bofh source is nil")
	}
	if bofh.Name == "" {
		t.Error("bofh source has no name")
	}
}

func TestGetBuiltInSources_Structure(t *testing.T) {
	sources := GetBuiltInSources()
	
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

