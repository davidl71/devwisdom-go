#!/usr/bin/env bash
# Install shared Cursor MCP config to ~/.cursor/mcp.json from single source of truth.
# Merges with existing config: template servers win (no duplicates), extra servers kept.
# Usage: ./install-mcp-config.sh [PROJECTS_BASE]
#   PROJECTS_BASE defaults to $HOME/Projects or set CURSOR_PROJECTS_BASE.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE="${SCRIPT_DIR}/mcp-servers.json.template"
DEST="${HOME}/.cursor/mcp.json"
PROJECTS_BASE="${1:-${CURSOR_PROJECTS_BASE:-$HOME/Projects}}"

if [[ ! -f "$TEMPLATE" ]]; then
  echo "Error: template not found: $TEMPLATE" >&2
  exit 1
fi

# Normalize path (no trailing slash)
PROJECTS_BASE="${PROJECTS_BASE%/}"

# Apply PROJECTS_BASE to template
APPLIED="$(sed "s|PROJECTS_BASE|${PROJECTS_BASE}|g" "$TEMPLATE")"

merge_and_write() {
  # Merge: template servers + existing servers not in template (template wins for same name = no duplicates)
  if command -v jq >/dev/null 2>&1 && [[ -f "$DEST" ]]; then
    local existing_servers merged
    existing_servers="$(jq -c '.mcpServers // {}' "$DEST" 2>/dev/null)" || existing_servers="{}"
    merged="$(echo "$APPLIED" | jq --argjson exist "$existing_servers" '
      .mcpServers as $tmpl |
      ($exist | to_entries | map(select(.key | IN($tmpl | keys) | not)) | from_entries) as $extra |
      # Standard (template) overrides existing for same name: extra + tmpl so tmpl wins
      { mcpServers: ($extra + $tmpl) }
    ')"
    echo "$merged" > "$DEST"
    echo "Merged with existing config (template wins for same name; extra servers kept)."
  else
    echo "$APPLIED" > "$DEST"
    if [[ -f "$DEST" ]] && ! command -v jq >/dev/null 2>&1; then
      echo "Note: jq not found; overwrote config. Install jq to enable merge (keeps your extra servers)." >&2
    fi
  fi
}

if [[ -f "$DEST" ]]; then
  BACKUP="${DEST}.backup.$(date +%Y%m%d%H%M%S)"
  cp "$DEST" "$BACKUP"
  echo "Backed up existing config to $BACKUP"
fi

mkdir -p "$(dirname "$DEST")"
merge_and_write
echo "Installed MCP config to $DEST (PROJECTS_BASE=$PROJECTS_BASE)"
echo "Restart Cursor (or reload window) for changes to take effect."
