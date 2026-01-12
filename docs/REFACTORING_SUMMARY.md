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
... (see docs/REFACTORING_PLAN.md for full details)
