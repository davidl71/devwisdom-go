# Todo2 Tasks: Apply exarp-go Patterns to devwisdom-go

This document contains Todo2 tasks for applying lessons learned from exarp-go to devwisdom-go.

**Created:** 2025-01-XX  
**Source:** Lessons learned from exarp-go project  
**Total Tasks:** 12 tasks across 4 phases

---

## Phase 1: Framework Abstraction & Architecture (High Priority)

### T-A1: Create Framework-Agnostic MCP Interface
**Status:** Todo  
**Priority:** High  
**Tags:** architecture, framework, refactor  
**Dependencies:** None

**Description:**
Create framework-agnostic interfaces for MCP protocol operations, following exarp-go's pattern. This enables future framework switching and improves testability.

**Tasks:**
- Create `internal/framework/server.go` with `MCPServer` interface
- Define `ToolHandler`, `PromptHandler`, `ResourceHandler` types
- Create `ToolSchema`, `TextContent`, `ToolInfo` types
- Move existing MCP implementation behind interface

**Acceptance Criteria:**
- [ ] `MCPServer` interface defined with all required methods
- [ ] Existing wisdom server implements the interface
- [ ] No breaking changes to existing functionality
- [ ] Unit tests verify interface implementation

**Research Required:**
- Review exarp-go `internal/framework/server.go` implementation
- Understand current devwisdom-go MCP server structure
- Identify required interface methods from current usage

---

### T-A2: Implement Factory Pattern for Server Creation
**Status:** Todo  
**Priority:** High  
**Tags:** architecture, factory, refactor  
**Dependencies:** T-A1

**Description:**
Implement factory pattern for server creation following exarp-go's `internal/factory/server.go` pattern. This centralizes server creation logic and enables configuration-driven creation.

**Tasks:**
- Create `internal/factory/server.go`
- Implement `NewServer()` and `NewServerFromConfig()` functions
- Update `cmd/server/main.go` to use factory
- Add configuration support

**Acceptance Criteria:**
- [ ] Factory functions created and tested
- [ ] Server creation uses factory pattern
- [ ] Configuration-driven server creation works
- [ ] Unit tests for factory functions

**Research Required:**
- Study exarp-go factory implementation
- Review devwisdom-go current server initialization
- Identify configuration needs

---

### T-A3: Add Configuration Management with Environment Variables
**Status:** Todo  
**Priority:** Medium  
**Tags:** configuration, environment, refactor  
**Dependencies:** None

**Description:**
Implement environment variable-based configuration following exarp-go's `internal/config/config.go` pattern. Support both environment variables and sensible defaults.

**Tasks:**
- Create `internal/config/config.go` if not exists
- Add `Config` struct with environment variable support
- Implement `Load()` function with environment override
- Add configuration validation
- Update server to use configuration

**Acceptance Criteria:**
- [ ] Configuration struct defined with environment variable tags
- [ ] `Load()` function reads from environment with defaults
- [ ] Validation ensures valid configuration
- [ ] Unit tests for configuration loading
- [ ] Documentation for environment variables

**Research Required:**
- Review exarp-go config implementation
- Identify devwisdom-go configuration needs
- Define environment variable naming convention

---

## Phase 2: CLI/MCP Dual Mode (High Priority)

### T-B1: Implement TTY Detection for Dual Mode
**Status:** Todo  
**Priority:** High  
**Tags:** cli, mcp, tty, refactor  
**Dependencies:** T-A1, T-A2

**Description:**
Implement TTY detection to enable single binary for both CLI and MCP server modes, following exarp-go's pattern in `cmd/server/main.go` and `internal/cli/cli.go`.

**Tasks:**
- Create `internal/cli/cli.go` with `IsTTY()` function
- Update `cmd/server/main.go` to detect TTY
- Route to CLI mode if TTY detected
- Route to MCP server mode if stdio
- Test both modes

**Acceptance Criteria:**
- [ ] `IsTTY()` function detects terminal correctly
- [ ] Main function routes to appropriate mode
- [ ] CLI mode works when run interactively
- [ ] MCP mode works when run via stdio
- [ ] Integration tests for both modes

**Research Required:**
- Study exarp-go TTY detection implementation
- Understand Go `golang.org/x/term` package
- Test TTY detection on different platforms

---

### T-B2: Create CLI Command Structure
**Status:** Todo  
**Priority:** High  
**Tags:** cli, commands, refactor  
**Dependencies:** T-B1

**Description:**
Implement CLI command structure following exarp-go's CLI pattern. Support tool execution, listing, and interactive mode.

**Tasks:**
- Implement `Run()` function in `internal/cli/cli.go`
- Add command parsing with flag package
- Implement `listAllTools()` function
- Implement `executeTool()` function
- Implement `runInteractive()` function (optional)
- Add usage and help messages

**Acceptance Criteria:**
- [ ] CLI can list all available tools
- [ ] CLI can execute tools with arguments
- [ ] Interactive mode works (if implemented)
- [ ] Help messages are clear and useful
- [ ] Error handling is robust

**Research Required:**
- Review exarp-go CLI implementation
- Understand current devwisdom-go tool structure
- Define CLI command structure

---

### T-B3: Reuse Server Infrastructure for CLI
**Status:** Todo  
**Priority:** Medium  
**Tags:** cli, refactor, code-reuse  
**Dependencies:** T-B1, T-B2

**Description:**
Refactor CLI to reuse server infrastructure for tool execution, following exarp-go's pattern where CLI uses server's `CallTool()` method.

**Tasks:**
- Update CLI to use server interface for tool execution
- Share tool registration between CLI and MCP modes
- Ensure consistent behavior between modes
- Update tests to cover shared infrastructure

**Acceptance Criteria:**
- [ ] CLI uses server interface for tool execution
- [ ] Tool registration is shared between modes
- [ ] Behavior is consistent between CLI and MCP
- [ ] No code duplication
- [ ] Tests verify shared infrastructure works

**Research Required:**
- Review exarp-go CLI server integration
- Identify code duplication opportunities
- Plan refactoring approach

---

## Phase 3: Tool Registration & Organization (Medium Priority)

### T-C1: Organize Tool Registration in Batches
**Status:** Todo  
**Priority:** Medium  
**Tags:** refactor, organization, tools  
**Dependencies:** T-A1

**Description:**
Organize tool registration in batches following exarp-go's pattern in `internal/tools/registry.go`. This improves maintainability as tools grow.

**Tasks:**
- Review current tool registration structure
- Organize tools into logical batches
- Implement batch registration functions
- Update `RegisterAllTools()` to use batches
- Add batch-level error handling

**Acceptance Criteria:**
- [ ] Tools organized into logical batches
- [ ] Batch registration functions implemented
- [ ] Error handling per batch
- [ ] Clear organization and maintainability
- [ ] No breaking changes

**Research Required:**
- Review exarp-go batch registration pattern
- Analyze current devwisdom-go tool structure
- Identify logical batch groupings

---

### T-C2: Improve Tool Schema Definitions
**Status:** Todo  
**Priority:** Medium  
**Tags:** tools, schema, refactor  
**Dependencies:** T-C1

**Description:**
Improve tool schema definitions following exarp-go's structured schema pattern. Use `ToolSchema` type with proper type definitions.

**Tasks:**
- Review current tool schemas
- Create structured `ToolSchema` definitions
- Ensure consistent schema format
- Add validation for schemas
- Update documentation

**Acceptance Criteria:**
- [ ] All tools use structured schemas
- [ ] Schemas are consistent across tools
- [ ] Schema validation works
- [ ] Documentation reflects schemas
- [ ] Type safety improved

**Research Required:**
- Review exarp-go schema patterns
- Study JSON Schema best practices
- Identify schema improvement opportunities

---

## Phase 4: Development Workflow & Testing (Medium Priority)

### T-D1: Create Hot Reload Development Scripts
**Status:** Todo  
**Priority:** Medium  
**Tags:** development, scripts, workflow  
**Dependencies:** None

**Description:**
Create development scripts with hot reload functionality following exarp-go's development workflow. Support file watching and auto-rebuild.

**Tasks:**
- Create `dev-go.sh` script with file watching
- Support `fswatch` (macOS) and `inotifywait` (Linux)
- Fallback to polling if no file watcher available
- Auto-rebuild on Go file changes
- Auto-rerun tests on changes (optional)
- Update Makefile with dev targets

**Acceptance Criteria:**
- [ ] Dev script watches relevant files
- [ ] Auto-rebuilds on Go file changes
- [ ] Works on macOS and Linux
- [ ] Falls back gracefully if no watcher
- [ ] Makefile targets work

**Research Required:**
- Review exarp-go dev scripts
- Understand file watching tools
- Test on different platforms

---

### T-D2: Enhance Makefile with Tool Detection
**Status:** Todo  
**Priority:** Low  
**Tags:** makefile, build, development  
**Dependencies:** T-D1

**Description:**
Enhance Makefile with tool detection following exarp-go's Makefile pattern. Detect available tools and optimize targets accordingly.

**Tasks:**
- Review exarp-go Makefile patterns
- Add `make config` target for tool detection
- Generate `.make.config` file
- Add conditional targets based on available tools
- Update help target

**Acceptance Criteria:**
- [ ] `make config` detects available tools
- [ ] Conditional targets work correctly
- [ ] Help target shows available targets
- [ ] Config file is generated correctly
- [ ] No hard dependencies on optional tools

**Research Required:**
- Study exarp-go Makefile implementation
- Understand Make conditional logic
- Identify tools to detect

---

### T-D3: Add Tool Testing Utilities
**Status:** Todo  
**Priority:** Medium  
**Tags:** testing, tools, utilities  
**Dependencies:** T-B2

**Description:**
Add tool testing utilities following exarp-go's `testToolExecution()` pattern. Enable easy testing of tools with example arguments.

**Tasks:**
- Implement `testToolExecution()` function
- Generate example arguments from schemas
- Add `generateExampleArgs()` helper
- Create CLI flag for tool testing
- Add documentation

**Acceptance Criteria:**
- [ ] Tool testing function works
- [ ] Example arguments generated from schemas
- [ ] CLI can test tools easily
- [ ] Tests verify utility functions
- [ ] Documentation added

**Research Required:**
- Review exarp-go testing utilities
- Understand schema-based test generation
- Identify testing needs

---

## Phase 5: Security & Path Validation (High Priority)

### T-E1: Implement Path Validation Security
**Status:** Todo  
**Priority:** High  
**Tags:** security, paths, validation  
**Dependencies:** None

**Description:**
Implement path validation security following exarp-go's `internal/security/path.go` pattern. Validate all paths before file operations.

**Tasks:**
- Create `internal/security/path.go`
- Implement `ValidatePath()` function
- Validate workspace root paths
- Validate file paths for operations
- Add security tests

**Acceptance Criteria:**
- [ ] Path validation function implemented
- [ ] All file operations use path validation
- [ ] Path traversal attacks prevented
- [ ] Security tests pass
- [ ] Clear error messages

**Research Required:**
- Review exarp-go security implementation
- Study path traversal attack vectors
- Understand Go path validation best practices

---

### T-E2: Add Workspace Root Detection with Fallbacks
**Status:** Todo  
**Priority:** Medium  
**Tags:** paths, detection, fallbacks  
**Dependencies:** T-E1

**Description:**
Implement robust workspace root detection with multiple fallback strategies following exarp-go's `getWorkspaceRoot()` pattern.

**Tasks:**
- Implement `getWorkspaceRoot()` function
- Check `PROJECT_ROOT` environment variable
- Fallback to executable location
- Fallback to runtime.Caller
- Handle placeholder values
- Add tests for all fallback paths

**Acceptance Criteria:**
- [ ] Workspace root detection works
- [ ] Multiple fallback strategies implemented
- [ ] Handles all deployment scenarios
- [ ] Tests cover all paths
- [ ] Graceful fallbacks

**Research Required:**
- Review exarp-go workspace root detection
- Understand deployment scenarios
- Test on different platforms

---

## Summary

### Task Breakdown by Phase

**Phase 1: Framework Abstraction & Architecture (3 tasks)**
- T-A1: Framework-Agnostic MCP Interface (High)
- T-A2: Factory Pattern for Server Creation (High)
- T-A3: Configuration Management (Medium)

**Phase 2: CLI/MCP Dual Mode (3 tasks)**
- T-B1: TTY Detection for Dual Mode (High)
- T-B2: CLI Command Structure (High)
- T-B3: Reuse Server Infrastructure (Medium)

**Phase 3: Tool Registration & Organization (2 tasks)**
- T-C1: Organize Tool Registration in Batches (Medium)
- T-C2: Improve Tool Schema Definitions (Medium)

**Phase 4: Development Workflow & Testing (3 tasks)**
- T-D1: Hot Reload Development Scripts (Medium)
- T-D2: Enhance Makefile with Tool Detection (Low)
- T-D3: Tool Testing Utilities (Medium)

**Phase 5: Security & Path Validation (2 tasks)**
- T-E1: Path Validation Security (High)
- T-E2: Workspace Root Detection (Medium)

### Priority Summary

- **High Priority:** 6 tasks
- **Medium Priority:** 5 tasks
- **Low Priority:** 1 task

### Dependencies

```
T-A1 (Interface) ──┐
                   ├──> T-A2 (Factory)
                   └──> T-C1 (Batch Registration)

T-A1 ──> T-B1 (TTY Detection)
T-A2 ──> T-B1
T-B1 ──> T-B2 (CLI Commands)
T-B2 ──> T-B3 (Reuse Infrastructure)
T-B2 ──> T-D3 (Testing Utilities)
```

---

## Implementation Notes

### Research Requirements
Each task requires research before starting:
1. Review exarp-go implementation
2. Understand current devwisdom-go structure
3. Identify specific changes needed
4. Plan implementation approach

### Migration Strategy
- Implement incrementally, one phase at a time
- Maintain backward compatibility
- Test thoroughly after each change
- Update documentation as changes are made

### Testing Strategy
- Unit tests for each component
- Integration tests for CLI/MCP modes
- Security tests for path validation
- Manual testing on different platforms

---

## Related Documentation

- [EXARP_GO_LESSONS.md](./EXARP_GO_LESSONS.md) - Complete lessons learned document
- [TODO.md](../TODO.md) - Current project TODO list
- [PROJECT_GOALS.md](../PROJECT_GOALS.md) - Project goals and phases

