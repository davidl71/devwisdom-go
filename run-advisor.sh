#!/bin/bash
# Wrapper script for devwisdom-go MCP server
# Ensures sources.json can be found by running from the correct directory

cd "$(dirname "$0")" || exit 1
exec ./devwisdom "$@"

