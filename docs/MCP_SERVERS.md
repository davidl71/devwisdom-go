# MCP Servers Configuration

This document describes the MCP (Model Context Protocol) servers configured for devwisdom-go.

## Overview

The `.cursor/mcp.json` file configures multiple MCP servers that provide complementary functionality to the devwisdom-go project. These servers are integrated into Cursor's AI assistant to provide enhanced capabilities.

## Configured Servers

### 1. **devwisdom** (Primary Server)
- **Command**: `/Users/davidl/Projects/devwisdom-go/devwisdom`
- **Description**: DevWisdom Go MCP Server - Wisdom quotes, trusted advisors, and inspirational guidance
- **Purpose**: Provides the core wisdom functionality for the project

### 2. **exarp_pma** (Project Management Automation)
- **Command**: `/Users/davidl/Projects/project-management-automation/exarp-uvx-wrapper.sh --mcp`
- **Description**: Exarp for Project Management Automation - auto-detects uvx location across platforms (FastMCP mode)
- **Purpose**: Provides project management automation tools including:
  - Documentation health checks
  - Todo2 alignment analysis
  - Duplicate task detection
  - Dependency security scanning
  - Automation opportunity discovery
  - Todo synchronization

### 3. **agentic-tools** (Task Management)
- **Command**: `uvx mcpower-proxy==0.0.87 --wrapped-config {...} @pimzino/agentic-tools-mcp`
- **Description**: Task management and agent memories with JSON file storage for Todo2 integration
- **Purpose**: Provides task management capabilities and agent memory storage

### 4. **tractatus_thinking** (Structural Analysis)
- **Command**: `uvx mcpower-proxy==0.0.87 --wrapped-config {...} tractatus_thinking`
- **Description**: Tractatus Thinking MCP server for structural analysis and logical decomposition
- **Purpose**: Assists in structural analysis and logical decomposition of complex problems
- **Workflow**: Use BEFORE exarp tools for structural analysis (WHAT)

### 5. **sequential_thinking** (Implementation Workflows)
- **Command**: `uvx mcpower-proxy==0.0.87 --wrapped-config {...} @modelcontextprotocol/server-sequential-thinking`
- **Description**: Sequential Thinking MCP server for implementation workflows
- **Purpose**: Converts structural understanding into implementation workflows
- **Workflow**: Use AFTER exarp analysis for implementation steps (HOW)

### 6. **context7** (Documentation Lookup)
- **Command**: `uvx mcpower-proxy==0.0.87 --wrapped-config {...} @upstash/context7-mcp`
- **Description**: Up-to-date documentation lookup for libraries and frameworks
- **Purpose**: Provides access to up-to-date documentation for various libraries and frameworks

### 7. **filesystem** (File Operations)
- **Command**: `uvx mcpower-proxy==0.0.87 --wrapped-config {...} @modelcontextprotocol/server-filesystem`
- **Description**: File system operations for reading, writing, and managing project files
- **Purpose**: Provides file system access within the project directory
- **Path**: Configured for `/Users/davidl/Projects/devwisdom-go`

### 8. **mcp-generic-tools** (Generic Tools)
- **Command**: `/Users/davidl/Projects/mcp-generic-tools/run_server.sh`
- **Description**: Generic MCP tools (context management, prompt tracking, recommendations)
- **Purpose**: Provides context management, prompt tracking, and recommendations

### 9. **mcp-stdio-tools** (Stdio Tools)
- **Command**: `/Users/davidl/Projects/mcp-stdio-tools/run_server.sh`
- **Description**: Stdio-based MCP tools (12 tools migrated from exarp_pma)
- **Purpose**: Provides stdio-based tools that work with the MCP protocol

### 10. **gotohuman** (Human-in-the-Loop)
- **Command**: `uvx mcpower-proxy==0.0.87 --wrapped-config {...} @gotohuman/mcp-server`
- **Description**: Human-in-the-loop platform - Allow AI agents and automations to send requests for approval to your gotoHuman inbox
- **Purpose**: Enables batch operations and critical actions to request human confirmation before execution
- **Use Cases**: 
  - Batch task approvals
  - Critical operation confirmations
  - Human oversight for automated workflows
- **Setup**: Requires a gotoHuman account at https://www.gotohuman.com

## Recommended Workflow

For optimal results when working with these MCP servers:

1. **tractatus_thinking** → Understand problem structure (WHAT)
2. **exarp_pma** → Analyze and automate project management tasks
3. **sequential_thinking** → Convert results into implementation steps (HOW)

## Configuration File Location

The MCP server configuration is located at:
```
.cursor/mcp.json
```

## Environment Variables

Several servers use environment variables:

- **PROJECT_ROOT**: Set to `/Users/davidl/Projects/devwisdom-go` for servers that need project context
- **EXARP_DEV_MODE**: Set to `1` for exarp_pma to enable development mode

## Dependencies

### Required Tools

- **uvx**: Python package runner (from `uv` toolchain)
- **npx**: Node.js package runner (for npm-based MCP servers)
- **mcpower-proxy**: MCP proxy wrapper (version 0.0.87)

### External Projects

- **project-management-automation**: Located at `/Users/davidl/Projects/project-management-automation`
- **mcp-generic-tools**: Located at `/Users/davidl/Projects/mcp-generic-tools`
- **mcp-stdio-tools**: Located at `/Users/davidl/Projects/mcp-stdio-tools`

## Troubleshooting

### Server Not Starting

1. Check that `uvx` is installed and in PATH
2. Verify paths to external projects are correct
3. Check that `devwisdom` binary exists and is executable
4. Review Cursor's MCP server logs for error messages

### Path Issues

If you're using this configuration on a different machine:
1. Update all absolute paths in `.cursor/mcp.json`
2. Update `PROJECT_ROOT` environment variables
3. Verify external project paths exist

### Permission Issues

Ensure that:
- `devwisdom` binary is executable: `chmod +x devwisdom`
- Script files in external projects are executable
- Cursor has permission to execute the configured commands

## Related Documentation

- **MCP Specification**: https://modelcontextprotocol.io/
- **PROJECT_GOALS.md**: Strategic phases and goals
- **README.md**: Project overview and quick start

## Updates

This configuration was created by integrating MCP servers from the project-management-automation project. The configuration has been adapted for devwisdom-go with appropriate path and environment variable updates.

**Latest Update**: Added `gotohuman` server for human-in-the-loop confirmations and batch operation approvals.

