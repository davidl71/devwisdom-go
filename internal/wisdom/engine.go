package wisdom

import (
	"fmt"
	"path/filepath"
	"sync"
)

// Engine is the main wisdom engine managing sources, advisors, and consultations
type Engine struct {
	sources   map[string]*Source
	advisors  *AdvisorRegistry
	config    *Config
	initialized bool
	mu        sync.RWMutex
}

// NewEngine creates a new wisdom engine instance
func NewEngine() *Engine {
	return &Engine{
		sources:  make(map[string]*Source),
		advisors: NewAdvisorRegistry(),
		config:   NewConfig(),
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

	// Initialize built-in wisdom sources
	if err := e.loadBuiltInSources(); err != nil {
		return fmt.Errorf("failed to load wisdom sources: %w", err)
	}

	// Initialize advisors
	e.advisors.Initialize()

	e.initialized = true
	return nil
}

// loadBuiltInSources loads all built-in wisdom sources
func (e *Engine) loadBuiltInSources() error {
	// Load sources from sources.go
	sources := GetBuiltInSources()
	
	for name, source := range sources {
		e.sources[name] = source
	}

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

	sources := make([]string, 0, len(e.sources))
	for name := range e.sources {
		sources = append(sources, name)
	}
	return sources
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	// Look for .exarp_wisdom_config in current directory or home
	// TODO: Implement proper config path resolution
	return filepath.Join(".", ".exarp_wisdom_config")
}
