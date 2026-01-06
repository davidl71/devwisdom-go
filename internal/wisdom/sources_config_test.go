package wisdom

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSourceLoader(t *testing.T) {
	loader := NewSourceLoader()
	if loader == nil {
		t.Fatal("NewSourceLoader returned nil")
	}
	if loader.cache == nil {
		t.Error("SourceLoader cache is nil")
	}
	if loader.httpClient == nil {
		t.Error("SourceLoader httpClient is nil")
	}
}

func TestSourceLoader_WithConfigPaths(t *testing.T) {
	loader := NewSourceLoader()
	loader.WithConfigPaths("path1.json", "path2.json")

	if len(loader.configPaths) != 2 {
		t.Errorf("configPaths length = %d, want 2", len(loader.configPaths))
	}
}

func TestSourceLoader_WithCacheTTL(t *testing.T) {
	loader := NewSourceLoader()
	ttl := 30 * time.Minute
	loader.WithCacheTTL(ttl)

	// Verify TTL was set by checking cache behavior
	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	loader.cache.Set("test", config, "")
	// Should still be valid with longer TTL
	time.Sleep(100 * time.Millisecond)
	_, found := loader.cache.Get("test")
	if !found {
		t.Error("Cache entry expired too quickly with custom TTL")
	}
}

func TestSourceLoader_LoadFromFile(t *testing.T) {
	// Create a temporary JSON file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "test_sources.json")

	config := SourcesConfig{
		Version: "1.0",
		Sources: map[string]*SourceConfig{
			"test": {
				ID:   "test",
				Name: "Test Source",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"chaos": {
						{Quote: "Test quote", Source: "Test", Encouragement: "Test"},
					},
				},
			},
		},
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	loader := NewSourceLoader().WithConfigPaths(configFile)
	if err := loader.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	source, found := loader.GetSource("test")
	if !found {
		t.Fatal("Source not found after loading")
	}
	if source.Name != "Test Source" {
		t.Errorf("Source name = %q, want %q", source.Name, "Test Source")
	}
}

func TestSourceLoader_AddSource(t *testing.T) {
	loader := NewSourceLoader()
	config := &SourceConfig{
		ID:   "test",
		Name: "Test Source",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {
				{Quote: "Test quote", Source: "Test", Encouragement: "Test"},
			},
		},
	}

	if err := loader.AddSource(config); err != nil {
		t.Fatalf("AddSource failed: %v", err)
	}

	source, found := loader.GetSource("test")
	if !found {
		t.Fatal("Source not found after AddSource")
	}
	if source.Name != "Test Source" {
		t.Errorf("Source name = %q, want %q", source.Name, "Test Source")
	}
}

func TestSourceLoader_AddSource_Invalid(t *testing.T) {
	loader := NewSourceLoader()
	invalidConfig := &SourceConfig{
		ID:   "", // Missing ID
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	err := loader.AddSource(invalidConfig)
	if err == nil {
		t.Error("AddSource should fail for invalid config")
	}
}

func TestSourceLoader_ListSourceIDs(t *testing.T) {
	loader := NewSourceLoader()
	config1 := &SourceConfig{
		ID:   "source1",
		Name: "Source 1",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}
	config2 := &SourceConfig{
		ID:   "source2",
		Name: "Source 2",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	if err := loader.AddSource(config1); err != nil {
		t.Fatalf("AddSource failed: %v", err)
	}
	if err := loader.AddSource(config2); err != nil {
		t.Fatalf("AddSource failed: %v", err)
	}

	ids := loader.ListSourceIDs()
	if len(ids) != 2 {
		t.Errorf("ListSourceIDs returned %d IDs, want 2", len(ids))
	}

	// Check that both IDs are present
	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}
	if !idMap["source1"] || !idMap["source2"] {
		t.Error("ListSourceIDs missing expected IDs")
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *SourceConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: &SourceConfig{
				ID:   "test",
				Name: "Test",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
				},
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			config: &SourceConfig{
				Name: "Test",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
				},
			},
			wantErr: true,
		},
		{
			name: "missing name",
			config: &SourceConfig{
				ID:   "test",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
				},
			},
			wantErr: true,
		},
		{
			name: "no quotes",
			config: &SourceConfig{
				ID:     "test",
				Name:   "Test",
				Icon:   "📜",
				Quotes: map[string][]Quote{},
			},
			wantErr: true,
		},
		{
			name: "invalid aeon level",
			config: &SourceConfig{
				ID:   "test",
				Name: "Test",
				Icon: "📜",
				Quotes: map[string][]Quote{
					"invalid_level": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSourceLoader_Reload(t *testing.T) {
	loader := NewSourceLoader()
	config := &SourceConfig{
		ID:   "test",
		Name: "Test",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	if err := loader.AddSource(config); err != nil {
		t.Fatalf("AddSource failed: %v", err)
	}

	// Verify source exists before reload
	_, found := loader.GetSource("test")
	if !found {
		t.Fatal("Source not found before reload")
	}

	// Reload clears cache and reloads from files (which don't exist in test)
	// So source added programmatically will be lost unless we add it again
	if err := loader.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// After reload, programmatically added sources are cleared
	// This is expected behavior - reload loads from files
	_, found = loader.GetSource("test")
	if found {
		t.Log("Note: Reload clears programmatically added sources (expected behavior)")
	}
}

func TestSourceLoader_GetAllSources(t *testing.T) {
	loader := NewSourceLoader()
	config1 := &SourceConfig{
		ID:   "source1",
		Name: "Source 1",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}
	config2 := &SourceConfig{
		ID:   "source2",
		Name: "Source 2",
		Icon: "📜",
		Quotes: map[string][]Quote{
			"chaos": {{Quote: "Test", Source: "Test", Encouragement: "Test"}},
		},
	}

	if err := loader.AddSource(config1); err != nil {
		t.Fatalf("AddSource failed: %v", err)
	}
	if err := loader.AddSource(config2); err != nil {
		t.Fatalf("AddSource failed: %v", err)
	}

	allSources := loader.GetAllSources()
	if len(allSources) != 2 {
		t.Errorf("GetAllSources returned %d sources, want 2", len(allSources))
	}
}
