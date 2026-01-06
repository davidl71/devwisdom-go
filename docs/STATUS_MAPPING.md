# Task Status Mapping Between Todo2 and Agentic-Tools

## Overview

Todo2 and Agentic-Tools MCP use **different status naming conventions**. This document explains the mismatch and provides the mapping between systems.

## Status Value Comparison

### Todo2 Status Values (Title Case)
- **`Todo`** - Task not started (capitalized, space-separated)
- **`In Progress`** - Task actively being worked on (capitalized, space-separated)
- **`Done`** - Task completed (capitalized)
- **`Review`** - Task awaiting review/approval (capitalized)

### Agentic-Tools Status Values (lowercase)
- **`pending`** - Task not started (lowercase)
- **`in-progress`** - Task actively being worked on (lowercase, hyphenated)
- **`done`** - Task completed (lowercase)

## Status Mapping

The mapping between systems (from `project-management-automation/scripts/automate_todo_sync.py`):

### Todo2 → Agentic-Tools
```
Todo        → pending
In Progress → in-progress
Done        → done
Review      → in-progress  (Review maps to in-progress in agentic-tools)
```

### Agentic-Tools → Todo2
```
pending     → Todo
in-progress → In Progress
done        → Done
```

## Why They Don't Match

**Different Systems, Different Conventions:**

1. **Todo2** (Cursor extension) uses **Title Case** with spaces:
   - Designed for human-readable display
   - Matches common task management UI conventions
   - Status values: `Todo`, `In Progress`, `Done`

2. **Agentic-Tools MCP** uses **lowercase** with hyphens:
   - Designed for programmatic access
   - Follows API/JSON conventions (lowercase, hyphenated)
   - Status values: `pending`, `in-progress`, `done`

## Synchronization

The `sync_todo_tasks` tool in project-management-automation handles the mapping automatically:

```python
# From automate_todo_sync.py
status_map = {
    'pending': 'Todo',
    'in_progress': 'In Progress',
    'completed': 'Done',
    'Todo': 'pending',
    'In Progress': 'in_progress',
    'Done': 'completed',
    'Review': 'in_progress',  # Review maps to in-progress
    'Cancelled': 'completed'
}
```

## Normalization Utilities

The `project-management-automation` project provides normalization utilities in `utils/todo2_utils.py`:

- `normalize_status(status)` - Normalizes to canonical lowercase form
- `normalize_status_to_title_case(status)` - Normalizes to Title Case for Todo2
- `is_pending_status(status)` - Checks if status is pending
- `is_completed_status(status)` - Checks if status is completed

## Current Status in devwisdom-go

**Todo2 Tasks:**
- `Todo`: 22 tasks
- `In Progress`: 4 tasks

**Agentic-Tools Tasks:**
- `pending`: 20 tasks
- `in-progress`: 2 tasks
- `done`: 8 tasks

## Recommendations

1. **Use sync tool** - The `sync_todo_tasks` tool automatically handles mapping
2. **Use normalization utilities** - When checking statuses, use `normalize_status()` for consistent comparison
3. **Don't manually sync** - Let the sync tool handle the conversion to avoid errors

## Related Documentation

- **TASK_STATUS_STANDARDIZATION.md** - Full status standardization guide in project-management-automation
- **todo2_utils.py** - Status normalization utilities
- **automate_todo_sync.py** - Sync implementation with status mapping

