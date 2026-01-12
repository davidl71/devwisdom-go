# Lessons Learned from exarp-go

This document captures key patterns, practices, and architectural decisions from the `exarp-go` project that can be applied to `devwisdom-go` and future Go MCP servers.

## 📋 Table of Contents

1. [Architecture Patterns](#architecture-patterns)
2. [CLI/MCP Dual Mode](#climcp-dual-mode)
3. [Framework Abstraction](#framework-abstraction)
4. [Python Bridge Pattern](#python-bridge-pattern)
5. [Tool Registration System](#tool-registration-system)
6. [Configuration Management](#configuration-management)
7. [Error Handling & Logging](#error-handling--logging)
8. [Development Workflow](#development-workflow)
9. [Testing Patterns](#testing-patterns)
10. [Security Practices](#security-practices)

---

## Architecture Patterns

### Framework-Agnostic Design

**Key Insight**: Use interfaces to abstract framework-specific implementations, enabling easy framework switching.

**exarp-go Pattern**:
```go
// internal/framework/server.go
type MCPServer interface {
    RegisterTool(name, description string, schema ToolSchema, handler ToolHandler) error
    RegisterPrompt(name, description string, handler PromptHandler) error
    RegisterResource(uri, name, description, mimeType string, handler ResourceHandler) error
    Run(ctx context.Context, transport Transport) error
    GetName() string
    CallTool(ctx context.Context, name string, args json.RawMessage) ([]TextContent, error)
    ListTools() []ToolInfo
}
```

**Benefits**:
- Easy to swap frameworks (Go SDK, custom implementation, etc.)
- Testable with mock implementations
- Clear separation of concerns
- Future-proof for framework changes

**Applied to devwisdom-go**:
- Consider abstracting the MCP protocol implementation
- Create interfaces for wisdom engine operations
- Enable pluggable wisdom sources

### Factory Pattern

**Key Insight**: Centralize server creation logic with factory functions.

**exarp-go Pattern**:
```go
// internal/factory/server.go
func NewServer(frameworkType config.FrameworkType, name, version string) (framework.MCPServer, error)
func NewServerFromConfig(cfg *config.Config) (framework.MCPServer, error)
```

**Benefits**:
- Single point of server creation
- Consistent initialization
- Easy configuration-driven creation

---

## CLI/MCP Dual Mode

### TTY Detection Pattern

**Key Insight**: Detect execution mode (CLI vs MCP server) by checking if stdin is a TTY.

**exarp-go Pattern**:
```go
// cmd/server/main.go
func main() {
    // Detect if running in CLI mode (TTY) or MCP server mode (stdio)
    if cli.IsTTY() {
        // CLI mode - run command line interface
        if err := cli.Run(); err != nil {
            log.Fatalf("CLI error: %v", err)
        }
        return
    }

    // MCP server mode - run as stdio server
    // ... MCP server initialization
}
```

**Implementation**:
```go
// internal/cli/cli.go
func IsTTY() bool {
    return term.IsTerminal(int(os.Stdin.Fd()))
}
```

**Benefits**:
- Single binary for both modes
- Automatic mode detection
- Clean separation of concerns

**Applied to devwisdom-go**:
- ✅ Already has CLI mode (separate binary)
- Consider unified binary with TTY detection
- Would simplify deployment

### CLI Tool Execution

**Key Insight**: Reuse server infrastructure for CLI tool execution.

**exarp-go Pattern**:
```go
// internal/cli/cli.go
func executeTool(server framework.MCPServer, toolName, argsJSON string) error {
    ctx := context.Background()
    // Parse arguments
    var args map[string]interface{}
    json.Unmarshal([]byte(argsJSON), &args)
    
    // Execute via server interface
    result, err := server.CallTool(ctx, toolName, argsBytes)
    // Display results
}
```

**Benefits**:
- Code reuse between CLI and MCP
- Consistent behavior
- Single source of truth for tool logic

---

## Framework Abstraction

### Adapter Pattern

**Key Insight**: Use adapters to wrap framework-specific implementations.

**exarp-go Structure**:
```
internal/
  framework/
    server.go          # Interface definitions
    factory.go         # Factory functions
    adapters/
      gosdk/           # Go SDK adapter implementation
```

**Benefits**:
- Framework implementation details hidden
- Easy to add new framework adapters
- Core logic framework-independent

---

## Python Bridge Pattern

### Workspace Root Detection

**Key Insight**: Robust workspace root detection for finding bridge scripts.

**exarp-go Pattern**:
```go
// internal/bridge/python.go
func getWorkspaceRoot() string {
    // 1. Check environment variable first
    projectRoot := os.Getenv("PROJECT_ROOT")
    
    // 2. Ignore placeholder values
    if strings.Contains(projectRoot, "{{PROJECT_ROOT}}") || projectRoot == "" {
        // 3. Try executable location
        execPath, err := os.Executable()
        if err == nil {
            execDir := filepath.Dir(execPath)
            if filepath.Base(execDir) == "bin" {
                return filepath.Dir(execDir)  // parent of bin/
            }
            return execDir
        }
        
        // 4. Fallback to runtime.Caller
        _, filename, _, ok := runtime.Caller(0)
        if ok {
            return filepath.Join(filepath.Dir(filename), "..", "..")
        }
    }
    
    return "."
}
```

**Benefits**:
- Works in various deployment scenarios
- Handles development and production paths
- Graceful fallbacks

### Security Validation

**Key Insight**: Always validate paths before using them with subprocess execution.

**exarp-go Pattern**:
```go
// Validate workspace root path
validatedRoot, err := security.ValidatePath(workspaceRoot, workspaceRoot)
if err != nil {
    return "", fmt.Errorf("invalid workspace root: %w", err)
}

// Validate bridge script path is within workspace root
_, err = security.ValidatePath(bridgeScript, workspaceRoot)
if err != nil {
    return "", fmt.Errorf("invalid bridge script path: %w", err)
}
```

**Benefits**:
- Prevents path traversal attacks
- Ensures scripts are in expected locations
- Clear error messages

---

## Tool Registration System

### Batch Registration Pattern

**Key Insight**: Organize tool registration in batches for maintainability.

**exarp-go Pattern**:
```go
// internal/tools/registry.go
func RegisterAllTools(server framework.MCPServer) error {
    // Batch 1: Simple tools
    if err := registerBatch1Tools(server); err != nil {
        return fmt.Errorf("failed to register Batch 1 tools: %w", err)
    }
    
    // Batch 2: Medium tools
    if err := registerBatch2Tools(server); err != nil {
        return fmt.Errorf("failed to register Batch 2 tools: %w", err)
    }
    
    // Batch 3: Advanced tools
    if err := registerBatch3Tools(server); err != nil {
        return fmt.Errorf("failed to register Batch 3 tools: %w", err)
    }
    
    return nil
}
```

**Benefits**:
- Clear organization
- Easy to add new tools
- Error handling per batch
- Maintainable as tool count grows

### Schema Definition

**Key Insight**: Use structured schema definitions for tool parameters.

**exarp-go Pattern**:
```go
framework.ToolSchema{
    Type: "object",
    Properties: map[string]interface{}{
        "action": map[string]interface{}{
            "type":    "string",
            "enum":    []string{"run", "analyze"},
            "default": "run",
        },
        "path": map[string]interface{}{
            "type": "string",
        },
    },
}
```

**Benefits**:
- Type-safe parameter definitions
- Validation built into schema
- Self-documenting
- IDE-friendly

---

## Configuration Management

### Environment-Based Configuration

**Key Insight**: Support both environment variables and defaults.

**exarp-go Pattern**:
```go
// internal/config/config.go
type Config struct {
    Framework FrameworkType `yaml:"framework" env:"MCP_FRAMEWORK"`
    Name      string        `yaml:"name" env:"MCP_SERVER_NAME"`
    Version   string        `yaml:"version" env:"MCP_VERSION"`
}

func Load() (*Config, error) {
    cfg := &Config{
        Framework: FrameworkGoSDK, // Default
        Name:      "exarp-go",
        Version:   "1.0.0",
    }
    
    // Override from environment
    if frameworkStr := os.Getenv("MCP_FRAMEWORK"); frameworkStr != "" {
        cfg.Framework = FrameworkType(frameworkStr)
    }
    // ... more overrides
    
    return cfg, nil
}
```

**Benefits**:
- Sensible defaults
- Environment variable override
- Easy deployment configuration
- No external config files required

---

## Error Handling & Logging

### Structured Logging

**Key Insight**: Use structured logging with request IDs for traceability.

**exarp-go Pattern** (from devwisdom-go, inspired by exarp-go):
```go
// Log request start
requestID := formatRequestID(req.ID)
startTime := time.Now()
s.appLogger.LogRequest(requestID, req.Method)

// Process request
resp := s.handleRequest(&req)

// Log request completion
duration := time.Since(startTime)
s.appLogger.LogRequestComplete(requestID, req.Method, duration)
```

**Benefits**:
- Request traceability
- Performance monitoring
- Debugging support
- Production-ready logging

### Error Context

**Key Insight**: Always provide context in error messages.

**exarp-go Pattern**:
```go
return fmt.Errorf("failed to register Batch 1 tools: %w", err)
return fmt.Errorf("invalid workspace root: %w", err)
```

**Benefits**:
- Easier debugging
- Clear error propagation
- Actionable error messages

---

## Development Workflow

### Makefile Patterns

**Key Insight**: Comprehensive Makefile with configuration detection.

**exarp-go Features**:
- `make config` - Detect available tools
- `make dev-full` - Full dev mode (auto-reload + auto-test)
- `make test-watch` - Test watch mode
- Tool detection (uv, pytest, file watchers)
- Conditional targets based on available tools

**Benefits**:
- One command development setup
- Automatic tool detection
- Developer-friendly workflow

### Hot Reload Scripts

**Key Insight**: Development scripts that watch files and auto-reload.

**exarp-go Approach**:
- Uses `fswatch` (macOS) or `inotifywait` (Linux)
- Falls back to polling if no file watcher available
- Auto-rebuilds on Go file changes
- Auto-reruns tests on changes

---

## Testing Patterns

### Tool Testing Utilities

**Key Insight**: Provide utilities for testing tools with example arguments.

**exarp-go Pattern**:
```go
// internal/cli/cli.go
func testToolExecution(server framework.MCPServer, toolName string) error {
    // Get tool info
    tools := server.ListTools()
    var toolInfo *framework.ToolInfo
    // Find tool...
    
    // Generate example arguments from schema
    exampleArgs := generateExampleArgs(toolInfo.Schema)
    
    // Execute with example arguments
    result, err := server.CallTool(context.Background(), toolName, argsBytes)
    // Display results
}
```

**Benefits**:
- Easy tool testing
- Schema-driven test generation
- Quick verification of tools

---

## Security Practices

### Path Validation

**Key Insight**: Always validate paths before file operations or subprocess execution.

**exarp-go Pattern**:
```go
// internal/security/path.go
func ValidatePath(path, root string) (string, error) {
    // Resolve absolute paths
    // Check path is within root
    // Prevent path traversal
}
```

**Benefits**:
- Prevents path traversal attacks
- Safe subprocess execution
- Clear security boundaries

---

## Recommendations for devwisdom-go

### Immediate Improvements

1. **Unified Binary**: Consider combining CLI and MCP server into single binary with TTY detection
2. **Framework Abstraction**: Create interfaces for MCP protocol operations
3. **Batch Tool Registration**: Organize tool registration in batches
4. **Configuration Management**: Add environment variable support for configuration

### Future Enhancements

1. **Python Bridge**: If adding Python tools, use the bridge pattern from exarp-go
2. **Development Scripts**: Add hot-reload scripts for development
3. **Testing Utilities**: Add tool testing utilities
4. **Path Security**: Add path validation for any file operations

### Architecture Considerations

1. **Modularity**: Keep core wisdom logic framework-independent
2. **Testability**: Use interfaces for all external dependencies
3. **Extensibility**: Design for easy addition of new wisdom sources
4. **Maintainability**: Follow the batch registration pattern as tools grow

---

## Conclusion

exarp-go demonstrates several proven patterns for building production-ready Go MCP servers:

- **Framework-agnostic design** enables flexibility
- **CLI/MCP dual mode** provides excellent developer experience
- **Python bridge pattern** enables gradual migration
- **Structured registration** scales well
- **Security-first approach** prevents vulnerabilities

These patterns can be incrementally adopted in devwisdom-go to improve maintainability, testability, and developer experience.

