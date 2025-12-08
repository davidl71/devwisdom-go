package wisdom

// GetBuiltInSources returns all built-in wisdom sources
// DEPRECATED: Use SourceLoader instead for configurable sources
// This function is kept for backward compatibility
func GetBuiltInSources() map[string]*Source {
	// Create a loader with default configuration
	loader := NewSourceLoader()
	
	// Try to load from default locations
	if err := loader.Load(); err == nil {
		return loader.GetAllSources()
	}

	// Fallback to minimal hard-coded source if loading fails
	sources := make(map[string]*Source)
	sources["bofh"] = &Source{
		Name:   "BOFH (Bastard Operator From Hell)",
		Icon:   "😈",
		Quotes: make(map[string][]Quote),
	}
	
	sources["bofh"].Quotes["chaos"] = []Quote{
		{
			Quote:        "It's not a bug, it's a feature.",
			Source:       "BOFH Excuse Calendar",
			Encouragement: "Document it and ship it.",
		},
	}

	return sources
}
