package wisdom

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"sync"
	"time"

	"github.com/davidl71/devwisdom-go/internal/config"
)

// Engine is the main wisdom engine managing sources, advisors, and consultations.
// It provides thread-safe access to wisdom sources and advisor consultations.
// The engine must be initialized before use by calling Initialize().
type Engine struct {
	sources     map[string]*Source
	loader      *SourceLoader
	advisors    *AdvisorRegistry
	config      *config.Config
	initialized bool
	mu          sync.RWMutex
}

// NewEngine creates a new wisdom engine instance.
// The engine is not initialized by default; call Initialize() before use.
func NewEngine() *Engine {
	return &Engine{
		sources:  make(map[string]*Source),
		loader:   NewSourceLoader(),
		advisors: NewAdvisorRegistry(),
		config:   config.NewConfig(),
	}
}

// Initialize loads wisdom sources and configuration.
// This method is idempotent and can be called multiple times safely.
// It loads sources from configuration files or falls back to built-in sources.
func (e *Engine) Initialize() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.initialized {
		return nil
	}

	// Load configuration
	if err := e.config.Load(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Configure source loader
	e.loader = NewSourceLoader().
		WithConfigPaths(
			"sources.json",
			"wisdom/sources.json",
			".wisdom/sources.json",
		)

	// Try to load from default locations
	if err := e.loader.Load(); err != nil {
		// Fallback to hard-coded sources if config loading fails
		e.sources = GetBuiltInSources()
	} else {
		// Use loaded sources
		e.sources = e.loader.GetAllSources()
	}

	// Initialize advisors
	e.advisors.Initialize()

	e.initialized = true
	return nil
}

// ReloadSources reloads sources from configuration files.
// This is useful when sources are updated externally and you want to refresh the engine.
func (e *Engine) ReloadSources() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.loader.Reload(); err != nil {
		return fmt.Errorf("failed to reload sources: %w", err)
	}

	e.sources = e.loader.GetAllSources()
	return nil
}

// GetWisdom retrieves a wisdom quote based on score and source.
// The score determines the aeon level, which selects appropriate quotes from the source.
// If source is "random", a date-seeded random source is selected for consistency.
func (e *Engine) GetWisdom(score float64, source string) (*Quote, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
	}

	// Handle "random" source selection
	if source == "random" {
		randomSource, err := e.getRandomSourceLocked(true)
		if err != nil {
			return nil, fmt.Errorf("failed to get random source: %w", err)
		}
		source = randomSource
	}

	src, exists := e.sources[source]
	if !exists {
		return nil, fmt.Errorf("unknown source: %s", source)
	}

	// Determine aeon level from score
	aeonLevel := GetAeonLevel(score)

	// Get quote from source based on aeon level
	quote := src.GetQuote(aeonLevel)
	return quote, nil
}

// GetRandomSource returns a random wisdom source ID.
// If seedDate is true, the same source will be returned for the entire day (date-seeded).
// This ensures consistent daily source selection across sessions.
func (e *Engine) GetRandomSource(seedDate bool) (string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.getRandomSourceLocked(seedDate)
}

// getRandomSourceLocked is the internal implementation (assumes RLock is held)
func (e *Engine) getRandomSourceLocked(seedDate bool) (string, error) {
	if !e.initialized {
		return "", fmt.Errorf("engine not initialized")
	}

	// Get all available source IDs (excluding Sefaria API sources for now)
	allSources := make([]string, 0, len(e.sources))
	for id := range e.sources {
		// Exclude Sefaria API sources (they require API integration)
		if id != "pirkei_avot" && id != "proverbs" && id != "ecclesiastes" && id != "psalms" {
			allSources = append(allSources, id)
		}
	}

	if len(allSources) == 0 {
		return "", fmt.Errorf("no sources available")
	}

	// Date-seeded random selection for consistency
	var seed int64
	if seedDate {
		now := time.Now()
		dateStr := now.Format("20060102") // YYYYMMDD format
		
		// Convert date string to int and add hash offset (matching Python implementation)
		var dateInt int64
		fmt.Sscanf(dateStr, "%d", &dateInt)
		
		// Hash "random_source" string for offset (matching Python hash("random_source"))
		h := fnv.New32a()
		h.Write([]byte("random_source"))
		hashOffset := int64(h.Sum32())
		
		seed = dateInt + hashOffset
	} else {
		seed = time.Now().UnixNano()
	}

	// Create seeded random generator
	rng := rand.New(rand.NewSource(seed))
	
	// Select random source
	selectedIndex := rng.Intn(len(allSources))
	return allSources[selectedIndex], nil
}

// ListSources returns all available wisdom source IDs.
// Returns an empty slice if the engine is not initialized.
func (e *Engine) ListSources() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.loader != nil {
		return e.loader.ListSourceIDs()
	}

	sources := make([]string, 0, len(e.sources))
	for name := range e.sources {
		sources = append(sources, name)
	}
	return sources
}

// GetSource returns a specific source by ID.
// The second return value indicates whether the source was found.
func (e *Engine) GetSource(id string) (*Source, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.loader != nil {
		return e.loader.GetSource(id)
	}

	source, exists := e.sources[id]
	return source, exists
}

// GetLoader returns the source loader for advanced usage.
// This allows direct access to source loading and configuration management.
func (e *Engine) GetLoader() *SourceLoader {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.loader
}

// GetAdvisors returns the advisor registry.
// This provides access to advisor mappings and consultation functionality.
func (e *Engine) GetAdvisors() *AdvisorRegistry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.advisors
}

// AddProjectSource adds a source and saves it to the project directory.
// The source is persisted to the project's sources.json file and immediately available.
func (e *Engine) AddProjectSource(config *SourceConfig) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.loader == nil {
		return fmt.Errorf("loader not initialized")
	}

	// Save to project directory
	if err := e.loader.SaveProjectSource(config); err != nil {
		return fmt.Errorf("failed to save project source: %w", err)
	}

	// Add to loader (this will also update the loader's internal sources)
	if err := e.loader.AddSource(config); err != nil {
		return fmt.Errorf("failed to add source: %w", err)
	}

	// Reload to ensure consistency
	if err := e.loader.Reload(); err != nil {
		return fmt.Errorf("failed to reload sources: %w", err)
	}

	// Update engine's sources map from loader
	e.sources = e.loader.GetAllSources()

	return nil
}
