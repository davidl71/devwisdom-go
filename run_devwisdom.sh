#!/bin/bash
# Wrapper script to run devwisdom-go MCP server
# Handles path resolution for Go binary
#
# Usage: ./run_devwisdom.sh
# Designed for Cursor IDE STDIO transport

set -e

# Find project root (devwisdom-go directory)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"

# Change to project root
cd "$PROJECT_ROOT"

# Build binary if it doesn't exist
if [[ ! -f "$PROJECT_ROOT/devwisdom" ]]; then
    echo "Building devwisdom binary..." >&2
    make build >&2 || {
        echo "Error: Failed to build devwisdom binary" >&2
        exit 1
    }
fi

# Run Go binary
exec "$PROJECT_ROOT/devwisdom"

