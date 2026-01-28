# Shared Cursor MCP config

Single source of truth for MCP servers used across devwisdom-go, exarp-go, and other workspaces. Avoids duplicate `.cursor/mcp.json` in each project.

## Quick install

From this directory or from devwisdom-go root:

```bash
# Default: merge with existing (template wins; extra servers kept)
./scripts/cursor/install-mcp-config.sh

# Replace entirely: only the 4 template servers (removes openmemory, duplicates, etc.)
./scripts/cursor/install-mcp-config.sh --replace

# Or via make (from devwisdom-go root)
make install-mcp-config              # merge
make install-mcp-config REPLACE=1    # replace (template only)

# Set PROJECTS_BASE explicitly
./scripts/cursor/install-mcp-config.sh /Users/you/Projects
./scripts/cursor/install-mcp-config.sh --replace /Users/you/Projects

# Or use env
export CURSOR_PROJECTS_BASE=/path/to/Projects
./scripts/cursor/install-mcp-config.sh
```

This writes `~/.cursor/mcp.json`. Cursor substitutes `{{PROJECT_ROOT}}` per workspace.

## What it does

- Reads `mcp-servers.json.template` (placeholder `PROJECTS_BASE`)
- Replaces `PROJECTS_BASE` with your projects directory
- **If `~/.cursor/mcp.json` exists and `jq` is installed:** merges with existing config so that:
  - **Standard overrides local:** servers defined in the template (devwisdom, exarp-go, tractatus_thinking, context7) always replace any same-named entry in existing config. Local/custom values for those names are ignored.
  - **Extra servers kept:** servers that exist only in your current config (not in the template) are left as-is (e.g. a custom or project-specific MCP server).

When this is driven by Ansible, the same rule applies: standard MCP servers from Ansible vars override existing local config for those names.
- Backs up existing `~/.cursor/mcp.json` before writing
- Writes `~/.cursor/mcp.json`

Without `jq`, the script overwrites the file (no merge). Install `jq` (e.g. `brew install jq`) to get merge/dedupe.

## Adding or changing servers

Edit `mcp-servers.json.template`, then run `install-mcp-config.sh` again. All workspaces will see the same set after you restart Cursor (or reload window).

## With gh (GitHub)

If this config lives in a repo you clone via `gh repo clone ...`, run the install script from the clone; use `CURSOR_PROJECTS_BASE` or the first argument so paths point to your machine.
