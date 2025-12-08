# Configurable Wisdom Sources

**Version**: 1.0  
**Date**: 2025-01-26

---

## Overview

The wisdom system now supports configurable sources loaded from JSON files instead of hard-coded Go structs. This makes it easy to:

- Add new wisdom sources without recompiling
- Customize existing sources
- Share source configurations
- Update quotes without code changes
- Support multiple source files

---

## Configuration File Format

### Basic Structure

```json
{
  "version": "1.0",
  "last_updated": "2025-01-26",
  "author": "devwisdom-go",
  "sources": {
    "source_id": {
      "id": "source_id",
      "name": "Source Display Name",
      "icon": "📜",
      "description": "Optional description",
      "language": "english",
      "quotes": {
        "chaos": [...],
        "lower_aeons": [...],
        "middle_aeons": [...],
        "upper_aeons": [...],
        "treasury": [...]
      }
    }
  }
}
```

### Quote Format

```json
{
  "quote": "The quote text",
  "source": "Source attribution (chapter, book, etc.)",
  "encouragement": "Encouraging message for developers"
}
```

### Source Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Unique identifier for the source |
| `name` | string | Yes | Display name |
| `icon` | string | Yes | Emoji or icon |
| `description` | string | No | Optional description |
| `language` | string | No | Language code (e.g., "hebrew", "english") |
| `quotes` | object | Yes | Quotes organized by aeon level |
| `sefaria_source` | string | No | For Sefaria API sources |
| `api_endpoint` | string | No | For future API-based sources |

### Aeon Levels

Quotes must be organized into these aeon levels:

- `chaos` - 0-30% health score
- `lower_aeons` - 31-50% health score
- `middle_aeons` - 51-70% health score
- `upper_aeons` - 71-85% health score
- `treasury` - 86-100% health score

---

## Configuration File Locations

The system searches for configuration files in this order (later files override earlier):

1. **Embedded sources** (if compiled in)
2. **Explicit paths** (via `WithConfigPaths()`)
3. **Current directory**: `sources.json`
4. **Current directory**: `wisdom/sources.json`
5. **Current directory**: `.wisdom/sources.json`
6. **Home directory**: `~/.wisdom/sources.json`
7. **Home directory**: `~/.exarp_wisdom/sources.json`
8. **XDG config**: `$XDG_CONFIG_HOME/wisdom/sources.json` or `~/.config/wisdom/sources.json`

---

## Usage Examples

### Basic Usage

```go
// Create engine (automatically loads from default locations)
engine := wisdom.NewEngine()
if err := engine.Initialize(); err != nil {
    log.Fatal(err)
}

// Get wisdom quote
quote, err := engine.GetWisdom(75.0, "stoic")
if err != nil {
    log.Fatal(err)
}

fmt.Println(quote.Quote)
```

### Custom Configuration Paths

```go
engine := wisdom.NewEngine()

// Configure custom paths
loader := wisdom.NewSourceLoader().
    WithConfigPaths(
        "/path/to/custom/sources.json",
        "/path/to/user/sources.json",
    )

// Initialize with custom loader
engine.loader = loader
if err := engine.Initialize(); err != nil {
    log.Fatal(err)
}
```

### Reloading Sources

```go
// Reload sources from files (useful for hot-reloading)
if err := engine.ReloadSources(); err != nil {
    log.Printf("Failed to reload: %v", err)
}
```

### Creating a New Source

```go
config := &wisdom.SourceConfig{
    ID:   "my_custom_source",
    Name: "My Custom Wisdom",
    Icon: "✨",
    Quotes: map[string][]wisdom.Quote{
        "chaos": {
            {
                Quote:        "Custom quote for chaos",
                Source:       "My Source",
                Encouragement: "Keep going!",
            },
        },
        // ... other aeon levels
    },
}

// Save to file
if err := wisdom.SaveSourceConfig("my_sources.json", config); err != nil {
    log.Fatal(err)
}
```

---

## Example Configuration File

See `examples/sources.json` for a complete example with multiple sources.

---

## Migration from Hard-Coded Sources

### Before (Hard-Coded)

```go
sources["bofh"] = &Source{
    Name:   "BOFH",
    Icon:   "😈",
    Quotes: make(map[string][]Quote),
}
sources["bofh"].Quotes["chaos"] = []Quote{
    {Quote: "...", Source: "...", Encouragement: "..."},
}
```

### After (Configurable)

```json
{
  "sources": {
    "bofh": {
      "id": "bofh",
      "name": "BOFH",
      "icon": "😈",
      "quotes": {
        "chaos": [
          {
            "quote": "...",
            "source": "...",
            "encouragement": "..."
          }
        ]
      }
    }
  }
}
```

---

## Benefits

1. **No Recompilation**: Add sources without rebuilding
2. **Easy Updates**: Update quotes via JSON files
3. **Modularity**: Split sources across multiple files
4. **Sharing**: Share source configs with others
5. **Version Control**: Track source changes in git
6. **Hot Reloading**: Reload sources at runtime
7. **Customization**: Users can add their own sources

---

## Validation

The system validates source configurations:

- Source ID is required
- Source name is required
- At least one quote must be present
- Aeon levels must be valid (`chaos`, `lower_aeons`, `middle_aeons`, `upper_aeons`, `treasury`)

Invalid configurations are rejected with descriptive error messages.

---

## Backward Compatibility

The old `GetBuiltInSources()` function is still available for backward compatibility but is deprecated. It now uses the configurable system internally with fallback to minimal hard-coded sources.

---

## Future Enhancements

- **Embedded Sources**: Compile default sources into binary using `embed`
- **Source Plugins**: Load sources from external plugins
- **API Sources**: Support for dynamic sources via API
- **Source Versioning**: Track source versions and updates
- **Source Metadata**: Additional metadata (tags, categories, etc.)

---

## See Also

- `examples/sources.json` - Example configuration
- `PYTHON_SOURCES_ANALYSIS.md` - Python source structure reference
- `internal/wisdom/sources_config.go` - Implementation details

