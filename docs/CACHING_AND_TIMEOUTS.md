# Caching and Timeout Configuration

**Guide**: How caching and timeouts work in the wisdom source system

---

## Overview

The wisdom system includes built-in caching and timeout support for:

- **File-based sources** - Cached to avoid repeated file I/O
- **API-based sources** - Timeout protection for remote requests
- **Automatic invalidation** - Cache cleared when files change
- **Configurable TTL** - Customizable cache duration

---

## Caching

### How It Works

1. **File Sources**: When a source is loaded from a JSON file, it's cached in memory
2. **Cache Key**: Based on file path and source ID
3. **TTL**: Default 5 minutes (configurable)
4. **Max Age**: Default 1 hour (configurable)
5. **File Watching**: Cache invalidated when file modification time changes

### Default Settings

```go
Default TTL:     5 minutes
Max Age:         1 hour
Cleanup Interval: 5 minutes (background cleanup)
```

### Configuration

```go
loader := wisdom.NewSourceLoader().
    WithCacheTTL(10 * time.Minute).      // Set cache TTL
    WithCacheMaxAge(2 * time.Hour).      // Set max cache age
    WithCacheEnabled(true)                // Enable/disable caching
```

### Cache Invalidation

**Automatic:**
- File modification time changes
- Cache entry exceeds TTL
- Cache entry exceeds max age
- Background cleanup runs every 5 minutes

**Manual:**
```go
loader.InvalidateCache()              // Clear all cache
loader.cache.Invalidate("source_id")  // Clear specific source
```

### Cache Benefits

- ✅ **Performance**: Avoid repeated file reads
- ✅ **Efficiency**: Reduced I/O operations
- ✅ **Smart Invalidation**: Auto-detects file changes
- ✅ **Memory Safe**: Automatic cleanup of expired entries

---

## Timeouts

### HTTP Client Timeout

For API-based sources (like Sefaria), the system uses an HTTP client with timeout protection:

```go
Default Timeout: 10 seconds
```

### Configuration

```go
loader := wisdom.NewSourceLoader().
    WithHTTPTimeout(30 * time.Second)  // Set API timeout
```

### API Source Loading

```go
apiLoader := wisdom.NewAPISourceLoader(
    "https://api.example.com",
    10 * time.Second,  // Timeout
)

// Load with timeout
config, err := apiLoader.LoadSourceWithTimeout("endpoint", 5*time.Second)

// Load with context (for cancellation)
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
config, err := apiLoader.LoadSource(ctx, "endpoint")

// Load with retry
config, err := apiLoader.LoadSourceWithRetry(ctx, "endpoint", 3)
```

---

## Usage Examples

### Basic Usage (Default Caching)

```go
engine := wisdom.NewEngine()
engine.Initialize()

// Sources are automatically cached
// File changes are detected automatically
```

### Custom Cache Settings

```go
loader := wisdom.NewSourceLoader().
    WithCacheTTL(15 * time.Minute).      // Longer cache
    WithCacheMaxAge(4 * time.Hour)       // Longer max age

engine.loader = loader
engine.Initialize()
```

### Disable Caching

```go
loader := wisdom.NewSourceLoader().
    WithCacheEnabled(false)  // Disable caching

engine.loader = loader
engine.Initialize()
```

### API Source with Timeout

```go
// Create API loader
apiLoader := wisdom.NewAPISourceLoader(
    "https://api.sefaria.org",
    10 * time.Second,  // 10 second timeout
)

// Load source with custom timeout
config, err := apiLoader.LoadSourceWithTimeout(
    "texts/Pirkei_Avot",
    5 * time.Second,  // 5 second timeout for this call
)
```

### API Source with Retry

```go
ctx := context.Background()
config, err := apiLoader.LoadSourceWithRetry(
    ctx,
    "texts/Pirkei_Avot",
    3,  // Max 3 retries with exponential backoff
)
```

---

## Cache Statistics

### Get Cache Size

```go
size := loader.cache.Size()
fmt.Printf("Cached sources: %d\n", size)
```

### Clear Expired Entries

```go
cleared := loader.cache.ClearExpired()
fmt.Printf("Cleared %d expired entries\n", cleared)
```

---

## File Change Detection

The cache automatically detects when source files are modified:

1. **File Modification Time**: Tracked when file is cached
2. **On Access**: Check modification time when retrieving from cache
3. **Auto-Invalidate**: If file changed, cache entry is removed
4. **Reload**: Next access will reload from file

### Example

```go
// File cached at 10:00:00
loader.Load()  // Loads sources.json, caches it

// File modified at 10:05:00
// ... user edits sources.json ...

// Next access at 10:06:00
loader.Load()  // Detects file change, invalidates cache, reloads
```

---

## Performance Considerations

### Cache Hit vs Miss

- **Cache Hit**: ~0.001ms (memory lookup)
- **Cache Miss**: ~1-5ms (file read + parse)

### When to Use Caching

✅ **Use Caching When:**
- Sources are loaded frequently
- Files don't change often
- Performance is important

❌ **Disable Caching When:**
- Files change very frequently
- Memory is constrained
- Always want fresh data

### Memory Usage

- **Per Source**: ~1-5 KB (depends on quote count)
- **100 Sources**: ~100-500 KB
- **Automatic Cleanup**: Prevents memory leaks

---

## Best Practices

### 1. Production Settings

```go
loader := wisdom.NewSourceLoader().
    WithCacheTTL(10 * time.Minute).      // Reasonable TTL
    WithCacheMaxAge(1 * time.Hour).       // Prevent stale data
    WithHTTPTimeout(10 * time.Second)     // Protect against slow APIs
```

### 2. Development Settings

```go
loader := wisdom.NewSourceLoader().
    WithCacheTTL(1 * time.Minute).        // Shorter TTL for testing
    WithCacheEnabled(true)                // Still cache for performance
```

### 3. Testing Settings

```go
loader := wisdom.NewSourceLoader().
    WithCacheEnabled(false)                // No cache for predictable tests
```

### 4. API Sources

```go
// Always use timeouts
apiLoader := wisdom.NewAPISourceLoader(
    baseURL,
    10 * time.Second,  // Reasonable timeout
)

// Use retry for unreliable APIs
config, err := apiLoader.LoadSourceWithRetry(ctx, endpoint, 3)
```

---

## Troubleshooting

### Cache Not Updating

**Problem**: Changes to source file not reflected

**Solutions**:
1. Check file modification time is changing
2. Manually invalidate cache: `loader.InvalidateCache()`
3. Reload: `engine.ReloadSources()`
4. Reduce TTL: `WithCacheTTL(1 * time.Minute)`

### Timeout Errors

**Problem**: API calls timing out

**Solutions**:
1. Increase timeout: `WithHTTPTimeout(30 * time.Second)`
2. Use retry: `LoadSourceWithRetry(ctx, endpoint, 3)`
3. Check network connectivity
4. Verify API endpoint is accessible

### Memory Usage

**Problem**: High memory usage

**Solutions**:
1. Reduce max age: `WithCacheMaxAge(30 * time.Minute)`
2. Disable cache: `WithCacheEnabled(false)`
3. Clear expired: `loader.cache.ClearExpired()`
4. Manual cleanup: `loader.InvalidateCache()`

---

## Configuration Summary

| Setting | Default | Description |
|---------|---------|-------------|
| `Cache TTL` | 5 minutes | How long entries stay valid |
| `Max Age` | 1 hour | Maximum cache age |
| `HTTP Timeout` | 10 seconds | API request timeout |
| `Cleanup Interval` | 5 minutes | Background cleanup frequency |
| `Cache Enabled` | true | Enable/disable caching |

---

## See Also

- `CONFIGURABLE_SOURCES.md` - Source configuration
- `ADDING_PROJECT_SOURCES.md` - Adding custom sources
- `sources_cache.go` - Cache implementation
- `sources_api.go` - API loader implementation

