# Todo2 Task Import Instructions

This document provides instructions for importing the exarp-go improvement tasks into Todo2.

## Created Files

1. **`TODO_EXARP_GO_IMPROVEMENTS.md`** - Markdown format with all 12 tasks (formatted for task discovery)
2. **`exarp_go_tasks_todo2.json`** - JSON format with structured task data
3. **`docs/EXARP_GO_IMPROVEMENT_TASKS.md`** - Detailed task documentation

## Import Methods

### Method 1: Using Task Discovery Tool (Recommended)

If you have the exarp-go MCP server configured with Python dependencies:

```bash
# From project root
cd /home/dlowes/projects/devwisdom-go

# Use task discovery tool to import tasks
# This requires the task_discovery tool from exarp-go
# The tool should parse TODO_EXARP_GO_IMPROVEMENTS.md
```

**Note:** The task_discovery tool requires Python dependencies (`project_management_automation` module). If these aren't available, use one of the manual methods below.

### Method 2: Manual Creation via MCP Tools

Use the exarp-go MCP tools directly through Cursor:

1. Open Cursor chat
2. Use the `task_discovery` tool:
   ```
   Discover tasks from TODO_EXARP_GO_IMPROVEMENTS.md and create them in Todo2
   ```
3. The tool should parse the markdown and create tasks automatically

### Method 3: Manual Import (Fallback)

If automated methods don't work, you can manually create tasks in Todo2:

1. Review `TODO_EXARP_GO_IMPROVEMENTS.md`
2. For each task (T-A1 through T-E2):
   - Create a new task in Todo2
   - Set name: e.g., "T-A1: Create Framework-Agnostic MCP Interface"
   - Set status: "Todo"
   - Set priority: High/Medium/Low
   - Add tags: from the tags field
   - Add description: from the description field
   - Add dependencies: from the dependencies field
   - Add subtasks: from the tasks array

## Task Summary

### Phase 1: Framework Abstraction & Architecture (3 tasks)
- T-A1: Create Framework-Agnostic MCP Interface (High)
- T-A2: Implement Factory Pattern (High)
- T-A3: Configuration Management (Medium)

### Phase 2: CLI/MCP Dual Mode (3 tasks)
- T-B1: TTY Detection for Dual Mode (High)
- T-B2: CLI Command Structure (High)
- T-B3: Reuse Server Infrastructure (Medium)

### Phase 3: Tool Registration & Organization (2 tasks)
- T-C1: Organize Tool Registration in Batches (Medium)
- T-C2: Improve Tool Schema Definitions (Medium)

### Phase 4: Development Workflow & Testing (3 tasks)
- T-D1: Hot Reload Development Scripts (Medium)
- T-D2: Enhance Makefile with Tool Detection (Low)
- T-D3: Tool Testing Utilities (Medium)

### Phase 5: Security & Path Validation (2 tasks)
- T-E1: Path Validation Security (High)
- T-E2: Workspace Root Detection (Medium)

## Dependencies Map

```
T-A1 (Interface) ──┐
                   ├──> T-A2 (Factory)
                   └──> T-C1 (Batch Registration)

T-A1 ──> T-B1 (TTY Detection)
T-A2 ──> T-B1
T-B1 ──> T-B2 (CLI Commands)
T-B2 ──> T-B3 (Reuse Infrastructure)
T-B2 ──> T-D3 (Testing Utilities)
T-C1 ──> T-C2 (Schema Definitions)
T-D1 ──> T-D2 (Makefile)
T-E1 ──> T-E2 (Workspace Root)
```

## Next Steps

1. **Import Tasks**: Use one of the methods above to import tasks into Todo2
2. **Add Research Comments**: Before starting each task, add research comments following exarp-go's pattern:
   - Section 1: Local Codebase Analysis
   - Section 2: Internet Research
   - Section 3: Synthesis from Batch Research
3. **Start with Phase 1**: Begin with T-A1 (Framework-Agnostic Interface) as it's the foundation
4. **Track Progress**: Update task status in Todo2 as you complete each task

## Related Documentation

- [EXARP_GO_LESSONS.md](./EXARP_GO_LESSONS.md) - Complete lessons learned document
- [EXARP_GO_IMPROVEMENT_TASKS.md](./EXARP_GO_IMPROVEMENT_TASKS.md) - Detailed task documentation
- [TODO_EXARP_GO_IMPROVEMENTS.md](../TODO_EXARP_GO_IMPROVEMENTS.md) - Source markdown file

