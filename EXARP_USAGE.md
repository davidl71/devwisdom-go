# Using Exarp Tools for devwisdom-go Planning

This guide shows how to use exarp MCP tools to plan and track devwisdom-go development.

---

## Quick Start

### Via Cursor (Recommended)

Ask the AI assistant in Cursor:

1. **Generate PRD**:
   ```
   Use exarp prd_generator tool to generate a PRD for devwisdom-go 
   using the content from devwisdom-go/PROJECT_GOALS.md
   ```

2. **Discover Tasks**:
   ```
   Use exarp task_discovery tool to extract tasks from 
   devwisdom-go/TODO.md and create structured tasks
   ```

3. **Analyze Tasks**:
   ```
   Use exarp task_analysis tool with action=hierarchy to analyze 
   the task structure for devwisdom-go
   ```

4. **Check Alignment**:
   ```
   Use exarp analyze_todo2_alignment tool to check if tasks align 
   with devwisdom-go/PROJECT_GOALS.md
   ```

---

## Available Exarp Tools for Planning

### 1. PRD Generation
**Tool**: `prd_generator`  
**Purpose**: Generate structured PRD from PROJECT_GOALS.md  
**Input**: PROJECT_GOALS.md content  
**Output**: Structured PRD with personas, features, workflows

**Usage**:
```python
# Via MCP tool call
prd_generator(
    project_root="/path/to/project-management-automation",
    goals_content=PROJECT_GOALS.md_content
)
```

---

### 2. Task Discovery
**Tool**: `task_discovery`  
**Purpose**: Extract tasks from markdown files  
**Input**: TODO.md or any markdown  
**Output**: Structured task list

**Usage**:
```python
# Via MCP tool call
task_discovery(
    action="markdown",
    doc_path="devwisdom-go/TODO.md",
    create_tasks=True  # Create Todo2 tasks
)
```

---

### 3. Task Analysis
**Tool**: `task_analysis`  
**Purpose**: Analyze task structure, hierarchy, duplicates, tags  
**Input**: Existing tasks  
**Output**: Analysis and recommendations

**Usage**:
```python
# Via MCP tool call
task_analysis(
    action="hierarchy",  # or "tags", "duplicates"
    output_path="devwisdom-go/task_analysis.md"
)
```

---

### 4. Alignment Analysis
**Tool**: `analyze_todo2_alignment`  
**Purpose**: Check task alignment with PROJECT_GOALS.md  
**Input**: Tasks + PROJECT_GOALS.md  
**Output**: Alignment score, misaligned tasks

**Usage**:
```python
# Via MCP tool call
analyze_todo2_alignment(
    create_followup_tasks=True,
    output_path="devwisdom-go/alignment_report.md"
)
```

---

## Recommended Workflow

### Step 1: Initial Planning
1. ✅ Create PROJECT_GOALS.md (done)
2. ✅ Create TODO.md (done)
3. 🔄 Generate PRD using `prd_generator`
4. 🔄 Discover tasks using `task_discovery`
5. 🔄 Create Todo2 tasks using agentic-tools MCP

### Step 2: Analysis
1. 🔄 Analyze task hierarchy using `task_analysis`
2. 🔄 Check alignment using `analyze_todo2_alignment`
3. 🔄 Detect duplicates using `task_analysis` (action=duplicates)
4. 🔄 Consolidate tags using `task_analysis` (action=tags)

### Step 3: Development Tracking
1. 🔄 Track progress with `project_scorecard`
2. 🔄 Run daily checks with `run_daily_automation`
3. 🔄 Monitor alignment periodically
4. 🔄 Update tasks as work progresses

---

## Task Creation Strategy

### Phase-Based Tasks
Create tasks aligned with PROJECT_GOALS.md phases:

```
Phase 2: Wisdom Data Porting
├── Port pistis_sophia source
├── Port stoic source
├── Port tao source
├── ... (21+ sources)
└── Implement random source selector

Phase 3: Advisor System
├── Complete metric advisor mappings
├── Complete tool advisor mappings
├── Complete stage advisor mappings
└── Implement mode-aware selection

Phase 4: MCP Protocol
├── Implement JSON-RPC 2.0 handler
├── Register consult_advisor tool
├── Register get_wisdom tool
├── ... (5 tools total)
└── Register resources (4 resources)
```

### Task Tags
Use tags for organization:
- `phase-2` - Wisdom data porting
- `phase-3` - Advisor system
- `phase-4` - MCP protocol
- `wisdom-source` - Wisdom source porting
- `mcp-tool` - MCP tool implementation
- `testing` - Test-related
- `documentation` - Docs-related

---

## Progress Tracking

### Daily Checks
Run `run_daily_automation` to:
- Check alignment with goals
- Detect duplicate tasks
- Update project health

### Weekly Reviews
Use `project_scorecard` to:
- Assess overall progress
- Identify blockers
- Review task completion

### Milestone Reviews
Use `analyze_todo2_alignment` to:
- Verify alignment with PROJECT_GOALS.md
- Identify misaligned work
- Adjust priorities

---

## Integration with Agentic-Tools MCP

Use agentic-tools MCP to create and manage tasks:

```python
# Create project first (if not exists)
create_project(
    name="devwisdom-go",
    description="Wisdom module extraction - Go proof of concept"
)

# Create phase tasks
create_task(
    projectId="devwisdom-go-id",
    name="Phase 2: Port pistis_sophia source",
    details="Port pistis_sophia wisdom source from Python to Go",
    status="pending",
    priority=8,
    tags=["phase-2", "wisdom-source"]
)

# Link dependencies
create_task(
    projectId="devwisdom-go-id",
    name="Phase 4: Implement JSON-RPC handler",
    details="Implement JSON-RPC 2.0 handler for MCP protocol",
    status="pending",
    priority=9,
    dependsOn=["phase-2-complete", "phase-3-complete"],
    tags=["phase-4", "mcp-protocol"]
)
```

---

**Last Updated**: 2025-01-26
