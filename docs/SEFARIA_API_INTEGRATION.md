# Sefaria API Integration Plan

**Date**: 2026-01-06  
**Status**: Research Complete, Ready for Implementation  
**Task**: T-25 (Phase 7.1)

---

## Overview

This document outlines the integration plan for Sefaria API to fetch Hebrew text sources for the devwisdom-go project. The Sefaria API provides access to Hebrew texts including Pirkei Avot, Proverbs, Ecclesiastes, and Psalms.

---

## API Documentation

### Base URL
- **Production**: `https://www.sefaria.org/api`
- **Protocol**: HTTPS
- **Authentication**: None required for public text access
- **CORS**: Enabled (`access-control-allow-origin: *`)

### Endpoint Patterns

1. **Full Book Access**:
   ```
   GET /api/texts/{Book}
   Example: /api/texts/Proverbs
   ```

2. **Chapter Access**:
   ```
   GET /api/texts/{Book}.{Chapter}
   Example: /api/texts/Pirkei_Avot.1
   ```

3. **Verse Access**:
   ```
   GET /api/texts/{Book}.{Chapter}.{Verse}
   Example: /api/texts/Pirkei_Avot.1.1
   ```

4. **Alternative Format** (with query parameter):
   ```
   GET /api?vv2=texts/{Book}.{Chapter}.{Verse}
   Example: /api?vv2=texts/Pirkei_Avot.1.1
   ```

### Response Format

**Example Response Structure** (from `/api/texts/Proverbs`):
```json
{
  "ref": "Proverbs 1",
  "heRef": "משלי א׳",
  "text": ["English verse 1", "English verse 2", ...],
  "he": ["Hebrew verse 1", "Hebrew verse 2", ...],
  "versions": [
    {
      "versionTitle": "THE JPS TANAKH: Gender-Sensitive Edition",
      "language": "en",
      "text": [...]
    },
    {
      "versionTitle": "Miqra according to the Masorah",
      "language": "he",
      "text": [...]
    }
  ],
  "metadata": {
    "book": "Proverbs",
    "heTitle": "משלי",
    "categories": ["Tanakh", "Writings"],
    ...
  }
}
```

**Key Fields**:
- `text[]`: Array of English verses
- `he[]`: Array of Hebrew verses
- `ref`: English reference (e.g., "Proverbs 1")
- `heRef`: Hebrew reference (e.g., "משלי א׳")
- `versions[]`: Multiple translation versions available

---

## Required Hebrew Sources

The following sources need API integration:

1. **pirkei_avot** (Pirkei Avot / Ethics of the Fathers)
   - Sefaria ID: `Pirkei_Avot`
   - Chapters: 6
   - Language: Hebrew + English

2. **proverbs** (Book of Proverbs / Mishlei)
   - Sefaria ID: `Proverbs`
   - Chapters: 31
   - Language: Hebrew + English

3. **ecclesiastes** (Kohelet / Ecclesiastes)
   - Sefaria ID: `Ecclesiastes`
   - Chapters: 12
   - Language: Hebrew + English

4. **psalms** (Tehillim / Psalms)
   - Sefaria ID: `Psalms`
   - Chapters: 150
   - Language: Hebrew + English

---

## Rate Limiting & Usage

- **Rate Limits**: Not explicitly documented, but Cloudflare protection is present
- **Recommendation**: 
  - Implement request throttling (max 1 request/second)
  - Use caching aggressively (24-hour TTL)
  - Implement exponential backoff on errors

---

## Integration Architecture

### Package Structure

```
internal/wisdom/sefaria/
├── client.go      # HTTP client wrapper
├── types.go       # Response data structures
├── cache.go       # Response caching
└── mapper.go      # Map API responses to Quote format
```

### Client Implementation

**File**: `internal/wisdom/sefaria/client.go`

```go
type Client struct {
    httpClient *http.Client
    baseURL    string
    cache      *Cache
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
        baseURL: "https://www.sefaria.org/api",
        cache:   NewCache(),
    }
}

func (c *Client) GetText(book string, chapter, verse int) (*TextResponse, error) {
    // Check cache first
    // Make API request
    // Parse response
    // Cache result
    // Return mapped Quote
}
```

### Data Structures

**File**: `internal/wisdom/sefaria/types.go`

```go
type TextResponse struct {
    Ref     string   `json:"ref"`
    HeRef   string   `json:"heRef"`
    Text    []string `json:"text"`
    He      []string `json:"he"`
    Versions []Version `json:"versions"`
}

type Version struct {
    VersionTitle string   `json:"versionTitle"`
    Language     string   `json:"language"`
    Text         []string `json:"text"`
}
```

### Caching Strategy

**File**: `internal/wisdom/sefaria/cache.go`

- **Cache Key Format**: `sefaria:{book}:{chapter}:{verse}`
- **TTL**: 24 hours (Hebrew texts don't change)
- **Storage**: In-memory with optional disk persistence
- **Cache Invalidation**: Manual refresh or TTL expiration

---

## Error Handling Strategy

### Network Failures
- **Action**: Fallback to embedded sources (if available)
- **Logging**: Log error with context
- **Retry**: Exponential backoff (3 attempts max)

### Rate Limiting
- **Detection**: HTTP 429 status code
- **Action**: Exponential backoff, wait before retry
- **Fallback**: Use cached data if available

### Invalid Responses
- **Detection**: JSON parse errors, missing required fields
- **Action**: Log error, use fallback
- **Validation**: Verify `text` and `he` arrays are present

### Timeout
- **Timeout**: 10 seconds (already configured in SourceLoader)
- **Action**: Return error, use fallback
- **Logging**: Log timeout with request details

---

## Testing Strategy

### Unit Tests
- Mock HTTP responses
- Test response parsing
- Test caching behavior
- Test error handling scenarios

### Integration Tests
- Test with real API (rate-limited)
- Test offline fallback
- Test timeout scenarios
- Test with invalid book/chapter/verse

### Test Data
- Use mock responses for unit tests
- Create test fixtures for common scenarios
- Test with actual API in CI (limited runs)

---

## Implementation Steps

1. **Create Sefaria Package** (`internal/wisdom/sefaria/`)
   - Implement HTTP client wrapper
   - Define response data structures
   - Implement caching layer

2. **Integrate with SourceLoader**
   - Detect `sefaria_source` field in SourceConfig
   - Call Sefaria API when source is API-based
   - Map API response to Quote format

3. **Update Engine**
   - Remove hardcoded exclusions for Sefaria sources
   - Enable Sefaria sources in random selection
   - Handle API failures gracefully

4. **Add Tests**
   - Unit tests with mocked responses
   - Integration tests with real API
   - Error handling tests

---

## Dependencies

- **Go Standard Library**: `net/http` (already in use)
- **No External Dependencies**: Use standard library only
- **Existing Infrastructure**: 
  - HTTP client in `SourceLoader` (can be reused)
  - Caching system (`SourceCache`) can be extended

---

## Security Considerations

- **HTTPS Only**: All API calls use HTTPS
- **No Authentication**: Public API, no credentials needed
- **Input Validation**: Validate book/chapter/verse parameters
- **Rate Limiting**: Implement client-side throttling
- **Error Messages**: Don't expose internal errors to users

---

## Performance Considerations

- **Caching**: Aggressive caching (24-hour TTL)
- **Concurrent Requests**: Limit concurrent API calls
- **Timeout**: 10-second timeout prevents hanging
- **Connection Pooling**: Reuse HTTP connections

---

## Future Enhancements

- Support for additional Sefaria sources
- Offline mode with pre-cached data
- Background refresh of cached data
- Support for commentary and translations

---

## References

- **Sefaria Website**: https://www.sefaria.org
- **API Base URL**: https://www.sefaria.org/api
- **Tested Endpoint**: `/api/texts/Proverbs` (verified working)
- **Response Format**: JSON with `text[]` and `he[]` arrays

---

## Next Steps

1. ✅ Research complete (T-25)
2. ⏳ Implement API client (T-26)
3. ⏳ Integrate with source loading
4. ⏳ Add tests
5. ⏳ Update documentation

