package wisdom

// GetBuiltInSources returns all built-in wisdom sources
// This will be populated with the actual wisdom data from Python version
func GetBuiltInSources() map[string]*Source {
	sources := make(map[string]*Source)

	// BOFH Source
	sources["bofh"] = &Source{
		Name:   "BOFH (Bastard Operator From Hell)",
		Icon:   "😈",
		Quotes: make(map[string][]Quote),
	}
	
	// Initialize with placeholder quotes (will be populated from Python data)
	sources["bofh"].Quotes["chaos"] = []Quote{
		{
			Quote:        "It's not a bug, it's a feature.",
			Source:       "BOFH Excuse Calendar",
			Encouragement: "Document it and ship it.",
		},
	}

	// TODO: Port all 21+ wisdom sources from Python version
	// Sources to port:
	// - pistis_sophia
	// - stoic
	// - tao
	// - art_of_war
	// - bible
	// - confucius
	// - tao_of_programming
	// - murphy
	// - shakespeare
	// - kybalion
	// - gracian
	// - enochian
	// - rebbe (Hebrew)
	// - tzaddik (Hebrew)
	// - chacham (Hebrew)
	// - pirkei_avot (Hebrew)
	// - proverbs (Hebrew)
	// - ecclesiastes (Hebrew)
	// - psalms (Hebrew)
	// - random (daily selection)

	return sources
}
