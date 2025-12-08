package wisdom

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SourceConfig represents a configurable wisdom source
type SourceConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Icon        string            `json:"icon"`
	Description string            `json:"description,omitempty"`
	Language    string            `json:"language,omitempty"` // "hebrew", "english", etc.
	Quotes      map[string][]Quote `json:"quotes"`            // Key: aeon level
	// Optional fields for API-based sources
	SefariaSource string `json:"sefaria_source,omitempty"` // For Sefaria API sources
	APIEndpoint   string `json:"api_endpoint,omitempty"`   // For future API sources
}

// SourcesConfig represents the complete sources configuration
type SourcesConfig struct {
	Version string                    `json:"version"`
	Sources map[string]*SourceConfig  `json:"sources"`
	// Metadata
	LastUpdated string `json:"last_updated,omitempty"`
	Author      string `json:"author,omitempty"`
}

// SourceLoader handles loading sources from various locations
type SourceLoader struct {
	mu            sync.RWMutex
	sources       map[string]*Source
	configPaths   []string
	embeddedFS    *embed.FS
	embeddedPath  string
	reloadEnabled bool
	projectRoot   string // Project root directory for project-specific sources
	cache         *SourceCache
	httpClient    *http.Client // For API-based sources with timeout
}

// NewSourceLoader creates a new source loader
func NewSourceLoader() *SourceLoader {
	loader := &SourceLoader{
		sources:       make(map[string]*Source),
		configPaths:   []string{},
		reloadEnabled: true,
		projectRoot:   findProjectRoot(),
		cache:         NewSourceCache(),
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Default timeout for API calls
		},
	}

	// Start cache cleanup every 5 minutes
	loader.cache.StartCleanup(5 * time.Minute)

	return loader
}

// WithEmbeddedFS sets embedded filesystem for default sources
func (sl *SourceLoader) WithEmbeddedFS(fs *embed.FS, path string) *SourceLoader {
	sl.embeddedFS = fs
	sl.embeddedPath = path
	return sl
}

// WithConfigPaths adds configuration file paths to search
func (sl *SourceLoader) WithConfigPaths(paths ...string) *SourceLoader {
	sl.configPaths = append(sl.configPaths, paths...)
	return sl
}

// WithProjectRoot sets the project root directory
func (sl *SourceLoader) WithProjectRoot(root string) *SourceLoader {
	sl.projectRoot = root
	return sl
}

// WithReload enables or disables reloading
func (sl *SourceLoader) WithReload(enabled bool) *SourceLoader {
	sl.reloadEnabled = enabled
	return sl
}

// WithCacheTTL sets the cache TTL
func (sl *SourceLoader) WithCacheTTL(ttl time.Duration) *SourceLoader {
	sl.cache.WithTTL(ttl)
	return sl
}

// WithCacheMaxAge sets the maximum cache age
func (sl *SourceLoader) WithCacheMaxAge(maxAge time.Duration) *SourceLoader {
	sl.cache.WithMaxAge(maxAge)
	return sl
}

// WithCacheEnabled enables or disables caching
func (sl *SourceLoader) WithCacheEnabled(enabled bool) *SourceLoader {
	sl.cache.Enable(enabled)
	return sl
}

// WithHTTPTimeout sets the HTTP client timeout for API-based sources
func (sl *SourceLoader) WithHTTPTimeout(timeout time.Duration) *SourceLoader {
	sl.httpClient.Timeout = timeout
	return sl
}

// InvalidateCache clears the cache
func (sl *SourceLoader) InvalidateCache() {
	sl.cache.InvalidateAll()
}

// Load loads all sources from configured locations
func (sl *SourceLoader) Load() error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Start with empty sources
	sl.sources = make(map[string]*Source)
	
	// Clear existing configs
	configsMu.Lock()
	configs = make(map[string]*SourceConfig)
	configsMu.Unlock()

	// 1. Load embedded default sources (if available)
	if sl.embeddedFS != nil && sl.embeddedPath != "" {
		if err := sl.loadFromEmbedded(); err != nil {
			// Log but don't fail - embedded sources are optional
			fmt.Printf("Warning: Failed to load embedded sources: %v\n", err)
		}
	}

	// 2. Load from explicit config paths (in order, later files override earlier)
	for _, path := range sl.configPaths {
		if err := sl.loadFromFile(path); err != nil {
			// Log but continue - config files are optional
			fmt.Printf("Warning: Failed to load sources from %s: %v\n", path, err)
		}
	}

	// 3. Load from default locations (project-specific first, then global)
	sl.loadFromDefaultLocations()

	// 4. Convert SourceConfig to Source
	for id, config := range sl.getConfigs() {
		source := sl.configToSource(id, config)
		sl.sources[id] = source
	}

	return nil
}

// Reload reloads sources from all configured locations
func (sl *SourceLoader) Reload() error {
	if !sl.reloadEnabled {
		return nil
	}

	// Invalidate cache before reloading
	sl.cache.InvalidateAll()

	return sl.Load()
}

// GetSource retrieves a source by ID
func (sl *SourceLoader) GetSource(id string) (*Source, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	source, exists := sl.sources[id]
	return source, exists
}

// GetAllSources returns all loaded sources
func (sl *SourceLoader) GetAllSources() map[string]*Source {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	
	// Return a copy to prevent external modification
	result := make(map[string]*Source)
	for id, source := range sl.sources {
		result[id] = source
	}
	return result
}

// ListSourceIDs returns all available source IDs
func (sl *SourceLoader) ListSourceIDs() []string {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	
	ids := make([]string, 0, len(sl.sources))
	for id := range sl.sources {
		ids = append(ids, id)
	}
	return ids
}

// AddSource adds a source programmatically (useful for runtime additions)
func (sl *SourceLoader) AddSource(config *SourceConfig) error {
	if err := ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid source config: %w", err)
	}

	sl.mu.Lock()
	defer sl.mu.Unlock()

	// Add to configs
	configsMu.Lock()
	configs[config.ID] = config
	configsMu.Unlock()

	// Convert and add to sources
	source := sl.configToSource(config.ID, config)
	sl.sources[config.ID] = source

	return nil
}

// SaveProjectSource saves a source configuration to the project directory
func (sl *SourceLoader) SaveProjectSource(config *SourceConfig) error {
	if sl.projectRoot == "" {
		return fmt.Errorf("project root not found - cannot save project source")
	}

	// Create .wisdom directory in project root if it doesn't exist
	wisdomDir := filepath.Join(sl.projectRoot, ".wisdom")
	if err := os.MkdirAll(wisdomDir, 0755); err != nil {
		return fmt.Errorf("failed to create .wisdom directory: %w", err)
	}

	// Save to project-specific sources file
	projectSourcesFile := filepath.Join(wisdomDir, "sources.json")
	return SaveSourceConfig(projectSourcesFile, config)
}

// GetProjectSourcesPath returns the path where project sources are stored
func (sl *SourceLoader) GetProjectSourcesPath() string {
	if sl.projectRoot == "" {
		return ""
	}
	return filepath.Join(sl.projectRoot, ".wisdom", "sources.json")
}

// Internal state for configs (before conversion to Sources)
var (
	configsMu sync.Mutex
	configs   = make(map[string]*SourceConfig)
)

func (sl *SourceLoader) getConfigs() map[string]*SourceConfig {
	configsMu.Lock()
	defer configsMu.Unlock()
	
	result := make(map[string]*SourceConfig)
	for id, config := range configs {
		result[id] = config
	}
	return result
}

func (sl *SourceLoader) addConfig(config *SourceConfig) {
	configsMu.Lock()
	defer configsMu.Unlock()
	configs[config.ID] = config
}

// findProjectRoot finds the project root directory by looking for common markers
func findProjectRoot() string {
	// Start from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	current := cwd
	for {
		// Check for project markers
		markers := []string{".git", ".todo2", "go.mod", "package.json", "CMakeLists.txt", "Makefile"}
		for _, marker := range markers {
			if _, err := os.Stat(filepath.Join(current, marker)); err == nil {
				return current
			}
		}

		// Check for .wisdom directory (indicates project wants wisdom sources)
		if _, err := os.Stat(filepath.Join(current, ".wisdom")); err == nil {
			return current
		}

		// Move up one directory
		parent := filepath.Dir(current)
		if parent == current {
			// Reached filesystem root
			break
		}
		current = parent
	}

	// Fallback to current working directory
	return cwd
}

// loadFromEmbedded loads sources from embedded filesystem
func (sl *SourceLoader) loadFromEmbedded() error {
	if sl.embeddedFS == nil {
		return fmt.Errorf("no embedded filesystem configured")
	}

	data, err := sl.embeddedFS.ReadFile(sl.embeddedPath)
	if err != nil {
		return fmt.Errorf("failed to read embedded file: %w", err)
	}

	var sourcesConfig SourcesConfig
	if err := json.Unmarshal(data, &sourcesConfig); err != nil {
		return fmt.Errorf("failed to parse embedded config: %w", err)
	}

	// Add all sources from embedded config
	for id, config := range sourcesConfig.Sources {
		config.ID = id // Ensure ID is set
		sl.addConfig(config)
	}

	return nil
}

// loadFromFile loads sources from a JSON file (with caching)
func (sl *SourceLoader) loadFromFile(path string) error {
	// Check cache first
	cacheKey := fmt.Sprintf("file:%s", path)
	if cached, found := sl.cache.Get(cacheKey); found {
		// Use cached config
		sl.addConfig(cached)
		return nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var sourcesConfig SourcesConfig
	if err := json.Unmarshal(data, &sourcesConfig); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Add/override sources from file
	for id, config := range sourcesConfig.Sources {
		config.ID = id // Ensure ID is set
		
		// Cache individual source configs
		sourceCacheKey := fmt.Sprintf("source:%s:%s", path, id)
		sl.cache.Set(sourceCacheKey, config, path)
		
		sl.addConfig(config)
	}

	// Cache the entire file config (for quick lookup)
	sl.cache.Set(cacheKey, nil, path) // nil means "file loaded successfully"

	return nil
}

// loadFromDefaultLocations loads from standard config locations
// Priority: Project-specific sources override global sources
func (sl *SourceLoader) loadFromDefaultLocations() {
	// PROJECT-SPECIFIC SOURCES (highest priority)
	// These are loaded first but can be overridden by explicit paths
	if sl.projectRoot != "" {
		// Project root .wisdom directory
		sl.tryLoadPath(filepath.Join(sl.projectRoot, ".wisdom", "sources.json"))
		// Project root directly
		sl.tryLoadPath(filepath.Join(sl.projectRoot, "sources.json"))
		sl.tryLoadPath(filepath.Join(sl.projectRoot, "wisdom", "sources.json"))
	}

	// Current working directory (if different from project root)
	cwd, _ := os.Getwd()
	if cwd != sl.projectRoot {
		sl.tryLoadPath(filepath.Join(cwd, "sources.json"))
		sl.tryLoadPath(filepath.Join(cwd, "wisdom", "sources.json"))
		sl.tryLoadPath(filepath.Join(cwd, ".wisdom", "sources.json"))
	}

	// GLOBAL SOURCES (lower priority)
	// Home directory
	if home, err := os.UserHomeDir(); err == nil {
		sl.tryLoadPath(filepath.Join(home, ".wisdom", "sources.json"))
		sl.tryLoadPath(filepath.Join(home, ".exarp_wisdom", "sources.json"))
	}

	// XDG config directory
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		sl.tryLoadPath(filepath.Join(xdgConfig, "wisdom", "sources.json"))
	} else if home, err := os.UserHomeDir(); err == nil {
		sl.tryLoadPath(filepath.Join(home, ".config", "wisdom", "sources.json"))
	}
}

func (sl *SourceLoader) tryLoadPath(path string) {
	if err := sl.loadFromFile(path); err != nil {
		// Silently ignore - file might not exist
		return
	}
}

// configToSource converts SourceConfig to Source
func (sl *SourceLoader) configToSource(id string, config *SourceConfig) *Source {
	source := &Source{
		Name:   config.Name,
		Icon:   config.Icon,
		Quotes: make(map[string][]Quote),
	}

	if config.Description != "" {
		source.Description = config.Description
	}

	// Copy quotes by aeon level
	for level, quotes := range config.Quotes {
		source.Quotes[level] = make([]Quote, len(quotes))
		copy(source.Quotes[level], quotes)
	}

	return source
}

// ValidateConfig validates a source configuration
func ValidateConfig(config *SourceConfig) error {
	if config.ID == "" {
		return fmt.Errorf("source ID is required")
	}
	if config.Name == "" {
		return fmt.Errorf("source name is required")
	}
	if len(config.Quotes) == 0 {
		return fmt.Errorf("source must have at least one quote")
	}

	// Validate aeon levels
	validLevels := map[string]bool{
		"chaos":        true,
		"lower_aeons":  true,
		"middle_aeons": true,
		"upper_aeons":  true,
		"treasury":     true,
	}

	for level := range config.Quotes {
		if !validLevels[level] {
			return fmt.Errorf("invalid aeon level: %s", level)
		}
	}

	return nil
}

// SaveSourceConfig saves a source configuration to a file
// If the file exists, it will merge the new source with existing sources
func SaveSourceConfig(path string, config *SourceConfig) error {
	if err := ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	var sourcesConfig SourcesConfig

	// Try to load existing config
	if data, err := os.ReadFile(path); err == nil {
		if err := json.Unmarshal(data, &sourcesConfig); err == nil {
			// Use existing config
		}
	}

	// Initialize if needed
	if sourcesConfig.Sources == nil {
		sourcesConfig.Sources = make(map[string]*SourceConfig)
	}
	if sourcesConfig.Version == "" {
		sourcesConfig.Version = "1.0"
	}

	// Add or update the source
	sourcesConfig.Sources[config.ID] = config

	data, err := json.MarshalIndent(sourcesConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
