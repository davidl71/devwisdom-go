# DevWisdom-Go Refactoring Plan

**Date:** 2026-01-12  
**Status:** 📋 Planning  
**Codebase:** ~9,500 lines of Go code, 44 Go files

---

## Executive Summary

This document outlines identified code duplication and refactoring opportunities for the devwisdom-go project. The plan prioritizes high-impact refactoring to improve maintainability and code quality.

**Analysis Tool:** `golangci-lint --enable=dupl`  
**Total Duplication Found:** ~300-400 lines across 3 major areas

---

## Identified Duplication Patterns

### 🔴 HIGH PRIORITY

#### 1. CLI Test Helper Duplication
- **Location:** `internal/cli/*_test.go` (4 files)
  - `briefing_test.go:37-93` 
  - `quote_test.go:28-84`
  - `sources_test.go:29-85`
  - `consult_test.go:48-99`
- **Issue:** JSON validation helper logic duplicated across all CLI test files
- **Duplication:** ~50 lines repeated 4 times (~200 lines total)
- **Impact:** High - affects test maintainability
- **Solution:** Extract JSON validation helper into `internal/cli/test_helpers.go`
- **Complexity:** Low
- **Estimated Effort:** 1-2 hours

**Pattern:**
```go
// Repeated in all test files:
check: func(output string) bool {
    // Extract JSON from output (may have warnings before it)
    lines := strings.Split(output, "\n")
    var jsonLines []string
    inJSON := false
    braceCount := 0
    bracketCount := 0
    // ... 40+ lines of identical JSON parsing logic
}
```

**Refactor to:**
```go
// internal/cli/test_helpers.go
func extractJSONFromOutput(output string) string
func validateJSONOutput(output string) bool
```

---

#### 2. SDK Adapter Tool Handler Duplication
- **Location:** `internal/mcp/sdk_adapter.go`
  - Lines 118-167 (consultAdvisorHandler)
  - Lines 191-237 (getWisdomHandler)
  - Lines 256-302 (getDailyBriefingHandler)
  - Lines 321-367 (getConsultationLogHandler)
- **Issue:** All 4 tool handlers follow identical pattern:
  1. Unmarshal arguments
  2. Call handler function
  3. Marshal result
  4. Return CallToolResult
- **Duplication:** ~50 lines repeated 4 times (~200 lines total, ~150 duplicated)
- **Impact:** Very High - affects core MCP tool functionality
- **Solution:** Create generic tool handler wrapper function
- **Complexity:** Medium
- **Estimated Effort:** 2-3 hours

**Pattern:**
```go
// Repeated 4 times:
handler := func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := make(map[string]interface{})
    if req.Params != nil && len(req.Params.Arguments) > 0 {
        if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
            return errorResult(err), nil
        }
    }
    result, err := handlers.handleXXX(args)
    if err != nil {
        return errorResult(err), nil
    }
    resultJSON, err := json.Marshal(result)
    if err != nil {
        return errorResult(err), nil
    }
    return successResult(resultJSON), nil
}
```

**Refactor to:**
```go
// internal/mcp/sdk_adapter.go
type ToolHandlerFunc func(map[string]interface{}) (interface{}, error)

func wrapToolHandler(handler ToolHandlerFunc) func(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Shared unmarshal/marshal/error handling logic
    }
}
```

---

#### 3. Tool List Definition Duplication
- **Location:** 
  - `internal/mcp/handlers.go:261-336` (HandleToolsResource)
  - `internal/mcp/server.go:202-277` (handleToolsList)
- **Issue:** Identical tool list definitions in two places
- **Duplication:** ~75 lines duplicated exactly
- **Impact:** High - violates DRY principle, hard to maintain
- **Solution:** Extract tool definitions into shared constant or function
- **Complexity:** Low
- **Estimated Effort:** 1 hour

**Pattern:**
```go
// Duplicated in handlers.go and server.go:
tools := []Tool{
    {
        Name:        "consult_advisor",
        Description: "...",
        InputSchema: map[string]interface{}{...},
    },
    // ... same 4 tools defined twice
}
```

**Refactor to:**
```go
// internal/mcp/tools.go (new file)
func GetToolDefinitions() []Tool {
    return []Tool{
        // Single source of truth
    }
}
```

---

### 🟡 MEDIUM PRIORITY

#### 4. Error Result Construction Pattern
- **Location:** `internal/mcp/sdk_adapter.go`, `internal/mcp/handlers.go`
- **Issue:** Error result construction pattern repeated multiple times
- **Duplication:** ~10 lines repeated ~8+ times
- **Impact:** Moderate - affects consistency
- **Solution:** Create helper functions for error/success result construction
- **Complexity:** Low
- **Estimated Effort:** 30 minutes

**Pattern:**
```go
// Repeated multiple times:
return &mcp.CallToolResult{
    IsError: true,
    Content: []mcp.Content{
        &mcp.TextContent{
            Text: fmt.Sprintf("Error message: %v", err),
        },
    },
}, nil
```

**Refactor to:**
```go
func newErrorResult(err error) *mcp.CallToolResult
func newSuccessResult(text string) *mcp.CallToolResult
```

---

#### 5. JSON Schema Definition Patterns
- **Location:** `internal/mcp/sdk_adapter.go`, `internal/mcp/handlers.go`
- **Issue:** Similar patterns for building JSON schemas
- **Duplication:** Moderate - similar structures, not identical
- **Impact:** Low-Medium - affects maintainability
- **Solution:** Create schema builder helpers (optional)
- **Complexity:** Medium
- **Estimated Effort:** 2-3 hours

---

### 🟢 LOW PRIORITY

#### 6. Test Setup/Teardown Patterns
- **Location:** Test files across `internal/wisdom/`, `internal/mcp/`, `internal/cli/`
- **Issue:** Similar test setup/teardown code
- **Duplication:** Low - acceptable test pattern
- **Impact:** Low
- **Solution:** Consider test helpers if pattern becomes more complex
- **Complexity:** Low
- **Estimated Effort:** 1-2 hours (if needed)

---

## Refactoring Strategy

### Phase 1: High Priority (Immediate)
1. ✅ Extract CLI test helpers (1-2 hours)
2. ✅ Extract tool handler wrapper in SDK adapter (2-3 hours)
3. ✅ Consolidate tool list definitions (1 hour)

**Total Phase 1:** ~4-6 hours, ~300 lines eliminated

### Phase 2: Medium Priority (Soon)
4. ⏳ Create error result helpers (30 minutes)
5. ⏳ Consider JSON schema builders (2-3 hours, optional)

**Total Phase 2:** ~3-4 hours, ~80 lines improved

### Phase 3: Low Priority (Backlog)
6. ⏳ Review test patterns (1-2 hours, if needed)

---

## Implementation Guidelines

### Before Refactoring
1. ✅ Run `golangci-lint --enable=dupl` to identify duplication
2. ✅ Measure code reduction potential
3. ✅ Review dependencies and test coverage
4. ✅ Create test cases if needed

### During Refactoring
1. ✅ Extract helpers incrementally
2. ✅ Update all call sites
3. ✅ Run tests frequently
4. ✅ Verify builds succeed

### After Refactoring
1. ✅ Run full test suite
2. ✅ Verify golangci-lint passes
3. ✅ Measure code reduction
4. ✅ Update documentation if needed
5. ✅ Commit with clear message

---

## Success Criteria

### Code Quality
- ✅ Duplication warnings reduced by 80%+
- ✅ Code lines reduced by ~300-400 lines
- ✅ Test coverage maintained or improved
- ✅ Builds succeed
- ✅ No performance regressions

### Maintainability
- ✅ Single source of truth for tool definitions
- ✅ Centralized test helpers
- ✅ Consistent error handling patterns
- ✅ Easier to add new tools/tests

---

## Comparison with exarp-go

### Similarities
- Both have MCP tool handler duplication patterns
- Both have test helper duplication
- Both benefit from centralized tool definitions

### Differences
- **devwisdom-go** is simpler (~9,500 lines vs ~28,000 lines)
- **devwisdom-go** has fewer packages (5 vs 10+)
- **devwisdom-go** duplication is more localized (3 main areas)
- **devwisdom-go** refactoring will be faster (~6-8 hours total vs ongoing)

---

## Next Steps

### Immediate (This Sprint)
1. ✅ Review and approve this plan
2. ✅ Start with Phase 1 (High Priority items)
3. ✅ Extract CLI test helpers first (easiest win)

### Short Term (Next Month)
1. ⏳ Complete Phase 2 refactoring
2. ⏳ Add golangci-lint config for ongoing duplication detection
3. ⏳ Document refactoring patterns for future reference

### Long Term (Future)
1. ⏳ Consider applying exarp-go patterns (see `TODO_EXARP_GO_IMPROVEMENTS.md`)
2. ⏳ Framework abstraction if needed
3. ⏳ Additional code quality improvements

---

## References

- **Codebase Size:** ~9,500 lines, 44 Go files
- **Duplication Tool:** `golangci-lint --enable=dupl`
- **Related Plans:** `TODO_EXARP_GO_IMPROVEMENTS.md`, `docs/EXARP_GO_LESSONS.md`
- **Patterns:** See `.cursor/rules/exarp-go-patterns.mdc`

---

**Last Updated:** 2026-01-12  
**Next Review:** After Phase 1 completion
