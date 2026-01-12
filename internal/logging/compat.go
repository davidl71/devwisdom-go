// Package logging provides compatibility layer for mcp-go-core logging
// while maintaining backward compatibility with DEVWISDOM_DEBUG env var
package logging

import (
	"os"

	"github.com/davidl71/mcp-go-core/pkg/mcp/logging"
)

// NewLogger creates a logger that supports both DEVWISDOM_DEBUG and MCP_DEBUG env vars
// This maintains backward compatibility while transitioning to mcp-go-core
func NewLogger() *logging.Logger {
	// Check DEVWISDOM_DEBUG first for backward compatibility
	if os.Getenv("DEVWISDOM_DEBUG") == "1" {
		os.Setenv("MCP_DEBUG", "1")
	}
	
	// Use mcp-go-core logger
	return logging.NewLogger()
}

// Re-export types and constants for backward compatibility
type LogLevel = logging.LogLevel

const (
	LevelDebug = logging.LevelDebug
	LevelInfo  = logging.LevelInfo
	LevelWarn  = logging.LevelWarn
	LevelError = logging.LevelError
)
