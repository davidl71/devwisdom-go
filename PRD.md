# PRD: devwisdom-go

*Generated: 2025-12-07T23:59:11.999386*

---

## 1. Overview

**Project:** devwisdom-go
**Type:** MCP Server / Developer Tools

### Vision

A comprehensive MCP server providing project management automation tools for AI-assisted development workflows. Enables documentation health checks, task alignment, duplicate detection, security scanning, and intelligent automation orchestration.

## 2. Problem Statement

A comprehensive MCP server providing project management automation tools for AI-assisted development workflows. Enables documentation health checks, task alignment, duplicate detection, security scanning, and intelligent automation orchestration.

## 3. Target Users / Personas

| Persona | Role | Trusted Advisor | Goal |
|---------|------|-----------------|------|
| **Developer** | Daily Contributor | 💻 Tao Of Programming | Write quality code, stay unblocked, contribute effectively |
| **Project Manager** | Delivery Focus | ⚔️ Art Of War | Track progress, remove blockers, ensure delivery |
| **Code Reviewer** | Quality Assurance | 🏛️ Stoic | Ensure code quality and standards compliance |
| **Architect** | System Design | 🔮 Enochian | Design scalable, maintainable systems |
| **Security Engineer** | Risk Management | 😈 Bofh | Identify and mitigate security risks |
| **QA Engineer** | Quality Assurance | 🏛️ Stoic | Ensure product quality through testing |
| **Executive/Stakeholder** | Status Overview | 📜 Pistis Sophia | Understand project health and progress at a glance |
| **Technical Writer** | Documentation | 🎓 Confucius | Create and maintain clear documentation |

### Persona Details

#### 👤 Developer (Daily Contributor)

**Goal:** Write quality code, stay unblocked, contribute effectively

**Trusted Advisor:** 💻 **Tao Of Programming**
> *The Tao of Programming teaches elegant flow - let code emerge naturally*

**Key Metrics:**
- Cyclomatic Complexity <10
- Test Coverage >80%
- Bandit Findings: 0 high/critical

**Workflows:**
- Morning Checkin
- Before Committing
- Before PR/Push
- Weekly Self-Review

*Relevance to this project: 6 keyword matches, Primary user of development tools*

#### 👤 Project Manager (Delivery Focus)

**Goal:** Track progress, remove blockers, ensure delivery

**Trusted Advisor:** ⚔️ **Art Of War**
> *Sun Tzu teaches strategy and decisive execution - sprints are campaigns*

**Key Metrics:**
- On-Time Delivery
- Blocked Tasks: 0
- Sprint Velocity

**Workflows:**
- Daily Standup Prep
- Sprint Planning
- Sprint Review
- Stakeholder Update

*Relevance to this project: 2 keyword matches, Todo2 task management present*

#### 👤 Code Reviewer (Quality Assurance)

**Goal:** Ensure code quality and standards compliance

**Trusted Advisor:** 🏛️ **Stoic**
> *Stoics accept harsh truths with equanimity - reviews reveal reality*

**Key Metrics:**
- Review Cycle Time <24h
- Defect Escape Rate <5%

**Workflows:**
- PR Review
- Architecture Review
- Security Review

*Relevance to this project: 5 keyword matches*

#### 👤 Architect (System Design)

**Goal:** Design scalable, maintainable systems

**Trusted Advisor:** 🔮 **Enochian**
> *Enochian mysticism reveals hidden structure and patterns in architecture*

**Key Metrics:**
- Avg Complexity <5
- Max Complexity <15
- Distance from Main Sequence <0.3

**Workflows:**
- Weekly Architecture Review
- Before Major Changes
- Tech Debt Prioritization

*Relevance to this project: 2 keyword matches*

#### 👤 Security Engineer (Risk Management)

**Goal:** Identify and mitigate security risks

**Trusted Advisor:** 😈 **Bofh**
> *BOFH is paranoid about security - expects users to break everything*

**Key Metrics:**
- Critical Vulns: 0
- High Vulns: 0
- Security Score >90%

**Workflows:**
- Daily Scan
- Weekly Deep Scan
- Security Audit

*Relevance to this project: 2 keyword matches, Project has security focus*

#### 👤 QA Engineer (Quality Assurance)

**Goal:** Ensure product quality through testing

**Trusted Advisor:** 🏛️ **Stoic**
> *Stoics teach discipline through adversity - tests reveal truth*

**Key Metrics:**
- Test Coverage >80%
- Tests Passing: 100%
- Defect Density <5/KLOC

**Workflows:**
- Daily Testing Status
- Sprint Testing Review
- Defect Analysis

*Relevance to this project: 4 keyword matches, Test suite present*

#### 👤 Executive/Stakeholder (Status Overview)

**Goal:** Understand project health and progress at a glance

**Trusted Advisor:** 📜 **Pistis Sophia**
> *Pistis Sophia's journey through aeons mirrors understanding project health stages*

**Key Metrics:**
- Overall Health Score
- On-Time Delivery %
- Risk Level

**Workflows:**
- Weekly Status Review
- Stakeholder Briefing
- Executive Dashboard

*Relevance to this project: 4 keyword matches*

#### 👤 Technical Writer (Documentation)

**Goal:** Create and maintain clear documentation

**Trusted Advisor:** 🎓 **Confucius**
> *Confucius emphasized teaching and transmitting wisdom to future generations*

**Key Metrics:**
- Broken Links: 0
- Stale Docs: 0
- Docstring Coverage >90%

**Workflows:**
- Weekly Doc Health
- Doc Update Cycle
- API Documentation

*Relevance to this project: 5 keyword matches, Documentation is emphasized*

## 4. User Stories

### US-1: Fix broken documentation links

**Priority:** 🟠 HIGH
**Status:** Done
**Tags:** links, fix, docs

*As a* Technical Writer (Documentation),
*I want* Fix broken documentation links,
*So that* I can understand and use the project effectively.

### US-2: Research: Exarp Cursor Extension Architecture

**Priority:** 🟢 LOW (Future Goal)
**Status:** Future Enhancement
**Tags:** research, architecture, extension, typescript, future

> **⚠️ FUTURE GOAL - NOT CURRENTLY BEING DEVELOPED**  
> This user story represents a potential future enhancement. The cursor extension is documented but not in active development. The MCP server works standalone and provides all core functionality.

*As a* Developer (Daily Contributor),
*I want* Research: Exarp Cursor Extension Architecture,
*So that* I can be confident the code works correctly.

**Details:**
🎯 **Objective:** Evaluate and plan the implementation of a Cursor extension to complement the Exarp MCP server. **This is a future enhancement, not currently being developed.**

📋 **Acceptance Criteria:**
- Review VS Code Extension API documentation
- Prototype minimal status bar integration
- Test Todo2 file watching from TypeScript
- Validate extension ↔ MCP communication patterns
- Document technical decisions

🚫 **Scope Boundaries:**
- **Included:** Phase 1 (status bar + basic commands) research
- **Excluded:** Full implementation, MCP Proxy
- **Clarifica

### US-3: Enhance Exarp to validate and use Todo2 project ownership

**Priority:** 🟠 HIGH
**Status:** Todo
**Tags:** enhancement, project-ownership, validation, one-system

*As a* Developer (Daily Contributor),
*I want* Enhance Exarp to validate and use Todo2 project ownership,
*So that* I can save time through automation.

**Details:**
Implement project ownership validation in Exarp to ensure it operates on the correct project and enable future "one system to track them all" capability.

## Changes Made
- Added `project` object to state.todo2.json root with: id, name, path, repository
- Added `project_id` field to every task (538 tasks across 2 projects)

## Required Exarp Enhancements

### 1. Project Validation (High Priority)
- On startup, compare PROJECT_ROOT env var with Todo2's project.path
- Warn or fail if mismatch dete

### US-4: Problems Advisor MCP Tool

**Priority:** 🟠 HIGH
**Status:** Done
**Tags:** mcp, tools, diagnostics

*As a* Developer (Daily Contributor),
*I want* Problems Advisor MCP Tool,
*So that* the project is secure and protected.

**Details:**
New MCP tool that analyzes IDE diagnostics and provides resolution hints.

## Features
- Categorizes problems (spelling, type_error, import_error, unused_code, etc.)
- Provides resolution hints for each problem
- Identifies auto-fixable issues
- Suggests quick fix commands
- Generates detailed reports

## Usage
```python
# Analyze problems from read_lints output
analyze_problems(problems_json="[{...}]", include_hints=True)

# List all categories
list_problem_categories()
```

## Categories Suppo

### US-5: Implement generate_prd tool for PRD generation from codebase

**Priority:** 🟠 HIGH
**Status:** Done
**Tags:** best-practices, new-tool, prd

*As a* Developer (Daily Contributor),
*I want* Implement generate_prd tool for PRD generation from codebase,
*So that* I can understand and use the project effectively.

**Details:**
Create a new MCP tool that generates Product Requirements Documents by analyzing the existing codebase.

**Objective:** Help users create/maintain PRD files as their "North Star" document.

## Research Summary

### Key PRD AI Tools Analyzed
1. **ChatPRD** - AI copilot specifically for PMs, $5/month, MCP integration available
2. **MakePRD** - Generate detailed PRDs from product descriptions
3. **Miro AI PRD Generator** - Visual PRD creation with collaboration
4. **Leiga** - AI-assisted drafts wit

### US-6: Implement recommend_workflow_mode tool for AGENT vs ASK guidance

**Priority:** 🟠 HIGH
**Status:** Done
**Tags:** best-practices, new-tool, workflow

*As a* Developer (Daily Contributor),
*I want* Implement recommend_workflow_mode tool for AGENT vs ASK guidance,
*So that* I can save time through automation.

**Details:**
Create a tool that recommends Cursor workflow mode (AGENT vs ASK) based on task description.

**Objective:** Help users choose the right Cursor mode for each task.

**Features:**
- Analyze task description for implementation vs analysis intent
- Map Exarp tools to recommended modes
- Suggest relevant Exarp tools for the task
- Integrate with tractatus → exarp → sequential workflow

**Mode Mapping:**
- ASK: project_scorecard, check_documentation_health, analyze_todo2_alignment
- AGENT: sprint_aut

### US-7: Implement prompt iteration tracking tools

**Priority:** 🟠 HIGH
**Status:** Done
**Tags:** best-practices, new-tool, prompts

*As a* Developer (Daily Contributor),
*I want* Implement prompt iteration tracking tools,
*So that* the project capabilities are enhanced.

**Details:**
Create tools for tracking and improving prompt iterations.

**Objective:** Embody the "iterate, iterate, iterate" philosophy from best practices.

**New Tools:**
1. record_prompt_iteration - Track original prompt, issue, refined prompt
2. suggest_prompt_refinement - Coach better prompts based on issue type

**Features:**
- Store iteration history in session_memory
- Build project-specific prompt patterns
- Integrate with advisor system for coaching
- Suggest refinements for common issues

**Issu

### US-8: Create unit tests for Git-inspired features

**Priority:** 🟠 HIGH
**Status:** Todo
**Tags:** testing, git-inspired, unit-tests, branch:git-features

*As a* QA Engineer (Quality Assurance),
*I want* Create unit tests for Git-inspired features,
*So that* I can be confident the code works correctly.

### US-9: Integrate Git-inspired tools into MCP server

**Priority:** 🟠 HIGH
**Status:** Todo
**Tags:** integration, git-inspired, mcp-server, branch:git-features

*As a* Developer (Daily Contributor),
*I want* Integrate Git-inspired tools into MCP server,
*So that* the project capabilities are enhanced.

### US-10: Fix broken documentation links

**Priority:** 🟠 HIGH
**Status:** Done
**Tags:** documentation, fix, links, completed

*As a* Technical Writer (Documentation),
*I want* Fix broken documentation links,
*So that* I can understand and use the project effectively.

### US-11: Fix FastMCP return type issues (dict vs string)

**Priority:** 🟠 HIGH
**Status:** Done
**Tags:** bugfix, fastmcp, tools, resources

*As a* Developer (Daily Contributor),
*I want* Fix FastMCP return type issues (dict vs string),
*So that* the project capabilities are enhanced.

### US-12: Review misaligned high-priority tasks

**Priority:** 🟠 HIGH
**Status:** todo
**Tags:** todo2, alignment, review

*As a* Code Reviewer (Quality Assurance),
*I want* Review misaligned high-priority tasks,
*So that* the project capabilities are enhanced.

### US-13: Automation: Sprint Automation

**Priority:** 🟡 MEDIUM
**Status:** Done
**Tags:** automation, sprint-automation

*As a* Project Manager (Delivery Focus),
*I want* Automation: Sprint Automation,
*So that* I can save time through automation.

### US-14: Automation: Documentation Health Analysis

**Priority:** 🟡 MEDIUM
**Status:** Done
**Tags:** automation, documentation-health-analysis

*As a* Technical Writer (Documentation),
*I want* Automation: Documentation Health Analysis,
*So that* I can save time through automation.

### US-15: Automation: Todo2 Alignment Analysis

**Priority:** 🟡 MEDIUM
**Status:** Done
**Tags:** automation, todo2-alignment-analysis

*As a* Developer (Daily Contributor),
*I want* Automation: Todo2 Alignment Analysis,
*So that* I can save time through automation.

### US-16: Automation: Todo2 Duplicate Detection

**Priority:** 🟡 MEDIUM
**Status:** Done
**Tags:** automation, todo2-duplicate-detection

*As a* Developer (Daily Contributor),
*I want* Automation: Todo2 Duplicate Detection,
*So that* I can save time through automation.

### US-17: Automation: Automation Opportunity Finder

**Priority:** 🟡 MEDIUM
**Status:** Done
**Tags:** automation, coverage-analyzer, test-runner, automation-finder

*As a* Developer (Daily Contributor),
*I want* Automation: Automation Opportunity Finder,
*So that* I can save time through automation.

### US-18: Automation: Test Runner

**Priority:** 🟡 MEDIUM
**Status:** Done
**Tags:** automation, test-runner

*As a* QA Engineer (Quality Assurance),
*I want* Automation: Test Runner,
*So that* I can save time through automation.

### US-19: Automation: Test Coverage Analyzer

**Priority:** 🟡 MEDIUM
**Status:** Done
**Tags:** automation, coverage-analyzer

*As a* QA Engineer (Quality Assurance),
*I want* Automation: Test Coverage Analyzer,
*So that* I can save time through automation.

### US-20: Update stale documentation

**Priority:** 🟡 MEDIUM
**Status:** Todo
**Tags:** update, docs

*As a* Technical Writer (Documentation),
*I want* Update stale documentation,
*So that* I can understand and use the project effectively.

*... and 109 more user stories*

## 5. Key Features

- Tool
- `server_status`
- `project_scorecard`
- `project_overview`
- Tool
- `check_documentation_health`
- `add_external_tool_hints`
- Tool
- `analyze_todo2_alignment`
- `detect_duplicate_tasks`
- `consolidate_tags`
- `task_hierarchy_analyzer`
- `batch_approve_tasks`
- `sync_todo_tasks`
- Tool
- `scan_dependency_security`
- Tool
- `run_daily_automation`
- `run_nightly_task_automation`
- `sprint_automation`
- `find_automation_opportunities`
- Tool
- `validate_ci_cd_workflow`
- `setup_git_hooks`
- `setup_pattern_triggers`
- `check_working_copy_health`
- Tool
- `run_tests`
- `analyze_test_coverage`
- Tool: Session Memory

## 6. Technical Requirements

**Language:** Python
**Frameworks:** FastMCP, Pydantic, pytest, FastMCP, Pydantic, pytest
**Patterns:** MCP Server Pattern, Todo2 Task Management, Unit Testing, CI/CD with GitHub Actions

### Dependencies

- fastmcp
- pydantic
- pyyaml
- pytest
- pytest-cov
- pytest-asyncio
- coverage

### Project Structure

- `extension/src/`
- `cursor-extension/src/`
- `.venv/lib/`
- `build/lib/`
- `tools/`
- `project_management_automation/tools/`
- `scripts/`
- `project_management_automation/scripts/`
- `tests/`
- `docs/`

## 7. Success Metrics

| Metric | Target | Current |
|--------|--------|---------|
| Tool Count | ≤30 | ~30 ✅ |
| Test Coverage | 80% | TBD |
| Doc Health Score | 90+ | TBD |
| MCP Compliance | Full | Partial |

## 8. Risks & Dependencies

- 🟡 Dependency updates may introduce breaking changes
- 🟡 External API changes may affect integrations

## 9. Timeline & Progress

### Phases

- **Phase 1:** Core Infrastructure
- **Phase 2:** Integration & Interoperability
- **Phase 3:** Automation & Intelligence
- **Phase 4:** Quality & Testing
- **Phase 5:** Documentation & Polish

### Current Progress

- **Total Tasks:** 129
- **Completed:** 93 (72%)
- **In Progress:** 1
- **Remaining:** 35
- **Estimated Hours Remaining:** 70h

---

*Generated by Exarp PRD Generator*

## How to Use This PRD

1. **Review and refine** - This is a starting point, customize for your needs
2. **Iterate with AI** - Use prompts like 'Review this PRD and suggest improvements'
3. **Keep updated** - Re-run generation as project evolves
4. **Align tasks** - Use `analyze_todo2_alignment` to verify task alignment with PRD
