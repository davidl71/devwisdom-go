# FastMCP Go Library Design Document

**Date:** 2026-01-09  
**Purpose:** Design specification for a standalone Go FastMCP library  
**Status:** Design Phase

---

## Overview

This document specifies the design for a **FastMCP-style wrapper library** around the official Model Context Protocol Go SDK. The library provides a simplified, fluent API similar to Python FastMCP while maintaining the reliability and spec compliance of the official SDK.

### Design Philosophy

**Best of Both Worlds:**
- ✅ **FastMCP Simplicity** - Fluent API, less boilerplate, automatic type inference
- ✅ **Official SDK Reliability** - Spec compliance, long-term support, proven implementation
- ✅ **Go Idioms** - Follows Go conventions and best practices
- ✅ **Flexibility** - Supports both automatic and explicit schema definition

---

## Package Structure

```
fastmcp/
├── app.go          # Main FastMCP app builder and core types
├── tool.go         # Tool registration with type inference
├── resource.go     # Resource registration
├── prompt.go       # Prompt registration
├── transport.go    # Transport implementations (STDIO, WebSocket, SSE)
├── schema.go       # Schema generation from Go types
└── internal/
    ├── reflect.go  # Reflection utilities for type inference
    └── types.go    # Internal type definitions
```

---

## Core API Design

### Main App Builder

```go
package fastmcp

import (
    "context"
    "github.com/modelcontextprotocol/go-sdk/mcp"
)

// App represents a FastMCP application
type App struct {
    server *mcp.Server
    name   string
    version string
}

// New creates a new FastMCP application
func New(name, version string) *App {
    return &App{
        server: mcp.NewServer(&mcp.Implementation{
            Name:    name,
            Version: version,
        }, nil),
        name:    name,
        version: version,
    }
}
```

### Fluent API Methods

```go
// Tool registers a tool with automatic schema generation
func (a *App) Tool(name string, handler interface{}, description string) *App

// Resource registers a resource handler
func (a *App) Resource(uri string, handler interface{}, description string) *App

// Prompt registers a prompt template
func (a *App) Prompt(name string, handler interface{}, description string) *App

// RunStdio runs the server on STDIO transport
func (a *App) RunStdio(ctx context.Context) error

// RunWebSocket runs the server on WebSocket transport
func (a *App) RunWebSocket(ctx context.Context, addr string) error

// RunSSE runs the server on Server-Sent Events transport
func (a *App) RunSSE(ctx context.Context, addr string) error
```

---

## Tool Registration

### Automatic Type Inference

The library automatically generates JSON schemas from Go function signatures:

```go
// Simple handler: func(string) string
app.Tool("greet", func(name string) string {
    return "Hello, " + name + "!"
}, "Greet a person")

// Handler with error: func(string) (string, error)
app.Tool("read", func(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err
    }
    return string(data), nil
}, "Read a file")

// Handler with struct input: func(Input) (Output, error)
type SearchParams struct {
    Query string `json:"query" jsonschema:"search query"`
    Limit int    `json:"limit" jsonschema:"result limit"`
}

type SearchResult struct {
    Results []string `json:"results"`
    Count   int      `json:"count"`
}

app.Tool("search", func(params SearchParams) (SearchResult, error) {
    // Implementation
    return SearchResult{
        Results: []string{"result1", "result2"},
        Count:   2,
    }, nil
}, "Search for items")
```

### Supported Handler Signatures

The library supports various handler function signatures:

```go
// Simple handlers
func() string
func(string) string
func(int) (int, error)

// Struct-based handlers
func(Input) Output
func(Input) (Output, error)
func(Input) (*Output, error)

// Context-aware handlers
func(context.Context, Input) (Output, error)
func(context.Context, *mcp.CallToolRequest, Input) (*mcp.CallToolResult, Output, error)
```

### Explicit Schema (Advanced)

For complex cases, explicit schema can be provided:

```go
app.ToolWithSchema("complex", &mcp.Tool{
    Name:        "complex",
    Description: "Complex tool with explicit schema",
    InputSchema: &jsonschema.Schema{
        Type: "object",
        Properties: map[string]*jsonschema.Schema{
            "nested": {
                Type: "object",
                Properties: map[string]*jsonschema.Schema{
                    "field": {Type: "string"},
                },
            },
        },
    },
}, func(ctx context.Context, req *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
    // Handler implementation
    return nil, map[string]interface{}{"result": "success"}, nil
})
```

---

## Resource Registration

Resources are registered with URI patterns and handlers:

```go
// Simple resource handler
app.Resource("files/{path}", func(path string) ([]byte, string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, "", err
    }
    return data, "text/plain", nil
}, "Access files")

// Resource with context
app.Resource("config/{key}", func(ctx context.Context, key string) ([]byte, string, error) {
    // Implementation
    return []byte("value"), "application/json", nil
}, "Get configuration")
```

### Resource Handler Signatures

```go
// Simple: func(string) ([]byte, string, error)
func(path string) ([]byte, string, error)

// With context: func(context.Context, string) ([]byte, string, error)
func(ctx context.Context, path string) ([]byte, string, error)

// With MCP request: func(context.Context, *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)
func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)
```

---

## Prompt Registration

Prompts are registered with template handlers:

```go
// Simple prompt handler
app.Prompt("greeting", func(name string) string {
    return fmt.Sprintf("Hello, %s! How can I help you today?", name)
}, "Generate a greeting")

// Prompt with arguments
type PromptArgs struct {
    UserName string `json:"user_name"`
    Context  string `json:"context"`
}

app.Prompt("custom", func(args PromptArgs) string {
    return fmt.Sprintf("Welcome %s! Context: %s", args.UserName, args.Context)
}, "Custom prompt with arguments")
```

---

## Transport Support

### STDIO Transport

```go
app := fastmcp.New("my-app", "v1.0.0").
    Tool("greet", func(name string) string {
        return "Hello, " + name + "!"
    }, "Greet a person").
    RunStdio(context.Background())
```

### WebSocket Transport

```go
app := fastmcp.New("my-app", "v1.0.0").
    Tool("greet", func(name string) string {
        return "Hello, " + name + "!"
    }, "Greet a person").
    RunWebSocket(context.Background(), ":8080")
```

### Server-Sent Events (SSE) Transport

```go
app := fastmcp.New("my-app", "v1.0.0").
    Tool("greet", func(name string) string {
        return "Hello, " + name + "!"
    }, "Greet a person").
    RunSSE(context.Background(), ":8080")
```

---

## Schema Generation

The library automatically generates JSON schemas from Go types using reflection:

### Type Mapping

| Go Type | JSON Schema Type |
|---------|------------------|
| `string` | `"type": "string"` |
| `int`, `int64` | `"type": "integer"` |
| `float64` | `"type": "number"` |
| `bool` | `"type": "boolean"` |
| `[]T` | `"type": "array"` |
| `struct` | `"type": "object"` |
| `map[string]T` | `"type": "object"` |

### JSON Schema Tags

Go struct tags are used to customize schema generation:

```go
type User struct {
    Name     string `json:"name" jsonschema:"User's full name"`
    Age      int    `json:"age" jsonschema:"User's age" jsonschema_minimum:"0" jsonschema_maximum:"150"`
    Email    string `json:"email" jsonschema:"User's email address" jsonschema_format:"email"`
    Optional string `json:"optional,omitempty" jsonschema:"Optional field"`
}
```

### Supported JSON Schema Tags

- `jsonschema` - Description
- `jsonschema_minimum` - Minimum value
- `jsonschema_maximum` - Maximum value
- `jsonschema_format` - Format (email, uri, date-time, etc.)
- `jsonschema_enum` - Enum values
- `jsonschema_default` - Default value
- `jsonschema_required` - Mark as required (default: true for non-pointer fields)

---

## Complete Example

```go
package main

import (
    "context"
    "log"
    "os"
    
    "github.com/yourorg/fastmcp"
)

func main() {
    app := fastmcp.New("my-mcp-server", "v1.0.0").
        // Register a simple tool
        Tool("greet", func(name string) string {
            return "Hello, " + name + "!"
        }, "Greet a person").
        
        // Register a tool with struct input
        Tool("search", func(params SearchParams) (SearchResult, error) {
            // Implementation
            return SearchResult{
                Results: []string{"result1", "result2"},
                Count:   2,
            }, nil
        }, "Search for items").
        
        // Register a resource
        Resource("files/{path}", func(path string) ([]byte, string, error) {
            data, err := os.ReadFile(path)
            if err != nil {
                return nil, "", err
            }
            return data, "text/plain", nil
        }, "Access files").
        
        // Register a prompt
        Prompt("greeting", func(name string) string {
            return fmt.Sprintf("Hello, %s! How can I help you?", name)
        }, "Generate a greeting").
        
        // Run on STDIO
        RunStdio(context.Background())
    
    if err := app; err != nil {
        log.Fatal(err)
    }
}

type SearchParams struct {
    Query string `json:"query" jsonschema:"search query"`
    Limit int    `json:"limit" jsonschema:"result limit" jsonschema_minimum:"1" jsonschema_maximum:"100"`
}

type SearchResult struct {
    Results []string `json:"results"`
    Count   int      `json:"count"`
}
```

---

## Comparison with Official SDK

### Official SDK Style

```go
server := mcp.NewServer(&mcp.Implementation{
    Name:    "my-server",
    Version: "v1.0.0",
}, nil)

type GreetArgs struct {
    Name string `json:"name" jsonschema:"the person to greet"`
}

mcp.AddTool(server, &mcp.Tool{
    Name:        "greet",
    Description: "Greet a person",
}, func(ctx context.Context, req *mcp.CallToolRequest, args GreetArgs) (*mcp.CallToolResult, any, error) {
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: "Hello, " + args.Name + "!"},
        },
    }, nil, nil
})

if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
    log.Fatal(err)
}
```

### FastMCP Style

```go
app := fastmcp.New("my-server", "v1.0.0").
    Tool("greet", func(name string) string {
        return "Hello, " + name + "!"
    }, "Greet a person").
    RunStdio(context.Background())
```

**Benefits:**
- ✅ **Less boilerplate** - ~10 lines vs ~30+ lines
- ✅ **Automatic schema** - No explicit schema definition needed
- ✅ **Fluent API** - Chainable methods
- ✅ **Simpler handlers** - Direct return values instead of MCP result types

---

## Comparison with Third-Party FastMCP Go

### Third-Party Package

```go
app := fastmcp.New("My App").
    Tool("greet", func(name string) string {
        return "Hello, " + name + "!"
    }, "Greet a person").
    RunStdio()
```

### Our Design

```go
app := fastmcp.New("My App", "v1.0.0").
    Tool("greet", func(name string) string {
        return "Hello, " + name + "!"
    }, "Greet a person").
    RunStdio(context.Background())
```

**Key Differences:**
- ✅ **Official SDK Backend** - Uses official SDK for reliability
- ✅ **Context Support** - Explicit context.Context parameter
- ✅ **Version Required** - Version parameter for better tracking
- ✅ **Go Idioms** - Follows Go conventions more closely

---

## Design Decisions

### 1. Official SDK Backend

**Decision:** Use official Model Context Protocol Go SDK as the backend.

**Rationale:**
- ✅ Guaranteed spec compliance
- ✅ Long-term support and maintenance
- ✅ Proven, production-ready implementation
- ✅ Active development and updates

**Trade-offs:**
- ⚠️ Slightly more complex internal implementation
- ✅ Better long-term reliability

### 2. Automatic Schema Generation

**Decision:** Automatically generate JSON schemas from Go function signatures.

**Rationale:**
- ✅ Reduces boilerplate significantly
- ✅ Type-safe at compile time
- ✅ Familiar to Go developers
- ✅ Matches FastMCP philosophy

**Trade-offs:**
- ⚠️ Less control for complex schemas
- ✅ Can provide explicit schema when needed

### 3. Fluent API

**Decision:** Use chainable methods for building servers.

**Rationale:**
- ✅ Readable and intuitive
- ✅ Matches FastMCP patterns
- ✅ Reduces nesting and complexity

**Trade-offs:**
- ⚠️ Less flexible than builder pattern
- ✅ Simpler for common use cases

### 4. Context Support

**Decision:** Require explicit `context.Context` for transport methods.

**Rationale:**
- ✅ Follows Go best practices
- ✅ Enables cancellation and timeouts
- ✅ Standard Go pattern

**Trade-offs:**
- ⚠️ Slightly more verbose than implicit context
- ✅ Better control and idiomatic Go

### 5. FastMCP Context Object

**Decision:** Provide a FastMCP-style Context object wrapping Go's `context.Context` with advanced features.

**Rationale:**
- ✅ Matches Python FastMCP patterns
- ✅ Provides logging, progress reporting, and resource access
- ✅ Enhances developer experience
- ✅ Enables interactive workflows

**Trade-offs:**
- ⚠️ Additional abstraction layer
- ✅ Significantly enhanced functionality

---

## Implementation Strategy

### Phase 1: Core Infrastructure

1. **App Builder** - Basic app structure and initialization
2. **Schema Generation** - Reflection-based schema generation from types
3. **Tool Registration** - Basic tool registration with type inference

### Phase 2: Transport Support

1. **STDIO Transport** - Primary transport for MCP servers
2. **WebSocket Transport** - HTTP-based WebSocket support
3. **SSE Transport** - Server-Sent Events support

### Phase 3: Advanced Features

1. **Resource Registration** - Resource handlers with URI patterns
2. **Prompt Registration** - Prompt template handlers
3. **Explicit Schema Support** - Advanced schema customization

### Phase 4: Polish

1. **Error Handling** - Comprehensive error messages
2. **Documentation** - Complete API documentation
3. **Examples** - Comprehensive usage examples
4. **Testing** - Full test coverage

---

## Future Enhancements

### Potential Features

1. **Middleware Support** - Request/response middleware
2. **Logging Integration** - Structured logging support
3. **Metrics** - Built-in metrics and observability
4. **Validation** - Enhanced input validation
5. **Code Generation** - Generate schemas at build time
6. **Plugin System** - Extensible plugin architecture

### Considerations

- Keep API simple and focused
- Maintain compatibility with official SDK
- Follow Go best practices
- Prioritize developer experience

---

## References

- [Official Model Context Protocol Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Python FastMCP (Original)](https://github.com/jlowin/fastmcp)
- [TypeScript FastMCP](https://github.com/punkpeye/fastmcp)
- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [FastMCP Go SDK Analysis](../exarp-go/docs/archive/analysis/FASTMCP_GO_ANALYSIS.md)

---

## Conclusion

This design provides a **FastMCP-style wrapper** around the official Model Context Protocol Go SDK, combining:

- ✅ **Simplicity** of FastMCP API
- ✅ **Reliability** of official SDK
- ✅ **Go Idioms** and best practices
- ✅ **Flexibility** for advanced use cases

The library enables developers to build MCP servers with minimal boilerplate while maintaining full compatibility with the MCP specification and official SDK features.

