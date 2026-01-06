# devwisdom-go Examples

This directory contains practical examples demonstrating how to use devwisdom-go.

## Contents

- **CLI Usage** (`cli_usage.md`) - Command-line interface examples
- **MCP Integration** (`mcp_integration.md`) - Model Context Protocol integration examples
- **Programmatic Usage** (`programmatic_usage.go`) - Go API usage examples
- **Add Project Source** (`add_project_source.go`) - Programmatic source addition example

## Quick Start

### CLI Examples

```bash
# Get a random wisdom quote
devwisdom quote

# Get quote from specific source with score context
devwisdom quote --source stoic --score 75

# Consult an advisor for a metric
devwisdom consult --metric security --score 40

# List all available sources
devwisdom sources
```

### MCP Integration

The MCP server runs over stdio and implements JSON-RPC 2.0. See `mcp_integration.md` for detailed examples.

### Programmatic API Usage

For Go developers integrating devwisdom-go into their applications:

```bash
# Run the programmatic examples
go run examples/programmatic_usage.go
```

See `programmatic_usage.go` for comprehensive examples including:
- Basic quote retrieval
- Random source selection
- Advisor consultations
- Listing sources and advisors
- Working with aeon levels
- Custom source configuration
- Consultation modes

## More Information

- [CLI Usage Guide](./cli_usage.md)
- [MCP Integration Guide](./mcp_integration.md)
- [Programmatic API Examples](./programmatic_usage.go)
- [Main README](../README.md)

