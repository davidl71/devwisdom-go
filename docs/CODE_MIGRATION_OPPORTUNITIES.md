# Code Migration Opportunities from devwisdom-go to mcp-go-core

**Date:** 2026-01-13  
**Task:** T-1768325426665  
**Status:** 📊 Analysis Complete

---

## Executive Summary

This document identifies code in devwisdom-go that could potentially be migrated to mcp-go-core for reuse across MCP server projects.

**Key Findings:**
- **High Value:** `convertToMap()` helper function
- **Medium Value:** Tool handler wrapper pattern (if generalized)
- **Low Value:** Most other code is project-specific

---

## Identified Migration Opportunities

### 🔴 High Priority

#### 1. `convertToMap()` Helper Function
**Location:** `internal/mcp/sdk_adapter.go:192-210`  
**Status:** ✅ **RECOMMENDED FOR MIGRATION**

**Current Implementation:**
```go
// convertToMap converts any result to map[string]interface{}
// Handles both maps and structs by marshaling/unmarshaling through JSON
func convertToMap(result interface{}) (map[string]interface{}, error) {
	// If already a map, return it
	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	// Marshal to JSON and unmarshal to map
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &resultMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result to map: %w", err)
	}

	return resultMap, nil
}
```

**Why Migrate:**
- ✅ **Generic utility** - Useful for any MCP server
- ✅ **Common pattern** - Converting structs to maps for JSON responses
- ✅ **No dependencies** - Only uses standard library
- ✅ **Well-tested** - Used in production code

**Proposed Location:** `mcp-go-core/pkg/mcp/response/convert.go`

**Migration Effort:** Low (1-2 hours)
- Move function to mcp-go-core
- Add tests
- Update devwisdom-go to use from mcp-go-core

**Benefits:**
- Reusable across all MCP servers
- Standardized result conversion
- Consistent error handling

---

### 🟡 Medium Priority

#### 2. Tool Handler Wrapper Pattern
**Location:** `internal/mcp/sdk_adapter.go:141-188`  
**Status:** ⚠️ **CONDITIONAL MIGRATION** (if generalized)

**Current Implementation:**
```go
func wrapToolHandler(handler ToolHandlerFunc) func(context.Context, *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Argument unmarshaling
		// Handler execution
		// Result conversion
		// Response formatting
	}
}
```

**Why Consider Migration:**
- ✅ **Common pattern** - Many MCP servers need tool handler wrapping
- ✅ **Reduces duplication** - Eliminates repetitive wrapper code
- ⚠️ **SDK-specific** - Currently tied to `modelcontextprotocol/go-sdk`

**Challenges:**
- Requires SDK-specific types (`mcp.CallToolRequest`, `mcp.CallToolResult`)
- May need framework abstraction layer
- Different MCP servers may have different wrapper needs

**Proposed Location:** `mcp-go-core/pkg/mcp/framework/adapters/gosdk/wrapper.go`

**Migration Effort:** Medium (3-4 hours)
- Generalize wrapper pattern
- Make it configurable
- Add to GoSDK adapter package

**Benefits:**
- Standardized tool handler wrapping
- Consistent error handling
- Reduced boilerplate

**Recommendation:** ⏸️ **Defer** - Only migrate if pattern is needed in other projects or if framework abstraction is implemented.

---

### 🟢 Low Priority / Not Recommended

#### 3. CLI Application Structure
**Location:** `internal/cli/app.go`  
**Status:** ❌ **DO NOT MIGRATE**

**Reason:**
- Project-specific command structure
- Wisdom-specific commands (quote, consult, briefing)
- Different from generic MCP CLI utilities

**Note:** mcp-go-core already has `pkg/mcp/cli` for TTY detection and argument parsing, which is sufficient.

---

#### 4. Configuration Management
**Location:** `internal/config/config.go`  
**Status:** ❌ **DO NOT MIGRATE**

**Reason:**
- Wisdom-specific configuration (sources, Hebrew options)
- Project-specific environment variables
- Domain-specific business logic

**Note:** mcp-go-core already has `pkg/mcp/config` for base server configuration, which is sufficient.

---

#### 5. Error Result Helpers
**Location:** `internal/mcp/sdk_adapter.go:111-132`  
**Status:** ❌ **DO NOT MIGRATE** (Already in mcp-go-core)

**Reason:**
- mcp-go-core already has error handling patterns
- SDK-specific implementation
- Simple enough to keep project-specific

---

#### 6. Resource Response Conversion
**Location:** `internal/mcp/sdk_adapter.go:408-432`  
**Status:** ❌ **DO NOT MIGRATE**

**Reason:**
- Project-specific response format
- Tightly coupled to devwisdom-go resource structure
- Not generic enough

---

## Migration Priority Summary

| Code | Priority | Effort | Value | Recommendation |
|------|----------|--------|-------|----------------|
| `convertToMap()` | 🔴 High | Low | High | ✅ **Migrate** |
| Tool handler wrapper | 🟡 Medium | Medium | Medium | ⏸️ **Defer** |
| CLI app structure | 🟢 Low | N/A | Low | ❌ **Don't migrate** |
| Config management | 🟢 Low | N/A | Low | ❌ **Don't migrate** |
| Error helpers | 🟢 Low | N/A | Low | ❌ **Already in mcp-go-core** |
| Resource conversion | 🟢 Low | N/A | Low | ❌ **Don't migrate** |

---

## Recommended Migration Plan

### Phase 1: High-Value Migration (Immediate)

#### Migrate `convertToMap()` Function

**Steps:**
1. Create `mcp-go-core/pkg/mcp/response/convert.go`
2. Move `convertToMap()` function
3. Add comprehensive tests
4. Update devwisdom-go to import from mcp-go-core
5. Remove local implementation

**Estimated Effort:** 1-2 hours  
**Benefits:** High reusability, standardized conversion

**Code to Migrate:**
```go
// mcp-go-core/pkg/mcp/response/convert.go
package response

import (
	"encoding/json"
	"fmt"
)

// ConvertToMap converts any result to map[string]interface{}
// Handles both maps and structs by marshaling/unmarshaling through JSON.
// This is useful for standardizing tool responses before formatting.
func ConvertToMap(result interface{}) (map[string]interface{}, error) {
	// If already a map, return it
	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	// Marshal to JSON and unmarshal to map
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &resultMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result to map: %w", err)
	}

	return resultMap, nil
}
```

**Update devwisdom-go:**
```go
// Before
resultMap, err := convertToMap(result)

// After
import mcpresponse "github.com/davidl71/mcp-go-core/pkg/mcp/response"
resultMap, err := mcpresponse.ConvertToMap(result)
```

---

### Phase 2: Conditional Migration (Future)

#### Tool Handler Wrapper Pattern

**Condition:** Only migrate if:
- Framework abstraction layer is implemented in mcp-go-core
- Other MCP servers need similar wrapper functionality
- Pattern becomes common across multiple projects

**Current Status:** ⏸️ **Defer** - Not needed yet

---

## Code Analysis by Package

### `internal/mcp/`
- ✅ `convertToMap()` - **MIGRATE** (high value)
- ⏸️ `wrapToolHandler()` - **DEFER** (conditional)
- ❌ `convertResourceResponse()` - **KEEP** (project-specific)
- ❌ `newErrorResult()` / `newSuccessResult()` - **KEEP** (SDK-specific, simple)

### `internal/cli/`
- ❌ All code - **KEEP** (project-specific CLI app)

### `internal/config/`
- ❌ All code - **KEEP** (wisdom-specific configuration)

### `internal/logging/`
- ❌ `ConsultationLogger` - **KEEP** (domain-specific, see separate evaluation)

### `internal/wisdom/`
- ❌ All code - **KEEP** (core business logic, project-specific)

---

## Migration Checklist

### High Priority
- [x] Migrate `convertToMap()` to `mcp-go-core/pkg/mcp/response/convert.go` ✅ **COMPLETE**
- [x] Add tests for `ConvertToMap()` in mcp-go-core ✅ **COMPLETE**
- [x] Update devwisdom-go to use `mcpresponse.ConvertToMap()` ✅ **COMPLETE**
- [x] Remove local `convertToMap()` implementation ✅ **COMPLETE**
- [x] Update documentation ✅ **COMPLETE**

### Medium Priority (Deferred)
- [ ] Evaluate tool handler wrapper pattern generalization
- [ ] Consider framework abstraction layer
- [ ] Assess need in other MCP projects

### Low Priority (Not Recommended)
- ❌ CLI app structure - Keep project-specific
- ❌ Config management - Keep wisdom-specific
- ❌ Resource conversion - Keep project-specific

---

## Benefits of Migration

### For mcp-go-core
- ✅ More reusable utilities
- ✅ Broader library coverage
- ✅ Standardized patterns

### For devwisdom-go
- ✅ Reduced code duplication
- ✅ Shared improvements
- ✅ Consistent patterns with other MCP servers

### For Other MCP Servers
- ✅ Ready-to-use utilities
- ✅ Proven patterns
- ✅ Consistent implementations

---

## Conclusion

**Immediate Action:** Migrate `convertToMap()` function to mcp-go-core

**Future Consideration:** Evaluate tool handler wrapper pattern if framework abstraction is implemented

**Keep Project-Specific:** CLI, config, domain logic, ConsultationLogger

---

**Last Updated:** 2026-01-13  
**Next Review:** After `convertToMap()` migration
