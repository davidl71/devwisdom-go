# mcp-go-core Integration Plan

**Date:** 2026-01-12  
**Status:** 📋 Planning  
**Library:** `github.com/davidl71/mcp-go-core`

---

## Overview

This document outlines the plan to integrate `mcp-go-core` into `devwisdom-go` to replace duplicated infrastructure code with the shared library.

**Benefits:**
- Reduce code duplication
- Share improvements across projects
- Standardize MCP infrastructure
- Easier maintenance

---

## Current Status

### devwisdom-go Internal Packages
- `internal/logging/` - Structured logging (similar to mcp-go-core)
- `internal/config/` - Configuration management (wisdom-specific)
- `internal/mcp/` - MCP server implementation (custom)
- `internal/cli/` - CLI utilities (TTY detection)

### mcp-go-core Packages Available
- `pkg/mcp/logging/` - Structured logging (uses MCP_DEBUG)
- `pkg/mcp/config/` - Base configuration (framework-agnostic)
- `pkg/mcp/framework/` - Framework abstraction (MCPServer interface)
- `pkg/mcp/security/` - Security utilities (path validation, project root)
- `pkg/mcp/cli/` - CLI utilities (TTY detection)
- `pkg/mcp/protocol/` - JSON-RPC types
- `pkg/mcp/types/` - Common types (TextContent, ToolInfo, etc.)
- `pkg/mcp/platform/` - Platform detection
- `pkg/mcp/factory/` - Server factory pattern

---

## Integration Strategy

### Phase 1: Low-Risk Integration (Logging & Security)

#### 1.1 Replace `internal/logging` with `mcp-go-core/pkg/mcp/logging`
- **Status:** ⏳ Pending
- **Impact:** Low risk - logging is self-contained
- **Changes:**
  - Update imports from `internal/logging` to `github.com/davidl71/mcp-go-core/pkg/mcp/logging`
  - Change env var from `DEVWISDOM_DEBUG` to `MCP_DEBUG` (or support both)
  - Update all logging calls
- **Complexity:** Low
- **Estimated Effort:** 1-2 hours
- **Files Affected:**
  - `internal/mcp/*.go` (server, handlers, sdk_adapter)
  - `internal/wisdom/*.go` (engine, sources)
  - `internal/cli/*.go` (all CLI commands)
  - Remove `internal/logging/` package

#### 1.2 Add `mcp-go-core/pkg/mcp/security` for Path Validation
- **Status:** ⏳ Pending
- **Impact:** Low risk - add new functionality
- **Changes:**
  - Import `mcp-go-core/pkg/mcp/security`
  - Use `security.GetProjectRoot()` and `security.ValidatePath()`
  - Keep existing wisdom-specific security if any
- **Complexity:** Low
- **Estimated Effort:** 30 minutes

### Phase 2: Framework Abstraction (Optional)

#### 2.1 Use `mcp-go-core/pkg/mcp/framework` Interface
- **Status:** ⏳ Pending (Optional)
- **Impact:** Medium risk - architectural change
- **Changes:**
  - Implement `framework.MCPServer` interface
  - Use framework abstraction instead of direct SDK usage
  - Could enable framework switching in future
- **Complexity:** Medium-High
- **Estimated Effort:** 4-6 hours
- **Note:** May not be needed if current SDK usage is sufficient

#### 2.2 Use `mcp-go-core/pkg/mcp/types` for Common Types
- **Status:** ⏳ Pending
- **Impact:** Low-Medium risk - type changes
- **Changes:**
  - Use `types.TextContent`, `types.ToolInfo`, `types.ToolSchema`
  - Update all type references
- **Complexity:** Medium
- **Estimated Effort:** 2-3 hours

### Phase 3: Advanced Integration (Future)

#### 3.1 Use `mcp-go-core/pkg/mcp/factory` for Server Creation
- **Status:** ⏳ Pending (Future)
- **Impact:** Medium risk - architectural change
- **Changes:**
  - Use factory pattern for server creation
  - Centralize server initialization
- **Complexity:** Medium
- **Estimated Effort:** 2-3 hours

#### 3.2 Use `mcp-go-core/pkg/mcp/cli` Utilities
- **Status:** ⏳ Pending (Future)
- **Impact:** Low risk - utility functions
- **Changes:**
  - Replace custom TTY detection with `cli.IsTTY()`
  - Use other CLI utilities as needed
- **Complexity:** Low
- **Estimated Effort:** 1 hour

---

## Detailed Integration Plan

### Step 1: Add mcp-go-core Dependency

```bash
cd /home/dlowes/projects/devwisdom-go
go get github.com/davidl71/mcp-go-core
go mod tidy
```

### Step 2: Replace Logging Package

**Before:**
```go
import "github.com/davidl71/devwisdom-go/internal/logging"
logger := logging.NewLogger()  // Uses DEVWISDOM_DEBUG
```

**After:**
```go
import "github.com/davidl71/mcp-go-core/pkg/mcp/logging"
logger := logging.NewLogger()  // Uses MCP_DEBUG
```

**Migration Steps:**
1. Add mcp-go-core dependency
2. Update all imports
3. Update env var handling (support both DEVWISDOM_DEBUG and MCP_DEBUG for backward compatibility)
4. Test all logging calls
5. Remove `internal/logging/` package

### Step 3: Add Security Utilities

```go
import "github.com/davidl71/mcp-go-core/pkg/mcp/security"

// Use for path validation
projectRoot, err := security.GetProjectRoot(".")
validatedPath, err := security.ValidatePath(path, projectRoot)
```

### Step 4: Use Common Types (Optional)

```go
import "github.com/davidl71/mcp-go-core/pkg/mcp/types"

// Use shared types
content := []types.TextContent{
    {Type: "text", Text: "result"},
}
```

---

## Compatibility Considerations

### Environment Variables
- **Current:** `DEVWISDOM_DEBUG=1`
- **mcp-go-core:** `MCP_DEBUG=1`
- **Solution:** Support both for backward compatibility

### API Compatibility
- Check if `mcp-go-core` logger API matches current usage
- Verify type compatibility
- Test all functionality after migration

### Breaking Changes
- Ensure no breaking changes to public API
- Maintain backward compatibility where possible
- Update tests accordingly

---

## Testing Strategy

### Unit Tests
1. Update all tests to use new imports
2. Verify logging functionality
3. Test path validation
4. Verify type compatibility

### Integration Tests
1. Test full server startup with new logging
2. Verify CLI functionality
3. Test MCP tool execution

### Backward Compatibility
1. Test with `DEVWISDOM_DEBUG` env var (should still work)
2. Test with `MCP_DEBUG` env var (new way)
3. Verify all existing functionality still works

---

## Implementation Checklist

### Phase 1: Logging & Security
- [ ] Add mcp-go-core dependency
- [ ] Create compatibility layer for env vars (DEVWISDOM_DEBUG → MCP_DEBUG)
- [ ] Update all logging imports
- [ ] Update all logging calls
- [ ] Add security utilities imports
- [ ] Test all logging functionality
- [ ] Remove `internal/logging/` package
- [ ] Run full test suite
- [ ] Update documentation

### Phase 2: Types & Framework (Optional)
- [ ] Review framework abstraction need
- [ ] If needed, implement MCPServer interface
- [ ] Update to use common types
- [ ] Test framework integration
- [ ] Update documentation

### Phase 3: Advanced (Future)
- [ ] Use factory pattern
- [ ] Use CLI utilities
- [ ] Additional integrations as needed

---

## Risks & Mitigation

### Risk 1: API Incompatibility
- **Risk:** mcp-go-core API doesn't match current usage
- **Mitigation:** Review APIs before integration, create compatibility layer if needed

### Risk 2: Breaking Changes
- **Risk:** Integration breaks existing functionality
- **Mitigation:** Extensive testing, maintain backward compatibility

### Risk 3: Dependency Management
- **Risk:** Adding external dependency increases complexity
- **Mitigation:** Pin version, test thoroughly, have rollback plan

---

## Success Criteria

- ✅ All tests pass with mcp-go-core
- ✅ No breaking changes to functionality
- ✅ Code duplication reduced
- ✅ Maintainability improved
- ✅ Backward compatibility maintained
- ✅ Documentation updated

---

## References

- **mcp-go-core:** `https://github.com/davidl71/mcp-go-core`
- **Current Logging:** `internal/logging/logger.go`
- **Current Config:** `internal/config/config.go`
- **Integration Plan:** This document

---

**Last Updated:** 2026-01-12  
**Next Review:** After Phase 1 completion

