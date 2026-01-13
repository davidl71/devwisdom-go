# Code Duplication and Streamlining Analysis

**Date:** 2026-01-13  
**Status:** 📊 Analysis Complete  
**Codebase:** ~9,500 lines of Go code, 44 Go files  
**Analysis Method:** Manual code review + pattern matching

---

## Executive Summary

This document identifies code duplication patterns and streamlining opportunities in the devwisdom-go codebase. The analysis focuses on eliminating redundancy while maintaining code clarity and functionality.

**Key Findings:**
- **High Priority:** 3 major duplication areas (~400 lines)
- **Medium Priority:** 4 streamlining opportunities (~150 lines)
- **Low Priority:** 2 minor improvements (~50 lines)

**Total Potential Reduction:** ~600 lines of duplicated/redundant code

---

## 🔴 HIGH PRIORITY - Major Duplication

### 1. Tool Schema Definitions Duplication

**Status:** ⚠️ **ACTIVE DUPLICATION**  
**Impact:** High - violates DRY, hard to maintain  
**Effort:** Low (1-2 hours)  
**Lines Affected:** ~200 lines duplicated

**Problem:**
Tool schemas are defined in two places:
1. `internal/mcp/handlers.go:getToolDefinitions()` - Shared function (lines 268-336)
2. `internal/mcp/sdk_adapter.go:registerTools()` - Inline definitions (lines 190-284)

**Current Pattern:**
```go
// In sdk_adapter.go (lines 190-218):
consultAdvisorTool := &mcp.Tool{
    Name:        "consult_advisor",
    Description: "Consult a wisdom advisor based on metric, tool, or stage",
    InputSchema: map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "metric": map[string]interface{}{
                "type":        "string",
                "description": "Metric name (e.g., 'security', 'testing')",
            },
            // ... 50+ lines of schema definition
        },
    },
}

// In handlers.go (lines 268-297):
func getToolDefinitions() []Tool {
    return []Tool{
        {
            Name:        "consult_advisor",
            Description: "Consult a wisdom advisor based on metric, tool, or stage",
            InputSchema: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "metric": map[string]interface{}{
                        "type":        "string",
                        "description": "Metric name (e.g., 'security', 'testing')",
                    },
                    // ... same 50+ lines duplicated
                },
            },
        },
        // ... 3 more tools duplicated
    }
}
```

**Solution:**
1. Convert `getToolDefinitions()` to return `[]*mcp.Tool` instead of `[]Tool`
2. Use `getToolDefinitions()` in `registerTools()` instead of inline definitions
3. Remove duplicate schema definitions from `sdk_adapter.go`

**Refactored Code:**
```go
// internal/mcp/tools.go (new file or in handlers.go)
func GetToolDefinitions() []*mcp.Tool {
    return []*mcp.Tool{
        {
            Name:        "consult_advisor",
            Description: "Consult a wisdom advisor based on metric, tool, or stage",
            InputSchema: map[string]interface{}{
                // ... single source of truth
            },
        },
        // ... other tools
    }
}

// In sdk_adapter.go:
func (s *WisdomServerSDK) registerTools() error {
    handlers := NewWisdomHandlers(s.wisdom, s.logger, s.appLogger)
    
    tools := GetToolDefinitions()
    for _, tool := range tools {
        // Map tool name to handler
        var handlerFunc ToolHandlerFunc
        switch tool.Name {
        case "consult_advisor":
            handlerFunc = handlers.handleConsultAdvisor
        case "get_wisdom":
            handlerFunc = handlers.handleGetWisdom
        case "get_daily_briefing":
            handlerFunc = handlers.handleGetDailyBriefing
        case "get_consultation_log":
            handlerFunc = handlers.handleGetConsultationLog
        default:
            return fmt.Errorf("unknown tool: %s", tool.Name)
        }
        
        s.server.AddTool(tool, wrapToolHandler(handlerFunc))
    }
    
    return nil
}
```

**Benefits:**
- Single source of truth for tool definitions
- Easier to add new tools (one place to update)
- Reduced maintenance burden
- Consistent tool definitions

---

### 2. CLI Test Output Capture Pattern Duplication

**Status:** ⚠️ **ACTIVE DUPLICATION**  
**Impact:** Medium-High - affects test maintainability  
**Effort:** Low (1 hour)  
**Lines Affected:** ~120 lines duplicated (30 lines × 4 test files)

**Problem:**
All CLI test files (`quote_test.go`, `consult_test.go`, `briefing_test.go`, `sources_test.go`) duplicate the same stdout capture pattern:

**Current Pattern (repeated 4 times):**
```go
// In each test file:
t.Run(tt.name, func(t *testing.T) {
    // Capture output using a pipe.
    r, w, _ := os.Pipe()
    oldStdout := os.Stdout
    os.Stdout = w

    var buf bytes.Buffer
    done := make(chan bool)
    go func() {
        _, err := buf.ReadFrom(r)
        if err != nil {
            t.Errorf("buf.ReadFrom failed: %v", err)
        }
        done <- true
    }()

    err := app.runXXX(tt.args)

    w.Close()
    os.Stdout = oldStdout
    <-done

    if (err != nil) != tt.wantErr {
        t.Errorf("runXXX() error = %v, wantErr %v", err, tt.wantErr)
        return
    }

    if !tt.wantErr && tt.check != nil {
        output := buf.String()
        if !tt.check(output) {
            t.Errorf("runXXX() output validation failed. Output: %s", output)
        }
    }
})
```

**Solution:**
Extract to a test helper function:

**Refactored Code:**
```go
// internal/cli/test_helpers.go (new file)
package cli

import (
    "bytes"
    "os"
    "testing"
)

// captureOutput captures stdout from a function call.
// Returns the captured output and any error.
func captureOutput(fn func() error) (string, error) {
    r, w, _ := os.Pipe()
    oldStdout := os.Stdout
    os.Stdout = w

    var buf bytes.Buffer
    done := make(chan bool)
    go func() {
        _, _ = buf.ReadFrom(r)
        done <- true
    }()

    err := fn()

    w.Close()
    os.Stdout = oldStdout
    <-done

    return buf.String(), err
}

// runCLITest runs a CLI test with output capture and validation.
func runCLITest(t *testing.T, name string, fn func() error, wantErr bool, check func(string) bool) {
    t.Run(name, func(t *testing.T) {
        output, err := captureOutput(fn)
        
        if (err != nil) != wantErr {
            t.Errorf("error = %v, wantErr %v", err, wantErr)
            return
        }

        if !wantErr && check != nil {
            if !check(output) {
                t.Errorf("output validation failed. Output: %s", output)
            }
        }
    })
}
```

**Usage:**
```go
// In quote_test.go:
runCLITest(t, "quote with json flag", 
    func() error { return app.runQuote([]string{"--json"}) },
    false,
    validateJSONOutput,
)
```

**Benefits:**
- Eliminates ~120 lines of duplicated test code
- Easier to maintain test infrastructure
- Consistent test pattern across all CLI tests
- Less error-prone (single implementation)

---

### 3. Resource Template Handler Pattern Duplication

**Status:** ⚠️ **ACTIVE DUPLICATION**  
**Impact:** Medium - affects resource handler maintainability  
**Effort:** Medium (2-3 hours)  
**Lines Affected:** ~80 lines duplicated

**Problem:**
Resource template handlers for `wisdom://advisor/{id}` and `wisdom://consultations/{days}` follow the same pattern:

**Current Pattern (duplicated in sdk_adapter.go lines 348-405):**
```go
// Pattern 1: Advisor template handler
advisorTemplateHandler := func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
    if req.Params == nil || req.Params.URI == "" {
        return nil, fmt.Errorf("resource URI is required")
    }
    uri := req.Params.URI
    
    // Extract advisor ID from URI
    if !strings.HasPrefix(uri, "wisdom://advisor/") {
        return nil, fmt.Errorf("invalid advisor URI format: %s", uri)
    }
    advisorID := strings.TrimPrefix(uri, "wisdom://advisor/")
    
    mockReq := &JSONRPCRequest{
        ID:     "resource",
        Method: "resources/read",
        Params: json.RawMessage(fmt.Sprintf(`{"uri": "%s"}`, uri)),
    }
    
    resp := handlers.HandleAdvisorResource(mockReq, advisorID)
    return s.convertResourceResponse(resp, uri)
}

// Pattern 2: Consultations template handler (same pattern, different extraction)
consultationsTemplateHandler := func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
    if req.Params == nil || req.Params.URI == "" {
        return nil, fmt.Errorf("resource URI is required")
    }
    uri := req.Params.URI
    
    // Extract days from URI
    if !strings.HasPrefix(uri, "wisdom://consultations/") {
        return nil, fmt.Errorf("invalid consultations URI format: %s", uri)
    }
    daysStr := strings.TrimPrefix(uri, "wisdom://consultations/")
    days := 7 // default
    if daysStr != "" {
        if d, err := strconv.Atoi(daysStr); err == nil {
            days = d
        }
    }
    
    mockReq := &JSONRPCRequest{
        ID:     "resource",
        Method: "resources/read",
        Params: json.RawMessage(fmt.Sprintf(`{"uri": "%s"}`, uri)),
    }
    
    resp := handlers.HandleConsultationsResource(mockReq, days)
    return s.convertResourceResponse(resp, uri)
}
```

**Solution:**
Create a generic resource template handler factory:

**Refactored Code:**
```go
// In sdk_adapter.go:
type ResourceTemplateHandler func(*JSONRPCRequest, string) *JSONRPCResponse

// createResourceTemplateHandler creates a handler for resource templates.
func (s *WisdomServerSDK) createResourceTemplateHandler(
    prefix string,
    extractParam func(string) (string, error),
    handlerFunc ResourceTemplateHandler,
) func(context.Context, *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
    return func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
        if req.Params == nil || req.Params.URI == "" {
            return nil, fmt.Errorf("resource URI is required")
        }
        uri := req.Params.URI
        
        if !strings.HasPrefix(uri, prefix) {
            return nil, fmt.Errorf("invalid URI format: %s (expected prefix: %s)", uri, prefix)
        }
        
        param, err := extractParam(uri)
        if err != nil {
            return nil, fmt.Errorf("failed to extract parameter from URI %s: %w", uri, err)
        }
        
        mockReq := &JSONRPCRequest{
            ID:     "resource",
            Method: "resources/read",
            Params: json.RawMessage(fmt.Sprintf(`{"uri": "%s"}`, uri)),
        }
        
        resp := handlerFunc(mockReq, param)
        return s.convertResourceResponse(resp, uri)
    }
}

// Usage:
advisorTemplateHandler := s.createResourceTemplateHandler(
    "wisdom://advisor/",
    func(uri string) (string, error) {
        return strings.TrimPrefix(uri, "wisdom://advisor/"), nil
    },
    handlers.HandleAdvisorResource,
)

consultationsTemplateHandler := s.createResourceTemplateHandler(
    "wisdom://consultations/",
    func(uri string) (string, error) {
        daysStr := strings.TrimPrefix(uri, "wisdom://consultations/")
        days := 7
        if daysStr != "" {
            var err error
            days, err = strconv.Atoi(daysStr)
            if err != nil {
                return "", fmt.Errorf("invalid days parameter: %w", err)
            }
        }
        return strconv.Itoa(days), nil
    },
    func(req *JSONRPCRequest, daysStr string) *JSONRPCResponse {
        days, _ := strconv.Atoi(daysStr)
        return handlers.HandleConsultationsResource(req, days)
    },
)
```

**Benefits:**
- Eliminates ~80 lines of duplicated handler code
- Consistent error handling pattern
- Easier to add new template resources
- Single place to fix bugs

---

## 🟡 MEDIUM PRIORITY - Streamlining Opportunities

### 4. Parameter Extraction Pattern Repetition

**Status:** ⚠️ **REPETITIVE PATTERN**  
**Impact:** Medium - affects handler maintainability  
**Effort:** Medium (2-3 hours)  
**Lines Affected:** ~100 lines

**Problem:**
All tool handlers repeat the same parameter extraction pattern:

**Current Pattern (repeated in handlers.go):**
```go
// In handleConsultAdvisor (lines 54-73):
var metric, tool, stage, context string
var score float64

if m, ok := params["metric"].(string); ok {
    metric = m
}
if t, ok := params["tool"].(string); ok {
    tool = t
}
if st, ok := params["stage"].(string); ok {
    stage = st
}
if c, ok := params["context"].(string); ok {
    context = c
}
if sc, ok := params["score"].(float64); ok {
    score = sc
} else if sc, ok := params["score"].(int); ok {
    score = float64(sc)
}
```

**Solution:**
Create parameter extraction helpers:

**Refactored Code:**
```go
// internal/mcp/params.go (new file)
package mcp

// GetString extracts a string parameter with optional default.
func GetString(params map[string]interface{}, key, defaultValue string) string {
    if v, ok := params[key].(string); ok && v != "" {
        return v
    }
    return defaultValue
}

// GetFloat64 extracts a float64 parameter with optional default.
// Handles both float64 and int types (JSON unmarshaling).
func GetFloat64(params map[string]interface{}, key string, defaultValue float64) float64 {
    if v, ok := params[key].(float64); ok {
        return v
    }
    if v, ok := params[key].(int); ok {
        return float64(v)
    }
    return defaultValue
}

// GetInt extracts an int parameter with optional default.
func GetInt(params map[string]interface{}, key string, defaultValue int) int {
    if v, ok := params[key].(int); ok {
        return v
    }
    if v, ok := params[key].(float64); ok {
        return int(v)
    }
    return defaultValue
}

// ClampScore clamps a score to 0-100 range.
func ClampScore(score float64) float64 {
    if score < 0 {
        return 0
    }
    if score > 100 {
        return 100
    }
    return score
}
```

**Usage:**
```go
// In handlers.go:
func (h *WisdomHandlers) handleConsultAdvisor(params map[string]interface{}) (interface{}, error) {
    metric := GetString(params, "metric", "")
    tool := GetString(params, "tool", "")
    stage := GetString(params, "stage", "")
    context := GetString(params, "context", "")
    score := ClampScore(GetFloat64(params, "score", 50))
    
    // ... rest of handler
}
```

**Benefits:**
- Eliminates ~100 lines of repetitive parameter extraction
- Consistent parameter handling
- Type-safe extraction with defaults
- Easier to add new parameters

---

### 5. Resource List Definition Duplication

**Status:** ⚠️ **ACTIVE DUPLICATION**  
**Impact:** Low-Medium - affects resource list consistency  
**Effort:** Low (30 minutes)  
**Lines Affected:** ~40 lines duplicated

**Problem:**
Resource list is defined in two places:
1. `internal/mcp/server.go:handleResourcesList()` (lines 278-309)
2. `internal/mcp/sdk_adapter.go:registerResources()` (implicitly through resource registration)

**Solution:**
Extract to shared constant or function:

**Refactored Code:**
```go
// internal/mcp/resources.go (new file)
package mcp

// GetResourceDefinitions returns the list of available resources.
func GetResourceDefinitions() []Resource {
    return []Resource{
        {
            URI:         "wisdom://tools",
            Name:        "Available Tools",
            Description: "List all available MCP tools with descriptions and parameters",
            MimeType:    "application/json",
        },
        // ... other resources
    }
}
```

**Benefits:**
- Single source of truth for resource definitions
- Consistent resource lists across implementations
- Easier to add new resources

---

### 6. Error Result Construction Pattern

**Status:** ✅ **PARTIALLY ADDRESSED**  
**Impact:** Low - already has helpers  
**Effort:** Low (30 minutes)  
**Lines Affected:** ~20 lines

**Current State:**
- `newErrorResult()` and `newSuccessResult()` already exist in `sdk_adapter.go`
- Some error construction still uses inline patterns

**Recommendation:**
- Review all error construction sites
- Ensure all use `newErrorResult()` helper
- Consider adding more specific error helpers (e.g., `newInvalidParamsError()`)

---

### 7. JSON Marshaling Pattern

**Status:** ✅ **MOSTLY ADDRESSED**  
**Impact:** Low - already has helpers  
**Effort:** Low (15 minutes)  
**Lines Affected:** ~10 lines

**Current State:**
- `mustMarshalJSONCompact()` exists in `handlers.go`
- Some JSON marshaling still uses inline `json.Marshal()`

**Recommendation:**
- Review all JSON marshaling sites
- Use `mustMarshalJSONCompact()` where appropriate
- Consider adding `mustMarshalJSONIndent()` for pretty-printed JSON

---

## 🟢 LOW PRIORITY - Minor Improvements

### 8. Test Setup Pattern

**Status:** ℹ️ **ACCEPTABLE PATTERN**  
**Impact:** Low - standard Go test pattern  
**Effort:** N/A  
**Lines Affected:** N/A

**Current State:**
- Test setup/teardown patterns are similar but acceptable
- No significant duplication

**Recommendation:**
- Keep as-is unless pattern becomes more complex
- Consider test helpers only if pattern evolves

---

### 9. Comment Style Inconsistency

**Status:** ℹ️ **STYLE ISSUE**  
**Impact:** Very Low - cosmetic  
**Effort:** Low (30 minutes)  
**Lines Affected:** ~50 lines

**Current State:**
- Some comments use periods, some don't
- Some use full sentences, some use fragments

**Recommendation:**
- Standardize comment style (prefer full sentences with periods)
- Run `gofmt` and `goimports` for consistency
- Consider adding golangci-lint rule for comment style

---

## Implementation Priority

### Phase 1: High Priority (Immediate - 4-6 hours)
1. ✅ Extract tool schema definitions (1-2 hours)
2. ✅ Extract CLI test helpers (1 hour)
3. ✅ Extract resource template handler pattern (2-3 hours)

**Total Phase 1:** ~400 lines eliminated

### Phase 2: Medium Priority (Soon - 3-4 hours)
4. ⏳ Create parameter extraction helpers (2-3 hours)
5. ⏳ Extract resource list definitions (30 minutes)
6. ⏳ Standardize error result construction (30 minutes)

**Total Phase 2:** ~150 lines improved

### Phase 3: Low Priority (Backlog - 1 hour)
7. ⏳ Review JSON marshaling patterns (15 minutes)
8. ⏳ Standardize comment style (30 minutes)

**Total Phase 3:** ~60 lines improved

---

## Success Metrics

### Code Quality
- ✅ Duplication reduced by 80%+
- ✅ Code lines reduced by ~600 lines
- ✅ Test coverage maintained or improved
- ✅ Builds succeed
- ✅ No performance regressions

### Maintainability
- ✅ Single source of truth for tool/resource definitions
- ✅ Centralized test helpers
- ✅ Consistent parameter handling
- ✅ Easier to add new tools/resources/tests

---

## Comparison with Previous Analysis

### Already Addressed (from REFACTORING_PLAN.md)
- ✅ `wrapToolHandler()` - Already implemented
- ✅ `newErrorResult()` / `newSuccessResult()` - Already implemented
- ✅ `newResourceResponse()` - Already implemented
- ✅ `getToolDefinitions()` - Already exists (but not fully utilized)

### New Findings
- ⚠️ Tool schema duplication (not fully addressed)
- ⚠️ CLI test output capture duplication
- ⚠️ Resource template handler duplication
- ⚠️ Parameter extraction pattern repetition

---

## Next Steps

### Immediate Actions
1. Review and approve this analysis
2. Start with Phase 1 (High Priority items)
3. Extract tool schema definitions first (biggest win)

### Short Term
1. Complete Phase 2 refactoring
2. Add golangci-lint config for ongoing duplication detection
3. Document refactoring patterns for future reference

### Long Term
1. Monitor for new duplication patterns
2. Consider applying exarp-go patterns where applicable
3. Regular code review for duplication

---

## References

- **Previous Analysis:** `docs/REFACTORING_PLAN.md`
- **Code Migration:** `docs/CODE_MIGRATION_OPPORTUNITIES.md`
- **Integration Status:** `docs/MCP_GO_CORE_INTEGRATION_STATUS.md`
- **Patterns:** `.cursor/rules/exarp-go-patterns.mdc`

---

**Last Updated:** 2026-01-13  
**Next Review:** After Phase 1 completion
