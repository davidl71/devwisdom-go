# Standalone CLI and Zsh Plugin Implementation Plan

**Created**: 2025-01-26  
**Project**: devwisdom-go  
**Goal**: Enable standalone CLI execution and zsh plugin integration

---

## Overview

Currently, devwisdom-go runs only as an MCP server via stdio. We need to add:
1. **Standalone CLI** - Direct command-line interface for shell usage
2. **Zsh Plugin** - Integration with zsh for convenient daily wisdom

---

## Requirements

### Standalone CLI

**Commands to implement:**
- `devwisdom quote [--source SOURCE] [--score SCORE]` - Get a wisdom quote
- `devwisdom consult [--metric METRIC] [--tool TOOL] [--stage STAGE] [--score SCORE]` - Consult an advisor
- `devwisdom sources` - List available sources
- `devwisdom advisors` - List available advisors
- `devwisdom briefing [--days DAYS]` - Get daily briefing
- `devwisdom version` - Show version
- `devwisdom help` - Show help

**Output formats:**
- Default: Human-readable text
- `--json` flag: JSON output
- `--quiet` flag: Minimal output (just the quote text)

### Zsh Plugin

**Features:**
- `devwisdom-daily` function - Show daily wisdom on shell startup
- `devwisdom-quote` function - Quick quote access
- `devwisdom-consult` function - Quick consultation
- Auto-completion support
- Configurable via environment variables

**Installation:**
- Plugin directory structure
- Installation script
- Documentation

---

## Implementation Plan

### Phase 1: CLI Structure (Foundation)
1. Create `cmd/cli/main.go` - New CLI entry point
2. Add CLI flag parsing (cobra or standard library)
3. Implement command structure
4. Add version and help commands
5. Create output formatting utilities

### Phase 2: Core CLI Commands
1. Implement `quote` command
2. Implement `consult` command  
3. Implement `sources` command
4. Implement `advisors` command
5. Implement `briefing` command

### Phase 3: CLI Polish
1. Add JSON output format
2. Add quiet mode
3. Add error handling and user-friendly messages
4. Add configuration file support (optional)
5. Update Makefile with CLI build target

### Phase 4: Zsh Plugin
1. Create `zsh/` directory structure
2. Create `devwisdom.plugin.zsh` main plugin file
3. Implement helper functions:
   - `devwisdom-daily`
   - `devwisdom-quote`
   - `devwisdom-consult`
4. Add auto-completion support
5. Create installation script

### Phase 5: Documentation & Testing
1. Update README with CLI usage
2. Add zsh plugin installation instructions
3. Create examples and use cases
4. Add CLI tests
5. Test zsh plugin integration

---

## Technical Decisions

### CLI Library Choice
**Decision**: Use standard library `flag` package (keep dependencies minimal)
**Alternative**: cobra (more features, but adds dependency)
**Rationale**: Project rules prefer standard library, CLI is simple enough

### Binary Names
- Main binary: `devwisdom` (MCP server mode by default, or CLI mode)
- Or: `devwisdom` (CLI) and `devwisdom-server` (MCP server)
- **Decision**: Single binary with mode detection:
  - If stdin is a TTY → CLI mode
  - If stdin is not a TTY → MCP server mode

### Zsh Plugin Location
- Standard: `~/.oh-my-zsh/custom/plugins/devwisdom/`
- Or: `~/.zsh/plugins/devwisdom/`
- **Decision**: Support both, default to oh-my-zsh

---

## File Structure

```
devwisdom-go/
├── cmd/
│   ├── server/          # MCP server (existing)
│   │   └── main.go
│   └── cli/              # NEW: CLI entry point
│       └── main.go
├── internal/
│   ├── cli/              # NEW: CLI command implementations
│   │   ├── commands.go
│   │   ├── output.go
│   │   └── format.go
│   └── ...
├── zsh/                  # NEW: Zsh plugin
│   ├── devwisdom.plugin.zsh
│   ├── _devwisdom        # Auto-completion
│   └── install.sh
└── Makefile              # Update with CLI targets
```

---

## Dependencies

- No new external dependencies (use standard library)
- Zsh plugin requires zsh 5.0+

---

## Testing Strategy

1. **CLI Tests**: Unit tests for each command
2. **Integration Tests**: End-to-end CLI usage
3. **Zsh Plugin Tests**: Manual testing in zsh environment
4. **Cross-platform**: Ensure works on macOS, Linux

---

## Success Criteria

- [ ] `devwisdom quote` works from command line
- [ ] `devwisdom consult` works with all parameters
- [ ] `devwisdom sources` lists all sources
- [ ] JSON output format works
- [ ] Zsh plugin installs correctly
- [ ] `devwisdom-daily` shows wisdom on shell startup
- [ ] Auto-completion works
- [ ] Documentation is complete
- [ ] All tests pass

---

## Estimated Effort

- Phase 1: 2-3 hours
- Phase 2: 4-6 hours
- Phase 3: 2-3 hours
- Phase 4: 3-4 hours
- Phase 5: 2-3 hours

**Total**: ~13-19 hours (2-3 days)
