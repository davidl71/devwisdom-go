package wisdom

import (
	"fmt"
	"sync"

	"github.com/davidl71/devwisdom-go/internal/config"
)

// Engine is the main wisdom engine managing sources, advisors, and consultations
type Engine struct {
	sources     map[string]*Source
	loader      *SourceLoader
	advisors    *AdvisorRegistry
	config      *config.Config
	initialized bool
	mu          sync.RWMutex
}

// NewEngine creates a new wisdom engine instance
func NewEngine() *Engine {
	return &Engine{
		sources:  make(map[string]*Source),
		loader:   NewSourceLoader(),
		advisors: NewAdvisorRegistry(),
		config:   config.NewConfig(),
	}
}

// Initialize loads wisdom sources and configuration
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

// ReloadSources reloads sources from configuration files
func (e *Engine) ReloadSources() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.loader.Reload(); err != nil {
		return fmt.Errorf("failed to reload sources: %w", err)
	}

	e.sources = e.loader.GetAllSources()
	return nil
}

// GetWisdom retrieves wisdom quote based on score and source
func (e *Engine) GetWisdom(score float64, source string) (*Quote, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.initialized {
		return nil, fmt.Errorf("engine not initialized")
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

// ListSources returns all available wisdom sources
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

// GetSource returns a specific source by ID
func (e *Engine) GetSource(id string) (*Source, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.loader != nil {
		return e.loader.GetSource(id)
	}

	source, exists := e.sources[id]
	return source, exists
}

// GetLoader returns the source loader (for advanced usage)
func (e *Engine) GetLoader() *SourceLoader {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.loader
}

// GetAdvisors returns the advisor registry
func (e *Engine) GetAdvisors() *AdvisorRegistry {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.advisors
}

// AddProjectSource adds a source and saves it to the project directory
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
