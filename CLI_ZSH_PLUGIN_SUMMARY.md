# Standalone CLI and Zsh Plugin - Planning Summary

**Created**: 2025-01-26  
**Status**: ✅ Planning Complete

---

## Overview

Planning complete for adding standalone CLI and zsh plugin support to devwisdom-go. All tasks have been created in todo2 and organized into a logical hierarchy.

---

## Plan Document

See `CLI_ZSH_PLUGIN_PLAN.md` for detailed implementation plan including:
- Requirements
- Technical decisions
- File structure
- Testing strategy
- Success criteria

---

## Tasks Created in Todo2

### Parent Task
- **Phase: Standalone CLI and Zsh Plugin - Foundation** (ID: `e815d0b2-81b5-411a-b53f-0f2976f63fb5`)
  - Priority: 8/10
  - Complexity: 4/10
  - Estimated: 3 hours
  - Status: Pending

### CLI Command Tasks (under Foundation)
1. **Implement CLI quote command** (ID: `7123c57d-272b-4554-842c-7a9b4d9be537`)
   - Priority: 8/10, Complexity: 3/10, Est: 2h
   - Command: `devwisdom quote [--source SOURCE] [--score SCORE]`

2. **Implement CLI consult command** (ID: `3d9bd991-baf6-4981-bf9f-61fc5ec3ada1`)
   - Priority: 8/10, Complexity: 4/10, Est: 3h
   - Command: `devwisdom consult [--metric METRIC] [--tool TOOL] [--stage STAGE] [--score SCORE]`

3. **Implement CLI sources and advisors commands** (ID: `e8826b08-cfce-4c10-90c4-a8294a37c189`)
   - Priority: 7/10, Complexity: 2/10, Est: 1h
   - Commands: `devwisdom sources`, `devwisdom advisors`

4. **Implement CLI briefing command** (ID: `132bc9ab-e5b5-4678-b0b0-33c8c8502a8b`)
   - Priority: 7/10, Complexity: 3/10, Est: 2h
   - Command: `devwisdom briefing [--days DAYS]`

### Zsh Plugin Tasks
5. **Create zsh plugin structure** (ID: `1d2050ca-63f9-4153-a032-ea7679ab6176`)
   - Priority: 7/10, Complexity: 4/10, Est: 4h
   - Functions: `devwisdom-daily`, `devwisdom-quote`, `devwisdom-consult`

6. **Create zsh plugin installation script** (ID: `9ff270b9-9e4a-49ba-bed1-1e5b4bf4a4b0`)
   - Priority: 6/10, Complexity: 2/10, Est: 1h
   - Support oh-my-zsh and standard zsh

### Documentation Task
7. **Update Makefile and documentation** (ID: `9b288e9d-6d95-410c-b3a6-f9fa18234631`)
   - Priority: 6/10, Complexity: 2/10, Est: 2h
   - Update README, Makefile, add examples

---

## Total Estimated Effort

**~18 hours** (2-3 days of focused work)

Breakdown:
- Foundation: 3h
- CLI Commands: 8h (quote: 2h, consult: 3h, sources/advisors: 1h, briefing: 2h)
- Zsh Plugin: 5h (structure: 4h, installation: 1h)
- Documentation: 2h

---

## Technical Decisions

1. **CLI Library**: Standard library `flag` package (no external dependencies)
2. **Binary Mode**: Single binary with TTY detection (CLI vs MCP server)
3. **Output Formats**: Human-readable (default), JSON (`--json`), Quiet (`--quiet`)
4. **Zsh Plugin Location**: Support both oh-my-zsh and standard zsh directories

---

## Next Steps

1. Start with **Foundation** task to set up CLI structure
2. Implement CLI commands in order (quote → consult → sources/advisors → briefing)
3. Create zsh plugin structure and functions
4. Add installation script and documentation
5. Test all functionality
6. Update README with usage examples

---

## Project ID

All tasks are tracked in project: `bfec727f-93a6-42e5-bfe6-66e01b495ddf` (devwisdom-go)

---

## Files Created

- `CLI_ZSH_PLUGIN_PLAN.md` - Detailed implementation plan
- `CLI_ZSH_PLUGIN_SUMMARY.md` - This summary document
