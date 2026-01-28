#!/usr/bin/env bash
# Install shared Cursor MCP config to ~/.cursor/mcp.json from single source of truth.
# Default: merge with existing (template wins for same name; extra servers kept).
# --replace: overwrite with template only (removes extra servers e.g. openmemory, duplicates).
# Usage: ./install-mcp-config.sh [--replace] [PROJECTS_BASE]
#   PROJECTS_BASE defaults to $HOME/Projects or set CURSOR_PROJECTS_BASE.

set -e

REPLACE=
while [[ $# -gt 0 ]]; do
  case "$1" in
    --replace) REPLACE=1; shift ;;
    *) break ;;
  esac
done

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
  if [[ -n "$REPLACE" ]]; then
    echo "$APPLIED" > "$DEST"
    echo "Replaced config with template only (removed extra servers and duplicates)."
    return
  fi
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
